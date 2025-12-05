package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/pricing/domain"
	"github.com/shopspring/decimal"
)

// PostgresPriceListItemRepository implements PriceListItemRepository using PostgreSQL
type PostgresPriceListItemRepository struct {
	db *sql.DB
}

// NewPostgresPriceListItemRepository creates a new PostgresPriceListItemRepository
func NewPostgresPriceListItemRepository(db *sql.DB) domain.PriceListItemRepository {
	return &PostgresPriceListItemRepository{db: db}
}

func (r *PostgresPriceListItemRepository) Save(ctx context.Context, item *domain.PriceListItem) error {
	if item.ID == 0 {
		return r.insert(ctx, item)
	}
	return r.update(ctx, item)
}

func (r *PostgresPriceListItemRepository) insert(ctx context.Context, item *domain.PriceListItem) error {
	query := `
		INSERT INTO blc_price_list_item (
			price_list_id, sku_id, product_id, price, compare_at_price,
			min_quantity, max_quantity, is_active, start_date, end_date,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		item.PriceListID,
		item.SKUID,
		item.ProductID,
		item.Price,
		item.CompareAtPrice,
		item.MinQuantity,
		item.MaxQuantity,
		item.IsActive,
		item.StartDate,
		item.EndDate,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(&item.ID)

	if err != nil {
		return fmt.Errorf("failed to insert price list item: %w", err)
	}
	return nil
}

func (r *PostgresPriceListItemRepository) update(ctx context.Context, item *domain.PriceListItem) error {
	query := `
		UPDATE blc_price_list_item
		SET price = $2, compare_at_price = $3, min_quantity = $4,
		    max_quantity = $5, is_active = $6, start_date = $7,
		    end_date = $8, updated_at = $9
		WHERE id = $1
	`

	item.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(
		ctx,
		query,
		item.ID,
		item.Price,
		item.CompareAtPrice,
		item.MinQuantity,
		item.MaxQuantity,
		item.IsActive,
		item.StartDate,
		item.EndDate,
		item.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update price list item: %w", err)
	}
	return nil
}

func (r *PostgresPriceListItemRepository) FindByID(ctx context.Context, id int64) (*domain.PriceListItem, error) {
	query := `
		SELECT id, price_list_id, sku_id, product_id, price, compare_at_price,
		       min_quantity, max_quantity, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_price_list_item
		WHERE id = $1
	`

	return r.scanPriceListItem(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresPriceListItemRepository) FindByPriceListID(ctx context.Context, priceListID int64) ([]*domain.PriceListItem, error) {
	query := `
		SELECT id, price_list_id, sku_id, product_id, price, compare_at_price,
		       min_quantity, max_quantity, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_price_list_item
		WHERE price_list_id = $1
		ORDER BY sku_id
	`

	return r.queryPriceListItems(ctx, query, priceListID)
}

func (r *PostgresPriceListItemRepository) FindBySKU(ctx context.Context, skuID string) ([]*domain.PriceListItem, error) {
	query := `
		SELECT id, price_list_id, sku_id, product_id, price, compare_at_price,
		       min_quantity, max_quantity, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_price_list_item
		WHERE sku_id = $1
		ORDER BY price_list_id
	`

	return r.queryPriceListItems(ctx, query, skuID)
}

func (r *PostgresPriceListItemRepository) FindBySKUAndPriceList(ctx context.Context, skuID string, priceListID int64) (*domain.PriceListItem, error) {
	query := `
		SELECT id, price_list_id, sku_id, product_id, price, compare_at_price,
		       min_quantity, max_quantity, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_price_list_item
		WHERE sku_id = $1 AND price_list_id = $2
	`

	return r.scanPriceListItem(r.db.QueryRowContext(ctx, query, skuID, priceListID))
}

func (r *PostgresPriceListItemRepository) FindActiveForSKU(ctx context.Context, skuID string, quantity int) ([]*domain.PriceListItem, error) {
	query := `
		SELECT id, price_list_id, sku_id, product_id, price, compare_at_price,
		       min_quantity, max_quantity, is_active, start_date, end_date,
		       created_at, updated_at
		FROM blc_price_list_item
		WHERE sku_id = $1
		  AND is_active = true
		  AND (start_date IS NULL OR start_date <= NOW())
		  AND (end_date IS NULL OR end_date >= NOW())
		  AND min_quantity <= $2
		  AND (max_quantity IS NULL OR max_quantity >= $2)
		ORDER BY price_list_id
	`

	return r.queryPriceListItems(ctx, query, skuID, quantity)
}

func (r *PostgresPriceListItemRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_price_list_item WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete price list item: %w", err)
	}
	return nil
}

func (r *PostgresPriceListItemRepository) DeleteByPriceListID(ctx context.Context, priceListID int64) error {
	query := `DELETE FROM blc_price_list_item WHERE price_list_id = $1`

	_, err := r.db.ExecContext(ctx, query, priceListID)
	if err != nil {
		return fmt.Errorf("failed to delete price list items: %w", err)
	}
	return nil
}

func (r *PostgresPriceListItemRepository) scanPriceListItem(row *sql.Row) (*domain.PriceListItem, error) {
	item := &domain.PriceListItem{}
	var price, compareAtPrice sql.NullString

	err := row.Scan(
		&item.ID,
		&item.PriceListID,
		&item.SKUID,
		&item.ProductID,
		&price,
		&compareAtPrice,
		&item.MinQuantity,
		&item.MaxQuantity,
		&item.IsActive,
		&item.StartDate,
		&item.EndDate,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan price list item: %w", err)
	}

	// Parse decimal values
	if price.Valid {
		item.Price, _ = decimal.NewFromString(price.String)
	}
	if compareAtPrice.Valid {
		cap, _ := decimal.NewFromString(compareAtPrice.String)
		item.CompareAtPrice = &cap
	}

	return item, nil
}

func (r *PostgresPriceListItemRepository) queryPriceListItems(ctx context.Context, query string, args ...interface{}) ([]*domain.PriceListItem, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query price list items: %w", err)
	}
	defer rows.Close()

	items := make([]*domain.PriceListItem, 0)
	for rows.Next() {
		item := &domain.PriceListItem{}
		var price, compareAtPrice sql.NullString

		err := rows.Scan(
			&item.ID,
			&item.PriceListID,
			&item.SKUID,
			&item.ProductID,
			&price,
			&compareAtPrice,
			&item.MinQuantity,
			&item.MaxQuantity,
			&item.IsActive,
			&item.StartDate,
			&item.EndDate,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan price list item: %w", err)
		}

		// Parse decimal values
		if price.Valid {
			item.Price, _ = decimal.NewFromString(price.String)
		}
		if compareAtPrice.Valid {
			cap, _ := decimal.NewFromString(compareAtPrice.String)
			item.CompareAtPrice = &cap
		}

		items = append(items, item)
	}

	return items, nil
}
