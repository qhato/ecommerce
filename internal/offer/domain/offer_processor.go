package domain

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// OfferQualification determines if an offer can be applied
type OfferQualification struct {
	Offer     *Offer
	Qualifies bool
	Reason    string
}

// OfferAdjustment represents a discount applied from an offer
type OfferAdjustment struct {
	OfferID        int64
	OfferName      string
	AdjustmentType OfferAdjustmentType
	Value          decimal.Decimal
	Applied        time.Time
}

// CandidateOffer represents an offer that might apply
type CandidateOffer struct {
	Offer           *Offer
	Priority        int
	DiscountAmount  decimal.Decimal
	QualifyingItems []string // Item IDs
	TargetItems     []string // Item IDs
}

// OfferContext holds the data needed for offer evaluation
type OfferContext struct {
	OrderTotal         decimal.Decimal
	OrderSubtotal      decimal.Decimal
	CustomerID         *string
	Items              []OfferItem
	AppliedOffers      []*OfferAdjustment
	AvailableOffers    []*Offer
	CustomerUsageCount map[int64]int // Offer ID -> usage count
}

// OfferItem represents an item for offer evaluation
type OfferItem struct {
	ItemID      string
	SKUID       string
	CategoryID  *string
	Price       decimal.Decimal
	SalePrice   *decimal.Decimal
	Quantity    int
	Subtotal    decimal.Decimal
	ProductID   *string
	Adjustments []OfferAdjustment
}

// GetEffectivePrice returns the price to use for calculations
func (item *OfferItem) GetEffectivePrice(applyToSalePrice bool) decimal.Decimal {
	if applyToSalePrice && item.SalePrice != nil {
		return *item.SalePrice
	}
	return item.Price
}

// GetAdjustedPrice returns the price after all adjustments
func (item *OfferItem) GetAdjustedPrice() decimal.Decimal {
	price := item.Subtotal
	for _, adj := range item.Adjustments {
		price = price.Sub(adj.Value)
	}
	return price
}

// OfferProcessor processes and applies offers to an order
type OfferProcessor struct {
	ruleEvaluator RuleEvaluator
}

// RuleEvaluator defines interface for evaluating offer rules
type RuleEvaluator interface {
	Evaluate(ruleExpression string, context map[string]interface{}) (bool, error)
}

// NewOfferProcessor creates a new OfferProcessor
func NewOfferProcessor(ruleEvaluator RuleEvaluator) *OfferProcessor {
	return &OfferProcessor{
		ruleEvaluator: ruleEvaluator,
	}
}

// QualifyOffer determines if an offer qualifies for the given context
func (p *OfferProcessor) QualifyOffer(offer *Offer, ctx *OfferContext) (*OfferQualification, error) {
	qualification := &OfferQualification{
		Offer:     offer,
		Qualifies: false,
		Reason:    "",
	}

	// Check if offer is active
	if offer.Archived {
		qualification.Reason = "Offer is archived"
		return qualification, nil
	}

	// Check date range
	now := time.Now()
	if now.Before(offer.StartDate) {
		qualification.Reason = "Offer has not started yet"
		return qualification, nil
	}
	if offer.EndDate != nil && now.After(*offer.EndDate) {
		qualification.Reason = "Offer has expired"
		return qualification, nil
	}

	// Check order minimum total
	if offer.OrderMinTotal > 0 {
		if ctx.OrderSubtotal.LessThan(decimal.NewFromFloat(offer.OrderMinTotal)) {
			qualification.Reason = fmt.Sprintf("Order subtotal below minimum of %.2f", offer.OrderMinTotal)
			return qualification, nil
		}
	}

	// Check max uses
	if offer.MaxUses != nil {
		// This would need to check against actual usage from repository
		// For now, we assume it's checked elsewhere
	}

	// Check max uses per customer
	if offer.MaxUsesPerCustomer != nil && ctx.CustomerID != nil {
		if usageCount, exists := ctx.CustomerUsageCount[offer.ID]; exists {
			if int64(usageCount) >= *offer.MaxUsesPerCustomer {
				qualification.Reason = "Customer has exceeded maximum uses for this offer"
				return qualification, nil
			}
		}
	}

	// Check if offer can be combined with already applied offers
	if !offer.CombinableWithOtherOffers && len(ctx.AppliedOffers) > 0 {
		qualification.Reason = "Offer cannot be combined with other offers"
		return qualification, nil
	}

	// Check qualifying item min total
	if offer.QualifyingItemMinTotal > 0 {
		qualifyingTotal := p.calculateQualifyingItemTotal(offer, ctx)
		if qualifyingTotal.LessThan(decimal.NewFromFloat(offer.QualifyingItemMinTotal)) {
			qualification.Reason = fmt.Sprintf("Qualifying items total below minimum of %.2f", offer.QualifyingItemMinTotal)
			return qualification, nil
		}
	}

	// Evaluate custom qualifier rule if present
	if offer.OfferItemQualifierRule != "" {
		qualifies, err := p.evaluateRule(offer.OfferItemQualifierRule, offer, ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate qualifier rule: %w", err)
		}
		if !qualifies {
			qualification.Reason = "Custom qualifier rule did not match"
			return qualification, nil
		}
	}

	// All checks passed
	qualification.Qualifies = true
	qualification.Reason = "Offer qualifies"
	return qualification, nil
}

// calculateQualifyingItemTotal calculates the total of items that qualify for the offer
func (p *OfferProcessor) calculateQualifyingItemTotal(offer *Offer, ctx *OfferContext) decimal.Decimal {
	total := decimal.Zero

	for _, item := range ctx.Items {
		// If there's a qualifier rule, evaluate it
		if offer.OfferItemQualifierRule != "" {
			itemCtx := map[string]interface{}{
				"item":  item,
				"order": ctx,
			}
			qualifies, err := p.ruleEvaluator.Evaluate(offer.OfferItemQualifierRule, itemCtx)
			if err != nil || !qualifies {
				continue
			}
		}

		effectivePrice := item.GetEffectivePrice(offer.ApplyToSalePrice)
		itemTotal := effectivePrice.Mul(decimal.NewFromInt(int64(item.Quantity)))
		total = total.Add(itemTotal)
	}

	return total
}

// CalculateDiscount calculates the discount amount for an offer
func (p *OfferProcessor) CalculateDiscount(offer *Offer, ctx *OfferContext) (decimal.Decimal, []string, error) {
	targetItems := p.findTargetItems(offer, ctx)
	if len(targetItems) == 0 {
		return decimal.Zero, nil, nil
	}

	var discountAmount decimal.Decimal
	targetItemIDs := make([]string, 0)

	switch offer.OfferDiscountType {
	case OfferDiscountTypePercentDiscount:
		// Apply percentage discount to target items
		targetTotal := decimal.Zero
		for _, item := range targetItems {
			effectivePrice := item.GetEffectivePrice(offer.ApplyToSalePrice)
			itemTotal := effectivePrice.Mul(decimal.NewFromInt(int64(item.Quantity)))
			targetTotal = targetTotal.Add(itemTotal)
			targetItemIDs = append(targetItemIDs, item.ItemID)
		}

		percentage := decimal.NewFromFloat(offer.OfferValue).Div(decimal.NewFromInt(100))
		discountAmount = targetTotal.Mul(percentage)

	case OfferDiscountTypeAmountOff:
		// Fixed amount off
		discountAmount = decimal.NewFromFloat(offer.OfferValue)
		for _, item := range targetItems {
			targetItemIDs = append(targetItemIDs, item.ItemID)
		}

	case OfferDiscountTypeFixPrice:
		// Set items to a fixed price
		targetTotal := decimal.Zero
		fixedPrice := decimal.NewFromFloat(offer.OfferValue)

		for _, item := range targetItems {
			currentPrice := item.GetEffectivePrice(offer.ApplyToSalePrice)
			if currentPrice.GreaterThan(fixedPrice) {
				priceDiff := currentPrice.Sub(fixedPrice)
				itemDiscount := priceDiff.Mul(decimal.NewFromInt(int64(item.Quantity)))
				targetTotal = targetTotal.Add(itemDiscount)
				targetItemIDs = append(targetItemIDs, item.ItemID)
			}
		}
		discountAmount = targetTotal

	default:
		return decimal.Zero, nil, fmt.Errorf("unsupported discount type: %s", offer.OfferDiscountType)
	}

	return discountAmount, targetItemIDs, nil
}

// findTargetItems finds items that are targets for the offer
func (p *OfferProcessor) findTargetItems(offer *Offer, ctx *OfferContext) []OfferItem {
	targetItems := make([]OfferItem, 0)

	for _, item := range ctx.Items {
		// If there's a target rule, evaluate it
		if offer.OfferItemTargetRule != "" {
			itemCtx := map[string]interface{}{
				"item":  item,
				"order": ctx,
			}
			isTarget, err := p.ruleEvaluator.Evaluate(offer.OfferItemTargetRule, itemCtx)
			if err != nil || !isTarget {
				continue
			}
		}

		targetItems = append(targetItems, item)
	}

	return targetItems
}

// evaluateRule evaluates a rule expression
func (p *OfferProcessor) evaluateRule(rule string, offer *Offer, ctx *OfferContext) (bool, error) {
	if p.ruleEvaluator == nil {
		return true, nil // No rule evaluator, pass by default
	}

	ruleCtx := map[string]interface{}{
		"offer": offer,
		"order": ctx,
	}

	return p.ruleEvaluator.Evaluate(rule, ruleCtx)
}

// ApplyOffer creates an adjustment from an offer
func (p *OfferProcessor) ApplyOffer(offer *Offer, discountAmount decimal.Decimal) *OfferAdjustment {
	return &OfferAdjustment{
		OfferID:        offer.ID,
		OfferName:      offer.Name,
		AdjustmentType: offer.AdjustmentType,
		Value:          discountAmount,
		Applied:        time.Now(),
	}
}

// SelectBestOffers selects the best combination of offers to apply
func (p *OfferProcessor) SelectBestOffers(candidates []*CandidateOffer) []*CandidateOffer {
	if len(candidates) == 0 {
		return candidates
	}

	// Sort by priority (lower number = higher priority), then by discount amount (descending)
	sortedCandidates := make([]*CandidateOffer, len(candidates))
	copy(sortedCandidates, candidates)

	// Simple selection: take highest priority and highest discount
	// In a real system, this would be more sophisticated (e.g., maximize total discount)
	selected := make([]*CandidateOffer, 0)

	for _, candidate := range sortedCandidates {
		canCombine := true

		// Check if this offer can be combined with already selected offers
		for _, selectedOffer := range selected {
			if !candidate.Offer.CombinableWithOtherOffers || !selectedOffer.Offer.CombinableWithOtherOffers {
				canCombine = false
				break
			}
		}

		if canCombine {
			selected = append(selected, candidate)
		}
	}

	return selected
}
