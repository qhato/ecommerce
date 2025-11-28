package domain

import (
	"context"
)

// ProductRepository defines the interface for product persistence
type ProductRepository interface {
	// Create creates a new product
	Create(ctx context.Context, product *Product) error

	// Update updates an existing product
	Update(ctx context.Context, product *Product) error

	// Delete deletes a product by ID (soft delete - marks as archived)
	Delete(ctx context.Context, id int64) error

	// FindByID retrieves a product by ID
	FindByID(ctx context.Context, id int64) (*Product, error)

	// FindByURL retrieves a product by URL
	FindByURL(ctx context.Context, url string) (*Product, error)

	// FindByURLKey retrieves a product by URL key
	FindByURLKey(ctx context.Context, urlKey string) (*Product, error)

	// FindAll retrieves all products with pagination
	FindAll(ctx context.Context, filter *ProductFilter) ([]*Product, int64, error)

	// FindByCategoryID retrieves products by category ID
	FindByCategoryID(ctx context.Context, categoryID int64, filter *ProductFilter) ([]*Product, int64, error)

	// Search searches products by query
	Search(ctx context.Context, query string, filter *ProductFilter) ([]*Product, int64, error)

	// AddToCategory adds a product to a category
	AddToCategory(ctx context.Context, productID, categoryID int64) error

	// RemoveFromCategory removes a product from a category
	RemoveFromCategory(ctx context.Context, productID, categoryID int64) error
}

// CategoryRepository defines the interface for category persistence
type CategoryRepository interface {
	// Create creates a new category
	Create(ctx context.Context, category *Category) error

	// Update updates an existing category
	Update(ctx context.Context, category *Category) error

	// Delete deletes a category by ID (soft delete - marks as archived)
	Delete(ctx context.Context, id int64) error

	// FindByID retrieves a category by ID
	FindByID(ctx context.Context, id int64) (*Category, error)

	// FindByURL retrieves a category by URL
	FindByURL(ctx context.Context, url string) (*Category, error)

	// FindByURLKey retrieves a category by URL key
	FindByURLKey(ctx context.Context, urlKey string) (*Category, error)

	// FindAll retrieves all categories with pagination
	FindAll(ctx context.Context, filter *CategoryFilter) ([]*Category, int64, error)

	// FindByParentID retrieves child categories by parent ID
	FindByParentID(ctx context.Context, parentID int64, filter *CategoryFilter) ([]*Category, int64, error)

	// FindRootCategories retrieves root categories (no parent)
	FindRootCategories(ctx context.Context, filter *CategoryFilter) ([]*Category, int64, error)

	// GetCategoryPath retrieves the full path from root to category
	GetCategoryPath(ctx context.Context, categoryID int64) ([]*Category, error)
}

// SKURepository defines the interface for SKU persistence
type SKURepository interface {
	// Create creates a new SKU
	Create(ctx context.Context, sku *SKU) error

	// Update updates an existing SKU
	Update(ctx context.Context, sku *SKU) error

	// Delete deletes a SKU by ID
	Delete(ctx context.Context, id int64) error

	// FindByID retrieves a SKU by ID
	FindByID(ctx context.Context, id int64) (*SKU, error)

	// FindByUPC retrieves a SKU by UPC
	FindByUPC(ctx context.Context, upc string) (*SKU, error)

	// FindByProductID retrieves SKUs by product ID
	FindByProductID(ctx context.Context, productID int64) ([]*SKU, error)

	// FindAll retrieves all SKUs with pagination
	FindAll(ctx context.Context, filter *SKUFilter) ([]*SKU, int64, error)

	// UpdateAvailability updates the availability of a SKU
	UpdateAvailability(ctx context.Context, id int64, available bool) error
}

// ProductFilter represents filtering and pagination options for products
type ProductFilter struct {
	Page          int
	PageSize      int
	IncludeArchived bool
	SortBy        string // "name", "created_at", "updated_at", "price"
	SortOrder     string // "asc", "desc"
}

// CategoryFilter represents filtering and pagination options for categories
type CategoryFilter struct {
	Page          int
	PageSize      int
	IncludeArchived bool
	ActiveOnly    bool
	SortBy        string // "name", "display_order", "created_at"
	SortOrder     string // "asc", "desc"
}

// SKUFilter represents filtering and pagination options for SKUs
type SKUFilter struct {
	Page          int
	PageSize      int
	AvailableOnly bool
	ActiveOnly    bool
	SortBy        string // "name", "price", "created_at"
	SortOrder     string // "asc", "desc"
}

// NewProductFilter creates a default product filter
func NewProductFilter() *ProductFilter {
	return &ProductFilter{
		Page:          1,
		PageSize:      20,
		IncludeArchived: false,
		SortBy:        "created_at",
		SortOrder:     "desc",
	}
}

// NewCategoryFilter creates a default category filter
func NewCategoryFilter() *CategoryFilter {
	return &CategoryFilter{
		Page:          1,
		PageSize:      20,
		IncludeArchived: false,
		ActiveOnly:    true,
		SortBy:        "display_order",
		SortOrder:     "asc",
	}
}

// NewSKUFilter creates a default SKU filter
func NewSKUFilter() *SKUFilter {
	return &SKUFilter{
		Page:          1,
		PageSize:      20,
		AvailableOnly: true,
		ActiveOnly:    true,
		SortBy:        "name",
		SortOrder:     "asc",
	}
}
