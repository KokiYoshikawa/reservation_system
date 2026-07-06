package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"reservation-system/backend/internal/repository"

	"github.com/jackc/pgx/v5"
)

type reservationSlotService struct {
	repo *repository.Repository
}

func NewReservationSlotService(repo *repository.Repository) ReservationSlotService {
	return &reservationSlotService{repo: repo}
}

func (s *reservationSlotService) SearchAvailableSlots(ctx context.Context, targetDate time.Time, serviceID int64) (*ReservationSlotSearchResult, error) {
	serviceItem, err := s.repo.FindServiceByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrServiceNotFound
		}

		return nil, fmt.Errorf("reservation slot service find service by id: %w", err)
	}

	if !serviceItem.IsActive {
		return &ReservationSlotSearchResult{
			Date:  targetDate.Format("2006-01-02"),
			Slots: []ReservationSlotAvailability{},
		}, nil
	}

	slots, err := s.repo.FindReservationSlotsByDate(ctx, targetDate)
	if err != nil {
		return nil, fmt.Errorf("reservation slot service find reservation slots by date: %w", err)
	}

	slotIDs := make([]int64, 0, len(slots))
	for _, slot := range slots {
		slotIDs = append(slotIDs, slot.ID)
	}

	reservedCounts, err := s.repo.CountReservedBySlotIDs(ctx, slotIDs)
	if err != nil {
		return nil, fmt.Errorf("reservation slot service count reserved by slot ids: %w", err)
	}

	result := &ReservationSlotSearchResult{
		Date:  targetDate.Format("2006-01-02"),
		Slots: make([]ReservationSlotAvailability, 0, len(slots)),
	}

	for _, slot := range slots {
		reservedCount := reservedCounts[slot.ID]
		result.Slots = append(result.Slots, ReservationSlotAvailability{
			SlotID:        slot.ID,
			StartTime:     slot.StartTime,
			EndTime:       slot.EndTime,
			Capacity:      slot.Capacity,
			ReservedCount: reservedCount,
			Available:     reservedCount < slot.Capacity,
		})
	}

	return result, nil
}
