package models

import (
	"time"
)

// UserGroup — 用户组（Annabel's 式的私密分组，不搞公开等级）
type UserGroup struct {
	ID          uint64    `gorm:"primaryKey;column:id" json:"id"`
	Name        string    `gorm:"column:name;not null;size:100" json:"name"`           // "VIP 2026" / "勐海古树 2026 批次买家"
	Description string    `gorm:"column:description;type:text" json:"description,omitempty"`
	AutoRule    JSONMap   `gorm:"column:auto_rule;type:jsonb" json:"auto_rule,omitempty"` // {spend_threshold: 5000, registered_before: "2025-01-01"}
	Privileges  JSONArray `gorm:"column:privileges;type:jsonb" json:"privileges,omitempty"` // ["preorder_priority", "offline_garden_tour:batch_2026_menghai", ...]
	Reciprocal  JSONArray `gorm:"column:reciprocal;type:jsonb" json:"reciprocal,omitempty"` // ["annabels", "soho_house"]
	CreatedByStaffID *uint64 `gorm:"column:created_by_staff_id" json:"created_by_staff_id,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"-"`

	// 成员列表（optional preload）
	Members []UserGroupMember `gorm:"foreignKey:GroupID" json:"members,omitempty"`
}

func (UserGroup) TableName() string {
	return "user_groups"
}

// UserGroupMember — 组<->用户关联表
type UserGroupMember struct {
	ID             uint64    `gorm:"primaryKey;column:id" json:"id"`
	GroupID        uint64    `gorm:"column:group_id;not null;uniqueIndex:grp_user" json:"group_id"`
	UserID         uint64    `gorm:"column:user_id;not null;uniqueIndex:grp_user" json:"user_id"`
	AddedAt        time.Time `gorm:"column:added_at;not null" json:"added_at"`
	AddedByStaffID *uint64   `gorm:"column:added_by_staff_id" json:"added_by_staff_id,omitempty"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (UserGroupMember) TableName() string {
	return "user_group_members"
}
