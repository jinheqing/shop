package models

import (
	"time"
)

// DeclarationLedger — 物流报关台账
type DeclarationLedger struct {
	ID                       uint64     `gorm:"primaryKey;column:id" json:"id"`
	OrderID                  uint64     `gorm:"column:order_id;not null" json:"order_id"`
	CustomsDeclarationNo     string     `gorm:"column:customs_declaration_no;size:50" json:"customs_declaration_no,omitempty"`
	HsCode                   string     `gorm:"column:hs_code;not null;size:20" json:"hs_code"`
	CommodityDesc            string     `gorm:"column:commodity_desc;not null;type:text" json:"commodity_desc"`
	GrossWeight              *float64   `gorm:"column:gross_weight;type:numeric(10,3)" json:"gross_weight,omitempty"`
	NetWeight                *float64   `gorm:"column:net_weight;type:numeric(10,3)" json:"net_weight,omitempty"`
	DeclaredValue            *float64   `gorm:"column:declared_value;type:numeric(10,2)" json:"declared_value,omitempty"`
	CustomsStatus            string     `gorm:"column:customs_status;size:30" json:"customs_status,omitempty"` // pending/submitted/cleared/held
	DeclarationDate          *time.Time `gorm:"column:declaration_date;type:date" json:"declaration_date,omitempty"`
	ClearedDate              *time.Time `gorm:"column:cleared_date;type:date" json:"cleared_date,omitempty"`
	Remarks                  string     `gorm:"column:remarks;type:text" json:"remarks,omitempty"`
	VoidedAt                 *time.Time `gorm:"column:voided_at" json:"voided_at,omitempty"`
	VoidedReason             string     `gorm:"column:voided_reason;type:text" json:"voided_reason,omitempty"`
	CreatedAt                time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt                time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`

	Order *Order `gorm:"-:migration;foreignKey:OrderID" json:"order,omitempty"`
}

func (DeclarationLedger) TableName() string {
	return "declaration_ledger"
}

// ForeignExchangeLedger — 收汇台账
type ForeignExchangeLedger struct {
	ID                   uint64     `gorm:"primaryKey;column:id" json:"id"`
	OrderID              uint64     `gorm:"column:order_id;not null" json:"order_id"`
	PaymentGateway       string     `gorm:"column:payment_gateway;not null;size:20" json:"payment_gateway"`
	GatewayTransactionID string     `gorm:"column:gateway_transaction_id;not null;size:100" json:"gateway_transaction_id"`
	AmountGBP            float64    `gorm:"column:amount_gbp;not null;type:numeric(10,2)" json:"amount_gbp"`
	ExchangeRate         *float64   `gorm:"column:exchange_rate;type:numeric(10,6)" json:"exchange_rate,omitempty"`
	AmountCNY            *float64   `gorm:"column:amount_cny;type:numeric(12,2)" json:"amount_cny,omitempty"`
	SettlementDate       *time.Time `gorm:"column:settlement_date;type:date" json:"settlement_date,omitempty"`
	BankAccount          string     `gorm:"column:bank_account;size:100" json:"bank_account,omitempty"`
	Remarks              string     `gorm:"column:remarks;type:text" json:"remarks,omitempty"`
	CreatedAt            time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`

	Order *Order `gorm:"-:migration;foreignKey:OrderID" json:"order,omitempty"`
}

func (ForeignExchangeLedger) TableName() string {
	return "foreign_exchange_ledger"
}
