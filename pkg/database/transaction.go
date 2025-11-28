package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/qhato/ecommerce/pkg/logger"
)

// TxFunc is a function that runs within a transaction
type TxFunc func(ctx context.Context, tx pgx.Tx) error

// WithTransaction executes a function within a database transaction
// It automatically commits on success or rolls back on error
func (db *DB) WithTransaction(ctx context.Context, fn TxFunc) error {
	// Begin transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure rollback on panic
	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				logger.WithError(rbErr).Error("Failed to rollback transaction after panic")
			}
			panic(p) // Re-throw panic after rollback
		}
	}()

	// Execute function
	if err := fn(ctx, tx); err != nil {
		// Rollback transaction on error
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			logger.WithError(rbErr).Error("Failed to rollback transaction")
			return fmt.Errorf("transaction error: %w (rollback error: %v)", err, rbErr)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithTransactionOptions executes a function within a database transaction with custom options
func (db *DB) WithTransactionOptions(ctx context.Context, txOptions pgx.TxOptions, fn TxFunc) error {
	// Begin transaction with options
	tx, err := db.BeginTx(ctx, txOptions)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure rollback on panic
	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				logger.WithError(rbErr).Error("Failed to rollback transaction after panic")
			}
			panic(p) // Re-throw panic after rollback
		}
	}()

	// Execute function
	if err := fn(ctx, tx); err != nil {
		// Rollback transaction on error
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			logger.WithError(rbErr).Error("Failed to rollback transaction")
			return fmt.Errorf("transaction error: %w (rollback error: %v)", err, rbErr)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Transactional is a helper to run multiple operations in a transaction
type Transactional struct {
	db *DB
}

// NewTransactional creates a new Transactional helper
func NewTransactional(db *DB) *Transactional {
	return &Transactional{db: db}
}

// Execute runs a transaction with the given function
func (t *Transactional) Execute(ctx context.Context, fn TxFunc) error {
	return t.db.WithTransaction(ctx, fn)
}

// ExecuteWithOptions runs a transaction with options
func (t *Transactional) ExecuteWithOptions(ctx context.Context, txOptions pgx.TxOptions, fn TxFunc) error {
	return t.db.WithTransactionOptions(ctx, txOptions, fn)
}
