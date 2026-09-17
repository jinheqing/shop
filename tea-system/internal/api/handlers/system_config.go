package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"tea-system/internal/models"
)

// SystemConfigHandler — 简单 KV 系统配置（用 site_contents 表 page='system_config'）
type SystemConfigHandler struct {
	DB *gorm.DB
}

func NewSystemConfigHandler(db *gorm.DB) *SystemConfigHandler {
	return &SystemConfigHandler{DB: db}
}

// Get — GET /system/config/:key
func (h *SystemConfigHandler) Get(c *gin.Context) {
	key := c.Param("key")
	var sc models.SiteContent
	err := h.DB.Where("page_key = ? AND section_key = ?", "system_config", key).First(&sc).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"key": key, "value": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": key, "value": sc.Content})
}

// Put — PUT /system/config/:key
func (h *SystemConfigHandler) Put(c *gin.Context) {
	key := c.Param("key")
	var body struct {
		Value models.JSONMap `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := uint64(0)
	if tmp := c.GetUint64("subject_id"); tmp > 0 {
		uid = tmp
	}
	var sc models.SiteContent
	err := h.DB.Where("page_key = ? AND section_key = ?", "system_config", key).First(&sc).Error
	if err == gorm.ErrRecordNotFound {
		sc = models.SiteContent{
			PageKey:          "system_config",
			SectionKey:       key,
			Content:          body.Value,
			UpdatedByStaffID: &uid,
		}
		if err := h.DB.Create(&sc).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
			return
		}
	} else {
		if err := h.DB.Model(&sc).Updates(map[string]interface{}{
			"content":            body.Value,
			"updated_by_staff_id": uid,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"key": key, "message": "saved"})
}

// List — GET /system/config (列出所有系统配置 key)
func (h *SystemConfigHandler) List(c *gin.Context) {
	var scs []models.SiteContent
	h.DB.Where("page_key = ?", "system_config").Find(&scs)
	type item struct {
		Key string `json:"key"`
	}
	out := make([]item, 0, len(scs))
	for _, s := range scs {
		out = append(out, item{Key: s.SectionKey})
	}
	c.JSON(http.StatusOK, gin.H{"items": out, "total": len(out)})
}
