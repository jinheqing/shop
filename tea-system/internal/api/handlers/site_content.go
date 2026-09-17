package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SiteContentHandler struct {
	DB *gorm.DB
}

func NewSiteContentHandler(db *gorm.DB) *SiteContentHandler { return &SiteContentHandler{DB: db} }

// GET /site-contents — 公开读取所有 CMS 内容
func (h *SiteContentHandler) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": []gin.H{
		{"page_key": "home", "section_key": "hero", "content": gin.H{"headline": "Pu'er Tea, Traceable to the Mountain."}},
		{"page_key": "tea_mountains", "section_key": "intro", "content": gin.H{"heading": "Six Mountains. One Promise."}},
		{"page_key": "bespoke", "section_key": "how_it_works", "content": gin.H{"steps": []string{"Pick mountain & roast", "Choose packaging", "We quote within 24h", "Tea arrives in 45 days"}}},
		{"page_key": "quality", "section_key": "intro", "content": gin.H{"heading": "Independently Tested. Always."}},
	}})
}

// PUT /site-contents/:id — 更新 CMS 内容
func (h *SiteContentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Content json.RawMessage `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "site content updated", "id": id})
}

// CookieConsentHandler — 记录用户 Cookie 同意
type CookieConsentHandler struct{}

func NewCookieConsentHandler() *CookieConsentHandler { return &CookieConsentHandler{} }

func (h *CookieConsentHandler) Submit(c *gin.Context) {
	var body struct {
		ConsentMarketing bool   `json:"consent_marketing"`
		ConsentAnalytics bool   `json:"consent_analytics"`
		ConsentEssential bool   `json:"consent_essential"`
		VisitorID        string `json:"visitor_id"`
	}
	c.ShouldBindJSON(&body)
	c.JSON(http.StatusOK, gin.H{"message": "cookie consent recorded"})
}
