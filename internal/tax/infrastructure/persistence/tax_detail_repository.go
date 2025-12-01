package persistence

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"

	"github.com/qhato/ecommerce/internal/tax/domain"
	"github.com/qhato/ecommerce/pkg/database"
	"github.com/qhato/ecommerce/pkg/errors"
)

// PostgresTaxDetailRepository implements the TaxDetailRepository interface
type PostgresTaxDetailRepository struct {
	db *database.DB
}

// NewPostgresTaxDetailRepository creates a new PostgresTaxDetailRepository
func NewPostgresTaxDetailRepository(db *database.DB) *PostgresTaxDetailRepository {
	return &PostgresTaxDetailRepository{db: db}
}

// Save stores a new tax detail or updates an existing one.
func (r *PostgresTaxDetailRepository) Save(ctx context.Context, taxDetail *domain.TaxDetail) error {
	if taxDetail.ID == 0 {
		return r.create(ctx, taxDetail)
	}
	return r.update(ctx, taxDetail)
}

func (r *PostgresTaxDetailRepository) create(ctx context.Context, taxDetail *domain.TaxDetail) error {
	query := `
		INSERT INTO blc_tax_detail (
			amount, tax_country, jurisdiction_name, rate, tax_region,
			tax_name, type, currency_code, module_config_id, date_created, date_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING tax_detail_id`

	err := r.db.QueryRow(ctx, query,
		taxDetail.Amount,
		taxDetail.TaxCountry,
		taxDetail.JurisdictionName,
		taxDetail.Rate,
		taxDetail.TaxRegion,
		taxDetail.TaxName,
		taxDetail.Type,
		taxDetail.CurrencyCode,
		taxDetail.ModuleConfigID,
		taxDetail.CreatedAt,
		taxDetail.UpdatedAt,
	).Scan(&taxDetail.ID)

	if err != nil {
		return errors.InternalWrap(err, "failed to create tax detail")
	}
	return nil
}

func (r *PostgresTaxDetailRepository) update(ctx context.Context, taxDetail *domain.TaxDetail) error {
	query := `
		UPDATE blc_tax_detail SET
			amount = $1, tax_country = $2, jurisdiction_name = $3, rate = $4,
			tax_region = $5, tax_name = $6, type = $7, currency_code = $8,
			module_config_id = $9, date_updated = $10
		WHERE tax_detail_id = $11`

	tag, err := r.db.Pool().Exec(ctx, query,
		taxDetail.Amount,
		taxDetail.TaxCountry,
		taxDetail.JurisdictionName,
		taxDetail.Rate,
		taxDetail.TaxRegion,
		taxDetail.TaxName,
		taxDetail.Type,
		taxDetail.CurrencyCode,
		taxDetail.ModuleConfigID,
		taxDetail.UpdatedAt,
		taxDetail.ID,
	)

	if err != nil {
		return errors.InternalWrap(err, "failed to update tax detail")
	}
	if tag.RowsAffected() == 0 {
		return errors.NotFound("tax detail not found")
	}
	return nil
}

// FindByID retrieves a tax detail by its unique identifier.
func (r *PostgresTaxDetailRepository) FindByID(ctx context.Context, id int64) (*domain.TaxDetail, error) {
	query := `
		SELECT
			tax_detail_id, amount, tax_country, jurisdiction_name, rate,
			tax_region, tax_name, type, currency_code, module_config_id,
			date_created, date_updated
		FROM blc_tax_detail
		WHERE tax_detail_id = $1`

	taxDetail := &domain.TaxDetail{}
	var moduleConfigID sql.NullInt64

	err := r.db.QueryRow(ctx, query, id).Scan(
		&taxDetail.ID,
		&taxDetail.Amount,
		&taxDetail.TaxCountry,
		&taxDetail.JurisdictionName,
		&taxDetail.Rate,
		&taxDetail.TaxRegion,
		&taxDetail.TaxName,
		&taxDetail.Type,
		&taxDetail.CurrencyCode,
		&moduleConfigID,
		&taxDetail.CreatedAt,
		&taxDetail.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find tax detail by ID")
	}

	if moduleConfigID.Valid {
		taxDetail.ModuleConfigID = &moduleConfigID.Int64
	}

	return taxDetail, nil
}

// FindApplicableTaxDetails retrieves tax details applicable to a given country, region, and type.
func (r *PostgresTaxDetailRepository) FindApplicableTaxDetails(ctx context.Context, taxCountry, taxRegion, taxType string) ([]*domain.TaxDetail, error) {
	query := `
		SELECT
			tax_detail_id, amount, tax_country, jurisdiction_name, rate,
			tax_region, tax_name, type, currency_code, module_config_id,
			date_created, date_updated
		FROM blc_tax_detail
		WHERE tax_country = $1 AND tax_region = $2 AND type = $3`

	rows, err := r.db.Query(ctx, query, taxCountry, taxRegion, taxType)
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find applicable tax details")
	}
	defer rows.Close()

	var taxDetails []*domain.TaxDetail
	for rows.Next() {
		taxDetail := &domain.TaxDetail{}
		var moduleConfigID sql.NullInt64

		err := rows.Scan(
			&taxDetail.ID,
			&taxDetail.Amount,
			&taxDetail.TaxCountry,
			&taxDetail.JurisdictionName,
			&taxDetail.Rate,
			&taxDetail.TaxRegion,
			&taxDetail.TaxName,
			&taxDetail.Type,
			&taxDetail.CurrencyCode,
			&moduleConfigID,
			&taxDetail.CreatedAt,
			&taxDetail.UpdatedAt,
		)
		if err != nil {
			return nil, errors.InternalWrap(err, "failed to scan tax detail")
		}
		if moduleConfigID.Valid {
			taxDetail.ModuleConfigID = &moduleConfigID.Int64
		}
		taxDetails = append(taxDetails, taxDetail)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.InternalWrap(err, "failed to iterate tax details")
	}

	return taxDetails, nil
}

// Delete removes a tax detail by its unique identifier.
func (r *PostgresTaxDetailRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_tax_detail WHERE tax_detail_id = $1`
	tag, err := r.db.Pool().Exec(ctx, query, id)
	if err != nil {
		return errors.InternalWrap(err, "failed to delete tax detail")
	}
	if tag.RowsAffected() == 0 {
		return errors.NotFound("tax detail not found")
	}
	return nil
}
