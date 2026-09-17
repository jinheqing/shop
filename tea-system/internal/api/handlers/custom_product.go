package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"tea-system/internal/middleware"
	"tea-system/internal/models"
	"tea-system/internal/repository"
	"tea-system/internal/service"
)

// ============================================================
// CustomProduct — 定制报价 CRUD
// ============================================================

type CustomProductHandler struct {
	repo *repository.CustomProductRepo
	svc  *service.CustomProductService
}

func NewCustomProductHandler(
	repo *repository.CustomProductRepo,
	svc *service.CustomProductService,
) *CustomProductHandler {
	return &CustomProductHandler{repo: repo, svc: svc}
}

// ==================== 请求 DTO ====================

// CustomProductCreateRequest — POST /custom-products
type CustomProductCreateRequest struct {
	// 必填
	Title             string  `json:"title" binding:"required"`
	RawTeaSource      string  `json:"raw_tea_source" binding:"required"`
	CustomRequirement string  `json:"custom_requirement" binding:"required"`
	TeaType           string  `json:"tea_type" binding:"required"`
	TeaShape          string  `json:"tea_shape" binding:"required"`
	InnerPackaging    string  `json:"inner_packaging" binding:"required"`
	OuterPackaging    string  `json:"outer_packaging" binding:"required"`
	QrCodePosition    string  `json:"qr_code_position" binding:"required"`
	UnitPrice         float64 `json:"unit_price" binding:"required,gt=0"`
	Quantity          int     `json:"quantity" binding:"required,gt=0"`
	ShippingCost      float64 `json:"shipping_cost" binding:"required,gte=0"`
	LeadTime          string  `json:"lead_time" binding:"required"`
	HarvestDate       string  `json:"harvest_date" binding:"required"`
	RoastingDate      string  `json:"roasting_date" binding:"required"`
	TeaGardenLocation  string  `json:"tea_garden_location" binding:"required"`
	MasterName        string  `json:"master_name" binding:"required"`
	StorageLocation   string  `json:"storage_location" binding:"required"`

	// 可选
	TeaShapeWeight    *int   `json:"tea_shape_weight,omitempty"`
	SmokedWithFlower  bool   `json:"smoked_with_flower"`
	FlowerType        string `json:"flower_type,omitempty"`
	ProductCardText   string `json:"product_card_text,omitempty"`
	ProductCardFormat string `json:"product_card_format,omitempty"`
	SKU               string `json:"sku,omitempty"`
	SgsReportID       *uint64 `json:"sgs_report_id,omitempty"`
        IsBespoke         *bool   `json:"is_bespoke,omitempty"`
        NonRefundable     *bool   `json:"non_refundable,omitempty"`

        // 直播定制字段（2026-09 补齐前后端漂移）
        IncludeCustomLive bool   `json:"include_custom_live"`
        LiveScheduledDate string `json:"live_scheduled_date,omitempty"`
}

// CustomProductUpdateRequest — PUT /custom-products/:id（所有字段可选 patch）
type CustomProductUpdateRequest struct {
	Title             *string  `json:"title,omitempty"`
	RawTeaSource      *string  `json:"raw_tea_source,omitempty"`
	CustomRequirement *string  `json:"custom_requirement,omitempty"`
	TeaType           *string  `json:"tea_type,omitempty"`
	TeaShape          *string  `json:"tea_shape,omitempty"`
	TeaShapeWeight    *int     `json:"tea_shape_weight,omitempty"`
	SmokedWithFlower  *bool    `json:"smoked_with_flower,omitempty"`
	FlowerType        *string  `json:"flower_type,omitempty"`
	InnerPackaging    *string  `json:"inner_packaging,omitempty"`
	OuterPackaging    *string  `json:"outer_packaging,omitempty"`
	ProductCardText   *string  `json:"product_card_text,omitempty"`
	ProductCardFormat *string  `json:"product_card_format,omitempty"`
	QrCodePosition    *string  `json:"qr_code_position,omitempty"`
	SKU               *string  `json:"sku,omitempty"`
	UnitPrice         *float64 `json:"unit_price,omitempty"`
	Quantity          *int     `json:"quantity,omitempty"`
	ShippingCost      *float64 `json:"shipping_cost,omitempty"`
	LeadTime          *string  `json:"lead_time,omitempty"`
	HarvestDate       *string  `json:"harvest_date,omitempty"`
	RoastingDate      *string  `json:"roasting_date,omitempty"`
	TeaGardenLocation  *string  `json:"tea_garden_location,omitempty"`
	MasterName        *string  `json:"master_name,omitempty"`
	StorageLocation   *string  `json:"storage_location,omitempty"`
	SgsReportID       *uint64  `json:"sgs_report_id,omitempty"`
	IsBespoke         *bool    `json:"is_bespoke,omitempty"`
	NonRefundable     *bool    `json:"non_refundable,omitempty"`

	// 直播定制字段（2026-09 补齐前后端漂移）
	IncludeCustomLive *bool   `json:"include_custom_live,omitempty"`
	LiveScheduledDate *string `json:"live_scheduled_date,omitempty"`
}

// ==================== handler 方法 ====================

// Create — POST /custom-products
func (h *CustomProductHandler) Create(c *gin.Context) {
	var req CustomProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	harvest, err := time.Parse("2006-01-02", req.HarvestDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid harvest_date (expect YYYY-MM-DD)"})
		return
	}
	roasting, err := time.Parse("2006-01-02", req.RoastingDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid roasting_date (expect YYYY-MM-DD)"})
		return
	}

	// live_scheduled_date — 可选
	var liveScheduled *time.Time
	if req.LiveScheduledDate != "" {
		if t, err := time.Parse("2006-01-02", req.LiveScheduledDate); err == nil {
			liveScheduled = &t
		}
	}

	isBespoke := true
	if req.IsBespoke != nil {
		isBespoke = *req.IsBespoke
	}
	nonRefundable := true
	if req.NonRefundable != nil {
		nonRefundable = *req.NonRefundable
	}

	staffID := middleware.GetSubjectID(c)

	p := &models.CustomProduct{
		ProductToken:     nil, // publish 时生成
		Version:          1,
		Status:           models.CustomProductStatusDraft,
		IsBespoke:        isBespoke,
		NonRefundable:    nonRefundable,
		Title:            req.Title,
		RawTeaSource:     req.RawTeaSource,
		CustomRequirement: req.CustomRequirement,
		TeaType:          req.TeaType,
		TeaShape:         req.TeaShape,
		TeaShapeWeight:   req.TeaShapeWeight,
		SmokedWithFlower: req.SmokedWithFlower,
		FlowerType:       req.FlowerType,
		InnerPackaging:   req.InnerPackaging,
		OuterPackaging:   req.OuterPackaging,
		ProductCardText:  req.ProductCardText,
		QrCodePosition:   req.QrCodePosition,
		SKU:              func() string { if req.SKU != "" { return req.SKU }; return fmt.Sprintf("SKU-%s-%d%04d-%03d", req.TeaType, time.Now().Year(), time.Now().YearDay(), time.Now().Nanosecond()%1000) }(),
		UnitPrice:        req.UnitPrice,
		Quantity:         req.Quantity,
		ShippingCost:     req.ShippingCost,
		TotalAmount:      h.svc.CalcTotalAmount(req.UnitPrice, req.Quantity, req.ShippingCost),
		LeadTime:         req.LeadTime,
		HarvestDate:      harvest,
		RoastingDate:     roasting,
		TeaGardenLocation: req.TeaGardenLocation,
		MasterName:       req.MasterName,
		StorageLocation:  req.StorageLocation,
		SgsReportID:      req.SgsReportID,
		IncludeCustomLive: req.IncludeCustomLive,
		LiveScheduledDate: liveScheduled,
	}
	if req.ProductCardFormat != "" {
		p.ProductCardFormat = req.ProductCardFormat
	} else {
		p.ProductCardFormat = "vertical"
	}
	if staffID > 0 {
		p.CreatedByStaffID = &staffID
	}

	ctx := c.Request.Context()
	if err := h.repo.Create(ctx, p); err != nil {
		log.Error().Err(err).Msg("custom_product: create failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to create custom product"})
		return
	}

	c.JSON(http.StatusCreated, p)
}

// List — GET /custom-products
func (h *CustomProductHandler) List(c *gin.Context) {
	status := c.Query("status")
	keyword := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	items, total, err := h.repo.List(c.Request.Context(), repository.ListFilter{
		Status:  status,
		Keyword: keyword,
		Page:    page,
		Size:    size,
	})
	if err != nil {
		log.Error().Err(err).Msg("custom_product: list failed")
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

// GetByID — GET /custom-products/:id
func (h *CustomProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	p, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrCustomProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "custom product not found"})
			return
		}
		log.Error().Err(err).Msg("custom_product: get failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// Update — PUT /custom-products/:id (version 自动 +1)
func (h *CustomProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var req CustomProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	patch, err := buildUpdatePatch(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	if len(patch) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "no fields to update"})
		return
	}

	// 如果价格相关字段有变化 —— 重新计算 total_amount
	if _, hasUnit := patch["unit_price"]; hasUnit {
		// 需要一起看 quantity
		var unit float64
		var qty int
		if v, ok := patch["unit_price"].(float64); ok {
			unit = v
		}
		// 拿当前的 quantity
		current, err := h.repo.GetByID(c.Request.Context(), id)
		if err == nil {
			qty = current.Quantity
			unitPrice := current.UnitPrice
			if v, ok := patch["unit_price"].(float64); ok {
				unitPrice = v
			}
			if v, ok := patch["quantity"].(int); ok {
				qty = v
			}
			shipping := current.ShippingCost
			if v, ok := patch["shipping_cost"].(float64); ok {
				shipping = v
			}
			patch["total_amount"] = h.svc.CalcTotalAmount(unitPrice, qty, shipping)
		}
		_ = unit // 无 lint 警告
	} else {
		// 如果只 quantity / shipping_cost 变了也要重算
		if _, hasQty := patch["quantity"]; hasQty || false {
			current, err := h.repo.GetByID(c.Request.Context(), id)
			if err == nil {
				unitPrice := current.UnitPrice
				qty := current.Quantity
				shipping := current.ShippingCost
				if v, ok := patch["unit_price"].(float64); ok {
					unitPrice = v
				}
				if v, ok := patch["quantity"].(int); ok {
					qty = v
				}
				if v, ok := patch["shipping_cost"].(float64); ok {
					shipping = v
				}
				patch["total_amount"] = h.svc.CalcTotalAmount(unitPrice, qty, shipping)
			}
		}
	}

	_ = h.repo.Update // 占位（真正 patch 用 Update 会 bump version）

	if err := h.repo.Update(c.Request.Context(), id, patch); err != nil {
		if errors.Is(err, repository.ErrCustomProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "custom product not found"})
			return
		}
		log.Error().Err(err).Msg("custom_product: update failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update failed"})
		return
	}

	// 返回最新状态
	p, _ := h.repo.GetByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, p)
}

// buildUpdatePatch — 把 CustomProductUpdateRequest 中所有非 nil 字段转成 map
func buildUpdatePatch(req CustomProductUpdateRequest) (map[string]interface{}, error) {
	patch := map[string]interface{}{}

	if v := req.Title; v != nil {
		patch["title"] = *v
	}
	if v := req.RawTeaSource; v != nil {
		patch["raw_tea_source"] = *v
	}
	if v := req.CustomRequirement; v != nil {
		patch["custom_requirement"] = *v
	}
	if v := req.TeaType; v != nil {
		patch["tea_type"] = *v
	}
	if v := req.TeaShape; v != nil {
		patch["tea_shape"] = *v
	}
	if v := req.TeaShapeWeight; v != nil {
		patch["tea_shape_weight"] = *v
	}
	if v := req.SmokedWithFlower; v != nil {
		patch["smoked_with_flower"] = *v
	}
	if v := req.FlowerType; v != nil {
		patch["flower_type"] = *v
	}
	if v := req.InnerPackaging; v != nil {
		patch["inner_packaging"] = *v
	}
	if v := req.OuterPackaging; v != nil {
		patch["outer_packaging"] = *v
	}
	if v := req.ProductCardText; v != nil {
		patch["product_card_text"] = *v
	}
	if v := req.ProductCardFormat; v != nil {
		patch["product_card_format"] = *v
	}
	if v := req.QrCodePosition; v != nil {
		patch["qr_code_position"] = *v
	}
	if v := req.SKU; v != nil {
		patch["sku"] = *v
	}
	if v := req.UnitPrice; v != nil {
		patch["unit_price"] = *v
	}
	if v := req.Quantity; v != nil {
		patch["quantity"] = *v
	}
	if v := req.ShippingCost; v != nil {
		patch["shipping_cost"] = *v
	}
	if v := req.LeadTime; v != nil {
		patch["lead_time"] = *v
	}
	if v := req.HarvestDate; v != nil {
		t, err := time.Parse("2006-01-02", *v)
		if err != nil {
			return nil, errors.New("invalid harvest_date (YYYY-MM-DD)")
		}
		patch["harvest_date"] = t
	}
	if v := req.RoastingDate; v != nil {
		t, err := time.Parse("2006-01-02", *v)
		if err != nil {
			return nil, errors.New("invalid roasting_date (YYYY-MM-DD)")
		}
		patch["roasting_date"] = t
	}
	if v := req.TeaGardenLocation; v != nil {
		patch["tea_garden_location"] = *v
	}
	if v := req.MasterName; v != nil {
		patch["master_name"] = *v
	}
	if v := req.StorageLocation; v != nil {
		patch["storage_location"] = *v
	}
	if v := req.SgsReportID; v != nil {
		patch["sgs_report_id"] = *v
	}
	if v := req.IsBespoke; v != nil {
		patch["is_bespoke"] = *v
	}
	if v := req.NonRefundable; v != nil {
		patch["non_refundable"] = *v
	}
	if v := req.IncludeCustomLive; v != nil {
		patch["include_custom_live"] = *v
	}
	if v := req.LiveScheduledDate; v != nil {
		if *v == "" {
			patch["live_scheduled_date"] = nil
		} else if t, err := time.Parse("2006-01-02", *v); err == nil {
			patch["live_scheduled_date"] = t
		}
	}

	return patch, nil
}

// Delete — DELETE /custom-products/:id（软删除 + 状态置 archived）
func (h *CustomProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	if err := h.repo.SoftDelete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrCustomProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "custom product not found"})
			return
		}
		log.Error().Err(err).Msg("custom_product: delete failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "archived"})
}

// Publish — POST /custom-products/:id/publish
func (h *CustomProductHandler) Publish(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	token, err := h.svc.GenerateProductToken()
	if err != nil {
		log.Error().Err(err).Msg("custom_product: token gen failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "token generation failed"})
		return
	}

	if err := h.repo.Publish(c.Request.Context(), id, token); err != nil {
		if errors.Is(err, repository.ErrCustomProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "custom product not found"})
			return
		}
		log.Error().Err(err).Msg("custom_product: publish failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "publish failed"})
		return
	}

	// 返回完整对象
	p, _ := h.repo.GetByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, p)
}

// Review — POST /custom-products/:id/review
func (h *CustomProductHandler) Review(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	staffID := middleware.GetSubjectID(c)

	if err := h.repo.Review(c.Request.Context(), id, staffID); err != nil {
		if errors.Is(err, repository.ErrCustomProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "custom product not found"})
			return
		}
		log.Error().Err(err).Msg("custom_product: review failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "review failed"})
		return
	}

	p, _ := h.repo.GetByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, p)
}

// GetByToken — GET /custom-products/by-token/:token（**公开接口，无 JWT**）
func (h *CustomProductHandler) GetByToken(c *gin.Context) {
	token := c.Param("token")
	if len(token) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "token required"})
		return
	}

	p, err := h.repo.GetByToken(c.Request.Context(), token)
	if err != nil {
		if errors.Is(err, repository.ErrCustomProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "custom product not found"})
			return
		}
		log.Error().Err(err).Msg("custom_product: get by token failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, p)
}
