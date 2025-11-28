package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/application"
	"github.com/qhato/ecommerce/internal/catalog/domain"
	"github.com/qhato/ecommerce/pkg/cache"
	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/logger"
)

// GetCategoryByIDQuery represents a query to get a category by ID
type GetCategoryByIDQuery struct {
	ID int64 `json:"id" validate:"required"`
}

// GetCategoryByURLQuery represents a query to get a category by URL
type GetCategoryByURLQuery struct {
	URL string `json:"url" validate:"required"`
}

// ListCategoriesQuery represents a query to list categories
type ListCategoriesQuery struct {
	Page            int    `json:"page" validate:"min=1"`
	PageSize        int    `json:"page_size" validate:"min=1,max=100"`
	IncludeArchived bool   `json:"include_archived"`
	ActiveOnly      bool   `json:"active_only"`
	SortBy          string `json:"sort_by"`
	SortOrder       string `json:"sort_order"`
}

// ListCategoriesByParentQuery represents a query to list categories by parent
type ListCategoriesByParentQuery struct {
	ParentID        int64  `json:"parent_id" validate:"required"`
	Page            int    `json:"page" validate:"min=1"`
	PageSize        int    `json:"page_size" validate:"min=1,max=100"`
	IncludeArchived bool   `json:"include_archived"`
	ActiveOnly      bool   `json:"active_only"`
	SortBy          string `json:"sort_by"`
	SortOrder       string `json:"sort_order"`
}

// ListRootCategoriesQuery represents a query to list root categories
type ListRootCategoriesQuery struct {
	Page            int    `json:"page" validate:"min=1"`
	PageSize        int    `json:"page_size" validate:"min=1,max=100"`
	IncludeArchived bool   `json:"include_archived"`
	ActiveOnly      bool   `json:"active_only"`
	SortBy          string `json:"sort_by"`
	SortOrder       string `json:"sort_order"`
}

// GetCategoryPathQuery represents a query to get the category path
type GetCategoryPathQuery struct {
	CategoryID int64 `json:"category_id" validate:"required"`
}

// CategoryQueryHandler handles category queries
type CategoryQueryHandler struct {
	repo   domain.CategoryRepository
	cache  cache.Cache
	logger *logger.Logger
}

// NewCategoryQueryHandler creates a new category query handler
func NewCategoryQueryHandler(
	repo domain.CategoryRepository,
	cache cache.Cache,
	logger *logger.Logger,
) *CategoryQueryHandler {
	return &CategoryQueryHandler{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

// HandleGetCategoryByID handles the get category by ID query
func (h *CategoryQueryHandler) HandleGetCategoryByID(ctx context.Context, query *GetCategoryByIDQuery) (*application.CategoryDTO, error) {
	// Try to get from cache first
	cacheKey := categoryCacheKey(query.ID)
	var category *domain.Category

	if cached, err := h.cache.Get(ctx, cacheKey); err == nil && len(cached) > 0 {
		if err := json.Unmarshal(cached, &category); err == nil {
			h.logger.WithField("category_id", query.ID).Debug("category found in cache")
			return application.ToCategoryDTO(category), nil
		}
	}

	// Get from repository
	category, err := h.repo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, errors.InternalWrap(err, "category not found")
	}

	// Cache the result
	if data, err := json.Marshal(category); err == nil {
		if err := h.cache.Set(ctx, cacheKey, data, 5*time.Minute); err != nil {
			h.logger.WithField("category_id", query.ID).WithError(err).Warn("failed to cache category")
		}
	}

	return application.ToCategoryDTO(category), nil
}

// HandleGetCategoryByURL handles the get category by URL query
func (h *CategoryQueryHandler) HandleGetCategoryByURL(ctx context.Context, query *GetCategoryByURLQuery) (*application.CategoryDTO, error) {
	category, err := h.repo.FindByURL(ctx, query.URL)
	if err != nil {
		return nil, errors.Wrap(err, "category not found")
	}

	// Cache the result
	// Cache the result
	cacheKey := categoryCacheKey(category.ID)
	if data, err := json.Marshal(category); err == nil {
		if err := h.cache.Set(ctx, cacheKey, data, 5*time.Minute); err != nil {
			h.logger.WithField("category_id", category.ID).WithError(err).Warn("failed to cache category")
		}
	}

	return application.ToCategoryDTO(category), nil
}

// HandleListCategories handles the list categories query
func (h *CategoryQueryHandler) HandleListCategories(ctx context.Context, query *ListCategoriesQuery) (*application.PaginatedResponse, error) {
	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	if query.SortBy == "" {
		query.SortBy = "display_order"
	}
	if query.SortOrder == "" {
		query.SortOrder = "asc"
	}

	// Create filter
	filter := &domain.CategoryFilter{
		Page:            query.Page,
		PageSize:        query.PageSize,
		IncludeArchived: query.IncludeArchived,
		ActiveOnly:      query.ActiveOnly,
		SortBy:          query.SortBy,
		SortOrder:       query.SortOrder,
	}

	// Get from repository
	categories, total, err := h.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to list categories")
	}

	// Convert to DTOs
	categoryDTOs := make([]*application.CategoryDTO, len(categories))
	for i, category := range categories {
		categoryDTOs[i] = application.ToCategoryDTO(category)
	}

	return application.NewPaginatedResponse(categoryDTOs, query.Page, query.PageSize, total), nil
}

// HandleListCategoriesByParent handles the list categories by parent query
func (h *CategoryQueryHandler) HandleListCategoriesByParent(ctx context.Context, query *ListCategoriesByParentQuery) (*application.PaginatedResponse, error) {
	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	if query.SortBy == "" {
		query.SortBy = "display_order"
	}
	if query.SortOrder == "" {
		query.SortOrder = "asc"
	}

	// Create filter
	filter := &domain.CategoryFilter{
		Page:            query.Page,
		PageSize:        query.PageSize,
		IncludeArchived: query.IncludeArchived,
		ActiveOnly:      query.ActiveOnly,
		SortBy:          query.SortBy,
		SortOrder:       query.SortOrder,
	}

	// Get from repository
	categories, total, err := h.repo.FindByParentID(ctx, query.ParentID, filter)
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to list categories by parent")
	}

	// Convert to DTOs
	categoryDTOs := make([]*application.CategoryDTO, len(categories))
	for i, category := range categories {
		categoryDTOs[i] = application.ToCategoryDTO(category)
	}

	return application.NewPaginatedResponse(categoryDTOs, query.Page, query.PageSize, total), nil
}

// HandleListRootCategories handles the list root categories query
func (h *CategoryQueryHandler) HandleListRootCategories(ctx context.Context, query *ListRootCategoriesQuery) (*application.PaginatedResponse, error) {
	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	if query.SortBy == "" {
		query.SortBy = "display_order"
	}
	if query.SortOrder == "" {
		query.SortOrder = "asc"
	}

	// Create filter
	filter := &domain.CategoryFilter{
		Page:            query.Page,
		PageSize:        query.PageSize,
		IncludeArchived: query.IncludeArchived,
		ActiveOnly:      query.ActiveOnly,
		SortBy:          query.SortBy,
		SortOrder:       query.SortOrder,
	}

	// Get from repository
	categories, total, err := h.repo.FindRootCategories(ctx, filter)
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to list root categories")
	}

	// Convert to DTOs
	categoryDTOs := make([]*application.CategoryDTO, len(categories))
	for i, category := range categories {
		categoryDTOs[i] = application.ToCategoryDTO(category)
	}

	return application.NewPaginatedResponse(categoryDTOs, query.Page, query.PageSize, total), nil
}

// HandleGetCategoryPath handles the get category path query
func (h *CategoryQueryHandler) HandleGetCategoryPath(ctx context.Context, query *GetCategoryPathQuery) ([]*application.CategoryDTO, error) {
	// Get category path from repository
	categories, err := h.repo.GetCategoryPath(ctx, query.CategoryID)
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to get category path")
	}

	// Convert to DTOs
	categoryDTOs := make([]*application.CategoryDTO, len(categories))
	for i, category := range categories {
		categoryDTOs[i] = application.ToCategoryDTO(category)
	}

	return categoryDTOs, nil
}

// categoryCacheKey generates a cache key for a category
	return fmt.Sprintf("catalog:category:%d", id)
