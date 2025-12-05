package domain

import "time"

// SearchRedirect represents a search redirect rule
// Business Logic: Permite redirigir búsquedas específicas a URLs predefinidas
// Ejemplo: Buscar "iphone" → redirect a /categories/smartphones/iphone
type SearchRedirect struct {
	ID             int64
	SearchTerm     string // Término de búsqueda que activa el redirect
	TargetURL      string // URL de destino
	Priority       int    // Prioridad (mayor = más prioritario)
	IsActive       bool
	ActivationDate *time.Time // Fecha de activación (para campañas)
	ExpirationDate *time.Time // Fecha de expiración
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// SearchRedirectRepository defines the interface for search redirect persistence
type SearchRedirectRepository interface {
	Create(redirect *SearchRedirect) error
	Update(redirect *SearchRedirect) error
	Delete(id int64) error
	FindByID(id int64) (*SearchRedirect, error)
	FindBySearchTerm(term string) (*SearchRedirect, error)
	FindAllActive() ([]*SearchRedirect, error)
}

// IsCurrentlyActive checks if the redirect is currently active
func (r *SearchRedirect) IsCurrentlyActive() bool {
	if !r.IsActive {
		return false
	}

	now := time.Now()

	// Check activation date
	if r.ActivationDate != nil && r.ActivationDate.After(now) {
		return false
	}

	// Check expiration date
	if r.ExpirationDate != nil && r.ExpirationDate.Before(now) {
		return false
	}

	return true
}

// Activate activates the redirect
func (r *SearchRedirect) Activate() {
	r.IsActive = true
	r.UpdatedAt = time.Now()
}

// Deactivate deactivates the redirect
func (r *SearchRedirect) Deactivate() {
	r.IsActive = false
	r.UpdatedAt = time.Now()
}

// SetSchedule sets activation and expiration dates
func (r *SearchRedirect) SetSchedule(activationDate, expirationDate *time.Time) {
	r.ActivationDate = activationDate
	r.ExpirationDate = expirationDate
	r.UpdatedAt = time.Now()
}
