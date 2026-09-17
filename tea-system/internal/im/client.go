package im

import (
	"context"
	"encoding/json"
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
	messageRepo    *repository.MessageRepo
	translateSvc   *service.TranslateService
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

// RegisterConversation — 让此 client 订阅某会话的广播
func (c *Client) RegisterConversation(id uint64) {
	c.conversationIDs[id] = true
}

// ReadPump — 读消息循环：从 WS 读 → 存 DB → 异步翻译 → 广播
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// 读超时（心跳）
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
		case "ping":
			// 客户端心跳，忽略（pong 由 handler 自动发）
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
				// hub 关闭了 channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			// 将队列中的消息一并写出减少 syscall
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

// handleSendMessage — 处理客户端发消息
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

	if p.MessageType == "" {
		p.MessageType = "text"
	}

	// 构造 Message
	msg := &models.Message{
		ConversationID:   p.ConversationID,
		SenderType:       c.userType,
		SenderID:         c.userID,
		MessageType:      p.MessageType,
		Content:          p.Content,
		TranslationStatus: "pending",
	}

	// 初始翻译方向（同步一次，失败不阻塞）
	go c.asyncTranslate(msg)

	ctx := context.Background()
	if err := c.messageRepo.Create(ctx, msg); err != nil {
		log.Error().Err(err).Msg("failed to create message")
		c.writeError(500, "failed to save message")
		return
	}

	// 构造广播消息
	chatMsg := ChatMessage{
		ID:               msg.ID,
		ConversationID:   msg.ConversationID,
		SenderType:       msg.SenderType,
		SenderID:         msg.SenderID,
		Content:          msg.Content,
		TranslationZH:    msg.TranslationZH,
		TranslationEN:    msg.TranslationEN,
		TranslationStatus: msg.TranslationStatus,
		CreatedAt:        msg.CreatedAt,
	}

	_ = c.hub.BroadcastChatMessage(chatMsg)
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
		_ = c.messageRepo.UpdateTranslation(ctx, msg.ID, "", "", "pending")
		return
	}

	// 更新 DB
	var zh, en string
	if field == "translation_zh" {
		zh = translated
		en = msg.TranslationEN
	} else {
		zh = msg.TranslationZH
		en = translated
	}
	_ = c.messageRepo.UpdateTranslation(ctx, msg.ID, zh, en, "translated")

	// 再广播一条翻译完成的更新消息
	msg.TranslationZH = zh
	msg.TranslationEN = en
	msg.TranslationStatus = "translated"
	chatMsg := ChatMessage{
		ID:               msg.ID,
		ConversationID:   msg.ConversationID,
		SenderType:       msg.SenderType,
		SenderID:         msg.SenderID,
		Content:          msg.Content,
		TranslationZH:    msg.TranslationZH,
		TranslationEN:    msg.TranslationEN,
		TranslationStatus: msg.TranslationStatus,
		CreatedAt:        msg.CreatedAt,
	}
	_ = c.hub.BroadcastChatMessage(chatMsg)
}

func (c *Client) writeError(code int, msg string) {
	env := IMMessage{
		Type:    "error",
		Payload: ErrorPayload{Code: code, Message: msg},
	}
	raw, _ := json.Marshal(env)
	select {
	case c.send <- raw:
	default:
	}
}
