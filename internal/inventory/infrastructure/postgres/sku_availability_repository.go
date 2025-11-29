package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/inventory/domain"
)

// InventoryRepository implements domain.InventoryRepository for PostgreSQL persistence.
type InventoryRepository struct {
	db *sql.DB
}

// NewInventoryRepository creates a new PostgreSQL inventory repository.
func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// Save stores a new SKU availability record or updates an existing one.
func (r *InventoryRepository) Save(ctx context.Context, availability *domain.SKUAvailability) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle nullable fields
	availabilityDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if availability.AvailabilityDate != nil {
		availabilityDate = sql.NullTime{Time: *availability.AvailabilityDate, Valid: true}
	}
	locationID := sql.NullInt64{Int64: 0, Valid: false}
	if availability.LocationID != nil {
		locationID = sql.NullInt64{Int64: *availability.LocationID, Valid: true}
	}

	if availability.ID == 0 {
		// Insert new SKU availability
		query := `
			INSERT INTO blc_sku_availability (
				sku_id, availability_date, availability_status, location_id, 
				qty_on_hand, reserve_qty, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8
			) RETURNING sku_availability_id`
		err = tx.QueryRowContext(ctx, query,
			availability.SkuID, availabilityDate, availability.AvailabilityStatus, locationID,
			availability.QtyOnHand, availability.ReserveQty, availability.CreatedAt, availability.UpdatedAt,
		).Scan(&availability.ID)
		if err != nil {
			return fmt.Errorf("failed to insert SKU availability: %w", err)
		}
	} else {
		// Update existing SKU availability
		query := `
			UPDATE blc_sku_availability SET
				sku_id = $1, availability_date = $2, availability_status = $3, 
				location_id = $4, qty_on_hand = $5, reserve_qty = $6, updated_at = $7
			WHERE sku_availability_id = $8`
		_, err = tx.ExecContext(ctx, query,
			availability.SkuID, availabilityDate, availability.AvailabilityStatus, locationID,
			availability.QtyOnHand, availability.ReserveQty, availability.UpdatedAt, availability.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update SKU availability: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a SKU availability record by its unique identifier.
func (r *InventoryRepository) FindByID(ctx context.Context, id int64) (*domain.SKUAvailability, error) {
	query := `
		SELECT
			sku_availability_id, sku_id, availability_date, availability_status, 
			location_id, qty_on_hand, reserve_qty, created_at, updated_at
		FROM blc_sku_availability WHERE sku_availability_id = $1`

	var availability domain.SKUAvailability
	var availabilityDate sql.NullTime
	var locationID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&availability.ID, &availability.SkuID, &availabilityDate, &availability.AvailabilityStatus,
		&locationID, &availability.QtyOnHand, &availability.ReserveQty, &availability.CreatedAt, &availability.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query SKU availability by ID: %w", err)
	}

	if availabilityDate.Valid {
		availability.AvailabilityDate = &availabilityDate.Time
	}
	if locationID.Valid {
		availability.LocationID = &locationID.Int64
	}

	return &availability, nil
}

// FindBySKUID retrieves a SKU availability record by SKU ID.
func (r *InventoryRepository) FindBySKUID(ctx context.Context, skuID int64) (*domain.SKUAvailability, error) {
	query := `
		SELECT
			sku_availability_id, sku_id, availability_date, availability_status, 
			location_id, qty_on_hand, reserve_qty, created_at, updated_at
		FROM blc_sku_availability WHERE sku_id = $1`

	var availability domain.SKUAvailability
	var availabilityDate sql.NullTime
	var locationID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, skuID)
	err := row.Scan(
		&availability.ID, &availability.SkuID, &availabilityDate, &availability.AvailabilityStatus,
		&locationID, &availability.QtyOnHand, &availability.ReserveQty, &availability.CreatedAt, &availability.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query SKU availability by SKU ID: %w", err)
	}

	if availabilityDate.Valid {
		availability.AvailabilityDate = &availabilityDate.Time
	}
	if locationID.Valid {
		availability.LocationID = &locationID.Int64
	}

	return &availability, nil
}

// Delete removes a SKU availability record by its unique identifier.
func (r *InventoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_sku_availability WHERE sku_availability_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete SKU availability: %w", err)
	}
	return nil
}

// DeleteBySKUID removes all SKU availability records for a given SKU ID.
func (r *InventoryRepository) DeleteBySKUID(ctx context.Context, skuID int64) error {
	query := `DELETE FROM blc_sku_availability WHERE sku_id = $1`
	_, err := r.db.ExecContext(ctx, query, skuID)
	if err != nil {
		return fmt.Errorf("failed to delete SKU availability by SKU ID: %w", err)
	}
	return nil
}
