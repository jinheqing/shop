package service

import (
	"fmt"
	"net/url"
	"time"

	"github.com/livekit/protocol/auth"
	"github.com/rs/zerolog/log"
)

// LiveKitService — LiveKit token 签发 + 房间管理
//
// 核心设计:
// 1) token 生成仅依赖 livekit/protocol/auth，不依赖 LiveKit 服务本身，
//    即使 LiveKit Server 没有启动也能返回合法 JWT，方便前端离线联调。
// 2) 房间管理（CreateRoom / DeleteRoom）为 no-op stub，返回 simulated，
//    以便后续真正接入 server-sdk-go RoomService API 时不破坏调用方接口。
type LiveKitService struct {
	url       string
	apiKey    string
	apiSecret string
}

// NewLiveKitService — 构造 service
func NewLiveKitService(url, apiKey, apiSecret string) *LiveKitService {
	if url == "" || apiKey == "" || apiSecret == "" {
		log.Warn().Msg("LiveKitService: url/apiKey/apiSecret 不全，token 生成仍可用，但房间 API 为 simulated")
	}
	return &LiveKitService{
		url:       url,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

// URL — LiveKit server URL
func (s *LiveKitService) URL() string { return s.url }

// APIKey — LiveKit API Key
func (s *LiveKitService) APIKey() string { return s.apiKey }

// IsConnected — 占位，始终返回 false（真正连 LiveKit RoomService 时再实现）
func (s *LiveKitService) IsConnected() bool { return false }

// GenerateToken — 生成 LiveKit JWT
// isHost=true:  可发布 + 订阅 + 发布 data + 管理房间
// isHost=false: 仅订阅
// token 有效期默认 1 小时（3600s）
func (s *LiveKitService) GenerateToken(roomName, identity string, isHost bool) (string, int, error) {
	if roomName == "" {
		return "", 0, fmt.Errorf("room_name is required")
	}
	if identity == "" {
		return "", 0, fmt.Errorf("identity is required")
	}

	at := auth.NewAccessToken(s.apiKey, s.apiSecret)

	grant := &auth.VideoGrant{
		Room:     roomName,
		RoomJoin: true,
	}

	if isHost {
		grant.SetCanPublish(true)
		grant.SetCanSubscribe(true)
		grant.SetCanPublishData(true)
		grant.RoomAdmin = true
	} else {
		grant.SetCanSubscribe(true)
	}

	at.SetVideoGrant(grant)
	at.SetIdentity(identity)

	expireSeconds := 3600
	at.SetValidFor(time.Duration(expireSeconds) * time.Second)

	token, err := at.ToJWT()
	if err != nil {
		return "", 0, fmt.Errorf("livekit token sign: %w", err)
	}
	return token, expireSeconds, nil
}

// CreateRoom — 占位 stub，始终返回 (false, nil) 表示 simulated
// 后续接入 server-sdk-go 的 RoomService.CreateRoom 时再替换真实实现
func (s *LiveKitService) CreateRoom(roomName string) (bool, error) {
	log.Debug().Str("room", roomName).Msg("LiveKitService.CreateRoom: simulated (no LiveKit API client wired)")
	return false, nil
}

// DeleteRoom — 占位 stub，始终返回 (false, nil) 表示 simulated
func (s *LiveKitService) DeleteRoom(roomName string) (bool, error) {
	log.Debug().Str("room", roomName).Msg("LiveKitService.DeleteRoom: simulated")
	return false, nil
}

// MediaMTXWebURL — 根据 baseWebURL 推导 RTSPtoWeb 访问地址
// 例: http://localhost:8889 → http://localhost:8889/{roomName}
func MediaMTXWebURL(baseWebURL, roomName string) string {
	u, err := url.Parse(baseWebURL)
	if err != nil || u.Host == "" {
		return fmt.Sprintf("%s/%s", baseWebURL, roomName)
	}
	u.Path = "/" + roomName
	return u.String()
}
