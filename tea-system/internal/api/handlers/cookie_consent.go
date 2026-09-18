package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"tea-system/internal/middleware"
	"tea-system/internal/models"
)

type CookieConsentHandler struct {
	DB *gorm.DB
}

func NewCookieConsentHandler(db *gorm.DB) *CookieConsentHandler {
	return &CookieConsentHandler{DB: db}
}

func (h *CookieConsentHandler) Submit(c *gin.Context) {
	var body struct {
		ConsentEssential bool   `json:"consent_essential"`
		ConsentAnalytics bool   `json:"consent_analytics"`
		ConsentMarketing bool   `json:"consent_marketing"`
		VisitorID        string `json:"visitor_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	body.ConsentEssential = true // GDPR 强制必须同意
	uid := middleware.GetSubjectID(c)
	var userID *uint64
	if middleware.GetSubjectType(c) != "" && uid > 0 {
		userID = &uid
	}
	logEntry := models.CookieConsentLog{
		VisitorID:        body.VisitorID,
		UserID:           userID,
		ConsentEssential: body.ConsentEssential,
		ConsentAnalytics: body.ConsentAnalytics,
		ConsentMarketing: body.ConsentMarketing,
		UserAgent:        c.GetHeader("User-Agent"),
		IPAddress:        c.ClientIP(),
		CreatedAt:        time.Now().UTC(),
	}
	if h.DB != nil {
		if err := h.DB.Create(&logEntry).Error; err != nil {
			log.Warn().Err(err).Msg("cookie_consent: create failed")
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "cookie consent recorded"})
}
