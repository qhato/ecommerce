package domain

import "context"

// CheckoutSessionRepository defines the interface for checkout session persistence
type CheckoutSessionRepository interface {
	// Create creates a new checkout session
	Create(ctx context.Context, session *CheckoutSession) error

	// Update updates an existing checkout session
	Update(ctx context.Context, session *CheckoutSession) error

	// FindByID finds a checkout session by ID
	FindByID(ctx context.Context, id string) (*CheckoutSession, error)

	// FindByOrderID finds a checkout session by order ID
	FindByOrderID(ctx context.Context, orderID int64) (*CheckoutSession, error)

	// FindByCustomerID finds checkout sessions by customer ID
	FindByCustomerID(ctx context.Context, customerID string, activeOnly bool) ([]*CheckoutSession, error)

	// FindActiveByEmail finds active checkout sessions by email
	FindActiveByEmail(ctx context.Context, email string) ([]*CheckoutSession, error)

	// FindExpiredSessions finds expired sessions
	FindExpiredSessions(ctx context.Context, limit int) ([]*CheckoutSession, error)

	// Delete deletes a checkout session
	Delete(ctx context.Context, id string) error

	// ExistsByOrderID checks if a checkout session exists for an order
	ExistsByOrderID(ctx context.Context, orderID int64) (bool, error)
}

// ShippingOptionRepository defines the interface for shipping option persistence
type ShippingOptionRepository interface {
	// FindByID finds a shipping option by ID
	FindByID(ctx context.Context, id string) (*ShippingOption, error)

	// FindAll finds all shipping options
	FindAll(ctx context.Context, activeOnly bool) ([]*ShippingOption, error)

	// FindByCarrier finds shipping options by carrier
	FindByCarrier(ctx context.Context, carrier string, activeOnly bool) ([]*ShippingOption, error)

	// FindAvailableForAddress finds available shipping options for an address
	FindAvailableForAddress(ctx context.Context, country, stateProvince, postalCode string) ([]*ShippingOption, error)
}
