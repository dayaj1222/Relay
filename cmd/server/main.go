package main

import (
	"log"
	"net/http"
	"os"
	"relay/internals/rooms"

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

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			clientURL,
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
	}))

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	r.POST("/rooms", rooms.CreateRoom)
	r.GET("/rooms", rooms.GetRooms)

	r.Run(":" + port)
}
