package application

import (
	"context"
	"fmt"

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
	CanonicalURL         string
	DisplayTemplate      string
	MetaDescription      string
	MetaTitle            string
	OverrideGeneratedURL bool
	DefaultCategoryID    *int64
	Attributes           map[string]string
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
	productRepo             domain.ProductRepository
	productAttributeRepo    domain.ProductAttributeRepository
	productOptionXrefRepo   domain.ProductOptionXrefRepository
	categoryProductXrefRepo domain.CategoryProductXrefRepository
}

// NewProductService creates a new instance of ProductService.
func NewProductService(
	productRepo domain.ProductRepository,
	productAttributeRepo domain.ProductAttributeRepository,
	productOptionXrefRepo domain.ProductOptionXrefRepository,
	categoryProductXrefRepo domain.CategoryProductXrefRepository,
) ProductService {
	return &productService{
		productRepo:             productRepo,
		productAttributeRepo:    productAttributeRepo,
		productOptionXrefRepo:   productOptionXrefRepo,
		categoryProductXrefRepo: categoryProductXrefRepo,
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

	product.CanonicalURL = cmd.CanonicalURL
	product.DisplayTemplate = cmd.DisplayTemplate
	product.MetaDescription = cmd.MetaDescription
	product.MetaTitle = cmd.MetaTitle
	product.OverrideGeneratedURL = cmd.OverrideGeneratedURL
	product.DefaultCategoryID = cmd.DefaultCategoryID

	err := s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to save product: %w", err)
	}
	
	if len(cmd.Attributes) > 0 {
		for name, value := range cmd.Attributes {
			attr, err := domain.NewProductAttribute(product.ID, name, value)
			if err != nil {
				// Log error but don't fail creation? Or fail? Failing is safer.
				return nil, fmt.Errorf("failed to create attribute: %w", err)
			}
			if err := s.productAttributeRepo.Save(ctx, attr); err != nil {
				return nil, fmt.Errorf("failed to save attribute: %w", err)
			}
		}
	}

	return ToProductDTO(product), nil
}

func (s *productService) GetProductByID(ctx context.Context, id int64) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find product by ID: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product with ID %d not found", id)
	}
	return ToProductDTO(product), nil
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

	err = s.productRepo.Update(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return ToProductDTO(product), nil
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
	err = s.productRepo.Update(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to archive product: %w", err)
	}
	return nil
}

func (s *productService) UnarchiveProduct(ctx context.Context, id int64) error {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find product by ID for unarchiving: %w", err)
	}
	if product == nil {
		return fmt.Errorf("product with ID %d not found for unarchiving", id)
	}

	product.Unarchive()
	err = s.productRepo.Update(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to unarchive product: %w", err)
	}
	return nil
}
