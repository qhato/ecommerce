package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/fulfillment/application"
	"github.com/qhato/ecommerce/internal/fulfillment/domain"
	httpPkg "github.com/qhato/ecommerce/pkg/http"
	"github.com/qhato/ecommerce/pkg/logger"
)

// StorefrontShipmentHandler handles storefront shipment HTTP requests
type StorefrontShipmentHandler struct {
	repo domain.ShipmentRepository
	log  *logger.Logger
}

// NewStorefrontShipmentHandler creates a new StorefrontShipmentHandler
func NewStorefrontShipmentHandler(
	repo domain.ShipmentRepository,
	log *logger.Logger,
) *StorefrontShipmentHandler {
	return &StorefrontShipmentHandler{
		repo: repo,
		log:  log,
	}
}

// RegisterRoutes registers storefront shipment routes
func (h *StorefrontShipmentHandler) RegisterRoutes(r chi.Router) {
	r.Route("/shipments", func(r chi.Router) {
		r.Get("/tracking/{trackingNumber}", h.TrackShipment)
		r.Get("/order/{orderId}", h.GetShipmentsByOrder)
	})
}

// TrackShipment tracks a shipment by tracking number
func (h *StorefrontShipmentHandler) TrackShipment(w http.ResponseWriter, r *http.Request) {
	trackingNumber := chi.URLParam(r, "trackingNumber")
	if trackingNumber == "" {
		httpPkg.RespondError(w, http.StatusBadRequest, "tracking number is required", nil)
		return
	}

	shipment, err := h.repo.FindByTrackingNumber(r.Context(), trackingNumber)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to track shipment", err)
		return
	}
	if shipment == nil {
		httpPkg.RespondError(w, http.StatusNotFound, "shipment not found", nil)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToShipmentDTO(shipment))
}

// GetShipmentsByOrder retrieves shipments for an order
func (h *StorefrontShipmentHandler) GetShipmentsByOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "orderId")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid order ID", err)
		return
	}

	// TODO: In production, verify that the authenticated user owns this order

	shipments, err := h.repo.FindByOrderID(r.Context(), orderID)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to list shipments", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToShipmentDTOs(shipments))
}
