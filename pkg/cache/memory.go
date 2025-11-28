package cache

import (
	"context"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// MemoryCache implements Cache interface using in-memory storage
// Useful for testing and development
type MemoryCache struct {
	cache *gocache.Cache
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache(defaultTTL, cleanupInterval time.Duration) *MemoryCache {
	return &MemoryCache{
		cache: gocache.New(defaultTTL, cleanupInterval),
	}
}

// Get retrieves a value from memory cache
func (mc *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, found := mc.cache.Get(key)
	if !found {
		return nil, &ErrCacheMiss{Key: key}
	}
	return val.([]byte), nil
}

// Set stores a value in memory cache with TTL
func (mc *MemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	mc.cache.Set(key, value, ttl)
	return nil
}

// Delete removes a value from memory cache
func (mc *MemoryCache) Delete(ctx context.Context, key string) error {
	mc.cache.Delete(key)
	return nil
}

// Exists checks if a key exists in memory cache
func (mc *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	_, found := mc.cache.Get(key)
	return found, nil
}

// Expire sets a TTL on an existing key
func (mc *MemoryCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	val, found := mc.cache.Get(key)
	if !found {
		return &ErrCacheMiss{Key: key}
	}
	mc.cache.Set(key, val, ttl)
	return nil
}

// Clear removes all keys from memory cache
func (mc *MemoryCache) Clear(ctx context.Context) error {
	mc.cache.Flush()
	return nil
}

// Close closes the memory cache (no-op for memory cache)
func (mc *MemoryCache) Close() error {
	return nil
}

// Health checks memory cache health (always healthy)
func (mc *MemoryCache) Health(ctx context.Context) error {
	return nil
}

// ItemCount returns the number of items in the cache
func (mc *MemoryCache) ItemCount() int {
	return mc.cache.ItemCount()
}
