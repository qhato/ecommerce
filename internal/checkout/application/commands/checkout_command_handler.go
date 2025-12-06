package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/checkout/application"
	"github.com/qhato/ecommerce/internal/checkout/domain"
)

// CheckoutCommandHandler handles checkout-related commands
type CheckoutCommandHandler struct {
	sessionRepo        domain.CheckoutSessionRepository
	shippingOptionRepo domain.ShippingOptionRepository
	orchestrator       *application.CheckoutOrchestrator
}

// NewCheckoutCommandHandler creates a new command handler
func NewCheckoutCommandHandler(
	sessionRepo domain.CheckoutSessionRepository,
	shippingOptionRepo domain.ShippingOptionRepository,
	orchestrator *application.CheckoutOrchestrator,
) *CheckoutCommandHandler {
	return &CheckoutCommandHandler{
		sessionRepo:        sessionRepo,
		shippingOptionRepo: shippingOptionRepo,
		orchestrator:       orchestrator,
	}
}

// HandleInitiateCheckout handles initiating a new checkout session
func (h *CheckoutCommandHandler) HandleInitiateCheckout(ctx context.Context, cmd InitiateCheckoutCommand) (*domain.CheckoutSession, error) {
	// Check if checkout session already exists for this order
	exists, err := h.sessionRepo.ExistsByOrderID(ctx, cmd.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing session: %w", err)
	}
	if exists {
		return nil, domain.ErrCheckoutSessionAlreadyExists
	}

	// Create new checkout session
	session, err := domain.NewCheckoutSession(cmd.OrderID, cmd.Email, cmd.IsGuest)
	if err != nil {
		return nil, err
	}

	// Set customer ID if provided
	if cmd.CustomerID != nil {
		session.CustomerID = cmd.CustomerID
	}

	// Save session
	if err := h.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return session, nil
}

// HandleAddCustomerInfo handles adding customer information
func (h *CheckoutCommandHandler) HandleAddCustomerInfo(ctx context.Context, cmd AddCustomerInfoCommand) (*domain.CheckoutSession, error) {
	// Find session
	session, err := h.sessionRepo.FindByID(ctx, cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrCheckoutSessionNotFound
	}

	// Check expiration
	if session.IsExpired() {
		return nil, domain.ErrCheckoutSessionExpired
	}

	// Set customer info
	if err := session.SetCustomerInfo(cmd.CustomerID, cmd.Email); err != nil {
		return nil, err
	}

	// Save session
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}

// HandleAddShippingAddress handles adding shipping address
func (h *CheckoutCommandHandler) HandleAddShippingAddress(ctx context.Context, cmd AddShippingAddressCommand) (*domain.CheckoutSession, error) {
	// Find session
	session, err := h.sessionRepo.FindByID(ctx, cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrCheckoutSessionNotFound
	}

	// Check expiration
	if session.IsExpired() {
		return nil, domain.ErrCheckoutSessionExpired
	}

	// TODO: Create address via address service and get address ID
	// For now, using a mock address ID
	addressID := int64(1)

	// Set shipping address
	if err := session.SetShippingAddress(addressID); err != nil {
		return nil, err
	}

	// Save session
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}

// HandleSelectShippingMethod handles selecting shipping method
func (h *CheckoutCommandHandler) HandleSelectShippingMethod(ctx context.Context, cmd SelectShippingMethodCommand) (*domain.CheckoutSession, error) {
	// Find session
	session, err := h.sessionRepo.FindByID(ctx, cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrCheckoutSessionNotFound
	}

	// Check expiration
	if session.IsExpired() {
		return nil, domain.ErrCheckoutSessionExpired
	}

	// Validate shipping method exists
	shippingOption, err := h.shippingOptionRepo.FindByID(ctx, cmd.ShippingMethodID)
	if err != nil {
		return nil, fmt.Errorf("failed to find shipping method: %w", err)
	}
	if shippingOption == nil {
		return nil, domain.ErrShippingMethodNotFound
	}

	// Check if shipping method is available
	if !shippingOption.IsActive {
		return nil, domain.ErrShippingMethodUnavailable
	}

	// Calculate actual shipping cost using orchestrator if available
	// Otherwise use base cost as fallback
	shippingCost := shippingOption.BaseCost
	if h.orchestrator != nil {
		// NOTE: In a real implementation, we would fetch order items, weight, and address
		// For now, we use the base cost. The orchestrator integration is ready when needed.
		// Example usage:
		// calculatedCost, err := h.orchestrator.CalculateShipping(ctx, cmd.ShippingMethodID, weight, orderTotal, country, zipCode)
		// if err == nil {
		//     shippingCost = calculatedCost
		// }
	}

	// Set shipping method
	if err := session.SetShippingMethod(cmd.ShippingMethodID, shippingCost); err != nil {
		return nil, err
	}

	// Save session
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}

// HandleAddBillingAddress handles adding billing address
func (h *CheckoutCommandHandler) HandleAddBillingAddress(ctx context.Context, cmd AddBillingAddressCommand) (*domain.CheckoutSession, error) {
	// Find session
	session, err := h.sessionRepo.FindByID(ctx, cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrCheckoutSessionNotFound
	}

	// Check expiration
	if session.IsExpired() {
		return nil, domain.ErrCheckoutSessionExpired
	}

	// Handle same as shipping
	if cmd.SameAsShipping {
		if err := session.UseSameAddressForBilling(); err != nil {
			return nil, err
		}
	} else {
		// TODO: Create address via address service and get address ID
		// For now, using a mock address ID
		addressID := int64(2)

		if err := session.SetBillingAddress(addressID); err != nil {
			return nil, err
		}
	}

	// Save session
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}

// HandleAddPaymentMethod handles adding payment method
func (h *CheckoutCommandHandler) HandleAddPaymentMethod(ctx context.Context, cmd AddPaymentMethodCommand) (*domain.CheckoutSession, error) {
	// Find session
	session, err := h.sessionRepo.FindByID(ctx, cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrCheckoutSessionNotFound
	}

	// Check expiration
	if session.IsExpired() {
		return nil, domain.ErrCheckoutSessionExpired
	}

	// TODO: Validate and tokenize payment method via payment service
	// For now, using a mock payment method ID
	paymentMethodID := int64(1)

	// Set payment method
	if err := session.SetPaymentMethod(paymentMethodID); err != nil {
		return nil, err
	}

	// Mark ready for submission
	if err := session.MarkReadyForSubmission(); err != nil {
		return nil, err
	}

	// Save session
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}

// HandleApplyCoupon handles applying a coupon code
func (h *CheckoutCommandHandler) HandleApplyCoupon(ctx context.Context, cmd ApplyCouponCommand) (*domain.CheckoutSession, error) {
	// Find session
	session, err := h.sessionRepo.FindByID(ctx, cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrCheckoutSessionNotFound
	}

	// Check expiration
	if session.IsExpired() {
		return nil, domain.ErrCheckoutSessionExpired
	}

	// Apply coupon
	if err := session.ApplyCouponCode(cmd.CouponCode); err != nil {
		return nil, err
	}

	// TODO: Validate coupon via offer service and calculate discount

	// Save session
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}

// HandleRemoveCoupon handles removing a coupon code
func (h *CheckoutCommandHandler) HandleRemoveCoupon(ctx context.Context, cmd RemoveCouponCommand) (*domain.CheckoutSession, error) {
	// Find session
	session, err := h.sessionRepo.FindByID(ctx, cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrCheckoutSessionNotFound
	}

	// Check expiration
	if session.IsExpired() {
		return nil, domain.ErrCheckoutSessionExpired
	}

	// Remove coupon
	session.RemoveCouponCode(cmd.CouponCode)

	// TODO: Recalculate totals

	// Save session
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}

// HandleSubmitCheckout handles submitting the checkout
func (h *CheckoutCommandHandler) HandleSubmitCheckout(ctx context.Context, cmd SubmitCheckoutCommand) (*domain.CheckoutSession, error) {
	// Find session
	session, err := h.sessionRepo.FindByID(ctx, cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrCheckoutSessionNotFound
	}

	// Check expiration
	if session.IsExpired() {
		return nil, domain.ErrCheckoutSessionExpired
	}

	// Submit
	if err := session.Submit(); err != nil {
		return nil, err
	}

	// TODO: Trigger payment processing workflow

	// Save session
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}

// HandleConfirmCheckout handles confirming checkout after payment
func (h *CheckoutCommandHandler) HandleConfirmCheckout(ctx context.Context, cmd ConfirmCheckoutCommand) (*domain.CheckoutSession, error) {
	// Find session
	session, err := h.sessionRepo.FindByID(ctx, cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrCheckoutSessionNotFound
	}

	// Confirm
	if err := session.Confirm(); err != nil {
		return nil, err
	}

	// TODO: Update order status, trigger fulfillment, send confirmation email

	// Save session
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}

// HandleCancelCheckout handles cancelling a checkout
func (h *CheckoutCommandHandler) HandleCancelCheckout(ctx context.Context, cmd CancelCheckoutCommand) error {
	// Find session
	session, err := h.sessionRepo.FindByID(ctx, cmd.SessionID)
	if err != nil {
		return fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return domain.ErrCheckoutSessionNotFound
	}

	// Cancel
	session.Cancel()

	// Save session
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

// HandleExtendSession handles extending session expiration
func (h *CheckoutCommandHandler) HandleExtendSession(ctx context.Context, cmd ExtendSessionCommand) (*domain.CheckoutSession, error) {
	// Find session
	session, err := h.sessionRepo.FindByID(ctx, cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrCheckoutSessionNotFound
	}

	// Extend expiration
	duration := time.Duration(cmd.Hours) * time.Hour
	session.ExtendExpiration(duration)

	// Save session
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}

// Helper methods for orchestrator integration

// ValidateInventoryForCheckout validates inventory availability for all order items
// This method is ready to be called from HandleSubmitCheckout when order items are available
func (h *CheckoutCommandHandler) ValidateInventoryForCheckout(ctx context.Context, orderItems []application.OrderItem) error {
	if h.orchestrator == nil {
		return nil // Skip validation if orchestrator not available
	}
	return h.orchestrator.ValidateInventoryAvailability(ctx, orderItems)
}

// ReserveInventoryForCheckout reserves inventory for the checkout session
// This method is ready to be called from HandleConfirmCheckout when order is confirmed
func (h *CheckoutCommandHandler) ReserveInventoryForCheckout(ctx context.Context, sessionID string, orderID int64, orderItems []application.OrderItem) error {
	if h.orchestrator == nil {
		return nil // Skip if orchestrator not available
	}
	return h.orchestrator.ReserveInventory(ctx, sessionID, orderID, orderItems)
}

// ReleaseInventoryForCheckout releases reserved inventory
// This method is ready to be called from HandleCancelCheckout
func (h *CheckoutCommandHandler) ReleaseInventoryForCheckout(ctx context.Context, sessionID string, orderID int64) error {
	if h.orchestrator == nil {
		return nil // Skip if orchestrator not available
	}
	return h.orchestrator.ReleaseInventory(ctx, sessionID, orderID)
}
