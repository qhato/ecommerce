package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReviewStatus represents the status of a review
type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "PENDING"
	ReviewStatusApproved ReviewStatus = "APPROVED"
	ReviewStatusRejected ReviewStatus = "REJECTED"
	ReviewStatusFlagged  ReviewStatus = "FLAGGED"
)

// Review represents a product review
type Review struct {
	ID              string
	ProductID       string
	CustomerID      string
	CustomerName    string
	OrderID         *string
	Rating          int // 1-5
	Title           string
	Comment         string
	Status          ReviewStatus
	IsVerifiedBuyer bool
	HelpfulCount    int
	NotHelpfulCount int
	ReviewerEmail   string
	ResponseText    *string // Store/admin response
	ResponseDate    *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewReview creates a new review
func NewReview(productID, customerID, customerName, reviewerEmail string, rating int, title, comment string) (*Review, error) {
	if productID == "" {
		return nil, ErrProductIDRequired
	}
	if customerID == "" {
		return nil, ErrCustomerIDRequired
	}
	if rating < 1 || rating > 5 {
		return nil, ErrInvalidRating
	}
	if comment == "" {
		return nil, ErrCommentRequired
	}

	now := time.Now()
	return &Review{
		ID:              uuid.New().String(),
		ProductID:       productID,
		CustomerID:      customerID,
		CustomerName:    customerName,
		ReviewerEmail:   reviewerEmail,
		Rating:          rating,
		Title:           title,
		Comment:         comment,
		Status:          ReviewStatusPending,
		IsVerifiedBuyer: false,
		HelpfulCount:    0,
		NotHelpfulCount: 0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Approve approves the review
func (r *Review) Approve() error {
	if r.Status == ReviewStatusApproved {
		return ErrReviewAlreadyApproved
	}
	r.Status = ReviewStatusApproved
	r.UpdatedAt = time.Now()
	return nil
}

// Reject rejects the review
func (r *Review) Reject() {
	r.Status = ReviewStatusRejected
	r.UpdatedAt = time.Now()
}

// Flag flags the review for moderation
func (r *Review) Flag() {
	r.Status = ReviewStatusFlagged
	r.UpdatedAt = time.Now()
}

// AddResponse adds a store/admin response
func (r *Review) AddResponse(responseText string) {
	r.ResponseText = &responseText
	now := time.Now()
	r.ResponseDate = &now
	r.UpdatedAt = now
}

// MarkHelpful increments helpful count
func (r *Review) MarkHelpful() {
	r.HelpfulCount++
	r.UpdatedAt = time.Now()
}

// MarkNotHelpful increments not helpful count
func (r *Review) MarkNotHelpful() {
	r.NotHelpfulCount++
	r.UpdatedAt = time.Now()
}

// SetVerifiedBuyer marks the review as from a verified buyer
func (r *Review) SetVerifiedBuyer(orderID string) {
	r.IsVerifiedBuyer = true
	r.OrderID = &orderID
	r.UpdatedAt = time.Now()
}

// UpdateReview updates review content
func (r *Review) UpdateReview(title, comment string, rating int) error {
	if rating < 1 || rating > 5 {
		return ErrInvalidRating
	}
	if comment == "" {
		return ErrCommentRequired
	}

	r.Title = title
	r.Comment = comment
	r.Rating = rating
	r.UpdatedAt = time.Now()
	return nil
}
