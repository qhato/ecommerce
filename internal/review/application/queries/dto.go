package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/review/domain"
)

type ReviewDTO struct {
	ID              string     `json:"id"`
	ProductID       string     `json:"product_id"`
	CustomerID      string     `json:"customer_id"`
	CustomerName    string     `json:"customer_name"`
	OrderID         *string    `json:"order_id,omitempty"`
	Rating          int        `json:"rating"`
	Title           string     `json:"title"`
	Comment         string     `json:"comment"`
	Status          string     `json:"status"`
	IsVerifiedBuyer bool       `json:"is_verified_buyer"`
	HelpfulCount    int        `json:"helpful_count"`
	NotHelpfulCount int        `json:"not_helpful_count"`
	ReviewerEmail   string     `json:"reviewer_email"`
	ResponseText    *string    `json:"response_text,omitempty"`
	ResponseDate    *time.Time `json:"response_date,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type RatingSummaryDTO struct {
	ProductID          string         `json:"product_id"`
	AverageRating      float64        `json:"average_rating"`
	TotalReviews       int64          `json:"total_reviews"`
	RatingDistribution map[int]int    `json:"rating_distribution"`
}

func ToReviewDTO(review *domain.Review) *ReviewDTO {
	return &ReviewDTO{
		ID:              review.ID,
		ProductID:       review.ProductID,
		CustomerID:      review.CustomerID,
		CustomerName:    review.CustomerName,
		OrderID:         review.OrderID,
		Rating:          review.Rating,
		Title:           review.Title,
		Comment:         review.Comment,
		Status:          string(review.Status),
		IsVerifiedBuyer: review.IsVerifiedBuyer,
		HelpfulCount:    review.HelpfulCount,
		NotHelpfulCount: review.NotHelpfulCount,
		ReviewerEmail:   review.ReviewerEmail,
		ResponseText:    review.ResponseText,
		ResponseDate:    review.ResponseDate,
		CreatedAt:       review.CreatedAt,
		UpdatedAt:       review.UpdatedAt,
	}
}
