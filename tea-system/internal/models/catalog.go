package models

import (
	"time"
)

// SgsReport — SGS 检测报告表
type SgsReport struct {
	ID        uint64    `gorm:"primaryKey;column:id" json:"id"`
	ReportNo  string    `gorm:"column:report_no;uniqueIndex;not null;size:100" json:"report_no"`
	BatchNo   string    `gorm:"column:batch_no;not null;size:100" json:"batch_no"`
	TeaType   string    `gorm:"column:tea_type;not null;size:30" json:"tea_type"`
	TestDate  time.Time `gorm:"column:test_date;not null;type:date" json:"test_date"`
	IssueDate time.Time `gorm:"column:issue_date;not null;type:date" json:"issue_date"`
	PdfURL    string    `gorm:"column:pdf_url;not null;size:500" json:"pdf_url"`
	TestItems JSONMap   `gorm:"column:test_items;type:jsonb" json:"test_items,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (SgsReport) TableName() string {
	return "sgs_reports"
}

// ⚠️ 安全约束：SSH 密码不存入此表。
// 部署期间仅在 Go 进程内存变量中存在，用完立即置空并 runtime.GC()。

// TranslationSession — 翻译会话表
type TranslationSession struct {
	ID          uint64    `gorm:"primaryKey;column:id" json:"id"`
	SourceType  string    `gorm:"column:source_type;not null;size:20" json:"source_type"` // im_text / live_subtitle
	SourceID    *uint64   `gorm:"column:source_id" json:"source_id,omitempty"`
	Direction   string    `gorm:"column:direction;not null;size:10" json:"direction"` // zh_en / en_zh / both
	Status      string    `gorm:"column:status;not null;size:20" json:"status"`        // pending/translating/completed/failed
	CreatedAt   time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (TranslationSession) TableName() string {
	return "translation_sessions"
}

// AuditLog — 审计日志表（**独立 PostgreSQL 实例**，无软删除）
type AuditLog struct {
	ID         uint64    `gorm:"primaryKey;column:id" json:"id"`
	StaffID    *uint64   `gorm:"column:staff_id" json:"staff_id,omitempty"`
	Action     string    `gorm:"column:action;not null;size:100" json:"action"`
	TargetType string    `gorm:"column:target_type;size:50" json:"target_type,omitempty"`
	TargetID   *uint64   `gorm:"column:target_id" json:"target_id,omitempty"`
	Detail     JSONMap   `gorm:"column:detail;type:jsonb" json:"detail,omitempty"`
	IPAddress  string    `gorm:"column:ip_address;size:50" json:"ip_address,omitempty"`
	UserAgent  string    `gorm:"column:user_agent;type:text" json:"user_agent,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// SiteContent — 公开网站 CMS 表
type SiteContent struct {
	ID               uint64    `gorm:"primaryKey;column:id" json:"id"`
	PageKey          string    `gorm:"column:page_key;uniqueIndex:idx_site_page_sec;not null;size:50" json:"page_key"`
	SectionKey       string    `gorm:"column:section_key;uniqueIndex:idx_site_page_sec;not null;size:50" json:"section_key"`
	Content          JSONMap   `gorm:"column:content;not null;type:jsonb" json:"content"`
	IsPublished      bool      `gorm:"column:is_published;default:true" json:"is_published"`
	UpdatedByStaffID *uint64   `gorm:"column:updated_by_staff_id" json:"updated_by_staff_id,omitempty"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (SiteContent) TableName() string {
	return "site_contents"
}
