package dto

import "reservation-system/backend/internal/domain"

type CreateUserRequest struct {
	Name         string
	Email        string
	PasswordHash string
	Role         domain.UserRole
}
