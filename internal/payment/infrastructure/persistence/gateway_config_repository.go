package persistence

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/qhato/ecommerce/internal/payment/domain"
)

type PostgresGatewayConfigRepository struct {
	db *sql.DB
}

func NewPostgresGatewayConfigRepository(db *sql.DB) *PostgresGatewayConfigRepository {
	return &PostgresGatewayConfigRepository{db: db}
}

func (r *PostgresGatewayConfigRepository) Create(ctx context.Context, config *domain.GatewayConfig) error {
	configJSON, err := json.Marshal(config.Config)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO blc_payment_gateway_config (gateway_name, enabled, priority, environment,
			api_key, api_secret, merchant_id, config)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = r.db.ExecContext(ctx, query,
		config.GatewayName, config.Enabled, config.Priority, config.Environment,
		config.APIKey, config.APISecret, config.MerchantID, configJSON,
	)
	return err
}

func (r *PostgresGatewayConfigRepository) Update(ctx context.Context, config *domain.GatewayConfig) error {
	configJSON, err := json.Marshal(config.Config)
	if err != nil {
		return err
	}

	query := `
		UPDATE blc_payment_gateway_config
		SET enabled = $1, priority = $2, environment = $3, api_key = $4, api_secret = $5,
		    merchant_id = $6, config = $7
		WHERE gateway_name = $8`

	_, err = r.db.ExecContext(ctx, query,
		config.Enabled, config.Priority, config.Environment, config.APIKey, config.APISecret,
		config.MerchantID, configJSON, config.GatewayName,
	)
	return err
}

func (r *PostgresGatewayConfigRepository) FindByName(ctx context.Context, gatewayName string) (*domain.GatewayConfig, error) {
	query := `
		SELECT gateway_name, enabled, priority, environment, api_key, api_secret, merchant_id, config
		FROM blc_payment_gateway_config
		WHERE gateway_name = $1`

	return r.scanGatewayConfig(r.db.QueryRowContext(ctx, query, gatewayName))
}

func (r *PostgresGatewayConfigRepository) FindAllEnabled(ctx context.Context) ([]*domain.GatewayConfig, error) {
	query := `
		SELECT gateway_name, enabled, priority, environment, api_key, api_secret, merchant_id, config
		FROM blc_payment_gateway_config
		WHERE enabled = true
		ORDER BY priority ASC`

	return r.queryGatewayConfigs(ctx, query)
}

func (r *PostgresGatewayConfigRepository) FindAll(ctx context.Context) ([]*domain.GatewayConfig, error) {
	query := `
		SELECT gateway_name, enabled, priority, environment, api_key, api_secret, merchant_id, config
		FROM blc_payment_gateway_config
		ORDER BY priority ASC`

	return r.queryGatewayConfigs(ctx, query)
}

func (r *PostgresGatewayConfigRepository) scanGatewayConfig(row interface {
	Scan(dest ...interface{}) error
}) (*domain.GatewayConfig, error) {
	config := &domain.GatewayConfig{}
	var configJSON []byte

	err := row.Scan(
		&config.GatewayName, &config.Enabled, &config.Priority, &config.Environment,
		&config.APIKey, &config.APISecret, &config.MerchantID, &configJSON,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if len(configJSON) > 0 {
		if err := json.Unmarshal(configJSON, &config.Config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

func (r *PostgresGatewayConfigRepository) queryGatewayConfigs(ctx context.Context, query string, args ...interface{}) ([]*domain.GatewayConfig, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*domain.GatewayConfig
	for rows.Next() {
		config := &domain.GatewayConfig{}
		var configJSON []byte

		if err := rows.Scan(
			&config.GatewayName, &config.Enabled, &config.Priority, &config.Environment,
			&config.APIKey, &config.APISecret, &config.MerchantID, &configJSON,
		); err != nil {
			return nil, err
		}

		if len(configJSON) > 0 {
			if err := json.Unmarshal(configJSON, &config.Config); err != nil {
				return nil, err
			}
		}

		configs = append(configs, config)
	}

	return configs, rows.Err()
}
