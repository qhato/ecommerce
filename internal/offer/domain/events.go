package domain

import "time"

// OfferCreatedEvent is published when a new offer is successfully created.
type OfferCreatedEvent struct {
	OfferID      int64
	Name         string
	OfferType    OfferType
	Value        float64
	CreationTime time.Time
}

// OfferUpdatedEvent is published when an existing offer is modified.
type OfferUpdatedEvent struct {
	OfferID    int64
	Name       string
	UpdateTime time.Time
}

// OfferActivatedEvent is published when an offer becomes active.
type OfferActivatedEvent struct {
	OfferID        int64
	ActivationTime time.Time
}

// OfferDeactivatedEvent is published when an offer becomes inactive.
type OfferDeactivatedEvent struct {
	OfferID          int64
	DeactivationTime time.Time
}

// OfferUsedEvent is published when an offer is successfully applied to an order.
type OfferUsedEvent struct {
	OfferID    int64
	OrderID    int64
	CustomerID int64
	UsageTime  time.Time
}

// OfferDeletedEvent is published when an offer is removed.
type OfferDeletedEvent struct {
	OfferID      int64
	DeletionTime time.Time
}
