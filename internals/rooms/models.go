package rooms

import "time"

type Room struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatorID int       `json:"creatorId" db:"creator_id"`
	IsPrivate bool      `json:"isPrivate" db:"is_private"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateRoomDTO struct {
	Name      string `json:"name" binding:"required,min=3,max=100"`
	IsPrivate bool   `json:"isPrivate"`
}
