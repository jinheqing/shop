package models

import (
	"time"
)

// StaffRole — 内部角色常量
const (
	RoleAdvisor    = "advisor"
	RoleSupervisor = "supervisor"
	RoleAdmin      = "admin"
	RoleTeaFarmer  = "tea_farmer"
	RoleOperations = "operations"
)

// Staff — 内部员工表
type Staff struct {
	ID               uint64    `gorm:"primaryKey;column:id" json:"id"`
	Name             string    `gorm:"column:name;not null;size:100" json:"name"`
	Email            string    `gorm:"column:email;uniqueIndex;not null;size:200" json:"email"`
	PasswordHash     string    `gorm:"column:password_hash;not null;size:255" json:"-"`
	Role             string    `gorm:"column:role;not null;size:30" json:"role"`
	Permissions      JSONMap   `gorm:"column:permissions;type:jsonb" json:"permissions"`
	MfaEnabled       bool      `gorm:"column:mfa_enabled;default:false" json:"mfa_enabled"`
	MfaSecret        string    `gorm:"column:mfa_secret;size:100" json:"-"`
	WorkTimezone     string    `gorm:"column:work_timezone;size:50;default:'Europe/London'" json:"work_timezone"`
	AssignedFarmID   *uint64   `gorm:"column:assigned_farm_id" json:"assigned_farm_id,omitempty"`
	IsActive         bool      `gorm:"column:is_active;default:true" json:"is_active"`
	LastLoginAt      *time.Time `gorm:"column:last_login_at" json:"last_login_at,omitempty"`
	CreatedAt        time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (Staff) TableName() string {
	return "staff"
}

// ForceMFA — admin/supervisor 角色强制 MFA
func (s *Staff) ForceMFA() bool {
	return s.Role == RoleAdmin || s.Role == RoleSupervisor
}
