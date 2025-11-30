package application

import (
	"context"
	"fmt"

	shippingApp "github.com/qhato/ecommerce/internal/fulfillment/application"
	"github.com/qhato/ecommerce/internal/order/domain"
)

// CheckoutService defines the application service for managing the order checkout workflow.
type CheckoutService interface {
	// StartCheckout initializes the checkout process for an order, moving it from PENDING to CUSTOMER_INFO.
	StartCheckout(ctx context.Context, orderID int64) (*OrderDTO, error)

	// UpdateCustomerInformation updates customer details for the order.
	UpdateCustomerInformation(ctx context.Context, orderID int64, cmd *UpdateCustomerInformationCommand) (*OrderDTO, error)

	// SelectShippingAddressAndMethod selects shipping address and method for the order.
	SelectShippingAddressAndMethod(ctx context.Context, orderID int64, cmd *SelectShippingCommand) (*OrderDTO, error)

	// SelectPaymentMethod selects payment method for the order and attempts authorization.
	SelectPaymentMethod(ctx context.Context, orderID int64, cmd *SelectPaymentMethodCommand) (*OrderDTO, error)

	// ConfirmOrder finalizes the order, moving it to SUBMITTED.
	ConfirmOrder(ctx context.Context, orderID int64) (*OrderDTO, error)

	// CancelCheckout cancels the checkout process, moving the order to CANCELLED.
	CancelCheckout(ctx context.Context, orderID int64) error
}

// UpdateCustomerInformationCommand represents the command to update customer details during checkout.
type UpdateCustomerInformationCommand struct {
	EmailAddress string
	FirstName    string
	LastName     string
	PhoneNumber  string
	// BillingAddressID int64 // Example: Reference to an existing address
}

// SelectShippingCommand represents the command to select shipping address and method.
type SelectShippingCommand struct {
	ShippingAddressID   int64
	ShippingMethod      string
	FulfillmentOptionID int64 // Reference to a fulfillment option
}

// SelectPaymentMethodCommand represents the command to select a payment method.
type SelectPaymentMethodCommand struct {
	PaymentMethodType string // e.g., "CREDIT_CARD", "PAYPAL"
	PaymentToken      string // Tokenized payment information
	SavePaymentMethod bool
	// Other payment details
}

type checkoutService struct {
	orderService    OrderService
	shippingService shippingApp.ShippingService
	// customerService  CustomerService // Dependency on Customer service
	// paymentService   PaymentService  // Dependency on Payment service
	// fulfillmentService FulfillmentService // Dependency on Fulfillment service
}

// NewCheckoutService creates a new instance of CheckoutService.
func NewCheckoutService(
	orderService OrderService,
	shippingService shippingApp.ShippingService,
	// customerService CustomerService,
	// paymentService PaymentService,
	// fulfillmentService FulfillmentService,
) CheckoutService {
	return &checkoutService{
		orderService:    orderService,
		shippingService: shippingService,
		// customerService:  customerService,
		// paymentService:   paymentService,
		// fulfillmentService: fulfillmentService,
	}
}

// StartCheckout initializes the checkout process for an order.
func (s *checkoutService) StartCheckout(ctx context.Context, orderID int64) (*OrderDTO, error) {
	order, err := s.orderService.HandleGetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order %d: %w", orderID, err)
	}

	if order.Status != domain.OrderStatusPending {
		return nil, fmt.Errorf("order %d is not in PENDING status, cannot start checkout (current status: %s)", orderID, order.Status)
	}

	err = s.orderService.UpdateOrderStatus(ctx, orderID, domain.OrderStatusCustomerInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to update order %d status to CUSTOMER_INFO: %w", orderID, err)
	}

	return s.orderService.HandleGetOrderByID(ctx, orderID)
}

// UpdateCustomerInformation updates customer details for the order.
func (s *checkoutService) UpdateCustomerInformation(ctx context.Context, orderID int64, cmd *UpdateCustomerInformationCommand) (*OrderDTO, error) {
	order, err := s.orderService.HandleGetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order %d: %w", orderID, err)
	}

	if order.Status != domain.OrderStatusCustomerInfo {
		return nil, fmt.Errorf("order %d is not in CUSTOMER_INFO status (current status: %s)", orderID, order.Status)
	}

	// In a real implementation, this would interact with a CustomerService
	// to update customer details or link addresses.
	// For now, we simulate success and move to next state.

	err = s.orderService.UpdateOrderStatus(ctx, orderID, domain.OrderStatusShipping)
	if err != nil {
		return nil, fmt.Errorf("failed to update order %d status to SHIPPING: %w", orderID, err)
	}

	return s.orderService.HandleGetOrderByID(ctx, orderID)
}

// SelectShippingAddressAndMethod selects shipping address and method for the order.
func (s *checkoutService) SelectShippingAddressAndMethod(ctx context.Context, orderID int64, cmd *SelectShippingCommand) (*OrderDTO, error) {
	order, err := s.orderService.HandleGetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order %d: %w", orderID, err)
	}

	if order.Status != domain.OrderStatusShipping {
		return nil, fmt.Errorf("order %d is not in SHIPPING status (current status: %s)", orderID, order.Status)
	}

	// 1. Validate shipping address
	addressValid, err := s.shippingService.ValidateShippingAddress(ctx, cmd.ShippingAddressID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate shipping address %d: %w", cmd.ShippingAddressID, err)
	}
	if !addressValid {
		return nil, fmt.Errorf("shipping address %d is invalid", cmd.ShippingAddressID)
	}

	// 2. Calculate shipping cost
	shippingCost, err := s.shippingService.CalculateShippingCost(ctx, orderID, cmd.ShippingAddressID, cmd.FulfillmentOptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate shipping cost for order %d: %w", orderID, err)
	}

	// 3. Update order's shipping total
	err = s.orderService.UpdateOrderShippingDetails(ctx, orderID, shippingCost)
	if err != nil {
		return nil, fmt.Errorf("failed to update order %d shipping details: %w", orderID, err)
	}

	// 4. Create fulfillment group (if not already created for this address/method)
	fgCmd := &CreateFulfillmentGroupCommand{
		Type:                "PHYSICAL_GOODS", // This should be determined dynamically
		AddressID:           &cmd.ShippingAddressID,
		FulfillmentOptionID: &cmd.FulfillmentOptionID,
		IsPrimary:           true, // Assuming one FG for now
		Status:              "PENDING",
	}
	_, err = s.orderService.CreateFulfillmentGroup(ctx, orderID, fgCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create fulfillment group for order %d: %w", orderID, err)
	}

	err = s.orderService.UpdateOrderStatus(ctx, orderID, domain.OrderStatusPayment)
	if err != nil {
		return nil, fmt.Errorf("failed to update order %d status to PAYMENT: %w", orderID, err)
	}

	return s.orderService.HandleGetOrderByID(ctx, orderID)
}

// SelectPaymentMethod selects payment method for the order and attempts authorization.
func (s *checkoutService) SelectPaymentMethod(ctx context.Context, orderID int64, cmd *SelectPaymentMethodCommand) (*OrderDTO, error) {
	order, err := s.orderService.HandleGetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order %d: %w", orderID, err)
	}

	if order.Status != domain.OrderStatusPayment {
		return nil, fmt.Errorf("order %d is not in PAYMENT status (current status: %s)", orderID, order.Status)
	}

	// In a real implementation, this would interact with a PaymentService
	// to authorize payment.
	// For now, simulate success and move to next state.

	// Simulate payment authorization
	// if cmd.PaymentToken == "" {
	// 	return nil, NewDomainError("Payment token is required")
	// }
	// paymentAuthorized := true // Simulate successful authorization

	// if !paymentAuthorized {
	// 	return nil, NewDomainError("Payment authorization failed")
	// }

	err = s.orderService.UpdateOrderStatus(ctx, orderID, domain.OrderStatusReview)
	if err != nil {
		return nil, fmt.Errorf("failed to update order %d status to REVIEW: %w", orderID, err)
	}

	return s.orderService.HandleGetOrderByID(ctx, orderID)
}

// ConfirmOrder finalizes the order.
func (s *checkoutService) ConfirmOrder(ctx context.Context, orderID int64) (*OrderDTO, error) {
	order, err := s.orderService.HandleGetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order %d: %w", orderID, err)
	}

	if order.Status != domain.OrderStatusReview {
		return nil, fmt.Errorf("order %d is not in REVIEW status (current status: %s)", orderID, order.Status)
	}

	err = s.orderService.SubmitOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to submit order %d: %w", orderID, err)
	}

	return s.orderService.HandleGetOrderByID(ctx, orderID)
}

// CancelCheckout cancels the checkout process.
func (s *checkoutService) CancelCheckout(ctx context.Context, orderID int64) error {
	order, err := s.orderService.HandleGetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order %d: %w", orderID, err)
	}

	if order.Status == domain.OrderStatusSubmitted || order.Status == domain.OrderStatusFulfilled || order.Status == domain.OrderStatusCancelled {
		return fmt.Errorf("order %d cannot be cancelled from status %s", orderID, order.Status)
	}

	return s.orderService.CancelOrder(ctx, orderID, "Customer cancelled checkout")
}
