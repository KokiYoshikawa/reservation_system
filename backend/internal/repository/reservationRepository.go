package repository

import (
	"context"
	"fmt"

	"reservation-system/backend/internal/domain"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) CountReservedBySlotIDs(ctx context.Context, slotIDs []int64) (map[int64]int, error) {
	if len(slotIDs) == 0 {
		return map[int64]int{}, nil
	}

	const query = `
		SELECT
			slot_id,
			COUNT(*)
		FROM reservations
		WHERE slot_id = ANY($1)
		  AND status = 'reserved'
		GROUP BY slot_id
	`

	rows, err := r.db.Query(ctx, query, slotIDs)
	if err != nil {
		return nil, fmt.Errorf("count reserved by slot ids: %w", err)
	}
	defer rows.Close()

	counts := make(map[int64]int, len(slotIDs))
	for rows.Next() {
		var slotID int64
		var count int64
		if err := rows.Scan(&slotID, &count); err != nil {
			return nil, fmt.Errorf("scan reserved counts by slot ids: %w", err)
		}

		counts[slotID] = int(count)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reserved counts by slot ids: %w", err)
	}

	return counts, nil
}

func (r *Repository) ExistsActiveReservationBySlotID(ctx context.Context, tx pgx.Tx, slotID int64) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM reservations
			WHERE slot_id = $1
			  AND status = 'reserved'
		)
	`

	var exists bool
	if err := tx.QueryRow(ctx, query, slotID).Scan(&exists); err != nil {
		return false, fmt.Errorf("exists active reservation by slot id: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateReservation(ctx context.Context, tx pgx.Tx, reservation *domain.Reservation) (*domain.Reservation, error) {
	const query = `
		INSERT INTO reservations (
			user_id,
			service_id,
			slot_id,
			status,
			note,
			reserved_at,
			cancelled_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING
			id,
			user_id,
			service_id,
			slot_id,
			status,
			note,
			reserved_at,
			cancelled_at,
			created_at,
			updated_at
	`

	createdReservation := &domain.Reservation{}
	if err := tx.QueryRow(
		ctx,
		query,
		reservation.UserID,
		reservation.ServiceID,
		reservation.SlotID,
		reservation.Status,
		reservation.Note,
		reservation.ReservedAt,
		reservation.CancelledAt,
	).Scan(
		&createdReservation.ID,
		&createdReservation.UserID,
		&createdReservation.ServiceID,
		&createdReservation.SlotID,
		&createdReservation.Status,
		&createdReservation.Note,
		&createdReservation.ReservedAt,
		&createdReservation.CancelledAt,
		&createdReservation.CreatedAt,
		&createdReservation.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("create reservation: %w", err)
	}

	return createdReservation, nil
}
