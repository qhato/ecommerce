package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// ProductAttributeRepository implements domain.ProductAttributeRepository for PostgreSQL persistence.
type ProductAttributeRepository struct {
	db *sql.DB
}

// NewProductAttributeRepository creates a new PostgreSQL product attribute repository.
func NewProductAttributeRepository(db *sql.DB) *ProductAttributeRepository {
	return &ProductAttributeRepository{db: db}
}

// Save stores a new product attribute or updates an existing one.
func (r *ProductAttributeRepository) Save(ctx context.Context, attribute *domain.ProductAttribute) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if attribute.ID == 0 {
		// Insert new product attribute
		query := `
			INSERT INTO blc_product_attribute (
				product_id, name, value, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5
			) RETURNING product_attribute_id`
		err = tx.QueryRowContext(ctx, query,
			attribute.ProductID, attribute.Name, attribute.Value, attribute.CreatedAt, attribute.UpdatedAt,
		).Scan(&attribute.ID)
		if err != nil {
			return fmt.Errorf("failed to insert product attribute: %w", err)
		}
	} else {
		// Update existing product attribute
		query := `
			UPDATE blc_product_attribute SET
				product_id = $1, name = $2, value = $3, updated_at = $4
			WHERE product_attribute_id = $5`
		_, err = tx.ExecContext(ctx, query,
			attribute.ProductID, attribute.Name, attribute.Value, attribute.UpdatedAt, attribute.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update product attribute: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a product attribute by its unique identifier.
func (r *ProductAttributeRepository) FindByID(ctx context.Context, id int64) (*domain.ProductAttribute, error) {
	query := `
		SELECT
			product_attribute_id, product_id, name, value, created_at, updated_at
		FROM blc_product_attribute WHERE product_attribute_id = $1`

	var attribute domain.ProductAttribute

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&attribute.ID, &attribute.ProductID, &attribute.Name, &attribute.Value, &attribute.CreatedAt, &attribute.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query product attribute by ID: %w", err)
	}

	return &attribute, nil
}

// FindByProductID retrieves all product attributes for a given product ID.
func (r *ProductAttributeRepository) FindByProductID(ctx context.Context, productID int64) ([]*domain.ProductAttribute, error) {
	query := `
		SELECT
			product_attribute_id, product_id, name, value, created_at, updated_at
		FROM blc_product_attribute WHERE product_id = $1`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query product attributes by product ID: %w", err)
	}
	defer rows.Close()

	var attributes []*domain.ProductAttribute
	for rows.Next() {
		var attribute domain.ProductAttribute
		err := rows.Scan(
			&attribute.ID, &attribute.ProductID, &attribute.Name, &attribute.Value, &attribute.CreatedAt, &attribute.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product attribute row: %w", err)
		}
		attributes = append(attributes, &attribute)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return attributes, nil
}

// Delete removes a product attribute by its unique identifier.
func (r *ProductAttributeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_product_attribute WHERE product_attribute_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product attribute: %w", err)
	}
	return nil
}

// DeleteByProductID removes all product attributes for a given product ID.
func (r *ProductAttributeRepository) DeleteByProductID(ctx context.Context, productID int64) error {
	query := `DELETE FROM blc_product_attribute WHERE product_id = $1`
	_, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product attributes by product ID: %w", err)
	}
	return nil
}
