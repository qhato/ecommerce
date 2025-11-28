package domain

import (
	"time"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusConfirmed  OrderStatus = "CONFIRMED"
	OrderStatusShipped    OrderStatus = "SHIPPED"
	OrderStatusDelivered  OrderStatus = "DELIVERED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
	OrderStatusRefunded   OrderStatus = "REFUNDED"
)

// Order represents an order entity
type Order struct {
	ID            int64
	OrderNumber   string
	CustomerID    int64
	EmailAddress  string
	Name          string
	Status        OrderStatus
	SubTotal      float64
	TotalTax      float64
	TotalShipping float64
	Total         float64
	CurrencyCode  string
	Items         []OrderItem
	SubmitDate    *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID             int64
	OrderID        int64
	SKUID          int64
	ProductName    string
	Quantity       int
	Price          float64
	TotalPrice     float64
	TaxAmount      float64
	ShippingAmount float64
}

// NewOrder creates a new order
func NewOrder(customerID int64, emailAddress, name, currencyCode string) *Order {
	now := time.Now()
	return &Order{
		CustomerID:   customerID,
		EmailAddress: emailAddress,
		Name:         name,
		Status:       OrderStatusPending,
		CurrencyCode: currencyCode,
		Items:        make([]OrderItem, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddItem adds an item to the order
func (o *Order) AddItem(skuID int64, productName string, quantity int, price float64) {
	item := OrderItem{
		OrderID:     o.ID,
		SKUID:       skuID,
		ProductName: productName,
		Quantity:    quantity,
		Price:       price,
		TotalPrice:  price * float64(quantity),
	}
	o.Items = append(o.Items, item)
	o.CalculateTotals()
}

// CalculateTotals recalculates order totals
func (o *Order) CalculateTotals() {
	o.SubTotal = 0
	for _, item := range o.Items {
		o.SubTotal += item.TotalPrice
	}
	o.Total = o.SubTotal + o.TotalTax + o.TotalShipping
	o.UpdatedAt = time.Now()
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(status OrderStatus) {
	o.Status = status
	o.UpdatedAt = time.Now()
}

// Submit submits the order
func (o *Order) Submit() {
	now := time.Now()
	o.SubmitDate = &now
	o.Status = OrderStatusProcessing
	o.UpdatedAt = now
}

// Cancel cancels the order
func (o *Order) Cancel() {
	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
}

// IsCancellable checks if order can be cancelled
func (o *Order) IsCancellable() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusProcessing
}
