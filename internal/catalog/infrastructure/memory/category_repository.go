package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// CategoryRepository implements domain.CategoryRepository for in-memory persistence.
type CategoryRepository struct {
	mu        sync.RWMutex
	categories map[int64]*domain.Category
	nextID    int64
}

// NewCategoryRepository creates a new in-memory category repository.
func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		categories: make(map[int64]*domain.Category),
		nextID:     1,
	}
}

// Save stores a new category or updates an existing one.
func (r *CategoryRepository) Save(ctx context.Context, category *domain.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if category.ID == 0 {
		category.ID = r.nextID
		r.nextID++
	}
	r.categories[category.ID] = category
	return nil
}

// FindByID retrieves a category by its unique identifier.
func (r *CategoryRepository) FindByID(ctx context.Context, id int64) (*domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	category, ok := r.categories[id]
	if !ok {
		return nil, nil // Or return a specific "not found" error
	}
	return category, nil
}

// FindAll retrieves all categories.
func (r *CategoryRepository) FindAll(ctx context.Context) ([]*domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]*domain.Category, 0, len(r.categories))
	for _, category := range r.categories {
		categories = append(categories, category)
	}
	return categories, nil
}

// Delete removes a category by its unique identifier.
func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.categories[id]; !ok {
		return fmt.Errorf("category with ID %d not found", id)
	}
	delete(r.categories, id)
	return nil
}
