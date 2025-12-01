package application

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/inventory/domain"
)

// InventoryService defines the application service for inventory-related operations.
type InventoryService interface {
	CreateInventoryLevel(ctx context.Context, cmd *CreateInventoryLevelCommand) (*InventoryLevelDTO, error)
	GetInventoryLevelByID(ctx context.Context, id string) (*InventoryLevelDTO, error)
	GetInventoryLevelBySKUID(ctx context.Context, skuID string) (*InventoryLevelDTO, error)
	IncrementInventory(ctx context.Context, id string, quantity int) (*InventoryLevelDTO, error)
	DecrementInventory(ctx context.Context, id string, quantity int) (*InventoryLevelDTO, error)
	ReserveInventory(ctx context.Context, id string, quantity int) (*InventoryLevelDTO, error)
	ReleaseInventory(ctx context.Context, id string, quantity int) (*InventoryLevelDTO, error)
	UpdateInventoryQuantities(ctx context.Context, id string, quantityOnHand, quantityReserved int) (*InventoryLevelDTO, error)
	DeleteInventoryLevel(ctx context.Context, id string) error
}

// InventoryLevelDTO represents a SKU availability data transfer object.
type InventoryLevelDTO struct {
	ID                  string
	SKUID               string
	WarehouseID         *string
	LocationID          *string
	QuantityAvailable   int
	QuantityReserved    int
	QuantityOnHand      int
	QuantityAllocated   int
	QuantityBackordered int
	QuantityInTransit   int
	QuantityDamaged     int
	ReorderPoint        int
	ReorderQuantity     int
	SafetyStock         int
	AllowBackorder      bool
	AllowPreorder       bool
	LastCountDate       *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// CreateInventoryLevelCommand is a command to create a new SKUAvailability.
type CreateInventoryLevelCommand struct {
	SKUID          string
	QuantityOnHand int
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

func (s *inventoryService) CreateInventoryLevel(ctx context.Context, cmd *CreateInventoryLevelCommand) (*InventoryLevelDTO, error) {
	level, err := domain.NewInventoryLevel(cmd.SKUID, cmd.QuantityOnHand)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory level domain entity: %w", err)
	}

	err = s.inventoryRepo.Save(ctx, level)
	if err != nil {
		return nil, fmt.Errorf("failed to save inventory level: %w", err)
	}

	return toInventoryLevelDTO(level), nil
}

func (s *inventoryService) GetInventoryLevelByID(ctx context.Context, id string) (*InventoryLevelDTO, error) {
	level, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory level by ID: %w", err)
	}
	if level == nil {
		return nil, fmt.Errorf("inventory level with ID %s not found", id)
	}
	return toInventoryLevelDTO(level), nil
}

func (s *inventoryService) GetInventoryLevelBySKUID(ctx context.Context, skuID string) (*InventoryLevelDTO, error) {
	level, err := s.inventoryRepo.FindBySKUID(ctx, skuID)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory level by SKU ID: %w", err)
	}
	if level == nil {
		return nil, fmt.Errorf("inventory level for SKU ID %s not found", skuID)
	}
	return toInventoryLevelDTO(level), nil
}

func (s *inventoryService) IncrementInventory(ctx context.Context, id string, quantity int) (*InventoryLevelDTO, error) {
	level, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory level by ID for update: %w", err)
	}
	if level == nil {
		return nil, fmt.Errorf("inventory level with ID %s not found for update", id)
	}

	err = level.Increment(quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to increment inventory: %w", err)
	}
	err = s.inventoryRepo.Save(ctx, level)
	if err != nil {
		return nil, fmt.Errorf("failed to save inventory level after increment: %w", err)
	}
	return toInventoryLevelDTO(level), nil
}

func (s *inventoryService) DecrementInventory(ctx context.Context, id string, quantity int) (*InventoryLevelDTO, error) {
	level, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory level by ID for update: %w", err)
	}
	if level == nil {
		return nil, fmt.Errorf("inventory level with ID %s not found for update", id)
	}

	err = level.Decrement(quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to decrement inventory: %w", err)
	}
	err = s.inventoryRepo.Save(ctx, level)
	if err != nil {
		return nil, fmt.Errorf("failed to save inventory level after decrement: %w", err)
	}
	return toInventoryLevelDTO(level), nil
}

func (s *inventoryService) ReserveInventory(ctx context.Context, id string, quantity int) (*InventoryLevelDTO, error) {
	level, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory level by ID for update: %w", err)
	}
	if level == nil {
		return nil, fmt.Errorf("inventory level with ID %s not found for update", id)
	}

	err = level.Reserve(quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to reserve inventory: %w", err)
	}
	err = s.inventoryRepo.Save(ctx, level)
	if err != nil {
		return nil, fmt.Errorf("failed to save inventory level after reservation: %w", err)
	}
	return toInventoryLevelDTO(level), nil
}

func (s *inventoryService) ReleaseInventory(ctx context.Context, id string, quantity int) (*InventoryLevelDTO, error) {
	level, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory level by ID for update: %w", err)
	}
	if level == nil {
		return nil, fmt.Errorf("inventory level with ID %s not found for update", id)
	}

	err = level.Release(quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to release inventory: %w", err)
	}
	err = s.inventoryRepo.Save(ctx, level)
	if err != nil {
		return nil, fmt.Errorf("failed to save inventory level after release: %w", err)
	}
	return toInventoryLevelDTO(level), nil
}

func (s *inventoryService) DeleteInventoryLevel(ctx context.Context, id string) error {
	err := s.inventoryRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete inventory level: %w", err)
	}
	return nil
}

func (s *inventoryService) UpdateInventoryQuantities(ctx context.Context, id string, quantityOnHand, quantityReserved int) (*InventoryLevelDTO, error) {
	level, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory level by ID for update: %w", err)
	}
	if level == nil {
		return nil, fmt.Errorf("inventory level with ID %s not found for update", id)
	}

	level.QuantityOnHand = quantityOnHand
	level.QuantityReserved = quantityReserved
	level.UpdatedAt = time.Now() // Update timestamp

	err = s.inventoryRepo.Save(ctx, level)
	if err != nil {
		return nil, fmt.Errorf("failed to save inventory level after quantity update: %w", err)
	}
	return toInventoryLevelDTO(level), nil
}

func toInventoryLevelDTO(level *domain.InventoryLevel) *InventoryLevelDTO {
	return &InventoryLevelDTO{
		ID:                  level.ID,
		SKUID:               level.SKUID,
		WarehouseID:         level.WarehouseID,
		LocationID:          level.LocationID,
		QuantityAvailable:   level.QuantityAvailable,
		QuantityReserved:    level.QuantityReserved,
		QuantityOnHand:      level.QuantityOnHand,
		QuantityAllocated:   level.QuantityAllocated,
		QuantityBackordered: level.QuantityBackordered,
		QuantityInTransit:   level.QuantityInTransit,
		QuantityDamaged:     level.QuantityDamaged,
		ReorderPoint:        level.ReorderPoint,
		ReorderQuantity:     level.ReorderQuantity,
		SafetyStock:         level.SafetyStock,
		AllowBackorder:      level.AllowBackorder,
		AllowPreorder:       level.AllowPreorder,
		LastCountDate:       level.LastCountDate,
		CreatedAt:           level.CreatedAt,
		UpdatedAt:           level.UpdatedAt,
	}
}
