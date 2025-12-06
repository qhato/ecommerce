package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/wishlist/domain"
)

type PostgresWishlistItemRepository struct {
	db *sql.DB
}

func NewPostgresWishlistItemRepository(db *sql.DB) *PostgresWishlistItemRepository {
	return &PostgresWishlistItemRepository{db: db}
}

func (r *PostgresWishlistItemRepository) Create(ctx context.Context, item *domain.WishlistItem) error {
	query := `
		INSERT INTO blc_wishlist_item (id, wishlist_id, product_id, sku_id, quantity, priority, notes, added_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.WishlistID, item.ProductID, item.SKUID,
		item.Quantity, item.Priority, item.Notes, item.AddedAt, item.UpdatedAt,
	)
	return err
}

func (r *PostgresWishlistItemRepository) Update(ctx context.Context, item *domain.WishlistItem) error {
	query := `
		UPDATE blc_wishlist_item
		SET wishlist_id = $1, sku_id = $2, quantity = $3, priority = $4, notes = $5, updated_at = $6
		WHERE id = $7`

	_, err := r.db.ExecContext(ctx, query,
		item.WishlistID, item.SKUID, item.Quantity, item.Priority, item.Notes, item.UpdatedAt, item.ID,
	)
	return err
}

func (r *PostgresWishlistItemRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM blc_wishlist_item WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresWishlistItemRepository) FindByID(ctx context.Context, id string) (*domain.WishlistItem, error) {
	query := `
		SELECT id, wishlist_id, product_id, sku_id, quantity, priority, notes, added_at, updated_at
		FROM blc_wishlist_item
		WHERE id = $1`

	return r.scanWishlistItem(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresWishlistItemRepository) FindByWishlistID(ctx context.Context, wishlistID string) ([]*domain.WishlistItem, error) {
	query := `
		SELECT id, wishlist_id, product_id, sku_id, quantity, priority, notes, added_at, updated_at
		FROM blc_wishlist_item
		WHERE wishlist_id = $1
		ORDER BY priority DESC, added_at DESC`

	return r.queryWishlistItems(ctx, query, wishlistID)
}

func (r *PostgresWishlistItemRepository) ExistsByWishlistAndProduct(ctx context.Context, wishlistID, productID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blc_wishlist_item WHERE wishlist_id = $1 AND product_id = $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, wishlistID, productID).Scan(&exists)
	return exists, err
}

func (r *PostgresWishlistItemRepository) scanWishlistItem(row interface {
	Scan(dest ...interface{}) error
}) (*domain.WishlistItem, error) {
	item := &domain.WishlistItem{}
	err := row.Scan(
		&item.ID, &item.WishlistID, &item.ProductID, &item.SKUID,
		&item.Quantity, &item.Priority, &item.Notes,
		&item.AddedAt, &item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *PostgresWishlistItemRepository) queryWishlistItems(ctx context.Context, query string, args ...interface{}) ([]*domain.WishlistItem, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.WishlistItem
	for rows.Next() {
		item := &domain.WishlistItem{}
		if err := rows.Scan(
			&item.ID, &item.WishlistID, &item.ProductID, &item.SKUID,
			&item.Quantity, &item.Priority, &item.Notes,
			&item.AddedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
