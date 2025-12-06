package domain

import "time"

// TieredOffer represents a tiered promotional offer (spend more, save more)
type TieredOffer struct {
	ID                 int64
	Name               string
	Description        string
	StartDate          time.Time
	EndDate            *time.Time
	IsActive           bool
	Tiers              []TierLevel
	CustomerSegmentIDs []int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// TierLevel represents a single tier in a tiered offer
type TierLevel struct {
	ID            int64
	TieredOfferID int64
	MinSpend      float64
	DiscountType  string // PERCENT, AMOUNT
	DiscountValue float64
	SortOrder     int
}

// NewTieredOffer creates a new tiered offer
func NewTieredOffer(name, description string, startDate time.Time, tiers []TierLevel) *TieredOffer {
	now := time.Now()
	return &TieredOffer{
		Name:        name,
		Description: description,
		StartDate:   startDate,
		IsActive:    true,
		Tiers:       tiers,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// GetApplicableTier returns the tier that applies for the given spend amount
func (t *TieredOffer) GetApplicableTier(spendAmount float64) *TierLevel {
	var applicableTier *TierLevel
	for i := range t.Tiers {
		tier := &t.Tiers[i]
		if spendAmount >= tier.MinSpend {
			if applicableTier == nil || tier.MinSpend > applicableTier.MinSpend {
				applicableTier = tier
			}
		}
	}
	return applicableTier
}

// CalculateDiscount calculates the discount for a given spend amount
func (t *TieredOffer) CalculateDiscount(spendAmount float64) float64 {
	tier := t.GetApplicableTier(spendAmount)
	if tier == nil {
		return 0
	}

	if tier.DiscountType == "PERCENT" {
		return spendAmount * (tier.DiscountValue / 100)
	}
	return tier.DiscountValue
}

// IsActive checks if the offer is currently active
func (t *TieredOffer) IsCurrentlyActive() bool {
	if !t.IsActive {
		return false
	}
	now := time.Now()
	if now.Before(t.StartDate) {
		return false
	}
	if t.EndDate != nil && now.After(*t.EndDate) {
		return false
	}
	return true
}
