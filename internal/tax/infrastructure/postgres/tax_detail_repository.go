package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/tax/domain"
)

// TaxDetailRepository implements domain.TaxDetailRepository for PostgreSQL persistence.
type TaxDetailRepository struct {
	db *sql.DB
}

// NewTaxDetailRepository creates a new PostgreSQL tax detail repository.
func NewTaxDetailRepository(db *sql.DB) *TaxDetailRepository {
	return &TaxDetailRepository{db: db}
}

// Save stores a new tax detail or updates an existing one.
func (r *TaxDetailRepository) Save(ctx context.Context, taxDetail *domain.TaxDetail) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle nullable fields
	amount := sql.NullFloat64{Float64: taxDetail.Amount, Valid: taxDetail.Amount != 0.0}
	taxCountry := sql.NullString{String: taxDetail.TaxCountry, Valid: taxDetail.TaxCountry != ""}
	jurisdictionName := sql.NullString{String: taxDetail.JurisdictionName, Valid: taxDetail.JurisdictionName != ""}
	rate := sql.NullFloat64{Float64: taxDetail.Rate, Valid: taxDetail.Rate != 0.0}
	taxRegion := sql.NullString{String: taxDetail.TaxRegion, Valid: taxDetail.TaxRegion != ""}
	taxName := sql.NullString{String: taxDetail.TaxName, Valid: taxDetail.TaxName != ""}
	taxType := sql.NullString{String: taxDetail.Type, Valid: taxDetail.Type != ""}
	currencyCode := sql.NullString{String: taxDetail.CurrencyCode, Valid: taxDetail.CurrencyCode != ""}
	moduleConfigID := sql.NullInt64{Int64: 0, Valid: false}
	if taxDetail.ModuleConfigID != nil {
		moduleConfigID = sql.NullInt64{Int64: *taxDetail.ModuleConfigID, Valid: true}
	}

	if taxDetail.ID == 0 {
		// Insert new tax detail
		query := `
			INSERT INTO blc_tax_detail (
				amount, tax_country, jurisdiction_name, rate, tax_region, tax_name, 
				type, currency_code, module_config_id, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
			) RETURNING tax_detail_id`
		_, err = tx.ExecContext(ctx, query,
			amount, taxCountry, jurisdictionName, rate, taxRegion, taxName,
			taxType, currencyCode, moduleConfigID, taxDetail.CreatedAt, taxDetail.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert tax detail: %w", err)
		}
	} else {
		// Update existing tax detail
		query := `
			UPDATE blc_tax_detail SET
				amount = $1, tax_country = $2, jurisdiction_name = $3, rate = $4, 
				tax_region = $5, tax_name = $6, type = $7, currency_code = $8, 
				module_config_id = $9, updated_at = $10
			WHERE tax_detail_id = $11`
		_, err = tx.ExecContext(ctx, query,
			amount, taxCountry, jurisdictionName, rate, taxRegion, taxName,
			taxType, currencyCode, moduleConfigID, taxDetail.UpdatedAt, taxDetail.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update tax detail: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a tax detail by its unique identifier.
func (r *TaxDetailRepository) FindByID(ctx context.Context, id int64) (*domain.TaxDetail, error) {
	query := `
		SELECT
			tax_detail_id, amount, tax_country, jurisdiction_name, rate, tax_region, 
			tax_name, type, currency_code, module_config_id, created_at, updated_at
		FROM blc_tax_detail WHERE tax_detail_id = $1`

	var taxDetail domain.TaxDetail
	var amount sql.NullFloat64
	var taxCountry sql.NullString
	var jurisdictionName sql.NullString
	var rate sql.NullFloat64
	var taxRegion sql.NullString
	var taxName sql.NullString
	var taxType sql.NullString
	var currencyCode sql.NullString
	var moduleConfigID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&taxDetail.ID, &amount, &taxCountry, &jurisdictionName, &rate, &taxRegion,
		&taxName, &taxType, &currencyCode, &moduleConfigID, &taxDetail.CreatedAt, &taxDetail.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query tax detail by ID: %w", err)
	}

	if amount.Valid {
		taxDetail.Amount = amount.Float64
	}
	if taxCountry.Valid {
		taxDetail.TaxCountry = taxCountry.String
	}
	if jurisdictionName.Valid {
		taxDetail.JurisdictionName = jurisdictionName.String
	}
	if rate.Valid {
		taxDetail.Rate = rate.Float64
	}
	if taxRegion.Valid {
		taxDetail.TaxRegion = taxRegion.String
	}
	if taxName.Valid {
		taxDetail.TaxName = taxName.String
	}
	if taxType.Valid {
		taxDetail.Type = taxType.String
	}
	if currencyCode.Valid {
		taxDetail.CurrencyCode = currencyCode.String
	}
	if moduleConfigID.Valid {
		taxDetail.ModuleConfigID = &moduleConfigID.Int64
	}

	return &taxDetail, nil
}

// FindApplicableTaxDetails retrieves tax details applicable to a given country, region, and type.
func (r *TaxDetailRepository) FindApplicableTaxDetails(ctx context.Context, taxCountry, taxRegion, taxType string) ([]*domain.TaxDetail, error) {
	query := `
		SELECT
			tax_detail_id, amount, tax_country, jurisdiction_name, rate, tax_region, 
			tax_name, type, currency_code, module_config_id, created_at, updated_at
		FROM blc_tax_detail 
		WHERE tax_country = $1 AND (tax_region IS NULL OR tax_region = $2) AND (type IS NULL OR type = $3)`

	rows, err := r.db.QueryContext(ctx, query, taxCountry, taxRegion, taxType)
	if err != nil {
		return nil, fmt.Errorf("failed to query applicable tax details: %w", err)
	}
	defer rows.Close()

	var taxDetails []*domain.TaxDetail
	for rows.Next() {
		var taxDetail domain.TaxDetail
		var amount sql.NullFloat64
		var tc sql.NullString // Use different name to avoid conflict
		var jn sql.NullString
		var rt sql.NullFloat64
		var tr sql.NullString
		var tn sql.NullString
		var ty sql.NullString
		var cc sql.NullString
		var mcid sql.NullInt64

		err := rows.Scan(
			&taxDetail.ID, &amount, &tc, &jn, &rt, &tr,
			&tn, &ty, &cc, &mcid, &taxDetail.CreatedAt, &taxDetail.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tax detail row: %w", err)
		}

		if amount.Valid {
			taxDetail.Amount = amount.Float64
		}
		if tc.Valid {
			taxDetail.TaxCountry = tc.String
		}
		if jn.Valid {
			taxDetail.JurisdictionName = jn.String
		}
		if rt.Valid {
			taxDetail.Rate = rt.Float64
		}
		if tr.Valid {
			taxDetail.TaxRegion = tr.String
		}
		if tn.Valid {
			taxDetail.TaxName = tn.String
		}
		if ty.Valid {
			taxDetail.Type = ty.String
		}
		if cc.Valid {
			taxDetail.CurrencyCode = cc.String
		}
		if mcid.Valid {
			taxDetail.ModuleConfigID = &mcid.Int64
		}
		taxDetails = append(taxDetails, &taxDetail)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for applicable tax details: %w", err)
	}

	return taxDetails, nil
}

// Delete removes a tax detail by its unique identifier.
func (r *TaxDetailRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_tax_detail WHERE tax_detail_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tax detail: %w", err)
	}
	return nil
}
