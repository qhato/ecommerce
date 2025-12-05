package queries

import (
	"context"
	"fmt"
	"sort"

	"github.com/qhato/ecommerce/internal/tax/domain"
	"github.com/shopspring/decimal"
)

// TaxCalculatorService implements the tax calculation business logic
type TaxCalculatorService struct {
	jurisdictionRepo domain.TaxJurisdictionRepository
	rateRepo         domain.TaxRateRepository
	exemptionRepo    domain.TaxExemptionRepository
}

// NewTaxCalculatorService creates a new tax calculator service
func NewTaxCalculatorService(
	jurisdictionRepo domain.TaxJurisdictionRepository,
	rateRepo domain.TaxRateRepository,
	exemptionRepo domain.TaxExemptionRepository,
) *TaxCalculatorService {
	return &TaxCalculatorService{
		jurisdictionRepo: jurisdictionRepo,
		rateRepo:         rateRepo,
		exemptionRepo:    exemptionRepo,
	}
}

// Calculate calculates taxes for a request
func (s *TaxCalculatorService) Calculate(ctx context.Context, request *domain.TaxCalculationRequest) (*domain.TaxCalculationResult, error) {
	// Validate request
	if err := s.validateRequest(request); err != nil {
		return nil, err
	}

	// Find applicable jurisdictions for the shipping address
	jurisdictions, err := s.findApplicableJurisdictions(ctx, request.ShippingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to find jurisdictions: %w", err)
	}

	if len(jurisdictions) == 0 {
		return nil, domain.ErrNoApplicableJurisdictions
	}

	// Get customer exemptions if customer ID is provided
	var exemptions []*domain.TaxExemption
	if request.CustomerID != nil {
		exemptions, err = s.exemptionRepo.FindActiveExemptions(ctx, *request.CustomerID)
		if err != nil {
			return nil, fmt.Errorf("failed to find exemptions: %w", err)
		}
	}

	// Calculate with exemptions
	return s.CalculateWithExemptions(ctx, request, exemptions)
}

// CalculateWithExemptions calculates taxes applying customer exemptions
func (s *TaxCalculatorService) CalculateWithExemptions(ctx context.Context, request *domain.TaxCalculationRequest, exemptions []*domain.TaxExemption) (*domain.TaxCalculationResult, error) {
	// Validate request
	if err := s.validateRequest(request); err != nil {
		return nil, err
	}

	// Find applicable jurisdictions
	jurisdictions, err := s.findApplicableJurisdictions(ctx, request.ShippingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to find jurisdictions: %w", err)
	}

	if len(jurisdictions) == 0 {
		return nil, domain.ErrNoApplicableJurisdictions
	}

	// Sort jurisdictions by priority (lower = applied first)
	sort.Slice(jurisdictions, func(i, j int) bool {
		return jurisdictions[i].Priority < jurisdictions[j].Priority
	})

	// Extract jurisdiction IDs
	jurisdictionIDs := make([]int64, len(jurisdictions))
	jurisdictionMap := make(map[int64]*domain.TaxJurisdiction)
	for i, j := range jurisdictions {
		jurisdictionIDs[i] = j.ID
		jurisdictionMap[j.ID] = j
	}

	// Create result
	result := domain.NewTaxCalculationResult()
	result.OrderID = request.OrderID

	// Calculate taxes for each item
	for _, item := range request.Items {
		taxedItem, err := s.calculateItemTax(ctx, item, jurisdictionIDs, jurisdictionMap, exemptions)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate tax for item %s: %w", item.ItemID, err)
		}
		result.AddTaxedItem(taxedItem)
	}

	// Calculate shipping tax if applicable
	if !request.ShippingAmount.IsZero() {
		shippingTax, err := s.calculateShippingTax(ctx, request.ShippingAmount, jurisdictionIDs, jurisdictionMap, exemptions)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate shipping tax: %w", err)
		}
		result.SetShippingTax(shippingTax)
	}

	// Track jurisdictions used
	for _, j := range jurisdictions {
		result.JurisdictionsUsed = append(result.JurisdictionsUsed, j.Code)
	}

	// Finalize result (calculate breakdowns and totals)
	result.Finalize()

	return result, nil
}

// EstimateTax provides a quick tax estimate without full calculation
func (s *TaxCalculatorService) EstimateTax(ctx context.Context, address domain.Address, subtotal decimal.Decimal) (decimal.Decimal, error) {
	// Find applicable jurisdictions
	jurisdictions, err := s.findApplicableJurisdictions(ctx, address)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to find jurisdictions: %w", err)
	}

	if len(jurisdictions) == 0 {
		return decimal.Zero, domain.ErrNoApplicableJurisdictions
	}

	// Extract jurisdiction IDs
	jurisdictionIDs := make([]int64, len(jurisdictions))
	for i, j := range jurisdictions {
		jurisdictionIDs[i] = j.ID
	}

	// Find applicable tax rates for GENERAL category
	rates, err := s.rateRepo.FindApplicableRates(ctx, jurisdictionIDs, domain.TaxCategoryGeneral, true)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to find tax rates: %w", err)
	}

	if len(rates) == 0 {
		return decimal.Zero, nil
	}

	// Calculate estimate using simple addition (no compounding for estimates)
	totalTax := decimal.Zero
	for _, rate := range rates {
		if rate.AppliesTo(domain.TaxCategoryGeneral, subtotal) {
			tax := rate.CalculateTax(subtotal, 1, decimal.Zero)
			totalTax = totalTax.Add(tax)
		}
	}

	return totalTax, nil
}

// ValidateAddress validates if an address has applicable tax jurisdictions
func (s *TaxCalculatorService) ValidateAddress(ctx context.Context, address domain.Address) (bool, error) {
	jurisdictions, err := s.findApplicableJurisdictions(ctx, address)
	if err != nil {
		return false, err
	}
	return len(jurisdictions) > 0, nil
}

// Private helper methods

func (s *TaxCalculatorService) validateRequest(request *domain.TaxCalculationRequest) error {
	if request.ShippingAddress.Country == "" {
		return domain.ErrShippingAddressRequired
	}
	if len(request.Items) == 0 {
		return domain.ErrNoItemsToCalculate
	}
	for _, item := range request.Items {
		if item.Quantity < 0 {
			return domain.ErrNegativeQuantity
		}
		if item.UnitPrice.IsNegative() {
			return domain.ErrNegativePrice
		}
	}
	return nil
}

func (s *TaxCalculatorService) findApplicableJurisdictions(ctx context.Context, address domain.Address) ([]*domain.TaxJurisdiction, error) {
	jurisdictions, err := s.jurisdictionRepo.FindByLocation(
		ctx,
		address.Country,
		address.StateProvince,
		address.County,
		address.City,
		address.PostalCode,
	)
	if err != nil {
		return nil, err
	}

	// Filter only active jurisdictions
	active := make([]*domain.TaxJurisdiction, 0)
	for _, j := range jurisdictions {
		if j.IsActive {
			active = append(active, j)
		}
	}

	return active, nil
}

func (s *TaxCalculatorService) calculateItemTax(
	ctx context.Context,
	item domain.TaxableItem,
	jurisdictionIDs []int64,
	jurisdictionMap map[int64]*domain.TaxJurisdiction,
	exemptions []*domain.TaxExemption,
) (domain.TaxedItem, error) {
	taxedItem := domain.TaxedItem{
		ItemID:      item.ItemID,
		SKU:         item.SKU,
		Quantity:    item.Quantity,
		UnitPrice:   item.UnitPrice,
		Subtotal:    item.Subtotal,
		TaxAmount:   decimal.Zero,
		TaxCategory: item.TaxCategory,
		Taxes:       make([]domain.AppliedTax, 0),
	}

	// If item is exempt, return with zero tax
	if item.IsExempt {
		return taxedItem, nil
	}

	// Find applicable tax rates
	rates, err := s.rateRepo.FindApplicableRates(ctx, jurisdictionIDs, item.TaxCategory, true)
	if err != nil {
		return taxedItem, fmt.Errorf("failed to find tax rates: %w", err)
	}

	// Sort rates by priority (lower = applied first)
	sort.Slice(rates, func(i, j int) bool {
		return rates[i].Priority < rates[j].Priority
	})

	// Calculate taxes
	cumulativeTax := decimal.Zero

	for _, rate := range rates {
		// Check if rate is currently active
		if !rate.IsCurrentlyActive() {
			continue
		}

		// Check if rate applies to this item
		if !rate.AppliesTo(item.TaxCategory, item.Subtotal) {
			continue
		}

		// Check if customer has an exemption for this rate
		if s.isExempt(rate, exemptions) {
			continue
		}

		// Calculate tax for this rate
		taxableAmount := item.Subtotal
		if rate.IsCompound {
			// Compound tax is calculated on subtotal + previous taxes
			taxableAmount = item.Subtotal.Add(cumulativeTax)
		}

		tax := rate.CalculateTax(taxableAmount, item.Quantity, cumulativeTax)
		cumulativeTax = cumulativeTax.Add(tax)

		// Get jurisdiction details
		jurisdiction := jurisdictionMap[rate.JurisdictionID]

		// Record applied tax
		appliedTax := domain.AppliedTax{
			JurisdictionCode: jurisdiction.Code,
			JurisdictionName: jurisdiction.Name,
			TaxRateName:      rate.Name,
			TaxType:          rate.TaxType,
			Rate:             rate.Rate,
			TaxableAmount:    taxableAmount,
			TaxAmount:        tax,
			IsCompound:       rate.IsCompound,
		}

		taxedItem.Taxes = append(taxedItem.Taxes, appliedTax)
	}

	taxedItem.TaxAmount = cumulativeTax
	return taxedItem, nil
}

func (s *TaxCalculatorService) calculateShippingTax(
	ctx context.Context,
	shippingAmount decimal.Decimal,
	jurisdictionIDs []int64,
	jurisdictionMap map[int64]*domain.TaxJurisdiction,
	exemptions []*domain.TaxExemption,
) (decimal.Decimal, error) {
	// Find applicable tax rates for SHIPPING category
	rates, err := s.rateRepo.FindApplicableRates(ctx, jurisdictionIDs, domain.TaxCategoryShipping, true)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to find shipping tax rates: %w", err)
	}

	// Sort rates by priority
	sort.Slice(rates, func(i, j int) bool {
		return rates[i].Priority < rates[j].Priority
	})

	totalTax := decimal.Zero

	for _, rate := range rates {
		// Check if rate is currently active
		if !rate.IsCurrentlyActive() {
			continue
		}

		// Check if shipping is taxable for this rate
		if !rate.IsShippingTaxable {
			continue
		}

		// Check if customer has an exemption
		if s.isExempt(rate, exemptions) {
			continue
		}

		// Calculate tax
		taxableAmount := shippingAmount
		if rate.IsCompound {
			taxableAmount = shippingAmount.Add(totalTax)
		}

		tax := rate.CalculateTax(taxableAmount, 1, totalTax)
		totalTax = totalTax.Add(tax)
	}

	return totalTax, nil
}

func (s *TaxCalculatorService) isExempt(rate *domain.TaxRate, exemptions []*domain.TaxExemption) bool {
	for _, exemption := range exemptions {
		if !exemption.IsCurrentlyActive() {
			continue
		}
		if exemption.AppliesTo(rate.JurisdictionID, rate.TaxCategory) {
			return true
		}
	}
	return false
}

// Query methods for retrieving jurisdictions, rates, and exemptions

// GetJurisdictionByID retrieves a jurisdiction by ID
func (s *TaxCalculatorService) GetJurisdictionByID(ctx context.Context, id int64) (*domain.TaxJurisdiction, error) {
	return s.jurisdictionRepo.FindByID(ctx, id)
}

// GetJurisdictionByCode retrieves a jurisdiction by code
func (s *TaxCalculatorService) GetJurisdictionByCode(ctx context.Context, code string) (*domain.TaxJurisdiction, error) {
	return s.jurisdictionRepo.FindByCode(ctx, code)
}

// GetAllJurisdictions retrieves all jurisdictions
func (s *TaxCalculatorService) GetAllJurisdictions(ctx context.Context, activeOnly bool) ([]*domain.TaxJurisdiction, error) {
	return s.jurisdictionRepo.FindAll(ctx, activeOnly)
}

// GetJurisdictionsByCountry retrieves jurisdictions for a country
func (s *TaxCalculatorService) GetJurisdictionsByCountry(ctx context.Context, country string, activeOnly bool) ([]*domain.TaxJurisdiction, error) {
	return s.jurisdictionRepo.FindByCountry(ctx, country, activeOnly)
}

// GetTaxRateByID retrieves a tax rate by ID
func (s *TaxCalculatorService) GetTaxRateByID(ctx context.Context, id int64) (*domain.TaxRate, error) {
	return s.rateRepo.FindByID(ctx, id)
}

// GetTaxRatesByJurisdiction retrieves tax rates for a jurisdiction
func (s *TaxCalculatorService) GetTaxRatesByJurisdiction(ctx context.Context, jurisdictionID int64, activeOnly bool) ([]*domain.TaxRate, error) {
	return s.rateRepo.FindByJurisdiction(ctx, jurisdictionID, activeOnly)
}

// GetAllTaxRates retrieves all tax rates
func (s *TaxCalculatorService) GetAllTaxRates(ctx context.Context, activeOnly bool) ([]*domain.TaxRate, error) {
	return s.rateRepo.FindAll(ctx, activeOnly)
}

// GetExemptionByID retrieves an exemption by ID
func (s *TaxCalculatorService) GetExemptionByID(ctx context.Context, id int64) (*domain.TaxExemption, error) {
	return s.exemptionRepo.FindByID(ctx, id)
}

// GetExemptionsByCustomer retrieves exemptions for a customer
func (s *TaxCalculatorService) GetExemptionsByCustomer(ctx context.Context, customerID string, activeOnly bool) ([]*domain.TaxExemption, error) {
	return s.exemptionRepo.FindByCustomerID(ctx, customerID, activeOnly)
}

// GetAllExemptions retrieves all exemptions
func (s *TaxCalculatorService) GetAllExemptions(ctx context.Context, activeOnly bool) ([]*domain.TaxExemption, error) {
	return s.exemptionRepo.FindAll(ctx, activeOnly)
}
