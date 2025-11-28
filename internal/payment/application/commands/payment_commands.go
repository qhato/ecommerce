package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/payment/domain"
	"github.com/qhato/ecommerce/pkg/errors"
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
	h.log.WithFields(logger.Fields{"orderID": orderID, "amount": amount}).Info("Creating new payment")

	// Validate amount
	if amount <= 0 {
		return nil, errors.ValidationError("payment amount must be greater than zero")
	}

	// Create payment
	payment := domain.NewPayment(orderID, customerID, paymentMethod, amount, currencyCode)

	// Save payment
	if err := h.repo.Create(ctx, payment); err != nil {
		h.log.WithError(err).Error("Failed to create payment")
		return nil, errors.InternalWrap(err, "failed to create payment")
	}

	// Publish event
	evt := domain.NewPaymentCreatedEvent(payment.ID, payment.OrderID, payment.CustomerID, payment.PaymentMethod, payment.Amount, payment.CurrencyCode)
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.log.WithError(err).Error("Failed to publish payment created event")
	}

	h.log.WithFields(logger.Fields{"paymentID": payment.ID, "orderID": payment.OrderID, "amount": payment.Amount}).Info("Payment created successfully")
	return payment, nil
}

// AuthorizePayment authorizes a payment
func (h *PaymentCommandHandler) AuthorizePayment(ctx context.Context, paymentID int64, authorizationCode, transactionID string) error {
	h.log.WithField("paymentID", paymentID).Info("Authorizing payment")

	// Find payment
	payment, err := h.repo.FindByID(ctx, paymentID)
	if err != nil {
		h.log.Error("Failed to find payment", "error", err)
		return err
	}
	if payment == nil {
		return errors.NotFound(fmt.Sprintf("payment %d", paymentID))
	}

	// Authorize payment
	payment.Authorize(authorizationCode, transactionID)

	// Save payment
	if err := h.repo.Update(ctx, payment); err != nil {
		h.log.WithError(err).Error("Failed to authorize payment")
		return errors.InternalWrap(err, "failed to authorize payment")
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
		h.log.WithError(err).Error("Failed to publish payment authorized event")
	}

	h.log.WithFields(logger.Fields{"paymentID": paymentID, "orderID": payment.OrderID, "transactionID": transactionID}).Info("Payment authorized successfully")
	return nil
}

// CapturePayment captures an authorized payment
func (h *PaymentCommandHandler) CapturePayment(ctx context.Context, paymentID int64, transactionID string) error {
	h.log.WithFields(logger.Fields{"paymentID": paymentID, "transactionID": transactionID}).Info("Capturing payment")

	// Find payment
	payment, err := h.repo.FindByID(ctx, paymentID)
	if err != nil {
		h.log.Error("Failed to find payment", "error", err)
		return err
	}
	if payment == nil {
		return errors.NotFound(fmt.Sprintf("payment %d", paymentID))
	}

	// Capture payment
	if err := payment.Capture(transactionID); err != nil {
		return errors.ValidationError(err.Error())
	}

	// Save payment
	if err := h.repo.Update(ctx, payment); err != nil {
		h.log.WithError(err).Error("Failed to capture payment")
		return errors.InternalWrap(err, "failed to capture payment")
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
		h.log.WithError(err).Error("Failed to publish payment captured event")
	}

	h.log.WithFields(logger.Fields{"paymentID": paymentID, "orderID": payment.OrderID, "transactionID": transactionID}).Info("Payment captured successfully")
	return nil
}

// CompletePayment completes a payment
func (h *PaymentCommandHandler) CompletePayment(ctx context.Context, paymentID int64, transactionID string) error {
	h.log.WithFields(logger.Fields{"paymentID": paymentID, "transactionID": transactionID}).Info("Completing payment")

	// Find payment
	payment, err := h.repo.FindByID(ctx, paymentID)
	if err != nil {
		h.log.Error("Failed to find payment", "error", err)
		return err
	}
	if payment == nil {
		return errors.NotFound(fmt.Sprintf("payment %d", paymentID))
	}

	// Complete payment
	payment.Complete(transactionID)

	// Save payment
	if err := h.repo.Update(ctx, payment); err != nil {
		h.log.WithError(err).Error("Failed to complete payment")
		return errors.InternalWrap(err, "failed to complete payment")
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
		h.log.WithError(err).Error("Failed to publish payment completed event")
	}

	h.log.WithFields(logger.Fields{"paymentID": paymentID, "orderID": payment.OrderID, "transactionID": transactionID}).Info("Payment completed successfully")
	return nil
}

// FailPayment marks a payment as failed
func (h *PaymentCommandHandler) FailPayment(ctx context.Context, paymentID int64, reason string) error {
	h.log.WithFields(logger.Fields{"paymentID": paymentID, "reason": reason}).Info("Failing payment")

	// Find payment
	payment, err := h.repo.FindByID(ctx, paymentID)
	if err != nil {
		h.log.Error("Failed to find payment", "error", err)
		return err
	}
	if payment == nil {
		return errors.NotFound(fmt.Sprintf("payment %d", paymentID))
	}

	// Fail payment
	payment.Fail(reason)

	// Save payment
	if err := h.repo.Update(ctx, payment); err != nil {
		h.log.WithError(err).Error("Failed to update payment status")
		return errors.InternalWrap(err, "failed to update payment status")
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
		h.log.WithError(err).Error("Failed to publish payment failed event")
	}

	h.log.WithFields(logger.Fields{"paymentID": paymentID, "orderID": payment.OrderID, "reason": reason}).Info("Payment marked as failed")
	return nil
}

// RefundPayment refunds a payment
func (h *PaymentCommandHandler) RefundPayment(ctx context.Context, paymentID int64, amount float64) error {
	h.log.WithFields(logger.Fields{"paymentID": paymentID, "amount": amount}).Info("Refunding payment")

	// Find payment
	payment, err := h.repo.FindByID(ctx, paymentID)
	if err != nil {
		h.log.Error("Failed to find payment", "error", err)
		return err
	}
	if payment == nil {
		return errors.NotFound(fmt.Sprintf("payment %d", paymentID))
	}

	// Validate refund
	if !payment.IsRefundable() {
		return errors.ValidationError("payment is not refundable")
	}

	// Refund payment
	if err := payment.Refund(amount); err != nil {
		return errors.ValidationError(err.Error())
	}

	// Save payment
	if err := h.repo.Update(ctx, payment); err != nil {
		h.log.WithError(err).Error("Failed to refund payment")
		return errors.InternalWrap(err, "failed to refund payment")
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
		h.log.WithError(err).Error("Failed to publish payment refunded event")
	}

	h.log.WithFields(logger.Fields{"paymentID": paymentID, "orderID": payment.OrderID, "refundAmount": amount, "totalRefunded": payment.RefundAmount}).Info("Payment refunded successfully")
	return nil
}

// CancelPayment cancels a payment
func (h *PaymentCommandHandler) CancelPayment(ctx context.Context, paymentID int64) error {
	h.log.WithField("paymentID", paymentID).Info("Cancelling payment")

	// Find payment
	payment, err := h.repo.FindByID(ctx, paymentID)
	if err != nil {
		h.log.Error("Failed to find payment", "error", err)
		return err
	}
	if payment == nil {
		return errors.NotFound(fmt.Sprintf("payment %d", paymentID))
	}

	// Validate cancellation
	if !payment.IsCancellable() {
		return errors.ValidationError("payment cannot be cancelled in current status")
	}

	// Cancel payment
	if err := payment.Cancel(); err != nil {
		return errors.ValidationError(err.Error())
	}

	// Save payment
	if err := h.repo.Update(ctx, payment); err != nil {
		h.log.WithError(err).Error("Failed to cancel payment")
		return errors.InternalWrap(err, "failed to cancel payment")
	}

	h.log.WithField("paymentID", paymentID).Info("Payment cancelled successfully")
	return nil
}
