package conversations

import (
	"context"
	"fmt"
)

// Service provides high-level conversation operations.
type Service struct {
	store *Store
}

// NewService creates a new conversation service.
func NewService(store *Store) *Service {
	return &Service{store: store}
}

// GetOrCreateDM gets an existing DM or creates a new one.
// Always returns a valid conversation ID.
func (s *Service) GetOrCreateDM(ctx context.Context, userID1, userID2 int) (*Conversation, error) {
	if userID1 == userID2 {
		return nil, fmt.Errorf("cannot create DM with yourself")
	}

	return s.store.CreateDM(ctx, userID1, userID2)
}

// GetOrCreateGroup gets an existing group conversation or creates a new one.
// Assumes the group already exists in the groups table.
func (s *Service) GetOrCreateGroup(ctx context.Context, groupID string) (*Conversation, error) {
	return s.store.CreateGroup(ctx, groupID)
}

// GetConversation fetches a conversation by ID.
func (s *Service) GetConversation(ctx context.Context, convID string) (*Conversation, error) {
	return s.store.GetByID(ctx, convID)
}

// ValidateAccess ensures a user can access a conversation.
// Returns an error if access is denied.
func (s *Service) ValidateAccess(ctx context.Context, userID int, convID string) error {
	hasAccess, err := s.store.CanUserAccess(ctx, userID, convID)
	if err != nil {
		return fmt.Errorf("access check failed: %w", err)
	}

	if !hasAccess {
		return fmt.Errorf("access denied: user %d cannot access conversation %s", userID, convID)
	}

	return nil
}

// GetConversationPoolType returns the pool type for a conversation.
// DMs use PoolTypeIndividual, groups use PoolTypeGroup.
func (s *Service) GetConversationPoolType(ctx context.Context, convID string) (string, error) {
	convType, err := s.store.GetConversationType(ctx, convID)
	if err != nil {
		return "", err
	}

	if convType == TypeDM {
		return "individual", nil
	}
	return "group", nil
}

// ListUserConversations fetches all conversations a user can access.
func (s *Service) ListUserConversations(ctx context.Context, userID int) ([]*Conversation, error) {
	return s.store.ListUserConversations(ctx, userID)
}

// UpdateActivity marks a conversation as recently active.
func (s *Service) UpdateActivity(ctx context.Context, convID string) error {
	return s.store.UpdateLastMessageAt(ctx, convID)
}
