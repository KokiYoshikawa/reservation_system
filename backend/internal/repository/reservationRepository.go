package repository

import (
	"context"
	"fmt"
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
