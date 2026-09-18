package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthClaims — JWT 统一声明结构（Staff 和 User 共用，靠 Role 区分）
type AuthClaims struct {
	SubjectType string `json:"sub_type"` // "staff" | "user"
	SubjectID   uint64 `json:"sub_id"`
	Email       string `json:"email"`
	Role        string `json:"role"` // staff: advisor/supervisor/admin/tea_farmer/operations; user: customer
	jwt.RegisteredClaims
}

// context key 常量
const (
	CtxKeySubjectType = "sub_type"
	CtxKeySubjectID   = "sub_id"
	CtxKeyEmail       = "email"
	CtxKeyRole        = "role"
	CtxKeyIsStaff     = "is_staff"
)

// JWTAuth — JWT 中间件（同时支持 Staff 和 User 双通道）
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "missing authorization header",
			})
			return
		}

		// 支持 "Bearer xxx" 或直接 "xxx"
		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		claims := &AuthClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "invalid or expired token",
			})
			return
		}

		// 检查是否过期
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token expired",
			})
			return
		}

		// 注入 context
		c.Set(CtxKeySubjectType, claims.SubjectType)
		c.Set(CtxKeySubjectID, claims.SubjectID)
		c.Set(CtxKeyEmail, claims.Email)
		c.Set(CtxKeyRole, claims.Role)
		c.Set(CtxKeyIsStaff, claims.SubjectType == "staff")

		c.Next()
	}
}

// RequireStaff — JWT 之后加这层，要求调用者是 Staff
func RequireStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		isStaff, exists := c.Get(CtxKeyIsStaff)
		if !exists || !isStaff.(bool) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "staff access required",
			})
			return
		}
		c.Next()
	}
}

// RequireRoles — 要求 Staff 有指定角色之一
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(CtxKeyRole)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "message": "forbidden"})
			return
		}
		role := roleVal.(string)
		for _, allowed := range roles {
			if role == allowed {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "role not permitted",
		})
	}
}

// 方便 handler 里取值的辅助函数
func GetSubjectID(c *gin.Context) uint64 {
	if v, ok := c.Get(CtxKeySubjectID); ok {
		return v.(uint64)
	}
	return 0
}

func GetSubjectType(c *gin.Context) string {
	if v, ok := c.Get(CtxKeySubjectType); ok {
		return v.(string)
	}
	return ""
}

func GetEmail(c *gin.Context) string {
	if v, ok := c.Get(CtxKeyEmail); ok {
		return v.(string)
	}
	return ""
}

func GetRole(c *gin.Context) string {
	if v, ok := c.Get(CtxKeyRole); ok {
		return v.(string)
	}
	return ""
}
