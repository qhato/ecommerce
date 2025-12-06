package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/catalog/domain"
	"github.com/shopspring/decimal"
)

type ProductBundleDTO struct {
	ID          int64                    `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	BundlePrice decimal.Decimal          `json:"bundle_price"`
	IsActive    bool                     `json:"is_active"`
	Items       []ProductBundleItemDTO   `json:"items"`
}

type ProductBundleItemDTO struct {
	ID        int64  `json:"id"`
	ProductID *int64 `json:"product_id,omitempty"`
	SKUID     *int64 `json:"sku_id,omitempty"`
	Quantity  int    `json:"quantity"`
	SortOrder int    `json:"sort_order"`
}

type ProductBundleQueryService struct {
	bundleRepo     domain.ProductBundleRepository
	bundleItemRepo domain.ProductBundleItemRepository
}

func NewProductBundleQueryService(
	bundleRepo domain.ProductBundleRepository,
	bundleItemRepo domain.ProductBundleItemRepository,
) *ProductBundleQueryService {
	return &ProductBundleQueryService{
		bundleRepo:     bundleRepo,
		bundleItemRepo: bundleItemRepo,
	}
}

func (s *ProductBundleQueryService) GetBundle(ctx context.Context, id int64) (*ProductBundleDTO, error) {
	bundle, err := s.bundleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find bundle: %w", err)
	}
	if bundle == nil {
		return nil, domain.ErrBundleNotFound
	}

	items, err := s.bundleItemRepo.FindByBundleID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find bundle items: %w", err)
	}

	return toBundleDTO(bundle, items), nil
}

func (s *ProductBundleQueryService) GetAllBundles(ctx context.Context, activeOnly bool) ([]*ProductBundleDTO, error) {
	bundles, err := s.bundleRepo.FindAll(ctx, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find bundles: %w", err)
	}

	dtos := make([]*ProductBundleDTO, len(bundles))
	for i, bundle := range bundles {
		items, err := s.bundleItemRepo.FindByBundleID(ctx, bundle.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to find bundle items: %w", err)
		}
		dtos[i] = toBundleDTO(bundle, items)
	}

	return dtos, nil
}

func toBundleDTO(bundle *domain.ProductBundle, items []*domain.ProductBundleItem) *ProductBundleDTO {
	itemDTOs := make([]ProductBundleItemDTO, len(items))
	for i, item := range items {
		itemDTOs[i] = ProductBundleItemDTO{
			ID:        item.ID,
			ProductID: item.ProductID,
			SKUID:     item.SKUID,
			Quantity:  item.Quantity,
			SortOrder: item.SortOrder,
		}
	}

	return &ProductBundleDTO{
		ID:          bundle.ID,
		Name:        bundle.Name,
		Description: bundle.Description,
		BundlePrice: bundle.BundlePrice,
		IsActive:    bundle.IsActive,
		Items:       itemDTOs,
	}
}
