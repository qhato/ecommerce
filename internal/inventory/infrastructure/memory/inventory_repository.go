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
	inventories map[string]*domain.InventoryLevel
	// A map to quickly find inventory by SKU ID
	skuIDIndex     map[string]string
	warehouseIndex map[string][]string
}

// NewInventoryRepository creates a new in-memory inventory repository.
func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{
		inventories:    make(map[string]*domain.InventoryLevel),
		skuIDIndex:     make(map[string]string),
		warehouseIndex: make(map[string][]string),
	}
}

// Save stores a new InventoryLevel record or updates an existing one.
func (r *InventoryRepository) Save(ctx context.Context, level *domain.InventoryLevel) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.inventories[level.ID] = level
	r.skuIDIndex[level.SKUID] = level.ID

	if level.WarehouseID != nil {
		warehouseID := *level.WarehouseID
		if _, exists := r.warehouseIndex[warehouseID]; !exists {
			r.warehouseIndex[warehouseID] = []string{}
		}
		r.warehouseIndex[warehouseID] = append(r.warehouseIndex[warehouseID], level.ID)
	}

	return nil
}

// FindByID retrieves an InventoryLevel record by its unique identifier.
func (r *InventoryRepository) FindByID(ctx context.Context, id string) (*domain.InventoryLevel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	level, ok := r.inventories[id]
	if !ok {
		return nil, nil
	}
	return level, nil
}

// FindBySKUID retrieves an InventoryLevel record by its associated SKU ID.
func (r *InventoryRepository) FindBySKUID(ctx context.Context, skuID string) (*domain.InventoryLevel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.skuIDIndex[skuID]
	if !ok {
		return nil, nil
	}
	return r.inventories[id], nil
}

// FindByWarehouse retrieves InventoryLevel records by warehouse.
func (r *InventoryRepository) FindByWarehouse(ctx context.Context, warehouseID string) ([]*domain.InventoryLevel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids, ok := r.warehouseIndex[warehouseID]
	if !ok {
		return []*domain.InventoryLevel{}, nil
	}

	var result []*domain.InventoryLevel
	for _, id := range ids {
		if level, exists := r.inventories[id]; exists {
			result = append(result, level)
		}
	}
	return result, nil
}

// Delete removes an InventoryLevel record by its unique identifier.
func (r *InventoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	level, ok := r.inventories[id]
	if !ok {
		return fmt.Errorf("inventory level with ID %s not found", id)
	}
	delete(r.inventories, id)
	delete(r.skuIDIndex, level.SKUID)

	if level.WarehouseID != nil {
		warehouseID := *level.WarehouseID
		if ids, exists := r.warehouseIndex[warehouseID]; exists {
			filtered := make([]string, 0)
			for _, wID := range ids {
				if wID != id {
					filtered = append(filtered, wID)
				}
			}
			r.warehouseIndex[warehouseID] = filtered
		}
	}

	return nil
}
