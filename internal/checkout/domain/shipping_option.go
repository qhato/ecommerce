package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// ShippingSpeed represents the delivery speed
type ShippingSpeed string

const (
	ShippingSpeedStandard   ShippingSpeed = "STANDARD"    // 5-7 business days
	ShippingSpeedExpedited  ShippingSpeed = "EXPEDITED"   // 2-3 business days
	ShippingSpeedOvernight  ShippingSpeed = "OVERNIGHT"   // Next business day
	ShippingSpeedTwoDay     ShippingSpeed = "TWO_DAY"     // 2 business days
	ShippingSpeedSameDay    ShippingSpeed = "SAME_DAY"    // Same day delivery
)

// ShippingOption represents a shipping method/option
// Business Logic: Define available shipping methods with costs and constraints
type ShippingOption struct {
	ID                   string
	Name                 string
	Description          string
	Carrier              string // e.g., "UPS", "FedEx", "USPS", "DHL"
	ServiceCode          string // Carrier-specific service code
	Speed                ShippingSpeed
	EstimatedDaysMin     int
	EstimatedDaysMax     int
	BaseCost             decimal.Decimal
	CostPerItem          decimal.Decimal
	CostPerWeight        decimal.Decimal // Per kg or lb
	FreeShippingThreshold *decimal.Decimal
	IsActive             bool
	IsInternational      bool
	RequiresSignature    bool
	AllowedCountries     []string // Empty = all countries
	ExcludedCountries    []string
	AllowedStates        []string // Empty = all states
	ExcludedStates       []string
	MaxWeight            *decimal.Decimal
	MaxDimensions        *Dimensions
	TrackingSupported    bool
	InsuranceIncluded    bool
	Priority             int // Display order
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// Dimensions represents package dimensions
type Dimensions struct {
	Length decimal.Decimal
	Width  decimal.Decimal
	Height decimal.Decimal
	Unit   string // "cm" or "in"
}

// NewShippingOption creates a new shipping option
func NewShippingOption(name, carrier string, speed ShippingSpeed, baseCost decimal.Decimal) (*ShippingOption, error) {
	if name == "" {
		return nil, ErrShippingMethodRequired
	}
	if carrier == "" {
		return nil, ErrShippingMethodRequired
	}
	if baseCost.IsNegative() {
		return nil, ErrShippingCostInvalid
	}

	now := time.Now()
	return &ShippingOption{
		ID:                   generateShippingOptionID(),
		Name:                 name,
		Carrier:              carrier,
		Speed:                speed,
		BaseCost:             baseCost,
		CostPerItem:          decimal.Zero,
		CostPerWeight:        decimal.Zero,
		IsActive:             true,
		IsInternational:      false,
		RequiresSignature:    false,
		AllowedCountries:     make([]string, 0),
		ExcludedCountries:    make([]string, 0),
		AllowedStates:        make([]string, 0),
		ExcludedStates:       make([]string, 0),
		TrackingSupported:    true,
		InsuranceIncluded:    false,
		Priority:             0,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

// CalculateCost calculates the shipping cost for given parameters
func (so *ShippingOption) CalculateCost(itemCount int, totalWeight decimal.Decimal, orderSubtotal decimal.Decimal) decimal.Decimal {
	// Check free shipping threshold
	if so.FreeShippingThreshold != nil && orderSubtotal.GreaterThanOrEqual(*so.FreeShippingThreshold) {
		return decimal.Zero
	}

	cost := so.BaseCost

	// Add per-item cost
	if !so.CostPerItem.IsZero() {
		itemCost := so.CostPerItem.Mul(decimal.NewFromInt(int64(itemCount)))
		cost = cost.Add(itemCost)
	}

	// Add per-weight cost
	if !so.CostPerWeight.IsZero() && !totalWeight.IsZero() {
		weightCost := so.CostPerWeight.Mul(totalWeight)
		cost = cost.Add(weightCost)
	}

	return cost
}

// IsAvailableForLocation checks if shipping is available for a location
func (so *ShippingOption) IsAvailableForLocation(country, stateProvince string) bool {
	if !so.IsActive {
		return false
	}

	// Check excluded countries
	for _, c := range so.ExcludedCountries {
		if c == country {
			return false
		}
	}

	// Check allowed countries (if specified)
	if len(so.AllowedCountries) > 0 {
		found := false
		for _, c := range so.AllowedCountries {
			if c == country {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check excluded states
	for _, s := range so.ExcludedStates {
		if s == stateProvince {
			return false
		}
	}

	// Check allowed states (if specified)
	if len(so.AllowedStates) > 0 {
		found := false
		for _, s := range so.AllowedStates {
			if s == stateProvince {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// CanHandlePackage checks if the option can handle the package size/weight
func (so *ShippingOption) CanHandlePackage(weight decimal.Decimal, dimensions *Dimensions) bool {
	// Check weight limit
	if so.MaxWeight != nil && weight.GreaterThan(*so.MaxWeight) {
		return false
	}

	// Check dimensions
	if so.MaxDimensions != nil && dimensions != nil {
		if dimensions.Length.GreaterThan(so.MaxDimensions.Length) ||
			dimensions.Width.GreaterThan(so.MaxDimensions.Width) ||
			dimensions.Height.GreaterThan(so.MaxDimensions.Height) {
			return false
		}
	}

	return true
}

// GetEstimatedDeliveryDays returns the estimated delivery time range
func (so *ShippingOption) GetEstimatedDeliveryDays() (int, int) {
	return so.EstimatedDaysMin, so.EstimatedDaysMax
}

// Activate activates the shipping option
func (so *ShippingOption) Activate() {
	so.IsActive = true
	so.UpdatedAt = time.Now()
}

// Deactivate deactivates the shipping option
func (so *ShippingOption) Deactivate() {
	so.IsActive = false
	so.UpdatedAt = time.Now()
}

// SetFreeShippingThreshold sets the free shipping threshold
func (so *ShippingOption) SetFreeShippingThreshold(threshold decimal.Decimal) {
	so.FreeShippingThreshold = &threshold
	so.UpdatedAt = time.Now()
}

// Private helper methods

func generateShippingOptionID() string {
	return "SHIP-" + time.Now().Format("20060102150405")
}
