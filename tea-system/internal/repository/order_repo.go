package repository

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"tea-system/internal/models"
)

// ============================================================
// Order
// ============================================================

var (
	ErrOrderNotFound           = errors.New("order not found")
	ErrPaymentTransactionDup   = errors.New("duplicate payment transaction (idempotent)")
	ErrPaymentTransactionNotFound = errors.New("payment transaction not found")
)

// OrderRepo — 订单 + 支付交易仓储（统一在一个文件里，因为紧密关联）
type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

// generateOrderNo — 生成订单号 ORD-YYYYMMDD-6位随机
// 遇到 UNIQUE 冲突会重试最多 5 次
func generateOrderNo(ctx context.Context, db *gorm.DB) (string, error) {
	const maxRetry = 5
	for i := 0; i < maxRetry; i++ {
		now := time.Now()
		randPart := fmt.Sprintf("%06d", rand.Int63n(1_000_000))
		no := fmt.Sprintf("ORD-%s-%s", now.Format("20060102"), randPart)

		var count int64
		if err := db.WithContext(ctx).Model(&models.Order{}).
			Where("order_no = ?", no).
			Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return no, nil
		}
	}
	return "", errors.New("failed to generate unique order_no")
}

// Create — 创建订单（自动生成 order_no）
// snapshot 是 CustomProduct 的 JSON 快照（调用方负责生成）
func (r *OrderRepo) Create(ctx context.Context, o *models.Order) error {
	if o.OrderNo == "" {
		no, err := generateOrderNo(ctx, r.db)
		if err != nil {
			return err
		}
		o.OrderNo = no
	}
	return r.db.WithContext(ctx).Create(o).Error
}

// GetByID — 按主键查
func (r *OrderRepo) GetByID(ctx context.Context, id uint64) (*models.Order, error) {
	var o models.Order
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// GetByOrderNo — 按业务订单号查
func (r *OrderRepo) GetByOrderNo(ctx context.Context, orderNo string) (*models.Order, error) {
	var o models.Order
	err := r.db.WithContext(ctx).
		Where("order_no = ?", orderNo).
		First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// OrderListFilter — 列表过滤条件
type OrderListFilter struct {
	UserID uint64 // 0 = 全部
	State  string // "" = 全部
	Page   int
	Size   int
}

// List — 订单列表（支持 user_id / state 过滤 + 分页）
func (r *OrderRepo) List(ctx context.Context, filter OrderListFilter) ([]models.Order, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Size <= 0 || filter.Size > 200 {
		filter.Size = 20
	}

	q := r.db.WithContext(ctx).Model(&models.Order{})
	if filter.UserID > 0 {
		q = q.Where("user_id = ?", filter.UserID)
	}
	if filter.State != "" {
		q = q.Where("state = ?", filter.State)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.Order
	err := q.Order("created_at DESC").
		Offset((filter.Page - 1) * filter.Size).
		Limit(filter.Size).
		Find(&items).Error
	return items, total, err
}

// UpdateState — 原子更新状态字段（供状态机调用方使用）
func (r *OrderRepo) UpdateState(ctx context.Context, id uint64, newState string) error {
	res := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"state":      newState,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrOrderNotFound
	}
	return nil
}

// SoftDelete — 软删除订单（GORM DeletedAt）
func (r *OrderRepo) SoftDelete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Delete(&models.Order{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrOrderNotFound
	}
	return nil
}

// Update — 通用 patch 更新
func (r *OrderRepo) Update(ctx context.Context, id uint64, patch map[string]interface{}) error {
	if _, has := patch["updated_at"]; !has {
		patch["updated_at"] = time.Now()
	}
	res := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("id = ?", id).
		Updates(patch)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrOrderNotFound
	}
	return nil
}

// ============================================================
// PaymentTransaction — 依赖 gateway_transaction_id UNIQUE 索引实现幂等
// ============================================================

// CreatePayment — 创建支付交易，若遇到 UNIQUE 冲突返回 ErrPaymentTransactionDup（不是错误）
// 调用方可据此判断"已处理过"，实现 webhook 幂等。
func (r *OrderRepo) CreatePayment(ctx context.Context, tx *models.PaymentTransaction) error {
	err := r.db.WithContext(ctx).Create(tx).Error
	if err != nil {
		// 检查是否是唯一约束冲突
		if isUniqueViolation(err) {
			return ErrPaymentTransactionDup
		}
		return err
	}
	return nil
}

// GetPaymentByGatewayID — 按网关交易 ID 查（用于回调幂等检查）
func (r *OrderRepo) GetPaymentByGatewayID(ctx context.Context, gatewayTransactionID string) (*models.PaymentTransaction, error) {
	var t models.PaymentTransaction
	err := r.db.WithContext(ctx).
		Where("gateway_transaction_id = ?", gatewayTransactionID).
		First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPaymentTransactionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ListPaymentsByOrderID — 查某订单的所有支付交易
func (r *OrderRepo) ListPaymentsByOrderID(ctx context.Context, orderID uint64) ([]models.PaymentTransaction, error) {
	var list []models.PaymentTransaction
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

// MarkPaymentSucceeded — 标记交易成功
func (r *OrderRepo) MarkPaymentSucceeded(ctx context.Context, id uint64, raw models.JSONMap) error {
	now := time.Now()
	patch := map[string]interface{}{
		"status":      models.PaymentStatusSuccess,
		"paid_at":     now,
		"raw_callback": raw,
		"updated_at":  now,
	}
	return r.db.WithContext(ctx).Model(&models.PaymentTransaction{}).
		Where("id = ?", id).
		Updates(patch).Error
}

// MarkPaymentFailed — 标记交易失败
func (r *OrderRepo) MarkPaymentFailed(ctx context.Context, id uint64, raw models.JSONMap) error {
	return r.db.WithContext(ctx).Model(&models.PaymentTransaction{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       models.PaymentStatusFailed,
			"raw_callback": raw,
			"updated_at":   time.Now(),
		}).Error
}

// MarkPaymentStatus — 通用状态变更
func (r *OrderRepo) MarkPaymentStatus(ctx context.Context, id uint64, status string) error {
	return r.db.WithContext(ctx).Model(&models.PaymentTransaction{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// isUniqueViolation — 判断 error 是否来自 PostgreSQL unique_violation
// 不同驱动错误码不同，这里用 errors.As 做兼容
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	// pg driver: 23505
	type pgErr interface {
		SQLState() string
	}
	var pe pgErr
	if errors.As(err, &pe) {
		return pe.SQLState() == "23505"
	}
	// 兜底：看字符串包含 duplicate key
	return containsAny(err.Error(), "duplicate key", "unique constraint")
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(s) >= len(sub) {
			found := false
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					found = true
					break
				}
			}
			if found {
				return true
			}
		}
	}
	return false
}

// ListPaymentTransactions — 通用支付交易列表（admin 用）
func (r *OrderRepo) ListPaymentTransactions(ctx context.Context, gateway, status string, page, size int) ([]models.PaymentTransaction, int64, error) {
	var list []models.PaymentTransaction
	var total int64
	q := r.db.WithContext(ctx).Model(&models.PaymentTransaction{})
	if gateway != "" {
		q = q.Where("payment_gateway = ?", gateway)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
