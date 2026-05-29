package messages

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	store *Store
}

func NewHTTPHandler(store *Store) *HTTPHandler {
	return &HTTPHandler{store: store}
}

type SendMessageDTO struct {
	Type    int             `json:"type"`
	Content json.RawMessage `json:"content" binding:"required"`
}

// SendMessage handles POST /api/conversations/:id/messages
func (h *HTTPHandler) SendMessage(c *gin.Context) {
	convID := c.Param("id")
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDInt, err := strconv.Atoi(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var dto SendMessageDTO
	if err := c.BindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if dto.Type < 0 || dto.Type > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message type"})
		return
	}

	msg, err := h.store.CreateMessage(c.Request.Context(), userIDInt, convID, MessageType(dto.Type), dto.Content)
	if err != nil {
		log.Printf("failed to create message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create message"})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

// GetMessages handles GET /api/conversations/:id/messages
func (h *HTTPHandler) GetMessages(c *gin.Context) {
	convID := c.Param("id")

	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 200 {
			limit = parsed
		}
	}

	offset := 0
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	msgs, err := h.store.GetMessagesByConversation(c.Request.Context(), convID, limit, offset)
	if err != nil {
		log.Printf("failed to get messages: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
		return
	}

	if msgs == nil {
		msgs = []MessageWithTimestamp{}
	}

	c.JSON(http.StatusOK, msgs)
}

// GetRecentMessages handles GET /api/conversations/:id/messages/recent
func (h *HTTPHandler) GetRecentMessages(c *gin.Context) {
	convID := c.Param("id")

	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 200 {
			limit = parsed
		}
	}

	msgs, err := h.store.GetRecentMessages(c.Request.Context(), convID, limit)
	if err != nil {
		log.Printf("failed to get messages: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
		return
	}

	if msgs == nil {
		msgs = []MessageWithTimestamp{}
	}

	c.JSON(http.StatusOK, msgs)
}
