package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tea-system/internal/repository"
)

// MessageHandler — 消息列表 HTTP
type MessageHandler struct {
	repo *repository.MessageRepo
}

func NewMessageHandler(repo *repository.MessageRepo) *MessageHandler {
	return &MessageHandler{repo: repo}
}

// List — GET /conversations/:id/messages?limit=50&before_id=xxx
func (h *MessageHandler) List(c *gin.Context) {
	convID, ok := pathUint64(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid conversation id"})
		return
	}

	// limit 默认 50，最大 200
	limit := 50
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > 200 {
		limit = 200
	}

	var beforeID *uint64
	if v := c.Query("before_id"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			beforeID = &n
		}
	}

	ctx := c.Request.Context()
	msgs, err := h.repo.ListByConversationID(ctx, convID, limit, beforeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to list messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"messages": msgs,
			"count":    len(msgs),
		},
	})
}
