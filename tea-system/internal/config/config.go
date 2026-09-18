package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config — 全局配置（viper 从 .env 加载）
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
	JWT      JWTConfig
	LiveKit  LiveKitConfig
	MediaMTX MediaMTXConfig
	Mail     MailConfig
	Payment  PaymentConfig
	Admin    AdminSeedConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	BusinessDSN string // 主业务库
	AuditDSN    string // 审计库（独立 PostgreSQL 实例）
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type RabbitMQConfig struct {
	URL string
}

type JWTConfig struct {
	Secret           string
	StaffExpireMin   int
	UserExpireMin    int
	RefreshExpireDay int
}

type LiveKitConfig struct {
	URL       string
	APIKey    string
	APISecret string
}

type MediaMTXConfig struct {
	RTMPPort  int
	WebRTCURL string
}

type MailConfig struct {
	Domain   string // 发信域名，如 https://ukteahouse.co.uk
	SMTPHost string
	FromAddr string
	APIKey   string // 生产用 Mailgun/Resend 等 API key（本地测试可空）
}

type PaymentConfig struct {
	BaseURL   string // 支付网关回调的基础 URL（mock 模式也用得到）
	APIKey    string
	APISecret string
	Gateway   string // "2checkout" | "paypal"
}

// AdminSeedConfig — 首次启动时自动创建 admin 账号，生产务必通过环境变量覆盖默认值
type AdminSeedConfig struct {
	Email    string
	Name     string
	Password string
}

// Load — 从 .env 或环境变量加载配置
func Load() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// ---- 安全敏感项：默认空串，生产必须显式设置 ----
	// JWT_SECRET / LIVEKIT_API_KEY / LIVEKIT_API_SECRET / MAILGUN_API_KEY / 数据库密码
	// 这些都 **没有** SetDefault，意味着 viper 读到空时会返回 ""
	// 启动时会在 Validate() 里 fail-fast

	// 非敏感默认值（开发环境友好）
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_MODE", "debug")
	viper.SetDefault("MEDIAMTX_RTMP_PORT", 1935)

	viper.SetDefault("DB_BUSINESS_DSN", "postgres://tea_system@localhost:5432/tea_system?sslmode=disable")
	viper.SetDefault("DB_AUDIT_DSN", "postgres://tea_system@localhost:5432/tea_audit?sslmode=disable")

	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	viper.SetDefault("JWT_STAFF_EXPIRE_MIN", 60)
	viper.SetDefault("JWT_USER_EXPIRE_MIN", 120)
	viper.SetDefault("JWT_REFRESH_EXPIRE_DAY", 30)

	viper.SetDefault("LIVEKIT_URL", "http://localhost:7880")
	viper.SetDefault("MEDIAMTX_WEBRTC_URL", "http://localhost:8889")

	viper.SetDefault("SMTP_HOST", "localhost")
	viper.SetDefault("MAIL_FROM", "tea@localhost")
	viper.SetDefault("APP_DOMAIN", "http://localhost:8080") // 魔法链接域名 + Mail 发信域名

	// Admin seed — 只在 staff 表为空时触发
	viper.SetDefault("ADMIN_SEED_EMAIL", "admin@ukteahouse.co.uk")
	viper.SetDefault("ADMIN_SEED_NAME", "Admin")
	viper.SetDefault("ADMIN_SEED_PASSWORD", "Admin!Tea2026") // ⚠️ 生产必须通过环境变量覆盖

	viper.SetDefault("PAYMENT_GATEWAY", "2checkout")
	viper.SetDefault("PAYMENT_BASE_URL", "")

	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	cfg := &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
			Mode: viper.GetString("SERVER_MODE"),
		},
		Database: DatabaseConfig{
			BusinessDSN: viper.GetString("DB_BUSINESS_DSN"),
			AuditDSN:    viper.GetString("DB_AUDIT_DSN"),
		},
		Redis: RedisConfig{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: viper.GetString("RABBITMQ_URL"),
		},
		JWT: JWTConfig{
			Secret:           viper.GetString("JWT_SECRET"),
			StaffExpireMin:   viper.GetInt("JWT_STAFF_EXPIRE_MIN"),
			UserExpireMin:    viper.GetInt("JWT_USER_EXPIRE_MIN"),
			RefreshExpireDay: viper.GetInt("JWT_REFRESH_EXPIRE_DAY"),
		},
		LiveKit: LiveKitConfig{
			URL:       viper.GetString("LIVEKIT_URL"),
			APIKey:    viper.GetString("LIVEKIT_API_KEY"),
			APISecret: viper.GetString("LIVEKIT_API_SECRET"),
		},
		MediaMTX: MediaMTXConfig{
			RTMPPort:  viper.GetInt("MEDIAMTX_RTMP_PORT"),
			WebRTCURL: viper.GetString("MEDIAMTX_WEBRTC_URL"),
		},
		Mail: MailConfig{
			Domain:   viper.GetString("APP_DOMAIN"),
			SMTPHost: viper.GetString("SMTP_HOST"),
			FromAddr: viper.GetString("MAIL_FROM"),
			APIKey:   viper.GetString("MAILGUN_API_KEY"),
		},
		Payment: PaymentConfig{
			BaseURL:   viper.GetString("PAYMENT_BASE_URL"),
			APIKey:    viper.GetString("PAYMENT_API_KEY"),
			APISecret: viper.GetString("PAYMENT_API_SECRET"),
			Gateway:   viper.GetString("PAYMENT_GATEWAY"),
		},
		Admin: AdminSeedConfig{
			Email:    viper.GetString("ADMIN_SEED_EMAIL"),
			Name:     viper.GetString("ADMIN_SEED_NAME"),
			Password: viper.GetString("ADMIN_SEED_PASSWORD"),
		},
	}

	if os.Getenv("DB_BUSINESS_DSN") != "" {
		cfg.Database.BusinessDSN = os.Getenv("DB_BUSINESS_DSN")
	}
	if os.Getenv("DB_AUDIT_DSN") != "" {
		cfg.Database.AuditDSN = os.Getenv("DB_AUDIT_DSN")
	}

	return cfg
}

// Validate — 生产模式下对敏感配置做 fail-fast 校验
func (c *Config) Validate() error {
	prod := c.Server.Mode == "release" || c.Server.Mode == "production"
	if prod {
		if c.JWT.Secret == "" {
			return fmt.Errorf("JWT_SECRET must be set in production")
		}
		if len(c.JWT.Secret) < 32 {
			return fmt.Errorf("JWT_SECRET must be at least 32 characters in production (got %d)", len(c.JWT.Secret))
		}
		if c.LiveKit.APIKey == "" || c.LiveKit.APISecret == "" {
			return fmt.Errorf("LIVEKIT_API_KEY / LIVEKIT_API_SECRET must be set in production")
		}
		if c.Admin.Password == "Admin!Tea2026" {
			return fmt.Errorf("ADMIN_SEED_PASSWORD must be overridden from default in production")
		}
		if c.Mail.Domain == "http://localhost:8080" {
			return fmt.Errorf("APP_DOMAIN must be set to your production domain")
		}
	}
	return nil
}
