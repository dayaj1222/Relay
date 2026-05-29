package messages

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

type MessageWithTimestamp struct {
	ID        int       `db:"id" json:"id"`
	SenderID  int       `db:"sender_id" json:"senderId"`
	ConvID    string    `db:"conversation_id" json:"conversationId"`
	Type      int       `db:"type" json:"type"`
	Content   JSONB     `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// CreateMessage inserts a new message
func (s *Store) CreateMessage(ctx context.Context, senderID int, convID string, msgType MessageType, content json.RawMessage) (*MessageWithTimestamp, error) {
	msg := &MessageWithTimestamp{
		SenderID:  senderID,
		ConvID:    convID,
		Type:      int(msgType),
		Content:   JSONB(content),
		CreatedAt: time.Now(),
	}

	err := s.db.QueryRowxContext(ctx,
		`INSERT INTO messages (sender_id, conversation_id, type, content, created_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, sender_id, conversation_id, type, content, created_at`,
		senderID, convID, msgType, JSONB(content), time.Now()).StructScan(msg)

	return msg, err
}

// GetMessagesByConversation fetches messages for a conversation, paginated
func (s *Store) GetMessagesByConversation(ctx context.Context, convID string, limit int, offset int) ([]MessageWithTimestamp, error) {
	var msgs []MessageWithTimestamp
	err := s.db.SelectContext(ctx, &msgs,
		`SELECT id, sender_id, conversation_id, type, content, created_at
		 FROM messages
		 WHERE conversation_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		convID, limit, offset)
	return msgs, err
}

// GetRecentMessages fetches last N messages
func (s *Store) GetRecentMessages(ctx context.Context, convID string, limit int) ([]MessageWithTimestamp, error) {
	var msgs []MessageWithTimestamp
	err := s.db.SelectContext(ctx, &msgs,
		`SELECT id, sender_id, conversation_id, type, content, created_at
		 FROM messages
		 WHERE conversation_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2`,
		convID, limit)

	// Reverse to show oldest first
	if len(msgs) > 1 {
		for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
			msgs[i], msgs[j] = msgs[j], msgs[i]
		}
	}

	return msgs, err
}
