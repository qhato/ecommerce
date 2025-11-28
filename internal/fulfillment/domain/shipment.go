package domain

import (
	"time"
)

// ShipmentStatus represents the status of a shipment
type ShipmentStatus string

const (
	ShipmentStatusPending    ShipmentStatus = "PENDING"
	ShipmentStatusProcessing ShipmentStatus = "PROCESSING"
	ShipmentStatusShipped    ShipmentStatus = "SHIPPED"
	ShipmentStatusInTransit  ShipmentStatus = "IN_TRANSIT"
	ShipmentStatusDelivered  ShipmentStatus = "DELIVERED"
	ShipmentStatusFailed     ShipmentStatus = "FAILED"
	ShipmentStatusCancelled  ShipmentStatus = "CANCELLED"
)

// Shipment represents a shipment entity
type Shipment struct {
	ID              int64
	OrderID         int64
	Status          ShipmentStatus
	TrackingNumber  string
	Carrier         string
	ShippingMethod  string
	ShippingCost    float64
	EstimatedDate   *time.Time
	ShippedDate     *time.Time
	DeliveredDate   *time.Time
	ShippingAddress Address
	Notes           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Address represents a shipping address
type Address struct {
	Name       string
	Line1      string
	Line2      string
	City       string
	State      string
	PostalCode string
	Country    string
	Phone      string
}

// NewShipment creates a new shipment
func NewShipment(orderID int64, carrier, shippingMethod string, shippingCost float64, address Address) *Shipment {
	now := time.Now()
	return &Shipment{
		OrderID:         orderID,
		Status:          ShipmentStatusPending,
		Carrier:         carrier,
		ShippingMethod:  shippingMethod,
		ShippingCost:    shippingCost,
		ShippingAddress: address,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Ship marks the shipment as shipped
func (s *Shipment) Ship(trackingNumber string) {
	now := time.Now()
	s.Status = ShipmentStatusShipped
	s.TrackingNumber = trackingNumber
	s.ShippedDate = &now
	s.UpdatedAt = now
}

// UpdateStatus updates the shipment status
func (s *Shipment) UpdateStatus(status ShipmentStatus) {
	s.Status = status
	s.UpdatedAt = time.Now()
}

// Deliver marks the shipment as delivered
func (s *Shipment) Deliver() {
	now := time.Now()
	s.Status = ShipmentStatusDelivered
	s.DeliveredDate = &now
	s.UpdatedAt = now
}

// Cancel cancels the shipment
func (s *Shipment) Cancel() error {
	if s.Status == ShipmentStatusDelivered {
		return NewFulfillmentError("cannot cancel delivered shipment")
	}
	s.Status = ShipmentStatusCancelled
	s.UpdatedAt = time.Now()
	return nil
}

// UpdateTracking updates tracking information
func (s *Shipment) UpdateTracking(trackingNumber, notes string) {
	s.TrackingNumber = trackingNumber
	s.Notes = notes
	s.UpdatedAt = time.Now()
}

// FulfillmentError represents a fulfillment domain error
type FulfillmentError struct {
	Message string
}

func (e *FulfillmentError) Error() string {
	return e.Message
}

// NewFulfillmentError creates a new fulfillment error
func NewFulfillmentError(message string) *FulfillmentError {
	return &FulfillmentError{Message: message}
}
