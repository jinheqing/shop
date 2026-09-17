package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"tea-system/internal/middleware"
	"tea-system/internal/models"
)

type SiteContentHandler struct {
	DB *gorm.DB
}

func NewSiteContentHandler(db *gorm.DB) *SiteContentHandler { return &SiteContentHandler{DB: db} }

// List — GET /site-contents
func (h *SiteContentHandler) List(c *gin.Context) {
	key := c.Query("page_key")
	var items []models.SiteContent
	q := h.DB.Model(&models.SiteContent{})
	if key != "" {
		q = q.Where("page_key = ?", key)
	}
	if err := q.Order("page_key, section_key").Limit(200).Find(&items).Error; err != nil {
		log.Error().Err(err).Msg("site_content: list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": len(items)})
}

// Update — PUT /site-contents/:id
func (h *SiteContentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Content models.JSONMap `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var sc models.SiteContent
	if err := h.DB.First(&sc, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "site content not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lookup failed"})
		return
	}
	staffID := middleware.GetSubjectID(c)
	h.DB.Model(&sc).Updates(map[string]interface{}{
		"content":             body.Content,
		"updated_by_staff_id": staffID,
		"updated_at":          time.Now().UTC(),
	})
	h.DB.First(&sc, id)
	c.JSON(http.StatusOK, sc)
}
