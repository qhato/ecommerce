package domain

import (
	"time"

	"github.com/qhato/ecommerce/pkg/event"
)

const (
	// Product events
	EventProductCreated  = "catalog.product.created"
	EventProductUpdated  = "catalog.product.updated"
	EventProductArchived = "catalog.product.archived"
	EventProductDeleted  = "catalog.product.deleted"

	// Category events
	EventCategoryCreated  = "catalog.category.created"
	EventCategoryUpdated  = "catalog.category.updated"
	EventCategoryArchived = "catalog.category.archived"
	EventCategoryDeleted  = "catalog.category.deleted"

	// SKU events
	EventSKUCreated            = "catalog.sku.created"
	EventSKUUpdated            = "catalog.sku.updated"
	EventSKUDeleted            = "catalog.sku.deleted"
	EventSKUAvailabilityChanged = "catalog.sku.availability_changed"
	EventSKUPriceChanged       = "catalog.sku.price_changed"
)

// ProductCreatedEvent is published when a product is created
type ProductCreatedEvent struct {
	event.BaseEvent
	ProductID   int64  `json:"product_id"`
	Model       string `json:"model"`
	Manufacture string `json:"manufacture"`
}

// NewProductCreatedEvent creates a new ProductCreatedEvent
func NewProductCreatedEvent(productID int64, model, manufacture string) *ProductCreatedEvent {
	return &ProductCreatedEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventProductCreated,
			Timestamp: time.Now(),
		},
		ProductID:   productID,
		Model:       model,
		Manufacture: manufacture,
	}
}

// Type returns the event type
func (e *ProductCreatedEvent) Type() string {
	return e.EventType
}

// ProductUpdatedEvent is published when a product is updated
type ProductUpdatedEvent struct {
	event.BaseEvent
	ProductID int64  `json:"product_id"`
	Changes   map[string]interface{} `json:"changes"`
}

// NewProductUpdatedEvent creates a new ProductUpdatedEvent
func NewProductUpdatedEvent(productID int64, changes map[string]interface{}) *ProductUpdatedEvent {
	return &ProductUpdatedEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventProductUpdated,
			Timestamp: time.Now(),
		},
		ProductID: productID,
		Changes:   changes,
	}
}

// Type returns the event type
func (e *ProductUpdatedEvent) Type() string {
	return e.EventType
}

// ProductArchivedEvent is published when a product is archived
type ProductArchivedEvent struct {
	event.BaseEvent
	ProductID int64 `json:"product_id"`
}

// NewProductArchivedEvent creates a new ProductArchivedEvent
func NewProductArchivedEvent(productID int64) *ProductArchivedEvent {
	return &ProductArchivedEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventProductArchived,
			Timestamp: time.Now(),
		},
		ProductID: productID,
	}
}

// Type returns the event type
func (e *ProductArchivedEvent) Type() string {
	return e.EventType
}

// CategoryCreatedEvent is published when a category is created
type CategoryCreatedEvent struct {
	event.BaseEvent
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	ParentID   *int64 `json:"parent_id,omitempty"`
}

// NewCategoryCreatedEvent creates a new CategoryCreatedEvent
func NewCategoryCreatedEvent(categoryID int64, name string, parentID *int64) *CategoryCreatedEvent {
	return &CategoryCreatedEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventCategoryCreated,
			Timestamp: time.Now(),
		},
		CategoryID: categoryID,
		Name:       name,
		ParentID:   parentID,
	}
}

// Type returns the event type
func (e *CategoryCreatedEvent) Type() string {
	return e.EventType
}

// CategoryUpdatedEvent is published when a category is updated
type CategoryUpdatedEvent struct {
	event.BaseEvent
	CategoryID int64                  `json:"category_id"`
	Changes    map[string]interface{} `json:"changes"`
}

// NewCategoryUpdatedEvent creates a new CategoryUpdatedEvent
func NewCategoryUpdatedEvent(categoryID int64, changes map[string]interface{}) *CategoryUpdatedEvent {
	return &CategoryUpdatedEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventCategoryUpdated,
			Timestamp: time.Now(),
		},
		CategoryID: categoryID,
		Changes:    changes,
	}
}

// Type returns the event type
func (e *CategoryUpdatedEvent) Type() string {
	return e.EventType
}

// SKUCreatedEvent is published when a SKU is created
type SKUCreatedEvent struct {
	event.BaseEvent
	SKUID     int64   `json:"sku_id"`
	ProductID *int64  `json:"product_id,omitempty"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
}

// NewSKUCreatedEvent creates a new SKUCreatedEvent
func NewSKUCreatedEvent(skuID int64, productID *int64, name string, price float64) *SKUCreatedEvent {
	return &SKUCreatedEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventSKUCreated,
			Timestamp: time.Now(),
		},
		SKUID:     skuID,
		ProductID: productID,
		Name:      name,
		Price:     price,
	}
}

// Type returns the event type
func (e *SKUCreatedEvent) Type() string {
	return e.EventType
}

// SKUAvailabilityChangedEvent is published when SKU availability changes
type SKUAvailabilityChangedEvent struct {
	event.BaseEvent
	SKUID     int64 `json:"sku_id"`
	Available bool  `json:"available"`
}

// NewSKUAvailabilityChangedEvent creates a new SKUAvailabilityChangedEvent
func NewSKUAvailabilityChangedEvent(skuID int64, available bool) *SKUAvailabilityChangedEvent {
	return &SKUAvailabilityChangedEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventSKUAvailabilityChanged,
			Timestamp: time.Now(),
		},
		SKUID:     skuID,
		Available: available,
	}
}

// Type returns the event type
func (e *SKUAvailabilityChangedEvent) Type() string {
	return e.EventType
}

// SKUPriceChangedEvent is published when SKU price changes
type SKUPriceChangedEvent struct {
	event.BaseEvent
	SKUID    int64   `json:"sku_id"`
	OldPrice float64 `json:"old_price"`
	NewPrice float64 `json:"new_price"`
}

// NewSKUPriceChangedEvent creates a new SKUPriceChangedEvent
func NewSKUPriceChangedEvent(skuID int64, oldPrice, newPrice float64) *SKUPriceChangedEvent {
	return &SKUPriceChangedEvent{
		BaseEvent: event.BaseEvent{
			EventType: EventSKUPriceChanged,
			Timestamp: time.Now(),
		},
		SKUID:    skuID,
		OldPrice: oldPrice,
		NewPrice: newPrice,
	}
}

// Type returns the event type
func (e *SKUPriceChangedEvent) Type() string {
	return e.EventType
}
