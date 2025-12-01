package event

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Event represents a domain event
type Event interface {
	// EventType returns the type/name of the event
	EventType() string

	// OccurredAt returns when the event occurred
	OccurredAt() time.Time

	// AggregateID returns the ID of the aggregate that generated this event
	AggregateID() string

	// EventID returns a unique identifier for this event
	EventID() string
}

// BaseEvent provides a base implementation of Event
type BaseEvent struct {
	Type       string    `json:"type"`
	ID         string    `json:"id"`
	Aggregate  string    `json:"aggregate_id"`
	OccurredOn time.Time `json:"occurred_at"`
	Payload    any       `json:"payload"`
}

// NewBaseEvent creates a new base event
func NewBaseEvent(eventType string, aggregateID string, payload any) BaseEvent {
	return BaseEvent{
		Type:       eventType,
		ID:         uuid.New().String(),
		Aggregate:  aggregateID,
		OccurredOn: time.Now().UTC(),
		Payload:    payload,
	}
}

func (e BaseEvent) EventType() string {
	return e.Type
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

func (e BaseEvent) AggregateID() string {
	return e.Aggregate
}

func (e BaseEvent) EventID() string {
	return e.ID
}

// Handler is a function that handles an event
type Handler func(ctx context.Context, event Event) error

// Bus defines the interface for event bus operations
type Bus interface {
	// Publish publishes an event to all subscribers
	Publish(ctx context.Context, event Event) error

	// Subscribe subscribes a handler to an event type
	Subscribe(eventType string, handler Handler) error

	// Unsubscribe removes a handler subscription
	Unsubscribe(eventType string, handler Handler) error

	// Close closes the event bus
	Close() error
}
