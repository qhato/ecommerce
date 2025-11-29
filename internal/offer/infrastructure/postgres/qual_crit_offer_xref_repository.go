package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/offer/domain"
)

// QualCritOfferXrefRepository implements domain.QualCritOfferXrefRepository for PostgreSQL persistence.
type QualCritOfferXrefRepository struct {
	db *sql.DB
}

// NewQualCritOfferXrefRepository creates a new PostgreSQL qualifying criteria xref repository.
func NewQualCritOfferXrefRepository(db *sql.DB) *QualCritOfferXrefRepository {
	return &QualCritOfferXrefRepository{db: db}
}

// Save stores a new qualifying criteria xref or updates an existing one.
func (r *QualCritOfferXrefRepository) Save(ctx context.Context, xref *domain.QualCritOfferXref) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	if xref.ID == 0 {
		// Insert new qualifying criteria xref
		query := `
			INSERT INTO blc_qual_crit_offer_xref (
				offer_id, offer_item_criteria_id, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4
			) RETURNING offer_qual_crit_id`
		err = tx.QueryRowContext(ctx, query,
			xref.OfferID, xref.OfferItemCriteriaID, xref.CreatedAt, xref.UpdatedAt,
		).Scan(&xref.ID)
		if err != nil {
			return fmt.Errorf("failed to insert qualifying criteria xref: %w", err)
		}
	} else {
		// Update existing qualifying criteria xref
		query := `
			UPDATE blc_qual_crit_offer_xref SET
				offer_id = $1, offer_item_criteria_id = $2, updated_at = $3
			WHERE offer_qual_crit_id = $4`
		_, err = tx.ExecContext(ctx, query,
			xref.OfferID, xref.OfferItemCriteriaID, xref.UpdatedAt, xref.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update qualifying criteria xref: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a qualifying criteria xref by its unique identifier.
func (r *QualCritOfferXrefRepository) FindByID(ctx context.Context, id int64) (*domain.QualCritOfferXref, error) {
	query := `
		SELECT
			offer_qual_crit_id, offer_id, offer_item_criteria_id, created_at, updated_at
		FROM blc_qual_crit_offer_xref WHERE offer_qual_crit_id = $1`

	var xref domain.QualCritOfferXref
	var offerItemCriteriaID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&xref.ID, &xref.OfferID, &offerItemCriteriaID, &xref.CreatedAt, &xref.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query qualifying criteria xref by ID: %w", err)
	}

	if offerItemCriteriaID.Valid {
		xref.OfferItemCriteriaID = offerItemCriteriaID.Int64
	}

	return &xref, nil
}

// FindByOfferID retrieves all qualifying criteria xrefs for a given offer ID.
func (r *QualCritOfferXrefRepository) FindByOfferID(ctx context.Context, offerID int64) ([]*domain.QualCritOfferXref, error) {
	query := `
		SELECT
			offer_qual_crit_id, offer_id, offer_item_criteria_id, created_at, updated_at
		FROM blc_qual_crit_offer_xref WHERE offer_id = $1`

	rows, err := r.db.QueryContext(ctx, query, offerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query qualifying criteria xrefs by offer ID: %w", err)
	}
	defer rows.Close()

	var xrefs []*domain.QualCritOfferXref
	for rows.Next() {
		var xref domain.QualCritOfferXref
		var offerItemCriteriaID sql.NullInt64

		err := rows.Scan(
			&xref.ID, &xref.OfferID, &offerItemCriteriaID, &xref.CreatedAt, &xref.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan qualifying criteria xref row: %w", err)
		}

		if offerItemCriteriaID.Valid {
			xref.OfferItemCriteriaID = offerItemCriteriaID.Int64
		}
		xrefs = append(xrefs, &xref)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for qualifying criteria xrefs: %w", err)
	}

	return xrefs, nil
}

// FindByOfferItemCriteriaID retrieves all qualifying criteria xrefs for a given offer item criteria ID.
func (r *QualCritOfferXrefRepository) FindByOfferItemCriteriaID(ctx context.Context, offerItemCriteriaID int64) ([]*domain.QualCritOfferXref, error) {
	query := `
		SELECT
			offer_qual_crit_id, offer_id, offer_item_criteria_id, created_at, updated_at
		FROM blc_qual_crit_offer_xref WHERE offer_item_criteria_id = $1`

	rows, err := r.db.QueryContext(ctx, query, offerItemCriteriaID)
	if err != nil {
		return nil, fmt.Errorf("failed to query qualifying criteria xrefs by offer item criteria ID: %w", err)
	}
	defer rows.Close()

	var xrefs []*domain.QualCritOfferXref
	for rows.Next() {
		var xref domain.QualCritOfferXref
		var offID sql.NullInt64 // Use different name to avoid conflict

		err := rows.Scan(
			&xref.ID, &offID, &xref.OfferItemCriteriaID, &xref.CreatedAt, &xref.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan qualifying criteria xref row: %w", err)
		}

		if offID.Valid {
			xref.OfferID = offID.Int64
		}
		xrefs = append(xrefs, &xref)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for qualifying criteria xrefs: %w", err)
	}

	return xrefs, nil
}


// Delete removes a qualifying criteria xref by its unique identifier.
func (r *QualCritOfferXrefRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_qual_crit_offer_xref WHERE offer_qual_crit_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete qualifying criteria xref: %w", err)
	}
	return nil
}

// DeleteByOfferID removes all qualifying criteria xrefs for a given offer ID.
func (r *QualCritOfferXrefRepository) DeleteByOfferID(ctx context.Context, offerID int64) error {
	query := `DELETE FROM blc_qual_crit_offer_xref WHERE offer_id = $1`
	_, err := r.db.ExecContext(ctx, query, offerID)
	if err != nil {
		return fmt.Errorf("failed to delete qualifying criteria xrefs by offer ID: %w", err)
	}
	return nil
}

// DeleteByOfferItemCriteriaID removes all qualifying criteria xrefs for a given offer item criteria ID.
func (r *QualCritOfferXrefRepository) DeleteByOfferItemCriteriaID(ctx context.Context, offerItemCriteriaID int64) error {
	query := `DELETE FROM blc_qual_crit_offer_xref WHERE offer_item_criteria_id = $1`
	_, err := r.db.ExecContext(ctx, query, offerItemCriteriaID)
	if err != nil {
		return fmt.Errorf("failed to delete qualifying criteria xrefs by offer item criteria ID: %w", err)
	}
	return nil
}

// RemoveQualCritOfferXref removes a specific qualifying criteria xref by offer ID and offer item criteria ID.
func (r *QualCritOfferXrefRepository) RemoveQualCritOfferXref(ctx context.Context, offerID, offerItemCriteriaID int64) error {
	query := `DELETE FROM blc_qual_crit_offer_xref WHERE offer_id = $1 AND offer_item_criteria_id = $2`
	_, err := r.db.ExecContext(ctx, query, offerID, offerItemCriteriaID)
	if err != nil {
		return fmt.Errorf("failed to remove qualifying criteria xref: %w", err)
	}
	return nil
}
