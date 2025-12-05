package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/offer/domain"
)

// PostgresOrderAdjustmentRepository implements OrderAdjustmentRepository using PostgreSQL
type PostgresOrderAdjustmentRepository struct {
	db *sql.DB
}

// NewPostgresOrderAdjustmentRepository creates a new PostgresOrderAdjustmentRepository
func NewPostgresOrderAdjustmentRepository(db *sql.DB) domain.OrderAdjustmentRepository {
	return &PostgresOrderAdjustmentRepository{db: db}
}

func (r *PostgresOrderAdjustmentRepository) CreateOrderAdjustment(adj *domain.OrderAdjustment) error {
	query := `
		INSERT INTO blc_order_adjustment (
			order_id, offer_id, offer_name, adjustment_value,
			adjustment_reason, applied_date, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		adj.OrderID,
		adj.OfferID,
		adj.OfferName,
		adj.AdjustmentValue,
		adj.AdjustmentReason,
		adj.AppliedDate,
		adj.CreatedAt,
	).Scan(&adj.ID)

	if err != nil {
		return fmt.Errorf("failed to create order adjustment: %w", err)
	}
	return nil
}

func (r *PostgresOrderAdjustmentRepository) CreateOrderItemAdjustment(adj *domain.OrderItemAdjustment) error {
	query := `
		INSERT INTO blc_order_item_adjustment (
			order_item_id, offer_id, offer_name, adjustment_value,
			quantity, applied_date, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		adj.OrderItemID,
		adj.OfferID,
		adj.OfferName,
		adj.AdjustmentValue,
		adj.Quantity,
		adj.AppliedDate,
		adj.CreatedAt,
	).Scan(&adj.ID)

	if err != nil {
		return fmt.Errorf("failed to create order item adjustment: %w", err)
	}
	return nil
}

func (r *PostgresOrderAdjustmentRepository) CreateFulfillmentAdjustment(adj *domain.FulfillmentGroupAdjustment) error {
	query := `
		INSERT INTO blc_fulfillment_group_adjustment (
			fulfillment_group_id, offer_id, offer_name, adjustment_value,
			adjustment_reason, applied_date, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		adj.FulfillmentGroupID,
		adj.OfferID,
		adj.OfferName,
		adj.AdjustmentValue,
		adj.AdjustmentReason,
		adj.AppliedDate,
		adj.CreatedAt,
	).Scan(&adj.ID)

	if err != nil {
		return fmt.Errorf("failed to create fulfillment group adjustment: %w", err)
	}
	return nil
}

func (r *PostgresOrderAdjustmentRepository) FindByOrderID(orderID int64) ([]*domain.OrderAdjustment, error) {
	query := `
		SELECT id, order_id, offer_id, offer_name, adjustment_value,
		       adjustment_reason, applied_date, created_at
		FROM blc_order_adjustment
		WHERE order_id = $1
		ORDER BY applied_date DESC
	`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order adjustments: %w", err)
	}
	defer rows.Close()

	adjustments := make([]*domain.OrderAdjustment, 0)
	for rows.Next() {
		adj := &domain.OrderAdjustment{}
		err := rows.Scan(
			&adj.ID,
			&adj.OrderID,
			&adj.OfferID,
			&adj.OfferName,
			&adj.AdjustmentValue,
			&adj.AdjustmentReason,
			&adj.AppliedDate,
			&adj.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order adjustment: %w", err)
		}
		adjustments = append(adjustments, adj)
	}

	return adjustments, nil
}

func (r *PostgresOrderAdjustmentRepository) FindItemAdjustmentsByOrderID(orderID int64) ([]*domain.OrderItemAdjustment, error) {
	query := `
		SELECT oia.id, oia.order_item_id, oia.offer_id, oia.offer_name,
		       oia.adjustment_value, oia.quantity, oia.applied_date, oia.created_at
		FROM blc_order_item_adjustment oia
		INNER JOIN blc_order_item oi ON oia.order_item_id = oi.id
		WHERE oi.order_id = $1
		ORDER BY oia.applied_date DESC
	`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order item adjustments: %w", err)
	}
	defer rows.Close()

	adjustments := make([]*domain.OrderItemAdjustment, 0)
	for rows.Next() {
		adj := &domain.OrderItemAdjustment{}
		err := rows.Scan(
			&adj.ID,
			&adj.OrderItemID,
			&adj.OfferID,
			&adj.OfferName,
			&adj.AdjustmentValue,
			&adj.Quantity,
			&adj.AppliedDate,
			&adj.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item adjustment: %w", err)
		}
		adjustments = append(adjustments, adj)
	}

	return adjustments, nil
}

func (r *PostgresOrderAdjustmentRepository) FindFulfillmentAdjustmentsByOrderID(orderID int64) ([]*domain.FulfillmentGroupAdjustment, error) {
	query := `
		SELECT fga.id, fga.fulfillment_group_id, fga.offer_id, fga.offer_name,
		       fga.adjustment_value, fga.adjustment_reason, fga.applied_date, fga.created_at
		FROM blc_fulfillment_group_adjustment fga
		INNER JOIN blc_fulfillment_group fg ON fga.fulfillment_group_id = fg.id
		WHERE fg.order_id = $1
		ORDER BY fga.applied_date DESC
	`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find fulfillment group adjustments: %w", err)
	}
	defer rows.Close()

	adjustments := make([]*domain.FulfillmentGroupAdjustment, 0)
	for rows.Next() {
		adj := &domain.FulfillmentGroupAdjustment{}
		err := rows.Scan(
			&adj.ID,
			&adj.FulfillmentGroupID,
			&adj.OfferID,
			&adj.OfferName,
			&adj.AdjustmentValue,
			&adj.AdjustmentReason,
			&adj.AppliedDate,
			&adj.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fulfillment group adjustment: %w", err)
		}
		adjustments = append(adjustments, adj)
	}

	return adjustments, nil
}

func (r *PostgresOrderAdjustmentRepository) DeleteByOrderID(orderID int64) error {
	ctx := context.Background()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete order-level adjustments
	_, err = tx.Exec("DELETE FROM blc_order_adjustment WHERE order_id = $1", orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order adjustments: %w", err)
	}

	// Delete order item adjustments
	_, err = tx.Exec(`
		DELETE FROM blc_order_item_adjustment
		WHERE order_item_id IN (
			SELECT id FROM blc_order_item WHERE order_id = $1
		)
	`, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order item adjustments: %w", err)
	}

	// Delete fulfillment group adjustments
	_, err = tx.Exec(`
		DELETE FROM blc_fulfillment_group_adjustment
		WHERE fulfillment_group_id IN (
			SELECT id FROM blc_fulfillment_group WHERE order_id = $1
		)
	`, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete fulfillment group adjustments: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
