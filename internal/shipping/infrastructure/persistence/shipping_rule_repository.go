package persistence

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/qhato/ecommerce/internal/shipping/domain"
)

type PostgresShippingRuleRepository struct {
	db *sql.DB
}

func NewPostgresShippingRuleRepository(db *sql.DB) *PostgresShippingRuleRepository {
	return &PostgresShippingRuleRepository{db: db}
}

func (r *PostgresShippingRuleRepository) Create(ctx context.Context, rule *domain.ShippingRule) error {
	countriesJSON, _ := json.Marshal(rule.Countries)
	excludedZipsJSON, _ := json.Marshal(rule.ExcludedZips)

	query := `INSERT INTO blc_shipping_rule 
		(name, description, rule_type, is_enabled, priority, min_order_value, countries, excluded_zips, discount_type, discount_value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) 
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		rule.Name, rule.Description, rule.RuleType, rule.IsEnabled, rule.Priority,
		rule.MinOrderValue, countriesJSON, excludedZipsJSON,
		rule.DiscountType, rule.DiscountValue, rule.CreatedAt, rule.UpdatedAt,
	).Scan(&rule.ID)
}

func (r *PostgresShippingRuleRepository) Update(ctx context.Context, rule *domain.ShippingRule) error {
	countriesJSON, _ := json.Marshal(rule.Countries)
	excludedZipsJSON, _ := json.Marshal(rule.ExcludedZips)

	query := `UPDATE blc_shipping_rule SET 
		name = $1, description = $2, rule_type = $3, is_enabled = $4, priority = $5,
		min_order_value = $6, countries = $7, excluded_zips = $8,
		discount_type = $9, discount_value = $10, updated_at = $11
		WHERE id = $12`

	_, err := r.db.ExecContext(ctx, query,
		rule.Name, rule.Description, rule.RuleType, rule.IsEnabled, rule.Priority,
		rule.MinOrderValue, countriesJSON, excludedZipsJSON,
		rule.DiscountType, rule.DiscountValue, rule.UpdatedAt, rule.ID,
	)
	return err
}

func (r *PostgresShippingRuleRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM blc_shipping_rule WHERE id = $1`, id)
	return err
}

func (r *PostgresShippingRuleRepository) FindByID(ctx context.Context, id int64) (*domain.ShippingRule, error) {
	rule := &domain.ShippingRule{}
	var countriesJSON, excludedZipsJSON []byte

	query := `SELECT id, name, description, rule_type, is_enabled, priority, min_order_value, countries, excluded_zips, discount_type, discount_value, created_at, updated_at 
		FROM blc_shipping_rule WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID, &rule.Name, &rule.Description, &rule.RuleType, &rule.IsEnabled,
		&rule.Priority, &rule.MinOrderValue, &countriesJSON, &excludedZipsJSON,
		&rule.DiscountType, &rule.DiscountValue, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(countriesJSON, &rule.Countries); err != nil {
		rule.Countries = make([]string, 0)
	}
	if err := json.Unmarshal(excludedZipsJSON, &rule.ExcludedZips); err != nil {
		rule.ExcludedZips = make([]string, 0)
	}

	return rule, nil
}

func (r *PostgresShippingRuleRepository) FindAllEnabled(ctx context.Context) ([]*domain.ShippingRule, error) {
	query := `SELECT id, name, description, rule_type, is_enabled, priority, min_order_value, countries, excluded_zips, discount_type, discount_value, created_at, updated_at 
		FROM blc_shipping_rule WHERE is_enabled = true ORDER BY priority DESC, name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*domain.ShippingRule
	for rows.Next() {
		rule := &domain.ShippingRule{}
		var countriesJSON, excludedZipsJSON []byte

		if err := rows.Scan(
			&rule.ID, &rule.Name, &rule.Description, &rule.RuleType, &rule.IsEnabled,
			&rule.Priority, &rule.MinOrderValue, &countriesJSON, &excludedZipsJSON,
			&rule.DiscountType, &rule.DiscountValue, &rule.CreatedAt, &rule.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(countriesJSON, &rule.Countries); err != nil {
			rule.Countries = make([]string, 0)
		}
		if err := json.Unmarshal(excludedZipsJSON, &rule.ExcludedZips); err != nil {
			rule.ExcludedZips = make([]string, 0)
		}

		rules = append(rules, rule)
	}

	return rules, rows.Err()
}
