package ws

import (
	"context"
	"encoding/json"
	"log"
	"relay/internals/messages"

	"github.com/coder/websocket"
)

const maxMessageSize = 1 * 1024 * 1024

const (
	logReadPumpVerbose  = false
	logWritePumpVerbose = false
)

func (c *Client) ReadPump(ctx context.Context, pm *PoolManager, id ConversationID) {

	log.Printf("[ReadPump %s] Started\n", c.ID)
	defer func() {
		pm.DisconnectClient(id, c)
	}()

	for {
		_, data, err := c.Conn.Read(ctx)
		if err != nil {
			log.Printf("[ReadPump %s] Read error: %v\n", c.ID, err)
			break
		}
		if logReadPumpVerbose {
			log.Printf("[ReadPump %s] Received message (%d bytes)\n", c.ID, len(data))
		}
		if len(data) > maxMessageSize {
			log.Printf("[ReadPump %s] Message too large (%d bytes)\n", c.ID, len(data))
			break
		}
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
				return
			}
			if logWritePumpVerbose {
				log.Printf("[WritePump %s] Sending message (%d bytes)\n", c.ID, len(message))
			}
			err := c.Conn.Write(ctx, websocket.MessageText, message)
			if err != nil {
				log.Printf("[WritePump %s] Write error: %v\n", c.ID, err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
