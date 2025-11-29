package application

import (
	"time"

	"github.com/qhato/ecommerce/internal/order/domain"
)

// OrderDTO represents order data for transfer
type OrderDTO struct {
	ID            int64          `json:"id"`
	OrderNumber   string         `json:"order_number"`
	CustomerID    int64          `json:"customer_id"`
	EmailAddress  string         `json:"email_address"`
	Name          string         `json:"name"`
	Status        string         `json:"status"`
	SubTotal      float64        `json:"sub_total"`
	TotalTax      float64        `json:"total_tax"`
	TotalShipping float64        `json:"total_shipping"`
	Total         float64        `json:"total"`
	CurrencyCode  string         `json:"currency_code"`
	Items         []OrderItemDTO `json:"items"`
	SubmitDate    *time.Time     `json:"submit_date,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// OrderItemDTO represents order item data for transfer
type OrderItemDTO struct {
	ID             int64   `json:"id"`
	OrderID        int64   `json:"order_id"`
	SKUID          int64   `json:"sku_id"`
	ProductName    string  `json:"product_name"`
	Quantity       int     `json:"quantity"`
	Price          float64 `json:"price"`
	TotalPrice     float64 `json:"total_price"`
	TaxAmount      float64 `json:"tax_amount"`
	ShippingAmount float64 `json:"shipping_amount"`
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

	items := make([]OrderItemDTO, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderItemDTO{
			ID:             item.ID,
			OrderID:        item.OrderID,
			SKUID:          item.SKUID,
			ProductName:    item.ProductName,
			Quantity:       item.Quantity,
			Price:          item.Price,
			TotalPrice:     item.TotalPrice,
			TaxAmount:      item.TaxAmount,
			ShippingAmount: item.ShippingAmount,
		}
	}

	return &OrderDTO{
		ID:            order.ID,
		OrderNumber:   order.OrderNumber,
		CustomerID:    order.CustomerID,
		EmailAddress:  order.EmailAddress,
		Name:          order.Name,
		Status:        string(order.Status),
		SubTotal:      order.SubTotal,
		TotalTax:      order.TotalTax,
		TotalShipping: order.TotalShipping,
		Total:         order.Total,
		CurrencyCode:  order.CurrencyCode,
		Items:         items,
		SubmitDate:    order.SubmitDate,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
	}
}

// ToOrderDTOs converts a slice of domain Orders to OrderDTOs
func ToOrderDTOs(orders []*domain.Order) []OrderDTO {
	dtos := make([]OrderDTO, len(orders))
	for i, order := range orders {
		dto := ToOrderDTO(order)
		if dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}
