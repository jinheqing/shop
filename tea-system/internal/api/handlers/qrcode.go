package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QRCodeHandler struct {
	DB *gorm.DB
}

func NewQRCodeHandler(db *gorm.DB) *QRCodeHandler { return &QRCodeHandler{DB: db} }

// POST /qrcodes/generate
func (h *QRCodeHandler) Generate(c *gin.Context) {
	var body struct {
		CustomProductID uint64 `json:"custom_product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 重新生成 product_token（如果为 null 也生成）
	var buf [48]byte
	rand.Read(buf[:])
	token := base64.RawURLEncoding.EncodeToString(buf[:])

	if err := h.DB.Model(map[string]interface{}{}).Table("custom_products").
		Where("id = ?", body.CustomProductID).
		Update("product_token", token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"product_token": token,
		"public_url":    "/trace/" + token,
	})
}

// GET /public/qrcodes/:token — 公开扫码页面（返回 trace 数据）
func (h *QRCodeHandler) GetTrace(c *gin.Context) {
	token := c.Param("token")
	var cp struct {
		Title               string `json:"title"`
		TeaGardenLocation    string `json:"tea_garden_location"`
		MasterName          string `json:"master_name"`
		HarvestDate         string `json:"harvest_date"`
		RoastingDate        string `json:"roasting_date"`
		StorageLocation     string `json:"storage_location"`
		TeaType             string `json:"tea_type"`
		TeaShape            string `json:"tea_shape"`
	}
	if err := h.DB.Table("custom_products").
		Select("title, tea_garden_location, master_name, harvest_date, roasting_date, storage_location, tea_type, tea_shape").
		Where("product_token = ?", token).
		Scan(&cp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trace not found"})
		return
	}
	c.JSON(http.StatusOK, cp)
}
