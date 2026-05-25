package ws

import (
	"context"
	"encoding/json"
	"log"
	"relay/internals/utils"

	"github.com/coder/websocket"
)

type MessageType int

const (
	TypeText MessageType = iota
	TypeVideo
	TypeAudio
	TypeDocument
)

type Message struct {
	ID         int         `json:"id,omitempty" db:"id"`
	SenderID   string      `json:"senderId" db:"sender_id"`
	ReceiverID string      `json:"receiverId,omitempty" db:"-"`
	PoolID     string      `json:"poolId" db:"pool_id"`
	Type       int         `json:"type" db:"type"`
	Content    utils.JSONB `json:"content" db:"content"`
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	Send chan []byte
}

func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close(websocket.StatusNormalClosure, "disconnected")
	}()

	for {
		_, data, err := c.Conn.Read(ctx)
		if err != nil {
			break
		}
		var msg Message
		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Printf("error unmarshalling message: %v\n", err)
			continue
		}

		// Put the known sender id and reencode for broadcasting
		msg.SenderID = c.ID
		secureData, err := json.Marshal(msg)
		if err != nil {
			log.Printf("error marshalling secure message: %v\n", err)
			continue
		}

		c.Pool.Broadcast <- secureData
	}
}

func (c *Client) WritePump(ctx context.Context) {
	defer func() {
		c.Conn.Close(websocket.StatusNormalClosure, "internal error")
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				return
			}
			err := c.Conn.Write(ctx, websocket.MessageText, message)
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
