package ws

import (
	"context"
	"sync"

	"github.com/coder/websocket"
)

type ConversationID string

type PoolType int

const (
	PoolTypeIndividual PoolType = iota
	PoolTypeGroup
)

// Client represents an active network connection.
// Fields are kept public so your read/write pump loops can access them.
type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

// hub handles raw channel communication.
type hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// pool bridges the metadata with the communication engine.
type pool struct {
	id      ConversationID
	pType   PoolType
	engine  *hub
	onEmpty func(ConversationID)
}

// PoolManager is the single public entry point for the package.
type PoolManager struct {
	pools          map[ConversationID]*pool
	defaultMaxConn int
	maxPools       int
	mu             sync.RWMutex
	appCtx         context.Context
}
