package hub

import (
	"encoding/json"
	"log"
	"sync"
	"umbra-backend/internal/crypto"

	"github.com/gofiber/websocket/v2"
)

// WSPayload mewakili isi payload pesan dengan tanda tangan digital
type WSPayload struct {
	Ciphertext string `json:"ciphertext"`
	IV         string `json:"iv"`
	Signature  string `json:"signature"`
	Timestamp  string `json:"timestamp"`
	PublicKey  string `json:"publicKey"` // ECDSA public key base64 dari frontend
}

// WSMessage digunakan untuk membaca roomId tujuan dan detail payload json frontend
type WSMessage struct {
	Event          string    `json:"event"`
	RoomID         string    `json:"roomId"`
	SenderMemberID string    `json:"senderMemberId"`
	Payload        WSPayload `json:"payload"`
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

		// 1. Parse JSON pesan untuk memverifikasi tanda tangan
		var wsMsg WSMessage
		if err := json.Unmarshal(msg, &wsMsg); err == nil && wsMsg.Event == "send_message" {
			// Gabungkan data sesuai dengan yang di-sign di frontend (ciphertext + iv + timestamp)
			dataToVerify := []byte(wsMsg.Payload.Ciphertext + wsMsg.Payload.IV + wsMsg.Payload.Timestamp)

			// Panggil verifikasi
			valid, err := crypto.VerifySignature(dataToVerify, wsMsg.Payload.Signature, wsMsg.Payload.PublicKey)
			if err != nil || !valid {
				log.Printf("🚨 ECDSA Verification failed for member %s: %v", wsMsg.SenderMemberID, err)

				// Kirim balik respon error code 16 ke client pengirim
				errResp := map[string]interface{}{
					"event":   "error",
					"code":    16,
					"message": "Security Alert: Tanda tangan pesan tidak valid (ECDSA Verify Failed)!",
				}
				errBytes, _ := json.Marshal(errResp)
				_ = c.WriteMessage(websocket.TextMessage, errBytes)

				continue // Lewati broadcast agar pesan tidak terkirim ke member lain
			}
		}

		h.broadcast <- msg
	}
}