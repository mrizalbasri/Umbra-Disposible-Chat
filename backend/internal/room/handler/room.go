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
	rooms map[string]*RoomData // key: roomCode
}

type RoomData struct {
	ID        string
	Code      string
	Members   map[string]string // memberID -> publicKey
	Status    string
	CreatedAt time.Time
}

type CreateRoomRequest struct {
	PublicKey string `json:"publicKey"`
}

type JoinRoomRequest struct {
	RoomCode  string `json:"roomCode"`
	PublicKey string `json:"publicKey"`
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
	if err := c.BodyParser(&req); err != nil || req.PublicKey == "" {
		return fail(c, 400, "12", "publicKey tidak boleh kosong")
	}

	roomID := uuid.New().String()
	memberID := uuid.New().String()
	code := roomCode()

	room := &RoomData{
		ID:        roomID,
		Code:      code,
		Members:   map[string]string{memberID: req.PublicKey},
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
	if err := c.BodyParser(&req); err != nil || req.RoomCode == "" || req.PublicKey == "" {
		return fail(c, 400, "12", "roomCode dan publicKey tidak boleh kosong")
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

	// grab creator's public key for ECDH exchange
	var peerPublicKey string
	for _, pk := range room.Members {
		peerPublicKey = pk
		break
	}

	memberID := uuid.New().String()
	room.Members[memberID] = req.PublicKey
	room.Status = "active"

	return ok(c, fiber.Map{
		"roomId":        room.ID,
		"memberId":      memberID,
		"peerPublicKey": peerPublicKey,
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
