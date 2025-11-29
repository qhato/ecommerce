package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// SkuRepository implements domain.SkuRepository for in-memory persistence.
type SkuRepository struct {
	mu   sync.RWMutex
	skus map[int64]*domain.SKU
	nextID int64
}

// NewSkuRepository creates a new in-memory SKU repository.
func NewSkuRepository() *SkuRepository {
	return &SkuRepository{
		skus:   make(map[int64]*domain.SKU),
		nextID: 1,
	}
}

// Save stores a new SKU or updates an existing one.
func (r *SkuRepository) Save(ctx context.Context, sku *domain.SKU) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if sku.ID == 0 {
		sku.ID = r.nextID
		r.nextID++
	}
	r.skus[sku.ID] = sku
	return nil
}

// FindByID retrieves a SKU by its unique identifier.
func (r *SkuRepository) FindByID(ctx context.Context, id int64) (*domain.SKU, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sku, ok := r.skus[id]
	if !ok {
		return nil, nil
	}
	return sku, nil
}

// Delete removes a SKU by its unique identifier.
func (r *SkuRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.skus[id]; !ok {
		return fmt.Errorf("SKU with ID %d not found", id)
	}
	delete(r.skus, id)
	return nil
}
