package models

import (
	"time"
)

// LiveRoom Types
const (
	RoomTypeSlowPreset          = "slow_preset"           // 24/7 茶山慢直播
	RoomTypeObsTasting          = "obs_tasting"           // OBS 品鉴直播
	RoomTypeOpenCalendar        = "open_calendar"         // 开放预约
	RoomTypeCustomerRequest     = "customer_request"      // 客户指定场景
	RoomTypeDeliveryInspection  = "delivery_inspection"   // 交付验货
	RoomTypeCustomPrivate       = "custom_private"        // 定制私密直播
)

// Push Source
const (
	PushSourceCameraRTMP = "camera_rtmp" // 慢直播专用
	PushSourceAppWebRTC  = "app_webrtc"  // APP 原生推流
	PushSourceOBSRTMP    = "obs_rtmp"    // OBS 推流
)

// LiveRoom Status
const (
	LiveStatusConfiguring = "configuring"
	LiveStatusScheduled   = "scheduled"
	LiveStatusLive        = "live"
	LiveStatusOffline     = "offline"
	LiveStatusEnded       = "ended"
)

// LiveRoom — 直播间表（6 种 room_type × 3 种 push_source）
type LiveRoom struct {
	ID                    uint64     `gorm:"primaryKey;column:id" json:"id"`
	RoomID                string     `gorm:"column:room_id;uniqueIndex;not null;size:64" json:"room_id"`
	RoomName              string     `gorm:"column:room_name;not null;size:200" json:"room_name"`
	RoomType              string     `gorm:"column:room_type;not null;size:30" json:"room_type"`
	Location              string     `gorm:"column:location;size:200" json:"location,omitempty"`
	CoverImage            string     `gorm:"column:cover_image;size:500" json:"cover_image,omitempty"`
	Description           string     `gorm:"column:description;type:text" json:"description,omitempty"`
	PushSource            string     `gorm:"column:push_source;not null;size:20" json:"push_source"`
	CameraRTMPURL         string     `gorm:"column:camera_rtmp_url;size:500" json:"camera_rtmp_url,omitempty"`
	ObsRTMPURL            string     `gorm:"column:obs_rtmp_url;size:500" json:"obs_rtmp_url,omitempty"`
	ObsRTMPKey            string     `gorm:"column:obs_rtmp_key;size:100" json:"obs_rtmp_key,omitempty"`
	OrderID               *uint64    `gorm:"column:order_id" json:"order_id,omitempty"`
	HostStaffID           *uint64    `gorm:"column:host_staff_id" json:"host_staff_id,omitempty"`
	LivekitTokenForHost   string     `gorm:"column:livekit_token_for_host;type:text" json:"-"`
	LivekitTokenForViewers string    `gorm:"column:livekit_token_for_viewers;type:text" json:"-"`
	Status                string     `gorm:"column:status;not null;size:20" json:"status"`
	ScheduledStart        *time.Time `gorm:"column:scheduled_start" json:"scheduled_start,omitempty"`
	StartedAt             *time.Time `gorm:"column:started_at" json:"started_at,omitempty"`
	EndedAt               *time.Time `gorm:"column:ended_at" json:"ended_at,omitempty"`
	PeakViewers           int        `gorm:"column:peak_viewers;default:0" json:"peak_viewers"`
	RecordingURL          string     `gorm:"column:recording_url;size:500" json:"recording_url,omitempty"`
	// ===== 2026-09 新增: 可见性 & 录制 & 类型 =====
	Type             string    `gorm:"column:type;size:30;default:'scheduled'" json:"type"` // slow_live / scheduled / advisor / admin
	Visibility       string    `gorm:"column:visibility;size:20;default:'registered'" json:"visibility"` // public / registered / restricted
	VisibleUserIDs   JSONArray `gorm:"column:visible_user_ids;type:jsonb" json:"visible_user_ids,omitempty"`
	VisibleGroupIDs  JSONArray `gorm:"column:visible_group_ids;type:jsonb" json:"visible_group_ids,omitempty"`
	EnableRecording  bool      `gorm:"column:enable_recording;default:true" json:"enable_recording"`
	RecordingID      *uint64   `gorm:"column:recording_id" json:"recording_id,omitempty"`
	// ===== 新增结束 =====
	TranslationSessionID  *uint64    `gorm:"column:translation_session_id" json:"translation_session_id,omitempty"`
	CreatedByStaffID      *uint64    `gorm:"column:created_by_staff_id" json:"created_by_staff_id,omitempty"`
	CreatedAt             time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt             *time.Time `gorm:"column:deleted_at" json:"-"`

	Order               *Order               `gorm:"-:migration;foreignKey:OrderID" json:"order,omitempty"`
	HostStaff           *Staff               `gorm:"-:migration;foreignKey:HostStaffID" json:"host_staff,omitempty"`
	TranslationSession  *TranslationSession  `gorm:"-:migration;foreignKey:TranslationSessionID" json:"translation_session,omitempty"`
}

func (LiveRoom) TableName() string {
	return "live_rooms"
}
