package handlers

import (
	"net/http"
	"strconv"
	"tea-system/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// VideoHandler — 视频发布系统 CRUD
type VideoHandler struct {
	db *gorm.DB
}

func NewVideoHandler(db *gorm.DB) *VideoHandler {
	return &VideoHandler{db: db}
}

// ===== 视频分类 =====

func (h *VideoHandler) ListCategories(c *gin.Context) {
	var cats []models.VideoCategory
	h.db.Order("sort_order ASC, id ASC").Find(&cats)
	c.JSON(http.StatusOK, cats)
}

func (h *VideoHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	cat := models.VideoCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}
	if err := h.db.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *VideoHandler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var cat models.VideoCategory
	if err := h.db.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}
	var req struct {
		Name        *string `json:"name"`
		Slug        *string `json:"slug"`
		Description *string `json:"description"`
		SortOrder   *int    `json:"sort_order"`
	}
	c.ShouldBindJSON(&req)
	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.Slug != nil {
		cat.Slug = *req.Slug
	}
	if req.Description != nil {
		cat.Description = *req.Description
	}
	if req.SortOrder != nil {
		cat.SortOrder = *req.SortOrder
	}
	h.db.Save(&cat)
	c.JSON(http.StatusOK, cat)
}

func (h *VideoHandler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	// 把该分类下的视频 category_id 置空
	h.db.Model(&models.Video{}).Where("category_id = ?", id).Update("category_id", nil)
	h.db.Delete(&models.VideoCategory{}, id)
	c.JSON(http.StatusOK, gin.H{"code": 0})
}

// ===== 视频 CRUD =====

func (h *VideoHandler) List(c *gin.Context) {
	var videos []models.Video
	query := h.db.Order("sort_order ASC, created_at DESC")
	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("category_id = ?", catID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	query.Preload("Category").Find(&videos)
	c.JSON(http.StatusOK, videos)
}

func (h *VideoHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var v models.Video
	if err := h.db.Preload("Category").First(&v, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *VideoHandler) Create(c *gin.Context) {
	var req struct {
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description"`
		VideoURL    string  `json:"video_url" binding:"required"`
		CoverURL    string  `json:"cover_url"`
		DurationSec int     `json:"duration_sec"`
		CategoryID  *uint64 `json:"category_id"`
		SortOrder   int     `json:"sort_order"`
		Status      string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	v := models.Video{
		Title:       req.Title,
		Description: req.Description,
		VideoURL:    req.VideoURL,
		CoverURL:    req.CoverURL,
		DurationSec: req.DurationSec,
		CategoryID:  req.CategoryID,
		SortOrder:    req.SortOrder,
		Status:       req.Status,
	}
	if v.Status == "" {
		v.Status = "draft"
	}
	if err := h.db.Create(&v).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	h.db.Preload("Category").First(&v, v.ID)
	c.JSON(http.StatusCreated, v)
}

func (h *VideoHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var v models.Video
	if err := h.db.First(&v, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}
	var req struct {
		Title       *string  `json:"title"`
		Description *string  `json:"description"`
		VideoURL    *string  `json:"video_url"`
		CoverURL    *string  `json:"cover_url"`
		DurationSec *int     `json:"duration_sec"`
		CategoryID  *uint64  `json:"category_id"`
		SortOrder   *int     `json:"sort_order"`
		Status      *string  `json:"status"`
	}
	c.ShouldBindJSON(&req)
	if req.Title != nil {
		v.Title = *req.Title
	}
	if req.Description != nil {
		v.Description = *req.Description
	}
	if req.VideoURL != nil {
		v.VideoURL = *req.VideoURL
	}
	if req.CoverURL != nil {
		v.CoverURL = *req.CoverURL
	}
	if req.DurationSec != nil {
		v.DurationSec = *req.DurationSec
	}
	if req.CategoryID != nil {
		v.CategoryID = req.CategoryID
	}
	if req.SortOrder != nil {
		v.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		v.Status = *req.Status
	}
	h.db.Save(&v)
	h.db.Preload("Category").First(&v, v.ID)
	c.JSON(http.StatusOK, v)
}

func (h *VideoHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Delete(&models.Video{}, id)
	c.JSON(http.StatusOK, gin.H{"code": 0})
}

// ===== 公开接口 =====

// PublishedList — 前台公开视频列表（只返回 published）
func (h *VideoHandler) PublishedList(c *gin.Context) {
	var videos []models.Video
	query := h.db.Where("status = ?", "published").Order("sort_order ASC, created_at DESC")
	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("category_id = ?", catID)
	}
	query.Preload("Category").Find(&videos)
	c.JSON(http.StatusOK, gin.H{"items": videos})
}

// PublishedCategories — 前台公开分类列表（只返回有已发布视频的分类）
func (h *VideoHandler) PublishedCategories(c *gin.Context) {
	var cats []models.VideoCategory
	h.db.Where("EXISTS (SELECT 1 FROM videos WHERE videos.category_id = video_categories.id AND videos.status = 'published')").
		Order("sort_order ASC, id ASC").Find(&cats)
	c.JSON(http.StatusOK, cats)
}
