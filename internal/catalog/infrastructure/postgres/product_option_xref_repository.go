package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// ProductOptionXrefRepository implements domain.ProductOptionXrefRepository for PostgreSQL persistence.
type ProductOptionXrefRepository struct {
	db *sql.DB
}

// NewProductOptionXrefRepository creates a new PostgreSQL product option xref repository.
func NewProductOptionXrefRepository(db *sql.DB) *ProductOptionXrefRepository {
	return &ProductOptionXrefRepository{db: db}
}

// Save stores a new product option cross-reference or updates an existing one.
func (r *ProductOptionXrefRepository) Save(ctx context.Context, xref *domain.ProductOptionXref) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if xref.ID == 0 {
		// Insert new product option xref
		query := `
			INSERT INTO blc_product_option_xref (
				product_id, product_option_id, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4
			) RETURNING product_option_xref_id`
		err = tx.QueryRowContext(ctx, query,
			xref.ProductID, xref.ProductOptionID, xref.CreatedAt, xref.UpdatedAt,
		).Scan(&xref.ID)
		if err != nil {
			return fmt.Errorf("failed to insert product option xref: %w", err)
		}
	} else {
		// Update existing product option xref
		query := `
			UPDATE blc_product_option_xref SET
				product_id = $1, product_option_id = $2, updated_at = $3
			WHERE product_option_xref_id = $4`
		_, err = tx.ExecContext(ctx, query,
			xref.ProductID, xref.ProductOptionID, xref.UpdatedAt, xref.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update product option xref: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a product option cross-reference by its unique identifier.
func (r *ProductOptionXrefRepository) FindByID(ctx context.Context, id int64) (*domain.ProductOptionXref, error) {
	query := `
		SELECT
			product_option_xref_id, product_id, product_option_id, created_at, updated_at
		FROM blc_product_option_xref WHERE product_option_xref_id = $1`

	var xref domain.ProductOptionXref

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&xref.ID, &xref.ProductID, &xref.ProductOptionID, &xref.CreatedAt, &xref.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query product option xref by ID: %w", err)
	}

	return &xref, nil
}

// FindByProductID retrieves all product option cross-references for a given product ID.
func (r *ProductOptionXrefRepository) FindByProductID(ctx context.Context, productID int64) ([]*domain.ProductOptionXref, error) {
	query := `
		SELECT
			product_option_xref_id, product_id, product_option_id, created_at, updated_at
		FROM blc_product_option_xref WHERE product_id = $1`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query product option xrefs by product ID: %w", err)
	}
	defer rows.Close()

	var xrefs []*domain.ProductOptionXref
	for rows.Next() {
		var xref domain.ProductOptionXref
		err := rows.Scan(
			&xref.ID, &xref.ProductID, &xref.ProductOptionID, &xref.CreatedAt, &xref.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product option xref row: %w", err)
		}
		xrefs = append(xrefs, &xref)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return xrefs, nil
}

// FindByProductOptionID retrieves all product option cross-references for a given product option ID.
func (r *ProductOptionXrefRepository) FindByProductOptionID(ctx context.Context, productOptionID int64) ([]*domain.ProductOptionXref, error) {
	query := `
		SELECT
			product_option_xref_id, product_id, product_option_id, created_at, updated_at
		FROM blc_product_option_xref WHERE product_option_id = $1`

	rows, err := r.db.QueryContext(ctx, query, productOptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query product option xrefs by product option ID: %w", err)
	}
	defer rows.Close()

	var xrefs []*domain.ProductOptionXref
	for rows.Next() {
		var xref domain.ProductOptionXref
		err := rows.Scan(
			&xref.ID, &xref.ProductID, &xref.ProductOptionID, &xref.CreatedAt, &xref.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product option xref row: %w", err)
		}
		xrefs = append(xrefs, &xref)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return xrefs, nil
}

// Delete removes a product option cross-reference by its unique identifier.
func (r *ProductOptionXrefRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_product_option_xref WHERE product_option_xref_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product option xref: %w", err)
	}
	return nil
}

// DeleteByProductID removes all product option cross-references for a given product ID.
func (r *ProductOptionXrefRepository) DeleteByProductID(ctx context.Context, productID int64) error {
	query := `DELETE FROM blc_product_option_xref WHERE product_id = $1`
	_, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product option xrefs by product ID: %w", err)
	}
	return nil
}

// DeleteByProductOptionID removes all product option cross-references for a given product option ID.
func (r *ProductOptionXrefRepository) DeleteByProductOptionID(ctx context.Context, productOptionID int64) error {
	query := `DELETE FROM blc_product_option_xref WHERE product_option_id = $1`
	_, err := r.db.ExecContext(ctx, query, productOptionID)
	if err != nil {
		return fmt.Errorf("failed to delete product option xrefs by product option ID: %w", err)
	}
	return nil
}

// RemoveProductOptionXref removes a specific product option cross-reference by product ID and product option ID.
func (r *ProductOptionXrefRepository) RemoveProductOptionXref(ctx context.Context, productID, productOptionID int64) error {
	query := `DELETE FROM blc_product_option_xref WHERE product_id = $1 AND product_option_id = $2`
	_, err := r.db.ExecContext(ctx, query, productID, productOptionID)
	if err != nil {
		return fmt.Errorf("failed to remove product option xref: %w", err)
	}
	return nil
}
