package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserAlreadyExists  = errors.New("username or email is already taken")
	ErrWeakPassword       = errors.New("password must be at least 8 characters long and contain uppercase, lowercase, and numbers")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterBusiness(ctx context.Context, dto *RegisterUserDTO) (*User, string, error) {
	if !(IsStrongPassword(dto.Password)) {
		return nil, "", ErrWeakPassword
	}

	existingUser, err := s.repo.GetByUsername(ctx, dto.Username)
	if err != nil {
		return nil, "", err
	}
	if existingUser != nil {
		return nil, "", ErrUserAlreadyExists
	}

	var displayName string
	if dto.DisplayName != nil && *dto.DisplayName != "" {
		displayName = *dto.DisplayName
	} else {
		displayName = dto.Username
	}

	passwordHash, err := HashPassword(dto.Password)
	if err != nil {
		return nil, "", err
	}

	rawAPIKey, apiKeyHash, err := generateAPIKey()
	if err != nil {
		return nil, "", err
	}

	newUser := &User{
		Username:     dto.Username,
		DisplayName:  &displayName,
		Email:        dto.Email,
		PasswordHash: passwordHash,
		APIKeyHash:   apiKeyHash,
		Status:       "active",
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		if isUniqueConstraintError(err) {
			return nil, "", ErrUserAlreadyExists
		}
		return nil, "", err
	}

	return newUser, rawAPIKey, nil
}

func (s *Service) LoginBusiness(ctx context.Context, dto *LoginUserDTO) (*User, string, error) {
	user, err := s.repo.GetByUsername(ctx, dto.Username)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}
	if !CheckPasswordHash(dto.Password, user.PasswordHash) {
		return nil, "", ErrInvalidCredentials
	}

	// 1. Generate a fresh API key pair
	rawAPIKey, apiKeyHash, err := generateAPIKey()
	if err != nil {
		return nil, "", err
	}

	// 2. Update the user's stored hash in the database
	user.APIKeyHash = apiKeyHash
	if err := s.repo.UpdateAPIKey(ctx, user.ID, apiKeyHash); err != nil {
		return nil, "", err
	}

	return user, rawAPIKey, nil
}

func (s *Service) FindUserIDByToken(ctx context.Context, rawToken string) (string, error) {
	tokenHash := HashAPIKey(rawToken)

	id, err := s.repo.GetUserIDByToken(ctx, tokenHash)
	if err != nil {
		return "", fmt.Errorf("token lookup failed: %w", err)
	}
	return id, nil
}

// GetByUsername returns a user by username, or an error if not found.
func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found: %s", username)
	}
	return user, nil
}

func isUniqueConstraintError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}


