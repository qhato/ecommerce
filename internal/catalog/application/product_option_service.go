package application

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// ProductOptionService defines the application service for product option related operations.
type ProductOptionService interface {
	// CreateProductOption creates a new product option.
	CreateProductOption(ctx context.Context, cmd *CreateProductOptionCommand) (*ProductOptionDTO, error)

	// GetProductOptionByID retrieves a product option by its ID.
	GetProductOptionByID(ctx context.Context, id int64) (*ProductOptionDTO, error)

	// UpdateProductOption updates an existing product option.
	UpdateProductOption(ctx context.Context, cmd *UpdateProductOptionCommand) (*ProductOptionDTO, error)

	// DeleteProductOption deletes a product option by its ID.
	DeleteProductOption(ctx context.Context, id int64) error

	// CreateProductOptionValue adds a value to a product option.
	CreateProductOptionValue(ctx context.Context, productOptionID int64, cmd *CreateProductOptionValueCommand) (*ProductOptionValueDTO, error)

	// GetProductOptionValueByID retrieves a product option value by its ID.
	GetProductOptionValueByID(ctx context.Context, id int64) (*ProductOptionValueDTO, error)

	// UpdateProductOptionValue updates an existing product option value.
	UpdateProductOptionValue(ctx context.Context, id int64, cmd *UpdateProductOptionValueCommand) (*ProductOptionValueDTO, error)

	// DeleteProductOptionValue deletes a product option value by its ID.
	DeleteProductOptionValue(ctx context.Context, id int64) error
}

// ProductOptionDTO represents a product option data transfer object.
type ProductOptionDTO struct {
	ID                     int64
	AttributeName          string
	DisplayOrder           int
	ErrorCode              string
	ErrorMessage           string
	Label                  string
	LongDescription        string
	Name                   string
	ValidationStrategyType string
	ValidationType         string
	Required               bool
	OptionType             string
	UseInSKUGeneration     bool
	ValidationString       string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// ProductOptionValueDTO represents a product option value data transfer object.
type ProductOptionValueDTO struct {
	ID              int64
	ProductOptionID int64
	AttributeValue  string
	DisplayOrder    int
	PriceAdjustment float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// CreateProductOptionCommand is a command to create a new product option.
type CreateProductOptionCommand struct {
	Name                   string
	Label                  string
	AttributeName          string
	Required               bool
	DisplayOrder           int
	LongDescription        string
	OptionType             string
	UseInSKUGeneration     bool
	ErrorCode              string
	ErrorMessage           string
	ValidationStrategyType string
	ValidationType         string
	ValidationString       string
}

// UpdateProductOptionCommand is a command to update an existing product option.
type UpdateProductOptionCommand struct {
	ID                     int64
	Name                   *string
	Label                  *string
	AttributeName          *string
	Required               *bool
	DisplayOrder           *int
	LongDescription        *string
	OptionType             *string
	UseInSKUGeneration     *bool
	ErrorCode              *string
	ErrorMessage           *string
	ValidationStrategyType *string
	ValidationType         *string
	ValidationString       *string
}

// CreateProductOptionValueCommand is a command to create a new product option value.
type CreateProductOptionValueCommand struct {
	AttributeValue  string
	DisplayOrder    int
	PriceAdjustment float64
}

// UpdateProductOptionValueCommand is a command to update an existing product option value.
type UpdateProductOptionValueCommand struct {
	ID              int64
	AttributeValue  *string
	DisplayOrder    *int
	PriceAdjustment *float64
}

type productOptionService struct {
	productOptionRepo      domain.ProductOptionRepository
	productOptionValueRepo domain.ProductOptionValueRepository
}

// NewProductOptionService creates a new instance of ProductOptionService.
func NewProductOptionService(
	productOptionRepo domain.ProductOptionRepository,
	productOptionValueRepo domain.ProductOptionValueRepository,
) ProductOptionService {
	return &productOptionService{
		productOptionRepo:      productOptionRepo,
		productOptionValueRepo: productOptionValueRepo,
	}
}

func (s *productOptionService) CreateProductOption(ctx context.Context, cmd *CreateProductOptionCommand) (*ProductOptionDTO, error) {
	option := domain.NewProductOption(cmd.Name, cmd.Label, cmd.AttributeName)
	option.Required = cmd.Required
	option.DisplayOrder = cmd.DisplayOrder
	option.LongDescription = cmd.LongDescription
	option.OptionType = cmd.OptionType
	option.UseInSKUGeneration = cmd.UseInSKUGeneration
	option.ErrorCode = cmd.ErrorCode
	option.ErrorMessage = cmd.ErrorMessage
	option.ValidationStrategyType = cmd.ValidationStrategyType
	option.ValidationType = cmd.ValidationType
	option.ValidationString = cmd.ValidationString

	err := s.productOptionRepo.Save(ctx, option)
	if err != nil {
		return nil, fmt.Errorf("failed to save product option: %w", err)
	}

	return toProductOptionDTO(option), nil
}

func (s *productOptionService) GetProductOptionByID(ctx context.Context, id int64) (*ProductOptionDTO, error) {
	option, err := s.productOptionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find product option by ID: %w", err)
	}
	if option == nil {
		return nil, fmt.Errorf("product option with ID %d not found", id)
	}

	return toProductOptionDTO(option), nil
}

func (s *productOptionService) UpdateProductOption(ctx context.Context, cmd *UpdateProductOptionCommand) (*ProductOptionDTO, error) {
	option, err := s.productOptionRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product option by ID for update: %w", err)
	}
	if option == nil {
		return nil, fmt.Errorf("product option with ID %d not found for update", cmd.ID)
	}

	if cmd.Name != nil {
		option.Name = *cmd.Name
	}
	if cmd.Label != nil {
		option.Label = *cmd.Label
	}
	if cmd.AttributeName != nil {
		option.AttributeName = *cmd.AttributeName
	}
	if cmd.Required != nil {
		option.Required = *cmd.Required
	}
	if cmd.DisplayOrder != nil {
		option.SetDisplayOrder(*cmd.DisplayOrder)
	}
	if cmd.LongDescription != nil {
		option.LongDescription = *cmd.LongDescription
	}
	if cmd.OptionType != nil {
		option.OptionType = *cmd.OptionType
	}
	if cmd.UseInSKUGeneration != nil {
		option.SetSKUGenerationFlag(*cmd.UseInSKUGeneration)
	}
	if cmd.ErrorCode != nil || cmd.ErrorMessage != nil || cmd.ValidationStrategyType != nil || cmd.ValidationType != nil || cmd.ValidationString != nil || cmd.Required != nil {
		// Consolidate updates for validation fields
		errorCode := option.ErrorCode
		if cmd.ErrorCode != nil {
			errorCode = *cmd.ErrorCode
		}
		errorMessage := option.ErrorMessage
		if cmd.ErrorMessage != nil {
			errorMessage = *cmd.ErrorMessage
		}
		validationStrategyType := option.ValidationStrategyType
		if cmd.ValidationStrategyType != nil {
			validationStrategyType = *cmd.ValidationStrategyType
		}
		validationType := option.ValidationType
		if cmd.ValidationType != nil {
			validationType = *cmd.ValidationType
		}
		validationString := option.ValidationString
		if cmd.ValidationString != nil {
			validationString = *cmd.ValidationString
		}
		required := option.Required
		if cmd.Required != nil {
			required = *cmd.Required
		}
		option.UpdateValidation(validationStrategyType, validationType, validationString, errorCode, errorMessage, required)
	}

	err = s.productOptionRepo.Save(ctx, option)
	if err != nil {
		return nil, fmt.Errorf("failed to update product option: %w", err)
	}

	return toProductOptionDTO(option), nil
}

func (s *productOptionService) DeleteProductOption(ctx context.Context, id int64) error {
	// First, delete all associated product option values
	err := s.productOptionValueRepo.DeleteByProductOptionID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product option values for option %d: %w", id, err)
	}

	// Then, delete the product option itself
	err = s.productOptionRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product option: %w", err)
	}
	return nil
}

func (s *productOptionService) CreateProductOptionValue(ctx context.Context, productOptionID int64, cmd *CreateProductOptionValueCommand) (*ProductOptionValueDTO, error) {
	value, err := domain.NewProductOptionValue(productOptionID, cmd.AttributeValue, cmd.DisplayOrder, cmd.PriceAdjustment)
	if err != nil {
		return nil, fmt.Errorf("failed to create new product option value: %w", err)
	}

	err = s.productOptionValueRepo.Save(ctx, value)
	if err != nil {
		return nil, fmt.Errorf("failed to save product option value: %w", err)
	}

	return toProductOptionValueDTO(value), nil
}

func (s *productOptionService) GetProductOptionValueByID(ctx context.Context, id int64) (*ProductOptionValueDTO, error) {
	value, err := s.productOptionValueRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find product option value by ID: %w", err)
	}
	if value == nil {
		return nil, fmt.Errorf("product option value with ID %d not found", id)
	}
	return toProductOptionValueDTO(value), nil
}

func (s *productOptionService) UpdateProductOptionValue(ctx context.Context, id int64, cmd *UpdateProductOptionValueCommand) (*ProductOptionValueDTO, error) {
	value, err := s.productOptionValueRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find product option value by ID for update: %w", err)
	}
	if value == nil {
		return nil, fmt.Errorf("product option value with ID %d not found for update", id)
	}

	attributeValue := value.AttributeValue
	if cmd.AttributeValue != nil {
		attributeValue = *cmd.AttributeValue
	}
	displayOrder := value.DisplayOrder
	if cmd.DisplayOrder != nil {
		displayOrder = *cmd.DisplayOrder
	}
	priceAdjustment := value.PriceAdjustment
	if cmd.PriceAdjustment != nil {
		priceAdjustment = *cmd.PriceAdjustment
	}

	value.UpdateValue(attributeValue, displayOrder, priceAdjustment)
	err = s.productOptionValueRepo.Save(ctx, value)
	if err != nil {
		return nil, fmt.Errorf("failed to update product option value: %w", err)
	}

	return toProductOptionValueDTO(value), nil
}

func (s *productOptionService) DeleteProductOptionValue(ctx context.Context, id int64) error {
	err := s.productOptionValueRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product option value: %w", err)
	}
	return nil
}

func toProductOptionDTO(option *domain.ProductOption) *ProductOptionDTO {
	return &ProductOptionDTO{
		ID:                     option.ID,
		AttributeName:          option.AttributeName,
		DisplayOrder:           option.DisplayOrder,
		ErrorCode:              option.ErrorCode,
		ErrorMessage:           option.ErrorMessage,
		Label:                  option.Label,
		LongDescription:        option.LongDescription,
		Name:                   option.Name,
		ValidationStrategyType: option.ValidationStrategyType,
		ValidationType:         option.ValidationType,
		Required:               option.Required,
		OptionType:             option.OptionType,
		UseInSKUGeneration:     option.UseInSKUGeneration,
		ValidationString:       option.ValidationString,
		CreatedAt:              option.CreatedAt,
		UpdatedAt:              option.UpdatedAt,
	}
}

func toProductOptionValueDTO(value *domain.ProductOptionValue) *ProductOptionValueDTO {
	return &ProductOptionValueDTO{
		ID:              value.ID,
		ProductOptionID: value.ProductOptionID,
		AttributeValue:  value.AttributeValue,
		DisplayOrder:    value.DisplayOrder,
		PriceAdjustment: value.PriceAdjustment,
		CreatedAt:       value.CreatedAt,
		UpdatedAt:       value.UpdatedAt,
	}
}