package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/pricing/domain"
	"github.com/shopspring/decimal"
)

// PriceListDTO represents a price list data transfer object
type PriceListDTO struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	Code             string    `json:"code"`
	PriceListType    string    `json:"price_list_type"`
	Currency         string    `json:"currency"`
	Priority         int       `json:"priority"`
	IsActive         bool      `json:"is_active"`
	StartDate        *time.Time `json:"start_date,omitempty"`
	EndDate          *time.Time `json:"end_date,omitempty"`
	Description      string    `json:"description"`
	CustomerSegments []string  `json:"customer_segments"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ToPriceListDTO converts a domain PriceList to a DTO
func ToPriceListDTO(priceList *domain.PriceList) *PriceListDTO {
	return &PriceListDTO{
		ID:               priceList.ID,
		Name:             priceList.Name,
		Code:             priceList.Code,
		PriceListType:    string(priceList.PriceListType),
		Currency:         priceList.Currency,
		Priority:         priceList.Priority,
		IsActive:         priceList.IsActive,
		StartDate:        priceList.StartDate,
		EndDate:          priceList.EndDate,
		Description:      priceList.Description,
		CustomerSegments: priceList.CustomerSegments,
		CreatedAt:        priceList.CreatedAt,
		UpdatedAt:        priceList.UpdatedAt,
	}
}

// PriceListItemDTO represents a price list item data transfer object
type PriceListItemDTO struct {
	ID              int64     `json:"id"`
	PriceListID     int64     `json:"price_list_id"`
	SKUID           string    `json:"sku_id"`
	ProductID       *string   `json:"product_id,omitempty"`
	Price           string    `json:"price"`
	CompareAtPrice  *string   `json:"compare_at_price,omitempty"`
	MinQuantity     int       `json:"min_quantity"`
	MaxQuantity     *int      `json:"max_quantity,omitempty"`
	IsActive        bool      `json:"is_active"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ToPriceListItemDTO converts a domain PriceListItem to a DTO
func ToPriceListItemDTO(item *domain.PriceListItem) *PriceListItemDTO {
	dto := &PriceListItemDTO{
		ID:          item.ID,
		PriceListID: item.PriceListID,
		SKUID:       item.SKUID,
		ProductID:   item.ProductID,
		Price:       item.Price.String(),
		MinQuantity: item.MinQuantity,
		MaxQuantity: item.MaxQuantity,
		IsActive:    item.IsActive,
		StartDate:   item.StartDate,
		EndDate:     item.EndDate,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	if item.CompareAtPrice != nil {
		compareAtPrice := item.CompareAtPrice.String()
		dto.CompareAtPrice = &compareAtPrice
	}

	return dto
}

// PricingRuleDTO represents a pricing rule data transfer object
type PricingRuleDTO struct {
	ID                   int64     `json:"id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	RuleType             string    `json:"rule_type"`
	Priority             int       `json:"priority"`
	IsActive             bool      `json:"is_active"`
	StartDate            *time.Time `json:"start_date,omitempty"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	ConditionExpression  string    `json:"condition_expression"`
	ActionType           string    `json:"action_type"`
	ActionValue          string    `json:"action_value"`
	ApplicableSKUs       []string  `json:"applicable_skus"`
	ApplicableCategories []string  `json:"applicable_categories"`
	CustomerSegments     []string  `json:"customer_segments"`
	MinQuantity          int       `json:"min_quantity"`
	MaxQuantity          *int      `json:"max_quantity,omitempty"`
	MinOrderValue        *string   `json:"min_order_value,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// ToPricingRuleDTO converts a domain PricingRule to a DTO
func ToPricingRuleDTO(rule *domain.PricingRule) *PricingRuleDTO {
	dto := &PricingRuleDTO{
		ID:                   rule.ID,
		Name:                 rule.Name,
		Description:          rule.Description,
		RuleType:             string(rule.RuleType),
		Priority:             rule.Priority,
		IsActive:             rule.IsActive,
		StartDate:            rule.StartDate,
		EndDate:              rule.EndDate,
		ConditionExpression:  rule.ConditionExpression,
		ActionType:           string(rule.ActionType),
		ActionValue:          rule.ActionValue.String(),
		ApplicableSKUs:       rule.ApplicableSKUs,
		ApplicableCategories: rule.ApplicableCategories,
		CustomerSegments:     rule.CustomerSegments,
		MinQuantity:          rule.MinQuantity,
		MaxQuantity:          rule.MaxQuantity,
		CreatedAt:            rule.CreatedAt,
		UpdatedAt:            rule.UpdatedAt,
	}

	if rule.MinOrderValue != nil {
		minOrderValue := rule.MinOrderValue.String()
		dto.MinOrderValue = &minOrderValue
	}

	return dto
}

// PricedItemDTO represents a priced item data transfer object
type PricedItemDTO struct {
	SKUID           string                `json:"sku_id"`
	ProductID       *string               `json:"product_id,omitempty"`
	Quantity        int                   `json:"quantity"`
	BasePrice       string                `json:"base_price"`
	SalePrice       *string               `json:"sale_price,omitempty"`
	FinalPrice      string                `json:"final_price"`
	CompareAtPrice  *string               `json:"compare_at_price,omitempty"`
	PriceListID     *int64                `json:"price_list_id,omitempty"`
	PriceListName   *string               `json:"price_list_name,omitempty"`
	DiscountAmount  string                `json:"discount_amount"`
	DiscountPercent string                `json:"discount_percent"`
	Subtotal        string                `json:"subtotal"`
	Savings         string                `json:"savings"`
	Currency        string                `json:"currency"`
	Adjustments     []PriceAdjustmentDTO  `json:"adjustments"`
	IsOnSale        bool                  `json:"is_on_sale"`
}

// ToPricedItemDTO converts a domain PricedItem to a DTO
func ToPricedItemDTO(item *domain.PricedItem) *PricedItemDTO {
	dto := &PricedItemDTO{
		SKUID:           item.SKUID,
		ProductID:       item.ProductID,
		Quantity:        item.Quantity,
		BasePrice:       item.BasePrice.String(),
		FinalPrice:      item.FinalPrice.String(),
		PriceListID:     item.PriceListID,
		PriceListName:   item.PriceListName,
		DiscountAmount:  item.DiscountAmount.String(),
		DiscountPercent: item.DiscountPercent.StringFixed(2),
		Subtotal:        item.Subtotal.String(),
		Savings:         item.GetSavings().String(),
		Currency:        item.Currency,
		Adjustments:     make([]PriceAdjustmentDTO, 0),
		IsOnSale:        item.IsOnSale,
	}

	if item.SalePrice != nil {
		salePrice := item.SalePrice.String()
		dto.SalePrice = &salePrice
	}

	if item.CompareAtPrice != nil {
		compareAtPrice := item.CompareAtPrice.String()
		dto.CompareAtPrice = &compareAtPrice
	}

	for _, adj := range item.Adjustments {
		dto.Adjustments = append(dto.Adjustments, PriceAdjustmentDTO{
			Type:        string(adj.Type),
			Amount:      adj.Amount.String(),
			Reason:      adj.Reason,
			Description: adj.Description,
			Priority:    adj.Priority,
		})
	}

	return dto
}

// PriceAdjustmentDTO represents a price adjustment data transfer object
type PriceAdjustmentDTO struct {
	Type        string `json:"type"`
	Amount      string `json:"amount"`
	Reason      string `json:"reason"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

// PricingResultDTO represents a pricing result data transfer object
type PricingResultDTO struct {
	Items        []*PricedItemDTO `json:"items"`
	Currency     string           `json:"currency"`
	TotalAmount  string           `json:"total_amount"`
	TotalSavings string           `json:"total_savings"`
	PricedAt     time.Time        `json:"priced_at"`
}

// ToPricingResultDTO converts a domain PricingResult to a DTO
func ToPricingResultDTO(result *domain.PricingResult) *PricingResultDTO {
	dto := &PricingResultDTO{
		Items:        make([]*PricedItemDTO, 0),
		Currency:     result.Currency,
		TotalAmount:  result.TotalAmount.String(),
		TotalSavings: result.GetTotalSavings().String(),
		PricedAt:     result.PricedAt,
	}

	for _, item := range result.Items {
		dto.Items = append(dto.Items, ToPricedItemDTO(item))
	}

	return dto
}

// CalculatePriceRequest represents a request to calculate prices
type CalculatePriceRequest struct {
	Currency        string                `json:"currency"`
	CustomerID      *string               `json:"customer_id,omitempty"`
	CustomerSegment *string               `json:"customer_segment,omitempty"`
	Items           []PricingRequestItem  `json:"items"`
}

// PricingRequestItem represents a single item in a pricing request
type PricingRequestItem struct {
	SKUID      string            `json:"sku_id"`
	ProductID  *string           `json:"product_id,omitempty"`
	Quantity   int               `json:"quantity"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// ToPricingContext converts a CalculatePriceRequest to a domain PricingContext
func ToPricingContext(req *CalculatePriceRequest) *domain.PricingContext {
	ctx := domain.NewPricingContext(req.Currency)
	ctx.CustomerID = req.CustomerID
	ctx.CustomerSegment = req.CustomerSegment

	for _, item := range req.Items {
		ctx.RequestedSKUs = append(ctx.RequestedSKUs, domain.PricingRequest{
			SKUID:      item.SKUID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			Attributes: item.Attributes,
		})
	}

	return ctx
}
