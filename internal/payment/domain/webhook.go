package domain

import (
	"time"

	"github.com/google/uuid"
)

// WebhookEvent represents a webhook event from a payment gateway
type WebhookEvent struct {
	ID          string
	GatewayName string
	EventType   WebhookEventType
	EventID     string // Gateway's event ID
	Payload     string // Raw JSON payload
	Status      WebhookStatus
	ProcessedAt *time.Time
	ErrorMsg    *string
	Signature   *string // Webhook signature for verification
	IPAddress   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// WebhookEventType represents the type of webhook event
type WebhookEventType string

const (
	WebhookEventPaymentSucceeded     WebhookEventType = "PAYMENT_SUCCEEDED"
	WebhookEventPaymentFailed        WebhookEventType = "PAYMENT_FAILED"
	WebhookEventPaymentRefunded      WebhookEventType = "PAYMENT_REFUNDED"
	WebhookEventPaymentCancelled     WebhookEventType = "PAYMENT_CANCELLED"
	WebhookEventChargebackCreated    WebhookEventType = "CHARGEBACK_CREATED"
	WebhookEventChargebackResolved   WebhookEventType = "CHARGEBACK_RESOLVED"
	WebhookEventSubscriptionCreated  WebhookEventType = "SUBSCRIPTION_CREATED"
	WebhookEventSubscriptionCancelled WebhookEventType = "SUBSCRIPTION_CANCELLED"
	WebhookEventDisputeCreated       WebhookEventType = "DISPUTE_CREATED"
	WebhookEventUnknown              WebhookEventType = "UNKNOWN"
)

// WebhookStatus represents the processing status of a webhook
type WebhookStatus string

const (
	WebhookStatusPending   WebhookStatus = "PENDING"
	WebhookStatusProcessed WebhookStatus = "PROCESSED"
	WebhookStatusFailed    WebhookStatus = "FAILED"
	WebhookStatusIgnored   WebhookStatus = "IGNORED" // Event type not handled
)

// NewWebhookEvent creates a new webhook event
func NewWebhookEvent(gatewayName, eventID, eventType, payload string) *WebhookEvent {
	now := time.Now()
	return &WebhookEvent{
		ID:          uuid.New().String(),
		GatewayName: gatewayName,
		EventID:     eventID,
		EventType:   WebhookEventType(eventType),
		Payload:     payload,
		Status:      WebhookStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MarkAsProcessed marks the webhook as successfully processed
func (w *WebhookEvent) MarkAsProcessed() {
	now := time.Now()
	w.Status = WebhookStatusProcessed
	w.ProcessedAt = &now
	w.UpdatedAt = now
}

// MarkAsFailed marks the webhook as failed
func (w *WebhookEvent) MarkAsFailed(errorMsg string) {
	w.Status = WebhookStatusFailed
	w.ErrorMsg = &errorMsg
	w.UpdatedAt = time.Now()
}

// MarkAsIgnored marks the webhook as ignored
func (w *WebhookEvent) MarkAsIgnored() {
	w.Status = WebhookStatusIgnored
	w.UpdatedAt = time.Now()
}

// SetSignature sets the webhook signature
func (w *WebhookEvent) SetSignature(signature string) {
	w.Signature = &signature
}

// SetIPAddress sets the source IP address
func (w *WebhookEvent) SetIPAddress(ip string) {
	w.IPAddress = &ip
}
