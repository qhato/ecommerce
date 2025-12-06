package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/payment/domain"
)

type PostgresPaymentTokenRepository struct {
	db *sql.DB
}

func NewPostgresPaymentTokenRepository(db *sql.DB) *PostgresPaymentTokenRepository {
	return &PostgresPaymentTokenRepository{db: db}
}

func (r *PostgresPaymentTokenRepository) Create(ctx context.Context, token *domain.PaymentToken) error {
	query := `
		INSERT INTO blc_payment_token (id, customer_id, token_type, token, gateway_name,
			last_4_digits, card_brand, expiry_month, expiry_year, is_default, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.CustomerID, token.TokenType, token.Token, token.GatewayName,
		token.Last4Digits, token.CardBrand, token.ExpiryMonth, token.ExpiryYear,
		token.IsDefault, token.IsActive, token.CreatedAt, token.UpdatedAt,
	)
	return err
}

func (r *PostgresPaymentTokenRepository) Update(ctx context.Context, token *domain.PaymentToken) error {
	query := `
		UPDATE blc_payment_token
		SET token = $1, last_4_digits = $2, card_brand = $3, expiry_month = $4, expiry_year = $5,
		    is_default = $6, is_active = $7, updated_at = $8
		WHERE id = $9`

	_, err := r.db.ExecContext(ctx, query,
		token.Token, token.Last4Digits, token.CardBrand, token.ExpiryMonth, token.ExpiryYear,
		token.IsDefault, token.IsActive, token.UpdatedAt, token.ID,
	)
	return err
}

func (r *PostgresPaymentTokenRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM blc_payment_token WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresPaymentTokenRepository) FindByID(ctx context.Context, id string) (*domain.PaymentToken, error) {
	query := `
		SELECT id, customer_id, token_type, token, gateway_name, last_4_digits, card_brand,
		       expiry_month, expiry_year, is_default, is_active, created_at, updated_at
		FROM blc_payment_token
		WHERE id = $1`

	return r.scanToken(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresPaymentTokenRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*domain.PaymentToken, error) {
	query := `
		SELECT id, customer_id, token_type, token, gateway_name, last_4_digits, card_brand,
		       expiry_month, expiry_year, is_default, is_active, created_at, updated_at
		FROM blc_payment_token
		WHERE customer_id = $1
		ORDER BY is_default DESC, created_at DESC`

	return r.queryTokens(ctx, query, customerID)
}

func (r *PostgresPaymentTokenRepository) FindDefaultByCustomerID(ctx context.Context, customerID string) (*domain.PaymentToken, error) {
	query := `
		SELECT id, customer_id, token_type, token, gateway_name, last_4_digits, card_brand,
		       expiry_month, expiry_year, is_default, is_active, created_at, updated_at
		FROM blc_payment_token
		WHERE customer_id = $1 AND is_default = true AND is_active = true
		LIMIT 1`

	return r.scanToken(r.db.QueryRowContext(ctx, query, customerID))
}

func (r *PostgresPaymentTokenRepository) FindActiveByCustomerID(ctx context.Context, customerID string) ([]*domain.PaymentToken, error) {
	query := `
		SELECT id, customer_id, token_type, token, gateway_name, last_4_digits, card_brand,
		       expiry_month, expiry_year, is_default, is_active, created_at, updated_at
		FROM blc_payment_token
		WHERE customer_id = $1 AND is_active = true
		ORDER BY is_default DESC, created_at DESC`

	return r.queryTokens(ctx, query, customerID)
}

func (r *PostgresPaymentTokenRepository) scanToken(row interface {
	Scan(dest ...interface{}) error
}) (*domain.PaymentToken, error) {
	token := &domain.PaymentToken{}
	err := row.Scan(
		&token.ID, &token.CustomerID, &token.TokenType, &token.Token, &token.GatewayName,
		&token.Last4Digits, &token.CardBrand, &token.ExpiryMonth, &token.ExpiryYear,
		&token.IsDefault, &token.IsActive, &token.CreatedAt, &token.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *PostgresPaymentTokenRepository) queryTokens(ctx context.Context, query string, args ...interface{}) ([]*domain.PaymentToken, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*domain.PaymentToken
	for rows.Next() {
		token := &domain.PaymentToken{}
		if err := rows.Scan(
			&token.ID, &token.CustomerID, &token.TokenType, &token.Token, &token.GatewayName,
			&token.Last4Digits, &token.CardBrand, &token.ExpiryMonth, &token.ExpiryYear,
			&token.IsDefault, &token.IsActive, &token.CreatedAt, &token.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}
