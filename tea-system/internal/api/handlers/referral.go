package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"tea-system/internal/models"
)

// ReferralHandler — 推荐人管理（staff 用）
//
// 公开 API：POST /user/magic-link/request 里已经带了推荐人字段，
// 这里只需要提供 staff 管理界面用到的 CRUD + 推荐人专属短链生成
type ReferralHandler struct {
	db *gorm.DB
}

func NewReferralHandler(db *gorm.DB) *ReferralHandler {
	return &ReferralHandler{db: db}
}

// ReferralListParams — GET /admin/referrals 查询参数
type ReferralListParams struct {
	ReferrerName  string `form:"referrer_name"`
	ReferredName  string `form:"referred_name"`
	RewardOnly    bool   `form:"reward_only"`    // 只看已触发奖励的
	RewardPending bool   `form:"reward_pending"` // 下单了但还没触发奖励的
	Page          int    `form:"page,default:1"`
	PageSize      int    `form:"page_size,default:50"`
}

// List — 推荐关系列表（staff 用，支持搜索）
func (h *ReferralHandler) List(c *gin.Context) {
	var p ReferralListParams
	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 || p.PageSize > 200 {
		p.PageSize = 50
	}

	q := h.db.WithContext(c.Request.Context()).Model(&models.Referral{})
	if p.ReferrerName != "" {
		q = q.Where("referrer_name ILIKE ?", "%"+p.ReferrerName+"%")
	}
	if p.ReferredName != "" {
		q = q.Where("referred_name ILIKE ?", "%"+p.ReferredName+"%")
	}
	if p.RewardOnly {
		q = q.Where("reward_triggered_at IS NOT NULL")
	}
	if p.RewardPending {
		q = q.Where("friend_order_amount > 0 AND reward_triggered_at IS NULL")
	}

	var total int64
	q.Count(&total)

	var items []models.Referral
	offset := (p.Page - 1) * p.PageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(p.PageSize).Find(&items).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"total":     total,
		"items":     items,
		"page":      p.Page,
		"page_size": p.PageSize,
	})
}

// Create — 顾问手动创建推荐关系（名字匹配不上系统里的人时用）
// POST /admin/referrals
type CreateReferralRequest struct {
	ReferredUserID uint64  `json:"referred_user_id" binding:"required"`
	ReferrerUserID *uint64 `json:"referrer_user_id"` // 可选，如果推荐人也在系统里
	ReferrerName   string  `json:"referrer_name" binding:"required"`
	Source         string  `json:"source"` // name_share / short_code / manual
	Notes          string  `json:"notes"`
}

func (h *ReferralHandler) Create(c *gin.Context) {
	var req CreateReferralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 防止重复
	var existing models.Referral
	err := h.db.WithContext(ctx).Where("referred_user_id = ?", req.ReferredUserID).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "this user already has a referrer"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 获取被推荐人的名字
	var referredUser models.User
	if err := h.db.WithContext(ctx).Select("name").First(&referredUser, req.ReferredUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "referred user not found"})
		return
	}

	staffID := uint64(0) // TODO: 从 JWT 里取
	if id, ok := c.Get("staff_id"); ok {
		staffID = id.(uint64)
	}

	var staffIDPtr *uint64
	if staffID > 0 {
		staffIDPtr = &staffID
	}
	r := &models.Referral{
		ReferrerUserID:   req.ReferrerUserID,
		ReferrerName:     req.ReferrerName,
		ReferredUserID:   req.ReferredUserID,
		ReferredName:     referredUser.Name,
		Source:           req.Source,
		CreatedByStaffID: staffIDPtr,
		Notes:            req.Notes,
	}
	if r.Source == "" {
		r.Source = "manual"
	}

	if err := h.db.WithContext(ctx).Create(r).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 同步更新 users 表
	updates := map[string]interface{}{
		"referrer_name":   req.ReferrerName,
		"referral_source": "friend",
	}
	if req.ReferrerUserID != nil {
		updates["referrer_user_id"] = *req.ReferrerUserID
	}
	h.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", req.ReferredUserID).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "item": r})
}

// UpdateReward — 顾问手动标记奖励已触发（比如手动安排了礼物之后）
// POST /admin/referrals/:id/trigger-reward
func (h *ReferralHandler) UpdateReward(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ctx := c.Request.Context()

	var r models.Referral
	if err := h.db.WithContext(ctx).First(&r, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}

	now := time.Now()
	r.RewardTriggeredAt = &now
	if err := h.db.WithContext(ctx).Save(&r).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "item": r})
}

// Delete — 删除推荐关系
func (h *ReferralHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ctx := c.Request.Context()

	var r models.Referral
	if err := h.db.WithContext(ctx).First(&r, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}

	h.db.WithContext(ctx).Delete(&r)

	// 同步清理 users 表里的推荐字段
	if r.ReferredUserID != 0 {
		h.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", r.ReferredUserID).Updates(map[string]interface{}{
			"referrer_user_id": nil,
			"referrer_name":    "",
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0})
}

// MyReferrals — 老客户看自己推荐了谁（公共 API，user 自己用）
// GET /user/referrals
func (h *ReferralHandler) MyReferrals(c *gin.Context) {
	userID := c.GetUint64("user_id")
	ctx := c.Request.Context()

	// 查所有以我为推荐人的推荐关系
	var items []models.Referral
	h.db.WithContext(ctx).Where("referrer_user_id = ?", userID).Order("created_at DESC").Find(&items)

	totalCount := int64(len(items))
	var validCount int64
	var rewardedCount int64
	for _, r := range items {
		if r.FriendOrderAmount > 0 {
			validCount++
		}
		if r.RewardTriggeredAt != nil {
			rewardedCount++
		}
	}

	// 生成专属短链 code（如果还没有的话）
	var user models.User
	h.db.WithContext(ctx).Select("referral_short_code").First(&user, userID)
	if user.ReferralShortCode == "" {
		code := genShortCode(userID)
		h.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("referral_short_code", code)
		user.ReferralShortCode = code
	}

	c.JSON(http.StatusOK, gin.H{
		"code":          0,
		"referral_code": user.ReferralShortCode,
		"referral_url":  "/r/" + user.ReferralShortCode,
		"stats": gin.H{
			"total_referred": totalCount,
			"valid_orders":   validCount,
			"rewarded":       rewardedCount,
		},
		"items": items,
	})
}

// genShortCode — 生成 8 位可读短码（base36，去掉易混字符）
func genShortCode(id uint64) string {
	const charset = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
	b := make([]byte, 8)
	n := id * 1315423911 // FNV-like mixing
	for i := range b {
		b[i] = charset[n%uint64(len(charset))]
		n /= uint64(len(charset))
	}
	return string(b)
}

// UpdateReferralOnOrderPaid — 订单支付成功后调用
// 1. 如果这个下单用户是被推荐人 → 更新 referrals 表的 friend_order_id/amount
// 2. 检查推荐人的累计有效推荐数 → 如果达到 AutoRule 阈值，自动加入 KOL 伙伴组
//
// 推荐门槛金额 (referral_min_order_amount) 和 进组阈值 (referral_count_to_vip)
// 可以通过 SystemSettings 配置；默认 £100 和 3
func UpdateReferralOnOrderPaid(db *gorm.DB, referredUserID uint64, orderID uint64, orderAmount float64) {
	ctx := context.Background()

	// 1. 查这个用户有没有被推荐
	var ref models.Referral
	if err := db.WithContext(ctx).Where("referred_user_id = ?", referredUserID).First(&ref).Error; err != nil {
		return // 没被推荐，跳过
	}

	// 已经有订单了，跳过（幂等）
	if ref.FriendOrderID != nil && *ref.FriendOrderID == orderID {
		return
	}

	// 2. 更新 Referral 记录
	ref.FriendOrderID = &orderID
	ref.FriendOrderAmount = orderAmount
	db.WithContext(ctx).Save(&ref)

	// 3. 查配置
	threshold := 100.0  // 默认 £100
	countThreshold := 3 // 默认 3 个
	var cfg struct {
		Value string `gorm:"column:value"`
	}
	if db.WithContext(ctx).Raw(
		`SELECT value FROM system_config WHERE key = ?`, "referral_min_order_amount",
	).Scan(&cfg).Error == nil && cfg.Value != "" {
		if v, err := strconv.ParseFloat(cfg.Value, 64); err == nil {
			threshold = v
		}
	}
	if db.WithContext(ctx).Raw(
		`SELECT value FROM system_config WHERE key = ?`, "referral_count_to_vip",
	).Scan(&cfg).Error == nil && cfg.Value != "" {
		if v, err := strconv.Atoi(cfg.Value); err == nil {
			countThreshold = v
		}
	}

	// 4. 订单金额没到门槛，结束（礼物和进组都不触发）
	if orderAmount < threshold {
		return
	}

	// 5. 数这个推荐人累计有多少个有效推荐（friend_order_amount >= threshold）
	if ref.ReferrerUserID == nil {
		return // 没有关联到系统里的用户
	}
	var validCount int64
	db.WithContext(ctx).Model(&models.Referral{}).
		Where("referrer_user_id = ? AND friend_order_amount >= ?", *ref.ReferrerUserID, threshold).
		Count(&validCount)

	// 6. 达到进组阈值 → 自动加进 KOL 伙伴组
	if validCount >= int64(countThreshold) {
		// 找 KOL 伙伴组
		var group models.UserGroup
		err := db.WithContext(ctx).Where("name = ?", "KOL 伙伴").First(&group).Error
		if err == nil {
			// 检查是不是已经在组里了
			var existing int64
			db.WithContext(ctx).Model(&models.UserGroupMember{}).
				Where("group_id = ? AND user_id = ?", group.ID, *ref.ReferrerUserID).
				Count(&existing)
			if existing == 0 {
				db.WithContext(ctx).Create(&models.UserGroupMember{
					GroupID: group.ID,
					UserID:  *ref.ReferrerUserID,
					AddedAt: time.Now(),
				})
			}
		}
	}
}
