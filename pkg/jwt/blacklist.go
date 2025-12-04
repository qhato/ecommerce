package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisBlacklist implements Blacklist using Redis
type RedisBlacklist struct {
	client *redis.Client
	prefix string
}

// NewRedisBlacklist creates a new Redis-based blacklist
func NewRedisBlacklist(client *redis.Client) *RedisBlacklist {
	return &RedisBlacklist{
		client: client,
		prefix: "jwt:blacklist:",
	}
}

// Add adds a token to the blacklist
func (b *RedisBlacklist) Add(ctx context.Context, tokenID string, expiry time.Duration) error {
	key := b.prefix + tokenID
	err := b.client.Set(ctx, key, "1", expiry).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}
	return nil
}

// IsBlacklisted checks if a token is blacklisted
func (b *RedisBlacklist) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := b.prefix + tokenID
	result, err := b.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check blacklist: %w", err)
	}
	return result > 0, nil
}

// MemoryBlacklist implements Blacklist using in-memory map (for testing)
type MemoryBlacklist struct {
	tokens map[string]time.Time
}

// NewMemoryBlacklist creates a new in-memory blacklist
func NewMemoryBlacklist() *MemoryBlacklist {
	return &MemoryBlacklist{
		tokens: make(map[string]time.Time),
	}
}

// Add adds a token to the blacklist
func (b *MemoryBlacklist) Add(ctx context.Context, tokenID string, expiry time.Duration) error {
	b.tokens[tokenID] = time.Now().Add(expiry)
	return nil
}

// IsBlacklisted checks if a token is blacklisted
func (b *MemoryBlacklist) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	expiryTime, exists := b.tokens[tokenID]
	if !exists {
		return false, nil
	}

	// Remove if expired
	if time.Now().After(expiryTime) {
		delete(b.tokens, tokenID)
		return false, nil
	}

	return true, nil
}

// Cleanup removes expired tokens (should be called periodically)
func (b *MemoryBlacklist) Cleanup() {
	now := time.Now()
	for tokenID, expiryTime := range b.tokens {
		if now.After(expiryTime) {
			delete(b.tokens, tokenID)
		}
	}
}