package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/promotionmsg/domain"
)

type PromotionMessageCommandHandler struct {
	messageRepo domain.PromotionMessageRepository
}

func NewPromotionMessageCommandHandler(messageRepo domain.PromotionMessageRepository) *PromotionMessageCommandHandler {
	return &PromotionMessageCommandHandler{messageRepo: messageRepo}
}

func (h *PromotionMessageCommandHandler) HandleCreatePromotionMessage(ctx context.Context, cmd CreatePromotionMessageCommand) (*domain.PromotionMessage, error) {
	message, err := domain.NewPromotionMessage(
		cmd.Name,
		cmd.Message,
		domain.MessageType(cmd.Type),
		cmd.Priority,
	)
	if err != nil {
		return nil, err
	}

	message.Description = cmd.Description
	message.StartDate = cmd.StartDate
	message.EndDate = cmd.EndDate
	message.MaxViews = cmd.MaxViews
	message.Placements = cmd.Placements

	if cmd.Rules != nil {
		message.Rules = make([]domain.MessageRule, len(cmd.Rules))
		for i, r := range cmd.Rules {
			message.Rules[i] = domain.MessageRule{
				Field:    r.Field,
				Operator: r.Operator,
				Value:    r.Value,
				Metadata: r.Metadata,
			}
		}
	}

	if cmd.Triggers != nil {
		message.Triggers = make([]domain.MessageTrigger, len(cmd.Triggers))
		for i, t := range cmd.Triggers {
			conditions := make([]domain.MessageRule, len(t.Conditions))
			for j, c := range t.Conditions {
				conditions[j] = domain.MessageRule{
					Field:    c.Field,
					Operator: c.Operator,
					Value:    c.Value,
					Metadata: c.Metadata,
				}
			}
			message.Triggers[i] = domain.MessageTrigger{
				Event:      t.Event,
				Conditions: conditions,
				Metadata:   t.Metadata,
			}
		}
	}

	if cmd.Metadata != nil {
		message.Metadata = cmd.Metadata
	}

	if err := h.messageRepo.Create(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to create promotion message: %w", err)
	}

	return message, nil
}

func (h *PromotionMessageCommandHandler) HandleUpdatePromotionMessage(ctx context.Context, cmd UpdatePromotionMessageCommand) (*domain.PromotionMessage, error) {
	message, err := h.messageRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find message: %w", err)
	}
	if message == nil {
		return nil, domain.ErrMessageNotFound
	}

	message.Name = cmd.Name
	message.Message = cmd.Message
	message.Description = cmd.Description
	message.Priority = cmd.Priority
	message.StartDate = cmd.StartDate
	message.EndDate = cmd.EndDate
	message.MaxViews = cmd.MaxViews
	message.Placements = cmd.Placements

	if cmd.Rules != nil {
		message.Rules = make([]domain.MessageRule, len(cmd.Rules))
		for i, r := range cmd.Rules {
			message.Rules[i] = domain.MessageRule{
				Field:    r.Field,
				Operator: r.Operator,
				Value:    r.Value,
				Metadata: r.Metadata,
			}
		}
	}

	if cmd.Triggers != nil {
		message.Triggers = make([]domain.MessageTrigger, len(cmd.Triggers))
		for i, t := range cmd.Triggers {
			conditions := make([]domain.MessageRule, len(t.Conditions))
			for j, c := range t.Conditions {
				conditions[j] = domain.MessageRule{
					Field:    c.Field,
					Operator: c.Operator,
					Value:    c.Value,
					Metadata: c.Metadata,
				}
			}
			message.Triggers[i] = domain.MessageTrigger{
				Event:      t.Event,
				Conditions: conditions,
				Metadata:   t.Metadata,
			}
		}
	}

	if cmd.Metadata != nil {
		message.Metadata = cmd.Metadata
	}

	message.UpdatedAt = time.Now()

	if err := h.messageRepo.Update(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to update message: %w", err)
	}

	return message, nil
}

func (h *PromotionMessageCommandHandler) HandleActivateMessage(ctx context.Context, cmd ActivateMessageCommand) (*domain.PromotionMessage, error) {
	message, err := h.messageRepo.FindByID(ctx, cmd.ID)
	if err != nil || message == nil {
		return nil, domain.ErrMessageNotFound
	}

	message.Activate()

	if err := h.messageRepo.Update(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to activate message: %w", err)
	}

	return message, nil
}

func (h *PromotionMessageCommandHandler) HandleDeactivateMessage(ctx context.Context, cmd DeactivateMessageCommand) (*domain.PromotionMessage, error) {
	message, err := h.messageRepo.FindByID(ctx, cmd.ID)
	if err != nil || message == nil {
		return nil, domain.ErrMessageNotFound
	}

	message.Deactivate()

	if err := h.messageRepo.Update(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to deactivate message: %w", err)
	}

	return message, nil
}

func (h *PromotionMessageCommandHandler) HandleIncrementView(ctx context.Context, cmd IncrementViewCommand) error {
	message, err := h.messageRepo.FindByID(ctx, cmd.ID)
	if err != nil || message == nil {
		return domain.ErrMessageNotFound
	}

	message.IncrementView()

	if err := h.messageRepo.Update(ctx, message); err != nil {
		return fmt.Errorf("failed to increment view: %w", err)
	}

	return nil
}

func (h *PromotionMessageCommandHandler) HandleIncrementClick(ctx context.Context, cmd IncrementClickCommand) error {
	message, err := h.messageRepo.FindByID(ctx, cmd.ID)
	if err != nil || message == nil {
		return domain.ErrMessageNotFound
	}

	message.IncrementClick()

	if err := h.messageRepo.Update(ctx, message); err != nil {
		return fmt.Errorf("failed to increment click: %w", err)
	}

	return nil
}

func (h *PromotionMessageCommandHandler) HandleDeleteMessage(ctx context.Context, cmd DeleteMessageCommand) error {
	return h.messageRepo.Delete(ctx, cmd.ID)
}
