package ws

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type WebSocketHandler struct {
	Manager *PoolManager
}

func NewWebSocketHandler(manager *PoolManager) *WebSocketHandler {
	return &WebSocketHandler{
		Manager: manager,
	}
}

func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	userID := "daya" //TODO

	roomID := ConversationID(r.URL.Query().Get("room"))
	if roomID == "" {
		http.Error(w, "Missing room parameter", http.StatusBadRequest)
		return
	}

	log.Printf("[WS] New connection attempt - User: %s, Room: %s\n", userID, roomID)

	options := &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	}
	conn, err := websocket.Accept(w, r, options)
	if err != nil {
		log.Printf("[WS] WebSocket upgrade failed: %v\n", err)
		return
	}

	log.Printf("[WS] WebSocket upgraded successfully - User: %s, Room: %s\n", userID, roomID)

	client := &Client{
		ID:   userID,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	// Use the request context only for initial connection setup
	err = h.Manager.HandleConnection(reqCtx, roomID, PoolTypeGroup, client)
	if err != nil {
		log.Printf("[WS] Failed to handle connection: %v\n", err)
		conn.Close(websocket.StatusPolicyViolation, err.Error())
		return
	}

	log.Printf("[WS] Client registered with pool - User: %s, Room: %s\n", userID, roomID)

	appCtx := h.Manager.appCtx
	go client.WritePump(appCtx)
	go client.ReadPump(appCtx, h.Manager, roomID)

	log.Printf("[WS] Goroutines spawned - WritePump and ReadPump - User: %s, Room: %s\n", userID, roomID)
}
