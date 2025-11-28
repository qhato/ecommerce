package commands

import (
	"context"
	"time"

	"github.com/qhato/ecommerce/internal/fulfillment/domain"
	"github.com/qhato/ecommerce/pkg/apperrors"
	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
)

// ShipmentCommandHandler handles shipment commands
type ShipmentCommandHandler struct {
	repo     domain.ShipmentRepository
	eventBus event.EventBus
	log      *logger.Logger
}

// NewShipmentCommandHandler creates a new ShipmentCommandHandler
func NewShipmentCommandHandler(repo domain.ShipmentRepository, eventBus event.EventBus, log *logger.Logger) *ShipmentCommandHandler {
	return &ShipmentCommandHandler{
		repo:     repo,
		eventBus: eventBus,
		log:      log,
	}
}

// CreateShipment creates a new shipment
func (h *ShipmentCommandHandler) CreateShipment(ctx context.Context, orderID int64, carrier, shippingMethod string, shippingCost float64, address domain.Address) (*domain.Shipment, error) {
	h.log.Info("Creating new shipment", "orderID", orderID, "carrier", carrier)

	// Create shipment
	shipment := domain.NewShipment(orderID, carrier, shippingMethod, shippingCost, address)

	// Save shipment
	if err := h.repo.Create(ctx, shipment); err != nil {
		h.log.Error("Failed to create shipment", "error", err)
		return nil, apperrors.NewInternalError("failed to create shipment", err)
	}

	// Publish event
	evt := domain.NewShipmentCreatedEvent(shipment.ID, shipment.OrderID, shipment.Carrier, shipment.ShippingMethod, shipment.ShippingCost)
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish shipment created event", "error", err)
	}

	h.log.Info("Shipment created successfully", "shipmentID", shipment.ID)
	return shipment, nil
}

// ShipShipment marks a shipment as shipped
func (h *ShipmentCommandHandler) ShipShipment(ctx context.Context, shipmentID int64, trackingNumber string) error {
	h.log.Info("Marking shipment as shipped", "shipmentID", shipmentID)

	// Find shipment
	shipment, err := h.repo.FindByID(ctx, shipmentID)
	if err != nil {
		h.log.Error("Failed to find shipment", "error", err)
		return err
	}
	if shipment == nil {
		return apperrors.NewNotFoundError("shipment", shipmentID)
	}

	// Ship shipment
	shipment.Ship(trackingNumber)

	// Save shipment
	if err := h.repo.Update(ctx, shipment); err != nil {
		h.log.Error("Failed to ship shipment", "error", err)
		return apperrors.NewInternalError("failed to ship shipment", err)
	}

	// Publish event
	evt := &domain.ShipmentShippedEvent{
		BaseEvent:      event.BaseEvent{EventType: domain.EventShipmentShipped, Timestamp: time.Now()},
		ShipmentID:     shipment.ID,
		OrderID:        shipment.OrderID,
		TrackingNumber: trackingNumber,
		Carrier:        shipment.Carrier,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish shipment shipped event", "error", err)
	}

	h.log.Info("Shipment marked as shipped", "shipmentID", shipmentID)
	return nil
}

// DeliverShipment marks a shipment as delivered
func (h *ShipmentCommandHandler) DeliverShipment(ctx context.Context, shipmentID int64) error {
	h.log.Info("Marking shipment as delivered", "shipmentID", shipmentID)

	// Find shipment
	shipment, err := h.repo.FindByID(ctx, shipmentID)
	if err != nil {
		h.log.Error("Failed to find shipment", "error", err)
		return err
	}
	if shipment == nil {
		return apperrors.NewNotFoundError("shipment", shipmentID)
	}

	// Deliver shipment
	shipment.Deliver()

	// Save shipment
	if err := h.repo.Update(ctx, shipment); err != nil {
		h.log.Error("Failed to deliver shipment", "error", err)
		return apperrors.NewInternalError("failed to deliver shipment", err)
	}

	// Publish event
	evt := &domain.ShipmentDeliveredEvent{
		BaseEvent:      event.BaseEvent{EventType: domain.EventShipmentDelivered, Timestamp: time.Now()},
		ShipmentID:     shipment.ID,
		OrderID:        shipment.OrderID,
		TrackingNumber: shipment.TrackingNumber,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish shipment delivered event", "error", err)
	}

	h.log.Info("Shipment marked as delivered", "shipmentID", shipmentID)
	return nil
}

// CancelShipment cancels a shipment
func (h *ShipmentCommandHandler) CancelShipment(ctx context.Context, shipmentID int64) error {
	h.log.Info("Cancelling shipment", "shipmentID", shipmentID)

	// Find shipment
	shipment, err := h.repo.FindByID(ctx, shipmentID)
	if err != nil {
		h.log.Error("Failed to find shipment", "error", err)
		return err
	}
	if shipment == nil {
		return apperrors.NewNotFoundError("shipment", shipmentID)
	}

	// Cancel shipment
	if err := shipment.Cancel(); err != nil {
		return apperrors.NewValidationError(err.Error())
	}

	// Save shipment
	if err := h.repo.Update(ctx, shipment); err != nil {
		h.log.Error("Failed to cancel shipment", "error", err)
		return apperrors.NewInternalError("failed to cancel shipment", err)
	}

	// Publish event
	evt := &domain.ShipmentCancelledEvent{
		BaseEvent:  event.BaseEvent{EventType: domain.EventShipmentCancelled, Timestamp: time.Now()},
		ShipmentID: shipment.ID,
		OrderID:    shipment.OrderID,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish shipment cancelled event", "error", err)
	}

	h.log.Info("Shipment cancelled successfully", "shipmentID", shipmentID)
	return nil
}

// UpdateTracking updates shipment tracking information
func (h *ShipmentCommandHandler) UpdateTracking(ctx context.Context, shipmentID int64, trackingNumber, notes string) error {
	h.log.Info("Updating shipment tracking", "shipmentID", shipmentID)

	// Find shipment
	shipment, err := h.repo.FindByID(ctx, shipmentID)
	if err != nil {
		h.log.Error("Failed to find shipment", "error", err)
		return err
	}
	if shipment == nil {
		return apperrors.NewNotFoundError("shipment", shipmentID)
	}

	// Update tracking
	shipment.UpdateTracking(trackingNumber, notes)

	// Save shipment
	if err := h.repo.Update(ctx, shipment); err != nil {
		h.log.Error("Failed to update tracking", "error", err)
		return apperrors.NewInternalError("failed to update tracking", err)
	}

	h.log.Info("Tracking updated successfully", "shipmentID", shipmentID)
	return nil
}
