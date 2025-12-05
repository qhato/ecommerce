package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/email/domain"
)

// EmailQueryService provides query operations for emails
type EmailQueryService struct {
	emailRepo domain.EmailRepository
}

// NewEmailQueryService creates a new email query service
func NewEmailQueryService(emailRepo domain.EmailRepository) *EmailQueryService {
	return &EmailQueryService{
		emailRepo: emailRepo,
	}
}

// GetEmailByID retrieves an email by ID
func (s *EmailQueryService) GetEmailByID(ctx context.Context, id int64) (*EmailDTO, error) {
	email, err := s.emailRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find email: %w", err)
	}

	return mapEmailToDTO(email), nil
}

// ListEmailsByStatus retrieves emails by status with pagination
func (s *EmailQueryService) ListEmailsByStatus(ctx context.Context, status string, offset, limit int) ([]*EmailDTO, error) {
	emails, err := s.emailRepo.FindByStatus(ctx, domain.EmailStatus(status), offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list emails by status: %w", err)
	}

	dtos := make([]*EmailDTO, len(emails))
	for i, email := range emails {
		dtos[i] = mapEmailToDTO(email)
	}

	return dtos, nil
}

// ListEmailsByType retrieves emails by type with pagination
func (s *EmailQueryService) ListEmailsByType(ctx context.Context, emailType string, offset, limit int) ([]*EmailDTO, error) {
	emails, err := s.emailRepo.FindByType(ctx, domain.EmailType(emailType), offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list emails by type: %w", err)
	}

	dtos := make([]*EmailDTO, len(emails))
	for i, email := range emails {
		dtos[i] = mapEmailToDTO(email)
	}

	return dtos, nil
}

// ListEmailsByOrderID retrieves emails associated with an order
func (s *EmailQueryService) ListEmailsByOrderID(ctx context.Context, orderID int64) ([]*EmailDTO, error) {
	emails, err := s.emailRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to list emails by order: %w", err)
	}

	dtos := make([]*EmailDTO, len(emails))
	for i, email := range emails {
		dtos[i] = mapEmailToDTO(email)
	}

	return dtos, nil
}

// ListEmailsByCustomerID retrieves emails associated with a customer
func (s *EmailQueryService) ListEmailsByCustomerID(ctx context.Context, customerID int64, offset, limit int) ([]*EmailDTO, error) {
	emails, err := s.emailRepo.FindByCustomerID(ctx, customerID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list emails by customer: %w", err)
	}

	dtos := make([]*EmailDTO, len(emails))
	for i, email := range emails {
		dtos[i] = mapEmailToDTO(email)
	}

	return dtos, nil
}

// GetEmailStats retrieves email statistics
func (s *EmailQueryService) GetEmailStats(ctx context.Context) (*EmailStatsDTO, error) {
	total, err := s.emailRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count total emails: %w", err)
	}

	pending, err := s.emailRepo.CountByStatus(ctx, domain.EmailStatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to count pending emails: %w", err)
	}

	queued, err := s.emailRepo.CountByStatus(ctx, domain.EmailStatusQueued)
	if err != nil {
		return nil, fmt.Errorf("failed to count queued emails: %w", err)
	}

	sent, err := s.emailRepo.CountByStatus(ctx, domain.EmailStatusSent)
	if err != nil {
		return nil, fmt.Errorf("failed to count sent emails: %w", err)
	}

	failed, err := s.emailRepo.CountByStatus(ctx, domain.EmailStatusFailed)
	if err != nil {
		return nil, fmt.Errorf("failed to count failed emails: %w", err)
	}

	return &EmailStatsDTO{
		Total:   total,
		Pending: pending,
		Queued:  queued,
		Sent:    sent,
		Failed:  failed,
	}, nil
}

// mapEmailToDTO maps domain email to DTO
func mapEmailToDTO(email *domain.Email) *EmailDTO {
	dto := &EmailDTO{
		ID:           email.ID,
		Type:         string(email.Type),
		Status:       string(email.Status),
		Priority:     int(email.Priority),
		From:         email.From,
		To:           email.To,
		CC:           email.CC,
		BCC:          email.BCC,
		ReplyTo:      email.ReplyTo,
		Subject:      email.Subject,
		TemplateName: email.TemplateName,
		MaxRetries:   email.MaxRetries,
		RetryCount:   email.RetryCount,
		ErrorMessage: email.ErrorMessage,
		OrderID:      email.OrderID,
		CustomerID:   email.CustomerID,
		CreatedAt:    email.CreatedAt,
		UpdatedAt:    email.UpdatedAt,
		ScheduledAt:  email.ScheduledAt,
		SentAt:       email.SentAt,
		FailedAt:     email.FailedAt,
	}

	if len(email.Attachments) > 0 {
		dto.AttachmentCount = len(email.Attachments)
		dto.HasAttachments = true
	}

	return dto
}
