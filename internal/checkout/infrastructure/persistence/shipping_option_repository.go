package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/qhato/ecommerce/internal/checkout/domain"
)

// PostgresShippingOptionRepository implements domain.ShippingOptionRepository
type PostgresShippingOptionRepository struct {
	db *sql.DB
}

// NewPostgresShippingOptionRepository creates a new repository
func NewPostgresShippingOptionRepository(db *sql.DB) *PostgresShippingOptionRepository {
	return &PostgresShippingOptionRepository{db: db}
}

// FindByID finds a shipping option by ID
func (r *PostgresShippingOptionRepository) FindByID(ctx context.Context, id string) (*domain.ShippingOption, error) {
	query := `
		SELECT id, name, description, carrier, service_code, speed,
			   estimated_days_min, estimated_days_max, base_cost, cost_per_item,
			   cost_per_weight, free_shipping_threshold, is_active, is_international,
			   requires_signature, allowed_countries, excluded_countries,
			   allowed_states, excluded_states, tracking_supported,
			   insurance_included, priority, created_at, updated_at
		FROM blc_shipping_option WHERE id = $1`

	option := &domain.ShippingOption{}
	var allowedCountries, excludedCountries, allowedStates, excludedStates pq.StringArray

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&option.ID, &option.Name, &option.Description, &option.Carrier,
		&option.ServiceCode, &option.Speed, &option.EstimatedDaysMin,
		&option.EstimatedDaysMax, &option.BaseCost, &option.CostPerItem,
		&option.CostPerWeight, &option.FreeShippingThreshold, &option.IsActive,
		&option.IsInternational, &option.RequiresSignature,
		&allowedCountries, &excludedCountries, &allowedStates, &excludedStates,
		&option.TrackingSupported, &option.InsuranceIncluded, &option.Priority,
		&option.CreatedAt, &option.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find shipping option: %w", err)
	}

	option.AllowedCountries = allowedCountries
	option.ExcludedCountries = excludedCountries
	option.AllowedStates = allowedStates
	option.ExcludedStates = excludedStates

	return option, nil
}

// FindAll finds all shipping options
func (r *PostgresShippingOptionRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.ShippingOption, error) {
	query := `
		SELECT id, name, description, carrier, service_code, speed,
			   estimated_days_min, estimated_days_max, base_cost, cost_per_item,
			   cost_per_weight, free_shipping_threshold, is_active, is_international,
			   requires_signature, allowed_countries, excluded_countries,
			   allowed_states, excluded_states, tracking_supported,
			   insurance_included, priority, created_at, updated_at
		FROM blc_shipping_option`

	if activeOnly {
		query += " WHERE is_active = true"
	}

	query += " ORDER BY priority ASC, name ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query shipping options: %w", err)
	}
	defer rows.Close()

	options := make([]*domain.ShippingOption, 0)
	for rows.Next() {
		option := &domain.ShippingOption{}
		var allowedCountries, excludedCountries, allowedStates, excludedStates pq.StringArray

		err := rows.Scan(
			&option.ID, &option.Name, &option.Description, &option.Carrier,
			&option.ServiceCode, &option.Speed, &option.EstimatedDaysMin,
			&option.EstimatedDaysMax, &option.BaseCost, &option.CostPerItem,
			&option.CostPerWeight, &option.FreeShippingThreshold, &option.IsActive,
			&option.IsInternational, &option.RequiresSignature,
			&allowedCountries, &excludedCountries, &allowedStates, &excludedStates,
			&option.TrackingSupported, &option.InsuranceIncluded, &option.Priority,
			&option.CreatedAt, &option.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan shipping option: %w", err)
		}

		option.AllowedCountries = allowedCountries
		option.ExcludedCountries = excludedCountries
		option.AllowedStates = allowedStates
		option.ExcludedStates = excludedStates

		options = append(options, option)
	}

	return options, nil
}

// FindByCarrier finds shipping options by carrier
func (r *PostgresShippingOptionRepository) FindByCarrier(ctx context.Context, carrier string, activeOnly bool) ([]*domain.ShippingOption, error) {
	query := `
		SELECT id, name, description, carrier, service_code, speed,
			   estimated_days_min, estimated_days_max, base_cost, cost_per_item,
			   cost_per_weight, free_shipping_threshold, is_active, is_international,
			   requires_signature, allowed_countries, excluded_countries,
			   allowed_states, excluded_states, tracking_supported,
			   insurance_included, priority, created_at, updated_at
		FROM blc_shipping_option
		WHERE carrier = $1`

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY priority ASC, name ASC"

	rows, err := r.db.QueryContext(ctx, query, carrier)
	if err != nil {
		return nil, fmt.Errorf("failed to query shipping options: %w", err)
	}
	defer rows.Close()

	options := make([]*domain.ShippingOption, 0)
	for rows.Next() {
		option := &domain.ShippingOption{}
		var allowedCountries, excludedCountries, allowedStates, excludedStates pq.StringArray

		err := rows.Scan(
			&option.ID, &option.Name, &option.Description, &option.Carrier,
			&option.ServiceCode, &option.Speed, &option.EstimatedDaysMin,
			&option.EstimatedDaysMax, &option.BaseCost, &option.CostPerItem,
			&option.CostPerWeight, &option.FreeShippingThreshold, &option.IsActive,
			&option.IsInternational, &option.RequiresSignature,
			&allowedCountries, &excludedCountries, &allowedStates, &excludedStates,
			&option.TrackingSupported, &option.InsuranceIncluded, &option.Priority,
			&option.CreatedAt, &option.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan shipping option: %w", err)
		}

		option.AllowedCountries = allowedCountries
		option.ExcludedCountries = excludedCountries
		option.AllowedStates = allowedStates
		option.ExcludedStates = excludedStates

		options = append(options, option)
	}

	return options, nil
}

// FindAvailableForAddress finds available shipping options for an address
func (r *PostgresShippingOptionRepository) FindAvailableForAddress(ctx context.Context, country, stateProvince, postalCode string) ([]*domain.ShippingOption, error) {
	// This is a simplified version - in production, you'd implement more complex logic
	// to check allowed/excluded countries and states
	query := `
		SELECT id, name, description, carrier, service_code, speed,
			   estimated_days_min, estimated_days_max, base_cost, cost_per_item,
			   cost_per_weight, free_shipping_threshold, is_active, is_international,
			   requires_signature, allowed_countries, excluded_countries,
			   allowed_states, excluded_states, tracking_supported,
			   insurance_included, priority, created_at, updated_at
		FROM blc_shipping_option
		WHERE is_active = true
		  AND (cardinality(allowed_countries) = 0 OR $1 = ANY(allowed_countries))
		  AND NOT ($1 = ANY(excluded_countries))
		ORDER BY priority ASC, name ASC`

	rows, err := r.db.QueryContext(ctx, query, country)
	if err != nil {
		return nil, fmt.Errorf("failed to query shipping options: %w", err)
	}
	defer rows.Close()

	options := make([]*domain.ShippingOption, 0)
	for rows.Next() {
		option := &domain.ShippingOption{}
		var allowedCountries, excludedCountries, allowedStates, excludedStates pq.StringArray

		err := rows.Scan(
			&option.ID, &option.Name, &option.Description, &option.Carrier,
			&option.ServiceCode, &option.Speed, &option.EstimatedDaysMin,
			&option.EstimatedDaysMax, &option.BaseCost, &option.CostPerItem,
			&option.CostPerWeight, &option.FreeShippingThreshold, &option.IsActive,
			&option.IsInternational, &option.RequiresSignature,
			&allowedCountries, &excludedCountries, &allowedStates, &excludedStates,
			&option.TrackingSupported, &option.InsuranceIncluded, &option.Priority,
			&option.CreatedAt, &option.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan shipping option: %w", err)
		}

		option.AllowedCountries = allowedCountries
		option.ExcludedCountries = excludedCountries
		option.AllowedStates = allowedStates
		option.ExcludedStates = excludedStates

		// Additional filtering using domain logic
		if option.IsAvailableForLocation(country, stateProvince) {
			options = append(options, option)
		}
	}

	return options, nil
}
