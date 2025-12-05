package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// PriceListCreatedEvent is published when a new price list is created
type PriceListCreatedEvent struct {
	PriceListID   int64
	Code          string
	Name          string
	Currency      string
	PriceListType PriceListType
	CreatedAt     time.Time
}

// PriceListActivatedEvent is published when a price list is activated
type PriceListActivatedEvent struct {
	PriceListID  int64
	Code         string
	ActivatedAt  time.Time
}

// PriceListDeactivatedEvent is published when a price list is deactivated
type PriceListDeactivatedEvent struct {
	PriceListID    int64
	Code           string
	DeactivatedAt  time.Time
}

// PriceListItemCreatedEvent is published when a price list item is created
type PriceListItemCreatedEvent struct {
	PriceListItemID int64
	PriceListID     int64
	SKUID           string
	Price           decimal.Decimal
	MinQuantity     int
	CreatedAt       time.Time
}

// PriceListItemUpdatedEvent is published when a price list item is updated
type PriceListItemUpdatedEvent struct {
	PriceListItemID int64
	PriceListID     int64
	SKUID           string
	OldPrice        decimal.Decimal
	NewPrice        decimal.Decimal
	UpdatedAt       time.Time
}

// PriceCalculatedEvent is published when prices are calculated
type PriceCalculatedEvent struct {
	SKUID          string
	Quantity       int
	BasePrice      decimal.Decimal
	FinalPrice     decimal.Decimal
	DiscountAmount decimal.Decimal
	PriceListID    *int64
	CustomerID     *string
	CalculatedAt   time.Time
}

// PricingRuleCreatedEvent is published when a pricing rule is created
type PricingRuleCreatedEvent struct {
	RuleID    int64
	Name      string
	RuleType  PricingRuleType
	Priority  int
	CreatedAt time.Time
}

// PricingRuleAppliedEvent is published when a pricing rule is applied
type PricingRuleAppliedEvent struct {
	RuleID      int64
	RuleName    string
	SKUID       string
	Adjustment  decimal.Decimal
	AppliedAt   time.Time
}
