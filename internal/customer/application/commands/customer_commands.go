package commands

import (
	"context"

	"github.com/qhato/ecommerce/internal/customer/domain"
	"github.com/qhato/ecommerce/pkg/auth"
	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/validator"
)

// RegisterCustomerCommand represents a command to register a new customer
type RegisterCustomerCommand struct {
	EmailAddress string `json:"email_address" validate:"required,email"`
	UserName     string `json:"user_name" validate:"required,min=3,max=50"`
	Password     string `json:"password" validate:"required,min=8"`
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	ReceiveEmail bool   `json:"receive_email"`
}

// UpdateCustomerCommand represents a command to update customer profile
type UpdateCustomerCommand struct {
	ID           int64             `json:"id" validate:"required"`
	FirstName    string            `json:"first_name,omitempty"`
	LastName     string            `json:"last_name,omitempty"`
	EmailAddress string            `json:"email_address,omitempty" validate:"omitempty,email"`
	ReceiveEmail *bool             `json:"receive_email,omitempty"`
	Attributes   map[string]string `json:"attributes,omitempty"`
}

// ChangePasswordCommand represents a command to change password
type ChangePasswordCommand struct {
	CustomerID  int64  `json:"customer_id" validate:"required"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// DeactivateCustomerCommand represents a command to deactivate a customer
type DeactivateCustomerCommand struct {
	ID int64 `json:"id" validate:"required"`
}

// ActivateCustomerCommand represents a command to activate a customer
type ActivateCustomerCommand struct {
	ID int64 `json:"id" validate:"required"`
}

// CustomerCommandHandler handles customer commands
type CustomerCommandHandler struct {
	repo      domain.CustomerRepository
	eventBus  event.Bus
	validator *validator.Validator
	logger    *logger.Logger
}

// NewCustomerCommandHandler creates a new customer command handler
func NewCustomerCommandHandler(
	repo domain.CustomerRepository,
	eventBus event.Bus,
	validator *validator.Validator,
	logger *logger.Logger,
) *CustomerCommandHandler {
	return &CustomerCommandHandler{
		repo:      repo,
		eventBus:  eventBus,
		validator: validator,
		logger:    logger,
	}
}

// HandleRegisterCustomer handles the register customer command
func (h *CustomerCommandHandler) HandleRegisterCustomer(ctx context.Context, cmd *RegisterCustomerCommand) (int64, error) {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return 0, errors.NewValidationError("invalid register customer command", err)
	}

	// Check if email already exists
	exists, err := h.repo.ExistsByEmail(ctx, cmd.EmailAddress)
	if err != nil {
		return 0, errors.Wrap(err, "failed to check email existence")
	}
	if exists {
		return 0, errors.NewBusinessError("email address already registered")
	}

	// Check if username already exists
	exists, err = h.repo.ExistsByUsername(ctx, cmd.UserName)
	if err != nil {
		return 0, errors.Wrap(err, "failed to check username existence")
	}
	if exists {
		return 0, errors.NewBusinessError("username already taken")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(cmd.Password)
	if err != nil {
		return 0, errors.Wrap(err, "failed to hash password")
	}

	// Create customer entity
	customer := domain.NewCustomer(
		cmd.EmailAddress,
		cmd.UserName,
		hashedPassword,
		cmd.FirstName,
		cmd.LastName,
	)
	customer.ReceiveEmail = cmd.ReceiveEmail

	// Save to repository
	if err := h.repo.Create(ctx, customer); err != nil {
		h.logger.Error("failed to register customer", "error", err)
		return 0, errors.Wrap(err, "failed to register customer")
	}

	// Publish domain event
	event := domain.NewCustomerRegisteredEvent(
		customer.ID,
		customer.EmailAddress,
		customer.UserName,
		customer.FirstName,
		customer.LastName,
	)
	if err := h.eventBus.Publish(ctx, event); err != nil {
		h.logger.Error("failed to publish customer registered event", "error", err)
	}

	h.logger.Info("customer registered", "customer_id", customer.ID, "email", customer.EmailAddress)
	return customer.ID, nil
}

// HandleUpdateCustomer handles the update customer command
func (h *CustomerCommandHandler) HandleUpdateCustomer(ctx context.Context, cmd *UpdateCustomerCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.NewValidationError("invalid update customer command", err)
	}

	// Find existing customer
	customer, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "customer not found")
	}

	if !customer.IsActive() {
		return errors.NewBusinessError("cannot update inactive customer")
	}

	// Track changes for event
	changes := make(map[string]interface{})

	// Update profile if provided
	if cmd.FirstName != "" || cmd.LastName != "" || cmd.EmailAddress != "" {
		firstName := cmd.FirstName
		if firstName == "" {
			firstName = customer.FirstName
		}
		lastName := cmd.LastName
		if lastName == "" {
			lastName = customer.LastName
		}
		emailAddress := cmd.EmailAddress
		if emailAddress == "" {
			emailAddress = customer.EmailAddress
		}

		// Check if new email is already taken
		if emailAddress != customer.EmailAddress {
			exists, err := h.repo.ExistsByEmail(ctx, emailAddress)
			if err != nil {
				return errors.Wrap(err, "failed to check email existence")
			}
			if exists {
				return errors.NewBusinessError("email address already in use")
			}
			changes["email"] = emailAddress
		}

		customer.UpdateProfile(firstName, lastName, emailAddress)
		if cmd.FirstName != "" {
			changes["first_name"] = firstName
		}
		if cmd.LastName != "" {
			changes["last_name"] = lastName
		}
	}

	// Update receive email preference
	if cmd.ReceiveEmail != nil {
		customer.ReceiveEmail = *cmd.ReceiveEmail
		changes["receive_email"] = *cmd.ReceiveEmail
	}

	// Update attributes
	if cmd.Attributes != nil {
		for name, value := range cmd.Attributes {
			customer.UpdateAttribute(name, value)
		}
		changes["attributes"] = true
	}

	// Save to repository
	if err := h.repo.Update(ctx, customer); err != nil {
		h.logger.Error("failed to update customer", "error", err, "customer_id", cmd.ID)
		return errors.Wrap(err, "failed to update customer")
	}

	// Publish domain event
	if len(changes) > 0 {
		event := domain.NewCustomerUpdatedEvent(customer.ID, changes)
		if err := h.eventBus.Publish(ctx, event); err != nil {
			h.logger.Error("failed to publish customer updated event", "error", err)
		}
	}

	h.logger.Info("customer updated", "customer_id", customer.ID)
	return nil
}

// HandleChangePassword handles the change password command
func (h *CustomerCommandHandler) HandleChangePassword(ctx context.Context, cmd *ChangePasswordCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.NewValidationError("invalid change password command", err)
	}

	// Find customer
	customer, err := h.repo.FindByID(ctx, cmd.CustomerID)
	if err != nil {
		return errors.Wrap(err, "customer not found")
	}

	// Verify old password
	if !auth.CheckPassword(cmd.OldPassword, customer.Password) {
		return errors.NewBusinessError("invalid old password")
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(cmd.NewPassword)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	// Update password
	if err := h.repo.UpdatePassword(ctx, customer.ID, hashedPassword); err != nil {
		h.logger.Error("failed to change password", "error", err, "customer_id", cmd.CustomerID)
		return errors.Wrap(err, "failed to change password")
	}

	// Publish domain event
	event := domain.NewCustomerPasswordChangedEvent(customer.ID)
	if err := h.eventBus.Publish(ctx, event); err != nil {
		h.logger.Error("failed to publish password changed event", "error", err)
	}

	h.logger.Info("password changed", "customer_id", customer.ID)
	return nil
}

// HandleDeactivateCustomer handles the deactivate customer command
func (h *CustomerCommandHandler) HandleDeactivateCustomer(ctx context.Context, cmd *DeactivateCustomerCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.NewValidationError("invalid deactivate customer command", err)
	}

	// Find customer
	customer, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "customer not found")
	}

	if customer.Deactivated {
		return errors.NewBusinessError("customer is already deactivated")
	}

	// Deactivate customer
	customer.Deactivate()

	// Save to repository
	if err := h.repo.Update(ctx, customer); err != nil {
		h.logger.Error("failed to deactivate customer", "error", err, "customer_id", cmd.ID)
		return errors.Wrap(err, "failed to deactivate customer")
	}

	// Publish domain event
	event := domain.NewCustomerDeactivatedEvent(customer.ID)
	if err := h.eventBus.Publish(ctx, event); err != nil {
		h.logger.Error("failed to publish customer deactivated event", "error", err)
	}

	h.logger.Info("customer deactivated", "customer_id", cmd.ID)
	return nil
}

// HandleActivateCustomer handles the activate customer command
func (h *CustomerCommandHandler) HandleActivateCustomer(ctx context.Context, cmd *ActivateCustomerCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.NewValidationError("invalid activate customer command", err)
	}

	// Find customer
	customer, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "customer not found")
	}

	if !customer.Deactivated {
		return errors.NewBusinessError("customer is already active")
	}

	// Activate customer
	customer.Activate()

	// Save to repository
	if err := h.repo.Update(ctx, customer); err != nil {
		h.logger.Error("failed to activate customer", "error", err, "customer_id", cmd.ID)
		return errors.Wrap(err, "failed to activate customer")
	}

	// Publish domain event
	event := domain.NewCustomerActivatedEvent(customer.ID)
	if err := h.eventBus.Publish(ctx, event); err != nil {
		h.logger.Error("failed to publish customer activated event", "error", err)
	}

	h.logger.Info("customer activated", "customer_id", cmd.ID)
	return nil
}
