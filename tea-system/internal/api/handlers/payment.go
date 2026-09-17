package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"tea-system/internal/models"
	"tea-system/internal/repository"
	"tea-system/internal/service"
)

// ============================================================
// PaymentHandler
// ============================================================

type PaymentHandler struct {
	repo         *repository.OrderRepo
	paymentSvc   *service.PaymentService
	stateMachine *service.OrderStateMachine
}

func NewPaymentHandler(
	repo *repository.OrderRepo,
	paymentSvc *service.PaymentService,
	stateMachine *service.OrderStateMachine,
) *PaymentHandler {
	return &PaymentHandler{
		repo:         repo,
		paymentSvc:   paymentSvc,
		stateMachine: stateMachine,
	}
}

// ==================== DTO ====================

// PaymentInitRequest — POST /orders/:id/payment/init
type PaymentInitRequest struct {
	Gateway  string `json:"gateway" binding:"required,oneof=2checkout paypal"`
	ReturnURL string `json:"return_url"`
	Currency string `json:"currency"` // 默认 GBP
}

// ==================== handlers ====================

// Init — POST /orders/:id/payment/init
func (h *PaymentHandler) Init(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid order id"})
		return
	}

	var req PaymentInitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	order, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "order not found"})
			return
		}
		log.Error().Err(err).Msg("payment: order lookup failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "GBP"
	}
	returnURL := req.ReturnURL
	if returnURL == "" {
		returnURL = "/"
	}

	gw := service.PaymentGateway(req.Gateway)
	result, err := h.paymentSvc.InitPayment(gw, service.InitPaymentRequest{
		OrderNo:   order.OrderNo,
		Amount:    order.TotalAmount,
		Currency:  currency,
		ReturnURL: returnURL,
	})
	if err != nil {
		log.Error().Err(err).Msg("payment: init failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "payment init failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"checkout_url": result.CheckoutURL,
		"gateway":      result.Gateway,
		"order_no":     order.OrderNo,
		"amount":       order.TotalAmount,
		"currency":     currency,
	})
}

// Handle2CheckoutWebhook — POST /webhooks/2checkout (公开，无 JWT)
func (h *PaymentHandler) Handle2CheckoutWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "failed to read body"})
		return
	}

	cb, err := h.paymentSvc.Handle2CheckoutCallback(body)
	if err != nil {
		log.Warn().Err(err).Msg("payment: 2co callback parse failed")
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	h.processCallback(c, cb, models.Gateway2Checkout)
}

// HandlePayPalWebhook — POST /webhooks/paypal (公开，无 JWT)
func (h *PaymentHandler) HandlePayPalWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "failed to read body"})
		return
	}

	cb, err := h.paymentSvc.HandlePayPalCallback(body)
	if err != nil {
		log.Warn().Err(err).Msg("payment: paypal callback parse failed")
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	h.processCallback(c, cb, models.GatewayPayPal)
}

// processCallback — 回调公共处理：幂等存交易 + 状态流转
func (h *PaymentHandler) processCallback(c *gin.Context, cb *service.CallbackResult, gateway string) {
	ctx := c.Request.Context()

	// 1. 按订单号查订单（callback.OrderNo 是我们传过去的 merchant ref）
	var order *models.Order
	var err error
	if cb.OrderNo != "" {
		order, err = h.repo.GetByOrderNo(ctx, cb.OrderNo)
	}
	if order == nil {
		// 2. 幂等：gateway_transaction_id 已存在就直接返回
		existing, lookErr := h.repo.GetPaymentByGatewayID(ctx, cb.GatewayTransactionID)
		if lookErr == nil {
			log.Info().Str("tx", cb.GatewayTransactionID).Msg("payment: idempotent duplicate (order lookup failed but tx exists)")
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "already processed",
				"tx":      existing.GatewayTransactionID,
				"status":  existing.Status,
			})
			return
		}
		if order == nil {
			log.Error().Err(err).Str("gateway", gateway).Str("tx", cb.GatewayTransactionID).Str("order_no", cb.OrderNo).Msg("payment: order lookup failed")
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "order not found for tx"})
			return
		}
	}

	// 3. 幂等：tx 已存在 → 直接返回 200
	existing, err := h.repo.GetPaymentByGatewayID(ctx, cb.GatewayTransactionID)
	if err == nil {
		log.Info().Str("tx", cb.GatewayTransactionID).Msg("payment: idempotent duplicate")
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "already processed",
			"tx":      existing.GatewayTransactionID,
			"status":  existing.Status,
		})
		return
	} else if !errors.Is(err, repository.ErrPaymentTransactionNotFound) {
		log.Error().Err(err).Msg("payment: tx lookup failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	// 4. 创建交易记录
	rawJSON, _ := json.Marshal(cb.RawBody)
	rawMap := models.JSONMap{}
	_ = json.Unmarshal(rawJSON, &rawMap)

	status := models.PaymentStatusPending
	switch cb.Status {
	case "success":
		status = models.PaymentStatusSuccess
	case "failed":
		status = models.PaymentStatusFailed
	}

	tx := &models.PaymentTransaction{
		OrderID:              order.ID,
		PaymentGateway:       gateway,
		GatewayTransactionID: cb.GatewayTransactionID,
		Amount:               cb.Amount,
		Currency:             cb.Currency,
		Status:               status,
		RawCallback:          rawMap,
	}

	if err := h.repo.CreatePayment(ctx, tx); err != nil {
		if errors.Is(err, repository.ErrPaymentTransactionDup) {
			c.JSON(http.StatusOK, gin.H{"code": 0, "message": "already processed (race)"})
			return
		}
		log.Error().Err(err).Msg("payment: create tx failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	// 5. 支付成功 → 订单状态流转（状态机校验）
	if status == models.PaymentStatusSuccess {
		if h.stateMachine.CanTransition(order.State, models.OrderStatePaid) {
			if err := h.repo.UpdateState(ctx, order.ID, models.OrderStatePaid); err != nil {
				log.Error().Err(err).Uint64("order_id", order.ID).Msg("payment: order state update failed")
			}
		} else {
			log.Warn().Str("from", order.State).Msg("payment: cannot transition to paid (invalid from state)")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":                    0,
		"gateway":                 gateway,
		"gateway_transaction_id":  cb.GatewayTransactionID,
		"order_no":                order.OrderNo,
		"amount":                  cb.Amount,
		"currency":                cb.Currency,
		"status":                  status,
	})
}

// ListTransactions — GET /payment/transactions (admin)
func (h *PaymentHandler) ListTransactions(c *gin.Context) {
	gateway := c.Query("gateway")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := h.repo.ListPaymentTransactions(c.Request.Context(), gateway, status, page, size)
	if err != nil {
		log.Error().Err(err).Msg("payment transactions list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list, "total": total, "page": page, "size": size})
}
