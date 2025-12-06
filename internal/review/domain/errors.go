package domain

import "errors"

var (
	ErrProductIDRequired      = errors.New("product ID is required")
	ErrCustomerIDRequired     = errors.New("customer ID is required")
	ErrInvalidRating          = errors.New("rating must be between 1 and 5")
	ErrCommentRequired        = errors.New("comment is required")
	ErrReviewNotFound         = errors.New("review not found")
	ErrReviewAlreadyApproved  = errors.New("review is already approved")
	ErrCustomerCannotReview   = errors.New("customer cannot review this product")
	ErrDuplicateReview        = errors.New("customer has already reviewed this product")
)
