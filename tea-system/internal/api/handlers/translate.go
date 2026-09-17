package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type TranslateHandler struct {
	FastAPIURL string
}

func NewTranslateHandler(fastAPIURL string) *TranslateHandler {
	return &TranslateHandler{FastAPIURL: fastAPIURL}
}

// POST /translate/text — 代理到 FastAPI
func (h *TranslateHandler) TranslateText(c *gin.Context) {
	if h.FastAPIURL == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "translate engine not configured"})
		return
	}
	body, _ := io.ReadAll(c.Request.Body)
	target, _ := url.Parse(h.FastAPIURL)
	target.Path = "/translate/text"
	proxy := httputil.NewSingleHostReverseProxy(target)
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	c.Request.URL = target
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		c.JSON(http.StatusBadGateway, gin.H{"error": "translate engine unavailable", "detail": err.Error()})
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

// GET /translate/status — 翻译引擎健康状态
func (h *TranslateHandler) Status(c *gin.Context) {
	if h.FastAPIURL == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_configured"})
		return
	}
	resp, err := http.Get(h.FastAPIURL + "/health")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unreachable", "detail": err.Error()})
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", data)
}

// GET /translate/asr — WebSocket ASR endpoint marker
func (h *TranslateHandler) ASR(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ws_endpoint": "/ws/translate/asr", "note": "upgrade connection in client"})
}
