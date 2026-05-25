package ws

import (
	"context"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
)

func AcceptConnection(pool *Pool, ctx *gin.Context) {
	userID := ctx.Query("userId")

	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	options := &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	}

	conn, err := websocket.Accept(ctx.Writer, ctx.Request, options)
	if err != nil {
		log.Printf("conn accept error: %v", err)
		return
	}

	client := &Client{
		ID:   userID,
		Conn: conn,
		Pool: pool,
		Send: make(chan []byte, 256),
	}
	pool.Register <- client

	bgCtx := context.Background()

	go client.WritePump(bgCtx)
	go client.ReadPump(bgCtx)
}
