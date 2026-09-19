package models

import (
	"time"
)

// OrderStateLog — 订单状态流转历史（每次状态变更写一条）
// 用于前端 Timeline 展示，让用户/顾问看到每一步发生了什么
type OrderStateLog struct {
	ID        uint64    `gorm:"primaryKey;column:id" json:"id"`
	OrderID   uint64    `gorm:"column:order_id;not null;index" json:"order_id"`
	FromState string    `gorm:"column:from_state;not null;size:30" json:"from_state"` // 变更前状态（第一次流转前是空字符串）
	ToState   string    `gorm:"column:to_state;not null;size:30" json:"to_state"`     // 变更后状态
	Reason    string    `gorm:"column:reason;size:500" json:"reason,omitempty"`       // 顾问填写的备注（物流号、集装箱号等）
	StaffID   uint64    `gorm:"column:staff_id;not null" json:"staff_id"`             // 操作的顾问
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (OrderStateLog) TableName() string {
	return "order_state_logs"
}
