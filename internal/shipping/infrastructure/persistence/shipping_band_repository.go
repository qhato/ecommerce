package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/shipping/domain"
)

type PostgresShippingBandRepository struct {
	db *sql.DB
}

func NewPostgresShippingBandRepository(db *sql.DB) *PostgresShippingBandRepository {
	return &PostgresShippingBandRepository{db: db}
}

func (r *PostgresShippingBandRepository) Create(ctx context.Context, band *domain.ShippingBand) error {
	query := `INSERT INTO blc_shipping_band 
		(method_id, band_type, min_value, max_value, price, percent_charge, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		band.MethodID, band.BandType, band.MinValue, band.MaxValue,
		band.Price, band.PercentCharge, band.CreatedAt,
	).Scan(&band.ID)
}

func (r *PostgresShippingBandRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM blc_shipping_band WHERE id = $1`, id)
	return err
}

func (r *PostgresShippingBandRepository) FindByMethodID(ctx context.Context, methodID int64) ([]*domain.ShippingBand, error) {
	query := `SELECT id, method_id, band_type, min_value, max_value, price, percent_charge, created_at 
		FROM blc_shipping_band WHERE method_id = $1 ORDER BY min_value`

	rows, err := r.db.QueryContext(ctx, query, methodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bands []*domain.ShippingBand
	for rows.Next() {
		band := &domain.ShippingBand{}
		if err := rows.Scan(
			&band.ID, &band.MethodID, &band.BandType, &band.MinValue,
			&band.MaxValue, &band.Price, &band.PercentCharge, &band.CreatedAt,
		); err != nil {
			return nil, err
		}
		bands = append(bands, band)
	}

	return bands, rows.Err()
}

func (r *PostgresShippingBandRepository) DeleteByMethodID(ctx context.Context, methodID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM blc_shipping_band WHERE method_id = $1`, methodID)
	return err
}
