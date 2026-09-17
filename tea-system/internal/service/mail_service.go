package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// MailService — 邮件发送（Mailgun EU 节点，GDPR 合规）
// 生产走 Mailgun API，dev 环境记录到日志（不真发）
type MailService struct {
	apiKey  string
	domain  string
	from    string
	devMode bool // true 时只打日志，不真调用 Mailgun
	httpCli *http.Client
}

func NewMailService(apiKey, domain, from string) *MailService {
	devMode := apiKey == "" || apiKey == "key-xxxxxxxxxxxxxxxxxxxxxxx" || apiKey == "CHANGE_ME"
	return &MailService{
		apiKey:  apiKey,
		domain:  domain,
		from:    from,
		devMode: devMode,
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendMagicLink — 发送魔法链接邮件
func (m *MailService) SendMagicLink(ctx context.Context, toEmail, magicURL string) error {
	subject := "🔐 登录您的普洱茶顾问账户"
	textBody := fmt.Sprintf(
		"您好！\n\n"+
			"点击以下链接登录（10 分钟有效，单次使用）：\n\n"+
			"%s\n\n"+
			"如果您没有请求此邮件，请忽略。\n\n"+
			"—— 普洱茶顾问团队",
		magicURL,
	)
	htmlBody := fmt.Sprintf(`
<div style="font-family: Arial, sans-serif; max-width: 480px; margin: 0 auto;">
  <h2 style="color: #1A1A1A;">🔐 登录您的普洱茶顾问账户</h2>
  <p>点击下方按钮登录（10 分钟有效，单次使用）：</p>
  <a href="%s" style="display: inline-block; padding: 12px 24px; background: #B8A47C; color: white; text-decoration: none; border-radius: 6px;">立即登录</a>
  <p style="margin-top: 16px; color: #666; font-size: 14px;">或复制粘贴：<br>%s</p>
  <p style="margin-top: 24px; color: #999; font-size: 12px;">如果您没有请求此邮件，请忽略。</p>
</div>`, magicURL, magicURL)

	return m.send(ctx, toEmail, subject, textBody, htmlBody)
}

// SendOrderNotification — 订单状态变更通知
func (m *MailService) SendOrderNotification(ctx context.Context, toEmail, orderNo, newState string) error {
	subject := fmt.Sprintf("📦 订单 %s 状态更新：%s", orderNo, newState)
	textBody := fmt.Sprintf("您的订单 %s 状态变更为：%s", orderNo, newState)
	return m.send(ctx, toEmail, subject, textBody, textBody)
}

// SendLiveInvitation — 直播预告
func (m *MailService) SendLiveInvitation(ctx context.Context, toEmail, roomName, dateTime string) error {
	subject := fmt.Sprintf("🎥 专属直播邀请：%s", roomName)
	textBody := fmt.Sprintf("您好！邀请您参加专属直播：\n\n%s\n\n时间：%s", roomName, dateTime)
	return m.send(ctx, toEmail, subject, textBody, textBody)
}

func (m *MailService) send(ctx context.Context, toEmail, subject, textBody, htmlBody string) error {
	if m.devMode {
		log.Warn().
			Str("mailgun_api_key", "(dev mode, not sending)").
			Str("to", toEmail).
			Str("subject", subject).
			Str("body_preview", textBody[:min(len(textBody), 80)]).
			Msg("📧 mail (dev mode — logged only, not sent)")
		return nil
	}

	payload := map[string]interface{}{
		"from":    m.from,
		"to":      toEmail,
		"subject": subject,
		"text":    textBody,
		"html":    htmlBody,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", m.domain),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	req.SetBasicAuth("api", m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("mailgun: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("mailgun returned %d", resp.StatusCode)
	}

	log.Info().Str("to", toEmail).Str("subject", subject).Msg("📧 mail sent")
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
