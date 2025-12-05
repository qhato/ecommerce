package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/tax/domain"
)

// PostgresTaxRateRepository implements domain.TaxRateRepository using PostgreSQL
type PostgresTaxRateRepository struct {
	db *sql.DB
}

// NewPostgresTaxRateRepository creates a new PostgreSQL repository
func NewPostgresTaxRateRepository(db *sql.DB) *PostgresTaxRateRepository {
	return &PostgresTaxRateRepository{db: db}
}

// Create creates a new tax rate
func (r *PostgresTaxRateRepository) Create(ctx context.Context, rate *domain.TaxRate) error {
	query := `
		INSERT INTO blc_tax_rate (
			jurisdiction_id, name, tax_type, rate, tax_category,
			is_compound, is_shipping_taxable, min_threshold, max_threshold,
			priority, is_active, start_date, end_date,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		rate.JurisdictionID,
		rate.Name,
		rate.TaxType,
		rate.Rate,
		rate.TaxCategory,
		rate.IsCompound,
		rate.IsShippingTaxable,
		rate.MinThreshold,
		rate.MaxThreshold,
		rate.Priority,
		rate.IsActive,
		rate.StartDate,
		rate.EndDate,
		rate.CreatedAt,
		rate.UpdatedAt,
	).Scan(&rate.ID)

	if err != nil {
		return fmt.Errorf("failed to insert tax rate: %w", err)
	}

	return nil
}

// Update updates an existing tax rate
func (r *PostgresTaxRateRepository) Update(ctx context.Context, rate *domain.TaxRate) error {
	query := `
		UPDATE blc_tax_rate
		SET name = $1, rate = $2, is_compound = $3, is_shipping_taxable = $4,
		    min_threshold = $5, max_threshold = $6, priority = $7,
		    is_active = $8, start_date = $9, end_date = $10, updated_at = $11
		WHERE id = $12`

	result, err := r.db.ExecContext(
		ctx,
		query,
		rate.Name,
		rate.Rate,
		rate.IsCompound,
		rate.IsShippingTaxable,
		rate.MinThreshold,
		rate.MaxThreshold,
		rate.Priority,
		rate.IsActive,
		rate.StartDate,
		rate.EndDate,
		rate.UpdatedAt,
		rate.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update tax rate: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrTaxRateNotFound
	}

	return nil
}

// FindByID finds a tax rate by ID
func (r *PostgresTaxRateRepository) FindByID(ctx context.Context, id int64) (*domain.TaxRate, error) {
	query := `
		SELECT id, jurisdiction_id, name, tax_type, rate, tax_category,
		       is_compound, is_shipping_taxable, min_threshold, max_threshold,
		       priority, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_tax_rate
		WHERE id = $1`

	rate := &domain.TaxRate{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rate.ID,
		&rate.JurisdictionID,
		&rate.Name,
		&rate.TaxType,
		&rate.Rate,
		&rate.TaxCategory,
		&rate.IsCompound,
		&rate.IsShippingTaxable,
		&rate.MinThreshold,
		&rate.MaxThreshold,
		&rate.Priority,
		&rate.IsActive,
		&rate.StartDate,
		&rate.EndDate,
		&rate.CreatedAt,
		&rate.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find tax rate: %w", err)
	}

	return rate, nil
}

// FindByJurisdiction finds all tax rates for a jurisdiction
func (r *PostgresTaxRateRepository) FindByJurisdiction(ctx context.Context, jurisdictionID int64, activeOnly bool) ([]*domain.TaxRate, error) {
	query := `
		SELECT id, jurisdiction_id, name, tax_type, rate, tax_category,
		       is_compound, is_shipping_taxable, min_threshold, max_threshold,
		       priority, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_tax_rate
		WHERE jurisdiction_id = $1`

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY priority ASC, name ASC"

	rows, err := r.db.QueryContext(ctx, query, jurisdictionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tax rates: %w", err)
	}
	defer rows.Close()

	rates := make([]*domain.TaxRate, 0)
	for rows.Next() {
		rate := &domain.TaxRate{}
		err := rows.Scan(
			&rate.ID,
			&rate.JurisdictionID,
			&rate.Name,
			&rate.TaxType,
			&rate.Rate,
			&rate.TaxCategory,
			&rate.IsCompound,
			&rate.IsShippingTaxable,
			&rate.MinThreshold,
			&rate.MaxThreshold,
			&rate.Priority,
			&rate.IsActive,
			&rate.StartDate,
			&rate.EndDate,
			&rate.CreatedAt,
			&rate.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tax rate: %w", err)
		}
		rates = append(rates, rate)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tax rates: %w", err)
	}

	return rates, nil
}

// FindByJurisdictionAndCategory finds tax rates for a jurisdiction and category
func (r *PostgresTaxRateRepository) FindByJurisdictionAndCategory(ctx context.Context, jurisdictionID int64, category domain.TaxCategory, activeOnly bool) ([]*domain.TaxRate, error) {
	query := `
		SELECT id, jurisdiction_id, name, tax_type, rate, tax_category,
		       is_compound, is_shipping_taxable, min_threshold, max_threshold,
		       priority, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_tax_rate
		WHERE jurisdiction_id = $1 AND tax_category = $2`

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY priority ASC, name ASC"

	rows, err := r.db.QueryContext(ctx, query, jurisdictionID, category)
	if err != nil {
		return nil, fmt.Errorf("failed to query tax rates: %w", err)
	}
	defer rows.Close()

	rates := make([]*domain.TaxRate, 0)
	for rows.Next() {
		rate := &domain.TaxRate{}
		err := rows.Scan(
			&rate.ID,
			&rate.JurisdictionID,
			&rate.Name,
			&rate.TaxType,
			&rate.Rate,
			&rate.TaxCategory,
			&rate.IsCompound,
			&rate.IsShippingTaxable,
			&rate.MinThreshold,
			&rate.MaxThreshold,
			&rate.Priority,
			&rate.IsActive,
			&rate.StartDate,
			&rate.EndDate,
			&rate.CreatedAt,
			&rate.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tax rate: %w", err)
		}
		rates = append(rates, rate)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tax rates: %w", err)
	}

	return rates, nil
}

// FindApplicableRates finds all applicable tax rates for a calculation
func (r *PostgresTaxRateRepository) FindApplicableRates(ctx context.Context, jurisdictionIDs []int64, category domain.TaxCategory, activeOnly bool) ([]*domain.TaxRate, error) {
	if len(jurisdictionIDs) == 0 {
		return []*domain.TaxRate{}, nil
	}

	query := `
		SELECT id, jurisdiction_id, name, tax_type, rate, tax_category,
		       is_compound, is_shipping_taxable, min_threshold, max_threshold,
		       priority, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_tax_rate
		WHERE jurisdiction_id = ANY($1) AND tax_category = $2`

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY priority ASC, name ASC"

	rows, err := r.db.QueryContext(ctx, query, jurisdictionIDs, category)
	if err != nil {
		return nil, fmt.Errorf("failed to query applicable tax rates: %w", err)
	}
	defer rows.Close()

	rates := make([]*domain.TaxRate, 0)
	for rows.Next() {
		rate := &domain.TaxRate{}
		err := rows.Scan(
			&rate.ID,
			&rate.JurisdictionID,
			&rate.Name,
			&rate.TaxType,
			&rate.Rate,
			&rate.TaxCategory,
			&rate.IsCompound,
			&rate.IsShippingTaxable,
			&rate.MinThreshold,
			&rate.MaxThreshold,
			&rate.Priority,
			&rate.IsActive,
			&rate.StartDate,
			&rate.EndDate,
			&rate.CreatedAt,
			&rate.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tax rate: %w", err)
		}
		rates = append(rates, rate)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tax rates: %w", err)
	}

	return rates, nil
}

// FindAll finds all tax rates with optional filters
func (r *PostgresTaxRateRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.TaxRate, error) {
	query := `
		SELECT id, jurisdiction_id, name, tax_type, rate, tax_category,
		       is_compound, is_shipping_taxable, min_threshold, max_threshold,
		       priority, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_tax_rate`

	if activeOnly {
		query += " WHERE is_active = true"
	}

	query += " ORDER BY priority ASC, name ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tax rates: %w", err)
	}
	defer rows.Close()

	rates := make([]*domain.TaxRate, 0)
	for rows.Next() {
		rate := &domain.TaxRate{}
		err := rows.Scan(
			&rate.ID,
			&rate.JurisdictionID,
			&rate.Name,
			&rate.TaxType,
			&rate.Rate,
			&rate.TaxCategory,
			&rate.IsCompound,
			&rate.IsShippingTaxable,
			&rate.MinThreshold,
			&rate.MaxThreshold,
			&rate.Priority,
			&rate.IsActive,
			&rate.StartDate,
			&rate.EndDate,
			&rate.CreatedAt,
			&rate.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tax rate: %w", err)
		}
		rates = append(rates, rate)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tax rates: %w", err)
	}

	return rates, nil
}

// Delete deletes a tax rate
func (r *PostgresTaxRateRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_tax_rate WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tax rate: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrTaxRateNotFound
	}

	return nil
}

// BulkCreate creates multiple tax rates in a transaction
func (r *PostgresTaxRateRepository) BulkCreate(ctx context.Context, rates []*domain.TaxRate) error {
	if len(rates) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO blc_tax_rate (
			jurisdiction_id, name, tax_type, rate, tax_category,
			is_compound, is_shipping_taxable, min_threshold, max_threshold,
			priority, is_active, start_date, end_date,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, rate := range rates {
		err := stmt.QueryRowContext(
			ctx,
			rate.JurisdictionID,
			rate.Name,
			rate.TaxType,
			rate.Rate,
			rate.TaxCategory,
			rate.IsCompound,
			rate.IsShippingTaxable,
			rate.MinThreshold,
			rate.MaxThreshold,
			rate.Priority,
			rate.IsActive,
			rate.StartDate,
			rate.EndDate,
			rate.CreatedAt,
			rate.UpdatedAt,
		).Scan(&rate.ID)

		if err != nil {
			return fmt.Errorf("failed to insert tax rate: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
