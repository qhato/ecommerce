package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/qhato/ecommerce/internal/order/domain"
)

// OrderRepository implements domain.OrderRepository for PostgreSQL persistence.
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new PostgreSQL order repository.
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create stores a new order.
func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle nullable fields
	orderNumber := sql.NullString{String: order.OrderNumber, Valid: order.OrderNumber != ""}
	name := sql.NullString{String: order.Name, Valid: order.Name != ""}
	emailAddress := sql.NullString{String: order.EmailAddress, Valid: order.EmailAddress != ""}
	localeCode := sql.NullString{String: order.LocaleCode, Valid: order.LocaleCode != ""}
	submitDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if order.SubmitDate != nil {
		submitDate = sql.NullTime{Time: *order.SubmitDate, Valid: true}
	}

	query := `
		INSERT INTO blc_order (
			customer_id, email_address, name, order_number, is_preview, order_status, 
			order_subtotal, submit_date, tax_override, order_total, total_shipping, 
			total_tax, currency_code, locale_code, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		) RETURNING order_id`
	err = tx.QueryRowContext(ctx, query,
		order.CustomerID, emailAddress, name, orderNumber, order.IsPreview, order.Status,
		order.OrderSubtotal, submitDate, order.TaxOverride, order.OrderTotal, order.TotalShipping,
		order.TotalTax, order.CurrencyCode, localeCode, order.CreatedAt, order.UpdatedAt,
	).Scan(&order.ID)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	return tx.Commit()
}

// Update updates an existing order.
func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle nullable fields
	orderNumber := sql.NullString{String: order.OrderNumber, Valid: order.OrderNumber != ""}
	name := sql.NullString{String: order.Name, Valid: order.Name != ""}
	emailAddress := sql.NullString{String: order.EmailAddress, Valid: order.EmailAddress != ""}
	localeCode := sql.NullString{String: order.LocaleCode, Valid: order.LocaleCode != ""}
	submitDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if order.SubmitDate != nil {
		submitDate = sql.NullTime{Time: *order.SubmitDate, Valid: true}
	}

	query := `
		UPDATE blc_order SET
			customer_id = $1, email_address = $2, name = $3, order_number = $4, is_preview = $5, 
			order_status = $6, order_subtotal = $7, submit_date = $8, tax_override = $9, 
			order_total = $10, total_shipping = $11, total_tax = $12, currency_code = $13, 
			locale_code = $14, updated_at = $15
		WHERE order_id = $16`
	_, err = tx.ExecContext(ctx, query,
		order.CustomerID, emailAddress, name, orderNumber, order.IsPreview,
		order.Status, order.OrderSubtotal, submitDate, order.TaxOverride,
		order.OrderTotal, order.TotalShipping, order.TotalTax, order.CurrencyCode,
		localeCode, order.UpdatedAt, order.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return tx.Commit()
}

// FindByID retrieves an order by its unique identifier.
func (r *OrderRepository) FindByID(ctx context.Context, id int64) (*domain.Order, error) {
	query := `
		SELECT
			order_id, customer_id, email_address, name, order_number, is_preview, 
			order_status, order_subtotal, submit_date, tax_override, order_total, 
			total_shipping, total_tax, currency_code, locale_code, created_at, updated_at
		FROM blc_order WHERE order_id = $1`

	var order domain.Order
	var orderNumber sql.NullString
	var name sql.NullString
	var emailAddress sql.NullString
	var submitDate sql.NullTime
	var localeCode sql.NullString

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&order.ID, &order.CustomerID, &emailAddress, &name, &orderNumber, &order.IsPreview,
		&order.Status, &order.OrderSubtotal, &submitDate, &order.TaxOverride, &order.OrderTotal,
		&order.TotalShipping, &order.TotalTax, &order.CurrencyCode, &localeCode, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query order by ID: %w", err)
	}

	if orderNumber.Valid {
		order.OrderNumber = orderNumber.String
	}
	if name.Valid {
		order.Name = name.String
	}
	if emailAddress.Valid {
		order.EmailAddress = emailAddress.String
	}
	if submitDate.Valid {
		order.SubmitDate = &submitDate.Time
	}
	if localeCode.Valid {
		order.LocaleCode = localeCode.String
	}

	return &order, nil
}

// FindByOrderNumber retrieves an order by its order number.
func (r *OrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	query := `
		SELECT
			order_id, customer_id, email_address, name, order_number, is_preview, 
			order_status, order_subtotal, submit_date, tax_override, order_total, 
			total_shipping, total_tax, currency_code, locale_code, created_at, updated_at
		FROM blc_order WHERE order_number = $1`

	var order domain.Order
	var orderNumberScan sql.NullString
	var name sql.NullString
	var emailAddress sql.NullString
	var submitDate sql.NullTime
	var localeCode sql.NullString

	row := r.db.QueryRowContext(ctx, query, orderNumber)
	err := row.Scan(
		&order.ID, &order.CustomerID, &emailAddress, &name, &orderNumberScan, &order.IsPreview,
		&order.Status, &order.OrderSubtotal, &submitDate, &order.TaxOverride, &order.OrderTotal,
		&order.TotalShipping, &order.TotalTax, &order.CurrencyCode, &localeCode, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query order by order number: %w", err)
	}

	if orderNumberScan.Valid {
		order.OrderNumber = orderNumberScan.String
	}
	if name.Valid {
		order.Name = name.String
	}
	if emailAddress.Valid {
		order.EmailAddress = emailAddress.String
	}
	if submitDate.Valid {
		order.SubmitDate = &submitDate.Time
	}
	if localeCode.Valid {
		order.LocaleCode = localeCode.String
	}

	return &order, nil
}

// FindByCustomerID retrieves orders by customer ID with pagination and filtering.
func (r *OrderRepository) FindByCustomerID(ctx context.Context, customerID int64, filter *domain.OrderFilter) ([]*domain.Order, int64, error) {
	// Base query
	countQuery := `SELECT COUNT(*) FROM blc_order WHERE customer_id = $1`
	query := `
		SELECT
			order_id, customer_id, email_address, name, order_number, is_preview, 
			order_status, order_subtotal, submit_date, tax_override, order_total, 
			total_shipping, total_tax, currency_code, locale_code, created_at, updated_at
		FROM blc_order WHERE customer_id = $1`

	var args []interface{}
	args = append(args, customerID)
	argIdx := 2

	// Build WHERE clause
	whereClauses := []string{}

	if filter.Status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("order_status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}

	// Apply WHERE clauses
	if len(whereClauses) > 0 {
		query += " AND " + strings.Join(whereClauses, " AND ")
		countQuery += " AND " + strings.Join(whereClauses, " AND ")
	}

	// Count total results
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders by customer ID: %w", err)
	}

	// Apply sorting
	if filter.SortBy != "" {
		orderBy := map[string]string{
			"date_created": "created_at",
			"status":       "order_status",
			"total":        "order_total",
		}
		sortColumn, ok := orderBy[filter.SortBy]
		if !ok {
			sortColumn = "created_at"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, strings.ToUpper(filter.SortOrder))
	}

	// Apply pagination
	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", argIdx, argIdx+1)
	args = append(args, (filter.Page-1)*filter.PageSize, filter.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query orders by customer ID: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		var orderNumberScan sql.NullString
		var name sql.NullString
		var emailAddress sql.NullString
		var submitDate sql.NullTime
		var localeCode sql.NullString

		err := rows.Scan(
			&order.ID, &order.CustomerID, &emailAddress, &name, &orderNumberScan, &order.IsPreview,
			&order.Status, &order.OrderSubtotal, &submitDate, &order.TaxOverride, &order.OrderTotal,
			&order.TotalShipping, &order.TotalTax, &order.CurrencyCode, &localeCode, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order row by customer ID: %w", err)
		}

		if orderNumberScan.Valid {
			order.OrderNumber = orderNumberScan.String
		}
		if name.Valid {
			order.Name = name.String
		}
		if emailAddress.Valid {
			order.EmailAddress = emailAddress.String
		}
		if submitDate.Valid {
			order.SubmitDate = &submitDate.Time
		}
		if localeCode.Valid {
			order.LocaleCode = localeCode.String
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during rows iteration for customer orders: %w", err)
	}

	return orders, totalCount, nil
}

// FindAll retrieves all orders with pagination and filtering.
func (r *OrderRepository) FindAll(ctx context.Context, filter *domain.OrderFilter) ([]*domain.Order, int64, error) {
	// Base query
	countQuery := `SELECT COUNT(*) FROM blc_order`
	query := `
		SELECT
			order_id, customer_id, email_address, name, order_number, is_preview, 
			order_status, order_subtotal, submit_date, tax_override, order_total, 
			total_shipping, total_tax, currency_code, locale_code, created_at, updated_at
		FROM blc_order`

	var args []interface{}
	argIdx := 1

	// Build WHERE clause
	whereClauses := []string{}

	if filter.CustomerID != 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("customer_id = $%d", argIdx))
		args = append(args, filter.CustomerID)
		argIdx++
	}
	if filter.Status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("order_status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}

	// Apply WHERE clauses
	if len(whereClauses) > 0 {
		countQuery += " WHERE " + strings.Join(whereClauses, " AND ")
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total results
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count all orders: %w", err)
	}

	// Apply sorting
	if filter.SortBy != "" {
		orderBy := map[string]string{
			"date_created": "created_at",
			"status":       "order_status",
			"total":        "order_total",
			"customer_id":  "customer_id",
		}
		sortColumn, ok := orderBy[filter.SortBy]
		if !ok {
			sortColumn = "created_at"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, strings.ToUpper(filter.SortOrder))
	}

	// Apply pagination
	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", argIdx, argIdx+1)
	args = append(args, (filter.Page-1)*filter.PageSize, filter.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query all orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		var orderNumber sql.NullString
		var name sql.NullString
		var emailAddress sql.NullString
		var submitDate sql.NullTime
		var localeCode sql.NullString

		err := rows.Scan(
			&order.ID, &order.CustomerID, &emailAddress, &name, &orderNumber, &order.IsPreview,
			&order.Status, &order.OrderSubtotal, &submitDate, &order.TaxOverride, &order.OrderTotal,
			&order.TotalShipping, &order.TotalTax, &order.CurrencyCode, &localeCode, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order row: %w", err)
		}

		if orderNumber.Valid {
			order.OrderNumber = orderNumber.String
		}
		if name.Valid {
			order.Name = name.String
		}
		if emailAddress.Valid {
			order.EmailAddress = emailAddress.String
		}
		if submitDate.Valid {
			order.SubmitDate = &submitDate.Time
		}
		if localeCode.Valid {
			order.LocaleCode = localeCode.String
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during rows iteration for all orders: %w", err)
	}

	return orders, totalCount, nil
}
