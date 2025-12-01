package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/order/application"
	"github.com/qhato/ecommerce/internal/order/application/commands"
	"github.com/qhato/ecommerce/internal/order/application/queries"
	"github.com/qhato/ecommerce/internal/order/domain"
	httpPkg "github.com/qhato/ecommerce/pkg/http"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/validator"
	"github.com/qhato/ecommerce/pkg/errors" // Import pkg/errors
)

// AdminOrderHandler handles admin order HTTP requests
type AdminOrderHandler struct {
	commandHandler *commands.OrderCommandHandler
	queryHandler   *queries.OrderQueryHandler
	validator      *validator.Validator
	log            *logger.Logger
}

// NewAdminOrderHandler creates a new AdminOrderHandler
func NewAdminOrderHandler(
	commandHandler *commands.OrderCommandHandler,
	queryHandler *queries.OrderQueryHandler,
	validator *validator.Validator,
	log *logger.Logger,
) *AdminOrderHandler {
	return &AdminOrderHandler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		validator:      validator,
		log:            log,
	}
}

// RegisterRoutes registers order routes
func (h *AdminOrderHandler) RegisterRoutes(r chi.Router) {
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", h.CreateOrder)
		r.Get("/", h.ListOrders)
		r.Get("/{id}", h.GetOrder)
		r.Put("/{id}/status", h.UpdateOrderStatus)
		r.Post("/{id}/submit", h.SubmitOrder)
		r.Post("/{id}/cancel", h.CancelOrder)
		r.Post("/{id}/items", h.AddOrderItem)
		r.Get("/number/{orderNumber}", h.GetOrderByNumber)
	})
}

// CreateOrder creates a new order
func (h *AdminOrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req application.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid request body").WithInternal(err))
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, errors.ValidationError("validation failed").WithInternal(err))
		return
	}

	cmd := &application.CreateOrderCommand{
		CustomerID: req.CustomerID,
		EmailAddress: req.EmailAddress,
		Name: req.Name,
		CurrencyCode: req.CurrencyCode,
		// Other fields as needed
	}

	order, err := h.commandHandler.HandleCreateOrder(
		r.Context(),
		cmd,
	)
	if err != nil {
		httpPkg.RespondError(w, errors.Internal("failed to create order").WithInternal(err))
		return
	}

	httpPkg.RespondJSON(w, http.StatusCreated, order) // Order is already DTO
}

// GetOrder retrieves an order by ID
func (h *AdminOrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
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

	httpPkg.RespondJSON(w, http.StatusOK, order) // Order is already DTO
}

// GetOrderByNumber retrieves an order by order number
func (h *AdminOrderHandler) GetOrderByNumber(w http.ResponseWriter, r *http.Request) {
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

	httpPkg.RespondJSON(w, http.StatusOK, order) // Order is already DTO
}

// ListOrders lists all orders with optional filtering
func (h *AdminOrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
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

	customerIDStr := r.URL.Query().Get("customer_id")
	var customerID *int64
	if customerIDStr != "" {
		id, err := strconv.ParseInt(customerIDStr, 10, 64)
		if err == nil {
			customerID = &id
		}
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.ListOrdersQuery{
		Page:       page,
		PageSize:   pageSize,
		Status:     status,
		CustomerID: customerID,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	result, err := h.queryHandler.HandleListOrders(r.Context(), query)
	if err != nil {
		httpPkg.RespondError(w, errors.Internal("failed to list orders").WithInternal(err))
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, result)
}

// UpdateOrderStatus updates the status of an order
func (h *AdminOrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid order ID").WithInternal(err))
		return
	}

	var req application.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid request body").WithInternal(err))
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, errors.ValidationError("validation failed").WithInternal(err))
		return
	}

	err = h.commandHandler.HandleUpdateOrderStatus(r.Context(), id, domain.OrderStatus(req.Status))
	if err != nil {
		if errors.IsNotFound(err) {
			httpPkg.RespondError(w, errors.NotFound(err.Error()))
		} else {
			httpPkg.RespondError(w, errors.Internal("failed to update order status").WithInternal(err))
		}
		return
	}

	// Invalidate cache
	h.queryHandler.InvalidateCache(r.Context(), id)

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "order status updated successfully"})
}

// SubmitOrder submits an order for processing
func (h *AdminOrderHandler) SubmitOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid order ID").WithInternal(err))
		return
	}

	err = h.commandHandler.HandleSubmitOrder(r.Context(), id)
	if err != nil {
		if errors.IsNotFound(err) {
			httpPkg.RespondError(w, errors.NotFound(err.Error()))
		} else {
			httpPkg.RespondError(w, errors.Internal("failed to submit order").WithInternal(err))
		}
		return
	}

	// Invalidate cache
	h.queryHandler.InvalidateCache(r.Context(), id)

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "order submitted successfully"})
}

// CancelOrder cancels an order
func (h *AdminOrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid order ID").WithInternal(err))
		return
	}

	err = h.commandHandler.HandleCancelOrder(r.Context(), id, "admin cancellation") // Added reason
	if err != nil {
		if errors.IsNotFound(err) {
			httpPkg.RespondError(w, errors.NotFound(err.Error()))
		} else {
			httpPkg.RespondError(w, errors.Internal("failed to cancel order").WithInternal(err))
		}
		return
	}

	// Invalidate cache
	h.queryHandler.InvalidateCache(r.Context(), id)

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "order cancelled successfully"})
}

// AddOrderItem adds an item to an existing order
func (h *AdminOrderHandler) AddOrderItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid order ID").WithInternal(err))
		return
	}

	var req application.AddItemToOrderCommand // Use correct command struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid request body").WithInternal(err))
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, errors.ValidationError("validation failed").WithInternal(err))
		return
	}

	item, err := h.commandHandler.HandleAddItemToOrder(
		r.Context(),
		id,
		&req, // Pass the command struct
	)
	if err != nil {
		httpPkg.RespondError(w, errors.Internal("failed to add item to order").WithInternal(err))
		return
	}

	// Invalidate cache
	h.queryHandler.InvalidateCache(r.Context(), id)

	httpPkg.RespondJSON(w, http.StatusOK, item)
}
