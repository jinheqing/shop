package config

import (
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
	Secret            string
	StaffExpireMin    int
	UserExpireMin     int
	RefreshExpireDay  int
}

type LiveKitConfig struct {
	URL    string
	APIKey string
	APISecret string
}

type MediaMTXConfig struct {
	RTMPPort  int
	WebRTCURL string
}

// Load — 从 .env 或环境变量加载配置
func Load() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// 设置默认值
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_MODE", "debug")

	viper.SetDefault("DB_BUSINESS_DSN", "postgres://postgres:test@localhost:5432/tea_system?sslmode=disable")
	viper.SetDefault("DB_AUDIT_DSN", "postgres://postgres:test@localhost:5432/tea_audit?sslmode=disable")

	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	viper.SetDefault("JWT_SECRET", "change-me-in-production-please-use-a-long-random-string")
	viper.SetDefault("JWT_STAFF_EXPIRE_MIN", 60)
	viper.SetDefault("JWT_USER_EXPIRE_MIN", 120)
	viper.SetDefault("JWT_REFRESH_EXPIRE_DAY", 30)

	viper.SetDefault("LIVEKIT_URL", "http://localhost:7880")
	viper.SetDefault("LIVEKIT_API_KEY", "dev-key")
	viper.SetDefault("LIVEKIT_API_SECRET", "dev-secret")

	viper.SetDefault("MEDIAMTX_RTMP_PORT", 1935)
	viper.SetDefault("MEDIAMTX_WEBRTC_URL", "http://localhost:8889")

	// 允许环境变量覆盖
	viper.AutomaticEnv()

	// 尝试读 .env，不存在也不报错（CI/生产用环境变量）
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
	}

	// 允许 os.Environ 兜底（某些部署环境 viper.AutomaticEnv 需要配合 SetEnvKeyReplacer）
	if os.Getenv("DB_BUSINESS_DSN") != "" {
		cfg.Database.BusinessDSN = os.Getenv("DB_BUSINESS_DSN")
	}
	if os.Getenv("DB_AUDIT_DSN") != "" {
		cfg.Database.AuditDSN = os.Getenv("DB_AUDIT_DSN")
	}

	return cfg
}
