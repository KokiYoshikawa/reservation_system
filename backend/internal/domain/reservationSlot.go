package domain

import "time"

type ReservationSlot struct {
	ID        int64
	StartTime time.Time
	EndTime   time.Time
	Capacity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
