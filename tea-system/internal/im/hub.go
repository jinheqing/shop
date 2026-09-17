package im

import (
	"encoding/json"
	"sync"

	"github.com/rs/zerolog/log"
)

// Hub — WebSocket 连接中心
type Hub struct {
	mu          sync.RWMutex
	clients     map[uint64]*Client        // key: user_id 或 staff_id
	broadcast   chan []byte              // 待广播的原始 JSON
	register    chan *Client
	unregister  chan *Client
}

// NewHub — 构造 Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint64]*Client),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Register — 外部注册 client（非阻塞，走 channel）
func (h *Hub) Register(c *Client) {
	select {
	case h.register <- c:
	default:
		log.Warn().Msg("hub register channel full")
	}
}

// Start — 后台 goroutine：处理注册/注销 + 广播
func (h *Hub) Start() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			log.Info().Uint64("user_id", client.userID).Msg("ws client registered")

		case client := <-h.unregister:
			h.mu.Lock()
			if existing, ok := h.clients[client.userID]; ok && existing == client {
				delete(h.clients, client.userID)
				close(client.send)
				log.Info().Uint64("user_id", client.userID).Msg("ws client unregistered")
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			// 解析 chat_message 的 conversation_id，只发给该会话参与者
			// 这里简单实现：先全量推，让客户端自己过滤；后续可优化按会话定向
			h.mu.RLock()
			for _, c := range h.clients {
				select {
				case c.send <- msg:
				default:
					// send buffer 满了，丢弃避免阻塞 hub
					log.Warn().Uint64("user_id", c.userID).Msg("ws client send buffer full, dropping")
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastChatMessage — 便捷方法：序列化 ChatMessage 并广播
func (h *Hub) BroadcastChatMessage(payload ChatMessage) error {
	env := IMMessage{Type: "chat_message", Payload: payload}
	raw, err := json.Marshal(env)
	if err != nil {
		return err
	}
	h.broadcast <- raw
	return nil
}

// ClientCount — 当前在线客户端数
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
