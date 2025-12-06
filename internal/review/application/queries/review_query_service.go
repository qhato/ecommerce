package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/review/domain"
)

type ReviewQueryService struct {
	reviewRepo domain.ReviewRepository
}

func NewReviewQueryService(reviewRepo domain.ReviewRepository) *ReviewQueryService {
	return &ReviewQueryService{reviewRepo: reviewRepo}
}

func (s *ReviewQueryService) GetReview(ctx context.Context, id string) (*ReviewDTO, error) {
	review, err := s.reviewRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find review: %w", err)
	}
	if review == nil {
		return nil, domain.ErrReviewNotFound
	}

	return ToReviewDTO(review), nil
}

func (s *ReviewQueryService) GetProductReviews(ctx context.Context, productID string, status *string, limit, offset int) ([]*ReviewDTO, error) {
	var reviewStatus *domain.ReviewStatus
	if status != nil {
		rs := domain.ReviewStatus(*status)
		reviewStatus = &rs
	}

	reviews, err := s.reviewRepo.FindByProductID(ctx, productID, reviewStatus, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find product reviews: %w", err)
	}

	dtos := make([]*ReviewDTO, len(reviews))
	for i, review := range reviews {
		dtos[i] = ToReviewDTO(review)
	}

	return dtos, nil
}

func (s *ReviewQueryService) GetCustomerReviews(ctx context.Context, customerID string) ([]*ReviewDTO, error) {
	reviews, err := s.reviewRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find customer reviews: %w", err)
	}

	dtos := make([]*ReviewDTO, len(reviews))
	for i, review := range reviews {
		dtos[i] = ToReviewDTO(review)
	}

	return dtos, nil
}

func (s *ReviewQueryService) GetPendingReviews(ctx context.Context) ([]*ReviewDTO, error) {
	reviews, err := s.reviewRepo.FindByStatus(ctx, domain.ReviewStatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to find pending reviews: %w", err)
	}

	dtos := make([]*ReviewDTO, len(reviews))
	for i, review := range reviews {
		dtos[i] = ToReviewDTO(review)
	}

	return dtos, nil
}

func (s *ReviewQueryService) GetProductRatingSummary(ctx context.Context, productID string) (*RatingSummaryDTO, error) {
	avgRating, err := s.reviewRepo.GetAverageRating(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get average rating: %w", err)
	}

	distribution, err := s.reviewRepo.GetRatingDistribution(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating distribution: %w", err)
	}

	count, err := s.reviewRepo.CountByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to count reviews: %w", err)
	}

	return &RatingSummaryDTO{
		ProductID:         productID,
		AverageRating:     avgRating,
		TotalReviews:      count,
		RatingDistribution: distribution,
	}, nil
}
