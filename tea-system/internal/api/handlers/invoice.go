package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"tea-system/internal/middleware"
	"tea-system/internal/repository"
	"tea-system/internal/service"
)

// ============================================================
// InvoiceHandler
// ============================================================

type InvoiceHandler struct {
	invoiceSvc *service.InvoiceService
	orderRepo  *repository.OrderRepo
	invRepo    *repository.InvoiceRepo
}

func NewInvoiceHandler(
	invoiceSvc *service.InvoiceService,
	orderRepo *repository.OrderRepo,
	invRepo *repository.InvoiceRepo,
) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceSvc: invoiceSvc,
		orderRepo:  orderRepo,
		invRepo:    invRepo,
	}
}

// ==================== handlers ====================

// GetByOrder — GET /orders/:id/invoice
func (h *InvoiceHandler) GetByOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid order id"})
		return
	}

	order, err := h.orderRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "order not found"})
			return
		}
		log.Error().Err(err).Msg("invoice: order lookup failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	// 用户只能看自己的发票
	subType := getUserSubType(c)
	if subType == "user" && (order.UserID == nil || *order.UserID != getUserID(c)) {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "forbidden"})
		return
	}

	inv, err := h.invRepo.GetInvoiceByOrderID(c.Request.Context(), order.ID)
	if err != nil {
		if errors.Is(err, repository.ErrInvoiceNotFound) {
			// 如果没有，自动生成一个
			inv, err = h.invoiceSvc.GenerateInvoice(c.Request.Context(), order.ID)
			if err != nil {
				log.Error().Err(err).Msg("invoice: generate failed")
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "invoice generate failed"})
				return
			}
		} else {
			log.Error().Err(err).Msg("invoice: repo lookup failed")
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
			return
		}
	}

	c.JSON(http.StatusOK, inv)
}

// DownloadPDF — GET /orders/:id/invoice/pdf
func (h *InvoiceHandler) DownloadPDF(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid order id"})
		return
	}

	order, err := h.orderRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "order not found"})
			return
		}
		log.Error().Err(err).Msg("invoice: order lookup failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	subType := getUserSubType(c)
	if subType == "user" && (order.UserID == nil || *order.UserID != getUserID(c)) {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "forbidden"})
		return
	}

	invRepo := h.invRepo
	inv, err := invRepo.GetInvoiceByOrderID(c.Request.Context(), order.ID)
	if err != nil {
		if errors.Is(err, repository.ErrInvoiceNotFound) {
			inv, err = h.invoiceSvc.GenerateInvoice(c.Request.Context(), order.ID)
			if err != nil {
				log.Error().Err(err).Msg("invoice: generate failed")
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "invoice generate failed"})
				return
			}
		} else {
			log.Error().Err(err).Msg("invoice: repo lookup failed")
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
			return
		}
	}

	pdfBytes, err := h.invoiceSvc.ReadPDF(inv.PdfURL)
	if err != nil {
		log.Error().Err(err).Str("path", inv.PdfURL).Msg("invoice: pdf read failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "pdf file not found"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", `attachment; filename="`+inv.InvoiceNo+`.pdf"`)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// Regenerate — POST /orders/:id/invoice/regenerate
func (h *InvoiceHandler) Regenerate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid order id"})
		return
	}

	order, err := h.orderRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "order not found"})
			return
		}
		log.Error().Err(err).Msg("invoice: order lookup failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	invRepo := h.invRepo
	inv, err := invRepo.GetInvoiceByOrderID(c.Request.Context(), order.ID)
	if errors.Is(err, repository.ErrInvoiceNotFound) {
		// 还不存在，直接 Generate
		inv, err = h.invoiceSvc.GenerateInvoice(c.Request.Context(), order.ID)
	} else if err == nil {
		inv, err = h.invoiceSvc.RegenerateInvoice(c.Request.Context(), inv.ID)
	}
	if err != nil {
		log.Error().Err(err).Msg("invoice: regenerate failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "invoice regenerate failed"})
		return
	}

	c.JSON(http.StatusOK, inv)
}

// ==================== 辅助函数（从 gin Context 取值） ====================
// 统一复用 middleware 包里的函数

func getUserID(c *gin.Context) uint64 {
	return middleware.GetSubjectID(c)
}

func getUserSubType(c *gin.Context) string {
	return middleware.GetSubjectType(c)
}
