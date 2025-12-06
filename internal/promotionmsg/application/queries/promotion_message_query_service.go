package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/promotionmsg/domain"
)

type PromotionMessageQueryService struct {
	messageRepo domain.PromotionMessageRepository
}

func NewPromotionMessageQueryService(messageRepo domain.PromotionMessageRepository) *PromotionMessageQueryService {
	return &PromotionMessageQueryService{messageRepo: messageRepo}
}

func (s *PromotionMessageQueryService) GetMessage(ctx context.Context, id int64) (*PromotionMessageDTO, error) {
	message, err := s.messageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find message: %w", err)
	}
	if message == nil {
		return nil, domain.ErrMessageNotFound
	}

	return ToPromotionMessageDTO(message), nil
}

func (s *PromotionMessageQueryService) GetMessagesByType(ctx context.Context, messageType string) ([]*PromotionMessageDTO, error) {
	messages, err := s.messageRepo.FindByType(ctx, domain.MessageType(messageType))
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}

	dtos := make([]*PromotionMessageDTO, len(messages))
	for i, msg := range messages {
		dtos[i] = ToPromotionMessageDTO(msg)
	}

	return dtos, nil
}

func (s *PromotionMessageQueryService) GetMessagesByStatus(ctx context.Context, status string) ([]*PromotionMessageDTO, error) {
	messages, err := s.messageRepo.FindByStatus(ctx, domain.MessageStatus(status))
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}

	dtos := make([]*PromotionMessageDTO, len(messages))
	for i, msg := range messages {
		dtos[i] = ToPromotionMessageDTO(msg)
	}

	return dtos, nil
}

func (s *PromotionMessageQueryService) GetActiveMessages(ctx context.Context, limit int) ([]*PromotionMessageDTO, error) {
	messages, err := s.messageRepo.FindActive(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find active messages: %w", err)
	}

	dtos := make([]*PromotionMessageDTO, 0)
	for _, msg := range messages {
		if msg.IsActive() {
			dtos = append(dtos, ToPromotionMessageDTO(msg))
		}
	}

	return dtos, nil
}

func (s *PromotionMessageQueryService) GetMessagesByPlacement(ctx context.Context, placement string) ([]*PromotionMessageDTO, error) {
	messages, err := s.messageRepo.FindByPlacement(ctx, placement)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}

	dtos := make([]*PromotionMessageDTO, 0)
	for _, msg := range messages {
		if msg.IsActive() {
			dtos = append(dtos, ToPromotionMessageDTO(msg))
		}
	}

	return dtos, nil
}

func (s *PromotionMessageQueryService) GetMessagesByEvent(ctx context.Context, event string, context map[string]interface{}) ([]*PromotionMessageDTO, error) {
	messages, err := s.messageRepo.FindByEvent(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}

	dtos := make([]*PromotionMessageDTO, 0)
	for _, msg := range messages {
		if msg.IsActive() && msg.MatchesTrigger(event, context) && msg.MatchesRules(context) {
			dtos = append(dtos, ToPromotionMessageDTO(msg))
		}
	}

	return dtos, nil
}

func (s *PromotionMessageQueryService) GetMatchingMessages(ctx context.Context, placement string, context map[string]interface{}) ([]*PromotionMessageDTO, error) {
	messages, err := s.messageRepo.FindByPlacement(ctx, placement)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}

	dtos := make([]*PromotionMessageDTO, 0)
	for _, msg := range messages {
		if msg.IsActive() && msg.MatchesRules(context) {
			dtos = append(dtos, ToPromotionMessageDTO(msg))
		}
	}

	// Sort by priority
	for i := 0; i < len(dtos)-1; i++ {
		for j := i + 1; j < len(dtos); j++ {
			if dtos[i].Priority < dtos[j].Priority {
				dtos[i], dtos[j] = dtos[j], dtos[i]
			}
		}
	}

	return dtos, nil
}
