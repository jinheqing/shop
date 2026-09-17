package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONMap is a helper type for jsonb columns
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONMap value")
	}
	return json.Unmarshal(bytes, j)
}

// User — 客户表
type User struct {
	ID                    uint64    `gorm:"primaryKey;column:id" json:"id"`
	Name                  string    `gorm:"column:name;not null;size:100" json:"name"`
	Email                 string    `gorm:"column:email;uniqueIndex;not null;size:200" json:"email"`
	Phone                 string    `gorm:"column:phone;size:30" json:"phone"`
	PasswordHash          string    `gorm:"column:password_hash;size:255" json:"-"`
	BillingAddress        JSONMap   `gorm:"column:billing_address;type:jsonb" json:"billing_address,omitempty"`
	DeliveryAddress       JSONMap   `gorm:"column:delivery_address;type:jsonb" json:"delivery_address,omitempty"`
	PreferredLanguage     string    `gorm:"column:preferred_language;size:10;default:'en'" json:"preferred_language"`
	PreferredTimezone     string    `gorm:"column:preferred_timezone;size:50;default:'Europe/London'" json:"preferred_timezone"`
	PreferredAdvisorID    *uint64   `gorm:"column:preferred_advisor_id" json:"preferred_advisor_id,omitempty"`
	ConsentMarketing      bool      `gorm:"column:consent_marketing;default:false" json:"consent_marketing"`
	ConsentAnalytics      bool      `gorm:"column:consent_analytics;default:false" json:"consent_analytics"`
	DataDeleteRequestedAt *time.Time `gorm:"column:data_delete_requested_at" json:"data_delete_requested_at,omitempty"`
	DataDeleteCompletedAt *time.Time `gorm:"column:data_delete_completed_at" json:"data_delete_completed_at,omitempty"`
	DsarRequestCount      int       `gorm:"column:dsar_request_count;default:0" json:"dsar_request_count"`
	LastLoginAt           *time.Time `gorm:"column:last_login_at" json:"last_login_at,omitempty"`
	CreatedAt             time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt             *time.Time `gorm:"column:deleted_at" json:"-"`

	// Associations
	PreferredAdvisor *Staff `gorm:"foreignKey:PreferredAdvisorID" json:"preferred_advisor,omitempty"`
}

func (User) TableName() string {
	return "users"
}
