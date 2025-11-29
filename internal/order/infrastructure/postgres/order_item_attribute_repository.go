package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/order/domain"
)

// OrderItemAttributeRepository implements domain.OrderItemAttributeRepository for PostgreSQL persistence.
type OrderItemAttributeRepository struct {
	db *sql.DB
}

// NewOrderItemAttributeRepository creates a new PostgreSQL order item attribute repository.
func NewOrderItemAttributeRepository(db *sql.DB) *OrderItemAttributeRepository {
	return &OrderItemAttributeRepository{db: db}
}

// Save stores a new order item attribute or updates an existing one.
// Since blc_order_item_add_attr has a composite primary key (order_item_id, name),
// we treat save as an UPSERT operation where the primary key is used to determine existence.
func (r *OrderItemAttributeRepository) Save(ctx context.Context, attribute *domain.OrderItemAttribute) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Check if the attribute already exists
	existingAttr, err := r.FindByOrderItemIDAndName(ctx, attribute.OrderItemID, attribute.Name)
	if err != nil {
		return fmt.Errorf("failed to check for existing order item attribute: %w", err)
	}

	if existingAttr == nil {
		// Insert new order item attribute
		query := `
			INSERT INTO blc_order_item_add_attr (
				order_item_id, name, value, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5
			)`
		_, err = tx.ExecContext(ctx, query,
			attribute.OrderItemID, attribute.Name, attribute.Value, attribute.CreatedAt, attribute.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert order item attribute: %w", err)
		}
	} else {
		// Update existing order item attribute
		query := `
			UPDATE blc_order_item_add_attr SET
				value = $1, updated_at = $2
			WHERE order_item_id = $3 AND name = $4`
		_, err = tx.ExecContext(ctx, query,
			attribute.Value, attribute.UpdatedAt, attribute.OrderItemID, attribute.Name,
		)
		if err != nil {
			return fmt.Errorf("failed to update order item attribute: %w", err)
		}
	}

	return tx.Commit()
}

// FindByOrderItemIDAndName retrieves an order item attribute by its composite primary key.
func (r *OrderItemAttributeRepository) FindByOrderItemIDAndName(ctx context.Context, orderItemID int64, name string) (*domain.OrderItemAttribute, error) {
	query := `
		SELECT
			order_item_id, name, value, created_at, updated_at
		FROM blc_order_item_add_attr WHERE order_item_id = $1 AND name = $2`

	var attribute domain.OrderItemAttribute

	row := r.db.QueryRowContext(ctx, query, orderItemID, name)
	err := row.Scan(
		&attribute.OrderItemID, &attribute.Name, &attribute.Value, &attribute.CreatedAt, &attribute.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query order item attribute by ID and name: %w", err)
	}

	return &attribute, nil
}

// FindByOrderItemID retrieves all order item attributes for a given order item ID.
func (r *OrderItemAttributeRepository) FindByOrderItemID(ctx context.Context, orderItemID int64) ([]*domain.OrderItemAttribute, error) {
	query := `
		SELECT
			order_item_id, name, value, created_at, updated_at
		FROM blc_order_item_add_attr WHERE order_item_id = $1`

	rows, err := r.db.QueryContext(ctx, query, orderItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order item attributes by order item ID: %w", err)
	}
	defer rows.Close()

	var attributes []*domain.OrderItemAttribute
	for rows.Next() {
		var attribute domain.OrderItemAttribute
		err := rows.Scan(
			&attribute.OrderItemID, &attribute.Name, &attribute.Value, &attribute.CreatedAt, &attribute.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item attribute row: %w", err)
		}
		attributes = append(attributes, &attribute)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for order item attributes: %w", err)
	}

	return attributes, nil
}

// Delete removes an order item attribute by its composite primary key.
func (r *OrderItemAttributeRepository) Delete(ctx context.Context, orderItemID int64, name string) error {
	query := `DELETE FROM blc_order_item_add_attr WHERE order_item_id = $1 AND name = $2`
	_, err := r.db.ExecContext(ctx, query, orderItemID, name)
	if err != nil {
		return fmt.Errorf("failed to delete order item attribute: %w", err)
	}
	return nil
}

// DeleteByOrderItemID removes all order item attributes for a given order item ID.
func (r *OrderItemAttributeRepository) DeleteByOrderItemID(ctx context.Context, orderItemID int64) error {
	query := `DELETE FROM blc_order_item_add_attr WHERE order_item_id = $1`
	_, err := r.db.ExecContext(ctx, query, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to delete order item attributes by order item ID: %w", err)
	}
	return nil
}
