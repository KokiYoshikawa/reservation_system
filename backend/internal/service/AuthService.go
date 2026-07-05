package service

import (
	"context"
	"errors"

	"reservation-system/backend/internal/domain"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrUserNotFound    = errors.New("user not found")
)

type AuthService interface {
	Login(ctx context.Context, email string) (string, *domain.User, error)
	Me(ctx context.Context, token string) (*domain.User, error)
	Logout(ctx context.Context, token string) error
}
