package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"tea-system/internal/models"
)

var ErrNodeNotFound = errors.New("node not found")

// NodeRepo — 节点 CRUD
type NodeRepo struct {
	db *gorm.DB
}

func NewNodeRepo(db *gorm.DB) *NodeRepo {
	return &NodeRepo{db: db}
}

// Create — 新增节点
func (r *NodeRepo) Create(ctx context.Context, node *models.Node) error {
	return r.db.WithContext(ctx).Create(node).Error
}

// GetByID — 按主键查
func (r *NodeRepo) GetByID(ctx context.Context, id uint64) (*models.Node, error) {
	var node models.Node
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&node).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNodeNotFound
	}
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// List — 全量列表
func (r *NodeRepo) List(ctx context.Context) ([]models.Node, error) {
	var nodes []models.Node
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&nodes).Error
	return nodes, err
}

// Update — patch 更新
func (r *NodeRepo) Update(ctx context.Context, id uint64, patch map[string]interface{}) error {
	patch["updated_at"] = time.Now()
	res := r.db.WithContext(ctx).Model(&models.Node{}).
		Where("id = ?", id).
		Updates(patch)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNodeNotFound
	}
	return nil
}

// SoftDelete — 软删除（置 offline + gorm soft delete）
func (r *NodeRepo) SoftDelete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Model(&models.Node{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      models.NodeStatusOffline,
			"error_message": "deleted",
			"updated_at":  time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNodeNotFound
	}
	return r.db.WithContext(ctx).Delete(&models.Node{}, id).Error
}

// UpdateStatus — 更新状态
func (r *NodeRepo) UpdateStatus(ctx context.Context, id uint64, status, errMsg string) error {
	patch := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if errMsg != "" {
		patch["error_message"] = errMsg
	} else {
		patch["error_message"] = ""
	}
	return r.Update(ctx, id, patch)
}

// UpdateHealthCheck — 更新健康检查时间
func (r *NodeRepo) UpdateHealthCheck(ctx context.Context, id uint64, ok bool, errMsg string) error {
	now := time.Now()
	patch := map[string]interface{}{
		"last_health_check": now,
		"updated_at":        now,
	}
	if ok {
		patch["status"] = models.NodeStatusOnline
		patch["error_message"] = ""
	} else {
		patch["status"] = models.NodeStatusOffline
		if errMsg != "" {
			patch["error_message"] = errMsg
		}
	}
	return r.Update(ctx, id, patch)
}
