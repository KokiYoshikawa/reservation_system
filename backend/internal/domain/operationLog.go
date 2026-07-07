package domain

import "time"

type OperationLogAction string

const (
	OperationLogActionReservationCreated  OperationLogAction = "reservation_created"
	OperationLogActionReservationUpdated  OperationLogAction = "reservation_updated"
	OperationLogActionReservationCanceled OperationLogAction = "reservation_canceled"
)

type OperationLog struct {
	ID            int64
	UserID        int64
	ReservationID *int64
	Action        OperationLogAction
	Detail        *string
	CreatedAt     time.Time
}
