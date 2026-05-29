package conversations

import "time"

type ConversationType string

const (
	TypeDM    ConversationType = "dm"
	TypeGroup ConversationType = "group"
)

// Conversation represents either a 1-to-1 DM or a group chat room.
// It acts as the bridge between persistent room data and ephemeral pool lifecycle.
type Conversation struct {
	ID               string            `json:"id" db:"id"`                      // Unique ID: "dm_5_8" or "room_abc123"
	Type             ConversationType  `json:"type" db:"type"`                  // "dm" or "group"
	GroupID          *string           `json:"groupId,omitempty" db:"group_id"` // FK to groups.id (only if type=group)
	UserID1          *int              `json:"userId1,omitempty" db:"user_id_1"`  // FK to users.id (only if type=dm, smaller ID)
	UserID2          *int              `json:"userId2,omitempty" db:"user_id_2"`  // FK to users.id (only if type=dm, larger ID)
	CreatedAt        time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time         `json:"updatedAt" db:"updated_at"`
	LastMessageAt    *time.Time        `json:"lastMessageAt,omitempty" db:"last_message_at"`
}

// CreateDMConversationDTO is the request body for initiating a DM.
type CreateDMConversationDTO struct {
	TargetUserID int `json:"targetUserId" binding:"required"`
}

// CreateGroupConversationDTO is the request body for creating a group.
type CreateGroupConversationDTO struct {
	Name      string `json:"name" binding:"required,min=3,max=100"`
	IsPrivate bool   `json:"isPrivate"`
}
