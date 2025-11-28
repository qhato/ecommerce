package domain

import (
	"context"
)

// CustomerRepository defines the interface for customer persistence
type CustomerRepository interface {
	// Create creates a new customer
	Create(ctx context.Context, customer *Customer) error

	// Update updates an existing customer
	Update(ctx context.Context, customer *Customer) error

	// Delete deletes a customer by ID (soft delete - marks as archived)
	Delete(ctx context.Context, id int64) error

	// FindByID retrieves a customer by ID
	FindByID(ctx context.Context, id int64) (*Customer, error)

	// FindByEmail retrieves a customer by email address
	FindByEmail(ctx context.Context, email string) (*Customer, error)

	// FindByUsername retrieves a customer by username
	FindByUsername(ctx context.Context, username string) (*Customer, error)

	// FindAll retrieves all customers with pagination
	FindAll(ctx context.Context, filter *CustomerFilter) ([]*Customer, int64, error)

	// ExistsByEmail checks if a customer exists with given email
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername checks if a customer exists with given username
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// UpdatePassword updates customer password
	UpdatePassword(ctx context.Context, customerID int64, hashedPassword string) error
}

// AddressRepository defines the interface for address persistence
type AddressRepository interface {
	// Create creates a new address
	Create(ctx context.Context, address *Address) error

	// Update updates an existing address
	Update(ctx context.Context, address *Address) error

	// Delete deletes an address by ID
	Delete(ctx context.Context, id int64) error

	// FindByID retrieves an address by ID
	FindByID(ctx context.Context, id int64) (*Address, error)

	// FindByCustomerID retrieves addresses by customer ID
	FindByCustomerID(ctx context.Context, customerID int64) ([]*CustomerAddress, error)
}

// CustomerFilter represents filtering and pagination options for customers
type CustomerFilter struct {
	Page            int
	PageSize        int
	IncludeArchived bool
	ActiveOnly      bool
	RegisteredOnly  bool
	SortBy          string // "name", "email", "created_at"
	SortOrder       string // "asc", "desc"
	SearchQuery     string
}

// NewCustomerFilter creates a default customer filter
func NewCustomerFilter() *CustomerFilter {
	return &CustomerFilter{
		Page:            1,
		PageSize:        20,
		IncludeArchived: false,
		ActiveOnly:      true,
		RegisteredOnly:  false,
		SortBy:          "created_at",
		SortOrder:       "desc",
	}
}
