package models

import "time"

const (
	DSARTypeAccess       = "access"
	DSARTypeErasure      = "erasure"
	DSARTypeRectification = "rectification"
	DSARTypePortability  = "portability"

	DSARStatusPending    = "pending"
	DSARStatusProcessing = "processing"
	DSARStatusCompleted  = "completed"
	DSARStatusOverdue    = "overdue"
)

type DSARRequest struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      uint64     `gorm:"column:user_id;not null;index" json:"user_id"`
	RequestType string     `gorm:"column:request_type;size:20;not null" json:"request_type"`
	Status      string     `gorm:"column:status;size:20;not null;default:pending" json:"status"`
	Reason      string     `gorm:"column:reason;size:500" json:"reason,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	DueAt       time.Time  `gorm:"column:due_at;not null" json:"due_at"`
	CompletedAt *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
}

func (DSARRequest) TableName() string { return "dsar_requests" }
