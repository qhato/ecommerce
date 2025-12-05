package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// PricingRuleType defines the type of pricing rule
type PricingRuleType string

const (
	PricingRuleTypeQuantityTiered   PricingRuleType = "QUANTITY_TIERED"   // Buy more, pay less per unit
	PricingRuleTypeVolumeDiscount   PricingRuleType = "VOLUME_DISCOUNT"   // Discount based on total volume
	PricingRuleTypeCustomerSegment  PricingRuleType = "CUSTOMER_SEGMENT"  // Pricing by customer segment
	PricingRuleTypeDynamic          PricingRuleType = "DYNAMIC"           // Dynamic pricing based on rules
	PricingRuleTypeTimeBasedDiscount PricingRuleType = "TIME_BASED"       // Happy hour, seasonal pricing
)

// PricingRuleActionType defines what action to take
type PricingRuleActionType string

const (
	PricingRuleActionTypeFixedPrice      PricingRuleActionType = "FIXED_PRICE"       // Set to fixed price
	PricingRuleActionTypePercentDiscount PricingRuleActionType = "PERCENT_DISCOUNT"  // Percentage discount
	PricingRuleActionTypeAmountDiscount  PricingRuleActionType = "AMOUNT_DISCOUNT"   // Fixed amount discount
	PricingRuleActionTypePercentSurcharge PricingRuleActionType = "PERCENT_SURCHARGE" // Percentage surcharge
	PricingRuleActionTypeAmountSurcharge PricingRuleActionType = "AMOUNT_SURCHARGE"  // Fixed amount surcharge
)

// PricingRule represents a rule for automatic price adjustments
// Business Logic: Reglas para ajustes automáticos de precio (descuentos por volumen, tiempo, segmento)
type PricingRule struct {
	ID              int64
	Name            string
	Description     string
	RuleType        PricingRuleType
	Priority        int // Higher priority rules are evaluated first
	IsActive        bool
	StartDate       *time.Time
	EndDate         *time.Time
	ConditionExpression string // Rule expression to evaluate
	ActionType      PricingRuleActionType
	ActionValue     decimal.Decimal
	ApplicableSKUs  []string // Empty means applies to all
	ApplicableCategories []string // Empty means applies to all
	CustomerSegments []string // Empty means applies to all customers
	MinQuantity     int
	MaxQuantity     *int
	MinOrderValue   *decimal.Decimal
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewPricingRule creates a new pricing rule
func NewPricingRule(name string, ruleType PricingRuleType, priority int) (*PricingRule, error) {
	if name == "" {
		return nil, ErrPricingRuleNameRequired
	}

	now := time.Now()
	return &PricingRule{
		Name:                 name,
		RuleType:             ruleType,
		Priority:             priority,
		IsActive:             true,
		ApplicableSKUs:       make([]string, 0),
		ApplicableCategories: make([]string, 0),
		CustomerSegments:     make([]string, 0),
		MinQuantity:          1,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

// IsCurrentlyActive checks if the rule is currently active
func (pr *PricingRule) IsCurrentlyActive() bool {
	if !pr.IsActive {
		return false
	}

	now := time.Now()
	if pr.StartDate != nil && now.Before(*pr.StartDate) {
		return false
	}
	if pr.EndDate != nil && now.After(*pr.EndDate) {
		return false
	}

	return true
}

// AppliesTo checks if this rule applies to a given context
func (pr *PricingRule) AppliesTo(skuID string, quantity int, customerSegment *string, orderValue *decimal.Decimal) bool {
	// Check SKU
	if len(pr.ApplicableSKUs) > 0 {
		found := false
		for _, sku := range pr.ApplicableSKUs {
			if sku == skuID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check quantity range
	if quantity < pr.MinQuantity {
		return false
	}
	if pr.MaxQuantity != nil && quantity > *pr.MaxQuantity {
		return false
	}

	// Check customer segment
	if len(pr.CustomerSegments) > 0 && customerSegment != nil {
		found := false
		for _, segment := range pr.CustomerSegments {
			if segment == *customerSegment {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check minimum order value
	if pr.MinOrderValue != nil && orderValue != nil {
		if orderValue.LessThan(*pr.MinOrderValue) {
			return false
		}
	}

	return true
}

// CalculateAdjustment calculates the price adjustment based on the rule
func (pr *PricingRule) CalculateAdjustment(basePrice decimal.Decimal) decimal.Decimal {
	switch pr.ActionType {
	case PricingRuleActionTypeFixedPrice:
		// Return the difference to reach the fixed price
		return basePrice.Sub(pr.ActionValue)

	case PricingRuleActionTypePercentDiscount:
		// Calculate percentage discount
		percentage := pr.ActionValue.Div(decimal.NewFromInt(100))
		return basePrice.Mul(percentage)

	case PricingRuleActionTypeAmountDiscount:
		// Fixed amount discount
		return pr.ActionValue

	case PricingRuleActionTypePercentSurcharge:
		// Calculate percentage surcharge (negative adjustment)
		percentage := pr.ActionValue.Div(decimal.NewFromInt(100))
		return basePrice.Mul(percentage).Neg()

	case PricingRuleActionTypeAmountSurcharge:
		// Fixed amount surcharge (negative adjustment)
		return pr.ActionValue.Neg()

	default:
		return decimal.Zero
	}
}

// SetAction sets the action type and value
func (pr *PricingRule) SetAction(actionType PricingRuleActionType, value decimal.Decimal) {
	pr.ActionType = actionType
	pr.ActionValue = value
	pr.UpdatedAt = time.Now()
}

// SetQuantityRange sets the quantity range for the rule
func (pr *PricingRule) SetQuantityRange(minQuantity int, maxQuantity *int) error {
	if minQuantity < 0 {
		return ErrMinQuantityCannotBeNegative
	}
	if maxQuantity != nil && *maxQuantity < minQuantity {
		return ErrMaxQuantityLessThanMin
	}

	pr.MinQuantity = minQuantity
	pr.MaxQuantity = maxQuantity
	pr.UpdatedAt = time.Now()
	return nil
}

// AddApplicableSKU adds a SKU that this rule applies to
func (pr *PricingRule) AddApplicableSKU(skuID string) {
	for _, sku := range pr.ApplicableSKUs {
		if sku == skuID {
			return // Already exists
		}
	}
	pr.ApplicableSKUs = append(pr.ApplicableSKUs, skuID)
	pr.UpdatedAt = time.Now()
}

// AddCustomerSegment adds a customer segment this rule applies to
func (pr *PricingRule) AddCustomerSegment(segment string) {
	for _, s := range pr.CustomerSegments {
		if s == segment {
			return // Already exists
		}
	}
	pr.CustomerSegments = append(pr.CustomerSegments, segment)
	pr.UpdatedAt = time.Now()
}

// Activate activates the pricing rule
func (pr *PricingRule) Activate() {
	pr.IsActive = true
	pr.UpdatedAt = time.Now()
}

// Deactivate deactivates the pricing rule
func (pr *PricingRule) Deactivate() {
	pr.IsActive = false
	pr.UpdatedAt = time.Now()
}
