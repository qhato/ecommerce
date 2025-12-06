package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

type ProductRelationshipDTO struct {
	ID               int64  `json:"id"`
	ProductID        int64  `json:"product_id"`
	RelatedProductID int64  `json:"related_product_id"`
	RelationshipType string `json:"relationship_type"`
	Sequence         int    `json:"sequence"`
	IsActive         bool   `json:"is_active"`
}

type ProductRelationshipQueryService struct {
	relationshipRepo domain.ProductRelationshipRepository
}

func NewProductRelationshipQueryService(
	relationshipRepo domain.ProductRelationshipRepository,
) *ProductRelationshipQueryService {
	return &ProductRelationshipQueryService{
		relationshipRepo: relationshipRepo,
	}
}

func (s *ProductRelationshipQueryService) GetCrossSellProducts(ctx context.Context, productID int64) ([]*ProductRelationshipDTO, error) {
	relationships, err := s.relationshipRepo.FindCrossSell(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find cross-sell products: %w", err)
	}

	return toRelationshipDTOs(relationships), nil
}

func (s *ProductRelationshipQueryService) GetUpSellProducts(ctx context.Context, productID int64) ([]*ProductRelationshipDTO, error) {
	relationships, err := s.relationshipRepo.FindUpSell(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find up-sell products: %w", err)
	}

	return toRelationshipDTOs(relationships), nil
}

func (s *ProductRelationshipQueryService) GetRelatedProducts(ctx context.Context, productID int64) ([]*ProductRelationshipDTO, error) {
	relationships, err := s.relationshipRepo.FindRelated(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find related products: %w", err)
	}

	return toRelationshipDTOs(relationships), nil
}

func toRelationshipDTOs(relationships []*domain.ProductRelationship) []*ProductRelationshipDTO {
	dtos := make([]*ProductRelationshipDTO, len(relationships))
	for i, rel := range relationships {
		dtos[i] = &ProductRelationshipDTO{
			ID:               rel.ID,
			ProductID:        rel.ProductID,
			RelatedProductID: rel.RelatedProductID,
			RelationshipType: string(rel.RelationshipType),
			Sequence:         rel.Sequence,
			IsActive:         rel.IsActive,
		}
	}
	return dtos
}
