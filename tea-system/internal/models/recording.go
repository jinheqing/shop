package models

import (
	"time"
)

// Recording — 直播回放（独立实体，一对一：一个直播 = 一个回放）
// visibility 三档同 LiveRoom: public / registered / restricted
type Recording struct {
	ID              uint64     `gorm:"primaryKey;column:id" json:"id"`
	LiveRoomID      *uint64    `gorm:"column:live_room_id;index" json:"live_room_id,omitempty"`
	Title           string     `gorm:"column:title;size:200" json:"title,omitempty"`
	FileURL         string     `gorm:"column:file_url;size:500;not null" json:"file_url"` // /uploads/recordings/xxx.mp4
	DurationSec     int        `gorm:"column:duration_sec" json:"duration_sec"`
	FileSize        int64      `gorm:"column:file_size" json:"file_size"`
	Visibility      string     `gorm:"column:visibility;size:20;default:'registered'" json:"visibility"`
	VisibleUserIDs  JSONArray  `gorm:"column:visible_user_ids;type:jsonb" json:"visible_user_ids,omitempty"`
	VisibleGroupIDs JSONArray  `gorm:"column:visible_group_ids;type:jsonb" json:"visible_group_ids,omitempty"`
	StartedAt       *time.Time `gorm:"column:started_at" json:"started_at,omitempty"`
	EndedAt         *time.Time `gorm:"column:ended_at" json:"ended_at,omitempty"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null" json:"created_at"`

	LiveRoom *LiveRoom `gorm:"-:migration;foreignKey:LiveRoomID" json:"live_room,omitempty"`
}

func (Recording) TableName() string {
	return "recordings"
}
