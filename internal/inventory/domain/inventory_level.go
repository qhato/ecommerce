package domain

import (
	"time"

	"github.com/google/uuid"
)

// InventoryLevel represents the current inventory level for a SKU
type InventoryLevel struct {
	ID                  string
	SKUID               string
	WarehouseID         *string
	LocationID          *string
	QuantityAvailable   int
	QuantityReserved    int
	QuantityOnHand      int // Physical inventory
	QuantityAllocated   int // Allocated for fulfillment
	QuantityBackordered int
	QuantityInTransit   int
	QuantityDamaged     int
	ReorderPoint        int
	ReorderQuantity     int
	SafetyStock         int
	AllowBackorder      bool
	AllowPreorder       bool
	LastCountDate       *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NewInventoryLevel creates a new inventory level
func NewInventoryLevel(skuID string, quantityOnHand int) (*InventoryLevel, error) {
	if skuID == "" {
		return nil, NewDomainError("SKUID is required")
	}

	now := time.Now()
	return &InventoryLevel{
		ID:                uuid.New().String(),
		SKUID:             skuID,
		QuantityAvailable: quantityOnHand,
		QuantityOnHand:    quantityOnHand,
		QuantityReserved:  0,
		QuantityAllocated: 0,
		AllowBackorder:    false,
		AllowPreorder:     false,
		ReorderPoint:      0,
		ReorderQuantity:   0,
		SafetyStock:       0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

// CanReserve checks if the requested quantity can be reserved
func (il *InventoryLevel) CanReserve(quantity int) bool {
	if il.AllowBackorder || il.AllowPreorder {
		return true
	}

	return il.QuantityAvailable >= quantity
}

// Reserve reserves inventory
func (il *InventoryLevel) Reserve(quantity int) error {
	if quantity <= 0 {
		return NewDomainError("Quantity must be positive")
	}

	if !il.CanReserve(quantity) {
		return NewDomainError("Insufficient inventory available")
	}

	il.QuantityReserved += quantity
	il.QuantityAvailable -= quantity

	// If we go negative, we're backordering
	if il.QuantityAvailable < 0 {
		il.QuantityBackordered += (-il.QuantityAvailable)
		il.QuantityAvailable = 0
	}

	il.UpdatedAt = time.Now()
	return nil
}

// Release releases reserved inventory
func (il *InventoryLevel) Release(quantity int) error {
	if quantity <= 0 {
		return NewDomainError("Quantity must be positive")
	}

	if il.QuantityReserved < quantity {
		return NewDomainError("Cannot release more than reserved")
	}

	il.QuantityReserved -= quantity
	il.QuantityAvailable += quantity

	il.UpdatedAt = time.Now()
	return nil
}

// Allocate allocates inventory for fulfillment
func (il *InventoryLevel) Allocate(quantity int) error {
	if quantity <= 0 {
		return NewDomainError("Quantity must be positive")
	}

	if il.QuantityReserved < quantity {
		return NewDomainError("Cannot allocate more than reserved")
	}

	il.QuantityReserved -= quantity
	il.QuantityAllocated += quantity
	il.UpdatedAt = time.Now()
	return nil
}

// Decrement decrements inventory (fulfilled/shipped)
func (il *InventoryLevel) Decrement(quantity int) error {
	if quantity <= 0 {
		return NewDomainError("Quantity must be positive")
	}

	if il.QuantityAllocated < quantity {
		return NewDomainError("Cannot decrement more than allocated")
	}

	il.QuantityAllocated -= quantity
	il.QuantityOnHand -= quantity
	il.UpdatedAt = time.Now()
	return nil
}

// Increment increments inventory (receiving stock)
func (il *InventoryLevel) Increment(quantity int) error {
	if quantity <= 0 {
		return NewDomainError("Quantity must be positive")
	}

	il.QuantityOnHand += quantity
	il.QuantityAvailable += quantity

	// If we had backorders, reduce them
	if il.QuantityBackordered > 0 {
		if quantity >= il.QuantityBackordered {
			il.QuantityBackordered = 0
		} else {
			il.QuantityBackordered -= quantity
		}
	}

	il.UpdatedAt = time.Now()
	return nil
}

// NeedsReorder checks if inventory needs reordering
func (il *InventoryLevel) NeedsReorder() bool {
	if il.ReorderPoint <= 0 {
		return false
	}

	effectiveInventory := il.QuantityOnHand + il.QuantityInTransit
	return effectiveInventory <= il.ReorderPoint
}

// SetReorderPoint sets the reorder point and quantity
func (il *InventoryLevel) SetReorderPoint(point, quantity int) {
	il.ReorderPoint = point
	il.ReorderQuantity = quantity
	il.UpdatedAt = time.Now()
}

// SetSafetyStock sets the safety stock level
func (il *InventoryLevel) SetSafetyStock(quantity int) {
	il.SafetyStock = quantity
	il.UpdatedAt = time.Now()
}

// EnableBackorder enables backorder for this SKU
func (il *InventoryLevel) EnableBackorder() {
	il.AllowBackorder = true
	il.UpdatedAt = time.Now()
}

// DisableBackorder disables backorder for this SKU
func (il *InventoryLevel) DisableBackorder() {
	il.AllowBackorder = false
	il.UpdatedAt = time.Now()
}

// RecordCount records a physical inventory count
func (il *InventoryLevel) RecordCount(counted int) {
	difference := counted - il.QuantityOnHand

	il.QuantityOnHand = counted
	il.QuantityAvailable += difference

	now := time.Now()
	il.LastCountDate = &now
	il.UpdatedAt = now
}
