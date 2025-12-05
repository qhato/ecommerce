package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/email/domain"
	"github.com/qhato/ecommerce/pkg/events"
)

// EmailCommandHandler handles email commands
type EmailCommandHandler struct {
	emailRepo    domain.EmailRepository
	emailSender  EmailSender
	eventBus     events.EventBus
	templateRenderer TemplateRenderer
}

// EmailSender defines the interface for sending emails
type EmailSender interface {
	Send(ctx context.Context, email *domain.Email) error
	Queue(ctx context.Context, email *domain.Email) error
}

// TemplateRenderer defines the interface for rendering email templates
type TemplateRenderer interface {
	Render(templateName string, data map[string]interface{}) (string, string, error)
}

// NewEmailCommandHandler creates a new email command handler
func NewEmailCommandHandler(
	emailRepo domain.EmailRepository,
	emailSender EmailSender,
	eventBus events.EventBus,
	templateRenderer TemplateRenderer,
) *EmailCommandHandler {
	return &EmailCommandHandler{
		emailRepo:        emailRepo,
		emailSender:      emailSender,
		eventBus:         eventBus,
		templateRenderer: templateRenderer,
	}
}

// HandleSendEmail handles sending an email
func (h *EmailCommandHandler) HandleSendEmail(ctx context.Context, cmd *SendEmailCommand) (int64, error) {
	email := domain.NewEmail(
		domain.EmailType(cmd.Type),
		cmd.From,
		cmd.To,
		cmd.Subject,
	)

	if len(cmd.CC) > 0 {
		for _, cc := range cmd.CC {
			email.AddCC(cc)
		}
	}

	if len(cmd.BCC) > 0 {
		for _, bcc := range cmd.BCC {
			email.AddBCC(bcc)
		}
	}

	if cmd.ReplyTo != "" {
		email.SetReplyTo(cmd.ReplyTo)
	}

	if cmd.Body != "" {
		email.SetBody(cmd.Body)
	}

	if cmd.HTMLBody != "" {
		email.SetHTMLBody(cmd.HTMLBody)
	}

	if cmd.TemplateName != "" {
		// Render template
		plainBody, htmlBody, err := h.templateRenderer.Render(cmd.TemplateName, cmd.TemplateData)
		if err != nil {
			return 0, fmt.Errorf("failed to render template: %w", err)
		}
		email.SetBody(plainBody)
		email.SetHTMLBody(htmlBody)
		email.SetTemplate(cmd.TemplateName, cmd.TemplateData)
	}

	if cmd.Priority > 0 {
		email.SetPriority(domain.EmailPriority(cmd.Priority))
	}

	if cmd.OrderID != nil {
		email.AssociateWithOrder(*cmd.OrderID)
	}

	if cmd.CustomerID != nil {
		email.AssociateWithCustomer(*cmd.CustomerID)
	}

	// Add attachments
	for _, att := range cmd.Attachments {
		email.AddAttachment(att.Filename, att.ContentType, att.Content)
	}

	// Validate email
	if err := email.Validate(); err != nil {
		return 0, fmt.Errorf("invalid email: %w", err)
	}

	// Save email
	if err := h.emailRepo.Create(ctx, email); err != nil {
		return 0, fmt.Errorf("failed to create email: %w", err)
	}

	// Queue email for sending
	if err := h.emailSender.Queue(ctx, email); err != nil {
		email.MarkAsFailed(err.Error())
		h.emailRepo.Update(ctx, email)
		return email.ID, fmt.Errorf("failed to queue email: %w", err)
	}

	email.MarkAsQueued()
	if err := h.emailRepo.Update(ctx, email); err != nil {
		return email.ID, fmt.Errorf("failed to update email status: %w", err)
	}

	// Publish event
	event := domain.NewEmailQueuedEvent(email)
	h.eventBus.Publish(ctx, "email.queued", event)

	return email.ID, nil
}

// HandleScheduleEmail handles scheduling an email
func (h *EmailCommandHandler) HandleScheduleEmail(ctx context.Context, cmd *ScheduleEmailCommand) (int64, error) {
	email := domain.NewEmail(
		domain.EmailType(cmd.Type),
		cmd.From,
		cmd.To,
		cmd.Subject,
	)

	if cmd.Body != "" {
		email.SetBody(cmd.Body)
	}

	if cmd.HTMLBody != "" {
		email.SetHTMLBody(cmd.HTMLBody)
	}

	if cmd.TemplateName != "" {
		email.SetTemplate(cmd.TemplateName, cmd.TemplateData)
	}

	if cmd.Priority > 0 {
		email.SetPriority(domain.EmailPriority(cmd.Priority))
	}

	// Schedule email
	if err := email.Schedule(cmd.ScheduledAt); err != nil {
		return 0, err
	}

	// Save email
	if err := h.emailRepo.Create(ctx, email); err != nil {
		return 0, fmt.Errorf("failed to create scheduled email: %w", err)
	}

	// Publish event
	event := domain.NewEmailScheduledEvent(email)
	h.eventBus.Publish(ctx, "email.scheduled", event)

	return email.ID, nil
}

// HandleSendOrderConfirmation handles sending order confirmation email
func (h *EmailCommandHandler) HandleSendOrderConfirmation(ctx context.Context, cmd *SendOrderConfirmationCommand) (int64, error) {
	templateData := map[string]interface{}{
		"OrderNumber":    cmd.OrderNumber,
		"OrderTotal":     cmd.OrderTotal,
		"OrderDate":      cmd.OrderDate,
		"Items":          cmd.Items,
		"ShippingAddr":   cmd.ShippingAddr,
		"BillingAddr":    cmd.BillingAddr,
	}

	sendCmd := &SendEmailCommand{
		Type:         string(domain.EmailTypeOrderConfirmation),
		To:           []string{cmd.To},
		Subject:      fmt.Sprintf("Order Confirmation - Order #%s", cmd.OrderNumber),
		TemplateName: "order_confirmation",
		TemplateData: templateData,
		Priority:     int(domain.EmailPriorityHigh),
		OrderID:      &cmd.OrderID,
		CustomerID:   &cmd.CustomerID,
	}

	return h.HandleSendEmail(ctx, sendCmd)
}

// HandleSendOrderShipped handles sending order shipped email
func (h *EmailCommandHandler) HandleSendOrderShipped(ctx context.Context, cmd *SendOrderShippedCommand) (int64, error) {
	templateData := map[string]interface{}{
		"OrderNumber":           cmd.OrderNumber,
		"TrackingNumber":        cmd.TrackingNumber,
		"Carrier":               cmd.Carrier,
		"ShippedAt":             cmd.ShippedAt,
		"EstimatedDeliveryDate": cmd.EstimatedDeliveryDate,
	}

	sendCmd := &SendEmailCommand{
		Type:         string(domain.EmailTypeOrderShipped),
		To:           []string{cmd.To},
		Subject:      fmt.Sprintf("Your Order Has Shipped - Order #%s", cmd.OrderNumber),
		TemplateName: "order_shipped",
		TemplateData: templateData,
		Priority:     int(domain.EmailPriorityHigh),
		OrderID:      &cmd.OrderID,
		CustomerID:   &cmd.CustomerID,
	}

	return h.HandleSendEmail(ctx, sendCmd)
}

// HandleSendPasswordReset handles sending password reset email
func (h *EmailCommandHandler) HandleSendPasswordReset(ctx context.Context, cmd *SendPasswordResetCommand) (int64, error) {
	templateData := map[string]interface{}{
		"ResetToken": cmd.ResetToken,
		"ExpiresAt":  cmd.ExpiresAt,
	}

	sendCmd := &SendEmailCommand{
		Type:         string(domain.EmailTypePasswordReset),
		To:           []string{cmd.To},
		Subject:      "Password Reset Request",
		TemplateName: "password_reset",
		TemplateData: templateData,
		Priority:     int(domain.EmailPriorityUrgent),
		CustomerID:   &cmd.CustomerID,
	}

	return h.HandleSendEmail(ctx, sendCmd)
}

// HandleSendWelcomeEmail handles sending welcome email
func (h *EmailCommandHandler) HandleSendWelcomeEmail(ctx context.Context, cmd *SendWelcomeEmailCommand) (int64, error) {
	templateData := map[string]interface{}{
		"FirstName": cmd.FirstName,
		"LastName":  cmd.LastName,
	}

	sendCmd := &SendEmailCommand{
		Type:         string(domain.EmailTypeWelcome),
		To:           []string{cmd.To},
		Subject:      "Welcome to Our Store!",
		TemplateName: "welcome",
		TemplateData: templateData,
		Priority:     int(domain.EmailPriorityNormal),
		CustomerID:   &cmd.CustomerID,
	}

	return h.HandleSendEmail(ctx, sendCmd)
}

// HandleSendCartAbandonment handles sending cart abandonment email
func (h *EmailCommandHandler) HandleSendCartAbandonment(ctx context.Context, cmd *SendCartAbandonmentCommand) (int64, error) {
	templateData := map[string]interface{}{
		"CartID":      cmd.CartID,
		"Items":       cmd.Items,
		"Total":       cmd.Total,
		"AbandonedAt": cmd.AbandonedAt,
	}

	sendCmd := &SendEmailCommand{
		Type:         string(domain.EmailTypeCartAbandonment),
		To:           []string{cmd.To},
		Subject:      "You Left Items in Your Cart",
		TemplateName: "cart_abandonment",
		TemplateData: templateData,
		Priority:     int(domain.EmailPriorityNormal),
		CustomerID:   &cmd.CustomerID,
	}

	return h.HandleSendEmail(ctx, sendCmd)
}

// HandleCancelEmail handles cancelling an email
func (h *EmailCommandHandler) HandleCancelEmail(ctx context.Context, cmd *CancelEmailCommand) error {
	email, err := h.emailRepo.FindByID(ctx, cmd.EmailID)
	if err != nil {
		return fmt.Errorf("failed to find email: %w", err)
	}

	if err := email.Cancel(); err != nil {
		return err
	}

	if err := h.emailRepo.Update(ctx, email); err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	// Publish event
	event := domain.NewEmailCancelledEvent(email)
	h.eventBus.Publish(ctx, "email.cancelled", event)

	return nil
}

// HandleRetryFailedEmail handles retrying a failed email
func (h *EmailCommandHandler) HandleRetryFailedEmail(ctx context.Context, cmd *RetryFailedEmailCommand) error {
	email, err := h.emailRepo.FindByID(ctx, cmd.EmailID)
	if err != nil {
		return fmt.Errorf("failed to find email: %w", err)
	}

	if !email.CanRetry() {
		return domain.ErrMaxRetriesExceeded
	}

	email.MarkAsRetrying()

	if err := h.emailRepo.Update(ctx, email); err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	// Queue email for retry
	if err := h.emailSender.Queue(ctx, email); err != nil {
		email.MarkAsFailed(err.Error())
		h.emailRepo.Update(ctx, email)
		return fmt.Errorf("failed to queue email for retry: %w", err)
	}

	return nil
}
