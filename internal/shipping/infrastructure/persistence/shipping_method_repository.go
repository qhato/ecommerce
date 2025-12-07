package persistence

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/qhato/ecommerce/internal/shipping/domain"
	"github.com/qhato/ecommerce/pkg/database"
)

type PostgresShippingMethodRepository struct {
	db *database.DB
}

func NewPostgresShippingMethodRepository(db *database.DB) *PostgresShippingMethodRepository {
	return &PostgresShippingMethodRepository{db: db}
}

func (r *PostgresShippingMethodRepository) Create(ctx context.Context, method *domain.ShippingMethod) error {
	query := `INSERT INTO blc_shipping_method 
		(carrier, name, description, service_code, estimated_days, pricing_type, flat_rate, is_enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING id`

	return r.db.QueryRow(ctx, query,
		method.Carrier, method.Name, method.Description, method.ServiceCode,
		method.EstimatedDays, method.PricingType, method.FlatRate,
		method.IsEnabled, method.CreatedAt, method.UpdatedAt,
	).Scan(&method.ID)
}

func (r *PostgresShippingMethodRepository) Update(ctx context.Context, method *domain.ShippingMethod) error {
	query := `UPDATE blc_shipping_method SET 
		name = $1, description = $2, service_code = $3, estimated_days = $4,
		pricing_type = $5, flat_rate = $6, is_enabled = $7, updated_at = $8
		WHERE id = $9`

	return r.db.Exec(ctx, query,
		method.Name, method.Description, method.ServiceCode, method.EstimatedDays,
		method.PricingType, method.FlatRate, method.IsEnabled, method.UpdatedAt, method.ID,
	)
}

func (r *PostgresShippingMethodRepository) Delete(ctx context.Context, id int64) error {
	return r.db.Exec(ctx, `DELETE FROM blc_shipping_method WHERE id = $1`, id)
}

func (r *PostgresShippingMethodRepository) FindByID(ctx context.Context, id int64) (*domain.ShippingMethod, error) {
	method := &domain.ShippingMethod{}

	query := `SELECT id, carrier, name, description, service_code, estimated_days, pricing_type, flat_rate, is_enabled, created_at, updated_at 
		FROM blc_shipping_method WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&method.ID, &method.Carrier, &method.Name, &method.Description,
		&method.ServiceCode, &method.EstimatedDays, &method.PricingType,
		&method.FlatRate, &method.IsEnabled, &method.CreatedAt, &method.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Load bands
	bands, err := r.loadBands(ctx, method.ID)
	if err != nil {
		return nil, err
	}
	method.Bands = bands

	return method, nil
}

func (r *PostgresShippingMethodRepository) FindByCarrier(ctx context.Context, carrier domain.ShippingCarrier) ([]*domain.ShippingMethod, error) {
	query := `SELECT id, carrier, name, description, service_code, estimated_days, pricing_type, flat_rate, is_enabled, created_at, updated_at 
		FROM blc_shipping_method WHERE carrier = $1 ORDER BY name`

	rows, err := r.db.Query(ctx, query, carrier)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMethods(ctx, rows)
}

func (r *PostgresShippingMethodRepository) FindAllEnabled(ctx context.Context) ([]*domain.ShippingMethod, error) {
	query := `SELECT id, carrier, name, description, service_code, estimated_days, pricing_type, flat_rate, is_enabled, created_at, updated_at
		FROM blc_shipping_method WHERE is_enabled = true ORDER BY carrier, name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMethods(ctx, rows)
}

func (r *PostgresShippingMethodRepository) scanMethods(ctx context.Context, rows pgx.Rows) ([]*domain.ShippingMethod, error) {
	var methods []*domain.ShippingMethod

	for rows.Next() {
		method := &domain.ShippingMethod{}
		if err := rows.Scan(
			&method.ID, &method.Carrier, &method.Name, &method.Description,
			&method.ServiceCode, &method.EstimatedDays, &method.PricingType,
			&method.FlatRate, &method.IsEnabled, &method.CreatedAt, &method.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Load bands for each method
		bands, err := r.loadBands(ctx, method.ID)
		if err != nil {
			return nil, err
		}
		method.Bands = bands

		methods = append(methods, method)
	}

	return methods, rows.Err()
}

func (r *PostgresShippingMethodRepository) loadBands(ctx context.Context, methodID int64) ([]domain.ShippingBand, error) {
	query := `SELECT id, method_id, band_type, min_value, max_value, price, percent_charge, created_at 
		FROM blc_shipping_band WHERE method_id = $1 ORDER BY min_value`

	rows, err := r.db.Query(ctx, query, methodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bands []domain.ShippingBand
	for rows.Next() {
		band := domain.ShippingBand{}
		if err := rows.Scan(
			&band.ID, &band.MethodID, &band.BandType, &band.MinValue,
			&band.MaxValue, &band.Price, &band.PercentCharge, &band.CreatedAt,
		); err != nil {
			return nil, err
		}
		bands = append(bands, band)
	}

	if bands == nil {
		bands = make([]domain.ShippingBand, 0)
	}

	return bands, rows.Err()
}
