package domain

import (
	"context"
)

// InventoryRepository provides an interface for managing inventory levels.
type InventoryRepository interface {
	// Save stores a new inventory level or updates an existing one.
	Save(ctx context.Context, level *InventoryLevel) error

	// FindByID retrieves an inventory level by its unique identifier.
	FindByID(ctx context.Context, id string) (*InventoryLevel, error)

	// FindBySKUID retrieves an inventory level by its associated SKU ID.
	FindBySKUID(ctx context.Context, skuID string) (*InventoryLevel, error)

	// FindByWarehouse retrieves inventory levels by warehouse.
	FindByWarehouse(ctx context.Context, warehouseID string) ([]*InventoryLevel, error)

	// Delete removes an inventory level by its unique identifier.
	Delete(ctx context.Context, id string) error
}

// InventoryReservationRepository provides an interface for managing reservations.
type InventoryReservationRepository interface {
	// Save stores a new reservation or updates an existing one.
	Save(ctx context.Context, reservation *InventoryReservation) error

	// FindByID retrieves a reservation by its unique identifier.
	FindByID(ctx context.Context, id string) (*InventoryReservation, error)

	// FindByOrderID retrieves all reservations for an order.
	FindByOrderID(ctx context.Context, orderID string) ([]*InventoryReservation, error)

	// FindExpired retrieves all expired reservations.
	FindExpired(ctx context.Context) ([]*InventoryReservation, error)

	// Delete removes a reservation by its unique identifier.
	Delete(ctx context.Context, id string) error
}
