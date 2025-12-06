package persistence

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/qhato/ecommerce/internal/shipping/domain"
)

type PostgresCarrierConfigRepository struct {
	db *sql.DB
}

func NewPostgresCarrierConfigRepository(db *sql.DB) *PostgresCarrierConfigRepository {
	return &PostgresCarrierConfigRepository{db: db}
}

func (r *PostgresCarrierConfigRepository) Create(ctx context.Context, config *domain.CarrierConfig) error {
	configJSON, err := json.Marshal(config.Config)
	if err != nil {
		return err
	}

	query := `INSERT INTO blc_carrier_config 
		(carrier, name, is_enabled, priority, api_key, api_secret, account_id, config, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		config.Carrier, config.Name, config.IsEnabled, config.Priority,
		config.APIKey, config.APISecret, config.AccountID, configJSON,
		config.CreatedAt, config.UpdatedAt,
	).Scan(&config.ID)
}

func (r *PostgresCarrierConfigRepository) Update(ctx context.Context, config *domain.CarrierConfig) error {
	configJSON, err := json.Marshal(config.Config)
	if err != nil {
		return err
	}

	query := `UPDATE blc_carrier_config SET 
		name = $1, is_enabled = $2, priority = $3, api_key = $4, 
		api_secret = $5, account_id = $6, config = $7, updated_at = $8
		WHERE id = $9`

	_, err = r.db.ExecContext(ctx, query,
		config.Name, config.IsEnabled, config.Priority, config.APIKey,
		config.APISecret, config.AccountID, configJSON, config.UpdatedAt, config.ID,
	)
	return err
}

func (r *PostgresCarrierConfigRepository) FindByID(ctx context.Context, id int64) (*domain.CarrierConfig, error) {
	config := &domain.CarrierConfig{}
	var configJSON []byte

	query := `SELECT id, carrier, name, is_enabled, priority, api_key, api_secret, account_id, config, created_at, updated_at 
		FROM blc_carrier_config WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&config.ID, &config.Carrier, &config.Name, &config.IsEnabled, &config.Priority,
		&config.APIKey, &config.APISecret, &config.AccountID, &configJSON,
		&config.CreatedAt, &config.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(configJSON, &config.Config); err != nil {
		config.Config = make(map[string]string)
	}

	return config, nil
}

func (r *PostgresCarrierConfigRepository) FindByCarrier(ctx context.Context, carrier domain.ShippingCarrier) (*domain.CarrierConfig, error) {
	config := &domain.CarrierConfig{}
	var configJSON []byte

	query := `SELECT id, carrier, name, is_enabled, priority, api_key, api_secret, account_id, config, created_at, updated_at 
		FROM blc_carrier_config WHERE carrier = $1`

	err := r.db.QueryRowContext(ctx, query, carrier).Scan(
		&config.ID, &config.Carrier, &config.Name, &config.IsEnabled, &config.Priority,
		&config.APIKey, &config.APISecret, &config.AccountID, &configJSON,
		&config.CreatedAt, &config.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(configJSON, &config.Config); err != nil {
		config.Config = make(map[string]string)
	}

	return config, nil
}

func (r *PostgresCarrierConfigRepository) FindAll(ctx context.Context, enabledOnly bool) ([]*domain.CarrierConfig, error) {
	query := `SELECT id, carrier, name, is_enabled, priority, api_key, api_secret, account_id, config, created_at, updated_at 
		FROM blc_carrier_config`
	
	if enabledOnly {
		query += " WHERE is_enabled = true"
	}
	query += " ORDER BY priority DESC, name ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*domain.CarrierConfig
	for rows.Next() {
		config := &domain.CarrierConfig{}
		var configJSON []byte

		if err := rows.Scan(
			&config.ID, &config.Carrier, &config.Name, &config.IsEnabled, &config.Priority,
			&config.APIKey, &config.APISecret, &config.AccountID, &configJSON,
			&config.CreatedAt, &config.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(configJSON, &config.Config); err != nil {
			config.Config = make(map[string]string)
		}

		configs = append(configs, config)
	}

	return configs, rows.Err()
}
