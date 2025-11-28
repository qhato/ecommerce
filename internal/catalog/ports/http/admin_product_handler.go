package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/catalog/application/commands"
	"github.com/qhato/ecommerce/internal/catalog/application/queries"
	pkghttp "github.com/qhato/ecommerce/pkg/http"
	"github.com/qhato/ecommerce/pkg/logger"
)

// AdminProductHandler handles admin product HTTP requests
type AdminProductHandler struct {
	commandHandler *commands.ProductCommandHandler
	queryHandler   *queries.ProductQueryHandler
	logger         *logger.Logger
}

// NewAdminProductHandler creates a new admin product handler
func NewAdminProductHandler(
	commandHandler *commands.ProductCommandHandler,
	queryHandler *queries.ProductQueryHandler,
	logger *logger.Logger,
) *AdminProductHandler {
	return &AdminProductHandler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		logger:         logger,
	}
}

// RegisterRoutes registers admin product routes
func (h *AdminProductHandler) RegisterRoutes(r chi.Router) {
	r.Route("/admin/products", func(r chi.Router) {
		r.Post("/", h.CreateProduct)
		r.Get("/", h.ListProducts)
		r.Get("/{id}", h.GetProduct)
		r.Put("/{id}", h.UpdateProduct)
		r.Delete("/{id}", h.DeleteProduct)
		r.Post("/{id}/archive", h.ArchiveProduct)
		r.Get("/search", h.SearchProducts)
	})
}

// CreateProduct creates a new product
func (h *AdminProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateProductCommand
	if err := pkghttp.DecodeJSON(r, &cmd); err != nil {
		pkghttp.RespondError(w, err)
		return
	}

	productID, err := h.commandHandler.HandleCreateProduct(r.Context(), &cmd)
	if err != nil {
		h.logger.Error("failed to create product", "error", err)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"id": productID,
	})
}

// ListProducts lists all products with pagination
func (h *AdminProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	includeArchived := r.URL.Query().Get("include_archived") == "true"
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.ListProductsQuery{
		Page:            page,
		PageSize:        pageSize,
		IncludeArchived: includeArchived,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	result, err := h.queryHandler.HandleListProducts(r.Context(), query)
	if err != nil {
		h.logger.Error("failed to list products", "error", err)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, result)
}

// GetProduct retrieves a product by ID
func (h *AdminProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid product ID"))
		return
	}

	query := &queries.GetProductByIDQuery{ID: id}
	product, err := h.queryHandler.HandleGetProductByID(r.Context(), query)
	if err != nil {
		h.logger.Error("failed to get product", "error", err, "product_id", id)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, product)
}

// UpdateProduct updates an existing product
func (h *AdminProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid product ID"))
		return
	}

	var cmd commands.UpdateProductCommand
	if err := pkghttp.DecodeJSON(r, &cmd); err != nil {
		pkghttp.RespondError(w, err)
		return
	}
	cmd.ID = id

	if err := h.commandHandler.HandleUpdateProduct(r.Context(), &cmd); err != nil {
		h.logger.Error("failed to update product", "error", err, "product_id", id)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "product updated successfully",
	})
}

// DeleteProduct deletes a product
func (h *AdminProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid product ID"))
		return
	}

	cmd := &commands.DeleteProductCommand{ID: id}
	if err := h.commandHandler.HandleDeleteProduct(r.Context(), cmd); err != nil {
		h.logger.Error("failed to delete product", "error", err, "product_id", id)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "product deleted successfully",
	})
}

// ArchiveProduct archives a product
func (h *AdminProductHandler) ArchiveProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid product ID"))
		return
	}

	cmd := &commands.ArchiveProductCommand{ID: id}
	if err := h.commandHandler.HandleArchiveProduct(r.Context(), cmd); err != nil {
		h.logger.Error("failed to archive product", "error", err, "product_id", id)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "product archived successfully",
	})
}

// SearchProducts searches for products
func (h *AdminProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
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

	includeArchived := r.URL.Query().Get("include_archived") == "true"
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.SearchProductsQuery{
		Query:           searchQuery,
		Page:            page,
		PageSize:        pageSize,
		IncludeArchived: includeArchived,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	result, err := h.queryHandler.HandleSearchProducts(r.Context(), query)
	if err != nil {
		h.logger.Error("failed to search products", "error", err)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, result)
}
