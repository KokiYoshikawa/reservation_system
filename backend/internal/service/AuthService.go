package service

import (
	"context"
	"errors"

	"reservation-system/backend/internal/domain"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenRevoked    = errors.New("token revoked")
)

type AuthClaims struct {
	Subject   string
	Email     string
	IssuedAt  int64
	ExpiresAt int64
}

type AuthService interface {
	Login(ctx context.Context, email string) (string, *domain.User, error)
	Me(ctx context.Context, email string) (*domain.User, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (*AuthClaims, error)
}
