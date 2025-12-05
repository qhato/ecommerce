package domain

import "time"

// SearchSynonym represents a search synonym mapping
// Business Logic: Permite que búsquedas de términos similares retornen los mismos resultados
// Ejemplo: "laptop" = "notebook" = "computadora portátil"
type SearchSynonym struct {
	ID        int64
	Term      string   // Término original
	Synonyms  []string // Lista de sinónimos
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SearchSynonymRepository defines the interface for search synonym persistence
type SearchSynonymRepository interface {
	Create(synonym *SearchSynonym) error
	Update(synonym *SearchSynonym) error
	Delete(id int64) error
	FindByID(id int64) (*SearchSynonym, error)
	FindByTerm(term string) (*SearchSynonym, error)
	FindAll() ([]*SearchSynonym, error)
	FindActive() ([]*SearchSynonym, error)
}

// MatchesTerm checks if a term matches this synonym
func (s *SearchSynonym) MatchesTerm(term string) bool {
	if s.Term == term {
		return true
	}
	for _, syn := range s.Synonyms {
		if syn == term {
			return true
		}
	}
	return false
}

// GetExpandedTerms returns all terms including the original and synonyms
func (s *SearchSynonym) GetExpandedTerms() []string {
	terms := make([]string, 0, len(s.Synonyms)+1)
	terms = append(terms, s.Term)
	terms = append(terms, s.Synonyms...)
	return terms
}

// Activate activates the synonym
func (s *SearchSynonym) Activate() {
	s.IsActive = true
	s.UpdatedAt = time.Now()
}

// Deactivate deactivates the synonym
func (s *SearchSynonym) Deactivate() {
	s.IsActive = false
	s.UpdatedAt = time.Now()
}
