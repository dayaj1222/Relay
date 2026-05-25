package ws

import (
	"context"
)

type Pool struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewPool() *Pool {
	return &Pool{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (p *Pool) Start(ctx context.Context) {
	for {
		select {
		case client := <-p.Register:
			p.Clients[client] = true
		case client := <-p.Unregister:
			if _, ok := p.Clients[client]; ok {
				delete(p.Clients, client)
				close(client.Send)
			}
		case message := <-p.Broadcast:
			for client := range p.Clients {
				select {
				case client.Send <- message:
				default:
					delete(p.Clients, client)
					close(client.Send)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
