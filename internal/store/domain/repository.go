package domain

import "context"

// StoreRepository defines the interface for store persistence
type StoreRepository interface {
	Create(ctx context.Context, store *Store) error
	Update(ctx context.Context, store *Store) error
	FindByID(ctx context.Context, id int64) (*Store, error)
	FindByCode(ctx context.Context, code string) (*Store, error)
	FindByStatus(ctx context.Context, status StoreStatus) ([]*Store, error)
	FindByType(ctx context.Context, storeType StoreType) ([]*Store, error)
	FindAll(ctx context.Context, limit int) ([]*Store, error)
	FindNearby(ctx context.Context, lat, lng float64, radiusKm float64, limit int) ([]*Store, error)
	Delete(ctx context.Context, id int64) error
}

// StoreInventoryRepository defines the interface for store inventory persistence
type StoreInventoryRepository interface {
	Create(ctx context.Context, inventory *StoreInventory) error
	Update(ctx context.Context, inventory *StoreInventory) error
	FindByStoreAndProduct(ctx context.Context, storeID, productID int64) (*StoreInventory, error)
	FindByStore(ctx context.Context, storeID int64) ([]*StoreInventory, error)
	FindByProduct(ctx context.Context, productID int64) ([]*StoreInventory, error)
	FindLowStock(ctx context.Context, storeID int64) ([]*StoreInventory, error)
	FindBySKU(ctx context.Context, sku string) ([]*StoreInventory, error)
	Delete(ctx context.Context, id int64) error
}

// TenantRepository defines the interface for tenant persistence
type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	Update(ctx context.Context, tenant *Tenant) error
	FindByID(ctx context.Context, id int64) (*Tenant, error)
	FindByCode(ctx context.Context, code string) (*Tenant, error)
	FindByStatus(ctx context.Context, status TenantStatus) ([]*Tenant, error)
	FindAll(ctx context.Context, limit int) ([]*Tenant, error)
	Delete(ctx context.Context, id int64) error
}

// StoreProductRepository defines the interface for store-product association persistence
type StoreProductRepository interface {
	Create(ctx context.Context, storeProduct *StoreProduct) error
	Update(ctx context.Context, storeProduct *StoreProduct) error
	FindByID(ctx context.Context, id int64) (*StoreProduct, error)
	FindByStoreAndProduct(ctx context.Context, storeID, productID int64) (*StoreProduct, error)
	FindByStore(ctx context.Context, storeID int64) ([]*StoreProduct, error)
	FindByProduct(ctx context.Context, productID int64) ([]*StoreProduct, error)
	FindAvailableByStore(ctx context.Context, storeID int64) ([]*StoreProduct, error)
	Delete(ctx context.Context, id int64) error
}

// PriceListRepository defines the interface for price list persistence
type PriceListRepository interface {
	Create(ctx context.Context, priceList *PriceList) error
	Update(ctx context.Context, priceList *PriceList) error
	FindByID(ctx context.Context, id int64) (*PriceList, error)
	FindByCode(ctx context.Context, code string) (*PriceList, error)
	FindActive(ctx context.Context) ([]*PriceList, error)
	FindAll(ctx context.Context) ([]*PriceList, error)
	Delete(ctx context.Context, id int64) error
}

// PriceListEntryRepository defines the interface for price list entry persistence
type PriceListEntryRepository interface {
	Create(ctx context.Context, entry *PriceListEntry) error
	Update(ctx context.Context, entry *PriceListEntry) error
	FindByID(ctx context.Context, id int64) (*PriceListEntry, error)
	FindByPriceList(ctx context.Context, priceListID int64) ([]*PriceListEntry, error)
	FindByPriceListAndProduct(ctx context.Context, priceListID, productID int64) (*PriceListEntry, error)
	Delete(ctx context.Context, id int64) error
}

// StorePaymentMethodRepository defines the interface for store payment method persistence
type StorePaymentMethodRepository interface {
	Create(ctx context.Context, method *StorePaymentMethod) error
	Update(ctx context.Context, method *StorePaymentMethod) error
	FindByStore(ctx context.Context, storeID int64) ([]*StorePaymentMethod, error)
	FindEnabledByStore(ctx context.Context, storeID int64) ([]*StorePaymentMethod, error)
	Delete(ctx context.Context, id int64) error
}

// StoreShippingMethodRepository defines the interface for store shipping method persistence
type StoreShippingMethodRepository interface {
	Create(ctx context.Context, method *StoreShippingMethod) error
	Update(ctx context.Context, method *StoreShippingMethod) error
	FindByStore(ctx context.Context, storeID int64) ([]*StoreShippingMethod, error)
	FindEnabledByStore(ctx context.Context, storeID int64) ([]*StoreShippingMethod, error)
	Delete(ctx context.Context, id int64) error
}
