package service

import (
	"context"
	"errors"

	"reservation-system/backend/internal/domain"
)

var ErrServiceNotFound = errors.New("service not found")

type ServiceService interface {
	FindServiceByID(ctx context.Context, serviceID int64) (*domain.Service, error)
	ListActiveServices(ctx context.Context) ([]domain.Service, error)
}
