package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/review/domain"
)

type ReviewCommandHandler struct {
	reviewRepo domain.ReviewRepository
}

func NewReviewCommandHandler(reviewRepo domain.ReviewRepository) *ReviewCommandHandler {
	return &ReviewCommandHandler{reviewRepo: reviewRepo}
}

func (h *ReviewCommandHandler) HandleCreateReview(ctx context.Context, cmd CreateReviewCommand) (*domain.Review, error) {
	// Check for duplicate review
	exists, err := h.reviewRepo.ExistsByCustomerAndProduct(ctx, cmd.CustomerID, cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate review: %w", err)
	}
	if exists {
		return nil, domain.ErrDuplicateReview
	}

	review, err := domain.NewReview(
		cmd.ProductID,
		cmd.CustomerID,
		cmd.CustomerName,
		cmd.ReviewerEmail,
		cmd.Rating,
		cmd.Title,
		cmd.Comment,
	)
	if err != nil {
		return nil, err
	}

	if cmd.OrderID != nil {
		review.SetVerifiedBuyer(*cmd.OrderID)
	}

	if err := h.reviewRepo.Create(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}

func (h *ReviewCommandHandler) HandleUpdateReview(ctx context.Context, cmd UpdateReviewCommand) (*domain.Review, error) {
	review, err := h.reviewRepo.FindByID(ctx, cmd.ID)
	if err != nil || review == nil {
		return nil, domain.ErrReviewNotFound
	}

	if err := review.UpdateReview(cmd.Title, cmd.Comment, cmd.Rating); err != nil {
		return nil, err
	}

	if err := h.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	return review, nil
}

func (h *ReviewCommandHandler) HandleApproveReview(ctx context.Context, cmd ApproveReviewCommand) (*domain.Review, error) {
	review, err := h.reviewRepo.FindByID(ctx, cmd.ID)
	if err != nil || review == nil {
		return nil, domain.ErrReviewNotFound
	}

	if err := review.Approve(); err != nil {
		return nil, err
	}

	if err := h.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to approve review: %w", err)
	}

	return review, nil
}

func (h *ReviewCommandHandler) HandleRejectReview(ctx context.Context, cmd RejectReviewCommand) (*domain.Review, error) {
	review, err := h.reviewRepo.FindByID(ctx, cmd.ID)
	if err != nil || review == nil {
		return nil, domain.ErrReviewNotFound
	}

	review.Reject()

	if err := h.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to reject review: %w", err)
	}

	return review, nil
}

func (h *ReviewCommandHandler) HandleFlagReview(ctx context.Context, cmd FlagReviewCommand) (*domain.Review, error) {
	review, err := h.reviewRepo.FindByID(ctx, cmd.ID)
	if err != nil || review == nil {
		return nil, domain.ErrReviewNotFound
	}

	review.Flag()

	if err := h.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to flag review: %w", err)
	}

	return review, nil
}

func (h *ReviewCommandHandler) HandleAddResponse(ctx context.Context, cmd AddResponseCommand) (*domain.Review, error) {
	review, err := h.reviewRepo.FindByID(ctx, cmd.ID)
	if err != nil || review == nil {
		return nil, domain.ErrReviewNotFound
	}

	review.AddResponse(cmd.ResponseText)

	if err := h.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to add response: %w", err)
	}

	return review, nil
}

func (h *ReviewCommandHandler) HandleMarkHelpful(ctx context.Context, cmd MarkHelpfulCommand) (*domain.Review, error) {
	review, err := h.reviewRepo.FindByID(ctx, cmd.ID)
	if err != nil || review == nil {
		return nil, domain.ErrReviewNotFound
	}

	review.MarkHelpful()

	if err := h.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to mark helpful: %w", err)
	}

	return review, nil
}

func (h *ReviewCommandHandler) HandleMarkNotHelpful(ctx context.Context, cmd MarkNotHelpfulCommand) (*domain.Review, error) {
	review, err := h.reviewRepo.FindByID(ctx, cmd.ID)
	if err != nil || review == nil {
		return nil, domain.ErrReviewNotFound
	}

	review.MarkNotHelpful()

	if err := h.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to mark not helpful: %w", err)
	}

	return review, nil
}

func (h *ReviewCommandHandler) HandleDeleteReview(ctx context.Context, cmd DeleteReviewCommand) error {
	return h.reviewRepo.Delete(ctx, cmd.ID)
}
