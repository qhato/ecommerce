package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/notification/domain"
)

type NotificationCommandHandler struct {
	notificationRepo domain.NotificationRepository
}

func NewNotificationCommandHandler(notificationRepo domain.NotificationRepository) *NotificationCommandHandler {
	return &NotificationCommandHandler{notificationRepo: notificationRepo}
}

func (h *NotificationCommandHandler) HandleCreateNotification(ctx context.Context, cmd CreateNotificationCommand) (*domain.Notification, error) {
	notification, err := domain.NewNotification(
		domain.NotificationType(cmd.Type),
		cmd.RecipientID,
		cmd.Subject,
		cmd.Body,
	)
	if err != nil {
		return nil, err
	}

	notification.RecipientEmail = cmd.RecipientEmail
	notification.RecipientPhone = cmd.RecipientPhone

	if cmd.Priority != "" {
		notification.Priority = domain.NotificationPriority(cmd.Priority)
	}

	if cmd.TemplateID != nil {
		notification.TemplateID = cmd.TemplateID
	}

	if cmd.TemplateData != nil {
		notification.TemplateData = cmd.TemplateData
	}

	if cmd.ScheduledFor != nil {
		notification.ScheduledFor = cmd.ScheduledFor
	}

	if err := h.notificationRepo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

func (h *NotificationCommandHandler) HandleSendNotification(ctx context.Context, cmd SendNotificationCommand) (*domain.Notification, error) {
	notification, err := h.notificationRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find notification: %w", err)
	}
	if notification == nil {
		return nil, domain.ErrNotificationNotFound
	}

	if !notification.IsReadyToSend() {
		return nil, fmt.Errorf("notification is not ready to send")
	}

	notification.MarkAsSending()
	if err := h.notificationRepo.Update(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to update notification: %w", err)
	}

	return notification, nil
}

func (h *NotificationCommandHandler) HandleMarkAsSent(ctx context.Context, cmd MarkAsSentCommand) (*domain.Notification, error) {
	notification, err := h.notificationRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find notification: %w", err)
	}
	if notification == nil {
		return nil, domain.ErrNotificationNotFound
	}

	notification.MarkAsSent()
	if err := h.notificationRepo.Update(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to update notification: %w", err)
	}

	return notification, nil
}

func (h *NotificationCommandHandler) HandleMarkAsFailed(ctx context.Context, cmd MarkAsFailedCommand) (*domain.Notification, error) {
	notification, err := h.notificationRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find notification: %w", err)
	}
	if notification == nil {
		return nil, domain.ErrNotificationNotFound
	}

	notification.MarkAsFailed(cmd.ErrorMsg)
	if err := h.notificationRepo.Update(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to update notification: %w", err)
	}

	return notification, nil
}

func (h *NotificationCommandHandler) HandleRetryNotification(ctx context.Context, cmd RetryNotificationCommand) (*domain.Notification, error) {
	notification, err := h.notificationRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find notification: %w", err)
	}
	if notification == nil {
		return nil, domain.ErrNotificationNotFound
	}

	if !notification.CanRetry() {
		return nil, fmt.Errorf("notification cannot be retried")
	}

	notification.Status = domain.NotificationStatusPending
	if err := h.notificationRepo.Update(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to retry notification: %w", err)
	}

	return notification, nil
}

func (h *NotificationCommandHandler) HandleCancelNotification(ctx context.Context, cmd CancelNotificationCommand) (*domain.Notification, error) {
	notification, err := h.notificationRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find notification: %w", err)
	}
	if notification == nil {
		return nil, domain.ErrNotificationNotFound
	}

	notification.Cancel()
	if err := h.notificationRepo.Update(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to cancel notification: %w", err)
	}

	return notification, nil
}

func (h *NotificationCommandHandler) HandleDeleteNotification(ctx context.Context, cmd DeleteNotificationCommand) error {
	return h.notificationRepo.Delete(ctx, cmd.ID)
}
