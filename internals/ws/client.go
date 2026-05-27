package ws

import (
	"context"
	"encoding/json"
	"log"
	"relay/internals/messages"

	"github.com/coder/websocket"
)

func (c *Client) ReadPump(ctx context.Context, pm *PoolManager, id ConversationID) {

	log.Printf("[ReadPump %s] Started\n", c.ID)
	defer func() {
		log.Printf("[ReadPump %s] Ending, disconnecting from pool\n", c.ID)
		pm.DisconnectClient(id, c)
	}()

	for {
		log.Printf("[ReadPump %s] Waiting for message...\n", c.ID)
		_, data, err := c.Conn.Read(ctx)
		if err != nil {
			log.Printf("[ReadPump %s] Read error: %v\n", c.ID, err)
			break
		}
		log.Printf("[ReadPump %s] Received message (%d bytes)\n", c.ID, len(data))
		var msg messages.Message
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

		if err := pm.BroadcastToPool(id, secureData); err != nil {
			log.Printf("broadcast error: %v\n", err)
		}
	}
}
func (c *Client) WritePump(ctx context.Context) {
	log.Printf("[WritePump %s] Started\n", c.ID)
	defer log.Printf("[WritePump %s] Ended\n", c.ID)

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				log.Printf("[WritePump %s] Send channel closed\n", c.ID)
				return
			}
			log.Printf("[WritePump %s] Sending message (%d bytes)\n", c.ID, len(message))
			err := c.Conn.Write(ctx, websocket.MessageText, message)
			if err != nil {
				log.Printf("[WritePump %s] Write error: %v\n", c.ID, err)
				return
			}
		case <-ctx.Done():
			log.Printf("[WritePump %s] Context done\n", c.ID)
			return
		}
	}
}
