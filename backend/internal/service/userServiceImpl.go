package service

import (
	"context"
	"fmt"

	"reservation-system/backend/internal/domain"
	"reservation-system/backend/internal/dto"
	"reservation-system/backend/internal/repository"
)

type userService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*domain.User, error) {
	user := &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.PasswordHash,
		Role:         req.Role,
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service create user: %w", err)
	}

	return createdUser, nil
}
