package models

import (
	"time"
)

// Conversation — IM 会话表
type Conversation struct {
	ID              uint64    `gorm:"primaryKey;column:id" json:"id"`
	ConversationType string   `gorm:"column:conversation_type;not null;size:20" json:"conversation_type"` // single / group
	Title           string    `gorm:"column:title;size:200" json:"title,omitempty"`
	CreatorID       uint64    `gorm:"column:creator_id;not null" json:"creator_id"`
	LastMessageAt   *time.Time `gorm:"column:last_message_at" json:"last_message_at,omitempty"`
	CreatedAt       time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"-"`

	Participants []ConversationParticipant `gorm:"foreignKey:ConversationID" json:"participants,omitempty"`
	Messages     []Message                `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

func (Conversation) TableName() string {
	return "conversations"
}

// ConversationParticipant — 会话参与方
type ConversationParticipant struct {
	ID             uint64    `gorm:"primaryKey;column:id" json:"id"`
	ConversationID uint64    `gorm:"column:conversation_id;not null;index" json:"conversation_id"`
	UserID         *uint64   `gorm:"column:user_id" json:"user_id,omitempty"`
	StaffID        *uint64   `gorm:"column:staff_id" json:"staff_id,omitempty"`
	Role           string    `gorm:"column:role;size:20;default:'member'" json:"role"` // owner / member
	LastReadAt     *time.Time `gorm:"column:last_read_at" json:"last_read_at,omitempty"`
	CreatedAt      time.Time `gorm:"column:created_at;not null" json:"created_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (ConversationParticipant) TableName() string {
	return "conversation_participants"
}

// Message — IM 消息表
type Message struct {
	ID               uint64    `gorm:"primaryKey;column:id" json:"id"`
	ConversationID   uint64    `gorm:"column:conversation_id;not null;index" json:"conversation_id"`
	SenderType       string    `gorm:"column:sender_type;not null;size:10" json:"sender_type"` // user / staff / system
	SenderID         uint64    `gorm:"column:sender_id;not null" json:"sender_id"`
	MessageType      string    `gorm:"column:message_type;not null;size:20" json:"message_type"` // text / image / quote_link / file
	Content          string    `gorm:"column:content;not null;type:text" json:"content"`
	TranslationZH    string    `gorm:"column:translation_zh;type:text" json:"translation_zh,omitempty"`
	TranslationEN    string    `gorm:"column:translation_en;type:text" json:"translation_en,omitempty"`
	TranslationStatus string   `gorm:"column:translation_status;size:20;default:'pending'" json:"translation_status"` // pending / translated / failed
	ReplyToID        *uint64   `gorm:"column:reply_to_id" json:"reply_to_id,omitempty"`
	AttachmentURLs   JSONMap   `gorm:"column:attachment_urls;type:jsonb" json:"attachment_urls,omitempty"`
	CreatedAt        time.Time `gorm:"column:created_at;not null;index" json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}
