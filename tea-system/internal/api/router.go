package api

import (
	"tea-system/internal/api/handlers"
	"tea-system/internal/config"
	"tea-system/internal/middleware"

	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Handlers — 所有 HTTP handler 聚合
type Handlers struct {
	Health        *handlers.HealthHandler
	StaffAuth     *handlers.StaffAuthHandler
	UserAuth      *handlers.UserAuthHandler
	CustomProduct *handlers.CustomProductHandler
	LiveKit       *handlers.LiveKitTokenHandler
	Conversation  *handlers.ConversationHandler
	Message       *handlers.MessageHandler
	IMWS          *handlers.IMWSHandler

	// Step 9+10
	Order       *handlers.OrderHandler
	Payment     *handlers.PaymentHandler
	Invoice     *handlers.InvoiceHandler
	Declaration *handlers.DeclarationHandler
	Ledger      *handlers.LedgerHandler
	SgsReport   *handlers.SgsReportHandler

	// Step 12+13+15
	SlowPreset *handlers.SlowPresetHandler
	LiveRoom   *handlers.LiveRoomHandler
	Node       *handlers.NodeHandler
}

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
	return &Router{cfg: cfg, db: db, audit: auditDB, rdb: rdb, h: h, engine: engine}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.Use(middleware.Recovery())
	r.engine.Use(middleware.Logging())
	r.engine.Use(middleware.CORS(nil))

	// ============ 公开路由 ============
	r.engine.GET("/health", r.h.Health.Health)
	r.engine.GET("/ready", r.h.Health.Ready)
	r.engine.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"service": "tea-system", "version": "0.4.0", "docs": "placeholder"})
	})

	// WebSocket（handler 内自己解析 JWT）
	r.engine.GET("/ws/im", r.h.IMWS.Serve)

	v1 := r.engine.Group("/api/v1")
	{
		// ---------- 认证 ----------
		v1.POST("/staff/login", r.h.StaffAuth.Login)
		v1.POST("/staff/mfa/verify", r.h.StaffAuth.MFAVerify)
		v1.POST("/staff/refresh", r.h.StaffAuth.Refresh)
		v1.POST("/user/magic-link/request", r.h.UserAuth.MagicLinkRequest)
		v1.POST("/user/magic-link/verify", r.h.UserAuth.MagicLinkVerify)
		v1.POST("/user/login", r.h.UserAuth.UserLogin)

		// ---------- 公开 API ----------
		v1.GET("/custom-products/by-token/:token", r.h.CustomProduct.GetByToken)
		v1.GET("/public/slow-presets", r.h.SlowPreset.PublicList)
		v1.POST("/live-rooms/customer-request", r.h.LiveRoom.CustomerRequest)
		v1.GET("/public/sgs-reports", r.h.SgsReport.PublicList)

		// ---------- 支付 Webhook ----------
		v1.POST("/webhooks/2checkout", r.h.Payment.Handle2CheckoutWebhook)
		v1.POST("/webhooks/paypal", r.h.Payment.HandlePayPalWebhook)

		// 需要 JWT 的
		auth := v1.Group("")
		auth.Use(middleware.JWTAuth(r.cfg.JWT.Secret))
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
			auth.POST("/staff/logout", r.h.StaffAuth.Logout)
			auth.GET("/staff/audit-logs", func(c *gin.Context) {
				type AuditLogRow struct {
					ID         uint64                 `gorm:"column:id" json:"id"`
					StaffID    *uint64                `gorm:"column:staff_id" json:"staff_id"`
					Action     string                 `gorm:"column:action" json:"action"`
					TargetType string                 `gorm:"column:target_type" json:"target_type"`
					TargetID   *uint64                `gorm:"column:target_id" json:"target_id"`
					Detail     map[string]interface{} `gorm:"column:detail" json:"detail"`
					IPAddress  string                 `gorm:"column:ip_address" json:"ip_address"`
					UserAgent  string                 `gorm:"column:user_agent" json:"user_agent"`
					CreatedAt  string                 `gorm:"column:created_at" json:"created_at"`
				}
				var rows []AuditLogRow
				if r.audit != nil {
					r.audit.Raw(`SELECT id, staff_id, action, target_type, target_id, detail, ip_address, user_agent, created_at FROM audit_logs ORDER BY created_at DESC LIMIT 50`).Scan(&rows)
				}
				c.JSON(http.StatusOK, gin.H{"items": rows, "total": len(rows)})
			})

			// Step 7: IM
			auth.GET("/conversations", r.h.Conversation.List)
			auth.POST("/conversations", r.h.Conversation.Create)
			auth.GET("/conversations/:id", r.h.Conversation.Get)
			auth.DELETE("/conversations/:id", r.h.Conversation.Delete)
			auth.GET("/conversations/:id/messages", r.h.Message.List)

			// Step 8: 定制报价
			auth.POST("/custom-products", r.h.CustomProduct.Create)
			auth.GET("/custom-products", r.h.CustomProduct.List)
			auth.GET("/custom-products/:id", r.h.CustomProduct.GetByID)
			auth.PUT("/custom-products/:id", r.h.CustomProduct.Update)
			auth.DELETE("/custom-products/:id", r.h.CustomProduct.Delete)
			auth.POST("/custom-products/:id/publish", r.h.CustomProduct.Publish)
			auth.POST("/custom-products/:id/review", r.h.CustomProduct.Review)

			// Step 11: LiveKit
			auth.POST("/livekit/token", r.h.LiveKit.Token)
			auth.POST("/livekit/token-for-obs", r.h.LiveKit.TokenForOBS)
			auth.POST("/livekit/rooms", r.h.LiveKit.CreateRoom)
			auth.DELETE("/livekit/rooms/:name", r.h.LiveKit.DeleteRoom)

			// Step 9: 订单
			auth.POST("/orders", r.h.Order.Create)
			auth.GET("/orders", r.h.Order.List)
			auth.GET("/orders/:id", r.h.Order.GetByID)
			auth.POST("/orders/:id/state", r.h.Order.UpdateState)
			auth.POST("/orders/:id/cancel", r.h.Order.Cancel)
			auth.POST("/orders/:id/payment/init", r.h.Payment.Init)

			// Step 10: 发票 + 台账 + SGS
			auth.GET("/orders/:id/invoice", r.h.Invoice.GetByOrder)
			auth.GET("/orders/:id/invoice/pdf", r.h.Invoice.DownloadPDF)
			auth.POST("/orders/:id/invoice/regenerate", r.h.Invoice.Regenerate)
			auth.POST("/declarations", r.h.Declaration.Create)
			auth.GET("/declarations", r.h.Declaration.List)
			auth.PUT("/declarations/:id", r.h.Declaration.Update)
			auth.DELETE("/declarations/:id", r.h.Declaration.Delete)
			auth.POST("/ledgers", r.h.Ledger.Create)
			auth.GET("/ledgers", r.h.Ledger.List)
			auth.PUT("/ledgers/:id", r.h.Ledger.Update)
			auth.DELETE("/ledgers/:id", r.h.Ledger.Delete)
			auth.POST("/sgs-reports", r.h.SgsReport.Create)
			auth.GET("/sgs-reports", r.h.SgsReport.List)
			auth.GET("/sgs-reports/:id", r.h.SgsReport.GetByID)
			auth.PUT("/sgs-reports/:id", r.h.SgsReport.Update)
			auth.DELETE("/sgs-reports/:id", r.h.SgsReport.Delete)

			// Step 12: 慢直播
			auth.POST("/slow-presets", r.h.SlowPreset.Create)
			auth.GET("/slow-presets", r.h.SlowPreset.List)
			auth.PUT("/slow-presets/:id", r.h.SlowPreset.Update)
			auth.DELETE("/slow-presets/:id", r.h.SlowPreset.Delete)
			auth.POST("/slow-presets/:id/enable", r.h.SlowPreset.Enable)
			auth.POST("/slow-presets/:id/disable", r.h.SlowPreset.Disable)

			// Step 13: 直播间
			auth.POST("/live-rooms", r.h.LiveRoom.Create)
			auth.GET("/live-rooms", r.h.LiveRoom.List)
			auth.GET("/live-rooms/:id", r.h.LiveRoom.GetByID)
			auth.PUT("/live-rooms/:id", r.h.LiveRoom.Update)
			auth.DELETE("/live-rooms/:id", r.h.LiveRoom.Delete)
			auth.POST("/live-rooms/:id/start", r.h.LiveRoom.Start)
			auth.POST("/live-rooms/:id/end", r.h.LiveRoom.End)
			auth.GET("/live-rooms/calendar", r.h.LiveRoom.Calendar)
			auth.GET("/live-rooms/customer-requests", r.h.LiveRoom.CustomerRequests)

			// Step 15: 节点
			auth.POST("/nodes", r.h.Node.Create)
			auth.POST("/nodes/deploy", r.h.Node.Deploy)
			auth.GET("/nodes", r.h.Node.List)
			auth.GET("/nodes/:id", r.h.Node.GetByID)
			auth.DELETE("/nodes/:id", r.h.Node.Delete)
			auth.GET("/nodes/:id/health", r.h.Node.Health)
		}
	}

	r.engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 404, "message": "route not found", "path": c.Request.URL.Path})
	})
	return r.engine
}

var _ = jwt.SigningMethodHS256
