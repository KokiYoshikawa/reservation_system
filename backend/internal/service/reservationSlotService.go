package service

import (
	"context"
	"time"
)

type ReservationSlotAvailability struct {
	SlotID        int64
	StartTime     time.Time
	EndTime       time.Time
	Capacity      int
	ReservedCount int
	Available     bool
}

type ReservationSlotSearchResult struct {
	Date  string
	Slots []ReservationSlotAvailability
}

type ReservationSlotService interface {
	SearchAvailableSlots(ctx context.Context, targetDate time.Time, serviceID int64) (*ReservationSlotSearchResult, error)
}
