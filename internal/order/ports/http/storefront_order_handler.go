package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/order/application"
	"github.com/qhato/ecommerce/internal/order/application/queries"
	"github.com/qhato/ecommerce/internal/order/domain"
	httpPkg "github.com/qhato/ecommerce/pkg/http"
	"github.com/qhato/ecommerce/pkg/logger"
)

// StorefrontOrderHandler handles storefront order HTTP requests
type StorefrontOrderHandler struct {
	queryHandler *queries.OrderQueryHandler
	log          *logger.Logger
}

// NewStorefrontOrderHandler creates a new StorefrontOrderHandler
func NewStorefrontOrderHandler(
	queryHandler *queries.OrderQueryHandler,
	log *logger.Logger,
) *StorefrontOrderHandler {
	return &StorefrontOrderHandler{
		queryHandler: queryHandler,
		log:          log,
	}
}

// RegisterRoutes registers storefront order routes
func (h *StorefrontOrderHandler) RegisterRoutes(r chi.Router) {
	r.Route("/orders", func(r chi.Router) {
		r.Get("/{id}", h.GetOrder)
		r.Get("/number/{orderNumber}", h.GetOrderByNumber)
		r.Get("/customer/{customerId}", h.ListCustomerOrders)
	})
}

// GetOrder retrieves an order by ID
func (h *StorefrontOrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid order ID", err)
		return
	}

	order, err := h.queryHandler.GetByID(r.Context(), id)
	if err != nil {
		httpPkg.RespondError(w, http.StatusNotFound, "order not found", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToOrderDTO(order))
}

// GetOrderByNumber retrieves an order by order number
func (h *StorefrontOrderHandler) GetOrderByNumber(w http.ResponseWriter, r *http.Request) {
	orderNumber := chi.URLParam(r, "orderNumber")
	if orderNumber == "" {
		httpPkg.RespondError(w, http.StatusBadRequest, "order number is required", nil)
		return
	}

	order, err := h.queryHandler.GetByOrderNumber(r.Context(), orderNumber)
	if err != nil {
		httpPkg.RespondError(w, http.StatusNotFound, "order not found", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToOrderDTO(order))
}

// ListCustomerOrders lists orders for a specific customer
func (h *StorefrontOrderHandler) ListCustomerOrders(w http.ResponseWriter, r *http.Request) {
	customerIDStr := chi.URLParam(r, "customerId")
	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid customer ID", err)
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	status := r.URL.Query().Get("status")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	filter := &domain.OrderFilter{
		Page:       page,
		PageSize:   pageSize,
		Status:     domain.OrderStatus(status),
		CustomerID: customerID,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	orders, total, err := h.queryHandler.ListByCustomer(r.Context(), customerID, filter)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to list customer orders", err)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response := application.PaginatedOrderResponse{
		Data:       application.ToOrderDTOs(orders),
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	httpPkg.RespondJSON(w, http.StatusOK, response)
}
