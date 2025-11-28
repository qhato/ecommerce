package application

import (
	"time"

	"github.com/qhato/ecommerce/internal/payment/domain"
)

// PaymentDTO represents payment data for transfer
type PaymentDTO struct {
	ID                int64      `json:"id"`
	OrderID           int64      `json:"order_id"`
	CustomerID        int64      `json:"customer_id"`
	PaymentMethod     string     `json:"payment_method"`
	Status            string     `json:"status"`
	Amount            float64    `json:"amount"`
	CurrencyCode      string     `json:"currency_code"`
	TransactionID     string     `json:"transaction_id,omitempty"`
	AuthorizationCode string     `json:"authorization_code,omitempty"`
	RefundAmount      float64    `json:"refund_amount"`
	FailureReason     string     `json:"failure_reason,omitempty"`
	ProcessedDate     *time.Time `json:"processed_date,omitempty"`
	AuthorizedDate    *time.Time `json:"authorized_date,omitempty"`
	CapturedDate      *time.Time `json:"captured_date,omitempty"`
	RefundedDate      *time.Time `json:"refunded_date,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// CreatePaymentRequest represents a request to create a payment
type CreatePaymentRequest struct {
	OrderID       int64   `json:"order_id" validate:"required"`
	CustomerID    int64   `json:"customer_id" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required,oneof=CREDIT_CARD DEBIT_CARD PAYPAL BANK_TRANSFER CASH"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	CurrencyCode  string  `json:"currency_code" validate:"required,len=3"`
}

// ProcessPaymentRequest represents a request to process a payment
type ProcessPaymentRequest struct {
	TransactionID string `json:"transaction_id" validate:"required"`
}

// AuthorizePaymentRequest represents a request to authorize a payment
type AuthorizePaymentRequest struct {
	AuthorizationCode string `json:"authorization_code" validate:"required"`
	TransactionID     string `json:"transaction_id" validate:"required"`
}

// CapturePaymentRequest represents a request to capture an authorized payment
type CapturePaymentRequest struct {
	TransactionID string `json:"transaction_id,omitempty"`
}

// RefundPaymentRequest represents a request to refund a payment
type RefundPaymentRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

// PaginatedPaymentResponse represents a paginated list of payments
type PaginatedPaymentResponse struct {
	Data       []PaymentDTO `json:"data"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalItems int64        `json:"total_items"`
	TotalPages int          `json:"total_pages"`
}

// ToPaymentDTO converts domain Payment to PaymentDTO
func ToPaymentDTO(payment *domain.Payment) *PaymentDTO {
	if payment == nil {
		return nil
	}

	return &PaymentDTO{
		ID:                payment.ID,
		OrderID:           payment.OrderID,
		CustomerID:        payment.CustomerID,
		PaymentMethod:     string(payment.PaymentMethod),
		Status:            string(payment.Status),
		Amount:            payment.Amount,
		CurrencyCode:      payment.CurrencyCode,
		TransactionID:     payment.TransactionID,
		AuthorizationCode: payment.AuthorizationCode,
		RefundAmount:      payment.RefundAmount,
		FailureReason:     payment.FailureReason,
		ProcessedDate:     payment.ProcessedDate,
		AuthorizedDate:    payment.AuthorizedDate,
		CapturedDate:      payment.CapturedDate,
		RefundedDate:      payment.RefundedDate,
		CreatedAt:         payment.CreatedAt,
		UpdatedAt:         payment.UpdatedAt,
	}
}

// ToPaymentDTOs converts a slice of domain Payments to PaymentDTOs
func ToPaymentDTOs(payments []*domain.Payment) []PaymentDTO {
	dtos := make([]PaymentDTO, len(payments))
	for i, payment := range payments {
		dto := ToPaymentDTO(payment)
		if dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}
