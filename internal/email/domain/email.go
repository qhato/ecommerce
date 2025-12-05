package domain

import (
	"fmt"
	"time"
)

// EmailStatus represents the status of an email
type EmailStatus string

const (
	EmailStatusPending   EmailStatus = "PENDING"
	EmailStatusQueued    EmailStatus = "QUEUED"
	EmailStatusSending   EmailStatus = "SENDING"
	EmailStatusSent      EmailStatus = "SENT"
	EmailStatusFailed    EmailStatus = "FAILED"
	EmailStatusRetrying  EmailStatus = "RETRYING"
	EmailStatusCancelled EmailStatus = "CANCELLED"
)

// EmailType represents the type of email
type EmailType string

const (
	EmailTypeOrderConfirmation      EmailType = "ORDER_CONFIRMATION"
	EmailTypeOrderShipped           EmailType = "ORDER_SHIPPED"
	EmailTypeOrderDelivered         EmailType = "ORDER_DELIVERED"
	EmailTypeOrderCancelled         EmailType = "ORDER_CANCELLED"
	EmailTypePasswordReset          EmailType = "PASSWORD_RESET"
	EmailTypeWelcome                EmailType = "WELCOME"
	EmailTypeCartAbandonment        EmailType = "CART_ABANDONMENT"
	EmailTypeProductBackInStock     EmailType = "PRODUCT_BACK_IN_STOCK"
	EmailTypePromotionalNewsletter  EmailType = "PROMOTIONAL_NEWSLETTER"
	EmailTypeTransactional          EmailType = "TRANSACTIONAL"
)

// EmailPriority represents the priority of an email
type EmailPriority int

const (
	EmailPriorityLow    EmailPriority = 1
	EmailPriorityNormal EmailPriority = 5
	EmailPriorityHigh   EmailPriority = 10
	EmailPriorityUrgent EmailPriority = 20
)

// Email represents an email entity
type Email struct {
	ID            int64
	Type          EmailType
	Status        EmailStatus
	Priority      EmailPriority
	From          string
	To            []string
	CC            []string
	BCC           []string
	ReplyTo       string
	Subject       string
	Body          string
	HTMLBody      string
	TemplateName  string
	TemplateData  map[string]interface{}
	Attachments   []EmailAttachment
	Headers       map[string]string
	MaxRetries    int
	RetryCount    int
	ScheduledAt   *time.Time
	SentAt        *time.Time
	FailedAt      *time.Time
	ErrorMessage  string
	OrderID       *int64
	CustomerID    *int64
	CreatedBy     int64
	UpdatedBy     int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	ID          int64
	EmailID     int64
	Filename    string
	ContentType string
	Content     []byte
	Size        int64
	CreatedAt   time.Time
}

// NewEmail creates a new email
func NewEmail(emailType EmailType, from string, to []string, subject string) *Email {
	now := time.Now()
	return &Email{
		Type:         emailType,
		Status:       EmailStatusPending,
		Priority:     EmailPriorityNormal,
		From:         from,
		To:           to,
		CC:           make([]string, 0),
		BCC:          make([]string, 0),
		Subject:      subject,
		TemplateData: make(map[string]interface{}),
		Attachments:  make([]EmailAttachment, 0),
		Headers:      make(map[string]string),
		MaxRetries:   3,
		RetryCount:   0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewTransactionalEmail creates a new transactional email
func NewTransactionalEmail(from string, to []string, subject string) *Email {
	email := NewEmail(EmailTypeTransactional, from, to, subject)
	email.Priority = EmailPriorityHigh
	return email
}

// SetTemplate sets the template for the email
func (e *Email) SetTemplate(templateName string, data map[string]interface{}) {
	e.TemplateName = templateName
	e.TemplateData = data
	e.UpdatedAt = time.Now()
}

// SetBody sets the plain text body
func (e *Email) SetBody(body string) {
	e.Body = body
	e.UpdatedAt = time.Now()
}

// SetHTMLBody sets the HTML body
func (e *Email) SetHTMLBody(htmlBody string) {
	e.HTMLBody = htmlBody
	e.UpdatedAt = time.Now()
}

// AddCC adds a CC recipient
func (e *Email) AddCC(email string) {
	e.CC = append(e.CC, email)
	e.UpdatedAt = time.Now()
}

// AddBCC adds a BCC recipient
func (e *Email) AddBCC(email string) {
	e.BCC = append(e.BCC, email)
	e.UpdatedAt = time.Now()
}

// SetReplyTo sets the reply-to address
func (e *Email) SetReplyTo(replyTo string) {
	e.ReplyTo = replyTo
	e.UpdatedAt = time.Now()
}

// AddAttachment adds an attachment to the email
func (e *Email) AddAttachment(filename, contentType string, content []byte) {
	attachment := EmailAttachment{
		Filename:    filename,
		ContentType: contentType,
		Content:     content,
		Size:        int64(len(content)),
		CreatedAt:   time.Now(),
	}
	e.Attachments = append(e.Attachments, attachment)
	e.UpdatedAt = time.Now()
}

// AddHeader adds a custom header
func (e *Email) AddHeader(key, value string) {
	e.Headers[key] = value
	e.UpdatedAt = time.Now()
}

// SetPriority sets the email priority
func (e *Email) SetPriority(priority EmailPriority) {
	e.Priority = priority
	e.UpdatedAt = time.Now()
}

// Schedule schedules the email for future sending
func (e *Email) Schedule(scheduledAt time.Time) error {
	if scheduledAt.Before(time.Now()) {
		return fmt.Errorf("scheduled time must be in the future")
	}
	e.ScheduledAt = &scheduledAt
	e.Status = EmailStatusPending
	e.UpdatedAt = time.Now()
	return nil
}

// MarkAsQueued marks the email as queued
func (e *Email) MarkAsQueued() {
	e.Status = EmailStatusQueued
	e.UpdatedAt = time.Now()
}

// MarkAsSending marks the email as being sent
func (e *Email) MarkAsSending() {
	e.Status = EmailStatusSending
	e.UpdatedAt = time.Now()
}

// MarkAsSent marks the email as sent
func (e *Email) MarkAsSent() {
	e.Status = EmailStatusSent
	now := time.Now()
	e.SentAt = &now
	e.UpdatedAt = now
}

// MarkAsFailed marks the email as failed
func (e *Email) MarkAsFailed(errorMessage string) {
	e.Status = EmailStatusFailed
	e.ErrorMessage = errorMessage
	now := time.Now()
	e.FailedAt = &now
	e.UpdatedAt = now
}

// MarkAsRetrying marks the email as retrying
func (e *Email) MarkAsRetrying() {
	e.Status = EmailStatusRetrying
	e.RetryCount++
	e.UpdatedAt = time.Now()
}

// Cancel cancels the email
func (e *Email) Cancel() error {
	if e.Status == EmailStatusSent {
		return fmt.Errorf("cannot cancel a sent email")
	}
	e.Status = EmailStatusCancelled
	e.UpdatedAt = time.Now()
	return nil
}

// CanRetry checks if the email can be retried
func (e *Email) CanRetry() bool {
	return e.Status == EmailStatusFailed && e.RetryCount < e.MaxRetries
}

// IsReadyToSend checks if the email is ready to be sent
func (e *Email) IsReadyToSend() bool {
	if e.Status != EmailStatusPending && e.Status != EmailStatusQueued {
		return false
	}

	// Check if scheduled
	if e.ScheduledAt != nil && e.ScheduledAt.After(time.Now()) {
		return false
	}

	return true
}

// Validate validates the email
func (e *Email) Validate() error {
	if e.From == "" {
		return fmt.Errorf("from address is required")
	}
	if len(e.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}
	if e.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	if e.Body == "" && e.HTMLBody == "" && e.TemplateName == "" {
		return fmt.Errorf("body, HTML body, or template is required")
	}
	return nil
}

// AssociateWithOrder associates the email with an order
func (e *Email) AssociateWithOrder(orderID int64) {
	e.OrderID = &orderID
	e.UpdatedAt = time.Now()
}

// AssociateWithCustomer associates the email with a customer
func (e *Email) AssociateWithCustomer(customerID int64) {
	e.CustomerID = &customerID
	e.UpdatedAt = time.Now()
}
