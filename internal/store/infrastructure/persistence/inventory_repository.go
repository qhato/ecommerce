package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/store/domain"
)

type PostgresStoreInventoryRepository struct {
	db *sql.DB
}

func NewPostgresStoreInventoryRepository(db *sql.DB) *PostgresStoreInventoryRepository {
	return &PostgresStoreInventoryRepository{db: db}
}

func (r *PostgresStoreInventoryRepository) Create(ctx context.Context, inventory *domain.StoreInventory) error {
	query := `INSERT INTO blc_store_inventory (
		store_id, product_id, sku, quantity_on_hand, reserved, available,
		reorder_point, reorder_quantity, last_restocked, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		inventory.StoreID, inventory.ProductID, inventory.SKU, inventory.QuantityOnHand,
		inventory.Reserved, inventory.Available, inventory.ReorderPoint,
		inventory.ReorderQuantity, inventory.LastRestocked, inventory.UpdatedAt,
	).Scan(&inventory.ID)
}

func (r *PostgresStoreInventoryRepository) Update(ctx context.Context, inventory *domain.StoreInventory) error {
	query := `UPDATE blc_store_inventory SET
		quantity_on_hand = $1, reserved = $2, available = $3, reorder_point = $4,
		reorder_quantity = $5, last_restocked = $6, updated_at = $7
	WHERE id = $8`

	_, err := r.db.ExecContext(ctx, query,
		inventory.QuantityOnHand, inventory.Reserved, inventory.Available,
		inventory.ReorderPoint, inventory.ReorderQuantity, inventory.LastRestocked,
		inventory.UpdatedAt, inventory.ID,
	)
	return err
}

func (r *PostgresStoreInventoryRepository) FindByStoreAndProduct(ctx context.Context, storeID, productID int64) (*domain.StoreInventory, error) {
	query := `SELECT id, store_id, product_id, sku, quantity_on_hand, reserved, available,
		reorder_point, reorder_quantity, last_restocked, updated_at
	FROM blc_store_inventory WHERE store_id = $1 AND product_id = $2`

	return r.scanInventory(r.db.QueryRowContext(ctx, query, storeID, productID))
}

func (r *PostgresStoreInventoryRepository) FindByStore(ctx context.Context, storeID int64) ([]*domain.StoreInventory, error) {
	query := `SELECT id, store_id, product_id, sku, quantity_on_hand, reserved, available,
		reorder_point, reorder_quantity, last_restocked, updated_at
	FROM blc_store_inventory WHERE store_id = $1 ORDER BY sku`

	return r.queryInventory(ctx, query, storeID)
}

func (r *PostgresStoreInventoryRepository) FindByProduct(ctx context.Context, productID int64) ([]*domain.StoreInventory, error) {
	query := `SELECT id, store_id, product_id, sku, quantity_on_hand, reserved, available,
		reorder_point, reorder_quantity, last_restocked, updated_at
	FROM blc_store_inventory WHERE product_id = $1 ORDER BY store_id`

	return r.queryInventory(ctx, query, productID)
}

func (r *PostgresStoreInventoryRepository) FindLowStock(ctx context.Context, storeID int64) ([]*domain.StoreInventory, error) {
	query := `SELECT id, store_id, product_id, sku, quantity_on_hand, reserved, available,
		reorder_point, reorder_quantity, last_restocked, updated_at
	FROM blc_store_inventory
	WHERE store_id = $1 AND available <= reorder_point
	ORDER BY available ASC`

	return r.queryInventory(ctx, query, storeID)
}

func (r *PostgresStoreInventoryRepository) FindBySKU(ctx context.Context, sku string) ([]*domain.StoreInventory, error) {
	query := `SELECT id, store_id, product_id, sku, quantity_on_hand, reserved, available,
		reorder_point, reorder_quantity, last_restocked, updated_at
	FROM blc_store_inventory WHERE sku = $1 ORDER BY store_id`

	return r.queryInventory(ctx, query, sku)
}

func (r *PostgresStoreInventoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_store_inventory WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresStoreInventoryRepository) scanInventory(row interface {
	Scan(dest ...interface{}) error
}) (*domain.StoreInventory, error) {
	inventory := &domain.StoreInventory{}

	err := row.Scan(
		&inventory.ID, &inventory.StoreID, &inventory.ProductID, &inventory.SKU,
		&inventory.QuantityOnHand, &inventory.Reserved, &inventory.Available,
		&inventory.ReorderPoint, &inventory.ReorderQuantity, &inventory.LastRestocked,
		&inventory.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return inventory, nil
}

func (r *PostgresStoreInventoryRepository) queryInventory(ctx context.Context, query string, args ...interface{}) ([]*domain.StoreInventory, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventory := make([]*domain.StoreInventory, 0)
	for rows.Next() {
		inv := &domain.StoreInventory{}

		if err := rows.Scan(
			&inv.ID, &inv.StoreID, &inv.ProductID, &inv.SKU,
			&inv.QuantityOnHand, &inv.Reserved, &inv.Available,
			&inv.ReorderPoint, &inv.ReorderQuantity, &inv.LastRestocked,
			&inv.UpdatedAt,
		); err != nil {
			return nil, err
		}

		inventory = append(inventory, inv)
	}

	return inventory, nil
}
