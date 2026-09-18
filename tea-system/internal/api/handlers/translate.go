package handlers

import (
	"bytes"
	"encoding/json"
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
// Python /health 返回 ok() 包装层 {"code":0, "data": {...}}
// TranslateStatus.vue 直接读 s.status / s.gpu_available / s.queue_depth ...
// 我们在这里统一解包并补字段，让前端不管 Python 有没有都能拿到
func (h *TranslateHandler) Status(c *gin.Context) {
	if h.FastAPIURL == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "not_configured",
			"engine":    "tea-translate",
			"endpoint":  h.FastAPIURL,
			"note":      "Configure TRANSLATE_SERVICE_URL env or System Settings → Translate",
		})
		return
	}

	resp, err := http.Get(h.FastAPIURL + "/health")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "unreachable",
			"endpoint": h.FastAPIURL,
			"detail":   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	// 先尝试 Python ok() 包装层
	var wrapped struct {
		Code    int                    `json:"code"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil && wrapped.Data != nil {
		out := gin.H{
			"status":              wrapped.Data["status"],
			"whisper_loaded":      wrapped.Data["whisper_loaded"],
			"nllb_loaded":         wrapped.Data["nllb_loaded"],
			"engine":              "tea-translate",
			"endpoint":            h.FastAPIURL,
			"gpu_available":       false, // Python health 暂未提供，默认 false
			"queue_depth":         0,
			"avg_latency_ms":      0,
			"languages_supported": 2,
			"total_translations_today": 0,
		}
		c.JSON(http.StatusOK, out)
		return
	}

	// fallback：裸响应
	c.Data(resp.StatusCode, "application/json", raw)
}

// GET /translate/asr — WebSocket ASR endpoint marker
// 前端直连 tea-translate 的 WebSocket（不走 Go 代理），这里只返回 marker + 完整地址
func (h *TranslateHandler) ASR(c *gin.Context) {
	if h.FastAPIURL == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "translate engine not configured"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ws_endpoint": h.FastAPIURL + "/translate/asr-stream",
		"note":        "connect WebSocket with binary int16 PCM audio",
		"sample_rate": 16000,
	})
}
