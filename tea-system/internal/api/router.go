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

		// Step 16+: 补齐设计文档缺失 handler
		Translate      *handlers.TranslateHandler
		QRCode         *handlers.QRCodeHandler
		DSAR           *handlers.DSARHandler
		SiteContent    *handlers.SiteContentHandler
		CookieConsent  *handlers.CookieConsentHandler
	SystemConfig    *handlers.SystemConfigHandler
	Upload          *handlers.UploadHandler
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
		// ---------- 认证（限流：防暴力破解/邮件轰炸） ----------
		v1.POST("/staff/login", middleware.RateLimit(2, 10), r.h.StaffAuth.Login)           // 2 req/s, burst 10
		v1.POST("/staff/mfa/verify", middleware.RateLimit(5, 20), r.h.StaffAuth.MFAVerify)
		v1.POST("/staff/refresh", r.h.StaffAuth.Refresh)
		v1.POST("/user/magic-link/request", middleware.RateLimit(1, 5), r.h.UserAuth.MagicLinkRequest) // 1 req/s, burst 5
		v1.POST("/user/magic-link/verify", middleware.RateLimit(5, 20), r.h.UserAuth.MagicLinkVerify)
		v1.POST("/user/login", middleware.RateLimit(2, 10), r.h.UserAuth.UserLogin)

		// ---------- 公开 API ----------
		v1.GET("/custom-products/by-token/:token", r.h.CustomProduct.GetByToken)
		v1.GET("/public/slow-presets", r.h.SlowPreset.PublicList)
		v1.POST("/live-rooms/customer-request", r.h.LiveRoom.CustomerRequest)
		v1.GET("/public/sgs-reports", r.h.SgsReport.PublicList)
                v1.GET("/public/live-rooms", r.h.LiveRoom.PublicLiveRooms)
                v1.GET("/custom-products/published", r.h.CustomProduct.Published)
                // Public bespoke submission — rate-limited, no auth required
                v1.POST("/public/custom-products", middleware.RateLimit(2, 10), r.h.CustomProduct.CreatePublic)

		// ---------- 支付 Webhook ----------
		v1.POST("/webhooks/2checkout", r.h.Payment.Handle2CheckoutWebhook)
		v1.POST("/webhooks/paypal", r.h.Payment.HandlePayPalWebhook)

		// 需要 JWT 的
		auth := v1.Group("")
		auth.Use(middleware.JWTAuth(r.cfg.JWT.Secret))
		{
			auth.POST("/staff/logout", r.h.StaffAuth.Logout)
			auth.GET("/staff/audit-logs", r.h.StaffAuth.AuditLogs)

			// ===== 新增缺失后端路由 =====
			// Staff CRUD (admin/supervisor)
			auth.GET("/staff", r.h.StaffAuth.StaffList)
			auth.POST("/staff", r.h.StaffAuth.StaffCreate)
			auth.POST("/staff/:id/toggle", r.h.StaffAuth.StaffToggle)
			auth.DELETE("/staff/:id", r.h.StaffAuth.StaffDelete)
			// Users (customers) list
			auth.GET("/users", middleware.RequireRole("admin", "supervisor"), r.h.StaffAuth.UserList)
			// System Config
			auth.GET("/system/config", middleware.RequireRole("admin", "supervisor"), r.h.SystemConfig.List)
			auth.GET("/system/config/:key", middleware.RequireRole("admin", "supervisor"), r.h.SystemConfig.Get)
			auth.PUT("/system/config/:key", middleware.RequireRole("admin", "supervisor"), r.h.SystemConfig.Put)
			// Payment Transactions
			auth.GET("/payment/transactions", middleware.RequireRole("admin", "supervisor"), r.h.Payment.ListTransactions)

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

			// === 设计文档缺失 API — 补齐 ===
			// 订单时间线
			auth.GET("/orders/:id/timeline", r.h.Order.Timeline)
			// PayPal 单独初始化
			auth.POST("/orders/:id/payment/paypal", r.h.Payment.Init)

			// 报关作废（标记 void 不可删除）
			auth.POST("/declarations/:id/void", r.h.Declaration.Void)
			// exchange-records 别名（设计文档命名）
			auth.GET("/exchange-records", r.h.Ledger.List)
			auth.POST("/exchange-records", r.h.Ledger.Create)

			// Live Room: schedule calendar + 客户申请审核 + system-create
			auth.GET("/live-rooms/schedule/calendar", r.h.LiveRoom.Calendar)
			auth.POST("/live-rooms/customer-requests/:id/approve", r.h.LiveRoom.ApproveRequest)
			auth.POST("/live-rooms/system-create", r.h.LiveRoom.SystemCreate)

			// QR Code
			auth.POST("/qrcodes/generate", r.h.QRCode.Generate)

			// 翻译引擎代理
			auth.POST("/translate/text", r.h.Translate.TranslateText)
			auth.GET("/translate/status", r.h.Translate.Status)
			auth.GET("/translate/asr", r.h.Translate.ASR)

			// CMS + GDPR + DSAR
			auth.GET("/site-contents", r.h.SiteContent.List)
			auth.PUT("/site-contents/:id", r.h.SiteContent.Update)
			auth.POST("/dsar/requests", r.h.DSAR.CreateRequest)
			auth.GET("/dsar/requests", r.h.DSAR.List)
			auth.POST("/dsar/requests/:id/export", r.h.DSAR.Export)
			auth.POST("/dsar/requests/:id/delete", r.h.DSAR.Delete)
			auth.GET("/audit-logs", r.h.StaffAuth.AuditLogs) // 别名 → 与 /staff/audit-logs 同实现

			// PayPal + 2Checkout callback 别名（设计文档路径）
			v1.POST("/payment/callback/2checkout", r.h.Payment.Handle2CheckoutWebhook)
			v1.POST("/payment/callback/paypal", r.h.Payment.HandlePayPalWebhook)
		}
	}

		// 公开路由：QR trace + cookie consent + site contents
		v1.GET("/public/qrcodes/:token", r.h.QRCode.GetTrace)
		v1.POST("/cookie-consent", r.h.CookieConsent.Submit)

                // 文件上传 + 静态文件服务
                r.engine.Static("/uploads", "./storage/uploads")
                v1.POST("/upload", r.h.Upload.Upload)


	r.engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 404, "message": "route not found", "path": c.Request.URL.Path})
	})
	return r.engine
}

var _ = jwt.SigningMethodHS256
