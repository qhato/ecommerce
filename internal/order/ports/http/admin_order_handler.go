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
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	// Convert request items to domain items
	items := make([]domain.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.OrderItem{
			SKUID:       item.SKUID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
		}
	}

	order, err := h.commandHandler.CreateOrder(
		r.Context(),
		req.CustomerID,
		req.EmailAddress,
		req.Name,
		req.CurrencyCode,
		items,
	)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to create order", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusCreated, application.ToOrderDTO(order))
}

// GetOrder retrieves an order by ID
func (h *AdminOrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
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
func (h *AdminOrderHandler) GetOrderByNumber(w http.ResponseWriter, r *http.Request) {
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

	status := r.URL.Query().Get("status")
	customerIDStr := r.URL.Query().Get("customer_id")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	var customerID int64
	if customerIDStr != "" {
		customerID, _ = strconv.ParseInt(customerIDStr, 10, 64)
	}

	filter := &domain.OrderFilter{
		Page:       page,
		PageSize:   pageSize,
		Status:     domain.OrderStatus(status),
		CustomerID: customerID,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	orders, total, err := h.queryHandler.List(r.Context(), filter)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to list orders", err)
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

// UpdateOrderStatus updates the status of an order
func (h *AdminOrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid order ID", err)
		return
	}

	var req application.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	err = h.commandHandler.UpdateOrderStatus(r.Context(), id, domain.OrderStatus(req.Status))
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to update order status", err)
		return
	}

	// Invalidate cache
	order, _ := h.queryHandler.GetByID(r.Context(), id)
	if order != nil {
		h.queryHandler.InvalidateCache(r.Context(), order.ID, order.OrderNumber)
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "order status updated successfully"})
}

// SubmitOrder submits an order for processing
func (h *AdminOrderHandler) SubmitOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid order ID", err)
		return
	}

	err = h.commandHandler.SubmitOrder(r.Context(), id)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to submit order", err)
		return
	}

	// Invalidate cache
	order, _ := h.queryHandler.GetByID(r.Context(), id)
	if order != nil {
		h.queryHandler.InvalidateCache(r.Context(), order.ID, order.OrderNumber)
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "order submitted successfully"})
}

// CancelOrder cancels an order
func (h *AdminOrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid order ID", err)
		return
	}

	err = h.commandHandler.CancelOrder(r.Context(), id)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to cancel order", err)
		return
	}

	// Invalidate cache
	order, _ := h.queryHandler.GetByID(r.Context(), id)
	if order != nil {
		h.queryHandler.InvalidateCache(r.Context(), order.ID, order.OrderNumber)
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "order cancelled successfully"})
}

// AddOrderItem adds an item to an existing order
func (h *AdminOrderHandler) AddOrderItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid order ID", err)
		return
	}

	var req application.AddOrderItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	err = h.commandHandler.AddOrderItem(
		r.Context(),
		id,
		req.SKUID,
		req.ProductName,
		req.Quantity,
		req.Price,
	)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to add item to order", err)
		return
	}

	// Invalidate cache
	order, _ := h.queryHandler.GetByID(r.Context(), id)
	if order != nil {
		h.queryHandler.InvalidateCache(r.Context(), order.ID, order.OrderNumber)
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "item added to order successfully"})
}
