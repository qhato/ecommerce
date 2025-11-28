package cache

import (
	"context"
	"time"
)

// Cache defines the interface for cache operations
type Cache interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value in cache with TTL
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)

	// Expire sets a TTL on an existing key
	Expire(ctx context.Context, key string, ttl time.Duration) error

	// Clear removes all keys (use with caution)
	Clear(ctx context.Context) error

	// Close closes the cache connection
	Close() error

	// Health checks cache health
	Health(ctx context.Context) error
}

// ErrCacheMiss is returned when a key is not found in cache
type ErrCacheMiss struct {
	Key string
}

func (e *ErrCacheMiss) Error() string {
	return "cache miss: key not found: " + e.Key
}

// IsCacheMiss checks if error is a cache miss
func IsCacheMiss(err error) bool {
	_, ok := err.(*ErrCacheMiss)
	return ok
}
