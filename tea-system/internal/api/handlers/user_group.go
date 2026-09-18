package handlers

import (
	"net/http"
	"strconv"
	"tea-system/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserGroupHandler struct {
	db *gorm.DB
}

func NewUserGroupHandler(db *gorm.DB) *UserGroupHandler {
	return &UserGroupHandler{db: db}
}

func (h *UserGroupHandler) Create(c *gin.Context) {
	var req struct {
		Name        string          `json:"name" binding:"required"`
		Description string          `json:"description"`
		Privileges  models.JSONArray `json:"privileges"`
		AutoRule    models.JSONMap  `json:"auto_rule"`
		Reciprocal  models.JSONArray `json:"reciprocal"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	staffID, _ := c.Get("staff_id")
	g := models.UserGroup{
		Name:        req.Name,
		Description: req.Description,
		Privileges:  req.Privileges,
		AutoRule:    req.AutoRule,
		Reciprocal:  req.Reciprocal,
	}
	if id, ok := staffID.(uint64); ok {
		g.CreatedByStaffID = &id
	}
	if err := h.db.Create(&g).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *UserGroupHandler) List(c *gin.Context) {
	var groups []models.UserGroup
	if err := h.db.Preload("Members").Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": groups})
}

func (h *UserGroupHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var g models.UserGroup
	if err := h.db.Preload("Members.User").First(&g, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "group not found"})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *UserGroupHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var g models.UserGroup
	if err := h.db.First(&g, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "group not found"})
		return
	}
	var req struct {
		Name        *string          `json:"name"`
		Description *string          `json:"description"`
		Privileges  *models.JSONArray `json:"privileges"`
		AutoRule    *models.JSONMap  `json:"auto_rule"`
		Reciprocal  *models.JSONArray `json:"reciprocal"`
	}
	c.ShouldBindJSON(&req)
	if req.Name != nil {
		g.Name = *req.Name
	}
	if req.Description != nil {
		g.Description = *req.Description
	}
	if req.Privileges != nil {
		g.Privileges = *req.Privileges
	}
	if req.AutoRule != nil {
		g.AutoRule = *req.AutoRule
	}
	if req.Reciprocal != nil {
		g.Reciprocal = *req.Reciprocal
	}
	if err := h.db.Save(&g).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *UserGroupHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Where("group_id = ?", id).Delete(&models.UserGroupMember{})
	if err := h.db.Delete(&models.UserGroup{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "deleted"})
}

func (h *UserGroupHandler) AddMember(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	staffID, _ := c.Get("staff_id")
	now := time.Now()
	m := models.UserGroupMember{GroupID: id, UserID: req.UserID, AddedAt: now}
	if sid, ok := staffID.(uint64); ok {
		m.AddedByStaffID = &sid
	}
	if err := h.db.Create(&m).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *UserGroupHandler) RemoveMember(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID, _ := strconv.ParseUint(c.Param("userId"), 10, 64)
	h.db.Where("group_id = ? AND user_id = ?", id, userID).Delete(&models.UserGroupMember{})
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "removed"})
}

func (h *UserGroupHandler) ListMembers(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var members []models.UserGroupMember
	h.db.Where("group_id = ?", id).Preload("User").Find(&members)
	c.JSON(http.StatusOK, members)
}

// BulkAdd 批量加成员（CSV / 数组都行）
func (h *UserGroupHandler) BulkAdd(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		UserIDs []uint64 `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	staffID, _ := c.Get("staff_id")
	now := time.Now()
	for _, uid := range req.UserIDs {
		var existing models.UserGroupMember
		if h.db.Where("group_id = ? AND user_id = ?", id, uid).First(&existing).Error == nil {
			continue // 已存在，跳过
		}
		m := models.UserGroupMember{GroupID: id, UserID: uid, AddedAt: now}
		if sid, ok := staffID.(uint64); ok {
			m.AddedByStaffID = &sid
		}
		h.db.Create(&m)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "bulk added"})
}

// AutoSync — 按 auto_rule 自动归类（简化版：只处理 spend_threshold）
func (h *UserGroupHandler) AutoSync(c *gin.Context) {
	var groups []models.UserGroup
	h.db.Where("auto_rule IS NOT NULL AND auto_rule != '{}'").Find(&groups)
	for _, g := range groups {
		if g.AutoRule == nil {
			continue
		}
		// spend_threshold: 自动把总消费 >= N 的用户加进来
		if threshold, ok := g.AutoRule["spend_threshold"]; ok {
			var thresholdFloat float64
			switch v := threshold.(type) {
			case float64:
				thresholdFloat = v
			case int:
				thresholdFloat = float64(v)
			}
			var users []models.User
			h.db.Find(&users)
			for _, u := range users {
				var total float64
				h.db.Model(&models.Order{}).
					Where("user_id = ? AND status IN ?", u.ID, []string{"paid", "delivered", "completed"}).
					Select("COALESCE(SUM(total_amount), 0)").Scan(&total)
				if total >= thresholdFloat {
					var existing models.UserGroupMember
					if h.db.Where("group_id = ? AND user_id = ?", g.ID, u.ID).First(&existing).Error != nil {
						h.db.Create(&models.UserGroupMember{GroupID: g.ID, UserID: u.ID, AddedAt: time.Now()})
					}
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "auto-sync done"})
}

// ListMyGroups — 用户看自己在哪些组里
func (h *UserGroupHandler) ListMyGroups(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if uid, ok := userID.(uint64); ok {
		var groups []models.UserGroup
		h.db.Joins("JOIN user_group_members m ON m.group_id = user_groups.id").
			Where("m.user_id = ?", uid).
			Find(&groups)
		c.JSON(http.StatusOK, groups)
		return
	}
	c.JSON(http.StatusOK, []models.UserGroup{})
}
