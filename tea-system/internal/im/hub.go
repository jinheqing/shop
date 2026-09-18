package im

import (
	"encoding/json"
	"sync"

	"github.com/rs/zerolog/log"
)

// Hub — WebSocket 连接中心
type Hub struct {
	mu         sync.RWMutex
	clients    map[uint64]*Client  // key: user_id 或 staff_id
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client

	// conversation → 订阅该会话的 user_id 集合（用于定向广播）
	convSubs   map[uint64]map[uint64]bool
	convSubsMu sync.RWMutex

	// room_id → 订阅该直播间弹幕的 user_id 集合
	barrageSubs   map[string]map[uint64]bool
	barrageSubsMu sync.RWMutex
}

// NewHub — 构造 Hub
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[uint64]*Client),
		broadcast:   make(chan []byte, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		convSubs:    make(map[uint64]map[uint64]bool),
		barrageSubs: make(map[string]map[uint64]bool),
	}
}

// SubscribeConversation — client 订阅某会话（定向广播）
func (h *Hub) SubscribeConversation(convID, userID uint64) {
	h.convSubsMu.Lock()
	defer h.convSubsMu.Unlock()
	if h.convSubs[convID] == nil {
		h.convSubs[convID] = make(map[uint64]bool)
	}
	h.convSubs[convID][userID] = true
}

// SubscribeBarrage — client 订阅某直播间弹幕
func (h *Hub) SubscribeBarrage(roomID string, userID uint64) {
	h.barrageSubsMu.Lock()
	defer h.barrageSubsMu.Unlock()
	if h.barrageSubs[roomID] == nil {
		h.barrageSubs[roomID] = make(map[uint64]bool)
	}
	h.barrageSubs[roomID][userID] = true
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
			h.mu.RLock()
			for _, c := range h.clients {
				select {
				case c.send <- msg:
				default:
					log.Warn().Uint64("user_id", c.userID).Msg("ws client send buffer full, dropping")
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastChatMessage — 序列化 ChatMessage 并定向广播给会话参与者
func (h *Hub) BroadcastChatMessage(payload ChatMessage) error {
	env := IMMessage{Type: "chat_message", Payload: payload}
	raw, err := json.Marshal(env)
	if err != nil {
		return err
	}

	// 定向：只发给该 conversation 的订阅者
	h.convSubsMu.RLock()
	userIDs := h.convSubs[payload.ConversationID]
	h.convSubsMu.RUnlock()

	if len(userIDs) > 0 {
		h.mu.RLock()
		for uid := range userIDs {
			if c, ok := h.clients[uid]; ok {
				select {
				case c.send <- raw:
				default:
				}
			}
		}
		h.mu.RUnlock()
	} else {
		// 无订阅者 → fallback 全量
		h.broadcast <- raw
	}
	return nil
}

// BroadcastBarrage — 定向弹幕广播给直播间订阅者
func (h *Hub) BroadcastBarrage(payload BarragePayload) error {
	env := IMMessage{Type: "barrage", Payload: payload}
	raw, err := json.Marshal(env)
	if err != nil {
		return err
	}

	h.barrageSubsMu.RLock()
	userIDs := h.barrageSubs[payload.RoomID]
	h.barrageSubsMu.RUnlock()

	h.mu.RLock()
	for uid := range userIDs {
		if c, ok := h.clients[uid]; ok {
			select {
			case c.send <- raw:
			default:
			}
		}
	}
	h.mu.RUnlock()
	return nil
}

// ClientCount — 当前在线客户端数
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
