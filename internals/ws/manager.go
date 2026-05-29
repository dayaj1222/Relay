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

func (pm *PoolManager) HandleConnection(ctx context.Context, id ConversationID, poolType PoolType, client *Client) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.pools[id]; !exists {
		if len(pm.pools) >= pm.maxPools {
			return errors.New("server capacity reached: cannot create more pools")
		}

		newPool := newPool(poolType, id, func(roomID ConversationID) {
			pm.mu.Lock()
			delete(pm.pools, roomID)
			pm.mu.Unlock()
		})

		pm.pools[id] = newPool
		go newPool.start(pm.appCtx)
	}

	p := pm.pools[id]
	if len(p.engine.clients) >= pm.defaultMaxConn {
		return errors.New("pool capacity reached")
	}

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
