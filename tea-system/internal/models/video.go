package models

import "time"

// VideoCategory — 视频分类（后台配置，前台展示）
type VideoCategory struct {
	ID          uint64    `gorm:"primaryKey;column:id" json:"id"`
	Name        string    `gorm:"column:name;size:100;not null" json:"name"`
	Slug        string    `gorm:"column:slug;size:100;uniqueIndex" json:"slug"`
	Description string    `gorm:"column:description;size:500" json:"description,omitempty"`
	SortOrder   int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (VideoCategory) TableName() string { return "video_categories" }

// Video — 独立视频实体（非直播回放，后台手动发布）
type Video struct {
	ID          uint64     `gorm:"primaryKey;column:id" json:"id"`
	CategoryID  *uint64    `gorm:"column:category_id;index" json:"category_id,omitempty"`
	Title       string     `gorm:"column:title;size:200;not null" json:"title"`
	Description string     `gorm:"column:description;type:text" json:"description,omitempty"`
	VideoURL    string     `gorm:"column:video_url;size:500;not null" json:"video_url"`
	CoverURL    string     `gorm:"column:cover_url;size:500" json:"cover_url,omitempty"`
	DurationSec int        `gorm:"column:duration_sec" json:"duration_sec"`
	Status      string     `gorm:"column:status;size:20;default:'draft'" json:"status"` // draft / published
	SortOrder   int        `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`

	Category *VideoCategory `gorm:"-:migration;foreignKey:CategoryID" json:"category,omitempty"`
}

func (Video) TableName() string { return "videos" }
