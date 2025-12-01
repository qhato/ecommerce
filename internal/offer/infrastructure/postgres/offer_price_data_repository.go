package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/offer/domain"
)

// OfferPriceDataRepository implements domain.OfferPriceDataRepository for PostgreSQL persistence.
type OfferPriceDataRepository struct {
	db *sql.DB
}

// NewOfferPriceDataRepository creates a new PostgreSQL offer price data repository.
func NewOfferPriceDataRepository(db *sql.DB) *OfferPriceDataRepository {
	return &OfferPriceDataRepository{db: db}
}

// Save stores new offer price data or updates an existing one.
func (r *OfferPriceDataRepository) Save(ctx context.Context, priceData *domain.OfferPriceData) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle BPCHAR(1) conversion
	archivedChar := "N"
	if priceData.Archived {
		archivedChar = "Y"
	}

	// Handle nullable fields
	endDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if priceData.EndDate != nil {
		endDate = sql.NullTime{Time: *priceData.EndDate, Valid: true}
	}
	startDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if priceData.StartDate != nil {
		startDate = sql.NullTime{Time: *priceData.StartDate, Valid: true}
	}

	identifierType := sql.NullString{String: priceData.IdentifierType, Valid: priceData.IdentifierType != ""}
	identifierValue := sql.NullString{String: priceData.IdentifierValue, Valid: priceData.IdentifierValue != ""}

	if priceData.ID == 0 {
		// Insert new offer price data
		query := `
			INSERT INTO blc_offer_price_data (
				offer_id, amount, discount_type, identifier_type, identifier_value, 
				quantity, start_date, end_date, archived, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
			) RETURNING offer_price_data_id`
		err = tx.QueryRowContext(ctx, query,
			priceData.OfferID, priceData.Amount, priceData.DiscountType, identifierType, identifierValue,
			priceData.Quantity, startDate, endDate, archivedChar, priceData.CreatedAt, priceData.UpdatedAt,
		).Scan(&priceData.ID)
		if err != nil {
			return fmt.Errorf("failed to insert offer price data: %w", err)
		}
	} else {
		// Update existing offer price data
		query := `
			UPDATE blc_offer_price_data SET
				offer_id = $1, amount = $2, discount_type = $3, identifier_type = $4, 
				identifier_value = $5, quantity = $6, start_date = $7, end_date = $8, 
				archived = $9, updated_at = $10
			WHERE offer_price_data_id = $11`
		_, err = tx.ExecContext(ctx, query,
			priceData.OfferID, priceData.Amount, priceData.DiscountType, identifierType, identifierValue,
			priceData.Quantity, startDate, endDate, archivedChar, priceData.UpdatedAt, priceData.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update offer price data: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves offer price data by its unique identifier.
func (r *OfferPriceDataRepository) FindByID(ctx context.Context, id int64) (*domain.OfferPriceData, error) {
	query := `
		SELECT
			offer_price_data_id, offer_id, amount, discount_type, identifier_type, 
			identifier_value, quantity, start_date, end_date, archived, created_at, updated_at
		FROM blc_offer_price_data WHERE offer_price_data_id = $1`

	var priceData domain.OfferPriceData
	var archivedChar string
	var identifierType sql.NullString
	var identifierValue sql.NullString
	var startDate sql.NullTime
	var endDate sql.NullTime

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&priceData.ID, &priceData.OfferID, &priceData.Amount, &priceData.DiscountType, &identifierType,
		&identifierValue, &priceData.Quantity, &startDate, &endDate, &archivedChar, &priceData.CreatedAt, &priceData.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query offer price data by ID: %w", err)
	}

	priceData.Archived = (archivedChar == "Y")
	if identifierType.Valid {
		priceData.IdentifierType = identifierType.String
	}
	if identifierValue.Valid {
		priceData.IdentifierValue = identifierValue.String
	}
	if startDate.Valid {
		priceData.StartDate = &startDate.Time
	}
	if endDate.Valid {
		priceData.EndDate = &endDate.Time
	}

	return &priceData, nil
}

// FindByOfferID retrieves all offer price data associated with a given offer ID.
func (r *OfferPriceDataRepository) FindByOfferID(ctx context.Context, offerID int64) ([]*domain.OfferPriceData, error) {
	query := `
		SELECT
			offer_price_data_id, offer_id, amount, discount_type, identifier_type, 
			identifier_value, quantity, start_date, end_date, archived, created_at, updated_at
		FROM blc_offer_price_data WHERE offer_id = $1`

	rows, err := r.db.QueryContext(ctx, query, offerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query offer price data by offer ID: %w", err)
	}
	defer rows.Close()

	var priceDataList []*domain.OfferPriceData
	for rows.Next() {
		var priceData domain.OfferPriceData
		var archivedChar string
		var identifierType sql.NullString
		var identifierValue sql.NullString
		var startDate sql.NullTime
		var endDate sql.NullTime

		err := rows.Scan(
			&priceData.ID, &priceData.OfferID, &priceData.Amount, &priceData.DiscountType, &identifierType,
			&identifierValue, &priceData.Quantity, &startDate, &endDate, &archivedChar, &priceData.CreatedAt, &priceData.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan offer price data row: %w", err)
		}

		priceData.Archived = (archivedChar == "Y")
		if identifierType.Valid {
			priceData.IdentifierType = identifierType.String
		}
		if identifierValue.Valid {
			priceData.IdentifierValue = identifierValue.String
		}
		if startDate.Valid {
			priceData.StartDate = &startDate.Time
		}
		if endDate.Valid {
			priceData.EndDate = &endDate.Time
		}
		priceDataList = append(priceDataList, &priceData)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for offer price data: %w", err)
	}

	return priceDataList, nil
}

// FindActiveByOfferID retrieves all currently active offer price data for a given offer ID. (Based on StartDate/EndDate)
func (r *OfferPriceDataRepository) FindActiveByOfferID(ctx context.Context, offerID int64) ([]*domain.OfferPriceData, error) {
	query := `
		SELECT
			offer_price_data_id, offer_id, amount, discount_type, identifier_type, 
			identifier_value, quantity, start_date, end_date, archived, created_at, updated_at
		FROM blc_offer_price_data 
		WHERE offer_id = $1 AND archived = 'N' AND start_date <= NOW() AND (end_date IS NULL OR end_date >= NOW())`

	rows, err := r.db.QueryContext(ctx, query, offerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query active offer price data by offer ID: %w", err)
	}
	defer rows.Close()

	var priceDataList []*domain.OfferPriceData
	for rows.Next() {
		var priceData domain.OfferPriceData
		var archivedChar string
		var identifierType sql.NullString
		var identifierValue sql.NullString
		var startDate sql.NullTime
		var endDate sql.NullTime

		err := rows.Scan(
			&priceData.ID, &priceData.OfferID, &priceData.Amount, &priceData.DiscountType, &identifierType,
			&identifierValue, &priceData.Quantity, &startDate, &endDate, &archivedChar, &priceData.CreatedAt, &priceData.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan active offer price data row: %w", err)
		}

		priceData.Archived = (archivedChar == "Y")
		if identifierType.Valid {
			priceData.IdentifierType = identifierType.String
		}
		if identifierValue.Valid {
			priceData.IdentifierValue = identifierValue.String
		}
		if startDate.Valid {
			priceData.StartDate = &startDate.Time
		}
		if endDate.Valid {
			priceData.EndDate = &endDate.Time
		}
		priceDataList = append(priceDataList, &priceData)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for active offer price data: %w", err)
	}

	return priceDataList, nil
}

// Delete removes offer price data by its unique identifier.
func (r *OfferPriceDataRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_offer_price_data WHERE offer_price_data_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer price data: %w", err)
	}
	return nil
}
