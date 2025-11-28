package queries

import (
	"context"

	"github.com/qhato/ecommerce/internal/catalog/application"
	"github.com/qhato/ecommerce/internal/catalog/domain"
	"github.com/qhato/ecommerce/pkg/cache"
	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/logger"
)

// GetSKUByIDQuery represents a query to get a SKU by ID
type GetSKUByIDQuery struct {
	ID int64 `json:"id" validate:"required"`
}

// GetSKUByUPCQuery represents a query to get a SKU by UPC
type GetSKUByUPCQuery struct {
	UPC string `json:"upc" validate:"required"`
}

// ListSKUsQuery represents a query to list SKUs
type ListSKUsQuery struct {
	Page          int    `json:"page" validate:"min=1"`
	PageSize      int    `json:"page_size" validate:"min=1,max=100"`
	AvailableOnly bool   `json:"available_only"`
	ActiveOnly    bool   `json:"active_only"`
	SortBy        string `json:"sort_by"`
	SortOrder     string `json:"sort_order"`
}

// ListSKUsByProductQuery represents a query to list SKUs by product
type ListSKUsByProductQuery struct {
	ProductID int64 `json:"product_id" validate:"required"`
}

// SKUQueryHandler handles SKU queries
type SKUQueryHandler struct {
	repo   domain.SKURepository
	cache  cache.Cache
	logger *logger.Logger
}

// NewSKUQueryHandler creates a new SKU query handler
func NewSKUQueryHandler(
	repo domain.SKURepository,
	cache cache.Cache,
	logger *logger.Logger,
) *SKUQueryHandler {
	return &SKUQueryHandler{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

// HandleGetSKUByID handles the get SKU by ID query
func (h *SKUQueryHandler) HandleGetSKUByID(ctx context.Context, query *GetSKUByIDQuery) (*application.SKUDTO, error) {
	// Try to get from cache first
	cacheKey := skuCacheKey(query.ID)
	var sku *domain.SKU

	if err := h.cache.Get(ctx, cacheKey, &sku); err == nil && sku != nil {
		h.logger.Debug("SKU found in cache", "sku_id", query.ID)
		return application.ToSKUDTO(sku), nil
	}

	// Get from repository
	sku, err := h.repo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, errors.Wrap(err, "SKU not found")
	}

	// Cache the result
	if err := h.cache.Set(ctx, cacheKey, sku, cache.DefaultTTL); err != nil {
		h.logger.Warn("failed to cache SKU", "error", err, "sku_id", query.ID)
	}

	return application.ToSKUDTO(sku), nil
}

// HandleGetSKUByUPC handles the get SKU by UPC query
func (h *SKUQueryHandler) HandleGetSKUByUPC(ctx context.Context, query *GetSKUByUPCQuery) (*application.SKUDTO, error) {
	sku, err := h.repo.FindByUPC(ctx, query.UPC)
	if err != nil {
		return nil, errors.Wrap(err, "SKU not found")
	}

	// Cache the result
	cacheKey := skuCacheKey(sku.ID)
	if err := h.cache.Set(ctx, cacheKey, sku, cache.DefaultTTL); err != nil {
		h.logger.Warn("failed to cache SKU", "error", err, "sku_id", sku.ID)
	}

	return application.ToSKUDTO(sku), nil
}

// HandleListSKUs handles the list SKUs query
func (h *SKUQueryHandler) HandleListSKUs(ctx context.Context, query *ListSKUsQuery) (*application.PaginatedResponse, error) {
	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	if query.SortBy == "" {
		query.SortBy = "name"
	}
	if query.SortOrder == "" {
		query.SortOrder = "asc"
	}

	// Create filter
	filter := &domain.SKUFilter{
		Page:          query.Page,
		PageSize:      query.PageSize,
		AvailableOnly: query.AvailableOnly,
		ActiveOnly:    query.ActiveOnly,
		SortBy:        query.SortBy,
		SortOrder:     query.SortOrder,
	}

	// Get from repository
	skus, total, err := h.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list SKUs")
	}

	// Convert to DTOs
	skuDTOs := make([]*application.SKUDTO, len(skus))
	for i, sku := range skus {
		skuDTOs[i] = application.ToSKUDTO(sku)
	}

	return application.NewPaginatedResponse(skuDTOs, query.Page, query.PageSize, total), nil
}

// HandleListSKUsByProduct handles the list SKUs by product query
func (h *SKUQueryHandler) HandleListSKUsByProduct(ctx context.Context, query *ListSKUsByProductQuery) ([]*application.SKUDTO, error) {
	// Get from repository
	skus, err := h.repo.FindByProductID(ctx, query.ProductID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list SKUs by product")
	}

	// Convert to DTOs
	skuDTOs := make([]*application.SKUDTO, len(skus))
	for i, sku := range skus {
		skuDTOs[i] = application.ToSKUDTO(sku)
	}

	return skuDTOs, nil
}

// skuCacheKey generates a cache key for a SKU
func skuCacheKey(id int64) string {
	return cache.Key("catalog", "sku", id)
}
