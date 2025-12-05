package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// OrderAdjustment represents an adjustment applied to an order
// Business Logic: Descuentos aplicados a nivel de orden completa
type OrderAdjustment struct {
	ID            int64
	OrderID       int64
	OfferID       int64
	OfferName     string
	AdjustmentValue decimal.Decimal // Monto del descuento
	AdjustmentReason string          // "OFFER_DISCOUNT", "MANUAL_ADJUSTMENT"
	AppliedDate   time.Time
	CreatedAt     time.Time
}

// OrderItemAdjustment represents an adjustment applied to an order item
// Business Logic: Descuentos aplicados a items específicos
type OrderItemAdjustment struct {
	ID              int64
	OrderItemID     int64
	OfferID         int64
	OfferName       string
	AdjustmentValue decimal.Decimal
	Quantity        int             // Cantidad de items afectados
	AppliedDate     time.Time
	CreatedAt       time.Time
}

// FulfillmentGroupAdjustment represents an adjustment applied to fulfillment (shipping)
// Business Logic: Descuentos en envío (free shipping, $5 off shipping)
type FulfillmentGroupAdjustment struct {
	ID                   int64
	FulfillmentGroupID   int64
	OfferID              int64
	OfferName            string
	AdjustmentValue      decimal.Decimal
	AdjustmentReason     string
	AppliedDate          time.Time
	CreatedAt            time.Time
}

// OrderAdjustmentRepository defines the interface for order adjustment persistence
type OrderAdjustmentRepository interface {
	CreateOrderAdjustment(adj *OrderAdjustment) error
	CreateOrderItemAdjustment(adj *OrderItemAdjustment) error
	CreateFulfillmentAdjustment(adj *FulfillmentGroupAdjustment) error
	FindByOrderID(orderID int64) ([]*OrderAdjustment, error)
	FindItemAdjustmentsByOrderID(orderID int64) ([]*OrderItemAdjustment, error)
	FindFulfillmentAdjustmentsByOrderID(orderID int64) ([]*FulfillmentGroupAdjustment, error)
	DeleteByOrderID(orderID int64) error
}

// GetAdjustmentValue returns the adjustment value as float64
func (oa *OrderAdjustment) GetAdjustmentValue() float64 {
	val, _ := oa.AdjustmentValue.Float64()
	return val
}

// GetAdjustmentValue returns the adjustment value as float64
func (oia *OrderItemAdjustment) GetAdjustmentValue() float64 {
	val, _ := oia.AdjustmentValue.Float64()
	return val
}

// GetAdjustmentValue returns the adjustment value as float64
func (fga *FulfillmentGroupAdjustment) GetAdjustmentValue() float64 {
	val, _ := fga.AdjustmentValue.Float64()
	return val
}
