package domain

import (
	"time"

	"github.com/google/uuid"
)

// PaymentToken represents a tokenized payment method
type PaymentToken struct {
	ID             string
	CustomerID     string
	TokenType      TokenType
	Token          string // Gateway-specific token
	GatewayName    string
	Last4Digits    *string
	CardBrand      *string
	ExpiryMonth    *int
	ExpiryYear     *int
	BillingAddress *Address
	IsDefault      bool
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TokenType represents the type of payment token
type TokenType string

const (
	TokenTypeCreditCard   TokenType = "CREDIT_CARD"
	TokenTypeDebitCard    TokenType = "DEBIT_CARD"
	TokenTypeBankAccount  TokenType = "BANK_ACCOUNT"
	TokenTypeDigitalWallet TokenType = "DIGITAL_WALLET"
)

// NewPaymentToken creates a new payment token
func NewPaymentToken(customerID, token, gatewayName string, tokenType TokenType) *PaymentToken {
	now := time.Now()
	return &PaymentToken{
		ID:          uuid.New().String(),
		CustomerID:  customerID,
		Token:       token,
		GatewayName: gatewayName,
		TokenType:   tokenType,
		IsDefault:   false,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// SetAsDefault marks this token as default
func (t *PaymentToken) SetAsDefault() {
	t.IsDefault = true
	t.UpdatedAt = time.Now()
}

// Deactivate deactivates the token
func (t *PaymentToken) Deactivate() {
	t.IsActive = false
	t.UpdatedAt = time.Now()
}

// Activate activates the token
func (t *PaymentToken) Activate() {
	t.IsActive = true
	t.UpdatedAt = time.Now()
}

// IsExpired checks if a card token is expired
func (t *PaymentToken) IsExpired() bool {
	if t.ExpiryMonth == nil || t.ExpiryYear == nil {
		return false
	}

	now := time.Now()
	expiryDate := time.Date(*t.ExpiryYear, time.Month(*t.ExpiryMonth), 1, 0, 0, 0, 0, time.UTC)
	// Card expires at end of month
	expiryDate = expiryDate.AddDate(0, 1, 0).Add(-time.Second)

	return now.After(expiryDate)
}
