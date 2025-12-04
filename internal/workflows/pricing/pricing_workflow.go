package pricing

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/qhato/ecommerce/pkg/workflow"
)

// PricingContext contains pricing workflow input/output
type PricingContext struct {
	ProductID     int64
	Quantity      int
	CustomerID    *int64
	BasePrice     decimal.Decimal
	Subtotal      decimal.Decimal
	Discounts     []Discount
	TotalDiscount decimal.Decimal
	TaxRate       decimal.Decimal
	TaxAmount     decimal.Decimal
	ShippingCost  decimal.Decimal
	FinalPrice    decimal.Decimal
	CurrencyCode  string
	Metadata      map[string]interface{}
}

// Discount represents a pricing discount
type Discount struct {
	ID          int64
	Name        string
	Type        string // "percentage", "fixed", "bogo"
	Value       decimal.Decimal
	Amount      decimal.Decimal
	Description string
}

// GetBasePriceActivity retrieves the base price for a product
type GetBasePriceActivity struct {
	workflow.BaseActivity
	priceService PriceService
}

// PriceService interface for price retrieval
type PriceService interface {
	GetProductPrice(ctx context.Context, productID int64) (decimal.Decimal, error)
}

// NewGetBasePriceActivity creates a new get base price activity
func NewGetBasePriceActivity(priceService PriceService) *GetBasePriceActivity {
	return &GetBasePriceActivity{
		BaseActivity: workflow.NewBaseActivity("GetBasePrice", "Retrieve product base price"),
		priceService: priceService,
	}
}

func (a *GetBasePriceActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	pricingCtx, ok := input.(*PricingContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type, expected *PricingContext")
	}

	basePrice, err := a.priceService.GetProductPrice(ctx, pricingCtx.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get base price: %w", err)
	}

	pricingCtx.BasePrice = basePrice
	pricingCtx.Subtotal = basePrice.Mul(decimal.NewFromInt(int64(pricingCtx.Quantity)))
	pricingCtx.FinalPrice = pricingCtx.Subtotal

	return pricingCtx, nil
}

func (a *GetBasePriceActivity) Compensate(ctx context.Context, input interface{}) error {
	// No compensation needed for read-only operation
	return nil
}

// ApplyPromotionsActivity applies promotions and discounts
type ApplyPromotionsActivity struct {
	workflow.BaseActivity
	promotionService PromotionService
}

// PromotionService interface for promotion logic
type PromotionService interface {
	GetApplicablePromotions(ctx context.Context, productID int64, customerID *int64, quantity int) ([]Discount, error)
}

// NewApplyPromotionsActivity creates a new apply promotions activity
func NewApplyPromotionsActivity(promotionService PromotionService) *ApplyPromotionsActivity {
	return &ApplyPromotionsActivity{
		BaseActivity:     workflow.NewBaseActivity("ApplyPromotions", "Apply promotions and discounts"),
		promotionService: promotionService,
	}
}

func (a *ApplyPromotionsActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	pricingCtx, ok := input.(*PricingContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type, expected *PricingContext")
	}

	promotions, err := a.promotionService.GetApplicablePromotions(
		ctx,
		pricingCtx.ProductID,
		pricingCtx.CustomerID,
		pricingCtx.Quantity,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get promotions: %w", err)
	}

	pricingCtx.Discounts = promotions
	pricingCtx.TotalDiscount = decimal.Zero

	// Calculate discount amounts
	for i, discount := range promotions {
		var amount decimal.Decimal

		switch discount.Type {
		case "percentage":
			amount = pricingCtx.Subtotal.Mul(discount.Value).Div(decimal.NewFromInt(100))
		case "fixed":
			amount = discount.Value.Mul(decimal.NewFromInt(int64(pricingCtx.Quantity)))
		case "bogo":
			// Buy one get one: discount 50% on half the items
			discountedQty := pricingCtx.Quantity / 2
			amount = pricingCtx.BasePrice.Mul(decimal.NewFromInt(int64(discountedQty)))
		default:
			continue
		}

		promotions[i].Amount = amount
		pricingCtx.TotalDiscount = pricingCtx.TotalDiscount.Add(amount)
	}

	pricingCtx.FinalPrice = pricingCtx.Subtotal.Sub(pricingCtx.TotalDiscount)
	if pricingCtx.FinalPrice.LessThan(decimal.Zero) {
		pricingCtx.FinalPrice = decimal.Zero
	}

	return pricingCtx, nil
}

func (a *ApplyPromotionsActivity) Compensate(ctx context.Context, input interface{}) error {
	// No compensation needed
	return nil
}

// CalculateTaxActivity calculates tax amount
type CalculateTaxActivity struct {
	workflow.BaseActivity
	taxService TaxService
}

// TaxService interface for tax calculation
type TaxService interface {
	GetTaxRate(ctx context.Context, productID int64, customerID *int64) (decimal.Decimal, error)
}

// NewCalculateTaxActivity creates a new calculate tax activity
func NewCalculateTaxActivity(taxService TaxService) *CalculateTaxActivity {
	return &CalculateTaxActivity{
		BaseActivity: workflow.NewBaseActivity("CalculateTax", "Calculate tax amount"),
		taxService:   taxService,
	}
}

func (a *CalculateTaxActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	pricingCtx, ok := input.(*PricingContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type, expected *PricingContext")
	}

	taxRate, err := a.taxService.GetTaxRate(ctx, pricingCtx.ProductID, pricingCtx.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax rate: %w", err)
	}

	pricingCtx.TaxRate = taxRate
	pricingCtx.TaxAmount = pricingCtx.FinalPrice.Mul(taxRate).Div(decimal.NewFromInt(100))
	pricingCtx.FinalPrice = pricingCtx.FinalPrice.Add(pricingCtx.TaxAmount)

	return pricingCtx, nil
}

func (a *CalculateTaxActivity) Compensate(ctx context.Context, input interface{}) error {
	// No compensation needed
	return nil
}

// CalculateShippingActivity calculates shipping cost
type CalculateShippingActivity struct {
	workflow.BaseActivity
	shippingService ShippingService
}

// ShippingService interface for shipping cost calculation
type ShippingService interface {
	CalculateShippingCost(ctx context.Context, productID int64, quantity int, customerID *int64) (decimal.Decimal, error)
}

// NewCalculateShippingActivity creates a new calculate shipping activity
func NewCalculateShippingActivity(shippingService ShippingService) *CalculateShippingActivity {
	return &CalculateShippingActivity{
		BaseActivity:    workflow.NewBaseActivity("CalculateShipping", "Calculate shipping cost"),
		shippingService: shippingService,
	}
}

func (a *CalculateShippingActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	pricingCtx, ok := input.(*PricingContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type, expected *PricingContext")
	}

	shippingCost, err := a.shippingService.CalculateShippingCost(
		ctx,
		pricingCtx.ProductID,
		pricingCtx.Quantity,
		pricingCtx.CustomerID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate shipping: %w", err)
	}

	pricingCtx.ShippingCost = shippingCost
	pricingCtx.FinalPrice = pricingCtx.FinalPrice.Add(shippingCost)

	return pricingCtx, nil
}

func (a *CalculateShippingActivity) Compensate(ctx context.Context, input interface{}) error {
	// No compensation needed
	return nil
}

// PricingWorkflow creates a pricing workflow
func PricingWorkflow(
	priceService PriceService,
	promotionService PromotionService,
	taxService TaxService,
	shippingService ShippingService,
) (*workflow.Workflow, error) {
	return workflow.NewWorkflowBuilder("pricing", "Pricing Workflow").
		Description("Calculate final price with promotions, tax, and shipping").
		AddActivities(
			NewGetBasePriceActivity(priceService),
			NewApplyPromotionsActivity(promotionService),
			NewCalculateTaxActivity(taxService),
			NewCalculateShippingActivity(shippingService),
		).
		MaxRetries(2).
		CompensateOnFail(false). // Pricing is read-only, no compensation needed
		Build()
}
