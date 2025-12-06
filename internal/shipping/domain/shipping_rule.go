package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// ShippingRule represents a shipping rule (free shipping, restrictions)
type ShippingRule struct {
	ID             int64
	Name           string
	Description    string
	RuleType       RuleType
	IsEnabled      bool
	Priority       int
	MinOrderValue  *decimal.Decimal
	Countries      []string
	ExcludedZips   []string
	DiscountType   DiscountType
	DiscountValue  decimal.Decimal
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RuleType string

const (
	RuleTypeFreeShipping       RuleType = "FREE_SHIPPING"
	RuleTypeDiscountShipping   RuleType = "DISCOUNT_SHIPPING"
	RuleTypeRestrictionCountry RuleType = "RESTRICTION_COUNTRY"
	RuleTypeRestrictionZip     RuleType = "RESTRICTION_ZIP"
)

type DiscountType string

const (
	DiscountTypeFixed      DiscountType = "FIXED"
	DiscountTypePercentage DiscountType = "PERCENTAGE"
)

// NewShippingRule creates a new shipping rule
func NewShippingRule(name string, ruleType RuleType) *ShippingRule {
	now := time.Now()
	return &ShippingRule{
		Name:          name,
		RuleType:      ruleType,
		IsEnabled:     false,
		Priority:      0,
		Countries:     make([]string, 0),
		ExcludedZips:  make([]string, 0),
		DiscountValue: decimal.Zero,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// AppliesTo checks if the rule applies to the given order
func (r *ShippingRule) AppliesTo(orderTotal decimal.Decimal, country, zip string) bool {
	if !r.IsEnabled {
		return false
	}

	if r.MinOrderValue != nil && orderTotal.LessThan(*r.MinOrderValue) {
		return false
	}

	if r.RuleType == RuleTypeRestrictionCountry && len(r.Countries) > 0 {
		allowed := false
		for _, c := range r.Countries {
			if c == country {
				allowed = true
				break
			}
		}
		return allowed
	}

	if r.RuleType == RuleTypeRestrictionZip && len(r.ExcludedZips) > 0 {
		for _, z := range r.ExcludedZips {
			if z == zip {
				return false
			}
		}
	}

	return true
}

// CalculateDiscount calculates the discount amount
func (r *ShippingRule) CalculateDiscount(shippingCost decimal.Decimal) decimal.Decimal {
	if r.RuleType == RuleTypeFreeShipping {
		return shippingCost
	}

	if r.RuleType == RuleTypeDiscountShipping {
		if r.DiscountType == DiscountTypePercentage {
			return shippingCost.Mul(r.DiscountValue).Div(decimal.NewFromInt(100))
		}
		return r.DiscountValue
	}

	return decimal.Zero
}
