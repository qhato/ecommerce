package queries

import (
	"context"
	"fmt"
	"sort"

	"github.com/qhato/ecommerce/internal/pricing/domain"
	"github.com/shopspring/decimal"
)

// PricingQueryService handles pricing queries and calculations
type PricingQueryService struct {
	priceListRepo     domain.PriceListRepository
	priceListItemRepo domain.PriceListItemRepository
	pricingRuleRepo   domain.PricingRuleRepository
}

// NewPricingQueryService creates a new PricingQueryService
func NewPricingQueryService(
	priceListRepo domain.PriceListRepository,
	priceListItemRepo domain.PriceListItemRepository,
	pricingRuleRepo domain.PricingRuleRepository,
) *PricingQueryService {
	return &PricingQueryService{
		priceListRepo:     priceListRepo,
		priceListItemRepo: priceListItemRepo,
		pricingRuleRepo:   pricingRuleRepo,
	}
}

// CalculatePrices calculates prices for all items in the context
func (s *PricingQueryService) CalculatePrices(ctx context.Context, pricingCtx *domain.PricingContext) (*domain.PricingResult, error) {
	if pricingCtx == nil {
		return nil, domain.ErrPricingContextRequired
	}

	result := domain.NewPricingResult(pricingCtx.Currency)

	// Determine effective price list
	priceList, err := s.getEffectivePriceList(ctx, pricingCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get effective price list: %w", err)
	}

	// Get active pricing rules
	pricingRules, err := s.pricingRuleRepo.FindActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pricing rules: %w", err)
	}

	// Sort rules by priority (higher priority first)
	sort.Slice(pricingRules, func(i, j int) bool {
		return pricingRules[i].Priority > pricingRules[j].Priority
	})

	// Calculate price for each SKU
	for _, request := range pricingCtx.RequestedSKUs {
		pricedItem, err := s.calculatePriceForSKU(ctx, request, priceList, pricingRules, pricingCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate price for SKU %s: %w", request.SKUID, err)
		}
		result.AddItem(pricedItem)
	}

	return result, nil
}

// GetPriceForSKU gets the price for a single SKU
func (s *PricingQueryService) GetPriceForSKU(ctx context.Context, skuID string, quantity int, currency string, customerSegment *string) (*domain.PricedItem, error) {
	// Create pricing context
	pricingCtx := domain.NewPricingContext(currency)
	pricingCtx.CustomerSegment = customerSegment
	pricingCtx.AddSKU(skuID, quantity)

	// Calculate prices
	result, err := s.CalculatePrices(ctx, pricingCtx)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, domain.ErrNoPriceFound
	}

	return result.Items[0], nil
}

// GetEffectivePriceList determines which price list to use for a customer
func (s *PricingQueryService) GetEffectivePriceList(ctx context.Context, currency string, customerSegment *string) (*domain.PriceList, error) {
	pricingCtx := &domain.PricingContext{
		Currency:        currency,
		CustomerSegment: customerSegment,
	}
	return s.getEffectivePriceList(ctx, pricingCtx)
}

// getEffectivePriceList internal method to determine effective price list
func (s *PricingQueryService) getEffectivePriceList(ctx context.Context, pricingCtx *domain.PricingContext) (*domain.PriceList, error) {
	// Get all active price lists by priority
	priceLists, err := s.priceListRepo.FindByPriority(ctx, pricingCtx.Currency)
	if err != nil {
		return nil, fmt.Errorf("failed to find price lists: %w", err)
	}

	if len(priceLists) == 0 {
		return nil, domain.ErrNoActivePriceList
	}

	// Find the highest priority price list that applies to this customer
	for _, priceList := range priceLists {
		if !priceList.IsCurrentlyActive() {
			continue
		}

		// Check if price list applies to customer segment
		if pricingCtx.CustomerSegment != nil {
			if priceList.AppliesTo(*pricingCtx.CustomerSegment) {
				return priceList, nil
			}
		} else {
			// If no customer segment, use first active standard price list
			if priceList.PriceListType == domain.PriceListTypeStandard {
				return priceList, nil
			}
		}
	}

	// If no specific match, return first active price list
	for _, priceList := range priceLists {
		if priceList.IsCurrentlyActive() {
			return priceList, nil
		}
	}

	return nil, domain.ErrNoActivePriceList
}

// calculatePriceForSKU calculates price for a single SKU
func (s *PricingQueryService) calculatePriceForSKU(
	ctx context.Context,
	request domain.PricingRequest,
	priceList *domain.PriceList,
	pricingRules []*domain.PricingRule,
	pricingCtx *domain.PricingContext,
) (*domain.PricedItem, error) {
	// Get price from price list
	priceListItem, err := s.priceListItemRepo.FindBySKUAndPriceList(ctx, request.SKUID, priceList.ID)
	if err != nil || priceListItem == nil {
		return nil, fmt.Errorf("no price found for SKU %s in price list %s: %w", request.SKUID, priceList.Code, domain.ErrNoPriceFound)
	}

	// Check if item is active and applies to quantity
	if !priceListItem.IsCurrentlyActive() {
		return nil, fmt.Errorf("price for SKU %s is not currently active", request.SKUID)
	}

	if !priceListItem.AppliesTo(request.Quantity) {
		return nil, fmt.Errorf("price for SKU %s does not apply to quantity %d", request.SKUID, request.Quantity)
	}

	// Create priced item
	pricedItem := &domain.PricedItem{
		SKUID:          request.SKUID,
		ProductID:      request.ProductID,
		Quantity:       request.Quantity,
		BasePrice:      priceListItem.Price,
		FinalPrice:     priceListItem.Price,
		CompareAtPrice: priceListItem.CompareAtPrice,
		PriceListID:    &priceList.ID,
		PriceListName:  &priceList.Name,
		DiscountAmount: decimal.Zero,
		DiscountPercent: decimal.Zero,
		Currency:       pricingCtx.Currency,
		Adjustments:    make([]domain.PriceAdjustment, 0),
		IsOnSale:       priceListItem.CompareAtPrice != nil,
	}

	// Apply pricing rules
	for _, rule := range pricingRules {
		if !rule.IsCurrentlyActive() {
			continue
		}

		// Check if rule applies
		if !rule.AppliesTo(request.SKUID, request.Quantity, pricingCtx.CustomerSegment, nil) {
			continue
		}

		// Calculate adjustment
		adjustmentAmount := rule.CalculateAdjustment(pricedItem.BasePrice)
		if adjustmentAmount.IsZero() {
			continue
		}

		// Create adjustment
		adjustment := domain.PriceAdjustment{
			Type:        s.mapRuleActionToAdjustmentType(rule.ActionType),
			Amount:      adjustmentAmount.Abs(),
			Reason:      string(rule.RuleType),
			Description: rule.Name,
			Priority:    rule.Priority,
		}

		pricedItem.AddAdjustment(adjustment)
	}

	// Calculate final price
	pricedItem.CalculateFinalPrice()

	return pricedItem, nil
}

// mapRuleActionToAdjustmentType maps rule action type to adjustment type
func (s *PricingQueryService) mapRuleActionToAdjustmentType(actionType domain.PricingRuleActionType) domain.PriceAdjustmentType {
	switch actionType {
	case domain.PricingRuleActionTypePercentDiscount, domain.PricingRuleActionTypeAmountDiscount, domain.PricingRuleActionTypeFixedPrice:
		return domain.PriceAdjustmentTypeDiscount
	case domain.PricingRuleActionTypePercentSurcharge, domain.PricingRuleActionTypeAmountSurcharge:
		return domain.PriceAdjustmentTypeSurcharge
	default:
		return domain.PriceAdjustmentTypeDiscount
	}
}

// GetPriceList retrieves a price list by ID
func (s *PricingQueryService) GetPriceList(ctx context.Context, id int64) (*domain.PriceList, error) {
	priceList, err := s.priceListRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find price list: %w", err)
	}
	if priceList == nil {
		return nil, domain.ErrPriceListNotFound
	}
	return priceList, nil
}

// GetPriceListByCode retrieves a price list by code
func (s *PricingQueryService) GetPriceListByCode(ctx context.Context, code string) (*domain.PriceList, error) {
	priceList, err := s.priceListRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to find price list: %w", err)
	}
	if priceList == nil {
		return nil, domain.ErrPriceListNotFound
	}
	return priceList, nil
}

// GetActivePriceLists retrieves all active price lists for a currency
func (s *PricingQueryService) GetActivePriceLists(ctx context.Context, currency string) ([]*domain.PriceList, error) {
	priceLists, err := s.priceListRepo.FindActive(ctx, currency)
	if err != nil {
		return nil, fmt.Errorf("failed to find active price lists: %w", err)
	}
	return priceLists, nil
}

// GetPriceListItems retrieves all items in a price list
func (s *PricingQueryService) GetPriceListItems(ctx context.Context, priceListID int64) ([]*domain.PriceListItem, error) {
	items, err := s.priceListItemRepo.FindByPriceListID(ctx, priceListID)
	if err != nil {
		return nil, fmt.Errorf("failed to find price list items: %w", err)
	}
	return items, nil
}

// GetPriceListItem retrieves a price list item by ID
func (s *PricingQueryService) GetPriceListItem(ctx context.Context, id int64) (*domain.PriceListItem, error) {
	item, err := s.priceListItemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find price list item: %w", err)
	}
	if item == nil {
		return nil, domain.ErrPriceListItemNotFound
	}
	return item, nil
}

// GetPricingRule retrieves a pricing rule by ID
func (s *PricingQueryService) GetPricingRule(ctx context.Context, id int64) (*domain.PricingRule, error) {
	rule, err := s.pricingRuleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find pricing rule: %w", err)
	}
	if rule == nil {
		return nil, domain.ErrPricingRuleNotFound
	}
	return rule, nil
}

// GetActivePricingRules retrieves all active pricing rules
func (s *PricingQueryService) GetActivePricingRules(ctx context.Context) ([]*domain.PricingRule, error) {
	rules, err := s.pricingRuleRepo.FindActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find active pricing rules: %w", err)
	}
	return rules, nil
}
