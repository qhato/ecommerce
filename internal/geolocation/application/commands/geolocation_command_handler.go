package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/geolocation/domain"
)

type GeolocationCommandHandler struct {
	locationRepo domain.LocationRepository
	geoZoneRepo  domain.GeoZoneRepository
}

func NewGeolocationCommandHandler(
	locationRepo domain.LocationRepository,
	geoZoneRepo domain.GeoZoneRepository,
) *GeolocationCommandHandler {
	return &GeolocationCommandHandler{
		locationRepo: locationRepo,
		geoZoneRepo:  geoZoneRepo,
	}
}

func (h *GeolocationCommandHandler) HandleCreateLocation(ctx context.Context, cmd CreateLocationCommand) (*domain.Location, error) {
	location, err := domain.NewLocation(
		cmd.Name,
		cmd.Address,
		cmd.City,
		cmd.State,
		cmd.Country,
		cmd.Latitude,
		cmd.Longitude,
	)
	if err != nil {
		return nil, err
	}

	location.PostalCode = cmd.PostalCode
	location.Timezone = cmd.Timezone
	if cmd.Metadata != nil {
		location.Metadata = cmd.Metadata
	}

	if err := h.locationRepo.Create(ctx, location); err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return location, nil
}

func (h *GeolocationCommandHandler) HandleUpdateLocation(ctx context.Context, cmd UpdateLocationCommand) (*domain.Location, error) {
	location, err := h.locationRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find location: %w", err)
	}
	if location == nil {
		return nil, domain.ErrLocationNotFound
	}

	location.Name = cmd.Name
	location.Address = cmd.Address
	location.City = cmd.City
	location.State = cmd.State
	location.Country = cmd.Country
	location.PostalCode = cmd.PostalCode
	location.Latitude = cmd.Latitude
	location.Longitude = cmd.Longitude
	location.Timezone = cmd.Timezone
	if cmd.Metadata != nil {
		location.Metadata = cmd.Metadata
	}
	location.UpdatedAt = time.Now()

	if err := h.locationRepo.Update(ctx, location); err != nil {
		return nil, fmt.Errorf("failed to update location: %w", err)
	}

	return location, nil
}

func (h *GeolocationCommandHandler) HandleDeleteLocation(ctx context.Context, cmd DeleteLocationCommand) error {
	return h.locationRepo.Delete(ctx, cmd.ID)
}

func (h *GeolocationCommandHandler) HandleCreateGeoZone(ctx context.Context, cmd CreateGeoZoneCommand) (*domain.GeoZone, error) {
	zone, err := domain.NewGeoZone(cmd.Name, cmd.Description, domain.GeoZoneType(cmd.Type))
	if err != nil {
		return nil, err
	}

	zone.Countries = cmd.Countries
	zone.States = cmd.States
	zone.PostalCodes = cmd.PostalCodes
	zone.Radius = cmd.Radius
	zone.CenterLat = cmd.CenterLat
	zone.CenterLng = cmd.CenterLng

	if cmd.Polygon != nil {
		zone.Polygon = make([]domain.Coordinate, len(cmd.Polygon))
		for i, p := range cmd.Polygon {
			zone.Polygon[i] = domain.Coordinate{
				Latitude:  p.Latitude,
				Longitude: p.Longitude,
			}
		}
	}

	if err := h.geoZoneRepo.Create(ctx, zone); err != nil {
		return nil, fmt.Errorf("failed to create geo zone: %w", err)
	}

	return zone, nil
}

func (h *GeolocationCommandHandler) HandleUpdateGeoZone(ctx context.Context, cmd UpdateGeoZoneCommand) (*domain.GeoZone, error) {
	zone, err := h.geoZoneRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find geo zone: %w", err)
	}
	if zone == nil {
		return nil, domain.ErrGeoZoneNotFound
	}

	zone.Name = cmd.Name
	zone.Description = cmd.Description
	zone.IsActive = cmd.IsActive
	zone.Countries = cmd.Countries
	zone.States = cmd.States
	zone.PostalCodes = cmd.PostalCodes
	zone.Radius = cmd.Radius
	zone.CenterLat = cmd.CenterLat
	zone.CenterLng = cmd.CenterLng

	if cmd.Polygon != nil {
		zone.Polygon = make([]domain.Coordinate, len(cmd.Polygon))
		for i, p := range cmd.Polygon {
			zone.Polygon[i] = domain.Coordinate{
				Latitude:  p.Latitude,
				Longitude: p.Longitude,
			}
		}
	}

	zone.UpdatedAt = time.Now()

	if err := h.geoZoneRepo.Update(ctx, zone); err != nil {
		return nil, fmt.Errorf("failed to update geo zone: %w", err)
	}

	return zone, nil
}

func (h *GeolocationCommandHandler) HandleDeleteGeoZone(ctx context.Context, cmd DeleteGeoZoneCommand) error {
	return h.geoZoneRepo.Delete(ctx, cmd.ID)
}
