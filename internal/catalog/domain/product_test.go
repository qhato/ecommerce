package domain

import (
	"testing"
	"time"

	"github.com/qhato/ecommerce/pkg/testutil"
)

func TestProduct_Archive(t *testing.T) {
	// Arrange
	product := &Product{
		ID:                    1,
		Name:                  "Test Product",
		ArchivedFlag:          false,
		CanSellWithoutOptions: true,
	}

	// Act
	err := product.Archive()

	// Assert
	testutil.AssertNoError(t, err, "Archive should not return error")
	testutil.AssertTrue(t, product.ArchivedFlag, "Product should be archived")
}

func TestProduct_Unarchive(t *testing.T) {
	// Arrange
	product := &Product{
		ID:           1,
		Name:         "Test Product",
		ArchivedFlag: true,
	}

	// Act
	err := product.Unarchive()

	// Assert
	testutil.AssertNoError(t, err, "Unarchive should not return error")
	testutil.AssertFalse(t, product.ArchivedFlag, "Product should not be archived")
}

func TestProduct_AddAttribute(t *testing.T) {
	// Arrange
	product := &Product{
		ID:   1,
		Name: "Test Product",
	}

	// Act
	product.AddAttribute("color", "red")
	product.AddAttribute("size", "large")

	// Assert
	testutil.AssertLen(t, product.Attributes, 2, "Should have 2 attributes")
	testutil.AssertEqual(t, product.Attributes["color"], "red", "Color attribute")
	testutil.AssertEqual(t, product.Attributes["size"], "large", "Size attribute")
}

func TestProduct_UpdateMetadata(t *testing.T) {
	// Arrange
	product := &Product{
		ID:   1,
		Name: "Test Product",
	}

	metadata := map[string]interface{}{
		"seo_title":       "SEO Title",
		"seo_description": "SEO Description",
	}

	// Act
	product.UpdateMetadata(metadata)

	// Assert
	testutil.AssertEqual(t, product.Metadata["seo_title"], "SEO Title", "SEO title")
	testutil.AssertEqual(t, product.Metadata["seo_description"], "SEO Description", "SEO description")
}

func TestProduct_IsActive(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	tests := []struct {
		name           string
		activeStartDate *time.Time
		activeEndDate   *time.Time
		want            bool
	}{
		{
			name:            "No dates set",
			activeStartDate: nil,
			activeEndDate:   nil,
			want:            true,
		},
		{
			name:            "Start date in past, no end date",
			activeStartDate: &yesterday,
			activeEndDate:   nil,
			want:            true,
		},
		{
			name:            "Start date in future",
			activeStartDate: &tomorrow,
			activeEndDate:   nil,
			want:            false,
		},
		{
			name:            "End date in past",
			activeStartDate: &yesterday,
			activeEndDate:   &yesterday,
			want:            false,
		},
		{
			name:            "Currently active",
			activeStartDate: &yesterday,
			activeEndDate:   &tomorrow,
			want:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &Product{
				ID:              1,
				Name:            "Test Product",
				ActiveStartDate: tt.activeStartDate,
				ActiveEndDate:   tt.activeEndDate,
			}

			got := product.IsActive()
			testutil.AssertEqual(t, got, tt.want, "IsActive result")
		})
	}
}

func TestProduct_Validation(t *testing.T) {
	tests := []struct {
		name    string
		product *Product
		wantErr bool
	}{
		{
			name: "Valid product",
			product: &Product{
				Name:                  "Test Product",
				CanSellWithoutOptions: true,
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			product: &Product{
				Name:                  "",
				CanSellWithoutOptions: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if tt.wantErr {
				testutil.AssertError(t, err, "Should return validation error")
			} else {
				testutil.AssertNoError(t, err, "Should not return validation error")
			}
		})
	}
}

func TestProduct_HasOptions(t *testing.T) {
	tests := []struct {
		name                  string
		canSellWithoutOptions bool
		want                  bool
	}{
		{
			name:                  "Can sell without options",
			canSellWithoutOptions: true,
			want:                  false,
		},
		{
			name:                  "Cannot sell without options (has options)",
			canSellWithoutOptions: false,
			want:                  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &Product{
				ID:                    1,
				Name:                  "Test Product",
				CanSellWithoutOptions: tt.canSellWithoutOptions,
			}

			got := product.HasOptions()
			testutil.AssertEqual(t, got, tt.want, "HasOptions result")
		})
	}
}