package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"tea-system/internal/models"
)

var ErrLiveRoomNotFound = errors.New("live room not found")

// LiveRoomListFilter — 直播间列表过滤条件
type LiveRoomListFilter struct {
	RoomType   string
	Status     string
	Visibility string
	Type       string // 2026-09 新增：slow_live / scheduled / advisor / admin
	OrderID    *uint64
	HostID     *uint64
	Keyword    string // 模糊匹配 room_name / room_id
	Page       int
	Size       int
}

// LiveRoomRepo — 直播间 CRUD
type LiveRoomRepo struct {
	db *gorm.DB
}

func NewLiveRoomRepo(db *gorm.DB) *LiveRoomRepo {
	return &LiveRoomRepo{db: db}
}

// Create — 新建直播间
func (r *LiveRoomRepo) Create(ctx context.Context, room *models.LiveRoom) error {
	return r.db.WithContext(ctx).Create(room).Error
}

// GetByID — 按主键查
func (r *LiveRoomRepo) GetByID(ctx context.Context, id uint64) (*models.LiveRoom, error) {
	var room models.LiveRoom
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrLiveRoomNotFound
	}
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// GetByRoomID — 按业务 room_id 查（UNIQUE）
func (r *LiveRoomRepo) GetByRoomID(ctx context.Context, roomID string) (*models.LiveRoom, error) {
	var room models.LiveRoom
	err := r.db.WithContext(ctx).Where("room_id = ?", roomID).First(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrLiveRoomNotFound
	}
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// List — 列表（支持过滤 + 分页）
func (r *LiveRoomRepo) List(ctx context.Context, f LiveRoomListFilter) ([]models.LiveRoom, int64, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 || f.Size > 200 {
		f.Size = 20
	}

	q := r.db.WithContext(ctx).Model(&models.LiveRoom{})
	if f.RoomType != "" {
		q = q.Where("room_type = ?", f.RoomType)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Visibility != "" {
		q = q.Where("visibility = ?", f.Visibility)
	}
	if f.Type != "" {
		q = q.Where("type = ?", f.Type)
	}
	if f.OrderID != nil {
		q = q.Where("order_id = ?", *f.OrderID)
	}
	if f.HostID != nil {
		q = q.Where("host_staff_id = ?", *f.HostID)
	}
	if f.Keyword != "" {
		q = q.Where("room_name ILIKE ? OR room_id ILIKE ?", "%"+f.Keyword+"%", "%"+f.Keyword+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.LiveRoom
	err := q.Order("created_at DESC").
		Offset((f.Page - 1) * f.Size).
		Limit(f.Size).
		Find(&items).Error
	return items, total, err
}

// Update — patch 更新
func (r *LiveRoomRepo) Update(ctx context.Context, id uint64, patch map[string]interface{}) error {
	patch["updated_at"] = time.Now()
	res := r.db.WithContext(ctx).Model(&models.LiveRoom{}).
		Where("id = ?", id).
		Updates(patch)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrLiveRoomNotFound
	}
	return nil
}

// SoftDelete — 软删除
func (r *LiveRoomRepo) SoftDelete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Model(&models.LiveRoom{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     models.LiveStatusOffline,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrLiveRoomNotFound
	}
	return r.db.WithContext(ctx).Delete(&models.LiveRoom{}, id).Error
}

// Start — 开始直播：status→live, started_at=now
func (r *LiveRoomRepo) Start(ctx context.Context, id uint64) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&models.LiveRoom{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     models.LiveStatusLive,
			"started_at":  now,
			"updated_at": now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrLiveRoomNotFound
	}
	return nil
}

// End — 结束直播：status→ended, ended_at=now
func (r *LiveRoomRepo) End(ctx context.Context, id uint64) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&models.LiveRoom{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     models.LiveStatusEnded,
			"ended_at":   now,
			"updated_at": now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrLiveRoomNotFound
	}
	return nil
}

// GetActiveByOrderID — 某订单是否有关联的 delivery_inspection 直播（status 不是 offline/ended）
func (r *LiveRoomRepo) GetActiveByOrderID(ctx context.Context, orderID uint64) (*models.LiveRoom, error) {
	var room models.LiveRoom
	err := r.db.WithContext(ctx).
		Where("order_id = ? AND room_type = ? AND status NOT IN ?",
			orderID, models.RoomTypeDeliveryInspection,
			[]string{models.LiveStatusOffline, models.LiveStatusEnded}).
		First(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &room, nil
}
