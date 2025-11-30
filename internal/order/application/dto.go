package application

import (
	"time"

	"github.com/qhato/ecommerce/internal/order/domain"
)

// OrderDTO represents an order data transfer object.
type OrderDTO struct {
	ID                      int64                     `json:"id"`
	OrderNumber             string                    `json:"order_number"`
	CustomerID              int64                     `json:"customer_id"`
	EmailAddress            string                    `json:"email_address"`
	Name                    string                    `json:"name"`
	Status                  domain.OrderStatus        `json:"status"`
	OrderSubtotal           float64                   `json:"order_subtotal"`
	TotalTax                float64                   `json:"total_tax"`
	TotalShipping           float64                   `json:"total_shipping"`
	OrderTotal              float64                   `json:"order_total"`
	CurrencyCode            string                    `json:"currency_code"`
	IsPreview               bool                      `json:"is_preview"`
	TaxOverride             bool                      `json:"tax_override"`
	LocaleCode              string                    `json:"locale_code"`
	SubmitDate              *time.Time                `json:"submit_date"`
	CreatedAt               time.Time                 `json:"created_at"`
	UpdatedAt               time.Time                 `json:"updated_at"`
	Items                   []*OrderItemDTO           `json:"items"`
	OrderAdjustments        []*OrderAdjustmentDTO     `json:"order_adjustments"`
	FulfillmentGroups       []*FulfillmentGroupDTO    `json:"fulfillment_groups"`
}

// OrderItemDTO represents an order item data transfer object.
type OrderItemDTO struct {
	ID                      int64     `json:"id"`
	OrderID                 int64     `json:"order_id"`
	SKUID                   int64     `json:"sku_id"`
	ProductID               int64     `json:"product_id"`
	Name                    string    `json:"name"`
	Quantity                int       `json:"quantity"`
	RetailPrice             float64   `json:"retail_price"`
	SalePrice               float64   `json:"sale_price"`
	Price                   float64   `json:"price"`
	TotalPrice              float64   `json:"total_price"`
	TaxAmount               float64   `json:"tax_amount"`
	TaxCategory             string    `json:"tax_category"`
	ShippingAmount          float64   `json:"shipping_amount"`
	DiscountsAllowed        bool      `json:"discounts_allowed"`
	HasValidationErrors     bool      `json:"has_validation_errors"`
	ItemTaxableFlag         bool      `json:"item_taxable_flag"`
	OrderItemType           string    `json:"order_item_type"`
	RetailPriceOverride     bool      `json:"retail_price_override"`
	SalePriceOverride       bool      `json:"sale_price_override"`
	CategoryID              *int64    `json:"category_id"`
	GiftWrapItemID          *int64    `json:"gift_wrap_item_id"`
	ParentOrderItemID       *int64    `json:"parent_order_item_id"`
	PersonalMessageID       *int64    `json:"personal_message_id"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// OrderAdjustmentDTO represents an order adjustment data transfer object.
type OrderAdjustmentDTO struct {
	ID               int64     `json:"id"`
	OrderID          int64     `json:"order_id"`
	OfferID          int64     `json:"offer_id"`
	AdjustmentReason string    `json:"adjustment_reason"`
	AdjustmentValue  float64   `json:"adjustment_value"`
	IsFutureCredit   bool      `json:"is_future_credit"`
	CreatedAt        time.Time `json:"created_at"`
}

// OrderItemAdjustmentDTO represents an order item adjustment data transfer object.
type OrderItemAdjustmentDTO struct {
	ID                 int64     `json:"id"`
	OrderItemID        int64     `json:"order_item_id"`
	OfferID            int64     `json:"offer_id"`
	AdjustmentReason   string    `json:"adjustment_reason"`
	AdjustmentValue    float64   `json:"adjustment_value"`
	AppliedToSalePrice bool      `json:"applied_to_sale_price"`
	CreatedAt          time.Time `json:"created_at"`
}

// OrderItemAttributeDTO represents a custom attribute for an order item.
type OrderItemAttributeDTO struct {
	OrderItemID int64     `json:"order_item_id"`
	Name        string    `json:"name"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FulfillmentGroupDTO represents a fulfillment group data transfer object.
type FulfillmentGroupDTO struct {
	ID                   int64     `json:"id"`
	OrderID              int64     `json:"order_id"`
	Type                 string    `json:"type"`
	ShippingPrice        float64   `json:"shipping_price"`
	ShippingPriceTaxable bool      `json:"shipping_price_taxable"`
	MerchandiseTotal     float64   `json:"merchandise_total"`
	Method               string    `json:"method"`
	IsPrimary            bool      `json:"is_primary"`
	ReferenceNumber      string    `json:"reference_number"`
	RetailPrice          float64   `json:"retail_price"`
	SalePrice            float64   `json:"sale_price"`
	Sequence             int       `json:"sequence"`
	Service              string    `json:"service"`
	ShippingOverride     bool      `json:"shipping_override"`
	Status               string    `json:"status"`
	Total                float64   `json:"total"`
	TotalFeeTax          float64   `json:"total_fee_tax"`
	TotalFgTax           float64   `json:"total_fg_tax"`
	TotalItemTax         float64   `json:"total_item_tax"`
	TotalTax             float64   `json:"total_tax"`
	AddressID            *int64    `json:"address_id"`
	FulfillmentOptionID  *int64    `json:"fulfillment_option_id"`
	PersonalMessageID    *int64    `json:"personal_message_id"`
	PhoneID              *int64    `json:"phone_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// CreateOrderRequest represents a request to create an order
type CreateOrderRequest struct {
	CustomerID   int64                    `json:"customer_id" validate:"required"`
	EmailAddress string                   `json:"email_address" validate:"required,email"`
	Name         string                   `json:"name" validate:"required"`
	CurrencyCode string                   `json:"currency_code" validate:"required,len=3"`
	Items        []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

// CreateOrderItemRequest represents a request to add an item to order
type CreateOrderItemRequest struct {
	SKUID       int64   `json:"sku_id" validate:"required"`
	ProductName string  `json:"product_name" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
	Price       float64 `json:"price" validate:"required,min=0"`
}

// UpdateOrderStatusRequest represents a request to update order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=PENDING PROCESSING CONFIRMED SHIPPED DELIVERED CANCELLED REFUNDED"`
}

// AddOrderItemRequest represents a request to add item to existing order
type AddOrderItemRequest struct {
	SKUID       int64   `json:"sku_id" validate:"required"`
	ProductName string  `json:"product_name" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
	Price       float64 `json:"price" validate:"required,min=0"`
}

// PaginatedOrderResponse represents a paginated list of orders
type PaginatedOrderResponse struct {
	Data       []OrderDTO `json:"data"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalItems int64      `json:"total_items"`
	TotalPages int        `json:"total_pages"`
}

// ToOrderDTO converts domain Order to OrderDTO
func ToOrderDTO(order *domain.Order) *OrderDTO {
	if order == nil {
		return nil
	}

	items := make([]*OrderItemDTO, len(order.Items))
	for i := range order.Items {
		items[i] = ToOrderItemDTO(&order.Items[i])
	}

	return &OrderDTO{
		ID:            order.ID,
		OrderNumber:   order.OrderNumber,
		CustomerID:    order.CustomerID,
		EmailAddress:  order.EmailAddress,
		Name:          order.Name,
		Status:        order.Status,
		OrderSubtotal: order.OrderSubtotal,
		TotalTax:      order.TotalTax,
		TotalShipping: order.TotalShipping,
		OrderTotal:    order.OrderTotal,
		CurrencyCode:  order.CurrencyCode,
		IsPreview:     order.IsPreview,
		TaxOverride:   order.TaxOverride,
		LocaleCode:    order.LocaleCode,
		SubmitDate:    order.SubmitDate,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
		Items:         items,
	}
}

// ToOrderItemDTO converts domain OrderItem to OrderItemDTO
func ToOrderItemDTO(item *domain.OrderItem) *OrderItemDTO {
	return &OrderItemDTO{
		ID:                  item.ID,
		OrderID:             item.OrderID,
		SKUID:               item.SKUID,
		ProductID:           item.ProductID,
		Name:                item.Name,
		Quantity:            item.Quantity,
		RetailPrice:         item.RetailPrice,
		SalePrice:           item.SalePrice,
		Price:               item.Price,
		TotalPrice:          item.TotalPrice,
		TaxAmount:           item.TaxAmount,
		TaxCategory:         item.TaxCategory,
		ShippingAmount:      item.ShippingAmount,
		DiscountsAllowed:    item.DiscountsAllowed,
		HasValidationErrors: item.HasValidationErrors,
		ItemTaxableFlag:     item.ItemTaxableFlag,
		OrderItemType:       item.OrderItemType,
		RetailPriceOverride: item.RetailPriceOverride,
		SalePriceOverride:   item.SalePriceOverride,
		CategoryID:          item.CategoryID,
		GiftWrapItemID:      item.GiftWrapItemID,
		ParentOrderItemID:   item.ParentOrderItemID,
		PersonalMessageID:   item.PersonalMessageID,
		CreatedAt:           item.CreatedAt,
		UpdatedAt:           item.UpdatedAt,
	}
}

func ToOrderAdjustmentDTO(adj *domain.OrderAdjustment) *OrderAdjustmentDTO {
	return &OrderAdjustmentDTO{
		ID:               adj.ID,
		OrderID:          adj.OrderID,
		OfferID:          adj.OfferID,
		AdjustmentReason: adj.AdjustmentReason,
		AdjustmentValue:  adj.AdjustmentValue,
		IsFutureCredit:   adj.IsFutureCredit,
		CreatedAt:        adj.CreatedAt,
	}
}

func ToFulfillmentGroupDTO(fg *domain.FulfillmentGroup) *FulfillmentGroupDTO {
	return &FulfillmentGroupDTO{
		ID:                   fg.ID,
		OrderID:              fg.OrderID,
		Type:                 fg.Type,
		ShippingPrice:        fg.ShippingPrice,
		ShippingPriceTaxable: fg.ShippingPriceTaxable,
		MerchandiseTotal:     fg.MerchandiseTotal,
		Method:               fg.Method,
		IsPrimary:            fg.IsPrimary,
		ReferenceNumber:      fg.ReferenceNumber,
		RetailPrice:          fg.RetailPrice,
		SalePrice:            fg.SalePrice,
		Sequence:             fg.Sequence,
		Service:              fg.Service,
		ShippingOverride:     fg.ShippingOverride,
		Status:               fg.Status,
		Total:                fg.Total,
		TotalFeeTax:          fg.TotalFeeTax,
		TotalFgTax:           fg.TotalFgTax,
		TotalItemTax:         fg.TotalItemTax,
		TotalTax:             fg.TotalTax,
		AddressID:            fg.AddressID,
		FulfillmentOptionID:  fg.FulfillmentOptionID,
		PersonalMessageID:    fg.PersonalMessageID,
		PhoneID:              fg.PhoneID,
		CreatedAt:            fg.CreatedAt,
		UpdatedAt:            fg.UpdatedAt,
	}
}

// PaginatedResponse represents a paginated response (generic)
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int64       `json:"total_pages"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, page, pageSize int, totalItems int64) *PaginatedResponse {
	totalPages := totalItems / int64(pageSize)
	if totalItems%int64(pageSize) > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
