package api

import (
	"tea-system/internal/api/handlers"
	"tea-system/internal/config"
	"tea-system/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Router — 路由注册器（依赖注入风格）
type Router struct {
	cfg    *config.Config
	db     *gorm.DB
	audit  *gorm.DB
	engine *gin.Engine
}

func NewRouter(cfg *config.Config, db *gorm.DB, auditDB *gorm.DB, jwtSecret string) *Router {
	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()
	return &Router{
		cfg:    cfg,
		db:     db,
		audit:  auditDB,
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
	r.engine.GET("/health", handlers.Health)
	r.engine.GET("/ready", handlers.Ready)
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
		// 公开 endpoint（魔法链接申请、public/sgs、public/slow-presets 等）
		// Step 6 实现
		// public.POST("/user/magic-link/request", ...)
		// public.POST("/user/magic-link/verify", ...)
		// public.GET("/public/sgs-reports", ...)
		// public.GET("/public/slow-presets", ...)

		// 需要 JWT 的 endpoint
		jwtSecret := r.cfg.JWT.Secret
		auth := v1.Group("")
		auth.Use(middleware.JWTAuth(jwtSecret))
		{
			// 测试端点（未带 JWT → 401；带有效 JWT → 404 因为还没注册具体路由）
			auth.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"sub_type": middleware.GetSubjectType(c),
					"sub_id":   middleware.GetSubjectID(c),
					"email":    middleware.GetEmail(c),
					"role":     middleware.GetRole(c),
				})
			})

			// IM（Step 7）
			// auth.GET("/conversations", ...)
			// auth.POST("/conversations", ...)
			// auth.GET("/conversations/:id", ...)

			// 定制报价（Step 8）
			// staffOnly := auth.Group("")
			// staffOnly.Use(middleware.RequireStaff())
			// staffOnly.POST("/custom-products", ...)

			// 订单（Step 9）
			// auth.POST("/orders", ...)
			// auth.GET("/orders", ...)

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

// 防止 jwt 包未被 import 时 go vet 报错（某些构建标签场景）
var _ = jwt.SigningMethodHS256
