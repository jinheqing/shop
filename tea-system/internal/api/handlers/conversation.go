package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tea-system/internal/middleware"
	"tea-system/internal/models"
	"tea-system/internal/repository"
)

// ConversationHandler — 会话 CRUD
type ConversationHandler struct {
	repo *repository.ConversationRepo
}

func NewConversationHandler(repo *repository.ConversationRepo) *ConversationHandler {
	return &ConversationHandler{repo: repo}
}

// CreateConversationRequest — POST /conversations
type CreateConversationRequest struct {
	// 对方（二选一）
	OtherUserID  *uint64 `json:"other_user_id,omitempty"`
	OtherStaffID *uint64 `json:"other_staff_id,omitempty"`
	Title        string  `json:"title,omitempty"`
}

// List — GET /conversations
func (h *ConversationHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	subID := middleware.GetSubjectID(c)
	subType := middleware.GetSubjectType(c)

	var userID, staffID *uint64
	switch subType {
	case "user":
		userID = &subID
	case "staff":
		staffID = &subID
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid subject type"})
		return
	}

	convs, err := h.repo.ListByParticipantID(ctx, userID, staffID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to list conversations"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": convs})
}

// Create — POST /conversations
func (h *ConversationHandler) Create(c *gin.Context) {
	var req CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 必须指定一个对方
	if req.OtherUserID == nil && req.OtherStaffID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "must specify other_user_id or other_staff_id"})
		return
	}
	if req.OtherUserID != nil && req.OtherStaffID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "only one of other_user_id/other_staff_id allowed"})
		return
	}

	ctx := c.Request.Context()
	subID := middleware.GetSubjectID(c)
	subType := middleware.GetSubjectType(c)

	// 当前用户 participant
	var p1UserID, p1StaffID *uint64
	switch subType {
	case "user":
		p1UserID = &subID
	case "staff":
		p1StaffID = &subID
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid subject type"})
		return
	}

	// 对方 participant
	var p2UserID, p2StaffID *uint64
	if req.OtherUserID != nil {
		p2UserID = req.OtherUserID
	} else {
		p2StaffID = req.OtherStaffID
	}

	conv := &models.Conversation{
		ConversationType: "single",
		Title:            req.Title,
		CreatorID:        subID,
	}

	participants := []models.ConversationParticipant{
		{UserID: p1UserID, StaffID: p1StaffID, Role: "owner"},
		{UserID: p2UserID, StaffID: p2StaffID, Role: "member"},
	}

	if err := h.repo.Create(ctx, conv, participants); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to create conversation"})
		return
	}

	// 重新取带 participants 的完整对象
	full, _ := h.repo.GetByID(ctx, conv.ID)
	c.JSON(http.StatusCreated, gin.H{"code": 0, "data": full})
}

// Get — GET /conversations/:id
func (h *ConversationHandler) Get(c *gin.Context) {
	id, ok := pathUint64(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	ctx := c.Request.Context()
	conv, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "conversation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": conv})
}

// Delete — DELETE /conversations/:id（软删除）
func (h *ConversationHandler) Delete(c *gin.Context) {
	id, ok := pathUint64(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	ctx := c.Request.Context()
	if err := h.repo.SoftDelete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "conversation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "deleted"})
}

// pathUint64 — 从 URL path 取 uint64
func pathUint64(c *gin.Context, key string) (uint64, bool) {
	v, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}
