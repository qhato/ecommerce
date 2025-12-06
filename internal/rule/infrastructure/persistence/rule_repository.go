package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/qhato/ecommerce/internal/rule/domain"
)

type PostgresRuleRepository struct {
	db *sql.DB
}

func NewPostgresRuleRepository(db *sql.DB) *PostgresRuleRepository {
	return &PostgresRuleRepository{db: db}
}

func (r *PostgresRuleRepository) Create(ctx context.Context, rule *domain.Rule) error {
	conditionsJSON, err := json.Marshal(rule.Conditions)
	if err != nil {
		return fmt.Errorf("failed to marshal conditions: %w", err)
	}

	actionsJSON, err := json.Marshal(rule.Actions)
	if err != nil {
		return fmt.Errorf("failed to marshal actions: %w", err)
	}

	contextJSON, err := json.Marshal(rule.Context)
	if err != nil {
		return fmt.Errorf("failed to marshal context: %w", err)
	}

	query := `INSERT INTO blc_rule (
		name, description, type, priority, is_active, conditions, actions,
		start_date, end_date, context, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		rule.Name, rule.Description, rule.Type, rule.Priority, rule.IsActive,
		conditionsJSON, actionsJSON, rule.StartDate, rule.EndDate,
		contextJSON, rule.CreatedAt, rule.UpdatedAt,
	).Scan(&rule.ID)
}

func (r *PostgresRuleRepository) Update(ctx context.Context, rule *domain.Rule) error {
	conditionsJSON, err := json.Marshal(rule.Conditions)
	if err != nil {
		return fmt.Errorf("failed to marshal conditions: %w", err)
	}

	actionsJSON, err := json.Marshal(rule.Actions)
	if err != nil {
		return fmt.Errorf("failed to marshal actions: %w", err)
	}

	contextJSON, err := json.Marshal(rule.Context)
	if err != nil {
		return fmt.Errorf("failed to marshal context: %w", err)
	}

	query := `UPDATE blc_rule SET
		name = $1, description = $2, priority = $3, is_active = $4,
		conditions = $5, actions = $6, start_date = $7, end_date = $8,
		context = $9, updated_at = $10
	WHERE id = $11`

	_, err = r.db.ExecContext(ctx, query,
		rule.Name, rule.Description, rule.Priority, rule.IsActive,
		conditionsJSON, actionsJSON, rule.StartDate, rule.EndDate,
		contextJSON, rule.UpdatedAt, rule.ID,
	)
	return err
}

func (r *PostgresRuleRepository) FindByID(ctx context.Context, id int64) (*domain.Rule, error) {
	query := `SELECT id, name, description, type, priority, is_active, conditions,
		actions, start_date, end_date, context, created_at, updated_at
	FROM blc_rule WHERE id = $1`

	return r.scanRule(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresRuleRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.Rule, error) {
	query := `SELECT id, name, description, type, priority, is_active, conditions,
		actions, start_date, end_date, context, created_at, updated_at
	FROM blc_rule`

	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY priority DESC, created_at DESC"

	return r.queryRules(ctx, query)
}

func (r *PostgresRuleRepository) FindByType(ctx context.Context, ruleType domain.RuleType, activeOnly bool) ([]*domain.Rule, error) {
	query := `SELECT id, name, description, type, priority, is_active, conditions,
		actions, start_date, end_date, context, created_at, updated_at
	FROM blc_rule WHERE type = $1`

	if activeOnly {
		query += " AND is_active = true"
	}
	query += " ORDER BY priority DESC"

	return r.queryRules(ctx, query, ruleType)
}

func (r *PostgresRuleRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM blc_rule WHERE id = $1`, id)
	return err
}

func (r *PostgresRuleRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM blc_rule WHERE name = $1)`, name).Scan(&exists)
	return exists, err
}

// Helper methods

func (r *PostgresRuleRepository) scanRule(row interface {
	Scan(dest ...interface{}) error
}) (*domain.Rule, error) {
	rule := &domain.Rule{}
	var conditionsJSON, actionsJSON, contextJSON []byte

	err := row.Scan(
		&rule.ID, &rule.Name, &rule.Description, &rule.Type, &rule.Priority,
		&rule.IsActive, &conditionsJSON, &actionsJSON, &rule.StartDate,
		&rule.EndDate, &contextJSON, &rule.CreatedAt, &rule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(conditionsJSON, &rule.Conditions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conditions: %w", err)
	}

	if err := json.Unmarshal(actionsJSON, &rule.Actions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal actions: %w", err)
	}

	if contextJSON != nil {
		if err := json.Unmarshal(contextJSON, &rule.Context); err != nil {
			return nil, fmt.Errorf("failed to unmarshal context: %w", err)
		}
	}

	return rule, nil
}

func (r *PostgresRuleRepository) queryRules(ctx context.Context, query string, args ...interface{}) ([]*domain.Rule, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := make([]*domain.Rule, 0)
	for rows.Next() {
		rule := &domain.Rule{}
		var conditionsJSON, actionsJSON, contextJSON []byte

		if err := rows.Scan(
			&rule.ID, &rule.Name, &rule.Description, &rule.Type, &rule.Priority,
			&rule.IsActive, &conditionsJSON, &actionsJSON, &rule.StartDate,
			&rule.EndDate, &contextJSON, &rule.CreatedAt, &rule.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(conditionsJSON, &rule.Conditions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal conditions: %w", err)
		}

		if err := json.Unmarshal(actionsJSON, &rule.Actions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal actions: %w", err)
		}

		if contextJSON != nil {
			if err := json.Unmarshal(contextJSON, &rule.Context); err != nil {
				return nil, fmt.Errorf("failed to unmarshal context: %w", err)
			}
		}

		rules = append(rules, rule)
	}

	return rules, nil
}
