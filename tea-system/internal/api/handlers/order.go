package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
// 修复 (2026-09-19): 最小化 required，admin 手动创建订单时允许只填 custom_product_id + quantity
type OrderCreateRequest struct {
	CustomProductID uint64          `json:"custom_product_id" binding:"required"`
	UnitPrice       *float64        `json:"unit_price,omitempty"`
	Quantity        *int            `json:"quantity,omitempty"`
	ShippingCost    *float64        `json:"shipping_cost,omitempty"`
	BillingAddress  json.RawMessage `json:"billing_address,omitempty"`
	DeliveryAddress json.RawMessage `json:"delivery_address,omitempty"`
	HsCode          string          `json:"hs_code,omitempty"`
	CountryOfOrigin string          `json:"country_of_origin,omitempty"`
}

// OrderStateRequest — POST /orders/:id/state
type OrderStateRequest struct {
	TargetState string `json:"target_state" binding:"required"`
	Reason      string `json:"reason,omitempty"`
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

	snapshot, err := toJSONMap(cp)
	if err != nil {
		log.Error().Err(err).Msg("order: snapshot failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to snapshot custom product"})
		return
	}

	// 地址 — 可选，空时用默认 JSON
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
		uid, _ := strconv.ParseUint(qid, 10, 64)
		if uid != 0 {
			userIDPtr = &uid
		}
	}

	staffID := uint64(0)
	if middleware.GetSubjectType(c) == "staff" {
		staffID = middleware.GetSubjectID(c)
	}

	// 从 CustomProduct 快照填默认值（admin 手动创建订单时前端可能不填）
	unitPrice := cp.UnitPrice
	if req.UnitPrice != nil && *req.UnitPrice > 0 {
		unitPrice = *req.UnitPrice
	}
	quantity := cp.Quantity
	if req.Quantity != nil && *req.Quantity > 0 {
		quantity = *req.Quantity
	}
	shippingCost := cp.ShippingCost
	if req.ShippingCost != nil {
		shippingCost = *req.ShippingCost
	}

	hsCode := req.HsCode
	if hsCode == "" { hsCode = "0902.10" }
	country := req.CountryOfOrigin
	if country == "" { country = "China" }

	total := unitPrice*float64(quantity) + shippingCost

	o := &models.Order{
		OrderNo:                 "",
		UserID:                  userIDPtr,
		StaffID:                 staffID,
		CustomProductID:         req.CustomProductID,
		CustomProductSnapshot:   snapshot,
		State:                   models.OrderStateOrdering,
		UnitPrice:               unitPrice,
		Quantity:                quantity,
		ShippingCost:            shippingCost,
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

	if err := h.repo.UpdateState(c.Request.Context(), id, req.TargetState); err != nil {
		log.Error().Err(err).Msg("order: update state failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update state failed"})
		return
	}

	o.State = req.TargetState
	c.JSON(http.StatusOK, o)
}

// Cancel — POST /orders/:id/cancel
// 语义：任何可取消状态 → cancelled（终态），内部复用状态机 transition 校验
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
// 聚合订单状态变化、审计日志等事件，按时间排序返回时间线
func (h *OrderHandler) Timeline(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	// 1. 订单本身（必须存在）
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

	// 2. 构建时间线事件
	type TimelineEvent struct {
		At      string                 `json:"at"`
		Type    string                 `json:"type"`     // order_state_change / audit / payment / declaration / live_room
		Detail  map[string]interface{} `json:"detail"`
	}

	var events []TimelineEvent

	// 事件 A — 订单创建
	events = append(events, TimelineEvent{
		At:   o.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		Type: "order_created",
		Detail: map[string]interface{}{
			"order_no": o.OrderNo,
			"initial_state": o.State,
			"staff_id": o.StaffID,
		},
	})

	// 事件 B — 审计日志（从独立 audit 库查）
	if h.auditDB != nil {
		var logs []models.AuditLog
		ptr := id
		h.auditDB.Where("target_type = ? AND target_id = ?", "order", ptr).
			Order("created_at ASC").
			Limit(100).
			Find(&logs)
		for _, l := range logs {
			detail := map[string]interface{}{
				"action":  l.Action,
				"staff_id": l.StaffID,
			}
			if l.Detail != nil {
				detail["extra"] = l.Detail
			}
			events = append(events, TimelineEvent{
				At:     l.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
				Type:   "audit",
				Detail: detail,
			})
		}
	}

	// 事件 C — 订单最后更新
	if o.UpdatedAt.After(o.CreatedAt) {
		events = append(events, TimelineEvent{
			At:   o.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
			Type: "order_updated",
			Detail: map[string]interface{}{
				"state": o.State,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":  id,
		"order_no":  o.OrderNo,
		"state":     o.State,
		"events":    events,
		"event_cnt": len(events),
	})
}
