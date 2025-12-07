package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/shipping/application/queries"
	"github.com/qhato/ecommerce/pkg/logger"
)

// StorefrontShippingHandler handles storefront shipping HTTP requests (read-only)
type StorefrontShippingHandler struct {
	queryService *queries.ShippingQueryService
	log          logger.Logger
}

// NewStorefrontShippingHandler creates a new storefront shipping HTTP handler
func NewStorefrontShippingHandler(
	queryService *queries.ShippingQueryService,
	log logger.Logger,
) *StorefrontShippingHandler {
	return &StorefrontShippingHandler{
		queryService: queryService,
		log:          log,
	}
}

// RegisterRoutes registers all storefront shipping routes
func (h *StorefrontShippingHandler) RegisterRoutes(r chi.Router) {
	r.Route("/shipping", func(r chi.Router) {
		// Public endpoints for storefront
		r.Post("/calculate-rates", h.CalculateRates)
		r.Get("/methods/available", h.GetAvailableShippingMethods)
	})
}

// CalculateRates calculates shipping rates for checkout
func (h *StorefrontShippingHandler) CalculateRates(w http.ResponseWriter, r *http.Request) {
	var req queries.CalculateShippingRatesQuery
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.Weight.IsZero() || req.Weight.IsNegative() {
		respondWithError(w, http.StatusBadRequest, "Weight must be greater than zero")
		return
	}
	if req.OrderTotal.IsNegative() {
		respondWithError(w, http.StatusBadRequest, "Order total cannot be negative")
		return
	}
	if req.Country == "" {
		respondWithError(w, http.StatusBadRequest, "Country is required")
		return
	}

	rates, err := h.queryService.CalculateShippingRates(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Error("Failed to calculate shipping rates")
		respondWithError(w, http.StatusInternalServerError, "Failed to calculate shipping rates")
		return
	}

	respondWithJSON(w, http.StatusOK, rates)
}

// GetAvailableShippingMethods gets available shipping methods for a location
func (h *StorefrontShippingHandler) GetAvailableShippingMethods(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	zipCode := r.URL.Query().Get("zip_code")

	if country == "" {
		respondWithError(w, http.StatusBadRequest, "Country is required")
		return
	}

	query := queries.GetAvailableShippingMethodsQuery{
		Country: country,
		ZipCode: zipCode,
	}

	methods, err := h.queryService.GetAvailableShippingMethods(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get available shipping methods")
		respondWithError(w, http.StatusInternalServerError, "Failed to get available shipping methods")
		return
	}

	respondWithJSON(w, http.StatusOK, methods)
}
