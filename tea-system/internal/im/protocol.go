package im

import (
	"time"
)

// IMMessage — WebSocket 统一消息信封
type IMMessage struct {
	Type    string      `json:"type"`    // send_message / chat_message / barrage / ack / error / presence / join_room
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

// BarragePayload — 弹幕（直播间）
type BarragePayload struct {
	RoomID    string `json:"room_id"`
	Content   string `json:"content"`
	Nickname  string `json:"nickname,omitempty"`
	Badge     string `json:"badge,omitempty"`
	UserID    uint64 `json:"user_id,omitempty"`
	Timestamp int64  `json:"ts"`
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
