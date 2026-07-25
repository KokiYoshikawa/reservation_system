package repository

import (
	"context"
	"fmt"
	"time"

	"reservation-system/backend/internal/domain"
)

func (r *Repository) FindReservationSlotByID(ctx context.Context, slotID int64) (*domain.ReservationSlot, error) {
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
	`

	slot := &domain.ReservationSlot{}
	if err := r.db.QueryRow(ctx, query, slotID).Scan(
		&slot.ID,
		&slot.StartTime,
		&slot.EndTime,
		&slot.Capacity,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("find reservation slot by id: %w", err)
	}

	return slot, nil
}

func (r *Repository) FindReservationSlotsByDate(ctx context.Context, targetDate time.Time) ([]domain.ReservationSlot, error) {
	startOfDay := time.Date(
		targetDate.Year(),
		targetDate.Month(),
		targetDate.Day(),
		0,
		0,
		0,
		0,
		targetDate.Location(),
	)
	endOfDay := startOfDay.Add(24 * time.Hour)

	const query = `
		SELECT
			id,
			start_time,
			end_time,
			capacity,
			created_at,
			updated_at
		FROM reservation_slots
		WHERE start_time >= $1
		  AND start_time < $2
		ORDER BY start_time ASC, id ASC
	`

	rows, err := r.db.Query(ctx, query, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("find reservation slots by date: %w", err)
	}
	defer rows.Close()

	slots := make([]domain.ReservationSlot, 0)
	for rows.Next() {
		var slot domain.ReservationSlot
		if err := rows.Scan(
			&slot.ID,
			&slot.StartTime,
			&slot.EndTime,
			&slot.Capacity,
			&slot.CreatedAt,
			&slot.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan reservation slot by date: %w", err)
		}

		slots = append(slots, slot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reservation slots by date: %w", err)
	}

	return slots, nil
}
