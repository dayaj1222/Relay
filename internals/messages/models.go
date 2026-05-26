package messages

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type MessageType int

const (
	TypeText MessageType = iota
	TypeVideo
	TypeAudio
	TypeDocument
)

// JSONB is a custom type to handle PostgreSQL jsonb compatibility seamlessly
type JSONB json.RawMessage

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan source for JSONB")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

type Message struct {
	ID         int         `json:"id,omitempty" db:"id"`
	SenderID   string      `json:"senderId" db:"sender_id"`
	ReceiverID string      `json:"receiverId,omitempty" db:"-"`
	PoolID     string      `json:"poolId" db:"pool_id"`
	Type       MessageType `json:"type" db:"type"`
	Content    JSONB       `json:"content" db:"content"`
}
