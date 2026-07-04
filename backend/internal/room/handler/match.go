package handler

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ponytail: match queue is in-memory, no persistence; add Redis/DB when matchmaking needs durability
var matchStore = &MatchStorage{queues: make(map[string]string)}

type MatchStorage struct {
	mu     sync.RWMutex
	queues map[string]string // queueId -> status
}

// MatchQueue — POST /v1/api/match/queue
func MatchQueue(c *fiber.Ctx) error {
    matchStore.mu.Lock()
    defer matchStore.mu.Unlock()

    // 1. Cari apakah ada user lain yang statusnya "waiting"
    var peerQueueID string
    for id, status := range matchStore.queues {
        if status == "waiting" {
            peerQueueID = id
            break
        }
    }

    // 2. Jika ada user lain yang siap dipasangkan
    if peerQueueID != "" {
        // Ubah status antrean lama menjadi matched (atau langsung hapus)
        delete(matchStore.queues, peerQueueID)
        
        // Buat Room ID baru untuk mereka berdua
        roomID := uuid.New().String()
        
        // (Opsional) Di sini nanti kamu panggil fungsi dari roomhandler untuk mendaftarkan roomID ini

        return ok(c, fiber.Map{
            "status":  "matched",
            "roomId":  roomID,
            "peerId":  peerQueueID,
        })
    }

    // 3. Jika tidak ada orang lain di antrean, masukkan user ini ke daftar "waiting"
    queueID := uuid.New().String()
    matchStore.queues[queueID] = "waiting"

    return ok(c, fiber.Map{"status": "waiting", "queueId": queueID})
}

// CancelQueue — DELETE /v1/api/match/queue/:queueId
func CancelQueue(c *fiber.Ctx) error {
	queueID := c.Params("queueId")

	matchStore.mu.Lock()
	defer matchStore.mu.Unlock()

	if _, exists := matchStore.queues[queueID]; !exists {
		return fail(c, 404, "11", "Queue tidak ditemukan")
	}

	delete(matchStore.queues, queueID)
	return ok(c, fiber.Map{"status": "cancelled"})
}
