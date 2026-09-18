package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole — 要求 staff 角色至少命中一个 roles 中的值
// 如果 JWT 不是 staff（sub_type != "staff"）或 role 不匹配，返回 403
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		subType, _ := c.Get(CtxKeySubjectType)
		if subType != "staff" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 403, "message": "staff only",
			})
			return
		}
		role, _ := c.Get(CtxKeyRole)
		roleStr, _ := role.(string)
		for _, r := range roles {
			if roleStr == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code": 403, "message": "forbidden: role " + roleStr + " cannot access this resource",
		})
	}
}
