package domain

import (
	"testing"
	"time"

	"github.com/qhato/ecommerce/pkg/testutil"
)

func TestSKU_UpdatePricing(t *testing.T) {
	// Arrange
	sku := &SKU{
		ID:          1,
		Name:        "Test SKU",
		SalePrice:   99.99,
		RetailPrice: 129.99,
	}

	// Act
	sku.UpdatePricing(89.99, 119.99)

	// Assert
	testutil.AssertEqual(t, sku.SalePrice, 89.99, "Sale price should be updated")
	testutil.AssertEqual(t, sku.RetailPrice, 119.99, "Retail price should be updated")
}

func TestSKU_GetEffectivePrice(t *testing.T) {
	tests := []struct {
		name        string
		salePrice   float64
		retailPrice float64
		want        float64
	}{
		{
			name:        "Sale price available",
			salePrice:   79.99,
			retailPrice: 99.99,
			want:        79.99,
		},
		{
			name:        "No sale price",
			salePrice:   0,
			retailPrice: 99.99,
			want:        99.99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sku := &SKU{
				ID:          1,
				Name:        "Test SKU",
				SalePrice:   tt.salePrice,
				RetailPrice: tt.retailPrice,
			}

			got := sku.GetEffectivePrice()
			testutil.AssertEqual(t, got, tt.want, "Effective price")
		})
	}
}

func TestSKU_SetAvailability(t *testing.T) {
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)

	// Arrange
	sku := &SKU{
		ID:     1,
		Name:   "Test SKU",
		Active: false,
	}

	// Act
	sku.SetAvailability(true, &now, &tomorrow)

	// Assert
	testutil.AssertTrue(t, sku.Active, "SKU should be active")
	testutil.AssertNotNil(t, sku.ActiveStartDate, "Start date should be set")
	testutil.AssertNotNil(t, sku.ActiveEndDate, "End date should be set")
}

func TestSKU_IsActive(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	tests := []struct {
		name           string
		active         bool
		activeStartDate *time.Time
		activeEndDate   *time.Time
		want            bool
	}{
		{
			name:            "Active with no dates",
			active:          true,
			activeStartDate: nil,
			activeEndDate:   nil,
			want:            true,
		},
		{
			name:            "Not active",
			active:          false,
			activeStartDate: nil,
			activeEndDate:   nil,
			want:            false,
		},
		{
			name:            "Active but start date in future",
			active:          true,
			activeStartDate: &tomorrow,
			activeEndDate:   nil,
			want:            false,
		},
		{
			name:            "Active but end date in past",
			active:          true,
			activeStartDate: &yesterday,
			activeEndDate:   &yesterday,
			want:            false,
		},
		{
			name:            "Active and currently valid",
			active:          true,
			activeStartDate: &yesterday,
			activeEndDate:   &tomorrow,
			want:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sku := &SKU{
				ID:              1,
				Name:            "Test SKU",
				Active:          tt.active,
				ActiveStartDate: tt.activeStartDate,
				ActiveEndDate:   tt.activeEndDate,
			}

			got := sku.IsActive()
			testutil.AssertEqual(t, got, tt.want, "IsActive result")
		})
	}
}

func TestSKU_Validation(t *testing.T) {
	tests := []struct {
		name    string
		sku     *SKU
		wantErr bool
	}{
		{
			name: "Valid SKU",
			sku: &SKU{
				Name:        "Test SKU",
				RetailPrice: 99.99,
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			sku: &SKU{
				Name:        "",
				RetailPrice: 99.99,
			},
			wantErr: true,
		},
		{
			name: "Negative retail price",
			sku: &SKU{
				Name:        "Test SKU",
				RetailPrice: -10.00,
			},
			wantErr: true,
		},
		{
			name: "Sale price higher than retail",
			sku: &SKU{
				Name:        "Test SKU",
				SalePrice:   150.00,
				RetailPrice: 99.99,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sku.Validate()
			if tt.wantErr {
				testutil.AssertError(t, err, "Should return validation error")
			} else {
				testutil.AssertNoError(t, err, "Should not return validation error")
			}
		})
	}
}

func TestSKU_CalculateDiscount(t *testing.T) {
	tests := []struct {
		name        string
		salePrice   float64
		retailPrice float64
		want        float64
	}{
		{
			name:        "20% discount",
			salePrice:   80.00,
			retailPrice: 100.00,
			want:        20.00,
		},
		{
			name:        "50% discount",
			salePrice:   50.00,
			retailPrice: 100.00,
			want:        50.00,
		},
		{
			name:        "No discount",
			salePrice:   0,
			retailPrice: 100.00,
			want:        0.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sku := &SKU{
				ID:          1,
				Name:        "Test SKU",
				SalePrice:   tt.salePrice,
				RetailPrice: tt.retailPrice,
			}

			got := sku.CalculateDiscount()
			// Allow small floating point differences
			diff := got - tt.want
			if diff < 0 {
				diff = -diff
			}
			testutil.AssertTrue(t, diff < 0.01, "Discount calculation")
		})
	}
}