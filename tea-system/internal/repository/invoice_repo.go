package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"tea-system/internal/models"
)

// ============================================================
// InvoiceRepo — 发票 / 报关台账 / 收汇台账 / SGS 报告
// ============================================================

var (
	ErrInvoiceNotFound           = errors.New("invoice not found")
	ErrDeclarationNotFound       = errors.New("declaration ledger not found")
	ErrForeignExchangeNotFound   = errors.New("foreign exchange ledger not found")
	ErrSgsReportNotFound         = errors.New("sgs report not found")
)

// InvoiceRepo — 发票 + 台账 + SGS 报告 仓储聚合
type InvoiceRepo struct {
	db *gorm.DB
}

func NewInvoiceRepo(db *gorm.DB) *InvoiceRepo {
	return &InvoiceRepo{db: db}
}

// ==================== Invoice ====================

func (r *InvoiceRepo) CreateInvoice(ctx context.Context, inv *models.Invoice) error {
	return r.db.WithContext(ctx).Create(inv).Error
}

func (r *InvoiceRepo) GetInvoiceByID(ctx context.Context, id uint64) (*models.Invoice, error) {
	var inv models.Invoice
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&inv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvoiceNotFound
	}
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *InvoiceRepo) GetInvoiceByOrderID(ctx context.Context, orderID uint64) (*models.Invoice, error) {
	var inv models.Invoice
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&inv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvoiceNotFound
	}
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *InvoiceRepo) UpdateInvoicePDFURL(ctx context.Context, id uint64, pdfURL string) error {
	res := r.db.WithContext(ctx).Model(&models.Invoice{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"pdf_url":    pdfURL,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrInvoiceNotFound
	}
	return nil
}

// ==================== DeclarationLedger ====================

func (r *InvoiceRepo) CreateDeclaration(ctx context.Context, d *models.DeclarationLedger) error {
	return r.db.WithContext(ctx).Create(d).Error
}

// DeclarationListFilter — 列表过滤
type DeclarationListFilter struct {
	OrderID       uint64
	CustomsStatus string
	Page          int
	Size          int
}

func (r *InvoiceRepo) ListDeclarations(ctx context.Context, filter DeclarationListFilter) ([]models.DeclarationLedger, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Size <= 0 || filter.Size > 200 {
		filter.Size = 20
	}
	q := r.db.WithContext(ctx).Model(&models.DeclarationLedger{}).Where("voided_at IS NULL")
	if filter.OrderID > 0 {
		q = q.Where("order_id = ?", filter.OrderID)
	}
	if filter.CustomsStatus != "" {
		q = q.Where("customs_status = ?", filter.CustomsStatus)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []models.DeclarationLedger
	err := q.Order("created_at DESC").
		Offset((filter.Page - 1) * filter.Size).
		Limit(filter.Size).
		Find(&list).Error
	return list, total, err
}

func (r *InvoiceRepo) UpdateDeclaration(ctx context.Context, id uint64, patch map[string]interface{}) error {
	if _, ok := patch["updated_at"]; !ok {
		patch["updated_at"] = time.Now()
	}
	res := r.db.WithContext(ctx).Model(&models.DeclarationLedger{}).
		Where("id = ? AND voided_at IS NULL", id).
		Updates(patch)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrDeclarationNotFound
	}
	return nil
}

// SoftDeleteDeclaration — 置 voided
func (r *InvoiceRepo) SoftDeleteDeclaration(ctx context.Context, id uint64, reason string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&models.DeclarationLedger{}).
		Where("id = ? AND voided_at IS NULL", id).
		Updates(map[string]interface{}{
			"voided_at":     now,
			"voided_reason": reason,
			"updated_at":    now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrDeclarationNotFound
	}
	return nil
}

func (r *InvoiceRepo) GetDeclarationByID(ctx context.Context, id uint64) (*models.DeclarationLedger, error) {
	var d models.DeclarationLedger
	err := r.db.WithContext(ctx).
		Where("id = ? AND voided_at IS NULL", id).
		First(&d).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrDeclarationNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// ==================== ForeignExchangeLedger ====================

func (r *InvoiceRepo) CreateForeignExchange(ctx context.Context, f *models.ForeignExchangeLedger) error {
	return r.db.WithContext(ctx).Create(f).Error
}

// ForexListFilter
type ForexListFilter struct {
	OrderID        uint64
	PaymentGateway string
	Page           int
	Size           int
}

func (r *InvoiceRepo) ListForeignExchange(ctx context.Context, filter ForexListFilter) ([]models.ForeignExchangeLedger, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Size <= 0 || filter.Size > 200 {
		filter.Size = 20
	}
	q := r.db.WithContext(ctx).Model(&models.ForeignExchangeLedger{})
	if filter.OrderID > 0 {
		q = q.Where("order_id = ?", filter.OrderID)
	}
	if filter.PaymentGateway != "" {
		q = q.Where("payment_gateway = ?", filter.PaymentGateway)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []models.ForeignExchangeLedger
	err := q.Order("created_at DESC").
		Offset((filter.Page - 1) * filter.Size).
		Limit(filter.Size).
		Find(&list).Error
	return list, total, err
}

func (r *InvoiceRepo) UpdateForeignExchange(ctx context.Context, id uint64, patch map[string]interface{}) error {
	if _, ok := patch["updated_at"]; !ok {
		patch["updated_at"] = time.Now()
	}
	res := r.db.WithContext(ctx).Model(&models.ForeignExchangeLedger{}).
		Where("id = ?", id).
		Updates(patch)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrForeignExchangeNotFound
	}
	return nil
}

func (r *InvoiceRepo) SoftDeleteForeignExchange(ctx context.Context, id uint64) error {
	// ForeignExchangeLedger 没有 DeletedAt，做 delete 即可（软删除语义：置空关键字段）
	res := r.db.WithContext(ctx).Delete(&models.ForeignExchangeLedger{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrForeignExchangeNotFound
	}
	return nil
}

func (r *InvoiceRepo) GetForeignExchangeByID(ctx context.Context, id uint64) (*models.ForeignExchangeLedger, error) {
	var f models.ForeignExchangeLedger
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&f).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrForeignExchangeNotFound
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// ==================== SgsReport ====================

func (r *InvoiceRepo) CreateSgsReport(ctx context.Context, s *models.SgsReport) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *InvoiceRepo) GetSgsReportByID(ctx context.Context, id uint64) (*models.SgsReport, error) {
	var s models.SgsReport
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSgsReportNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetSgsReportByReportNo — 公开查询（用 ReportNo 作为公开 token，因为 ReportNo 是 uniqueIndex）
func (r *InvoiceRepo) GetSgsReportByReportNo(ctx context.Context, reportNo string) (*models.SgsReport, error) {
	var s models.SgsReport
	err := r.db.WithContext(ctx).Where("report_no = ?", reportNo).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSgsReportNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// SgsListFilter
type SgsListFilter struct {
	TeaType  string
	BatchNo  string
	Page     int
	Size     int
}

func (r *InvoiceRepo) ListSgsReports(ctx context.Context, filter SgsListFilter) ([]models.SgsReport, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Size <= 0 || filter.Size > 200 {
		filter.Size = 20
	}
	q := r.db.WithContext(ctx).Model(&models.SgsReport{})
	if filter.TeaType != "" {
		q = q.Where("tea_type = ?", filter.TeaType)
	}
	if filter.BatchNo != "" {
		q = q.Where("batch_no = ?", filter.BatchNo)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []models.SgsReport
	err := q.Order("created_at DESC").
		Offset((filter.Page - 1) * filter.Size).
		Limit(filter.Size).
		Find(&list).Error
	return list, total, err
}

func (r *InvoiceRepo) UpdateSgsReport(ctx context.Context, id uint64, patch map[string]interface{}) error {
	if _, ok := patch["updated_at"]; !ok {
		patch["updated_at"] = time.Now()
	}
	res := r.db.WithContext(ctx).Model(&models.SgsReport{}).
		Where("id = ?", id).
		Updates(patch)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrSgsReportNotFound
	}
	return nil
}

// DeleteSgsReport — SgsReport 没有 DeletedAt，物理删除
func (r *InvoiceRepo) DeleteSgsReport(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Delete(&models.SgsReport{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrSgsReportNotFound
	}
	return nil
}
