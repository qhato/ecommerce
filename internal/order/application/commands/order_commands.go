package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/order/application" // Import order application package
	"github.com/qhato/ecommerce/internal/order/domain"
	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/validator"
)

// OrderCommandHandler handles order-related commands.
type OrderCommandHandler struct {
	orderService application.OrderService // Dependency on the application service
	eventBus     event.Bus
	logger       *logger.Logger
	validator    *validator.Validator
}

// NewOrderCommandHandler creates a new OrderCommandHandler.
func NewOrderCommandHandler(
	orderService application.OrderService,
	eventBus event.Bus,
	logger *logger.Logger,
	validator *validator.Validator,
) *OrderCommandHandler {
	return &OrderCommandHandler{
		orderService: orderService,
		eventBus:     eventBus,
		logger:       logger,
		validator:    validator,
	}
}

// HandleCreateOrder handles the creation of a new order.
func (h *OrderCommandHandler) HandleCreateOrder(ctx context.Context, cmd *application.CreateOrderCommand) (*application.OrderDTO, error) {
	if err := h.validator.Validate(cmd); err != nil {
		return nil, err // Let the application service handle validation errors
	}

	orderDTO, err := h.orderService.CreateOrder(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Publish an event for order creation if needed
	// h.eventBus.Publish(ctx, domain.NewOrderCreatedEvent(orderDTO.ID))

	return orderDTO, nil
}

// HandleUpdateOrderStatus handles updating the status of an order.
func (h *OrderCommandHandler) HandleUpdateOrderStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error {
	err := h.orderService.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}

// HandleAddItemToOrder handles adding an item to an order.
func (h *OrderCommandHandler) HandleAddItemToOrder(ctx context.Context, orderID int64, cmd *application.AddItemToOrderCommand) (*application.OrderItemDTO, error) {
	if err := h.validator.Validate(cmd); err != nil {
		return nil, err
	}

	itemDTO, err := h.orderService.AddItemToOrder(ctx, orderID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to add item to order: %w", err)
	}
	return itemDTO, nil
}

// HandleUpdateOrderItemQuantity handles updating the quantity of an order item.
func (h *OrderCommandHandler) HandleUpdateOrderItemQuantity(ctx context.Context, orderItemID int64, newQuantity int) (*application.OrderItemDTO, error) {
	itemDTO, err := h.orderService.UpdateOrderItemQuantity(ctx, orderItemID, newQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to update order item quantity: %w", err)
	}
	return itemDTO, nil
}

// HandleRemoveOrderItem handles removing an item from an order.
func (h *OrderCommandHandler) HandleRemoveOrderItem(ctx context.Context, orderItemID int64) error {
	err := h.orderService.RemoveOrderItem(ctx, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to remove order item: %w", err)
	}
	return nil
}

// HandleSubmitOrder handles submitting an order.
func (h *OrderCommandHandler) HandleSubmitOrder(ctx context.Context, orderID int64) error {
	err := h.orderService.SubmitOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to submit order: %w", err)
	}
	return nil
}

// HandleCancelOrder handles canceling an order.
func (h *OrderCommandHandler) HandleCancelOrder(ctx context.Context, orderID int64, reason string) error {
	err := h.orderService.CancelOrder(ctx, orderID, reason)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	return nil
}