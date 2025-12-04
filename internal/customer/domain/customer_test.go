package domain

import (
	"testing"

	"github.com/qhato/ecommerce/pkg/testutil"
	"golang.org/x/crypto/bcrypt"
)

func TestCustomer_UpdateProfile(t *testing.T) {
	// Arrange
	customer := &Customer{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Act
	customer.UpdateProfile("Jane", "Smith", "jane@example.com")

	// Assert
	testutil.AssertEqual(t, customer.FirstName, "Jane", "First name should be updated")
	testutil.AssertEqual(t, customer.LastName, "Smith", "Last name should be updated")
	testutil.AssertEqual(t, customer.Email, "jane@example.com", "Email should be updated")
}

func TestCustomer_ChangePassword(t *testing.T) {
	// Arrange
	oldPassword := "oldpassword123"
	newPassword := "newpassword456"

	hashedOld, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
	customer := &Customer{
		ID:       1,
		Username: "testuser",
		Password: string(hashedOld),
	}

	// Act
	err := customer.ChangePassword(oldPassword, newPassword)

	// Assert
	testutil.AssertNoError(t, err, "Password change should succeed")
	
	// Verify new password works
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(newPassword))
	testutil.AssertNoError(t, err, "New password should be valid")
}

func TestCustomer_ChangePassword_WrongOldPassword(t *testing.T) {
	// Arrange
	oldPassword := "oldpassword123"
	hashedOld, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
	customer := &Customer{
		ID:       1,
		Username: "testuser",
		Password: string(hashedOld),
	}

	// Act
	err := customer.ChangePassword("wrongpassword", "newpassword456")

	// Assert
	testutil.AssertError(t, err, "Should fail with wrong old password")
	testutil.AssertErrorContains(t, err, "invalid", "Error should mention invalid password")
}

func TestCustomer_Deactivate(t *testing.T) {
	// Arrange
	customer := &Customer{
		ID:       1,
		Username: "testuser",
		IsActive: true,
	}

	// Act
	customer.Deactivate()

	// Assert
	testutil.AssertFalse(t, customer.IsActive, "Customer should be deactivated")
}

func TestCustomer_Activate(t *testing.T) {
	// Arrange
	customer := &Customer{
		ID:       1,
		Username: "testuser",
		IsActive: false,
	}

	// Act
	customer.Activate()

	// Assert
	testutil.AssertTrue(t, customer.IsActive, "Customer should be activated")
}

func TestCustomer_Archive(t *testing.T) {
	// Arrange
	customer := &Customer{
		ID:           1,
		Username:     "testuser",
		ArchivedFlag: false,
	}

	// Act
	customer.Archive()

	// Assert
	testutil.AssertTrue(t, customer.ArchivedFlag, "Customer should be archived")
	testutil.AssertFalse(t, customer.IsActive, "Archived customer should be deactivated")
}

func TestCustomer_AddAttribute(t *testing.T) {
	// Arrange
	customer := &Customer{
		ID:       1,
		Username: "testuser",
	}

	// Act
	customer.AddAttribute("preferred_language", "en")
	customer.AddAttribute("newsletter", "true")

	// Assert
	testutil.AssertLen(t, customer.Attributes, 2, "Should have 2 attributes")
	testutil.AssertEqual(t, customer.Attributes["preferred_language"], "en", "Language attribute")
	testutil.AssertEqual(t, customer.Attributes["newsletter"], "true", "Newsletter attribute")
}

func TestCustomer_AddRole(t *testing.T) {
	// Arrange
	customer := &Customer{
		ID:       1,
		Username: "testuser",
		Roles:    make([]string, 0),
	}

	// Act
	customer.AddRole("CUSTOMER")
	customer.AddRole("VIP")

	// Assert
	testutil.AssertLen(t, customer.Roles, 2, "Should have 2 roles")
	testutil.AssertTrue(t, customer.HasRole("CUSTOMER"), "Should have CUSTOMER role")
	testutil.AssertTrue(t, customer.HasRole("VIP"), "Should have VIP role")
}

func TestCustomer_HasRole(t *testing.T) {
	// Arrange
	customer := &Customer{
		ID:       1,
		Username: "testuser",
		Roles:    []string{"CUSTOMER", "VIP"},
	}

	// Test existing role
	testutil.AssertTrue(t, customer.HasRole("CUSTOMER"), "Should have CUSTOMER role")
	testutil.AssertTrue(t, customer.HasRole("VIP"), "Should have VIP role")

	// Test non-existing role
	testutil.AssertFalse(t, customer.HasRole("ADMIN"), "Should not have ADMIN role")
}

func TestCustomer_Validation(t *testing.T) {
	tests := []struct {
		name     string
		customer *Customer
		wantErr  bool
	}{
		{
			name: "Valid customer",
			customer: &Customer{
				Username:  "testuser",
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: false,
		},
		{
			name: "Empty username",
			customer: &Customer{
				Username:  "",
				Email:     "test@example.com",
				FirstName: "John",
			},
			wantErr: true,
		},
		{
			name: "Invalid email",
			customer: &Customer{
				Username:  "testuser",
				Email:     "invalid-email",
				FirstName: "John",
			},
			wantErr: true,
		},
		{
			name: "Empty email",
			customer: &Customer{
				Username:  "testuser",
				Email:     "",
				FirstName: "John",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.customer.Validate()
			if tt.wantErr {
				testutil.AssertError(t, err, "Should return validation error")
			} else {
				testutil.AssertNoError(t, err, "Should not return validation error")
			}
		})
	}
}