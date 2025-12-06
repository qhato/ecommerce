package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

type PostgresProductBundleItemRepository struct {
	db *sql.DB
}

func NewPostgresProductBundleItemRepository(db *sql.DB) *PostgresProductBundleItemRepository {
	return &PostgresProductBundleItemRepository{db: db}
}

func (r *PostgresProductBundleItemRepository) Create(ctx context.Context, item *domain.ProductBundleItem) error {
	query := `
		INSERT INTO blc_product_bundle_item (bundle_id, product_id, sku_id, quantity, sort_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		item.BundleID, item.ProductID, item.SKUID, item.Quantity, item.SortOrder, item.CreatedAt,
	).Scan(&item.ID)
}

func (r *PostgresProductBundleItemRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_product_bundle_item WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresProductBundleItemRepository) FindByBundleID(ctx context.Context, bundleID int64) ([]*domain.ProductBundleItem, error) {
	query := `
		SELECT id, bundle_id, product_id, sku_id, quantity, sort_order, created_at
		FROM blc_product_bundle_item
		WHERE bundle_id = $1
		ORDER BY sort_order ASC`

	return r.queryBundleItems(ctx, query, bundleID)
}

func (r *PostgresProductBundleItemRepository) DeleteByBundleID(ctx context.Context, bundleID int64) error {
	query := `DELETE FROM blc_product_bundle_item WHERE bundle_id = $1`
	_, err := r.db.ExecContext(ctx, query, bundleID)
	return err
}

func (r *PostgresProductBundleItemRepository) queryBundleItems(ctx context.Context, query string, args ...interface{}) ([]*domain.ProductBundleItem, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.ProductBundleItem
	for rows.Next() {
		item := &domain.ProductBundleItem{}
		if err := rows.Scan(
			&item.ID, &item.BundleID, &item.ProductID, &item.SKUID,
			&item.Quantity, &item.SortOrder, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
