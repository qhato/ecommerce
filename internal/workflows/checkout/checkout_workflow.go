package checkout

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/qhato/ecommerce/pkg/workflow"
)

// CheckoutContext contains checkout workflow input/output
type CheckoutContext struct {
	CustomerID     int64
	CartID         int64
	CartItems      []CartItem
	ShippingAddress Address
	BillingAddress  Address
	PaymentMethodID int64
	
	// Workflow state
	OrderID           *int64
	InventoryReserved bool
	PricingCalculated bool
	OrderCreated      bool
	
	// Pricing
	Subtotal      decimal.Decimal
	TaxAmount     decimal.Decimal
	ShippingCost  decimal.Decimal
	Total         decimal.Decimal
	
	// Metadata
	Metadata map[string]interface{}
}

// CartItem represents an item in the cart
type CartItem struct {
	ProductID int64
	SKUID     int64
	Quantity  int
	Price     decimal.Decimal
}

// Address represents a shipping/billing address
type Address struct {
	FirstName   string
	LastName    string
	AddressLine1 string
	AddressLine2 string
	City        string
	State       string
	PostalCode  string
	Country     string
	Phone       string
}

// ValidateCartActivity validates the shopping cart
type ValidateCartActivity struct {
	workflow.BaseActivity
	cartService CartService
}

// CartService interface for cart operations
type CartService interface {
	GetCart(ctx context.Context, cartID int64) ([]CartItem, error)
	ValidateCart(ctx context.Context, cartID int64) error
}

// NewValidateCartActivity creates a new validate cart activity
func NewValidateCartActivity(cartService CartService) *ValidateCartActivity {
	return &ValidateCartActivity{
		BaseActivity: workflow.NewBaseActivity("ValidateCart", "Validate shopping cart"),
		cartService:  cartService,
	}
}

func (a *ValidateCartActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	checkoutCtx, ok := input.(*CheckoutContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type, expected *CheckoutContext")
	}

	// Validate cart
	if err := a.cartService.ValidateCart(ctx, checkoutCtx.CartID); err != nil {
		return nil, fmt.Errorf("cart validation failed: %w", err)
	}

	// Get cart items
	items, err := a.cartService.GetCart(ctx, checkoutCtx.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	checkoutCtx.CartItems = items
	return checkoutCtx, nil
}

func (a *ValidateCartActivity) Compensate(ctx context.Context, input interface{}) error {
	// No compensation needed for validation
	return nil
}

// CheckInventoryActivity checks inventory availability
type CheckInventoryActivity struct {
	workflow.BaseActivity
	inventoryService InventoryService
}

// InventoryService interface for inventory operations
type InventoryService interface {
	CheckAvailability(ctx context.Context, skuID int64, quantity int) (bool, error)
	ReserveInventory(ctx context.Context, skuID int64, quantity int) error
	ReleaseInventory(ctx context.Context, skuID int64, quantity int) error
}

// NewCheckInventoryActivity creates a new check inventory activity
func NewCheckInventoryActivity(inventoryService InventoryService) *CheckInventoryActivity {
	return &CheckInventoryActivity{
		BaseActivity:     workflow.NewBaseActivity("CheckInventory", "Check inventory availability"),
		inventoryService: inventoryService,
	}
}

func (a *CheckInventoryActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	checkoutCtx, ok := input.(*CheckoutContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type, expected *CheckoutContext")
	}

	// Check availability for all items
	for _, item := range checkoutCtx.CartItems {
		available, err := a.inventoryService.CheckAvailability(ctx, item.SKUID, item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to check inventory for SKU %d: %w", item.SKUID, err)
		}

		if !available {
			return nil, fmt.Errorf("insufficient inventory for SKU %d", item.SKUID)
		}
	}

	// Reserve inventory
	for _, item := range checkoutCtx.CartItems {
		if err := a.inventoryService.ReserveInventory(ctx, item.SKUID, item.Quantity); err != nil {
			// Rollback previous reservations
			for i := range checkoutCtx.CartItems {
				if checkoutCtx.CartItems[i].SKUID == item.SKUID {
					break
				}
				_ = a.inventoryService.ReleaseInventory(ctx, checkoutCtx.CartItems[i].SKUID, checkoutCtx.CartItems[i].Quantity)
			}
			return nil, fmt.Errorf("failed to reserve inventory for SKU %d: %w", item.SKUID, err)
		}
	}

	checkoutCtx.InventoryReserved = true
	return checkoutCtx, nil
}

func (a *CheckInventoryActivity) Compensate(ctx context.Context, input interface{}) error {
	checkoutCtx, ok := input.(*CheckoutContext)
	if !ok {
		return fmt.Errorf("invalid input type, expected *CheckoutContext")
	}

	if !checkoutCtx.InventoryReserved {
		return nil // Nothing to compensate
	}

	// Release all reserved inventory
	for _, item := range checkoutCtx.CartItems {
		if err := a.inventoryService.ReleaseInventory(ctx, item.SKUID, item.Quantity); err != nil {
			// Log error but continue compensation
			continue
		}
	}

	checkoutCtx.InventoryReserved = false
	return nil
}

// CalculatePricingActivity calculates final pricing
type CalculatePricingActivity struct {
	workflow.BaseActivity
	pricingService PricingService
}

// PricingService interface for pricing calculation
type PricingService interface {
	CalculateOrderPricing(ctx context.Context, items []CartItem, customerID int64, shippingAddress Address) (*PricingResult, error)
}

// PricingResult contains pricing calculation results
type PricingResult struct {
	Subtotal     decimal.Decimal
	TaxAmount    decimal.Decimal
	ShippingCost decimal.Decimal
	Total        decimal.Decimal
}

// NewCalculatePricingActivity creates a new calculate pricing activity
func NewCalculatePricingActivity(pricingService PricingService) *CalculatePricingActivity {
	return &CalculatePricingActivity{
		BaseActivity:   workflow.NewBaseActivity("CalculatePricing", "Calculate order pricing"),
		pricingService: pricingService,
	}
}

func (a *CalculatePricingActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	checkoutCtx, ok := input.(*CheckoutContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type, expected *CheckoutContext")
	}

	pricing, err := a.pricingService.CalculateOrderPricing(
		ctx,
		checkoutCtx.CartItems,
		checkoutCtx.CustomerID,
		checkoutCtx.ShippingAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate pricing: %w", err)
	}

	checkoutCtx.Subtotal = pricing.Subtotal
	checkoutCtx.TaxAmount = pricing.TaxAmount
	checkoutCtx.ShippingCost = pricing.ShippingCost
	checkoutCtx.Total = pricing.Total
	checkoutCtx.PricingCalculated = true

	return checkoutCtx, nil
}

func (a *CalculatePricingActivity) Compensate(ctx context.Context, input interface{}) error {
	// No compensation needed for calculation
	return nil
}

// CreateOrderActivity creates the order
type CreateOrderActivity struct {
	workflow.BaseActivity
	orderService OrderService
}

// OrderService interface for order operations
type OrderService interface {
	CreateOrder(ctx context.Context, customerID int64, items []CartItem, shippingAddress, billingAddress Address, total decimal.Decimal) (int64, error)
	CancelOrder(ctx context.Context, orderID int64) error
}

// NewCreateOrderActivity creates a new create order activity
func NewCreateOrderActivity(orderService OrderService) *CreateOrderActivity {
	return &CreateOrderActivity{
		BaseActivity: workflow.NewBaseActivity("CreateOrder", "Create order"),
		orderService: orderService,
	}
}

func (a *CreateOrderActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	checkoutCtx, ok := input.(*CheckoutContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type, expected *CheckoutContext")
	}

	orderID, err := a.orderService.CreateOrder(
		ctx,
		checkoutCtx.CustomerID,
		checkoutCtx.CartItems,
		checkoutCtx.ShippingAddress,
		checkoutCtx.BillingAddress,
		checkoutCtx.Total,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	checkoutCtx.OrderID = &orderID
	checkoutCtx.OrderCreated = true

	return checkoutCtx, nil
}

func (a *CreateOrderActivity) Compensate(ctx context.Context, input interface{}) error {
	checkoutCtx, ok := input.(*CheckoutContext)
	if !ok {
		return fmt.Errorf("invalid input type, expected *CheckoutContext")
	}

	if !checkoutCtx.OrderCreated || checkoutCtx.OrderID == nil {
		return nil // Nothing to compensate
	}

	// Cancel the order
	if err := a.orderService.CancelOrder(ctx, *checkoutCtx.OrderID); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	checkoutCtx.OrderCreated = false
	return nil
}

// CheckoutWorkflow creates a checkout workflow
func CheckoutWorkflow(
	cartService CartService,
	inventoryService InventoryService,
	pricingService PricingService,
	orderService OrderService,
) (*workflow.Workflow, error) {
	return workflow.NewWorkflowBuilder("checkout", "Checkout Workflow").
		Description("Complete order checkout process with validation and compensation").
		AddActivities(
			NewValidateCartActivity(cartService),
			NewCheckInventoryActivity(inventoryService),
			NewCalculatePricingActivity(pricingService),
			NewCreateOrderActivity(orderService),
		).
		MaxRetries(2).
		CompensateOnFail(true). // Enable compensation for checkout
		Build()
}