package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/order/domain"
	"github.com/qhato/ecommerce/pkg/apperrors"
	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
)

// OrderCommandHandler handles order commands
type OrderCommandHandler struct {
	repo     domain.OrderRepository
	eventBus event.EventBus
	log      *logger.Logger
}

// NewOrderCommandHandler creates a new OrderCommandHandler
func NewOrderCommandHandler(repo domain.OrderRepository, eventBus event.EventBus, log *logger.Logger) *OrderCommandHandler {
	return &OrderCommandHandler{
		repo:     repo,
		eventBus: eventBus,
		log:      log,
	}
}

// CreateOrder creates a new order
func (h *OrderCommandHandler) CreateOrder(ctx context.Context, customerID int64, emailAddress, name, currencyCode string, items []domain.OrderItem) (*domain.Order, error) {
	h.log.Info("Creating new order", "customerID", customerID, "email", emailAddress)

	// Validate items
	if len(items) == 0 {
		return nil, apperrors.NewValidationError("order must have at least one item")
	}

	// Create order
	order := domain.NewOrder(customerID, emailAddress, name, currencyCode)

	// Add items
	for _, item := range items {
		order.AddItem(item.SKUID, item.ProductName, item.Quantity, item.Price)
	}

	// Generate order number (simple implementation - can be improved)
	order.OrderNumber = generateOrderNumber()

	// Save order
	if err := h.repo.Create(ctx, order); err != nil {
		h.log.Error("Failed to create order", "error", err)
		return nil, apperrors.NewInternalError("failed to create order", err)
	}

	// Publish event
	evt := domain.NewOrderCreatedEvent(order.ID, order.OrderNumber, order.CustomerID, order.Total, order.CurrencyCode)
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish order created event", "error", err)
	}

	h.log.Info("Order created successfully", "orderID", order.ID, "orderNumber", order.OrderNumber)
	return order, nil
}

// UpdateOrderStatus updates the status of an order
func (h *OrderCommandHandler) UpdateOrderStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error {
	h.log.Info("Updating order status", "orderID", orderID, "status", status)

	// Find order
	order, err := h.repo.FindByID(ctx, orderID)
	if err != nil {
		h.log.Error("Failed to find order", "error", err)
		return err
	}
	if order == nil {
		return apperrors.NewNotFoundError("order", orderID)
	}

	// Update status
	order.UpdateStatus(status)

	// Save order
	if err := h.repo.Update(ctx, order); err != nil {
		h.log.Error("Failed to update order status", "error", err)
		return apperrors.NewInternalError("failed to update order status", err)
	}

	h.log.Info("Order status updated successfully", "orderID", orderID, "status", status)
	return nil
}

// SubmitOrder submits an order for processing
func (h *OrderCommandHandler) SubmitOrder(ctx context.Context, orderID int64) error {
	h.log.Info("Submitting order", "orderID", orderID)

	// Find order
	order, err := h.repo.FindByID(ctx, orderID)
	if err != nil {
		h.log.Error("Failed to find order", "error", err)
		return err
	}
	if order == nil {
		return apperrors.NewNotFoundError("order", orderID)
	}

	// Validate order can be submitted
	if order.Status != domain.OrderStatusPending {
		return apperrors.NewValidationError("only pending orders can be submitted")
	}

	// Submit order
	order.Submit()

	// Save order
	if err := h.repo.Update(ctx, order); err != nil {
		h.log.Error("Failed to submit order", "error", err)
		return apperrors.NewInternalError("failed to submit order", err)
	}

	// Publish event
	evt := &domain.OrderSubmittedEvent{
		BaseEvent:   event.BaseEvent{EventType: domain.EventOrderSubmitted, Timestamp: time.Now()},
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		CustomerID:  order.CustomerID,
		Total:       order.Total,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish order submitted event", "error", err)
	}

	h.log.Info("Order submitted successfully", "orderID", orderID)
	return nil
}

// CancelOrder cancels an order
func (h *OrderCommandHandler) CancelOrder(ctx context.Context, orderID int64) error {
	h.log.Info("Cancelling order", "orderID", orderID)

	// Find order
	order, err := h.repo.FindByID(ctx, orderID)
	if err != nil {
		h.log.Error("Failed to find order", "error", err)
		return err
	}
	if order == nil {
		return apperrors.NewNotFoundError("order", orderID)
	}

	// Validate order can be cancelled
	if !order.IsCancellable() {
		return apperrors.NewValidationError("order cannot be cancelled in current status")
	}

	// Cancel order
	order.Cancel()

	// Save order
	if err := h.repo.Update(ctx, order); err != nil {
		h.log.Error("Failed to cancel order", "error", err)
		return apperrors.NewInternalError("failed to cancel order", err)
	}

	// Publish event
	evt := &domain.OrderCancelledEvent{
		BaseEvent:   event.BaseEvent{EventType: domain.EventOrderCancelled, Timestamp: time.Now()},
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		CustomerID:  order.CustomerID,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish order cancelled event", "error", err)
	}

	h.log.Info("Order cancelled successfully", "orderID", orderID)
	return nil
}

// AddOrderItem adds an item to an existing order
func (h *OrderCommandHandler) AddOrderItem(ctx context.Context, orderID, skuID int64, productName string, quantity int, price float64) error {
	h.log.Info("Adding item to order", "orderID", orderID, "skuID", skuID)

	// Find order
	order, err := h.repo.FindByID(ctx, orderID)
	if err != nil {
		h.log.Error("Failed to find order", "error", err)
		return err
	}
	if order == nil {
		return apperrors.NewNotFoundError("order", orderID)
	}

	// Validate order is in editable state
	if order.Status != domain.OrderStatusPending {
		return apperrors.NewValidationError("items can only be added to pending orders")
	}

	// Add item
	order.AddItem(skuID, productName, quantity, price)

	// Save order
	if err := h.repo.Update(ctx, order); err != nil {
		h.log.Error("Failed to add item to order", "error", err)
		return apperrors.NewInternalError("failed to add item to order", err)
	}

	h.log.Info("Item added to order successfully", "orderID", orderID, "skuID", skuID)
	return nil
}

// generateOrderNumber generates a unique order number
// Simple implementation - in production, use a more robust approach
func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano()/1000000)
}
