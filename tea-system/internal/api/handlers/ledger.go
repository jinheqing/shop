package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"tea-system/internal/models"
	"tea-system/internal/repository"
)

// ============================================================
// LedgerHandler — 收汇台账 CRUD
// ============================================================

type LedgerHandler struct {
	repo *repository.InvoiceRepo
}

func NewLedgerHandler(repo *repository.InvoiceRepo) *LedgerHandler {
	return &LedgerHandler{repo: repo}
}

// LedgerCreateRequest — POST /ledgers
type LedgerCreateRequest struct {
	OrderID              uint64  `json:"order_id" binding:"required"`
	PaymentGateway       string  `json:"payment_gateway" binding:"required"`
	GatewayTransactionID string  `json:"gateway_transaction_id" binding:"required"`
	AmountGBP            float64 `json:"amount_gbp" binding:"required,gt=0"`
	ExchangeRate         *float64 `json:"exchange_rate,omitempty"`
	AmountCNY            *float64 `json:"amount_cny,omitempty"`
	BankAccount          string  `json:"bank_account,omitempty"`
	Remarks              string  `json:"remarks,omitempty"`
}

// LedgerUpdateRequest — PUT /ledgers/:id
type LedgerUpdateRequest struct {
	ExchangeRate   *float64 `json:"exchange_rate,omitempty"`
	AmountCNY      *float64 `json:"amount_cny,omitempty"`
	BankAccount    *string  `json:"bank_account,omitempty"`
	Remarks        *string  `json:"remarks,omitempty"`
}

// Create — POST /ledgers
func (h *LedgerHandler) Create(c *gin.Context) {
	var req LedgerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	l := &models.ForeignExchangeLedger{
		OrderID:              req.OrderID,
		PaymentGateway:       req.PaymentGateway,
		GatewayTransactionID: req.GatewayTransactionID,
		AmountGBP:            req.AmountGBP,
		ExchangeRate:         req.ExchangeRate,
		AmountCNY:            req.AmountCNY,
		BankAccount:          req.BankAccount,
		Remarks:              req.Remarks,
	}

	if err := h.repo.CreateForeignExchange(c.Request.Context(), l); err != nil {
		log.Error().Err(err).Msg("ledger: create failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, l)
}

// List — GET /ledgers
func (h *LedgerHandler) List(c *gin.Context) {
	orderID := uint64(0)
	if qid := c.Query("order_id"); qid != "" {
		if parsed, err := strconv.ParseUint(qid, 10, 64); err == nil {
			orderID = parsed
		}
	}
	gateway := c.Query("payment_gateway")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list, total, err := h.repo.ListForeignExchange(c.Request.Context(), repository.ForexListFilter{
		OrderID:        orderID,
		PaymentGateway: gateway,
		Page:           page,
		Size:           size,
	})
	if err != nil {
		log.Error().Err(err).Msg("ledger: list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": list,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// Update — PUT /ledgers/:id
func (h *LedgerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var req LedgerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	patch := map[string]interface{}{}
	if v := req.ExchangeRate; v != nil {
		patch["exchange_rate"] = *v
	}
	if v := req.AmountCNY; v != nil {
		patch["amount_cny"] = *v
	}
	if v := req.BankAccount; v != nil {
		patch["bank_account"] = *v
	}
	if v := req.Remarks; v != nil {
		patch["remarks"] = *v
	}

	if len(patch) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "no fields to update"})
		return
	}

	if err := h.repo.UpdateForeignExchange(c.Request.Context(), id, patch); err != nil {
		if errors.Is(err, repository.ErrForeignExchangeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "ledger not found"})
			return
		}
		log.Error().Err(err).Msg("ledger: update failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update failed"})
		return
	}

	l, _ := h.repo.GetForeignExchangeByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, l)
}

// Delete — DELETE /ledgers/:id
func (h *LedgerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	if err := h.repo.SoftDeleteForeignExchange(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrForeignExchangeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "ledger not found"})
			return
		}
		log.Error().Err(err).Msg("ledger: delete failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "deleted"})
}
