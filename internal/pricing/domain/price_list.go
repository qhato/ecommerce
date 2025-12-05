package domain

import (
	"time"
)

// PriceListType defines the type of price list
type PriceListType string

const (
	PriceListTypeStandard  PriceListType = "STANDARD"  // Regular pricing
	PriceListTypePromotion PriceListType = "PROMOTION" // Promotional pricing
	PriceListTypeCustomer  PriceListType = "CUSTOMER"  // Customer-specific pricing
	PriceListTypeSegment   PriceListType = "SEGMENT"   // Segment-based pricing
)

// PriceList represents a collection of prices for products
// Business Logic: Permite tener múltiples listas de precios (wholesale, retail, vip, etc.)
type PriceList struct {
	ID               int64
	Name             string
	Code             string // Unique code for the price list
	PriceListType    PriceListType
	Currency         string // ISO 4217 currency code (USD, EUR, etc.)
	Priority         int    // Higher priority lists take precedence
	IsActive         bool
	StartDate        *time.Time
	EndDate          *time.Time
	Description      string
	CustomerSegments []string // Customer segments this price list applies to
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewPriceList creates a new PriceList
func NewPriceList(name, code string, priceListType PriceListType, currency string, priority int) (*PriceList, error) {
	if name == "" {
		return nil, ErrPriceListNameRequired
	}
	if code == "" {
		return nil, ErrPriceListCodeRequired
	}
	if currency == "" {
		return nil, ErrCurrencyRequired
	}

	now := time.Now()
	return &PriceList{
		Name:             name,
		Code:             code,
		PriceListType:    priceListType,
		Currency:         currency,
		Priority:         priority,
		IsActive:         true,
		CustomerSegments: make([]string, 0),
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// IsCurrentlyActive checks if the price list is currently active
func (pl *PriceList) IsCurrentlyActive() bool {
	if !pl.IsActive {
		return false
	}

	now := time.Now()
	if pl.StartDate != nil && now.Before(*pl.StartDate) {
		return false
	}
	if pl.EndDate != nil && now.After(*pl.EndDate) {
		return false
	}

	return true
}

// Activate sets the price list to active
func (pl *PriceList) Activate() {
	pl.IsActive = true
	pl.UpdatedAt = time.Now()
}

// Deactivate sets the price list to inactive
func (pl *PriceList) Deactivate() {
	pl.IsActive = false
	pl.UpdatedAt = time.Now()
}

// SetDateRange sets the start and end dates for the price list
func (pl *PriceList) SetDateRange(startDate, endDate *time.Time) {
	pl.StartDate = startDate
	pl.EndDate = endDate
	pl.UpdatedAt = time.Now()
}

// AddCustomerSegment adds a customer segment to the price list
func (pl *PriceList) AddCustomerSegment(segment string) {
	for _, s := range pl.CustomerSegments {
		if s == segment {
			return // Already exists
		}
	}
	pl.CustomerSegments = append(pl.CustomerSegments, segment)
	pl.UpdatedAt = time.Now()
}

// RemoveCustomerSegment removes a customer segment from the price list
func (pl *PriceList) RemoveCustomerSegment(segment string) {
	for i, s := range pl.CustomerSegments {
		if s == segment {
			pl.CustomerSegments = append(pl.CustomerSegments[:i], pl.CustomerSegments[i+1:]...)
			pl.UpdatedAt = time.Now()
			return
		}
	}
}

// AppliesTo checks if this price list applies to a given customer segment
func (pl *PriceList) AppliesTo(customerSegment string) bool {
	if len(pl.CustomerSegments) == 0 {
		return true // Applies to all if no segments specified
	}

	for _, segment := range pl.CustomerSegments {
		if segment == customerSegment {
			return true
		}
	}
	return false
}
