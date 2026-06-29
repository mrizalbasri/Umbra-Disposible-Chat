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
	queueID := uuid.New().String()

	matchStore.mu.Lock()
	matchStore.queues[queueID] = "waiting"
	matchStore.mu.Unlock()

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
