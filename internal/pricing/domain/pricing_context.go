package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// PricingContext holds all the information needed to determine prices
// Business Logic: Contexto completo para evaluación de precios (customer, currency, date, quantity)
type PricingContext struct {
	CustomerID      *string
	CustomerSegment *string
	Currency        string
	Locale          string
	PriceDate       time.Time // Date to use for price evaluation
	RequestedSKUs   []PricingRequest
}

// PricingRequest represents a request for pricing a specific SKU
type PricingRequest struct {
	SKUID      string
	ProductID  *string
	Quantity   int
	Attributes map[string]string // Additional attributes for dynamic pricing
}

// PricedItem represents the result of pricing a single item
type PricedItem struct {
	SKUID            string
	ProductID        *string
	Quantity         int
	BasePrice        decimal.Decimal // Original list price
	SalePrice        *decimal.Decimal // Sale price if applicable
	FinalPrice       decimal.Decimal // Final price after all adjustments
	CompareAtPrice   *decimal.Decimal // Price to compare against for showing savings
	PriceListID      *int64 // ID of price list used
	PriceListName    *string // Name of price list used
	DiscountAmount   decimal.Decimal // Total discount amount
	DiscountPercent  decimal.Decimal // Discount percentage
	Subtotal         decimal.Decimal // FinalPrice * Quantity
	Currency         string
	Adjustments      []PriceAdjustment
	IsOnSale         bool
	AvailableFrom    *time.Time
	AvailableUntil   *time.Time
}

// PriceAdjustment represents a price adjustment applied to an item
type PriceAdjustment struct {
	Type        PriceAdjustmentType
	Amount      decimal.Decimal
	Reason      string
	Description string
	Priority    int
}

// PriceAdjustmentType defines the type of price adjustment
type PriceAdjustmentType string

const (
	PriceAdjustmentTypeDiscount       PriceAdjustmentType = "DISCOUNT"
	PriceAdjustmentTypeSurcharge      PriceAdjustmentType = "SURCHARGE"
	PriceAdjustmentTypeQuantityTier   PriceAdjustmentType = "QUANTITY_TIER"
	PriceAdjustmentTypeCustomerGroup  PriceAdjustmentType = "CUSTOMER_GROUP"
	PriceAdjustmentTypePromotional    PriceAdjustmentType = "PROMOTIONAL"
	PriceAdjustmentTypeDynamic        PriceAdjustmentType = "DYNAMIC"
)

// NewPricingContext creates a new pricing context
func NewPricingContext(currency string) *PricingContext {
	return &PricingContext{
		Currency:      currency,
		PriceDate:     time.Now(),
		RequestedSKUs: make([]PricingRequest, 0),
	}
}

// AddSKU adds a SKU to be priced
func (pc *PricingContext) AddSKU(skuID string, quantity int) {
	pc.RequestedSKUs = append(pc.RequestedSKUs, PricingRequest{
		SKUID:      skuID,
		Quantity:   quantity,
		Attributes: make(map[string]string),
	})
}

// SetCustomer sets the customer information for personalized pricing
func (pc *PricingContext) SetCustomer(customerID string, segment *string) {
	pc.CustomerID = &customerID
	pc.CustomerSegment = segment
}

// SetPriceDate sets the date to use for price evaluation
func (pc *PricingContext) SetPriceDate(date time.Time) {
	pc.PriceDate = date
}

// GetSavings calculates the total savings from compare at price
func (pi *PricedItem) GetSavings() decimal.Decimal {
	if pi.CompareAtPrice == nil {
		return decimal.Zero
	}

	if pi.CompareAtPrice.LessThanOrEqual(pi.FinalPrice) {
		return decimal.Zero
	}

	return pi.CompareAtPrice.Sub(pi.FinalPrice).Mul(decimal.NewFromInt(int64(pi.Quantity)))
}

// AddAdjustment adds a price adjustment to the item
func (pi *PricedItem) AddAdjustment(adjustment PriceAdjustment) {
	pi.Adjustments = append(pi.Adjustments, adjustment)

	// Update discount amount
	if adjustment.Type == PriceAdjustmentTypeDiscount {
		pi.DiscountAmount = pi.DiscountAmount.Add(adjustment.Amount)
	}
}

// CalculateFinalPrice calculates the final price after all adjustments
func (pi *PricedItem) CalculateFinalPrice() {
	finalPrice := pi.BasePrice

	// Apply sale price if exists
	if pi.SalePrice != nil {
		finalPrice = *pi.SalePrice
		pi.IsOnSale = true
	}

	// Apply adjustments
	for _, adj := range pi.Adjustments {
		switch adj.Type {
		case PriceAdjustmentTypeDiscount:
			finalPrice = finalPrice.Sub(adj.Amount)
		case PriceAdjustmentTypeSurcharge:
			finalPrice = finalPrice.Add(adj.Amount)
		}
	}

	// Ensure price doesn't go below zero
	if finalPrice.LessThan(decimal.Zero) {
		finalPrice = decimal.Zero
	}

	pi.FinalPrice = finalPrice
	pi.Subtotal = finalPrice.Mul(decimal.NewFromInt(int64(pi.Quantity)))

	// Calculate discount percentage
	if pi.BasePrice.GreaterThan(decimal.Zero) && pi.DiscountAmount.GreaterThan(decimal.Zero) {
		pi.DiscountPercent = pi.DiscountAmount.Div(pi.BasePrice).Mul(decimal.NewFromInt(100))
	}
}

// PricingResult represents the complete result of a pricing request
type PricingResult struct {
	Items       []*PricedItem
	Currency    string
	TotalAmount decimal.Decimal
	PricedAt    time.Time
}

// NewPricingResult creates a new pricing result
func NewPricingResult(currency string) *PricingResult {
	return &PricingResult{
		Items:       make([]*PricedItem, 0),
		Currency:    currency,
		TotalAmount: decimal.Zero,
		PricedAt:    time.Now(),
	}
}

// AddItem adds a priced item to the result
func (pr *PricingResult) AddItem(item *PricedItem) {
	pr.Items = append(pr.Items, item)
	pr.TotalAmount = pr.TotalAmount.Add(item.Subtotal)
}

// GetTotalSavings calculates total savings across all items
func (pr *PricingResult) GetTotalSavings() decimal.Decimal {
	totalSavings := decimal.Zero
	for _, item := range pr.Items {
		totalSavings = totalSavings.Add(item.GetSavings())
	}
	return totalSavings
}
