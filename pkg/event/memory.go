package event

import (
	"context"
	"fmt"
	"sync"

	"github.com/qhato/ecommerce/pkg/logger"
)

// MemoryBus implements Bus interface using in-memory pub/sub
// Suitable for development and testing, or when running as a single instance
type MemoryBus struct {
	subscribers map[string][]Handler
	mu          sync.RWMutex
}

// NewMemoryBus creates a new in-memory event bus
func NewMemoryBus() *MemoryBus {
	return &MemoryBus{
		subscribers: make(map[string][]Handler),
	}
}

// Publish publishes an event to all subscribers
func (mb *MemoryBus) Publish(ctx context.Context, event Event) error {
	mb.mu.RLock()
	handlers, exists := mb.subscribers[event.EventType()]
	mb.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		logger.WithFields(logger.Fields{
			"event_type": event.EventType(),
			"event_id":   event.EventID(),
		}).Debug("No subscribers for event")
		return nil
	}

	logger.WithFields(logger.Fields{
		"event_type":   event.EventType(),
		"event_id":     event.EventID(),
		"aggregate_id": event.AggregateID(),
		"subscribers":  len(handlers),
	}).Debug("Publishing event")

	// Execute handlers concurrently
	var wg sync.WaitGroup
	errChan := make(chan error, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(h Handler) {
			defer wg.Done()
			if err := h(ctx, event); err != nil {
				logger.WithError(err).WithFields(logger.Fields{
					"event_type": event.EventType(),
					"event_id":   event.EventID(),
				}).Error("Event handler failed")
				errChan <- err
			}
		}(handler)
	}

	// Wait for all handlers to complete
	wg.Wait()
	close(errChan)

	// Collect errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("event handlers failed: %d errors occurred", len(errors))
	}

	return nil
}

// Subscribe subscribes a handler to an event type
func (mb *MemoryBus) Subscribe(eventType string, handler Handler) error {
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	mb.mu.Lock()
	defer mb.mu.Unlock()

	mb.subscribers[eventType] = append(mb.subscribers[eventType], handler)

	logger.WithField("event_type", eventType).Debug("Handler subscribed to event")
	return nil
}

// Unsubscribe removes a handler subscription (not fully implemented for memory bus)
func (mb *MemoryBus) Unsubscribe(eventType string, handler Handler) error {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	// Simple implementation: clear all handlers for the event type
	// In production, you'd want to track handler IDs for precise removal
	delete(mb.subscribers, eventType)

	logger.WithField("event_type", eventType).Debug("Handler unsubscribed from event")
	return nil
}

// Close closes the event bus
func (mb *MemoryBus) Close() error {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	mb.subscribers = make(map[string][]Handler)
	logger.Info("Memory event bus closed")
	return nil
}

// SubscriberCount returns the number of subscribers for an event type
func (mb *MemoryBus) SubscriberCount(eventType string) int {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	return len(mb.subscribers[eventType])
}
