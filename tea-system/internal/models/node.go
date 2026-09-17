package models

import "time"

// NodeType 常量
const (
	NodeTypeLiveKitEdge = "livekit_edge"
	NodeTypeMediaMTX    = "mediamtx"
	NodeTypeCDN         = "cdn"
)

// NodeStatus 常量
const (
	NodeStatusOnline    = "online"
	NodeStatusOffline   = "offline"
	NodeStatusDeploying = "deploying"
	NodeStatusFailed    = "failed"
)

// Node — 媒体节点表（对齐 migration 已有 schema）
type Node struct {
	ID                  uint64     `gorm:"primaryKey;column:id" json:"id"`
	NodeName            string     `gorm:"column:node_name;not null;size:100" json:"node_name"`
	NodeType            string     `gorm:"column:node_type;not null;size:30" json:"node_type"`
	PublicIP            string     `gorm:"column:public_ip;not null;size:50" json:"public_ip"`
	WireguardIP         string     `gorm:"column:wireguard_ip;not null;size:50" json:"wireguard_ip"`
	WireguardPublicKey  string     `gorm:"column:wireguard_public_key;not null;size:200" json:"wireguard_public_key"`
	Status              string     `gorm:"column:status;not null;size:20" json:"status"`
	LastHealthCheck     *time.Time `gorm:"column:last_health_check" json:"last_health_check,omitempty"`
	LastCPUPercent      int        `gorm:"column:last_cpu_percent" json:"last_cpu_percent,omitempty"`
	LastMemoryPercent   int        `gorm:"column:last_memory_percent" json:"last_memory_percent,omitempty"`
	LastBandwidthMbps   int        `gorm:"column:last_bandwidth_mbps" json:"last_bandwidth_mbps,omitempty"`
	LastViewersCount    int        `gorm:"column:last_viewers_count" json:"last_viewers_count,omitempty"`
	DeployedByStaffID   *uint64    `gorm:"column:deployed_by_staff_id" json:"deployed_by_staff_id,omitempty"`
	DeployedAt          *time.Time `gorm:"column:deployed_at" json:"deployed_at,omitempty"`
        ErrorMsg            string     `gorm:"column:error_message;size:500" json:"error_message,omitempty"`
        CreatedAt           time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt           *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (Node) TableName() string {
	return "nodes"
}
