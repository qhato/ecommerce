package domain

import "context"

// LocationRepository defines the interface for location persistence
type LocationRepository interface {
	Create(ctx context.Context, location *Location) error
	Update(ctx context.Context, location *Location) error
	FindByID(ctx context.Context, id int64) (*Location, error)
	FindByCoordinates(ctx context.Context, lat, lng float64, radiusKm float64) ([]*Location, error)
	FindByCity(ctx context.Context, city, state, country string) ([]*Location, error)
	FindByPostalCode(ctx context.Context, postalCode, country string) ([]*Location, error)
	FindNearby(ctx context.Context, lat, lng float64, limit int) ([]*Location, error)
	Search(ctx context.Context, query string, limit int) ([]*Location, error)
	Delete(ctx context.Context, id int64) error
}

// GeoZoneRepository defines the interface for geo zone persistence
type GeoZoneRepository interface {
	Create(ctx context.Context, zone *GeoZone) error
	Update(ctx context.Context, zone *GeoZone) error
	FindByID(ctx context.Context, id int64) (*GeoZone, error)
	FindByType(ctx context.Context, zoneType GeoZoneType) ([]*GeoZone, error)
	FindActive(ctx context.Context, limit int) ([]*GeoZone, error)
	FindByLocation(ctx context.Context, location *Location) ([]*GeoZone, error)
	Delete(ctx context.Context, id int64) error
}
