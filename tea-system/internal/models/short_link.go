package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// ShortLink — 短链表，code → target_url 映射
// GET /s/{code} → 302 重定向
type ShortLink struct {
	ID          uint64    `gorm:"primaryKey;column:id" json:"id"`
	Code        string    `gorm:"column:code;uniqueIndex;size:8;not null" json:"code"` // auto-gen 6 chars
	TargetURL   string    `gorm:"column:target_url;size:500;not null" json:"target_url"`
	OwnerUserID *uint64   `gorm:"column:owner_user_id" json:"owner_user_id,omitempty"`
	ClickCount  int       `gorm:"column:click_count;default:0" json:"click_count"`
	ExpiresAt   *time.Time `gorm:"column:expires_at" json:"expires_at,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (ShortLink) TableName() string {
	return "short_links"
}

// GenerateCode 自动生成 6 字符短链 code（不含易混字符 0/O/1/I/l）
func GenerateShortCode() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
	b := make([]byte, 6)
	rand.Read(b)
	for i, v := range b {
		b[i] = chars[int(v)%len(chars)]
	}
	_ = hex.EncodeToString(b) // keep crypto import satisfied
	return fmt.Sprintf("%s", string(b))
}
