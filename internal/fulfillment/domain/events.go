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
		BaseEvent:      event.BaseEvent{Type: EventShipmentCreated, OccurredOn: time.Now()},
		ShipmentID:     shipmentID,
		OrderID:        orderID,
		Carrier:        carrier,
		ShippingMethod: shippingMethod,
		ShippingCost:   shippingCost,
	}
}



type ShipmentShippedEvent struct {
	event.BaseEvent
	ShipmentID     int64  `json:"shipment_id"`
	OrderID        int64  `json:"order_id"`
	TrackingNumber string `json:"tracking_number"`
	Carrier        string `json:"carrier"`
}



type ShipmentDeliveredEvent struct {
	event.BaseEvent
	ShipmentID     int64  `json:"shipment_id"`
	OrderID        int64  `json:"order_id"`
	TrackingNumber string `json:"tracking_number"`
}



type ShipmentCancelledEvent struct {
	event.BaseEvent
	ShipmentID int64 `json:"shipment_id"`
	OrderID    int64 `json:"order_id"`
}


