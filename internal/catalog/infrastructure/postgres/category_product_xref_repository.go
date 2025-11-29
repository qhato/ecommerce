package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// CategoryProductXrefRepository implements domain.CategoryProductXrefRepository for PostgreSQL persistence.
type CategoryProductXrefRepository struct {
	db *sql.DB
}

// NewCategoryProductXrefRepository creates a new PostgreSQL category-product xref repository.
func NewCategoryProductXrefRepository(db *sql.DB) *CategoryProductXrefRepository {
	return &CategoryProductXrefRepository{db: db}
}

// Save stores a new category-product cross-reference or updates an existing one.
func (r *CategoryProductXrefRepository) Save(ctx context.Context, xref *domain.CategoryProductXref) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	defaultReference := sql.NullBool{Bool: xref.DefaultReference, Valid: true}
	displayOrder := sql.NullFloat64{Float64: xref.DisplayOrder, Valid: true}


	if xref.ID == 0 {
		// Insert new category-product xref
		query := `
			INSERT INTO blc_category_product_xref (
				category_id, product_id, default_reference, display_order, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6
			) RETURNING category_product_id`
		err = tx.QueryRowContext(ctx, query,
			xref.CategoryID, xref.ProductID, defaultReference, displayOrder, xref.CreatedAt, xref.UpdatedAt,
		).Scan(&xref.ID)
		if err != nil {
			return fmt.Errorf("failed to insert category-product xref: %w", err)
		}
	} else {
		// Update existing category-product xref
		query := `
			UPDATE blc_category_product_xref SET
				category_id = $1, product_id = $2, default_reference = $3, display_order = $4, updated_at = $5
			WHERE category_product_id = $6`
		_, err = tx.ExecContext(ctx, query,
			xref.CategoryID, xref.ProductID, defaultReference, displayOrder, xref.UpdatedAt, xref.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update category-product xref: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a category-product cross-reference by its unique identifier.
func (r *CategoryProductXrefRepository) FindByID(ctx context.Context, id int64) (*domain.CategoryProductXref, error) {
	query := `
		SELECT
			category_product_id, category_id, product_id, default_reference, display_order, created_at, updated_at
		FROM blc_category_product_xref WHERE category_product_id = $1`

	var xref domain.CategoryProductXref
	var defaultReference sql.NullBool
	var displayOrder sql.NullFloat64

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&xref.ID, &xref.CategoryID, &xref.ProductID, &defaultReference, &displayOrder, &xref.CreatedAt, &xref.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query category-product xref by ID: %w", err)
	}

	if defaultReference.Valid {
		xref.DefaultReference = defaultReference.Bool
	}
	if displayOrder.Valid {
		xref.DisplayOrder = displayOrder.Float64
	}

	return &xref, nil
}

// FindByCategoryID retrieves all category-product cross-references for a given category ID.
func (r *CategoryProductXrefRepository) FindByCategoryID(ctx context.Context, categoryID int64) ([]*domain.CategoryProductXref, error) {
	query := `
		SELECT
			category_product_id, category_id, product_id, default_reference, display_order, created_at, updated_at
		FROM blc_category_product_xref WHERE category_id = $1`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query category-product xrefs by category ID: %w", err)
	}
	defer rows.Close()

	var xrefs []*domain.CategoryProductXref
	for rows.Next() {
		var xref domain.CategoryProductXref
		var defaultReference sql.NullBool
		var displayOrder sql.NullFloat64

		err := rows.Scan(
			&xref.ID, &xref.CategoryID, &xref.ProductID, &defaultReference, &displayOrder, &xref.CreatedAt, &xref.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category-product xref row: %w", err)
		}
		if defaultReference.Valid {
			xref.DefaultReference = defaultReference.Bool
		}
		if displayOrder.Valid {
			xref.DisplayOrder = displayOrder.Float64
		}
		xrefs = append(xrefs, &xref)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return xrefs, nil
}

// FindByProductID retrieves all category-product cross-references for a given product ID.
func (r *CategoryProductXrefRepository) FindByProductID(ctx context.Context, productID int64) ([]*domain.CategoryProductXref, error) {
	query := `
		SELECT
			category_product_id, category_id, product_id, default_reference, display_order, created_at, updated_at
		FROM blc_category_product_xref WHERE product_id = $1`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query category-product xrefs by product ID: %w", err)
	}
	defer rows.Close()

	var xrefs []*domain.CategoryProductXref
	for rows.Next() {
		var xref domain.CategoryProductXref
		var defaultReference sql.NullBool
		var displayOrder sql.NullFloat64

		err := rows.Scan(
			&xref.ID, &xref.CategoryID, &xref.ProductID, &defaultReference, &displayOrder, &xref.CreatedAt, &xref.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category-product xref row: %w", err)
		}
		if defaultReference.Valid {
			xref.DefaultReference = defaultReference.Bool
		}
		if displayOrder.Valid {
			xref.DisplayOrder = displayOrder.Float64
		}
		xrefs = append(xrefs, &xref)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return xrefs, nil
}

// Delete removes a category-product cross-reference by its unique identifier.
func (r *CategoryProductXrefRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_category_product_xref WHERE category_product_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category-product xref: %w", err)
	}
	return nil
}

// RemoveCategoryProductXref removes a specific category-product cross-reference by category ID and product ID.
func (r *CategoryProductXrefRepository) RemoveCategoryProductXref(ctx context.Context, categoryID, productID int64) error {
	query := `DELETE FROM blc_category_product_xref WHERE category_id = $1 AND product_id = $2`
	_, err := r.db.ExecContext(ctx, query, categoryID, productID)
	if err != nil {
		return fmt.Errorf("failed to remove category-product xref: %w", err)
	}
	return nil
}
