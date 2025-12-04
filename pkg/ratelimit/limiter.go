package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Limiter is an interface for rate limiting
type Limiter interface {
	Allow(ctx context.Context, key string) (bool, error)
	Reset(ctx context.Context, key string) error
}

// Config contains rate limiter configuration
type Config struct {
	RequestsPerWindow int           // Number of requests allowed per window
	WindowSize        time.Duration // Time window for rate limiting
	KeyPrefix         string        // Prefix for Redis keys
}

// RedisLimiter implements rate limiting using Redis
type RedisLimiter struct {
	client *redis.Client
	config Config
}

// NewRedisLimiter creates a new Redis-based rate limiter
func NewRedisLimiter(client *redis.Client, config Config) *RedisLimiter {
	if config.KeyPrefix == "" {
		config.KeyPrefix = "ratelimit:"
	}
	return &RedisLimiter{
		client: client,
		config: config,
	}
}

// Allow checks if a request is allowed under the rate limit
func (l *RedisLimiter) Allow(ctx context.Context, key string) (bool, error) {
	redisKey := l.config.KeyPrefix + key
	now := time.Now().Unix()

	// Use Redis sorted set to track requests in the time window
	pipe := l.client.Pipeline()

	// Remove old entries outside the window
	windowStart := now - int64(l.config.WindowSize.Seconds())
	pipe.ZRemRangeByScore(ctx, redisKey, "0", fmt.Sprintf("%d", windowStart))

	// Count current entries
	countCmd := pipe.ZCard(ctx, redisKey)

	// Add current request
	pipe.ZAdd(ctx, redisKey, redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	})

	// Set expiry
	pipe.Expire(ctx, redisKey, l.config.WindowSize)

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to execute rate limit check: %w", err)
	}

	// Check if under limit
	count := countCmd.Val()
	return count < int64(l.config.RequestsPerWindow), nil
}

// Reset resets the rate limit for a key
func (l *RedisLimiter) Reset(ctx context.Context, key string) error {
	redisKey := l.config.KeyPrefix + key
	return l.client.Del(ctx, redisKey).Err()
}

// TokenBucketLimiter implements token bucket algorithm using Redis
type TokenBucketLimiter struct {
	client       *redis.Client
	capacity     int           // Maximum tokens
	refillRate   int           // Tokens added per second
	refillPeriod time.Duration // How often to refill
	keyPrefix    string
}

// NewTokenBucketLimiter creates a new token bucket rate limiter
func NewTokenBucketLimiter(client *redis.Client, capacity, refillRate int, refillPeriod time.Duration) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		client:       client,
		capacity:     capacity,
		refillRate:   refillRate,
		refillPeriod: refillPeriod,
		keyPrefix:    "tokenbucket:",
	}
}

// Allow checks if a request is allowed
func (l *TokenBucketLimiter) Allow(ctx context.Context, key string) (bool, error) {
	redisKey := l.keyPrefix + key
	now := time.Now().Unix()

	// Lua script for atomic token bucket operation
	script := `
		local key = KEYS[1]
		local capacity = tonumber(ARGV[1])
		local refill_rate = tonumber(ARGV[2])
		local now = tonumber(ARGV[3])
		
		local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
		local tokens = tonumber(bucket[1])
		local last_refill = tonumber(bucket[2])
		
		if tokens == nil then
			tokens = capacity
			last_refill = now
		else
			local elapsed = now - last_refill
			local refill = math.floor(elapsed * refill_rate)
			tokens = math.min(capacity, tokens + refill)
			if refill > 0 then
				last_refill = now
			end
		end
		
		if tokens >= 1 then
			tokens = tokens - 1
			redis.call('HMSET', key, 'tokens', tokens, 'last_refill', last_refill)
			redis.call('EXPIRE', key, 3600)
			return 1
		else
			return 0
		end
	`

	result, err := l.client.Eval(ctx, script, []string{redisKey},
		l.capacity, float64(l.refillRate)/l.refillPeriod.Seconds(), now).Result()
	if err != nil {
		return false, fmt.Errorf("failed to execute token bucket: %w", err)
	}

	return result.(int64) == 1, nil
}

// Reset resets the token bucket
func (l *TokenBucketLimiter) Reset(ctx context.Context, key string) error {
	redisKey := l.keyPrefix + key
	return l.client.Del(ctx, redisKey).Err()
}

// MemoryLimiter implements in-memory rate limiting (for testing)
type MemoryLimiter struct {
	requests map[string][]time.Time
	config   Config
}

// NewMemoryLimiter creates a new in-memory rate limiter
func NewMemoryLimiter(config Config) *MemoryLimiter {
	return &MemoryLimiter{
		requests: make(map[string][]time.Time),
		config:   config,
	}
}

// Allow checks if a request is allowed
func (l *MemoryLimiter) Allow(ctx context.Context, key string) (bool, error) {
	now := time.Now()
	windowStart := now.Add(-l.config.WindowSize)

	// Get existing requests
	requests, exists := l.requests[key]
	if !exists {
		requests = []time.Time{}
	}

	// Filter requests within window
	validRequests := []time.Time{}
	for _, t := range requests {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	// Check limit
	if len(validRequests) >= l.config.RequestsPerWindow {
		return false, nil
	}

	// Add current request
	validRequests = append(validRequests, now)
	l.requests[key] = validRequests

	return true, nil
}

// Reset resets the rate limit for a key
func (l *MemoryLimiter) Reset(ctx context.Context, key string) error {
	delete(l.requests, key)
	return nil
}

// Cleanup removes expired entries
func (l *MemoryLimiter) Cleanup() {
	now := time.Now()
	windowStart := now.Add(-l.config.WindowSize)

	for key, requests := range l.requests {
		validRequests := []time.Time{}
		for _, t := range requests {
			if t.After(windowStart) {
				validRequests = append(validRequests, t)
			}
		}

		if len(validRequests) == 0 {
			delete(l.requests, key)
		} else {
			l.requests[key] = validRequests
		}
	}
}
