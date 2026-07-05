package service

import (
	"context"

	"reservation-system/backend/internal/domain"
)

type ServiceService interface {
	FindServiceByID(ctx context.Context, serviceID int64) (*domain.Service, error)
	ListActiveServices(ctx context.Context) ([]domain.Service, error)
}
