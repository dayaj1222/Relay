package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"relay/internals/conversations"
	"relay/internals/db"
	"relay/internals/messages"
	"relay/internals/middleware"
	"relay/internals/uploads"
	"relay/internals/users"
	"relay/internals/ws"
	"syscall"
	"time"

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
	connectionString := os.Getenv("POSTGRES_URL")

	database, err := db.Connect(connectionString)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer database.Close()

	// Handle graceful termination signals for background workers
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize WebSocket Infrastructure
	manager := ws.NewPoolManager(ctx, 256, 1000)
	wsHandler := ws.NewWebSocketHandler(manager)

	userRepo := users.NewRepository(database)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService)

	// Initialize Conversations Infrastructure
	convStore := conversations.NewStore(database)
	convService := conversations.NewService(convStore)
	convHandler := conversations.NewHTTPHandler(convService)

	// Initialize Messages Infrastructure
	msgStore := messages.NewStore(database)
	msgHandler := messages.NewHTTPHandler(msgStore, manager)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
			"Upgrade", "Connection", "Sec-WebSocket-Key", "Sec-WebSocket-Version",
		},
		ExposeHeaders:    []string{"Upgrade", "Connection"},
		AllowCredentials: true,
	}))

	// Public Routes
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	r.POST("/api/register", userHandler.Register)
	r.POST("/api/login", userHandler.Login)

	// Protected Routes Group
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired(userService))
	{

		protected.GET("/users", userHandler.GetUserByUsername)
		// Conversations (DMs + Group management)
		protected.POST("/conversations/dm", convHandler.CreateDM)
		protected.POST("/conversations/group", convHandler.CreateGroup)
		protected.GET("/conversations", convHandler.ListConversations)
		protected.GET("/conversations/:id", convHandler.GetConversation)

		// Messages
		protected.POST("/conversations/:id/messages", msgHandler.SendMessage)
		protected.GET("/conversations/:id/messages", msgHandler.GetMessages)
		protected.GET("/conversations/:id/messages/recent", msgHandler.GetRecentMessages)

		// Upload (Moved inside protection)
		protected.POST("/upload", uploads.UploadHandler)

		// Websocket (Moved inside protection so auth context runs)
		protected.GET("/ws", func(c *gin.Context) {
			// Middleware already validated token and set "user_id" into contexts
			wsHandler.ServeHTTP(c.Writer, c.Request)
		})
	}

	r.Static("/uploads", UploadDir)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Server running on %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to run server: %v", err)
	}
}
