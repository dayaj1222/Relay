package ws

import (
	"context"
	"errors"
)

func NewPoolManager(ctx context.Context, defaultMaxConn int, maxPools int) *PoolManager {
	return &PoolManager{
		pools:          make(map[ConversationID]*pool),
		defaultMaxConn: defaultMaxConn,
		maxPools:       maxPools,
		appCtx:         ctx,
	}
}

// Gets or creates pools when new connection occurs
func (pm *PoolManager) ensurePool(id ConversationID, poolType PoolType) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.pools[id]; exists {
		return nil
	}

	if len(pm.pools) >= pm.maxPools {
		return errors.New("server capacity reached: cannot create more pools")
	}

	// Self deletion fucntion is passed to each pool
	newPool := newPool(poolType, id, func(roomID ConversationID) {
		pm.mu.Lock()
		delete(pm.pools, roomID)
		pm.mu.Unlock()
	})

	pm.pools[id] = newPool
	go newPool.start(pm.appCtx)

	return nil
}

func (pm *PoolManager) HandleConnection(rewCtx context.Context, id ConversationID, poolType PoolType, client *Client) error {
	if err := pm.ensurePool(id, poolType); err != nil {
		return err
	}

	pm.mu.Lock()
	p := pm.pools[id]

	if len(p.engine.clients) >= pm.defaultMaxConn {
		pm.mu.Unlock()
		return errors.New("pool capacity reached")
	}
	pm.mu.Unlock()

	p.engine.register <- client
	return nil
}

func (pm *PoolManager) DisconnectClient(id ConversationID, client *Client) {
	pm.mu.RLock()
	p, exists := pm.pools[id]
	pm.mu.RUnlock()

	if exists {
		p.engine.unregister <- client
	}
}

func (pm *PoolManager) BroadcastToPool(id ConversationID, message []byte) error {
	pm.mu.RLock()
	p, exists := pm.pools[id]
	pm.mu.RUnlock()

	if !exists {
		return errors.New("cannot broadcast: target pool is not active in memory")
	}

	p.engine.broadcast <- message
	return nil
}

func (pm *PoolManager) KickClient(id ConversationID, clientID string) error {
	pm.mu.RLock()
	p, exists := pm.pools[id]
	pm.mu.RUnlock()

	if !exists {
		return errors.New("target pool is not active")
	}

	if kicked := p.kickClient(clientID); !kicked {
		return errors.New("client not found in the specified pool")
	}

	return nil
}
