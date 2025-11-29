package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// ProductOptionRepository implements domain.ProductOptionRepository for in-memory persistence.
type ProductOptionRepository struct {
	mu      sync.RWMutex
	options map[int64]*domain.ProductOption
	nextID  int64
}

// NewProductOptionRepository creates a new in-memory product option repository.
func NewProductOptionRepository() *ProductOptionRepository {
	return &ProductOptionRepository{
		options: make(map[int64]*domain.ProductOption),
		nextID:  1,
	}
}

// Save stores a new product option or updates an existing one.
func (r *ProductOptionRepository) Save(ctx context.Context, option *domain.ProductOption) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if option.ID == 0 {
		option.ID = r.nextID
		r.nextID++
	}
	r.options[option.ID] = option
	return nil
}

// FindByID retrieves a product option by its unique identifier.
func (r *ProductOptionRepository) FindByID(ctx context.Context, id int64) (*domain.ProductOption, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	option, ok := r.options[id]
	if !ok {
		return nil, nil
	}
	return option, nil
}

// Delete removes a product option by its unique identifier.
func (r *ProductOptionRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.options[id]; !ok {
		return fmt.Errorf("product option with ID %d not found", id)
	}
	delete(r.options, id)
	return nil
}

// ProductOptionValueRepository implements domain.ProductOptionValueRepository for in-memory persistence.
type ProductOptionValueRepository struct {
	mu     sync.RWMutex
	values map[int64]*domain.ProductOptionValue
	nextID int64
}

// NewProductOptionValueRepository creates a new in-memory product option value repository.
func NewProductOptionValueRepository() *ProductOptionValueRepository {
	return &ProductOptionValueRepository{
		values: make(map[int64]*domain.ProductOptionValue),
		nextID: 1,
	}
}

// Save stores a new product option value or updates an existing one.
func (r *ProductOptionValueRepository) Save(ctx context.Context, value *domain.ProductOptionValue) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if value.ID == 0 {
		value.ID = r.nextID
		r.nextID++
	}
	r.values[value.ID] = value
	return nil
}

// FindByID retrieves a product option value by its unique identifier.
func (r *ProductOptionValueRepository) FindByID(ctx context.Context, id int64) (*domain.ProductOptionValue, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.values[id]
	if !ok {
		return nil, nil
	}
	return value, nil
}

// FindByProductOptionID retrieves all product option values for a given product option ID.
func (r *ProductOptionValueRepository) FindByProductOptionID(ctx context.Context, optionID int64) ([]*domain.ProductOptionValue, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.ProductOptionValue
	for _, value := range r.values {
		if value.ProductOptionID == optionID {
			result = append(result, value)
		}
	}
	return result, nil
}

// Delete removes a product option value by its unique identifier.
func (r *ProductOptionValueRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.values[id]; !ok {
		return fmt.Errorf("product option value with ID %d not found", id)
	}
	delete(r.values, id)
	return nil
}
