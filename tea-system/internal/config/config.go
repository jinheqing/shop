package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Config — 全局配置（viper 从 .env 加载 + DB site_contents 运行时覆盖）
type Config struct {
	Server             ServerConfig
	Database           DatabaseConfig
	Redis              RedisConfig
	RabbitMQ           RabbitMQConfig
	JWT                JWTConfig
	LiveKit            LiveKitConfig
	MediaMTX           MediaMTXConfig
	Mail               MailConfig
	Payment            PaymentConfig
	Admin              AdminSeedConfig
	Node               NodeConfig
	App                AppConfig
	TranslateServiceURL string // translate 引擎地址（从 TRANSLATE_SERVICE_URL 或 DB section "translate" 读取）
}

type ServerConfig struct {
	Port    string
	Mode    string
	Version string
}

type DatabaseConfig struct {
	BusinessDSN string
	AuditDSN    string
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
	Domain   string // 魔法链接 / 邮件发信域名
	SMTPHost string
	FromAddr string
	APIKey   string
}

type PaymentConfig struct {
	BaseURL   string
	APIKey    string
	APISecret string
	Gateway   string
}

type AdminSeedConfig struct {
	Email    string
	Name     string
	Password string
}

type NodeConfig struct {
	DefaultWGIPStart string // 新建 node 默认 WG IP（替换 node.go 里的 "10.10.0.99"）
	WGIPPrefix       string
}

type AppConfig struct {
	Domain       string
	ServiceName  string
	ContactEmail string
}

// Load — 从 .env / 环境变量加载
func Load() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// ---- 安全敏感项：无 SetDefault → viper 返回 "" → 生产 Validate() fail-fast ----
	// JWT_SECRET / LIVEKIT_API_KEY / LIVEKIT_API_SECRET / MAILGUN_API_KEY / DB 密码

	// 非敏感默认值
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_MODE", "debug")
	viper.SetDefault("SERVER_VERSION", "0.5.0")

	viper.SetDefault("MEDIAMTX_RTMP_PORT", 1935)
	viper.SetDefault("MEDIAMTX_WEBRTC_URL", "http://localhost:8889")

	viper.SetDefault("DB_BUSINESS_DSN", "postgres://tea_system@localhost:5432/tea_system?sslmode=disable")
	viper.SetDefault("DB_AUDIT_DSN", "postgres://tea_system@localhost:5432/tea_audit?sslmode=disable")

	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	viper.SetDefault("JWT_STAFF_EXPIRE_MIN", 60)
	viper.SetDefault("JWT_USER_EXPIRE_MIN", 120)
	viper.SetDefault("JWT_REFRESH_EXPIRE_DAY", 30)

	viper.SetDefault("LIVEKIT_URL", "http://localhost:7880")

	viper.SetDefault("SMTP_HOST", "localhost")
	viper.SetDefault("MAIL_FROM", "tea@localhost")
	viper.SetDefault("APP_DOMAIN", "http://localhost:8080")
	viper.SetDefault("APP_SERVICE_NAME", "UK Tea House")
	viper.SetDefault("APP_CONTACT_EMAIL", "tea@ukteahouse.co.uk")

	viper.SetDefault("ADMIN_SEED_EMAIL", "admin@ukteahouse.co.uk")
	viper.SetDefault("ADMIN_SEED_NAME", "Admin")
	viper.SetDefault("ADMIN_SEED_PASSWORD", "Admin!Tea2026")

	viper.SetDefault("PAYMENT_GATEWAY", "2checkout")

	viper.SetDefault("NODE_WG_IP_START", "10.10.0.10")
	viper.SetDefault("NODE_WG_IP_PREFIX", "24")

	viper.SetDefault("TRANSLATE_SERVICE_URL", "http://localhost:8090")

	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	cfg := &Config{
		Server: ServerConfig{
			Port:    viper.GetString("SERVER_PORT"),
			Mode:    viper.GetString("SERVER_MODE"),
			Version: viper.GetString("SERVER_VERSION"),
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
		RabbitMQ: RabbitMQConfig{URL: viper.GetString("RABBITMQ_URL")},
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
		Node: NodeConfig{
			DefaultWGIPStart: viper.GetString("NODE_WG_IP_START"),
			WGIPPrefix:       viper.GetString("NODE_WG_IP_PREFIX"),
		},
		App: AppConfig{
			Domain:       viper.GetString("APP_DOMAIN"),
			ServiceName:  viper.GetString("APP_SERVICE_NAME"),
			ContactEmail: viper.GetString("APP_CONTACT_EMAIL"),
		},
		TranslateServiceURL: viper.GetString("TRANSLATE_SERVICE_URL"),
	}

	if os.Getenv("DB_BUSINESS_DSN") != "" {
		cfg.Database.BusinessDSN = os.Getenv("DB_BUSINESS_DSN")
	}
	if os.Getenv("DB_AUDIT_DSN") != "" {
		cfg.Database.AuditDSN = os.Getenv("DB_AUDIT_DSN")
	}

	return cfg
}


// LoadRuntimeOverrides — DB 连上后调用一次
// 从 site_contents 里读 system_config 覆盖到 cfg
// 这样 admin System Settings 页面改的值会在下一次启动生效
func (c *Config) LoadRuntimeOverrides(db *gorm.DB) {
	type row struct {
		SectionKey string
		Content    []byte
	}
	var rows []row
	if err := db.Raw(
		`SELECT section_key, content FROM site_contents WHERE page_key = 'system_config'`,
	).Scan(&rows).Error; err != nil {
		log.Warn().Err(err).Msg("config: LoadRuntimeOverrides skipped (no site_contents)")
		return
	}

	for _, r := range rows {
		var vals map[string]any
		if len(r.Content) > 0 {
			_ = json.Unmarshal(r.Content, &vals)
		}
		if vals == nil {
			vals = map[string]any{}
		}
		applyOverrides(c, r.SectionKey, vals)
	}

	if len(rows) > 0 {
		log.Info().Int("entries", len(rows)).Msg("config: applied runtime overrides from site_contents")
	}
}

func applyOverrides(c *Config, key string, v map[string]any) {
	str := func(k string) string {
		if s, ok := v[k].(string); ok {
			return s
		}
		return ""
	}
	switch key {
	case "app":
		if d := str("domain"); d != "" {
			c.App.Domain = d
			c.Mail.Domain = d
		}
		if s := str("service_name"); s != "" {
			c.App.ServiceName = s
		}
		if e := str("contact_email"); e != "" {
			c.App.ContactEmail = e
		}
	case "mail":
		if h := str("smtp_host"); h != "" {
			c.Mail.SMTPHost = h
		}
		if f := str("from_addr"); f != "" {
			c.Mail.FromAddr = f
		}
		if a := str("api_key"); a != "" {
			c.Mail.APIKey = a
		}
	case "translate":
		if u := str("service_url"); u != "" {
			c.TranslateServiceURL = u
		}
	case "payment":
		if b := str("base_url"); b != "" {
			c.Payment.BaseURL = b
		}
	case "nodes":
		if s := str("wg_ip_start"); s != "" {
			c.Node.DefaultWGIPStart = s
		}
		if p := str("wg_ip_prefix"); p != "" {
			c.Node.WGIPPrefix = p
		}
	case "server":
		if v := str("version"); v != "" {
			c.Server.Version = v
		}
	case "seed_admin":
		if e := str("email"); e != "" {
			c.Admin.Email = e
		}
		if n := str("name"); n != "" {
			c.Admin.Name = n
		}
		if p := str("password"); p != "" {
			c.Admin.Password = p
		}
	}
}

// Validate — 生产模式下对敏感配置做 fail-fast 校验
func (c *Config) Validate() error {
	prod := c.Server.Mode == "release" || c.Server.Mode == "production"
	if prod {
		if c.JWT.Secret == "" {
			return fmt.Errorf("JWT_SECRET must be set in production")
		}
		if len(c.JWT.Secret) < 32 {
			return fmt.Errorf("JWT_SECRET must be at least 32 characters (got %d)", len(c.JWT.Secret))
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
