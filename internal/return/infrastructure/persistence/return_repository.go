package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/qhato/ecommerce/internal/return/domain"
)

type PostgresReturnRepository struct {
	db *sql.DB
}

func NewPostgresReturnRepository(db *sql.DB) *PostgresReturnRepository {
	return &PostgresReturnRepository{db: db}
}

func (r *PostgresReturnRepository) Create(ctx context.Context, returnReq *domain.ReturnRequest) error {
	itemsJSON, err := json.Marshal(returnReq.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	query := `INSERT INTO blc_return_request (
		rma, order_id, customer_id, status, reason, reason_details, refund_amount, refund_method,
		items, notes, tracking_number, created_at, updated_at, approved_at, received_at,
		inspected_at, refunded_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		returnReq.RMA, returnReq.OrderID, returnReq.CustomerID, returnReq.Status,
		returnReq.Reason, returnReq.ReasonDetails, returnReq.RefundAmount, returnReq.RefundMethod,
		itemsJSON, returnReq.Notes, returnReq.TrackingNumber, returnReq.CreatedAt, returnReq.UpdatedAt,
		returnReq.ApprovedAt, returnReq.ReceivedAt, returnReq.InspectedAt, returnReq.RefundedAt,
	).Scan(&returnReq.ID)
}

func (r *PostgresReturnRepository) Update(ctx context.Context, returnReq *domain.ReturnRequest) error {
	itemsJSON, err := json.Marshal(returnReq.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	query := `UPDATE blc_return_request SET
		status = $1, reason = $2, reason_details = $3, refund_amount = $4, refund_method = $5,
		items = $6, notes = $7, tracking_number = $8, updated_at = $9,
		approved_at = $10, received_at = $11, inspected_at = $12, refunded_at = $13
	WHERE id = $14`

	_, err = r.db.ExecContext(ctx, query,
		returnReq.Status, returnReq.Reason, returnReq.ReasonDetails, returnReq.RefundAmount,
		returnReq.RefundMethod, itemsJSON, returnReq.Notes, returnReq.TrackingNumber,
		returnReq.UpdatedAt, returnReq.ApprovedAt, returnReq.ReceivedAt, returnReq.InspectedAt,
		returnReq.RefundedAt, returnReq.ID,
	)
	return err
}

func (r *PostgresReturnRepository) FindByID(ctx context.Context, id int64) (*domain.ReturnRequest, error) {
	query := `SELECT id, rma, order_id, customer_id, status, reason, reason_details,
		refund_amount, refund_method, items, notes, tracking_number, created_at, updated_at,
		approved_at, received_at, inspected_at, refunded_at
	FROM blc_return_request WHERE id = $1`

	return r.scanReturn(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresReturnRepository) FindByRMA(ctx context.Context, rma string) (*domain.ReturnRequest, error) {
	query := `SELECT id, rma, order_id, customer_id, status, reason, reason_details,
		refund_amount, refund_method, items, notes, tracking_number, created_at, updated_at,
		approved_at, received_at, inspected_at, refunded_at
	FROM blc_return_request WHERE rma = $1`

	return r.scanReturn(r.db.QueryRowContext(ctx, query, rma))
}

func (r *PostgresReturnRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*domain.ReturnRequest, error) {
	query := `SELECT id, rma, order_id, customer_id, status, reason, reason_details,
		refund_amount, refund_method, items, notes, tracking_number, created_at, updated_at,
		approved_at, received_at, inspected_at, refunded_at
	FROM blc_return_request WHERE customer_id = $1 ORDER BY created_at DESC`

	return r.queryReturns(ctx, query, customerID)
}

func (r *PostgresReturnRepository) FindByStatus(ctx context.Context, status domain.ReturnStatus, limit int) ([]*domain.ReturnRequest, error) {
	query := `SELECT id, rma, order_id, customer_id, status, reason, reason_details,
		refund_amount, refund_method, items, notes, tracking_number, created_at, updated_at,
		approved_at, received_at, inspected_at, refunded_at
	FROM blc_return_request WHERE status = $1 ORDER BY created_at DESC LIMIT $2`

	return r.queryReturns(ctx, query, status, limit)
}

func (r *PostgresReturnRepository) FindByOrderID(ctx context.Context, orderID int64) ([]*domain.ReturnRequest, error) {
	query := `SELECT id, rma, order_id, customer_id, status, reason, reason_details,
		refund_amount, refund_method, items, notes, tracking_number, created_at, updated_at,
		approved_at, received_at, inspected_at, refunded_at
	FROM blc_return_request WHERE order_id = $1 ORDER BY created_at DESC`

	return r.queryReturns(ctx, query, orderID)
}

func (r *PostgresReturnRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_return_request WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresReturnRepository) GetReturnItems(ctx context.Context, returnID int64) ([]domain.ReturnItem, error) {
	returnReq, err := r.FindByID(ctx, returnID)
	if err != nil {
		return nil, err
	}
	if returnReq == nil {
		return nil, domain.ErrReturnNotFound
	}
	return returnReq.Items, nil
}

func (r *PostgresReturnRepository) scanReturn(row interface {
	Scan(dest ...interface{}) error
}) (*domain.ReturnRequest, error) {
	returnReq := &domain.ReturnRequest{}
	var itemsJSON []byte

	err := row.Scan(
		&returnReq.ID, &returnReq.RMA, &returnReq.OrderID, &returnReq.CustomerID,
		&returnReq.Status, &returnReq.Reason, &returnReq.ReasonDetails, &returnReq.RefundAmount,
		&returnReq.RefundMethod, &itemsJSON, &returnReq.Notes, &returnReq.TrackingNumber,
		&returnReq.CreatedAt, &returnReq.UpdatedAt, &returnReq.ApprovedAt,
		&returnReq.ReceivedAt, &returnReq.InspectedAt, &returnReq.RefundedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(itemsJSON, &returnReq.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return returnReq, nil
}

func (r *PostgresReturnRepository) queryReturns(ctx context.Context, query string, args ...interface{}) ([]*domain.ReturnRequest, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	returns := make([]*domain.ReturnRequest, 0)
	for rows.Next() {
		returnReq := &domain.ReturnRequest{}
		var itemsJSON []byte

		if err := rows.Scan(
			&returnReq.ID, &returnReq.RMA, &returnReq.OrderID, &returnReq.CustomerID,
			&returnReq.Status, &returnReq.Reason, &returnReq.ReasonDetails, &returnReq.RefundAmount,
			&returnReq.RefundMethod, &itemsJSON, &returnReq.Notes, &returnReq.TrackingNumber,
			&returnReq.CreatedAt, &returnReq.UpdatedAt, &returnReq.ApprovedAt,
			&returnReq.ReceivedAt, &returnReq.InspectedAt, &returnReq.RefundedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(itemsJSON, &returnReq.Items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %w", err)
		}

		returns = append(returns, returnReq)
	}

	return returns, nil
}
