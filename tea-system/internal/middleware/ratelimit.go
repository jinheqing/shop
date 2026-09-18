package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ratelimit 简单令牌桶实现（进程内，后续可替换 Redis 版本）

type tokenBucket struct {
	tokens    float64
	lastRefill time.Time
}

var (
	rlBuckets   = make(map[string]*tokenBucket)
	rlMu        sync.Mutex
)

// RateLimit — 令牌桶限流
// rate = 每秒允许多少请求，burst = 桶容量
func RateLimit(rate float64, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + c.Request.URL.Path

		rlMu.Lock()
		b, exists := rlBuckets[key]
		if !exists {
			b = &tokenBucket{tokens: float64(burst), lastRefill: time.Now()}
			rlBuckets[key] = b
		}
		// 补充令牌
		now := time.Now()
		elapsed := now.Sub(b.lastRefill).Seconds()
		b.tokens += elapsed * rate
		if b.tokens > float64(burst) {
			b.tokens = float64(burst)
		}
		b.lastRefill = now

		if b.tokens < 1 {
			rlMu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "rate limit exceeded",
			})
			return
		}
		b.tokens -= 1
		rlMu.Unlock()

		c.Next()
	}
}

// RateLimitPerIP — 每个 IP 独立桶，默认 60 req/min
func RateLimitPerIP(rpm int) gin.HandlerFunc {
	return RateLimit(float64(rpm)/60.0, rpm)
}
