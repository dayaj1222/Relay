package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUsernameExists = errors.New("username already exists")
	ErrEmailExists    = errors.New("email already exists")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, display_name, email, password_hash, api_key_hash, status)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx,
		query,
		user.Username,
		user.DisplayName,
		user.Email,
		user.PasswordHash,
		user.APIKeyHash,
		user.Status).Scan(&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	query := `SELECT id, username, display_name, email, password_hash, api_key_hash, status, created_at, updated_at FROM users WHERE username = $1 LIMIT 1`

	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Update(ctx context.Context, user *User) error {
	query := `UPDATE users
			 SET display_name = $1,
				password_hash= $2,
				api_key_hash = $3,
				status = $4,
				updated_at = CURRENT_TIMESTAMP
	         WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, user.DisplayName, user.PasswordHash, user.APIKeyHash, user.Status, user.ID)

	return err
}

func (r *Repository) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	query := `
		SELECT CAST(id AS TEXT)
		FROM users
		WHERE api_key_hash = $1
		LIMIT 1`

	var userID string
	err := r.db.QueryRowContext(ctx, query, token).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
