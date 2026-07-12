package hub

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn      *websocket.Conn
	RoomID    string
	PublicKey string // Tambahkan field ini untuk mencatat public key client saat connect
	Send      chan []byte
}

type WSMessageWrapper struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

type SendMessagePayload struct {
	Ciphertext string `json:"ciphertext"`
	IV         string `json:"iv"`
	Signature  string `json:"signature"`
	Timestamp  string `json:"timestamp"`
}

type Hub struct {
	rooms      map[string]map[*Client]bool
	broadcast  chan MessageBroadcast
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type MessageBroadcast struct {
	RoomID  string
	Sender  *Client
	Message []byte
}

func New() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		broadcast:  make(chan MessageBroadcast),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// Jika room belum ada, buat map baru
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			
			// Ambil total member SEBELUM client baru ini dimasukkan
			currentMemberCount := len(h.rooms[client.RoomID])

			// Masukkan client baru ke dalam room
			h.rooms[client.RoomID][client] = true
			
			// Bagian 4.3: Jika sebelumnya sudah ada 1 orang (artinya client ini adalah orang ke-2),
			// kirim notifikasi peer_joined ke orang pertama yang sudah standby di room.
			if currentMemberCount == 1 {
				for peer := range h.rooms[client.RoomID] {
					// Kirim ke peer (orang pertama), bukan ke client yang baru masuk
					if peer != client {
						joinedNotify, _ := json.Marshal(map[string]interface{}{
							"event": "peer_joined",
							"payload": map[string]string{
								"publicKey": client.PublicKey,
							},
						})
						peer.Send <- joinedNotify
					}
				}
			}
			h.mu.Unlock()
			log.Printf("🔌 Client terhubung ke Room: %s", client.RoomID)

		case client := <-h.unregister:
			h.mu.Lock()
			if connections, exists := h.rooms[client.RoomID]; exists {
				if _, ok := connections[client]; ok {
					// 1. Hapus client yang keluar dari memori room
					delete(connections, client)
					log.Printf("❌ Client dengan transmisi kunci %s keluar dari Room: %s", client.PublicKey, client.RoomID)

					// 2. Hitung sisa orang di dalam room SEBELAHNYA
					remainingCount := len(connections)

					if remainingCount == 1 {
						// ⚠️ KONDISI A: Hanya sisa 1 orang -> Hancurkan room total demi privasi
						for peer := range connections {
							destroyedNotify, _ := json.Marshal(map[string]interface{}{
								"event": "room_destroyed",
								"payload": map[string]string{
									"reason": "all members left",
								},
							})
							_ = peer.Conn.WriteMessage(websocket.TextMessage, destroyedNotify)
							_ = peer.Conn.Close()
							delete(connections, peer)
						}
					} else if remainingCount > 1 {
						// 💬 KONDISI B: Masih ada lebih dari 1 orang -> Kirim peer_left (untuk skala group chat)
						for peer := range connections {
							leftNotify, _ := json.Marshal(map[string]interface{}{
								"event": "peer_left",
								"payload": map[string]string{
									"publicKey": client.PublicKey, // Kasih tahu public key siapa yang pergi
								},
							})
							select {
							case peer.Send <- leftNotify:
							default:
								close(peer.Send)
								delete(connections, peer)
							}
						}
					}
				}

				// 3. Bersihkan sisa map room jika benar-benar sudah kosong
				if len(connections) == 0 {
					delete(h.rooms, client.RoomID)
				}
			}
			h.mu.Unlock()

		case bReq := <-h.broadcast:
			h.mu.RLock()
			connections := h.rooms[bReq.RoomID]
			for client := range connections {
				if client != bReq.Sender {
					select {
					case client.Send <- bReq.Message:
					default:
						close(client.Send)
						delete(connections, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) HandleWS(c *websocket.Conn) {
	roomID := c.Query("roomId")
	pubKey := c.Query("publicKey") // Ambil public key dari query params untuk event peer_joined
	
	if roomID == "" {
		h.sendError(c, "11", "roomId diperlukan di query parameter")
		_ = c.Close()
		return
	}

	client := &Client{
		Conn:      c,
		RoomID:    roomID,
		PublicKey: pubKey,
		Send:      make(chan []byte, 256),
	}

	h.register <- client

	go client.writePump()

	defer func() {
		h.unregister <- client
		_ = c.Close()
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		var wrapper WSMessageWrapper
		if err := json.Unmarshal(msg, &wrapper); err != nil {
			h.sendError(c, "12", "Format JSON wrapper tidak valid")
			break 
		}

		if wrapper.Event == "send_message" {
			var payload SendMessagePayload
			if err := json.Unmarshal(wrapper.Payload, &payload); err != nil {
				h.sendError(c, "12", "Format payload tidak sesuai dokumen kontrak")
				break 
			}

			_, errTimestamp := time.Parse(time.RFC3339, payload.Timestamp)
			if payload.Ciphertext == "" || payload.IV == "" || payload.Signature == "" || errTimestamp != nil {
				h.sendError(c, "16", "ECDSA Signature tidak valid atau payload cacat")
				break 
			}

			// FIX: Lakukan marshalling payload yang valid agar tidak membungkus ganda json bytes
			validPayloadBytes, err := json.Marshal(payload)
			if err != nil {
				continue
			}

			// Map event terluar menjadi "message" sesuai spesifikasi Halaman 16
			outboundMessage := WSMessageWrapper{
				Event:   "message", 
				Payload: validPayloadBytes,
			}

			jsonBytes, err := json.Marshal(outboundMessage)
			if err != nil {
				continue
			}

			h.broadcast <- MessageBroadcast{
				RoomID:  client.RoomID,
				Sender:  client,
				Message: jsonBytes,
			}
		}
	}
}

func (h *Hub) sendError(c *websocket.Conn, code string, msg string) {
	_ = c.WriteJSON(map[string]interface{}{
		"responseMessage": msg,
		"responseCode":    code,
		"data":            nil,
	})
}

func (c *Client) writePump() {
	defer func() {
		_ = c.Conn.Close()
	}()
	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}