package handlers

import (
	"net/http"
	"strconv"
	"tea-system/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShortLinkHandler struct {
	db *gorm.DB
}

func NewShortLinkHandler(db *gorm.DB) *ShortLinkHandler {
	return &ShortLinkHandler{db: db}
}

func (h *ShortLinkHandler) Create(c *gin.Context) {
	var req struct {
		TargetURL string     `json:"target_url" binding:"required"`
		ExpiresAt *time.Time `json:"expires_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	ownerID, _ := c.Get("user_id")
	code := models.GenerateShortCode()
	// 碰撞重试（6 字符空间够大，基本不会）
	for i := 0; i < 3; i++ {
		var existing models.ShortLink
		if h.db.Where("code = ?", code).First(&existing).Error != nil {
			break
		}
		code = models.GenerateShortCode()
	}
	sl := models.ShortLink{Code: code, TargetURL: req.TargetURL, ExpiresAt: req.ExpiresAt}
	if uid, ok := ownerID.(uint64); ok {
		sl.OwnerUserID = &uid
	}
	if err := h.db.Create(&sl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sl)
}

// Redirect — GET /s/:code → 302 到 target
func (h *ShortLinkHandler) Redirect(c *gin.Context) {
	code := c.Param("code")
	var sl models.ShortLink
	if err := h.db.Where("code = ?", code).First(&sl).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "short link not found"})
		return
	}
	if sl.ExpiresAt != nil && time.Now().After(*sl.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"code": 410, "message": "link expired"})
		return
	}
	h.db.Model(&sl).Update("click_count", gorm.Expr("click_count + 1"))
	c.Redirect(http.StatusFound, sl.TargetURL)
}

func (h *ShortLinkHandler) List(c *gin.Context) {
	var list []models.ShortLink
	h.db.Order("created_at DESC").Limit(200).Find(&list)
	c.JSON(http.StatusOK, list)
}

func (h *ShortLinkHandler) Get(c *gin.Context) {
	code := c.Param("code")
	var sl models.ShortLink
	if err := h.db.Where("code = ?", code).First(&sl).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}
	c.JSON(http.StatusOK, sl)
}

func (h *ShortLinkHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Delete(&models.ShortLink{}, id)
	c.JSON(http.StatusOK, gin.H{"code": 0})
}

// TopStats — 被分享最多的前 10
func (h *ShortLinkHandler) TopStats(c *gin.Context) {
	var list []models.ShortLink
	h.db.Order("click_count DESC").Limit(10).Find(&list)
	c.JSON(http.StatusOK, list)
}
