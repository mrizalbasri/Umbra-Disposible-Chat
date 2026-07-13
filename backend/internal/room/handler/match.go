package handler

import (
	"sync"
	"time"
	"umbra-backend/internal/crypto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Menampung data user yang sedang mengantre di pool
type QueueItem struct {
	QueueID   string
	RoomID    string
	PublicKey string
	Nickname  string
}

// matchStore menggunakan in-memory storage
var matchStore = &MatchStorage{
	queues: make(map[string]QueueItem),
}

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
	if !crypto.ValidatePublicKey(req.PublicKey) {
		return fail(c, 400, "15", "Format publicKey tidak valid") // Sesuai Error Code 15
	}

	matchStore.mu.Lock()
	defer matchStore.mu.Unlock()

	// 1. Cari apakah ada user lain di pool antrean
	var peerQueueID string
	var peerData QueueItem
	for id, item := range matchStore.queues {
		peerQueueID = id
		peerData = item
		break
	}

	// 2. JIKA ADA USER LAIN: Langsung pasangkan (Matched)
	if peerQueueID != "" {
		// Hapus pasangan dari antrean pool karena sudah match
		delete(matchStore.queues, peerQueueID)

		// Gunakan room yang sudah dibuat saat user pertama masuk waiting
		roomID := peerData.RoomID

		// Update room yang sudah ada menjadi active dan tambahkan member kedua
		roomStore.mu.Lock()
		for _, room := range roomStore.rooms {
			if room.ID == roomID {
				room.Members[uuid.New().String()] = MemberCryptoData{
					PublicKey: req.PublicKey,
					Nickname:  req.Nickname,
				}
				room.Status = "active"
				break
			}
		}
		roomStore.mu.Unlock()

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
	roomID := newQueueID // room sementara untuk menunggu peer (dipakai juga saat matched)
	code := roomCode()

	roomStore.mu.Lock()
	roomStore.rooms[code] = &RoomData{
		ID:   roomID,
		Code: code,
		Members: map[string]MemberCryptoData{
			uuid.New().String(): {
				PublicKey: req.PublicKey,
				Nickname:  req.Nickname,
			},
		},
		Status:    "waiting",
		CreatedAt: time.Now(),
	}
	roomStore.mu.Unlock()

	matchStore.queues[newQueueID] = QueueItem{
		QueueID:   newQueueID,
		RoomID:    roomID,
		PublicKey: req.PublicKey,
		Nickname:  req.Nickname,
	}

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

	queue, exists := matchStore.queues[queueID]
	if !exists {
		return fail(c, 404, "11", "queueId tidak ditemukan atau sudah tidak aktif")
	}

	delete(matchStore.queues, queueID)
	DeleteRoom(queue.RoomID)

	return c.JSON(CancelQueueStandardResponse{
		ResponseMessage: "Pencarian Random Match berhasil dibatalkan.",
		ResponseCode:    "00",
		Data: CancelQueueResponseData{
			Status: "cancelled",
		},
	})
}
