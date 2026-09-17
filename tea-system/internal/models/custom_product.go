package models

import (
	"time"
)

// CustomProductStatus
const (
	CustomProductStatusDraft     = "draft"
	CustomProductStatusPublished = "published"
	CustomProductStatusArchived  = "archived"
)

// CustomProduct — 定制报价商品表（完整 26+ 字段，含溯源）
type CustomProduct struct {
	ID                  uint64    `gorm:"primaryKey;column:id" json:"id"`
	ProductToken        *string   `gorm:"column:product_token;uniqueIndex;size:64" json:"product_token,omitempty"`
	Version             int       `gorm:"column:version;default:1" json:"version"`
	IsBespoke           bool      `gorm:"column:is_bespoke;default:true" json:"is_bespoke"`
	NonRefundable       bool      `gorm:"column:non_refundable;default:true" json:"non_refundable"`
	Title               string    `gorm:"column:title;not null;size:200" json:"title"`
	RawTeaSource        string    `gorm:"column:raw_tea_source;not null;type:text" json:"raw_tea_source"`
	CustomRequirement   string    `gorm:"column:custom_requirement;not null;type:text" json:"custom_requirement"`
	TeaType             string    `gorm:"column:tea_type;not null;size:20" json:"tea_type"` // raw_puer / ripe_puer / ancient_tree / vintage
	TeaShape            string    `gorm:"column:tea_shape;not null;size:20" json:"tea_shape"` // loose / cake / brick / tuo
	TeaShapeWeight      *int      `gorm:"column:tea_shape_weight" json:"tea_shape_weight,omitempty"`
	SmokedWithFlower    bool      `gorm:"column:smoked_with_flower;default:false" json:"smoked_with_flower"`
	FlowerType          string    `gorm:"column:flower_type;size:30" json:"flower_type,omitempty"` // jasmine / osmanthus / orchid / custom
	InnerPackaging      string    `gorm:"column:inner_packaging;not null;size:100" json:"inner_packaging"`
	OuterPackaging      string    `gorm:"column:outer_packaging;not null;size:100" json:"outer_packaging"`
	ProductCardText     string    `gorm:"column:product_card_text;type:text" json:"product_card_text,omitempty"`
	ProductCardFormat   string    `gorm:"column:product_card_format;size:20;default:'vertical'" json:"product_card_format"` // vertical / horizontal
	QrCodePosition      string    `gorm:"column:qr_code_position;not null;size:20" json:"qr_code_position"` // outer_front / outer_back / inner_front / hidden
	SKU                 string    `gorm:"column:sku;uniqueIndex;size:50" json:"sku,omitempty"`
	UnitPrice           float64   `gorm:"column:unit_price;not null;type:numeric(10,2)" json:"unit_price"`
	Quantity            int       `gorm:"column:quantity;not null" json:"quantity"`
	ShippingCost        float64   `gorm:"column:shipping_cost;not null;type:numeric(10,2)" json:"shipping_cost"`
	TotalAmount         float64   `gorm:"column:total_amount;not null;type:numeric(10,2)" json:"total_amount"`
	LeadTime            string    `gorm:"column:lead_time;not null;size:100" json:"lead_time"`
	Status              string    `gorm:"column:status;not null;size:20" json:"status"` // draft / published / archived
	CreatedByStaffID    *uint64   `gorm:"column:created_by_staff_id" json:"created_by_staff_id,omitempty"`
	ReviewedByStaffID   *uint64   `gorm:"column:reviewed_by_staff_id" json:"reviewed_by_staff_id,omitempty"`
	ReviewedAt          *time.Time `gorm:"column:reviewed_at" json:"reviewed_at,omitempty"`

	// 溯源字段
	HarvestDate     time.Time `gorm:"column:harvest_date;not null;type:date" json:"harvest_date"`
	RoastingDate    time.Time `gorm:"column:roasting_date;not null;type:date" json:"roasting_date"`
	MountainLocation string   `gorm:"column:mountain_location;not null;size:200" json:"mountain_location"`
	MasterName      string    `gorm:"column:master_name;not null;size:100" json:"master_name"`
	StorageLocation string    `gorm:"column:storage_location;not null;size:200" json:"storage_location"`
	SgsReportID     *uint64   `gorm:"column:sgs_report_id" json:"sgs_report_id,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`

	CreatedByStaff *Staff     `gorm:"foreignKey:CreatedByStaffID" json:"created_by_staff,omitempty"`
	ReviewedByStaff *Staff    `gorm:"foreignKey:ReviewedByStaffID" json:"reviewed_by_staff,omitempty"`
	SgsReport      *SgsReport `gorm:"foreignKey:SgsReportID" json:"sgs_report,omitempty"`
}

func (CustomProduct) TableName() string {
	return "custom_products"
}
