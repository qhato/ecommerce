package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/order/domain"
	"github.com/qhato/ecommerce/pkg/errors"
)

// PostgresOrderRepository implements the OrderRepository interface using PostgreSQL
type PostgresOrderRepository struct {
	db *sql.DB
}

// NewPostgresOrderRepository creates a new PostgresOrderRepository
func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

// Create creates a new order
func (r *PostgresOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.InternalWrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	// Insert order
	query := `
		INSERT INTO blc_order (
			order_number, customer_id, email_address, name, order_status,
			order_subtotal, total_tax, total_shipping, order_total, currency_code,
			submit_date, date_created, date_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING order_id
	`

	err = tx.QueryRowContext(ctx, query,
		order.OrderNumber,
		order.CustomerID,
		order.EmailAddress,
		order.Name,
		order.Status,
		order.SubTotal,
		order.TotalTax,
		order.TotalShipping,
		order.Total,
		order.CurrencyCode,
		order.SubmitDate,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&order.ID)

	if err != nil {
		return errors.InternalWrap(err, "failed to insert order")
	}

	// Insert order items
	if len(order.Items) > 0 {
		itemQuery := `
			INSERT INTO blc_order_item (
				order_id, sku_id, name, quantity, price, total_price,
				tax_amount, shipping_amount
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING order_item_id
		`

		for i := range order.Items {
			item := &order.Items[i]
			item.OrderID = order.ID

			err = tx.QueryRowContext(ctx, itemQuery,
				item.OrderID,
				item.SKUID,
				item.ProductName,
				item.Quantity,
				item.Price,
				item.TotalPrice,
				item.TaxAmount,
				item.ShippingAmount,
			).Scan(&item.ID)

			if err != nil {
				return errors.InternalWrap(err, "failed to insert order item")
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.InternalWrap(err, "failed to commit transaction")
	}

	return nil
}

// Update updates an existing order
func (r *PostgresOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.InternalWrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	// Update order
	query := `
		UPDATE blc_order
		SET order_number = $1, customer_id = $2, email_address = $3, name = $4,
			order_status = $5, order_subtotal = $6, total_tax = $7, total_shipping = $8,
			order_total = $9, currency_code = $10, submit_date = $11, date_updated = $12
		WHERE order_id = $13
	`

	result, err := tx.ExecContext(ctx, query,
		order.OrderNumber,
		order.CustomerID,
		order.EmailAddress,
		order.Name,
		order.Status,
		order.SubTotal,
		order.TotalTax,
		order.TotalShipping,
		order.Total,
		order.CurrencyCode,
		order.SubmitDate,
		order.UpdatedAt,
		order.ID,
	)

	if err != nil {
		return errors.InternalWrap(err, "failed to update order")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.InternalWrap(err, "failed to get rows affected")
	}
	if rowsAffected == 0 {
		return errors.NotFound(fmt.Sprintf("order %d", order.ID))
	}

	// Delete existing items and re-insert
	// This is a simple approach; a more sophisticated one would update/insert/delete as needed
	_, err = tx.ExecContext(ctx, "DELETE FROM blc_order_item WHERE order_id = $1", order.ID)
	if err != nil {
		return errors.InternalWrap(err, "failed to delete order items")
	}

	// Insert order items
	if len(order.Items) > 0 {
		itemQuery := `
			INSERT INTO blc_order_item (
				order_id, sku_id, name, quantity, price, total_price,
				tax_amount, shipping_amount
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING order_item_id
		`

		for i := range order.Items {
			item := &order.Items[i]
			item.OrderID = order.ID

			err = tx.QueryRowContext(ctx, itemQuery,
				item.OrderID,
				item.SKUID,
				item.ProductName,
				item.Quantity,
				item.Price,
				item.TotalPrice,
				item.TaxAmount,
				item.ShippingAmount,
			).Scan(&item.ID)

			if err != nil {
				return errors.InternalWrap(err, "failed to insert order item")
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.InternalWrap(err, "failed to commit transaction")
	}

	return nil
}

// FindByID finds an order by ID
func (r *PostgresOrderRepository) FindByID(ctx context.Context, id int64) (*domain.Order, error) {
	query := `
		SELECT order_id, order_number, customer_id, email_address, name, order_status,
			   order_subtotal, total_tax, total_shipping, order_total, currency_code,
			   submit_date, date_created, date_updated
		FROM blc_order
		WHERE order_id = $1
	`

	order := &domain.Order{}
	var submitDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.OrderNumber,
		&order.CustomerID,
		&order.EmailAddress,
		&order.Name,
		&order.Status,
		&order.SubTotal,
		&order.TotalTax,
		&order.TotalShipping,
		&order.Total,
		&order.CurrencyCode,
		&submitDate,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find order by ID")
	}

	if submitDate.Valid {
		order.SubmitDate = &submitDate.Time
	}

	// Load order items
	items, err := r.findOrderItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

// FindByOrderNumber finds an order by order number
func (r *PostgresOrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	query := `
		SELECT order_id, order_number, customer_id, email_address, name, order_status,
			   order_subtotal, total_tax, total_shipping, order_total, currency_code,
			   submit_date, date_created, date_updated
		FROM blc_order
		WHERE order_number = $1
	`

	order := &domain.Order{}
	var submitDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, orderNumber).Scan(
		&order.ID,
		&order.OrderNumber,
		&order.CustomerID,
		&order.EmailAddress,
		&order.Name,
		&order.Status,
		&order.SubTotal,
		&order.TotalTax,
		&order.TotalShipping,
		&order.Total,
		&order.CurrencyCode,
		&submitDate,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find order by order number")
	}

	if submitDate.Valid {
		order.SubmitDate = &submitDate.Time
	}

	// Load order items
	items, err := r.findOrderItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

// FindByCustomerID finds orders by customer ID
func (r *PostgresOrderRepository) FindByCustomerID(ctx context.Context, customerID int64, filter *domain.OrderFilter) ([]*domain.Order, int64, error) {
	// Build query
	query := `
		SELECT order_id, order_number, customer_id, email_address, name, order_status,
			   order_subtotal, total_tax, total_shipping, order_total, currency_code,
			   submit_date, date_created, date_updated
		FROM blc_order
		WHERE customer_id = $1
	`

	args := []interface{}{customerID}
	argIndex := 2

	// Add status filter if provided
	if filter != nil && filter.Status != "" {
		query += fmt.Sprintf(" AND order_status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM blc_order WHERE customer_id = $1"
	countArgs := []interface{}{customerID}
	if filter != nil && filter.Status != "" {
		countQuery += " AND order_status = $2"
		countArgs = append(countArgs, filter.Status)
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to count orders")
	}

	// Add sorting
	if filter != nil && filter.SortBy != "" {
		sortOrder := "ASC"
		if filter.SortOrder == "DESC" {
			sortOrder = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", filter.SortBy, sortOrder)
	} else {
		query += " ORDER BY date_created DESC"
	}

	// Add pagination
	if filter != nil && filter.PageSize > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)
	}

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to find orders by customer")
	}
	defer rows.Close()

	orders := make([]*domain.Order, 0)
	for rows.Next() {
		order := &domain.Order{}
		var submitDate sql.NullTime

		err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&order.CustomerID,
			&order.EmailAddress,
			&order.Name,
			&order.Status,
			&order.SubTotal,
			&order.TotalTax,
			&order.TotalShipping,
			&order.Total,
			&order.CurrencyCode,
			&submitDate,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.InternalWrap(err, "failed to scan order")
		}

		if submitDate.Valid {
			order.SubmitDate = &submitDate.Time
		}

		// Load order items
		items, err := r.findOrderItems(ctx, order.ID)
		if err != nil {
			return nil, 0, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to iterate orders")
	}

	return orders, total, nil
}

// FindAll finds all orders
func (r *PostgresOrderRepository) FindAll(ctx context.Context, filter *domain.OrderFilter) ([]*domain.Order, int64, error) {
	// Build query
	query := `
		SELECT order_id, order_number, customer_id, email_address, name, order_status,
			   order_subtotal, total_tax, total_shipping, order_total, currency_code,
			   submit_date, date_created, date_updated
		FROM blc_order
		WHERE 1=1
	`

	args := make([]interface{}, 0)
	argIndex := 1

	// Add filters
	if filter != nil {
		if filter.Status != "" {
			query += fmt.Sprintf(" AND order_status = $%d", argIndex)
			args = append(args, filter.Status)
			argIndex++
		}
		if filter.CustomerID > 0 {
			query += fmt.Sprintf(" AND customer_id = $%d", argIndex)
			args = append(args, filter.CustomerID)
			argIndex++
		}
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM blc_order WHERE 1=1"
	countArgs := make([]interface{}, 0)
	countArgIndex := 1
	if filter != nil {
		if filter.Status != "" {
			countQuery += fmt.Sprintf(" AND order_status = $%d", countArgIndex)
			countArgs = append(countArgs, filter.Status)
			countArgIndex++
		}
		if filter.CustomerID > 0 {
			countQuery += fmt.Sprintf(" AND customer_id = $%d", countArgIndex)
			countArgs = append(countArgs, filter.CustomerID)
		}
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to count orders")
	}

	// Add sorting
	if filter != nil && filter.SortBy != "" {
		sortOrder := "ASC"
		if filter.SortOrder == "DESC" {
			sortOrder = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", filter.SortBy, sortOrder)
	} else {
		query += " ORDER BY date_created DESC"
	}

	// Add pagination
	if filter != nil && filter.PageSize > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)
	}

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to find all orders")
	}
	defer rows.Close()

	orders := make([]*domain.Order, 0)
	for rows.Next() {
		order := &domain.Order{}
		var submitDate sql.NullTime

		err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&order.CustomerID,
			&order.EmailAddress,
			&order.Name,
			&order.Status,
			&order.SubTotal,
			&order.TotalTax,
			&order.TotalShipping,
			&order.Total,
			&order.CurrencyCode,
			&submitDate,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.InternalWrap(err, "failed to scan order")
		}

		if submitDate.Valid {
			order.SubmitDate = &submitDate.Time
		}

		// Load order items
		items, err := r.findOrderItems(ctx, order.ID)
		if err != nil {
			return nil, 0, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to iterate orders")
	}

	return orders, total, nil
}

// findOrderItems finds all items for an order
func (r *PostgresOrderRepository) findOrderItems(ctx context.Context, orderID int64) ([]domain.OrderItem, error) {
	query := `
		SELECT order_item_id, order_id, sku_id, name, quantity, price, total_price,
			   tax_amount, shipping_amount
		FROM blc_order_item
		WHERE order_id = $1
		ORDER BY order_item_id
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find order items")
	}
	defer rows.Close()

	items := make([]domain.OrderItem, 0)
	for rows.Next() {
		var item domain.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.SKUID,
			&item.ProductName,
			&item.Quantity,
			&item.Price,
			&item.TotalPrice,
			&item.TaxAmount,
			&item.ShippingAmount,
		)
		if err != nil {
			return nil, errors.InternalWrap(err, "failed to scan order item")
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.InternalWrap(err, "failed to iterate order items")
	}

	return items, nil
}
