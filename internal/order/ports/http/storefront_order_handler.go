package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/order/application/queries"
	"github.com/qhato/ecommerce/internal/order/domain"
	httpPkg "github.com/qhato/ecommerce/pkg/http"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/errors" // Import pkg/errors
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
		httpPkg.RespondError(w, errors.BadRequest("invalid order ID").WithInternal(err))
		return
	}

	query := &queries.GetOrderByIDQuery{ID: id}
	order, err := h.queryHandler.HandleGetOrderByID(r.Context(), query)
	if err != nil {
		if errors.IsNotFound(err) {
			httpPkg.RespondError(w, errors.NotFound(err.Error()))
		} else {
			httpPkg.RespondError(w, errors.Internal("failed to get order").WithInternal(err))
		}
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, order)
}

// GetOrderByNumber retrieves an order by order number
func (h *StorefrontOrderHandler) GetOrderByNumber(w http.ResponseWriter, r *http.Request) {
	orderNumber := chi.URLParam(r, "orderNumber")
	if orderNumber == "" {
		httpPkg.RespondError(w, errors.BadRequest("order number is required"))
		return
	}

	query := &queries.GetOrderByOrderNumberQuery{OrderNumber: orderNumber}
	order, err := h.queryHandler.HandleGetOrderByOrderNumber(r.Context(), query)
	if err != nil {
		if errors.IsNotFound(err) {
			httpPkg.RespondError(w, errors.NotFound(err.Error()))
		} else {
			httpPkg.RespondError(w, errors.Internal("failed to get order by number").WithInternal(err))
		}
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, order)
}

// ListCustomerOrders lists orders for a specific customer
func (h *StorefrontOrderHandler) ListCustomerOrders(w http.ResponseWriter, r *http.Request) {
	customerIDStr := chi.URLParam(r, "customerId")
	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid customer ID").WithInternal(err))
		return
	}
	customerIDPtr := &customerID // Use pointer for CustomerID in filter

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	statusStr := r.URL.Query().Get("status")
	var status *domain.OrderStatus
	if statusStr != "" {
		s := domain.OrderStatus(statusStr)
		status = &s
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.ListOrdersQuery{
		Page:       page,
		PageSize:   pageSize,
		Status:     status,
		CustomerID: customerIDPtr, // Use pointer
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	result, err := h.queryHandler.HandleListOrders(r.Context(), query)
	if err != nil {
		httpPkg.RespondError(w, errors.Internal("failed to list customer orders").WithInternal(err))
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, result)
}
