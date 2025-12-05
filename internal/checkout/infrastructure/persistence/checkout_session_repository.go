package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
	"github.com/qhato/ecommerce/internal/checkout/domain"
)

// PostgresCheckoutSessionRepository implements domain.CheckoutSessionRepository
type PostgresCheckoutSessionRepository struct {
	db *sql.DB
}

// NewPostgresCheckoutSessionRepository creates a new repository
func NewPostgresCheckoutSessionRepository(db *sql.DB) *PostgresCheckoutSessionRepository {
	return &PostgresCheckoutSessionRepository{db: db}
}

// Create creates a new checkout session
func (r *PostgresCheckoutSessionRepository) Create(ctx context.Context, session *domain.CheckoutSession) error {
	sessionDataJSON, _ := json.Marshal(session.SessionData)

	query := `
		INSERT INTO blc_checkout_session (
			id, order_id, customer_id, email, is_guest_checkout, state,
			current_step, completed_steps, shipping_address_id, billing_address_id,
			shipping_method_id, payment_method_id, subtotal, shipping_cost,
			tax_amount, discount_amount, total_amount, coupon_codes,
			customer_notes, session_data, expires_at, last_activity_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)`

	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.OrderID, session.CustomerID, session.Email, session.IsGuestCheckout,
		session.State, session.CurrentStep, pq.Array(session.CompletedSteps),
		session.ShippingAddressID, session.BillingAddressID, session.ShippingMethodID,
		session.PaymentMethodID, session.Subtotal, session.ShippingCost, session.TaxAmount,
		session.DiscountAmount, session.TotalAmount, pq.Array(session.CouponCodes),
		session.CustomerNotes, sessionDataJSON, session.ExpiresAt, session.LastActivityAt,
		session.CreatedAt, session.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create checkout session: %w", err)
	}

	return nil
}

// Update updates an existing checkout session
func (r *PostgresCheckoutSessionRepository) Update(ctx context.Context, session *domain.CheckoutSession) error {
	sessionDataJSON, _ := json.Marshal(session.SessionData)

	query := `
		UPDATE blc_checkout_session SET
			customer_id = $1, email = $2, state = $3, current_step = $4,
			completed_steps = $5, shipping_address_id = $6, billing_address_id = $7,
			shipping_method_id = $8, payment_method_id = $9, subtotal = $10,
			shipping_cost = $11, tax_amount = $12, discount_amount = $13,
			total_amount = $14, coupon_codes = $15, customer_notes = $16,
			session_data = $17, expires_at = $18, last_activity_at = $19,
			updated_at = $20, submitted_at = $21, confirmed_at = $22
		WHERE id = $23`

	_, err := r.db.ExecContext(ctx, query,
		session.CustomerID, session.Email, session.State, session.CurrentStep,
		pq.Array(session.CompletedSteps), session.ShippingAddressID, session.BillingAddressID,
		session.ShippingMethodID, session.PaymentMethodID, session.Subtotal,
		session.ShippingCost, session.TaxAmount, session.DiscountAmount,
		session.TotalAmount, pq.Array(session.CouponCodes), session.CustomerNotes,
		sessionDataJSON, session.ExpiresAt, session.LastActivityAt,
		session.UpdatedAt, session.SubmittedAt, session.ConfirmedAt, session.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update checkout session: %w", err)
	}

	return nil
}

// FindByID finds a checkout session by ID
func (r *PostgresCheckoutSessionRepository) FindByID(ctx context.Context, id string) (*domain.CheckoutSession, error) {
	query := `
		SELECT id, order_id, customer_id, email, is_guest_checkout, state,
			   current_step, completed_steps, shipping_address_id, billing_address_id,
			   shipping_method_id, payment_method_id, subtotal, shipping_cost,
			   tax_amount, discount_amount, total_amount, coupon_codes,
			   customer_notes, session_data, expires_at, last_activity_at,
			   created_at, updated_at, submitted_at, confirmed_at
		FROM blc_checkout_session WHERE id = $1`

	session := &domain.CheckoutSession{}
	var sessionDataJSON []byte
	var completedSteps, couponCodes pq.StringArray

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID, &session.OrderID, &session.CustomerID, &session.Email,
		&session.IsGuestCheckout, &session.State, &session.CurrentStep,
		&completedSteps, &session.ShippingAddressID, &session.BillingAddressID,
		&session.ShippingMethodID, &session.PaymentMethodID, &session.Subtotal,
		&session.ShippingCost, &session.TaxAmount, &session.DiscountAmount,
		&session.TotalAmount, &couponCodes, &session.CustomerNotes,
		&sessionDataJSON, &session.ExpiresAt, &session.LastActivityAt,
		&session.CreatedAt, &session.UpdatedAt, &session.SubmittedAt, &session.ConfirmedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find checkout session: %w", err)
	}

	session.CompletedSteps = completedSteps
	session.CouponCodes = couponCodes
	json.Unmarshal(sessionDataJSON, &session.SessionData)

	return session, nil
}

// FindByOrderID finds a checkout session by order ID
func (r *PostgresCheckoutSessionRepository) FindByOrderID(ctx context.Context, orderID int64) (*domain.CheckoutSession, error) {
	query := `
		SELECT id, order_id, customer_id, email, is_guest_checkout, state,
			   current_step, completed_steps, shipping_address_id, billing_address_id,
			   shipping_method_id, payment_method_id, subtotal, shipping_cost,
			   tax_amount, discount_amount, total_amount, coupon_codes,
			   customer_notes, session_data, expires_at, last_activity_at,
			   created_at, updated_at, submitted_at, confirmed_at
		FROM blc_checkout_session WHERE order_id = $1 ORDER BY created_at DESC LIMIT 1`

	session := &domain.CheckoutSession{}
	var sessionDataJSON []byte
	var completedSteps, couponCodes pq.StringArray

	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&session.ID, &session.OrderID, &session.CustomerID, &session.Email,
		&session.IsGuestCheckout, &session.State, &session.CurrentStep,
		&completedSteps, &session.ShippingAddressID, &session.BillingAddressID,
		&session.ShippingMethodID, &session.PaymentMethodID, &session.Subtotal,
		&session.ShippingCost, &session.TaxAmount, &session.DiscountAmount,
		&session.TotalAmount, &couponCodes, &session.CustomerNotes,
		&sessionDataJSON, &session.ExpiresAt, &session.LastActivityAt,
		&session.CreatedAt, &session.UpdatedAt, &session.SubmittedAt, &session.ConfirmedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find checkout session: %w", err)
	}

	session.CompletedSteps = completedSteps
	session.CouponCodes = couponCodes
	json.Unmarshal(sessionDataJSON, &session.SessionData)

	return session, nil
}

// FindByCustomerID finds checkout sessions by customer ID
func (r *PostgresCheckoutSessionRepository) FindByCustomerID(ctx context.Context, customerID string, activeOnly bool) ([]*domain.CheckoutSession, error) {
	query := `
		SELECT id, order_id, customer_id, email, is_guest_checkout, state,
			   current_step, completed_steps, shipping_address_id, billing_address_id,
			   shipping_method_id, payment_method_id, subtotal, shipping_cost,
			   tax_amount, discount_amount, total_amount, coupon_codes,
			   customer_notes, session_data, expires_at, last_activity_at,
			   created_at, updated_at, submitted_at, confirmed_at
		FROM blc_checkout_session WHERE customer_id = $1`

	if activeOnly {
		query += " AND state NOT IN ('CONFIRMED', 'CANCELLED', 'EXPIRED')"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query checkout sessions: %w", err)
	}
	defer rows.Close()

	sessions := make([]*domain.CheckoutSession, 0)
	for rows.Next() {
		session := &domain.CheckoutSession{}
		var sessionDataJSON []byte
		var completedSteps, couponCodes pq.StringArray

		err := rows.Scan(
			&session.ID, &session.OrderID, &session.CustomerID, &session.Email,
			&session.IsGuestCheckout, &session.State, &session.CurrentStep,
			&completedSteps, &session.ShippingAddressID, &session.BillingAddressID,
			&session.ShippingMethodID, &session.PaymentMethodID, &session.Subtotal,
			&session.ShippingCost, &session.TaxAmount, &session.DiscountAmount,
			&session.TotalAmount, &couponCodes, &session.CustomerNotes,
			&sessionDataJSON, &session.ExpiresAt, &session.LastActivityAt,
			&session.CreatedAt, &session.UpdatedAt, &session.SubmittedAt, &session.ConfirmedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan checkout session: %w", err)
		}

		session.CompletedSteps = completedSteps
		session.CouponCodes = couponCodes
		json.Unmarshal(sessionDataJSON, &session.SessionData)

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// FindActiveByEmail finds active checkout sessions by email
func (r *PostgresCheckoutSessionRepository) FindActiveByEmail(ctx context.Context, email string) ([]*domain.CheckoutSession, error) {
	query := `
		SELECT id, order_id, customer_id, email, is_guest_checkout, state,
			   current_step, completed_steps, shipping_address_id, billing_address_id,
			   shipping_method_id, payment_method_id, subtotal, shipping_cost,
			   tax_amount, discount_amount, total_amount, coupon_codes,
			   customer_notes, session_data, expires_at, last_activity_at,
			   created_at, updated_at, submitted_at, confirmed_at
		FROM blc_checkout_session
		WHERE email = $1 AND state NOT IN ('CONFIRMED', 'CANCELLED', 'EXPIRED')
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to query checkout sessions: %w", err)
	}
	defer rows.Close()

	sessions := make([]*domain.CheckoutSession, 0)
	for rows.Next() {
		session := &domain.CheckoutSession{}
		var sessionDataJSON []byte
		var completedSteps, couponCodes pq.StringArray

		err := rows.Scan(
			&session.ID, &session.OrderID, &session.CustomerID, &session.Email,
			&session.IsGuestCheckout, &session.State, &session.CurrentStep,
			&completedSteps, &session.ShippingAddressID, &session.BillingAddressID,
			&session.ShippingMethodID, &session.PaymentMethodID, &session.Subtotal,
			&session.ShippingCost, &session.TaxAmount, &session.DiscountAmount,
			&session.TotalAmount, &couponCodes, &session.CustomerNotes,
			&sessionDataJSON, &session.ExpiresAt, &session.LastActivityAt,
			&session.CreatedAt, &session.UpdatedAt, &session.SubmittedAt, &session.ConfirmedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan checkout session: %w", err)
		}

		session.CompletedSteps = completedSteps
		session.CouponCodes = couponCodes
		json.Unmarshal(sessionDataJSON, &session.SessionData)

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// FindExpiredSessions finds expired sessions
func (r *PostgresCheckoutSessionRepository) FindExpiredSessions(ctx context.Context, limit int) ([]*domain.CheckoutSession, error) {
	query := `
		SELECT id, order_id, customer_id, email, is_guest_checkout, state,
			   current_step, completed_steps, shipping_address_id, billing_address_id,
			   shipping_method_id, payment_method_id, subtotal, shipping_cost,
			   tax_amount, discount_amount, total_amount, coupon_codes,
			   customer_notes, session_data, expires_at, last_activity_at,
			   created_at, updated_at, submitted_at, confirmed_at
		FROM blc_checkout_session
		WHERE expires_at < NOW() AND state NOT IN ('CONFIRMED', 'CANCELLED', 'EXPIRED')
		ORDER BY expires_at ASC
		LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query expired sessions: %w", err)
	}
	defer rows.Close()

	sessions := make([]*domain.CheckoutSession, 0)
	for rows.Next() {
		session := &domain.CheckoutSession{}
		var sessionDataJSON []byte
		var completedSteps, couponCodes pq.StringArray

		err := rows.Scan(
			&session.ID, &session.OrderID, &session.CustomerID, &session.Email,
			&session.IsGuestCheckout, &session.State, &session.CurrentStep,
			&completedSteps, &session.ShippingAddressID, &session.BillingAddressID,
			&session.ShippingMethodID, &session.PaymentMethodID, &session.Subtotal,
			&session.ShippingCost, &session.TaxAmount, &session.DiscountAmount,
			&session.TotalAmount, &couponCodes, &session.CustomerNotes,
			&sessionDataJSON, &session.ExpiresAt, &session.LastActivityAt,
			&session.CreatedAt, &session.UpdatedAt, &session.SubmittedAt, &session.ConfirmedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan checkout session: %w", err)
		}

		session.CompletedSteps = completedSteps
		session.CouponCodes = couponCodes
		json.Unmarshal(sessionDataJSON, &session.SessionData)

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// Delete deletes a checkout session
func (r *PostgresCheckoutSessionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM blc_checkout_session WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete checkout session: %w", err)
	}

	return nil
}

// ExistsByOrderID checks if a checkout session exists for an order
func (r *PostgresCheckoutSessionRepository) ExistsByOrderID(ctx context.Context, orderID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blc_checkout_session WHERE order_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check session existence: %w", err)
	}

	return exists, nil
}
