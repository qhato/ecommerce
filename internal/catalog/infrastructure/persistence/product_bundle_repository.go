package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

type PostgresProductBundleRepository struct {
	db *sql.DB
}

func NewPostgresProductBundleRepository(db *sql.DB) *PostgresProductBundleRepository {
	return &PostgresProductBundleRepository{db: db}
}

func (r *PostgresProductBundleRepository) Create(ctx context.Context, bundle *domain.ProductBundle) error {
	query := `
		INSERT INTO blc_product_bundle (name, description, bundle_price, is_active, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		bundle.Name, bundle.Description, bundle.BundlePrice, bundle.IsActive,
		bundle.Priority, bundle.CreatedAt, bundle.UpdatedAt,
	).Scan(&bundle.ID)
}

func (r *PostgresProductBundleRepository) Update(ctx context.Context, bundle *domain.ProductBundle) error {
	query := `
		UPDATE blc_product_bundle
		SET name = $1, description = $2, bundle_price = $3, is_active = $4, priority = $5, updated_at = $6
		WHERE id = $7`

	_, err := r.db.ExecContext(ctx, query,
		bundle.Name, bundle.Description, bundle.BundlePrice, bundle.IsActive,
		bundle.Priority, bundle.UpdatedAt, bundle.ID,
	)
	return err
}

func (r *PostgresProductBundleRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_product_bundle WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresProductBundleRepository) FindByID(ctx context.Context, id int64) (*domain.ProductBundle, error) {
	query := `
		SELECT id, name, description, bundle_price, is_active, priority, created_at, updated_at
		FROM blc_product_bundle
		WHERE id = $1`

	return r.scanBundle(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresProductBundleRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.ProductBundle, error) {
	query := `
		SELECT id, name, description, bundle_price, is_active, priority, created_at, updated_at
		FROM blc_product_bundle`

	if activeOnly {
		query += " WHERE is_active = true"
	}

	query += " ORDER BY priority DESC, name ASC"

	return r.queryBundles(ctx, query)
}

func (r *PostgresProductBundleRepository) FindByProduct(ctx context.Context, productID int64) ([]*domain.ProductBundle, error) {
	query := `
		SELECT DISTINCT b.id, b.name, b.description, b.bundle_price, b.is_active, b.priority, b.created_at, b.updated_at
		FROM blc_product_bundle b
		INNER JOIN blc_product_bundle_item i ON b.id = i.bundle_id
		WHERE i.product_id = $1 AND b.is_active = true
		ORDER BY b.priority DESC, b.name ASC`

	return r.queryBundles(ctx, query, productID)
}

func (r *PostgresProductBundleRepository) scanBundle(row interface {
	Scan(dest ...interface{}) error
}) (*domain.ProductBundle, error) {
	bundle := &domain.ProductBundle{}
	err := row.Scan(
		&bundle.ID, &bundle.Name, &bundle.Description, &bundle.BundlePrice,
		&bundle.IsActive, &bundle.Priority, &bundle.CreatedAt, &bundle.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	bundle.Items = make([]domain.ProductBundleItem, 0)
	return bundle, nil
}

func (r *PostgresProductBundleRepository) queryBundles(ctx context.Context, query string, args ...interface{}) ([]*domain.ProductBundle, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bundles []*domain.ProductBundle
	for rows.Next() {
		bundle := &domain.ProductBundle{}
		if err := rows.Scan(
			&bundle.ID, &bundle.Name, &bundle.Description, &bundle.BundlePrice,
			&bundle.IsActive, &bundle.Priority, &bundle.CreatedAt, &bundle.UpdatedAt,
		); err != nil {
			return nil, err
		}
		bundle.Items = make([]domain.ProductBundleItem, 0)
		bundles = append(bundles, bundle)
	}

	return bundles, rows.Err()
}
