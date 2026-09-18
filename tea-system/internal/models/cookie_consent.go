package models

import "time"

type CookieConsentLog struct {
	ID               uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	VisitorID        string    `gorm:"column:visitor_id;size:100;index" json:"visitor_id,omitempty"`
	UserID           *uint64   `gorm:"column:user_id" json:"user_id,omitempty"`
	ConsentEssential bool      `gorm:"column:consent_essential;not null;default:true" json:"consent_essential"`
	ConsentAnalytics bool      `gorm:"column:consent_analytics;not null;default:false" json:"consent_analytics"`
	ConsentMarketing bool      `gorm:"column:consent_marketing;not null;default:false" json:"consent_marketing"`
	UserAgent        string    `gorm:"column:user_agent;size:500" json:"-"`
	IPAddress        string    `gorm:"column:ip_address;size:64" json:"-"`
	CreatedAt        time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (CookieConsentLog) TableName() string { return "cookie_consents" }
