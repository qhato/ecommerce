package application

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// CategoryService defines the application service for category-related operations.
type CategoryService interface {
	// CreateCategory creates a new category.
	CreateCategory(ctx context.Context, cmd *CreateCategoryCommand) (*CategoryDTO, error)

	// GetCategoryByID retrieves a category by its ID.
	GetCategoryByID(ctx context.Context, id int64) (*CategoryDTO, error)

	// UpdateCategory updates an existing category.
	UpdateCategory(ctx context.Context, cmd *UpdateCategoryCommand) (*CategoryDTO, error)

	// ArchiveCategory archives a category, making it inactive.
	ArchiveCategory(ctx context.Context, id int64) error

	// UnarchiveCategory unarchives a category, making it active.
	UnarchiveCategory(ctx context.Context, id int64) error

	// SetDefaultParentCategory sets the default parent category for a category.
	SetDefaultParentCategory(ctx context.Context, categoryID, parentCategoryID int64) error

	// RemoveDefaultParentCategory removes the default parent category association.
	RemoveDefaultParentCategory(ctx context.Context, categoryID int64) error

	// AddCategoryAttribute adds a custom attribute to a category.
	AddCategoryAttribute(ctx context.Context, categoryID int64, name, value string) (*CategoryAttributeDTO, error)

	// UpdateCategoryAttribute updates an existing category attribute.
	UpdateCategoryAttribute(ctx context.Context, categoryAttributeID int64, name, value string) (*CategoryAttributeDTO, error)

	// RemoveCategoryAttribute removes a category attribute by its ID.
	RemoveCategoryAttribute(ctx context.Context, categoryAttributeID int64) error
}

// CategoryDTO represents a category data transfer object.
// CategoryAttributeDTO represents a category attribute data transfer object.
type CategoryAttributeDTO struct {
	ID         int64
	Name       string
	Value      string
	CategoryID int64
}

// CreateCategoryCommand is a command to create a new category.
type CreateCategoryCommand struct {
	Name                 string
	Description          string
	LongDescription      string
	URL                  string
	URLKey               string
	ActiveStartDate      *time.Time
	ActiveEndDate        *time.Time
	MetaTitle            string
	MetaDescription      string
	DisplayTemplate      string
	ExternalID           string
	FulfillmentType      string
	InventoryType        string
	OverrideGeneratedURL bool
	ProductDescPattern   string
	ProductTitlePattern  string
	RootDisplayOrder     float64
	TaxCode              string
}

// UpdateCategoryCommand is a command to update an existing category.
type UpdateCategoryCommand struct {
	ID                   int64
	Name                 *string
	Description          *string
	LongDescription      *string
	ActiveStartDate      *time.Time
	ActiveEndDate        *time.Time
	MetaDescription      *string
	MetaTitle            *string
	URL                  *string
	URLKey               *string
	DisplayTemplate      *string
	ExternalID           *string
	FulfillmentType      *string
	InventoryType        *string
	OverrideGeneratedURL *bool
	ProductDescPattern   *string
	ProductTitlePattern  *string
	RootDisplayOrder     *float64
	TaxCode              *string
}

type categoryService struct {
	categoryRepo          domain.CategoryRepository
	categoryAttributeRepo domain.CategoryAttributeRepository
}

// NewCategoryService creates a new instance of CategoryService.
func NewCategoryService(
	categoryRepo domain.CategoryRepository,
	categoryAttributeRepo domain.CategoryAttributeRepository,
) CategoryService {
	return &categoryService{
		categoryRepo:          categoryRepo,
		categoryAttributeRepo: categoryAttributeRepo,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, cmd *CreateCategoryCommand) (*CategoryDTO, error) {
	category := domain.NewCategory(cmd.Name, cmd.Description, cmd.URL, cmd.URLKey)

	category.LongDescription = cmd.LongDescription
	category.ActiveStartDate = cmd.ActiveStartDate
	category.ActiveEndDate = cmd.ActiveEndDate
	category.MetaTitle = cmd.MetaTitle
	category.MetaDescription = cmd.MetaDescription
	category.DisplayTemplate = cmd.DisplayTemplate
	category.ExternalID = cmd.ExternalID
	category.FulfillmentType = cmd.FulfillmentType
	category.InventoryType = cmd.InventoryType
	category.OverrideGeneratedURL = cmd.OverrideGeneratedURL
	category.ProductDescPattern = cmd.ProductDescPattern
	category.ProductTitlePattern = cmd.ProductTitlePattern
	category.RootDisplayOrder = cmd.RootDisplayOrder
	category.TaxCode = cmd.TaxCode

	err := s.categoryRepo.Save(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to save category: %w", err)
	}

	return toCategoryDTO(category), nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id int64) (*CategoryDTO, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find category by ID: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category with ID %d not found", id)
	}
	return toCategoryDTO(category), nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, cmd *UpdateCategoryCommand) (*CategoryDTO, error) {
	category, err := s.categoryRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find category by ID for update: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category with ID %d not found for update", cmd.ID)
	}

	if cmd.Name != nil {
		category.Name = *cmd.Name
	}
	if cmd.Description != nil {
		category.Description = *cmd.Description
	}
	if cmd.LongDescription != nil {
		category.LongDescription = *cmd.LongDescription
	}
	if cmd.ActiveStartDate != nil || cmd.ActiveEndDate != nil {
		category.SetActiveDate(cmd.ActiveStartDate, cmd.ActiveEndDate)
	}
	if cmd.MetaTitle != nil {
		category.MetaTitle = *cmd.MetaTitle
	}
	if cmd.MetaDescription != nil {
		category.MetaDescription = *cmd.MetaDescription
	}
	if cmd.URL != nil {
		category.URL = *cmd.URL
	}
	if cmd.URLKey != nil {
		category.URLKey = *cmd.URLKey
	}
	if cmd.DisplayTemplate != nil {
		category.DisplayTemplate = *cmd.DisplayTemplate
	}
	if cmd.ExternalID != nil {
		category.ExternalID = *cmd.ExternalID
	}
	if cmd.FulfillmentType != nil {
		category.FulfillmentType = *cmd.FulfillmentType
	}
	if cmd.InventoryType != nil {
		category.InventoryType = *cmd.InventoryType
	}
	if cmd.OverrideGeneratedURL != nil {
		category.OverrideGeneratedURL = *cmd.OverrideGeneratedURL
	}
	if cmd.ProductDescPattern != nil {
		category.ProductDescPattern = *cmd.ProductDescPattern
	}
	if cmd.ProductTitlePattern != nil {
		category.ProductTitlePattern = *cmd.ProductTitlePattern
	}
	if cmd.RootDisplayOrder != nil {
		category.SetDisplayOrder(*cmd.RootDisplayOrder)
	}
	if cmd.TaxCode != nil {
		category.TaxCode = *cmd.TaxCode
	}

	err = s.categoryRepo.Save(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return toCategoryDTO(category), nil
}

func (s *categoryService) ArchiveCategory(ctx context.Context, id int64) error {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find category by ID for archiving: %w", err)
	}
	if category == nil {
		return fmt.Errorf("category with ID %d not found for archiving", id)
	}

	category.Archive()
	err = s.categoryRepo.Save(ctx, category)
	if err != nil {
		return fmt.Errorf("failed to archive category: %w", err)
	}
	return nil
}

func (s *categoryService) UnarchiveCategory(ctx context.Context, id int64) error {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find category by ID for unarchiving: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category with ID %d not found for unarchiving", id)
	}

	category.Unarchive()
	err = s.categoryRepo.Save(ctx, category)
	if err != nil {
		return fmt.Errorf("failed to unarchive category: %w", err)
	}
	return nil
}

func (s *categoryService) SetDefaultParentCategory(ctx context.Context, categoryID, parentCategoryID int64) error {
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("failed to find category by ID: %w", err)
	}
	if category == nil {
		return fmt.Errorf("category with ID %d not found", categoryID)
	}

	// Validate parent category exists
	parent, err := s.categoryRepo.FindByID(ctx, parentCategoryID)
	if err != nil {
		return fmt.Errorf("failed to find parent category by ID: %w", err)
	}
	if parent == nil {
		return fmt.Errorf("parent category with ID %d not found", parentCategoryID)
	}

	category.SetParentCategory(parentCategoryID)
	err = s.categoryRepo.Save(ctx, category)
	if err != nil {
		return fmt.Errorf("failed to set parent category: %w", err)
	}
	return nil
}

func (s *categoryService) RemoveDefaultParentCategory(ctx context.Context, categoryID int64) error {
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("failed to find category by ID: %w", err)
	}
	if category == nil {
		return fmt.Errorf("category with ID %d not found", categoryID)
	}

	category.RemoveDefaultParentCategory()
	err = s.categoryRepo.Save(ctx, category)
	if err != nil {
		return fmt.Errorf("failed to remove default parent category: %w", err)
	}
	return nil
}

func (s *categoryService) AddCategoryAttribute(ctx context.Context, categoryID int64, name, value string) (*CategoryAttributeDTO, error) {
	attribute, err := domain.NewCategoryAttribute(categoryID, name, value)
	if err != nil {
		return nil, fmt.Errorf("failed to create new category attribute: %w", err)
	}

	err = s.categoryAttributeRepo.Save(ctx, attribute)
	if err != nil {
		return nil, fmt.Errorf("failed to save category attribute: %w", err)
	}

	return toCategoryAttributeDTO(attribute), nil
}

func (s *categoryService) UpdateCategoryAttribute(ctx context.Context, categoryAttributeID int64, name, value string) (*CategoryAttributeDTO, error) {
	attribute, err := s.categoryAttributeRepo.FindByID(ctx, categoryAttributeID)
	if err != nil {
		return nil, fmt.Errorf("failed to find category attribute by ID for update: %w", err)
	}
	if attribute == nil {
		return nil, fmt.Errorf("category attribute with ID %d not found for update", categoryAttributeID)
	}

	attribute.UpdateValue(value)
	err = s.categoryAttributeRepo.Save(ctx, attribute)
	if err != nil {
		return nil, fmt.Errorf("failed to update category attribute: %w", err)
	}

	return toCategoryAttributeDTO(attribute), nil
}

func (s *categoryService) RemoveCategoryAttribute(ctx context.Context, categoryAttributeID int64) error {
	err := s.categoryAttributeRepo.Delete(ctx, categoryAttributeID)
	if err != nil {
		return fmt.Errorf("failed to remove category attribute: %w", err)
	}
	return nil
}

func toCategoryDTO(category *domain.Category) *CategoryDTO {
	return &CategoryDTO{
		ID:                      category.ID,
		Name:                    category.Name,
		Description:             category.Description,
		LongDescription:         category.LongDescription,
		ActiveStartDate:         category.ActiveStartDate,
		ActiveEndDate:           category.ActiveEndDate,
		Archived:                category.Archived,
		DisplayTemplate:         category.DisplayTemplate,
		ExternalID:              category.ExternalID,
		FulfillmentType:         category.FulfillmentType,
		InventoryType:           category.InventoryType,
		MetaDescription:         category.MetaDescription,
		MetaTitle:               category.MetaTitle,
		OverrideGeneratedURL:    category.OverrideGeneratedURL,
		ProductDescPattern:      category.ProductDescPattern,
		ProductTitlePattern:     category.ProductTitlePattern,
		RootDisplayOrder:        category.RootDisplayOrder,
		TaxCode:                 category.TaxCode,
		URL:                     category.URL,
		URLKey:                  category.URLKey,
		DefaultParentCategoryID: category.DefaultParentCategoryID,
		CreatedAt:               category.CreatedAt,
		UpdatedAt:               category.UpdatedAt,
	}
}

func toCategoryAttributeDTO(attribute *domain.CategoryAttribute) *CategoryAttributeDTO {
	return &CategoryAttributeDTO{
		ID:         attribute.ID,
		Name:       attribute.Name,
		Value:      attribute.Value,
		CategoryID: attribute.CategoryID,
	}
}
