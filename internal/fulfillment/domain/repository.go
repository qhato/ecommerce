package domain

import (
	"context"
)

// ShipmentRepository defines the interface for shipment persistence
type ShipmentRepository interface {
	Create(ctx context.Context, shipment *Shipment) error
	Update(ctx context.Context, shipment *Shipment) error
	FindByID(ctx context.Context, id int64) (*Shipment, error)
	FindByOrderID(ctx context.Context, orderID int64) ([]*Shipment, error)
	FindByTrackingNumber(ctx context.Context, trackingNumber string) (*Shipment, error)
	FindAll(ctx context.Context, filter *ShipmentFilter) ([]*Shipment, int64, error)
}

// ShipmentFilter represents filtering options for shipments
type ShipmentFilter struct {
	Page       int
	PageSize   int
	Status     ShipmentStatus
	Carrier    string
	OrderID    int64
	SortBy     string
	SortOrder  string
}
