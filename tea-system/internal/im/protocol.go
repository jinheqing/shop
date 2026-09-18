package im

import (
	"time"
)

// IMMessage — WebSocket 统一消息信封
type IMMessage struct {
	Type    string      `json:"type"`    // send_message / chat_message / barrage / ack / error / presence / join_room / link_mic / subtitle
	Payload interface{} `json:"payload"` // 具体消息体
}

// Attachment — 附件（图片 / 视频 / 文件）
type Attachment struct {
	URL  string `json:"url"`
	Type string `json:"type"` // image / video / file
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
}

// SendMessagePayload — 客户端 → 服务器：发送消息
type SendMessagePayload struct {
	ConversationID uint64       `json:"conversation_id" binding:"required"`
	Content        string       `json:"content"`
	MessageType    string       `json:"message_type"`              // text / emoji / image / video / file / quote_card / order_card
	Attachments    []Attachment `json:"attachments,omitempty"`     // 图片/视频/文件：content 可空，附件走这里
	CardPayload    interface{}  `json:"card_payload,omitempty"`    // 卡片型消息的结构化 payload
	ClientMsgID    string       `json:"client_msg_id,omitempty"`   // 客户端生成去重用
}

// ChatMessage — 服务器 → 客户端：广播新消息
type ChatMessage struct {
	ID                uint64       `json:"id"`
	ConversationID    uint64       `json:"conversation_id"`
	SenderType        string       `json:"sender_type"` // user / staff / system
	SenderID          uint64       `json:"sender_id"`
	SenderName        string       `json:"sender_name,omitempty"`
	MessageType       string       `json:"message_type"`
	Content           string       `json:"content"`
	Attachments       []Attachment `json:"attachments,omitempty"`
	CardPayload       interface{}  `json:"card_payload,omitempty"`
	TranslationZH     string       `json:"translation_zh,omitempty"`
	TranslationEN     string       `json:"translation_en,omitempty"`
	TranslationStatus string       `json:"translation_status"` // pending / translated / failed
	CreatedAt         time.Time    `json:"created_at"`
}

// BarragePayload — 弹幕（直播间通用载体，subtype 区分业务）
//
// subtype 取值:
//   chat         — 普通弹幕聊天（原 barrage）
//   link_invite  — 观众请求连麦 / 主播邀请观众连麦
//   link_accept  — 被邀请者接受（请求者自己发布音视频）
//   link_reject  — 被邀请者拒绝
//   link_end     — 任一方结束连麦（双方都停推）
//   subtitle     — 主播/连麦者的字幕（ASR → 翻译 → 广播）
//
type BarragePayload struct {
	RoomID    string `json:"room_id"`
	Subtype   string `json:"subtype,omitempty"` // chat / link_invite / link_accept / link_reject / link_end / subtitle
	Content   string `json:"content"`
	Nickname  string `json:"nickname,omitempty"`
	Badge     string `json:"badge,omitempty"`
	UserID    uint64 `json:"user_id,omitempty"`
	Timestamp int64  `json:"ts"`

	// === link_mic 专属字段 ===
	TargetUserID  uint64 `json:"target_user_id,omitempty"`  // 邀请/拒绝谁
	LinkSessionID string `json:"link_session_id,omitempty"` // 本次连麦 session（双方一致）

	// === subtitle 专属字段 ===
	Speaker    string `json:"speaker,omitempty"`            // host / guest_{id}
	Text       string `json:"text,omitempty"`               // ASR 原文
	Translation string `json:"translation,omitempty"`       // 译文
	Lang       string `json:"lang,omitempty"`               // zh / en / auto
	TranscriptType string `json:"transcript_type,omitempty"` // partial / final
}

// AckPayload — 服务器 → 客户端：消息送达确认
type AckPayload struct {
	ClientMsgID string `json:"client_msg_id,omitempty"`
	MessageID   uint64 `json:"message_id"`
}

// ErrorPayload — 服务器 → 客户端：错误
type ErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
