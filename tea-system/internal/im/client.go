package im

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"tea-system/internal/models"
	"tea-system/internal/repository"
	"tea-system/internal/service"
)

// Upgrader — WebSocket 升级器（允许所有 origin）
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client — 单个 WebSocket 连接
type Client struct {
	hub             *Hub
	conn            *websocket.Conn
	send            chan []byte
	userID          uint64
	userType        string // "user" / "staff"
	conversationIDs map[uint64]bool

	// 依赖
	messageRepo  *repository.MessageRepo
	translateSvc *service.TranslateService
}

// NewClient — 构造 Client
func NewClient(
	hub *Hub,
	conn *websocket.Conn,
	userID uint64,
	userType string,
	messageRepo *repository.MessageRepo,
	translateSvc *service.TranslateService,
) *Client {
	return &Client{
		hub:             hub,
		conn:            conn,
		send:            make(chan []byte, 256),
		userID:          userID,
		userType:        userType,
		conversationIDs: make(map[uint64]bool),
		messageRepo:     messageRepo,
		translateSvc:    translateSvc,
	}
}

// ReadPump — 读消息循环：从 WS 读 → 分发处理
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(65536)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warn().Err(err).Uint64("user_id", c.userID).Msg("ws read error")
			}
			return
		}

		var env IMMessage
		if err := json.Unmarshal(raw, &env); err != nil {
			c.writeError(400, "invalid message envelope")
			continue
		}

		switch env.Type {
		case "send_message":
			c.handleSendMessage(env.Payload)
		case "barrage":
			c.handleBarrage(env.Payload)
		case "join_conversation":
			c.handleJoinConversation(env.Payload)
		case "leave_conversation":
			c.handleLeaveConversation(env.Payload)
		case "join_room":
			c.handleJoinRoom(env.Payload)
		case "leave_room":
			c.handleLeaveRoom(env.Payload)
		case "ping":
			// 心跳忽略
		default:
			c.writeError(400, "unknown message type: "+env.Type)
		}
	}
}

// WritePump — 写消息循环：从 hub.broadcast 收消息 → 写 WS
func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleJoinConversation — 订阅某会话
func (c *Client) handleJoinConversation(payload interface{}) {
	var p struct {
		ConversationID uint64 `json:"conversation_id"`
	}
	raw, _ := json.Marshal(payload)
	if err := json.Unmarshal(raw, &p); err != nil || p.ConversationID == 0 {
		return
	}
	c.conversationIDs[p.ConversationID] = true
	c.hub.SubscribeConversation(p.ConversationID, c.userID)
}

// handleLeaveConversation — 取消订阅会话
func (c *Client) handleLeaveConversation(payload interface{}) {
	var p struct {
		ConversationID uint64 `json:"conversation_id"`
	}
	raw, _ := json.Marshal(payload)
	if err := json.Unmarshal(raw, &p); err != nil {
		return
	}
	delete(c.conversationIDs, p.ConversationID)
}

// handleJoinRoom — 订阅直播间弹幕
func (c *Client) handleJoinRoom(payload interface{}) {
	var p struct {
		RoomID string `json:"room_id"`
	}
	raw, _ := json.Marshal(payload)
	if err := json.Unmarshal(raw, &p); err != nil || p.RoomID == "" {
		return
	}
	c.hub.SubscribeBarrage(p.RoomID, c.userID)
}

// handleLeaveRoom — 取消订阅直播间弹幕
func (c *Client) handleLeaveRoom(payload interface{}) {
	// no-op, hub handles subscriptions in-memory; leave is optional
}

// handleBarrage — 处理弹幕消息（存入临时弹幕表或直接广播）
func (c *Client) handleBarrage(payload interface{}) {
	raw, _ := json.Marshal(payload)
	var p BarragePayload
	if err := json.Unmarshal(raw, &p); err != nil {
		c.writeError(400, "invalid barrage payload")
		return
	}
	p.UserID = c.userID
	p.Timestamp = time.Now().UnixMilli()
	_ = c.hub.BroadcastBarrage(p)
}

// handleSendMessage — 处理客户端发消息（核心增强）
func (c *Client) handleSendMessage(payload interface{}) {
	raw, err := json.Marshal(payload)
	if err != nil {
		c.writeError(400, "invalid payload")
		return
	}

	var p SendMessagePayload
	if err := json.Unmarshal(raw, &p); err != nil {
		c.writeError(400, "invalid send_message payload")
		return
	}

	// 自动订阅会话（首次发消息时）
	c.conversationIDs[p.ConversationID] = true
	c.hub.SubscribeConversation(p.ConversationID, c.userID)

	if p.MessageType == "" {
		p.MessageType = "text"
	}

	// 构造 DB model
	msg := &models.Message{
		ConversationID:    p.ConversationID,
		SenderType:        c.userType,
		SenderID:          c.userID,
		MessageType:       p.MessageType,
		Content:           p.Content,
		TranslationStatus: "pending",
	}

	// 序列化 attachments → JSONB map
	if len(p.Attachments) > 0 {
		attBytes, _ := json.Marshal(p.Attachments)
		msg.AttachmentURLs = models.JSONMap{
			"items": json.RawMessage(attBytes),
		}
	}

	// card_payload → 存到 Content 的 JSON 字符串（兼容现有字段）
	if p.CardPayload != nil {
		cardBytes, _ := json.Marshal(p.CardPayload)
		msg.Content = string(cardBytes)
	}

	ctx := context.Background()
	if err := c.messageRepo.Create(ctx, msg); err != nil {
		log.Error().Err(err).Msg("failed to create message")
		c.writeError(500, "failed to save message")
		return
	}

	// 翻译（仅 text/emoji）
	if p.MessageType == "text" || p.MessageType == "emoji" {
		go c.asyncTranslate(msg)
	}

	// 构造广播消息
	attachmentsBroadcast := p.Attachments
	if msg.AttachmentURLs != nil {
		// 从 JSONMap 解回来
		if raw, ok := msg.AttachmentURLs["items"]; ok {
			var atts []Attachment
			_ = json.Unmarshal(raw.(json.RawMessage), &atts)
			attachmentsBroadcast = atts
		}
	}

	chatMsg := ChatMessage{
		ID:             msg.ID,
		ConversationID: msg.ConversationID,
		SenderType:     msg.SenderType,
		SenderID:       msg.SenderID,
		MessageType:    p.MessageType,
		Content:        p.Content, // 发送原 content，card 类型前端自行渲染
		Attachments:    attachmentsBroadcast,
		CardPayload:    p.CardPayload,
		CreatedAt:      msg.CreatedAt,
	}

	_ = c.hub.BroadcastChatMessage(chatMsg)

	// ACK（含 client_msg_id）
	if p.ClientMsgID != "" {
		ack := IMMessage{
			Type:    "ack",
			Payload: AckPayload{ClientMsgID: p.ClientMsgID, MessageID: msg.ID},
		}
		rawAck, _ := json.Marshal(ack)
		select {
		case c.send <- rawAck:
		default:
		}
	}
}

// asyncTranslate — 异步翻译并更新 DB
func (c *Client) asyncTranslate(msg *models.Message) {
	lang := service.DetectLang(msg.Content)
	var targetLang, field string

	if lang == "zh" {
		targetLang = "en"
		field = "translation_en"
	} else {
		targetLang = "zh"
		field = "translation_zh"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	translated, err := c.translateSvc.Translate(ctx, msg.Content, lang, targetLang)
	if err != nil {
		log.Warn().Err(err).Msg("translate failed, keep pending")
		return
	}

	var zh, en string
	if field == "translation_zh" {
		zh = translated
		en = msg.TranslationEN
	} else {
		zh = msg.TranslationZH
		en = translated
	}
	_ = c.messageRepo.UpdateTranslation(ctx, msg.ID, zh, en, "translated")

	chatMsg := ChatMessage{
		ID:               msg.ID,
		ConversationID:   msg.ConversationID,
		SenderType:       msg.SenderType,
		SenderID:         msg.SenderID,
		Content:          msg.Content,
		TranslationZH:    zh,
		TranslationEN:    en,
		TranslationStatus: "translated",
		CreatedAt:        msg.CreatedAt,
	}
	_ = c.hub.BroadcastChatMessage(chatMsg)
}

func (c *Client) writeError(code int, msg string) {
	env := IMMessage{
		Type:    "error",
		Payload: ErrorPayload{Code: code, Message: fmt.Sprintf("%s", msg)},
	}
	raw, _ := json.Marshal(env)
	select {
	case c.send <- raw:
	default:
	}
}
