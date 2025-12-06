package domain

import "time"

// BundleOffer represents a bundle promotional offer (buy specific products together for a discount)
type BundleOffer struct {
	ID                 int64
	Name               string
	Description        string
	StartDate          time.Time
	EndDate            *time.Time
	IsActive           bool
	DiscountType       string // PERCENT, AMOUNT, FIX_PRICE
	DiscountValue      float64
	BundleItems        []BundleItem
	CustomerSegmentIDs []int64
	MaxUsesPerCustomer *int64
	TotalUses          int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// BundleItem represents a product in a bundle
type BundleItem struct {
	ID            int64
	BundleOfferID int64
	ProductID     int64
	Quantity      int
	SortOrder     int
}

// NewBundleOffer creates a new bundle offer
func NewBundleOffer(name, description string, startDate time.Time, discountType string, discountValue float64, items []BundleItem) *BundleOffer {
	now := time.Now()
	return &BundleOffer{
		Name:          name,
		Description:   description,
		StartDate:     startDate,
		IsActive:      true,
		DiscountType:  discountType,
		DiscountValue: discountValue,
		BundleItems:   items,
		TotalUses:     0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// ValidateBundle checks if the given products match the bundle requirements
func (b *BundleOffer) ValidateBundle(productQuantities map[int64]int) bool {
	for _, item := range b.BundleItems {
		qty, exists := productQuantities[item.ProductID]
		if !exists || qty < item.Quantity {
			return false
		}
	}
	return true
}

// CalculateDiscount calculates the discount for the bundle
func (b *BundleOffer) CalculateDiscount(bundleTotal float64) float64 {
	switch b.DiscountType {
	case "PERCENT":
		return bundleTotal * (b.DiscountValue / 100)
	case "AMOUNT":
		return b.DiscountValue
	case "FIX_PRICE":
		discount := bundleTotal - b.DiscountValue
		if discount < 0 {
			return 0
		}
		return discount
	default:
		return 0
	}
}

// IsCurrentlyActive checks if the bundle offer is currently active
func (b *BundleOffer) IsCurrentlyActive() bool {
	if !b.IsActive {
		return false
	}
	now := time.Now()
	if now.Before(b.StartDate) {
		return false
	}
	if b.EndDate != nil && now.After(*b.EndDate) {
		return false
	}
	return true
}
