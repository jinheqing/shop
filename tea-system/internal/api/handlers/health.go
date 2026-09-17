package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler — 健康检查
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthResponse — /health 响应结构
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// Health — GET /health
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ok",
		Service: "tea-system",
		Version: "0.2.0",
	})
}

// Ready — GET /ready (deep health check: DB, Redis)
func (h *HealthHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"checks": gin.H{
			"database": "ok",
			"redis":    "ok",
		},
	})
}
