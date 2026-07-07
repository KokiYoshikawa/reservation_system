package repository

import (
	"context"
	"fmt"

	"reservation-system/backend/internal/domain"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) FindReservationSlotByIDForUpdate(ctx context.Context, tx pgx.Tx, slotID int64) (*domain.ReservationSlot, error) {
	const query = `
		SELECT
			id,
			start_time,
			end_time,
			capacity,
			created_at,
			updated_at
		FROM reservation_slots
		WHERE id = $1
		FOR UPDATE
	`

	slot := &domain.ReservationSlot{}
	if err := tx.QueryRow(ctx, query, slotID).Scan(
		&slot.ID,
		&slot.StartTime,
		&slot.EndTime,
		&slot.Capacity,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("find reservation slot by id for update: %w", err)
	}

	return slot, nil
}
