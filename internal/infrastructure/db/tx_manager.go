package db

import (
	"context"
	"database/sql"
	"fmt"

	"pr-service/internal/app"
)

type txManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) app.TxManager {
	return &txManager{db: db}
}

func (tm *txManager) Do(ctx context.Context, operation func(ctx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	operationErr := operation(ctx)
	if operationErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("operation failed: %w, rollback error: %v", operationErr, rollbackErr)
		}
		return fmt.Errorf("operation failed: %w", operationErr)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
