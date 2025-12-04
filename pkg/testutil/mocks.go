package testutil

import (
	"context"
	"sync"

	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
)

// MockEventBus is a mock implementation of event.Bus for testing
type MockEventBus struct {
	mu        sync.Mutex
	events    []event.Event
	Published int
}

// NewMockEventBus creates a new mock event bus
func NewMockEventBus() *MockEventBus {
	return &MockEventBus{
		events: make([]event.Event, 0),
	}
}

// Publish records the event
func (m *MockEventBus) Publish(ctx context.Context, e event.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.events = append(m.events, e)
	m.Published++
	return nil
}

// Subscribe is a no-op for testing
func (m *MockEventBus) Subscribe(eventType string, handler event.Handler) error {
	return nil
}

// GetEvents returns all published events
func (m *MockEventBus) GetEvents() []event.Event {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]event.Event{}, m.events...)
}

// GetEventsByType returns events of a specific type
func (m *MockEventBus) GetEventsByType(eventType string) []event.Event {
	m.mu.Lock()
	defer m.mu.Unlock()

	filtered := make([]event.Event, 0)
	for _, e := range m.events {
		if e.Type() == eventType {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// Reset clears all recorded events
func (m *MockEventBus) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.events = make([]event.Event, 0)
	m.Published = 0
}

// MockLogger is a mock logger for testing
type MockLogger struct {
	mu      sync.Mutex
	logs    []string
	InfoMsg []string
	ErrMsg  []string
}

// NewMockLogger creates a new mock logger
func NewMockLogger() *MockLogger {
	return &MockLogger{
		logs:    make([]string, 0),
		InfoMsg: make([]string, 0),
		ErrMsg:  make([]string, 0),
	}
}

// Info records an info log
func (m *MockLogger) Info(msg string, fields ...logger.Field) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = append(m.logs, msg)
	m.InfoMsg = append(m.InfoMsg, msg)
}

// Error records an error log
func (m *MockLogger) Error(msg string, fields ...logger.Field) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = append(m.logs, msg)
	m.ErrMsg = append(m.ErrMsg, msg)
}

// Warn records a warning log
func (m *MockLogger) Warn(msg string, fields ...logger.Field) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = append(m.logs, msg)
}

// Debug records a debug log
func (m *MockLogger) Debug(msg string, fields ...logger.Field) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = append(m.logs, msg)
}

// GetLogs returns all recorded logs
func (m *MockLogger) GetLogs() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]string{}, m.logs...)
}

// Reset clears all recorded logs
func (m *MockLogger) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = make([]string, 0)
	m.InfoMsg = make([]string, 0)
	m.ErrMsg = make([]string, 0)
}
