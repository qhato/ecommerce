package domain

import "time"

// FacetType represents the type of facet
type FacetType string

const (
	FacetTypeField FacetType = "FIELD" // Facet por campo (ej: brand, color)
	FacetTypeRange FacetType = "RANGE" // Facet por rango (ej: precio $0-50, $50-100)
)

// SearchFacetConfig represents the configuration for a search facet
// Business Logic: Configura qué facets (filtros) se muestran en los resultados de búsqueda
type SearchFacetConfig struct {
	ID                int64
	Name              string    // Nombre interno (ej: "brand", "price")
	Label             string    // Etiqueta visible (ej: "Marca", "Precio")
	FieldName         string    // Campo en el índice de búsqueda
	FacetType         FacetType // Tipo de facet
	IsActive          bool
	ShowInResults     bool  // Mostrar en resultados de búsqueda
	ShowInNavigation  bool  // Mostrar en navegación lateral
	Priority          int   // Orden de visualización
	MinDocCount       int   // Mínimo de documentos para mostrar un valor
	MaxValues         int   // Máximo de valores a mostrar
	Ranges            []FacetRange // Solo para FacetTypeRange
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// FacetRange represents a range for range facets
type FacetRange struct {
	Label string  // Etiqueta (ej: "$0 - $50")
	Min   *float64 // Valor mínimo (nil = sin límite inferior)
	Max   *float64 // Valor máximo (nil = sin límite superior)
}

// SearchFacetConfigRepository defines the interface for facet config persistence
type SearchFacetConfigRepository interface {
	Create(config *SearchFacetConfig) error
	Update(config *SearchFacetConfig) error
	Delete(id int64) error
	FindByID(id int64) (*SearchFacetConfig, error)
	FindByName(name string) (*SearchFacetConfig, error)
	FindActive() ([]*SearchFacetConfig, error)
	FindForResults() ([]*SearchFacetConfig, error)
	FindForNavigation() ([]*SearchFacetConfig, error)
}

// IsFieldFacet checks if this is a field facet
func (c *SearchFacetConfig) IsFieldFacet() bool {
	return c.FacetType == FacetTypeField
}

// IsRangeFacet checks if this is a range facet
func (c *SearchFacetConfig) IsRangeFacet() bool {
	return c.FacetType == FacetTypeRange
}

// Activate activates the facet config
func (c *SearchFacetConfig) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Now()
}

// Deactivate deactivates the facet config
func (c *SearchFacetConfig) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now()
}

// AddRange adds a range to the facet (only for range facets)
func (c *SearchFacetConfig) AddRange(label string, min, max *float64) {
	if c.FacetType != FacetTypeRange {
		return
	}
	c.Ranges = append(c.Ranges, FacetRange{
		Label: label,
		Min:   min,
		Max:   max,
	})
	c.UpdatedAt = time.Now()
}
