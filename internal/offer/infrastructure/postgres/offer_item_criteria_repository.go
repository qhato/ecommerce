package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/offer/domain"
)

// OfferItemCriteriaRepository implements domain.OfferItemCriteriaRepository for PostgreSQL persistence.
type OfferItemCriteriaRepository struct {
	db *sql.DB
}

// NewOfferItemCriteriaRepository creates a new PostgreSQL offer item criteria repository.
func NewOfferItemCriteriaRepository(db *sql.DB) *OfferItemCriteriaRepository {
	return &OfferItemCriteriaRepository{db: db}
}

// Save stores a new offer item criteria or updates an existing one.
func (r *OfferItemCriteriaRepository) Save(ctx context.Context, criteria *domain.OfferItemCriteria) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle nullable fields
	orderItemMatchRule := sql.NullString{String: criteria.OrderItemMatchRule, Valid: criteria.OrderItemMatchRule != ""}

	if criteria.ID == 0 {
		// Insert new offer item criteria
		query := `
			INSERT INTO blc_offer_item_criteria (
				quantity, order_item_match_rule, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4
			) RETURNING offer_item_criteria_id`
		err = tx.QueryRowContext(ctx, query,
			criteria.Quantity, orderItemMatchRule, criteria.CreatedAt, criteria.UpdatedAt,
		).Scan(&criteria.ID)
		if err != nil {
			return fmt.Errorf("failed to insert offer item criteria: %w", err)
		}
	} else {
		// Update existing offer item criteria
		query := `
			UPDATE blc_offer_item_criteria SET
				quantity = $1, order_item_match_rule = $2, updated_at = $3
			WHERE offer_item_criteria_id = $4`
		_, err = tx.ExecContext(ctx, query,
			criteria.Quantity, orderItemMatchRule, criteria.UpdatedAt, criteria.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update offer item criteria: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves an offer item criteria by its unique identifier.
func (r *OfferItemCriteriaRepository) FindByID(ctx context.Context, id int64) (*domain.OfferItemCriteria, error) {
	query := `
		SELECT
			offer_item_criteria_id, quantity, order_item_match_rule, created_at, updated_at
		FROM blc_offer_item_criteria WHERE offer_item_criteria_id = $1`

	var criteria domain.OfferItemCriteria
	var orderItemMatchRule sql.NullString

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&criteria.ID, &criteria.Quantity, &orderItemMatchRule, &criteria.CreatedAt, &criteria.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query offer item criteria by ID: %w", err)
	}

	if orderItemMatchRule.Valid {
		criteria.OrderItemMatchRule = orderItemMatchRule.String
	}

	return &criteria, nil
}

// FindAll retrieves all offer item criteria.
func (r *OfferItemCriteriaRepository) FindAll(ctx context.Context) ([]*domain.OfferItemCriteria, error) {
	query := `
		SELECT
			offer_item_criteria_id, quantity, order_item_match_rule, created_at, updated_at
		FROM blc_offer_item_criteria`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all offer item criteria: %w", err)
	}
	defer rows.Close()

	var criteriaList []*domain.OfferItemCriteria
	for rows.Next() {
		var criteria domain.OfferItemCriteria
		var orderItemMatchRule sql.NullString

		err := rows.Scan(
			&criteria.ID, &criteria.Quantity, &orderItemMatchRule, &criteria.CreatedAt, &criteria.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan offer item criteria row: %w", err)
		}

		if orderItemMatchRule.Valid {
			criteria.OrderItemMatchRule = orderItemMatchRule.String
		}
		criteriaList = append(criteriaList, &criteria)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for offer item criteria: %w", err)
	}

	return criteriaList, nil
}

// Delete removes an offer item criteria by its unique identifier.
func (r *OfferItemCriteriaRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_offer_item_criteria WHERE offer_item_criteria_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer item criteria: %w", err)
	}
	return nil
}
