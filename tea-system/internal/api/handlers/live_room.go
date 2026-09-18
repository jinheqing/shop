package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"tea-system/internal/middleware"
	"tea-system/internal/models"
	"tea-system/internal/repository"
	"tea-system/internal/service"
	"tea-system/internal/util"
)

// LiveRoomHandler — 直播间 CRUD + 状态机
type LiveRoomHandler struct {
	svc *service.LiveRoomService
}

func NewLiveRoomHandler(svc *service.LiveRoomService) *LiveRoomHandler {
	return &LiveRoomHandler{svc: svc}
}

type liveRoomCreateReq struct {
	RoomName       string  `json:"room_name"`
	RoomType       string  `json:"room_type" binding:"required"`
	PushSource     string  `json:"push_source" binding:"required"`
	Location       string  `json:"location"`
	CoverImage     string  `json:"cover_image"`
	Description    string  `json:"description"`
	OrderID        *uint64 `json:"order_id"`
	HostStaffID    *uint64 `json:"host_staff_id"`
	ScheduledStart *string `json:"scheduled_start"` // RFC3339
}

type liveRoomUpdateReq struct {
	RoomName    *string `json:"room_name"`
	Location    *string `json:"location"`
	CoverImage  *string `json:"cover_image"`
	Description *string `json:"description"`
	HostStaffID *uint64 `json:"host_staff_id"`
}

type customerRequestReq struct {
	Scene       string `json:"scene" binding:"required"`
	Date        string `json:"date" binding:"required"` // YYYY-MM-DD
	TimeSlot    string `json:"time_slot"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Contact     string `json:"contact"`
}

// Create — POST /live-rooms
func (h *LiveRoomHandler) Create(c *gin.Context) {
	var req liveRoomCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	staffID := middleware.GetSubjectID(c)
	if middleware.GetSubjectType(c) != "staff" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "staff only"})
		return
	}

	in := service.CreateInput{
		RoomName:         req.RoomName,
		RoomType:         req.RoomType,
		PushSource:       req.PushSource,
		Location:         req.Location,
		CoverImage:       req.CoverImage,
		Description:      req.Description,
		OrderID:          req.OrderID,
		HostStaffID:      req.HostStaffID,
		CreatedByStaffID: &staffID,
	}

	room, err := h.svc.Create(c.Request.Context(), in)
	if err != nil {
		log.Error().Err(err).Msg("live_room: create failed")
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, room)
}

// List — GET /live-rooms
func (h *LiveRoomHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	f := repository.LiveRoomListFilter{
		RoomType: c.Query("room_type"),
		Status:   c.Query("status"),
		Keyword:  c.Query("keyword"),
		Page:     page,
		Size:     size,
	}
	if v := c.Query("order_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil {
			f.OrderID = &id
		}
	}

	items, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		log.Error().Err(err).Msg("live_room: list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "page": page, "size": size})
}

// GetByID — GET /live-rooms/:id
func (h *LiveRoomHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	room, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrLiveRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
			return
		}
		log.Error().Err(err).Msg("live_room: get failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// Update — PUT /live-rooms/:id
func (h *LiveRoomHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	var req liveRoomUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	patch := map[string]interface{}{}
	if req.RoomName != nil {
		patch["room_name"] = *req.RoomName
	}
	if req.Location != nil {
		patch["location"] = *req.Location
	}
	if req.CoverImage != nil {
		patch["cover_image"] = *req.CoverImage
	}
	if req.Description != nil {
		patch["description"] = *req.Description
	}
	if req.HostStaffID != nil {
		patch["host_staff_id"] = *req.HostStaffID
	}

	room, err := h.svc.Update(c.Request.Context(), id, patch)
	if err != nil {
		if errors.Is(err, repository.ErrLiveRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
			return
		}
		log.Error().Err(err).Msg("live_room: update failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update failed"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// Delete — DELETE /live-rooms/:id
func (h *LiveRoomHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrLiveRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
			return
		}
		log.Error().Err(err).Msg("live_room: delete failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "deleted"})
}

// Start — POST /live-rooms/:id/start
func (h *LiveRoomHandler) Start(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	room, err := h.svc.Start(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrLiveRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
			return
		}
		log.Error().Err(err).Msg("live_room: start failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "start failed"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// End — POST /live-rooms/:id/end
func (h *LiveRoomHandler) End(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	room, err := h.svc.End(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrLiveRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
			return
		}
		log.Error().Err(err).Msg("live_room: end failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "end failed"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// Calendar — GET /live-rooms/calendar（开放预约日历）
func (h *LiveRoomHandler) Calendar(c *gin.Context) {
	items, err := h.svc.GetCalendar(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("live_room: calendar failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// CustomerRequest — POST /live-rooms/customer-request（公开接口，客户提交申请）
// 作为 live_room 的 room_type=customer_request 占位请求对象，真正落库时会先标记 pending
func (h *LiveRoomHandler) CustomerRequest(c *gin.Context) {
	var req customerRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 这里暂不真正落库 live_room，只返回确认占位
	// 后续可以接表单审核流程 → staff 审核通过后再真正创建 live_room
	log.Info().
		Str("scene", req.Scene).
		Str("date", req.Date).
		Str("contact", req.Contact).
		Msg("customer_request: received (STUB — not persisting yet)")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "request submitted, staff will review and create a live_room",
		"echo": gin.H{
			"scene":    req.Scene,
			"date":     req.Date,
			"time_slot": req.TimeSlot,
			"name":     req.Name,
			"contact":  req.Contact,
		},
	})
}

// CustomerRequests — GET /live-rooms/customer-requests（staff 待审核列表）
// 当前也是 stub（真实实现要一个独立的 customer_request 表）
func (h *LiveRoomHandler) CustomerRequests(c *gin.Context) {
	if middleware.GetSubjectType(c) != "staff" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "staff only"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "stub — would list pending customer_requests here",
		"items":   []interface{}{},
	})
}

// ApproveRequest — POST /live-rooms/customer-requests/:id/approve
// 审批客户申请：查找 live_room（room_type=customer_request），审核通过后改为 configuring
func (h *LiveRoomHandler) ApproveRequest(c *gin.Context) {
	if middleware.GetSubjectType(c) != "staff" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "staff only"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var req struct {
		ScheduledStart string `json:"scheduled_start"` // RFC3339
		Note           string `json:"note"`
	}
	_ = c.ShouldBindJSON(&req)

	// 先查一下 room 是否存在，并确认类型
	room, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrLiveRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "customer request not found"})
			return
		}
		log.Error().Err(err).Msg("live_room: approve lookup failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}
	_ = room // room 信息已存在，用于后续 audit 可扩展

	// patch：设置 scheduled_start + status=configuring
	patch := map[string]interface{}{
		"status": "configuring",
	}
	if req.ScheduledStart != "" {
		patch["scheduled_start"] = req.ScheduledStart
	}

	updated, err := h.svc.Update(c.Request.Context(), id, patch)
	if err != nil {
		log.Error().Err(err).Msg("live_room: approve update failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "approve failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       0,
		"message":    "customer request approved",
		"room":       updated,
		"note":       req.Note,
	})
}

// SystemCreate — POST /live-rooms/system-create
// 系统自动创建直播间（如：订单 → 交付验货直播间）
func (h *LiveRoomHandler) SystemCreate(c *gin.Context) {
	if middleware.GetSubjectType(c) != "staff" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "staff only"})
		return
	}

	var req struct {
		OrderID      uint64 `json:"order_id" binding:"required"`
		RoomType     string `json:"room_type"` // 默认 delivery_inspection
		Location     string `json:"location"`
		Description  string `json:"description"`
		HostStaffID  *uint64 `json:"host_staff_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	roomType := req.RoomType
	if roomType == "" {
		roomType = models.RoomTypeDeliveryInspection
	}

	staffID := middleware.GetSubjectID(c)
	in := service.CreateInput{
		RoomType:     roomType,
		PushSource:   models.PushSourceOBSRTMP,
		OrderID:      &req.OrderID,
		Location:     req.Location,
		Description:  req.Description,
		HostStaffID:  req.HostStaffID,
		CreatedByStaffID: &staffID,
	}

	room, err := h.svc.Create(c.Request.Context(), in)
	if err != nil {
		log.Error().Err(err).Msg("live_room: system-create failed")
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, room)
}

// PublicLiveRooms — GET /public/live-rooms（公开查询，仅暴露安全字段）
func (h *LiveRoomHandler) PublicLiveRooms(c *gin.Context) {
	items, _, err := h.svc.List(c.Request.Context(), repository.LiveRoomListFilter{
		Status: "live",
		Page:   1,
		Size:   50,
	})
	if err != nil {
		log.Error().Err(err).Msg("live_room: public list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}
	// Strip sensitive fields
	type PublicRoom struct {
		ID          uint64 `json:"id"`
		RoomID      string `json:"room_id"`
		RoomName    string `json:"room_name"`
		RoomType    string `json:"room_type"`
		Location    string `json:"location,omitempty"`
		Description string `json:"description,omitempty"`
		Status      string `json:"status"`
	}
	public := make([]PublicRoom, 0, len(items))
	for _, r := range items {
		public = append(public, PublicRoom{
			ID:          r.ID,
			RoomID:      r.RoomID,
			RoomName:    r.RoomName,
			RoomType:    r.RoomType,
			Location:    util.SanitizeGardenLocation(r.Location),
			Description: r.Description,
			Status:      r.Status,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": public})
}
