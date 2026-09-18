package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logging — zerolog 请求日志（替代 gin.Logger）
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// 跳过 health 检查的 INFO 日志（减少噪音）
		level := zerolog.InfoLevel
		if status >= 500 {
			level = zerolog.ErrorLevel
		} else if status >= 400 {
			level = zerolog.WarnLevel
		}
		if path == "/health" && status == 200 {
			level = zerolog.DebugLevel
		}

		log.WithLevel(level).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", raw).
			Int("status", status).
			Dur("latency", latency).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP request")
	}
}
