package domain

import (
	"time"

	"github.com/qhato/ecommerce/pkg/event"
)

const (
	EventCustomerRegistered      = "customer.registered"
	EventCustomerUpdated         = "customer.updated"
	EventCustomerDeactivated     = "customer.deactivated"
	EventCustomerActivated       = "customer.activated"
	EventCustomerPasswordChanged = "customer.password_changed"
	EventCustomerArchived        = "customer.archived"
)

// CustomerRegisteredEvent is published when a customer registers
type CustomerRegisteredEvent struct {
	event.BaseEvent
	CustomerID   int64  `json:"customer_id"`
	EmailAddress string `json:"email_address"`
	UserName     string `json:"user_name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

// NewCustomerRegisteredEvent creates a new CustomerRegisteredEvent
func NewCustomerRegisteredEvent(customerID int64, email, username, firstName, lastName string) *CustomerRegisteredEvent {
	return &CustomerRegisteredEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventCustomerRegistered,
			Timestamp: time.Now(),
		},
		CustomerID:   customerID,
		EmailAddress: email,
		UserName:     username,
		FirstName:    firstName,
		LastName:     lastName,
	}
}

// Type returns the event type
func (e *CustomerRegisteredEvent) Type() string {
	return e.EventType
}

// CustomerUpdatedEvent is published when a customer is updated
type CustomerUpdatedEvent struct {
	event.BaseEvent
	CustomerID int64                  `json:"customer_id"`
	Changes    map[string]interface{} `json:"changes"`
}

// NewCustomerUpdatedEvent creates a new CustomerUpdatedEvent
func NewCustomerUpdatedEvent(customerID int64, changes map[string]interface{}) *CustomerUpdatedEvent {
	return &CustomerUpdatedEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventCustomerUpdated,
			Timestamp: time.Now(),
		},
		CustomerID: customerID,
		Changes:    changes,
	}
}

// Type returns the event type
func (e *CustomerUpdatedEvent) Type() string {
	return e.EventType
}

// CustomerDeactivatedEvent is published when a customer is deactivated
type CustomerDeactivatedEvent struct {
	event.BaseEvent
	CustomerID int64 `json:"customer_id"`
}

// NewCustomerDeactivatedEvent creates a new CustomerDeactivatedEvent
func NewCustomerDeactivatedEvent(customerID int64) *CustomerDeactivatedEvent {
	return &CustomerDeactivatedEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventCustomerDeactivated,
			Timestamp: time.Now(),
		},
		CustomerID: customerID,
	}
}

// Type returns the event type
func (e *CustomerDeactivatedEvent) Type() string {
	return e.EventType
}

// CustomerActivatedEvent is published when a customer is activated
type CustomerActivatedEvent struct {
	event.BaseEvent
	CustomerID int64 `json:"customer_id"`
}

// NewCustomerActivatedEvent creates a new CustomerActivatedEvent
func NewCustomerActivatedEvent(customerID int64) *CustomerActivatedEvent {
	return &CustomerActivatedEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventCustomerActivated,
			Timestamp: time.Now(),
		},
		CustomerID: customerID,
	}
}

// Type returns the event type
func (e *CustomerActivatedEvent) Type() string {
	return e.EventType
}

// CustomerPasswordChangedEvent is published when password is changed
type CustomerPasswordChangedEvent struct {
	event.BaseEvent
	CustomerID int64 `json:"customer_id"`
}

// NewCustomerPasswordChangedEvent creates a new CustomerPasswordChangedEvent
func NewCustomerPasswordChangedEvent(customerID int64) *CustomerPasswordChangedEvent {
	return &CustomerPasswordChangedEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventCustomerPasswordChanged,
			Timestamp: time.Now(),
		},
		CustomerID: customerID,
	}
}

// Type returns the event type
func (e *CustomerPasswordChangedEvent) Type() string {
	return e.EventType
}
