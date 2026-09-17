package repository

import (
	"context"

	"gorm.io/gorm"

	"tea-system/internal/models"
)

// MessageRepo — 消息表操作
type MessageRepo struct {
	db *gorm.DB
}

func NewMessageRepo(db *gorm.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

// Create — 创建一条消息
func (r *MessageRepo) Create(ctx context.Context, msg *models.Message) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

// ListByConversationID — 按会话 ID 列出消息（分页，倒序）
// beforeID 用于游标翻页（返回该 ID 之前的 N 条）
func (r *MessageRepo) ListByConversationID(ctx context.Context, convID uint64, limit int, beforeID *uint64) ([]models.Message, error) {
	var msgs []models.Message
	q := r.db.WithContext(ctx).
		Where("conversation_id = ?", convID).
		Order("id DESC").
		Limit(limit)

	if beforeID != nil && *beforeID > 0 {
		q = q.Where("id < ?", *beforeID)
	}

	if err := q.Find(&msgs).Error; err != nil {
		return nil, err
	}
	// 反转时间线为正序
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}

// UpdateTranslation — 更新翻译结果
func (r *MessageRepo) UpdateTranslation(ctx context.Context, id uint64, zh, en, status string) error {
	updates := map[string]interface{}{
		"translation_zh":     zh,
		"translation_en":     en,
		"translation_status": status,
	}
	return r.db.WithContext(ctx).Model(&models.Message{}).Where("id = ?", id).Updates(updates).Error
}
