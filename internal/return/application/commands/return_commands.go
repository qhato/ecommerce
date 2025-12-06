package commands

type CreateReturnCommand struct {
	OrderID    int64           `json:"order_id"`
	CustomerID int64           `json:"customer_id"`
	Reason     string          `json:"reason"`
	Items      []ReturnItemDTO `json:"items"`
}

type ReturnItemDTO struct {
	OrderItemID int64  `json:"order_item_id"`
	Quantity    int    `json:"quantity"`
	Reason      string `json:"reason"`
}

type UpdateReturnCommand struct {
	ID     int64  `json:"id"`
	Reason string `json:"reason"`
}

type ApproveReturnCommand struct {
	ID int64 `json:"id"`
}

type RejectReturnCommand struct {
	ID     int64  `json:"id"`
	Reason string `json:"reason"`
}

type ReceiveReturnCommand struct {
	ID int64 `json:"id"`
}

type InspectReturnCommand struct {
	ID             int64   `json:"id"`
	InspectionNote string  `json:"inspection_note"`
	ApprovedAmount float64 `json:"approved_amount"`
}

type ProcessRefundCommand struct {
	ID           int64   `json:"id"`
	RefundAmount float64 `json:"refund_amount"`
	RefundMethod string  `json:"refund_method"`
}

type CancelReturnCommand struct {
	ID int64 `json:"id"`
}
