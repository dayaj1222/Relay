package messages

import "encoding/json"

type MessageType int

const (
	TypeText MessageType = iota
	TypeVideo
	TypeAudio
	TypeDocument
)

// Message is the wire format used over WebSocket.
// Content is stored as PostgreSQL jsonb and serialized as a JSON object on the wire.
type Message struct {
	ID         int             `json:"id,omitempty" db:"id"`
	SenderID   string          `json:"senderId" db:"sender_id"`
	ReceiverID string          `json:"receiverId,omitempty" db:"-"`
	PoolID     string          `json:"poolId" db:"pool_id"`
	Type       MessageType     `json:"type" db:"type"`
	Content    json.RawMessage `json:"content" db:"content"`
}