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
// 修复 (2026-09-19): 仅 title + unit_price 为必填，其余全部可选并在 handler 里填默认值
type CustomProductCreateRequest struct {
	Title             string   `json:"title" binding:"required"`
	UnitPrice         float64  `json:"unit_price" binding:"required,gt=0"`

	// 以下全部可选，handler 填默认值
	RawTeaSource      *string  `json:"raw_tea_source,omitempty"`
	CustomRequirement *string  `json:"custom_requirement,omitempty"`
	TeaType           *string  `json:"tea_type,omitempty"`
	TeaShape          *string  `json:"tea_shape,omitempty"`
	InnerPackaging    *string  `json:"inner_packaging,omitempty"`
	OuterPackaging    *string  `json:"outer_packaging,omitempty"`
	QrCodePosition    *string  `json:"qr_code_position,omitempty"`
	Quantity          *int     `json:"quantity,omitempty"`
	ShippingCost      *float64 `json:"shipping_cost,omitempty"`
	LeadTime          *string  `json:"lead_time,omitempty"`
	HarvestDate       *string  `json:"harvest_date,omitempty"`
	RoastingDate      *string  `json:"roasting_date,omitempty"`
	TeaGardenLocation  *string  `json:"tea_garden_location,omitempty"`
	MasterName        *string  `json:"master_name,omitempty"`
	StorageLocation   *string  `json:"storage_location,omitempty"`

	TeaShapeWeight    *int    `json:"tea_shape_weight,omitempty"`
	SmokedWithFlower  *bool   `json:"smoked_with_flower,omitempty"`
	FlowerType        string  `json:"flower_type,omitempty"`
	ProductCardText   string  `json:"product_card_text,omitempty"`
	ProductCardFormat string  `json:"product_card_format,omitempty"`
	SKU               string  `json:"sku,omitempty"`
	SgsReportID       *uint64 `json:"sgs_report_id,omitempty"`
	IsBespoke         *bool   `json:"is_bespoke,omitempty"`
	NonRefundable     *bool   `json:"non_refundable,omitempty"`

	IncludeCustomLive  *bool   `json:"include_custom_live,omitempty"`
	LiveScheduledDate *string `json:"live_scheduled_date,omitempty"`
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

	now := time.Now()
	today := now.Format("2006-01-02")
	teaType := valOr(req.TeaType, "Pu'er")
	quantity := valOrInt(req.Quantity, 100)
	shippingCost := valOrFloat(req.ShippingCost, 0)

	harvestStr := valOr(req.HarvestDate, today)
	roastingStr := valOr(req.RoastingDate, today)
	harvest, _ := time.Parse("2006-01-02", harvestStr)
	roasting, _ := time.Parse("2006-01-02", roastingStr)

	liveScheduled := parseOptionalDate(req.LiveScheduledDate)

	isBespoke := true
	if req.IsBespoke != nil {
		isBespoke = *req.IsBespoke
	}
	nonRefundable := true
	if req.NonRefundable != nil {
		nonRefundable = *req.NonRefundable
	}
	smoked := false
	if req.SmokedWithFlower != nil {
		smoked = *req.SmokedWithFlower
	}
	includeLive := false
	if req.IncludeCustomLive != nil {
		includeLive = *req.IncludeCustomLive
	}

	sku := req.SKU
	if sku == "" {
		sku = fmt.Sprintf("SKU-%s-%d%04d-%03d", teaType, now.Year(), now.YearDay(), now.Nanosecond()%1000)
	}

	staffID := middleware.GetSubjectID(c)

	p := &models.CustomProduct{
		ProductToken:      nil,
		Version:           1,
		Status:            models.CustomProductStatusDraft,
		IsBespoke:         isBespoke,
		NonRefundable:     nonRefundable,
		Title:             req.Title,
		RawTeaSource:      valOr(req.RawTeaSource, ""),
		CustomRequirement: valOr(req.CustomRequirement, ""),
		TeaType:           teaType,
		TeaShape:          valOr(req.TeaShape, "cake"),
		TeaShapeWeight:    req.TeaShapeWeight,
		SmokedWithFlower:  smoked,
		FlowerType:        req.FlowerType,
		InnerPackaging:    valOr(req.InnerPackaging, ""),
		OuterPackaging:    valOr(req.OuterPackaging, ""),
		ProductCardText:   req.ProductCardText,
		ProductCardFormat: func() string { if req.ProductCardFormat != "" { return req.ProductCardFormat }; return "vertical" }(),
		QrCodePosition:    valOr(req.QrCodePosition, "right-bottom"),
		SKU:               sku,
		UnitPrice:         req.UnitPrice,
		Quantity:          quantity,
		ShippingCost:      shippingCost,
		TotalAmount:       h.svc.CalcTotalAmount(req.UnitPrice, quantity, shippingCost),
		LeadTime:          valOr(req.LeadTime, "14-21 days"),
		HarvestDate:       harvest,
		RoastingDate:      roasting,
		TeaGardenLocation:  valOr(req.TeaGardenLocation, ""),
		MasterName:        valOr(req.MasterName, ""),
		StorageLocation:   valOr(req.StorageLocation, ""),
		SgsReportID:       req.SgsReportID,
		IncludeCustomLive: includeLive,
		LiveScheduledDate: liveScheduled,
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

// ===== helpers =====
func valOr(p *string, def string) string { if p != nil && *p != "" { return *p }; return def }
func valOrInt(p *int, def int) int { if p != nil && *p > 0 { return *p }; return def }
func valOrFloat(p *float64, def float64) float64 { if p != nil { return *p }; return def }
func parseOptionalDate(p *string) *time.Time {
	if p == nil || *p == "" { return nil }
	if t, err := time.Parse("2006-01-02", *p); err == nil { return &t }
	return nil
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

// Published — GET /custom-products/published（公开查询已发布定制茶）
func (h *CustomProductHandler) Published(c *gin.Context) {
	items, total, err := h.repo.List(c.Request.Context(), repository.ListFilter{
		Status: "published",
		Page:   1,
		Size:   20,
	})
	if err != nil {
		log.Error().Err(err).Msg("custom_product: published list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total})
}

// CreatePublic — POST /public/custom-products（公开提交定制请求，无需登录）
// 轻量版本：允许缺失某些 admin 才知道的必填字段（自动填默认值）
func (h *CustomProductHandler) CreatePublic(c *gin.Context) {
	var req struct {
		Title              string  `json:"title"`
		TeaType            string  `json:"tea_type"`
		TeaShape           string  `json:"tea_shape"`
		UnitPrice          float64 `json:"unit_price"`
		Quantity           int     `json:"quantity"`
		ShippingCost       float64 `json:"shipping_cost"`
		InnerPackaging     string  `json:"inner_packaging"`
		OuterPackaging     string  `json:"outer_packaging"`
		TeaGardenLocation  string  `json:"tea_garden_location"`
		MasterName         string  `json:"master_name"`
		RawTeaSource       string  `json:"raw_tea_source"`
		CustomRequirement  string  `json:"custom_requirement"`
		CustomerName       string  `json:"customer_name"`
		CustomerEmail      string  `json:"customer_email" binding:"omitempty,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	if req.Title == "" {
		req.Title = "Bespoke Pu'er Request"
	}
	if req.TeaType == "" {
		req.TeaType = "sheng_puer"
	}
	if req.TeaShape == "" {
		req.TeaShape = "cake_357g"
	}
	if req.TeaGardenLocation == "" {
		req.TeaGardenLocation = "Master Selection"
	}
	if req.RawTeaSource == "" {
		req.RawTeaSource = req.TeaGardenLocation
	}

	// 生成 product_token 便于后续报价/支付
	token, _ := h.svc.GenerateProductToken()
	now := time.Now()

	p := &models.CustomProduct{
		ProductToken:      &token,
		Version:           1,
		Status:            models.CustomProductStatusDraft,
		IsBespoke:         true,
		NonRefundable:     true,
		Title:             req.Title,
		RawTeaSource:      req.RawTeaSource,
		CustomRequirement: req.CustomRequirement,
		TeaType:           req.TeaType,
		TeaShape:          req.TeaShape,
		InnerPackaging:    req.InnerPackaging,
		OuterPackaging:    req.OuterPackaging,
		UnitPrice:         req.UnitPrice,
		Quantity:          req.Quantity,
		ShippingCost:      req.ShippingCost,
		TeaGardenLocation: req.TeaGardenLocation,
		MasterName:        req.MasterName,
		LeadTime:          "45 days from confirmation",
		HarvestDate:       now,
		RoastingDate:      now,
		StorageLocation:   "London (temporary)",
	}
	if err := h.repo.Create(c.Request.Context(), p); err != nil {
		log.Error().Err(err).Msg("custom_product: public create failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "bespoke request submitted",
		"token":   token,
		"product": p,
	})
}
