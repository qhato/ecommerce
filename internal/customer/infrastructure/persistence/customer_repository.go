package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/qhato/ecommerce/internal/customer/domain"
	"github.com/qhato/ecommerce/pkg/database"
	"github.com/qhato/ecommerce/pkg/errors"
)

// PostgresCustomerRepository implements the CustomerRepository interface using PostgreSQL
type PostgresCustomerRepository struct {
	db *database.DB
}

// NewPostgresCustomerRepository creates a new PostgresCustomerRepository
func NewPostgresCustomerRepository(db *database.DB) *PostgresCustomerRepository {
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

	err := r.db.QueryRow(ctx, query,
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
		return errors.InternalWrap(err, "failed to create customer")
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

	// Using Pool().Exec to get RowsAffected
	tag, err := r.db.Pool().Exec(ctx, query,
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
		return errors.InternalWrap(err, "failed to update customer")
	}

	if tag.RowsAffected() == 0 {
		return errors.NotFound(fmt.Sprintf("customer %d", customer.ID))
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

	err := r.db.QueryRow(ctx, query, id).Scan(
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

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find customer by ID")
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

	err := r.db.QueryRow(ctx, query, email).Scan(
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

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find customer by email")
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

	err := r.db.QueryRow(ctx, query, username).Scan(
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

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find customer by username")
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
		if filter.ActiveOnly {
			query += " AND deactivated = false"
		}
		if !filter.IncludeArchived {
			query += " AND archived = false"
		}
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM blc_customer WHERE 1=1"
	countArgs := make([]interface{}, 0)
	if filter != nil {
		if filter.ActiveOnly {
			countQuery += " AND deactivated = false"
		}
		if !filter.IncludeArchived {
			countQuery += " AND archived = false"
		}
	}

	var total int64
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to count customers")
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

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to find customers")
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
			return nil, 0, errors.InternalWrap(err, "failed to scan customer")
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
		return nil, 0, errors.InternalWrap(err, "failed to iterate customers")
	}

	return customers, total, nil
}

// ExistsByEmail checks if a customer exists by email
func (r *PostgresCustomerRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM blc_customer WHERE email_address = $1)"
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, errors.InternalWrap(err, "failed to check customer by email")
	}
	return exists, nil
}

// ExistsByUsername checks if a customer exists by username
func (r *PostgresCustomerRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM blc_customer WHERE user_name = $1)"
	var exists bool
	err := r.db.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, errors.InternalWrap(err, "failed to check customer by username")
	}
	return exists, nil
}

func (r *PostgresCustomerRepository) UpdatePassword(ctx context.Context, customerID int64, hashedPassword string) error {
	query := `UPDATE blc_customer SET password = $1 WHERE customer_id = $2`
	tag, err := r.db.Pool().Exec(ctx, query, hashedPassword, customerID)
	if err != nil {
		return errors.InternalWrap(err, "failed to update password")
	}
	if tag.RowsAffected() == 0 {
		return errors.NotFound(fmt.Sprintf("customer %d", customerID))
	}
	return nil
}

// Delete soft deletes a customer by setting the archived flag.
func (r *PostgresCustomerRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE blc_customer SET archived = true WHERE customer_id = $1`
	tag, err := r.db.Pool().Exec(ctx, query, id)
	if err != nil {
		return errors.InternalWrap(err, "failed to soft delete customer")
	}
	if tag.RowsAffected() == 0 {
		return errors.NotFound(fmt.Sprintf("customer %d", id))
	}
	return nil
}