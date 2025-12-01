package domain

import "time"

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusAuthorized PaymentStatus = "AUTHORIZED"
	PaymentStatusCaptured   PaymentStatus = "CAPTURED"
	PaymentStatusCompleted  PaymentStatus = "COMPLETED"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusCancelled  PaymentStatus = "CANCELLED"
	PaymentStatusRefunded   PaymentStatus = "REFUNDED"
)

// PaymentMethod represents the method of payment
type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "CREDIT_CARD"
	PaymentMethodDebitCard    PaymentMethod = "DEBIT_CARD"
	PaymentMethodPayPal       PaymentMethod = "PAYPAL"
	PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
	PaymentMethodCash         PaymentMethod = "CASH"
)

// Payment represents a payment entity
type Payment struct {
	ID                int64
	OrderID           int64
	CustomerID        int64
	PaymentMethod     PaymentMethod
	Status            PaymentStatus
	Amount            float64
	CurrencyCode      string
	TransactionID     string
	GatewayResponse   string
	AuthorizationCode string
	RefundAmount      float64
	FailureReason     string
	ProcessedDate     *time.Time
	AuthorizedDate    *time.Time
	CapturedDate      *time.Time
	RefundedDate      *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NewPayment creates a new payment
func NewPayment(orderID, customerID int64, paymentMethod PaymentMethod, amount float64, currencyCode string) *Payment {
	now := time.Now()
	return &Payment{
		OrderID:       orderID,
		CustomerID:    customerID,
		PaymentMethod: paymentMethod,
		Status:        PaymentStatusPending,
		Amount:        amount,
		CurrencyCode:  currencyCode,
		RefundAmount:  0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// Authorize authorizes the payment
func (p *Payment) Authorize(authorizationCode, transactionID string) {
	now := time.Now()
	p.Status = PaymentStatusAuthorized
	p.AuthorizationCode = authorizationCode
	p.TransactionID = transactionID
	p.AuthorizedDate = &now
	p.UpdatedAt = now
}

// Capture captures an authorized payment
func (p *Payment) Capture(transactionID string) error {
	if p.Status != PaymentStatusAuthorized {
		return NewPaymentError("payment must be authorized before capture")
	}
	now := time.Now()
	p.Status = PaymentStatusCaptured
	if transactionID != "" {
		p.TransactionID = transactionID
	}
	p.CapturedDate = &now
	p.UpdatedAt = now
	return nil
}

// Complete completes the payment
func (p *Payment) Complete(transactionID string) {
	now := time.Now()
	p.Status = PaymentStatusCompleted
	if transactionID != "" {
		p.TransactionID = transactionID
	}
	p.ProcessedDate = &now
	p.UpdatedAt = now
}

// Fail marks the payment as failed
func (p *Payment) Fail(reason string) {
	p.Status = PaymentStatusFailed
	p.FailureReason = reason
	p.UpdatedAt = time.Now()
}

// Cancel cancels the payment
func (p *Payment) Cancel() error {
	if p.Status == PaymentStatusCompleted || p.Status == PaymentStatusRefunded {
		return NewPaymentError("cannot cancel completed or refunded payment")
	}
	p.Status = PaymentStatusCancelled
	p.UpdatedAt = time.Now()
	return nil
}

// Refund refunds the payment
func (p *Payment) Refund(amount float64) error {
	if p.Status != PaymentStatusCompleted && p.Status != PaymentStatusCaptured {
		return NewPaymentError("only completed or captured payments can be refunded")
	}
	if amount <= 0 || amount > (p.Amount-p.RefundAmount) {
		return NewPaymentError("invalid refund amount")
	}

	now := time.Now()
	p.RefundAmount += amount

	// If fully refunded, update status
	if p.RefundAmount >= p.Amount {
		p.Status = PaymentStatusRefunded
	}

	p.RefundedDate = &now
	p.UpdatedAt = now
	return nil
}

// IsRefundable checks if payment can be refunded
func (p *Payment) IsRefundable() bool {
	return (p.Status == PaymentStatusCompleted || p.Status == PaymentStatusCaptured) &&
		p.RefundAmount < p.Amount
}

// IsCancellable checks if payment can be cancelled
func (p *Payment) IsCancellable() bool {
	return p.Status == PaymentStatusPending ||
		p.Status == PaymentStatusProcessing ||
		p.Status == PaymentStatusAuthorized
}

// UpdateStatus updates the payment status
func (p *Payment) UpdateStatus(status PaymentStatus) {
	p.Status = status
	p.UpdatedAt = time.Now()
}

// PaymentError represents a payment domain error
type PaymentError struct {
	Message string
}

func (e *PaymentError) Error() string {
	return e.Message
}

// NewPaymentError creates a new payment error
func NewPaymentError(message string) *PaymentError {
	return &PaymentError{Message: message}
}
