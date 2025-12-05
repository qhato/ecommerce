package commands

import (
	"time"

	"github.com/qhato/ecommerce/internal/pricing/domain"
	"github.com/shopspring/decimal"
)

// CreatePriceListCommand is a command to create a new price list
type CreatePriceListCommand struct {
	Name             string
	Code             string
	PriceListType    domain.PriceListType
	Currency         string
	Priority         int
	Description      string
	StartDate        *time.Time
	EndDate          *time.Time
	CustomerSegments []string
}

// UpdatePriceListCommand is a command to update an existing price list
type UpdatePriceListCommand struct {
	ID               int64
	Name             *string
	Priority         *int
	IsActive         *bool
	Description      *string
	StartDate        *time.Time
	EndDate          *time.Time
	CustomerSegments []string
}

// CreatePriceListItemCommand is a command to create a new price list item
type CreatePriceListItemCommand struct {
	PriceListID    int64
	SKUID          string
	ProductID      *string
	Price          decimal.Decimal
	CompareAtPrice *decimal.Decimal
	MinQuantity    int
	MaxQuantity    *int
	StartDate      *time.Time
	EndDate        *time.Time
}

// UpdatePriceListItemCommand is a command to update an existing price list item
type UpdatePriceListItemCommand struct {
	ID             int64
	Price          *decimal.Decimal
	CompareAtPrice *decimal.Decimal
	MinQuantity    *int
	MaxQuantity    *int
	IsActive       *bool
	StartDate      *time.Time
	EndDate        *time.Time
}

// BulkCreatePriceListItemsCommand is a command to create multiple price list items at once
type BulkCreatePriceListItemsCommand struct {
	PriceListID int64
	Items       []BulkPriceListItem
}

// BulkPriceListItem represents a single item in a bulk creation
type BulkPriceListItem struct {
	SKUID          string
	ProductID      *string
	Price          decimal.Decimal
	CompareAtPrice *decimal.Decimal
	MinQuantity    int
	MaxQuantity    *int
}

// CreatePricingRuleCommand is a command to create a new pricing rule
type CreatePricingRuleCommand struct {
	Name                 string
	Description          string
	RuleType             domain.PricingRuleType
	Priority             int
	ConditionExpression  string
	ActionType           domain.PricingRuleActionType
	ActionValue          decimal.Decimal
	ApplicableSKUs       []string
	ApplicableCategories []string
	CustomerSegments     []string
	MinQuantity          int
	MaxQuantity          *int
	MinOrderValue        *decimal.Decimal
	StartDate            *time.Time
	EndDate              *time.Time
}

// UpdatePricingRuleCommand is a command to update an existing pricing rule
type UpdatePricingRuleCommand struct {
	ID                   int64
	Name                 *string
	Description          *string
	Priority             *int
	IsActive             *bool
	ConditionExpression  *string
	ActionType           *domain.PricingRuleActionType
	ActionValue          *decimal.Decimal
	ApplicableSKUs       []string
	ApplicableCategories []string
	CustomerSegments     []string
	MinQuantity          *int
	MaxQuantity          *int
	MinOrderValue        *decimal.Decimal
	StartDate            *time.Time
	EndDate              *time.Time
}
