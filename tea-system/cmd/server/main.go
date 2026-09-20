// Package main — UK Tea House Go API Server 入口
//
// 启动流程:
//  1. 加载配置 (.env + viper)
//  2. 连接 PostgreSQL (业务库 + 审计库)
//  3. 连接 Redis
//  4. AutoMigrate 建表
//  5. 创建默认管理员（如果不存在）
//  6. 构造所有 Repository → Service → Handler
//  7. 启动 IM WebSocket Hub
//  8. 启动 HTTP Server
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"tea-system/internal/api"
	"tea-system/internal/api/handlers"
	"tea-system/internal/config"
	"tea-system/internal/im"
	"tea-system/internal/models"
	"tea-system/internal/repository"
	"tea-system/internal/service"
)

func main() {
	// ── 0. 日志初始化 ──
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("🚀 UK Tea House API starting...")

	// ── 1. 加载配置 ──
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal().Err(err).Msg("❌ config validation failed — refuse to start")
	}

	log.Info().
		Str("mode", cfg.Server.Mode).
		Str("version", cfg.Server.Version).
		Str("port", cfg.Server.Port).
		Msg("📋 Config loaded")

	// ── 2. 连接 PostgreSQL (业务库 + 审计库) ──
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Failed to connect PostgreSQL")
	}
	log.Info().Msg("📦 Both PostgreSQL databases connected")

	// ── 3. 连接 Redis ──
	rdb := newRedis(cfg)

	// ── 4. AutoMigrate 建表 ──
	log.Info().Msg("🔨 Running AutoMigrate...")
	// PostgreSQL: 同事务内 FK 约束检查绕过
	db.Business.Exec("SET session_replication_role = 'replica'")
	db.Audit.Exec("SET session_replication_role = 'replica'")
	if err := db.Business.AutoMigrate(models.AllModels...); err != nil {
		log.Fatal().Err(err).Msg("❌ Business DB AutoMigrate failed")
	}
	if err := db.Audit.AutoMigrate(models.AuditModels...); err != nil {
		log.Fatal().Err(err).Msg("❌ Audit DB AutoMigrate failed")
	}
	db.Business.Exec("SET session_replication_role = 'origin'")
	db.Audit.Exec("SET session_replication_role = 'origin'")
	log.Info().Msg("✅ AutoMigrate complete")

	// ── 5. 从 DB 加载运行时覆盖 + 创建默认管理员 ──
	cfg.LoadRuntimeOverrides(db.Business)

	// ── 6. 构造 Repository 层 ──
	staffRepo := repository.NewStaffRepo(db.Business)
	userRepo := repository.NewUserRepo(db.Business)
	orderRepo := repository.NewOrderRepo(db.Business)
	customProductRepo := repository.NewCustomProductRepo(db.Business)
	invoiceRepo := repository.NewInvoiceRepo(db.Business)
	liveRoomRepo := repository.NewLiveRoomRepo(db.Business)
	conversationRepo := repository.NewConversationRepo(db.Business)
	messageRepo := repository.NewMessageRepo(db.Business)
	nodeRepo := repository.NewNodeRepo(db.Business)

	// 确保默认管理员存在（用 repo 查询 + 创建）
	passwordSvc := service.NewPasswordService()
	ensureAdminUser(staffRepo, passwordSvc, cfg)

	// ── 7. 构造 Service 层 ──
	jwtSvc := service.NewJWTService(
		cfg.JWT.Secret, cfg.JWT.StaffExpireMin, cfg.JWT.UserExpireMin, cfg.JWT.RefreshExpireDay,
	)
	mfaSvc := service.NewMFAService(staffRepo)
	mailSvc := service.NewMailService(cfg.Mail.APIKey, cfg.App.Domain, cfg.Mail.FromAddr)
	magicLinkSvc := service.NewMagicLinkService(rdb, mailSvc, cfg.App.Domain)

	lkSvc := service.NewLiveKitService(cfg.LiveKit.URL, cfg.LiveKit.APIKey, cfg.LiveKit.APISecret)
	translateSvc := service.NewTranslateService(cfg.TranslateServiceURL)
	paymentSvc := service.NewPaymentService(cfg.Payment.BaseURL, cfg.Payment.APIKey, cfg.Payment.APISecret)

	orderSM := service.NewOrderStateMachine()
	customProductSvc := service.NewCustomProductService()
	liveRoomSvc := service.NewLiveRoomService(liveRoomRepo, cfg, lkSvc)
	slowPresetSvc := service.NewSlowPresetService(liveRoomRepo, cfg)

	storageDir := filepath.Join(".", "storage", "invoices")
	invoiceSvc := service.NewInvoiceService(orderRepo, invoiceRepo, userRepo, storageDir, cfg.RabbitMQ.URL)

	// WireGuard 边缘节点管理（UK 主节点禁用，国内边缘节点启用）
	nodeDeployer := service.NewNodeDeployer(nodeRepo, nil)
	healthChecker := service.NewHealthChecker(nodeRepo, 30*time.Second)

	// IM WebSocket Hub
	imHub := im.NewHub()
	go imHub.Start()

	// ── 8. 构造 Handlers ──
	h := &api.Handlers{
		Health:    handlers.NewHealthHandler(cfg.Server.Version),
		StaffAuth: handlers.NewStaffAuthHandler(staffRepo, userRepo, passwordSvc, jwtSvc, mfaSvc, rdb, db.Audit),
		UserAuth:  handlers.NewUserAuthHandler(userRepo, passwordSvc, jwtSvc, magicLinkSvc, db.Business),

		CustomProduct: handlers.NewCustomProductHandler(customProductRepo, customProductSvc),
		LiveKit:       handlers.NewLiveKitTokenHandler(lkSvc, cfg),
		Conversation:  handlers.NewConversationHandler(conversationRepo),
		Message:       handlers.NewMessageHandler(messageRepo),
		IMWS:          handlers.NewIMWSHandler(imHub, messageRepo, translateSvc, cfg.JWT.Secret),

		Order:       handlers.NewOrderHandler(orderRepo, customProductRepo, userRepo, orderSM, mailSvc, db.Audit),
		Payment:     handlers.NewPaymentHandler(orderRepo, paymentSvc, orderSM),
		Invoice:     handlers.NewInvoiceHandler(invoiceSvc, orderRepo, invoiceRepo),
		Declaration: handlers.NewDeclarationHandler(invoiceRepo),
		Ledger:      handlers.NewLedgerHandler(invoiceRepo),
		SgsReport:   handlers.NewSgsReportHandler(invoiceRepo),

		SlowPreset: handlers.NewSlowPresetHandler(slowPresetSvc),
		LiveRoom:   handlers.NewLiveRoomHandler(liveRoomSvc),

		Node: handlers.NewNodeHandler(nodeDeployer, healthChecker, nodeRepo, cfg.Node.DefaultWGIPStart, cfg.Node.WGIPPrefix),

		Translate:     handlers.NewTranslateHandler(cfg.TranslateServiceURL),
		QRCode:        handlers.NewQRCodeHandler(db.Business),
		DSAR:          handlers.NewDSARHandler(db.Business),
		SiteContent:   handlers.NewSiteContentHandler(db.Business),
		CookieConsent: handlers.NewCookieConsentHandler(db.Business),
		SystemConfig:  handlers.NewSystemConfigHandler(db.Business),
		Upload:        handlers.NewUploadHandler("./storage/uploads", "/uploads"),

		UserGroup:  handlers.NewUserGroupHandler(db.Business),
		ShortLink:  handlers.NewShortLinkHandler(db.Business),
		Recording:  handlers.NewRecordingHandler(db.Business),
		Video:      handlers.NewVideoHandler(db.Business),
	}

	// ── 9. 构造 Router + 启动 HTTP ──
	router := api.NewRouter(cfg, db.Business, db.Audit, rdb, h)
	engine := router.Setup()

	addr := ":" + cfg.Server.Port
	srv := &http.Server{
		Addr:              addr,
		Handler:           engine,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 优雅关闭
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Info().Msg("🛑 Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("shutdown error")
		}
		rdb.Close()
		log.Info().Msg("👋 Server stopped")
	}()

	log.Info().Str("addr", addr).Msg("🌶️  HTTP server listening")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("❌ Server fatal error")
	}
}

// ─── 辅助函数 ───

func newRedis(cfg *config.Config) *redis.Client {
	var opts *redis.Options
	if cfg.Redis.Password != "" {
		opts = &redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password, DB: cfg.Redis.DB}
	} else {
		opts = &redis.Options{Addr: cfg.Redis.Addr, DB: cfg.Redis.DB}
	}
	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("❌ Redis ping failed")
	}
	log.Info().Msg("🔴 Redis connected")
	return rdb
}

func ensureAdminUser(repo *repository.StaffRepo, pwdSvc *service.PasswordService, cfg *config.Config) {
	ctx := context.Background()

	existing, err := repo.GetByEmail(ctx, cfg.Admin.Email)
	if err == nil && existing != nil {
		log.Info().Str("email", cfg.Admin.Email).Msg("👤 Admin user already exists, skip seed")
		return
	}

	log.Info().Str("email", cfg.Admin.Email).Msg("🌱 Creating default admin user...")

	hashed, err := pwdSvc.Hash(cfg.Admin.Password)
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Failed to hash admin password")
	}

	staff := &models.Staff{
		Name:         cfg.Admin.Name,
		Email:        cfg.Admin.Email,
		PasswordHash: hashed,
		Role:         "admin",
		IsActive:     true,
		WorkTimezone: "Europe/London",
	}

	if err := repo.Create(ctx, staff); err != nil {
		log.Fatal().Err(err).Msg("❌ Failed to create admin user")
	}

	log.Info().
		Str("email", cfg.Admin.Email).
		Str("name", cfg.Admin.Name).
		Msg("✅ Admin user created — login with this email + the password you set in .env ADMIN_SEED_PASSWORD")

	// 占位格式: fmt 用到了但 import 里可能不需要
	_ = fmt.Sprintf
}
