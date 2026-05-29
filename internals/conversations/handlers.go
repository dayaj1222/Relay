package conversations

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HTTPHandler provides HTTP endpoints for conversation management.
type HTTPHandler struct {
	service *Service
}

// NewHTTPHandler creates a new conversation HTTP handler.
func NewHTTPHandler(service *Service) *HTTPHandler {
	return &HTTPHandler{service: service}
}

// CreateDM initiates a DM conversation with another user.
// POST /conversations/dm
// Body: { "targetUserId": 5 }
func (h *HTTPHandler) CreateDM(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currentUserID, err := strconv.Atoi(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req CreateDMConversationDTO
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	conv, err := h.service.GetOrCreateDM(ctx.Request.Context(), currentUserID, req.TargetUserID)
	if err != nil {
		log.Printf("[Conversations] CreateDM error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, conv)
}

// CreateGroup creates a new group conversation.
// POST /conversations/group
// Body: { "name": "Project Alpha", "isPrivate": false }
func (h *HTTPHandler) CreateGroup(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateGroupConversationDTO
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// TODO: This will be called after group is created in groups table.
	// For now, it just creates the conversation entry.
	// The actual group creation should be in a separate groups service.

	log.Printf("[Conversations] CreateGroup by user %s: %s\n", userID, req.Name)
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "use groups endpoint to create groups"})
}

// GetConversation retrieves a single conversation.
// GET /conversations/:id
func (h *HTTPHandler) GetConversation(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currentUserID, err := strconv.Atoi(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	convID := ctx.Param("id")

	// Validate access
	if err := h.service.ValidateAccess(ctx.Request.Context(), currentUserID, convID); err != nil {
		log.Printf("[Conversations] Access denied for user %d: %v\n", currentUserID, err)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	conv, err := h.service.GetConversation(ctx.Request.Context(), convID)
	if err != nil {
		log.Printf("[Conversations] GetConversation error: %v\n", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	ctx.JSON(http.StatusOK, conv)
}

// ListConversations lists all conversations for the current user.
// GET /conversations
func (h *HTTPHandler) ListConversations(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	log.Printf("[DEBUG] user_id from context: %v (type: %T)", userID, userID)

	currentUserID, err := strconv.Atoi(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	convs, err := h.service.ListUserConversations(ctx.Request.Context(), currentUserID)
	if err != nil {
		log.Printf("[Conversations] ListConversations error: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch conversations"})
		return
	}

	if convs == nil {
		convs = []*Conversation{} // Return empty array, not null
	}

	ctx.JSON(http.StatusOK, convs)
}
