package ws

import (
	"context"
	"errors"
	"log"
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
	if _, exists := pm.pools[id]; !exists {
		if len(pm.pools) >= pm.maxPools {
			pm.mu.Unlock()
			return errors.New("server capacity reached: cannot create more pools")
		}
		poolInstance := newPool(poolType, id)
		pm.pools[id] = poolInstance
		log.Printf("[MANAGER] Pool created and stored: %s, total pools: %d", id, len(pm.pools))
		go poolInstance.start(pm.appCtx)
	}
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
	log.Printf("[BROADCAST] Looking for pool '%s', found: %v, all pools: %v", id, exists, func() []string {
		keys := make([]string, 0, len(pm.pools))
		for k := range pm.pools {
			keys = append(keys, string(k))
		}
		return keys
	}())
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
