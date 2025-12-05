package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// TaxCalculationRequest represents a request to calculate taxes
// Business Logic: Contexto completo para cálculo de impuestos
type TaxCalculationRequest struct {
	OrderID         *int64
	CustomerID      *string
	ShippingAddress Address
	BillingAddress  Address
	Items           []TaxableItem
	ShippingAmount  decimal.Decimal
	CalculationDate time.Time
}

// Address represents a postal address for tax calculation
type Address struct {
	Country       string
	StateProvince string
	County        string
	City          string
	PostalCode    string
	AddressLine1  string
	AddressLine2  string
}

// TaxableItem represents an item to be taxed
type TaxableItem struct {
	ItemID      string
	SKU         string
	Description string
	Quantity    int
	UnitPrice   decimal.Decimal
	Subtotal    decimal.Decimal
	TaxCategory TaxCategory
	IsExempt    bool
}

// TaxCalculationResult represents the result of a tax calculation
type TaxCalculationResult struct {
	OrderID         *int64
	Items           []TaxedItem
	ShippingTax     decimal.Decimal
	TotalTax        decimal.Decimal
	Subtotal        decimal.Decimal
	TotalAmount     decimal.Decimal
	Breakdowns      []TaxBreakdown
	CalculatedAt    time.Time
	JurisdictionsUsed []string // List of jurisdiction codes used
}

// TaxedItem represents an item with calculated taxes
type TaxedItem struct {
	ItemID      string
	SKU         string
	Quantity    int
	UnitPrice   decimal.Decimal
	Subtotal    decimal.Decimal
	TaxAmount   decimal.Decimal
	TaxCategory TaxCategory
	Taxes       []AppliedTax
}

// AppliedTax represents a tax that was applied
type AppliedTax struct {
	JurisdictionCode string
	JurisdictionName string
	TaxRateName      string
	TaxType          TaxRateType
	Rate             decimal.Decimal
	TaxableAmount    decimal.Decimal
	TaxAmount        decimal.Decimal
	IsCompound       bool
}

// TaxBreakdown represents taxes grouped by jurisdiction
type TaxBreakdown struct {
	JurisdictionCode string
	JurisdictionName string
	JurisdictionType TaxJurisdictionType
	TotalTaxAmount   decimal.Decimal
	Rates            []AppliedTax
}

// NewTaxCalculationRequest creates a new tax calculation request
func NewTaxCalculationRequest(shippingAddress Address) *TaxCalculationRequest {
	return &TaxCalculationRequest{
		ShippingAddress: shippingAddress,
		Items:           make([]TaxableItem, 0),
		ShippingAmount:  decimal.Zero,
		CalculationDate: time.Now(),
	}
}

// AddItem adds an item to the calculation request
func (req *TaxCalculationRequest) AddItem(item TaxableItem) {
	req.Items = append(req.Items, item)
}

// SetShippingAmount sets the shipping amount
func (req *TaxCalculationRequest) SetShippingAmount(amount decimal.Decimal) {
	req.ShippingAmount = amount
}

// GetSubtotal calculates the subtotal of all items
func (req *TaxCalculationRequest) GetSubtotal() decimal.Decimal {
	subtotal := decimal.Zero
	for _, item := range req.Items {
		subtotal = subtotal.Add(item.Subtotal)
	}
	return subtotal
}

// NewTaxCalculationResult creates a new tax calculation result
func NewTaxCalculationResult() *TaxCalculationResult {
	return &TaxCalculationResult{
		Items:             make([]TaxedItem, 0),
		ShippingTax:       decimal.Zero,
		TotalTax:          decimal.Zero,
		Subtotal:          decimal.Zero,
		TotalAmount:       decimal.Zero,
		Breakdowns:        make([]TaxBreakdown, 0),
		CalculatedAt:      time.Now(),
		JurisdictionsUsed: make([]string, 0),
	}
}

// AddTaxedItem adds a taxed item to the result
func (result *TaxCalculationResult) AddTaxedItem(item TaxedItem) {
	result.Items = append(result.Items, item)
	result.Subtotal = result.Subtotal.Add(item.Subtotal)
	result.TotalTax = result.TotalTax.Add(item.TaxAmount)
}

// SetShippingTax sets the shipping tax
func (result *TaxCalculationResult) SetShippingTax(amount decimal.Decimal) {
	result.ShippingTax = amount
	result.TotalTax = result.TotalTax.Add(amount)
}

// Finalize finalizes the calculation result
func (result *TaxCalculationResult) Finalize() {
	result.TotalAmount = result.Subtotal.Add(result.TotalTax)

	// Group taxes by jurisdiction for breakdowns
	breakdownMap := make(map[string]*TaxBreakdown)

	for _, item := range result.Items {
		for _, tax := range item.Taxes {
			if breakdown, exists := breakdownMap[tax.JurisdictionCode]; exists {
				breakdown.TotalTaxAmount = breakdown.TotalTaxAmount.Add(tax.TaxAmount)
				breakdown.Rates = append(breakdown.Rates, tax)
			} else {
				breakdownMap[tax.JurisdictionCode] = &TaxBreakdown{
					JurisdictionCode: tax.JurisdictionCode,
					JurisdictionName: tax.JurisdictionName,
					TotalTaxAmount:   tax.TaxAmount,
					Rates:            []AppliedTax{tax},
				}
			}
		}
	}

	// Convert map to slice
	result.Breakdowns = make([]TaxBreakdown, 0, len(breakdownMap))
	for _, breakdown := range breakdownMap {
		result.Breakdowns = append(result.Breakdowns, *breakdown)
	}
}

// GetEffectiveTaxRate calculates the effective tax rate (total tax / subtotal)
func (result *TaxCalculationResult) GetEffectiveTaxRate() decimal.Decimal {
	if result.Subtotal.IsZero() {
		return decimal.Zero
	}
	return result.TotalTax.Div(result.Subtotal)
}

// TaxExemption represents a tax exemption for a customer or item
type TaxExemption struct {
	ID                int64
	CustomerID        *string
	ExemptionCertificate string
	JurisdictionID    *int64 // Null = all jurisdictions
	TaxCategory       *TaxCategory // Null = all categories
	Reason            string
	IsActive          bool
	StartDate         *time.Time
	EndDate           *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NewTaxExemption creates a new tax exemption
func NewTaxExemption(customerID, certificate, reason string) (*TaxExemption, error) {
	if customerID == "" {
		return nil, ErrCustomerIDRequired
	}
	if certificate == "" {
		return nil, ErrExemptionCertificateRequired
	}

	now := time.Now()
	return &TaxExemption{
		CustomerID:           &customerID,
		ExemptionCertificate: certificate,
		Reason:               reason,
		IsActive:             true,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

// IsCurrentlyActive checks if the exemption is currently active
func (te *TaxExemption) IsCurrentlyActive() bool {
	if !te.IsActive {
		return false
	}

	now := time.Now()
	if te.StartDate != nil && now.Before(*te.StartDate) {
		return false
	}
	if te.EndDate != nil && now.After(*te.EndDate) {
		return false
	}

	return true
}

// AppliesTo checks if this exemption applies to a jurisdiction and category
func (te *TaxExemption) AppliesTo(jurisdictionID int64, category TaxCategory) bool {
	// Check jurisdiction
	if te.JurisdictionID != nil && *te.JurisdictionID != jurisdictionID {
		return false
	}

	// Check category
	if te.TaxCategory != nil && *te.TaxCategory != category {
		return false
	}

	return true
}
