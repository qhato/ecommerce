package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/inventory/domain"
)

type InventoryCommandHandler struct {
	inventoryRepo   domain.InventoryRepository
	reservationRepo domain.InventoryReservationRepository
}

func NewInventoryCommandHandler(
	inventoryRepo domain.InventoryRepository,
	reservationRepo domain.InventoryReservationRepository,
) *InventoryCommandHandler {
	return &InventoryCommandHandler{
		inventoryRepo:   inventoryRepo,
		reservationRepo: reservationRepo,
	}
}

// Inventory Level Commands

func (h *InventoryCommandHandler) HandleCreateInventoryLevel(ctx context.Context, cmd CreateInventoryLevelCommand) (*domain.InventoryLevel, error) {
	level, err := domain.NewInventoryLevel(cmd.SKUID, cmd.QuantityOnHand)
	if err != nil {
		return nil, err
	}

	level.WarehouseID = &cmd.WarehouseID
	level.LocationID = &cmd.LocationID
	level.ReorderPoint = cmd.ReorderPoint
	level.ReorderQuantity = cmd.ReorderQty
	level.SafetyStock = cmd.SafetyStock
	level.AllowBackorder = cmd.AllowBackorder
	level.AllowPreorder = cmd.AllowPreorder

	if err := h.inventoryRepo.Save(ctx, level); err != nil {
		return nil, fmt.Errorf("failed to create inventory level: %w", err)
	}

	return level, nil
}

func (h *InventoryCommandHandler) HandleUpdateInventoryLevel(ctx context.Context, cmd UpdateInventoryLevelCommand) (*domain.InventoryLevel, error) {
	level, err := h.inventoryRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if level == nil {
		return nil, fmt.Errorf("inventory level not found")
	}

	level.RecordCount(cmd.QuantityOnHand)
	level.ReorderPoint = cmd.ReorderPoint
	level.ReorderQuantity = cmd.ReorderQty
	level.SafetyStock = cmd.SafetyStock
	level.AllowBackorder = cmd.AllowBackorder
	level.AllowPreorder = cmd.AllowPreorder

	if err := h.inventoryRepo.Save(ctx, level); err != nil {
		return nil, fmt.Errorf("failed to update inventory level: %w", err)
	}

	return level, nil
}

func (h *InventoryCommandHandler) HandleAdjustInventory(ctx context.Context, cmd AdjustInventoryCommand) (*domain.InventoryLevel, error) {
	level, err := h.inventoryRepo.FindBySKUID(ctx, cmd.SKUID)
	if err != nil {
		return nil, err
	}
	if level == nil {
		return nil, fmt.Errorf("inventory level not found for SKU: %s", cmd.SKUID)
	}

	if cmd.Adjustment > 0 {
		if err := level.Increment(cmd.Adjustment); err != nil {
			return nil, err
		}
	} else if cmd.Adjustment < 0 {
		// For negative adjustments, allocate and decrement
		absAdj := -cmd.Adjustment
		if err := level.Allocate(absAdj); err == nil {
			level.Decrement(absAdj)
		} else {
			// If can't allocate (not enough reserved), just update on hand directly
			level.QuantityOnHand -= absAdj
			if level.QuantityOnHand < 0 {
				level.QuantityOnHand = 0
			}
		}
	}

	if err := h.inventoryRepo.Save(ctx, level); err != nil {
		return nil, fmt.Errorf("failed to adjust inventory: %w", err)
	}

	return level, nil
}

func (h *InventoryCommandHandler) HandleSetInventory(ctx context.Context, cmd SetInventoryCommand) (*domain.InventoryLevel, error) {
	level, err := h.inventoryRepo.FindBySKUID(ctx, cmd.SKUID)
	if err != nil {
		return nil, err
	}
	if level == nil {
		return nil, fmt.Errorf("inventory level not found for SKU: %s", cmd.SKUID)
	}

	level.RecordCount(cmd.NewQuantity)

	if err := h.inventoryRepo.Save(ctx, level); err != nil {
		return nil, fmt.Errorf("failed to set inventory: %w", err)
	}

	return level, nil
}

func (h *InventoryCommandHandler) HandleDeleteInventoryLevel(ctx context.Context, cmd DeleteInventoryLevelCommand) error {
	if err := h.inventoryRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete inventory level: %w", err)
	}
	return nil
}

// Reservation Commands

func (h *InventoryCommandHandler) HandleReserveInventory(ctx context.Context, cmd ReserveInventoryCommand) (*domain.InventoryReservation, error) {
	// Create reservation
	reservation, err := domain.NewInventoryReservation(cmd.SKUID, cmd.OrderID, cmd.OrderItemID, cmd.Quantity, cmd.TTL)
	if err != nil {
		return nil, err
	}

	// Find inventory level
	level, err := h.inventoryRepo.FindBySKUID(ctx, cmd.SKUID)
	if err != nil {
		return nil, err
	}
	if level == nil {
		return nil, fmt.Errorf("inventory not found for SKU: %s", cmd.SKUID)
	}

	// Reserve inventory
	if err := level.Reserve(cmd.Quantity); err != nil {
		return nil, err
	}

	// Save both
	if err := h.inventoryRepo.Save(ctx, level); err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	if err := h.reservationRepo.Save(ctx, reservation); err != nil {
		// Rollback inventory reservation
		level.Release(cmd.Quantity)
		h.inventoryRepo.Save(ctx, level)
		return nil, fmt.Errorf("failed to save reservation: %w", err)
	}

	return reservation, nil
}

func (h *InventoryCommandHandler) HandleConfirmReservation(ctx context.Context, cmd ConfirmReservationCommand) (*domain.InventoryReservation, error) {
	reservation, err := h.reservationRepo.FindByID(ctx, cmd.ReservationID)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, fmt.Errorf("reservation not found")
	}

	if err := reservation.Confirm(); err != nil {
		return nil, err
	}

	if err := h.reservationRepo.Save(ctx, reservation); err != nil {
		return nil, fmt.Errorf("failed to confirm reservation: %w", err)
	}

	return reservation, nil
}

func (h *InventoryCommandHandler) HandleReleaseReservation(ctx context.Context, cmd ReleaseReservationCommand) error {
	reservation, err := h.reservationRepo.FindByID(ctx, cmd.ReservationID)
	if err != nil {
		return err
	}
	if reservation == nil {
		return fmt.Errorf("reservation not found")
	}

	if err := reservation.Release(); err != nil {
		return err
	}

	// Release inventory
	level, err := h.inventoryRepo.FindBySKUID(ctx, reservation.SKUID)
	if err != nil {
		return err
	}
	if level != nil {
		level.Release(reservation.Quantity)
		if err := h.inventoryRepo.Save(ctx, level); err != nil {
			return fmt.Errorf("failed to release inventory: %w", err)
		}
	}

	if err := h.reservationRepo.Save(ctx, reservation); err != nil {
		return fmt.Errorf("failed to save reservation: %w", err)
	}

	return nil
}

func (h *InventoryCommandHandler) HandleFulfillReservation(ctx context.Context, cmd FulfillReservationCommand) error {
	reservation, err := h.reservationRepo.FindByID(ctx, cmd.ReservationID)
	if err != nil {
		return err
	}
	if reservation == nil {
		return fmt.Errorf("reservation not found")
	}

	if err := reservation.Fulfill(); err != nil {
		return err
	}

	// Allocate inventory (reduce both on-hand and reserved)
	level, err := h.inventoryRepo.FindBySKUID(ctx, reservation.SKUID)
	if err != nil {
		return err
	}
	if level != nil {
		if err := level.Allocate(reservation.Quantity); err != nil {
			return err
		}
		if err := h.inventoryRepo.Save(ctx, level); err != nil {
			return fmt.Errorf("failed to allocate inventory: %w", err)
		}
	}

	if err := h.reservationRepo.Save(ctx, reservation); err != nil {
		return fmt.Errorf("failed to save reservation: %w", err)
	}

	return nil
}

func (h *InventoryCommandHandler) HandleExtendReservation(ctx context.Context, cmd ExtendReservationCommand) (*domain.InventoryReservation, error) {
	reservation, err := h.reservationRepo.FindByID(ctx, cmd.ReservationID)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, fmt.Errorf("reservation not found")
	}

	if err := reservation.ExtendExpiration(cmd.AdditionalTime); err != nil {
		return nil, err
	}

	if err := h.reservationRepo.Save(ctx, reservation); err != nil {
		return nil, fmt.Errorf("failed to extend reservation: %w", err)
	}

	return reservation, nil
}

func (h *InventoryCommandHandler) HandleExpireReservations(ctx context.Context, cmd ExpireReservationsCommand) (int, error) {
	expired, err := h.reservationRepo.FindExpired(ctx)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, reservation := range expired {
		if err := reservation.Expire(); err != nil {
			continue
		}

		// Release inventory
		level, err := h.inventoryRepo.FindBySKUID(ctx, reservation.SKUID)
		if err == nil && level != nil {
			level.Release(reservation.Quantity)
			h.inventoryRepo.Save(ctx, level)
		}

		if err := h.reservationRepo.Save(ctx, reservation); err != nil {
			continue
		}

		count++
	}

	return count, nil
}

func (h *InventoryCommandHandler) HandleReleaseOrderReservations(ctx context.Context, cmd ReleaseOrderReservationsCommand) error {
	reservations, err := h.reservationRepo.FindByOrderID(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	for _, reservation := range reservations {
		if reservation.Status == domain.ReservationStatusPending || reservation.Status == domain.ReservationStatusConfirmed {
			h.HandleReleaseReservation(ctx, ReleaseReservationCommand{ReservationID: reservation.ID})
		}
	}

	return nil
}

// Bulk Operations

func (h *InventoryCommandHandler) HandleBulkAdjustInventory(ctx context.Context, cmd BulkAdjustInventoryCommand) error {
	for _, adj := range cmd.Adjustments {
		level, err := h.inventoryRepo.FindBySKUID(ctx, adj.SKUID)
		if err != nil || level == nil {
			continue
		}

		if adj.Adjustment > 0 {
			level.Increment(adj.Adjustment)
		} else if adj.Adjustment < 0 {
			absAdj := -adj.Adjustment
			if err := level.Allocate(absAdj); err == nil {
				level.Decrement(absAdj)
			}
		}
		h.inventoryRepo.Save(ctx, level)
	}

	return nil
}

func (h *InventoryCommandHandler) HandleTransferInventory(ctx context.Context, cmd TransferInventoryCommand) error {
	// Find source warehouse inventory
	fromLevels, err := h.inventoryRepo.FindByWarehouse(ctx, cmd.FromWarehouseID)
	if err != nil {
		return err
	}

	var fromLevel *domain.InventoryLevel
	for _, level := range fromLevels {
		if level.SKUID == cmd.SKUID {
			fromLevel = level
			break
		}
	}

	if fromLevel == nil {
		return fmt.Errorf("source inventory not found")
	}

	// Allocate and decrement from source
	if err := fromLevel.Allocate(cmd.Quantity); err != nil {
		return fmt.Errorf("cannot allocate from source: %w", err)
	}
	if err := fromLevel.Decrement(cmd.Quantity); err != nil {
		return fmt.Errorf("cannot decrement from source: %w", err)
	}

	// Find or create destination inventory
	toLevels, err := h.inventoryRepo.FindByWarehouse(ctx, cmd.ToWarehouseID)
	if err != nil {
		return err
	}

	var toLevel *domain.InventoryLevel
	for _, level := range toLevels {
		if level.SKUID == cmd.SKUID {
			toLevel = level
			break
		}
	}

	if toLevel == nil {
		// Create new inventory level at destination
		toLevel, err = domain.NewInventoryLevel(cmd.SKUID, cmd.Quantity)
		if err != nil {
			return err
		}
		toLevel.WarehouseID = &cmd.ToWarehouseID
	} else {
		toLevel.Increment(cmd.Quantity)
	}

	// Save both
	if err := h.inventoryRepo.Save(ctx, fromLevel); err != nil {
		return fmt.Errorf("failed to update source inventory: %w", err)
	}

	if err := h.inventoryRepo.Save(ctx, toLevel); err != nil {
		// Rollback source
		fromLevel.Increment(cmd.Quantity)
		h.inventoryRepo.Save(ctx, fromLevel)
		return fmt.Errorf("failed to update destination inventory: %w", err)
	}

	return nil
}
