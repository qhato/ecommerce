package domain

import (
	"errors"
	"math"
	"time"
)

// Location represents a geographic location
type Location struct {
	ID          int64
	Name        string
	Address     string
	City        string
	State       string
	Country     string
	PostalCode  string
	Latitude    float64
	Longitude   float64
	Timezone    string
	Metadata    map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GeoZone represents a geographic zone/region
type GeoZone struct {
	ID          int64
	Name        string
	Description string
	Type        GeoZoneType
	IsActive    bool
	Countries   []string
	States      []string
	PostalCodes []string
	Radius      float64 // For circular zones in kilometers
	CenterLat   *float64
	CenterLng   *float64
	Polygon     []Coordinate // For polygon zones
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GeoZoneType represents the type of geographic zone
type GeoZoneType string

const (
	GeoZoneTypeCountry    GeoZoneType = "COUNTRY"
	GeoZoneTypeState      GeoZoneType = "STATE"
	GeoZoneTypePostalCode GeoZoneType = "POSTAL_CODE"
	GeoZoneTypeRadius     GeoZoneType = "RADIUS"
	GeoZoneTypePolygon    GeoZoneType = "POLYGON"
)

// Coordinate represents a geographic coordinate
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// DistanceUnit represents the unit for distance calculations
type DistanceUnit string

const (
	DistanceUnitKilometers DistanceUnit = "KM"
	DistanceUnitMiles      DistanceUnit = "MI"
	DistanceUnitMeters     DistanceUnit = "M"
)

// NewLocation creates a new location
func NewLocation(name, address, city, state, country string, lat, lng float64) (*Location, error) {
	if lat < -90 || lat > 90 {
		return nil, errors.New("invalid latitude: must be between -90 and 90")
	}
	if lng < -180 || lng > 180 {
		return nil, errors.New("invalid longitude: must be between -180 and 180")
	}

	now := time.Now()
	return &Location{
		Name:      name,
		Address:   address,
		City:      city,
		State:     state,
		Country:   country,
		Latitude:  lat,
		Longitude: lng,
		Metadata:  make(map[string]interface{}),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// NewGeoZone creates a new geographic zone
func NewGeoZone(name, description string, zoneType GeoZoneType) (*GeoZone, error) {
	now := time.Now()
	return &GeoZone{
		Name:        name,
		Description: description,
		Type:        zoneType,
		IsActive:    true,
		Countries:   make([]string, 0),
		States:      make([]string, 0),
		PostalCodes: make([]string, 0),
		Polygon:     make([]Coordinate, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// DistanceTo calculates the distance to another location using Haversine formula
func (l *Location) DistanceTo(other *Location, unit DistanceUnit) float64 {
	return CalculateDistance(l.Latitude, l.Longitude, other.Latitude, other.Longitude, unit)
}

// IsWithinRadius checks if location is within a radius of another location
func (l *Location) IsWithinRadius(centerLat, centerLng, radiusKm float64) bool {
	distance := CalculateDistance(l.Latitude, l.Longitude, centerLat, centerLng, DistanceUnitKilometers)
	return distance <= radiusKm
}

// ContainsLocation checks if a zone contains a location
func (z *GeoZone) ContainsLocation(loc *Location) bool {
	switch z.Type {
	case GeoZoneTypeCountry:
		return contains(z.Countries, loc.Country)
	case GeoZoneTypeState:
		return contains(z.States, loc.State)
	case GeoZoneTypePostalCode:
		return contains(z.PostalCodes, loc.PostalCode)
	case GeoZoneTypeRadius:
		if z.CenterLat != nil && z.CenterLng != nil {
			return loc.IsWithinRadius(*z.CenterLat, *z.CenterLng, z.Radius)
		}
		return false
	case GeoZoneTypePolygon:
		return pointInPolygon(loc.Latitude, loc.Longitude, z.Polygon)
	default:
		return false
	}
}

// CalculateDistance calculates distance between two coordinates using Haversine formula
func CalculateDistance(lat1, lon1, lat2, lon2 float64, unit DistanceUnit) float64 {
	const earthRadiusKm = 6371.0

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distanceKm := earthRadiusKm * c

	switch unit {
	case DistanceUnitMiles:
		return distanceKm * 0.621371
	case DistanceUnitMeters:
		return distanceKm * 1000
	default:
		return distanceKm
	}
}

// pointInPolygon checks if a point is inside a polygon using ray casting algorithm
func pointInPolygon(lat, lng float64, polygon []Coordinate) bool {
	if len(polygon) < 3 {
		return false
	}

	inside := false
	j := len(polygon) - 1

	for i := 0; i < len(polygon); i++ {
		xi, yi := polygon[i].Longitude, polygon[i].Latitude
		xj, yj := polygon[j].Longitude, polygon[j].Latitude

		intersect := ((yi > lat) != (yj > lat)) &&
			(lng < (xj-xi)*(lat-yi)/(yj-yi)+xi)
		if intersect {
			inside = !inside
		}
		j = i
	}

	return inside
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
