package domain

import (
	"testing"
	"time"

	"github.com/qhato/ecommerce/pkg/testutil"
)

func TestCategory_SetParentCategory(t *testing.T) {
	// Arrange
	parent := &Category{
		ID:   1,
		Name: "Parent Category",
	}

	child := &Category{
		ID:   2,
		Name: "Child Category",
	}

	// Act
	child.SetParentCategory(parent)

	// Assert
	testutil.AssertNotNil(t, child.ParentCategoryID, "Parent ID should be set")
	testutil.AssertEqual(t, *child.ParentCategoryID, parent.ID, "Parent ID should match")
}

func TestCategory_IsActive(t *testing.T) {
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
			name:            "Currently active",
			activeStartDate: &yesterday,
			activeEndDate:   &tomorrow,
			want:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := &Category{
				ID:              1,
				Name:            "Test Category",
				ActiveStartDate: tt.activeStartDate,
				ActiveEndDate:   tt.activeEndDate,
			}

			got := category.IsActive()
			testutil.AssertEqual(t, got, tt.want, "IsActive result")
		})
	}
}

func TestCategory_AddAttribute(t *testing.T) {
	// Arrange
	category := &Category{
		ID:   1,
		Name: "Test Category",
	}

	// Act
	category.AddAttribute("featured", "true")
	category.AddAttribute("priority", "high")

	// Assert
	testutil.AssertLen(t, category.Attributes, 2, "Should have 2 attributes")
	testutil.AssertEqual(t, category.Attributes["featured"], "true", "Featured attribute")
	testutil.AssertEqual(t, category.Attributes["priority"], "high", "Priority attribute")
}

func TestCategory_Validation(t *testing.T) {
	tests := []struct {
		name     string
		category *Category
		wantErr  bool
	}{
		{
			name: "Valid category",
			category: &Category{
				Name:   "Test Category",
				URLKey: "test-category",
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			category: &Category{
				Name:   "",
				URLKey: "test",
			},
			wantErr: true,
		},
		{
			name: "Empty URL key",
			category: &Category{
				Name:   "Test",
				URLKey: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.category.Validate()
			if tt.wantErr {
				testutil.AssertError(t, err, "Should return validation error")
			} else {
				testutil.AssertNoError(t, err, "Should not return validation error")
			}
		})
	}
}

func TestCategory_GetHierarchy(t *testing.T) {
	// Arrange
	grandParent := &Category{
		ID:   1,
		Name: "Electronics",
	}

	parent := &Category{
		ID:               2,
		Name:             "Computers",
		ParentCategoryID: &grandParent.ID,
	}

	child := &Category{
		ID:               3,
		Name:             "Laptops",
		ParentCategoryID: &parent.ID,
	}

	// Test hierarchy structure
	testutil.AssertNotNil(t, child.ParentCategoryID, "Child should have parent")
	testutil.AssertNotNil(t, parent.ParentCategoryID, "Parent should have grandparent")
	testutil.AssertNil(t, grandParent.ParentCategoryID, "Grandparent should have no parent")
}