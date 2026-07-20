package handler

import (
	"log"
	"sync"
	"time"
	"umbra-backend/internal/crypto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	queueTTL      = 5 * time.Minute
	cleanInterval = 60 * time.Second
	maxRoomMembers = 5
)

// Menampung data user yang sedang mengantre di pool
type QueueItem struct {
	QueueID     string
	RoomID      string
	PublicKey   string
	Nickname    string
	EnqueuedAt  time.Time
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

	// 1. Cari peer yang BERBEDA publicKey (fix: cegah self-match)
	var peerQueueID string
	var peerData QueueItem
	for id, item := range matchStore.queues {
		if item.PublicKey == req.PublicKey {
			continue // skip diri sendiri
		}
		peerQueueID = id
		peerData = item
		break
	}

	// 2. JIKA ADA USER LAIN: Langsung pasangkan (Matched)
	if peerQueueID != "" {
		// Hapus pasangan dari antrean pool karena sudah match
		delete(matchStore.queues, peerQueueID)
		matchStore.mu.Unlock() // fix: unlock sebelum acquire roomStore lock

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

	matchStore.queues[newQueueID] = QueueItem{
		QueueID:    newQueueID,
		RoomID:     roomID,
		PublicKey:  req.PublicKey,
		Nickname:   req.Nickname,
		EnqueuedAt: time.Now(),
	}
	matchStore.mu.Unlock() // fix: unlock sebelum acquire roomStore lock

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
		Status:     "waiting",
		Type:       "match",
		MaxMembers: 2,
		CreatedAt:  time.Now(),
	}
	roomStore.mu.Unlock()

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
	queue, exists := matchStore.queues[queueID]
	if !exists {
		matchStore.mu.Unlock()
		return fail(c, 404, "11", "queueId tidak ditemukan atau sudah tidak aktif")
	}
	delete(matchStore.queues, queueID)
	matchStore.mu.Unlock()

	DeleteRoom(queue.RoomID)

	return c.JSON(CancelQueueStandardResponse{
		ResponseMessage: "Pencarian Random Match berhasil dibatalkan.",
		ResponseCode:    "00",
		Data: CancelQueueResponseData{
			Status: "cancelled",
		},
	})
}

// StartQueueCleaner menjalankan goroutine yang membersihkan ghost queue entry secara berkala
// ponytail: simple ticker sweep, upgrade ke priority queue kalau queue bisa ribuan entry
func StartQueueCleaner() {
	go func() {
		ticker := time.NewTicker(cleanInterval)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			var expiredRoomIDs []string

			matchStore.mu.Lock()
			for id, item := range matchStore.queues {
				if now.Sub(item.EnqueuedAt) > queueTTL {
					expiredRoomIDs = append(expiredRoomIDs, item.RoomID)
					delete(matchStore.queues, id)
					log.Printf("🧹 Queue expired, dihapus: queueId=%s roomId=%s", id, item.RoomID)
				}
			}
			matchStore.mu.Unlock()

			for _, roomID := range expiredRoomIDs {
				DeleteRoom(roomID)
			}
		}
	}()
}
