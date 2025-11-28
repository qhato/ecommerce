package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/payment/domain"
	"github.com/qhato/ecommerce/pkg/apperrors"
	"github.com/qhato/ecommerce/pkg/cache"
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
	h.log.Debug("Fetching payment by ID", "id", id)

	// Try cache first
	cacheKey := fmt.Sprintf("payment:id:%d", id)
	if cached, err := h.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var payment domain.Payment
		if err := json.Unmarshal([]byte(cached), &payment); err == nil {
			h.log.Debug("Payment found in cache", "id", id)
			return &payment, nil
		}
	}

	// Fetch from repository
	payment, err := h.repo.FindByID(ctx, id)
	if err != nil {
		h.log.Error("Failed to fetch payment by ID", "error", err)
		return nil, err
	}
	if payment == nil {
		return nil, apperrors.NewNotFoundError("payment", id)
	}

	// Cache result
	if data, err := json.Marshal(payment); err == nil {
		_ = h.cache.Set(ctx, cacheKey, string(data), 5*time.Minute)
	}

	return payment, nil
}

// GetByTransactionID retrieves a payment by transaction ID
func (h *PaymentQueryHandler) GetByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error) {
	h.log.Debug("Fetching payment by transaction ID", "transactionID", transactionID)

	// Try cache first
	cacheKey := fmt.Sprintf("payment:txn:%s", transactionID)
	if cached, err := h.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var payment domain.Payment
		if err := json.Unmarshal([]byte(cached), &payment); err == nil {
			h.log.Debug("Payment found in cache", "transactionID", transactionID)
			return &payment, nil
		}
	}

	// Fetch from repository
	payment, err := h.repo.FindByTransactionID(ctx, transactionID)
	if err != nil {
		h.log.Error("Failed to fetch payment by transaction ID", "error", err)
		return nil, err
	}
	if payment == nil {
		return nil, apperrors.NewNotFoundError("payment with transaction ID", transactionID)
	}

	// Cache result
	if data, err := json.Marshal(payment); err == nil {
		_ = h.cache.Set(ctx, cacheKey, string(data), 5*time.Minute)
	}

	return payment, nil
}

// ListByOrder retrieves payments for an order
func (h *PaymentQueryHandler) ListByOrder(ctx context.Context, orderID int64) ([]*domain.Payment, error) {
	h.log.Debug("Fetching payments by order", "orderID", orderID)

	payments, err := h.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		h.log.Error("Failed to fetch payments by order", "error", err)
		return nil, err
	}

	return payments, nil
}

// ListByCustomer retrieves payments for a customer
func (h *PaymentQueryHandler) ListByCustomer(ctx context.Context, customerID int64, filter *domain.PaymentFilter) ([]*domain.Payment, int64, error) {
	h.log.Debug("Fetching payments by customer", "customerID", customerID)

	payments, total, err := h.repo.FindByCustomerID(ctx, customerID, filter)
	if err != nil {
		h.log.Error("Failed to fetch payments by customer", "error", err)
		return nil, 0, err
	}

	return payments, total, nil
}

// List retrieves all payments with optional filtering
func (h *PaymentQueryHandler) List(ctx context.Context, filter *domain.PaymentFilter) ([]*domain.Payment, int64, error) {
	h.log.Debug("Fetching all payments with filter", "filter", filter)

	payments, total, err := h.repo.FindAll(ctx, filter)
	if err != nil {
		h.log.Error("Failed to fetch payments", "error", err)
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
