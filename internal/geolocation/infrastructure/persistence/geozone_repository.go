package persistence

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/qhato/ecommerce/internal/geolocation/domain"
)

type PostgresGeoZoneRepository struct {
	db *sql.DB
}

func NewPostgresGeoZoneRepository(db *sql.DB) *PostgresGeoZoneRepository {
	return &PostgresGeoZoneRepository{db: db}
}

func (r *PostgresGeoZoneRepository) Create(ctx context.Context, zone *domain.GeoZone) error {
	countriesJSON, _ := json.Marshal(zone.Countries)
	statesJSON, _ := json.Marshal(zone.States)
	postalCodesJSON, _ := json.Marshal(zone.PostalCodes)
	polygonJSON, _ := json.Marshal(zone.Polygon)

	query := `INSERT INTO blc_geo_zone (
		name, description, type, is_active, countries, states, postal_codes,
		radius, center_lat, center_lng, polygon, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		zone.Name, zone.Description, zone.Type, zone.IsActive, countriesJSON, statesJSON,
		postalCodesJSON, zone.Radius, zone.CenterLat, zone.CenterLng, polygonJSON,
		zone.CreatedAt, zone.UpdatedAt,
	).Scan(&zone.ID)
}

func (r *PostgresGeoZoneRepository) Update(ctx context.Context, zone *domain.GeoZone) error {
	countriesJSON, _ := json.Marshal(zone.Countries)
	statesJSON, _ := json.Marshal(zone.States)
	postalCodesJSON, _ := json.Marshal(zone.PostalCodes)
	polygonJSON, _ := json.Marshal(zone.Polygon)

	query := `UPDATE blc_geo_zone SET
		name = $1, description = $2, is_active = $3, countries = $4, states = $5,
		postal_codes = $6, radius = $7, center_lat = $8, center_lng = $9, polygon = $10,
		updated_at = $11
	WHERE id = $12`

	_, err := r.db.ExecContext(ctx, query,
		zone.Name, zone.Description, zone.IsActive, countriesJSON, statesJSON,
		postalCodesJSON, zone.Radius, zone.CenterLat, zone.CenterLng, polygonJSON,
		zone.UpdatedAt, zone.ID,
	)
	return err
}

func (r *PostgresGeoZoneRepository) FindByID(ctx context.Context, id int64) (*domain.GeoZone, error) {
	query := `SELECT id, name, description, type, is_active, countries, states, postal_codes,
		radius, center_lat, center_lng, polygon, created_at, updated_at
	FROM blc_geo_zone WHERE id = $1`

	return r.scanGeoZone(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresGeoZoneRepository) FindByType(ctx context.Context, zoneType domain.GeoZoneType) ([]*domain.GeoZone, error) {
	query := `SELECT id, name, description, type, is_active, countries, states, postal_codes,
		radius, center_lat, center_lng, polygon, created_at, updated_at
	FROM blc_geo_zone WHERE type = $1 ORDER BY name`

	return r.queryGeoZones(ctx, query, zoneType)
}

func (r *PostgresGeoZoneRepository) FindActive(ctx context.Context, limit int) ([]*domain.GeoZone, error) {
	query := `SELECT id, name, description, type, is_active, countries, states, postal_codes,
		radius, center_lat, center_lng, polygon, created_at, updated_at
	FROM blc_geo_zone WHERE is_active = true ORDER BY name LIMIT $1`

	return r.queryGeoZones(ctx, query, limit)
}

func (r *PostgresGeoZoneRepository) FindByLocation(ctx context.Context, location *domain.Location) ([]*domain.GeoZone, error) {
	query := `SELECT id, name, description, type, is_active, countries, states, postal_codes,
		radius, center_lat, center_lng, polygon, created_at, updated_at
	FROM blc_geo_zone WHERE is_active = true ORDER BY name`

	zones, err := r.queryGeoZones(ctx, query)
	if err != nil {
		return nil, err
	}

	// Filter zones that contain the location
	matchingZones := make([]*domain.GeoZone, 0)
	for _, zone := range zones {
		if zone.ContainsLocation(location) {
			matchingZones = append(matchingZones, zone)
		}
	}

	return matchingZones, nil
}

func (r *PostgresGeoZoneRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_geo_zone WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresGeoZoneRepository) scanGeoZone(row interface {
	Scan(dest ...interface{}) error
}) (*domain.GeoZone, error) {
	zone := &domain.GeoZone{}
	var countriesJSON, statesJSON, postalCodesJSON, polygonJSON []byte

	err := row.Scan(
		&zone.ID, &zone.Name, &zone.Description, &zone.Type, &zone.IsActive,
		&countriesJSON, &statesJSON, &postalCodesJSON, &zone.Radius,
		&zone.CenterLat, &zone.CenterLng, &polygonJSON, &zone.CreatedAt, &zone.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(countriesJSON, &zone.Countries); err != nil {
		zone.Countries = make([]string, 0)
	}
	if err := json.Unmarshal(statesJSON, &zone.States); err != nil {
		zone.States = make([]string, 0)
	}
	if err := json.Unmarshal(postalCodesJSON, &zone.PostalCodes); err != nil {
		zone.PostalCodes = make([]string, 0)
	}
	if err := json.Unmarshal(polygonJSON, &zone.Polygon); err != nil {
		zone.Polygon = make([]domain.Coordinate, 0)
	}

	return zone, nil
}

func (r *PostgresGeoZoneRepository) queryGeoZones(ctx context.Context, query string, args ...interface{}) ([]*domain.GeoZone, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	zones := make([]*domain.GeoZone, 0)
	for rows.Next() {
		zone := &domain.GeoZone{}
		var countriesJSON, statesJSON, postalCodesJSON, polygonJSON []byte

		if err := rows.Scan(
			&zone.ID, &zone.Name, &zone.Description, &zone.Type, &zone.IsActive,
			&countriesJSON, &statesJSON, &postalCodesJSON, &zone.Radius,
			&zone.CenterLat, &zone.CenterLng, &polygonJSON, &zone.CreatedAt, &zone.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(countriesJSON, &zone.Countries); err != nil {
			zone.Countries = make([]string, 0)
		}
		if err := json.Unmarshal(statesJSON, &zone.States); err != nil {
			zone.States = make([]string, 0)
		}
		if err := json.Unmarshal(postalCodesJSON, &zone.PostalCodes); err != nil {
			zone.PostalCodes = make([]string, 0)
		}
		if err := json.Unmarshal(polygonJSON, &zone.Polygon); err != nil {
			zone.Polygon = make([]domain.Coordinate, 0)
		}

		zones = append(zones, zone)
	}

	return zones, nil
}
