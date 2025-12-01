package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/offer/domain"
)

// TarCritOfferXrefRepository implements domain.TarCritOfferXrefRepository for PostgreSQL persistence.
type TarCritOfferXrefRepository struct {
	db *sql.DB
}

// NewTarCritOfferXrefRepository creates a new PostgreSQL target criteria xref repository.
func NewTarCritOfferXrefRepository(db *sql.DB) *TarCritOfferXrefRepository {
	return &TarCritOfferXrefRepository{db: db}
}

// Save stores a new target criteria xref or updates an existing one.
func (r *TarCritOfferXrefRepository) Save(ctx context.Context, xref *domain.TarCritOfferXref) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	if xref.ID == 0 {
		// Insert new target criteria xref
		query := `
			INSERT INTO blc_tar_crit_offer_xref (
				offer_id, offer_item_criteria_id, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4
			) RETURNING offer_tar_crit_id`
		err = tx.QueryRowContext(ctx, query,
			xref.OfferID, xref.OfferItemCriteriaID, xref.CreatedAt, xref.UpdatedAt,
		).Scan(&xref.ID)
		if err != nil {
			return fmt.Errorf("failed to insert target criteria xref: %w", err)
		}
	} else {
		// Update existing target criteria xref
		query := `
			UPDATE blc_tar_crit_offer_xref SET
				offer_id = $1, offer_item_criteria_id = $2, updated_at = $3
			WHERE offer_tar_crit_id = $4`
		_, err = tx.ExecContext(ctx, query,
			xref.OfferID, xref.OfferItemCriteriaID, xref.UpdatedAt, xref.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update target criteria xref: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a target criteria xref by its unique identifier.
func (r *TarCritOfferXrefRepository) FindByID(ctx context.Context, id int64) (*domain.TarCritOfferXref, error) {
	query := `
		SELECT
			offer_tar_crit_id, offer_id, offer_item_criteria_id, created_at, updated_at
		FROM blc_tar_crit_offer_xref WHERE offer_tar_crit_id = $1`

	var xref domain.TarCritOfferXref
	var offerItemCriteriaID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&xref.ID, &xref.OfferID, &offerItemCriteriaID, &xref.CreatedAt, &xref.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query target criteria xref by ID: %w", err)
	}

	if offerItemCriteriaID.Valid {
		xref.OfferItemCriteriaID = offerItemCriteriaID.Int64
	}

	return &xref, nil
}

// FindByOfferID retrieves all target criteria xrefs for a given offer ID.
func (r *TarCritOfferXrefRepository) FindByOfferID(ctx context.Context, offerID int64) ([]*domain.TarCritOfferXref, error) {
	query := `
		SELECT
			offer_tar_crit_id, offer_id, offer_item_criteria_id, created_at, updated_at
		FROM blc_tar_crit_offer_xref WHERE offer_id = $1`

	rows, err := r.db.QueryContext(ctx, query, offerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query target criteria xrefs by offer ID: %w", err)
	}
	defer rows.Close()

	var xrefs []*domain.TarCritOfferXref
	for rows.Next() {
		var xref domain.TarCritOfferXref
		var offerItemCriteriaID sql.NullInt64

		err := rows.Scan(
			&xref.ID, &xref.OfferID, &offerItemCriteriaID, &xref.CreatedAt, &xref.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan target criteria xref row: %w", err)
		}

		if offerItemCriteriaID.Valid {
			xref.OfferItemCriteriaID = offerItemCriteriaID.Int64
		}
		xrefs = append(xrefs, &xref)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for target criteria xrefs: %w", err)
	}

	return xrefs, nil
}

// FindByOfferItemCriteriaID retrieves all target criteria xrefs for a given offer item criteria ID.
func (r *TarCritOfferXrefRepository) FindByOfferItemCriteriaID(ctx context.Context, offerItemCriteriaID int64) ([]*domain.TarCritOfferXref, error) {
	query := `
		SELECT
			offer_tar_crit_id, offer_id, offer_item_criteria_id, created_at, updated_at
		FROM blc_tar_crit_offer_xref WHERE offer_item_criteria_id = $1`

	rows, err := r.db.QueryContext(ctx, query, offerItemCriteriaID)
	if err != nil {
		return nil, fmt.Errorf("failed to query target criteria xrefs by offer item criteria ID: %w", err)
	}
	defer rows.Close()

	var xrefs []*domain.TarCritOfferXref
	for rows.Next() {
		var xref domain.TarCritOfferXref
		var offID sql.NullInt64 // Use different name to avoid conflict

		err := rows.Scan(
			&xref.ID, &offID, &xref.OfferItemCriteriaID, &xref.CreatedAt, &xref.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan target criteria xref row: %w", err)
		}

		if offID.Valid {
			xref.OfferID = offID.Int64
		}
		xrefs = append(xrefs, &xref)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for target criteria xrefs: %w", err)
	}

	return xrefs, nil
}

// Delete removes a target criteria xref by its unique identifier.
func (r *TarCritOfferXrefRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_tar_crit_offer_xref WHERE offer_tar_crit_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete target criteria xref: %w", err)
	}
	return nil
}

// DeleteByOfferID removes all target criteria xrefs for a given offer ID.
func (r *TarCritOfferXrefRepository) DeleteByOfferID(ctx context.Context, offerID int64) error {
	query := `DELETE FROM blc_tar_crit_offer_xref WHERE offer_id = $1`
	_, err := r.db.ExecContext(ctx, query, offerID)
	if err != nil {
		return fmt.Errorf("failed to delete target criteria xrefs by offer ID: %w", err)
	}
	return nil
}

// DeleteByOfferItemCriteriaID removes all target criteria xrefs for a given offer item criteria ID.
func (r *TarCritOfferXrefRepository) DeleteByOfferItemCriteriaID(ctx context.Context, offerItemCriteriaID int64) error {
	query := `DELETE FROM blc_tar_crit_offer_xref WHERE offer_item_criteria_id = $1`
	_, err := r.db.ExecContext(ctx, query, offerItemCriteriaID)
	if err != nil {
		return fmt.Errorf("failed to delete target criteria xrefs by offer item criteria ID: %w", err)
	}
	return nil
}

// RemoveTarCritOfferXref removes a specific target criteria xref by offer ID and offer item criteria ID.
func (r *TarCritOfferXrefRepository) RemoveTarCritOfferXref(ctx context.Context, offerID, offerItemCriteriaID int64) error {
	query := `DELETE FROM blc_tar_crit_offer_xref WHERE offer_id = $1 AND offer_item_criteria_id = $2`
	_, err := r.db.ExecContext(ctx, query, offerID, offerItemCriteriaID)
	if err != nil {
		return fmt.Errorf("failed to remove target criteria xref: %w", err)
	}
	return nil
}
