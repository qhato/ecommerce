package commands

// CreateStoreCommand creates a new store
type CreateStoreCommand struct {
	Code          string                 `json:"code"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Type          string                 `json:"type"`
	Email         string                 `json:"email"`
	Phone         string                 `json:"phone"`
	Website       string                 `json:"website"`
	Address       AddressCommand         `json:"address"`
	Timezone      string                 `json:"timezone"`
	Currency      string                 `json:"currency"`
	Locale        string                 `json:"locale"`
	TaxID         string                 `json:"tax_id,omitempty"`
	Settings      StoreSettingsCommand   `json:"settings"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	ParentStoreID *int64                 `json:"parent_store_id,omitempty"`
}

// UpdateStoreCommand updates a store
type UpdateStoreCommand struct {
	ID            int64                  `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Email         string                 `json:"email"`
	Phone         string                 `json:"phone"`
	Website       string                 `json:"website"`
	Address       AddressCommand         `json:"address"`
	Timezone      string                 `json:"timezone"`
	Currency      string                 `json:"currency"`
	Locale        string                 `json:"locale"`
	TaxID         string                 `json:"tax_id,omitempty"`
	Settings      StoreSettingsCommand   `json:"settings"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// ActivateStoreCommand activates a store
type ActivateStoreCommand struct {
	ID int64 `json:"id"`
}

// DeactivateStoreCommand deactivates a store
type DeactivateStoreCommand struct {
	ID int64 `json:"id"`
}

// CloseStoreCommand closes a store
type CloseStoreCommand struct {
	ID int64 `json:"id"`
}

// DeleteStoreCommand deletes a store
type DeleteStoreCommand struct {
	ID int64 `json:"id"`
}

// UpdateInventoryCommand updates store inventory
type UpdateInventoryCommand struct {
	StoreID        int64  `json:"store_id"`
	ProductID      int64  `json:"product_id"`
	SKU            string `json:"sku"`
	QuantityChange int    `json:"quantity_change"`
}

// ReserveInventoryCommand reserves inventory
type ReserveInventoryCommand struct {
	StoreID   int64 `json:"store_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// ReleaseInventoryCommand releases reserved inventory
type ReleaseInventoryCommand struct {
	StoreID   int64 `json:"store_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// AddressCommand represents an address
type AddressCommand struct {
	Street1    string  `json:"street1"`
	Street2    string  `json:"street2,omitempty"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
	PostalCode string  `json:"postal_code"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
}

// StoreSettingsCommand represents store settings
type StoreSettingsCommand struct {
	AllowPickup           bool               `json:"allow_pickup"`
	AllowShipping         bool               `json:"allow_shipping"`
	AllowBackorder        bool               `json:"allow_backorder"`
	InventoryTracking     bool               `json:"inventory_tracking"`
	DefaultShippingCost   float64            `json:"default_shipping_cost"`
	FreeShippingThreshold float64            `json:"free_shipping_threshold,omitempty"`
	MinOrderAmount        float64            `json:"min_order_amount,omitempty"`
	MaxOrderAmount        float64            `json:"max_order_amount,omitempty"`
	BusinessHours         []BusinessHourCmd  `json:"business_hours,omitempty"`
}

// BusinessHourCmd represents business hours
type BusinessHourCmd struct {
	DayOfWeek int    `json:"day_of_week"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
	IsClosed  bool   `json:"is_closed"`
}
