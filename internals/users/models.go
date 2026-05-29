// Package users
package users

import "time"

type User struct {
	ID           int     `json:"id" db:"id"`
	Username     string  `json:"username" db:"username"`
	DisplayName  *string `json:"displayName" db:"display_name"`
	Email        string  `json:"email" db:"email"`
	PasswordHash string  `json:"-" db:"password_hash"`
	APIKeyHash   string  `json:"-" db:"api_key_hash"`
	Status       string  `json:"status" db:"status"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterUserDTO struct {
	Username    string  `json:"username" binding:"required,min=3,max=30"`
	Email       string  `json:"email"`
	Password    string  `json:"password" binding:"required"`
	DisplayName *string `json:"displayName"`
}

type LoginUserDTO struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required"`
}
