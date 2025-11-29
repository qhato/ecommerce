package domain

import (
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

// TaxableItem represents an item that can be taxed
type TaxableItem struct {
	ItemID     string
	ProductID  string
	CategoryID string
	Amount     decimal.Decimal
	Quantity   int
}

// TaxAddress represents an address for tax calculation
type TaxAddress struct {
	Country    string
	Region     string // State/Province
	PostalCode string
	City       string
	County     string
}

// TaxCalculationContext holds all data needed for tax calculation
type TaxCalculationContext struct {
	Items            []TaxableItem
	ShippingAddress  *TaxAddress
	BillingAddress   *TaxAddress
	CustomerID       *string
	CalculationDate  time.Time
	ShippingAmount   decimal.Decimal
	HandlingAmount   decimal.Decimal
}

// TaxCalculationResult represents the result of tax calculation
type TaxCalculationResult struct {
	TotalTax    decimal.Decimal
	TaxDetails  []*TaxDetail
	Exemptions  []*AppliedExemption
	ItemTaxes   map[string]decimal.Decimal // ItemID -> tax amount
}

// AppliedExemption represents an exemption that was applied
type AppliedExemption struct {
	ExemptionCode string
	Reason        string
	TaxType       TaxType
	ItemID        *string
	Amount        decimal.Decimal // Tax amount exempted
}

// TaxCalculator calculates taxes based on jurisdiction and rates
type TaxCalculator struct {
	rateRepository       TaxRateRepository
	jurisdictionRepo     TaxJurisdictionRepository
	exemptionRepo        TaxExemptionRepository
}

// TaxRateRepository defines repository for tax rates
type TaxRateRepository interface {
	FindByJurisdiction(country, region string, date time.Time) ([]*TaxRate, error)
	FindActiveRates(date time.Time) ([]*TaxRate, error)
}

// TaxJurisdictionRepository defines repository for tax jurisdictions
type TaxJurisdictionRepository interface {
	FindByAddress(country, region, postalCode, city, county string) ([]*TaxJurisdiction, error)
}

// TaxExemptionRepository defines repository for tax exemptions
type TaxExemptionRepository interface {
	FindByCustomer(customerID string, date time.Time) ([]*TaxExemption, error)
	FindByProduct(productID string, date time.Time) ([]*TaxExemption, error)
	FindByCategory(categoryID string, date time.Time) ([]*TaxExemption, error)
}

// NewTaxCalculator creates a new TaxCalculator
func NewTaxCalculator(
	rateRepo TaxRateRepository,
	jurisdictionRepo TaxJurisdictionRepository,
	exemptionRepo TaxExemptionRepository,
) *TaxCalculator {
	return &TaxCalculator{
		rateRepository:   rateRepo,
		jurisdictionRepo: jurisdictionRepo,
		exemptionRepo:    exemptionRepo,
	}
}

// Calculate calculates taxes for the given context
func (tc *TaxCalculator) Calculate(ctx *TaxCalculationContext) (*TaxCalculationResult, error) {
	if ctx.ShippingAddress == nil {
		return nil, NewDomainError("Shipping address is required for tax calculation")
	}

	result := &TaxCalculationResult{
		TotalTax:   decimal.Zero,
		TaxDetails: make([]*TaxDetail, 0),
		Exemptions: make([]*AppliedExemption, 0),
		ItemTaxes:  make(map[string]decimal.Decimal),
	}

	// Get applicable tax rates based on jurisdiction
	taxRates, err := tc.rateRepository.FindByJurisdiction(
		ctx.ShippingAddress.Country,
		ctx.ShippingAddress.Region,
		ctx.CalculationDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax rates: %w", err)
	}

	// Filter to only effective rates
	effectiveRates := make([]*TaxRate, 0)
	for _, rate := range taxRates {
		if rate.IsEffective(ctx.CalculationDate) {
			effectiveRates = append(effectiveRates, rate)
		}
	}

	// Sort rates by priority
	sort.Slice(effectiveRates, func(i, j int) bool {
		return effectiveRates[i].Priority < effectiveRates[j].Priority
	})

	// Get applicable exemptions
	var customerExemptions []*TaxExemption
	if ctx.CustomerID != nil {
		customerExemptions, err = tc.exemptionRepo.FindByCustomer(*ctx.CustomerID, ctx.CalculationDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get customer exemptions: %w", err)
		}
	}

	// Calculate tax for each item
	for _, item := range ctx.Items {
		itemTax, itemDetails, itemExemptions := tc.calculateItemTax(
			item,
			effectiveRates,
			customerExemptions,
			ctx,
		)

		result.ItemTaxes[item.ItemID] = itemTax
		result.TotalTax = result.TotalTax.Add(itemTax)
		result.TaxDetails = append(result.TaxDetails, itemDetails...)
		result.Exemptions = append(result.Exemptions, itemExemptions...)
	}

	// Calculate tax on shipping if applicable
	if ctx.ShippingAmount.GreaterThan(decimal.Zero) {
		shippingTax, shippingDetails := tc.calculateShippingTax(
			ctx.ShippingAmount,
			effectiveRates,
			ctx,
		)
		result.TotalTax = result.TotalTax.Add(shippingTax)
		result.TaxDetails = append(result.TaxDetails, shippingDetails...)
	}

	return result, nil
}

// calculateItemTax calculates tax for a single item
func (tc *TaxCalculator) calculateItemTax(
	item TaxableItem,
	rates []*TaxRate,
	customerExemptions []*TaxExemption,
	ctx *TaxCalculationContext,
) (decimal.Decimal, []*TaxDetail, []*AppliedExemption) {

	totalTax := decimal.Zero
	details := make([]*TaxDetail, 0)
	exemptions := make([]*AppliedExemption, 0)

	// Get product and category exemptions
	var productExemptions, categoryExemptions []*TaxExemption
	var err error

	if tc.exemptionRepo != nil {
		productExemptions, err = tc.exemptionRepo.FindByProduct(item.ProductID, ctx.CalculationDate)
		if err == nil && len(productExemptions) == 0 {
			categoryExemptions, _ = tc.exemptionRepo.FindByCategory(item.CategoryID, ctx.CalculationDate)
		}
	}

	// Combine all exemptions
	allExemptions := make([]*TaxExemption, 0)
	allExemptions = append(allExemptions, customerExemptions...)
	allExemptions = append(allExemptions, productExemptions...)
	allExemptions = append(allExemptions, categoryExemptions...)

	// Calculate tax for each rate
	for _, rate := range rates {
		// Check if this tax type is exempt
		isExempt := false
		var appliedExemption *TaxExemption

		for _, exemption := range allExemptions {
			if exemption.IsEffective(ctx.CalculationDate) &&
				exemption.ExemptsTaxType(rate.TaxType) {
				isExempt = true
				appliedExemption = exemption
				break
			}
		}

		if isExempt && appliedExemption != nil {
			// Calculate what the tax would have been
			taxAmount := item.Amount.Mul(decimal.NewFromFloat(rate.Rate))

			// Record exemption
			itemIDCopy := item.ItemID
			exemptions = append(exemptions, &AppliedExemption{
				ExemptionCode: appliedExemption.ExemptionCode,
				Reason:        appliedExemption.Reason,
				TaxType:       rate.TaxType,
				ItemID:        &itemIDCopy,
				Amount:        taxAmount,
			})
			continue
		}

		// Calculate tax
		taxAmount := item.Amount.Mul(decimal.NewFromFloat(rate.Rate))
		totalTax = totalTax.Add(taxAmount)

		// Create tax detail
		taxAmountFloat, _ := taxAmount.Float64()
		detail, _ := NewTaxDetail(
			taxAmountFloat,
			rate.Country,
			rate.JurisdictionName,
			rate.Rate,
			rate.Region,
			rate.TaxName,
			string(rate.TaxType),
			"USD", // TODO: get from context
		)
		if detail != nil {
			details = append(details, detail)
		}
	}

	return totalTax, details, exemptions
}

// calculateShippingTax calculates tax on shipping
func (tc *TaxCalculator) calculateShippingTax(
	shippingAmount decimal.Decimal,
	rates []*TaxRate,
	ctx *TaxCalculationContext,
) (decimal.Decimal, []*TaxDetail) {

	totalTax := decimal.Zero
	details := make([]*TaxDetail, 0)

	// Only apply sales/VAT/GST taxes to shipping (not excise, etc.)
	for _, rate := range rates {
		if rate.TaxType == TaxTypeSales || rate.TaxType == TaxTypeVAT || rate.TaxType == TaxTypeGST {
			taxAmount := shippingAmount.Mul(decimal.NewFromFloat(rate.Rate))
			totalTax = totalTax.Add(taxAmount)

			taxAmountFloat, _ := taxAmount.Float64()
			detail, _ := NewTaxDetail(
				taxAmountFloat,
				rate.Country,
				rate.JurisdictionName,
				rate.Rate,
				rate.Region,
				fmt.Sprintf("%s (Shipping)", rate.TaxName),
				string(rate.TaxType),
				"USD",
			)
			if detail != nil {
				details = append(details, detail)
			}
		}
	}

	return totalTax, details
}

// EstimateTax provides a quick estimate without full calculation
func (tc *TaxCalculator) EstimateTax(amount decimal.Decimal, address *TaxAddress, date time.Time) (decimal.Decimal, error) {
	if address == nil {
		return decimal.Zero, NewDomainError("Address is required for tax estimation")
	}

	// Get applicable tax rates
	taxRates, err := tc.rateRepository.FindByJurisdiction(
		address.Country,
		address.Region,
		date,
	)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get tax rates: %w", err)
	}

	// Sum all effective rates
	totalRate := 0.0
	for _, rate := range taxRates {
		if rate.IsEffective(date) {
			totalRate += rate.Rate
		}
	}

	estimatedTax := amount.Mul(decimal.NewFromFloat(totalRate))
	return estimatedTax, nil
}
