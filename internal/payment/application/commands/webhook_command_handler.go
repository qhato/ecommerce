package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/payment/domain"
)

type WebhookCommandHandler struct {
	webhookRepo domain.WebhookEventRepository
	paymentRepo domain.PaymentRepository
}

func NewWebhookCommandHandler(
	webhookRepo domain.WebhookEventRepository,
	paymentRepo domain.PaymentRepository,
) *WebhookCommandHandler {
	return &WebhookCommandHandler{
		webhookRepo: webhookRepo,
		paymentRepo: paymentRepo,
	}
}

func (h *WebhookCommandHandler) HandleProcessWebhook(ctx context.Context, cmd ProcessWebhookCommand) (*domain.WebhookEvent, error) {
	// Check if event already exists (idempotency)
	existingEvent, err := h.webhookRepo.FindByEventID(ctx, cmd.GatewayName, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing webhook: %w", err)
	}
	if existingEvent != nil {
		// Event already processed, return existing
		return existingEvent, nil
	}

	// Create new webhook event
	event := domain.NewWebhookEvent(cmd.GatewayName, cmd.EventID, cmd.EventType, cmd.Payload)

	if cmd.Signature != nil {
		event.SetSignature(*cmd.Signature)
	}
	if cmd.IPAddress != nil {
		event.SetIPAddress(*cmd.IPAddress)
	}

	if err := h.webhookRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create webhook event: %w", err)
	}

	// Process the webhook based on event type
	if err := h.processWebhookEvent(ctx, event); err != nil {
		event.MarkAsFailed(err.Error())
		if updateErr := h.webhookRepo.Update(ctx, event); updateErr != nil {
			return nil, fmt.Errorf("failed to mark webhook as failed: %w", updateErr)
		}
		return event, err
	}

	event.MarkAsProcessed()
	if err := h.webhookRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to mark webhook as processed: %w", err)
	}

	return event, nil
}

func (h *WebhookCommandHandler) HandleRetryWebhook(ctx context.Context, cmd RetryWebhookCommand) (*domain.WebhookEvent, error) {
	event, err := h.webhookRepo.FindByID(ctx, cmd.WebhookID)
	if err != nil {
		return nil, fmt.Errorf("failed to find webhook: %w", err)
	}
	if event == nil {
		return nil, fmt.Errorf("webhook not found")
	}

	if event.Status == domain.WebhookStatusProcessed {
		return event, nil // Already processed
	}

	// Retry processing
	if err := h.processWebhookEvent(ctx, event); err != nil {
		event.MarkAsFailed(err.Error())
		if updateErr := h.webhookRepo.Update(ctx, event); updateErr != nil {
			return nil, fmt.Errorf("failed to mark webhook as failed: %w", updateErr)
		}
		return event, err
	}

	event.MarkAsProcessed()
	if err := h.webhookRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to mark webhook as processed: %w", err)
	}

	return event, nil
}

func (h *WebhookCommandHandler) processWebhookEvent(ctx context.Context, event *domain.WebhookEvent) error {
	// Process webhook based on event type
	switch event.EventType {
	case domain.WebhookEventPaymentSucceeded:
		return h.handlePaymentSucceeded(ctx, event)
	case domain.WebhookEventPaymentFailed:
		return h.handlePaymentFailed(ctx, event)
	case domain.WebhookEventPaymentRefunded:
		return h.handlePaymentRefunded(ctx, event)
	case domain.WebhookEventPaymentCancelled:
		return h.handlePaymentCancelled(ctx, event)
	case domain.WebhookEventChargebackCreated:
		return h.handleChargebackCreated(ctx, event)
	default:
		// Unknown event type, mark as ignored
		event.MarkAsIgnored()
		return nil
	}
}

func (h *WebhookCommandHandler) handlePaymentSucceeded(ctx context.Context, event *domain.WebhookEvent) error {
	// TODO: Parse payload and update payment status
	// This is a placeholder implementation
	return nil
}

func (h *WebhookCommandHandler) handlePaymentFailed(ctx context.Context, event *domain.WebhookEvent) error {
	// TODO: Parse payload and mark payment as failed
	return nil
}

func (h *WebhookCommandHandler) handlePaymentRefunded(ctx context.Context, event *domain.WebhookEvent) error {
	// TODO: Parse payload and process refund
	return nil
}

func (h *WebhookCommandHandler) handlePaymentCancelled(ctx context.Context, event *domain.WebhookEvent) error {
	// TODO: Parse payload and cancel payment
	return nil
}

func (h *WebhookCommandHandler) handleChargebackCreated(ctx context.Context, event *domain.WebhookEvent) error {
	// TODO: Handle chargeback creation
	return nil
}
