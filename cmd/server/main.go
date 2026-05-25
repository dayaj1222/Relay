package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"relay/internals/rooms"
	"relay/internals/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load env vars

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println(".env not found")
	}

	port := os.Getenv("PORT")
	clientURL := os.Getenv("CLIENT_URL")

	pool := ws.NewPool()

	bgCtx := context.Background()
	go pool.Start(bgCtx)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			clientURL,
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
			"Upgrade", "Connection", "Sec-WebSocket-Key", "Sec-WebSocket-Version",
		},
		ExposeHeaders: []string{
			"Upgrade", "Connection",
		},
	}))

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Rooms
	r.POST("/api/rooms", rooms.CreateRoom)
	r.GET("/api/rooms", rooms.GetRooms)

	r.GET("/api/ws", func(ctx *gin.Context) {
		ws.AcceptConnection(pool, ctx)
	})

	log.Println("Server running on %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
