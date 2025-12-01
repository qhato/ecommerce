package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/offer/domain"
)

// OfferCodeRepository implements domain.OfferCodeRepository for PostgreSQL persistence.
type OfferCodeRepository struct {
	db *sql.DB
}

// NewOfferCodeRepository creates a new PostgreSQL offer code repository.
func NewOfferCodeRepository(db *sql.DB) *OfferCodeRepository {
	return &OfferCodeRepository{db: db}
}

// Save stores a new offer code or updates an existing one.
func (r *OfferCodeRepository) Save(ctx context.Context, offerCode *domain.OfferCode) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle BPCHAR(1) conversion
	archivedChar := "N"
	if offerCode.Archived {
		archivedChar = "Y"
	}

	// Handle nullable fields
	maxUses := sql.NullInt32{Int32: 0, Valid: false}
	if offerCode.MaxUses != nil {
		maxUses = sql.NullInt32{Int32: int32(*offerCode.MaxUses), Valid: true}
	}
	emailAddress := sql.NullString{String: "", Valid: false}
	if offerCode.EmailAddress != nil {
		emailAddress = sql.NullString{String: *offerCode.EmailAddress, Valid: true}
	}
	startDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if offerCode.StartDate != nil {
		startDate = sql.NullTime{Time: *offerCode.StartDate, Valid: true}
	}
	endDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if offerCode.EndDate != nil {
		endDate = sql.NullTime{Time: *offerCode.EndDate, Valid: true}
	}

	if offerCode.ID == 0 {
		// Insert new offer code
		query := `
			INSERT INTO blc_offer_code (
				offer_id, archived, email_address, max_uses, offer_code, 
				end_date, start_date, uses, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
			) RETURNING offer_code_id`
		err = tx.QueryRowContext(ctx, query,
			offerCode.OfferID, archivedChar, emailAddress, maxUses, offerCode.Code,
			endDate, startDate, offerCode.Uses, offerCode.CreatedAt, offerCode.UpdatedAt,
		).Scan(&offerCode.ID)
		if err != nil {
			return fmt.Errorf("failed to insert offer code: %w", err)
		}
	} else {
		// Update existing offer code
		query := `
			UPDATE blc_offer_code SET
				offer_id = $1, archived = $2, email_address = $3, max_uses = $4, 
				offer_code = $5, end_date = $6, start_date = $7, uses = $8, updated_at = $9
			WHERE offer_code_id = $10`
		_, err = tx.ExecContext(ctx, query,
			offerCode.OfferID, archivedChar, emailAddress, maxUses, offerCode.Code,
			endDate, startDate, offerCode.Uses, offerCode.UpdatedAt, offerCode.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update offer code: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves an offer code by its unique identifier.
func (r *OfferCodeRepository) FindByID(ctx context.Context, id int64) (*domain.OfferCode, error) {
	query := `
		SELECT
			offer_code_id, offer_id, archived, email_address, max_uses, offer_code, 
			end_date, start_date, uses, created_at, updated_at
		FROM blc_offer_code WHERE offer_code_id = $1`

	var offerCode domain.OfferCode
	var archivedChar string
	var emailAddress sql.NullString
	var maxUses sql.NullInt32
	var endDate sql.NullTime
	var startDate sql.NullTime

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&offerCode.ID, &offerCode.OfferID, &archivedChar, &emailAddress, &maxUses, &offerCode.Code,
		&endDate, &startDate, &offerCode.Uses, &offerCode.CreatedAt, &offerCode.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query offer code by ID: %w", err)
	}

	offerCode.Archived = (archivedChar == "Y")
	if emailAddress.Valid {
		offerCode.EmailAddress = &emailAddress.String
	}
	if maxUses.Valid {
		intMaxUses := int(maxUses.Int32)
		offerCode.MaxUses = &intMaxUses
	}
	if endDate.Valid {
		offerCode.EndDate = &endDate.Time
	}
	if startDate.Valid {
		offerCode.StartDate = &startDate.Time
	}

	return &offerCode, nil
}

// FindByCode retrieves an offer code by its code string.
func (r *OfferCodeRepository) FindByCode(ctx context.Context, code string) (*domain.OfferCode, error) {
	query := `
		SELECT
			offer_code_id, offer_id, archived, email_address, max_uses, offer_code, 
			end_date, start_date, uses, created_at, updated_at
		FROM blc_offer_code WHERE offer_code = $1`

	var offerCode domain.OfferCode
	var archivedChar string
	var emailAddress sql.NullString
	var maxUses sql.NullInt32
	var endDate sql.NullTime
	var startDate sql.NullTime

	row := r.db.QueryRowContext(ctx, query, code)
	err := row.Scan(
		&offerCode.ID, &offerCode.OfferID, &archivedChar, &emailAddress, &maxUses, &offerCode.Code,
		&endDate, &startDate, &offerCode.Uses, &offerCode.CreatedAt, &offerCode.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query offer code by code: %w", err)
	}

	offerCode.Archived = (archivedChar == "Y")
	if emailAddress.Valid {
		offerCode.EmailAddress = &emailAddress.String
	}
	if maxUses.Valid {
		intMaxUses := int(maxUses.Int32)
		offerCode.MaxUses = &intMaxUses
	}
	if endDate.Valid {
		offerCode.EndDate = &endDate.Time
	}
	if startDate.Valid {
		offerCode.StartDate = &startDate.Time
	}

	return &offerCode, nil
}

// FindByOfferID retrieves all offer codes associated with a given offer ID.
func (r *OfferCodeRepository) FindByOfferID(ctx context.Context, offerID int64) ([]*domain.OfferCode, error) {
	query := `
		SELECT
			offer_code_id, offer_id, archived, email_address, max_uses, offer_code, 
			end_date, start_date, uses, created_at, updated_at
		FROM blc_offer_code WHERE offer_id = $1`

	rows, err := r.db.QueryContext(ctx, query, offerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query offer codes by offer ID: %w", err)
	}
	defer rows.Close()

	var offerCodes []*domain.OfferCode
	for rows.Next() {
		var offerCode domain.OfferCode
		var archivedChar string
		var emailAddress sql.NullString
		var maxUses sql.NullInt32
		var endDate sql.NullTime
		var startDate sql.NullTime

		err := rows.Scan(
			&offerCode.ID, &offerCode.OfferID, &archivedChar, &emailAddress, &maxUses, &offerCode.Code,
			&endDate, &startDate, &offerCode.Uses, &offerCode.CreatedAt, &offerCode.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan offer code row: %w", err)
		}

		offerCode.Archived = (archivedChar == "Y")
		if emailAddress.Valid {
			offerCode.EmailAddress = &emailAddress.String
		}
		if maxUses.Valid {
			intMaxUses := int(maxUses.Int32)
			offerCode.MaxUses = &intMaxUses
		}
		if endDate.Valid {
			offerCode.EndDate = &endDate.Time
		}
		if startDate.Valid {
			offerCode.StartDate = &startDate.Time
		}
		offerCodes = append(offerCodes, &offerCode)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for offer codes: %w", err)
	}

	return offerCodes, nil
}

// Delete removes an offer code by its unique identifier.
func (r *OfferCodeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_offer_code WHERE offer_code_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer code: %w", err)
	}
	return nil
}

// DeleteByOfferID removes all offer codes associated with a given offer ID.
func (r *OfferCodeRepository) DeleteByOfferID(ctx context.Context, offerID int64) error {
	query := `DELETE FROM blc_offer_code WHERE offer_id = $1`
	_, err := r.db.ExecContext(ctx, query, offerID)
	if err != nil {
		return fmt.Errorf("failed to delete offer codes by offer ID: %w", err)
	}
	return nil
}
