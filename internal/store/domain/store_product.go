package domain

import (
	"errors"
	"time"
)

// StoreProduct represents a product's availability and configuration in a specific store
type StoreProduct struct {
	ID            int64
	StoreID       int64
	ProductID     int64
	IsAvailable   bool
	PriceOverride *StorePriceOverride
	InventoryMode InventoryMode
	SortOrder     int
	Metadata      map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// InventoryMode represents how inventory is managed for this product in the store
type InventoryMode string

const (
	InventoryModeTracked   InventoryMode = "TRACKED"   // Track inventory quantities
	InventoryModeUntracked InventoryMode = "UNTRACKED" // Don't track (always available)
	InventoryModeBackorder InventoryMode = "BACKORDER" // Allow backorders
)

// StorePriceOverride represents store-specific pricing
type StorePriceOverride struct {
	ID              int64
	StoreProductID  int64
	BasePrice       float64
	SalePrice       *float64
	CostPrice       *float64
	Currency        string
	TaxIncluded     bool
	ValidFrom       time.Time
	ValidUntil      *time.Time
	PriceListID     *int64 // Reference to a price list if using price lists
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PriceList represents a named price list that can be assigned to stores
type PriceList struct {
	ID          int64
	Code        string
	Name        string
	Description string
	Currency    string
	IsActive    bool
	Priority    int
	ValidFrom   time.Time
	ValidUntil  *time.Time
	Metadata    map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PriceListEntry represents a product price in a price list
type PriceListEntry struct {
	ID          int64
	PriceListID int64
	ProductID   int64
	BasePrice   float64
	SalePrice   *float64
	MinQuantity int
	MaxQuantity *int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// StorePaymentMethod represents payment methods available in a store
type StorePaymentMethod struct {
	ID              int64
	StoreID         int64
	PaymentMethodID int64
	IsEnabled       bool
	Priority        int
	Settings        map[string]interface{}
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// StoreShippingMethod represents shipping methods available in a store
type StoreShippingMethod struct {
	ID               int64
	StoreID          int64
	ShippingMethodID int64
	IsEnabled        bool
	Priority         int
	Settings         map[string]interface{}
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewStoreProduct creates a new store-product association
func NewStoreProduct(storeID, productID int64) (*StoreProduct, error) {
	if storeID == 0 {
		return nil, errors.New("store ID is required")
	}
	if productID == 0 {
		return nil, errors.New("product ID is required")
	}

	now := time.Now()
	return &StoreProduct{
		StoreID:       storeID,
		ProductID:     productID,
		IsAvailable:   true,
		InventoryMode: InventoryModeTracked,
		SortOrder:     0,
		Metadata:      make(map[string]interface{}),
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// SetPriceOverride sets a price override for this store-product
func (sp *StoreProduct) SetPriceOverride(basePrice float64, currency string) {
	sp.PriceOverride = &StorePriceOverride{
		BasePrice:  basePrice,
		Currency:   currency,
		ValidFrom:  time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	sp.UpdatedAt = time.Now()
}

// GetEffectivePrice returns the effective price (sale price if available, otherwise base price)
func (sp *StoreProduct) GetEffectivePrice() float64 {
	if sp.PriceOverride == nil {
		return 0
	}
	if sp.PriceOverride.SalePrice != nil && *sp.PriceOverride.SalePrice > 0 {
		return *sp.PriceOverride.SalePrice
	}
	return sp.PriceOverride.BasePrice
}

// IsOnSale checks if the product is currently on sale in this store
func (sp *StoreProduct) IsOnSale() bool {
	if sp.PriceOverride == nil || sp.PriceOverride.SalePrice == nil {
		return false
	}
	now := time.Now()
	if sp.PriceOverride.ValidUntil != nil && now.After(*sp.PriceOverride.ValidUntil) {
		return false
	}
	return *sp.PriceOverride.SalePrice > 0 && *sp.PriceOverride.SalePrice < sp.PriceOverride.BasePrice
}

// Activate makes the product available in this store
func (sp *StoreProduct) Activate() {
	sp.IsAvailable = true
	sp.UpdatedAt = time.Now()
}

// Deactivate makes the product unavailable in this store
func (sp *StoreProduct) Deactivate() {
	sp.IsAvailable = false
	sp.UpdatedAt = time.Now()
}

// NewPriceList creates a new price list
func NewPriceList(code, name, currency string) (*PriceList, error) {
	if code == "" {
		return nil, errors.New("price list code is required")
	}
	if name == "" {
		return nil, errors.New("price list name is required")
	}

	now := time.Now()
	return &PriceList{
		Code:      code,
		Name:      name,
		Currency:  currency,
		IsActive:  true,
		Priority:  0,
		ValidFrom: now,
		Metadata:  make(map[string]interface{}),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Activate activates the price list
func (pl *PriceList) Activate() {
	pl.IsActive = true
	pl.UpdatedAt = time.Now()
}

// Deactivate deactivates the price list
func (pl *PriceList) Deactivate() {
	pl.IsActive = false
	pl.UpdatedAt = time.Now()
}

// IsCurrentlyValid checks if the price list is currently valid
func (pl *PriceList) IsCurrentlyValid() bool {
	if !pl.IsActive {
		return false
	}
	now := time.Now()
	if now.Before(pl.ValidFrom) {
		return false
	}
	if pl.ValidUntil != nil && now.After(*pl.ValidUntil) {
		return false
	}
	return true
}
