package domain

import "time"

// OrderItem represents an item in an order, capturing SKU details at the time of order
type OrderItem struct {
	ID                  int64
	OrderID             int64
	SKUID               int64
	ProductID           int64  // Reference to the product it belongs to
	Name                string // Product name at the time of order
	Quantity            int
	RetailPrice         float64 // Original retail price of the SKU
	SalePrice           float64 // Sale price of the SKU (if any) at the time of order
	Price               float64 // The actual price charged for the item (after item-level discounts)
	TotalPrice          float64 // Price * Quantity (after item-level discounts)
	TaxAmount           float64
	TaxCategory         string  // New: For tax calculations at the item level
	ShippingAmount      float64 // From blc_order_item.shipping_amount (not directly in blc_order_item, but often related)
	DiscountsAllowed    bool    // From blc_order_item.discounts_allowed
	HasValidationErrors bool    // From blc_order_item.has_validation_errors
	ItemTaxableFlag     bool    // From blc_order_item.item_taxable_flag
	OrderItemType       string  // From blc_order_item.order_item_type
	RetailPriceOverride bool    // From blc_order_item.retail_price_override
	SalePriceOverride   bool    // From blc_order_item.sale_price_override

	CategoryID        *int64 // From blc_order_item.category_id
	GiftWrapItemID    *int64 // From blc_order_item.gift_wrap_item_id
	ParentOrderItemID *int64 // From blc_order_item.parent_order_item_id
	PersonalMessageID *int64 // From blc_order_item.personal_message_id

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewOrderItem creates a new order item
func NewOrderItem(
	orderID, skuID, productID int64,
	name string,
	quantity int,
	retailPrice, salePrice float64,
	taxCategory string,
) (*OrderItem, error) {
	if orderID == 0 {
		return nil, NewDomainError("OrderID cannot be zero for OrderItem")
	}
	if skuID == 0 {
		return nil, NewDomainError("SKUID cannot be zero for OrderItem")
	}
	if productID == 0 {
		return nil, NewDomainError("ProductID cannot be zero for OrderItem")
	}
	if name == "" {
		return nil, NewDomainError("Name cannot be empty for OrderItem")
	}
	if quantity <= 0 {
		return nil, NewDomainError("Quantity must be greater than zero for OrderItem")
	}

	now := time.Now()
	itemPrice := retailPrice
	if salePrice > 0 && salePrice < retailPrice {
		itemPrice = salePrice
	}

	return &OrderItem{
		OrderID:             orderID,
		SKUID:               skuID,
		ProductID:           productID,
		Name:                name,
		Quantity:            quantity,
		RetailPrice:         retailPrice,
		SalePrice:           salePrice,
		Price:               itemPrice, // Initial price before adjustments
		TotalPrice:          itemPrice * float64(quantity),
		TaxCategory:         taxCategory,
		ShippingAmount:      0.0,       // Default
		DiscountsAllowed:    true,      // Default
		HasValidationErrors: false,     // Default
		ItemTaxableFlag:     true,      // Default
		OrderItemType:       "DEFAULT", // Default
		RetailPriceOverride: false,     // Default
		SalePriceOverride:   false,     // Default
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

// UpdateQuantity updates the quantity of the order item
func (oi *OrderItem) UpdateQuantity(newQuantity int) error {
	if newQuantity <= 0 {
		return NewDomainError("New quantity must be greater than zero")
	}
	oi.Quantity = newQuantity
	oi.TotalPrice = oi.Price * float64(newQuantity)
	oi.UpdatedAt = time.Now()
	return nil
}

// UpdatePrices updates the pricing of the order item
func (oi *OrderItem) UpdatePrices(retailPrice, salePrice, currentPrice float64) {
	oi.RetailPrice = retailPrice
	oi.SalePrice = salePrice
	oi.Price = currentPrice // This is the final price after any item-level adjustments
	oi.TotalPrice = currentPrice * float64(oi.Quantity)
	oi.UpdatedAt = time.Now()
}

// SetTaxAmount sets the tax amount for the order item
func (oi *OrderItem) SetTaxAmount(taxAmount float64) {
	oi.TaxAmount = taxAmount
	oi.UpdatedAt = time.Now()
}

// SetShippingAmount sets the shipping amount for the order item
func (oi *OrderItem) SetShippingAmount(shippingAmount float64) {
	oi.ShippingAmount = shippingAmount
	oi.UpdatedAt = time.Now()
}

// SetCategoryID sets the category ID for the order item
func (oi *OrderItem) SetCategoryID(categoryID int64) {
	oi.CategoryID = &categoryID
	oi.UpdatedAt = time.Now()
}

// SetGiftWrapItemID sets the gift wrap item ID for the order item
func (oi *OrderItem) SetGiftWrapItemID(giftWrapItemID int64) {
	oi.GiftWrapItemID = &giftWrapItemID
	oi.UpdatedAt = time.Now()
}

// SetParentOrderItemID sets the parent order item ID for bundled items
func (oi *OrderItem) SetParentOrderItemID(parentItemID int64) {
	oi.ParentOrderItemID = &parentItemID
	oi.UpdatedAt = time.Now()
}

// SetPersonalMessageID sets the personal message ID for the order item
func (oi *OrderItem) SetPersonalMessageID(personalMessageID int64) {
	oi.PersonalMessageID = &personalMessageID
	oi.UpdatedAt = time.Now()
}
