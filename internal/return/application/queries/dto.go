package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/return/domain"
)

type ReturnRequestDTO struct {
	ID             int64           `json:"id"`
	RMA            string          `json:"rma"`
	OrderID        int64           `json:"order_id"`
	CustomerID     string          `json:"customer_id"`
	Status         string          `json:"status"`
	Reason         string          `json:"reason"`
	ReasonDetails  string          `json:"reason_details,omitempty"`
	RefundAmount   float64         `json:"refund_amount"`
	RefundMethod   string          `json:"refund_method"`
	Notes          string          `json:"notes,omitempty"`
	TrackingNumber string          `json:"tracking_number,omitempty"`
	Items          []ReturnItemDTO `json:"items"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	ApprovedBy     *int64          `json:"approved_by,omitempty"`
	ApprovedAt     *time.Time      `json:"approved_at,omitempty"`
	ReceivedAt     *time.Time      `json:"received_at,omitempty"`
	InspectedAt    *time.Time      `json:"inspected_at,omitempty"`
	RefundedAt     *time.Time      `json:"refunded_at,omitempty"`
}

type ReturnItemDTO struct {
	ID         int64   `json:"id"`
	ReturnID   int64   `json:"return_id"`
	ProductID  int64   `json:"product_id"`
	SKU        string  `json:"sku"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
	Condition  string  `json:"condition,omitempty"`
}

func ToReturnRequestDTO(r *domain.ReturnRequest) *ReturnRequestDTO {
	items := make([]ReturnItemDTO, len(r.Items))
	for i, item := range r.Items {
		items[i] = ReturnItemDTO{
			ID:         item.ID,
			ReturnID:   item.ReturnID,
			ProductID:  item.ProductID,
			SKU:        item.SKU,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.TotalPrice,
			Condition:  item.Condition,
		}
	}

	return &ReturnRequestDTO{
		ID:             r.ID,
		RMA:            r.RMA,
		OrderID:        r.OrderID,
		CustomerID:     r.CustomerID,
		Status:         string(r.Status),
		Reason:         string(r.Reason),
		ReasonDetails:  r.ReasonDetails,
		RefundAmount:   r.RefundAmount,
		RefundMethod:   r.RefundMethod,
		Notes:          r.Notes,
		TrackingNumber: r.TrackingNumber,
		Items:          items,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
		ApprovedBy:     r.ApprovedBy,
		ApprovedAt:     r.ApprovedAt,
		ReceivedAt:     r.ReceivedAt,
		InspectedAt:    r.InspectedAt,
		RefundedAt:     r.RefundedAt,
	}
}
