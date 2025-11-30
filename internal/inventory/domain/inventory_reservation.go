package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReservationStatus represents the status of an inventory reservation
type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "PENDING"
	ReservationStatusConfirmed ReservationStatus = "CONFIRMED"
	ReservationStatusReleased  ReservationStatus = "RELEASED"
	ReservationStatusExpired   ReservationStatus = "EXPIRED"
	ReservationStatusFulfilled ReservationStatus = "FULFILLED"
)

// InventoryReservation represents a reservation of inventory for an order
type InventoryReservation struct {
	ID             string
	SKUID          string
	Quantity       int
	OrderID        string
	OrderItemID    string
	Status         ReservationStatus
	ReservedAt     time.Time
	ExpiresAt      *time.Time
	ReleasedAt     *time.Time
	FulfilledAt    *time.Time
	ReservationRef string // Reference ID for external systems
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewInventoryReservation creates a new inventory reservation
func NewInventoryReservation(skuID, orderID, orderItemID string, quantity int, ttl time.Duration) (*InventoryReservation, error) {
	if skuID == "" || orderID == "" || orderItemID == "" {
		return nil, NewDomainError("SKUID, OrderID, and OrderItemID are required")
	}

	if quantity <= 0 {
		return nil, NewDomainError("Quantity must be positive")
	}

	now := time.Now()
	var expiresAt *time.Time
	if ttl > 0 {
		expiry := now.Add(ttl)
		expiresAt = &expiry
	}

	return &InventoryReservation{
		ID:             uuid.New().String(),
		SKUID:          skuID,
		Quantity:       quantity,
		OrderID:        orderID,
		OrderItemID:    orderItemID,
		Status:         ReservationStatusPending,
		ReservedAt:     now,
		ExpiresAt:      expiresAt,
		ReservationRef: uuid.New().String(),
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Confirm confirms the reservation
func (r *InventoryReservation) Confirm() error {
	if r.Status != ReservationStatusPending {
		return NewDomainError("Can only confirm pending reservations")
	}

	if r.IsExpired() {
		return NewDomainError("Cannot confirm expired reservation")
	}

	r.Status = ReservationStatusConfirmed
	r.UpdatedAt = time.Now()
	return nil
}

// Release releases the reserved inventory
func (r *InventoryReservation) Release() error {
	if r.Status == ReservationStatusReleased || r.Status == ReservationStatusFulfilled {
		return NewDomainError("Reservation already released or fulfilled")
	}

	now := time.Now()
	r.Status = ReservationStatusReleased
	r.ReleasedAt = &now
	r.UpdatedAt = now
	return nil
}

// Fulfill marks the reservation as fulfilled
func (r *InventoryReservation) Fulfill() error {
	if r.Status != ReservationStatusConfirmed && r.Status != ReservationStatusPending {
		return NewDomainError("Can only fulfill confirmed or pending reservations")
	}

	now := time.Now()
	r.Status = ReservationStatusFulfilled
	r.FulfilledAt = &now
	r.UpdatedAt = now
	return nil
}

// IsExpired checks if the reservation has expired
func (r *InventoryReservation) IsExpired() bool {
	if r.ExpiresAt == nil {
		return false
	}

	return time.Now().After(*r.ExpiresAt)
}

// Expire marks the reservation as expired
func (r *InventoryReservation) Expire() error {
	if r.Status == ReservationStatusFulfilled {
		return NewDomainError("Cannot expire fulfilled reservation")
	}

	if !r.IsExpired() {
		return NewDomainError("Reservation has not expired yet")
	}

	r.Status = ReservationStatusExpired
	r.UpdatedAt = time.Now()
	return nil
}

// ExtendExpiration extends the expiration time
func (r *InventoryReservation) ExtendExpiration(additionalTime time.Duration) error {
	if r.Status != ReservationStatusPending && r.Status != ReservationStatusConfirmed {
		return NewDomainError("Can only extend active reservations")
	}

	if r.ExpiresAt == nil {
		// No expiration set, set one
		expiry := time.Now().Add(additionalTime)
		r.ExpiresAt = &expiry
	} else {
		// Extend existing expiration
		newExpiry := r.ExpiresAt.Add(additionalTime)
		r.ExpiresAt = &newExpiry
	}

	r.UpdatedAt = time.Now()
	return nil
}