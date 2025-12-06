package queries

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/payment/domain"
)

type WebhookEventDTO struct {
	ID          string     `json:"id"`
	GatewayName string     `json:"gateway_name"`
	EventType   string     `json:"event_type"`
	EventID     string     `json:"event_id"`
	Status      string     `json:"status"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
	ErrorMsg    *string    `json:"error_msg,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type WebhookQueryService struct {
	webhookRepo domain.WebhookEventRepository
}

func NewWebhookQueryService(webhookRepo domain.WebhookEventRepository) *WebhookQueryService {
	return &WebhookQueryService{
		webhookRepo: webhookRepo,
	}
}

func (s *WebhookQueryService) GetWebhook(ctx context.Context, id string) (*WebhookEventDTO, error) {
	webhook, err := s.webhookRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find webhook: %w", err)
	}
	if webhook == nil {
		return nil, fmt.Errorf("webhook not found")
	}

	return toWebhookEventDTO(webhook), nil
}

func (s *WebhookQueryService) GetPendingWebhooks(ctx context.Context, limit int) ([]*WebhookEventDTO, error) {
	webhooks, err := s.webhookRepo.FindPending(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find pending webhooks: %w", err)
	}

	dtos := make([]*WebhookEventDTO, len(webhooks))
	for i, webhook := range webhooks {
		dtos[i] = toWebhookEventDTO(webhook)
	}

	return dtos, nil
}

func (s *WebhookQueryService) GetWebhooksByStatus(ctx context.Context, status string, limit int) ([]*WebhookEventDTO, error) {
	webhooks, err := s.webhookRepo.FindByStatus(ctx, domain.WebhookStatus(status), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find webhooks by status: %w", err)
	}

	dtos := make([]*WebhookEventDTO, len(webhooks))
	for i, webhook := range webhooks {
		dtos[i] = toWebhookEventDTO(webhook)
	}

	return dtos, nil
}

func toWebhookEventDTO(webhook *domain.WebhookEvent) *WebhookEventDTO {
	return &WebhookEventDTO{
		ID:          webhook.ID,
		GatewayName: webhook.GatewayName,
		EventType:   string(webhook.EventType),
		EventID:     webhook.EventID,
		Status:      string(webhook.Status),
		ProcessedAt: webhook.ProcessedAt,
		ErrorMsg:    webhook.ErrorMsg,
		CreatedAt:   webhook.CreatedAt,
		UpdatedAt:   webhook.UpdatedAt,
	}
}
