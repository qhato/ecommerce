package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/checkout/domain"
	"github.com/shopspring/decimal"
)

// CheckoutSessionDTO represents a checkout session response
type CheckoutSessionDTO struct {
	ID                string            `json:"id"`
	OrderID           int64             `json:"orderId"`
	CustomerID        *string           `json:"customerId,omitempty"`
	Email             string            `json:"email"`
	IsGuestCheckout   bool              `json:"isGuestCheckout"`
	State             string            `json:"state"`
	CurrentStep       int               `json:"currentStep"`
	CompletedSteps    []string          `json:"completedSteps"`
	Progress          int               `json:"progress"`
	ShippingAddressID *int64            `json:"shippingAddressId,omitempty"`
	BillingAddressID  *int64            `json:"billingAddressId,omitempty"`
	ShippingMethodID  *string           `json:"shippingMethodId,omitempty"`
	PaymentMethodID   *int64            `json:"paymentMethodId,omitempty"`
	Subtotal          decimal.Decimal   `json:"subtotal"`
	ShippingCost      decimal.Decimal   `json:"shippingCost"`
	TaxAmount         decimal.Decimal   `json:"taxAmount"`
	DiscountAmount    decimal.Decimal   `json:"discountAmount"`
	TotalAmount       decimal.Decimal   `json:"totalAmount"`
	CouponCodes       []string          `json:"couponCodes"`
	CustomerNotes     string            `json:"customerNotes,omitempty"`
	ExpiresAt         time.Time         `json:"expiresAt"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
	SubmittedAt       *time.Time        `json:"submittedAt,omitempty"`
	ConfirmedAt       *time.Time        `json:"confirmedAt,omitempty"`
}

// ShippingOptionDTO represents a shipping option response
type ShippingOptionDTO struct {
	ID                    string           `json:"id"`
	Name                  string           `json:"name"`
	Description           string           `json:"description"`
	Carrier               string           `json:"carrier"`
	Speed                 string           `json:"speed"`
	EstimatedDaysMin      int              `json:"estimatedDaysMin"`
	EstimatedDaysMax      int              `json:"estimatedDaysMax"`
	BaseCost              decimal.Decimal  `json:"baseCost"`
	CalculatedCost        *decimal.Decimal `json:"calculatedCost,omitempty"`
	FreeShippingThreshold *decimal.Decimal `json:"freeShippingThreshold,omitempty"`
	IsActive              bool             `json:"isActive"`
	TrackingSupported     bool             `json:"trackingSupported"`
}

// ToCheckoutSessionDTO converts domain model to DTO
func ToCheckoutSessionDTO(session *domain.CheckoutSession) CheckoutSessionDTO {
	return CheckoutSessionDTO{
		ID:                session.ID,
		OrderID:           session.OrderID,
		CustomerID:        session.CustomerID,
		Email:             session.Email,
		IsGuestCheckout:   session.IsGuestCheckout,
		State:             string(session.State),
		CurrentStep:       session.CurrentStep,
		CompletedSteps:    session.CompletedSteps,
		Progress:          session.GetProgress(),
		ShippingAddressID: session.ShippingAddressID,
		BillingAddressID:  session.BillingAddressID,
		ShippingMethodID:  session.ShippingMethodID,
		PaymentMethodID:   session.PaymentMethodID,
		Subtotal:          session.Subtotal,
		ShippingCost:      session.ShippingCost,
		TaxAmount:         session.TaxAmount,
		DiscountAmount:    session.DiscountAmount,
		TotalAmount:       session.TotalAmount,
		CouponCodes:       session.CouponCodes,
		CustomerNotes:     session.CustomerNotes,
		ExpiresAt:         session.ExpiresAt,
		CreatedAt:         session.CreatedAt,
		UpdatedAt:         session.UpdatedAt,
		SubmittedAt:       session.SubmittedAt,
		ConfirmedAt:       session.ConfirmedAt,
	}
}

// ToShippingOptionDTO converts domain model to DTO
func ToShippingOptionDTO(option *domain.ShippingOption, calculatedCost *decimal.Decimal) ShippingOptionDTO {
	return ShippingOptionDTO{
		ID:                    option.ID,
		Name:                  option.Name,
		Description:           option.Description,
		Carrier:               option.Carrier,
		Speed:                 string(option.Speed),
		EstimatedDaysMin:      option.EstimatedDaysMin,
		EstimatedDaysMax:      option.EstimatedDaysMax,
		BaseCost:              option.BaseCost,
		CalculatedCost:        calculatedCost,
		FreeShippingThreshold: option.FreeShippingThreshold,
		IsActive:              option.IsActive,
		TrackingSupported:     option.TrackingSupported,
	}
}
