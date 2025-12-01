package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/catalog/application"
	"github.com/qhato/ecommerce/internal/catalog/application/queries"
	pkghttp "github.com/qhato/ecommerce/pkg/http"
	"github.com/qhato/ecommerce/pkg/logger"
)

// StorefrontCatalogHandler handles public storefront catalog HTTP requests (read-only)
type StorefrontCatalogHandler struct {
	productQueryHandler  *queries.ProductQueryHandler
	categoryQueryHandler *queries.CategoryQueryHandler
	skuQueryHandler      *queries.SKUQueryHandler
	logger               *logger.Logger
}

// NewStorefrontCatalogHandler creates a new storefront catalog handler
func NewStorefrontCatalogHandler(
	productQueryHandler *queries.ProductQueryHandler,
	categoryQueryHandler *queries.CategoryQueryHandler,
	skuQueryHandler *queries.SKUQueryHandler,
	logger *logger.Logger,
) *StorefrontCatalogHandler {
	return &StorefrontCatalogHandler{
		productQueryHandler:  productQueryHandler,
		categoryQueryHandler: categoryQueryHandler,
		skuQueryHandler:      skuQueryHandler,
		logger:               logger,
	}
}

// RegisterRoutes registers storefront catalog routes
func (h *StorefrontCatalogHandler) RegisterRoutes(r chi.Router) {
	r.Route("/catalog", func(r chi.Router) {
		// Product routes
		r.Get("/products", h.ListProducts)
		r.Get("/products/{id}", h.GetProduct)
		r.Get("/products/url/{url}", h.GetProductByURL)
		r.Get("/products/search", h.SearchProducts)

		// Category routes
		r.Get("/categories", h.ListRootCategories)
		r.Get("/categories/{id}", h.GetCategory)
		r.Get("/categories/url/{url}", h.GetCategoryByURL)
		r.Get("/categories/{id}/children", h.ListChildCategories)
		r.Get("/categories/{id}/products", h.ListProductsByCategory)
		r.Get("/categories/{id}/path", h.GetCategoryPath)

		// SKU routes
		r.Get("/skus", h.ListSKUs)
		r.Get("/skus/{id}", h.GetSKU)
		r.Get("/skus/upc/{upc}", h.GetSKUByUPC)
		r.Get("/skus/product/{product_id}", h.ListSKUsByProduct)
	})
}

// Product Handlers

// ListProducts lists all active products with pagination
func (h *StorefrontCatalogHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.ListProductsQuery{
		Page:            page,
		PageSize:        pageSize,
		IncludeArchived: false, // Storefront never shows archived products
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	result, err := h.productQueryHandler.HandleListProducts(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).Error("failed to list products")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, result)
}

// GetProduct retrieves a product by ID
func (h *StorefrontCatalogHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid product ID"))
		return
	}

	query := &queries.GetProductByIDQuery{ID: id}
	product, err := h.productQueryHandler.HandleGetProductByID(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("product_id", id).Error("failed to get product")
		pkghttp.RespondError(w, err)
		return
	}

	// Check if archived (storefront shouldn't show archived products)
	if product.Archived {
		pkghttp.RespondError(w, pkghttp.NewNotFoundError("product not found"))
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, product)
}

// GetProductByURL retrieves a product by URL
func (h *StorefrontCatalogHandler) GetProductByURL(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")
	if url == "" {
		pkghttp.RespondError(w, pkghttp.NewValidationError("URL is required"))
		return
	}

	query := &queries.GetProductByURLQuery{URL: url}
	product, err := h.productQueryHandler.HandleGetProductByURL(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("url", url).Error("failed to get product by URL")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, product)
}

// SearchProducts searches for products
func (h *StorefrontCatalogHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")
	if searchQuery == "" {
		pkghttp.RespondError(w, pkghttp.NewValidationError("search query is required"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.SearchProductsQuery{
		Query:           searchQuery,
		Page:            page,
		PageSize:        pageSize,
		IncludeArchived: false, // Storefront never shows archived products
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	result, err := h.productQueryHandler.HandleSearchProducts(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).Error("failed to search products")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, result)
}

// ListProductsByCategory lists products by category
func (h *StorefrontCatalogHandler) ListProductsByCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid category ID"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.ListProductsByCategoryQuery{
		CategoryID:      id,
		Page:            page,
		PageSize:        pageSize,
		IncludeArchived: false,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	result, err := h.productQueryHandler.HandleListProductsByCategory(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("category_id", id).Error("failed to list products by category")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, result)
}

// Category Handlers

// ListRootCategories lists active root categories
func (h *StorefrontCatalogHandler) ListRootCategories(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.ListRootCategoriesQuery{
		Page:            page,
		PageSize:        pageSize,
		IncludeArchived: false, // Storefront never shows archived categories
		ActiveOnly:      true,  // Only active categories
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	result, err := h.categoryQueryHandler.HandleListRootCategories(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).Error("failed to list root categories")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, result)
}

// GetCategory retrieves a category by ID
func (h *StorefrontCatalogHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid category ID"))
		return
	}

	query := &queries.GetCategoryByIDQuery{ID: id}
	category, err := h.categoryQueryHandler.HandleGetCategoryByID(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("category_id", id).Error("failed to get category")
		pkghttp.RespondError(w, err)
		return
	}

	// Check if archived or inactive (storefront shouldn't show them)
	if category.Archived || !category.IsActive {
		pkghttp.RespondError(w, pkghttp.NewNotFoundError("category not found"))
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, category)
}

// GetCategoryByURL retrieves a category by URL
func (h *StorefrontCatalogHandler) GetCategoryByURL(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")
	if url == "" {
		pkghttp.RespondError(w, pkghttp.NewValidationError("URL is required"))
		return
	}

	query := &queries.GetCategoryByURLQuery{URL: url}
	category, err := h.categoryQueryHandler.HandleGetCategoryByURL(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("url", url).Error("failed to get category by URL")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, category)
}

// ListChildCategories lists active child categories
func (h *StorefrontCatalogHandler) ListChildCategories(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid category ID"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.ListCategoriesByParentQuery{
		ParentID:        id,
		Page:            page,
		PageSize:        pageSize,
		IncludeArchived: false,
		ActiveOnly:      true,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	result, err := h.categoryQueryHandler.HandleListCategoriesByParent(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("parent_id", id).Error("failed to list child categories")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, result)
}

// GetCategoryPath retrieves the full path from root to category
func (h *StorefrontCatalogHandler) GetCategoryPath(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid category ID"))
		return
	}

	query := &queries.GetCategoryPathQuery{CategoryID: id}
	path, err := h.categoryQueryHandler.HandleGetCategoryPath(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("category_id", id).Error("failed to get category path")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, path)
}

// SKU Handlers

// ListSKUs lists all active and available SKUs with pagination
func (h *StorefrontCatalogHandler) ListSKUs(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.ListSKUsQuery{
		Page:          page,
		PageSize:      pageSize,
		AvailableOnly: true, // Storefront only shows available SKUs
		ActiveOnly:    true, // Only active SKUs
		SortBy:        sortBy,
		SortOrder:     sortOrder,
	}

	result, err := h.skuQueryHandler.HandleListSKUs(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).Error("failed to list SKUs")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, result)
}

// GetSKU retrieves a SKU by ID
func (h *StorefrontCatalogHandler) GetSKU(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid SKU ID"))
		return
	}

	query := &queries.GetSKUByIDQuery{ID: id}
	sku, err := h.skuQueryHandler.HandleGetSKUByID(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("sku_id", id).Error("failed to get SKU")
		pkghttp.RespondError(w, err)
		return
	}

	// Check if available and active (storefront shouldn't show unavailable SKUs)
	if !sku.Available || !sku.IsActive {
		pkghttp.RespondError(w, pkghttp.NewNotFoundError("SKU not found"))
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, sku)
}

// GetSKUByUPC retrieves a SKU by UPC
func (h *StorefrontCatalogHandler) GetSKUByUPC(w http.ResponseWriter, r *http.Request) {
	upc := chi.URLParam(r, "upc")
	if upc == "" {
		pkghttp.RespondError(w, pkghttp.NewValidationError("UPC is required"))
		return
	}

	query := &queries.GetSKUByUPCQuery{UPC: upc}
	sku, err := h.skuQueryHandler.HandleGetSKUByUPC(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("upc", upc).Error("failed to get SKU by UPC")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, sku)
}

// ListSKUsByProduct lists SKUs by product ID
func (h *StorefrontCatalogHandler) ListSKUsByProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid product ID"))
		return
	}

	query := &queries.ListSKUsByProductQuery{ProductID: productID}
	skus, err := h.skuQueryHandler.HandleListSKUsByProduct(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("product_id", productID).Error("failed to list SKUs by product")
		pkghttp.RespondError(w, err)
		return
	}

	// Filter only available and active SKUs for storefront
	var availableSKUs []*application.SkuDTO
	for _, sku := range skus {
		if sku.Available && sku.IsActive {
			availableSKUs = append(availableSKUs, sku)
		}
	}

	pkghttp.RespondJSON(w, http.StatusOK, availableSKUs)
}