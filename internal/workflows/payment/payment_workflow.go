package payment

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/qhato/ecommerce/pkg/workflow"
)

// PaymentContext contains payment workflow input/output
type PaymentContext struct {
	OrderID         int64
	CustomerID      int64
	Amount          decimal.Decimal
	PaymentMethodID int64
	CurrencyCode    string
	
	// Workflow state
	AuthorizationID *string
	CaptureID       *string
	Authorized      bool
	Captured        bool
	Refunded        bool
	
	Metadata map[string]interface{}
}

// ValidatePaymentActivity validates payment details
type ValidatePaymentActivity struct {
	workflow.BaseActivity
	paymentService PaymentService
}

type PaymentService interface {
	ValidatePaymentMethod(ctx context.Context, paymentMethodID int64, customerID int64) error
	AuthorizePayment(ctx context.Context, paymentMethodID int64, amount decimal.Decimal) (string, error)
	CapturePayment(ctx context.Context, authorizationID string, amount decimal.Decimal) (string, error)
	VoidAuthorization(ctx context.Context, authorizationID string) error
	RefundPayment(ctx context.Context, captureID string, amount decimal.Decimal) error
}

func NewValidatePaymentActivity(paymentService PaymentService) *ValidatePaymentActivity {
	return &ValidatePaymentActivity{
		BaseActivity:   workflow.NewBaseActivity("ValidatePayment", "Validate payment method"),
		paymentService: paymentService,
	}
}

func (a *ValidatePaymentActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	paymentCtx, ok := input.(*PaymentContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	if err := a.paymentService.ValidatePaymentMethod(ctx, paymentCtx.PaymentMethodID, paymentCtx.CustomerID); err != nil {
		return nil, fmt.Errorf("payment validation failed: %w", err)
	}

	return paymentCtx, nil
}

func (a *ValidatePaymentActivity) Compensate(ctx context.Context, input interface{}) error {
	return nil
}

// AuthorizePaymentActivity authorizes the payment
type AuthorizePaymentActivity struct {
	workflow.BaseActivity
	paymentService PaymentService
}

func NewAuthorizePaymentActivity(paymentService PaymentService) *AuthorizePaymentActivity {
	return &AuthorizePaymentActivity{
		BaseActivity:   workflow.NewBaseActivity("AuthorizePayment", "Authorize payment"),
		paymentService: paymentService,
	}
}

func (a *AuthorizePaymentActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	paymentCtx, ok := input.(*PaymentContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	authID, err := a.paymentService.AuthorizePayment(ctx, paymentCtx.PaymentMethodID, paymentCtx.Amount)
	if err != nil {
		return nil, fmt.Errorf("authorization failed: %w", err)
	}

	paymentCtx.AuthorizationID = &authID
	paymentCtx.Authorized = true

	return paymentCtx, nil
}

func (a *AuthorizePaymentActivity) Compensate(ctx context.Context, input interface{}) error {
	paymentCtx, ok := input.(*PaymentContext)
	if !ok || !paymentCtx.Authorized || paymentCtx.AuthorizationID == nil {
		return nil
	}

	return a.paymentService.VoidAuthorization(ctx, *paymentCtx.AuthorizationID)
}

// CapturePaymentActivity captures the authorized payment
type CapturePaymentActivity struct {
	workflow.BaseActivity
	paymentService PaymentService
}

func NewCapturePaymentActivity(paymentService PaymentService) *CapturePaymentActivity {
	return &CapturePaymentActivity{
		BaseActivity:   workflow.NewBaseActivity("CapturePayment", "Capture payment"),
		paymentService: paymentService,
	}
}

func (a *CapturePaymentActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	paymentCtx, ok := input.(*PaymentContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	if !paymentCtx.Authorized || paymentCtx.AuthorizationID == nil {
		return nil, fmt.Errorf("payment not authorized")
	}

	captureID, err := a.paymentService.CapturePayment(ctx, *paymentCtx.AuthorizationID, paymentCtx.Amount)
	if err != nil {
		return nil, fmt.Errorf("capture failed: %w", err)
	}

	paymentCtx.CaptureID = &captureID
	paymentCtx.Captured = true

	return paymentCtx, nil
}

func (a *CapturePaymentActivity) Compensate(ctx context.Context, input interface{}) error {
	paymentCtx, ok := input.(*PaymentContext)
	if !ok || !paymentCtx.Captured || paymentCtx.CaptureID == nil {
		return nil
	}

	return a.paymentService.RefundPayment(ctx, *paymentCtx.CaptureID, paymentCtx.Amount)
}

// PaymentWorkflow creates a payment workflow
func PaymentWorkflow(paymentService PaymentService) (*workflow.Workflow, error) {
	return workflow.NewWorkflowBuilder("payment", "Payment Workflow").
		Description("Process payment with authorization and capture").
		AddActivities(
			NewValidatePaymentActivity(paymentService),
			NewAuthorizePaymentActivity(paymentService),
			NewCapturePaymentActivity(paymentService),
		).
		MaxRetries(2).
		CompensateOnFail(true).
		Build()
}