package handler

import (
	"crypto/rand"
	"log"
	"regexp"
	"sync"
	"time"
	"umbra-backend/internal/crypto"

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
	PublicKey string `json:"publicKey"` // Sesuai API Contract: Single string base64 ECDH
	Nickname  string `json:"nickname"`  // Sesuai API Contract: Penampung display name opsional
}

type RoomData struct {
	ID         string
	Code       string
	Members    map[string]MemberCryptoData // memberID -> sepasang Public Key & Nickname
	Status     string
	Type       string                      // "custom" atau "match"
	MaxMembers int                         // 5 untuk custom, 2 untuk match
	CreatedAt  time.Time
}

type CreateRoomRequest struct {
	PublicKey string `json:"publicKey"` // Sesuai API Contract Halaman 5
	Nickname  string `json:"nickname"`  // Sesuai API Contract Halaman 5
}

type JoinRoomRequest struct {
	RoomCode  string `json:"roomCode"`  // Sesuai API Contract Halaman 7
	PublicKey string `json:"publicKey"` // Sesuai API Contract Halaman 7
	Nickname  string `json:"nickname"`  // Sesuai API Contract Halaman 7
}

// --- STRUCT RESPON KHUSUS UNTUK MEMAKSA URUTAN KAKU JSON (STRICT SORTING) ---

type CreateRoomResponseData struct {
	RoomCode string `json:"roomCode"`
	RoomID   string `json:"roomId"`
	Status   string `json:"status"`
}

type CreateRoomStandardResponse struct {
	ResponseMessage string                 `json:"responseMessage"`
	ResponseCode    string                 `json:"responseCode"`
	Data            CreateRoomResponseData `json:"data"`
}

type JoinRoomResponseData struct {
	RoomID        string `json:"roomId"`
	PeerPublicKey string `json:"peerPublicKey"`
	Status        string `json:"status"`
}

type JoinRoomStandardResponse struct {
	ResponseMessage string               `json:"responseMessage"`
	ResponseCode    string               `json:"responseCode"`
	Data            JoinRoomResponseData `json:"data"`
}

type RoomStatusResponseData struct {
	RoomID      string `json:"roomId"`
	MemberCount int    `json:"memberCount"`
	Status      string `json:"status"`
}

type RoomStatusStandardResponse struct {
	ResponseMessage string                 `json:"responseMessage"`
	ResponseCode    string                 `json:"responseCode"`
	Data            RoomStatusResponseData `json:"data"`
}

// fail formatter — disesuaikan agar format output JSON seragam sesuai spesifikasi aplikasi
func fail(c *fiber.Ctx, status int, code, msg string) error {
	return c.Status(status).JSON(fiber.Map{
		"responseMessage": msg,
		"responseCode":    code,
		"data":            nil,
	})
}

// CreateRoom — POST /v1/api/room/create
func CreateRoom(c *fiber.Ctx) error {
	var req CreateRoomRequest
	if err := c.BodyParser(&req); err != nil || req.PublicKey == "" {
		return fail(c, 400, "12", "publicKey tidak boleh kosong") // Sesuai Error Code 12
	}
	if !crypto.ValidatePublicKey(req.PublicKey) {
		return fail(c, 400, "15", "Format publicKey tidak valid") // Sesuai Error Code 15
	}

	roomID := uuid.New().String()
	memberID := uuid.New().String()
	code := roomCode()

	room := &RoomData{
		ID:   roomID,
		Code: code,
		Members: map[string]MemberCryptoData{
			memberID: {
				PublicKey: req.PublicKey,
				Nickname:  req.Nickname,
			},
		},
		Status:     "waiting", // Status awal: waiting
		Type:       "custom",  // Custom 1-on-1 chat
		MaxMembers: 2,         // Maksimal 2 anggota (ECDH pairwise E2EE)
		CreatedAt:  time.Now(),
	}

	roomStore.mu.Lock()
	roomStore.rooms[code] = room
	roomStore.mu.Unlock()

	// Memaksa cetakan JSON urut rapi dari atas ke bawah (Halaman 6 dokumen kontrak)
	return c.JSON(CreateRoomStandardResponse{
		ResponseMessage: "Room berhasil dibuat.",
		ResponseCode:    "00",
		Data: CreateRoomResponseData{
			RoomCode: code,
			RoomID:   roomID,
			Status:   "waiting",
		},
	})
}

// JoinRoom — POST /v1/api/room/join
func JoinRoom(c *fiber.Ctx) error {
	var req JoinRoomRequest
	if err := c.BodyParser(&req); err != nil || req.RoomCode == "" || req.PublicKey == "" {
		return fail(c, 400, "12", "roomCode/publicKey tidak boleh kosong") // Sesuai Error Code 12
	}
	if matched, _ := regexp.MatchString(`^(?i)[A-Z0-9]{4}-[A-Z0-9]{2}$`, req.RoomCode); !matched {
		return fail(c, 400, "13", "Format kode room tidak valid. Gunakan format XXXX-XX") // Sesuai Error Code 13
	}
	if !crypto.ValidatePublicKey(req.PublicKey) {
		return fail(c, 400, "15", "Format publicKey tidak valid") // Sesuai Error Code 15
	}

	roomStore.mu.Lock()
	defer roomStore.mu.Unlock()

	room, exists := roomStore.rooms[req.RoomCode]
	if !exists {
		return fail(c, 404, "11", "Room tidak ditemukan") // Sesuai Error Code 11
	}
	if room.Type == "match" {
		return fail(c, 403, "18", "Kode room ini khusus untuk Random Match")
	}
	maxCap := room.MaxMembers
	if maxCap == 0 {
		maxCap = 2
	}
	if len(room.Members) >= maxCap {
		return fail(c, 409, "14", "Room sudah penuh, maksimal 2 anggota") // Sesuai Error Code 14
	}

	// Ambil data publicKey milik user pertama (peer) untuk proses key exchange
	var peerPublicKey string
	for _, cryptoData := range room.Members {
		peerPublicKey = cryptoData.PublicKey
		break
	}

	memberID := uuid.New().String()
	room.Members[memberID] = MemberCryptoData{
		PublicKey: req.PublicKey,
		Nickname:  req.Nickname,
	}
	room.Status = "active" // Status berubah menjadi active

	// Memaksa cetakan JSON urut rapi dari atas ke bawah (Halaman 8 dokumen kontrak)
	return c.JSON(JoinRoomStandardResponse{
		ResponseMessage: "Berhasil bergabung ke room.",
		ResponseCode:    "00",
		Data: JoinRoomResponseData{
			RoomID:        room.ID,
			PeerPublicKey: peerPublicKey,
			Status:        "active",
		},
	})
}

// RoomStatus — GET /v1/api/room/:roomId/status
func RoomStatus(c *fiber.Ctx) error {
	roomID := c.Params("roomId")

	roomStore.mu.RLock()
	defer roomStore.mu.RUnlock()

	for _, room := range roomStore.rooms {
		if room.ID == roomID {
			// Memaksa cetakan JSON urut rapi dari atas ke bawah (Halaman 10 dokumen kontrak)
			return c.JSON(RoomStatusStandardResponse{
				ResponseMessage: "Data room ditemukan.",
				ResponseCode:    "00",
				Data: RoomStatusResponseData{
					RoomID:      room.ID,
					MemberCount: len(room.Members),
					Status:      room.Status,
				},
			})
		}
	}

	// DIUBAH DI SINI: Disamakan persis dengan "Penjelasan Response" pada dokumen contract
	return fail(c, 404, "11", "Room tidak ditemukan atau sudah dihapus dari memori") // Sesuai Error Code 11
}

// GetRoomType mengembalikan tipe room ("custom" atau "match") berdasarkan roomID
func GetRoomType(roomID string) string {
	roomStore.mu.RLock()
	defer roomStore.mu.RUnlock()

	for _, room := range roomStore.rooms {
		if room.ID == roomID {
			return room.Type
		}
	}
	return "custom"
}

// DeleteRoom menghapus ruangan dari roomStore berdasarkan roomID (untuk mencegah memory leak)
func DeleteRoom(roomID string) {
	roomStore.mu.Lock()
	defer roomStore.mu.Unlock()

	for code, room := range roomStore.rooms {
		if room.ID == roomID {
			delete(roomStore.rooms, code)
			log.Printf("🗑️ Room %s (code: %s) dihapus total dari memori", roomID, code)
			break
		}
	}
}

// StartRoomCleaner menjalankan goroutine yang membersihkan room yang sudah kosong atau expired secara berkala
func StartRoomCleaner() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			var expiredCodes []string

			roomStore.mu.Lock()
			for code, room := range roomStore.rooms {
				// Hapus jika room sudah berumur lebih dari 30 menit (ephemeral TTL)
				if now.Sub(room.CreatedAt) > 30*time.Minute {
					expiredCodes = append(expiredCodes, code)
				}
			}

			for _, code := range expiredCodes {
				delete(roomStore.rooms, code)
				log.Printf("🧹 Room expired (>30 menit) dihapus: code=%s", code)
			}
			roomStore.mu.Unlock()
		}
	}()
}

// roomCode generates a XXXX-XX code using crypto/rand (not time-based)
func roomCode() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 6)
	_, _ = rand.Read(b) // ponytail: ignoring error — rand.Read on crypto/rand never fails on modern OS
	for i, v := range b {
		b[i] = chars[int(v)%len(chars)]
	}
	return string(b[:4]) + "-" + string(b[4:])
}
