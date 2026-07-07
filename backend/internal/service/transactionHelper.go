package service

import (
	"context"
	"errors"
	"fmt"

	"reservation-system/backend/internal/repository"

	"github.com/jackc/pgx/v5"
)

func executeInTx(
	ctx context.Context,
	repo *repository.Repository,
	fn func(tx pgx.Tx) error,
) (err error) {
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err == nil {
			return
		}

		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			err = fmt.Errorf("%w; rollback tx: %v", err, rollbackErr)
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
