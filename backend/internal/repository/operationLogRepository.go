package repository

import (
	"context"
	"fmt"

	"reservation-system/backend/internal/domain"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateOperationLog(ctx context.Context, tx pgx.Tx, log *domain.OperationLog) (*domain.OperationLog, error) {
	const query = `
		INSERT INTO operation_logs (
			user_id,
			reservation_id,
			action,
			detail,
			created_at
		)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING
			id,
			user_id,
			reservation_id,
			action,
			detail,
			created_at
	`

	createdLog := &domain.OperationLog{}
	if err := tx.QueryRow(
		ctx,
		query,
		log.UserID,
		log.ReservationID,
		log.Action,
		log.Detail,
	).Scan(
		&createdLog.ID,
		&createdLog.UserID,
		&createdLog.ReservationID,
		&createdLog.Action,
		&createdLog.Detail,
		&createdLog.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("create operation log: %w", err)
	}

	return createdLog, nil
}
