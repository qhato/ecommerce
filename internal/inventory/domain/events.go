package domain

import (
	"time"
)

// SkuInventoryCreatedEvent is published when a new SkuInventory record is created.
type SkuInventoryCreatedEvent struct {
	InventoryID    int64
	SKUID          int64
	QuantityOnHand int
	CreationTime   time.Time
}

// QuantityOnHandUpdatedEvent is published when the quantity on hand for a SKU changes.
type QuantityOnHandUpdatedEvent struct {
	InventoryID    int64
	SKUID          int64
	OldQuantity    int
	NewQuantity    int
	UpdateTime     time.Time
}

// QuantityAllocatedEvent is published when inventory is allocated for an order.
type QuantityAllocatedEvent struct {
	InventoryID    int64
	SKUID          int64
	Quantity       int
	OrderID        int64
	AllocationTime time.Time
}

// QuantityDeallocatedEvent is published when allocated inventory is released.
type QuantityDeallocatedEvent struct {
	InventoryID      int64
	SKUID            int64
	Quantity         int
	OrderID          int64
	DeallocationTime time.Time
}

// InventoryFulfilledEvent is published when allocated inventory is fulfilled (e.g., shipped).
type InventoryFulfilledEvent struct {
	InventoryID    int64
	SKUID          int64
	Quantity       int
	OrderID        int64
	FulfillmentTime time.Time
}

// InventoryStatusUpdatedEvent is published when the inventory status of a SKU changes.
type InventoryStatusUpdatedEvent struct {
	InventoryID int64
	SKUID       int64
	OldStatus   InventoryStatus
	NewStatus   InventoryStatus
	UpdateTime  time.Time
}
