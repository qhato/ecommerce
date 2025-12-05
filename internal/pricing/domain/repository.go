package domain

import "context"

// PriceListRepository provides an interface for managing PriceLists
type PriceListRepository interface {
	// Save stores a new price list or updates an existing one
	Save(ctx context.Context, priceList *PriceList) error

	// FindByID retrieves a price list by its unique identifier
	FindByID(ctx context.Context, id int64) (*PriceList, error)

	// FindByCode retrieves a price list by its code
	FindByCode(ctx context.Context, code string) (*PriceList, error)

	// FindActive retrieves all currently active price lists
	FindActive(ctx context.Context, currency string) ([]*PriceList, error)

	// FindByPriority retrieves active price lists ordered by priority (highest first)
	FindByPriority(ctx context.Context, currency string) ([]*PriceList, error)

	// FindByCustomerSegment retrieves price lists for a customer segment
	FindByCustomerSegment(ctx context.Context, segment string, currency string) ([]*PriceList, error)

	// Delete removes a price list by its unique identifier
	Delete(ctx context.Context, id int64) error
}

// PriceListItemRepository provides an interface for managing PriceListItems
type PriceListItemRepository interface {
	// Save stores a new price list item or updates an existing one
	Save(ctx context.Context, item *PriceListItem) error

	// FindByID retrieves a price list item by its unique identifier
	FindByID(ctx context.Context, id int64) (*PriceListItem, error)

	// FindByPriceListID retrieves all items in a price list
	FindByPriceListID(ctx context.Context, priceListID int64) ([]*PriceListItem, error)

	// FindBySKU retrieves all price list items for a SKU across all price lists
	FindBySKU(ctx context.Context, skuID string) ([]*PriceListItem, error)

	// FindBySKUAndPriceList retrieves a price list item for a specific SKU in a specific price list
	FindBySKUAndPriceList(ctx context.Context, skuID string, priceListID int64) (*PriceListItem, error)

	// FindActiveForSKU retrieves all currently active price list items for a SKU
	FindActiveForSKU(ctx context.Context, skuID string, quantity int) ([]*PriceListItem, error)

	// Delete removes a price list item by its unique identifier
	Delete(ctx context.Context, id int64) error

	// DeleteByPriceListID removes all items in a price list
	DeleteByPriceListID(ctx context.Context, priceListID int64) error
}

// PricingRuleRepository provides an interface for managing PricingRules
type PricingRuleRepository interface {
	// Save stores a new pricing rule or updates an existing one
	Save(ctx context.Context, rule *PricingRule) error

	// FindByID retrieves a pricing rule by its unique identifier
	FindByID(ctx context.Context, id int64) (*PricingRule, error)

	// FindActive retrieves all currently active pricing rules ordered by priority
	FindActive(ctx context.Context) ([]*PricingRule, error)

	// FindBySKU retrieves all pricing rules applicable to a SKU
	FindBySKU(ctx context.Context, skuID string) ([]*PricingRule, error)

	// FindByCustomerSegment retrieves pricing rules for a customer segment
	FindByCustomerSegment(ctx context.Context, segment string) ([]*PricingRule, error)

	// Delete removes a pricing rule by its unique identifier
	Delete(ctx context.Context, id int64) error
}

// PricingService defines the interface for pricing operations
type PricingService interface {
	// CalculatePrices calculates prices for all items in the context
	CalculatePrices(ctx context.Context, pricingContext *PricingContext) (*PricingResult, error)

	// GetPriceForSKU gets the price for a single SKU
	GetPriceForSKU(ctx context.Context, skuID string, quantity int, currency string, customerSegment *string) (*PricedItem, error)

	// GetEffectivePriceList determines which price list to use for a customer
	GetEffectivePriceList(ctx context.Context, currency string, customerSegment *string) (*PriceList, error)
}
