package users

import (
	"context"
	"errors"
)

var (
	ErrUserAlreadyExists  = errors.New("username is already taken")
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
		return nil, "", err
	}

	return newUser, rawAPIKey, nil
}

func (s *Service) LoginBusiness(ctx context.Context, dto *LoginUserDTO) (*User, error) {
	user, err := s.repo.GetByUsername(ctx, dto.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !CheckPasswordHash(dto.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
