package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// ProductBundle represents a bundle of products sold together
type ProductBundle struct {
	ID            int64
	Name          string
	Description   string
	BundlePrice   decimal.Decimal
	IsActive      bool
	Priority      int
	Items         []ProductBundleItem
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ProductBundleItem represents an item in a product bundle
type ProductBundleItem struct {
	ID           int64
	BundleID     int64
	ProductID    *int64 // Product or SKU, one must be set
	SKUID        *int64
	Quantity     int
	SortOrder    int
	CreatedAt    time.Time
}

// NewProductBundle creates a new product bundle
func NewProductBundle(name, description string, bundlePrice decimal.Decimal) *ProductBundle {
	now := time.Now()
	return &ProductBundle{
		Name:        name,
		Description: description,
		BundlePrice: bundlePrice,
		IsActive:    false,
		Priority:    0,
		Items:       make([]ProductBundleItem, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddItem adds an item to the bundle
func (b *ProductBundle) AddItem(productID, skuID *int64, quantity, sortOrder int) error {
	if productID == nil && skuID == nil {
		return ErrInvalidBundleItem
	}
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	item := ProductBundleItem{
		BundleID:  b.ID,
		ProductID: productID,
		SKUID:     skuID,
		Quantity:  quantity,
		SortOrder: sortOrder,
		CreatedAt: time.Now(),
	}

	b.Items = append(b.Items, item)
	b.UpdatedAt = time.Now()
	return nil
}

// Activate activates the bundle
func (b *ProductBundle) Activate() {
	b.IsActive = true
	b.UpdatedAt = time.Now()
}

// Deactivate deactivates the bundle
func (b *ProductBundle) Deactivate() {
	b.IsActive = false
	b.UpdatedAt = time.Now()
}

// UpdatePrice updates the bundle price
func (b *ProductBundle) UpdatePrice(price decimal.Decimal) error {
	if price.LessThan(decimal.Zero) {
		return ErrInvalidPrice
	}
	b.BundlePrice = price
	b.UpdatedAt = time.Now()
	return nil
}
