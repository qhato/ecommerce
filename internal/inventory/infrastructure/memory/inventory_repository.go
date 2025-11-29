package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/qhato/ecommerce/internal/inventory/domain"
)

// InventoryRepository implements domain.InventoryRepository for in-memory persistence.
type InventoryRepository struct {
	mu          sync.RWMutex
	inventories map[int64]*domain.SkuInventory
	nextID      int64
	// A map to quickly find inventory by SKU ID
	skuIDIndex  map[int64]int64
}

// NewInventoryRepository creates a new in-memory inventory repository.
func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{
		inventories: make(map[int64]*domain.SkuInventory),
		nextID:      1,
		skuIDIndex:  make(map[int64]int64),
	}
}

// Save stores a new SkuInventory record or updates an existing one.
func (r *InventoryRepository) Save(ctx context.Context, inventory *domain.SkuInventory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if inventory.ID == 0 {
		inventory.ID = r.nextID
		r.nextID++
	}
	r.inventories[inventory.ID] = inventory
	r.skuIDIndex[inventory.SKUID] = inventory.ID
	return nil
}

// FindByID retrieves a SkuInventory record by its unique identifier.
func (r *InventoryRepository) FindByID(ctx context.Context, id int64) (*domain.SkuInventory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	inventory, ok := r.inventories[id]
	if !ok {
		return nil, nil
	}
	return inventory, nil
}

// FindBySKUID retrieves a SkuInventory record by its associated SKU ID.
func (r *InventoryRepository) FindBySKUID(ctx context.Context, skuID int64) (*domain.SkuInventory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.skuIDIndex[skuID]
	if !ok {
		return nil, nil
	}
	return r.inventories[id], nil
}

// FindByFulfillmentLocation retrieves SkuInventory records by fulfillment location.
func (r *InventoryRepository) FindByFulfillmentLocation(ctx context.Context, location string) ([]*domain.SkuInventory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.SkuInventory
	for _, inventory := range r.inventories {
		if inventory.FulfillmentLocation == location {
			result = append(result, inventory)
		}
	}
	return result, nil
}

// Delete removes a SkuInventory record by its unique identifier.
func (r *InventoryRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	inventory, ok := r.inventories[id]
	if !ok {
		return fmt.Errorf("SKU inventory with ID %d not found", id)
	}
	delete(r.inventories, id)
	delete(r.skuIDIndex, inventory.SKUID)
	return nil
}
