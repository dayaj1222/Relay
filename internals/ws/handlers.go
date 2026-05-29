package ws

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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

var conversationIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

func isValidConversationID(id ConversationID) bool {
	return conversationIDPattern.MatchString(string(id))
}

// getPoolType determines the pool type based on conversation ID format.
// "dm_X_Y" format → PoolTypeIndividual
// "room_*" format → PoolTypeGroup
func getPoolType(convID ConversationID) PoolType {
	if strings.HasPrefix(string(convID), "dm_") {
		return PoolTypeIndividual
	}
	return PoolTypeGroup
}

func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	convID := ConversationID(r.URL.Query().Get("room"))
	if convID == "" {
		http.Error(w, "Missing room parameter", http.StatusBadRequest)
		return
	}

	if !isValidConversationID(convID) {
		http.Error(w, "Invalid conversation ID format", http.StatusBadRequest)
		return
	}

	// Additional validation for DMs: ensure user is one of the two participants
	if strings.HasPrefix(string(convID), "dm_") {
		if !isDMParticipant(string(convID), userID) {
			log.Printf("[WS] Unauthorized DM access attempt - User: %s, Conv: %s\n", userID, convID)
			http.Error(w, "Forbidden: not a participant in this DM", http.StatusForbidden)
			return
		}
	}

	log.Printf("[WS] New connection attempt - User: %s, Conv: %s\n", userID, convID)

	options := &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	}
	conn, err := websocket.Accept(w, r, options)
	if err != nil {
		log.Printf("[WS] WebSocket upgrade failed: %v\n", err)
		return
	}

	log.Printf("[WS] WebSocket upgraded successfully - User: %s, Conv: %s\n", userID, convID)

	client := &Client{
		ID:   userID,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	// Determine pool type from conversation ID
	poolType := getPoolType(convID)

	// Use the request context only for initial connection setup
	err = h.Manager.HandleConnection(reqCtx, convID, poolType, client)
	if err != nil {
		log.Printf("[WS] Failed to handle connection: %v\n", err)
		conn.Close(websocket.StatusPolicyViolation, err.Error())
		return
	}

	log.Printf("[WS] Client registered with pool - User: %s, Conv: %s, PoolType: %d\n", userID, convID, poolType)

	appCtx := h.Manager.appCtx
	go client.WritePump(appCtx)
	go client.ReadPump(appCtx, h.Manager, convID)

	log.Printf("[WS] Goroutines spawned - WritePump and ReadPump - User: %s, Conv: %s\n", userID, convID)
}

// isDMParticipant extracts user IDs from DM format and checks if userID is one of them.
// Format: "dm_X_Y" where X and Y are user IDs.
func isDMParticipant(convID, userID string) bool {
	parts := strings.Split(convID, "_")
	if len(parts) != 3 || parts[0] != "dm" {
		return false
	}

	id1, err1 := strconv.Atoi(parts[1])
	id2, err2 := strconv.Atoi(parts[2])
	userIDInt, err3 := strconv.Atoi(userID)

	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}

	return userIDInt == id1 || userIDInt == id2
}
