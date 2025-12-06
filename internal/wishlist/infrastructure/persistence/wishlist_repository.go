package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/wishlist/domain"
)

type PostgresWishlistRepository struct {
	db *sql.DB
}

func NewPostgresWishlistRepository(db *sql.DB) *PostgresWishlistRepository {
	return &PostgresWishlistRepository{db: db}
}

func (r *PostgresWishlistRepository) Create(ctx context.Context, wishlist *domain.Wishlist) error {
	query := `
		INSERT INTO blc_wishlist (id, customer_id, name, is_default, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		wishlist.ID, wishlist.CustomerID, wishlist.Name, wishlist.IsDefault,
		wishlist.IsPublic, wishlist.CreatedAt, wishlist.UpdatedAt,
	)
	return err
}

func (r *PostgresWishlistRepository) Update(ctx context.Context, wishlist *domain.Wishlist) error {
	query := `
		UPDATE blc_wishlist
		SET name = $1, is_default = $2, is_public = $3, updated_at = $4
		WHERE id = $5`

	_, err := r.db.ExecContext(ctx, query,
		wishlist.Name, wishlist.IsDefault, wishlist.IsPublic, wishlist.UpdatedAt, wishlist.ID,
	)
	return err
}

func (r *PostgresWishlistRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM blc_wishlist WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresWishlistRepository) FindByID(ctx context.Context, id string) (*domain.Wishlist, error) {
	query := `
		SELECT id, customer_id, name, is_default, is_public, created_at, updated_at
		FROM blc_wishlist
		WHERE id = $1`

	return r.scanWishlist(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresWishlistRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*domain.Wishlist, error) {
	query := `
		SELECT id, customer_id, name, is_default, is_public, created_at, updated_at
		FROM blc_wishlist
		WHERE customer_id = $1
		ORDER BY is_default DESC, created_at DESC`

	return r.queryWishlists(ctx, query, customerID)
}

func (r *PostgresWishlistRepository) FindDefaultByCustomerID(ctx context.Context, customerID string) (*domain.Wishlist, error) {
	query := `
		SELECT id, customer_id, name, is_default, is_public, created_at, updated_at
		FROM blc_wishlist
		WHERE customer_id = $1 AND is_default = true`

	return r.scanWishlist(r.db.QueryRowContext(ctx, query, customerID))
}

func (r *PostgresWishlistRepository) FindPublicByID(ctx context.Context, id string) (*domain.Wishlist, error) {
	query := `
		SELECT id, customer_id, name, is_default, is_public, created_at, updated_at
		FROM blc_wishlist
		WHERE id = $1 AND is_public = true`

	return r.scanWishlist(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresWishlistRepository) scanWishlist(row interface {
	Scan(dest ...interface{}) error
}) (*domain.Wishlist, error) {
	wishlist := &domain.Wishlist{}
	err := row.Scan(
		&wishlist.ID, &wishlist.CustomerID, &wishlist.Name,
		&wishlist.IsDefault, &wishlist.IsPublic,
		&wishlist.CreatedAt, &wishlist.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	wishlist.Items = make([]domain.WishlistItem, 0)
	return wishlist, nil
}

func (r *PostgresWishlistRepository) queryWishlists(ctx context.Context, query string, args ...interface{}) ([]*domain.Wishlist, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlists []*domain.Wishlist
	for rows.Next() {
		wishlist := &domain.Wishlist{}
		if err := rows.Scan(
			&wishlist.ID, &wishlist.CustomerID, &wishlist.Name,
			&wishlist.IsDefault, &wishlist.IsPublic,
			&wishlist.CreatedAt, &wishlist.UpdatedAt,
		); err != nil {
			return nil, err
		}
		wishlist.Items = make([]domain.WishlistItem, 0)
		wishlists = append(wishlists, wishlist)
	}

	return wishlists, rows.Err()
}
