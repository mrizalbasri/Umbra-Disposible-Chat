package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"

	roomhandler "umbra-backend/internal/room/handler"
	"umbra-backend/internal/hub"
)

func main() {
	h := hub.New()
	go h.Run()

	app := fiber.New(fiber.Config{AppName: "Umbra Backend v1.0.0"})

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} ${method} ${path} ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
		AllowMethods: "GET,POST,DELETE",
		AllowHeaders: "Content-Type",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "app": "Umbra", "version": "1.0.0"})
	})

	api := app.Group("/v1/api")

	room := api.Group("/room")
	room.Post("/create", roomhandler.CreateRoom)
	room.Post("/join", roomhandler.JoinRoom)
	room.Get("/:roomId/status", roomhandler.RoomStatus)

	match := api.Group("/match")
	match.Post("/queue", roomhandler.MatchQueue)
	match.Delete("/queue/:queueId", roomhandler.CancelQueue)

	// WebSocket upgrade guard
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws", websocket.New(h.HandleWS))

	port := getEnv("PORT", "8080")
	log.Printf("🌑 Umbra server starting on port %s", port)
	log.Printf("📡 WebSocket: ws://localhost:%s/ws", port)
	log.Printf("🔌 REST API: http://localhost:%s/v1/api", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
