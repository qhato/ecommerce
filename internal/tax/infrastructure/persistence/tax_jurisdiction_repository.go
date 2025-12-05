package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/tax/domain"
)

// PostgresTaxJurisdictionRepository implements domain.TaxJurisdictionRepository using PostgreSQL
type PostgresTaxJurisdictionRepository struct {
	db *sql.DB
}

// NewPostgresTaxJurisdictionRepository creates a new PostgreSQL repository
func NewPostgresTaxJurisdictionRepository(db *sql.DB) *PostgresTaxJurisdictionRepository {
	return &PostgresTaxJurisdictionRepository{db: db}
}

// Create creates a new tax jurisdiction
func (r *PostgresTaxJurisdictionRepository) Create(ctx context.Context, jurisdiction *domain.TaxJurisdiction) error {
	query := `
		INSERT INTO blc_tax_jurisdiction (
			code, name, jurisdiction_type, parent_id, country,
			state_province, county, city, postal_code,
			is_active, priority, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		jurisdiction.Code,
		jurisdiction.Name,
		jurisdiction.JurisdictionType,
		jurisdiction.ParentID,
		jurisdiction.Country,
		jurisdiction.StateProvince,
		jurisdiction.County,
		jurisdiction.City,
		jurisdiction.PostalCode,
		jurisdiction.IsActive,
		jurisdiction.Priority,
		jurisdiction.CreatedAt,
		jurisdiction.UpdatedAt,
	).Scan(&jurisdiction.ID)

	if err != nil {
		return fmt.Errorf("failed to insert jurisdiction: %w", err)
	}

	return nil
}

// Update updates an existing tax jurisdiction
func (r *PostgresTaxJurisdictionRepository) Update(ctx context.Context, jurisdiction *domain.TaxJurisdiction) error {
	query := `
		UPDATE blc_tax_jurisdiction
		SET name = $1, parent_id = $2, state_province = $3, county = $4,
		    city = $5, postal_code = $6, is_active = $7, priority = $8, updated_at = $9
		WHERE id = $10`

	result, err := r.db.ExecContext(
		ctx,
		query,
		jurisdiction.Name,
		jurisdiction.ParentID,
		jurisdiction.StateProvince,
		jurisdiction.County,
		jurisdiction.City,
		jurisdiction.PostalCode,
		jurisdiction.IsActive,
		jurisdiction.Priority,
		jurisdiction.UpdatedAt,
		jurisdiction.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update jurisdiction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrJurisdictionNotFound
	}

	return nil
}

// FindByID finds a jurisdiction by ID
func (r *PostgresTaxJurisdictionRepository) FindByID(ctx context.Context, id int64) (*domain.TaxJurisdiction, error) {
	query := `
		SELECT id, code, name, jurisdiction_type, parent_id, country,
		       state_province, county, city, postal_code,
		       is_active, priority, created_at, updated_at
		FROM blc_tax_jurisdiction
		WHERE id = $1`

	jurisdiction := &domain.TaxJurisdiction{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&jurisdiction.ID,
		&jurisdiction.Code,
		&jurisdiction.Name,
		&jurisdiction.JurisdictionType,
		&jurisdiction.ParentID,
		&jurisdiction.Country,
		&jurisdiction.StateProvince,
		&jurisdiction.County,
		&jurisdiction.City,
		&jurisdiction.PostalCode,
		&jurisdiction.IsActive,
		&jurisdiction.Priority,
		&jurisdiction.CreatedAt,
		&jurisdiction.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find jurisdiction: %w", err)
	}

	return jurisdiction, nil
}

// FindByCode finds a jurisdiction by code
func (r *PostgresTaxJurisdictionRepository) FindByCode(ctx context.Context, code string) (*domain.TaxJurisdiction, error) {
	query := `
		SELECT id, code, name, jurisdiction_type, parent_id, country,
		       state_province, county, city, postal_code,
		       is_active, priority, created_at, updated_at
		FROM blc_tax_jurisdiction
		WHERE code = $1`

	jurisdiction := &domain.TaxJurisdiction{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&jurisdiction.ID,
		&jurisdiction.Code,
		&jurisdiction.Name,
		&jurisdiction.JurisdictionType,
		&jurisdiction.ParentID,
		&jurisdiction.Country,
		&jurisdiction.StateProvince,
		&jurisdiction.County,
		&jurisdiction.City,
		&jurisdiction.PostalCode,
		&jurisdiction.IsActive,
		&jurisdiction.Priority,
		&jurisdiction.CreatedAt,
		&jurisdiction.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find jurisdiction: %w", err)
	}

	return jurisdiction, nil
}

// FindByLocation finds all jurisdictions matching a location
func (r *PostgresTaxJurisdictionRepository) FindByLocation(ctx context.Context, country, stateProvince, county, city, postalCode string) ([]*domain.TaxJurisdiction, error) {
	query := `
		SELECT id, code, name, jurisdiction_type, parent_id, country,
		       state_province, county, city, postal_code,
		       is_active, priority, created_at, updated_at
		FROM blc_tax_jurisdiction
		WHERE country = $1
		  AND (state_province IS NULL OR state_province = $2)
		  AND (county IS NULL OR county = $3)
		  AND (city IS NULL OR city = $4)
		  AND (postal_code IS NULL OR postal_code = $5)
		  AND is_active = true
		ORDER BY priority ASC`

	rows, err := r.db.QueryContext(ctx, query, country, stateProvince, county, city, postalCode)
	if err != nil {
		return nil, fmt.Errorf("failed to query jurisdictions: %w", err)
	}
	defer rows.Close()

	jurisdictions := make([]*domain.TaxJurisdiction, 0)
	for rows.Next() {
		jurisdiction := &domain.TaxJurisdiction{}
		err := rows.Scan(
			&jurisdiction.ID,
			&jurisdiction.Code,
			&jurisdiction.Name,
			&jurisdiction.JurisdictionType,
			&jurisdiction.ParentID,
			&jurisdiction.Country,
			&jurisdiction.StateProvince,
			&jurisdiction.County,
			&jurisdiction.City,
			&jurisdiction.PostalCode,
			&jurisdiction.IsActive,
			&jurisdiction.Priority,
			&jurisdiction.CreatedAt,
			&jurisdiction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan jurisdiction: %w", err)
		}
		jurisdictions = append(jurisdictions, jurisdiction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating jurisdictions: %w", err)
	}

	return jurisdictions, nil
}

// FindAll finds all jurisdictions with optional filters
func (r *PostgresTaxJurisdictionRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.TaxJurisdiction, error) {
	query := `
		SELECT id, code, name, jurisdiction_type, parent_id, country,
		       state_province, county, city, postal_code,
		       is_active, priority, created_at, updated_at
		FROM blc_tax_jurisdiction`

	if activeOnly {
		query += " WHERE is_active = true"
	}

	query += " ORDER BY priority ASC, code ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query jurisdictions: %w", err)
	}
	defer rows.Close()

	jurisdictions := make([]*domain.TaxJurisdiction, 0)
	for rows.Next() {
		jurisdiction := &domain.TaxJurisdiction{}
		err := rows.Scan(
			&jurisdiction.ID,
			&jurisdiction.Code,
			&jurisdiction.Name,
			&jurisdiction.JurisdictionType,
			&jurisdiction.ParentID,
			&jurisdiction.Country,
			&jurisdiction.StateProvince,
			&jurisdiction.County,
			&jurisdiction.City,
			&jurisdiction.PostalCode,
			&jurisdiction.IsActive,
			&jurisdiction.Priority,
			&jurisdiction.CreatedAt,
			&jurisdiction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan jurisdiction: %w", err)
		}
		jurisdictions = append(jurisdictions, jurisdiction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating jurisdictions: %w", err)
	}

	return jurisdictions, nil
}

// FindByCountry finds all jurisdictions in a country
func (r *PostgresTaxJurisdictionRepository) FindByCountry(ctx context.Context, country string, activeOnly bool) ([]*domain.TaxJurisdiction, error) {
	query := `
		SELECT id, code, name, jurisdiction_type, parent_id, country,
		       state_province, county, city, postal_code,
		       is_active, priority, created_at, updated_at
		FROM blc_tax_jurisdiction
		WHERE country = $1`

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY priority ASC, code ASC"

	rows, err := r.db.QueryContext(ctx, query, country)
	if err != nil {
		return nil, fmt.Errorf("failed to query jurisdictions: %w", err)
	}
	defer rows.Close()

	jurisdictions := make([]*domain.TaxJurisdiction, 0)
	for rows.Next() {
		jurisdiction := &domain.TaxJurisdiction{}
		err := rows.Scan(
			&jurisdiction.ID,
			&jurisdiction.Code,
			&jurisdiction.Name,
			&jurisdiction.JurisdictionType,
			&jurisdiction.ParentID,
			&jurisdiction.Country,
			&jurisdiction.StateProvince,
			&jurisdiction.County,
			&jurisdiction.City,
			&jurisdiction.PostalCode,
			&jurisdiction.IsActive,
			&jurisdiction.Priority,
			&jurisdiction.CreatedAt,
			&jurisdiction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan jurisdiction: %w", err)
		}
		jurisdictions = append(jurisdictions, jurisdiction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating jurisdictions: %w", err)
	}

	return jurisdictions, nil
}

// FindChildren finds all child jurisdictions of a parent
func (r *PostgresTaxJurisdictionRepository) FindChildren(ctx context.Context, parentID int64) ([]*domain.TaxJurisdiction, error) {
	query := `
		SELECT id, code, name, jurisdiction_type, parent_id, country,
		       state_province, county, city, postal_code,
		       is_active, priority, created_at, updated_at
		FROM blc_tax_jurisdiction
		WHERE parent_id = $1
		ORDER BY priority ASC, code ASC`

	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query child jurisdictions: %w", err)
	}
	defer rows.Close()

	jurisdictions := make([]*domain.TaxJurisdiction, 0)
	for rows.Next() {
		jurisdiction := &domain.TaxJurisdiction{}
		err := rows.Scan(
			&jurisdiction.ID,
			&jurisdiction.Code,
			&jurisdiction.Name,
			&jurisdiction.JurisdictionType,
			&jurisdiction.ParentID,
			&jurisdiction.Country,
			&jurisdiction.StateProvince,
			&jurisdiction.County,
			&jurisdiction.City,
			&jurisdiction.PostalCode,
			&jurisdiction.IsActive,
			&jurisdiction.Priority,
			&jurisdiction.CreatedAt,
			&jurisdiction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan jurisdiction: %w", err)
		}
		jurisdictions = append(jurisdictions, jurisdiction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating jurisdictions: %w", err)
	}

	return jurisdictions, nil
}

// Delete deletes a jurisdiction
func (r *PostgresTaxJurisdictionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_tax_jurisdiction WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete jurisdiction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrJurisdictionNotFound
	}

	return nil
}

// ExistsByCode checks if a jurisdiction exists with the given code
func (r *PostgresTaxJurisdictionRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blc_tax_jurisdiction WHERE code = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check jurisdiction existence: %w", err)
	}

	return exists, nil
}
