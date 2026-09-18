package service

import (
	"bytes"
	"context"
	"encoding/json"
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
type translateResponse struct {
	Translation string `json:"translation"`
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

	if resp.StatusCode != http.StatusOK {
		return "", io.EOF // 非 200 视为失败
	}

	var out translateResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
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
