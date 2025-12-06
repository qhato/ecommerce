package commands

// CreateLocationCommand creates a new location
type CreateLocationCommand struct {
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
}

// UpdateLocationCommand updates a location
type UpdateLocationCommand struct {
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
}

// DeleteLocationCommand deletes a location
type DeleteLocationCommand struct {
	ID int64 `json:"id"`
}

// CreateGeoZoneCommand creates a new geo zone
type CreateGeoZoneCommand struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Countries   []string  `json:"countries,omitempty"`
	States      []string  `json:"states,omitempty"`
	PostalCodes []string  `json:"postal_codes,omitempty"`
	Radius      float64   `json:"radius,omitempty"`
	CenterLat   *float64  `json:"center_lat,omitempty"`
	CenterLng   *float64  `json:"center_lng,omitempty"`
	Polygon     []LatLng  `json:"polygon,omitempty"`
}

// UpdateGeoZoneCommand updates a geo zone
type UpdateGeoZoneCommand struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	Countries   []string  `json:"countries,omitempty"`
	States      []string  `json:"states,omitempty"`
	PostalCodes []string  `json:"postal_codes,omitempty"`
	Radius      float64   `json:"radius,omitempty"`
	CenterLat   *float64  `json:"center_lat,omitempty"`
	CenterLng   *float64  `json:"center_lng,omitempty"`
	Polygon     []LatLng  `json:"polygon,omitempty"`
}

// DeleteGeoZoneCommand deletes a geo zone
type DeleteGeoZoneCommand struct {
	ID int64 `json:"id"`
}

// LatLng represents a coordinate
type LatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
