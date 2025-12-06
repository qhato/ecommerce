package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/store/domain"
)

type StoreCommandHandler struct {
	storeRepo     domain.StoreRepository
	inventoryRepo domain.StoreInventoryRepository
}

func NewStoreCommandHandler(
	storeRepo domain.StoreRepository,
	inventoryRepo domain.StoreInventoryRepository,
) *StoreCommandHandler {
	return &StoreCommandHandler{
		storeRepo:     storeRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (h *StoreCommandHandler) HandleCreateStore(ctx context.Context, cmd CreateStoreCommand) (*domain.Store, error) {
	// Check if code is already taken
	existing, _ := h.storeRepo.FindByCode(ctx, cmd.Code)
	if existing != nil {
		return nil, domain.ErrStoreCodeTaken
	}

	store, err := domain.NewStore(cmd.Code, cmd.Name, domain.StoreType(cmd.Type))
	if err != nil {
		return nil, err
	}

	store.Description = cmd.Description
	store.Email = cmd.Email
	store.Phone = cmd.Phone
	store.Website = cmd.Website
	store.Timezone = cmd.Timezone
	store.Currency = cmd.Currency
	store.Locale = cmd.Locale
	store.TaxID = cmd.TaxID
	store.ParentStoreID = cmd.ParentStoreID

	store.Address = domain.Address{
		Street1:    cmd.Address.Street1,
		Street2:    cmd.Address.Street2,
		City:       cmd.Address.City,
		State:      cmd.Address.State,
		Country:    cmd.Address.Country,
		PostalCode: cmd.Address.PostalCode,
		Latitude:   cmd.Address.Latitude,
		Longitude:  cmd.Address.Longitude,
	}

	businessHours := make([]domain.BusinessHour, len(cmd.Settings.BusinessHours))
	for i, bh := range cmd.Settings.BusinessHours {
		businessHours[i] = domain.BusinessHour{
			DayOfWeek: bh.DayOfWeek,
			OpenTime:  bh.OpenTime,
			CloseTime: bh.CloseTime,
			IsClosed:  bh.IsClosed,
		}
	}

	store.Settings = domain.StoreSettings{
		AllowPickup:           cmd.Settings.AllowPickup,
		AllowShipping:         cmd.Settings.AllowShipping,
		AllowBackorder:        cmd.Settings.AllowBackorder,
		InventoryTracking:     cmd.Settings.InventoryTracking,
		DefaultShippingCost:   cmd.Settings.DefaultShippingCost,
		FreeShippingThreshold: cmd.Settings.FreeShippingThreshold,
		MinOrderAmount:        cmd.Settings.MinOrderAmount,
		MaxOrderAmount:        cmd.Settings.MaxOrderAmount,
		BusinessHours:         businessHours,
	}

	if cmd.Metadata != nil {
		store.Metadata = cmd.Metadata
	}

	if err := h.storeRepo.Create(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return store, nil
}

func (h *StoreCommandHandler) HandleUpdateStore(ctx context.Context, cmd UpdateStoreCommand) (*domain.Store, error) {
	store, err := h.storeRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find store: %w", err)
	}
	if store == nil {
		return nil, domain.ErrStoreNotFound
	}

	store.Name = cmd.Name
	store.Description = cmd.Description
	store.Email = cmd.Email
	store.Phone = cmd.Phone
	store.Website = cmd.Website
	store.Timezone = cmd.Timezone
	store.Currency = cmd.Currency
	store.Locale = cmd.Locale
	store.TaxID = cmd.TaxID

	store.Address = domain.Address{
		Street1:    cmd.Address.Street1,
		Street2:    cmd.Address.Street2,
		City:       cmd.Address.City,
		State:      cmd.Address.State,
		Country:    cmd.Address.Country,
		PostalCode: cmd.Address.PostalCode,
		Latitude:   cmd.Address.Latitude,
		Longitude:  cmd.Address.Longitude,
	}

	businessHours := make([]domain.BusinessHour, len(cmd.Settings.BusinessHours))
	for i, bh := range cmd.Settings.BusinessHours {
		businessHours[i] = domain.BusinessHour{
			DayOfWeek: bh.DayOfWeek,
			OpenTime:  bh.OpenTime,
			CloseTime: bh.CloseTime,
			IsClosed:  bh.IsClosed,
		}
	}

	store.Settings = domain.StoreSettings{
		AllowPickup:           cmd.Settings.AllowPickup,
		AllowShipping:         cmd.Settings.AllowShipping,
		AllowBackorder:        cmd.Settings.AllowBackorder,
		InventoryTracking:     cmd.Settings.InventoryTracking,
		DefaultShippingCost:   cmd.Settings.DefaultShippingCost,
		FreeShippingThreshold: cmd.Settings.FreeShippingThreshold,
		MinOrderAmount:        cmd.Settings.MinOrderAmount,
		MaxOrderAmount:        cmd.Settings.MaxOrderAmount,
		BusinessHours:         businessHours,
	}

	if cmd.Metadata != nil {
		store.Metadata = cmd.Metadata
	}

	store.UpdatedAt = time.Now()

	if err := h.storeRepo.Update(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to update store: %w", err)
	}

	return store, nil
}

func (h *StoreCommandHandler) HandleActivateStore(ctx context.Context, cmd ActivateStoreCommand) (*domain.Store, error) {
	store, err := h.storeRepo.FindByID(ctx, cmd.ID)
	if err != nil || store == nil {
		return nil, domain.ErrStoreNotFound
	}

	store.Activate()

	if err := h.storeRepo.Update(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to activate store: %w", err)
	}

	return store, nil
}

func (h *StoreCommandHandler) HandleDeactivateStore(ctx context.Context, cmd DeactivateStoreCommand) (*domain.Store, error) {
	store, err := h.storeRepo.FindByID(ctx, cmd.ID)
	if err != nil || store == nil {
		return nil, domain.ErrStoreNotFound
	}

	store.Deactivate()

	if err := h.storeRepo.Update(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to deactivate store: %w", err)
	}

	return store, nil
}

func (h *StoreCommandHandler) HandleCloseStore(ctx context.Context, cmd CloseStoreCommand) (*domain.Store, error) {
	store, err := h.storeRepo.FindByID(ctx, cmd.ID)
	if err != nil || store == nil {
		return nil, domain.ErrStoreNotFound
	}

	store.Close()

	if err := h.storeRepo.Update(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to close store: %w", err)
	}

	return store, nil
}

func (h *StoreCommandHandler) HandleDeleteStore(ctx context.Context, cmd DeleteStoreCommand) error {
	return h.storeRepo.Delete(ctx, cmd.ID)
}

func (h *StoreCommandHandler) HandleUpdateInventory(ctx context.Context, cmd UpdateInventoryCommand) (*domain.StoreInventory, error) {
	inventory, err := h.inventoryRepo.FindByStoreAndProduct(ctx, cmd.StoreID, cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory: %w", err)
	}

	if inventory == nil {
		// Create new inventory record
		inventory = &domain.StoreInventory{
			StoreID:   cmd.StoreID,
			ProductID: cmd.ProductID,
			SKU:       cmd.SKU,
			UpdatedAt: time.Now(),
		}
		inventory.UpdateInventory(cmd.QuantityChange)

		if err := h.inventoryRepo.Create(ctx, inventory); err != nil {
			return nil, fmt.Errorf("failed to create inventory: %w", err)
		}
	} else {
		inventory.UpdateInventory(cmd.QuantityChange)

		if err := h.inventoryRepo.Update(ctx, inventory); err != nil {
			return nil, fmt.Errorf("failed to update inventory: %w", err)
		}
	}

	return inventory, nil
}

func (h *StoreCommandHandler) HandleReserveInventory(ctx context.Context, cmd ReserveInventoryCommand) (*domain.StoreInventory, error) {
	inventory, err := h.inventoryRepo.FindByStoreAndProduct(ctx, cmd.StoreID, cmd.ProductID)
	if err != nil || inventory == nil {
		return nil, domain.ErrInventoryNotFound
	}

	if err := inventory.Reserve(cmd.Quantity); err != nil {
		return nil, err
	}

	if err := h.inventoryRepo.Update(ctx, inventory); err != nil {
		return nil, fmt.Errorf("failed to reserve inventory: %w", err)
	}

	return inventory, nil
}

func (h *StoreCommandHandler) HandleReleaseInventory(ctx context.Context, cmd ReleaseInventoryCommand) (*domain.StoreInventory, error) {
	inventory, err := h.inventoryRepo.FindByStoreAndProduct(ctx, cmd.StoreID, cmd.ProductID)
	if err != nil || inventory == nil {
		return nil, domain.ErrInventoryNotFound
	}

	inventory.Release(cmd.Quantity)

	if err := h.inventoryRepo.Update(ctx, inventory); err != nil {
		return nil, fmt.Errorf("failed to release inventory: %w", err)
	}

	return inventory, nil
}
