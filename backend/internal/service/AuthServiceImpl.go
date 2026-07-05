package service

import (
	"context"
	"errors"
	"fmt"

	"reservation-system/backend/internal/domain"
	"reservation-system/backend/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type authService struct {
	repo *repository.Repository
}

func NewAuthService(repo *repository.Repository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Login(ctx context.Context, email string) (string, *domain.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil, ErrUserNotFound
		}

		return "", nil, fmt.Errorf("service login get user by email: %w", err)
	}

	token, err := s.repo.CreateSession(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("service login create session: %w", err)
	}

	return token, user, nil
}

func (s *authService) Me(ctx context.Context, token string) (*domain.User, error) {
	email, err := s.repo.GetSessionEmail(ctx, token)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrSessionNotFound
		}

		return nil, fmt.Errorf("service me get session email: %w", err)
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("service me get user by email: %w", err)
	}

	return user, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	if err := s.repo.DeleteSession(ctx, token); err != nil {
		return fmt.Errorf("service logout delete session: %w", err)
	}

	return nil
}
