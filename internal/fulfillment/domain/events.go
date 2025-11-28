package domain

import (
	"time"

	"github.com/qhato/ecommerce/pkg/event"
)

const (
	EventShipmentCreated   = "shipment.created"
	EventShipmentShipped   = "shipment.shipped"
	EventShipmentDelivered = "shipment.delivered"
	EventShipmentCancelled = "shipment.cancelled"
)

type ShipmentCreatedEvent struct {
	event.BaseEvent
	ShipmentID     int64   `json:"shipment_id"`
	OrderID        int64   `json:"order_id"`
	Carrier        string  `json:"carrier"`
	ShippingMethod string  `json:"shipping_method"`
	ShippingCost   float64 `json:"shipping_cost"`
}

func NewShipmentCreatedEvent(shipmentID, orderID int64, carrier, shippingMethod string, shippingCost float64) *ShipmentCreatedEvent {
	return &ShipmentCreatedEvent{
		BaseEvent:      event.BaseEvent{EventType: EventShipmentCreated, Timestamp: time.Now()},
		ShipmentID:     shipmentID,
		OrderID:        orderID,
		Carrier:        carrier,
		ShippingMethod: shippingMethod,
		ShippingCost:   shippingCost,
	}
}

func (e *ShipmentCreatedEvent) Type() string {
	return e.EventType
}

type ShipmentShippedEvent struct {
	event.BaseEvent
	ShipmentID     int64  `json:"shipment_id"`
	OrderID        int64  `json:"order_id"`
	TrackingNumber string `json:"tracking_number"`
	Carrier        string `json:"carrier"`
}

func (e *ShipmentShippedEvent) Type() string {
	return e.EventType
}

type ShipmentDeliveredEvent struct {
	event.BaseEvent
	ShipmentID     int64  `json:"shipment_id"`
	OrderID        int64  `json:"order_id"`
	TrackingNumber string `json:"tracking_number"`
}

func (e *ShipmentDeliveredEvent) Type() string {
	return e.EventType
}

type ShipmentCancelledEvent struct {
	event.BaseEvent
	ShipmentID int64 `json:"shipment_id"`
	OrderID    int64 `json:"order_id"`
}

func (e *ShipmentCancelledEvent) Type() string {
	return e.EventType
}
