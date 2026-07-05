package hub

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// WSMessage digunakan untuk membaca roomId tujuan dari payload json frontend
type WSMessage struct {
	RoomID   string `json:"roomId"`
	SenderID string `json:"senderId"`
	Payload  string `json:"payload"`
}

// Hub mengelola koneksi WebSocket aktif berdasarkan ruangan (Room)
type Hub struct {
	mu        sync.RWMutex
	// Memetakan RoomID -> Daftar koneksi websocket aktif di dalam room tersebut
	rooms     map[string]map[*websocket.Conn]struct{}
	broadcast chan []byte
}

func New() *Hub {
	return &Hub{
		rooms:     make(map[string]map[*websocket.Conn]struct{}),
		broadcast: make(chan []byte, 256),
	}
}

func (h *Hub) Run() {
	for msgBytes := range h.broadcast {
		// Bongkar pesan JSON sebentar untuk tahu pesan ini milik room mana
		var msg WSMessage
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("ws unmarshal error: %v", err)
			continue
		}

		h.mu.RLock()
		// Hanya ambil koneksi yang berada di room tujuan (msg.RoomID)
		if conns, exists := h.rooms[msg.RoomID]; exists {
			for conn := range conns {
				if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
					log.Printf("ws write error: %v", err)
				}
			}
		}
		h.mu.RUnlock()
	}
}

func (h *Hub) HandleWS(c *websocket.Conn) {
	// Ambil roomId dari query parameter URL koneksi (Contoh: ws://localhost:8080/ws?roomId=ABC)
	roomID := c.Query("roomId")
	if roomID == "" {
		c.WriteMessage(websocket.TextMessage, []byte(`{"error": "roomId wajib diisi"}`))
		c.Close()
		return
	}

	h.mu.Lock()
	// Jika ruangan belum terdaftar di memori RAM, buat map baru untuk ruangan tersebut
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*websocket.Conn]struct{})
	}
	h.rooms[roomID][c] = struct{}{}
	h.mu.Unlock()

	log.Printf("👥 Koneksi baru masuk ke Room: %s", roomID)

	defer func() {
		h.mu.Lock()
		if conns, exists := h.rooms[roomID]; exists {
			delete(conns, c)
			// Jika ruangan sudah kosong melompong, hapus sekalian dari memori biar hemat RAM
			if len(conns) == 0 {
				delete(h.rooms, roomID)
			}
		}
		h.mu.Unlock()
		c.Close()
		log.Printf("❌ Koneksi keluar dari Room: %s", roomID)
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break // client disconnected
		}
		h.broadcast <- msg
	}
}