package models

import (
	"time"
)

// PaymentTransaction statuses
const (
	PaymentStatusPending    = "pending"
	PaymentStatusProcessing = "processing"
	PaymentStatusSuccess    = "success"
	PaymentStatusFailed     = "failed"
	PaymentStatusRefunded   = "refunded"
	PaymentStatusDisputed   = "disputed"
)

// PaymentGateway
const (
	Gateway2Checkout = "2checkout"
	GatewayPayPal    = "paypal"
)

// PaymentTransaction — 支付交易表（gateway_transaction_id 必须 UNIQUE → 数据库层幂等）
type PaymentTransaction struct {
	ID                     uint64    `gorm:"primaryKey;column:id" json:"id"`
	OrderID                uint64    `gorm:"column:order_id;not null;index" json:"order_id"`
	PaymentGateway         string    `gorm:"column:payment_gateway;not null;size:20" json:"payment_gateway"`
	GatewayTransactionID   string    `gorm:"column:gateway_transaction_id;uniqueIndex;not null;size:100" json:"gateway_transaction_id"`
	Amount                 float64   `gorm:"column:amount;not null;type:numeric(10,2)" json:"amount"`
	Currency               string    `gorm:"column:currency;size:3;default:'GBP'" json:"currency"`
	Status                 string    `gorm:"column:status;not null;size:20" json:"status"`
	RawCallback            JSONMap   `gorm:"column:raw_callback;type:jsonb" json:"raw_callback,omitempty"`
	PaidAt                 *time.Time `gorm:"column:paid_at" json:"paid_at,omitempty"`
	CreatedAt              time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at;not null" json:"updated_at"`

	Order *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (PaymentTransaction) TableName() string {
	return "payment_transactions"
}

// Invoice — 商业发票表
type Invoice struct {
	ID             uint64    `gorm:"primaryKey;column:id" json:"id"`
	OrderID        uint64    `gorm:"column:order_id;uniqueIndex;not null" json:"order_id"`
	InvoiceNo      string    `gorm:"column:invoice_no;uniqueIndex;not null;size:30" json:"invoice_no"`
	IssueDate      time.Time `gorm:"column:issue_date;not null;type:date" json:"issue_date"`
	PdfURL         string    `gorm:"column:pdf_url;not null;size:500" json:"pdf_url"`
	Locked         bool      `gorm:"column:locked;default:true" json:"locked"`
	HsCode         string    `gorm:"column:hs_code;not null;size:20" json:"hs_code"`
	CountryOfOrigin string   `gorm:"column:country_of_origin;not null;size:50" json:"country_of_origin"`
	TotalAmount    float64   `gorm:"column:total_amount;not null;type:numeric(10,2)" json:"total_amount"`
	CreatedAt      time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null" json:"updated_at"`

	Order *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (Invoice) TableName() string {
	return "invoices"
}
