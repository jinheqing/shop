package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// TranslateService — 调 FastAPI /translate/text 接口
// 环境变量 TRANSLATE_SERVICE_URL 默认 http://localhost:8090
// 如果翻译服务不可用，返回空 string（允许降级，不 fatal）
type TranslateService struct {
	baseURL string
	client  *http.Client
}

func NewTranslateService(baseURL string) *TranslateService {
	if baseURL == "" {
		baseURL = os.Getenv("TRANSLATE_SERVICE_URL")
		baseURL = "http://localhost:8090"
	}
	return &TranslateService{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

// translateRequest — FastAPI /translate/text 请求体
type translateRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

// translateResponse — FastAPI /translate/text 响应体
// Python 侧 ok() 会把数据包一层: {"code":0, "message":"ok", "data": {...}}
// 但我们也兼容直接返回 {"translation":"xxx"} 的裸响应（防后续有人改了包装）
type translateResponse struct {
	Translation string `json:"translation"`

	// Python ok() 包装层
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		Translation string `json:"translation"`
	} `json:"data"`
}

// Translate — 翻译文本。失败时返回空 string + error（调用方决定是否降级）
func (s *TranslateService) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	if text == "" {
		return "", nil
	}

	body, _ := json.Marshal(translateRequest{
		Text:       text,
		SourceLang: sourceLang,
		TargetLang: targetLang,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/translate/text", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("translate http %d: %s", resp.StatusCode, string(raw))
	}

	// 先尝试 Python ok() 包装层
	var wrapped struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    *struct {
			Translation string `json:"translation"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil && wrapped.Data != nil && wrapped.Data.Translation != "" {
		return wrapped.Data.Translation, nil
	}

	// 再尝试裸响应 {"translation":"xxx"}
	var out struct {
		Translation string `json:"translation"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", fmt.Errorf("translate json parse: %w", err)
	}
	return out.Translation, nil
}

// DetectLang — 简单启发式判断语言（后续可替换为 FastAPI /detect）
// 规则：包含 CJK 字符 → zh；否则 → en
func DetectLang(text string) string {
	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FFF {
			return "zh"
		}
	}
	return "en"
}
