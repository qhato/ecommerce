package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// SkuProductOptionValueXrefRepository implements domain.SkuProductOptionValueXrefRepository for PostgreSQL persistence.
type SkuProductOptionValueXrefRepository struct {
	db *sql.DB
}

// NewSkuProductOptionValueXrefRepository creates a new PostgreSQL SKU product option value xref repository.
func NewSkuProductOptionValueXrefRepository(db *sql.DB) *SkuProductOptionValueXrefRepository {
	return &SkuProductOptionValueXrefRepository{db: db}
}

// Save stores a new SKU product option value cross-reference or updates an existing one.
func (r *SkuProductOptionValueXrefRepository) Save(ctx context.Context, xref *domain.SkuProductOptionValueXref) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if xref.ID == 0 {
		// Insert new SKU product option value xref
		query := `
			INSERT INTO blc_sku_option_value_xref (
				sku_id, product_option_value_id, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4
			) RETURNING sku_option_value_xref_id`
		err = tx.QueryRowContext(ctx, query,
			xref.SKUID, xref.ProductOptionValueID, xref.CreatedAt, xref.UpdatedAt,
		).Scan(&xref.ID)
		if err != nil {
			return fmt.Errorf("failed to insert SKU product option value xref: %w", err)
		}
	} else {
		// Update existing SKU product option value xref
		query := `
			UPDATE blc_sku_option_value_xref SET
				sku_id = $1, product_option_value_id = $2, updated_at = $3
			WHERE sku_option_value_xref_id = $4`
		_, err = tx.ExecContext(ctx, query,
			xref.SKUID, xref.ProductOptionValueID, xref.UpdatedAt, xref.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update SKU product option value xref: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a SKU product option value cross-reference by its unique identifier.
func (r *SkuProductOptionValueXrefRepository) FindByID(ctx context.Context, id int64) (*domain.SkuProductOptionValueXref, error) {
	query := `
		SELECT
			sku_option_value_xref_id, sku_id, product_option_value_id, created_at, updated_at
		FROM blc_sku_option_value_xref WHERE sku_option_value_xref_id = $1`

	var xref domain.SkuProductOptionValueXref

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&xref.ID, &xref.SKUID, &xref.ProductOptionValueID, &xref.CreatedAt, &xref.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query SKU product option value xref by ID: %w", err)
	}

	return &xref, nil
}

// FindBySKUID retrieves all SKU product option value cross-references for a given SKU ID.
func (r *SkuProductOptionValueXrefRepository) FindBySKUID(ctx context.Context, skuID int64) ([]*domain.SkuProductOptionValueXref, error) {
	query := `
		SELECT
			sku_option_value_xref_id, sku_id, product_option_value_id, created_at, updated_at
		FROM blc_sku_option_value_xref WHERE sku_id = $1`

	rows, err := r.db.QueryContext(ctx, query, skuID)
	if err != nil {
		return nil, fmt.Errorf("failed to query SKU product option value xrefs by SKU ID: %w", err)
	}
	defer rows.Close()

	var xrefs []*domain.SkuProductOptionValueXref
	for rows.Next() {
		var xref domain.SkuProductOptionValueXref
		err := rows.Scan(
			&xref.ID, &xref.SKUID, &xref.ProductOptionValueID, &xref.CreatedAt, &xref.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan SKU product option value xref row: %w", err)
		}
		xrefs = append(xrefs, &xref)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return xrefs, nil
}

// FindByProductOptionValueID retrieves all SKU product option value cross-references for a given product option value ID.
func (r *SkuProductOptionValueXrefRepository) FindByProductOptionValueID(ctx context.Context, productOptionValueID int64) ([]*domain.SkuProductOptionValueXref, error) {
	query := `
		SELECT
			sku_option_value_xref_id, sku_id, product_option_value_id, created_at, updated_at
		FROM blc_sku_option_value_xref WHERE product_option_value_id = $1`

	rows, err := r.db.QueryContext(ctx, query, productOptionValueID)
	if err != nil {
		return nil, fmt.Errorf("failed to query SKU product option value xrefs by product option value ID: %w", err)
	}
	defer rows.Close()

	var xrefs []*domain.SkuProductOptionValueXref
	for rows.Next() {
		var xref domain.SkuProductOptionValueXref
		err := rows.Scan(
			&xref.ID, &xref.SKUID, &xref.ProductOptionValueID, &xref.CreatedAt, &xref.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan SKU product option value xref row: %w", err)
		}
		xrefs = append(xrefs, &xref)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return xrefs, nil
}

// Delete removes a SKU product option value cross-reference by its unique identifier.
func (r *SkuProductOptionValueXrefRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_sku_option_value_xref WHERE sku_option_value_xref_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete SKU product option value xref: %w", err)
	}
	return nil
}

// DeleteBySKUID removes all SKU product option value cross-references for a given SKU ID.
func (r *SkuProductOptionValueXrefRepository) DeleteBySKUID(ctx context.Context, skuID int64) error {
	query := `DELETE FROM blc_sku_option_value_xref WHERE sku_id = $1`
	_, err := r.db.ExecContext(ctx, query, skuID)
	if err != nil {
		return fmt.Errorf("failed to delete SKU product option value xrefs by SKU ID: %w", err)
	}
	return nil
}

// DeleteByProductOptionValueID removes all SKU product option value cross-references for a given product option value ID.
func (r *SkuProductOptionValueXrefRepository) DeleteByProductOptionValueID(ctx context.Context, productOptionValueID int64) error {
	query := `DELETE FROM blc_sku_option_value_xref WHERE product_option_value_id = $1`
	_, err := r.db.ExecContext(ctx, query, productOptionValueID)
	if err != nil {
		return fmt.Errorf("failed to delete SKU product option value xrefs by product option value ID: %w", err)
	}
	return nil
}

// RemoveSkuProductOptionValueXref removes a specific SKU product option value cross-reference by SKU ID and product option value ID.
func (r *SkuProductOptionValueXrefRepository) RemoveSkuProductOptionValueXref(ctx context.Context, skuID, productOptionValueID int64) error {
	query := `DELETE FROM blc_sku_option_value_xref WHERE sku_id = $1 AND product_option_value_id = $2`
	_, err := r.db.ExecContext(ctx, query, skuID, productOptionValueID)
	if err != nil {
		return fmt.Errorf("failed to remove SKU product option value xref: %w", err)
	}
	return nil
}
