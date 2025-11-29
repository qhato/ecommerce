package domain

import (
	"time"
)

// InventoryType defines how inventory is managed for a SKU. (Kept for potential future use or other domain structs)
type InventoryType string

const (
	InventoryTypeAlwaysAvailable InventoryType = "ALWAYS_AVAILABLE"
	InventoryTypeCheckQuantity   InventoryType = "CHECK_QUANTITY"
	InventoryTypeProductBundle   InventoryType = "PRODUCT_BUNDLE"
)

// SKUAvailability represents the availability details for a specific SKU in a given location.
type SKUAvailability struct {
	ID                  int64
	SkuID               int64      // From blc_sku_availability.sku_id
	AvailabilityDate    *time.Time // From blc_sku_availability.availability_date
	AvailabilityStatus  string     // From blc_sku_availability.availability_status
	LocationID          *int64     // From blc_sku_availability.location_id
	QtyOnHand           int        // From blc_sku_availability.qty_on_hand
	ReserveQty          int        // From blc_sku_availability.reserve_qty
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NewSKUAvailability creates a new SKUAvailability record.
func NewSKUAvailability(skuID int64, qtyOnHand, reserveQty int, availabilityStatus string) (*SKUAvailability, error) {
	if skuID == 0 {
		return nil, NewDomainError("SKU ID cannot be zero")
	}

	if qtyOnHand < 0 {
		return nil, NewDomainError("Quantity on hand cannot be negative")
	}
	if reserveQty < 0 {
		return nil, NewDomainError("Reserve quantity cannot be negative")
	}
	
	// Basic validation for availabilityStatus
	if availabilityStatus == "" {
		availabilityStatus = "UNKNOWN" // Default status
	}

	now := time.Now()
	return &SKUAvailability{
		SkuID:               skuID,
		QtyOnHand:           qtyOnHand,
		ReserveQty:          reserveQty,
		AvailabilityStatus:  availabilityStatus,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

// UpdateQuantities updates the quantity on hand and reserved quantity.
func (sa *SKUAvailability) UpdateQuantities(qtyOnHand, reserveQty int) error {
	if qtyOnHand < 0 {
		return NewDomainError("Quantity on hand cannot be negative")
	}
	if reserveQty < 0 {
		return NewDomainError("Reserve quantity cannot be negative")
	}
	sa.QtyOnHand = qtyOnHand
	sa.ReserveQty = reserveQty
	sa.UpdatedAt = time.Now()
	return nil
}

// UpdateStatus updates the availability status.
func (sa *SKUAvailability) UpdateStatus(newStatus string) {
	sa.AvailabilityStatus = newStatus
	sa.UpdatedAt = time.Now()
}

// SetAvailabilityDate sets the availability date.
func (sa *SKUAvailability) SetAvailabilityDate(date *time.Time) {
	sa.AvailabilityDate = date
	sa.UpdatedAt = time.Now()
}

// SetLocation sets the location ID.
func (sa *SKUAvailability) SetLocation(locationID int64) {
	sa.LocationID = &locationID
	sa.UpdatedAt = time.Now()
}

// DomainError represents a business rule validation error within the domain.
type DomainError struct {
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new DomainError.
func NewDomainError(message string) error {
	return &DomainError{Message: message}
}
