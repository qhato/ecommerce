package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// ProductOptionValueRepository implements domain.ProductOptionValueRepository for PostgreSQL persistence.
type ProductOptionValueRepository struct {
	db *sql.DB
}

// NewProductOptionValueRepository creates a new PostgreSQL product option value repository.
func NewProductOptionValueRepository(db *sql.DB) *ProductOptionValueRepository {
	return &ProductOptionValueRepository{db: db}
}

// Save stores a new product option value or updates an existing one.
func (r *ProductOptionValueRepository) Save(ctx context.Context, value *domain.ProductOptionValue) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle nullable fields
	attributeValue := sql.NullString{String: value.AttributeValue, Valid: value.AttributeValue != ""}
	displayOrder := sql.NullInt32{Int32: int32(value.DisplayOrder), Valid: true}
	// If displayOrder is 0, Broadleaf's default might be NULL, adjust Valid accordingly
	if value.DisplayOrder == 0 {
		displayOrder.Valid = false
	}
	priceAdjustment := sql.NullFloat64{Float64: value.PriceAdjustment, Valid: true}
	// If priceAdjustment is 0.0, Broadleaf's default might be NULL, adjust Valid accordingly
	if value.PriceAdjustment == 0.0 {
		priceAdjustment.Valid = false
	}

	if value.ID == 0 {
		// Insert new product option value
		query := `
			INSERT INTO blc_product_option_value (
				product_option_id, attribute_value, display_order, price_adjustment, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6
			) RETURNING product_option_value_id`
		err = tx.QueryRowContext(ctx, query,
			value.ProductOptionID, attributeValue, displayOrder, priceAdjustment, value.CreatedAt, value.UpdatedAt,
		).Scan(&value.ID)
		if err != nil {
			return fmt.Errorf("failed to insert product option value: %w", err)
		}
	} else {
		// Update existing product option value
		query := `
			UPDATE blc_product_option_value SET
				product_option_id = $1, attribute_value = $2, display_order = $3, 
				price_adjustment = $4, updated_at = $5
			WHERE product_option_value_id = $6`
		_, err = tx.ExecContext(ctx, query,
			value.ProductOptionID, attributeValue, displayOrder, priceAdjustment, value.UpdatedAt, value.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update product option value: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a product option value by its unique identifier.
func (r *ProductOptionValueRepository) FindByID(ctx context.Context, id int64) (*domain.ProductOptionValue, error) {
	query := `
		SELECT
			product_option_value_id, product_option_id, attribute_value, display_order, 
			price_adjustment, created_at, updated_at
		FROM blc_product_option_value WHERE product_option_value_id = $1`

	var value domain.ProductOptionValue
	var productOptionID sql.NullInt64
	var attributeValue sql.NullString
	var displayOrder sql.NullInt32
	var priceAdjustment sql.NullFloat64

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&value.ID, &productOptionID, &attributeValue, &displayOrder,
		&priceAdjustment, &value.CreatedAt, &value.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query product option value by ID: %w", err)
	}

	if productOptionID.Valid {
		value.ProductOptionID = productOptionID.Int64
	}
	if attributeValue.Valid {
		value.AttributeValue = attributeValue.String
	}
	if displayOrder.Valid {
		value.DisplayOrder = int(displayOrder.Int32)
	}
	if priceAdjustment.Valid {
		value.PriceAdjustment = priceAdjustment.Float64
	}

	return &value, nil
}

// FindByProductOptionID retrieves all product option values for a given product option ID.
func (r *ProductOptionValueRepository) FindByProductOptionID(ctx context.Context, productOptionID int64) ([]*domain.ProductOptionValue, error) {
	query := `
		SELECT
			product_option_value_id, product_option_id, attribute_value, display_order, 
			price_adjustment, created_at, updated_at
		FROM blc_product_option_value WHERE product_option_id = $1`

	rows, err := r.db.QueryContext(ctx, query, productOptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query product option values by product option ID: %w", err)
	}
	defer rows.Close()

	var values []*domain.ProductOptionValue
	for rows.Next() {
		var value domain.ProductOptionValue
		var poID sql.NullInt64
		var attributeValue sql.NullString
		var displayOrder sql.NullInt32
		var priceAdjustment sql.NullFloat64

		err := rows.Scan(
			&value.ID, &poID, &attributeValue, &displayOrder,
			&priceAdjustment, &value.CreatedAt, &value.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product option value row: %w", err)
		}

		if poID.Valid {
			value.ProductOptionID = poID.Int64
		}
		if attributeValue.Valid {
			value.AttributeValue = attributeValue.String
		}
		if displayOrder.Valid {
			value.DisplayOrder = int(displayOrder.Int32)
		}
		if priceAdjustment.Valid {
			value.PriceAdjustment = priceAdjustment.Float64
		}
		values = append(values, &value)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return values, nil
}

// Delete removes a product option value by its unique identifier.
func (r *ProductOptionValueRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_product_option_value WHERE product_option_value_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product option value: %w", err)
	}
	return nil
}

// DeleteByProductOptionID removes all product option values for a given product option ID.
func (r *ProductOptionValueRepository) DeleteByProductOptionID(ctx context.Context, productOptionID int64) error {
	query := `DELETE FROM blc_product_option_value WHERE product_option_id = $1`
	_, err := r.db.ExecContext(ctx, query, productOptionID)
	if err != nil {
		return fmt.Errorf("failed to delete product option values by product option ID: %w", err)
	}
	return nil
}
