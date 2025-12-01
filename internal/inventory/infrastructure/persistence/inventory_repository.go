package persistence

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"

	"github.com/qhato/ecommerce/internal/inventory/domain"
	"github.com/qhato/ecommerce/pkg/database"
	"github.com/qhato/ecommerce/pkg/errors"
)

// PostgresInventoryRepository implements the InventoryRepository interface
type PostgresInventoryRepository struct {
	db *database.DB
}

// NewPostgresInventoryRepository creates a new PostgresInventoryRepository
func NewPostgresInventoryRepository(db *database.DB) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{db: db}
}

// Save stores a new inventory level or updates an existing one.
func (r *PostgresInventoryRepository) Save(ctx context.Context, level *domain.InventoryLevel) error {
	if level.CreatedAt.IsZero() {
		return r.create(ctx, level)
	}
	return r.update(ctx, level)
}

func (r *PostgresInventoryRepository) create(ctx context.Context, level *domain.InventoryLevel) error {
	query := `
		INSERT INTO blc_inventory_level (
			id, sku_id, warehouse_id, location_id, qty_on_hand, qty_reserved,
			qty_available, qty_allocated, qty_backordered, qty_in_transit,
			qty_damaged, reorder_point, reorder_qty, safety_stock,
			allow_backorder, allow_preorder, last_count_date,
			date_created, date_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	err := r.db.Exec(ctx, query,
		level.ID,
		level.SKUID,
		level.WarehouseID,
		level.LocationID,
		level.QuantityOnHand,
		level.QuantityReserved,
		level.QuantityAvailable,
		level.QuantityAllocated,
		level.QuantityBackordered,
		level.QuantityInTransit,
		level.QuantityDamaged,
		level.ReorderPoint,
		level.ReorderQuantity,
		level.SafetyStock,
		level.AllowBackorder,
		level.AllowPreorder,
		level.LastCountDate,
		level.CreatedAt,
		level.UpdatedAt,
	)

	if err != nil {
		return errors.InternalWrap(err, "failed to create inventory level")
	}
	return nil
}

func (r *PostgresInventoryRepository) update(ctx context.Context, level *domain.InventoryLevel) error {
	query := `
		UPDATE blc_inventory_level SET
			sku_id = $2, warehouse_id = $3, location_id = $4, qty_on_hand = $5,
			qty_reserved = $6, qty_available = $7, qty_allocated = $8,
			qty_backordered = $9, qty_in_transit = $10, qty_damaged = $11,
			reorder_point = $12, reorder_qty = $13, safety_stock = $14,
			allow_backorder = $15, allow_preorder = $16, last_count_date = $17,
			date_updated = $18
		WHERE id = $1`

	tag, err := r.db.Pool().Exec(ctx, query,
		level.ID,
		level.SKUID,
		level.WarehouseID,
		level.LocationID,
		level.QuantityOnHand,
		level.QuantityReserved,
		level.QuantityAvailable,
		level.QuantityAllocated,
		level.QuantityBackordered,
		level.QuantityInTransit,
		level.QuantityDamaged,
		level.ReorderPoint,
		level.ReorderQuantity,
		level.SafetyStock,
		level.AllowBackorder,
		level.AllowPreorder,
		level.LastCountDate,
		level.UpdatedAt,
	)

	if err != nil {
		return errors.InternalWrap(err, "failed to update inventory level")
	}
	if tag.RowsAffected() == 0 {
		return errors.NotFound("inventory level not found")
	}
	return nil
}

// FindByID retrieves an inventory level by its unique identifier.
func (r *PostgresInventoryRepository) FindByID(ctx context.Context, id string) (*domain.InventoryLevel, error) {
	query := `
		SELECT
			id, sku_id, warehouse_id, location_id, qty_on_hand, qty_reserved,
			qty_available, qty_allocated, qty_backordered, qty_in_transit,
			qty_damaged, reorder_point, reorder_qty, safety_stock,
			allow_backorder, allow_preorder, last_count_date,
			date_created, date_updated
		FROM blc_inventory_level
		WHERE id = $1`

	level := &domain.InventoryLevel{}
	var (
		warehouseID     sql.NullString
		locationID      sql.NullString
		lastCountDate   sql.NullTime
	)

	err := r.db.QueryRow(ctx, query, id).Scan(
		&level.ID,
		&level.SKUID,
		&warehouseID,
		&locationID,
		&level.QuantityOnHand,
		&level.QuantityReserved,
		&level.QuantityAvailable,
		&level.QuantityAllocated,
		&level.QuantityBackordered,
		&level.QuantityInTransit,
		&level.QuantityDamaged,
		&level.ReorderPoint,
		&level.ReorderQuantity,
		&level.SafetyStock,
		&level.AllowBackorder,
		&level.AllowPreorder,
		&lastCountDate,
		&level.CreatedAt,
		&level.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find inventory level by ID")
	}

	if warehouseID.Valid {
		level.WarehouseID = &warehouseID.String
	}
	if locationID.Valid {
		level.LocationID = &locationID.String
	}
	if lastCountDate.Valid {
		level.LastCountDate = &lastCountDate.Time
	}

	return level, nil
}

// FindBySKUID retrieves an inventory level by its associated SKU ID.
func (r *PostgresInventoryRepository) FindBySKUID(ctx context.Context, skuID string) (*domain.InventoryLevel, error) {
	query := `
		SELECT
			id, sku_id, warehouse_id, location_id, qty_on_hand, qty_reserved,
			qty_available, qty_allocated, qty_backordered, qty_in_transit,
			qty_damaged, reorder_point, reorder_qty, safety_stock,
			allow_backorder, allow_preorder, last_count_date,
			date_created, date_updated
		FROM blc_inventory_level
		WHERE sku_id = $1`

	level := &domain.InventoryLevel{}
	var (
		warehouseID     sql.NullString
		locationID      sql.NullString
		lastCountDate   sql.NullTime
	)

	err := r.db.QueryRow(ctx, query, skuID).Scan(
		&level.ID,
		&level.SKUID,
		&warehouseID,
		&locationID,
		&level.QuantityOnHand,
		&level.QuantityReserved,
		&level.QuantityAvailable,
		&level.QuantityAllocated,
		&level.QuantityBackordered,
		&level.QuantityInTransit,
		&level.QuantityDamaged,
		&level.ReorderPoint,
		&level.ReorderQuantity,
		&level.SafetyStock,
		&level.AllowBackorder,
		&level.AllowPreorder,
		&lastCountDate,
		&level.CreatedAt,
		&level.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find inventory level by SKU ID")
	}

	if warehouseID.Valid {
		level.WarehouseID = &warehouseID.String
	}
	if locationID.Valid {
		level.LocationID = &locationID.String
	}
	if lastCountDate.Valid {
		level.LastCountDate = &lastCountDate.Time
	}

	return level, nil
}

// FindByWarehouse retrieves inventory levels by warehouse.
func (r *PostgresInventoryRepository) FindByWarehouse(ctx context.Context, warehouseID string) ([]*domain.InventoryLevel, error) {
	query := `
		SELECT
			id, sku_id, warehouse_id, location_id, qty_on_hand, qty_reserved,
			qty_available, qty_allocated, qty_backordered, qty_in_transit,
			qty_damaged, reorder_point, reorder_qty, safety_stock,
			allow_backorder, allow_preorder, last_count_date,
			date_created, date_updated
		FROM blc_inventory_level
		WHERE warehouse_id = $1`

	rows, err := r.db.Query(ctx, query, warehouseID)
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find inventory levels by warehouse")
	}
	defer rows.Close()

	var levels []*domain.InventoryLevel
	for rows.Next() {
		level := &domain.InventoryLevel{}
		var (
			whID          sql.NullString
			locID         sql.NullString
			lastCountDate sql.NullTime
		)

		err := rows.Scan(
			&level.ID,
			&level.SKUID,
			&whID,
			&locID,
			&level.QuantityOnHand,
			&level.QuantityReserved,
			&level.QuantityAvailable,
			&level.QuantityAllocated,
			&level.QuantityBackordered,
			&level.QuantityInTransit,
			&level.QuantityDamaged,
			&level.ReorderPoint,
			&level.ReorderQuantity,
			&level.SafetyStock,
			&level.AllowBackorder,
			&level.AllowPreorder,
			&lastCountDate,
			&level.CreatedAt,
			&level.UpdatedAt,
		)
		if err != nil {
			return nil, errors.InternalWrap(err, "failed to scan inventory level")
		}

		if whID.Valid {
			level.WarehouseID = &whID.String
		}
		if locID.Valid {
			level.LocationID = &locID.String
		}
		if lastCountDate.Valid {
			level.LastCountDate = &lastCountDate.Time
		}
		levels = append(levels, level)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.InternalWrap(err, "failed to iterate inventory levels")
	}

	return levels, nil
}

// Delete removes an inventory level by its unique identifier.
func (r *PostgresInventoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM blc_inventory_level WHERE id = $1`
	tag, err := r.db.Pool().Exec(ctx, query, id)
	if err != nil {
		return errors.InternalWrap(err, "failed to delete inventory level")
	}
	if tag.RowsAffected() == 0 {
		return errors.NotFound("inventory level not found")
	}
	return nil
}
