package application

import (
	"context"

	"github.com/qhato/ecommerce/internal/email/application/commands"
	"github.com/qhato/ecommerce/internal/email/application/queries"
	"github.com/qhato/ecommerce/internal/email/domain"
	"github.com/qhato/ecommerce/pkg/events"
)

// EmailService provides the main application service for emails
type EmailService struct {
	commandHandler *commands.EmailCommandHandler
	queryService   *queries.EmailQueryService
}

// NewEmailService creates a new email service
func NewEmailService(
	emailRepo domain.EmailRepository,
	emailSender commands.EmailSender,
	eventBus events.EventBus,
	templateRenderer commands.TemplateRenderer,
) *EmailService {
	return &EmailService{
		commandHandler: commands.NewEmailCommandHandler(emailRepo, emailSender, eventBus, templateRenderer),
		queryService:   queries.NewEmailQueryService(emailRepo),
	}
}

// Commands

// SendEmail sends an email
func (s *EmailService) SendEmail(ctx context.Context, cmd *commands.SendEmailCommand) (int64, error) {
	return s.commandHandler.HandleSendEmail(ctx, cmd)
}

// ScheduleEmail schedules an email for future sending
func (s *EmailService) ScheduleEmail(ctx context.Context, cmd *commands.ScheduleEmailCommand) (int64, error) {
	return s.commandHandler.HandleScheduleEmail(ctx, cmd)
}

// SendOrderConfirmation sends an order confirmation email
func (s *EmailService) SendOrderConfirmation(ctx context.Context, cmd *commands.SendOrderConfirmationCommand) (int64, error) {
	return s.commandHandler.HandleSendOrderConfirmation(ctx, cmd)
}

// SendOrderShipped sends an order shipped email
func (s *EmailService) SendOrderShipped(ctx context.Context, cmd *commands.SendOrderShippedCommand) (int64, error) {
	return s.commandHandler.HandleSendOrderShipped(ctx, cmd)
}

// SendPasswordReset sends a password reset email
func (s *EmailService) SendPasswordReset(ctx context.Context, cmd *commands.SendPasswordResetCommand) (int64, error) {
	return s.commandHandler.HandleSendPasswordReset(ctx, cmd)
}

// SendWelcomeEmail sends a welcome email
func (s *EmailService) SendWelcomeEmail(ctx context.Context, cmd *commands.SendWelcomeEmailCommand) (int64, error) {
	return s.commandHandler.HandleSendWelcomeEmail(ctx, cmd)
}

// SendCartAbandonment sends a cart abandonment email
func (s *EmailService) SendCartAbandonment(ctx context.Context, cmd *commands.SendCartAbandonmentCommand) (int64, error) {
	return s.commandHandler.HandleSendCartAbandonment(ctx, cmd)
}

// CancelEmail cancels a pending email
func (s *EmailService) CancelEmail(ctx context.Context, cmd *commands.CancelEmailCommand) error {
	return s.commandHandler.HandleCancelEmail(ctx, cmd)
}

// RetryFailedEmail retries a failed email
func (s *EmailService) RetryFailedEmail(ctx context.Context, cmd *commands.RetryFailedEmailCommand) error {
	return s.commandHandler.HandleRetryFailedEmail(ctx, cmd)
}

// Queries

// GetEmailByID retrieves an email by ID
func (s *EmailService) GetEmailByID(ctx context.Context, id int64) (*queries.EmailDTO, error) {
	return s.queryService.GetEmailByID(ctx, id)
}

// ListEmailsByStatus retrieves emails by status with pagination
func (s *EmailService) ListEmailsByStatus(ctx context.Context, status string, offset, limit int) ([]*queries.EmailDTO, error) {
	return s.queryService.ListEmailsByStatus(ctx, status, offset, limit)
}

// ListEmailsByType retrieves emails by type with pagination
func (s *EmailService) ListEmailsByType(ctx context.Context, emailType string, offset, limit int) ([]*queries.EmailDTO, error) {
	return s.queryService.ListEmailsByType(ctx, emailType, offset, limit)
}

// ListEmailsByOrderID retrieves emails associated with an order
func (s *EmailService) ListEmailsByOrderID(ctx context.Context, orderID int64) ([]*queries.EmailDTO, error) {
	return s.queryService.ListEmailsByOrderID(ctx, orderID)
}

// ListEmailsByCustomerID retrieves emails associated with a customer
func (s *EmailService) ListEmailsByCustomerID(ctx context.Context, customerID int64, offset, limit int) ([]*queries.EmailDTO, error) {
	return s.queryService.ListEmailsByCustomerID(ctx, customerID, offset, limit)
}

// GetEmailStats retrieves email statistics
func (s *EmailService) GetEmailStats(ctx context.Context) (*queries.EmailStatsDTO, error) {
	return s.queryService.GetEmailStats(ctx)
}
