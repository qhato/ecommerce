package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/review/domain"
)

type PostgresReviewRepository struct {
	db *sql.DB
}

func NewPostgresReviewRepository(db *sql.DB) *PostgresReviewRepository {
	return &PostgresReviewRepository{db: db}
}

func (r *PostgresReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	query := `
		INSERT INTO blc_review (id, product_id, customer_id, customer_name, order_id, rating,
			title, comment, status, is_verified_buyer, helpful_count, not_helpful_count,
			reviewer_email, response_text, response_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	_, err := r.db.ExecContext(ctx, query,
		review.ID, review.ProductID, review.CustomerID, review.CustomerName, review.OrderID,
		review.Rating, review.Title, review.Comment, review.Status, review.IsVerifiedBuyer,
		review.HelpfulCount, review.NotHelpfulCount, review.ReviewerEmail, review.ResponseText,
		review.ResponseDate, review.CreatedAt, review.UpdatedAt,
	)
	return err
}

func (r *PostgresReviewRepository) Update(ctx context.Context, review *domain.Review) error {
	query := `
		UPDATE blc_review
		SET title = $1, comment = $2, rating = $3, status = $4, is_verified_buyer = $5,
		    helpful_count = $6, not_helpful_count = $7, response_text = $8, response_date = $9,
		    updated_at = $10
		WHERE id = $11`

	_, err := r.db.ExecContext(ctx, query,
		review.Title, review.Comment, review.Rating, review.Status, review.IsVerifiedBuyer,
		review.HelpfulCount, review.NotHelpfulCount, review.ResponseText, review.ResponseDate,
		review.UpdatedAt, review.ID,
	)
	return err
}

func (r *PostgresReviewRepository) FindByID(ctx context.Context, id string) (*domain.Review, error) {
	query := `
		SELECT id, product_id, customer_id, customer_name, order_id, rating, title, comment,
		       status, is_verified_buyer, helpful_count, not_helpful_count, reviewer_email,
		       response_text, response_date, created_at, updated_at
		FROM blc_review WHERE id = $1`

	return r.scanReview(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresReviewRepository) FindByProductID(ctx context.Context, productID string, status *domain.ReviewStatus, limit, offset int) ([]*domain.Review, error) {
	var query string
	var args []interface{}

	if status != nil {
		query = `
			SELECT id, product_id, customer_id, customer_name, order_id, rating, title, comment,
			       status, is_verified_buyer, helpful_count, not_helpful_count, reviewer_email,
			       response_text, response_date, created_at, updated_at
			FROM blc_review
			WHERE product_id = $1 AND status = $2
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4`
		args = []interface{}{productID, *status, limit, offset}
	} else {
		query = `
			SELECT id, product_id, customer_id, customer_name, order_id, rating, title, comment,
			       status, is_verified_buyer, helpful_count, not_helpful_count, reviewer_email,
			       response_text, response_date, created_at, updated_at
			FROM blc_review
			WHERE product_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3`
		args = []interface{}{productID, limit, offset}
	}

	return r.queryReviews(ctx, query, args...)
}

func (r *PostgresReviewRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*domain.Review, error) {
	query := `
		SELECT id, product_id, customer_id, customer_name, order_id, rating, title, comment,
		       status, is_verified_buyer, helpful_count, not_helpful_count, reviewer_email,
		       response_text, response_date, created_at, updated_at
		FROM blc_review
		WHERE customer_id = $1
		ORDER BY created_at DESC`

	return r.queryReviews(ctx, query, customerID)
}

func (r *PostgresReviewRepository) FindByStatus(ctx context.Context, status domain.ReviewStatus) ([]*domain.Review, error) {
	query := `
		SELECT id, product_id, customer_id, customer_name, order_id, rating, title, comment,
		       status, is_verified_buyer, helpful_count, not_helpful_count, reviewer_email,
		       response_text, response_date, created_at, updated_at
		FROM blc_review
		WHERE status = $1
		ORDER BY created_at DESC`

	return r.queryReviews(ctx, query, status)
}

func (r *PostgresReviewRepository) GetAverageRating(ctx context.Context, productID string) (float64, error) {
	query := `SELECT COALESCE(AVG(rating), 0) FROM blc_review WHERE product_id = $1 AND status = 'APPROVED'`
	var avg float64
	err := r.db.QueryRowContext(ctx, query, productID).Scan(&avg)
	return avg, err
}

func (r *PostgresReviewRepository) GetRatingDistribution(ctx context.Context, productID string) (map[int]int, error) {
	query := `
		SELECT rating, COUNT(*) as count
		FROM blc_review
		WHERE product_id = $1 AND status = 'APPROVED'
		GROUP BY rating`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	distribution := make(map[int]int)
	for i := 1; i <= 5; i++ {
		distribution[i] = 0
	}

	for rows.Next() {
		var rating, count int
		if err := rows.Scan(&rating, &count); err != nil {
			return nil, err
		}
		distribution[rating] = count
	}

	return distribution, rows.Err()
}

func (r *PostgresReviewRepository) CountByProductID(ctx context.Context, productID string) (int64, error) {
	query := `SELECT COUNT(*) FROM blc_review WHERE product_id = $1 AND status = 'APPROVED'`
	var count int64
	err := r.db.QueryRowContext(ctx, query, productID).Scan(&count)
	return count, err
}

func (r *PostgresReviewRepository) ExistsByCustomerAndProduct(ctx context.Context, customerID, productID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blc_review WHERE customer_id = $1 AND product_id = $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, customerID, productID).Scan(&exists)
	return exists, err
}

func (r *PostgresReviewRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM blc_review WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresReviewRepository) scanReview(row interface {
	Scan(dest ...interface{}) error
}) (*domain.Review, error) {
	review := &domain.Review{}
	err := row.Scan(
		&review.ID, &review.ProductID, &review.CustomerID, &review.CustomerName,
		&review.OrderID, &review.Rating, &review.Title, &review.Comment,
		&review.Status, &review.IsVerifiedBuyer, &review.HelpfulCount, &review.NotHelpfulCount,
		&review.ReviewerEmail, &review.ResponseText, &review.ResponseDate,
		&review.CreatedAt, &review.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *PostgresReviewRepository) queryReviews(ctx context.Context, query string, args ...interface{}) ([]*domain.Review, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		review := &domain.Review{}
		if err := rows.Scan(
			&review.ID, &review.ProductID, &review.CustomerID, &review.CustomerName,
			&review.OrderID, &review.Rating, &review.Title, &review.Comment,
			&review.Status, &review.IsVerifiedBuyer, &review.HelpfulCount, &review.NotHelpfulCount,
			&review.ReviewerEmail, &review.ResponseText, &review.ResponseDate,
			&review.CreatedAt, &review.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, rows.Err()
}
