package application

import (
	"time"

	"github.com/qhato/ecommerce/internal/fulfillment/domain"
)

// ShipmentDTO represents shipment data for transfer
type ShipmentDTO struct {
	ID              int64      `json:"id"`
	OrderID         int64      `json:"order_id"`
	Status          string     `json:"status"`
	TrackingNumber  string     `json:"tracking_number,omitempty"`
	Carrier         string     `json:"carrier"`
	ShippingMethod  string     `json:"shipping_method"`
	ShippingCost    float64    `json:"shipping_cost"`
	EstimatedDate   *time.Time `json:"estimated_date,omitempty"`
	ShippedDate     *time.Time `json:"shipped_date,omitempty"`
	DeliveredDate   *time.Time `json:"delivered_date,omitempty"`
	ShippingAddress AddressDTO `json:"shipping_address"`
	Notes           string     `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// AddressDTO represents address data for transfer
type AddressDTO struct {
	Name       string `json:"name"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	Phone      string `json:"phone,omitempty"`
}

// CreateShipmentRequest represents a request to create a shipment
type CreateShipmentRequest struct {
	OrderID         int64      `json:"order_id" validate:"required"`
	Carrier         string     `json:"carrier" validate:"required"`
	ShippingMethod  string     `json:"shipping_method" validate:"required"`
	ShippingCost    float64    `json:"shipping_cost" validate:"required,min=0"`
	ShippingAddress AddressDTO `json:"shipping_address" validate:"required"`
}

// ShipRequest represents a request to mark shipment as shipped
type ShipRequest struct {
	TrackingNumber string `json:"tracking_number" validate:"required"`
}

// UpdateTrackingRequest represents a request to update tracking
type UpdateTrackingRequest struct {
	TrackingNumber string `json:"tracking_number"`
	Notes          string `json:"notes"`
}

// ToShipmentDTO converts domain Shipment to ShipmentDTO
func ToShipmentDTO(shipment *domain.Shipment) *ShipmentDTO {
	if shipment == nil {
		return nil
	}

	return &ShipmentDTO{
		ID:             shipment.ID,
		OrderID:        shipment.OrderID,
		Status:         string(shipment.Status),
		TrackingNumber: shipment.TrackingNumber,
		Carrier:        shipment.Carrier,
		ShippingMethod: shipment.ShippingMethod,
		ShippingCost:   shipment.ShippingCost,
		EstimatedDate:  shipment.EstimatedDate,
		ShippedDate:    shipment.ShippedDate,
		DeliveredDate:  shipment.DeliveredDate,
		ShippingAddress: AddressDTO{
			Name:       shipment.ShippingAddress.Name,
			Line1:      shipment.ShippingAddress.Line1,
			Line2:      shipment.ShippingAddress.Line2,
			City:       shipment.ShippingAddress.City,
			State:      shipment.ShippingAddress.State,
			PostalCode: shipment.ShippingAddress.PostalCode,
			Country:    shipment.ShippingAddress.Country,
			Phone:      shipment.ShippingAddress.Phone,
		},
		Notes:     shipment.Notes,
		CreatedAt: shipment.CreatedAt,
		UpdatedAt: shipment.UpdatedAt,
	}
}

// ToShipmentDTOs converts a slice of domain Shipments to ShipmentDTOs
func ToShipmentDTOs(shipments []*domain.Shipment) []ShipmentDTO {
	dtos := make([]ShipmentDTO, len(shipments))
	for i, shipment := range shipments {
		dto := ToShipmentDTO(shipment)
		if dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}
