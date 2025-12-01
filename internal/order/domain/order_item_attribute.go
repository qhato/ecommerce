package domain

import "time"

// OrderItemAttribute represents a custom attribute for an order item
type OrderItemAttribute struct {
	OrderItemID int64
	Name        string
	Value       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewOrderItemAttribute creates a new OrderItemAttribute
func NewOrderItemAttribute(orderItemID int64, name, value string) (*OrderItemAttribute, error) {
	if orderItemID == 0 {
		return nil, NewDomainError("OrderItemID cannot be zero for OrderItemAttribute")
	}
	if name == "" {
		return nil, NewDomainError("Name cannot be empty for OrderItemAttribute")
	}

	now := time.Now()
	return &OrderItemAttribute{
		OrderItemID: orderItemID,
		Name:        name,
		Value:       value,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UpdateValue updates the value of the order item attribute
func (oia *OrderItemAttribute) UpdateValue(value string) {
	oia.Value = value
	oia.UpdatedAt = time.Now()
}
