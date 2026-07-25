package service

import (
	"context"
	"errors"
	"time"

	"reservation-system/backend/internal/domain"
	"reservation-system/backend/internal/dto"
)

var (
	ErrReservationNotFound     = errors.New("reservation not found")
	ErrReservationSlotNotFound = errors.New("reservation slot not found")
	ErrReservationConflict     = errors.New("reservation conflict")
	ErrReservationSlotPast     = errors.New("reservation slot is in the past")
)

type ReservationResult struct {
	ReservationID int64
	Status        domain.ReservationStatus
	ServiceID     int64
	ServiceName   string
	StartTime     time.Time
	EndTime       time.Time
	SlotID        int64
	Note          *string
	ReservedAt    time.Time
	CancelledAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ReservationService interface {
	CreateReservation(ctx context.Context, userID int64, req dto.CreateReservationRequest) (*ReservationResult, error)
	GetReservationsByUser(ctx context.Context, userID int64) ([]ReservationResult, error)
	GetReservationDetail(ctx context.Context, userID, reservationID int64) (*ReservationResult, error)
}
