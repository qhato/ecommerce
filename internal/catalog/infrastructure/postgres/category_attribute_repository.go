package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// CategoryAttributeRepository implements domain.CategoryAttributeRepository for PostgreSQL persistence.
type CategoryAttributeRepository struct {
	db *sql.DB
}

// NewCategoryAttributeRepository creates a new PostgreSQL category attribute repository.
func NewCategoryAttributeRepository(db *sql.DB) *CategoryAttributeRepository {
	return &CategoryAttributeRepository{db: db}
}

// Save stores a new category attribute or updates an existing one.
func (r *CategoryAttributeRepository) Save(ctx context.Context, attribute *domain.CategoryAttribute) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if attribute.ID == 0 {
		// Insert new category attribute
		query := `
			INSERT INTO blc_category_attribute (
				category_id, name, value, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5
			) RETURNING category_attribute_id`
		err = tx.QueryRowContext(ctx, query,
			attribute.CategoryID, attribute.Name, attribute.Value, attribute.CreatedAt, attribute.UpdatedAt,
		).Scan(&attribute.ID)
		if err != nil {
			return fmt.Errorf("failed to insert category attribute: %w", err)
		}
	} else {
		// Update existing category attribute
		query := `
			UPDATE blc_category_attribute SET
				category_id = $1, name = $2, value = $3, updated_at = $4
			WHERE category_attribute_id = $5`
		_, err = tx.ExecContext(ctx, query,
			attribute.CategoryID, attribute.Name, attribute.Value, attribute.UpdatedAt, attribute.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update category attribute: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a category attribute by its unique identifier.
func (r *CategoryAttributeRepository) FindByID(ctx context.Context, id int64) (*domain.CategoryAttribute, error) {
	query := `
		SELECT
			category_attribute_id, category_id, name, value, created_at, updated_at
		FROM blc_category_attribute WHERE category_attribute_id = $1`

	var attribute domain.CategoryAttribute

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&attribute.ID, &attribute.CategoryID, &attribute.Name, &attribute.Value, &attribute.CreatedAt, &attribute.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query category attribute by ID: %w", err)
	}

	return &attribute, nil
}

// FindByCategoryID retrieves all category attributes for a given category ID.
func (r *CategoryAttributeRepository) FindByCategoryID(ctx context.Context, categoryID int64) ([]*domain.CategoryAttribute, error) {
	query := `
		SELECT
			category_attribute_id, category_id, name, value, created_at, updated_at
		FROM blc_category_attribute WHERE category_id = $1`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query category attributes by category ID: %w", err)
	}
	defer rows.Close()

	var attributes []*domain.CategoryAttribute
	for rows.Next() {
		var attribute domain.CategoryAttribute
		err := rows.Scan(
			&attribute.ID, &attribute.CategoryID, &attribute.Name, &attribute.Value, &attribute.CreatedAt, &attribute.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category attribute row: %w", err)
		}
		attributes = append(attributes, &attribute)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return attributes, nil
}

// Delete removes a category attribute by its unique identifier.
func (r *CategoryAttributeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_category_attribute WHERE category_attribute_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category attribute: %w", err)
	}
	return nil
}

// DeleteByCategoryID removes all category attributes for a given category ID.
func (r *CategoryAttributeRepository) DeleteByCategoryID(ctx context.Context, categoryID int64) error {
	query := `DELETE FROM blc_category_attribute WHERE category_id = $1`
	_, err := r.db.ExecContext(ctx, query, categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete category attributes by category ID: %w", err)
	}
	return nil
}
