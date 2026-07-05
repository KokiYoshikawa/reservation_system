package repository

import (
	"context"
	"fmt"

	"reservation-system/backend/internal/domain"
)

func (r *Repository) FindServiceByID(ctx context.Context, serviceID int64) (*domain.Service, error) {
	const query = `
		SELECT
			id,
			name,
			duration_minutes,
			price,
			is_active,
			created_at,
			updated_at
		FROM services
		WHERE id = $1
	`

	service := &domain.Service{}
	if err := r.db.QueryRow(ctx, query, serviceID).Scan(
		&service.ID,
		&service.Name,
		&service.DurationMinutes,
		&service.Price,
		&service.IsActive,
		&service.CreatedAt,
		&service.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("find service by id: %w", err)
	}

	return service, nil
}

func (r *Repository) ListActiveServices(ctx context.Context) ([]domain.Service, error) {
	const query = `
		SELECT
			id,
			name,
			duration_minutes,
			price,
			is_active,
			created_at,
			updated_at
		FROM services
		WHERE is_active = TRUE
		ORDER BY id ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list active services: %w", err)
	}
	defer rows.Close()

	services := make([]domain.Service, 0)
	for rows.Next() {
		var service domain.Service
		if err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.DurationMinutes,
			&service.Price,
			&service.IsActive,
			&service.CreatedAt,
			&service.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan active service: %w", err)
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate active services: %w", err)
	}

	return services, nil
}
