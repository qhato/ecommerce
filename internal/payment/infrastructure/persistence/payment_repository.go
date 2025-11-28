package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/payment/domain"
	"github.com/qhato/ecommerce/pkg/apperrors"
)

// PostgresPaymentRepository implements the PaymentRepository interface using PostgreSQL
type PostgresPaymentRepository struct {
	db *sql.DB
}

// NewPostgresPaymentRepository creates a new PostgresPaymentRepository
func NewPostgresPaymentRepository(db *sql.DB) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{db: db}
}

// Create creates a new payment
func (r *PostgresPaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO blc_order_payment (
			order_id, customer_id, type, amount, currency_code,
			transaction_id, gateway_response_code, authorization_code,
			refund_amount, failure_reason, processed_date, authorized_date,
			captured_date, refunded_date, date_created, date_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING payment_id
	`

	err := r.db.QueryRowContext(ctx, query,
		payment.OrderID,
		payment.CustomerID,
		payment.PaymentMethod,
		payment.Amount,
		payment.CurrencyCode,
		payment.TransactionID,
		payment.GatewayResponse,
		payment.AuthorizationCode,
		payment.RefundAmount,
		payment.FailureReason,
		payment.ProcessedDate,
		payment.AuthorizedDate,
		payment.CapturedDate,
		payment.RefundedDate,
		payment.CreatedAt,
		payment.UpdatedAt,
	).Scan(&payment.ID)

	if err != nil {
		return apperrors.NewInternalError("failed to create payment", err)
	}

	return nil
}

// Update updates an existing payment
func (r *PostgresPaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	query := `
		UPDATE blc_order_payment
		SET order_id = $1, customer_id = $2, type = $3, amount = $4,
			currency_code = $5, transaction_id = $6, gateway_response_code = $7,
			authorization_code = $8, refund_amount = $9, failure_reason = $10,
			processed_date = $11, authorized_date = $12, captured_date = $13,
			refunded_date = $14, date_updated = $15
		WHERE payment_id = $16
	`

	result, err := r.db.ExecContext(ctx, query,
		payment.OrderID,
		payment.CustomerID,
		payment.PaymentMethod,
		payment.Amount,
		payment.CurrencyCode,
		payment.TransactionID,
		payment.GatewayResponse,
		payment.AuthorizationCode,
		payment.RefundAmount,
		payment.FailureReason,
		payment.ProcessedDate,
		payment.AuthorizedDate,
		payment.CapturedDate,
		payment.RefundedDate,
		payment.UpdatedAt,
		payment.ID,
	)

	if err != nil {
		return apperrors.NewInternalError("failed to update payment", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewInternalError("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return apperrors.NewNotFoundError("payment", payment.ID)
	}

	return nil
}

// FindByID finds a payment by ID
func (r *PostgresPaymentRepository) FindByID(ctx context.Context, id int64) (*domain.Payment, error) {
	query := `
		SELECT payment_id, order_id, customer_id, type, amount, currency_code,
			   transaction_id, gateway_response_code, authorization_code, refund_amount,
			   failure_reason, processed_date, authorized_date, captured_date, refunded_date,
			   date_created, date_updated
		FROM blc_order_payment
		WHERE payment_id = $1
	`

	payment := &domain.Payment{}
	var (
		processedDate  sql.NullTime
		authorizedDate sql.NullTime
		capturedDate   sql.NullTime
		refundedDate   sql.NullTime
		transactionID  sql.NullString
		gatewayResponse sql.NullString
		authCode       sql.NullString
		failureReason  sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.CustomerID,
		&payment.PaymentMethod,
		&payment.Amount,
		&payment.CurrencyCode,
		&transactionID,
		&gatewayResponse,
		&authCode,
		&payment.RefundAmount,
		&failureReason,
		&processedDate,
		&authorizedDate,
		&capturedDate,
		&refundedDate,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, apperrors.NewInternalError("failed to find payment by ID", err)
	}

	// Handle nullable fields
	if transactionID.Valid {
		payment.TransactionID = transactionID.String
	}
	if gatewayResponse.Valid {
		payment.GatewayResponse = gatewayResponse.String
	}
	if authCode.Valid {
		payment.AuthorizationCode = authCode.String
	}
	if failureReason.Valid {
		payment.FailureReason = failureReason.String
	}
	if processedDate.Valid {
		payment.ProcessedDate = &processedDate.Time
	}
	if authorizedDate.Valid {
		payment.AuthorizedDate = &authorizedDate.Time
	}
	if capturedDate.Valid {
		payment.CapturedDate = &capturedDate.Time
	}
	if refundedDate.Valid {
		payment.RefundedDate = &refundedDate.Time
	}

	return payment, nil
}

// FindByOrderID finds payments by order ID
func (r *PostgresPaymentRepository) FindByOrderID(ctx context.Context, orderID int64) ([]*domain.Payment, error) {
	query := `
		SELECT payment_id, order_id, customer_id, type, amount, currency_code,
			   transaction_id, gateway_response_code, authorization_code, refund_amount,
			   failure_reason, processed_date, authorized_date, captured_date, refunded_date,
			   date_created, date_updated
		FROM blc_order_payment
		WHERE order_id = $1
		ORDER BY date_created DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to find payments by order", err)
	}
	defer rows.Close()

	return r.scanPayments(rows)
}

// FindByCustomerID finds payments by customer ID
func (r *PostgresPaymentRepository) FindByCustomerID(ctx context.Context, customerID int64, filter *domain.PaymentFilter) ([]*domain.Payment, int64, error) {
	query := `
		SELECT payment_id, order_id, customer_id, type, amount, currency_code,
			   transaction_id, gateway_response_code, authorization_code, refund_amount,
			   failure_reason, processed_date, authorized_date, captured_date, refunded_date,
			   date_created, date_updated
		FROM blc_order_payment
		WHERE customer_id = $1
	`

	args := []interface{}{customerID}
	argIndex := 2

	// Add filters
	if filter != nil && filter.PaymentMethod != "" {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, filter.PaymentMethod)
		argIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM blc_order_payment WHERE customer_id = $1"
	countArgs := []interface{}{customerID}
	if filter != nil && filter.PaymentMethod != "" {
		countQuery += " AND type = $2"
		countArgs = append(countArgs, filter.PaymentMethod)
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("failed to count payments", err)
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

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("failed to find payments by customer", err)
	}
	defer rows.Close()

	payments, err := r.scanPayments(rows)
	return payments, total, err
}

// FindByTransactionID finds a payment by transaction ID
func (r *PostgresPaymentRepository) FindByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error) {
	query := `
		SELECT payment_id, order_id, customer_id, type, amount, currency_code,
			   transaction_id, gateway_response_code, authorization_code, refund_amount,
			   failure_reason, processed_date, authorized_date, captured_date, refunded_date,
			   date_created, date_updated
		FROM blc_order_payment
		WHERE transaction_id = $1
	`

	payment := &domain.Payment{}
	var (
		processedDate   sql.NullTime
		authorizedDate  sql.NullTime
		capturedDate    sql.NullTime
		refundedDate    sql.NullTime
		txnID           sql.NullString
		gatewayResponse sql.NullString
		authCode        sql.NullString
		failureReason   sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, transactionID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.CustomerID,
		&payment.PaymentMethod,
		&payment.Amount,
		&payment.CurrencyCode,
		&txnID,
		&gatewayResponse,
		&authCode,
		&payment.RefundAmount,
		&failureReason,
		&processedDate,
		&authorizedDate,
		&capturedDate,
		&refundedDate,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, apperrors.NewInternalError("failed to find payment by transaction ID", err)
	}

	// Handle nullable fields
	if txnID.Valid {
		payment.TransactionID = txnID.String
	}
	if gatewayResponse.Valid {
		payment.GatewayResponse = gatewayResponse.String
	}
	if authCode.Valid {
		payment.AuthorizationCode = authCode.String
	}
	if failureReason.Valid {
		payment.FailureReason = failureReason.String
	}
	if processedDate.Valid {
		payment.ProcessedDate = &processedDate.Time
	}
	if authorizedDate.Valid {
		payment.AuthorizedDate = &authorizedDate.Time
	}
	if capturedDate.Valid {
		payment.CapturedDate = &capturedDate.Time
	}
	if refundedDate.Valid {
		payment.RefundedDate = &refundedDate.Time
	}

	return payment, nil
}

// FindAll finds all payments
func (r *PostgresPaymentRepository) FindAll(ctx context.Context, filter *domain.PaymentFilter) ([]*domain.Payment, int64, error) {
	query := `
		SELECT payment_id, order_id, customer_id, type, amount, currency_code,
			   transaction_id, gateway_response_code, authorization_code, refund_amount,
			   failure_reason, processed_date, authorized_date, captured_date, refunded_date,
			   date_created, date_updated
		FROM blc_order_payment
		WHERE 1=1
	`

	args := make([]interface{}, 0)
	argIndex := 1

	// Add filters
	if filter != nil {
		if filter.PaymentMethod != "" {
			query += fmt.Sprintf(" AND type = $%d", argIndex)
			args = append(args, filter.PaymentMethod)
			argIndex++
		}
		if filter.CustomerID > 0 {
			query += fmt.Sprintf(" AND customer_id = $%d", argIndex)
			args = append(args, filter.CustomerID)
			argIndex++
		}
		if filter.OrderID > 0 {
			query += fmt.Sprintf(" AND order_id = $%d", argIndex)
			args = append(args, filter.OrderID)
			argIndex++
		}
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM blc_order_payment WHERE 1=1"
	countArgs := make([]interface{}, 0)
	countArgIndex := 1
	if filter != nil {
		if filter.PaymentMethod != "" {
			countQuery += fmt.Sprintf(" AND type = $%d", countArgIndex)
			countArgs = append(countArgs, filter.PaymentMethod)
			countArgIndex++
		}
		if filter.CustomerID > 0 {
			countQuery += fmt.Sprintf(" AND customer_id = $%d", countArgIndex)
			countArgs = append(countArgs, filter.CustomerID)
			countArgIndex++
		}
		if filter.OrderID > 0 {
			countQuery += fmt.Sprintf(" AND order_id = $%d", countArgIndex)
			countArgs = append(countArgs, filter.OrderID)
		}
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("failed to count payments", err)
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

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("failed to find all payments", err)
	}
	defer rows.Close()

	payments, err := r.scanPayments(rows)
	return payments, total, err
}

// scanPayments scans payment rows
func (r *PostgresPaymentRepository) scanPayments(rows *sql.Rows) ([]*domain.Payment, error) {
	payments := make([]*domain.Payment, 0)

	for rows.Next() {
		payment := &domain.Payment{}
		var (
			processedDate   sql.NullTime
			authorizedDate  sql.NullTime
			capturedDate    sql.NullTime
			refundedDate    sql.NullTime
			transactionID   sql.NullString
			gatewayResponse sql.NullString
			authCode        sql.NullString
			failureReason   sql.NullString
		)

		err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.CustomerID,
			&payment.PaymentMethod,
			&payment.Amount,
			&payment.CurrencyCode,
			&transactionID,
			&gatewayResponse,
			&authCode,
			&payment.RefundAmount,
			&failureReason,
			&processedDate,
			&authorizedDate,
			&capturedDate,
			&refundedDate,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to scan payment", err)
		}

		// Handle nullable fields
		if transactionID.Valid {
			payment.TransactionID = transactionID.String
		}
		if gatewayResponse.Valid {
			payment.GatewayResponse = gatewayResponse.String
		}
		if authCode.Valid {
			payment.AuthorizationCode = authCode.String
		}
		if failureReason.Valid {
			payment.FailureReason = failureReason.String
		}
		if processedDate.Valid {
			payment.ProcessedDate = &processedDate.Time
		}
		if authorizedDate.Valid {
			payment.AuthorizedDate = &authorizedDate.Time
		}
		if capturedDate.Valid {
			payment.CapturedDate = &capturedDate.Time
		}
		if refundedDate.Valid {
			payment.RefundedDate = &refundedDate.Time
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.NewInternalError("failed to iterate payments", err)
	}

	return payments, nil
}
