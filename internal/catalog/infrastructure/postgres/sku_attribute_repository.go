package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// SKUAttributeRepository implements domain.SKUAttributeRepository for PostgreSQL persistence.
type SKUAttributeRepository struct {
	db *sql.DB
}

// NewSKUAttributeRepository creates a new PostgreSQL SKU attribute repository.
func NewSKUAttributeRepository(db *sql.DB) *SKUAttributeRepository {
	return &SKUAttributeRepository{db: db}
}

// Save stores a new SKU attribute or updates an existing one.
func (r *SKUAttributeRepository) Save(ctx context.Context, attribute *domain.SKUAttribute) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if attribute.ID == 0 {
		// Insert new SKU attribute
		query := `
			INSERT INTO blc_sku_attribute (
				sku_id, name, value, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5
			) RETURNING sku_attr_id`
		err = tx.QueryRowContext(ctx, query,
			attribute.SKUID, attribute.Name, attribute.Value, attribute.CreatedAt, attribute.UpdatedAt,
		).Scan(&attribute.ID)
		if err != nil {
			return fmt.Errorf("failed to insert SKU attribute: %w", err)
		}
	} else {
		// Update existing SKU attribute
		query := `
			UPDATE blc_sku_attribute SET
				sku_id = $1, name = $2, value = $3, updated_at = $4
			WHERE sku_attr_id = $5`
		_, err = tx.ExecContext(ctx, query,
			attribute.SKUID, attribute.Name, attribute.Value, attribute.UpdatedAt, attribute.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update SKU attribute: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a SKU attribute by its unique identifier.
func (r *SKUAttributeRepository) FindByID(ctx context.Context, id int64) (*domain.SKUAttribute, error) {
	query := `
		SELECT
			sku_attr_id, sku_id, name, value, created_at, updated_at
		FROM blc_sku_attribute WHERE sku_attr_id = $1`

	var attribute domain.SKUAttribute

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&attribute.ID, &attribute.SKUID, &attribute.Name, &attribute.Value, &attribute.CreatedAt, &attribute.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query SKU attribute by ID: %w", err)
	}

	return &attribute, nil
}

// FindBySKUID retrieves all SKU attributes for a given SKU ID.
func (r *SKUAttributeRepository) FindBySKUID(ctx context.Context, skuID int64) ([]*domain.SKUAttribute, error) {
	query := `
		SELECT
			sku_attr_id, sku_id, name, value, created_at, updated_at
		FROM blc_sku_attribute WHERE sku_id = $1`

	rows, err := r.db.QueryContext(ctx, query, skuID)
	if err != nil {
		return nil, fmt.Errorf("failed to query SKU attributes by SKU ID: %w", err)
	}
	defer rows.Close()

	var attributes []*domain.SKUAttribute
	for rows.Next() {
		var attribute domain.SKUAttribute
		err := rows.Scan(
			&attribute.ID, &attribute.SKUID, &attribute.Name, &attribute.Value, &attribute.CreatedAt, &attribute.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan SKU attribute row: %w", err)
		}
		attributes = append(attributes, &attribute)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return attributes, nil
}

// Delete removes a SKU attribute by its unique identifier.
func (r *SKUAttributeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_sku_attribute WHERE sku_attr_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete SKU attribute: %w", err)
	}
	return nil
}

// DeleteBySKUID removes all SKU attributes for a given SKU ID.
func (r *SKUAttributeRepository) DeleteBySKUID(ctx context.Context, skuID int64) error {
	query := `DELETE FROM blc_sku_attribute WHERE sku_id = $1`
	_, err := r.db.ExecContext(ctx, query, skuID)
	if err != nil {
		return fmt.Errorf("failed to delete SKU attributes by SKU ID: %w", err)
	}
	return nil
}
