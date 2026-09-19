package handlers

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"tea-system/internal/config"
	"tea-system/internal/service"

	"github.com/gin-gonic/gin"
)

// extractHost — 从 URL 提取 host 部分，失败返回 fallback
func extractHost(rawURL, fallback string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return fallback
	}
	host := u.Host
	// 去掉端口
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	if host == "" {
		return fallback
	}
	return host
}

// LiveKitTokenHandler — token 签发 + 房间管理
type LiveKitTokenHandler struct {
	svc  *service.LiveKitService
	cfg  *config.Config
}

func NewLiveKitTokenHandler(svc *service.LiveKitService, cfg *config.Config) *LiveKitTokenHandler {
	return &LiveKitTokenHandler{svc: svc, cfg: cfg}
}

// ---------------- Request Bodies ----------------

type tokenRequest struct {
	RoomName string `json:"room_name" binding:"required"`
	Identity string `json:"identity" binding:"required"`
}

type roomCreateRequest struct {
	RoomName string `json:"room_name" binding:"required"`
}

// ---------------- Responses ----------------

type tokenResponse struct {
	Token    string `json:"token"`
	RoomName string `json:"room_name"`
	Expires  int    `json:"expires"`
}

type hostTokenResponse struct {
	Token     string `json:"token"`
	RoomName  string `json:"room_name"`
	Expires   int    `json:"expires"`
	OBSRTMPURL string `json:"obs_rtmp_url"`
	OBSRTMPKey string `json:"obs_rtmp_key"`
	WebURL    string `json:"web_url"`
}

type roomResponse struct {
	RoomName string `json:"room_name"`
	Status   string `json:"status"`
	ByAPI    bool   `json:"by_api,omitempty"`
}

// ---------------- Handlers ----------------

// Token — POST /livekit/token
// 生成观众 token（只读订阅权限）
func (h *LiveKitTokenHandler) Token(c *gin.Context) {
	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, expire, err := h.svc.GenerateToken(req.RoomName, req.Identity, false /*isHost*/)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		Token:    token,
		RoomName: req.RoomName,
		Expires:  expire,
	})
}

// TokenForOBS — POST /livekit/token-for-obs
// 生成 host token（可发布 + 订阅 + 管理），并返回 MediaMTX RTMP 推流地址
func (h *LiveKitTokenHandler) TokenForOBS(c *gin.Context) {
	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, expire, err := h.svc.GenerateToken(req.RoomName, req.Identity, true /*isHost*/)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// RTMP 地址: 优先使用后台配置的 OBS 推流域名
	rtmpPort := h.cfg.MediaMTX.RTMPPort
	if rtmpPort == 0 {
		rtmpPort = 1935
	}
	rtmpHost := "localhost"
	if h.cfg.LiveKit.URL != "" {
		rtmpHost = extractHost(h.cfg.LiveKit.URL, rtmpHost)
	}

	var obsRTMPURL string
	if h.cfg.OBS.PushDomain != "" {
		// 后台配置了推流域名，直接使用
		obsRTMPURL = h.cfg.OBS.PushDomain + "/" + req.RoomName
	} else {
		obsRTMPURL = "rtmp://" + rtmpHost + ":" + strconv.Itoa(rtmpPort) + "/live/" + req.RoomName
	}
	obsRTMPKey := req.RoomName

	var webURL string
	if h.cfg.OBS.PullDomain != "" {
		// 后台配置了拉流域名
		webURL = h.cfg.OBS.PullDomain + "/" + req.RoomName
	} else {
		webURL = service.MediaMTXWebURL(h.cfg.MediaMTX.WebRTCURL, req.RoomName)
	}

	c.JSON(http.StatusOK, hostTokenResponse{
		Token:      token,
		RoomName:   req.RoomName,
		Expires:    expire,
		OBSRTMPURL: obsRTMPURL,
		OBSRTMPKey: obsRTMPKey,
		WebURL:     webURL,
	})
}

// CreateRoom — POST /livekit/rooms
// 若 LiveKit API 可用则真实创建，否则返回 "simulated"
func (h *LiveKitTokenHandler) CreateRoom(c *gin.Context) {
	var req roomCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	byAPI, err := h.svc.CreateRoom(req.RoomName)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	status := "simulated"
	if byAPI {
		status = "created"
	}
	c.JSON(http.StatusOK, roomResponse{
		RoomName: req.RoomName,
		Status:   status,
		ByAPI:    byAPI,
	})
}

// DeleteRoom — DELETE /livekit/rooms/:name
func (h *LiveKitTokenHandler) DeleteRoom(c *gin.Context) {
	roomName := c.Param("name")
	if roomName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room name is required"})
		return
	}

	byAPI, err := h.svc.DeleteRoom(roomName)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	status := "simulated"
	if byAPI {
		status = "deleted"
	}
	c.JSON(http.StatusOK, roomResponse{
		RoomName: roomName,
		Status:   status,
		ByAPI:    byAPI,
	})
}
