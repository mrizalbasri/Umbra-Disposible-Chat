package hub

import (
	"encoding/json"
	"log"
	"sync"
	"time"
	"umbra-backend/internal/crypto"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn           *websocket.Conn
	RoomID         string
	MemberID       string
	EcdhPublicKey  string
	EcdsaPublicKey string
	PublicKey      string
	Send           chan []byte
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
		broadcast:  make(chan MessageBroadcast, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}

			currentMemberCount := len(h.rooms[client.RoomID])
			h.rooms[client.RoomID][client] = true

			if currentMemberCount >= 1 {
				for peer := range h.rooms[client.RoomID] {
					if peer != client {
						joinedNotify, _ := json.Marshal(map[string]interface{}{
							"event": "peer_joined",
							"payload": map[string]string{
								"publicKey": client.PublicKey,
							},
						})
						select {
						case peer.Send <- joinedNotify:
						default:
							close(peer.Send)
							delete(h.rooms[client.RoomID], peer)
						}
					}
				}
			}
			h.mu.Unlock()

			log.Printf("🔌 Client terhubung ke Room: %s", client.RoomID)

		case client := <-h.unregister:
			h.mu.Lock()
			if connections, exists := h.rooms[client.RoomID]; exists {
				if _, ok := connections[client]; ok {
					delete(connections, client)
					close(client.Send)
					log.Printf("❌ Client keluar dari Room: %s (memberId=%s)", client.RoomID, client.MemberID)

					remainingCount := len(connections)

					if remainingCount == 1 {
						for peer := range connections {
							destroyedNotify, _ := json.Marshal(map[string]interface{}{
								"event": "room_destroyed",
								"payload": map[string]string{
									"reason": "all_members_left",
								},
							})
							_ = peer.Conn.WriteMessage(websocket.TextMessage, destroyedNotify)
							_ = peer.Conn.Close()
							close(peer.Send)
							delete(connections, peer)
						}
					} else if remainingCount > 1 {
						for peer := range connections {
							leftNotify, _ := json.Marshal(map[string]interface{}{
								"event": "peer_left",
								"payload": map[string]string{
									"publicKey": client.PublicKey,
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

				if len(connections) == 0 {
					delete(h.rooms, client.RoomID)
				}
			}
			h.mu.Unlock()

		case bReq := <-h.broadcast:
			h.mu.Lock()
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
			h.mu.Unlock()
		}
	}
}

func (h *Hub) HandleWS(c *websocket.Conn) {
	roomID := c.Query("roomId")
	pubKey := c.Query("publicKey")

	if roomID == "" {
		h.sendError(c, "11", "roomId diperlukan di query parameter")
		_ = c.Close()
		return
	}
	if pubKey == "" {
		h.sendError(c, "12", "publicKey diperlukan di query parameter")
		_ = c.Close()
		return
	}

	client := &Client{
		Conn:           c,
		RoomID:         roomID,
		EcdhPublicKey:  pubKey,
		EcdsaPublicKey: pubKey,
		PublicKey:      pubKey,
		Send:           make(chan []byte, 256),
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

			publicKeyForVerification := client.PublicKey
			if publicKeyForVerification == "" {
				h.sendError(c, "16", "Public key tidak tersedia untuk verifikasi tanda tangan")
				continue
			}

			dataToVerify := []byte(payload.Ciphertext + payload.IV + payload.Timestamp)
			validSig, verifyErr := crypto.VerifySignature(dataToVerify, payload.Signature, publicKeyForVerification)
			if verifyErr != nil || !validSig {
				log.Printf("🚨 ECDSA verification failed (room=%s member=%s): %v", client.RoomID, client.MemberID, verifyErr)
				h.sendError(c, "16", "ECDSA Signature tidak valid")
				continue
			}

			validPayloadBytes, err := json.Marshal(payload)
			if err != nil {
				continue
			}

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
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}
