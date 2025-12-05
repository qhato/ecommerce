package commands

import "time"

// SendEmailCommand represents a command to send an email
type SendEmailCommand struct {
	Type         string
	From         string
	To           []string
	CC           []string
	BCC          []string
	ReplyTo      string
	Subject      string
	Body         string
	HTMLBody     string
	TemplateName string
	TemplateData map[string]interface{}
	Priority     int
	OrderID      *int64
	CustomerID   *int64
	Attachments  []AttachmentData
}

// AttachmentData represents attachment data for command
type AttachmentData struct {
	Filename    string
	ContentType string
	Content     []byte
}

// ScheduleEmailCommand represents a command to schedule an email
type ScheduleEmailCommand struct {
	SendEmailCommand
	ScheduledAt time.Time
}

// SendOrderConfirmationCommand represents a command to send order confirmation email
type SendOrderConfirmationCommand struct {
	OrderID      int64
	CustomerID   int64
	To           string
	OrderNumber  string
	OrderTotal   float64
	OrderDate    time.Time
	Items        []OrderItemData
	ShippingAddr AddressData
	BillingAddr  AddressData
}

// OrderItemData represents order item data for email
type OrderItemData struct {
	ProductName string
	SKU         string
	Quantity    int
	Price       float64
	Total       float64
}

// AddressData represents address data for email
type AddressData struct {
	FirstName string
	LastName  string
	Line1     string
	Line2     string
	City      string
	State     string
	ZipCode   string
	Country   string
}

// SendOrderShippedCommand represents a command to send order shipped email
type SendOrderShippedCommand struct {
	OrderID        int64
	CustomerID     int64
	To             string
	OrderNumber    string
	TrackingNumber string
	Carrier        string
	ShippedAt      time.Time
	EstimatedDeliveryDate *time.Time
}

// SendPasswordResetCommand represents a command to send password reset email
type SendPasswordResetCommand struct {
	CustomerID int64
	To         string
	ResetToken string
	ExpiresAt  time.Time
}

// SendWelcomeEmailCommand represents a command to send welcome email
type SendWelcomeEmailCommand struct {
	CustomerID int64
	To         string
	FirstName  string
	LastName   string
}

// SendCartAbandonmentCommand represents a command to send cart abandonment email
type SendCartAbandonmentCommand struct {
	CustomerID  int64
	To          string
	CartID      int64
	Items       []OrderItemData
	Total       float64
	AbandonedAt time.Time
}

// CancelEmailCommand represents a command to cancel a pending email
type CancelEmailCommand struct {
	EmailID int64
}

// RetryFailedEmailCommand represents a command to retry a failed email
type RetryFailedEmailCommand struct {
	EmailID int64
}
