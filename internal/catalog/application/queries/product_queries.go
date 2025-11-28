package queries

import (
	"context"

	"github.com/qhato/ecommerce/internal/catalog/application"
	"github.com/qhato/ecommerce/internal/catalog/domain"
	"github.com/qhato/ecommerce/pkg/cache"
	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/logger"
)

// GetProductByIDQuery represents a query to get a product by ID
type GetProductByIDQuery struct {
	ID int64 `json:"id" validate:"required"`
}

// GetProductByURLQuery represents a query to get a product by URL
type GetProductByURLQuery struct {
	URL string `json:"url" validate:"required"`
}

// ListProductsQuery represents a query to list products
type ListProductsQuery struct {
	Page            int    `json:"page" validate:"min=1"`
	PageSize        int    `json:"page_size" validate:"min=1,max=100"`
	IncludeArchived bool   `json:"include_archived"`
	SortBy          string `json:"sort_by"`
	SortOrder       string `json:"sort_order"`
}

// ListProductsByCategoryQuery represents a query to list products by category
type ListProductsByCategoryQuery struct {
	CategoryID      int64  `json:"category_id" validate:"required"`
	Page            int    `json:"page" validate:"min=1"`
	PageSize        int    `json:"page_size" validate:"min=1,max=100"`
	IncludeArchived bool   `json:"include_archived"`
	SortBy          string `json:"sort_by"`
	SortOrder       string `json:"sort_order"`
}

// SearchProductsQuery represents a query to search products
type SearchProductsQuery struct {
	Query           string `json:"query" validate:"required"`
	Page            int    `json:"page" validate:"min=1"`
	PageSize        int    `json:"page_size" validate:"min=1,max=100"`
	IncludeArchived bool   `json:"include_archived"`
	SortBy          string `json:"sort_by"`
	SortOrder       string `json:"sort_order"`
}

// ProductQueryHandler handles product queries
type ProductQueryHandler struct {
	repo   domain.ProductRepository
	cache  cache.Cache
	logger *logger.Logger
}

// NewProductQueryHandler creates a new product query handler
func NewProductQueryHandler(
	repo domain.ProductRepository,
	cache cache.Cache,
	logger *logger.Logger,
) *ProductQueryHandler {
	return &ProductQueryHandler{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

// HandleGetProductByID handles the get product by ID query
func (h *ProductQueryHandler) HandleGetProductByID(ctx context.Context, query *GetProductByIDQuery) (*application.ProductDTO, error) {
	// Try to get from cache first
	cacheKey := productCacheKey(query.ID)
	var product *domain.Product

	if err := h.cache.Get(ctx, cacheKey, &product); err == nil && product != nil {
		h.logger.Debug("product found in cache", "product_id", query.ID)
		return application.ToProductDTO(product), nil
	}

	// Get from repository
	product, err := h.repo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, errors.Wrap(err, "product not found")
	}

	// Cache the result
	if err := h.cache.Set(ctx, cacheKey, product, cache.DefaultTTL); err != nil {
		h.logger.Warn("failed to cache product", "error", err, "product_id", query.ID)
	}

	return application.ToProductDTO(product), nil
}

// HandleGetProductByURL handles the get product by URL query
func (h *ProductQueryHandler) HandleGetProductByURL(ctx context.Context, query *GetProductByURLQuery) (*application.ProductDTO, error) {
	product, err := h.repo.FindByURL(ctx, query.URL)
	if err != nil {
		return nil, errors.Wrap(err, "product not found")
	}

	// Cache the result
	cacheKey := productCacheKey(product.ID)
	if err := h.cache.Set(ctx, cacheKey, product, cache.DefaultTTL); err != nil {
		h.logger.Warn("failed to cache product", "error", err, "product_id", product.ID)
	}

	return application.ToProductDTO(product), nil
}

// HandleListProducts handles the list products query
func (h *ProductQueryHandler) HandleListProducts(ctx context.Context, query *ListProductsQuery) (*application.PaginatedResponse, error) {
	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	if query.SortBy == "" {
		query.SortBy = "created_at"
	}
	if query.SortOrder == "" {
		query.SortOrder = "desc"
	}

	// Create filter
	filter := &domain.ProductFilter{
		Page:            query.Page,
		PageSize:        query.PageSize,
		IncludeArchived: query.IncludeArchived,
		SortBy:          query.SortBy,
		SortOrder:       query.SortOrder,
	}

	// Get from repository
	products, total, err := h.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list products")
	}

	// Convert to DTOs
	productDTOs := make([]*application.ProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = application.ToProductDTO(product)
	}

	return application.NewPaginatedResponse(productDTOs, query.Page, query.PageSize, total), nil
}

// HandleListProductsByCategory handles the list products by category query
func (h *ProductQueryHandler) HandleListProductsByCategory(ctx context.Context, query *ListProductsByCategoryQuery) (*application.PaginatedResponse, error) {
	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	if query.SortBy == "" {
		query.SortBy = "created_at"
	}
	if query.SortOrder == "" {
		query.SortOrder = "desc"
	}

	// Create filter
	filter := &domain.ProductFilter{
		Page:            query.Page,
		PageSize:        query.PageSize,
		IncludeArchived: query.IncludeArchived,
		SortBy:          query.SortBy,
		SortOrder:       query.SortOrder,
	}

	// Get from repository
	products, total, err := h.repo.FindByCategoryID(ctx, query.CategoryID, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list products by category")
	}

	// Convert to DTOs
	productDTOs := make([]*application.ProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = application.ToProductDTO(product)
	}

	return application.NewPaginatedResponse(productDTOs, query.Page, query.PageSize, total), nil
}

// HandleSearchProducts handles the search products query
func (h *ProductQueryHandler) HandleSearchProducts(ctx context.Context, query *SearchProductsQuery) (*application.PaginatedResponse, error) {
	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	if query.SortBy == "" {
		query.SortBy = "created_at"
	}
	if query.SortOrder == "" {
		query.SortOrder = "desc"
	}

	// Create filter
	filter := &domain.ProductFilter{
		Page:            query.Page,
		PageSize:        query.PageSize,
		IncludeArchived: query.IncludeArchived,
		SortBy:          query.SortBy,
		SortOrder:       query.SortOrder,
	}

	// Search from repository
	products, total, err := h.repo.Search(ctx, query.Query, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to search products")
	}

	// Convert to DTOs
	productDTOs := make([]*application.ProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = application.ToProductDTO(product)
	}

	return application.NewPaginatedResponse(productDTOs, query.Page, query.PageSize, total), nil
}

// productCacheKey generates a cache key for a product
func productCacheKey(id int64) string {
	return cache.Key("catalog", "product", id)
}
