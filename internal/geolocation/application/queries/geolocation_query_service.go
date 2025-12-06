package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/geolocation/domain"
)

type GeolocationQueryService struct {
	locationRepo domain.LocationRepository
	geoZoneRepo  domain.GeoZoneRepository
}

func NewGeolocationQueryService(
	locationRepo domain.LocationRepository,
	geoZoneRepo domain.GeoZoneRepository,
) *GeolocationQueryService {
	return &GeolocationQueryService{
		locationRepo: locationRepo,
		geoZoneRepo:  geoZoneRepo,
	}
}

func (s *GeolocationQueryService) GetLocation(ctx context.Context, id int64) (*LocationDTO, error) {
	location, err := s.locationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find location: %w", err)
	}
	if location == nil {
		return nil, domain.ErrLocationNotFound
	}

	return ToLocationDTO(location), nil
}

func (s *GeolocationQueryService) FindLocationsByCity(ctx context.Context, city, state, country string) ([]*LocationDTO, error) {
	locations, err := s.locationRepo.FindByCity(ctx, city, state, country)
	if err != nil {
		return nil, fmt.Errorf("failed to find locations: %w", err)
	}

	dtos := make([]*LocationDTO, len(locations))
	for i, loc := range locations {
		dtos[i] = ToLocationDTO(loc)
	}

	return dtos, nil
}

func (s *GeolocationQueryService) FindLocationsByPostalCode(ctx context.Context, postalCode, country string) ([]*LocationDTO, error) {
	locations, err := s.locationRepo.FindByPostalCode(ctx, postalCode, country)
	if err != nil {
		return nil, fmt.Errorf("failed to find locations: %w", err)
	}

	dtos := make([]*LocationDTO, len(locations))
	for i, loc := range locations {
		dtos[i] = ToLocationDTO(loc)
	}

	return dtos, nil
}

func (s *GeolocationQueryService) FindNearbyLocations(ctx context.Context, lat, lng float64, limit int) ([]*LocationDTO, error) {
	locations, err := s.locationRepo.FindNearby(ctx, lat, lng, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearby locations: %w", err)
	}

	dtos := make([]*LocationDTO, len(locations))
	for i, loc := range locations {
		dtos[i] = ToLocationDTO(loc)
	}

	return dtos, nil
}

func (s *GeolocationQueryService) FindLocationsInRadius(ctx context.Context, lat, lng, radiusKm float64) ([]*LocationDTO, error) {
	locations, err := s.locationRepo.FindByCoordinates(ctx, lat, lng, radiusKm)
	if err != nil {
		return nil, fmt.Errorf("failed to find locations in radius: %w", err)
	}

	dtos := make([]*LocationDTO, len(locations))
	for i, loc := range locations {
		dtos[i] = ToLocationDTO(loc)
	}

	return dtos, nil
}

func (s *GeolocationQueryService) SearchLocations(ctx context.Context, query string, limit int) ([]*LocationDTO, error) {
	locations, err := s.locationRepo.Search(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search locations: %w", err)
	}

	dtos := make([]*LocationDTO, len(locations))
	for i, loc := range locations {
		dtos[i] = ToLocationDTO(loc)
	}

	return dtos, nil
}

func (s *GeolocationQueryService) CalculateDistance(ctx context.Context, fromID, toID int64, unit string) (*DistanceResultDTO, error) {
	fromLoc, err := s.locationRepo.FindByID(ctx, fromID)
	if err != nil || fromLoc == nil {
		return nil, domain.ErrLocationNotFound
	}

	toLoc, err := s.locationRepo.FindByID(ctx, toID)
	if err != nil || toLoc == nil {
		return nil, domain.ErrLocationNotFound
	}

	distanceUnit := domain.DistanceUnit(unit)
	distance := fromLoc.DistanceTo(toLoc, distanceUnit)

	return &DistanceResultDTO{
		FromLocation: *ToLocationDTO(fromLoc),
		ToLocation:   *ToLocationDTO(toLoc),
		Distance:     distance,
		Unit:         unit,
	}, nil
}

func (s *GeolocationQueryService) GetGeoZone(ctx context.Context, id int64) (*GeoZoneDTO, error) {
	zone, err := s.geoZoneRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find geo zone: %w", err)
	}
	if zone == nil {
		return nil, domain.ErrGeoZoneNotFound
	}

	return ToGeoZoneDTO(zone), nil
}

func (s *GeolocationQueryService) GetGeoZonesByType(ctx context.Context, zoneType string) ([]*GeoZoneDTO, error) {
	zones, err := s.geoZoneRepo.FindByType(ctx, domain.GeoZoneType(zoneType))
	if err != nil {
		return nil, fmt.Errorf("failed to find geo zones: %w", err)
	}

	dtos := make([]*GeoZoneDTO, len(zones))
	for i, zone := range zones {
		dtos[i] = ToGeoZoneDTO(zone)
	}

	return dtos, nil
}

func (s *GeolocationQueryService) GetActiveGeoZones(ctx context.Context, limit int) ([]*GeoZoneDTO, error) {
	zones, err := s.geoZoneRepo.FindActive(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find active geo zones: %w", err)
	}

	dtos := make([]*GeoZoneDTO, len(zones))
	for i, zone := range zones {
		dtos[i] = ToGeoZoneDTO(zone)
	}

	return dtos, nil
}

func (s *GeolocationQueryService) GetGeoZonesForLocation(ctx context.Context, locationID int64) ([]*GeoZoneDTO, error) {
	location, err := s.locationRepo.FindByID(ctx, locationID)
	if err != nil || location == nil {
		return nil, domain.ErrLocationNotFound
	}

	zones, err := s.geoZoneRepo.FindByLocation(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("failed to find geo zones: %w", err)
	}

	dtos := make([]*GeoZoneDTO, len(zones))
	for i, zone := range zones {
		dtos[i] = ToGeoZoneDTO(zone)
	}

	return dtos, nil
}
