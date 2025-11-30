package domain

import "time"

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending      OrderStatus = "PENDING"
	OrderStatusCustomerInfo OrderStatus = "CUSTOMER_INFO"
	OrderStatusShipping     OrderStatus = "SHIPPING"
	OrderStatusPayment      OrderStatus = "PAYMENT"
	OrderStatusReview       OrderStatus = "REVIEW"
	OrderStatusSubmitted    OrderStatus = "SUBMITTED"
	OrderStatusProcessing   OrderStatus = "PROCESSING"
	OrderStatusConfirmed    OrderStatus = "CONFIRMED"
	OrderStatusShipped      OrderStatus = "SHIPPED"
	OrderStatusDelivered    OrderStatus = "DELIVERED"
	OrderStatusCancelled    OrderStatus = "CANCELLED"
	OrderStatusRefunded     OrderStatus = "REFUNDED"
	OrderStatusFulfilled    OrderStatus = "FULFILLED"
)

// Order represents an order entity
type Order struct {
	ID            int64
	OrderNumber   string
	CustomerID    int64
	EmailAddress  string
	Name          string
	Status        OrderStatus
	OrderSubtotal float64 // From blc_order.order_subtotal
	TotalTax      float64
	TotalShipping float64
	OrderTotal    float64 // From blc_order.order_total
	CurrencyCode  string
	IsPreview     bool   // From blc_order.is_preview
	TaxOverride   bool   // From blc_order.tax_override
	LocaleCode    string // From blc_order.locale_code
	SubmitDate    *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Items         []OrderItem
}

// NewOrder creates a new order
func NewOrder(customerID int64, emailAddress, name, currencyCode, localeCode string) *Order {
	now := time.Now()
	return &Order{
		CustomerID:    customerID,
		EmailAddress:  emailAddress,
		Name:          name,
		Status:        OrderStatusPending,
		CurrencyCode:  currencyCode,
		LocaleCode:    localeCode,
		OrderSubtotal: 0.0,
		TotalTax:      0.0,
		TotalShipping: 0.0,
		OrderTotal:    0.0,
		IsPreview:     false,
		TaxOverride:   false,
		CreatedAt:     now,
		UpdatedAt:     now,
		Items:         make([]OrderItem, 0),
	}
}

// AddItem adds an item to the order
func (o *Order) AddItem(item OrderItem) {
	o.Items = append(o.Items, item)
	o.recalculateSubtotal()
	o.UpdatedAt = time.Now()
}

func (o *Order) recalculateSubtotal() {
	var subtotal float64
	for _, item := range o.Items {
		subtotal += item.TotalPrice
	}
	o.OrderSubtotal = subtotal
	o.recalculateTotal()
}

func (o *Order) recalculateTotal() {
	o.OrderTotal = o.OrderSubtotal + o.TotalTax + o.TotalShipping
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(status OrderStatus) {
	o.Status = status
	o.UpdatedAt = time.Now()
}

// Submit submits the order
func (o *Order) Submit() error {
	// Item check is now handled by application service ensuring items are added
	now := time.Now()
	o.SubmitDate = &now
	o.Status = OrderStatusProcessing
	o.UpdatedAt = now
	return nil
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

// OrderFilter represents filtering and pagination options for orders
type OrderFilter struct {
	Page       int
	PageSize   int
	CustomerID *int64
	Status     *OrderStatus
	SortBy     string
	SortOrder  string
}

// DomainError represents a business rule validation error within the domain.
type DomainError struct {
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new DomainError.
func NewDomainError(message string) error {
	return &DomainError{Message: message}
}
