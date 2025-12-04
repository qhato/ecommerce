package domain

import (
	"testing"

	"github.com/qhato/ecommerce/pkg/testutil"
)

func TestOrder_AddItem(t *testing.T) {
	// Arrange
	order := &Order{
		ID:     1,
		Status: OrderStatusPending,
		Items:  make([]OrderItem, 0),
	}

	item := OrderItem{
		SKUID:    1,
		Name:     "Test Product",
		Quantity: 2,
		Price:    99.99,
	}

	// Act
	order.AddItem(item)

	// Assert
	testutil.AssertLen(t, order.Items, 1, "Should have 1 item")
	testutil.AssertEqual(t, order.Items[0].Name, "Test Product", "Item name")
	testutil.AssertEqual(t, order.Items[0].Quantity, 2, "Item quantity")
}

func TestOrder_CalculateTotals(t *testing.T) {
	// Arrange
	order := &Order{
		ID:     1,
		Status: OrderStatusPending,
		Items: []OrderItem{
			{
				SKUID:    1,
				Name:     "Product 1",
				Quantity: 2,
				Price:    50.00,
			},
			{
				SKUID:    2,
				Name:     "Product 2",
				Quantity: 1,
				Price:    30.00,
			},
		},
		TaxTotal:      13.00,
		ShippingTotal: 10.00,
	}

	// Act
	order.CalculateTotals()

	// Assert
	// Subtotal = (2 * 50.00) + (1 * 30.00) = 130.00
	testutil.AssertEqual(t, order.Subtotal, 130.00, "Subtotal calculation")
	// Total = 130.00 + 13.00 + 10.00 = 153.00
	testutil.AssertEqual(t, order.Total, 153.00, "Total calculation")
}

func TestOrder_Submit(t *testing.T) {
	// Arrange
	order := &Order{
		ID:     1,
		Status: OrderStatusPending,
		Items: []OrderItem{
			{SKUID: 1, Name: "Product 1", Quantity: 1, Price: 99.99},
		},
	}

	// Act
	err := order.Submit()

	// Assert
	testutil.AssertNoError(t, err, "Submit should succeed")
	testutil.AssertEqual(t, order.Status, OrderStatusSubmitted, "Status should be SUBMITTED")
	testutil.AssertNotNil(t, order.SubmitDate, "Submit date should be set")
}

func TestOrder_Submit_EmptyOrder(t *testing.T) {
	// Arrange
	order := &Order{
		ID:     1,
		Status: OrderStatusPending,
		Items:  make([]OrderItem, 0),
	}

	// Act
	err := order.Submit()

	// Assert
	testutil.AssertError(t, err, "Should fail for empty order")
	testutil.AssertErrorContains(t, err, "empty", "Error should mention empty order")
}

func TestOrder_Submit_AlreadySubmitted(t *testing.T) {
	// Arrange
	order := &Order{
		ID:     1,
		Status: OrderStatusSubmitted,
		Items: []OrderItem{
			{SKUID: 1, Name: "Product 1", Quantity: 1, Price: 99.99},
		},
	}

	// Act
	err := order.Submit()

	// Assert
	testutil.AssertError(t, err, "Should fail for already submitted order")
}

func TestOrder_Cancel(t *testing.T) {
	// Arrange
	order := &Order{
		ID:     1,
		Status: OrderStatusPending,
	}

	// Act
	err := order.Cancel()

	// Assert
	testutil.AssertNoError(t, err, "Cancel should succeed")
	testutil.AssertEqual(t, order.Status, OrderStatusCancelled, "Status should be CANCELLED")
}

func TestOrder_IsCancellable(t *testing.T) {
	tests := []struct {
		name   string
		status OrderStatus
		want   bool
	}{
		{
			name:   "Pending order can be cancelled",
			status: OrderStatusPending,
			want:   true,
		},
		{
			name:   "Submitted order can be cancelled",
			status: OrderStatusSubmitted,
			want:   true,
		},
		{
			name:   "Processing order can be cancelled",
			status: OrderStatusProcessing,
			want:   true,
		},
		{
			name:   "Shipped order cannot be cancelled",
			status: OrderStatusShipped,
			want:   false,
		},
		{
			name:   "Delivered order cannot be cancelled",
			status: OrderStatusDelivered,
			want:   false,
		},
		{
			name:   "Already cancelled order cannot be cancelled again",
			status: OrderStatusCancelled,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := &Order{
				ID:     1,
				Status: tt.status,
			}

			got := order.IsCancellable()
			testutil.AssertEqual(t, got, tt.want, "IsCancellable result")
		})
	}
}

func TestOrder_UpdateStatus(t *testing.T) {
	// Arrange
	order := &Order{
		ID:     1,
		Status: OrderStatusPending,
	}

	// Act
	order.UpdateStatus(OrderStatusProcessing)

	// Assert
	testutil.AssertEqual(t, order.Status, OrderStatusProcessing, "Status should be updated")
}

func TestOrder_GetItemCount(t *testing.T) {
	// Arrange
	order := &Order{
		ID:     1,
		Status: OrderStatusPending,
		Items: []OrderItem{
			{SKUID: 1, Quantity: 2},
			{SKUID: 2, Quantity: 3},
			{SKUID: 3, Quantity: 1},
		},
	}

	// Act
	count := order.GetItemCount()

	// Assert
	testutil.AssertEqual(t, count, 6, "Total item count should be 6 (2+3+1)")
}

func TestOrder_Validation(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		wantErr bool
	}{
		{
			name: "Valid order",
			order: &Order{
				CustomerID: 1,
				Status:     OrderStatusPending,
			},
			wantErr: false,
		},
		{
			name: "No customer ID",
			order: &Order{
				CustomerID: 0,
				Status:     OrderStatusPending,
			},
			wantErr: true,
		},
		{
			name: "Invalid status",
			order: &Order{
				CustomerID: 1,
				Status:     OrderStatus("INVALID"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.Validate()
			if tt.wantErr {
				testutil.AssertError(t, err, "Should return validation error")
			} else {
				testutil.AssertNoError(t, err, "Should not return validation error")
			}
		})
	}
}
