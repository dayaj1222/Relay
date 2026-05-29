package ws

import (
	"context"
	"log"

	"github.com/coder/websocket"
)

func newPool(poolType PoolType, conversationID ConversationID, onEmpty func(ConversationID)) *pool {
	h := &hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	return &pool{
		id:      conversationID,
		engine:  h,
		pType:   poolType,
		onEmpty: onEmpty,
	}
}

func (p *pool) start(ctx context.Context) {
	for {
		select {
		case client := <-p.engine.register:
			log.Printf("[POOL %s] Client registered: %s (total: %d)\n", p.id, client.ID, len(p.engine.clients)+1)
			p.engine.clients[client] = true

		case client := <-p.engine.unregister:
			if _, ok := p.engine.clients[client]; ok {
				log.Printf("[POOL %s] Client unregistered: %s (remaining: %d)\n", p.id, client.ID, len(p.engine.clients)-1)
				p.removeClient(client)
				if len(p.engine.clients) == 0 {
					log.Printf("[POOL %s] Pool empty, shutting down\n", p.id)
					p.onEmpty(p.id)
					return
				}
			}

		case message := <-p.engine.broadcast:
			for client := range p.engine.clients {
				select {
				case client.Send <- message:
				default:
					p.removeClient(client)
					log.Printf("[POOL %s] Dropped slow client: %s\n", p.id, client.ID)
				}
			}

		case <-ctx.Done():
			log.Printf("[POOL %s] Context cancelled, shutting down pool with %d clients\n", p.id, len(p.engine.clients))
			for client := range p.engine.clients {
				p.removeClient(client)
			}
			p.onEmpty(p.id)
			return
		}
	}
}

func (p *pool) removeClient(client *Client) {
	delete(p.engine.clients, client)
	close(client.Send)
	if client.Conn != nil {
		client.Conn.Close(websocket.StatusNormalClosure, "connection closed by pool")
	}
}

func (p *pool) kickClient(clientID string) bool {
	for client := range p.engine.clients {
		if client.ID == clientID {
			p.removeClient(client)
			return true
		}
	}
	return false
}
