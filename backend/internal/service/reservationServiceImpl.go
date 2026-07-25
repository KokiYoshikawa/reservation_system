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

func (s *reservationService) GetReservationsByUser(ctx context.Context, userID int64) ([]ReservationResult, error) {
	reservations, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("reservation service find by user id: %w", err)
	}

	results := make([]ReservationResult, 0, len(reservations))
	for i := range reservations {
		result, err := s.toReservationResult(ctx, &reservations[i])
		if err != nil {
			return nil, err
		}

		results = append(results, *result)
	}

	return results, nil
}

func (s *reservationService) GetReservationDetail(ctx context.Context, userID, reservationID int64) (*ReservationResult, error) {
	reservation, err := s.repo.FindByID(ctx, reservationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrReservationNotFound
		}

		return nil, fmt.Errorf("reservation service find by id: %w", err)
	}

	// Return the same result for a missing reservation and another user's
	// reservation so callers cannot discover reservation IDs they do not own.
	if reservation.UserID != userID {
		return nil, ErrReservationNotFound
	}

	return s.toReservationResult(ctx, reservation)
}

func (s *reservationService) toReservationResult(ctx context.Context, reservation *domain.Reservation) (*ReservationResult, error) {
	serviceItem, err := s.repo.FindServiceByID(ctx, reservation.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("reservation service find related service: %w", err)
	}

	slot, err := s.repo.FindReservationSlotByID(ctx, reservation.SlotID)
	if err != nil {
		return nil, fmt.Errorf("reservation service find related slot: %w", err)
	}

	return &ReservationResult{
		ReservationID: reservation.ID,
		Status:        reservation.Status,
		ServiceID:     serviceItem.ID,
		ServiceName:   serviceItem.Name,
		SlotID:        slot.ID,
		StartTime:     slot.StartTime,
		EndTime:       slot.EndTime,
		Note:          reservation.Note,
		ReservedAt:    reservation.ReservedAt,
		CancelledAt:   reservation.CancelledAt,
		CreatedAt:     reservation.CreatedAt,
		UpdatedAt:     reservation.UpdatedAt,
	}, nil
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
