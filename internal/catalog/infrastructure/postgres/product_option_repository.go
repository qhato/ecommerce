package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// ProductOptionRepository implements domain.ProductOptionRepository for PostgreSQL persistence.
type ProductOptionRepository struct {
	db *sql.DB
}

// NewProductOptionRepository creates a new PostgreSQL product option repository.
func NewProductOptionRepository(db *sql.DB) *ProductOptionRepository {
	return &ProductOptionRepository{db: db}
}

// Save stores a new product option or updates an existing one.
func (r *ProductOptionRepository) Save(ctx context.Context, option *domain.ProductOption) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle nullable fields
	attributeName := sql.NullString{String: option.AttributeName, Valid: option.AttributeName != ""}
	displayOrder := sql.NullInt32{Int32: int32(option.DisplayOrder), Valid: option.DisplayOrder != 0}
	errorCode := sql.NullString{String: option.ErrorCode, Valid: option.ErrorCode != ""}
	errorMessage := sql.NullString{String: option.ErrorMessage, Valid: option.ErrorMessage != ""}
	label := sql.NullString{String: option.Label, Valid: option.Label != ""}
	longDescription := sql.NullString{String: option.LongDescription, Valid: option.LongDescription != ""}
	name := sql.NullString{String: option.Name, Valid: option.Name != ""}
	validationStrategyType := sql.NullString{String: option.ValidationStrategyType, Valid: option.ValidationStrategyType != ""}
	validationType := sql.NullString{String: option.ValidationType, Valid: option.ValidationType != ""}
	required := sql.NullBool{Bool: option.Required, Valid: true}
	optionType := sql.NullString{String: option.OptionType, Valid: option.OptionType != ""}
	useInSkuGeneration := sql.NullBool{Bool: option.UseInSKUGeneration, Valid: true}
	validationString := sql.NullString{String: option.ValidationString, Valid: option.ValidationString != ""}

	if option.ID == 0 {
		// Insert new product option
		query := `
			INSERT INTO blc_product_option (
				attribute_name, display_order, error_code, error_message, label, 
				long_description, name, validation_strategy_type, validation_type, 
				required, option_type, use_in_sku_generation, validation_string, 
				created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
			) RETURNING product_option_id`
		err = tx.QueryRowContext(ctx, query,
			attributeName, displayOrder, errorCode, errorMessage, label,
			longDescription, name, validationStrategyType, validationType,
			required, optionType, useInSkuGeneration, validationString,
			option.CreatedAt, option.UpdatedAt,
		).Scan(&option.ID)
		if err != nil {
			return fmt.Errorf("failed to insert product option: %w", err)
		}
	} else {
		// Update existing product option
		query := `
			UPDATE blc_product_option SET
				attribute_name = $1, display_order = $2, error_code = $3, error_message = $4, 
				label = $5, long_description = $6, name = $7, 
				validation_strategy_type = $8, validation_type = $9, required = $10, 
				option_type = $11, use_in_sku_generation = $12, validation_string = $13, 
				updated_at = $14
			WHERE product_option_id = $15`
		_, err = tx.ExecContext(ctx, query,
			attributeName, displayOrder, errorCode, errorMessage, label,
			longDescription, name, validationStrategyType, validationType,
			required, optionType, useInSkuGeneration, validationString,
			option.UpdatedAt, option.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update product option: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a product option by its unique identifier.
func (r *ProductOptionRepository) FindByID(ctx context.Context, id int64) (*domain.ProductOption, error) {
	query := `
		SELECT
			product_option_id, attribute_name, display_order, error_code, error_message, 
			label, long_description, name, validation_strategy_type, validation_type, 
			required, option_type, use_in_sku_generation, validation_string, 
			created_at, updated_at
		FROM blc_product_option WHERE product_option_id = $1`

	var option domain.ProductOption
	var attributeName sql.NullString
	var displayOrder sql.NullInt32
	var errorCode sql.NullString
	var errorMessage sql.NullString
	var label sql.NullString
	var longDescription sql.NullString
	var name sql.NullString
	var validationStrategyType sql.NullString
	var validationType sql.NullString
	var required sql.NullBool
	var optionType sql.NullString
	var useInSkuGeneration sql.NullBool
	var validationString sql.NullString

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&option.ID, &attributeName, &displayOrder, &errorCode, &errorMessage,
		&label, &longDescription, &name, &validationStrategyType, &validationType,
		&required, &optionType, &useInSkuGeneration, &validationString,
		&option.CreatedAt, &option.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query product option by ID: %w", err)
	}

	if attributeName.Valid {
		option.AttributeName = attributeName.String
	}
	if displayOrder.Valid {
		option.DisplayOrder = int(displayOrder.Int32)
	}
	if errorCode.Valid {
		option.ErrorCode = errorCode.String
	}
	if errorMessage.Valid {
		option.ErrorMessage = errorMessage.String
	}
	if label.Valid {
		option.Label = label.String
	}
	if longDescription.Valid {
		option.LongDescription = longDescription.String
	}
	if name.Valid {
		option.Name = name.String
	}
	if validationStrategyType.Valid {
		option.ValidationStrategyType = validationStrategyType.String
	}
	if validationType.Valid {
		option.ValidationType = validationType.String
	}
	if required.Valid {
		option.Required = required.Bool
	}
	if optionType.Valid {
		option.OptionType = optionType.String
	}
	if useInSkuGeneration.Valid {
		option.UseInSKUGeneration = useInSkuGeneration.Bool
	}
	if validationString.Valid {
		option.ValidationString = validationString.String
	}

	return &option, nil
}

// FindAll retrieves all product options with pagination.
func (r *ProductOptionRepository) FindAll(ctx context.Context, filter *domain.ProductOptionFilter) ([]*domain.ProductOption, int64, error) {
	// Base query
	countQuery := `SELECT COUNT(*) FROM blc_product_option`
	query := `SELECT product_option_id, attribute_name, display_order, error_code, error_message, 
			label, long_description, name, validation_strategy_type, validation_type, 
			required, option_type, use_in_sku_generation, validation_string, 
			created_at, updated_at
		FROM blc_product_option`

	var args []interface{}
	argIdx := 1

	// Build WHERE clause (currently no filters in ProductOptionFilter, but keep structure)
	whereClauses := []string{}

	if len(whereClauses) > 0 {
		countQuery += " WHERE " + strings.Join(whereClauses, " AND ")
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total results
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count product options: %w", err)
	}

	// Apply sorting
	if filter.SortBy != "" {
		orderBy := map[string]string{
			"name":          "name",
			"display_order": "display_order",
			"created_at":    "created_at",
			"updated_at":    "updated_at",
		}
		sortColumn, ok := orderBy[filter.SortBy]
		if !ok {
			sortColumn = "name"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, strings.ToUpper(filter.SortOrder))
	}

	// Apply pagination
	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", argIdx, argIdx+1)
	args = append(args, (filter.Page-1)*filter.PageSize, filter.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query all product options: %w", err)
	}
	defer rows.Close()

	var options []*domain.ProductOption
	for rows.Next() {
		var option domain.ProductOption
		var attributeName sql.NullString
		var displayOrder sql.NullInt32
		var errorCode sql.NullString
		var errorMessage sql.NullString
		var label sql.NullString
		var longDescription sql.NullString
		var name sql.NullString
		var validationStrategyType sql.NullString
		var validationType sql.NullString
		var required sql.NullBool
		var optionType sql.NullString
		var useInSkuGeneration sql.NullBool
		var validationString sql.NullString

		err := rows.Scan(
			&option.ID, &attributeName, &displayOrder, &errorCode, &errorMessage,
			&label, &longDescription, &name, &validationStrategyType, &validationType,
			&required, &optionType, &useInSkuGeneration, &validationString,
			&option.CreatedAt, &option.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product option row: %w", err)
		}

		if attributeName.Valid {
			option.AttributeName = attributeName.String
		}
		if displayOrder.Valid {
			option.DisplayOrder = int(displayOrder.Int32)
		}
		if errorCode.Valid {
			option.ErrorCode = errorCode.String
		}
		if errorMessage.Valid {
			option.ErrorMessage = errorMessage.String
		}
		if label.Valid {
			option.Label = label.String
		}
		if longDescription.Valid {
			option.LongDescription = longDescription.String
		}
		if name.Valid {
			option.Name = name.String
		}
		if validationStrategyType.Valid {
			option.ValidationStrategyType = validationStrategyType.String
		}
		if validationType.Valid {
			option.ValidationType = validationType.String
		}
		if required.Valid {
			option.Required = required.Bool
		}
		if optionType.Valid {
			option.OptionType = optionType.String
		}
		if useInSkuGeneration.Valid {
			option.UseInSKUGeneration = useInSkuGeneration.Bool
		}
		if validationString.Valid {
			option.ValidationString = validationString.String
		}
		options = append(options, &option)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during rows iteration: %w", err)
	}

	return options, totalCount, nil
}

// Delete removes a product option by its unique identifier.
func (r *ProductOptionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_product_option WHERE product_option_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product option: %w", err)
	}
	return nil
}
