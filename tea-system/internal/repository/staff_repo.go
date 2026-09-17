package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"tea-system/internal/models"
)

var ErrStaffNotFound = errors.New("staff not found")

// StaffRepo — Staff CRUD
type StaffRepo struct {
	db *gorm.DB
}

func NewStaffRepo(db *gorm.DB) *StaffRepo {
	return &StaffRepo{db: db}
}

// GetByEmail — 登录时按 email 查（不软删除过滤：需要能找到被禁用的用户然后返回错误）
func (r *StaffRepo) GetByEmail(ctx context.Context, email string) (*models.Staff, error) {
	var staff models.Staff
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&staff).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrStaffNotFound
	}
	if err != nil {
		return nil, err
	}
	return &staff, nil
}

// GetByID — 按 id 查
func (r *StaffRepo) GetByID(ctx context.Context, id uint64) (*models.Staff, error) {
	var staff models.Staff
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&staff).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrStaffNotFound
	}
	return &staff, err
}

// UpdateLastLogin — 更新最后登录时间
func (r *StaffRepo) UpdateLastLogin(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&models.Staff{}).
		Where("id = ?", id).
		Update("last_login_at", gorm.Expr("NOW()")).Error
}

// Create — 创建新 Staff
func (r *StaffRepo) Create(ctx context.Context, staff *models.Staff) error {
	return r.db.WithContext(ctx).Create(staff).Error
}

// List — Staff 列表（分页）
func (r *StaffRepo) List(ctx context.Context, page, size int) ([]models.Staff, int64, error) {
	var staffs []models.Staff
	var total int64

	r.db.WithContext(ctx).Model(&models.Staff{}).Count(&total)
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&staffs).Error

	return staffs, total, err
}

// ToggleActive — 启用/禁用 Staff
func (r *StaffRepo) ToggleActive(ctx context.Context, id uint64, active bool) error {
	return r.db.WithContext(ctx).Model(&models.Staff{}).
		Where("id = ?", id).
		Update("is_active", active).Error
}

// SoftDelete — 软删除 Staff（置 is_active=false + email 加 deleted_ 前缀）
func (r *StaffRepo) SoftDelete(ctx context.Context, id uint64) error {
	var s models.Staff
	if err := r.db.WithContext(ctx).First(&s, id).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&s).Updates(map[string]interface{}{
		"is_active": false,
		"email":     "deleted_" + s.Email,
	}).Error
}
