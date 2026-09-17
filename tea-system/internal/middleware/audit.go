package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"tea-system/internal/models"
)

// Audit — 审计日志中间件（写入独立 audit_logs 库）
// 只对写操作记录（POST/PUT/PATCH/DELETE），GET 跳过
// 生产环境这应该是异步的（消息队列），这里先同步实现
func Audit(auditDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 只读请求跳过
		method := c.Request.Method
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			return
		}

		// 从 JWT 上下文取 Staff 信息
		var staffID *uint64
		subType, _ := c.Get(CtxKeySubjectType)
		if subType == "staff" {
			if id, ok := c.Get(CtxKeySubjectID); ok {
				sid := id.(uint64)
				staffID = &sid
			}
		}

		// 推断 action: POST /custom-products → custom_products.create
		path := c.Request.URL.Path
		action := mapMethodPathToAction(method, path)

		entry := models.AuditLog{
			StaffID:    staffID,
			Action:     action,
			TargetType: inferTargetType(path),
			TargetID:   extractTargetID(path),
			IPAddress:  c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			CreatedAt:  time.Now().UTC(),
		}

		if err := auditDB.Create(&entry).Error; err != nil {
			log.Warn().Err(err).Msg("audit log write failed (non-blocking)")
		}
	}
}

func mapMethodPathToAction(method, path string) string {
	// 简化版: method + path → action
	// POST /custom-products → custom_products.create
	// PUT /custom-products/:id → custom_products.update
	// POST /staff/login → staff.login
	return method + " " + path
}

func inferTargetType(path string) string {
	// /api/v1/custom-products → custom_products
	segments := splitPath(path)
	if len(segments) >= 2 {
		return segments[len(segments)-2]
	}
	return ""
}

func extractTargetID(path string) *uint64 {
	// 简化实现：不解析
	return nil
}

func splitPath(p string) []string {
	var out []string
	cur := ""
	for _, r := range p {
		if r == '/' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
		} else {
			cur += string(r)
		}
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}
