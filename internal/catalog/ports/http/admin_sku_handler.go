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

// AdminSKUHandler handles admin SKU HTTP requests
type AdminSKUHandler struct {
	commandHandler *commands.SKUCommandHandler
	queryHandler   *queries.SKUQueryHandler
	logger         *logger.Logger
}

// NewAdminSKUHandler creates a new admin SKU handler
func NewAdminSKUHandler(
	commandHandler *commands.SKUCommandHandler,
	queryHandler *queries.SKUQueryHandler,
	logger *logger.Logger,
) *AdminSKUHandler {
	return &AdminSKUHandler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		logger:         logger,
	}
}

// RegisterRoutes registers admin SKU routes
func (h *AdminSKUHandler) RegisterRoutes(r chi.Router) {
	r.Route("/admin/skus", func(r chi.Router) {
		r.Post("/", h.CreateSKU)
		r.Get("/", h.ListSKUs)
		r.Get("/{id}", h.GetSKU)
		r.Put("/{id}", h.UpdateSKU)
		r.Delete("/{id}", h.DeleteSKU)
		r.Put("/{id}/pricing", h.UpdateSKUPricing)
		r.Put("/{id}/availability", h.UpdateSKUAvailability)
		r.Get("/upc/{upc}", h.GetSKUByUPC)
		r.Get("/product/{product_id}", h.ListSKUsByProduct)
	})
}

// CreateSKU creates a new SKU
func (h *AdminSKUHandler) CreateSKU(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateSKUCommand
	if err := pkghttp.DecodeJSON(r, &cmd); err != nil {
		pkghttp.RespondError(w, err)
		return
	}

	skuID, err := h.commandHandler.HandleCreateSKU(r.Context(), &cmd)
	if err != nil {
		h.logger.Error("failed to create SKU", "error", err)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"id": skuID,
	})
}

// ListSKUs lists all SKUs with pagination
func (h *AdminSKUHandler) ListSKUs(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	availableOnly := r.URL.Query().Get("available_only") == "true"
	activeOnly := r.URL.Query().Get("active_only") == "true"
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.ListSKUsQuery{
		Page:          page,
		PageSize:      pageSize,
		AvailableOnly: availableOnly,
		ActiveOnly:    activeOnly,
		SortBy:        sortBy,
		SortOrder:     sortOrder,
	}

	result, err := h.queryHandler.HandleListSKUs(r.Context(), query)
	if err != nil {
		h.logger.Error("failed to list SKUs", "error", err)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, result)
}

// GetSKU retrieves a SKU by ID
func (h *AdminSKUHandler) GetSKU(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid SKU ID"))
		return
	}

	query := &queries.GetSKUByIDQuery{ID: id}
	sku, err := h.queryHandler.HandleGetSKUByID(r.Context(), query)
	if err != nil {
		h.logger.Error("failed to get SKU", "error", err, "sku_id", id)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, sku)
}

// UpdateSKU updates an existing SKU
func (h *AdminSKUHandler) UpdateSKU(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid SKU ID"))
		return
	}

	var cmd commands.UpdateSKUCommand
	if err := pkghttp.DecodeJSON(r, &cmd); err != nil {
		pkghttp.RespondError(w, err)
		return
	}
	cmd.ID = id

	if err := h.commandHandler.HandleUpdateSKU(r.Context(), &cmd); err != nil {
		h.logger.Error("failed to update SKU", "error", err, "sku_id", id)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "SKU updated successfully",
	})
}

// DeleteSKU deletes a SKU
func (h *AdminSKUHandler) DeleteSKU(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid SKU ID"))
		return
	}

	cmd := &commands.DeleteSKUCommand{ID: id}
	if err := h.commandHandler.HandleDeleteSKU(r.Context(), cmd); err != nil {
		h.logger.Error("failed to delete SKU", "error", err, "sku_id", id)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "SKU deleted successfully",
	})
}

// UpdateSKUPricing updates SKU pricing
func (h *AdminSKUHandler) UpdateSKUPricing(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid SKU ID"))
		return
	}

	var cmd commands.UpdateSKUPricingCommand
	if err := pkghttp.DecodeJSON(r, &cmd); err != nil {
		pkghttp.RespondError(w, err)
		return
	}
	cmd.ID = id

	if err := h.commandHandler.HandleUpdateSKUPricing(r.Context(), &cmd); err != nil {
		h.logger.Error("failed to update SKU pricing", "error", err, "sku_id", id)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "SKU pricing updated successfully",
	})
}

// UpdateSKUAvailability updates SKU availability
func (h *AdminSKUHandler) UpdateSKUAvailability(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid SKU ID"))
		return
	}

	var cmd commands.UpdateSKUAvailabilityCommand
	if err := pkghttp.DecodeJSON(r, &cmd); err != nil {
		pkghttp.RespondError(w, err)
		return
	}
	cmd.ID = id

	if err := h.commandHandler.HandleUpdateSKUAvailability(r.Context(), &cmd); err != nil {
		h.logger.Error("failed to update SKU availability", "error", err, "sku_id", id)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "SKU availability updated successfully",
	})
}

// GetSKUByUPC retrieves a SKU by UPC
func (h *AdminSKUHandler) GetSKUByUPC(w http.ResponseWriter, r *http.Request) {
	upc := chi.URLParam(r, "upc")
	if upc == "" {
		pkghttp.RespondError(w, pkghttp.NewValidationError("UPC is required"))
		return
	}

	query := &queries.GetSKUByUPCQuery{UPC: upc}
	sku, err := h.queryHandler.HandleGetSKUByUPC(r.Context(), query)
	if err != nil {
		h.logger.Error("failed to get SKU by UPC", "error", err, "upc", upc)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, sku)
}

// ListSKUsByProduct lists SKUs by product ID
func (h *AdminSKUHandler) ListSKUsByProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid product ID"))
		return
	}

	query := &queries.ListSKUsByProductQuery{ProductID: productID}
	skus, err := h.queryHandler.HandleListSKUsByProduct(r.Context(), query)
	if err != nil {
		h.logger.Error("failed to list SKUs by product", "error", err, "product_id", productID)
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, skus)
}
