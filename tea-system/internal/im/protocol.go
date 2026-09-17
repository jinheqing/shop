package im

import (
	"time"
)

// IMMessage — WebSocket 统一消息信封
type IMMessage struct {
	Type    string      `json:"type"`    // send_message / chat_message / ack / error / presence
	Payload interface{} `json:"payload"` // 具体消息体
}

// SendMessagePayload — 客户端 → 服务器：发送消息
type SendMessagePayload struct {
	ConversationID uint64 `json:"conversation_id" binding:"required"`
	Content        string `json:"content" binding:"required"`
	MessageType    string `json:"message_type"` // text / image / quote_link / file（默认 text）
}

// ChatMessage — 服务器 → 客户端：广播新消息
type ChatMessage struct {
	ID               uint64    `json:"id"`
	ConversationID   uint64    `json:"conversation_id"`
	SenderType       string    `json:"sender_type"` // user / staff / system
	SenderID         uint64    `json:"sender_id"`
	Content          string    `json:"content"`
	TranslationZH    string    `json:"translation_zh,omitempty"`
	TranslationEN    string    `json:"translation_en,omitempty"`
	TranslationStatus string   `json:"translation_status"` // pending / translated / failed
	CreatedAt        time.Time `json:"created_at"`
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
