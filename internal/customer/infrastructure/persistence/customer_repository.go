package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/customer/domain"
	"github.com/qhato/ecommerce/pkg/apperrors"
)

// PostgresCustomerRepository implements the CustomerRepository interface using PostgreSQL
type PostgresCustomerRepository struct {
	db *sql.DB
}

// NewPostgresCustomerRepository creates a new PostgresCustomerRepository
func NewPostgresCustomerRepository(db *sql.DB) *PostgresCustomerRepository {
	return &PostgresCustomerRepository{db: db}
}

// Create creates a new customer
func (r *PostgresCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO blc_customer (
			archived, challenge_answer, deactivated, email_address, external_id,
			first_name, is_tax_exempt, last_name, password, password_change_required,
			is_preview, receive_email, is_registered, tax_exemption_code, user_name,
			challenge_question_id, locale_code, created_by, date_created, date_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING customer_id
	`

	err := r.db.QueryRowContext(ctx, query,
		customer.Archived,
		customer.ChallengeAnswer,
		customer.Deactivated,
		customer.EmailAddress,
		customer.ExternalID,
		customer.FirstName,
		customer.IsTaxExempt,
		customer.LastName,
		customer.Password,
		customer.PasswordChangeRequired,
		customer.IsPreview,
		customer.ReceiveEmail,
		customer.IsRegistered,
		customer.TaxExemptionCode,
		customer.UserName,
		customer.ChallengeQuestionID,
		customer.LocaleCode,
		customer.CreatedBy,
		customer.CreatedAt,
		customer.UpdatedAt,
	).Scan(&customer.ID)

	if err != nil {
		return apperrors.NewInternalError("failed to create customer", err)
	}

	return nil
}

// Update updates an existing customer
func (r *PostgresCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	query := `
		UPDATE blc_customer
		SET archived = $1, challenge_answer = $2, deactivated = $3, email_address = $4,
			external_id = $5, first_name = $6, is_tax_exempt = $7, last_name = $8,
			password = $9, password_change_required = $10, is_preview = $11,
			receive_email = $12, is_registered = $13, tax_exemption_code = $14,
			user_name = $15, challenge_question_id = $16, locale_code = $17,
			updated_by = $18, date_updated = $19
		WHERE customer_id = $20
	`

	result, err := r.db.ExecContext(ctx, query,
		customer.Archived,
		customer.ChallengeAnswer,
		customer.Deactivated,
		customer.EmailAddress,
		customer.ExternalID,
		customer.FirstName,
		customer.IsTaxExempt,
		customer.LastName,
		customer.Password,
		customer.PasswordChangeRequired,
		customer.IsPreview,
		customer.ReceiveEmail,
		customer.IsRegistered,
		customer.TaxExemptionCode,
		customer.UserName,
		customer.ChallengeQuestionID,
		customer.LocaleCode,
		customer.UpdatedBy,
		customer.UpdatedAt,
		customer.ID,
	)

	if err != nil {
		return apperrors.NewInternalError("failed to update customer", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewInternalError("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return apperrors.NewNotFoundError("customer", customer.ID)
	}

	return nil
}

// FindByID finds a customer by ID
func (r *PostgresCustomerRepository) FindByID(ctx context.Context, id int64) (*domain.Customer, error) {
	query := `
		SELECT customer_id, archived, challenge_answer, deactivated, email_address,
			   external_id, first_name, is_tax_exempt, last_name, password,
			   password_change_required, is_preview, receive_email, is_registered,
			   tax_exemption_code, user_name, challenge_question_id, locale_code,
			   created_by, updated_by, date_created, date_updated
		FROM blc_customer
		WHERE customer_id = $1
	`

	customer := &domain.Customer{}
	var (
		challengeAnswer     sql.NullString
		externalID          sql.NullString
		taxExemptionCode    sql.NullString
		challengeQuestionID sql.NullInt64
		localeCode          sql.NullString
		createdBy           sql.NullInt64
		updatedBy           sql.NullInt64
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.ID,
		&customer.Archived,
		&challengeAnswer,
		&customer.Deactivated,
		&customer.EmailAddress,
		&externalID,
		&customer.FirstName,
		&customer.IsTaxExempt,
		&customer.LastName,
		&customer.Password,
		&customer.PasswordChangeRequired,
		&customer.IsPreview,
		&customer.ReceiveEmail,
		&customer.IsRegistered,
		&taxExemptionCode,
		&customer.UserName,
		&challengeQuestionID,
		&localeCode,
		&createdBy,
		&updatedBy,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, apperrors.NewInternalError("failed to find customer by ID", err)
	}

	// Handle nullable fields
	if challengeAnswer.Valid {
		customer.ChallengeAnswer = challengeAnswer.String
	}
	if externalID.Valid {
		customer.ExternalID = externalID.String
	}
	if taxExemptionCode.Valid {
		customer.TaxExemptionCode = taxExemptionCode.String
	}
	if challengeQuestionID.Valid {
		customer.ChallengeQuestionID = &challengeQuestionID.Int64
	}
	if localeCode.Valid {
		customer.LocaleCode = localeCode.String
	}
	if createdBy.Valid {
		customer.CreatedBy = createdBy.Int64
	}
	if updatedBy.Valid {
		customer.UpdatedBy = updatedBy.Int64
	}

	// Initialize slices
	customer.Addresses = make([]domain.CustomerAddress, 0)
	customer.Phones = make([]domain.CustomerPhone, 0)
	customer.Attributes = make([]domain.CustomerAttribute, 0)
	customer.Roles = make([]domain.CustomerRole, 0)

	return customer, nil
}

// FindByEmail finds a customer by email address
func (r *PostgresCustomerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	query := `
		SELECT customer_id, archived, challenge_answer, deactivated, email_address,
			   external_id, first_name, is_tax_exempt, last_name, password,
			   password_change_required, is_preview, receive_email, is_registered,
			   tax_exemption_code, user_name, challenge_question_id, locale_code,
			   created_by, updated_by, date_created, date_updated
		FROM blc_customer
		WHERE email_address = $1
	`

	customer := &domain.Customer{}
	var (
		challengeAnswer     sql.NullString
		externalID          sql.NullString
		taxExemptionCode    sql.NullString
		challengeQuestionID sql.NullInt64
		localeCode          sql.NullString
		createdBy           sql.NullInt64
		updatedBy           sql.NullInt64
	)

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&customer.ID,
		&customer.Archived,
		&challengeAnswer,
		&customer.Deactivated,
		&customer.EmailAddress,
		&externalID,
		&customer.FirstName,
		&customer.IsTaxExempt,
		&customer.LastName,
		&customer.Password,
		&customer.PasswordChangeRequired,
		&customer.IsPreview,
		&customer.ReceiveEmail,
		&customer.IsRegistered,
		&taxExemptionCode,
		&customer.UserName,
		&challengeQuestionID,
		&localeCode,
		&createdBy,
		&updatedBy,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, apperrors.NewInternalError("failed to find customer by email", err)
	}

	// Handle nullable fields
	if challengeAnswer.Valid {
		customer.ChallengeAnswer = challengeAnswer.String
	}
	if externalID.Valid {
		customer.ExternalID = externalID.String
	}
	if taxExemptionCode.Valid {
		customer.TaxExemptionCode = taxExemptionCode.String
	}
	if challengeQuestionID.Valid {
		customer.ChallengeQuestionID = &challengeQuestionID.Int64
	}
	if localeCode.Valid {
		customer.LocaleCode = localeCode.String
	}
	if createdBy.Valid {
		customer.CreatedBy = createdBy.Int64
	}
	if updatedBy.Valid {
		customer.UpdatedBy = updatedBy.Int64
	}

	// Initialize slices
	customer.Addresses = make([]domain.CustomerAddress, 0)
	customer.Phones = make([]domain.CustomerPhone, 0)
	customer.Attributes = make([]domain.CustomerAttribute, 0)
	customer.Roles = make([]domain.CustomerRole, 0)

	return customer, nil
}

// FindByUsername finds a customer by username
func (r *PostgresCustomerRepository) FindByUsername(ctx context.Context, username string) (*domain.Customer, error) {
	query := `
		SELECT customer_id, archived, challenge_answer, deactivated, email_address,
			   external_id, first_name, is_tax_exempt, last_name, password,
			   password_change_required, is_preview, receive_email, is_registered,
			   tax_exemption_code, user_name, challenge_question_id, locale_code,
			   created_by, updated_by, date_created, date_updated
		FROM blc_customer
		WHERE user_name = $1
	`

	customer := &domain.Customer{}
	var (
		challengeAnswer     sql.NullString
		externalID          sql.NullString
		taxExemptionCode    sql.NullString
		challengeQuestionID sql.NullInt64
		localeCode          sql.NullString
		createdBy           sql.NullInt64
		updatedBy           sql.NullInt64
	)

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&customer.ID,
		&customer.Archived,
		&challengeAnswer,
		&customer.Deactivated,
		&customer.EmailAddress,
		&externalID,
		&customer.FirstName,
		&customer.IsTaxExempt,
		&customer.LastName,
		&customer.Password,
		&customer.PasswordChangeRequired,
		&customer.IsPreview,
		&customer.ReceiveEmail,
		&customer.IsRegistered,
		&taxExemptionCode,
		&customer.UserName,
		&challengeQuestionID,
		&localeCode,
		&createdBy,
		&updatedBy,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, apperrors.NewInternalError("failed to find customer by username", err)
	}

	// Handle nullable fields
	if challengeAnswer.Valid {
		customer.ChallengeAnswer = challengeAnswer.String
	}
	if externalID.Valid {
		customer.ExternalID = externalID.String
	}
	if taxExemptionCode.Valid {
		customer.TaxExemptionCode = taxExemptionCode.String
	}
	if challengeQuestionID.Valid {
		customer.ChallengeQuestionID = &challengeQuestionID.Int64
	}
	if localeCode.Valid {
		customer.LocaleCode = localeCode.String
	}
	if createdBy.Valid {
		customer.CreatedBy = createdBy.Int64
	}
	if updatedBy.Valid {
		customer.UpdatedBy = updatedBy.Int64
	}

	// Initialize slices
	customer.Addresses = make([]domain.CustomerAddress, 0)
	customer.Phones = make([]domain.CustomerPhone, 0)
	customer.Attributes = make([]domain.CustomerAttribute, 0)
	customer.Roles = make([]domain.CustomerRole, 0)

	return customer, nil
}

// FindAll finds all customers
func (r *PostgresCustomerRepository) FindAll(ctx context.Context, filter *domain.CustomerFilter) ([]*domain.Customer, int64, error) {
	query := `
		SELECT customer_id, archived, challenge_answer, deactivated, email_address,
			   external_id, first_name, is_tax_exempt, last_name, password,
			   password_change_required, is_preview, receive_email, is_registered,
			   tax_exemption_code, user_name, challenge_question_id, locale_code,
			   created_by, updated_by, date_created, date_updated
		FROM blc_customer
		WHERE 1=1
	`

	args := make([]interface{}, 0)
	argIndex := 1

	// Add filters
	if filter != nil {
		if filter.Deactivated != nil {
			query += fmt.Sprintf(" AND deactivated = $%d", argIndex)
			args = append(args, *filter.Deactivated)
			argIndex++
		}
		if filter.Archived != nil {
			query += fmt.Sprintf(" AND archived = $%d", argIndex)
			args = append(args, *filter.Archived)
			argIndex++
		}
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM blc_customer WHERE 1=1"
	countArgs := make([]interface{}, 0)
	countArgIndex := 1
	if filter != nil {
		if filter.Deactivated != nil {
			countQuery += fmt.Sprintf(" AND deactivated = $%d", countArgIndex)
			countArgs = append(countArgs, *filter.Deactivated)
			countArgIndex++
		}
		if filter.Archived != nil {
			countQuery += fmt.Sprintf(" AND archived = $%d", countArgIndex)
			countArgs = append(countArgs, *filter.Archived)
		}
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("failed to count customers", err)
	}

	// Add sorting
	if filter != nil && filter.SortBy != "" {
		sortOrder := "ASC"
		if filter.SortOrder == "DESC" {
			sortOrder = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", filter.SortBy, sortOrder)
	} else {
		query += " ORDER BY date_created DESC"
	}

	// Add pagination
	if filter != nil && filter.PageSize > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("failed to find customers", err)
	}
	defer rows.Close()

	customers := make([]*domain.Customer, 0)
	for rows.Next() {
		customer := &domain.Customer{}
		var (
			challengeAnswer     sql.NullString
			externalID          sql.NullString
			taxExemptionCode    sql.NullString
			challengeQuestionID sql.NullInt64
			localeCode          sql.NullString
			createdBy           sql.NullInt64
			updatedBy           sql.NullInt64
		)

		err := rows.Scan(
			&customer.ID,
			&customer.Archived,
			&challengeAnswer,
			&customer.Deactivated,
			&customer.EmailAddress,
			&externalID,
			&customer.FirstName,
			&customer.IsTaxExempt,
			&customer.LastName,
			&customer.Password,
			&customer.PasswordChangeRequired,
			&customer.IsPreview,
			&customer.ReceiveEmail,
			&customer.IsRegistered,
			&taxExemptionCode,
			&customer.UserName,
			&challengeQuestionID,
			&localeCode,
			&createdBy,
			&updatedBy,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, 0, apperrors.NewInternalError("failed to scan customer", err)
		}

		// Handle nullable fields
		if challengeAnswer.Valid {
			customer.ChallengeAnswer = challengeAnswer.String
		}
		if externalID.Valid {
			customer.ExternalID = externalID.String
		}
		if taxExemptionCode.Valid {
			customer.TaxExemptionCode = taxExemptionCode.String
		}
		if challengeQuestionID.Valid {
			customer.ChallengeQuestionID = &challengeQuestionID.Int64
		}
		if localeCode.Valid {
			customer.LocaleCode = localeCode.String
		}
		if createdBy.Valid {
			customer.CreatedBy = createdBy.Int64
		}
		if updatedBy.Valid {
			customer.UpdatedBy = updatedBy.Int64
		}

		// Initialize slices
		customer.Addresses = make([]domain.CustomerAddress, 0)
		customer.Phones = make([]domain.CustomerPhone, 0)
		customer.Attributes = make([]domain.CustomerAttribute, 0)
		customer.Roles = make([]domain.CustomerRole, 0)

		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.NewInternalError("failed to iterate customers", err)
	}

	return customers, total, nil
}

// ExistsByEmail checks if a customer exists by email
func (r *PostgresCustomerRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM blc_customer WHERE email_address = $1)"
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, apperrors.NewInternalError("failed to check customer by email", err)
	}
	return exists, nil
}

// ExistsByUsername checks if a customer exists by username
func (r *PostgresCustomerRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM blc_customer WHERE user_name = $1)"
	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, apperrors.NewInternalError("failed to check customer by username", err)
	}
	return exists, nil
}
