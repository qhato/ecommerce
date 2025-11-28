package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/order/domain"
	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
)

// OrderCommandHandler handles order commands
type OrderCommandHandler struct {
	repo     domain.OrderRepository
	eventBus event.Bus
	log      *logger.Logger
}

// NewOrderCommandHandler creates a new OrderCommandHandler
func NewOrderCommandHandler(repo domain.OrderRepository, eventBus event.Bus, log *logger.Logger) *OrderCommandHandler {
	return &OrderCommandHandler{
		repo:     repo,
		eventBus: eventBus,
		log:      log,
	}
}

// CreateOrder creates a new order
func (h *OrderCommandHandler) CreateOrder(ctx context.Context, customerID int64, emailAddress, name, currencyCode string, items []domain.OrderItem) (*domain.Order, error) {
	h.log.WithFields(map[string]interface{}{
		"customerID": customerID,
		"email":      emailAddress,
	}).Info("Creating new order")

	// Validate items
	if len(items) == 0 {
		return nil, errors.ValidationError("order must have at least one item")
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
		h.log.WithError(err).Error("Failed to create order")
		return nil, errors.InternalWrap(err, "failed to create order")
	}

	// Publish event
	evt := domain.NewOrderCreatedEvent(order.ID, order.OrderNumber, order.CustomerID, order.Total, order.CurrencyCode)
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.WithError(err).Error("Failed to publish order created event")
	}

	h.log.WithFields(map[string]interface{}{
		"orderID":     order.ID,
		"orderNumber": order.OrderNumber,
	}).Info("Order created successfully")
	return order, nil
}

// UpdateOrderStatus updates the status of an order
func (h *OrderCommandHandler) UpdateOrderStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error {
	h.log.WithFields(map[string]interface{}{
		"orderID": orderID,
		"status":  status,
	}).Info("Updating order status")

	// Find order
	order, err := h.repo.FindByID(ctx, orderID)
	if err != nil {
		h.log.WithError(err).Error("Failed to find order")
		return err
	}
	if order == nil {
		return errors.NotFound(fmt.Sprintf("order %d", orderID))
	}

	// Update status
	order.UpdateStatus(status)

	// Save order
	if err := h.repo.Update(ctx, order); err != nil {
		h.log.WithError(err).Error("Failed to update order status")
		return errors.InternalWrap(err, "failed to update order status")
	}

	h.log.WithFields(map[string]interface{}{
		"orderID": orderID,
		"status":  status,
	}).Info("Order status updated successfully")
	return nil
}

// SubmitOrder submits an order for processing
func (h *OrderCommandHandler) SubmitOrder(ctx context.Context, orderID int64) error {
	h.log.WithField("orderID", orderID).Info("Submitting order")

	// Find order
	order, err := h.repo.FindByID(ctx, orderID)
	if err != nil {
		h.log.WithError(err).Error("Failed to find order")
		return err
	}
	if order == nil {
		return errors.NotFound(fmt.Sprintf("order %d", orderID))
	}

	// Validate order can be submitted
	if order.Status != domain.OrderStatusPending {
		return errors.ValidationError("only pending orders can be submitted")
	}

	// Submit order
	order.Submit()

	// Save order
	if err := h.repo.Update(ctx, order); err != nil {
		h.log.WithError(err).Error("Failed to submit order")
		return errors.InternalWrap(err, "failed to submit order")
	}

	// Publish event
	evt := &domain.OrderSubmittedEvent{
		BaseEvent:   event.BaseEvent{Type: domain.EventOrderSubmitted, OccurredOn: time.Now()},
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		CustomerID:  order.CustomerID,
		Total:       order.Total,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.WithError(err).Error("Failed to publish order submitted event")
	}

	h.log.WithField("orderID", orderID).Info("Order submitted successfully")
	return nil
}

// CancelOrder cancels an order
func (h *OrderCommandHandler) CancelOrder(ctx context.Context, orderID int64) error {
	h.log.WithField("orderID", orderID).Info("Cancelling order")

	// Find order
	order, err := h.repo.FindByID(ctx, orderID)
	if err != nil {
		h.log.WithError(err).Error("Failed to find order")
		return err
	}
	if order == nil {
		return errors.NotFound(fmt.Sprintf("order %d", orderID))
	}

	// Validate order can be cancelled
	if !order.IsCancellable() {
		return errors.ValidationError("order cannot be cancelled in current status")
	}

	// Cancel order
	order.Cancel()

	// Save order
	if err := h.repo.Update(ctx, order); err != nil {
		h.log.WithError(err).Error("Failed to cancel order")
		return errors.InternalWrap(err, "failed to cancel order")
	}

	// Publish event
	evt := &domain.OrderCancelledEvent{
		BaseEvent:   event.BaseEvent{Type: domain.EventOrderCancelled, OccurredOn: time.Now()},
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		CustomerID:  order.CustomerID,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.WithError(err).Error("Failed to publish order cancelled event")
	}

	h.log.WithField("orderID", orderID).Info("Order cancelled successfully")
	return nil
}

// AddOrderItem adds an item to an existing order
func (h *OrderCommandHandler) AddOrderItem(ctx context.Context, orderID, skuID int64, productName string, quantity int, price float64) error {
	h.log.WithFields(map[string]interface{}{
		"orderID": orderID,
		"skuID":   skuID,
	}).Info("Adding item to order")

	// Find order
	order, err := h.repo.FindByID(ctx, orderID)
	if err != nil {
		h.log.WithError(err).Error("Failed to find order")
		return err
	}
	if order == nil {
		return errors.NotFound(fmt.Sprintf("order %d", orderID))
	}

	// Validate order is in editable state
	if order.Status != domain.OrderStatusPending {
		return errors.ValidationError("items can only be added to pending orders")
	}

	// Add item
	order.AddItem(skuID, productName, quantity, price)

	// Save order
	if err := h.repo.Update(ctx, order); err != nil {
		h.log.WithError(err).Error("Failed to add item to order")
		return errors.InternalWrap(err, "failed to add item to order")
	}

	h.log.WithFields(map[string]interface{}{
		"orderID": orderID,
		"skuID":   skuID,
	}).Info("Item added to order successfully")
	return nil
}

// generateOrderNumber generates a unique order number
// Simple implementation - in production, use a more robust approach
func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano()/1000000)
}
