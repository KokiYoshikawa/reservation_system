package service

import (
	"context"

	"reservation-system/backend/internal/domain"
	"reservation-system/backend/internal/dto"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*domain.User, error)
}
