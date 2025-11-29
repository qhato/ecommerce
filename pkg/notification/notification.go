package notification

import (
	"context"
	"fmt"
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeEmail NotificationType = "EMAIL"
	NotificationTypeSMS   NotificationType = "SMS"
	NotificationTypePush  NotificationType = "PUSH"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "PENDING"
	NotificationStatusSent      NotificationStatus = "SENT"
	NotificationStatusFailed    NotificationStatus = "FAILED"
	NotificationStatusDelivered NotificationStatus = "DELIVERED"
)

// Notification represents a notification to be sent
type Notification struct {
	ID          string
	Type        NotificationType
	Recipient   string
	Subject     string
	Body        string
	TemplateID  *string
	TemplateData map[string]interface{}
	Status      NotificationStatus
	Error       *string
	SentAt      *time.Time
	DeliveredAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NotificationSender defines the interface for sending notifications
type NotificationSender interface {
	Send(ctx context.Context, notification *Notification) error
	GetType() NotificationType
}

// NotificationService manages sending notifications
type NotificationService struct {
	senders map[NotificationType]NotificationSender
}

// NewNotificationService creates a new notification service
func NewNotificationService() *NotificationService {
	return &NotificationService{
		senders: make(map[NotificationType]NotificationSender),
	}
}

// RegisterSender registers a sender for a notification type
func (s *NotificationService) RegisterSender(sender NotificationSender) {
	s.senders[sender.GetType()] = sender
}

// Send sends a notification
func (s *NotificationService) Send(ctx context.Context, notification *Notification) error {
	sender, exists := s.senders[notification.Type]
	if !exists {
		return fmt.Errorf("no sender registered for notification type: %s", notification.Type)
	}

	notification.Status = NotificationStatusPending

	err := sender.Send(ctx, notification)
	if err != nil {
		notification.Status = NotificationStatusFailed
		errStr := err.Error()
		notification.Error = &errStr
		return err
	}

	notification.Status = NotificationStatusSent
	now := time.Now()
	notification.SentAt = &now

	return nil
}

// SendEmail sends an email notification
func (s *NotificationService) SendEmail(ctx context.Context, to, subject, body string) error {
	notification := &Notification{
		Type:      NotificationTypeEmail,
		Recipient: to,
		Subject:   subject,
		Body:      body,
		CreatedAt: time.Now(),
	}

	return s.Send(ctx, notification)
}

// SendSMS sends an SMS notification
func (s *NotificationService) SendSMS(ctx context.Context, to, body string) error {
	notification := &Notification{
		Type:      NotificationTypeSMS,
		Recipient: to,
		Body:      body,
		CreatedAt: time.Now(),
	}

	return s.Send(ctx, notification)
}

// SendFromTemplate sends a notification using a template
func (s *NotificationService) SendFromTemplate(
	ctx context.Context,
	notifType NotificationType,
	recipient, templateID string,
	data map[string]interface{},
) error {
	notification := &Notification{
		Type:         notifType,
		Recipient:    recipient,
		TemplateID:   &templateID,
		TemplateData: data,
		CreatedAt:    time.Now(),
	}

	return s.Send(ctx, notification)
}

// Email Sender Implementation (SMTP)
type EmailSender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// NewEmailSender creates a new email sender
func NewEmailSender(host string, port int, username, password, from string) *EmailSender {
	return &EmailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *EmailSender) GetType() NotificationType {
	return NotificationTypeEmail
}

func (s *EmailSender) Send(ctx context.Context, notification *Notification) error {
	// TODO: Implement actual SMTP sending
	// For now, just log
	fmt.Printf("Sending email to %s: %s - %s\n", notification.Recipient, notification.Subject, notification.Body)
	return nil
}

// SMS Sender Implementation (Twilio/etc)
type SMSSender struct {
	accountSID string
	authToken  string
	fromNumber string
}

// NewSMSSender creates a new SMS sender
func NewSMSSender(accountSID, authToken, fromNumber string) *SMSSender {
	return &SMSSender{
		accountSID: accountSID,
		authToken:  authToken,
		fromNumber: fromNumber,
	}
}

func (s *SMSSender) GetType() NotificationType {
	return NotificationTypeSMS
}

func (s *SMSSender) Send(ctx context.Context, notification *Notification) error {
	// TODO: Implement actual SMS sending (Twilio API)
	// For now, just log
	fmt.Printf("Sending SMS to %s: %s\n", notification.Recipient, notification.Body)
	return nil
}

// Common notification templates
const (
	TemplateOrderConfirmation   = "order_confirmation"
	TemplateOrderShipped        = "order_shipped"
	TemplateOrderDelivered      = "order_delivered"
	TemplatePasswordReset       = "password_reset"
	TemplateWelcome             = "welcome"
	TemplatePaymentConfirmation = "payment_confirmation"
)
