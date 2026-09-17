package api

import (
	"tea-system/internal/api/handlers"
	"tea-system/internal/config"
	"tea-system/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Handlers — 所有 HTTP handler 聚合（依赖注入风格）
type Handlers struct {
	Health    *handlers.HealthHandler
	StaffAuth *handlers.StaffAuthHandler
	UserAuth  *handlers.UserAuthHandler
}

// Router — 路由注册器
type Router struct {
	cfg    *config.Config
	db     *gorm.DB
	audit  *gorm.DB
	rdb    *redis.Client
	h      *Handlers
	engine *gin.Engine
}

func NewRouter(cfg *config.Config, db *gorm.DB, auditDB *gorm.DB, rdb *redis.Client, h *Handlers) *Router {
	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()
	return &Router{
		cfg:    cfg,
		db:     db,
		audit:  auditDB,
		rdb:    rdb,
		h:      h,
		engine: engine,
	}
}

// Setup — 注册所有路由 + 全局中间件
func (r *Router) Setup() *gin.Engine {
	// ==================== 全局中间件 ====================
	r.engine.Use(middleware.Recovery())
	r.engine.Use(middleware.Logging())
	r.engine.Use(middleware.CORS(nil)) // dev: 允许所有; prod: 传具体域名
	// r.engine.Use(middleware.RateLimitPerIP(120)) // 生产打开

	// ==================== 公开路由（不需要 JWT） ====================
	r.engine.GET("/health", r.h.Health.Health)
	r.engine.GET("/ready", r.h.Health.Ready)
	r.engine.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "tea-system",
			"version": "0.2.0",
			"docs":    "placeholder",
		})
	})

	// API v1
	v1 := r.engine.Group("/api/v1")
	{
		// ---------- 公开 endpoint ----------
		v1.POST("/staff/login", r.h.StaffAuth.Login)
		v1.POST("/staff/mfa/verify", r.h.StaffAuth.MFAVerify)
		v1.POST("/staff/refresh", r.h.StaffAuth.Refresh)

		v1.POST("/user/magic-link/request", r.h.UserAuth.MagicLinkRequest)
		v1.POST("/user/magic-link/verify", r.h.UserAuth.MagicLinkVerify)
		v1.POST("/user/login", r.h.UserAuth.UserLogin)

		// 需要 JWT 的 endpoint
		jwtSecret := r.cfg.JWT.Secret
		auth := v1.Group("")
		auth.Use(middleware.JWTAuth(jwtSecret))
		{
			// 测试端点
			auth.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"sub_type": middleware.GetSubjectType(c),
					"sub_id":   middleware.GetSubjectID(c),
					"email":    middleware.GetEmail(c),
					"role":     middleware.GetRole(c),
				})
			})

			// Staff 登出
			auth.POST("/staff/logout", r.h.StaffAuth.Logout)

			// IM（Step 7）
			// 定制报价（Step 8）
			// 订单（Step 9）
			// 慢直播（Step 12）
			// 直播间（Step 13）
			// LiveKit token（Step 11）
			// 节点（Step 15）
			// CMS（Step 16）
			// DSAR（Step 18）
		}
	}

	// ==================== 404 ====================
	r.engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "route not found",
			"path":    c.Request.URL.Path,
		})
	})

	return r.engine
}

// 防止 jwt 包未被 import 时 go vet 报错
var _ = jwt.SigningMethodHS256
