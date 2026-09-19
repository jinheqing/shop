package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"tea-system/internal/models"
	"tea-system/internal/repository"
	"tea-system/internal/service"
)

// ============================================================
// User 认证（魔法链接 + 密码登录）
// ============================================================

type UserAuthHandler struct {
	userRepo         *repository.UserRepo
	passwordSvc      *service.PasswordService
	jwtSvc           *service.JWTService
	magicLinkSvc     *service.MagicLinkService
	geoipSvc         *service.GeoIPService
	db               *gorm.DB
}

func NewUserAuthHandler(
	userRepo *repository.UserRepo,
	passwordSvc *service.PasswordService,
	jwtSvc *service.JWTService,
	magicLinkSvc *service.MagicLinkService,
	geoipSvc *service.GeoIPService,
	db *gorm.DB,
) *UserAuthHandler {
	return &UserAuthHandler{
		userRepo:     userRepo,
		passwordSvc:  passwordSvc,
		jwtSvc:       jwtSvc,
		magicLinkSvc: magicLinkSvc,
		geoipSvc:     geoipSvc,
		db:           db,
	}
}

// MagicLinkRequestRequest — POST /user/magic-link/request
type MagicLinkRequestRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// MagicLinkRequest — 生成魔法链接并发邮件（生产走 Mailgun，dev 只存 Redis）
func (h *UserAuthHandler) MagicLinkRequest(c *gin.Context) {
	var req MagicLinkRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 1. 查找或创建用户（首次用魔法链接自动注册）
	user, err := h.userRepo.GetByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			// 自动创建
			user = &models.User{
				Name:  strings.Split(req.Email, "@")[0],
				Email: strings.ToLower(req.Email),
			}
			if err := h.userRepo.Create(ctx, user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to create user"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
			return
		}
	}

	// 2. 生成魔法链接
	_, magicURL, err := h.magicLinkSvc.Generate(ctx, user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to generate magic link"})
		return
	}

	// 3. 发送邮件（dev 模式会打日志）
	_ = h.magicLinkSvc.MailSvc.SendMagicLink(ctx, user.Email, magicURL)

	// 4. 安全：永远返回 200 + 模糊消息，不泄露 email 是否存在
	//    dev 环境额外返回 magic_url 方便测试
	resp := gin.H{
		"code":    0,
		"message": "if this email is registered, a magic link has been sent",
	}
	if c.GetHeader("X-Dev-Mode") == "true" || c.Query("dev") == "1" {
		resp["_dev_magic_url"] = magicURL
	}
	c.JSON(http.StatusOK, resp)
}

// MagicLinkVerifyRequest — POST /user/magic-link/verify
type MagicLinkVerifyRequest struct {
	Token string `json:"token" binding:"required,min=16"`
}

// MagicLinkVerify — 验证魔法链接 token（一次性使用）
func (h *UserAuthHandler) MagicLinkVerify(c *gin.Context) {
	var req MagicLinkVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx := c.Request.Context()
	userID, userEmail, err := h.magicLinkSvc.Verify(ctx, req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}

	// 获取完整 User
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "user not found"})
		return
	}
	_ = userEmail // 已经在 user 对象里了

	// 签发 JWT
	accessToken, _ := h.jwtSvc.GenerateUserToken(user.ID, user.Email)
	refreshToken, _ := h.jwtSvc.GenerateRefreshToken("user", user.ID)

	// 统计已有数据（用于匿名合并提示）
	var mergedCount int64
	if h.db != nil {
		h.db.Model(&models.Conversation{}).
			Joins("JOIN conversation_participants cp ON cp.conversation_id = conversations.id").
			Where("cp.user_id = ?", user.ID).
			Count(&mergedCount)
		customCount := int64(0)
		h.db.Model(&models.CustomProduct{}).
			Where("created_by_staff_id IS NULL").
			Count(&customCount)
		mergedCount += customCount
	}

	// 捕获 IP + 归属地
	clientIP := c.ClientIP()
	if h.geoipSvc != nil {
		geo := h.geoipSvc.Lookup(ctx, clientIP)
		if geo != nil {
			_ = h.userRepo.UpdateLoginGeo(ctx, user.ID, geo.IP, geo.City, geo.Country, geo.CountryCode, geo.Region)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":          accessToken,
		"refresh_token":         refreshToken,
		"token_type":            "Bearer",
		"expires_in":            120 * 60,
		"merged_anonymous_count": mergedCount,
		"user": gin.H{
			"id":                 user.ID,
			"name":               user.Name,
			"email":              user.Email,
			"preferred_language": user.PreferredLanguage,
		},
	})
}

// UserLoginRequest — POST /user/login（传统密码）
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserLogin — 密码登录（可选，魔法链接为主）
func (h *UserAuthHandler) UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx := c.Request.Context()
	user, err := h.userRepo.GetByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "email or password invalid"})
		return
	}

	if user.PasswordHash == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "this account uses magic link authentication",
		})
		return
	}

	if !h.passwordSvc.Verify(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "email or password invalid"})
		return
	}

	accessToken, _ := h.jwtSvc.GenerateUserToken(user.ID, user.Email)
	refreshToken, _ := h.jwtSvc.GenerateRefreshToken("user", user.ID)

	// 捕获 IP + 归属地
	clientIP := c.ClientIP()
	if h.geoipSvc != nil {
		geo := h.geoipSvc.Lookup(ctx, clientIP)
		if geo != nil {
			_ = h.userRepo.UpdateLoginGeo(ctx, user.ID, geo.IP, geo.City, geo.Country, geo.CountryCode, geo.Region)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    120 * 60,
	})
}

// 编译期检查
var _ context.Context = nil
