package workflow

import (
	"context"
	"sync"
)

// ProcessContext holds the state and data for a workflow execution
type ProcessContext interface {
	// Context returns the underlying context.Context
	Context() context.Context

	// SetContext sets a new context
	SetContext(ctx context.Context)

	// Get retrieves a value from the context by key
	Get(key string) (interface{}, bool)

	// Set stores a value in the context by key
	Set(key string, value interface{})

	// Delete removes a value from the context
	Delete(key string)

	// StopProcess signals that the workflow should stop
	StopProcess()

	// IsStopped returns whether the workflow has been stopped
	IsStopped() bool

	// SetError sets an error in the context
	SetError(err error)

	// GetError returns the error if any
	GetError() error

	// SeedData returns initial data for the workflow
	SeedData() interface{}

	// SetSeedData sets the initial data
	SetSeedData(data interface{})
}

// DefaultProcessContext is the default implementation of ProcessContext
type DefaultProcessContext struct {
	ctx      context.Context
	data     map[string]interface{}
	seedData interface{}
	stopped  bool
	err      error
	mu       sync.RWMutex
}

// NewProcessContext creates a new ProcessContext
func NewProcessContext(ctx context.Context, seedData interface{}) ProcessContext {
	return &DefaultProcessContext{
		ctx:      ctx,
		data:     make(map[string]interface{}),
		seedData: seedData,
		stopped:  false,
	}
}

func (pc *DefaultProcessContext) Context() context.Context {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.ctx
}

func (pc *DefaultProcessContext) SetContext(ctx context.Context) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.ctx = ctx
}

func (pc *DefaultProcessContext) Get(key string) (interface{}, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	val, ok := pc.data[key]
	return val, ok
}

func (pc *DefaultProcessContext) Set(key string, value interface{}) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.data[key] = value
}

func (pc *DefaultProcessContext) Delete(key string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	delete(pc.data, key)
}

func (pc *DefaultProcessContext) StopProcess() {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.stopped = true
}

func (pc *DefaultProcessContext) IsStopped() bool {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.stopped
}

func (pc *DefaultProcessContext) SetError(err error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.err = err
}

func (pc *DefaultProcessContext) GetError() error {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.err
}

func (pc *DefaultProcessContext) SeedData() interface{} {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.seedData
}

func (pc *DefaultProcessContext) SetSeedData(data interface{}) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.seedData = data
}
