package application

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// ProductService defines the application service for product-related operations.
type ProductService interface {
	// CreateProduct creates a new product.
	CreateProduct(ctx context.Context, cmd *CreateProductCommand) (*ProductDTO, error)

	// GetProductByID retrieves a product by its ID.
	GetProductByID(ctx context.Context, id int64) (*ProductDTO, error)

	// UpdateProduct updates an existing product.
	UpdateProduct(ctx context.Context, cmd *UpdateProductCommand) (*ProductDTO, error)

	// ArchiveProduct archives a product, making it inactive.
	ArchiveProduct(ctx context.Context, id int64) error

	// UnarchiveProduct unarchives a product, making it active.
	UnarchiveProduct(ctx context.Context, id int64) error

	// AddProductAttribute adds a custom attribute to a product.
	AddProductAttribute(ctx context.Context, productID int64, name, value string) (*ProductAttributeDTO, error)

	// UpdateProductAttribute updates an existing product attribute.
	UpdateProductAttribute(ctx context.Context, productAttributeID int64, name, value string) (*ProductAttributeDTO, error)

	// RemoveProductAttribute removes a product attribute by its ID.
	RemoveProductAttribute(ctx context.Context, productAttributeID int64) error

	// AddProductOption adds a product option (via xref) to a product.
	AddProductOption(ctx context.Context, productID, productOptionID int64) (*ProductOptionXrefDTO, error)

	// RemoveProductOption removes a product option (via xref) from a product.
	RemoveProductOption(ctx context.Context, productID, productOptionID int64) error

	// AddCategoryToProduct associates a product with a category.
	AddCategoryToProduct(ctx context.Context, productID, categoryID int64, defaultReference bool, displayOrder float64) (*CategoryProductXrefDTO, error)

	// RemoveCategoryFromProduct disassociates a product from a category.
	RemoveCategoryFromProduct(ctx context.Context, productID, categoryID int64) error
}

// ProductDTO represents a product data transfer object.
type ProductDTO struct {
	ID                          int64
	Manufacture                 string
	Model                       string
	URL                         string
	URLKey                      string
	Archived                    bool
	MetaTitle                   string
	MetaDescription             string
	DefaultSkuID                *int64
	CanonicalURL                string
	DisplayTemplate             string
	EnableDefaultSKUInInventory bool
	CanSellWithoutOptions       bool
	OverrideGeneratedURL        bool
	DefaultCategoryID           *int64
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
}

// ProductAttributeDTO represents a product attribute data transfer object.
type ProductAttributeDTO struct {
	ID        int64
	Name      string
	Value     string
	ProductID int64
}

// ProductOptionXrefDTO represents a product option cross-reference data transfer object.
type ProductOptionXrefDTO struct {
	ID              int64
	ProductID       int64
	ProductOptionID int64
}

// CategoryProductXrefDTO represents a category product cross-reference data transfer object.
type CategoryProductXrefDTO struct {
	ID               int64
	CategoryID       int64
	ProductID        int64
	DefaultReference bool
	DisplayOrder     float64
}

// CreateProductCommand is a command to create a new product.
type CreateProductCommand struct {
	Manufacture                 string
	Model                       string
	URL                         string
	URLKey                      string
	CanSellWithoutOptions       bool
	EnableDefaultSKUInInventory bool
	// Other fields for initial product creation
}

// UpdateProductCommand is a command to update an existing product.
type UpdateProductCommand struct {
	ID                          int64
	Manufacture                 *string
	Model                       *string
	URL                         *string
	URLKey                      *string
	MetaTitle                   *string
	MetaDescription             *string
	DefaultSkuID                *int64
	CanonicalURL                *string
	DisplayTemplate             *string
	EnableDefaultSKUInInventory *bool
	CanSellWithoutOptions       *bool
	OverrideGeneratedURL        *bool
	DefaultCategoryID           *int64
}

type productService struct {
	productRepo              domain.ProductRepository
	productAttributeRepo     domain.ProductAttributeRepository
	productOptionXrefRepo    domain.ProductOptionXrefRepository
	categoryProductXrefRepo  domain.CategoryProductXrefRepository
}

// NewProductService creates a new instance of ProductService.
func NewProductService(
	productRepo domain.ProductRepository,
	productAttributeRepo domain.ProductAttributeRepository,
	productOptionXrefRepo domain.ProductOptionXrefRepository,
	categoryProductXrefRepo domain.CategoryProductXrefRepository,
) ProductService {
	return &productService{
		productRepo:              productRepo,
		productAttributeRepo:     productAttributeRepo,
		productOptionXrefRepo:    productOptionXrefRepo,
		categoryProductXrefRepo:  categoryProductXrefRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, cmd *CreateProductCommand) (*ProductDTO, error) {
	product := domain.NewProduct(
		cmd.Manufacture,
		cmd.Model,
		cmd.URL,
		cmd.URLKey,
		cmd.CanSellWithoutOptions,
		cmd.EnableDefaultSKUInInventory,
	)

	err := s.productRepo.Save(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to save product: %w", err)
	}

	return toProductDTO(product), nil
}

func (s *productService) GetProductByID(ctx context.Context, id int64) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find product by ID: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product with ID %d not found", id)
	}
	return toProductDTO(product), nil
}

func (s *productService) UpdateProduct(ctx context.Context, cmd *UpdateProductCommand) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product by ID for update: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product with ID %d not found for update", cmd.ID)
	}

	if cmd.Manufacture != nil {
		product.Manufacture = *cmd.Manufacture
	}
	if cmd.Model != nil {
		product.Model = *cmd.Model
	}
	if cmd.URL != nil {
		product.URL = *cmd.URL
	}
	if cmd.URLKey != nil {
		product.URLKey = *cmd.URLKey
	}
	if cmd.MetaTitle != nil {
		product.MetaTitle = *cmd.MetaTitle
	}
	if cmd.MetaDescription != nil {
		product.MetaDescription = *cmd.MetaDescription
	}
	if cmd.DefaultSkuID != nil {
		product.SetDefaultSKU(*cmd.DefaultSkuID)
	}
	if cmd.CanonicalURL != nil {
		product.CanonicalURL = *cmd.CanonicalURL
	}
	if cmd.DisplayTemplate != nil {
		product.DisplayTemplate = *cmd.DisplayTemplate
	}
	if cmd.EnableDefaultSKUInInventory != nil {
		product.EnableDefaultSKUInInventory = *cmd.EnableDefaultSKUInInventory
	}
	if cmd.CanSellWithoutOptions != nil {
		product.CanSellWithoutOptions = *cmd.CanSellWithoutOptions
	}
	if cmd.OverrideGeneratedURL != nil {
		product.OverrideGeneratedURL = *cmd.OverrideGeneratedURL
	}
	if cmd.DefaultCategoryID != nil {
		product.DefaultCategoryID = cmd.DefaultCategoryID
	}

	err = s.productRepo.Save(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return toProductDTO(product), nil
}

func (s *productService) ArchiveProduct(ctx context.Context, id int64) error {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find product by ID for archiving: %w", err)
	}
	if product == nil {
		return fmt.Errorf("product with ID %d not found for archiving", id)
	}

	product.Archive()
	err = s.productRepo.Save(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to archive product: %w", err)
	}
	return nil
}

func (s *productService) UnarchiveProduct(ctx context.Context, id int64) error {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find product by ID for unarchiving: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product with ID %d not found for unarchiving", id)
	}

	product.Unarchive()
	err = s.productRepo.Save(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to unarchive product: %w", err)
	}
	return nil
}

func (s *productService) AddProductAttribute(ctx context.Context, productID int64, name, value string) (*ProductAttributeDTO, error) {
	attribute, err := domain.NewProductAttribute(productID, name, value)
	if err != nil {
		return nil, fmt.Errorf("failed to create new product attribute: %w", err)
	}

	err = s.productAttributeRepo.Save(ctx, attribute)
	if err != nil {
		return nil, fmt.Errorf("failed to save product attribute: %w", err)
	}

	return toProductAttributeDTO(attribute), nil
}

func (s *productService) UpdateProductAttribute(ctx context.Context, productAttributeID int64, name, value string) (*ProductAttributeDTO, error) {
	attribute, err := s.productAttributeRepo.FindByID(ctx, productAttributeID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product attribute by ID for update: %w", err)
	}
	if attribute == nil {
		return nil, fmt.Errorf("product attribute with ID %d not found for update", productAttributeID)
	}

	attribute.UpdateValue(value)
	err = s.productAttributeRepo.Save(ctx, attribute)
	if err != nil {
		return nil, fmt.Errorf("failed to update product attribute: %w", err)
	}

	return toProductAttributeDTO(attribute), nil
}

func (s *productService) RemoveProductAttribute(ctx context.Context, productAttributeID int64) error {
	err := s.productAttributeRepo.Delete(ctx, productAttributeID)
	if err != nil {
		return fmt.Errorf("failed to remove product attribute: %w", err)
	}
	return nil
}

func (s *productService) AddProductOption(ctx context.Context, productID, productOptionID int64) (*ProductOptionXrefDTO, error) {
	xref, err := domain.NewProductOptionXref(productID, productOptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new product option xref: %w", err)
	}

	err = s.productOptionXrefRepo.Save(ctx, xref)
	if err != nil {
		return nil, fmt.Errorf("failed to save product option xref: %w", err)
	}

	return toProductOptionXrefDTO(xref), nil
}

func (s *productService) RemoveProductOption(ctx context.Context, productID, productOptionID int64) error {
	err := s.productOptionXrefRepo.RemoveProductOptionXref(ctx, productID, productOptionID)
	if err != nil {
		return fmt.Errorf("failed to remove product option xref: %w", err)
	}
	return nil
}

func (s *productService) AddCategoryToProduct(ctx context.Context, productID, categoryID int64, defaultReference bool, displayOrder float64) (*CategoryProductXrefDTO, error) {
	xref, err := domain.NewCategoryProductXref(categoryID, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new category product xref: %w", err)
	}
	xref.SetDefaultReference(defaultReference)
	xref.SetDisplayOrder(displayOrder)

	err = s.categoryProductXrefRepo.Save(ctx, xref)
	if err != nil {
		return nil, fmt.Errorf("failed to save category product xref: %w", err)
	}

	return toCategoryProductXrefDTO(xref), nil
}

func (s *productService) RemoveCategoryFromProduct(ctx context.Context, productID, categoryID int64) error {
	err := s.categoryProductXrefRepo.RemoveCategoryProductXref(ctx, categoryID, productID)
	if err != nil {
		return fmt.Errorf("failed to remove category product xref: %w", err)
	}
	return nil
}

func toProductDTO(product *domain.Product) *ProductDTO {
	return &ProductDTO{
		ID:                          product.ID,
		Manufacture:                 product.Manufacture,
		Model:                       product.Model,
		URL:                         product.URL,
		URLKey:                      product.URLKey,
		Archived:                    product.Archived,
		MetaTitle:                   product.MetaTitle,
		MetaDescription:             product.MetaDescription,
		DefaultSkuID:                product.DefaultSkuID,
		CanonicalURL:                product.CanonicalURL,
		DisplayTemplate:             product.DisplayTemplate,
		EnableDefaultSKUInInventory: product.EnableDefaultSKUInInventory,
		CanSellWithoutOptions:       product.CanSellWithoutOptions,
		OverrideGeneratedURL:        product.OverrideGeneratedURL,
		DefaultCategoryID:           product.DefaultCategoryID,
		CreatedAt:                   product.CreatedAt,
		UpdatedAt:                   product.UpdatedAt,
	}
}

func toProductAttributeDTO(attribute *domain.ProductAttribute) *ProductAttributeDTO {
	return &ProductAttributeDTO{
		ID:        attribute.ID,
		Name:      attribute.Name,
		Value:     attribute.Value,
		ProductID: attribute.ProductID,
	}
}

func toProductOptionXrefDTO(xref *domain.ProductOptionXref) *ProductOptionXrefDTO {
	return &ProductOptionXrefDTO{
		ID:              xref.ID,
		ProductID:       xref.ProductID,
		ProductOptionID: xref.ProductOptionID,
	}
}

func toCategoryProductXrefDTO(xref *domain.CategoryProductXref) *CategoryProductXrefDTO {
	return &CategoryProductXrefDTO{
		ID:               xref.ID,
		CategoryID:       xref.CategoryID,
		ProductID:        xref.ProductID,
		DefaultReference: xref.DefaultReference,
		DisplayOrder:     xref.DisplayOrder,
	}
}