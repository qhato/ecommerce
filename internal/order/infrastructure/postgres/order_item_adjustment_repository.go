package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/order/domain"
)

// OrderItemAdjustmentRepository implements domain.OrderItemAdjustmentRepository for PostgreSQL persistence.
type OrderItemAdjustmentRepository struct {
	db *sql.DB
}

// NewOrderItemAdjustmentRepository creates a new PostgreSQL order item adjustment repository.
func NewOrderItemAdjustmentRepository(db *sql.DB) *OrderItemAdjustmentRepository {
	return &OrderItemAdjustmentRepository{db: db}
}

// Save stores a new order item adjustment or updates an existing one.
func (r *OrderItemAdjustmentRepository) Save(ctx context.Context, adjustment *domain.OrderItemAdjustment) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle nullable fields
	appliedToSalePrice := sql.NullBool{Bool: adjustment.AppliedToSalePrice, Valid: true}
	offerID := sql.NullInt64{Int64: adjustment.OfferID, Valid: adjustment.OfferID != 0}
	orderItemID := sql.NullInt64{Int64: adjustment.OrderItemID, Valid: adjustment.OrderItemID != 0}

	if adjustment.ID == 0 {
		// Insert new order item adjustment
		query := `
			INSERT INTO blc_order_item_adjustment (
				order_item_id, offer_id, adjustment_reason, adjustment_value, applied_to_sale_price, created_at
			) VALUES (
				$1, $2, $3, $4, $5, $6
			) RETURNING order_item_adjustment_id`
		err = tx.QueryRowContext(ctx, query,
			orderItemID, offerID, adjustment.AdjustmentReason, adjustment.AdjustmentValue, appliedToSalePrice, adjustment.CreatedAt,
		).Scan(&adjustment.ID)
		if err != nil {
			return fmt.Errorf("failed to insert order item adjustment: %w", err)
		}
	} else {
		// Update existing order item adjustment
		query := `
			UPDATE blc_order_item_adjustment SET
				order_item_id = $1, offer_id = $2, adjustment_reason = $3, adjustment_value = $4, 
				applied_to_sale_price = $5
			WHERE order_item_adjustment_id = $6`
		_, err = tx.ExecContext(ctx, query,
			orderItemID, offerID, adjustment.AdjustmentReason, adjustment.AdjustmentValue,
			appliedToSalePrice, adjustment.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update order item adjustment: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves an order item adjustment by its unique identifier.
func (r *OrderItemAdjustmentRepository) FindByID(ctx context.Context, id int64) (*domain.OrderItemAdjustment, error) {
	query := `
		SELECT
			order_item_adjustment_id, order_item_id, offer_id, adjustment_reason, adjustment_value, 
			applied_to_sale_price, created_at
		FROM blc_order_item_adjustment WHERE order_item_adjustment_id = $1`

	var adjustment domain.OrderItemAdjustment
	var offerID sql.NullInt64
	var orderItemID sql.NullInt64
	var appliedToSalePrice sql.NullBool

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&adjustment.ID, &orderItemID, &offerID, &adjustment.AdjustmentReason, &adjustment.AdjustmentValue,
		&appliedToSalePrice, &adjustment.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query order item adjustment by ID: %w", err)
	}

	if offerID.Valid {
		adjustment.OfferID = offerID.Int64
	}
	if orderItemID.Valid {
		adjustment.OrderItemID = orderItemID.Int64
	}
	if appliedToSalePrice.Valid {
		adjustment.AppliedToSalePrice = appliedToSalePrice.Bool
	}

	return &adjustment, nil
}

// FindByOrderItemID retrieves all order item adjustments for a given order item ID.
func (r *OrderItemAdjustmentRepository) FindByOrderItemID(ctx context.Context, orderItemID int64) ([]*domain.OrderItemAdjustment, error) {
	query := `
		SELECT
			order_item_adjustment_id, order_item_id, offer_id, adjustment_reason, adjustment_value, 
			applied_to_sale_price, created_at
		FROM blc_order_item_adjustment WHERE order_item_id = $1`

	rows, err := r.db.QueryContext(ctx, query, orderItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order item adjustments by order item ID: %w", err)
	}
	defer rows.Close()

	var adjustments []*domain.OrderItemAdjustment
	for rows.Next() {
		var adjustment domain.OrderItemAdjustment
		var offerID sql.NullInt64
		var ordItemID sql.NullInt64 // Use different name to avoid conflict
		var appliedToSalePrice sql.NullBool

		err := rows.Scan(
			&adjustment.ID, &ordItemID, &offerID, &adjustment.AdjustmentReason, &adjustment.AdjustmentValue,
			&appliedToSalePrice, &adjustment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item adjustment row: %w", err)
		}

		if offerID.Valid {
			adjustment.OfferID = offerID.Int64
		}
		if ordItemID.Valid {
			adjustment.OrderItemID = ordItemID.Int64
		}
		if appliedToSalePrice.Valid {
			adjustment.AppliedToSalePrice = appliedToSalePrice.Bool
		}
		adjustments = append(adjustments, &adjustment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for order item adjustments: %w", err)
	}

	return adjustments, nil
}

// Delete removes an order item adjustment by its unique identifier.
func (r *OrderItemAdjustmentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_order_item_adjustment WHERE order_item_adjustment_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order item adjustment: %w", err)
	}
	return nil
}

// DeleteByOrderItemID removes all order item adjustments for a given order item ID.
func (r *OrderItemAdjustmentRepository) DeleteByOrderItemID(ctx context.Context, orderItemID int64) error {
	query := `DELETE FROM blc_order_item_adjustment WHERE order_item_id = $1`
	_, err := r.db.ExecContext(ctx, query, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to delete order item adjustments by order item ID: %w", err)
	}
	return nil
}
