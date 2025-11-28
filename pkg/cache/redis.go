package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/redis/go-redis/v9"
)

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client *redis.Client
	prefix string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Database int
	PoolSize int
	Prefix   string // Key prefix for namespacing
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(cfg RedisConfig) (*RedisCache, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.Database,
		PoolSize: cfg.PoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis cache connected successfully")

	return &RedisCache{
		client: client,
		prefix: cfg.Prefix,
	}, nil
}

// prefixKey adds prefix to key
func (rc *RedisCache) prefixKey(key string) string {
	if rc.prefix == "" {
		return key
	}
	return rc.prefix + ":" + key
}

// Get retrieves a value from Redis
func (rc *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := rc.client.Get(ctx, rc.prefixKey(key)).Bytes()
	if err == redis.Nil {
		return nil, &ErrCacheMiss{Key: key}
	}
	if err != nil {
		return nil, fmt.Errorf("redis get error: %w", err)
	}
	return val, nil
}

// Set stores a value in Redis with TTL
func (rc *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if err := rc.client.Set(ctx, rc.prefixKey(key), value, ttl).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}
	return nil
}

// Delete removes a value from Redis
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	if err := rc.client.Del(ctx, rc.prefixKey(key)).Err(); err != nil {
		return fmt.Errorf("redis delete error: %w", err)
	}
	return nil
}

// Exists checks if a key exists in Redis
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := rc.client.Exists(ctx, rc.prefixKey(key)).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists error: %w", err)
	}
	return result > 0, nil
}

// Expire sets a TTL on an existing key
func (rc *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if err := rc.client.Expire(ctx, rc.prefixKey(key), ttl).Err(); err != nil {
		return fmt.Errorf("redis expire error: %w", err)
	}
	return nil
}

// Clear removes all keys with the configured prefix
func (rc *RedisCache) Clear(ctx context.Context) error {
	pattern := rc.prefix + ":*"
	if rc.prefix == "" {
		pattern = "*"
	}

	iter := rc.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := rc.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("redis clear error: %w", err)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("redis scan error: %w", err)
	}

	return nil
}

// Close closes the Redis connection
func (rc *RedisCache) Close() error {
	if rc.client != nil {
		if err := rc.client.Close(); err != nil {
			return fmt.Errorf("failed to close redis connection: %w", err)
		}
		logger.Info("Redis cache connection closed")
	}
	return nil
}

// Health checks Redis health
func (rc *RedisCache) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := rc.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis unhealthy: %w", err)
	}

	return nil
}

// GetClient returns the underlying Redis client
func (rc *RedisCache) GetClient() *redis.Client {
	return rc.client
}
