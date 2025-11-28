package commands

import (
	"context"
	"time"

	"github.com/qhato/ecommerce/internal/payment/domain"
	"github.com/qhato/ecommerce/pkg/apperrors"
	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
)

// PaymentCommandHandler handles payment commands
type PaymentCommandHandler struct {
	repo     domain.PaymentRepository
	eventBus event.EventBus
	log      *logger.Logger
}

// NewPaymentCommandHandler creates a new PaymentCommandHandler
func NewPaymentCommandHandler(repo domain.PaymentRepository, eventBus event.EventBus, log *logger.Logger) *PaymentCommandHandler {
	return &PaymentCommandHandler{
		repo:     repo,
		eventBus: eventBus,
		log:      log,
	}
}

// CreatePayment creates a new payment
func (h *PaymentCommandHandler) CreatePayment(ctx context.Context, orderID, customerID int64, paymentMethod domain.PaymentMethod, amount float64, currencyCode string) (*domain.Payment, error) {
	h.log.Info("Creating new payment", "orderID", orderID, "amount", amount)

	// Validate amount
	if amount <= 0 {
		return nil, apperrors.NewValidationError("payment amount must be greater than zero")
	}

	// Create payment
	payment := domain.NewPayment(orderID, customerID, paymentMethod, amount, currencyCode)

	// Save payment
	if err := h.repo.Create(ctx, payment); err != nil {
		h.log.Error("Failed to create payment", "error", err)
		return nil, apperrors.NewInternalError("failed to create payment", err)
	}

	// Publish event
	evt := domain.NewPaymentCreatedEvent(payment.ID, payment.OrderID, payment.CustomerID, payment.PaymentMethod, payment.Amount, payment.CurrencyCode)
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish payment created event", "error", err)
	}

	h.log.Info("Payment created successfully", "paymentID", payment.ID)
	return payment, nil
}

// AuthorizePayment authorizes a payment
func (h *PaymentCommandHandler) AuthorizePayment(ctx context.Context, paymentID int64, authorizationCode, transactionID string) error {
	h.log.Info("Authorizing payment", "paymentID", paymentID)

	// Find payment
	payment, err := h.repo.FindByID(ctx, paymentID)
	if err != nil {
		h.log.Error("Failed to find payment", "error", err)
		return err
	}
	if payment == nil {
		return apperrors.NewNotFoundError("payment", paymentID)
	}

	// Authorize payment
	payment.Authorize(authorizationCode, transactionID)

	// Save payment
	if err := h.repo.Update(ctx, payment); err != nil {
		h.log.Error("Failed to authorize payment", "error", err)
		return apperrors.NewInternalError("failed to authorize payment", err)
	}

	// Publish event
	evt := &domain.PaymentAuthorizedEvent{
		BaseEvent:         event.BaseEvent{EventType: domain.EventPaymentAuthorized, Timestamp: time.Now()},
		PaymentID:         payment.ID,
		OrderID:           payment.OrderID,
		TransactionID:     transactionID,
		AuthorizationCode: authorizationCode,
		Amount:            payment.Amount,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish payment authorized event", "error", err)
	}

	h.log.Info("Payment authorized successfully", "paymentID", paymentID)
	return nil
}

// CapturePayment captures an authorized payment
func (h *PaymentCommandHandler) CapturePayment(ctx context.Context, paymentID int64, transactionID string) error {
	h.log.Info("Capturing payment", "paymentID", paymentID)

	// Find payment
	payment, err := h.repo.FindByID(ctx, paymentID)
	if err != nil {
		h.log.Error("Failed to find payment", "error", err)
		return err
	}
	if payment == nil {
		return apperrors.NewNotFoundError("payment", paymentID)
	}

	// Capture payment
	if err := payment.Capture(transactionID); err != nil {
		return apperrors.NewValidationError(err.Error())
	}

	// Save payment
	if err := h.repo.Update(ctx, payment); err != nil {
		h.log.Error("Failed to capture payment", "error", err)
		return apperrors.NewInternalError("failed to capture payment", err)
	}

	// Publish event
	evt := &domain.PaymentCapturedEvent{
		BaseEvent:     event.BaseEvent{EventType: domain.EventPaymentCaptured, Timestamp: time.Now()},
		PaymentID:     payment.ID,
		OrderID:       payment.OrderID,
		TransactionID: payment.TransactionID,
		Amount:        payment.Amount,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish payment captured event", "error", err)
	}

	h.log.Info("Payment captured successfully", "paymentID", paymentID)
	return nil
}

// CompletePayment completes a payment
func (h *PaymentCommandHandler) CompletePayment(ctx context.Context, paymentID int64, transactionID string) error {
	h.log.Info("Completing payment", "paymentID", paymentID)

	// Find payment
	payment, err := h.repo.FindByID(ctx, paymentID)
	if err != nil {
		h.log.Error("Failed to find payment", "error", err)
		return err
	}
	if payment == nil {
		return apperrors.NewNotFoundError("payment", paymentID)
	}

	// Complete payment
	payment.Complete(transactionID)

	// Save payment
	if err := h.repo.Update(ctx, payment); err != nil {
		h.log.Error("Failed to complete payment", "error", err)
		return apperrors.NewInternalError("failed to complete payment", err)
	}

	// Publish event
	evt := &domain.PaymentCompletedEvent{
		BaseEvent:     event.BaseEvent{EventType: domain.EventPaymentCompleted, Timestamp: time.Now()},
		PaymentID:     payment.ID,
		OrderID:       payment.OrderID,
		TransactionID: payment.TransactionID,
		Amount:        payment.Amount,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish payment completed event", "error", err)
	}

	h.log.Info("Payment completed successfully", "paymentID", paymentID)
	return nil
}

// FailPayment marks a payment as failed
func (h *PaymentCommandHandler) FailPayment(ctx context.Context, paymentID int64, reason string) error {
	h.log.Info("Failing payment", "paymentID", paymentID, "reason", reason)

	// Find payment
	payment, err := h.repo.FindByID(ctx, paymentID)
	if err != nil {
		h.log.Error("Failed to find payment", "error", err)
		return err
	}
	if payment == nil {
		return apperrors.NewNotFoundError("payment", paymentID)
	}

	// Fail payment
	payment.Fail(reason)

	// Save payment
	if err := h.repo.Update(ctx, payment); err != nil {
		h.log.Error("Failed to update payment status", "error", err)
		return apperrors.NewInternalError("failed to update payment status", err)
	}

	// Publish event
	evt := &domain.PaymentFailedEvent{
		BaseEvent:     event.BaseEvent{EventType: domain.EventPaymentFailed, Timestamp: time.Now()},
		PaymentID:     payment.ID,
		OrderID:       payment.OrderID,
		FailureReason: reason,
		Amount:        payment.Amount,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish payment failed event", "error", err)
	}

	h.log.Info("Payment marked as failed", "paymentID", paymentID)
	return nil
}

// RefundPayment refunds a payment
func (h *PaymentCommandHandler) RefundPayment(ctx context.Context, paymentID int64, amount float64) error {
	h.log.Info("Refunding payment", "paymentID", paymentID, "amount", amount)

	// Find payment
	payment, err := h.repo.FindByID(ctx, paymentID)
	if err != nil {
		h.log.Error("Failed to find payment", "error", err)
		return err
	}
	if payment == nil {
		return apperrors.NewNotFoundError("payment", paymentID)
	}

	// Validate refund
	if !payment.IsRefundable() {
		return apperrors.NewValidationError("payment is not refundable")
	}

	// Refund payment
	if err := payment.Refund(amount); err != nil {
		return apperrors.NewValidationError(err.Error())
	}

	// Save payment
	if err := h.repo.Update(ctx, payment); err != nil {
		h.log.Error("Failed to refund payment", "error", err)
		return apperrors.NewInternalError("failed to refund payment", err)
	}

	// Publish event
	evt := &domain.PaymentRefundedEvent{
		BaseEvent:     event.BaseEvent{EventType: domain.EventPaymentRefunded, Timestamp: time.Now()},
		PaymentID:     payment.ID,
		OrderID:       payment.OrderID,
		RefundAmount:  amount,
		TotalRefunded: payment.RefundAmount,
	}
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.Error("Failed to publish payment refunded event", "error", err)
	}

	h.log.Info("Payment refunded successfully", "paymentID", paymentID, "amount", amount)
	return nil
}

// CancelPayment cancels a payment
func (h *PaymentCommandHandler) CancelPayment(ctx context.Context, paymentID int64) error {
	h.log.Info("Cancelling payment", "paymentID", paymentID)

	// Find payment
	payment, err := h.repo.FindByID(ctx, paymentID)
	if err != nil {
		h.log.Error("Failed to find payment", "error", err)
		return err
	}
	if payment == nil {
		return apperrors.NewNotFoundError("payment", paymentID)
	}

	// Validate cancellation
	if !payment.IsCancellable() {
		return apperrors.NewValidationError("payment cannot be cancelled in current status")
	}

	// Cancel payment
	if err := payment.Cancel(); err != nil {
		return apperrors.NewValidationError(err.Error())
	}

	// Save payment
	if err := h.repo.Update(ctx, payment); err != nil {
		h.log.Error("Failed to cancel payment", "error", err)
		return apperrors.NewInternalError("failed to cancel payment", err)
	}

	h.log.Info("Payment cancelled successfully", "paymentID", paymentID)
	return nil
}
