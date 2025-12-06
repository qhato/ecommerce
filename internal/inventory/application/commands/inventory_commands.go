package commands

import "time"

// Inventory Level Commands
type CreateInventoryLevelCommand struct {
	SKUID            string
	WarehouseID      string
	LocationID       string
	QuantityOnHand   int
	ReorderPoint     int
	ReorderQty       int
	SafetyStock      int
	AllowBackorder   bool
	AllowPreorder    bool
}

type UpdateInventoryLevelCommand struct {
	ID               string
	QuantityOnHand   int
	ReorderPoint     int
	ReorderQty       int
	SafetyStock      int
	AllowBackorder   bool
	AllowPreorder    bool
}

type AdjustInventoryCommand struct {
	SKUID       string
	WarehouseID string
	Adjustment  int  // Positive for increase, negative for decrease
	Reason      string
}

type SetInventoryCommand struct {
	SKUID       string
	WarehouseID string
	NewQuantity int
	Reason      string
}

type DeleteInventoryLevelCommand struct {
	ID string
}

// Inventory Reservation Commands
type ReserveInventoryCommand struct {
	SKUID       string
	Quantity    int
	OrderID     string
	OrderItemID string
	TTL         time.Duration // Time until reservation expires
}

type ConfirmReservationCommand struct {
	ReservationID string
}

type ReleaseReservationCommand struct {
	ReservationID string
}

type FulfillReservationCommand struct {
	ReservationID string
}

type ExtendReservationCommand struct {
	ReservationID    string
	AdditionalTime   time.Duration
}

type ExpireReservationsCommand struct {
	// Expires all reservations that have passed their expiration time
}

type ReleaseOrderReservationsCommand struct{
	OrderID string
}

// Bulk Operations
type BulkAdjustInventoryCommand struct {
	Adjustments []struct {
		SKUID       string
		WarehouseID string
		Adjustment  int
	}
	Reason string
}

type TransferInventoryCommand struct {
	SKUID           string
	FromWarehouseID string
	ToWarehouseID   string
	Quantity        int
	Reason          string
}
