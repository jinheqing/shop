package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"tea-system/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

// UserRepo — User CRUD
type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &u, err
}

func (r *UserRepo) GetByID(ctx context.Context, id uint64) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &u, err
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// UpdateLoginGeo — 更新用户最近登录的 IP 和归属地
func (r *UserRepo) UpdateLoginGeo(ctx context.Context, userID uint64, ip, city, country, countryCode, region string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_ip":           ip,
			"last_login_city":         city,
			"last_login_country":      country,
			"last_login_country_code": countryCode,
			"last_login_region":       region,
			"last_login_at":           now,
		}).Error
}

// List — User 列表（分页，admin 用）
func (r *UserRepo) List(ctx context.Context, page, size int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.WithContext(ctx).Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&users).Error
	return users, total, err
}
