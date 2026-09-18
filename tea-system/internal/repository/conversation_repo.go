package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"tea-system/internal/models"
)

// ConversationRepo — 会话表操作
type ConversationRepo struct {
	db *gorm.DB
}

func NewConversationRepo(db *gorm.DB) *ConversationRepo {
	return &ConversationRepo{db: db}
}

// ErrConversationNotFound — 会话不存在
var ErrConversationNotFound = errors.New("conversation not found")

// Create — 创建会话（同时创建两个 participant）
func (r *ConversationRepo) Create(ctx context.Context, conv *models.Conversation, participants []models.ConversationParticipant) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(conv).Error; err != nil {
			return err
		}
		for i := range participants {
			participants[i].ConversationID = conv.ID
			if err := tx.Create(&participants[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByID — 按 ID 取会话（带 participants）
func (r *ConversationRepo) GetByID(ctx context.Context, id uint64) (*models.Conversation, error) {
	var conv models.Conversation
	if err := r.db.WithContext(ctx).
		Preload("Participants").
		First(&conv, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}
	return &conv, nil
}

// ListByParticipantID — 列出某用户/员工参与的所有会话
func (r *ConversationRepo) ListByParticipantID(ctx context.Context, userID, staffID *uint64) ([]models.Conversation, error) {
	var convs []models.Conversation
	q := r.db.WithContext(ctx).
		Joins("JOIN conversation_participants cp ON cp.conversation_id = conversations.id AND cp.deleted_at IS NULL")
	if userID != nil {
		q = q.Where("cp.user_id = ?", *userID)
	} else if staffID != nil {
		q = q.Where("cp.staff_id = ?", *staffID)
	} else {
		return convs, nil
	}

	if err := q.Preload("Participants").
		Order("conversations.updated_at DESC").
		Find(&convs).Error; err != nil {
		return nil, err
	}
	return convs, nil
}

// SoftDelete — 软删除会话
func (r *ConversationRepo) SoftDelete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Delete(&models.Conversation{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrConversationNotFound
	}
	return nil
}

// AddParticipant — 给会话增加参与方
func (r *ConversationRepo) AddParticipant(ctx context.Context, convID uint64, userID, staffID *uint64, role string) error {
	return r.db.WithContext(ctx).Create(&models.ConversationParticipant{
		ConversationID: convID,
		UserID:         userID,
		StaffID:        staffID,
		Role:           role,
	}).Error
}
