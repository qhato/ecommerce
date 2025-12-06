package domain

import "errors"

var (
	ErrLocationNotFound     = errors.New("location not found")
	ErrGeoZoneNotFound      = errors.New("geo zone not found")
	ErrInvalidCoordinates   = errors.New("invalid coordinates")
	ErrInvalidRadius        = errors.New("invalid radius")
	ErrInvalidPolygon       = errors.New("invalid polygon")
	ErrGeocodingFailed      = errors.New("geocoding failed")
	ErrReverseGeocodeFailed = errors.New("reverse geocoding failed")
)
