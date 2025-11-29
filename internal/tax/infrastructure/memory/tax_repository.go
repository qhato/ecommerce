package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/qhato/ecommerce/internal/tax/domain"
)

// TaxRateRepository implements domain.TaxRateRepository for in-memory persistence.
type TaxRateRepository struct {
	mu      sync.RWMutex
	taxRates map[int64]*domain.TaxRate
	nextID  int64
}

// NewTaxRateRepository creates a new in-memory tax rate repository.
func NewTaxRateRepository() *TaxRateRepository {
	return &TaxRateRepository{
		taxRates: make(map[int64]*domain.TaxRate),
		nextID:   1,
	}
}

// Save stores a new tax rate or updates an existing one.
func (r *TaxRateRepository) Save(ctx context.Context, taxRate *domain.TaxRate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if taxRate.ID == 0 {
		taxRate.ID = r.nextID
		r.nextID++
	}
	r.taxRates[taxRate.ID] = taxRate
	return nil
}

// FindByID retrieves a tax rate by its unique identifier.
func (r *TaxRateRepository) FindByID(ctx context.Context, id int64) (*domain.TaxRate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	taxRate, ok := r.taxRates[id]
	if !ok {
		return nil, nil
	}
	return taxRate, nil
}

// FindApplicableTaxRates retrieves tax rates applicable to a given jurisdiction and category.
func (r *TaxRateRepository) FindApplicableTaxRates(ctx context.Context, jurisdiction, category string) ([]*domain.TaxRate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var applicableRates []*domain.TaxRate
	now := time.Now()

	for _, rate := range r.taxRates {
		if rate.IsApplicable() &&
			rate.Jurisdiction == jurisdiction &&
			rate.Category == category &&
			(rate.StartDate.Before(now) || rate.StartDate.Equal(now)) &&
			(rate.EndDate == nil || rate.EndDate.After(now) || rate.EndDate.Equal(now)) {
			applicableRates = append(applicableRates, rate)
		}
	}
	return applicableRates, nil
}

// Delete removes a tax rate by its unique identifier.
func (r *TaxRateRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.taxRates[id]; !ok {
		return fmt.Errorf("tax rate with ID %d not found", id)
	}
	delete(r.taxRates, id)
	return nil
}
