package application

import (
	"context"
)

// ShippingService defines the application service for shipping-related operations.
type ShippingService interface {
	// CalculateShippingCost calculates the shipping cost for a given order based on address and fulfillment option.
	CalculateShippingCost(ctx context.Context, orderID int64, shippingAddressID int64, fulfillmentOptionID int64) (float64, error)

	// ValidateShippingAddress validates a given shipping address.
	ValidateShippingAddress(ctx context.Context, addressID int64) (bool, error)

	// GetShippingMethods retrieves available shipping methods for an order/address combination.
	// This would typically return a list of shipping options with their costs.
	GetShippingMethods(ctx context.Context, orderID int64, shippingAddressID int64) ([]*ShippingMethodDTO, error)
}

// ShippingMethodDTO represents a shipping method data transfer object.
type ShippingMethodDTO struct {
	ID                  int64
	Name                string
	Description         string
	Cost                float64
	DeliveryEstimate    string
	FulfillmentOptionID int64
}
