package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/order/domain"
	"github.com/qhato/ecommerce/pkg/cache"
	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/logger"
)

// OrderQueryHandler handles order queries
type OrderQueryHandler struct {
	repo  domain.OrderRepository
	cache cache.Cache
	log   *logger.Logger
}

// NewOrderQueryHandler creates a new OrderQueryHandler
func NewOrderQueryHandler(repo domain.OrderRepository, cache cache.Cache, log *logger.Logger) *OrderQueryHandler {
	return &OrderQueryHandler{
		repo:  repo,
		cache: cache,
		log:   log,
	}
}

// GetByID retrieves an order by ID
func (h *OrderQueryHandler) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	h.log.Debug("Fetching order by ID", "id", id)

	// Try cache first
	cacheKey := fmt.Sprintf("order:id:%d", id)
	if cached, err := h.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var order domain.Order
		if err := json.Unmarshal([]byte(cached), &order); err == nil {
			h.log.Debug("Order found in cache", "id", id)
			return &order, nil
		}
	}

	// Fetch from repository
	order, err := h.repo.FindByID(ctx, id)
	if err != nil {
		h.log.Error("Failed to fetch order by ID", "error", err)
		return nil, err
	}
	if order == nil {
		return nil, errors.NotFound(fmt.Sprintf("order %d", id))
	}

	// Cache result
	if data, err := json.Marshal(order); err == nil {
		_ = h.cache.Set(ctx, cacheKey, string(data), 5*time.Minute)
	}

	return order, nil
}

// GetByOrderNumber retrieves an order by order number
func (h *OrderQueryHandler) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	h.log.Debug("Fetching order by order number", "orderNumber", orderNumber)

	// Try cache first
	cacheKey := fmt.Sprintf("order:number:%s", orderNumber)
	if cached, err := h.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var order domain.Order
		if err := json.Unmarshal([]byte(cached), &order); err == nil {
			h.log.Debug("Order found in cache", "orderNumber", orderNumber)
			return &order, nil
		}
	}

	// Fetch from repository
	order, err := h.repo.FindByOrderNumber(ctx, orderNumber)
	if err != nil {
		h.log.Error("Failed to fetch order by order number", "error", err)
		return nil, err
	}
	if order == nil {
		return nil, errors.NotFound(fmt.Sprintf("order with number %s", orderNumber))
	}

	// Cache result
	if data, err := json.Marshal(order); err == nil {
		_ = h.cache.Set(ctx, cacheKey, string(data), 5*time.Minute)
	}

	return order, nil
}

// ListByCustomer retrieves orders for a customer
func (h *OrderQueryHandler) ListByCustomer(ctx context.Context, customerID int64, filter *domain.OrderFilter) ([]*domain.Order, int64, error) {
	h.log.Debug("Fetching orders by customer", "customerID", customerID)

	orders, total, err := h.repo.FindByCustomerID(ctx, customerID, filter)
	if err != nil {
		h.log.Error("Failed to fetch orders by customer", "error", err)
		return nil, 0, err
	}

	return orders, total, nil
}

// List retrieves all orders with optional filtering
func (h *OrderQueryHandler) List(ctx context.Context, filter *domain.OrderFilter) ([]*domain.Order, int64, error) {
	h.log.Debug("Fetching all orders with filter", "filter", filter)

	orders, total, err := h.repo.FindAll(ctx, filter)
	if err != nil {
		h.log.Error("Failed to fetch orders", "error", err)
		return nil, 0, err
	}

	return orders, total, nil
}

// InvalidateCache invalidates the cache for an order
func (h *OrderQueryHandler) InvalidateCache(ctx context.Context, orderID int64, orderNumber string) {
	cacheKey1 := fmt.Sprintf("order:id:%d", orderID)
	cacheKey2 := fmt.Sprintf("order:number:%s", orderNumber)
	_ = h.cache.Delete(ctx, cacheKey1)
	_ = h.cache.Delete(ctx, cacheKey2)
}
