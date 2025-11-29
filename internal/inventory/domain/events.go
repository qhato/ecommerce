package domain

import "time"

// InventoryLevelCreatedEvent is published when a new inventory level is created.
type InventoryLevelCreatedEvent struct {
	InventoryID    string
	SKUID          string
	QuantityOnHand int
	CreationTime   time.Time
}

// QuantityOnHandUpdatedEvent is published when the quantity on hand for a SKU changes.
type QuantityOnHandUpdatedEvent struct {
	InventoryID string
	SKUID       string
	OldQuantity int
	NewQuantity int
	UpdateTime  time.Time
}

// InventoryReservedEvent is published when inventory is reserved.
type InventoryReservedEvent struct {
	ReservationID string
	SKUID         string
	Quantity      int
	OrderID       string
	OrderItemID   string
	ReservedTime  time.Time
}

// InventoryReleasedEvent is published when reserved inventory is released.
type InventoryReleasedEvent struct {
	ReservationID string
	SKUID         string
	Quantity      int
	OrderID       string
	ReleasedTime  time.Time
}

// InventoryFulfilledEvent is published when allocated inventory is fulfilled (e.g., shipped).
type InventoryFulfilledEvent struct {
	ReservationID   string
	SKUID           string
	Quantity        int
	OrderID         string
	FulfillmentTime time.Time
}

// InventoryReceivedEvent is published when new inventory is received.
type InventoryReceivedEvent struct {
	InventoryID  string
	SKUID        string
	Quantity     int
	WarehouseID  *string
	ReceivedTime time.Time
}

// ReorderPointReachedEvent is published when inventory reaches reorder point.
type ReorderPointReachedEvent struct {
	InventoryID     string
	SKUID           string
	CurrentQuantity int
	ReorderPoint    int
	ReorderQuantity int
	EventTime       time.Time
}
