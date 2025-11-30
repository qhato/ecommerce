package queries

import (
	"context"
	"encoding/json" // Added json import
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/order/application" // Import order application package
	"github.com/qhato/ecommerce/internal/order/domain"
	"github.com/qhato/ecommerce/pkg/cache"
	"github.com/qhato/ecommerce/pkg/logger"
)

// GetOrderByIDQuery represents a query to get an order by ID.
type GetOrderByIDQuery struct {
	ID int64 `json:"id" validate:"required"`
}

// ListOrdersQuery represents a query to list orders with filters.
type ListOrdersQuery struct {
	Page       int               `json:"page" validate:"min=1"`
	PageSize   int               `json:"page_size" validate:"min=1,max=100"`
	CustomerID *int64            `json:"customer_id,omitempty"`
	Status     *domain.OrderStatus `json:"status,omitempty"`
	SortBy     string            `json:"sort_by"`
	SortOrder  string            `json:"sort_order"`
}

// GetOrderByOrderNumberQuery represents a query to get an order by order number.
type GetOrderByOrderNumberQuery struct {
	OrderNumber string `json:"order_number" validate:"required"`
}

// OrderQueryHandler handles order-related queries.
type OrderQueryHandler struct {
	orderService application.OrderService // Dependency on the application service
	cache        cache.Cache
	logger       *logger.Logger
}

// NewOrderQueryHandler creates a new OrderQueryHandler.
func NewOrderQueryHandler(
	orderService application.OrderService,
	cache cache.Cache,
	logger *logger.Logger,
) *OrderQueryHandler {
	return &OrderQueryHandler{
		orderService: orderService,
		cache:        cache,
		logger:       logger,
	}
}

// HandleGetOrderByID handles the GetOrderByIDQuery.
func (h *OrderQueryHandler) HandleGetOrderByID(ctx context.Context, query *GetOrderByIDQuery) (*application.OrderDTO, error) {
	// Try to get from cache first
	cacheKey := orderCacheKey(query.ID)
	cached, err := h.cache.Get(ctx, cacheKey)
	if err == nil && cached != nil {
		var orderDTO application.OrderDTO
		if err := json.Unmarshal(cached, &orderDTO); err == nil {
			h.logger.WithField("order_id", query.ID).Debug("order found in cache")
			return &orderDTO, nil
		}
	}

	orderDTO, err := h.orderService.HandleGetOrderByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	serialized, err := json.Marshal(orderDTO)
	if err == nil {
		if err := h.cache.Set(ctx, cacheKey, serialized, 5*time.Minute); err != nil {
			h.logger.WithError(err).WithField("order_id", query.ID).Warn("failed to cache order")
		}
	} else {
		h.logger.WithError(err).WithField("order_id", query.ID).Warn("failed to serialize order for caching")
	}

	return orderDTO, nil
}

// HandleGetOrderByOrderNumber handles the GetOrderByOrderNumberQuery.
func (h *OrderQueryHandler) HandleGetOrderByOrderNumber(ctx context.Context, query *GetOrderByOrderNumberQuery) (*application.OrderDTO, error) {
	// Try to get from cache first (using order number as key)
	cacheKey := orderCacheKeyByNumber(query.OrderNumber)
	cached, err := h.cache.Get(ctx, cacheKey)
	if err == nil && cached != nil {
		var orderDTO application.OrderDTO
		if err := json.Unmarshal(cached, &orderDTO); err == nil {
			h.logger.WithField("order_number", query.OrderNumber).Debug("order found in cache by number")
			return &orderDTO, nil
		}
	}

	// Get from repository
	orderDTO, err := h.orderService.GetOrderByOrderNumber(ctx, query.OrderNumber) // Assuming GetOrderByOrderNumber method exists in OrderService
	if err != nil {
		return nil, err
	}

	// Cache the result
	serialized, err := json.Marshal(orderDTO)
	if err == nil {
		if err := h.cache.Set(ctx, cacheKey, serialized, 5*time.Minute); err != nil {
			h.logger.WithError(err).WithField("order_number", query.OrderNumber).Warn("failed to cache order by number")
		}
	} else {
		h.logger.WithError(err).WithField("order_number", query.OrderNumber).Warn("failed to serialize order for caching by number")
	}

	return orderDTO, nil
}

// HandleListOrders handles the ListOrdersQuery.
func (h *OrderQueryHandler) HandleListOrders(ctx context.Context, query *ListOrdersQuery) (*application.PaginatedResponse, error) {
	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	if query.SortBy == "" {
		query.SortBy = "created_at"
	}
	if query.SortOrder == "" {
		query.SortOrder = "desc"
	}

	filter := &domain.OrderFilter{
		Page:       query.Page,
		PageSize:   query.PageSize,
		CustomerID: query.CustomerID,
		Status:     query.Status,
		SortBy:     query.SortBy,
		SortOrder:  query.SortOrder,
	}

	orders, total, err := h.orderService.ListOrders(ctx, filter) // Assuming ListOrders method exists in OrderService
	if err != nil {
		return nil, err
	}

	orderDTOs := make([]application.OrderDTO, len(orders))
	for i, order := range orders {
		orderDTOs[i] = *application.ToOrderDTO(order)
	}

	return application.NewPaginatedResponse(orderDTOs, query.Page, query.PageSize, total), nil
}

// InvalidateCache invalidates the cache for a specific order ID.
func (h *OrderQueryHandler) InvalidateCache(ctx context.Context, orderID int64) {
	cacheKey := orderCacheKey(orderID)
	if err := h.cache.Delete(ctx, cacheKey); err != nil {
		h.logger.WithError(err).WithField("order_id", orderID).Warn("failed to invalidate order cache")
	}
}

// orderCacheKey generates a cache key for an order.
func orderCacheKey(id int64) string {
	return fmt.Sprintf("order:%d", id)
}

// orderCacheKeyByNumber generates a cache key for an order by its order number.
func orderCacheKeyByNumber(orderNumber string) string {
	return fmt.Sprintf("order:number:%s", orderNumber)
}