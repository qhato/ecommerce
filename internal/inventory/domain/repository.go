package domain

import (
	"context"
)

// InventoryRepository provides an interface for managing SkuInventory.
type InventoryRepository interface {
	// Save stores a new SkuInventory record or updates an existing one.
	Save(ctx context.Context, inventory *SkuInventory) error

	// FindByID retrieves a SkuInventory record by its unique identifier.
	FindByID(ctx context.Context, id int64) (*SkuInventory, error)

	// FindBySKUID retrieves a SkuInventory record by its associated SKU ID.
	FindBySKUID(ctx context.Context, skuID int64) (*SkuInventory, error)

	// FindByFulfillmentLocation retrieves SkuInventory records by fulfillment location.
	FindByFulfillmentLocation(ctx context.Context, location string) ([]*SkuInventory, error)

	// Delete removes a SkuInventory record by its unique identifier.
	Delete(ctx context.Context, id int64) error
}
