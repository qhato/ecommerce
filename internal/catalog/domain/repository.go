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

	// FindByCategoryID retrieves products by category ID
	FindByCategoryID(ctx context.Context, categoryID int64, filter *ProductFilter) ([]*Product, int64, error)

	// FindAll retrieves all products with pagination
	FindAll(ctx context.Context, filter *ProductFilter) ([]*Product, int64, error)

	// Search searches products by query
	Search(ctx context.Context, query string, filter *ProductFilter) ([]*Product, int64, error)
}

// ProductAttributeRepository defines the interface for product attribute persistence
type ProductAttributeRepository interface {
	// Save stores a new product attribute or updates an existing one.
	Save(ctx context.Context, attribute *ProductAttribute) error

	// FindByID retrieves a product attribute by its unique identifier.
	FindByID(ctx context.Context, id int64) (*ProductAttribute, error)

	// FindByProductID retrieves all product attributes for a given product ID.
	FindByProductID(ctx context.Context, productID int64) ([]*ProductAttribute, error)

	// Delete removes a product attribute by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteByProductID removes all product attributes for a given product ID.
	DeleteByProductID(ctx context.Context, productID int64) error
}

// ProductOptionXrefRepository defines the interface for product option cross-reference persistence
type ProductOptionXrefRepository interface {
	// Save stores a new product option cross-reference.
	Save(ctx context.Context, xref *ProductOptionXref) error

	// FindByID retrieves a product option cross-reference by its unique identifier.
	FindByID(ctx context.Context, id int64) (*ProductOptionXref, error)

	// FindByProductID retrieves all product option cross-references for a given product ID.
	FindByProductID(ctx context.Context, productID int64) ([]*ProductOptionXref, error)

	// FindByProductOptionID retrieves all product option cross-references for a given product option ID.
	FindByProductOptionID(ctx context.Context, productOptionID int64) ([]*ProductOptionXref, error)

	// Delete removes a product option cross-reference by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteByProductID removes all product option cross-references for a given product ID.
	DeleteByProductID(ctx context.Context, productID int64) error

	// DeleteByProductOptionID removes all product option cross-references for a given product option ID.
	DeleteByProductOptionID(ctx context.Context, productOptionID int64) error

	// RemoveProductOptionXref removes a specific product option cross-reference by product ID and product option ID.
	RemoveProductOptionXref(ctx context.Context, productID, productOptionID int64) error
}

// CategoryProductXrefRepository defines the interface for category-product cross-reference persistence
type CategoryProductXrefRepository interface {
	// Save stores a new category-product cross-reference.
	Save(ctx context.Context, xref *CategoryProductXref) error

	// FindByID retrieves a category-product cross-reference by its unique identifier.
	FindByID(ctx context.Context, id int64) (*CategoryProductXref, error)

	// FindByCategoryID retrieves all category-product cross-references for a given category ID.
	FindByCategoryID(ctx context.Context, categoryID int64) ([]*CategoryProductXref, error)

	// FindByProductID retrieves all category-product cross-references for a given product ID.
	FindByProductID(ctx context.Context, productID int64) ([]*CategoryProductXref, error)

	// Delete removes a category-product cross-reference by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// RemoveCategoryProductXref removes a specific category-product cross-reference by category ID and product ID.
	RemoveCategoryProductXref(ctx context.Context, categoryID, productID int64) error
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

// CategoryAttributeRepository defines the interface for category attribute persistence
type CategoryAttributeRepository interface {
	// Save stores a new category attribute or updates an existing one.
	Save(ctx context.Context, attribute *CategoryAttribute) error

	// FindByID retrieves a category attribute by its unique identifier.
	FindByID(ctx context.Context, id int64) (*CategoryAttribute, error)

	// FindByCategoryID retrieves all category attributes for a given category ID.
	FindByCategoryID(ctx context.Context, categoryID int64) ([]*CategoryAttribute, error)

	// Delete removes a category attribute by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteByCategoryID removes all category attributes for a given category ID.
	DeleteByCategoryID(ctx context.Context, categoryID int64) error
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

// SKUAttributeRepository defines the interface for SKU attribute persistence
type SKUAttributeRepository interface {
	// Save stores a new SKU attribute or updates an existing one.
	Save(ctx context.Context, attribute *SKUAttribute) error

	// FindByID retrieves a SKU attribute by its unique identifier.
	FindByID(ctx context.Context, id int64) (*SKUAttribute, error)

	// FindBySKUID retrieves all SKU attributes for a given SKU ID.
	FindBySKUID(ctx context.Context, skuID int64) ([]*SKUAttribute, error)

	// Delete removes a SKU attribute by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteBySKUID removes all SKU attributes for a given SKU ID.
	DeleteBySKUID(ctx context.Context, skuID int64) error
}

// SkuProductOptionValueXrefRepository defines the interface for SKU product option value cross-reference persistence
type SkuProductOptionValueXrefRepository interface {
	// Save stores a new SKU product option value cross-reference.
	Save(ctx context.Context, xref *SkuProductOptionValueXref) error

	// FindByID retrieves a SKU product option value cross-reference by its unique identifier.
	FindByID(ctx context.Context, id int64) (*SkuProductOptionValueXref, error)

	// FindBySKUID retrieves all SKU product option value cross-references for a given SKU ID.
	FindBySKUID(ctx context.Context, skuID int64) ([]*SkuProductOptionValueXref, error)

	// FindByProductOptionValueID retrieves all SKU product option value cross-references for a given product option value ID.
	FindByProductOptionValueID(ctx context.Context, productOptionValueID int64) ([]*SkuProductOptionValueXref, error)

	// Delete removes a SKU product option value cross-reference by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteBySKUID removes all SKU product option value cross-references for a given SKU ID.
	DeleteBySKUID(ctx context.Context, skuID int64) error

	// DeleteByProductOptionValueID removes all SKU product option value cross-references for a given product option value ID.
	DeleteByProductOptionValueID(ctx context.Context, productOptionValueID int64) error

	// RemoveSkuProductOptionValueXref removes a specific SKU product option value cross-reference by SKU ID and product option value ID.
	RemoveSkuProductOptionValueXref(ctx context.Context, skuID, productOptionValueID int64) error
}

// ProductOptionRepository defines the interface for ProductOption persistence
type ProductOptionRepository interface {
	// Save stores a new product option or updates an existing one.
	Save(ctx context.Context, option *ProductOption) error

	// FindByID retrieves a product option by its unique identifier.
	FindByID(ctx context.Context, id int64) (*ProductOption, error)

	// FindAll retrieves all product options with pagination.
	FindAll(ctx context.Context, filter *ProductOptionFilter) ([]*ProductOption, int64, error)

	// Delete removes a product option by its unique identifier.
	Delete(ctx context.Context, id int64) error
}

// ProductOptionValueRepository defines the interface for ProductOptionValue persistence
type ProductOptionValueRepository interface {
	// Save stores a new product option value or updates an existing one.
	Save(ctx context.Context, value *ProductOptionValue) error

	// FindByID retrieves a product option value by its unique identifier.
	FindByID(ctx context.Context, id int64) (*ProductOptionValue, error)

	// FindByProductOptionID retrieves all product option values for a given product option ID.
	FindByProductOptionID(ctx context.Context, productOptionID int64) ([]*ProductOptionValue, error)

	// Delete removes a product option value by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteByProductOptionID removes all product option values for a given product option ID.
	DeleteByProductOptionID(ctx context.Context, productOptionID int64) error
}

// ProductFilter represents filtering and pagination options for products
type ProductFilter struct {
	Page            int
	PageSize        int
	IncludeArchived bool
	SortBy          string // "name", "created_at", "updated_at", "price"
	SortOrder       string // "asc", "desc"
}

// CategoryFilter represents filtering and pagination options for categories
type CategoryFilter struct {
	Page            int
	PageSize        int
	IncludeArchived bool
	ActiveOnly      bool
	SortBy          string // "name", "display_order", "created_at"
	SortOrder       string // "asc", "desc"
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

// ProductOptionFilter represents filtering and pagination options for product options
type ProductOptionFilter struct {
	Page      int
	PageSize  int
	SortBy    string // "name", "display_order", "created_at"
	SortOrder string // "asc", "desc"
}

// NewProductFilter creates a default product filter
func NewProductFilter() *ProductFilter {
	return &ProductFilter{
		Page:            1,
		PageSize:        20,
		IncludeArchived: false,
		SortBy:          "created_at",
		SortOrder:       "desc",
	}
}

// NewCategoryFilter creates a default category filter
func NewCategoryFilter() *CategoryFilter {
	return &CategoryFilter{
		Page:            1,
		PageSize:        20,
		IncludeArchived: false,
		ActiveOnly:      true,
		SortBy:          "display_order",
		SortOrder:       "asc",
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

// NewProductOptionFilter creates a default product option filter
func NewProductOptionFilter() *ProductOptionFilter {
	return &ProductOptionFilter{
		Page:      1,
		PageSize:  20,
		SortBy:    "display_order",
		SortOrder: "asc",
	}
}

// ProductBundleRepository defines the interface for product bundle persistence
type ProductBundleRepository interface {
	Create(ctx context.Context, bundle *ProductBundle) error
	Update(ctx context.Context, bundle *ProductBundle) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*ProductBundle, error)
	FindAll(ctx context.Context, activeOnly bool) ([]*ProductBundle, error)
	FindByProduct(ctx context.Context, productID int64) ([]*ProductBundle, error)
}

// ProductBundleItemRepository defines the interface for product bundle item persistence
type ProductBundleItemRepository interface {
	Create(ctx context.Context, item *ProductBundleItem) error
	Delete(ctx context.Context, id int64) error
	FindByBundleID(ctx context.Context, bundleID int64) ([]*ProductBundleItem, error)
	DeleteByBundleID(ctx context.Context, bundleID int64) error
}

// ProductRelationshipRepository defines the interface for product relationship persistence
type ProductRelationshipRepository interface {
	Create(ctx context.Context, relationship *ProductRelationship) error
	Update(ctx context.Context, relationship *ProductRelationship) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*ProductRelationship, error)
	FindByProductID(ctx context.Context, productID int64, relationshipType *ProductRelationshipType) ([]*ProductRelationship, error)
	FindCrossSell(ctx context.Context, productID int64) ([]*ProductRelationship, error)
	FindUpSell(ctx context.Context, productID int64) ([]*ProductRelationship, error)
	FindRelated(ctx context.Context, productID int64) ([]*ProductRelationship, error)
	ExistsByProducts(ctx context.Context, productID, relatedProductID int64, relationshipType ProductRelationshipType) (bool, error)
}
