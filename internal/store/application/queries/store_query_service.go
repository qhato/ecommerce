package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/store/domain"
)

type StoreQueryService struct {
	storeRepo     domain.StoreRepository
	inventoryRepo domain.StoreInventoryRepository
}

func NewStoreQueryService(
	storeRepo domain.StoreRepository,
	inventoryRepo domain.StoreInventoryRepository,
) *StoreQueryService {
	return &StoreQueryService{
		storeRepo:     storeRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *StoreQueryService) GetStore(ctx context.Context, id int64) (*StoreDTO, error) {
	store, err := s.storeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find store: %w", err)
	}
	if store == nil {
		return nil, domain.ErrStoreNotFound
	}

	return ToStoreDTO(store), nil
}

func (s *StoreQueryService) GetStoreByCode(ctx context.Context, code string) (*StoreDTO, error) {
	store, err := s.storeRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to find store: %w", err)
	}
	if store == nil {
		return nil, domain.ErrStoreNotFound
	}

	return ToStoreDTO(store), nil
}

func (s *StoreQueryService) GetStoresByStatus(ctx context.Context, status string) ([]*StoreDTO, error) {
	stores, err := s.storeRepo.FindByStatus(ctx, domain.StoreStatus(status))
	if err != nil {
		return nil, fmt.Errorf("failed to find stores: %w", err)
	}

	dtos := make([]*StoreDTO, len(stores))
	for i, store := range stores {
		dtos[i] = ToStoreDTO(store)
	}

	return dtos, nil
}

func (s *StoreQueryService) GetStoresByType(ctx context.Context, storeType string) ([]*StoreDTO, error) {
	stores, err := s.storeRepo.FindByType(ctx, domain.StoreType(storeType))
	if err != nil {
		return nil, fmt.Errorf("failed to find stores: %w", err)
	}

	dtos := make([]*StoreDTO, len(stores))
	for i, store := range stores {
		dtos[i] = ToStoreDTO(store)
	}

	return dtos, nil
}

func (s *StoreQueryService) GetAllStores(ctx context.Context, limit int) ([]*StoreDTO, error) {
	stores, err := s.storeRepo.FindAll(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find stores: %w", err)
	}

	dtos := make([]*StoreDTO, len(stores))
	for i, store := range stores {
		dtos[i] = ToStoreDTO(store)
	}

	return dtos, nil
}

func (s *StoreQueryService) GetNearbyStores(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*StoreDTO, error) {
	stores, err := s.storeRepo.FindNearby(ctx, lat, lng, radiusKm, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearby stores: %w", err)
	}

	dtos := make([]*StoreDTO, len(stores))
	for i, store := range stores {
		dtos[i] = ToStoreDTO(store)
	}

	return dtos, nil
}

func (s *StoreQueryService) GetStoreInventory(ctx context.Context, storeID, productID int64) (*StoreInventoryDTO, error) {
	inventory, err := s.inventoryRepo.FindByStoreAndProduct(ctx, storeID, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory: %w", err)
	}
	if inventory == nil {
		return nil, domain.ErrInventoryNotFound
	}

	return ToStoreInventoryDTO(inventory), nil
}

func (s *StoreQueryService) GetInventoryByStore(ctx context.Context, storeID int64) ([]*StoreInventoryDTO, error) {
	inventory, err := s.inventoryRepo.FindByStore(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory: %w", err)
	}

	dtos := make([]*StoreInventoryDTO, len(inventory))
	for i, inv := range inventory {
		dtos[i] = ToStoreInventoryDTO(inv)
	}

	return dtos, nil
}

func (s *StoreQueryService) GetInventoryByProduct(ctx context.Context, productID int64) ([]*StoreInventoryDTO, error) {
	inventory, err := s.inventoryRepo.FindByProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory: %w", err)
	}

	dtos := make([]*StoreInventoryDTO, len(inventory))
	for i, inv := range inventory {
		dtos[i] = ToStoreInventoryDTO(inv)
	}

	return dtos, nil
}

func (s *StoreQueryService) GetLowStockItems(ctx context.Context, storeID int64) ([]*StoreInventoryDTO, error) {
	inventory, err := s.inventoryRepo.FindLowStock(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to find low stock items: %w", err)
	}

	dtos := make([]*StoreInventoryDTO, len(inventory))
	for i, inv := range inventory {
		dtos[i] = ToStoreInventoryDTO(inv)
	}

	return dtos, nil
}

func (s *StoreQueryService) GetInventoryBySKU(ctx context.Context, sku string) ([]*StoreInventoryDTO, error) {
	inventory, err := s.inventoryRepo.FindBySKU(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory: %w", err)
	}

	dtos := make([]*StoreInventoryDTO, len(inventory))
	for i, inv := range inventory {
		dtos[i] = ToStoreInventoryDTO(inv)
	}

	return dtos, nil
}
