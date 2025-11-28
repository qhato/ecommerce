package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/fulfillment/application"
	"github.com/qhato/ecommerce/internal/fulfillment/application/commands"
	"github.com/qhato/ecommerce/internal/fulfillment/domain"
	httpPkg "github.com/qhato/ecommerce/pkg/http"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/validator"
)

// AdminShipmentHandler handles admin shipment HTTP requests
type AdminShipmentHandler struct {
	commandHandler *commands.ShipmentCommandHandler
	repo           domain.ShipmentRepository
	validator      *validator.Validator
	log            *logger.Logger
}

// NewAdminShipmentHandler creates a new AdminShipmentHandler
func NewAdminShipmentHandler(
	commandHandler *commands.ShipmentCommandHandler,
	repo domain.ShipmentRepository,
	validator *validator.Validator,
	log *logger.Logger,
) *AdminShipmentHandler {
	return &AdminShipmentHandler{
		commandHandler: commandHandler,
		repo:           repo,
		validator:      validator,
		log:            log,
	}
}

// RegisterRoutes registers shipment routes
func (h *AdminShipmentHandler) RegisterRoutes(r chi.Router) {
	r.Route("/shipments", func(r chi.Router) {
		r.Post("/", h.CreateShipment)
		r.Get("/", h.ListShipments)
		r.Get("/{id}", h.GetShipment)
		r.Post("/{id}/ship", h.ShipShipment)
		r.Post("/{id}/deliver", h.DeliverShipment)
		r.Post("/{id}/cancel", h.CancelShipment)
		r.Put("/{id}/tracking", h.UpdateTracking)
		r.Get("/order/{orderId}", h.GetShipmentsByOrder)
		r.Get("/tracking/{trackingNumber}", h.GetShipmentByTracking)
	})
}

// CreateShipment creates a new shipment
func (h *AdminShipmentHandler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	var req application.CreateShipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	// Convert DTO address to domain address
	address := domain.Address{
		Name:       req.ShippingAddress.Name,
		Line1:      req.ShippingAddress.Line1,
		Line2:      req.ShippingAddress.Line2,
		City:       req.ShippingAddress.City,
		State:      req.ShippingAddress.State,
		PostalCode: req.ShippingAddress.PostalCode,
		Country:    req.ShippingAddress.Country,
		Phone:      req.ShippingAddress.Phone,
	}

	shipment, err := h.commandHandler.CreateShipment(
		r.Context(),
		req.OrderID,
		req.Carrier,
		req.ShippingMethod,
		req.ShippingCost,
		address,
	)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to create shipment", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusCreated, application.ToShipmentDTO(shipment))
}

// GetShipment retrieves a shipment by ID
func (h *AdminShipmentHandler) GetShipment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid shipment ID", err)
		return
	}

	shipment, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to get shipment", err)
		return
	}
	if shipment == nil {
		httpPkg.RespondError(w, http.StatusNotFound, "shipment not found", nil)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToShipmentDTO(shipment))
}

// GetShipmentByTracking retrieves a shipment by tracking number
func (h *AdminShipmentHandler) GetShipmentByTracking(w http.ResponseWriter, r *http.Request) {
	trackingNumber := chi.URLParam(r, "trackingNumber")

	shipment, err := h.repo.FindByTrackingNumber(r.Context(), trackingNumber)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to get shipment", err)
		return
	}
	if shipment == nil {
		httpPkg.RespondError(w, http.StatusNotFound, "shipment not found", nil)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToShipmentDTO(shipment))
}

// GetShipmentsByOrder retrieves shipments for an order
func (h *AdminShipmentHandler) GetShipmentsByOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "orderId")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid order ID", err)
		return
	}

	shipments, err := h.repo.FindByOrderID(r.Context(), orderID)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to list shipments", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToShipmentDTOs(shipments))
}

// ListShipments lists all shipments
func (h *AdminShipmentHandler) ListShipments(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	status := r.URL.Query().Get("status")
	carrier := r.URL.Query().Get("carrier")
	orderIDStr := r.URL.Query().Get("order_id")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	var orderID int64
	if orderIDStr != "" {
		orderID, _ = strconv.ParseInt(orderIDStr, 10, 64)
	}

	filter := &domain.ShipmentFilter{
		Page:      page,
		PageSize:  pageSize,
		Status:    domain.ShipmentStatus(status),
		Carrier:   carrier,
		OrderID:   orderID,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	shipments, total, err := h.repo.FindAll(r.Context(), filter)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to list shipments", err)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response := map[string]interface{}{
		"data":        application.ToShipmentDTOs(shipments),
		"page":        page,
		"page_size":   pageSize,
		"total_items": total,
		"total_pages": totalPages,
	}

	httpPkg.RespondJSON(w, http.StatusOK, response)
}

// ShipShipment marks a shipment as shipped
func (h *AdminShipmentHandler) ShipShipment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid shipment ID", err)
		return
	}

	var req application.ShipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	err = h.commandHandler.ShipShipment(r.Context(), id, req.TrackingNumber)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to ship shipment", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "shipment marked as shipped successfully"})
}

// DeliverShipment marks a shipment as delivered
func (h *AdminShipmentHandler) DeliverShipment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid shipment ID", err)
		return
	}

	err = h.commandHandler.DeliverShipment(r.Context(), id)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to deliver shipment", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "shipment marked as delivered successfully"})
}

// CancelShipment cancels a shipment
func (h *AdminShipmentHandler) CancelShipment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid shipment ID", err)
		return
	}

	err = h.commandHandler.CancelShipment(r.Context(), id)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to cancel shipment", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "shipment cancelled successfully"})
}

// UpdateTracking updates shipment tracking information
func (h *AdminShipmentHandler) UpdateTracking(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid shipment ID", err)
		return
	}

	var req application.UpdateTrackingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	err = h.commandHandler.UpdateTracking(r.Context(), id, req.TrackingNumber, req.Notes)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to update tracking", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "tracking updated successfully"})
}
