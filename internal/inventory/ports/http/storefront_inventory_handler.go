package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/inventory/application/queries"
	"github.com/qhato/ecommerce/pkg/logger"
)

// StorefrontInventoryHandler handles storefront inventory HTTP requests (read-only)
type StorefrontInventoryHandler struct {
	queryService *queries.InventoryQueryService
	log          logger.Logger
}

// NewStorefrontInventoryHandler creates a new storefront inventory HTTP handler
func NewStorefrontInventoryHandler(
	queryService *queries.InventoryQueryService,
	log logger.Logger,
) *StorefrontInventoryHandler {
	return &StorefrontInventoryHandler{
		queryService: queryService,
		log:          log,
	}
}

// RegisterRoutes registers all storefront inventory routes
func (h *StorefrontInventoryHandler) RegisterRoutes(r chi.Router) {
	r.Route("/inventory", func(r chi.Router) {
		// Public endpoints for storefront
		r.Post("/check-availability", h.CheckAvailability)
		r.Get("/sku/{skuId}/availability", h.GetSKUAvailability)
	})
}

// CheckAvailability checks if requested quantity is available
func (h *StorefrontInventoryHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SKUID    string `json:"sku_id"`
		Quantity int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.SKUID == "" {
		respondWithError(w, http.StatusBadRequest, "SKU ID is required")
		return
	}
	if req.Quantity <= 0 {
		respondWithError(w, http.StatusBadRequest, "Quantity must be greater than zero")
		return
	}

	query := queries.CheckInventoryAvailabilityQuery{
		SKUID:    req.SKUID,
		Quantity: req.Quantity,
	}

	availability, err := h.queryService.CheckInventoryAvailability(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to check inventory availability")
		respondWithError(w, http.StatusInternalServerError, "Failed to check availability")
		return
	}

	respondWithJSON(w, http.StatusOK, availability)
}

// GetSKUAvailability gets availability info for a SKU
func (h *StorefrontInventoryHandler) GetSKUAvailability(w http.ResponseWriter, r *http.Request) {
	skuID := chi.URLParam(r, "skuId")

	if skuID == "" {
		respondWithError(w, http.StatusBadRequest, "SKU ID is required")
		return
	}

	query := queries.GetInventoryBySKUQuery{SKUID: skuID}
	level, err := h.queryService.GetInventoryBySKU(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get inventory by SKU")
		respondWithError(w, http.StatusNotFound, "Inventory not found")
		return
	}

	// Return limited info for storefront (don't expose warehouse details)
	response := map[string]interface{}{
		"sku_id":         level.SKUID,
		"available":      level.QuantityAvailable,
		"can_backorder":  level.AllowBackorder,
		"can_preorder":   level.AllowPreorder,
		"is_in_stock":    level.QuantityAvailable > 0,
	}

	respondWithJSON(w, http.StatusOK, response)
}
