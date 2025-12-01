package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/fulfillment/domain"
	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
)

// ShipmentCommandHandler handles shipment commands
type ShipmentCommandHandler struct {
	repo     domain.ShipmentRepository
	eventBus event.Bus
	log      *logger.Logger
}

// NewShipmentCommandHandler creates a new ShipmentCommandHandler
func NewShipmentCommandHandler(repo domain.ShipmentRepository, eventBus event.Bus, log *logger.Logger) *ShipmentCommandHandler {
	return &ShipmentCommandHandler{
		repo:     repo,
		eventBus: eventBus,
		log:      log,
	}
}

// CreateShipment creates a new shipment
func (h *ShipmentCommandHandler) CreateShipment(ctx context.Context, orderID int64, carrier, shippingMethod string, shippingCost float64, address domain.Address) (*domain.Shipment, error) {
	h.log.WithFields(map[string]interface{}{
		"orderID": orderID,
		"carrier": carrier,
	}).Info("Creating new shipment")

	// Create shipment
	shipment := domain.NewShipment(orderID, carrier, shippingMethod, shippingCost, address)

	// Save shipment
	if err := h.repo.Create(ctx, shipment); err != nil {
		h.log.WithError(err).Error("Failed to create shipment")
		return nil, errors.InternalWrap(err, "failed to create shipment")
	}

	// Publish event
	evt := domain.NewShipmentCreatedEvent(shipment.ID, shipment.OrderID, shipment.Carrier, shipment.ShippingMethod, shipment.ShippingCost)
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.WithError(err).Error("Failed to publish shipment created event")
	}

	h.log.WithField("shipmentID", shipment.ID).Info("Shipment created successfully")
	return shipment, nil
}

// ShipShipment marks a shipment as shipped
func (h *ShipmentCommandHandler) ShipShipment(ctx context.Context, shipmentID int64, trackingNumber string) error {
	h.log.WithField("shipmentID", shipmentID).Info("Marking shipment as shipped")

	// Find shipment
	shipment, err := h.repo.FindByID(ctx, shipmentID)
	if err != nil {
		h.log.WithError(err).Error("Failed to find shipment")
		return err
	}
	if shipment == nil {
		return errors.NotFound(fmt.Sprintf("shipment %d", shipmentID))
	}

	// Ship shipment
	shipment.Ship(trackingNumber)

	// Save shipment
	if err := h.repo.Update(ctx, shipment); err != nil {
		h.log.WithError(err).Error("Failed to ship shipment")
		return errors.InternalWrap(err, "failed to ship shipment")
	}

	// Publish event
	evt := &domain.ShipmentShippedEvent{
		BaseEvent:      event.BaseEvent{Type: domain.EventShipmentShipped, OccurredOn: time.Now()},
		ShipmentID:     shipment.ID,
		OrderID:        shipment.OrderID,
		TrackingNumber: trackingNumber,
		Carrier:        shipment.Carrier,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.WithError(err).Error("Failed to publish shipment shipped event")
	}

	h.log.WithField("shipmentID", shipmentID).Info("Shipment marked as shipped")
	return nil
}

// DeliverShipment marks a shipment as delivered
func (h *ShipmentCommandHandler) DeliverShipment(ctx context.Context, shipmentID int64) error {
	h.log.WithField("shipmentID", shipmentID).Info("Marking shipment as delivered")

	// Find shipment
	shipment, err := h.repo.FindByID(ctx, shipmentID)
	if err != nil {
		h.log.WithError(err).Error("Failed to find shipment")
		return err
	}
	if shipment == nil {
		return errors.NotFound(fmt.Sprintf("shipment %d", shipmentID))
	}

	// Deliver shipment
	shipment.Deliver()

	// Save shipment
	if err := h.repo.Update(ctx, shipment); err != nil {
		h.log.WithError(err).Error("Failed to deliver shipment")
		return errors.InternalWrap(err, "failed to deliver shipment")
	}

	// Publish event
	evt := &domain.ShipmentDeliveredEvent{
		BaseEvent:      event.BaseEvent{Type: domain.EventShipmentDelivered, OccurredOn: time.Now()},
		ShipmentID:     shipment.ID,
		OrderID:        shipment.OrderID,
		TrackingNumber: shipment.TrackingNumber,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.WithError(err).Error("Failed to publish shipment delivered event")
	}

	h.log.WithField("shipmentID", shipmentID).Info("Shipment marked as delivered")
	return nil
}

// CancelShipment cancels a shipment
func (h *ShipmentCommandHandler) CancelShipment(ctx context.Context, shipmentID int64) error {
	h.log.WithField("shipmentID", shipmentID).Info("Cancelling shipment")

	// Find shipment
	shipment, err := h.repo.FindByID(ctx, shipmentID)
	if err != nil {
		h.log.WithError(err).Error("Failed to find shipment")
		return err
	}
	if shipment == nil {
		return errors.NotFound(fmt.Sprintf("shipment %d", shipmentID))
	}

	// Cancel shipment
	if err := shipment.Cancel(); err != nil {
		return errors.ValidationError(err.Error())
	}

	// Save shipment
	if err := h.repo.Update(ctx, shipment); err != nil {
		h.log.WithError(err).Error("Failed to cancel shipment")
		return errors.InternalWrap(err, "failed to cancel shipment")
	}

	// Publish event
	evt := &domain.ShipmentCancelledEvent{
		BaseEvent:  event.BaseEvent{Type: domain.EventShipmentCancelled, OccurredOn: time.Now()},
		ShipmentID: shipment.ID,
		OrderID:    shipment.OrderID,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.WithError(err).Error("Failed to publish shipment cancelled event")
	}

	h.log.WithField("shipmentID", shipmentID).Info("Shipment cancelled successfully")
	return nil
}

// UpdateTracking updates shipment tracking information
func (h *ShipmentCommandHandler) UpdateTracking(ctx context.Context, shipmentID int64, trackingNumber, notes string) error {
	h.log.WithField("shipmentID", shipmentID).Info("Updating shipment tracking")

	// Find shipment
	shipment, err := h.repo.FindByID(ctx, shipmentID)
	if err != nil {
		h.log.WithError(err).Error("Failed to find shipment")
		return err
	}
	if shipment == nil {
		return errors.NotFound(fmt.Sprintf("shipment %d", shipmentID))
	}

	// Update tracking
	shipment.UpdateTracking(trackingNumber, notes)

	// Save shipment
	if err := h.repo.Update(ctx, shipment); err != nil {
		h.log.WithError(err).Error("Failed to update tracking")
		return errors.InternalWrap(err, "failed to update tracking")
	}

	h.log.WithField("shipmentID", shipmentID).Info("Tracking updated successfully")
	return nil
}
