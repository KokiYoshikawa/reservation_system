package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"reservation-system/backend/internal/domain"
	"reservation-system/backend/internal/dto"
	"reservation-system/backend/internal/repository"

	"github.com/jackc/pgx/v5"
)

type reservationService struct {
	repo *repository.Repository
}

func NewReservationService(repo *repository.Repository) ReservationService {
	return &reservationService{repo: repo}
}

func (s *reservationService) CreateReservation(ctx context.Context, userID int64, req dto.CreateReservationRequest) (*ReservationResult, error) {
	serviceItem, err := s.repo.FindServiceByID(ctx, req.ServiceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrServiceNotFound
		}

		return nil, fmt.Errorf("reservation service find service by id: %w", err)
	}

	if !serviceItem.IsActive {
		return nil, ErrServiceNotFound
	}

	var note *string
	if req.Note != "" {
		note = &req.Note
	}

	var (
		slot               *domain.ReservationSlot
		createdReservation *domain.Reservation
	)

	if err := executeInTx(ctx, s.repo, func(tx pgx.Tx) error {
		slot, err = s.repo.FindReservationSlotByIDForUpdate(ctx, tx, req.SlotID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrReservationSlotNotFound
			}

			return fmt.Errorf("reservation service find slot for update: %w", err)
		}

		if !slot.StartTime.After(time.Now()) {
			return ErrReservationSlotPast
		}

		reserved, err := s.repo.ExistsActiveReservationBySlotID(ctx, tx, req.SlotID)
		if err != nil {
			return fmt.Errorf("reservation service check slot conflict: %w", err)
		}

		if reserved {
			return ErrReservationConflict
		}

		reservation := &domain.Reservation{
			UserID:     userID,
			ServiceID:  req.ServiceID,
			SlotID:     req.SlotID,
			Status:     domain.ReservationStatusReserved,
			Note:       note,
			ReservedAt: time.Now(),
		}

		createdReservation, err = s.repo.CreateReservation(ctx, tx, reservation)
		if err != nil {
			return fmt.Errorf("reservation service create reservation: %w", err)
		}

		if err := s.createReservationOperationLog(ctx, tx, userID, createdReservation); err != nil {
			return fmt.Errorf("reservation service create operation log: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &ReservationResult{
		ReservationID: createdReservation.ID,
		Status:        createdReservation.Status,
		ServiceID:     serviceItem.ID,
		ServiceName:   serviceItem.Name,
		StartTime:     slot.StartTime,
		EndTime:       slot.EndTime,
		ReservedAt:    createdReservation.ReservedAt,
	}, nil
}
