package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/store/domain"
)

type StoreDTO struct {
	ID            int64                  `json:"id"`
	Code          string                 `json:"code"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Type          string                 `json:"type"`
	Status        string                 `json:"status"`
	Email         string                 `json:"email"`
	Phone         string                 `json:"phone"`
	Website       string                 `json:"website"`
	Address       AddressDTO             `json:"address"`
	Timezone      string                 `json:"timezone"`
	Currency      string                 `json:"currency"`
	Locale        string                 `json:"locale"`
	TaxID         string                 `json:"tax_id,omitempty"`
	Settings      StoreSettingsDTO       `json:"settings"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	ParentStoreID *int64                 `json:"parent_store_id,omitempty"`
	IsOpen        bool                   `json:"is_open"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type AddressDTO struct {
	Street1    string  `json:"street1"`
	Street2    string  `json:"street2,omitempty"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
	PostalCode string  `json:"postal_code"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
}

type StoreSettingsDTO struct {
	AllowPickup           bool              `json:"allow_pickup"`
	AllowShipping         bool              `json:"allow_shipping"`
	AllowBackorder        bool              `json:"allow_backorder"`
	InventoryTracking     bool              `json:"inventory_tracking"`
	DefaultShippingCost   float64           `json:"default_shipping_cost"`
	FreeShippingThreshold float64           `json:"free_shipping_threshold,omitempty"`
	MinOrderAmount        float64           `json:"min_order_amount,omitempty"`
	MaxOrderAmount        float64           `json:"max_order_amount,omitempty"`
	BusinessHours         []BusinessHourDTO `json:"business_hours,omitempty"`
}

type BusinessHourDTO struct {
	DayOfWeek int    `json:"day_of_week"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
	IsClosed  bool   `json:"is_closed"`
}

type StoreInventoryDTO struct {
	ID              int64      `json:"id"`
	StoreID         int64      `json:"store_id"`
	StoreName       string     `json:"store_name,omitempty"`
	ProductID       int64      `json:"product_id"`
	SKU             string     `json:"sku"`
	QuantityOnHand  int        `json:"quantity_on_hand"`
	Reserved        int        `json:"reserved"`
	Available       int        `json:"available"`
	ReorderPoint    int        `json:"reorder_point"`
	ReorderQuantity int        `json:"reorder_quantity"`
	NeedsReorder    bool       `json:"needs_reorder"`
	LastRestocked   *time.Time `json:"last_restocked,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func ToStoreDTO(s *domain.Store) *StoreDTO {
	businessHours := make([]BusinessHourDTO, len(s.Settings.BusinessHours))
	for i, bh := range s.Settings.BusinessHours {
		businessHours[i] = BusinessHourDTO{
			DayOfWeek: bh.DayOfWeek,
			OpenTime:  bh.OpenTime,
			CloseTime: bh.CloseTime,
			IsClosed:  bh.IsClosed,
		}
	}

	return &StoreDTO{
		ID:          s.ID,
		Code:        s.Code,
		Name:        s.Name,
		Description: s.Description,
		Type:        string(s.Type),
		Status:      string(s.Status),
		Email:       s.Email,
		Phone:       s.Phone,
		Website:     s.Website,
		Address: AddressDTO{
			Street1:    s.Address.Street1,
			Street2:    s.Address.Street2,
			City:       s.Address.City,
			State:      s.Address.State,
			Country:    s.Address.Country,
			PostalCode: s.Address.PostalCode,
			Latitude:   s.Address.Latitude,
			Longitude:  s.Address.Longitude,
		},
		Timezone:      s.Timezone,
		Currency:      s.Currency,
		Locale:        s.Locale,
		TaxID:         s.TaxID,
		Settings: StoreSettingsDTO{
			AllowPickup:           s.Settings.AllowPickup,
			AllowShipping:         s.Settings.AllowShipping,
			AllowBackorder:        s.Settings.AllowBackorder,
			InventoryTracking:     s.Settings.InventoryTracking,
			DefaultShippingCost:   s.Settings.DefaultShippingCost,
			FreeShippingThreshold: s.Settings.FreeShippingThreshold,
			MinOrderAmount:        s.Settings.MinOrderAmount,
			MaxOrderAmount:        s.Settings.MaxOrderAmount,
			BusinessHours:         businessHours,
		},
		Metadata:      s.Metadata,
		ParentStoreID: s.ParentStoreID,
		IsOpen:        s.IsOpen(),
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

func ToStoreInventoryDTO(si *domain.StoreInventory) *StoreInventoryDTO {
	return &StoreInventoryDTO{
		ID:              si.ID,
		StoreID:         si.StoreID,
		ProductID:       si.ProductID,
		SKU:             si.SKU,
		QuantityOnHand:  si.QuantityOnHand,
		Reserved:        si.Reserved,
		Available:       si.Available,
		ReorderPoint:    si.ReorderPoint,
		ReorderQuantity: si.ReorderQuantity,
		NeedsReorder:    si.NeedsReorder(),
		LastRestocked:   si.LastRestocked,
		UpdatedAt:       si.UpdatedAt,
	}
}
