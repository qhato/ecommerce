package fulfillment

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/pkg/workflow"
)

// FulfillmentContext contains fulfillment workflow input/output
type FulfillmentContext struct {
	OrderID    int64
	CustomerID int64
	Items      []FulfillmentItem
	ShippingAddress Address
	
	// Workflow state
	ShipmentID        *int64
	TrackingNumber    *string
	ShippingLabelURL  *string
	InventoryAllocated bool
	ShipmentCreated   bool
	LabelGenerated    bool
	
	EstimatedDelivery *time.Time
	Metadata          map[string]interface{}
}

type FulfillmentItem struct {
	SKUID    int64
	Quantity int
}

type Address struct {
	FirstName    string
	LastName     string
	AddressLine1 string
	City         string
	State        string
	PostalCode   string
	Country      string
}

// AllocateInventoryActivity allocates inventory for fulfillment
type AllocateInventoryActivity struct {
	workflow.BaseActivity
	inventoryService InventoryService
}

type InventoryService interface {
	AllocateInventory(ctx context.Context, skuID int64, quantity int) error
	ReleaseInventory(ctx context.Context, skuID int64, quantity int) error
}

func NewAllocateInventoryActivity(inventoryService InventoryService) *AllocateInventoryActivity {
	return &AllocateInventoryActivity{
		BaseActivity:     workflow.NewBaseActivity("AllocateInventory", "Allocate inventory for shipment"),
		inventoryService: inventoryService,
	}
}

func (a *AllocateInventoryActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	fulfillmentCtx, ok := input.(*FulfillmentContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	for _, item := range fulfillmentCtx.Items {
		if err := a.inventoryService.AllocateInventory(ctx, item.SKUID, item.Quantity); err != nil {
			return nil, fmt.Errorf("failed to allocate inventory for SKU %d: %w", item.SKUID, err)
		}
	}

	fulfillmentCtx.InventoryAllocated = true
	return fulfillmentCtx, nil
}

func (a *AllocateInventoryActivity) Compensate(ctx context.Context, input interface{}) error {
	fulfillmentCtx, ok := input.(*FulfillmentContext)
	if !ok || !fulfillmentCtx.InventoryAllocated {
		return nil
	}

	for _, item := range fulfillmentCtx.Items {
		_ = a.inventoryService.ReleaseInventory(ctx, item.SKUID, item.Quantity)
	}

	return nil
}

// CreateShipmentActivity creates a shipment record
type CreateShipmentActivity struct {
	workflow.BaseActivity
	shipmentService ShipmentService
}

type ShipmentService interface {
	CreateShipment(ctx context.Context, orderID int64, items []FulfillmentItem, address Address) (int64, error)
	CancelShipment(ctx context.Context, shipmentID int64) error
	GenerateShippingLabel(ctx context.Context, shipmentID int64) (string, string, error) // labelURL, trackingNumber, error
}

func NewCreateShipmentActivity(shipmentService ShipmentService) *CreateShipmentActivity {
	return &CreateShipmentActivity{
		BaseActivity:    workflow.NewBaseActivity("CreateShipment", "Create shipment record"),
		shipmentService: shipmentService,
	}
}

func (a *CreateShipmentActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	fulfillmentCtx, ok := input.(*FulfillmentContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	shipmentID, err := a.shipmentService.CreateShipment(ctx, fulfillmentCtx.OrderID, fulfillmentCtx.Items, fulfillmentCtx.ShippingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	fulfillmentCtx.ShipmentID = &shipmentID
	fulfillmentCtx.ShipmentCreated = true

	return fulfillmentCtx, nil
}

func (a *CreateShipmentActivity) Compensate(ctx context.Context, input interface{}) error {
	fulfillmentCtx, ok := input.(*FulfillmentContext)
	if !ok || !fulfillmentCtx.ShipmentCreated || fulfillmentCtx.ShipmentID == nil {
		return nil
	}

	return a.shipmentService.CancelShipment(ctx, *fulfillmentCtx.ShipmentID)
}

// GenerateShippingLabelActivity generates shipping label
type GenerateShippingLabelActivity struct {
	workflow.BaseActivity
	shipmentService ShipmentService
}

func NewGenerateShippingLabelActivity(shipmentService ShipmentService) *GenerateShippingLabelActivity {
	return &GenerateShippingLabelActivity{
		BaseActivity:    workflow.NewBaseActivity("GenerateShippingLabel", "Generate shipping label"),
		shipmentService: shipmentService,
	}
}

func (a *GenerateShippingLabelActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	fulfillmentCtx, ok := input.(*FulfillmentContext)
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	if fulfillmentCtx.ShipmentID == nil {
		return nil, fmt.Errorf("shipment not created")
	}

	labelURL, trackingNumber, err := a.shipmentService.GenerateShippingLabel(ctx, *fulfillmentCtx.ShipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate shipping label: %w", err)
	}

	fulfillmentCtx.ShippingLabelURL = &labelURL
	fulfillmentCtx.TrackingNumber = &trackingNumber
	fulfillmentCtx.LabelGenerated = true

	return fulfillmentCtx, nil
}

func (a *GenerateShippingLabelActivity) Compensate(ctx context.Context, input interface{}) error {
	// Label generation is typically not reversible
	// The shipment cancellation will handle cleanup
	return nil
}

// FulfillmentWorkflow creates a fulfillment workflow
func FulfillmentWorkflow(
	inventoryService InventoryService,
	shipmentService ShipmentService,
) (*workflow.Workflow, error) {
	return workflow.NewWorkflowBuilder("fulfillment", "Fulfillment Workflow").
		Description("Process order fulfillment with inventory and shipping").
		AddActivities(
			NewAllocateInventoryActivity(inventoryService),
			NewCreateShipmentActivity(shipmentService),
			NewGenerateShippingLabelActivity(shipmentService),
		).
		MaxRetries(2).
		CompensateOnFail(true).
		Build()
}