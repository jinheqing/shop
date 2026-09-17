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

// Node — 分发节点表（SSH 密码**不持久化**，部署期间仅内存存在）
type Node struct {
	ID                 uint64     `gorm:"primaryKey;column:id" json:"id"`
	NodeName           string     `gorm:"column:node_name;not null;size:100" json:"node_name"`
	NodeType           string     `gorm:"column:node_type;not null;size:30" json:"node_type"` // livekit_edge / video_cdn
	PublicIP           string     `gorm:"column:public_ip;not null;size:50" json:"public_ip"`
	WireguardIP        string     `gorm:"column:wireguard_ip;not null;size:50" json:"wireguard_ip"` // 10.10.0.xx
	WireguardPublicKey string     `gorm:"column:wireguard_public_key;not null;size:200" json:"wireguard_public_key"`
	Status             string     `gorm:"column:status;not null;size:20" json:"status"` // deploying/online/offline/error
	LastHealthCheck    *time.Time `gorm:"column:last_health_check" json:"last_health_check,omitempty"`
	LastCPUPercent     *int       `gorm:"column:last_cpu_percent" json:"last_cpu_percent,omitempty"`
	LastMemoryPercent  *int       `gorm:"column:last_memory_percent" json:"last_memory_percent,omitempty"`
	LastBandwidthMbps  *int       `gorm:"column:last_bandwidth_mbps" json:"last_bandwidth_mbps,omitempty"`
	LastViewersCount   *int       `gorm:"column:last_viewers_count" json:"last_viewers_count,omitempty"`
	DeployedByStaffID  *uint64    `gorm:"column:deployed_by_staff_id" json:"deployed_by_staff_id,omitempty"`
	DeployedAt         *time.Time `gorm:"column:deployed_at" json:"deployed_at,omitempty"`
	DeletedAt          *time.Time `gorm:"column:deleted_at" json:"-"`
	CreatedAt          time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`

	DeployedByStaff *Staff `gorm:"foreignKey:DeployedByStaffID" json:"deployed_by_staff,omitempty"`
}

func (Node) TableName() string {
	return "nodes"
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
	ID            uint64    `gorm:"primaryKey;column:id" json:"id"`
	PageKey       string    `gorm:"column:page_key;uniqueIndex;not null;size:50" json:"page_key"`
	SectionKey    string    `gorm:"column:section_key;not null;size:50" json:"section_key"`
	Content       JSONMap   `gorm:"column:content;not null;type:jsonb" json:"content"`
	UpdatedByStaffID *uint64 `gorm:"column:updated_by_staff_id" json:"updated_by_staff_id,omitempty"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (SiteContent) TableName() string {
	return "site_contents"
}
