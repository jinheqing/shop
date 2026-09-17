package repository

import (
	"context"
	"errors"

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
