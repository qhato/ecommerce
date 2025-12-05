package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/email/domain"
	"github.com/qhato/ecommerce/pkg/events"
	"github.com/qhato/ecommerce/pkg/logger"
)

// EmailSender defines the interface for sending emails
type EmailSender interface {
	Send(ctx context.Context, email *domain.Email) error
}

// EmailRepository defines the interface for email persistence
type EmailRepository interface {
	Update(ctx context.Context, email *domain.Email) error
}

// Worker processes emails from the queue
type Worker struct {
	queue      *RedisQueue
	sender     EmailSender
	repository EmailRepository
	eventBus   events.EventBus
	logger     logger.Logger
	stopCh     chan struct{}
	doneCh     chan struct{}
}

// WorkerConfig holds worker configuration
type WorkerConfig struct {
	PollInterval       time.Duration
	RecoveryInterval   time.Duration
	StalledEmailMaxAge time.Duration
	MaxConcurrency     int
}

// NewWorker creates a new email worker
func NewWorker(
	queue *RedisQueue,
	sender EmailSender,
	repository EmailRepository,
	eventBus events.EventBus,
	logger logger.Logger,
) *Worker {
	return &Worker{
		queue:      queue,
		sender:     sender,
		repository: repository,
		eventBus:   eventBus,
		logger:     logger,
		stopCh:     make(chan struct{}),
		doneCh:     make(chan struct{}),
	}
}

// Start starts the worker
func (w *Worker) Start(ctx context.Context, config *WorkerConfig) {
	if config == nil {
		config = &WorkerConfig{
			PollInterval:       5 * time.Second,
			RecoveryInterval:   5 * time.Minute,
			StalledEmailMaxAge: 15 * time.Minute,
			MaxConcurrency:     5,
		}
	}

	w.logger.Info("Starting email worker",
		logger.Field{Key: "poll_interval", Value: config.PollInterval},
		logger.Field{Key: "max_concurrency", Value: config.MaxConcurrency},
	)

	// Start recovery goroutine
	go w.runRecovery(ctx, config.RecoveryInterval, config.StalledEmailMaxAge)

	// Start worker goroutines
	for i := 0; i < config.MaxConcurrency; i++ {
		go w.processQueue(ctx, config.PollInterval, i)
	}

	<-w.stopCh
	close(w.doneCh)
}

// Stop stops the worker
func (w *Worker) Stop() {
	w.logger.Info("Stopping email worker")
	close(w.stopCh)
	<-w.doneCh
}

// processQueue processes emails from the queue
func (w *Worker) processQueue(ctx context.Context, pollInterval time.Duration, workerID int) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	w.logger.Info("Worker started",
		logger.Field{Key: "worker_id", Value: workerID},
	)

	for {
		select {
		case <-w.stopCh:
			w.logger.Info("Worker stopped",
				logger.Field{Key: "worker_id", Value: workerID},
			)
			return
		case <-ticker.C:
			if err := w.processNextEmail(ctx); err != nil {
				w.logger.Error("Failed to process email",
					logger.Field{Key: "worker_id", Value: workerID},
					logger.Field{Key: "error", Value: err.Error()},
				)
			}
		}
	}
}

// processNextEmail processes the next email in the queue
func (w *Worker) processNextEmail(ctx context.Context) error {
	// Dequeue next email
	email, err := w.queue.Dequeue(ctx)
	if err != nil {
		return fmt.Errorf("failed to dequeue: %w", err)
	}

	if email == nil {
		// Queue is empty
		return nil
	}

	w.logger.Info("Processing email",
		logger.Field{Key: "email_id", Value: email.ID},
		logger.Field{Key: "type", Value: email.Type},
		logger.Field{Key: "retry_count", Value: email.RetryCount},
	)

	// Check if email is ready to send (scheduled)
	if !email.IsReadyToSend() {
		// Re-queue if scheduled for future
		if err := w.queue.Enqueue(ctx, email); err != nil {
			w.logger.Error("Failed to re-queue scheduled email",
				logger.Field{Key: "email_id", Value: email.ID},
				logger.Field{Key: "error", Value: err.Error()},
			)
		}
		w.queue.MarkAsProcessed(ctx, email)
		return nil
	}

	// Send email
	if err := w.sender.Send(ctx, email); err != nil {
		w.logger.Error("Failed to send email",
			logger.Field{Key: "email_id", Value: email.ID},
			logger.Field{Key: "error", Value: err.Error()},
		)

		// Mark as failed and potentially re-queue
		email.MarkAsFailed(err.Error())
		if err := w.queue.MarkAsFailed(ctx, email); err != nil {
			w.logger.Error("Failed to mark email as failed",
				logger.Field{Key: "email_id", Value: email.ID},
				logger.Field{Key: "error", Value: err.Error()},
			)
		}

		// Update repository
		if err := w.repository.Update(ctx, email); err != nil {
			w.logger.Error("Failed to update failed email",
				logger.Field{Key: "email_id", Value: email.ID},
				logger.Field{Key: "error", Value: err.Error()},
			)
		}

		// Publish failed event
		event := domain.NewEmailFailedEvent(email)
		w.eventBus.Publish(ctx, "email.failed", event)

		return nil
	}

	// Mark as processed
	if err := w.queue.MarkAsProcessed(ctx, email); err != nil {
		w.logger.Error("Failed to mark email as processed",
			logger.Field{Key: "email_id", Value: email.ID},
			logger.Field{Key: "error", Value: err.Error()},
		)
	}

	// Update repository
	if err := w.repository.Update(ctx, email); err != nil {
		w.logger.Error("Failed to update sent email",
			logger.Field{Key: "email_id", Value: email.ID},
			logger.Field{Key: "error", Value: err.Error()},
		)
	}

	// Publish sent event
	event := domain.NewEmailSentEvent(email)
	w.eventBus.Publish(ctx, "email.sent", event)

	w.logger.Info("Email sent successfully",
		logger.Field{Key: "email_id", Value: email.ID},
		logger.Field{Key: "to", Value: email.To},
	)

	return nil
}

// runRecovery periodically recovers stalled emails
func (w *Worker) runRecovery(ctx context.Context, interval, maxAge time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	w.logger.Info("Starting recovery routine",
		logger.Field{Key: "interval", Value: interval},
		logger.Field{Key: "max_age", Value: maxAge},
	)

	for {
		select {
		case <-w.stopCh:
			w.logger.Info("Recovery routine stopped")
			return
		case <-ticker.C:
			recovered, err := w.queue.RecoverStalledEmails(ctx, maxAge)
			if err != nil {
				w.logger.Error("Failed to recover stalled emails",
					logger.Field{Key: "error", Value: err.Error()},
				)
				continue
			}

			if recovered > 0 {
				w.logger.Info("Recovered stalled emails",
					logger.Field{Key: "count", Value: recovered},
				)
			}
		}
	}
}
