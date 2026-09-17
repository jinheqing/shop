package config

import (
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database — 两个独立 PostgreSQL 连接（业务库 + 审计库）
type Database struct {
	Business *gorm.DB // 主业务库（tea_system）
	Audit    *gorm.DB // 审计库（tea_audit，独立实例）
}

// NewDatabase — 建立双连接
func NewDatabase(cfg *Config) (*Database, error) {
	db := &Database{}
	var err error

	// 业务库
	db.Business, err = connectPostgres(cfg.Database.BusinessDSN)
	if err != nil {
		return nil, err
	}
	log.Info().Str("dsn", maskDSN(cfg.Database.BusinessDSN)).Msg("📦 Business DB connected")

	// 审计库（独立实例）
	db.Audit, err = connectPostgres(cfg.Database.AuditDSN)
	if err != nil {
		return nil, err
	}
	log.Info().Str("dsn", maskDSN(cfg.Database.AuditDSN)).Msg("📦 Audit DB connected")

	return db, nil
}

func connectPostgres(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
}

// maskDSN — 把 password 打码（日志安全）
func maskDSN(dsn string) string {
	// 简单替换：postgres://user:PASSWORD@host/... → postgres://user:***@host/...
	// 实际生产建议用 net/url 解析
	out := []rune(dsn)
	inPwd := false
	for i, r := range out {
		if r == ':' && i > 0 && (out[i-1] == '/' || out[i-1] == '@' || i == 0) {
			// 不是 password 的那个冒号
		}
		if r == ':' && i >= 6 && string(out[i-6:i]) == "postgres" {
			inPwd = true
			continue
		}
		if inPwd {
			if r == '@' {
				inPwd = false
			} else {
				out[i] = '*'
			}
		}
	}
	return string(out)
}
