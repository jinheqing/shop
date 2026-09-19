package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"tea-system/internal/middleware"
	"tea-system/internal/models"
	"tea-system/internal/repository"
	"tea-system/internal/service"
)

// ============================================================
// OrderHandler
// ============================================================

type OrderHandler struct {
	repo         *repository.OrderRepo
	customRepo   *repository.CustomProductRepo
	stateMachine *service.OrderStateMachine
	auditDB      *gorm.DB
}

func NewOrderHandler(
	repo *repository.OrderRepo,
	customRepo *repository.CustomProductRepo,
	stateMachine *service.OrderStateMachine,
	auditDB *gorm.DB,
) *OrderHandler {
	return &OrderHandler{
		repo:         repo,
		customRepo:   customRepo,
		stateMachine: stateMachine,
		auditDB:      auditDB,
	}
}

// ==================== DTO ====================

// OrderCreateRequest — POST /orders
type OrderCreateRequest struct {
	CustomProductID         uint64          `json:"custom_product_id" binding:"required"`
	UnitPrice               float64         `json:"unit_price" binding:"required,gt=0"`
	Quantity                int             `json:"quantity" binding:"required,gt=0"`
	ShippingCost            float64         `json:"shipping_cost" binding:"required,gte=0"`
	BillingAddress          json.RawMessage `json:"billing_address" binding:"required"`
	DeliveryAddress         json.RawMessage `json:"delivery_address" binding:"required"`
	HsCode                  string          `json:"hs_code"`
	CountryOfOrigin         string          `json:"country_of_origin"`
}

// OrderStateRequest — POST /orders/:id/state
type OrderStateRequest struct {
	TargetState string `json:"target_state" binding:"required"`
	Reason      string `json:"reason,omitempty"` // 顾问填写的备注（物流号、集装箱号等）

	// 可选的物流信息（在 shipped 状态转换时填入）
	Courier    string `json:"courier,omitempty"`    // "FedEx" / "DHL" / "EMS" / "Private"
	TrackingNo string `json:"tracking_no,omitempty"` // FedEx 追踪号 / 集装箱号 / AWB 号
	ShippedAt  string `json:"shipped_at,omitempty"`  // RFC3339 时间字符串
	EtaAt      string `json:"eta_at,omitempty"`      // RFC3339 时间字符串
}

// ==================== handlers ====================

// Create — POST /orders
func (h *OrderHandler) Create(c *gin.Context) {
	var req OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 查 CustomProduct 并做快照
	cp, err := h.customRepo.GetByID(ctx, req.CustomProductID)
	if err != nil {
		if errors.Is(err, repository.ErrCustomProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "custom product not found"})
			return
		}
		log.Error().Err(err).Msg("order: custom_product lookup failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "custom product lookup failed"})
		return
	}

	// 快照：把 CustomProduct 的 JSON map 序列化再反序列化到 models.JSONMap
	snapshot, err := toJSONMap(cp)
	if err != nil {
		log.Error().Err(err).Msg("order: snapshot failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to snapshot custom product"})
		return
	}

	billingSnap, err := rawToJSONMap(req.BillingAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid billing_address json"})
		return
	}
	deliverySnap, err := rawToJSONMap(req.DeliveryAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid delivery_address json"})
		return
	}

	userID := middleware.GetSubjectID(c)
	isUser := middleware.GetSubjectType(c) == "user"
	if !isUser {
		// staff 代用户下单（通过 query param user_id）
		if qid, ok := c.GetQuery("user_id"); ok {
			if parsed, perr := strconv.ParseUint(qid, 10, 64); perr == nil {
				userID = parsed
			}
		}
	}

	var userIDPtr *uint64
	if isUser {
		uid := userID
		userIDPtr = &uid
	} else if qid, ok := c.GetQuery("user_id"); ok {
		// staff 明确代指定 user 下单
		uid, _ := strconv.ParseUint(qid, 10, 64)
		if uid != 0 {
			userIDPtr = &uid
		}
	}

	staffID := uint64(0)
	if middleware.GetSubjectType(c) == "staff" {
		staffID = middleware.GetSubjectID(c)
	}

	hsCode := req.HsCode
	if hsCode == "" {
		hsCode = "0902.10"
	}
	country := req.CountryOfOrigin
	if country == "" {
		country = "China"
	}

	total := req.UnitPrice*float64(req.Quantity) + req.ShippingCost

	o := &models.Order{
		OrderNo:                 "", // repo 自动生成
		UserID:                  userIDPtr,
		StaffID:                 staffID,
		CustomProductID:         req.CustomProductID,
		CustomProductSnapshot:   snapshot,
		State:                   models.OrderStateOrdering,
		UnitPrice:               req.UnitPrice,
		Quantity:                req.Quantity,
		ShippingCost:            req.ShippingCost,
		TotalAmount:             total,
		HsCode:                  hsCode,
		CountryOfOrigin:         country,
		BillingAddressSnapshot:  billingSnap,
		DeliveryAddressSnapshot: deliverySnap,
	}

	if err := h.repo.Create(ctx, o); err != nil {
		log.Error().Err(err).Msg("order: create failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to create order"})
		return
	}
	c.JSON(http.StatusCreated, o)
}

// List — GET /orders  支持 ?state= & ?user_id= & 分页
func (h *OrderHandler) List(c *gin.Context) {
	state := c.Query("state")
	userID := uint64(0)
	if qid := c.Query("user_id"); qid != "" {
		if parsed, err := strconv.ParseUint(qid, 10, 64); err == nil {
			userID = parsed
		}
	}

	// 如果是 user JWT，强制过滤自己的订单
	subType := middleware.GetSubjectType(c)
	subID := middleware.GetSubjectID(c)
	if subType == "user" {
		userID = subID
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	items, total, err := h.repo.List(c.Request.Context(), repository.OrderListFilter{
		UserID: userID,
		State:  state,
		Page:   page,
		Size:   size,
	})
	if err != nil {
		log.Error().Err(err).Msg("order: list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetByID — GET /orders/:id
func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	o, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "order not found"})
			return
		}
		log.Error().Err(err).Msg("order: get failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	// 用户只能看自己的订单
	if middleware.GetSubjectType(c) == "user" && (o.UserID == nil || *o.UserID != middleware.GetSubjectID(c)) {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "forbidden"})
		return
	}

	c.JSON(http.StatusOK, o)
}

// UpdateState — POST /orders/:id/state
func (h *OrderHandler) UpdateState(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var req OrderStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	o, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "order not found"})
			return
		}
		log.Error().Err(err).Msg("order: get failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	// 状态机校验
	if err := h.stateMachine.CheckTransition(o.State, req.TargetState); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"from":    o.State,
			"to":      req.TargetState,
		})
		return
	}

	// 从 JWT 取 staffID
	staffID := uint64(0)
	if middleware.GetSubjectType(c) == "staff" {
		staffID = middleware.GetSubjectID(c)
	}
	if staffID == 0 {
		staffID = o.StaffID // fallback：用订单原来的顾问
	}

	// 构建 shippingPatch（只有在 shipped 状态转换时才填）
	shippingPatch := map[string]interface{}{}
	if req.TargetState == models.OrderStateShipped {
		if req.Courier != "" {
			shippingPatch["courier"] = req.Courier
		}
		if req.TrackingNo != "" {
			shippingPatch["tracking_no"] = req.TrackingNo
		}
		if req.ShippedAt != "" {
			if t, err := time.Parse(time.RFC3339, req.ShippedAt); err == nil {
				shippingPatch["shipped_at"] = t
			}
		}
		if req.EtaAt != "" {
			if t, err := time.Parse(time.RFC3339, req.EtaAt); err == nil {
				shippingPatch["eta_at"] = t
			}
		}
	}

	// 事务里写 state_log + 更新订单
	if err := h.repo.LogStateChange(c.Request.Context(), id, o.State, req.TargetState, req.Reason, staffID, shippingPatch); err != nil {
		log.Error().Err(err).Msg("order: log state change failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update state failed"})
		return
	}

	// 返回更新后的订单
	updated, _ := h.repo.GetByID(c.Request.Context(), id)
	if updated == nil {
		c.JSON(http.StatusOK, gin.H{"id": id, "state": req.TargetState})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	o, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "order not found"})
			return
		}
		log.Error().Err(err).Msg("order: get failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	// 如果已经终态 cancelled 或 completed，不允许再 cancel
	if h.stateMachine.IsTerminal(o.State) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "order already in terminal state: " + o.State,
		})
		return
	}

	if err := h.repo.UpdateState(c.Request.Context(), id, models.OrderStateCancelled); err != nil {
		log.Error().Err(err).Msg("order: cancel failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "cancel failed"})
		return
	}
	o.State = models.OrderStateCancelled
	c.JSON(http.StatusOK, o)
}

// ==================== 工具 ====================

// toJSONMap — 把任何 struct 转成 models.JSONMap
func toJSONMap(v interface{}) (models.JSONMap, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m models.JSONMap
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// rawToJSONMap — 把 json.RawMessage 转 models.JSONMap
func rawToJSONMap(raw json.RawMessage) (models.JSONMap, error) {
	if len(raw) == 0 {
		return models.JSONMap{}, nil
	}
	var m models.JSONMap
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Timeline — GET /orders/:id/timeline
// 聚合订单状态变化（state_logs）、审计日志、物流信息等事件，按时间排序返回时间线
func (h *OrderHandler) Timeline(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	o, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "order not found"})
			return
		}
		log.Error().Err(err).Msg("order: timeline get failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	type TimelineEvent struct {
		At     string                 `json:"at"`
		Type   string                 `json:"type"`   // order_created / state_change / audit / order_updated / shipping
		Detail map[string]interface{} `json:"detail"`
	}
	type rawEvent struct {
		at    time.Time
		event TimelineEvent
	}
	var raw []rawEvent

	// 事件 1 — 订单创建
	raw = append(raw, rawEvent{at: o.CreatedAt, event: TimelineEvent{
		At:   o.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		Type: "order_created",
		Detail: map[string]interface{}{
			"order_no":    o.OrderNo,
			"initial_state": o.State,
			"staff_id":    o.StaffID,
		},
	}})

	// 事件 2 — 状态流转日志（核心！顾问每次改状态都会写一条）
	stateLogs, _ := h.repo.GetStateLogs(c.Request.Context(), id)
	for _, sl := range stateLogs {
		detail := map[string]interface{}{
			"from_state": sl.FromState,
			"to_state":   sl.ToState,
			"staff_id":   sl.StaffID,
		}
		if sl.Reason != "" {
			detail["reason"] = sl.Reason
		}
		raw = append(raw, rawEvent{at: sl.CreatedAt, event: TimelineEvent{
			At:     sl.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
			Type:   "state_change",
			Detail: detail,
		}})
	}

	// 事件 3 — 物流发出（如果有 shipped_at 且状态流转里没有体现）
	if o.ShippedAt != nil {
		shippingDetail := map[string]interface{}{
			"shipped_at": o.ShippedAt.UTC().Format("2006-01-02T15:04:05Z"),
		}
		if o.Courier != "" {
			shippingDetail["courier"] = o.Courier
		}
		if o.TrackingNo != "" {
			shippingDetail["tracking_no"] = o.TrackingNo
		}
		if o.EtaAt != nil {
			shippingDetail["eta_at"] = o.EtaAt.UTC().Format("2006-01-02T15:04:05Z")
		}
		raw = append(raw, rawEvent{at: *o.ShippedAt, event: TimelineEvent{
			At:     o.ShippedAt.UTC().Format("2006-01-02T15:04:05Z"),
			Type:   "shipping",
			Detail: shippingDetail,
		}})
	}

	// 事件 4 — 审计日志
	if h.auditDB != nil {
		var logs []models.AuditLog
		ptr := id
		h.auditDB.Where("target_type = ? AND target_id = ?", "order", ptr).
			Order("created_at ASC").
			Limit(100).
			Find(&logs)
		for _, l := range logs {
			detail := map[string]interface{}{
				"action":   l.Action,
				"staff_id": l.StaffID,
			}
			if l.Detail != nil {
				detail["extra"] = l.Detail
			}
			raw = append(raw, rawEvent{at: l.CreatedAt, event: TimelineEvent{
				At:     l.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
				Type:   "audit",
				Detail: detail,
			}})
		}
	}

	// 事件 5 — 订单最后更新（兜底）
	if o.UpdatedAt.After(o.CreatedAt) {
		raw = append(raw, rawEvent{at: o.UpdatedAt, event: TimelineEvent{
			At:   o.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
			Type: "order_updated",
			Detail: map[string]interface{}{
				"state": o.State,
			},
		}})
	}

	// 按时间升序排序
	for i := 1; i < len(raw); i++ {
		key := raw[i]
		j := i - 1
		for j >= 0 && raw[j].at.After(key.at) {
			raw[j+1] = raw[j]
			j--
		}
		raw[j+1] = key
	}

	events := make([]TimelineEvent, len(raw))
	for i, r := range raw {
		events[i] = r.event
	}

	// 附带物流快照
	shipping := map[string]interface{}{}
	if o.Courier != "" {
		shipping["courier"] = o.Courier
	}
	if o.TrackingNo != "" {
		shipping["tracking_no"] = o.TrackingNo
	}
	if o.ShippedAt != nil {
		shipping["shipped_at"] = o.ShippedAt.UTC().Format("2006-01-02T15:04:05Z")
	}
	if o.EtaAt != nil {
		shipping["eta_at"] = o.EtaAt.UTC().Format("2006-01-02T15:04:05Z")
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":  id,
		"order_no":  o.OrderNo,
		"state":     o.State,
		"events":    events,
		"event_cnt": len(events),
		"shipping":  shipping,
	})
}

// UpdateShipping — PUT /orders/:id/shipping
// 独立更新物流信息（不改变订单状态，顾问可在 shipped 之后补填 courier/追踪号）
func (h *OrderHandler) UpdateShipping(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var req struct {
		Courier    string `json:"courier,omitempty"`
		TrackingNo string `json:"tracking_no,omitempty"`
		ShippedAt  string `json:"shipped_at,omitempty"`
		EtaAt      string `json:"eta_at,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var shippedAt, etaAt *time.Time
	if req.ShippedAt != "" {
		if t, perr := time.Parse(time.RFC3339, req.ShippedAt); perr == nil {
			shippedAt = &t
		}
	}
	if req.EtaAt != "" {
		if t, perr := time.Parse(time.RFC3339, req.EtaAt); perr == nil {
			etaAt = &t
		}
	}

	if err := h.repo.UpdateShippingInfo(c.Request.Context(), id, req.Courier, req.TrackingNo, shippedAt, etaAt); err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "order not found"})
			return
		}
		log.Error().Err(err).Msg("order: update shipping failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update shipping failed"})
		return
	}

	o, _ := h.repo.GetByID(c.Request.Context(), id)
	if o != nil {
		c.JSON(http.StatusOK, o)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// StateLogs — GET /orders/:id/state-logs
// 返回订单所有状态流转记录
func (h *OrderHandler) StateLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	logs, err := h.repo.GetStateLogs(c.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("order: state logs failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": logs, "cnt": len(logs)})
}
