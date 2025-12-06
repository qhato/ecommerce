package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/notification/domain"
)

type NotificationQueryService struct {
	notificationRepo domain.NotificationRepository
}

func NewNotificationQueryService(notificationRepo domain.NotificationRepository) *NotificationQueryService {
	return &NotificationQueryService{notificationRepo: notificationRepo}
}

func (s *NotificationQueryService) GetNotification(ctx context.Context, id int64) (*NotificationDTO, error) {
	notification, err := s.notificationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find notification: %w", err)
	}
	if notification == nil {
		return nil, domain.ErrNotificationNotFound
	}

	return ToNotificationDTO(notification), nil
}

func (s *NotificationQueryService) GetNotificationsByRecipient(ctx context.Context, recipientID string, limit int) ([]*NotificationDTO, error) {
	notifications, err := s.notificationRepo.FindByRecipientID(ctx, recipientID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find notifications: %w", err)
	}

	dtos := make([]*NotificationDTO, len(notifications))
	for i, n := range notifications {
		dtos[i] = ToNotificationDTO(n)
	}

	return dtos, nil
}

func (s *NotificationQueryService) GetNotificationsByStatus(ctx context.Context, status string, limit int) ([]*NotificationDTO, error) {
	notifications, err := s.notificationRepo.FindByStatus(ctx, domain.NotificationStatus(status), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find notifications: %w", err)
	}

	dtos := make([]*NotificationDTO, len(notifications))
	for i, n := range notifications {
		dtos[i] = ToNotificationDTO(n)
	}

	return dtos, nil
}

func (s *NotificationQueryService) GetPendingNotifications(ctx context.Context, limit int) ([]*NotificationDTO, error) {
	notifications, err := s.notificationRepo.FindPending(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find pending notifications: %w", err)
	}

	dtos := make([]*NotificationDTO, len(notifications))
	for i, n := range notifications {
		dtos[i] = ToNotificationDTO(n)
	}

	return dtos, nil
}

func (s *NotificationQueryService) GetScheduledNotifications(ctx context.Context, limit int) ([]*NotificationDTO, error) {
	notifications, err := s.notificationRepo.FindScheduled(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find scheduled notifications: %w", err)
	}

	dtos := make([]*NotificationDTO, len(notifications))
	for i, n := range notifications {
		dtos[i] = ToNotificationDTO(n)
	}

	return dtos, nil
}

func (s *NotificationQueryService) GetFailedNotifications(ctx context.Context, limit int) ([]*NotificationDTO, error) {
	notifications, err := s.notificationRepo.FindFailed(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find failed notifications: %w", err)
	}

	dtos := make([]*NotificationDTO, len(notifications))
	for i, n := range notifications {
		dtos[i] = ToNotificationDTO(n)
	}

	return dtos, nil
}
