package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// PriceListItem represents a price entry for a specific product/SKU in a price list
// Business Logic: Precios específicos por SKU con soporte para cantidad mínima
type PriceListItem struct {
	ID              int64
	PriceListID     int64
	SKUID           string
	ProductID       *string
	Price           decimal.Decimal
	CompareAtPrice  *decimal.Decimal // Original price for comparison (for showing discounts)
	MinQuantity     int              // Minimum quantity required for this price
	MaxQuantity     *int             // Maximum quantity for this price (for tiered pricing)
	IsActive        bool
	StartDate       *time.Time
	EndDate         *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewPriceListItem creates a new PriceListItem
func NewPriceListItem(priceListID int64, skuID string, price decimal.Decimal, minQuantity int) (*PriceListItem, error) {
	if priceListID == 0 {
		return nil, ErrPriceListIDRequired
	}
	if skuID == "" {
		return nil, ErrSKUIDRequired
	}
	if price.LessThan(decimal.Zero) {
		return nil, ErrPriceCannotBeNegative
	}
	if minQuantity < 0 {
		return nil, ErrMinQuantityCannotBeNegative
	}

	now := time.Now()
	return &PriceListItem{
		PriceListID: priceListID,
		SKUID:       skuID,
		Price:       price,
		MinQuantity: minQuantity,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// IsCurrentlyActive checks if the price list item is currently active
func (pli *PriceListItem) IsCurrentlyActive() bool {
	if !pli.IsActive {
		return false
	}

	now := time.Now()
	if pli.StartDate != nil && now.Before(*pli.StartDate) {
		return false
	}
	if pli.EndDate != nil && now.After(*pli.EndDate) {
		return false
	}

	return true
}

// AppliesTo checks if this price applies to the given quantity
func (pli *PriceListItem) AppliesTo(quantity int) bool {
	if quantity < pli.MinQuantity {
		return false
	}
	if pli.MaxQuantity != nil && quantity > *pli.MaxQuantity {
		return false
	}
	return true
}

// GetDiscountPercentage calculates the discount percentage if compare at price is set
func (pli *PriceListItem) GetDiscountPercentage() decimal.Decimal {
	if pli.CompareAtPrice == nil {
		return decimal.Zero
	}

	if pli.CompareAtPrice.LessThanOrEqual(pli.Price) {
		return decimal.Zero
	}

	discount := pli.CompareAtPrice.Sub(pli.Price)
	percentage := discount.Div(*pli.CompareAtPrice).Mul(decimal.NewFromInt(100))
	return percentage
}

// SetCompareAtPrice sets the compare at price for showing discounts
func (pli *PriceListItem) SetCompareAtPrice(price decimal.Decimal) {
	pli.CompareAtPrice = &price
	pli.UpdatedAt = time.Now()
}

// SetQuantityRange sets the quantity range for this price
func (pli *PriceListItem) SetQuantityRange(minQuantity int, maxQuantity *int) error {
	if minQuantity < 0 {
		return ErrMinQuantityCannotBeNegative
	}
	if maxQuantity != nil && *maxQuantity < minQuantity {
		return ErrMaxQuantityLessThanMin
	}

	pli.MinQuantity = minQuantity
	pli.MaxQuantity = maxQuantity
	pli.UpdatedAt = time.Now()
	return nil
}

// SetDateRange sets the start and end dates for the price list item
func (pli *PriceListItem) SetDateRange(startDate, endDate *time.Time) {
	pli.StartDate = startDate
	pli.EndDate = endDate
	pli.UpdatedAt = time.Now()
}

// Activate sets the price list item to active
func (pli *PriceListItem) Activate() {
	pli.IsActive = true
	pli.UpdatedAt = time.Now()
}

// Deactivate sets the price list item to inactive
func (pli *PriceListItem) Deactivate() {
	pli.IsActive = false
	pli.UpdatedAt = time.Now()
}

// UpdatePrice updates the price
func (pli *PriceListItem) UpdatePrice(price decimal.Decimal) error {
	if price.LessThan(decimal.Zero) {
		return ErrPriceCannotBeNegative
	}
	pli.Price = price
	pli.UpdatedAt = time.Now()
	return nil
}
