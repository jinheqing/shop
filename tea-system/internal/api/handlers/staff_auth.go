package handlers

import (
	"context"
	"os"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"tea-system/internal/middleware"
	"tea-system/internal/models"
	"tea-system/internal/repository"
	"tea-system/internal/service"
)

// ============================================================
// Staff 认证
// ============================================================

type StaffAuthHandler struct {
	staffRepo    *repository.StaffRepo
	userRepo     *repository.UserRepo
	passwordSvc  *service.PasswordService
	jwtSvc       *service.JWTService
	mfaSvc       *service.MFAService
	rdb          *redis.Client
	auditDB      *gorm.DB
}

func NewStaffAuthHandler(
	staffRepo *repository.StaffRepo,
	userRepo *repository.UserRepo,
	passwordSvc *service.PasswordService,
	jwtSvc *service.JWTService,
	mfaSvc *service.MFAService,
	rdb *redis.Client,
	auditDB *gorm.DB,
) *StaffAuthHandler {
	return &StaffAuthHandler{
		staffRepo:   staffRepo,
		userRepo:    userRepo,
		passwordSvc: passwordSvc,
		jwtSvc:      jwtSvc,
		mfaSvc:      mfaSvc,
		rdb:         rdb,
		auditDB:     auditDB,
	}
}

// StaffLoginRequest — POST /staff/login
type StaffLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// StaffLoginResponse — 登录响应
type StaffLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // 秒
	MFARequired  bool   `json:"mfa_required"`
	StaffID      uint64 `json:"staff_id"`
	Role         string `json:"role"`
	Email        string `json:"email"`
}

// Login — POST /staff/login
func (h *StaffAuthHandler) Login(c *gin.Context) {
	var req StaffLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 1. 查找 Staff
	staff, err := h.staffRepo.GetByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		if errors.Is(err, repository.ErrStaffNotFound) {
			// 安全：统一返回 "email or password invalid"，不要泄露哪个不对
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "email or password invalid"})
			return
		}
		log.Error().Err(err).Msg("staff login: db error")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	// 2. 检查是否被禁用
	if !staff.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "account disabled"})
		return
	}

	// 3. 验证密码
	if !h.passwordSvc.Verify(staff.PasswordHash, req.Password) {
		// TODO: 登录失败计数 + 5 次锁定 15 分钟
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "email or password invalid"})
		return
	}

	// 4. 更新 last_login_at
	_ = h.staffRepo.UpdateLastLogin(ctx, staff.ID)

	// 5. 判断是否需要 MFA（admin/supervisor 强制 + 任何角色开启了 mfa_enabled）
	//    dev 模式（GIN_MODE=debug）下跳过 ForceMFA，方便测试 admin 账号
	mfaRequired := staff.MfaEnabled
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = os.Getenv("SERVER_MODE")
	}
	if ginMode != "debug" && ginMode != "" {
		mfaRequired = mfaRequired || staff.ForceMFA()
	}

	// 6. 如果需要 MFA 但还没配 secret，先生成一个（用户首次登录时设置）
	if mfaRequired && staff.MfaSecret == "" {
		if _, err := h.mfaSvc.GenerateSecret(); err != nil {
			log.Warn().Err(err).Msgf("staff %s: failed to generate MFA secret", staff.Email)
		}
		log.Warn().Msgf("staff %s requires MFA but has no secret — please set up TOTP", staff.Email)
	}

	// 7. 如果需要 MFA，返回临时 token（只有 short TTL，且标记需 MFA）
	if mfaRequired {
		mfaToken, err := h.jwtSvc.GenerateStaffToken(staff.ID, staff.Email, staff.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "token generation failed"})
			return
		}
		// 临时 token 存 Redis 标记 "pending_mfa"
		h.rdb.Set(ctx, "staff:mfa_pending:"+staff.Email, mfaToken, 10*time.Minute)

		c.JSON(http.StatusOK, StaffLoginResponse{
			MFARequired: true,
			StaffID:     staff.ID,
			Role:        staff.Role,
			Email:       staff.Email,
		})
		return
	}

	// 8. 正常签发完整 token
	accessToken, err := h.jwtSvc.GenerateStaffToken(staff.ID, staff.Email, staff.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "token generation failed"})
		return
	}
	refreshToken, err := h.jwtSvc.GenerateRefreshToken("staff", staff.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, StaffLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    60 * 60, // 60 min
		StaffID:      staff.ID,
		Role:         staff.Role,
		Email:        staff.Email,
	})
}

// MFAVerifyRequest — POST /staff/mfa/verify
type MFAVerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// MFAVerify — POST /staff/mfa/verify
func (h *StaffAuthHandler) MFAVerify(c *gin.Context) {
	var req MFAVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx := c.Request.Context()
	staff, err := h.staffRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid"})
		return
	}

	// 验证 TOTP
	if !h.mfaSvc.Verify(staff.MfaSecret, req.Code) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid mfa code"})
		return
	}

	// 签发完整 token
	accessToken, _ := h.jwtSvc.GenerateStaffToken(staff.ID, staff.Email, staff.Role)
	refreshToken, _ := h.jwtSvc.GenerateRefreshToken("staff", staff.ID)

	c.JSON(http.StatusOK, StaffLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    60 * 60,
		StaffID:      staff.ID,
		Role:         staff.Role,
		Email:        staff.Email,
	})
}

// RefreshRequest — POST /staff/refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh — POST /staff/refresh
func (h *StaffAuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	claims, err := h.jwtSvc.Parse(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid refresh token"})
		return
	}

	// Staff 还是 User？
	if claims.SubjectType == "staff" {
		staff, err := h.staffRepo.GetByID(c.Request.Context(), claims.SubjectID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "staff not found"})
			return
		}
		accessToken, _ := h.jwtSvc.GenerateStaffToken(staff.ID, staff.Email, staff.Role)
		newRefresh, _ := h.jwtSvc.GenerateRefreshToken("staff", staff.ID)
		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": newRefresh,
			"token_type":    "Bearer",
			"expires_in":    60 * 60,
		})
		return
	}

	// User refresh（Step 6 后半部分完整实现在 user_auth.go，这里先占位）
	c.JSON(http.StatusNotImplemented, gin.H{"code": 501, "message": "user refresh not implemented yet"})
}

// Logout — POST /staff/logout（当前 JWT → 黑名单）
func (h *StaffAuthHandler) Logout(c *gin.Context) {
	// 简化实现：返回 200 让客户端删除 token
	// 生产应该把 token 的 jti 加到 Redis 黑名单
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "logged out"})
}

// AuditLogs — GET /staff/audit-logs（从独立 audit PostgreSQL 实例查询）
func (h *StaffAuthHandler) AuditLogs(c *gin.Context) {
	// 解析过滤参数
	staffIDStr := c.Query("staff_id")
	action := c.Query("action")
	targetType := c.Query("target_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "50"))
	if size > 200 {
		size = 200
	}
	if page < 1 {
		page = 1
	}

	q := h.auditDB.Model(&models.AuditLog{})
	if staffIDStr != "" {
		if id, err := strconv.ParseUint(staffIDStr, 10, 64); err == nil {
			q = q.Where("staff_id = ?", id)
		}
	}
	if action != "" {
		q = q.Where("action LIKE ?", "%"+action+"%")
	}
	if targetType != "" {
		q = q.Where("target_type = ?", targetType)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		log.Warn().Err(err).Msg("audit_logs: count failed")
		// auditDB 可能未连接，返回空数组
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}, "total": 0, "page": page, "size": size})
		return
	}

	var items []models.AuditLog
	if err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		log.Warn().Err(err).Msg("audit_logs: query failed")
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}, "total": 0, "page": page, "size": size})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "page": page, "size": size})
}

// 编译期接口检查
var _ context.Context = nil // 确保 import 不警告

// 方便后续取 claims（middleware 层已注入）
var _ = middleware.GetSubjectID
var _ = models.Staff{}

// ==================== Staff 管理（admin/supervisor 用） ====================

// StaffList — GET /staff
func (h *StaffAuthHandler) StaffList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := h.staffRepo.List(c.Request.Context(), page, size)
	if err != nil {
		log.Error().Err(err).Msg("staff list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list, "total": total, "page": page, "size": size})
}

// StaffCreate — POST /staff
func (h *StaffAuthHandler) StaffCreate(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Role        string `json:"role" binding:"required,oneof=admin supervisor advisor tea_farmer operations"`
		WorkTZ      string `json:"work_timezone"`
		AssignFarm  *uint64 `json:"assigned_farm_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 生成初始密码
	initialPwd := strings.ReplaceAll(time.Now().Format("2006"), "-", "") + "Tea!" + strconv.FormatInt(int64(time.Now().UnixNano()%1000), 10)
	hash, _ := h.passwordSvc.Hash(initialPwd)

	tz := req.WorkTZ
	if tz == "" {
		tz = "Europe/London"
	}

	s := models.Staff{
		Name:           req.Name,
		Email:          strings.ToLower(req.Email),
		PasswordHash:   hash,
		Role:           req.Role,
		WorkTimezone:   tz,
		AssignedFarmID: req.AssignFarm,
		IsActive:       true,
	}
	if err := h.staffRepo.Create(c.Request.Context(), &s); err != nil {
		log.Error().Err(err).Msg("staff create failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed (email might exist)"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"staff": s, "initial_password": initialPwd, "note": "share this password securely — require reset on first login"})
}

// StaffToggle — POST /staff/:id/toggle (启用/禁用)
func (h *StaffAuthHandler) StaffToggle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct{ Active *bool `json:"active"` }
	if err := c.ShouldBindJSON(&body); err != nil || body.Active == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide {active: true/false}"})
		return
	}
	if err := h.staffRepo.ToggleActive(c.Request.Context(), id, *body.Active); err != nil {
		log.Error().Err(err).Msg("staff toggle failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "toggle failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "staff updated", "active": *body.Active})
}

// StaffDelete — DELETE /staff/:id (软删除)
func (h *StaffAuthHandler) StaffDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.staffRepo.SoftDelete(c.Request.Context(), id); err != nil {
		log.Error().Err(err).Msg("staff delete failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "staff soft-deleted"})
}

// ==================== Users 管理 ====================

// UserList — GET /users (admin)
func (h *StaffAuthHandler) UserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := h.userRepo.List(c.Request.Context(), page, size)
	if err != nil {
		log.Error().Err(err).Msg("user list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list, "total": total, "page": page, "size": size})
}
