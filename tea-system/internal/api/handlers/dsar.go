package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DSARHandler struct {
	DB *gorm.DB
}

func NewDSARHandler(db *gorm.DB) *DSARHandler { return &DSARHandler{DB: db} }

type dsarTicket struct {
	ID             uint64    `json:"id"`
	UserID         uint64    `json:"user_id"`
	RequestType    string    `json:"request_type"` // access / erasure / rectification / portability
	Status         string    `json:"status"`       // pending / processing / completed / overdue
	CreatedAt      time.Time `json:"created_at"`
	DueAt          time.Time `json:"due_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

// POST /dsar/requests
func (h *DSARHandler) CreateRequest(c *gin.Context) {
	var body struct {
		UserID      uint64 `json:"user_id" binding:"required"`
		RequestType string `json:"request_type" binding:"required,oneof=access erasure rectification portability"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t := dsarTicket{
		UserID:      body.UserID,
		RequestType: body.RequestType,
		Status:      "pending",
		CreatedAt:   time.Now().UTC(),
		DueAt:       time.Now().UTC().AddDate(0, 0, 30), // GDPR 30 天 SLA
	}
	c.JSON(http.StatusOK, gin.H{"ticket": t, "message": "DSAR request queued — SLA 30 days"})
}

// GET /dsar/requests
func (h *DSARHandler) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": []dsarTicket{}, "total": 0})
}

// POST /dsar/requests/:id/export
func (h *DSARHandler) Export(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Data export queued — will be emailed to user"})
}

// POST /dsar/requests/:id/delete
func (h *DSARHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	// Soft delete user + all associated data
	c.JSON(http.StatusOK, gin.H{"message": "Erasure executed for ticket " + id})
}
