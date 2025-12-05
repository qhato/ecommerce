package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// TaxRateType defines the type of tax rate calculation
type TaxRateType string

const (
	TaxRateTypePercentage TaxRateType = "PERCENTAGE" // Percentage-based tax
	TaxRateTypeFlat       TaxRateType = "FLAT"       // Flat amount tax
	TaxRateTypeCompound   TaxRateType = "COMPOUND"   // Compound tax (tax on tax)
)

// TaxCategory defines categories of taxable items
type TaxCategory string

const (
	TaxCategoryGeneral   TaxCategory = "GENERAL"   // General merchandise
	TaxCategoryFood      TaxCategory = "FOOD"      // Food and beverages
	TaxCategoryClothing  TaxCategory = "CLOTHING"  // Clothing and apparel
	TaxCategoryDigital   TaxCategory = "DIGITAL"   // Digital goods and services
	TaxCategoryShipping  TaxCategory = "SHIPPING"  // Shipping and delivery
	TaxCategoryService   TaxCategory = "SERVICE"   // Services
	TaxCategoryExempt    TaxCategory = "EXEMPT"    // Exempt items
)

// TaxRate represents a tax rate for a jurisdiction
// Business Logic: Define tax rates with categories, thresholds, and compound support
type TaxRate struct {
	ID                int64
	JurisdictionID    int64
	Name              string
	TaxType           TaxRateType
	Rate              decimal.Decimal // For percentage rates, 0.0825 = 8.25%
	TaxCategory       TaxCategory
	IsCompound        bool // Whether this tax is calculated on subtotal + previous taxes
	IsShippingTaxable bool
	MinThreshold      *decimal.Decimal // Minimum amount for tax to apply
	MaxThreshold      *decimal.Decimal // Maximum amount for tax to apply
	Priority          int              // Lower = applied first
	IsActive          bool
	StartDate         *time.Time
	EndDate           *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NewTaxRate creates a new TaxRate
func NewTaxRate(jurisdictionID int64, name string, taxType TaxRateType, rate decimal.Decimal, category TaxCategory) (*TaxRate, error) {
	if jurisdictionID == 0 {
		return nil, ErrJurisdictionIDRequired
	}
	if name == "" {
		return nil, ErrTaxRateNameRequired
	}
	if rate.IsNegative() {
		return nil, ErrTaxRateCannotBeNegative
	}

	now := time.Now()
	return &TaxRate{
		JurisdictionID:    jurisdictionID,
		Name:              name,
		TaxType:           taxType,
		Rate:              rate,
		TaxCategory:       category,
		IsCompound:        false,
		IsShippingTaxable: true,
		Priority:          0,
		IsActive:          true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

// CalculateTax calculates the tax amount for a given price and quantity
func (tr *TaxRate) CalculateTax(price decimal.Decimal, quantity int, existingTaxes decimal.Decimal) decimal.Decimal {
	// Calculate total amount
	amount := price.Mul(decimal.NewFromInt(int64(quantity)))

	// Check thresholds
	if tr.MinThreshold != nil && amount.LessThan(*tr.MinThreshold) {
		return decimal.Zero
	}
	if tr.MaxThreshold != nil && amount.GreaterThan(*tr.MaxThreshold) {
		return decimal.Zero
	}

	// Calculate tax based on type
	switch tr.TaxType {
	case TaxRateTypePercentage:
		taxableAmount := amount
		if tr.IsCompound {
			// Compound tax: calculate on subtotal + existing taxes
			taxableAmount = amount.Add(existingTaxes)
		}
		return taxableAmount.Mul(tr.Rate)

	case TaxRateTypeFlat:
		// Flat tax per quantity
		return tr.Rate.Mul(decimal.NewFromInt(int64(quantity)))

	case TaxRateTypeCompound:
		// Compound tax: calculate on subtotal + existing taxes
		taxableAmount := amount.Add(existingTaxes)
		return taxableAmount.Mul(tr.Rate)

	default:
		return decimal.Zero
	}
}

// AppliesTo checks if this rate applies to a given category and amount
func (tr *TaxRate) AppliesTo(category TaxCategory, subtotal decimal.Decimal) bool {
	// Check category match
	if tr.TaxCategory != category {
		return false
	}

	// Check thresholds
	if tr.MinThreshold != nil && subtotal.LessThan(*tr.MinThreshold) {
		return false
	}
	if tr.MaxThreshold != nil && subtotal.GreaterThan(*tr.MaxThreshold) {
		return false
	}

	return true
}

// IsCurrentlyActive checks if the rate is currently active
func (tr *TaxRate) IsCurrentlyActive() bool {
	if !tr.IsActive {
		return false
	}

	now := time.Now()

	// Check start date
	if tr.StartDate != nil && now.Before(*tr.StartDate) {
		return false
	}

	// Check end date
	if tr.EndDate != nil && now.After(*tr.EndDate) {
		return false
	}

	return true
}

// Activate activates the tax rate
func (tr *TaxRate) Activate() {
	tr.IsActive = true
	tr.UpdatedAt = time.Now()
}

// Deactivate deactivates the tax rate
func (tr *TaxRate) Deactivate() {
	tr.IsActive = false
	tr.UpdatedAt = time.Now()
}

// SetThresholds sets the minimum and maximum thresholds
func (tr *TaxRate) SetThresholds(min, max *decimal.Decimal) error {
	if min != nil && max != nil && max.LessThan(*min) {
		return ErrInvalidThresholdRange
	}
	tr.MinThreshold = min
	tr.MaxThreshold = max
	tr.UpdatedAt = time.Now()
	return nil
}

// UpdateRate updates the tax rate value
func (tr *TaxRate) UpdateRate(rate decimal.Decimal) error {
	if rate.IsNegative() {
		return ErrTaxRateCannotBeNegative
	}
	tr.Rate = rate
	tr.UpdatedAt = time.Now()
	return nil
}
