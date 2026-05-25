package ws

import (
	"context"

	"github.com/coder/websocket"
)

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
		c.Pool.Broadcast <- data
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
