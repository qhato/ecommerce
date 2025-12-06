package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

type ProductRelationshipCommandHandler struct {
	relationshipRepo domain.ProductRelationshipRepository
}

func NewProductRelationshipCommandHandler(
	relationshipRepo domain.ProductRelationshipRepository,
) *ProductRelationshipCommandHandler {
	return &ProductRelationshipCommandHandler{
		relationshipRepo: relationshipRepo,
	}
}

func (h *ProductRelationshipCommandHandler) HandleCreateProductRelationship(ctx context.Context, cmd CreateProductRelationshipCommand) (*domain.ProductRelationship, error) {
	if cmd.ProductID == cmd.RelatedProductID {
		return nil, domain.ErrSelfRelationship
	}

	relType := domain.ProductRelationshipType(cmd.RelationshipType)
	exists, err := h.relationshipRepo.ExistsByProducts(ctx, cmd.ProductID, cmd.RelatedProductID, relType)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing relationship: %w", err)
	}
	if exists {
		return nil, domain.ErrDuplicateRelationship
	}

	relationship := domain.NewProductRelationship(cmd.ProductID, cmd.RelatedProductID, relType)
	relationship.Sequence = cmd.Sequence

	if err := h.relationshipRepo.Create(ctx, relationship); err != nil {
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}

	return relationship, nil
}

func (h *ProductRelationshipCommandHandler) HandleUpdateRelationshipSequence(ctx context.Context, cmd UpdateRelationshipSequenceCommand) error {
	relationship, err := h.relationshipRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find relationship: %w", err)
	}
	if relationship == nil {
		return domain.ErrRelationshipNotFound
	}

	relationship.UpdateSequence(cmd.Sequence)
	return h.relationshipRepo.Update(ctx, relationship)
}

func (h *ProductRelationshipCommandHandler) HandleActivateRelationship(ctx context.Context, cmd ActivateRelationshipCommand) error {
	relationship, err := h.relationshipRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find relationship: %w", err)
	}
	if relationship == nil {
		return domain.ErrRelationshipNotFound
	}

	relationship.Activate()
	return h.relationshipRepo.Update(ctx, relationship)
}

func (h *ProductRelationshipCommandHandler) HandleDeactivateRelationship(ctx context.Context, cmd DeactivateRelationshipCommand) error {
	relationship, err := h.relationshipRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find relationship: %w", err)
	}
	if relationship == nil {
		return domain.ErrRelationshipNotFound
	}

	relationship.Deactivate()
	return h.relationshipRepo.Update(ctx, relationship)
}

func (h *ProductRelationshipCommandHandler) HandleDeleteRelationship(ctx context.Context, cmd DeleteRelationshipCommand) error {
	return h.relationshipRepo.Delete(ctx, cmd.ID)
}
