package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/qhato/ecommerce/internal/pricing/domain"
	"github.com/shopspring/decimal"
)

// PostgresPricingRuleRepository implements PricingRuleRepository using PostgreSQL
type PostgresPricingRuleRepository struct {
	db *sql.DB
}

// NewPostgresPricingRuleRepository creates a new PostgresPricingRuleRepository
func NewPostgresPricingRuleRepository(db *sql.DB) domain.PricingRuleRepository {
	return &PostgresPricingRuleRepository{db: db}
}

func (r *PostgresPricingRuleRepository) Save(ctx context.Context, rule *domain.PricingRule) error {
	if rule.ID == 0 {
		return r.insert(ctx, rule)
	}
	return r.update(ctx, rule)
}

func (r *PostgresPricingRuleRepository) insert(ctx context.Context, rule *domain.PricingRule) error {
	query := `
		INSERT INTO blc_pricing_rule (
			name, description, rule_type, priority, is_active,
			start_date, end_date, condition_expression, action_type, action_value,
			applicable_skus, applicable_categories, customer_segments,
			min_quantity, max_quantity, min_order_value,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		rule.Name,
		rule.Description,
		rule.RuleType,
		rule.Priority,
		rule.IsActive,
		rule.StartDate,
		rule.EndDate,
		rule.ConditionExpression,
		rule.ActionType,
		rule.ActionValue,
		pq.Array(rule.ApplicableSKUs),
		pq.Array(rule.ApplicableCategories),
		pq.Array(rule.CustomerSegments),
		rule.MinQuantity,
		rule.MaxQuantity,
		rule.MinOrderValue,
		rule.CreatedAt,
		rule.UpdatedAt,
	).Scan(&rule.ID)

	if err != nil {
		return fmt.Errorf("failed to insert pricing rule: %w", err)
	}
	return nil
}

func (r *PostgresPricingRuleRepository) update(ctx context.Context, rule *domain.PricingRule) error {
	query := `
		UPDATE blc_pricing_rule
		SET name = $2, description = $3, priority = $4, is_active = $5,
		    start_date = $6, end_date = $7, condition_expression = $8,
		    action_type = $9, action_value = $10, applicable_skus = $11,
		    applicable_categories = $12, customer_segments = $13,
		    min_quantity = $14, max_quantity = $15, min_order_value = $16,
		    updated_at = $17
		WHERE id = $1
	`

	rule.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(
		ctx,
		query,
		rule.ID,
		rule.Name,
		rule.Description,
		rule.Priority,
		rule.IsActive,
		rule.StartDate,
		rule.EndDate,
		rule.ConditionExpression,
		rule.ActionType,
		rule.ActionValue,
		pq.Array(rule.ApplicableSKUs),
		pq.Array(rule.ApplicableCategories),
		pq.Array(rule.CustomerSegments),
		rule.MinQuantity,
		rule.MaxQuantity,
		rule.MinOrderValue,
		rule.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update pricing rule: %w", err)
	}
	return nil
}

func (r *PostgresPricingRuleRepository) FindByID(ctx context.Context, id int64) (*domain.PricingRule, error) {
	query := `
		SELECT id, name, description, rule_type, priority, is_active,
		       start_date, end_date, condition_expression, action_type, action_value,
		       applicable_skus, applicable_categories, customer_segments,
		       min_quantity, max_quantity, min_order_value,
		       created_at, updated_at
		FROM blc_pricing_rule
		WHERE id = $1
	`

	return r.scanPricingRule(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresPricingRuleRepository) FindActive(ctx context.Context) ([]*domain.PricingRule, error) {
	query := `
		SELECT id, name, description, rule_type, priority, is_active,
		       start_date, end_date, condition_expression, action_type, action_value,
		       applicable_skus, applicable_categories, customer_segments,
		       min_quantity, max_quantity, min_order_value,
		       created_at, updated_at
		FROM blc_pricing_rule
		WHERE is_active = true
		  AND (start_date IS NULL OR start_date <= NOW())
		  AND (end_date IS NULL OR end_date >= NOW())
		ORDER BY priority DESC
	`

	return r.queryPricingRules(ctx, query)
}

func (r *PostgresPricingRuleRepository) FindBySKU(ctx context.Context, skuID string) ([]*domain.PricingRule, error) {
	query := `
		SELECT id, name, description, rule_type, priority, is_active,
		       start_date, end_date, condition_expression, action_type, action_value,
		       applicable_skus, applicable_categories, customer_segments,
		       min_quantity, max_quantity, min_order_value,
		       created_at, updated_at
		FROM blc_pricing_rule
		WHERE is_active = true
		  AND (start_date IS NULL OR start_date <= NOW())
		  AND (end_date IS NULL OR end_date >= NOW())
		  AND ($1 = ANY(applicable_skus) OR applicable_skus = '{}')
		ORDER BY priority DESC
	`

	return r.queryPricingRules(ctx, query, skuID)
}

func (r *PostgresPricingRuleRepository) FindByCustomerSegment(ctx context.Context, segment string) ([]*domain.PricingRule, error) {
	query := `
		SELECT id, name, description, rule_type, priority, is_active,
		       start_date, end_date, condition_expression, action_type, action_value,
		       applicable_skus, applicable_categories, customer_segments,
		       min_quantity, max_quantity, min_order_value,
		       created_at, updated_at
		FROM blc_pricing_rule
		WHERE is_active = true
		  AND (start_date IS NULL OR start_date <= NOW())
		  AND (end_date IS NULL OR end_date >= NOW())
		  AND ($1 = ANY(customer_segments) OR customer_segments = '{}')
		ORDER BY priority DESC
	`

	return r.queryPricingRules(ctx, query, segment)
}

func (r *PostgresPricingRuleRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_pricing_rule WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete pricing rule: %w", err)
	}
	return nil
}

func (r *PostgresPricingRuleRepository) scanPricingRule(row *sql.Row) (*domain.PricingRule, error) {
	rule := &domain.PricingRule{}
	var applicableSKUs, applicableCategories, customerSegments pq.StringArray
	var actionValue, minOrderValue sql.NullString

	err := row.Scan(
		&rule.ID,
		&rule.Name,
		&rule.Description,
		&rule.RuleType,
		&rule.Priority,
		&rule.IsActive,
		&rule.StartDate,
		&rule.EndDate,
		&rule.ConditionExpression,
		&rule.ActionType,
		&actionValue,
		&applicableSKUs,
		&applicableCategories,
		&customerSegments,
		&rule.MinQuantity,
		&rule.MaxQuantity,
		&minOrderValue,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan pricing rule: %w", err)
	}

	rule.ApplicableSKUs = applicableSKUs
	rule.ApplicableCategories = applicableCategories
	rule.CustomerSegments = customerSegments

	if actionValue.Valid {
		rule.ActionValue, _ = decimal.NewFromString(actionValue.String)
	}
	if minOrderValue.Valid {
		mov, _ := decimal.NewFromString(minOrderValue.String)
		rule.MinOrderValue = &mov
	}

	return rule, nil
}

func (r *PostgresPricingRuleRepository) queryPricingRules(ctx context.Context, query string, args ...interface{}) ([]*domain.PricingRule, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query pricing rules: %w", err)
	}
	defer rows.Close()

	rules := make([]*domain.PricingRule, 0)
	for rows.Next() {
		rule := &domain.PricingRule{}
		var applicableSKUs, applicableCategories, customerSegments pq.StringArray
		var actionValue, minOrderValue sql.NullString

		err := rows.Scan(
			&rule.ID,
			&rule.Name,
			&rule.Description,
			&rule.RuleType,
			&rule.Priority,
			&rule.IsActive,
			&rule.StartDate,
			&rule.EndDate,
			&rule.ConditionExpression,
			&rule.ActionType,
			&actionValue,
			&applicableSKUs,
			&applicableCategories,
			&customerSegments,
			&rule.MinQuantity,
			&rule.MaxQuantity,
			&minOrderValue,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pricing rule: %w", err)
		}

		rule.ApplicableSKUs = applicableSKUs
		rule.ApplicableCategories = applicableCategories
		rule.CustomerSegments = customerSegments

		if actionValue.Valid {
			rule.ActionValue, _ = decimal.NewFromString(actionValue.String)
		}
		if minOrderValue.Valid {
			mov, _ := decimal.NewFromString(minOrderValue.String)
			rule.MinOrderValue = &mov
		}

		rules = append(rules, rule)
	}

	return rules, nil
}
