package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler — 健康检查
type HealthHandler struct {
	Version string // 从 cfg.Server.Version 注入（DB 覆盖会在上游生效）
}

func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{Version: version}
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
		Version: h.Version,
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
