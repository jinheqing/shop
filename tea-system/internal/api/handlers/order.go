package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

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
}

func NewOrderHandler(
	repo *repository.OrderRepo,
	customRepo *repository.CustomProductRepo,
	stateMachine *service.OrderStateMachine,
) *OrderHandler {
	return &OrderHandler{
		repo:         repo,
		customRepo:   customRepo,
		stateMachine: stateMachine,
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
