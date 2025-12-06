package application

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/checkout/domain"
	invDomain "github.com/qhato/ecommerce/internal/inventory/domain"
	shipDomain "github.com/qhato/ecommerce/internal/shipping/domain"
	taxDomain "github.com/qhato/ecommerce/internal/tax/domain"
	"github.com/shopspring/decimal"
)

// CheckoutOrchestrator coordinates checkout flow with other services
type CheckoutOrchestrator struct {
	sessionRepo        domain.CheckoutSessionRepository
	inventoryRepo      invDomain.InventoryRepository
	reservationRepo    invDomain.InventoryReservationRepository
	shippingMethodRepo shipDomain.ShippingMethodRepository
	shippingRuleRepo   shipDomain.ShippingRuleRepository
	taxJurisdictionRepo taxDomain.TaxJurisdictionRepository
	taxRateRepo        taxDomain.TaxRateRepository
	taxExemptionRepo   taxDomain.TaxExemptionRepository
}

func NewCheckoutOrchestrator(
	sessionRepo domain.CheckoutSessionRepository,
	inventoryRepo invDomain.InventoryRepository,
	reservationRepo invDomain.InventoryReservationRepository,
	shippingMethodRepo shipDomain.ShippingMethodRepository,
	shippingRuleRepo shipDomain.ShippingRuleRepository,
	taxJurisdictionRepo taxDomain.TaxJurisdictionRepository,
	taxRateRepo taxDomain.TaxRateRepository,
	taxExemptionRepo taxDomain.TaxExemptionRepository,
) *CheckoutOrchestrator {
	return &CheckoutOrchestrator{
		sessionRepo:         sessionRepo,
		inventoryRepo:       inventoryRepo,
		reservationRepo:     reservationRepo,
		shippingMethodRepo:  shippingMethodRepo,
		shippingRuleRepo:    shippingRuleRepo,
		taxJurisdictionRepo: taxJurisdictionRepo,
		taxRateRepo:         taxRateRepo,
		taxExemptionRepo:    taxExemptionRepo,
	}
}

// ValidateInventoryAvailability validates inventory for all order items
func (o *CheckoutOrchestrator) ValidateInventoryAvailability(ctx context.Context, orderItems []OrderItem) error {
	for _, item := range orderItems {
		level, err := o.inventoryRepo.FindBySKUID(ctx, item.SKUID)
		if err != nil {
			return fmt.Errorf("failed to check inventory for SKU %s: %w", item.SKUID, err)
		}
		
		if level == nil {
			return fmt.Errorf("inventory not found for SKU: %s", item.SKUID)
		}

		if !level.CanReserve(item.Quantity) {
			return fmt.Errorf("insufficient inventory for SKU %s: need %d, available %d", 
				item.SKUID, item.Quantity, level.QuantityAvailable)
		}
	}
	
	return nil
}

// ReserveInventory reserves inventory for checkout
func (o *CheckoutOrchestrator) ReserveInventory(ctx context.Context, sessionID string, orderID int64, orderItems []OrderItem) error {
	for _, item := range orderItems {
		// Reserve inventory
		level, err := o.inventoryRepo.FindBySKUID(ctx, item.SKUID)
		if err != nil || level == nil {
			return fmt.Errorf("failed to find inventory for SKU %s", item.SKUID)
		}

		if err := level.Reserve(item.Quantity); err != nil {
			// Rollback previous reservations
			o.ReleaseInventory(ctx, sessionID, orderID)
			return fmt.Errorf("failed to reserve inventory: %w", err)
		}

		if err := o.inventoryRepo.Save(ctx, level); err != nil {
			o.ReleaseInventory(ctx, sessionID, orderID)
			return fmt.Errorf("failed to save inventory: %w", err)
		}

		// Create reservation record
		orderIDStr := fmt.Sprintf("%d", orderID)
		orderItemIDStr := fmt.Sprintf("%d-%s", orderID, item.SKUID)
		reservation, err := invDomain.NewInventoryReservation(
			item.SKUID, orderIDStr, orderItemIDStr, item.Quantity, 24 * 60 * 60 * 1000000000, // 24 hours
		)
		if err != nil {
			return err
		}

		if err := o.reservationRepo.Save(ctx, reservation); err != nil {
			return fmt.Errorf("failed to save reservation: %w", err)
		}
	}

	return nil
}

// ReleaseInventory releases reserved inventory
func (o *CheckoutOrchestrator) ReleaseInventory(ctx context.Context, sessionID string, orderID int64) error {
	orderIDStr := fmt.Sprintf("%d", orderID)
	reservations, err := o.reservationRepo.FindByOrderID(ctx, orderIDStr)
	if err != nil {
		return err
	}

	for _, res := range reservations {
		if res.Status == invDomain.ReservationStatusPending || res.Status == invDomain.ReservationStatusConfirmed {
			res.Release()
			o.reservationRepo.Save(ctx, res)

			// Release from inventory level
			level, _ := o.inventoryRepo.FindBySKUID(ctx, res.SKUID)
			if level != nil {
				level.Release(res.Quantity)
				o.inventoryRepo.Save(ctx, level)
			}
		}
	}

	return nil
}

// CalculateShipping calculates shipping cost
func (o *CheckoutOrchestrator) CalculateShipping(ctx context.Context, methodID string, weight decimal.Decimal, orderTotal decimal.Decimal, country, zipCode string) (decimal.Decimal, error) {
	// Get shipping method
	methodIDInt := int64(0)
	fmt.Sscanf(methodID, "%d", &methodIDInt)
	
	method, err := o.shippingMethodRepo.FindByID(ctx, methodIDInt)
	if err != nil || method == nil {
		return decimal.Zero, fmt.Errorf("shipping method not found")
	}

	if !method.IsEnabled {
		return decimal.Zero, fmt.Errorf("shipping method is disabled")
	}

	// Calculate base cost
	cost := method.CalculateRate(weight, orderTotal, 1)

	// Apply shipping rules (discounts/free shipping)
	rules, err := o.shippingRuleRepo.FindAllEnabled(ctx)
	if err != nil {
		return cost, nil // Return base cost if rules fail
	}

	for _, rule := range rules {
		if rule.AppliesTo(orderTotal, country, zipCode) {
			discount := rule.CalculateDiscount(cost)
			cost = cost.Sub(discount)
		}
	}

	if cost.LessThan(decimal.Zero) {
		cost = decimal.Zero
	}

	return cost, nil
}

// GetAvailableShippingMethods returns available shipping methods
func (o *CheckoutOrchestrator) GetAvailableShippingMethods(ctx context.Context, country, zipCode string) ([]*shipDomain.ShippingMethod, error) {
	methods, err := o.shippingMethodRepo.FindAllEnabled(ctx)
	if err != nil {
		return nil, err
	}

	// Filter by shipping rules restrictions
	rules, err := o.shippingRuleRepo.FindAllEnabled(ctx)
	if err != nil {
		return methods, nil // Return all if rules fail
	}

	available := make([]*shipDomain.ShippingMethod, 0)
	for _, method := range methods {
		allowed := true
		for _, rule := range rules {
			if rule.RuleType == shipDomain.RuleTypeRestrictionCountry || rule.RuleType == shipDomain.RuleTypeRestrictionZip {
				if !rule.AppliesTo(decimal.Zero, country, zipCode) {
					allowed = false
					break
				}
			}
		}
		if allowed {
			available = append(available, method)
		}
	}

	return available, nil
}

// OrderItem represents an item in the order
type OrderItem struct {
	SKUID       string
	SKU         string
	Quantity    int
	UnitPrice   decimal.Decimal
	Subtotal    decimal.Decimal
	TaxCategory string
	IsExempt    bool
}

// CalculateTax calculates taxes for the checkout
func (o *CheckoutOrchestrator) CalculateTax(ctx context.Context, customerID *string, orderItems []OrderItem, shippingAddress taxDomain.Address, shippingAmount decimal.Decimal) (*taxDomain.TaxCalculationResult, error) {
	// Find applicable jurisdictions
	jurisdictions, err := o.findApplicableJurisdictions(ctx, shippingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to find tax jurisdictions: %w", err)
	}

	if len(jurisdictions) == 0 {
		// No tax jurisdictions found, return zero tax
		result := taxDomain.NewTaxCalculationResult()
		return result, nil
	}

	// Get customer exemptions if customer ID provided
	var exemptions []*taxDomain.TaxExemption
	if customerID != nil {
		exemptions, err = o.taxExemptionRepo.FindActiveExemptions(ctx, *customerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get tax exemptions: %w", err)
		}
	}

	// Create tax calculation request
	request := taxDomain.NewTaxCalculationRequest(shippingAddress)
	request.CustomerID = customerID
	request.ShippingAmount = shippingAmount

	// Add items to request
	for _, item := range orderItems {
		taxableItem := taxDomain.TaxableItem{
			ItemID:      item.SKUID,
			SKU:         item.SKU,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Subtotal:    item.Subtotal,
			TaxCategory: taxDomain.TaxCategory(item.TaxCategory),
			IsExempt:    item.IsExempt,
		}
		request.AddItem(taxableItem)
	}

	// Calculate taxes
	result, err := o.calculateTaxWithExemptions(ctx, request, jurisdictions, exemptions)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate tax: %w", err)
	}

	return result, nil
}

// EstimateTax provides a quick tax estimate
func (o *CheckoutOrchestrator) EstimateTax(ctx context.Context, address taxDomain.Address, subtotal decimal.Decimal) (decimal.Decimal, error) {
	// Find applicable jurisdictions
	jurisdictions, err := o.findApplicableJurisdictions(ctx, address)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to find jurisdictions: %w", err)
	}

	if len(jurisdictions) == 0 {
		return decimal.Zero, nil
	}

	// Extract jurisdiction IDs
	jurisdictionIDs := make([]int64, len(jurisdictions))
	for i, j := range jurisdictions {
		jurisdictionIDs[i] = j.ID
	}

	// Find applicable tax rates for GENERAL category
	rates, err := o.taxRateRepo.FindApplicableRates(ctx, jurisdictionIDs, taxDomain.TaxCategoryGeneral, true)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to find tax rates: %w", err)
	}

	if len(rates) == 0 {
		return decimal.Zero, nil
	}

	// Calculate estimate
	totalTax := decimal.Zero
	for _, rate := range rates {
		if rate.AppliesTo(taxDomain.TaxCategoryGeneral, subtotal) {
			tax := rate.CalculateTax(subtotal, 1, decimal.Zero)
			totalTax = totalTax.Add(tax)
		}
	}

	return totalTax, nil
}

// Private helper methods

func (o *CheckoutOrchestrator) findApplicableJurisdictions(ctx context.Context, address taxDomain.Address) ([]*taxDomain.TaxJurisdiction, error) {
	jurisdictions, err := o.taxJurisdictionRepo.FindByLocation(
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
	active := make([]*taxDomain.TaxJurisdiction, 0)
	for _, j := range jurisdictions {
		if j.IsActive {
			active = append(active, j)
		}
	}

	return active, nil
}

func (o *CheckoutOrchestrator) calculateTaxWithExemptions(ctx context.Context, request *taxDomain.TaxCalculationRequest, jurisdictions []*taxDomain.TaxJurisdiction, exemptions []*taxDomain.TaxExemption) (*taxDomain.TaxCalculationResult, error) {
	// Extract jurisdiction IDs and create map
	jurisdictionIDs := make([]int64, len(jurisdictions))
	jurisdictionMap := make(map[int64]*taxDomain.TaxJurisdiction)
	for i, j := range jurisdictions {
		jurisdictionIDs[i] = j.ID
		jurisdictionMap[j.ID] = j
	}

	// Create result
	result := taxDomain.NewTaxCalculationResult()
	result.OrderID = request.OrderID

	// Calculate taxes for each item
	for _, item := range request.Items {
		taxedItem, err := o.calculateItemTax(ctx, item, jurisdictionIDs, jurisdictionMap, exemptions)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate tax for item %s: %w", item.ItemID, err)
		}
		result.AddTaxedItem(taxedItem)
	}

	// Calculate shipping tax if applicable
	if !request.ShippingAmount.IsZero() {
		shippingTax, err := o.calculateShippingTax(ctx, request.ShippingAmount, jurisdictionIDs, exemptions)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate shipping tax: %w", err)
		}
		result.SetShippingTax(shippingTax)
	}

	// Track jurisdictions used
	for _, j := range jurisdictions {
		result.JurisdictionsUsed = append(result.JurisdictionsUsed, j.Code)
	}

	// Finalize result
	result.Finalize()

	return result, nil
}

func (o *CheckoutOrchestrator) calculateItemTax(
	ctx context.Context,
	item taxDomain.TaxableItem,
	jurisdictionIDs []int64,
	jurisdictionMap map[int64]*taxDomain.TaxJurisdiction,
	exemptions []*taxDomain.TaxExemption,
) (taxDomain.TaxedItem, error) {
	taxedItem := taxDomain.TaxedItem{
		ItemID:      item.ItemID,
		SKU:         item.SKU,
		Quantity:    item.Quantity,
		UnitPrice:   item.UnitPrice,
		Subtotal:    item.Subtotal,
		TaxAmount:   decimal.Zero,
		TaxCategory: item.TaxCategory,
		Taxes:       make([]taxDomain.AppliedTax, 0),
	}

	// If item is exempt, return with zero tax
	if item.IsExempt {
		return taxedItem, nil
	}

	// Find applicable tax rates
	rates, err := o.taxRateRepo.FindApplicableRates(ctx, jurisdictionIDs, item.TaxCategory, true)
	if err != nil {
		return taxedItem, fmt.Errorf("failed to find tax rates: %w", err)
	}

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
		if o.isExempt(rate, exemptions) {
			continue
		}

		// Calculate tax for this rate
		taxableAmount := item.Subtotal
		if rate.IsCompound {
			taxableAmount = item.Subtotal.Add(cumulativeTax)
		}

		tax := rate.CalculateTax(taxableAmount, item.Quantity, cumulativeTax)
		cumulativeTax = cumulativeTax.Add(tax)

		// Get jurisdiction details
		jurisdiction := jurisdictionMap[rate.JurisdictionID]

		// Record applied tax
		appliedTax := taxDomain.AppliedTax{
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

func (o *CheckoutOrchestrator) calculateShippingTax(
	ctx context.Context,
	shippingAmount decimal.Decimal,
	jurisdictionIDs []int64,
	exemptions []*taxDomain.TaxExemption,
) (decimal.Decimal, error) {
	// Find applicable tax rates for SHIPPING category
	rates, err := o.taxRateRepo.FindApplicableRates(ctx, jurisdictionIDs, taxDomain.TaxCategoryShipping, true)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to find shipping tax rates: %w", err)
	}

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
		if o.isExempt(rate, exemptions) {
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

func (o *CheckoutOrchestrator) isExempt(rate *taxDomain.TaxRate, exemptions []*taxDomain.TaxExemption) bool {
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
