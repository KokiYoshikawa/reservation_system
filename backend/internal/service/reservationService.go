package service

import (
	"context"
	"errors"
	"time"

	"reservation-system/backend/internal/domain"
	"reservation-system/backend/internal/dto"
)

var (
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
	ReservedAt    time.Time
}

type ReservationService interface {
	CreateReservation(ctx context.Context, userID int64, req dto.CreateReservationRequest) (*ReservationResult, error)
}
