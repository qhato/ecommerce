package domain

import "time"

// ReturnStatus represents the status of a return request
type ReturnStatus string

const (
	ReturnStatusRequested ReturnStatus = "REQUESTED"
	ReturnStatusApproved  ReturnStatus = "APPROVED"
	ReturnStatusRejected  ReturnStatus = "REJECTED"
	ReturnStatusReceived  ReturnStatus = "RECEIVED"
	ReturnStatusInspected ReturnStatus = "INSPECTED"
	ReturnStatusRefunded  ReturnStatus = "REFUNDED"
	ReturnStatusCancelled ReturnStatus = "CANCELLED"
)

// ReturnReason represents the reason for return
type ReturnReason string

const (
	ReturnReasonDefective     ReturnReason = "DEFECTIVE"
	ReturnReasonWrongItem     ReturnReason = "WRONG_ITEM"
	ReturnReasonNotAsDescribed ReturnReason = "NOT_AS_DESCRIBED"
	ReturnReasonChangedMind   ReturnReason = "CHANGED_MIND"
	ReturnReasonOther         ReturnReason = "OTHER"
)

// ReturnRequest represents a product return request
type ReturnRequest struct {
	ID              int64
	RMA             string // Return Merchandise Authorization number
	OrderID         int64
	CustomerID      string
	Status          ReturnStatus
	Reason          ReturnReason
	ReasonDetails   string
	Items           []ReturnItem
	RefundAmount    float64
	RefundMethod    string
	ApprovedBy      *int64
	ApprovedAt      *time.Time
	ReceivedAt      *time.Time
	InspectedAt     *time.Time
	RefundedAt      *time.Time
	TrackingNumber  string
	Notes           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ReturnItem represents an item in a return request
type ReturnItem struct {
	ID         int64
	ReturnID   int64
	ProductID  int64
	SKU        string
	Quantity   int
	UnitPrice  float64
	TotalPrice float64
	Condition  string
}

// NewReturnRequest creates a new return request
func NewReturnRequest(orderID int64, customerID string, reason ReturnReason) (*ReturnRequest, error) {
	now := time.Now()
	return &ReturnRequest{
		RMA:        generateRMA(),
		OrderID:    orderID,
		CustomerID: customerID,
		Status:     ReturnStatusRequested,
		Reason:     reason,
		Items:      make([]ReturnItem, 0),
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// Approve approves the return request
func (r *ReturnRequest) Approve(approvedBy int64) {
	now := time.Now()
	r.Status = ReturnStatusApproved
	r.ApprovedBy = &approvedBy
	r.ApprovedAt = &now
	r.UpdatedAt = now
}

// Reject rejects the return request
func (r *ReturnRequest) Reject(reason string) {
	r.Status = ReturnStatusRejected
	r.Notes = reason
	r.UpdatedAt = time.Now()
}

// MarkAsReceived marks the return as received
func (r *ReturnRequest) MarkAsReceived() {
	now := time.Now()
	r.Status = ReturnStatusReceived
	r.ReceivedAt = &now
	r.UpdatedAt = now
}

// MarkAsInspected marks the return as inspected
func (r *ReturnRequest) MarkAsInspected() {
	now := time.Now()
	r.Status = ReturnStatusInspected
	r.InspectedAt = &now
	r.UpdatedAt = now
}

// MarkAsRefunded marks the return as refunded
func (r *ReturnRequest) MarkAsRefunded(refundAmount float64, refundMethod string) {
	now := time.Now()
	r.Status = ReturnStatusRefunded
	r.RefundAmount = refundAmount
	r.RefundMethod = refundMethod
	r.RefundedAt = &now
	r.UpdatedAt = now
}

// Cancel cancels the return request
func (r *ReturnRequest) Cancel() {
	r.Status = ReturnStatusCancelled
	r.UpdatedAt = time.Now()
}

func generateRMA() string {
	return "RMA-" + time.Now().Format("20060102150405")
}
