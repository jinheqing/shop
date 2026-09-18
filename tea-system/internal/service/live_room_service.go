package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/zerolog/log"

	"tea-system/internal/config"
	"tea-system/internal/models"
	"tea-system/internal/repository"
)

// LiveRoomService — 直播间业务逻辑（CRUD + 按 room_type/push_source 自动派生字段）
type LiveRoomService struct {
	repo    *repository.LiveRoomRepo
	cfg     *config.Config
	lk      *LiveKitService
}

// NewLiveRoomService — 构造
func NewLiveRoomService(repo *repository.LiveRoomRepo, cfg *config.Config, lk *LiveKitService) *LiveRoomService {
	return &LiveRoomService{repo: repo, cfg: cfg, lk: lk}
}

// CreateInput — 创建直播间入参
type CreateInput struct {
	RoomName     string
	RoomType     string
	PushSource   string
	Location     string
	CoverImage   string
	Description  string
	OrderID      *uint64
	HostStaffID  *uint64
	ScheduledStart   *time.Time
	CreatedByStaffID *uint64
	// ===== 2026-09 新增 =====
	Type            string
	Visibility      string
	VisibleUserIDs  models.JSONArray
	VisibleGroupIDs models.JSONArray
	EnableRecording bool
}

// Create — 创建直播间（按 room_type + push_source 自动生成字段）
func (s *LiveRoomService) Create(ctx context.Context, in CreateInput) (*models.LiveRoom, error) {
	if in.RoomType == "" {
		return nil, errors.New("room_type is required")
	}
	if in.PushSource == "" {
		return nil, errors.New("push_source is required")
	}
	if in.RoomName == "" {
		in.RoomName = defaultRoomName(in.RoomType)
	}

	roomID := fmt.Sprintf("%s-%s", in.RoomType, shortuuid.New())

	now := time.Now()
	room := &models.LiveRoom{
		RoomID:       roomID,
		RoomName:     in.RoomName,
		RoomType:     in.RoomType,
		PushSource:   in.PushSource,
		Location:     in.Location,
		CoverImage:   in.CoverImage,
		Description:  in.Description,
		OrderID:      in.OrderID,
		HostStaffID:  in.HostStaffID,
		Status:       models.LiveStatusConfiguring,
		ScheduledStart:   in.ScheduledStart,
		CreatedByStaffID: in.CreatedByStaffID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	// 2026-09: 新字段 (visibility defaults to "registered", enable_recording defaults true)
	if in.Type != "" {
		room.Type = in.Type
	} else {
		room.Type = "scheduled"
	}
	if in.Visibility != "" {
		room.Visibility = in.Visibility
	} else {
		room.Visibility = "registered"
	}
	room.VisibleUserIDs = in.VisibleUserIDs
	room.VisibleGroupIDs = in.VisibleGroupIDs
	room.EnableRecording = in.EnableRecording || true

	mediamtxHost := s.mediamtxHost()

	switch in.PushSource {
	case models.PushSourceOBSRTMP:
		room.ObsRTMPURL = fmt.Sprintf("rtmp://%s:%d/live/%s", mediamtxHost, s.cfg.MediaMTX.RTMPPort, roomID)
		room.ObsRTMPKey = roomID
	case models.PushSourceCameraRTMP:
		room.CameraRTMPURL = fmt.Sprintf("rtmp://%s:%d/slow/%s", mediamtxHost, s.cfg.MediaMTX.RTMPPort, roomID)
	case models.PushSourceAppWebRTC:
		// 生成 host/viewer JWT（不阻塞保存）
		if s.lk != nil {
			if tk, _, err := s.lk.GenerateToken(roomID, "host-"+roomID, true); err == nil {
				room.LivekitTokenForHost = tk
			}
			if tk, _, err := s.lk.GenerateToken(roomID, "viewer-"+roomID, false); err == nil {
				room.LivekitTokenForViewers = tk
			}
		}
	default:
		return nil, fmt.Errorf("unknown push_source: %s", in.PushSource)
	}

	if err := s.repo.Create(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}

// GetByID — 按主键查
func (s *LiveRoomService) GetByID(ctx context.Context, id uint64) (*models.LiveRoom, error) {
	return s.repo.GetByID(ctx, id)
}

// List — 列表
func (s *LiveRoomService) List(ctx context.Context, f repository.LiveRoomListFilter) ([]models.LiveRoom, int64, error) {
	return s.repo.List(ctx, f)
}

// Update — patch 更新（先取再合并零值安全）
func (s *LiveRoomService) Update(ctx context.Context, id uint64, patch map[string]interface{}) (*models.LiveRoom, error) {
	if err := s.repo.Update(ctx, id, patch); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

// Delete — 软删除
func (s *LiveRoomService) Delete(ctx context.Context, id uint64) error {
	return s.repo.SoftDelete(ctx, id)
}

// Start — 开始直播
// 状态机: configuring/scheduled → live
// 副作用: 触发 LiveKitService.CreateRoom + 确保 host token 已生成
func (s *LiveRoomService) Start(ctx context.Context, id uint64) (*models.LiveRoom, error) {
	room, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if room.Status != models.LiveStatusConfiguring && room.Status != models.LiveStatusScheduled {
		return nil, fmt.Errorf("room %d cannot start: current status=%s (expect configuring/scheduled)", id, room.Status)
	}

	// 1) 如果 push_source 是 app_webrtc / obs_rtmp，确保 host viewer token 都已生成
	if room.PushSource == models.PushSourceAppWebRTC && s.lk != nil {
		if room.LivekitTokenForHost == "" {
			if tk, _, e := s.lk.GenerateToken(room.RoomID, "host-"+room.RoomID, true); e == nil {
				room.LivekitTokenForHost = tk
				_ = s.repo.Update(ctx, id, map[string]interface{}{"livekit_token_for_host": tk})
			} else {
				log.Warn().Err(e).Uint64("room_id", id).Msg("generate host token failed, continue without")
			}
		}
	}

	// 2) 触发 LiveKit 真正创建房间（若服务可用）
	if s.lk != nil {
		if _, e := s.lk.CreateRoom(room.RoomID); e != nil {
			log.Warn().Err(e).Str("room", room.RoomID).Msg("livekit create room failed (non-blocking)")
		}
	}

	// 3) 状态流转
	if err := s.repo.Start(ctx, id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

// End — 结束直播
func (s *LiveRoomService) End(ctx context.Context, id uint64) (*models.LiveRoom, error) {
	if err := s.repo.End(ctx, id); err != nil {
		return nil, err
	}
	room, _ := s.repo.GetByID(ctx, id)

	// 触发交付验货自动流转（占位）
	if room != nil && room.RoomType == models.RoomTypeDeliveryInspection && room.OrderID != nil {
		OnLiveRoomEnded(room.RoomID)
	}
	return room, nil
}

// GetCalendar — 开放预约日历（status=scheduled/live 的 open_calendar）
func (s *LiveRoomService) GetCalendar(ctx context.Context) ([]models.LiveRoom, error) {
	items, _, err := s.repo.List(ctx, repository.LiveRoomListFilter{
		RoomType: models.RoomTypeOpenCalendar,
		Page:     1,
		Size:     100,
	})
	return items, err
}

// mediamtxHost — 从 MediaMTXWebRTCURL 推导 host（默认 localhost）
func (s *LiveRoomService) mediamtxHost() string {
	if s.cfg == nil {
		return "localhost"
	}
	// 简单解析 WebRTCURL 或直接用默认
	if s.cfg.MediaMTX.WebRTCURL != "" {
		// 从 http://host:port 提取 host
		for i := 0; i < len(s.cfg.MediaMTX.WebRTCURL); i++ {
			if s.cfg.MediaMTX.WebRTCURL[i] == '/' && i > 6 {
				potential := s.cfg.MediaMTX.WebRTCURL[i+1:]
				for j := 0; j < len(potential); j++ {
					if potential[j] == ':' || potential[j] == '/' {
						return potential[:j]
					}
				}
				if potential != "" {
					return potential
				}
			}
		}
	}
	return "localhost"
}

// defaultRoomName — 按 room_type 给一个默认名称
func defaultRoomName(roomType string) string {
	switch roomType {
	case models.RoomTypeSlowPreset:
		return "茶山慢直播"
	case models.RoomTypeObsTasting:
		return "品鉴直播"
	case models.RoomTypeOpenCalendar:
		return "开放预约"
	case models.RoomTypeCustomerRequest:
		return "客户指定场景"
	case models.RoomTypeDeliveryInspection:
		return "交付验货"
	case models.RoomTypeCustomPrivate:
		return "定制私密直播"
	default:
		return "直播间-" + strconv.Itoa(int(time.Now().Unix()%1000))
	}
}

// logWarn — 包内工具，避免重复引用
func logWarn(msg string) { log.Warn().Msg(msg) }
