package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/tax/domain"
)

// PostgresTaxExemptionRepository implements domain.TaxExemptionRepository using PostgreSQL
type PostgresTaxExemptionRepository struct {
	db *sql.DB
}

// NewPostgresTaxExemptionRepository creates a new PostgreSQL repository
func NewPostgresTaxExemptionRepository(db *sql.DB) *PostgresTaxExemptionRepository {
	return &PostgresTaxExemptionRepository{db: db}
}

// Create creates a new tax exemption
func (r *PostgresTaxExemptionRepository) Create(ctx context.Context, exemption *domain.TaxExemption) error {
	query := `
		INSERT INTO blc_tax_exemption (
			customer_id, exemption_certificate, jurisdiction_id, tax_category,
			reason, is_active, start_date, end_date,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		exemption.CustomerID,
		exemption.ExemptionCertificate,
		exemption.JurisdictionID,
		exemption.TaxCategory,
		exemption.Reason,
		exemption.IsActive,
		exemption.StartDate,
		exemption.EndDate,
		exemption.CreatedAt,
		exemption.UpdatedAt,
	).Scan(&exemption.ID)

	if err != nil {
		return fmt.Errorf("failed to insert tax exemption: %w", err)
	}

	return nil
}

// Update updates an existing tax exemption
func (r *PostgresTaxExemptionRepository) Update(ctx context.Context, exemption *domain.TaxExemption) error {
	query := `
		UPDATE blc_tax_exemption
		SET jurisdiction_id = $1, tax_category = $2, reason = $3,
		    is_active = $4, start_date = $5, end_date = $6, updated_at = $7
		WHERE id = $8`

	result, err := r.db.ExecContext(
		ctx,
		query,
		exemption.JurisdictionID,
		exemption.TaxCategory,
		exemption.Reason,
		exemption.IsActive,
		exemption.StartDate,
		exemption.EndDate,
		exemption.UpdatedAt,
		exemption.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update tax exemption: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrExemptionNotFound
	}

	return nil
}

// FindByID finds an exemption by ID
func (r *PostgresTaxExemptionRepository) FindByID(ctx context.Context, id int64) (*domain.TaxExemption, error) {
	query := `
		SELECT id, customer_id, exemption_certificate, jurisdiction_id, tax_category,
		       reason, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_tax_exemption
		WHERE id = $1`

	exemption := &domain.TaxExemption{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&exemption.ID,
		&exemption.CustomerID,
		&exemption.ExemptionCertificate,
		&exemption.JurisdictionID,
		&exemption.TaxCategory,
		&exemption.Reason,
		&exemption.IsActive,
		&exemption.StartDate,
		&exemption.EndDate,
		&exemption.CreatedAt,
		&exemption.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find tax exemption: %w", err)
	}

	return exemption, nil
}

// FindByCustomerID finds all exemptions for a customer
func (r *PostgresTaxExemptionRepository) FindByCustomerID(ctx context.Context, customerID string, activeOnly bool) ([]*domain.TaxExemption, error) {
	query := `
		SELECT id, customer_id, exemption_certificate, jurisdiction_id, tax_category,
		       reason, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_tax_exemption
		WHERE customer_id = $1`

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tax exemptions: %w", err)
	}
	defer rows.Close()

	exemptions := make([]*domain.TaxExemption, 0)
	for rows.Next() {
		exemption := &domain.TaxExemption{}
		err := rows.Scan(
			&exemption.ID,
			&exemption.CustomerID,
			&exemption.ExemptionCertificate,
			&exemption.JurisdictionID,
			&exemption.TaxCategory,
			&exemption.Reason,
			&exemption.IsActive,
			&exemption.StartDate,
			&exemption.EndDate,
			&exemption.CreatedAt,
			&exemption.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tax exemption: %w", err)
		}
		exemptions = append(exemptions, exemption)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tax exemptions: %w", err)
	}

	return exemptions, nil
}

// FindActiveExemptions finds all currently active exemptions for a customer
func (r *PostgresTaxExemptionRepository) FindActiveExemptions(ctx context.Context, customerID string) ([]*domain.TaxExemption, error) {
	query := `
		SELECT id, customer_id, exemption_certificate, jurisdiction_id, tax_category,
		       reason, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_tax_exemption
		WHERE customer_id = $1
		  AND is_active = true
		  AND (start_date IS NULL OR start_date <= NOW())
		  AND (end_date IS NULL OR end_date >= NOW())
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query active tax exemptions: %w", err)
	}
	defer rows.Close()

	exemptions := make([]*domain.TaxExemption, 0)
	for rows.Next() {
		exemption := &domain.TaxExemption{}
		err := rows.Scan(
			&exemption.ID,
			&exemption.CustomerID,
			&exemption.ExemptionCertificate,
			&exemption.JurisdictionID,
			&exemption.TaxCategory,
			&exemption.Reason,
			&exemption.IsActive,
			&exemption.StartDate,
			&exemption.EndDate,
			&exemption.CreatedAt,
			&exemption.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tax exemption: %w", err)
		}
		exemptions = append(exemptions, exemption)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tax exemptions: %w", err)
	}

	return exemptions, nil
}

// FindByCustomerAndJurisdiction finds exemptions for a customer in a jurisdiction
func (r *PostgresTaxExemptionRepository) FindByCustomerAndJurisdiction(ctx context.Context, customerID string, jurisdictionID int64, activeOnly bool) ([]*domain.TaxExemption, error) {
	query := `
		SELECT id, customer_id, exemption_certificate, jurisdiction_id, tax_category,
		       reason, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_tax_exemption
		WHERE customer_id = $1 AND (jurisdiction_id = $2 OR jurisdiction_id IS NULL)`

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, customerID, jurisdictionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tax exemptions: %w", err)
	}
	defer rows.Close()

	exemptions := make([]*domain.TaxExemption, 0)
	for rows.Next() {
		exemption := &domain.TaxExemption{}
		err := rows.Scan(
			&exemption.ID,
			&exemption.CustomerID,
			&exemption.ExemptionCertificate,
			&exemption.JurisdictionID,
			&exemption.TaxCategory,
			&exemption.Reason,
			&exemption.IsActive,
			&exemption.StartDate,
			&exemption.EndDate,
			&exemption.CreatedAt,
			&exemption.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tax exemption: %w", err)
		}
		exemptions = append(exemptions, exemption)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tax exemptions: %w", err)
	}

	return exemptions, nil
}

// FindByCertificate finds an exemption by certificate number
func (r *PostgresTaxExemptionRepository) FindByCertificate(ctx context.Context, certificate string) (*domain.TaxExemption, error) {
	query := `
		SELECT id, customer_id, exemption_certificate, jurisdiction_id, tax_category,
		       reason, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_tax_exemption
		WHERE exemption_certificate = $1`

	exemption := &domain.TaxExemption{}
	err := r.db.QueryRowContext(ctx, query, certificate).Scan(
		&exemption.ID,
		&exemption.CustomerID,
		&exemption.ExemptionCertificate,
		&exemption.JurisdictionID,
		&exemption.TaxCategory,
		&exemption.Reason,
		&exemption.IsActive,
		&exemption.StartDate,
		&exemption.EndDate,
		&exemption.CreatedAt,
		&exemption.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find tax exemption: %w", err)
	}

	return exemption, nil
}

// FindAll finds all exemptions with optional filters
func (r *PostgresTaxExemptionRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.TaxExemption, error) {
	query := `
		SELECT id, customer_id, exemption_certificate, jurisdiction_id, tax_category,
		       reason, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_tax_exemption`

	if activeOnly {
		query += " WHERE is_active = true"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tax exemptions: %w", err)
	}
	defer rows.Close()

	exemptions := make([]*domain.TaxExemption, 0)
	for rows.Next() {
		exemption := &domain.TaxExemption{}
		err := rows.Scan(
			&exemption.ID,
			&exemption.CustomerID,
			&exemption.ExemptionCertificate,
			&exemption.JurisdictionID,
			&exemption.TaxCategory,
			&exemption.Reason,
			&exemption.IsActive,
			&exemption.StartDate,
			&exemption.EndDate,
			&exemption.CreatedAt,
			&exemption.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tax exemption: %w", err)
		}
		exemptions = append(exemptions, exemption)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tax exemptions: %w", err)
	}

	return exemptions, nil
}

// Delete deletes an exemption
func (r *PostgresTaxExemptionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_tax_exemption WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tax exemption: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrExemptionNotFound
	}

	return nil
}

// ExistsByCertificate checks if an exemption exists with the given certificate
func (r *PostgresTaxExemptionRepository) ExistsByCertificate(ctx context.Context, certificate string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blc_tax_exemption WHERE exemption_certificate = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, certificate).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check exemption certificate existence: %w", err)
	}

	return exists, nil
}
