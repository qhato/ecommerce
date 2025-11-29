package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/offer/domain"
)

// OfferRuleRepository implements domain.OfferRuleRepository for PostgreSQL persistence.
type OfferRuleRepository struct {
	db *sql.DB
}

// NewOfferRuleRepository creates a new PostgreSQL offer rule repository.
func NewOfferRuleRepository(db *sql.DB) *OfferRuleRepository {
	return &OfferRuleRepository{db: db}
}

// Save stores a new offer rule or updates an existing one.
func (r *OfferRuleRepository) Save(ctx context.Context, rule *domain.OfferRule) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle nullable fields
	matchRule := sql.NullString{String: rule.MatchRule, Valid: rule.MatchRule != ""}

	if rule.ID == 0 {
		// Insert new offer rule
		query := `
			INSERT INTO blc_offer_rule (
				match_rule, created_at, updated_at
			) VALUES (
				$1, $2, $3
			) RETURNING offer_rule_id`
		err = tx.QueryRowContext(ctx, query,
			matchRule, rule.CreatedAt, rule.UpdatedAt,
		).Scan(&rule.ID)
		if err != nil {
			return fmt.Errorf("failed to insert offer rule: %w", err)
		}
	} else {
		// Update existing offer rule
		query := `
			UPDATE blc_offer_rule SET
				match_rule = $1, updated_at = $2
			WHERE offer_rule_id = $3`
		_, err = tx.ExecContext(ctx, query,
			matchRule, rule.UpdatedAt, rule.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update offer rule: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves an offer rule by its unique identifier.
func (r *OfferRuleRepository) FindByID(ctx context.Context, id int64) (*domain.OfferRule, error) {
	query := `
		SELECT
			offer_rule_id, match_rule, created_at, updated_at
		FROM blc_offer_rule WHERE offer_rule_id = $1`

	var rule domain.OfferRule
	var matchRule sql.NullString

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&rule.ID, &matchRule, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query offer rule by ID: %w", err)
	}

	if matchRule.Valid {
		rule.MatchRule = matchRule.String
	}

	return &rule, nil
}

// FindAll retrieves all offer rules.
func (r *OfferRuleRepository) FindAll(ctx context.Context) ([]*domain.OfferRule, error) {
	query := `
		SELECT
			offer_rule_id, match_rule, created_at, updated_at
		FROM blc_offer_rule`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all offer rules: %w", err)
	}
	defer rows.Close()

	var rules []*domain.OfferRule
	for rows.Next() {
		var rule domain.OfferRule
		var matchRule sql.NullString

		err := rows.Scan(
			&rule.ID, &matchRule, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan offer rule row: %w", err)
		}

		if matchRule.Valid {
			rule.MatchRule = matchRule.String
		}
		rules = append(rules, &rule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for offer rules: %w", err)
	}

	return rules, nil
}

// Delete removes an offer rule by its unique identifier.
func (r *OfferRuleRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_offer_rule WHERE offer_rule_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer rule: %w", err)
	}
	return nil
}
