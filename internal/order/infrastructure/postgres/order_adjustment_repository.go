package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/order/domain"
)

// OrderAdjustmentRepository implements domain.OrderAdjustmentRepository for PostgreSQL persistence.
type OrderAdjustmentRepository struct {
	db *sql.DB
}

// NewOrderAdjustmentRepository creates a new PostgreSQL order adjustment repository.
func NewOrderAdjustmentRepository(db *sql.DB) *OrderAdjustmentRepository {
	return &OrderAdjustmentRepository{db: db}
}

// Save stores a new order adjustment or updates an existing one.
func (r *OrderAdjustmentRepository) Save(ctx context.Context, adjustment *domain.OrderAdjustment) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle nullable fields
	isFutureCredit := sql.NullBool{Bool: adjustment.IsFutureCredit, Valid: true}
	offerID := sql.NullInt64{Int64: adjustment.OfferID, Valid: adjustment.OfferID != 0}
	orderID := sql.NullInt64{Int64: adjustment.OrderID, Valid: adjustment.OrderID != 0}

	if adjustment.ID == 0 {
		// Insert new order adjustment
		query := `
			INSERT INTO blc_order_adjustment (
				order_id, offer_id, adjustment_reason, adjustment_value, is_future_credit, created_at
			) VALUES (
				$1, $2, $3, $4, $5, $6
			) RETURNING order_adjustment_id`
		err = tx.QueryRowContext(ctx, query,
			orderID, offerID, adjustment.AdjustmentReason, adjustment.AdjustmentValue, isFutureCredit, adjustment.CreatedAt,
		).Scan(&adjustment.ID)
		if err != nil {
			return fmt.Errorf("failed to insert order adjustment: %w", err)
		}
	} else {
		// Update existing order adjustment
		query := `
			UPDATE blc_order_adjustment SET
				order_id = $1, offer_id = $2, adjustment_reason = $3, adjustment_value = $4, 
				is_future_credit = $5
			WHERE order_adjustment_id = $6`
		_, err = tx.ExecContext(ctx, query,
			orderID, offerID, adjustment.AdjustmentReason, adjustment.AdjustmentValue,
			isFutureCredit, adjustment.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update order adjustment: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves an order adjustment by its unique identifier.
func (r *OrderAdjustmentRepository) FindByID(ctx context.Context, id int64) (*domain.OrderAdjustment, error) {
	query := `
		SELECT
			order_adjustment_id, order_id, offer_id, adjustment_reason, adjustment_value, 
			is_future_credit, created_at
		FROM blc_order_adjustment WHERE order_adjustment_id = $1`

	var adjustment domain.OrderAdjustment
	var offerID sql.NullInt64
	var orderID sql.NullInt64
	var isFutureCredit sql.NullBool

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&adjustment.ID, &orderID, &offerID, &adjustment.AdjustmentReason, &adjustment.AdjustmentValue,
		&isFutureCredit, &adjustment.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query order adjustment by ID: %w", err)
	}

	if offerID.Valid {
		adjustment.OfferID = offerID.Int64
	}
	if orderID.Valid {
		adjustment.OrderID = orderID.Int64
	}
	if isFutureCredit.Valid {
		adjustment.IsFutureCredit = isFutureCredit.Bool
	}

	return &adjustment, nil
}

// FindByOrderID retrieves all order adjustments for a given order ID.
func (r *OrderAdjustmentRepository) FindByOrderID(ctx context.Context, orderID int64) ([]*domain.OrderAdjustment, error) {
	query := `
		SELECT
			order_adjustment_id, order_id, offer_id, adjustment_reason, adjustment_value, 
			is_future_credit, created_at
		FROM blc_order_adjustment WHERE order_id = $1`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order adjustments by order ID: %w", err)
	}
	defer rows.Close()

	var adjustments []*domain.OrderAdjustment
	for rows.Next() {
		var adjustment domain.OrderAdjustment
		var offerID sql.NullInt64
		var ordID sql.NullInt64 // Use different name to avoid conflict
		var isFutureCredit sql.NullBool

		err := rows.Scan(
			&adjustment.ID, &ordID, &offerID, &adjustment.AdjustmentReason, &adjustment.AdjustmentValue,
			&isFutureCredit, &adjustment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order adjustment row: %w", err)
		}

		if offerID.Valid {
			adjustment.OfferID = offerID.Int64
		}
		if ordID.Valid {
			adjustment.OrderID = ordID.Int64
		}
		if isFutureCredit.Valid {
			adjustment.IsFutureCredit = isFutureCredit.Bool
		}
		adjustments = append(adjustments, &adjustment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for order adjustments: %w", err)
	}

	return adjustments, nil
}

// Delete removes an order adjustment by its unique identifier.
func (r *OrderAdjustmentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_order_adjustment WHERE order_adjustment_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order adjustment: %w", err)
	}
	return nil
}

// DeleteByOrderID removes all order adjustments for a given order ID.
func (r *OrderAdjustmentRepository) DeleteByOrderID(ctx context.Context, orderID int64) error {
	query := `DELETE FROM blc_order_adjustment WHERE order_id = $1`
	_, err := r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order adjustments by order ID: %w", err)
	}
	return nil
}
