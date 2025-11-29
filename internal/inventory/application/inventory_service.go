package application

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/inventory/domain"
)

// InventoryService defines the application service for inventory-related operations.
type InventoryService interface {
	// CreateSKUAvailability creates a new SKU availability record.
	CreateSKUAvailability(ctx context.Context, cmd *CreateSKUAvailabilityCommand) (*SKUAvailabilityDTO, error)

	// GetSKUAvailabilityByID retrieves a SKU availability record by its ID.
	GetSKUAvailabilityByID(ctx context.Context, id int64) (*SKUAvailabilityDTO, error)

	// GetSKUAvailabilityBySKUID retrieves a SKU availability record by SKU ID.
	GetSKUAvailabilityBySKUID(ctx context.Context, skuID int64) (*SKUAvailabilityDTO, error)

	// UpdateSKUAvailabilityQuantities updates the quantities for a SKU availability record.
	UpdateSKUAvailabilityQuantities(ctx context.Context, id int64, qtyOnHand, reserveQty int) (*SKUAvailabilityDTO, error)

	// UpdateSKUAvailabilityStatus updates the status for a SKU availability record.
	UpdateSKUAvailabilityStatus(ctx context.Context, id int64, newStatus string) (*SKUAvailabilityDTO, error)

	// SetSKUAvailabilityLocation sets the location for a SKU availability record.
	SetSKUAvailabilityLocation(ctx context.Context, id int64, locationID int64) (*SKUAvailabilityDTO, error)

	// SetSKUAvailabilityDate sets the availability date for a SKU availability record.
	SetSKUAvailabilityDate(ctx context.Context, id int64, date *time.Time) (*SKUAvailabilityDTO, error)

	// DeleteSKUAvailability deletes a SKU availability record by its ID.
	DeleteSKUAvailability(ctx context.Context, id int64) error
}

// SKUAvailabilityDTO represents a SKU availability data transfer object.
type SKUAvailabilityDTO struct {
	ID                 int64
	SkuID              int64
	AvailabilityDate   *time.Time
	AvailabilityStatus string
	LocationID         *int64
	QtyOnHand          int
	ReserveQty         int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// CreateSKUAvailabilityCommand is a command to create a new SKUAvailability.
type CreateSKUAvailabilityCommand struct {
	SkuID              int64
	QtyOnHand          int
	ReserveQty         int
	AvailabilityStatus string
	AvailabilityDate   *time.Time
	LocationID         *int64
}

// UpdateSKUAvailabilityCommand is a command to update an existing SKUAvailability.
type UpdateSKUAvailabilityCommand struct {
	ID                 int64
	QtyOnHand          *int
	ReserveQty         *int
	AvailabilityStatus *string
	AvailabilityDate   *time.Time
	LocationID         *int64
}

type inventoryService struct {
	inventoryRepo domain.InventoryRepository
}

// NewInventoryService creates a new instance of InventoryService.
func NewInventoryService(inventoryRepo domain.InventoryRepository) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
	}
}

func (s *inventoryService) CreateSKUAvailability(ctx context.Context, cmd *CreateSKUAvailabilityCommand) (*SKUAvailabilityDTO, error) {
	availability, err := domain.NewSKUAvailability(cmd.SkuID, cmd.QtyOnHand, cmd.ReserveQty, cmd.AvailabilityStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to create SKU availability domain entity: %w", err)
	}

	if cmd.AvailabilityDate != nil {
		availability.SetAvailabilityDate(cmd.AvailabilityDate)
	}
	if cmd.LocationID != nil {
		availability.SetLocation(*cmd.LocationID)
	}

	err = s.inventoryRepo.Save(ctx, availability)
	if err != nil {
		return nil, fmt.Errorf("failed to save SKU availability: %w", err)
	}

	return toSKUAvailabilityDTO(availability), nil
}

func (s *inventoryService) GetSKUAvailabilityByID(ctx context.Context, id int64) (*SKUAvailabilityDTO, error) {
	availability, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find SKU availability by ID: %w", err)
	}
	if availability == nil {
		return nil, fmt.Errorf("SKU availability with ID %d not found", id)
	}
	return toSKUAvailabilityDTO(availability), nil
}

func (s *inventoryService) GetSKUAvailabilityBySKUID(ctx context.Context, skuID int64) (*SKUAvailabilityDTO, error) {
	availability, err := s.inventoryRepo.FindBySKUID(ctx, skuID)
	if err != nil {
		return nil, fmt.Errorf("failed to find SKU availability by SKU ID: %w", err)
	}
	if availability == nil {
		return nil, fmt.Errorf("SKU availability for SKU ID %d not found", skuID)
	}
	return toSKUAvailabilityDTO(availability), nil
}

func (s *inventoryService) UpdateSKUAvailabilityQuantities(ctx context.Context, id int64, qtyOnHand, reserveQty int) (*SKUAvailabilityDTO, error) {
	availability, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find SKU availability by ID for update: %w", err)
	}
	if availability == nil {
		return nil, fmt.Errorf("SKU availability with ID %d not found for update", id)
	}

	err = availability.UpdateQuantities(qtyOnHand, reserveQty)
	if err != nil {
		return nil, fmt.Errorf("failed to update quantities for SKU availability: %w", err)
	}
	err = s.inventoryRepo.Save(ctx, availability)
	if err != nil {
		return nil, fmt.Errorf("failed to save SKU availability after quantity update: %w", err)
	}
	return toSKUAvailabilityDTO(availability), nil
}

func (s *inventoryService) UpdateSKUAvailabilityStatus(ctx context.Context, id int64, newStatus string) (*SKUAvailabilityDTO, error) {
	availability, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find SKU availability by ID for update: %w", err)
	}
	if availability == nil {
		return nil, fmt.Errorf("SKU availability with ID %d not found for update", id)
	}

	availability.UpdateStatus(newStatus)
	err = s.inventoryRepo.Save(ctx, availability)
	if err != nil {
		return nil, fmt.Errorf("failed to save SKU availability after status update: %w", err)
	}
	return toSKUAvailabilityDTO(availability), nil
}

func (s *inventoryService) SetSKUAvailabilityLocation(ctx context.Context, id int64, locationID int64) (*SKUAvailabilityDTO, error) {
	availability, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find SKU availability by ID for update: %w", err)
	}
	if availability == nil {
		return nil, fmt.Errorf("SKU availability with ID %d not found for update", id)
	}

	availability.SetLocation(locationID)
	err = s.inventoryRepo.Save(ctx, availability)
	if err != nil {
		return nil, fmt.Errorf("failed to save SKU availability after location update: %w", err)
	}
	return toSKUAvailabilityDTO(availability), nil
}

func (s *inventoryService) SetSKUAvailabilityDate(ctx context.Context, id int64, date *time.Time) (*SKUAvailabilityDTO, error) {
	availability, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find SKU availability by ID for update: %w", err)
	}
	if availability == nil {
		return nil, fmt.Errorf("SKU availability with ID %d not found for update", id)
	}

	availability.SetAvailabilityDate(date)
	err = s.inventoryRepo.Save(ctx, availability)
	if err != nil {
		return nil, fmt.Errorf("failed to save SKU availability after date update: %w", err)
	}
	return toSKUAvailabilityDTO(availability), nil
}

func (s *inventoryService) DeleteSKUAvailability(ctx context.Context, id int64) error {
	err := s.inventoryRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete SKU availability: %w", err)
	}
	return nil
}

func toSKUAvailabilityDTO(availability *domain.SKUAvailability) *SKUAvailabilityDTO {
	return &SKUAvailabilityDTO{
		ID:                 availability.ID,
		SkuID:              availability.SkuID,
		AvailabilityDate:   availability.AvailabilityDate,
		AvailabilityStatus: availability.AvailabilityStatus,
		LocationID:         availability.LocationID,
		QtyOnHand:          availability.QtyOnHand,
		ReserveQty:         availability.ReserveQty,
		CreatedAt:          availability.CreatedAt,
		UpdatedAt:          availability.UpdatedAt,
	}
}
