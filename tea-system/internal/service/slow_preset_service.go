package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"tea-system/internal/config"
	"tea-system/internal/models"
	"tea-system/internal/repository"
)

// SlowPresetService — 慢直播专用业务逻辑
// 慢直播实际存为 LiveRoom（room_type=slow_preset, push_source=camera_rtmp），
// 本 service 提供独立的接口组，避免和普通 live_rooms 路由混用。
type SlowPresetService struct {
	repo *repository.LiveRoomRepo
	cfg  *config.Config
}

// NewSlowPresetService — 构造
func NewSlowPresetService(repo *repository.LiveRoomRepo, cfg *config.Config) *SlowPresetService {
	return &SlowPresetService{repo: repo, cfg: cfg}
}

// Create — 创建慢直播（自动 room_type=slow_preset, push_source=camera_rtmp）
// 初始状态 configuring。
func (s *SlowPresetService) Create(ctx context.Context, name, location, description string) (*models.LiveRoom, error) {
	if name == "" {
		name = "茶山慢直播"
	}

	in := CreateInput{
		RoomName:   name,
		RoomType:   models.RoomTypeSlowPreset,
		PushSource: models.PushSourceCameraRTMP,
		Location:   location,
		Description: description,
	}

	// 直接复用 LiveRoomService 的逻辑，避免重复代码
	// 但我们需要自己组装 live_room，因为 LiveRoomService 依赖 LiveKitService
	// 所以这里直接调用 Create 分支
	lrSvc := NewLiveRoomService(s.repo, s.cfg, nil)
	room, err := lrSvc.Create(ctx, in)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// GetByID — 按主键查
func (s *SlowPresetService) GetByID(ctx context.Context, id uint64) (*models.LiveRoom, error) {
	room, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if room.RoomType != models.RoomTypeSlowPreset {
		return nil, repository.ErrLiveRoomNotFound
	}
	return room, nil
}

// List — 慢直播列表（status != offline 的）
func (s *SlowPresetService) List(ctx context.Context, enabledOnly bool) ([]models.LiveRoom, error) {
	f := repository.LiveRoomListFilter{
		RoomType: models.RoomTypeSlowPreset,
		Page:     1,
		Size:     200,
	}
	if enabledOnly {
		f.Status = models.LiveStatusLive
	}
	items, _, err := s.repo.List(ctx, f)
	return items, err
}

// PublicList — 公开接口返回的慢直播列表（只返回 enabled 的）
func (s *SlowPresetService) PublicList(ctx context.Context) ([]models.LiveRoom, error) {
	items, err := s.List(ctx, true)
	return items, err
}

// Update — 更新名称/位置/描述
func (s *SlowPresetService) Update(ctx context.Context, id uint64, name, location, description string) (*models.LiveRoom, error) {
	room, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	patch := map[string]interface{}{}
	if name != "" {
		patch["room_name"] = name
	}
	if location != "" {
		patch["location"] = location
	}
	if description != "" {
		patch["description"] = description
	}
	if len(patch) == 0 {
		return room, nil
	}
	if err := s.repo.Update(ctx, id, patch); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

// Delete — 软删除
func (s *SlowPresetService) Delete(ctx context.Context, id uint64) error {
	room, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if room.Status == models.LiveStatusLive {
		return errors.New("please disable before deleting")
	}
	return s.repo.SoftDelete(ctx, id)
}

// Enable — 启用慢直播：configuring → scheduled → live，生成 camera_rtmp_url
func (s *SlowPresetService) Enable(ctx context.Context, id uint64) (*models.LiveRoom, error) {
	room, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 确保 camera_rtmp_url 已生成
	mediamtxHost := s.mediamtxHost()
	newURL := fmt.Sprintf("rtmp://%s:%d/slow/%s", mediamtxHost, s.cfg.MediaMTX.RTMPPort, room.RoomID)

	patch := map[string]interface{}{
		"camera_rtmp_url": newURL,
		"status":          models.LiveStatusLive,
		"started_at":      time.Now(),
	}
	if err := s.repo.Update(ctx, id, patch); err != nil {
		return nil, err
	}

	log.Info().Uint64("id", id).Str("room_id", room.RoomID).Msg("slow_preset: enabled → live")
	return s.repo.GetByID(ctx, id)
}

// Disable — 禁用慢直播：status → offline
func (s *SlowPresetService) Disable(ctx context.Context, id uint64) (*models.LiveRoom, error) {
	room, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	patch := map[string]interface{}{
		"status":    models.LiveStatusOffline,
		"ended_at":  time.Now(),
	}
	if err := s.repo.Update(ctx, id, patch); err != nil {
		return nil, err
	}
	log.Info().Uint64("id", id).Str("room_id", room.RoomID).Msg("slow_preset: disabled → offline")
	return s.repo.GetByID(ctx, id)
}

// mediamtxHost — 复用 LiveRoomService 里的实现
func (s *SlowPresetService) mediamtxHost() string {
	if s.cfg == nil {
		return "localhost"
	}
	if s.cfg.MediaMTX.WebRTCURL != "" {
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
