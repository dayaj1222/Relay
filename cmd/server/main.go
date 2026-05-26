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

	// FIX 1: Rename the local instance to 'database' to prevent shadowing package 'db'
	database, err := db.Connect(connectionString)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer database.Close()

	// FIX 2: Handle graceful termination signals for background workers
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool := ws.NewPool()
	go pool.Start(ctx) // Shuts down loop automatically when context closes

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

	// Websocket
	r.GET("/api/ws", func(ctx *gin.Context) {
		ws.AcceptConnection(pool, ctx)
	})

	// Upload
	r.POST("/api/upload", uploads.UploadHandler)
	r.Static("/uploads", UploadDir)

	log.Printf("Server running on %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
