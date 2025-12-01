package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/payment/domain"
	"github.com/qhato/ecommerce/pkg/cache"
	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/logger"
)

// PaymentQueryHandler handles payment queries
type PaymentQueryHandler struct {
	repo  domain.PaymentRepository
	cache cache.Cache
	log   *logger.Logger
}

// NewPaymentQueryHandler creates a new PaymentQueryHandler
func NewPaymentQueryHandler(repo domain.PaymentRepository, cache cache.Cache, log *logger.Logger) *PaymentQueryHandler {
	return &PaymentQueryHandler{
		repo:  repo,
		cache: cache,
		log:   log,
	}
}

// GetByID retrieves a payment by ID
func (h *PaymentQueryHandler) GetByID(ctx context.Context, id int64) (*domain.Payment, error) {
	h.log.WithField("id", id).Debug("Fetching payment by ID")

	// Try cache first
	cacheKey := fmt.Sprintf("payment:id:%d", id)
	if cached, err := h.cache.Get(ctx, cacheKey); err == nil && len(cached) > 0 {
		var payment domain.Payment
		if err := json.Unmarshal(cached, &payment); err == nil {
			h.log.WithField("id", id).Debug("Payment found in cache")
			return &payment, nil
		}
	}

	// Fetch from repository
	payment, err := h.repo.FindByID(ctx, id)
	if err != nil {
		h.log.WithError(err).Error("Failed to fetch payment by ID")
		return nil, err
	}
	if payment == nil {
		return nil, errors.NotFound(fmt.Sprintf("payment %d", id))
	}

	// Cache result
	if data, err := json.Marshal(payment); err == nil {
		_ = h.cache.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return payment, nil
}

// GetByTransactionID retrieves a payment by transaction ID
func (h *PaymentQueryHandler) GetByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error) {
	h.log.WithField("transactionID", transactionID).Debug("Fetching payment by transaction ID")

	// Try cache first
	cacheKey := fmt.Sprintf("payment:txn:%s", transactionID)
	if cached, err := h.cache.Get(ctx, cacheKey); err == nil && len(cached) > 0 {
		var payment domain.Payment
		if err := json.Unmarshal(cached, &payment); err == nil {
			h.log.WithField("transactionID", transactionID).Debug("Payment found in cache")
			return &payment, nil
		}
	}

	// Fetch from repository
	payment, err := h.repo.FindByTransactionID(ctx, transactionID)
	if err != nil {
		h.log.WithError(err).Error("Failed to fetch payment by transaction ID")
		return nil, err
	}
	if payment == nil {
		return nil, errors.NotFound(fmt.Sprintf("payment with transaction ID %s", transactionID))
	}

	// Cache result
	if data, err := json.Marshal(payment); err == nil {
		_ = h.cache.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return payment, nil
}

// ListByOrder retrieves payments for an order
func (h *PaymentQueryHandler) ListByOrder(ctx context.Context, orderID int64) ([]*domain.Payment, error) {
	h.log.WithField("orderID", orderID).Debug("Fetching payments by order")

	payments, err := h.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		h.log.WithError(err).Error("Failed to fetch payments by order")
		return nil, err
	}

	return payments, nil
}

// ListByCustomer retrieves payments for a customer
func (h *PaymentQueryHandler) ListByCustomer(ctx context.Context, customerID int64, filter *domain.PaymentFilter) ([]*domain.Payment, int64, error) {
	h.log.WithField("customerID", customerID).Debug("Fetching payments by customer")

	payments, total, err := h.repo.FindByCustomerID(ctx, customerID, filter)
	if err != nil {
		h.log.WithError(err).Error("Failed to fetch payments by customer")
		return nil, 0, err
	}

	return payments, total, nil
}

// List retrieves all payments with optional filtering
func (h *PaymentQueryHandler) List(ctx context.Context, filter *domain.PaymentFilter) ([]*domain.Payment, int64, error) {
	h.log.WithField("filter", filter).Debug("Fetching all payments with filter")

	payments, total, err := h.repo.FindAll(ctx, filter)
	if err != nil {
		h.log.WithError(err).Error("Failed to fetch payments")
		return nil, 0, err
	}

	return payments, total, nil
}

// InvalidateCache invalidates the cache for a payment
func (h *PaymentQueryHandler) InvalidateCache(ctx context.Context, paymentID int64, transactionID string) {
	cacheKey1 := fmt.Sprintf("payment:id:%d", paymentID)
	_ = h.cache.Delete(ctx, cacheKey1)

	if transactionID != "" {
		cacheKey2 := fmt.Sprintf("payment:txn:%s", transactionID)
		_ = h.cache.Delete(ctx, cacheKey2)
	}
}
