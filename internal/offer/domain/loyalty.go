package domain

import (
	"errors"
	"time"
)

// LoyaltyAccount represents a customer's loyalty points account
type LoyaltyAccount struct {
	ID             int64
	CustomerID     int64
	CurrentPoints  int64
	LifetimePoints int64
	Tier           string // BRONZE, SILVER, GOLD, PLATINUM
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// LoyaltyTransaction represents a points transaction
type LoyaltyTransaction struct {
	ID                int64
	LoyaltyAccountID  int64
	Points            int64  // Can be positive (earn) or negative (redeem)
	TransactionType   string // EARN, REDEEM, EXPIRE, ADJUSTMENT
	OrderID           *int64
	Description       string
	ExpiresAt         *time.Time
	CreatedAt         time.Time
}

// LoyaltyRule represents rules for earning loyalty points
type LoyaltyRule struct {
	ID                int64
	Name              string
	Description       string
	PointsPerDollar   float64
	MinPurchaseAmount float64
	ProductIDs        []int64
	CategoryIDs       []int64
	CustomerTier      *string
	IsActive          bool
	StartDate         time.Time
	EndDate           *time.Time
	Priority          int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NewLoyaltyAccount creates a new loyalty account
func NewLoyaltyAccount(customerID int64) *LoyaltyAccount {
	now := time.Now()
	return &LoyaltyAccount{
		CustomerID:     customerID,
		CurrentPoints:  0,
		LifetimePoints: 0,
		Tier:           "BRONZE",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// AwardPoints adds points to the account
func (l *LoyaltyAccount) AwardPoints(points int64) error {
	if points <= 0 {
		return errors.New("points must be greater than 0")
	}
	l.CurrentPoints += points
	l.LifetimePoints += points
	l.UpdatedAt = time.Now()
	l.updateTier()
	return nil
}

// RedeemPoints redeems points from the account
func (l *LoyaltyAccount) RedeemPoints(points int64) error {
	if points <= 0 {
		return errors.New("points must be greater than 0")
	}
	if points > l.CurrentPoints {
		return errors.New("insufficient points")
	}
	l.CurrentPoints -= points
	l.UpdatedAt = time.Now()
	return nil
}

// updateTier updates the customer's tier based on lifetime points
func (l *LoyaltyAccount) updateTier() {
	switch {
	case l.LifetimePoints >= 10000:
		l.Tier = "PLATINUM"
	case l.LifetimePoints >= 5000:
		l.Tier = "GOLD"
	case l.LifetimePoints >= 1000:
		l.Tier = "SILVER"
	default:
		l.Tier = "BRONZE"
	}
}

// NewLoyaltyRule creates a new loyalty rule
func NewLoyaltyRule(name, description string, pointsPerDollar, minPurchaseAmount float64, startDate time.Time) *LoyaltyRule {
	now := time.Now()
	return &LoyaltyRule{
		Name:              name,
		Description:       description,
		PointsPerDollar:   pointsPerDollar,
		MinPurchaseAmount: minPurchaseAmount,
		IsActive:          true,
		StartDate:         startDate,
		Priority:          0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// CalculatePoints calculates points earned for a purchase amount
func (r *LoyaltyRule) CalculatePoints(amount float64) int64 {
	if amount < r.MinPurchaseAmount {
		return 0
	}
	return int64(amount * r.PointsPerDollar)
}

// IsCurrentlyActive checks if the rule is currently active
func (r *LoyaltyRule) IsCurrentlyActive() bool {
	if !r.IsActive {
		return false
	}
	now := time.Now()
	if now.Before(r.StartDate) {
		return false
	}
	if r.EndDate != nil && now.After(*r.EndDate) {
		return false
	}
	return true
}

// ConvertPointsToDiscount converts loyalty points to a discount amount
func ConvertPointsToDiscount(points int64, conversionRate float64) float64 {
	return float64(points) * conversionRate
}
