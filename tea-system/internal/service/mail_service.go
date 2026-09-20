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
	subject := fmt.Sprintf("Order %s — Status Update: %s", orderNo, newState)
	textBody := fmt.Sprintf("Your order %s has been updated to: %s", orderNo, newState)
	return m.send(ctx, toEmail, subject, textBody, textBody)
}

// SendShippingNotification — 发货/物流更新通知（含物流单号与承运商）
func (m *MailService) SendShippingNotification(ctx context.Context, toEmail, orderNo, carrier, trackingNo string) error {
	subject := fmt.Sprintf("Your Order %s Has Been Dispatched", orderNo)

	trackingLine := ""
	if trackingNo != "" {
		trackingLine = fmt.Sprintf("Tracking Number: %s", trackingNo)
	}
	carrierLine := ""
	if carrier != "" {
		carrierLine = fmt.Sprintf("Carrier: %s", carrier)
	}

	textBody := fmt.Sprintf(
		"Dear Member,\n\n"+
			"Your bespoke order %s has been dispatched from our atelier.\n\n"+
			"%s\n%s\n\n"+
			"Should you wish to discuss your delivery, your advisor remains at your service.\n\n"+
			"— UK Tea House",
		orderNo, carrierLine, trackingLine,
	)

	htmlBody := fmt.Sprintf(`
<div style="font-family: Georgia, 'Times New Roman', serif; max-width: 520px; margin: 0 auto; color: #0B0A09;">
  <div style="border-bottom: 1px solid #C5A572; padding-bottom: 16px; margin-bottom: 24px;">
    <h2 style="margin: 0; font-weight: 400; letter-spacing: 0.05em;">UK Tea House</h2>
  </div>
  <p style="font-size: 15px; line-height: 1.7;">Dear Member,</p>
  <p style="font-size: 15px; line-height: 1.7;">
    Your bespoke order <strong>%s</strong> has been dispatched from our atelier.
  </p>
  <div style="background: #F8F5EF; border-left: 2px solid #C5A572; padding: 16px 20px; margin: 20px 0;">
    %s<br>
    %s
  </div>
  <p style="font-size: 14px; line-height: 1.7; color: #6b6459;">
    Should you wish to discuss your delivery, your advisor remains at your service.
  </p>
  <p style="font-size: 13px; color: #8a8578; margin-top: 32px; border-top: 1px solid #C5A572; padding-top: 16px;">
    — UK Tea House
  </p>
</div>`, orderNo, carrierLine, trackingLine)

	return m.send(ctx, toEmail, subject, textBody, htmlBody)
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
