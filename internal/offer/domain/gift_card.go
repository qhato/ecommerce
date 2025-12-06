package domain

import (
	"errors"
	"time"
)

// GiftCard represents a gift card
type GiftCard struct {
	ID             int64
	Code           string
	InitialAmount  float64
	CurrentBalance float64
	PurchasedBy    *int64
	RecipientEmail string
	Message        string
	IsActive       bool
	ExpiresAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewGiftCard creates a new gift card
func NewGiftCard(code string, amount float64, purchasedBy *int64, recipientEmail, message string, expiresAt *time.Time) (*GiftCard, error) {
	if amount <= 0 {
		return nil, errors.New("gift card amount must be greater than 0")
	}
	if code == "" {
		return nil, errors.New("gift card code is required")
	}

	now := time.Now()
	return &GiftCard{
		Code:           code,
		InitialAmount:  amount,
		CurrentBalance: amount,
		PurchasedBy:    purchasedBy,
		RecipientEmail: recipientEmail,
		Message:        message,
		IsActive:       true,
		ExpiresAt:      expiresAt,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Redeem redeems an amount from the gift card
func (g *GiftCard) Redeem(amount float64) error {
	if !g.IsActive {
		return errors.New("gift card is not active")
	}
	if g.ExpiresAt != nil && time.Now().After(*g.ExpiresAt) {
		return errors.New("gift card has expired")
	}
	if amount > g.CurrentBalance {
		return errors.New("insufficient balance on gift card")
	}
	if amount <= 0 {
		return errors.New("redemption amount must be greater than 0")
	}

	g.CurrentBalance -= amount
	g.UpdatedAt = time.Now()
	return nil
}

// Deactivate deactivates the gift card
func (g *GiftCard) Deactivate() {
	g.IsActive = false
	g.UpdatedAt = time.Now()
}

// IsExpired checks if the gift card is expired
func (g *GiftCard) IsExpired() bool {
	return g.ExpiresAt != nil && time.Now().After(*g.ExpiresAt)
}

// IsValid checks if the gift card is valid for use
func (g *GiftCard) IsValid() bool {
	return g.IsActive && !g.IsExpired() && g.CurrentBalance > 0
}
