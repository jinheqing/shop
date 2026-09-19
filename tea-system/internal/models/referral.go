package models

import (
	"time"
)

// Referral — 推荐关系表（顾问管理，私密信息，不公开）
//
// 触发两种方式：
//  1. 新客户注册/结账时填了 "朋友推荐" + 推荐人名字 → 系统尝试匹配，匹配上了自动建记录
//  2. 顾问在后台手动创建
//
// 奖励触发条件：
//
//	friend_order_amount >= threshold（配置项，默认 £100）→ 给推荐人发感谢邮件 + 累计 referral_count
//	referral_count >= 3（配置项）→ 自动把推荐人拉进 "KOL 伙伴" 用户组
type Referral struct {
	ID             uint64  `gorm:"primaryKey;column:id" json:"id"`
	ReferrerUserID *uint64 `gorm:"column:referrer_user_id" json:"referrer_user_id,omitempty"`   // 推荐人用户 ID（可能为空如果只知道名字）
	ReferrerName   string  `gorm:"column:referrer_name;size:100" json:"referrer_name"`          // 推荐人名字（冗余存，方便搜索）
	ReferredUserID uint64  `gorm:"column:referred_user_id;uniqueIndex" json:"referred_user_id"` // 被推荐人用户 ID（唯一，一个人只能有一个推荐人）
	ReferredName   string  `gorm:"column:referred_name;size:100" json:"referred_name"`          // 被推荐人名字（冗余存）
	Source         string  `gorm:"column:source;size:20;default:'name_share'" json:"source"`    // name_share / short_code / manual

	// 状态追踪
	FriendOrderID     *uint64    `gorm:"column:friend_order_id" json:"friend_order_id,omitempty"`         // 朋友下单的订单 ID
	FriendOrderAmount float64    `gorm:"column:friend_order_amount;default:0" json:"friend_order_amount"` // 朋友下单金额
	RewardTriggeredAt *time.Time `gorm:"column:reward_triggered_at" json:"reward_triggered_at,omitempty"` // 奖励触发时间

	CreatedByStaffID *uint64 `gorm:"column:created_by_staff_id" json:"created_by_staff_id,omitempty"`
	Notes            string  `gorm:"column:notes;type:text" json:"notes,omitempty"`

	CreatedAt time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`

	// Associations
	ReferrerUser *User `gorm:"-:migration;foreignKey:ReferrerUserID" json:"referrer_user,omitempty"`
	ReferredUser *User `gorm:"-:migration;foreignKey:ReferredUserID" json:"referred_user,omitempty"`
}

func (Referral) TableName() string {
	return "referrals"
}

// IsRewardEligible 朋友下单后调用，判断是否达到奖励门槛
func (r *Referral) IsRewardEligible(threshold float64) bool {
	return r.FriendOrderAmount >= threshold && r.RewardTriggeredAt == nil
}
