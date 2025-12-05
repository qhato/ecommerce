package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/qhato/ecommerce/internal/pricing/domain"
)

// PostgresPriceListRepository implements PriceListRepository using PostgreSQL
type PostgresPriceListRepository struct {
	db *sql.DB
}

// NewPostgresPriceListRepository creates a new PostgresPriceListRepository
func NewPostgresPriceListRepository(db *sql.DB) domain.PriceListRepository {
	return &PostgresPriceListRepository{db: db}
}

func (r *PostgresPriceListRepository) Save(ctx context.Context, priceList *domain.PriceList) error {
	if priceList.ID == 0 {
		return r.insert(ctx, priceList)
	}
	return r.update(ctx, priceList)
}

func (r *PostgresPriceListRepository) insert(ctx context.Context, priceList *domain.PriceList) error {
	query := `
		INSERT INTO blc_price_list (
			name, code, price_list_type, currency, priority,
			is_active, start_date, end_date, description,
			customer_segments, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		priceList.Name,
		priceList.Code,
		priceList.PriceListType,
		priceList.Currency,
		priceList.Priority,
		priceList.IsActive,
		priceList.StartDate,
		priceList.EndDate,
		priceList.Description,
		pq.Array(priceList.CustomerSegments),
		priceList.CreatedAt,
		priceList.UpdatedAt,
	).Scan(&priceList.ID)

	if err != nil {
		return fmt.Errorf("failed to insert price list: %w", err)
	}
	return nil
}

func (r *PostgresPriceListRepository) update(ctx context.Context, priceList *domain.PriceList) error {
	query := `
		UPDATE blc_price_list
		SET name = $2, priority = $3, is_active = $4,
		    start_date = $5, end_date = $6, description = $7,
		    customer_segments = $8, updated_at = $9
		WHERE id = $1
	`

	priceList.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(
		ctx,
		query,
		priceList.ID,
		priceList.Name,
		priceList.Priority,
		priceList.IsActive,
		priceList.StartDate,
		priceList.EndDate,
		priceList.Description,
		pq.Array(priceList.CustomerSegments),
		priceList.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update price list: %w", err)
	}
	return nil
}

func (r *PostgresPriceListRepository) FindByID(ctx context.Context, id int64) (*domain.PriceList, error) {
	query := `
		SELECT id, name, code, price_list_type, currency, priority,
		       is_active, start_date, end_date, description,
		       customer_segments, created_at, updated_at
		FROM blc_price_list
		WHERE id = $1
	`

	priceList := &domain.PriceList{}
	var customerSegments pq.StringArray

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&priceList.ID,
		&priceList.Name,
		&priceList.Code,
		&priceList.PriceListType,
		&priceList.Currency,
		&priceList.Priority,
		&priceList.IsActive,
		&priceList.StartDate,
		&priceList.EndDate,
		&priceList.Description,
		&customerSegments,
		&priceList.CreatedAt,
		&priceList.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find price list by ID: %w", err)
	}

	priceList.CustomerSegments = customerSegments
	return priceList, nil
}

func (r *PostgresPriceListRepository) FindByCode(ctx context.Context, code string) (*domain.PriceList, error) {
	query := `
		SELECT id, name, code, price_list_type, currency, priority,
		       is_active, start_date, end_date, description,
		       customer_segments, created_at, updated_at
		FROM blc_price_list
		WHERE code = $1
	`

	priceList := &domain.PriceList{}
	var customerSegments pq.StringArray

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&priceList.ID,
		&priceList.Name,
		&priceList.Code,
		&priceList.PriceListType,
		&priceList.Currency,
		&priceList.Priority,
		&priceList.IsActive,
		&priceList.StartDate,
		&priceList.EndDate,
		&priceList.Description,
		&customerSegments,
		&priceList.CreatedAt,
		&priceList.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find price list by code: %w", err)
	}

	priceList.CustomerSegments = customerSegments
	return priceList, nil
}

func (r *PostgresPriceListRepository) FindActive(ctx context.Context, currency string) ([]*domain.PriceList, error) {
	query := `
		SELECT id, name, code, price_list_type, currency, priority,
		       is_active, start_date, end_date, description,
		       customer_segments, created_at, updated_at
		FROM blc_price_list
		WHERE is_active = true
		  AND currency = $1
		  AND (start_date IS NULL OR start_date <= NOW())
		  AND (end_date IS NULL OR end_date >= NOW())
		ORDER BY priority DESC, created_at DESC
	`

	return r.queryPriceLists(ctx, query, currency)
}

func (r *PostgresPriceListRepository) FindByPriority(ctx context.Context, currency string) ([]*domain.PriceList, error) {
	query := `
		SELECT id, name, code, price_list_type, currency, priority,
		       is_active, start_date, end_date, description,
		       customer_segments, created_at, updated_at
		FROM blc_price_list
		WHERE is_active = true
		  AND currency = $1
		  AND (start_date IS NULL OR start_date <= NOW())
		  AND (end_date IS NULL OR end_date >= NOW())
		ORDER BY priority DESC
	`

	return r.queryPriceLists(ctx, query, currency)
}

func (r *PostgresPriceListRepository) FindByCustomerSegment(ctx context.Context, segment string, currency string) ([]*domain.PriceList, error) {
	query := `
		SELECT id, name, code, price_list_type, currency, priority,
		       is_active, start_date, end_date, description,
		       customer_segments, created_at, updated_at
		FROM blc_price_list
		WHERE is_active = true
		  AND currency = $1
		  AND (start_date IS NULL OR start_date <= NOW())
		  AND (end_date IS NULL OR end_date >= NOW())
		  AND ($2 = ANY(customer_segments) OR customer_segments = '{}')
		ORDER BY priority DESC
	`

	return r.queryPriceLists(ctx, query, currency, segment)
}

func (r *PostgresPriceListRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_price_list WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete price list: %w", err)
	}
	return nil
}

func (r *PostgresPriceListRepository) queryPriceLists(ctx context.Context, query string, args ...interface{}) ([]*domain.PriceList, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query price lists: %w", err)
	}
	defer rows.Close()

	priceLists := make([]*domain.PriceList, 0)
	for rows.Next() {
		priceList := &domain.PriceList{}
		var customerSegments pq.StringArray

		err := rows.Scan(
			&priceList.ID,
			&priceList.Name,
			&priceList.Code,
			&priceList.PriceListType,
			&priceList.Currency,
			&priceList.Priority,
			&priceList.IsActive,
			&priceList.StartDate,
			&priceList.EndDate,
			&priceList.Description,
			&customerSegments,
			&priceList.CreatedAt,
			&priceList.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan price list: %w", err)
		}

		priceList.CustomerSegments = customerSegments
		priceLists = append(priceLists, priceList)
	}

	return priceLists, nil
}
