package handler

import (
	"github.com/gofiber/fiber/v2"
)

type RoomHandler struct{}

func NewRoomHandler() *RoomHandler {
	return &RoomHandler{}
}

// RegisterRoutes mendaftarkan seluruh endpoint REST API Umbra
func (h *RoomHandler) RegisterRoutes(router fiber.Router) {
	// 1. Grouping untuk rute /v1/api/room
	roomGroup := router.Group("/room")
	roomGroup.Post("/create", CreateRoom)        // POST /v1/api/room/create
	roomGroup.Post("/join", JoinRoom)            // POST /v1/api/room/join
	roomGroup.Get("/:roomId/status", RoomStatus) // GET /v1/api/room/:roomId/status

	// 2. Grouping untuk rute /v1/api/match
	matchGroup := router.Group("/match")
	matchGroup.Post("/queue", MatchQueue)             // POST /v1/api/match/queue
	matchGroup.Delete("/queue/:queueId", CancelQueue) // DELETE /v1/api/match/queue/:queueId
}
