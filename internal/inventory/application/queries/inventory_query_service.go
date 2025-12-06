package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/inventory/domain"
)

type InventoryAvailability struct {
	SKUID           string
	Available       int
	OnHand          int
	Reserved        int
	Allocated       int
	Backordered     int
	InTransit       int
	CanBackorder    bool
	CanPreorder     bool
	IsAvailable     bool
}

type LowStockItem struct {
	SKUID          string
	WarehouseID    string
	Available      int
	ReorderPoint   int
	ReorderQty     int
	NeedsReorder   bool
}

type InventoryQueryService struct {
	inventoryRepo   domain.InventoryRepository
	reservationRepo domain.InventoryReservationRepository
}

func NewInventoryQueryService(
	inventoryRepo domain.InventoryRepository,
	reservationRepo domain.InventoryReservationRepository,
) *InventoryQueryService {
	return &InventoryQueryService{
		inventoryRepo:   inventoryRepo,
		reservationRepo: reservationRepo,
	}
}

// Inventory Level Queries

func (s *InventoryQueryService) GetInventoryLevel(ctx context.Context, query GetInventoryLevelQuery) (*domain.InventoryLevel, error) {
	level, err := s.inventoryRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory level: %w", err)
	}
	return level, nil
}

func (s *InventoryQueryService) GetInventoryBySKU(ctx context.Context, query GetInventoryBySKUQuery) (*domain.InventoryLevel, error) {
	level, err := s.inventoryRepo.FindBySKUID(ctx, query.SKUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory by SKU: %w", err)
	}
	return level, nil
}

func (s *InventoryQueryService) GetInventoryByWarehouse(ctx context.Context, query GetInventoryByWarehouseQuery) ([]*domain.InventoryLevel, error) {
	levels, err := s.inventoryRepo.FindByWarehouse(ctx, query.WarehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory by warehouse: %w", err)
	}
	return levels, nil
}

func (s *InventoryQueryService) CheckInventoryAvailability(ctx context.Context, query CheckInventoryAvailabilityQuery) (*InventoryAvailability, error) {
	level, err := s.inventoryRepo.FindBySKUID(ctx, query.SKUID)
	if err != nil {
		return nil, fmt.Errorf("failed to check availability: %w", err)
	}

	if level == nil {
		return &InventoryAvailability{
			SKUID:       query.SKUID,
			Available:   0,
			IsAvailable: false,
		}, nil
	}

	return &InventoryAvailability{
		SKUID:        level.SKUID,
		Available:    level.QuantityAvailable,
		OnHand:       level.QuantityOnHand,
		Reserved:     level.QuantityReserved,
		Allocated:    level.QuantityAllocated,
		Backordered:  level.QuantityBackordered,
		InTransit:    level.QuantityInTransit,
		CanBackorder: level.AllowBackorder,
		CanPreorder:  level.AllowPreorder,
		IsAvailable:  level.CanReserve(query.Quantity),
	}, nil
}

func (s *InventoryQueryService) GetLowStockItems(ctx context.Context, query GetLowStockItemsQuery) ([]*LowStockItem, error) {
	var levels []*domain.InventoryLevel
	var err error

	if query.WarehouseID != nil {
		levels, err = s.inventoryRepo.FindByWarehouse(ctx, *query.WarehouseID)
	} else {
		// This would need a FindAll method on the repository
		return nil, fmt.Errorf("warehouse ID required for low stock query")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get inventory levels: %w", err)
	}

	lowStock := make([]*LowStockItem, 0)
	for _, level := range levels {
		if level.NeedsReorder() {
			warehouseID := ""
			if level.WarehouseID != nil {
				warehouseID = *level.WarehouseID
			}
			lowStock = append(lowStock, &LowStockItem{
				SKUID:        level.SKUID,
				WarehouseID:  warehouseID,
				Available:    level.QuantityAvailable,
				ReorderPoint: level.ReorderPoint,
				ReorderQty:   level.ReorderQuantity,
				NeedsReorder: true,
			})

			if query.Limit > 0 && len(lowStock) >= query.Limit {
				break
			}
		}
	}

	return lowStock, nil
}

func (s *InventoryQueryService) GetBackorderableItems(ctx context.Context, query GetBackorderableItemsQuery) ([]*domain.InventoryLevel, error) {
	var levels []*domain.InventoryLevel
	var err error

	if query.WarehouseID != nil {
		levels, err = s.inventoryRepo.FindByWarehouse(ctx, *query.WarehouseID)
	} else {
		return nil, fmt.Errorf("warehouse ID required for backorderable items query")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get inventory levels: %w", err)
	}

	backorderable := make([]*domain.InventoryLevel, 0)
	for _, level := range levels {
		if level.AllowBackorder {
			backorderable = append(backorderable, level)
		}
	}

	return backorderable, nil
}

// Reservation Queries

func (s *InventoryQueryService) GetReservation(ctx context.Context, query GetReservationQuery) (*domain.InventoryReservation, error) {
	reservation, err := s.reservationRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}
	return reservation, nil
}

func (s *InventoryQueryService) GetReservationsByOrder(ctx context.Context, query GetReservationsByOrderQuery) ([]*domain.InventoryReservation, error) {
	reservations, err := s.reservationRepo.FindByOrderID(ctx, query.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations by order: %w", err)
	}
	return reservations, nil
}

func (s *InventoryQueryService) GetExpiredReservations(ctx context.Context, query GetExpiredReservationsQuery) ([]*domain.InventoryReservation, error) {
	reservations, err := s.reservationRepo.FindExpired(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired reservations: %w", err)
	}
	return reservations, nil
}

func (s *InventoryQueryService) GetActiveReservations(ctx context.Context, query GetActiveReservationsQuery) ([]*domain.InventoryReservation, error) {
	// This would require a new repository method
	// For now, we can get all reservations for a specific SKU by order ID
	return nil, fmt.Errorf("not implemented - requires FindActiveBySKU repository method")
}
