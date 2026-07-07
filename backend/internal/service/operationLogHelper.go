package service

import (
	"context"
	"fmt"

	"reservation-system/backend/internal/domain"

	"github.com/jackc/pgx/v5"
)

func (s *reservationService) createReservationOperationLog(
	ctx context.Context,
	tx pgx.Tx,
	userID int64,
	reservation *domain.Reservation,
) error {
	logDetail := fmt.Sprintf(
		"reservation created: reservation_id=%d service_id=%d slot_id=%d status=%s",
		reservation.ID,
		reservation.ServiceID,
		reservation.SlotID,
		reservation.Status,
	)

	operationLog := &domain.OperationLog{
		UserID:        userID,
		ReservationID: &reservation.ID,
		Action:        domain.OperationLogActionReservationCreated,
		Detail:        &logDetail,
	}

	if _, err := s.repo.CreateOperationLog(ctx, tx, operationLog); err != nil {
		return fmt.Errorf("create reservation operation log: %w", err)
	}

	return nil
}
