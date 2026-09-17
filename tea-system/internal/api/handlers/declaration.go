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
// DeclarationHandler — 报关台账 CRUD
// ============================================================

type DeclarationHandler struct {
	repo *repository.InvoiceRepo
}

func NewDeclarationHandler(repo *repository.InvoiceRepo) *DeclarationHandler {
	return &DeclarationHandler{repo: repo}
}

// DeclarationCreateRequest — POST /declarations
type DeclarationCreateRequest struct {
	OrderID              uint64  `json:"order_id" binding:"required"`
	CustomsDeclarationNo string  `json:"customs_declaration_no"`
	HsCode               string  `json:"hs_code" binding:"required"`
	CommodityDesc        string  `json:"commodity_desc" binding:"required"`
	GrossWeight          *float64 `json:"gross_weight,omitempty"`
	NetWeight            *float64 `json:"net_weight,omitempty"`
	DeclaredValue        *float64 `json:"declared_value,omitempty"`
	CustomsStatus        string  `json:"customs_status"`
	DeclarationDate      string  `json:"declaration_date,omitempty"`
	ClearedDate          string  `json:"cleared_date,omitempty"`
	Remarks              string  `json:"remarks,omitempty"`
}

// DeclarationUpdateRequest — PUT /declarations/:id
type DeclarationUpdateRequest struct {
	CustomsDeclarationNo *string  `json:"customs_declaration_no,omitempty"`
	HsCode               *string  `json:"hs_code,omitempty"`
	CommodityDesc        *string  `json:"commodity_desc,omitempty"`
	GrossWeight          *float64 `json:"gross_weight,omitempty"`
	NetWeight            *float64 `json:"net_weight,omitempty"`
	DeclaredValue        *float64 `json:"declared_value,omitempty"`
	CustomsStatus        *string  `json:"customs_status,omitempty"`
	DeclarationDate      *string  `json:"declaration_date,omitempty"`
	ClearedDate          *string  `json:"cleared_date,omitempty"`
	Remarks              *string  `json:"remarks,omitempty"`
}

// Create — POST /declarations
func (h *DeclarationHandler) Create(c *gin.Context) {
	var req DeclarationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	d := &models.DeclarationLedger{
		OrderID:                  req.OrderID,
		CustomsDeclarationNo:     req.CustomsDeclarationNo,
		HsCode:                   req.HsCode,
		CommodityDesc:            req.CommodityDesc,
		GrossWeight:              req.GrossWeight,
		NetWeight:                req.NetWeight,
		DeclaredValue:            req.DeclaredValue,
		CustomsStatus:            req.CustomsStatus,
		Remarks:                  req.Remarks,
	}

	if err := h.repo.CreateDeclaration(c.Request.Context(), d); err != nil {
		log.Error().Err(err).Msg("declaration: create failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, d)
}

// List — GET /declarations
func (h *DeclarationHandler) List(c *gin.Context) {
	orderID := uint64(0)
	if qid := c.Query("order_id"); qid != "" {
		if parsed, err := strconv.ParseUint(qid, 10, 64); err == nil {
			orderID = parsed
		}
	}
	status := c.Query("customs_status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list, total, err := h.repo.ListDeclarations(c.Request.Context(), repository.DeclarationListFilter{
		OrderID:       orderID,
		CustomsStatus: status,
		Page:          page,
		Size:          size,
	})
	if err != nil {
		log.Error().Err(err).Msg("declaration: list failed")
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

// Update — PUT /declarations/:id
func (h *DeclarationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var req DeclarationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	patch := map[string]interface{}{}
	if v := req.CustomsDeclarationNo; v != nil {
		patch["customs_declaration_no"] = *v
	}
	if v := req.HsCode; v != nil {
		patch["hs_code"] = *v
	}
	if v := req.CommodityDesc; v != nil {
		patch["commodity_desc"] = *v
	}
	if v := req.GrossWeight; v != nil {
		patch["gross_weight"] = *v
	}
	if v := req.NetWeight; v != nil {
		patch["net_weight"] = *v
	}
	if v := req.DeclaredValue; v != nil {
		patch["declared_value"] = *v
	}
	if v := req.CustomsStatus; v != nil {
		patch["customs_status"] = *v
	}
	if v := req.Remarks; v != nil {
		patch["remarks"] = *v
	}

	if len(patch) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "no fields to update"})
		return
	}

	if err := h.repo.UpdateDeclaration(c.Request.Context(), id, patch); err != nil {
		if errors.Is(err, repository.ErrDeclarationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "declaration not found"})
			return
		}
		log.Error().Err(err).Msg("declaration: update failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update failed"})
		return
	}

	d, _ := h.repo.GetDeclarationByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, d)
}

// Delete — DELETE /declarations/:id (软删除 = 置 voided_at)
func (h *DeclarationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	reason := c.Query("reason")
	if err := h.repo.SoftDeleteDeclaration(c.Request.Context(), id, reason); err != nil {
		if errors.Is(err, repository.ErrDeclarationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "declaration not found"})
			return
		}
		log.Error().Err(err).Msg("declaration: delete failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "voided"})
}
