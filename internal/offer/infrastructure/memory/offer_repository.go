package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/qhato/ecommerce/internal/offer/domain"
)

// OfferRepository implements domain.OfferRepository for in-memory persistence.
type OfferRepository struct {
	mu     sync.RWMutex
	offers map[int64]*domain.Offer
	nextID int64
}

// NewOfferRepository creates a new in-memory offer repository.
func NewOfferRepository() *OfferRepository {
	return &OfferRepository{
		offers: make(map[int64]*domain.Offer),
		nextID: 1,
	}
}

// Save stores a new offer or updates an existing one.
func (r *OfferRepository) Save(ctx context.Context, offer *domain.Offer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if offer.ID == 0 {
		offer.ID = r.nextID
		r.nextID++
	}
	r.offers[offer.ID] = offer
	return nil
}

// FindByID retrieves an offer by its unique identifier.
func (r *OfferRepository) FindByID(ctx context.Context, id int64) (*domain.Offer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	offer, ok := r.offers[id]
	if !ok {
		return nil, nil
	}
	return offer, nil
}

// FindByCode retrieves an offer by its promotional code.
// Note: Codes are now managed separately in OfferCode entity
func (r *OfferRepository) FindByCode(ctx context.Context, code string) (*domain.Offer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// This method is deprecated - use OfferCodeRepository instead
	return nil, nil
}

// FindActiveOffers retrieves all currently active offers.
func (r *OfferRepository) FindActiveOffers(ctx context.Context) ([]*domain.Offer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var activeOffers []*domain.Offer
	for _, offer := range r.offers {
		if !offer.Archived {
			activeOffers = append(activeOffers, offer)
		}
	}
	return activeOffers, nil
}

// Delete removes an offer by its unique identifier.
func (r *OfferRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.offers[id]; !ok {
		return fmt.Errorf("offer with ID %d not found", id)
	}
	delete(r.offers, id)
	return nil
}
