package memory

import (
	"context"
	"sync"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// ProductRepository implements domain.ProductRepository for in-memory persistence.
type ProductRepository struct {
	mu       sync.RWMutex
	products map[int64]*domain.Product
	nextID   int64
}

// NewProductRepository creates a new in-memory product repository.
func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: make(map[int64]*domain.Product),
		nextID:   1,
	}
}

// Save stores a new product or updates an existing one.
func (r *ProductRepository) Save(ctx context.Context, product *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if product.ID == 0 {
		product.ID = r.nextID
		r.nextID++
	}
	r.products[product.ID] = product
	return nil
}

// FindByID retrieves a product by its unique identifier.
func (r *ProductRepository) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[id]
	if !ok {
		return nil, nil // Or return a specific "not found" error
	}
	return product, nil
}

// FindAll retrieves all products (for demonstration purposes, in a real app this would be paginated).
func (r *ProductRepository) FindAll(ctx context.Context) ([]*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*domain.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}
	return products, nil
}

// Delete removes a product by its unique identifier.
func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		return fmt.Errorf("product with ID %d not found", id)
	}
	delete(r.products, id)
	return nil
}
