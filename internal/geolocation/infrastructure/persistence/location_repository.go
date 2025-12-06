package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/qhato/ecommerce/internal/geolocation/domain"
)

type PostgresLocationRepository struct {
	db *sql.DB
}

func NewPostgresLocationRepository(db *sql.DB) *PostgresLocationRepository {
	return &PostgresLocationRepository{db: db}
}

func (r *PostgresLocationRepository) Create(ctx context.Context, location *domain.Location) error {
	metadataJSON, err := json.Marshal(location.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `INSERT INTO blc_location (
		name, address, city, state, country, postal_code, latitude, longitude,
		timezone, metadata, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		location.Name, location.Address, location.City, location.State, location.Country,
		location.PostalCode, location.Latitude, location.Longitude, location.Timezone,
		metadataJSON, location.CreatedAt, location.UpdatedAt,
	).Scan(&location.ID)
}

func (r *PostgresLocationRepository) Update(ctx context.Context, location *domain.Location) error {
	metadataJSON, err := json.Marshal(location.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `UPDATE blc_location SET
		name = $1, address = $2, city = $3, state = $4, country = $5, postal_code = $6,
		latitude = $7, longitude = $8, timezone = $9, metadata = $10, updated_at = $11
	WHERE id = $12`

	_, err = r.db.ExecContext(ctx, query,
		location.Name, location.Address, location.City, location.State, location.Country,
		location.PostalCode, location.Latitude, location.Longitude, location.Timezone,
		metadataJSON, location.UpdatedAt, location.ID,
	)
	return err
}

func (r *PostgresLocationRepository) FindByID(ctx context.Context, id int64) (*domain.Location, error) {
	query := `SELECT id, name, address, city, state, country, postal_code, latitude, longitude,
		timezone, metadata, created_at, updated_at
	FROM blc_location WHERE id = $1`

	return r.scanLocation(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresLocationRepository) FindByCoordinates(ctx context.Context, lat, lng float64, radiusKm float64) ([]*domain.Location, error) {
	// Using Haversine formula in PostgreSQL
	query := `SELECT id, name, address, city, state, country, postal_code, latitude, longitude,
		timezone, metadata, created_at, updated_at,
		( 6371 * acos( cos( radians($1) ) * cos( radians( latitude ) ) *
		  cos( radians( longitude ) - radians($2) ) + sin( radians($1) ) *
		  sin( radians( latitude ) ) ) ) AS distance
	FROM blc_location
	WHERE ( 6371 * acos( cos( radians($1) ) * cos( radians( latitude ) ) *
		  cos( radians( longitude ) - radians($2) ) + sin( radians($1) ) *
		  sin( radians( latitude ) ) ) ) <= $3
	ORDER BY distance`

	return r.queryLocationsWithDistance(ctx, query, lat, lng, radiusKm)
}

func (r *PostgresLocationRepository) FindByCity(ctx context.Context, city, state, country string) ([]*domain.Location, error) {
	query := `SELECT id, name, address, city, state, country, postal_code, latitude, longitude,
		timezone, metadata, created_at, updated_at
	FROM blc_location WHERE city = $1 AND state = $2 AND country = $3 ORDER BY name`

	return r.queryLocations(ctx, query, city, state, country)
}

func (r *PostgresLocationRepository) FindByPostalCode(ctx context.Context, postalCode, country string) ([]*domain.Location, error) {
	query := `SELECT id, name, address, city, state, country, postal_code, latitude, longitude,
		timezone, metadata, created_at, updated_at
	FROM blc_location WHERE postal_code = $1 AND country = $2 ORDER BY name`

	return r.queryLocations(ctx, query, postalCode, country)
}

func (r *PostgresLocationRepository) FindNearby(ctx context.Context, lat, lng float64, limit int) ([]*domain.Location, error) {
	query := `SELECT id, name, address, city, state, country, postal_code, latitude, longitude,
		timezone, metadata, created_at, updated_at,
		( 6371 * acos( cos( radians($1) ) * cos( radians( latitude ) ) *
		  cos( radians( longitude ) - radians($2) ) + sin( radians($1) ) *
		  sin( radians( latitude ) ) ) ) AS distance
	FROM blc_location
	ORDER BY distance
	LIMIT $3`

	return r.queryLocationsWithDistance(ctx, query, lat, lng, limit)
}

func (r *PostgresLocationRepository) Search(ctx context.Context, query string, limit int) ([]*domain.Location, error) {
	sqlQuery := `SELECT id, name, address, city, state, country, postal_code, latitude, longitude,
		timezone, metadata, created_at, updated_at
	FROM blc_location
	WHERE name ILIKE $1 OR address ILIKE $1 OR city ILIKE $1
	ORDER BY name
	LIMIT $2`

	searchTerm := "%" + query + "%"
	return r.queryLocations(ctx, sqlQuery, searchTerm, limit)
}

func (r *PostgresLocationRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_location WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresLocationRepository) scanLocation(row interface {
	Scan(dest ...interface{}) error
}) (*domain.Location, error) {
	location := &domain.Location{}
	var metadataJSON []byte

	err := row.Scan(
		&location.ID, &location.Name, &location.Address, &location.City, &location.State,
		&location.Country, &location.PostalCode, &location.Latitude, &location.Longitude,
		&location.Timezone, &metadataJSON, &location.CreatedAt, &location.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(metadataJSON, &location.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return location, nil
}

func (r *PostgresLocationRepository) queryLocations(ctx context.Context, query string, args ...interface{}) ([]*domain.Location, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := make([]*domain.Location, 0)
	for rows.Next() {
		location := &domain.Location{}
		var metadataJSON []byte

		if err := rows.Scan(
			&location.ID, &location.Name, &location.Address, &location.City, &location.State,
			&location.Country, &location.PostalCode, &location.Latitude, &location.Longitude,
			&location.Timezone, &metadataJSON, &location.CreatedAt, &location.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(metadataJSON, &location.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		locations = append(locations, location)
	}

	return locations, nil
}

func (r *PostgresLocationRepository) queryLocationsWithDistance(ctx context.Context, query string, args ...interface{}) ([]*domain.Location, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := make([]*domain.Location, 0)
	for rows.Next() {
		location := &domain.Location{}
		var metadataJSON []byte
		var distance float64

		if err := rows.Scan(
			&location.ID, &location.Name, &location.Address, &location.City, &location.State,
			&location.Country, &location.PostalCode, &location.Latitude, &location.Longitude,
			&location.Timezone, &metadataJSON, &location.CreatedAt, &location.UpdatedAt, &distance,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(metadataJSON, &location.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		locations = append(locations, location)
	}

	return locations, nil
}
