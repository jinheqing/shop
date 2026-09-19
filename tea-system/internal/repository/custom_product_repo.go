package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"tea-system/internal/models"
)

var ErrCustomProductNotFound = errors.New("custom product not found")

// CustomProductRepo — 定制报价表 CRUD
type CustomProductRepo struct {
	db *gorm.DB
}

func NewCustomProductRepo(db *gorm.DB) *CustomProductRepo {
	return &CustomProductRepo{db: db}
}

// Create — 新建报价
func (r *CustomProductRepo) Create(ctx context.Context, p *models.CustomProduct) error {
	return r.db.WithContext(ctx).Create(p).Error
}

// GetByID — 按主键查询（默认过滤软删除）
func (r *CustomProductRepo) GetByID(ctx context.Context, id uint64) (*models.CustomProduct, error) {
	var p models.CustomProduct
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCustomProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetByToken — 按 product_token 查询（公开接口用，不限状态）
func (r *CustomProductRepo) GetByToken(ctx context.Context, token string) (*models.CustomProduct, error) {
	var p models.CustomProduct
	err := r.db.WithContext(ctx).
		Where("product_token = ?", token).
		First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCustomProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListFilter — 列表过滤条件
type ListFilter struct {
	Status  string // "" 表示全部
	Keyword string // 模糊匹配 title
	Page    int
	Size    int
}

// List — staff 列表（支持分页 + 状态过滤 + 关键词）
func (r *CustomProductRepo) List(ctx context.Context, filter ListFilter) ([]models.CustomProduct, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Size <= 0 || filter.Size > 200 {
		filter.Size = 20
	}

	q := r.db.WithContext(ctx).Model(&models.CustomProduct{})
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.Keyword != "" {
		q = q.Where("title ILIKE ?", "%"+filter.Keyword+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.CustomProduct
	err := q.Order("created_at DESC").
		Offset((filter.Page - 1) * filter.Size).
		Limit(filter.Size).
		Find(&items).Error
	return items, total, err
}

// Update — 字段 patch 更新（零值安全，用 map）
// 同时 bump version，刷新 updated_at
func (r *CustomProductRepo) Update(ctx context.Context, id uint64, patch map[string]interface{}) error {
	// 先取当前 version
	var current models.CustomProduct
	if err := r.db.WithContext(ctx).Select("version").Where("id = ?", id).First(&current).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCustomProductNotFound
		}
		return err
	}

	patch["version"] = current.Version + 1
	if _, has := patch["updated_at"]; !has {
		patch["updated_at"] = time.Now()
	}

	res := r.db.WithContext(ctx).Model(&models.CustomProduct{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrCustomProductNotFound
	}
	return nil
}

// SoftDelete — 软删除（GORM DeletedAt 生效）+ 同时标记 archived
func (r *CustomProductRepo) SoftDelete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Model(&models.CustomProduct{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     models.CustomProductStatusArchived,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrCustomProductNotFound
	}
	// 物理软删
	return r.db.WithContext(ctx).Delete(&models.CustomProduct{}, id).Error
}

// Publish — 生成 token + 状态置 published
func (r *CustomProductRepo) Publish(ctx context.Context, id uint64, token string) error {
	res := r.db.WithContext(ctx).Model(&models.CustomProduct{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"product_token": &token,
			"status":        models.CustomProductStatusPublished,
			"version":       gorm.Expr("version + 1"),
			"updated_at":    time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrCustomProductNotFound
	}
	return nil
}

// Review — 审核通过
func (r *CustomProductRepo) Review(ctx context.Context, id uint64, staffID uint64) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&models.CustomProduct{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"reviewed_by_staff_id": staffID,
			"reviewed_at":          now,
			"updated_at":           now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrCustomProductNotFound
	}
	return nil
}
