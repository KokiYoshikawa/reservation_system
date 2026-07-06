package service

import (
	"context"
	"errors"
	"fmt"

	"reservation-system/backend/internal/domain"
	"reservation-system/backend/internal/repository"

	"github.com/jackc/pgx/v5"
)

type serviceService struct {
	repo *repository.Repository
}

func NewServiceService(repo *repository.Repository) ServiceService {
	return &serviceService{repo: repo}
}

func (s *serviceService) FindServiceByID(ctx context.Context, serviceID int64) (*domain.Service, error) {
	service, err := s.repo.FindServiceByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrServiceNotFound
		}

		return nil, fmt.Errorf("service find service by id: %w", err)
	}

	return service, nil
}

func (s *serviceService) ListActiveServices(ctx context.Context) ([]domain.Service, error) {
	services, err := s.repo.ListActiveServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("service list active services: %w", err)
	}

	return services, nil
}
