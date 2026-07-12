package hub

import (
	"encoding/json"
	"log"
	"sync"
	"umbra-backend/internal/crypto"
	roomhandler "umbra-backend/internal/room/handler"

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

// Client menyimpan detail koneksi dan kunci publik untuk E2EE
type Client struct {
	Conn           *websocket.Conn
	MemberID       string
	EcdhPublicKey  string
	EcdsaPublicKey string
}

// Hub mengelola koneksi WebSocket aktif berdasarkan ruangan (Room)
type Hub struct {
	mu        sync.RWMutex
	// Memetakan RoomID -> Daftar client aktif di dalam room tersebut
	rooms     map[string]map[*websocket.Conn]*Client
	broadcast chan []byte
}

func New() *Hub {
	return &Hub{
		rooms:     make(map[string]map[*websocket.Conn]*Client),
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
	// Ambil metadata dari query parameter URL koneksi
	roomID := c.Query("roomId")
	if roomID == "" {
		c.WriteMessage(websocket.TextMessage, []byte(`{"error": "roomId wajib diisi"}`))
		c.Close()
		return
	}

	memberID := c.Query("memberId")
	ecdhPublicKey := c.Query("ecdhPublicKey")
	ecdsaPublicKey := c.Query("ecdsaPublicKey")

	client := &Client{
		Conn:           c,
		MemberID:       memberID,
		EcdhPublicKey:  ecdhPublicKey,
		EcdsaPublicKey: ecdsaPublicKey,
	}

	h.mu.Lock()
	// Jika ruangan belum terdaftar di memori RAM, buat map baru untuk ruangan tersebut
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*websocket.Conn]*Client)
	}
	h.rooms[roomID][c] = client

	// Jika ada client lain di room ini, kirim info peer_joined secara timbal balik
	if len(h.rooms[roomID]) > 1 && memberID != "" {
		for otherConn, otherClient := range h.rooms[roomID] {
			if otherConn != c {
				// 1. Kirim data client baru ke client lama
				peerJoinedEvent := map[string]interface{}{
					"event":  "peer_joined",
					"roomId": roomID,
					"payload": map[string]string{
						"memberId":       memberID,
						"ecdhPublicKey":  ecdhPublicKey,
						"ecdsaPublicKey": ecdsaPublicKey,
					},
				}
				eventBytes, _ := json.Marshal(peerJoinedEvent)
				_ = otherConn.WriteMessage(websocket.TextMessage, eventBytes)

				// 2. Kirim data client lama ke client baru
				peerJoinedEventSelf := map[string]interface{}{
					"event":  "peer_joined",
					"roomId": roomID,
					"payload": map[string]string{
						"memberId":       otherClient.MemberID,
						"ecdhPublicKey":  otherClient.EcdhPublicKey,
						"ecdsaPublicKey": otherClient.EcdsaPublicKey,
					},
				}
				eventBytesSelf, _ := json.Marshal(peerJoinedEventSelf)
				_ = c.WriteMessage(websocket.TextMessage, eventBytesSelf)
			}
		}
	}
	h.mu.Unlock()

	log.Printf("👥 Koneksi baru masuk ke Room: %s (Member: %s)", roomID, memberID)

	defer func() {
		h.mu.Lock()
		if conns, exists := h.rooms[roomID]; exists {
			delete(conns, c)
			
			// Jika masih ada sisa client di room, kirim event peer_left ke mereka
			if len(conns) > 0 {
				peerLeftEvent := map[string]interface{}{
					"event":  "peer_left",
					"roomId": roomID,
					"payload": map[string]string{
						"memberId": memberID,
					},
				}
				eventBytes, _ := json.Marshal(peerLeftEvent)
				for otherConn := range conns {
					_ = otherConn.WriteMessage(websocket.TextMessage, eventBytes)
				}
			} else {
				// Jika ruangan sudah kosong melompong, hapus sekalian dari memori biar hemat RAM
				delete(h.rooms, roomID)
				roomhandler.DeleteRoom(roomID)
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

			// Ubah event ke "message" agar bisa dikenali oleh frontend penerima
			wsMsg.Event = "message"
			if marshaledMsg, err := json.Marshal(wsMsg); err == nil {
				msg = marshaledMsg
			}
		}

		h.broadcast <- msg
	}
}