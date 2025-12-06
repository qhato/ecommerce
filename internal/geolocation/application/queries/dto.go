package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/geolocation/domain"
)

type LocationDTO struct {
	ID         int64                  `json:"id"`
	Name       string                 `json:"name"`
	Address    string                 `json:"address"`
	City       string                 `json:"city"`
	State      string                 `json:"state"`
	Country    string                 `json:"country"`
	PostalCode string                 `json:"postal_code"`
	Latitude   float64                `json:"latitude"`
	Longitude  float64                `json:"longitude"`
	Timezone   string                 `json:"timezone,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type GeoZoneDTO struct {
	ID          int64                `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        string               `json:"type"`
	IsActive    bool                 `json:"is_active"`
	Countries   []string             `json:"countries,omitempty"`
	States      []string             `json:"states,omitempty"`
	PostalCodes []string             `json:"postal_codes,omitempty"`
	Radius      float64              `json:"radius,omitempty"`
	CenterLat   *float64             `json:"center_lat,omitempty"`
	CenterLng   *float64             `json:"center_lng,omitempty"`
	Polygon     []CoordinateDTO      `json:"polygon,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type CoordinateDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DistanceResultDTO struct {
	FromLocation LocationDTO `json:"from_location"`
	ToLocation   LocationDTO `json:"to_location"`
	Distance     float64     `json:"distance"`
	Unit         string      `json:"unit"`
}

func ToLocationDTO(l *domain.Location) *LocationDTO {
	return &LocationDTO{
		ID:         l.ID,
		Name:       l.Name,
		Address:    l.Address,
		City:       l.City,
		State:      l.State,
		Country:    l.Country,
		PostalCode: l.PostalCode,
		Latitude:   l.Latitude,
		Longitude:  l.Longitude,
		Timezone:   l.Timezone,
		Metadata:   l.Metadata,
		CreatedAt:  l.CreatedAt,
		UpdatedAt:  l.UpdatedAt,
	}
}

func ToGeoZoneDTO(z *domain.GeoZone) *GeoZoneDTO {
	polygon := make([]CoordinateDTO, len(z.Polygon))
	for i, p := range z.Polygon {
		polygon[i] = CoordinateDTO{
			Latitude:  p.Latitude,
			Longitude: p.Longitude,
		}
	}

	return &GeoZoneDTO{
		ID:          z.ID,
		Name:        z.Name,
		Description: z.Description,
		Type:        string(z.Type),
		IsActive:    z.IsActive,
		Countries:   z.Countries,
		States:      z.States,
		PostalCodes: z.PostalCodes,
		Radius:      z.Radius,
		CenterLat:   z.CenterLat,
		CenterLng:   z.CenterLng,
		Polygon:     polygon,
		CreatedAt:   z.CreatedAt,
		UpdatedAt:   z.UpdatedAt,
	}
}
