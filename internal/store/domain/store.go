package domain

import (
	"errors"
	"time"
)

// Store represents a physical or virtual store
type Store struct {
	ID            int64
	TenantID      int64  // Multi-tenancy support
	Code          string
	Name          string
	Description   string
	Type          StoreType
	Status        StoreStatus
	Email         string
	Phone         string
	Website       string
	Address       Address
	Timezone      string
	Currency      string
	Locale        string
	TaxID         string
	Settings      StoreSettings
	Metadata      map[string]interface{}
	ParentStoreID *int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// StoreType represents the type of store
type StoreType string

const (
	StoreTypePhysical StoreType = "PHYSICAL"
	StoreTypeOnline   StoreType = "ONLINE"
	StoreTypeHybrid   StoreType = "HYBRID"
)

// StoreStatus represents the status of a store
type StoreStatus string

const (
	StoreStatusActive   StoreStatus = "ACTIVE"
	StoreStatusInactive StoreStatus = "INACTIVE"
	StoreStatusClosed   StoreStatus = "CLOSED"
)

// Address represents a store address
type Address struct {
	Street1    string  `json:"street1"`
	Street2    string  `json:"street2,omitempty"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
	PostalCode string  `json:"postal_code"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
}

// StoreSettings represents store-specific settings
type StoreSettings struct {
	AllowPickup          bool    `json:"allow_pickup"`
	AllowShipping        bool    `json:"allow_shipping"`
	AllowBackorder       bool    `json:"allow_backorder"`
	InventoryTracking    bool    `json:"inventory_tracking"`
	DefaultShippingCost  float64 `json:"default_shipping_cost"`
	FreeShippingThreshold float64 `json:"free_shipping_threshold,omitempty"`
	MinOrderAmount       float64 `json:"min_order_amount,omitempty"`
	MaxOrderAmount       float64 `json:"max_order_amount,omitempty"`
	BusinessHours        []BusinessHour `json:"business_hours,omitempty"`
}

// BusinessHour represents operating hours for a day
type BusinessHour struct {
	DayOfWeek int    `json:"day_of_week"` // 0=Sunday, 1=Monday, etc.
	OpenTime  string `json:"open_time"`   // Format: "09:00"
	CloseTime string `json:"close_time"`  // Format: "18:00"
	IsClosed  bool   `json:"is_closed"`
}

// StoreInventory represents inventory at a specific store
type StoreInventory struct {
	ID             int64
	StoreID        int64
	ProductID      int64
	SKU            string
	QuantityOnHand int
	Reserved       int
	Available      int
	ReorderPoint   int
	ReorderQuantity int
	LastRestocked  *time.Time
	UpdatedAt      time.Time
}

// NewStore creates a new store
func NewStore(code, name string, storeType StoreType) (*Store, error) {
	if code == "" {
		return nil, errors.New("store code is required")
	}
	if name == "" {
		return nil, errors.New("store name is required")
	}

	now := time.Now()
	return &Store{
		Code:     code,
		Name:     name,
		Type:     storeType,
		Status:   StoreStatusActive,
		Metadata: make(map[string]interface{}),
		Settings: StoreSettings{
			AllowPickup:       true,
			AllowShipping:     true,
			InventoryTracking: true,
			BusinessHours:     make([]BusinessHour, 0),
		},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Activate activates the store
func (s *Store) Activate() {
	s.Status = StoreStatusActive
	s.UpdatedAt = time.Now()
}

// Deactivate deactivates the store
func (s *Store) Deactivate() {
	s.Status = StoreStatusInactive
	s.UpdatedAt = time.Now()
}

// Close closes the store permanently
func (s *Store) Close() {
	s.Status = StoreStatusClosed
	s.UpdatedAt = time.Now()
}

// IsOpen checks if store is currently open
func (s *Store) IsOpen() bool {
	if s.Status != StoreStatusActive {
		return false
	}

	now := time.Now()
	dayOfWeek := int(now.Weekday())

	for _, bh := range s.Settings.BusinessHours {
		if bh.DayOfWeek == dayOfWeek && !bh.IsClosed {
			// TODO: Add time comparison logic
			return true
		}
	}

	return len(s.Settings.BusinessHours) == 0 // Open if no hours configured
}

// UpdateInventory updates store inventory
func (si *StoreInventory) UpdateInventory(quantityChange int) {
	si.QuantityOnHand += quantityChange
	si.Available = si.QuantityOnHand - si.Reserved
	si.UpdatedAt = time.Now()
}

// Reserve reserves inventory
func (si *StoreInventory) Reserve(quantity int) error {
	if quantity > si.Available {
		return ErrInsufficientInventory
	}
	si.Reserved += quantity
	si.Available = si.QuantityOnHand - si.Reserved
	si.UpdatedAt = time.Now()
	return nil
}

// Release releases reserved inventory
func (si *StoreInventory) Release(quantity int) {
	if quantity > si.Reserved {
		quantity = si.Reserved
	}
	si.Reserved -= quantity
	si.Available = si.QuantityOnHand - si.Reserved
	si.UpdatedAt = time.Now()
}

// NeedsReorder checks if inventory needs reordering
func (si *StoreInventory) NeedsReorder() bool {
	return si.Available <= si.ReorderPoint
}
