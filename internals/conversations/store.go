package conversations

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Store handles all conversation database operations.
type Store struct {
	db *sqlx.DB
}

// NewStore creates a new conversation store.
func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

// CreateDM creates a new DM conversation between two users.
// Ensures user IDs are sorted (smaller first) for uniqueness.
func (s *Store) CreateDM(ctx context.Context, userID1, userID2 int) (*Conversation, error) {
	// Sort user IDs for consistency
	if userID1 > userID2 {
		userID1, userID2 = userID2, userID1
	}

	convID := fmt.Sprintf("dm_%d_%d", userID1, userID2)

	// Check if already exists
	existing, err := s.GetByID(ctx, convID)
	if err == nil && existing != nil {
		return existing, nil // Already exists, reuse
	}

	conv := &Conversation{
		ID:        convID,
		Type:      TypeDM,
		UserID1:   &userID1,
		UserID2:   &userID2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO conversations (id, type, user_id_1, user_id_2, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO NOTHING
		RETURNING id, type, user_id_1, user_id_2, group_id, created_at, updated_at, last_message_at
	`

	err = s.db.QueryRowContext(ctx, query, conv.ID, conv.Type, conv.UserID1, conv.UserID2, conv.CreatedAt, conv.UpdatedAt).
		Scan(&conv.ID, &conv.Type, &conv.UserID1, &conv.UserID2, &conv.GroupID, &conv.CreatedAt, &conv.UpdatedAt, &conv.LastMessageAt)

	if err != nil {
		// If conflict (already exists), fetch and return it
		return s.GetByID(ctx, convID)
	}

	return conv, nil
}

// CreateGroup creates a new group conversation.
func (s *Store) CreateGroup(ctx context.Context, groupID string) (*Conversation, error) {
	convID := fmt.Sprintf("room_%s", groupID)

	conv := &Conversation{
		ID:        convID,
		Type:      TypeGroup,
		GroupID:   &groupID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO conversations (id, type, group_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO NOTHING
		RETURNING id, type, user_id_1, user_id_2, group_id, created_at, updated_at, last_message_at
	`

	err := s.db.QueryRowContext(ctx, query, conv.ID, conv.Type, conv.GroupID, conv.CreatedAt, conv.UpdatedAt).
		Scan(&conv.ID, &conv.Type, &conv.UserID1, &conv.UserID2, &conv.GroupID, &conv.CreatedAt, &conv.UpdatedAt, &conv.LastMessageAt)

	if err != nil {
		// If conflict (already exists), fetch and return it
		return s.GetByID(ctx, convID)
	}

	return conv, nil
}

// GetByID fetches a conversation by its ID.
func (s *Store) GetByID(ctx context.Context, convID string) (*Conversation, error) {
	conv := &Conversation{}

	query := `
		SELECT id, type, user_id_1, user_id_2, group_id, created_at, updated_at, last_message_at
		FROM conversations
		WHERE id = $1
	`

	err := s.db.QueryRowContext(ctx, query, convID).
		Scan(&conv.ID, &conv.Type, &conv.UserID1, &conv.UserID2, &conv.GroupID, &conv.CreatedAt, &conv.UpdatedAt, &conv.LastMessageAt)

	if err != nil {
		return nil, err
	}

	return conv, nil
}

// ListUserConversations fetches all conversations (DMs + groups) for a user.
func (s *Store) ListUserConversations(ctx context.Context, userID int) ([]*Conversation, error) {
	query := `
		SELECT id, type, user_id_1, user_id_2, group_id, created_at, updated_at, last_message_at
		FROM conversations
		WHERE
			(type = 'dm' AND (user_id_1 = $1 OR user_id_2 = $1))
			OR
			(type = 'group' AND group_id IN (
				SELECT id FROM groups WHERE creator_id = $1 OR is_private = false
			))
		ORDER BY last_message_at DESC NULLS LAST
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []*Conversation
	for rows.Next() {
		conv := &Conversation{}
		err := rows.Scan(&conv.ID, &conv.Type, &conv.UserID1, &conv.UserID2, &conv.GroupID, &conv.CreatedAt, &conv.UpdatedAt, &conv.LastMessageAt)
		if err != nil {
			return nil, err
		}
		convs = append(convs, conv)
	}

	return convs, rows.Err()
}

// UpdateLastMessageAt updates the last_message_at timestamp for a conversation.
func (s *Store) UpdateLastMessageAt(ctx context.Context, convID string) error {
	query := `
		UPDATE conversations
		SET last_message_at = $1, updated_at = $1
		WHERE id = $2
	`

	_, err := s.db.ExecContext(ctx, query, time.Now(), convID)
	return err
}

// CanUserAccess checks if a user has permission to access a conversation.
func (s *Store) CanUserAccess(ctx context.Context, userID int, convID string) (bool, error) {
	conv, err := s.GetByID(ctx, convID)
	if err != nil {
		return false, err
	}

	if conv.Type == TypeDM {
		// Only the two participants can access a DM
		return (conv.UserID1 != nil && *conv.UserID1 == userID) ||
			(conv.UserID2 != nil && *conv.UserID2 == userID), nil
	}

	// Type is group
	if conv.GroupID == nil {
		return false, nil
	}

	// Check if user is member of group (or group is public)
	query := `
		SELECT EXISTS(
			SELECT 1 FROM groups g
			WHERE g.id = $1 AND (
				g.creator_id = $2 OR
				g.is_private = false OR
				EXISTS (SELECT 1 FROM group_members gm WHERE gm.group_id = g.id AND gm.user_id = $2)
			)
		)
	`

	var hasAccess bool
	err = s.db.QueryRowContext(ctx, query, *conv.GroupID, userID).Scan(&hasAccess)
	return hasAccess, err
}

// GetConversationType returns the type of the conversation (DM or Group).
func (s *Store) GetConversationType(ctx context.Context, convID string) (ConversationType, error) {
	conv, err := s.GetByID(ctx, convID)
	if err != nil {
		return "", err
	}
	return conv.Type, nil
}
