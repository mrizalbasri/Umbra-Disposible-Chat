package hub

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Hub manages all active WebSocket connections.
// ponytail: broadcast-to-all model; add room-scoped routing when needed
type Hub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]struct{}
	broadcast chan []byte
}

func New() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]struct{}),
		broadcast: make(chan []byte, 256),
	}
}

func (h *Hub) Run() {
	for msg := range h.broadcast {
		h.mu.RLock()
		for conn := range h.clients {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("ws write error: %v", err)
			}
		}
		h.mu.RUnlock()
	}
}

func (h *Hub) HandleWS(c *websocket.Conn) {
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, c)
		h.mu.Unlock()
		c.Close()
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break // client disconnected
		}
		h.broadcast <- msg
	}
}
