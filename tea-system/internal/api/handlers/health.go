package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse — /health 响应结构
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// Health — GET /health
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ok",
		Service: "tea-system",
		Version: "0.2.0",
	})
}

// Ready — GET /ready (deep health check: DB, Redis)
func Ready(c *gin.Context) {
	// 生产环境这里要实际 ping DB + Redis，现在先占位
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"checks": gin.H{
			"database": "ok",
			"redis":    "ok",
		},
	})
}
