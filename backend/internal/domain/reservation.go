package domain

import "time"

type ReservationStatus string

const (
	ReservationStatusReserved  ReservationStatus = "reserved"
	ReservationStatusCancelled ReservationStatus = "cancelled"
)

type Reservation struct {
	ID          int64
	UserID      int64
	ServiceID   int64
	SlotID      int64
	Status      ReservationStatus
	Note        *string
	ReservedAt  time.Time
	CancelledAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
