package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"tea-system/internal/middleware"
	"tea-system/internal/models"
)

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

// Put — PUT /system/config/:key  (Upsert via ON CONFLICT page_key+section_key)
func (h *SystemConfigHandler) Put(c *gin.Context) {
	key := c.Param("key")
	var body struct {
		Value models.JSONMap `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var staffID *uint64
	if id := middleware.GetSubjectID(c); id > 0 {
		staffID = &id
	}

	sc := models.SiteContent{
		PageKey:          "system_config",
		SectionKey:       key,
		Content:          body.Value,
		UpdatedByStaffID: staffID,
	}

	result := h.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "page_key"}, {Name: "section_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"content", "updated_by_staff_id", "updated_at"}),
	}).Create(&sc)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("system_config upsert failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed: " + result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": key, "message": "saved"})
}

// List — GET /system/config
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
