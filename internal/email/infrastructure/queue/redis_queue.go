package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/qhato/ecommerce/internal/email/domain"
	"github.com/qhato/ecommerce/pkg/logger"
)

const (
	emailQueueKey         = "email:queue"
	emailProcessingSetKey = "email:processing"
	emailQueueStatsKey    = "email:queue:stats"
)

// RedisQueue implements email queue using Redis
type RedisQueue struct {
	client *redis.Client
	logger logger.Logger
}

// NewRedisQueue creates a new Redis-based email queue
func NewRedisQueue(client *redis.Client, logger logger.Logger) *RedisQueue {
	return &RedisQueue{
		client: client,
		logger: logger,
	}
}

// Enqueue adds an email to the queue
func (q *RedisQueue) Enqueue(ctx context.Context, email *domain.Email) error {
	// Serialize email
	data, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("failed to marshal email: %w", err)
	}

	// Add to queue with priority (higher priority = lower score)
	score := calculateScore(email)

	if err := q.client.ZAdd(ctx, emailQueueKey, redis.Z{
		Score:  score,
		Member: string(data),
	}).Err(); err != nil {
		return fmt.Errorf("failed to add to queue: %w", err)
	}

	// Update stats
	q.incrementStat(ctx, "enqueued")

	q.logger.Info("Email enqueued",
		logger.Field{Key: "email_id", Value: email.ID},
		logger.Field{Key: "priority", Value: email.Priority},
		logger.Field{Key: "score", Value: score},
	)

	return nil
}

// Dequeue retrieves the next email from the queue
func (q *RedisQueue) Dequeue(ctx context.Context) (*domain.Email, error) {
	// Get the item with the lowest score (highest priority)
	results, err := q.client.ZPopMin(ctx, emailQueueKey, 1).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Queue is empty
		}
		return nil, fmt.Errorf("failed to dequeue: %w", err)
	}

	if len(results) == 0 {
		return nil, nil // Queue is empty
	}

	// Deserialize email
	var email domain.Email
	data := results[0].Member.(string)
	if err := json.Unmarshal([]byte(data), &email); err != nil {
		return nil, fmt.Errorf("failed to unmarshal email: %w", err)
	}

	// Move to processing set (for crash recovery)
	if err := q.client.SAdd(ctx, emailProcessingSetKey, data).Err(); err != nil {
		q.logger.Warn("Failed to add to processing set",
			logger.Field{Key: "email_id", Value: email.ID},
			logger.Field{Key: "error", Value: err.Error()},
		)
	}

	// Update stats
	q.incrementStat(ctx, "dequeued")

	q.logger.Info("Email dequeued",
		logger.Field{Key: "email_id", Value: email.ID},
		logger.Field{Key: "type", Value: email.Type},
	)

	return &email, nil
}

// MarkAsProcessed removes an email from the processing set
func (q *RedisQueue) MarkAsProcessed(ctx context.Context, email *domain.Email) error {
	data, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("failed to marshal email: %w", err)
	}

	if err := q.client.SRem(ctx, emailProcessingSetKey, string(data)).Err(); err != nil {
		return fmt.Errorf("failed to remove from processing set: %w", err)
	}

	q.incrementStat(ctx, "processed")

	return nil
}

// MarkAsFailed marks an email as failed and optionally re-queues it
func (q *RedisQueue) MarkAsFailed(ctx context.Context, email *domain.Email) error {
	data, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("failed to marshal email: %w", err)
	}

	// Remove from processing set
	if err := q.client.SRem(ctx, emailProcessingSetKey, string(data)).Err(); err != nil {
		return fmt.Errorf("failed to remove from processing set: %w", err)
	}

	// Re-queue if retries are available
	if email.CanRetry() {
		// Add exponential backoff
		delay := calculateRetryDelay(email.RetryCount)
		score := float64(time.Now().Add(delay).Unix())

		if err := q.client.ZAdd(ctx, emailQueueKey, redis.Z{
			Score:  score,
			Member: string(data),
		}).Err(); err != nil {
			return fmt.Errorf("failed to re-queue email: %w", err)
		}

		q.logger.Info("Email re-queued for retry",
			logger.Field{Key: "email_id", Value: email.ID},
			logger.Field{Key: "retry_count", Value: email.RetryCount},
			logger.Field{Key: "delay", Value: delay.String()},
		)

		q.incrementStat(ctx, "requeued")
	} else {
		q.incrementStat(ctx, "failed")
	}

	return nil
}

// Size returns the number of emails in the queue
func (q *RedisQueue) Size(ctx context.Context) (int, error) {
	count, err := q.client.ZCard(ctx, emailQueueKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue size: %w", err)
	}
	return int(count), nil
}

// ProcessingCount returns the number of emails being processed
func (q *RedisQueue) ProcessingCount(ctx context.Context) (int, error) {
	count, err := q.client.SCard(ctx, emailProcessingSetKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get processing count: %w", err)
	}
	return int(count), nil
}

// GetStats returns queue statistics
func (q *RedisQueue) GetStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Get all stats fields
	fields := []string{"enqueued", "dequeued", "processed", "requeued", "failed"}
	for _, field := range fields {
		val, err := q.client.HGet(ctx, emailQueueStatsKey, field).Int64()
		if err != nil && err != redis.Nil {
			return nil, fmt.Errorf("failed to get stat %s: %w", field, err)
		}
		stats[field] = val
	}

	// Add current queue size
	size, err := q.Size(ctx)
	if err != nil {
		return nil, err
	}
	stats["queue_size"] = int64(size)

	// Add processing count
	processing, err := q.ProcessingCount(ctx)
	if err != nil {
		return nil, err
	}
	stats["processing"] = int64(processing)

	return stats, nil
}

// ResetStats resets queue statistics
func (q *RedisQueue) ResetStats(ctx context.Context) error {
	return q.client.Del(ctx, emailQueueStatsKey).Err()
}

// Clear clears the entire queue (use with caution!)
func (q *RedisQueue) Clear(ctx context.Context) error {
	pipe := q.client.Pipeline()
	pipe.Del(ctx, emailQueueKey)
	pipe.Del(ctx, emailProcessingSetKey)
	_, err := pipe.Exec(ctx)
	return err
}

// RecoverStalledEmails recovers emails that were being processed but failed
func (q *RedisQueue) RecoverStalledEmails(ctx context.Context, maxAge time.Duration) (int, error) {
	// Get all emails in processing set
	members, err := q.client.SMembers(ctx, emailProcessingSetKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get processing set: %w", err)
	}

	recovered := 0
	for _, member := range members {
		var email domain.Email
		if err := json.Unmarshal([]byte(member), &email); err != nil {
			q.logger.Warn("Failed to unmarshal stalled email",
				logger.Field{Key: "error", Value: err.Error()},
			)
			continue
		}

		// Check if email has been processing for too long
		if time.Since(email.UpdatedAt) > maxAge {
			// Remove from processing set
			q.client.SRem(ctx, emailProcessingSetKey, member)

			// Re-queue
			email.MarkAsRetrying()
			score := calculateScore(&email)
			if err := q.client.ZAdd(ctx, emailQueueKey, redis.Z{
				Score:  score,
				Member: member,
			}).Err(); err != nil {
				q.logger.Warn("Failed to re-queue stalled email",
					logger.Field{Key: "email_id", Value: email.ID},
					logger.Field{Key: "error", Value: err.Error()},
				)
				continue
			}

			recovered++
			q.logger.Info("Recovered stalled email",
				logger.Field{Key: "email_id", Value: email.ID},
				logger.Field{Key: "age", Value: time.Since(email.UpdatedAt).String()},
			)
		}
	}

	return recovered, nil
}

// calculateScore calculates the priority score for an email
// Lower score = higher priority
func calculateScore(email *domain.Email) float64 {
	now := time.Now().Unix()

	// If scheduled, use scheduled time
	if email.ScheduledAt != nil && email.ScheduledAt.After(time.Now()) {
		return float64(email.ScheduledAt.Unix())
	}

	// Otherwise, use priority and current time
	// Higher priority emails get lower scores (processed first)
	priorityOffset := float64(1000 - email.Priority)
	return float64(now) - priorityOffset
}

// calculateRetryDelay calculates exponential backoff delay for retries
func calculateRetryDelay(retryCount int) time.Duration {
	// Exponential backoff: 1min, 5min, 15min
	delays := []time.Duration{
		1 * time.Minute,
		5 * time.Minute,
		15 * time.Minute,
	}

	if retryCount < len(delays) {
		return delays[retryCount]
	}

	return delays[len(delays)-1]
}

// incrementStat increments a stat counter
func (q *RedisQueue) incrementStat(ctx context.Context, field string) {
	if err := q.client.HIncrBy(ctx, emailQueueStatsKey, field, 1).Err(); err != nil {
		q.logger.Warn("Failed to increment stat",
			logger.Field{Key: "field", Value: field},
			logger.Field{Key: "error", Value: err.Error()},
		)
	}
}
