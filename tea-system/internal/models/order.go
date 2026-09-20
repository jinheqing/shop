package models

import (
	"time"
)

// Order states
const (
	OrderStateOrdering            = "ordering"
	OrderStatePaid                = "paid"
	OrderStatePendingDeclaration  = "pending_declaration"
	OrderStateProducing           = "producing"
	OrderStateReadyForProduction  = "ready_for_production"
	OrderStatePendingCustoms      = "pending_customs"
	OrderStateCustomsClear        = "customs_clear"
	OrderStateShipped             = "shipped"
	OrderStateCompleted           = "completed"
	OrderStateCancelled           = "cancelled"
	OrderStateDisputed            = "disputed"
)

// Order — 订单表
type Order struct {
	ID                    uint64    `gorm:"primaryKey;column:id" json:"id"`
	OrderNo               string    `gorm:"column:order_no;uniqueIndex;not null;size:30" json:"order_no"`
	UserID                *uint64   `gorm:"column:user_id;index" json:"user_id"`
	StaffID               uint64    `gorm:"column:staff_id;not null" json:"staff_id"`
	CustomProductID       uint64    `gorm:"column:custom_product_id;not null" json:"custom_product_id"`
	CustomProductSnapshot JSONMap   `gorm:"column:custom_product_snapshot;not null;type:jsonb" json:"custom_product_snapshot"`
	State                 string    `gorm:"column:state;not null;size:30" json:"state"`
	UnitPrice             float64   `gorm:"column:unit_price;not null;type:numeric(10,2)" json:"unit_price"`
	Quantity              int       `gorm:"column:quantity;not null" json:"quantity"`
	ShippingCost          float64   `gorm:"column:shipping_cost;not null;type:numeric(10,2)" json:"shipping_cost"`
	TotalAmount           float64   `gorm:"column:total_amount;not null;type:numeric(10,2)" json:"total_amount"`
	HsCode                string    `gorm:"column:hs_code;not null;size:20" json:"hs_code"`
	CountryOfOrigin       string    `gorm:"column:country_of_origin;not null;size:50" json:"country_of_origin"`
	BillingAddressSnapshot  JSONMap  `gorm:"column:billing_address_snapshot;not null;type:jsonb" json:"billing_address_snapshot"`
	DeliveryAddressSnapshot JSONMap  `gorm:"column:delivery_address_snapshot;not null;type:jsonb" json:"delivery_address_snapshot"`
	LiveRoomID            *uint64   `gorm:"column:live_room_id" json:"live_room_id,omitempty"`
	TrackingNumber        string    `gorm:"column:tracking_number;size:60" json:"tracking_number,omitempty"`
	ShippingCarrier       string    `gorm:"column:shipping_carrier;size:50" json:"shipping_carrier,omitempty"`
	ShippedAt             *time.Time `gorm:"column:shipped_at" json:"shipped_at,omitempty"`
	CreatedAt             time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt             *time.Time `gorm:"column:deleted_at" json:"-"`

	User          *User         `gorm:"-:migration;foreignKey:UserID" json:"user,omitempty"`
	Staff         *Staff        `gorm:"-:migration;foreignKey:StaffID" json:"staff,omitempty"`
	CustomProduct *CustomProduct `gorm:"-:migration;foreignKey:CustomProductID" json:"custom_product,omitempty"`
	LiveRoom      *LiveRoom     `gorm:"-:migration;foreignKey:LiveRoomID" json:"live_room,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}
