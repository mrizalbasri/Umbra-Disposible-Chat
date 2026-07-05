package handler

import (
	"crypto/rand"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RoomStorage — in-memory, data hilang saat server restart
// ponytail: no DB, ephemeral rooms only; add Redis when persistence needed
var roomStore = &RoomStorage{rooms: make(map[string]*RoomData)}

type RoomStorage struct {
	mu    sync.RWMutex
	rooms map[string]*RoomData // key: roomCode (Contoh: "A8B2-9F")
}

type MemberCryptoData struct {
	EcdhPublicKey  string `json:"ecdhPublicKey"`  // Dari ecdh.ts -> untuk kunci enkripsi
	EcdsaPublicKey string `json:"ecdsaPublicKey"` // Dari ecdsa.ts -> untuk tanda tangan pesan
}

type RoomData struct {
	ID        string
	Code      string
	Members   map[string]MemberCryptoData // memberID -> sepasang Public Key
	Status    string
	CreatedAt time.Time
}

type CreateRoomRequest struct {
	EcdhPublicKey  string `json:"ecdhPublicKey"`
	EcdsaPublicKey string `json:"ecdsaPublicKey"`
}

type JoinRoomRequest struct {
	RoomCode       string `json:"roomCode"`
	EcdhPublicKey  string `json:"ecdhPublicKey"`
	EcdsaPublicKey string `json:"ecdsaPublicKey"`
}

func ok(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{"responseCode": "00", "responseMessage": "Berhasil", "data": data})
}

func fail(c *fiber.Ctx, status int, code, msg string) error {
	return c.Status(status).JSON(fiber.Map{"responseCode": code, "responseMessage": msg, "data": nil})
}

// CreateRoom — POST /v1/api/room/create
func CreateRoom(c *fiber.Ctx) error {
	var req CreateRoomRequest
	if err := c.BodyParser(&req); err != nil || req.EcdhPublicKey == "" || req.EcdsaPublicKey == "" {
		return fail(c, 400, "12", "ecdhPublicKey dan ecdsaPublicKey tidak boleh kosong")
	}

	roomID := uuid.New().String()
	memberID := uuid.New().String()
	code := roomCode()

	room := &RoomData{
		ID:   roomID,
		Code: code,
		Members: map[string]MemberCryptoData{
			memberID: {
				EcdhPublicKey:  req.EcdhPublicKey,
				EcdsaPublicKey: req.EcdsaPublicKey,
			},
		},
		Status:    "waiting",
		CreatedAt: time.Now(),
	}

	roomStore.mu.Lock()
	roomStore.rooms[code] = room
	roomStore.mu.Unlock()

	return ok(c, fiber.Map{
		"roomCode": code,
		"roomId":   roomID,
		"memberId": memberID,
		"status":   "waiting",
	})
}

// JoinRoom — POST /v1/api/room/join
func JoinRoom(c *fiber.Ctx) error {
	var req JoinRoomRequest
	if err := c.BodyParser(&req); err != nil || req.RoomCode == "" || req.EcdhPublicKey == "" || req.EcdsaPublicKey == "" {
		return fail(c, 400, "12", "roomCode, ecdhPublicKey, dan ecdsaPublicKey tidak boleh kosong")
	}

	roomStore.mu.Lock()
	defer roomStore.mu.Unlock()

	room, exists := roomStore.rooms[req.RoomCode]
	if !exists {
		return fail(c, 404, "11", "Room tidak ditemukan")
	}
	if len(room.Members) >= 2 {
		return fail(c, 409, "14", "Room sudah penuh")
	}

	// Ambil kunci publik milik pembuat room (User 1) untuk dilempar ke pendaftar baru (User 2)
	var peerCrypto MemberCryptoData
	for _, cryptoData := range room.Members {
		peerCrypto = cryptoData
		break
	}

	memberID := uuid.New().String()
	room.Members[memberID] = MemberCryptoData{
		EcdhPublicKey:  req.EcdhPublicKey,
		EcdsaPublicKey: req.EcdsaPublicKey,
	}
	room.Status = "active"

	return ok(c, fiber.Map{
		"roomId":        room.ID,
		"memberId":      memberID,
		"peerPublicKey": peerCrypto.EcdhPublicKey,  // Kunci ECDH lawan untuk enkripsi
		"peerSignKey":   peerCrypto.EcdsaPublicKey, // Kunci ECDSA lawan untuk verifikasi tanda tangan
		"status":        "active",
	})
}

// RoomStatus — GET /v1/api/room/:roomId/status
func RoomStatus(c *fiber.Ctx) error {
	roomID := c.Params("roomId")

	roomStore.mu.RLock()
	defer roomStore.mu.RUnlock()

	for _, room := range roomStore.rooms {
		if room.ID == roomID {
			return ok(c, fiber.Map{
				"roomId":      room.ID,
				"memberCount": len(room.Members),
				"status":      room.Status,
			})
		}
	}
	return fail(c, 404, "11", "Room tidak ditemukan")
}

// roomCode generates a XXXX-XX code using crypto/rand (not time-based)
func roomCode() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 6)
	rand.Read(b) // ponytail: ignoring error — rand.Read on crypto/rand never fails on modern OS
	for i, v := range b {
		b[i] = chars[int(v)%len(chars)]
	}
	return string(b[:4]) + "-" + string(b[4:])
}

