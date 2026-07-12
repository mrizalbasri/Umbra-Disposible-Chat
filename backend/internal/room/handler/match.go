package handler

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Menampung data user yang sedang mengantre di pool
type QueueItem struct {
	QueueID   string
	PublicKey string
	Nickname  string
}

// matchStore menggunakan in-memory storage
var matchStore = &MatchStorage{queues: make(map[string]QueueItem)}

type MatchStorage struct {
	mu     sync.RWMutex
	queues map[string]QueueItem // key: queueId -> value: objek QueueItem
}

// --- STRUCT REQUEST SESUAI API CONTRACT HALAMAN 11 ---

type MatchQueueRequest struct {
	PublicKey string `json:"publicKey"` // required [cite: 157]
	Nickname  string `json:"nickname"`  // optional [cite: 157]
}

// --- STRUCT RESPON KHUSUS UNTUK MEMAKSA URUTAN KAKU JSON (Halaman 12) ---

type MatchMatchedResponseData struct {
	Status        string `json:"status"`
	RoomID        string `json:"roomId"`
	PeerPublicKey string `json:"peerPublicKey"`
}

type MatchMatchedStandardResponse struct {
	ResponseMessage string                   `json:"responseMessage"`
	ResponseCode    string                   `json:"responseCode"`
	Data            MatchMatchedResponseData `json:"data"`
}

type MatchWaitingResponseData struct {
	Status  string `json:"status"`  // Urutan pertama di dalam objek data 
	QueueID string `json:"queueId"` // Urutan kedua di dalam objek data 
}

type MatchWaitingStandardResponse struct {
	ResponseMessage string                   `json:"responseMessage"`
	ResponseCode    string                   `json:"responseCode"`
	Data            MatchWaitingResponseData `json:"data"`
}

type CancelQueueResponseData struct {
	Status string `json:"status"`
}

type CancelQueueStandardResponse struct {
	ResponseMessage string                  `json:"responseMessage"`
	ResponseCode    string                  `json:"responseCode"`
	Data            CancelQueueResponseData `json:"data"`
}

// MatchQueue — POST /v1/api/match/queue
func MatchQueue(c *fiber.Ctx) error {
	var req MatchQueueRequest
	if err := c.BodyParser(&req); err != nil || req.PublicKey == "" {
		return fail(c, 400, "12", "publicKey tidak boleh kosong") // Sesuai Error Code 12 [cite: 17, 161]
	}

	matchStore.mu.Lock()
	defer matchStore.mu.Unlock()

	// 1. Cari apakah ada user lain di pool antrean
	var peerQueueID string
	var peerData QueueItem
	for id, item := range matchStore.queues {
		// Mengambil item pertama yang sedang menunggu di map
		peerQueueID = id
		peerData = item
		break
	}

	// 2. JIKA ADA USER LAIN: Langsung pasangkan (Matched)
	if peerQueueID != "" {
		// Hapus pasangan dari antrean pool karena sudah match
		delete(matchStore.queues, peerQueueID)

		// Buat Room ID baru berupa UUID
		roomID := uuid.New().String()

		// --- INTEGRASI KE ROOM STORE ---
		// Membuat roomCode otomatis menggunakan fungsi roomCode() dari room.go
		code := roomCode() 
		
		roomStore.mu.Lock()
		roomStore.rooms[code] = &RoomData{
			ID:   roomID,
			Code: code,
			Members: map[string]MemberCryptoData{
				uuid.New().String(): {PublicKey: peerData.PublicKey, Nickname: peerData.Nickname}, // User Pertama (Peer)
				uuid.New().String(): {PublicKey: req.PublicKey, Nickname: req.Nickname},          // User Kedua (Dirimu)
			},
			Status:    "active", // Status langsung aktif karena sudah berdua
			CreatedAt: c.Context().ConnTime(),
		}
		roomStore.mu.Unlock()

		// Memaksa cetakan JSON urut rapi (Halaman 12 dokumen kontrak - MATCHED) 
		return c.JSON(MatchMatchedStandardResponse{
			ResponseMessage: "Berhasil dipasangkan dengan user lain.",
			ResponseCode:    "00",
			Data: MatchMatchedResponseData{
				Status:        "matched",
				RoomID:        roomID,
				PeerPublicKey: peerData.PublicKey,
			},
		})
	}

	// 3. JIKA BELUM ADA USER LAIN: Masukkan user ini ke waiting pool
	newQueueID := uuid.New().String()
	matchStore.queues[newQueueID] = QueueItem{
		QueueID:   newQueueID,
		PublicKey: req.PublicKey,
		Nickname:  req.Nickname,
	}

	// Memaksa cetakan JSON urut rapi (Halaman 12 dokumen kontrak - WAITING dengan Code 17) 
	return c.JSON(MatchWaitingStandardResponse{
		ResponseMessage: "Masuk waiting pool, menunggu pasangan.",
		ResponseCode:    "17",
		Data: MatchWaitingResponseData{
			Status:  "waiting",
			QueueID: newQueueID,
		},
	})
}

// CancelQueue — DELETE /v1/api/match/queue/:queueId
func CancelQueue(c *fiber.Ctx) error {
	queueID := c.Params("queueId")

	matchStore.mu.Lock()
	defer matchStore.mu.Unlock()

	if _, exists := matchStore.queues[queueID]; !exists {
		return fail(c, 404, "11", "queueId tidak ditemukan atau sudah tidak aktif") // Sesuai Error Code 11 [cite: 17, 162]
	}

	delete(matchStore.queues, queueID)

	// Memaksa cetakan JSON urut rapi (Halaman 14 dokumen kontrak - CANCELLED) [cite: 164]
	return c.JSON(CancelQueueStandardResponse{
		ResponseMessage: "Pencarian Random Match berhasil dibatalkan.",
		ResponseCode:    "00",
		Data: CancelQueueResponseData{
			Status: "cancelled",
		},
	})
}