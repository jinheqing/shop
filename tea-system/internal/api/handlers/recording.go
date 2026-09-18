package handlers

import (
	"net/http"
	"strconv"
	"tea-system/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecordingHandler struct {
	db *gorm.DB
}

func NewRecordingHandler(db *gorm.DB) *RecordingHandler {
	return &RecordingHandler{db: db}
}

func (h *RecordingHandler) List(c *gin.Context) {
	var list []models.Recording
	h.db.Order("created_at DESC").Preload("LiveRoom").Find(&list)
	c.JSON(http.StatusOK, list)
}

func (h *RecordingHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var r models.Recording
	if err := h.db.Preload("LiveRoom").First(&r, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}
	if !h.canAccess(c, &r) {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "not accessible"})
		return
	}
	c.JSON(http.StatusOK, r)
}

func (h *RecordingHandler) UpdateVisibility(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var r models.Recording
	if err := h.db.First(&r, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}
	var req struct {
		Visibility      *string          `json:"visibility"`
		VisibleUserIDs  *models.JSONArray `json:"visible_user_ids"`
		VisibleGroupIDs *models.JSONArray `json:"visible_group_ids"`
	}
	c.ShouldBindJSON(&req)
	if req.Visibility != nil {
		r.Visibility = *req.Visibility
	}
	if req.VisibleUserIDs != nil {
		r.VisibleUserIDs = *req.VisibleUserIDs
	}
	if req.VisibleGroupIDs != nil {
		r.VisibleGroupIDs = *req.VisibleGroupIDs
	}
	h.db.Save(&r)
	c.JSON(http.StatusOK, r)
}

func (h *RecordingHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Delete(&models.Recording{}, id)
	c.JSON(http.StatusOK, gin.H{"code": 0})
}

func (h *RecordingHandler) ListMine(c *gin.Context) {
	userIDRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "login required"})
		return
	}
	userID, _ := userIDRaw.(uint64)

	var list []models.Recording
	h.db.Preload("LiveRoom").Raw(`
		SELECT * FROM recordings
		WHERE visibility = 'public'
		   OR visibility = 'registered'
		   OR (visibility = 'restricted'
		       AND (visible_user_ids @> ? OR EXISTS (
		           SELECT 1 FROM user_group_members m
		           WHERE m.user_id = ?
		           AND m.group_id = ANY(string_to_array(array_to_string(visible_group_ids, ','), ',')::bigint[])
		       ))
		   )
		ORDER BY created_at DESC
	`, []string{strconv.FormatUint(userID, 10)}, userID).Scan(&list)
	c.JSON(http.StatusOK, list)
}

func (h *RecordingHandler) canAccess(c *gin.Context, r *models.Recording) bool {
	if r.Visibility == "public" {
		return true
	}
	if r.Visibility == "registered" {
		_, ok := c.Get("user_id")
		return ok
	}
	if r.Visibility != "restricted" {
		return false
	}
	userIDRaw, ok := c.Get("user_id")
	if !ok {
		return false
	}
	userID, _ := userIDRaw.(uint64)

	for _, uid := range r.VisibleUserIDs {
		if v, err := strconv.ParseUint(uid, 10, 64); err == nil && v == userID {
			return true
		}
	}
	var userGroups []models.UserGroupMember
	h.db.Where("user_id = ?", userID).Find(&userGroups)
	set := map[uint64]bool{}
	for _, m := range userGroups {
		set[m.GroupID] = true
	}
	for _, gid := range r.VisibleGroupIDs {
		if v, err := strconv.ParseUint(gid, 10, 64); err == nil && set[v] {
			return true
		}
	}
	return false
}
