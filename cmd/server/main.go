package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"relay/internals/db"
	"relay/internals/rooms"
	"relay/internals/uploads"
	"relay/internals/users"
	"relay/internals/ws"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const UploadDir = "./uploads"

func main() {
	// Load env vars
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println(".env file not found, reading from system environment")
	}

	port := os.Getenv("PORT")
	clientURL := os.Getenv("CLIENT_URL")
	connectionString := os.Getenv("POSTGRES_URL")

	database, err := db.Connect(connectionString)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer database.Close()

	// Handle graceful termination signals for background workers
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize WebSocket Infrastructure (e.g., max 256 connections per room, 1000 max active pools)
	manager := ws.NewPoolManager(ctx, 256, 1000)
	wsHandler := ws.NewWebSocketHandler(manager)

	userRepo := users.NewRepository(database)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{clientURL},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
			"Upgrade", "Connection", "Sec-WebSocket-Key", "Sec-WebSocket-Version",
		},
		ExposeHeaders: []string{"Upgrade", "Connection"},
	}))

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// User
	r.POST("/api/register", userHandler.Register)
	r.POST("/api/login", userHandler.Login)

	// Rooms
	r.POST("/api/rooms", rooms.CreateRoom)
	r.GET("/api/rooms", rooms.GetRooms)

	// Websocket - Integrated with clean architecture
	// (Ensure an authentication middleware runs before this endpoint to set "user_id" into the gin context)
	r.GET("/api/ws", func(c *gin.Context) {
		// 1. Get user identity from auth middleware context strings
		userID := c.GetString("user_id")
		if userID == "" {
			// Fallback/Safety block if your middleware naming differs
			userID = "anonymous_fallback"
		}

		// 2. Wrap it inside standard context.Context to cross boundary cleanly
		reqCtx := context.WithValue(c.Request.Context(), "user_id", userID)
		c.Request = c.Request.WithContext(reqCtx)

		// 3. Directly run the standard interface handler
		wsHandler.ServeHTTP(c.Writer, c.Request)
	})

	// Upload
	r.POST("/api/upload", uploads.UploadHandler)
	r.Static("/uploads", UploadDir)

	log.Printf("Server running on %s", port)

	// Pass the lifecycle context to handle smooth termination sequences
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
