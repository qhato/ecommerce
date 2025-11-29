package domain

import (
	"time"
)

// FulfillmentGroup represents a group of order items to be fulfilled together (e.g., for different shipping addresses)
type FulfillmentGroup struct {
	ID           int64
	OrderID      int64
	Type         string    // e.g., "PHYSICAL_GOODS", "DIGITAL_GOODS" (from blc_fulfillment_group.type)
	ShippingPrice float64   // From blc_fulfillment_group.price (repurposed as shipping price)
	ShippingPriceTaxable bool // From blc_fulfillment_group.shipping_price_taxable
	MerchandiseTotal float64 // From blc_fulfillment_group.merchandise_total
	Method       string    // From blc_fulfillment_group.method
	IsPrimary    bool      // From blc_fulfillment_group.is_primary
	ReferenceNumber string   // From blc_fulfillment_group.reference_number
	RetailPrice  float64   // From blc_fulfillment_group.retail_price
	SalePrice    float64   // From blc_fulfillment_group.sale_price
	Sequence     int       // From blc_fulfillment_group.fulfillment_group_sequnce
	Service      string    // From blc_fulfillment_group.service
	ShippingOverride bool    // From blc_fulfillment_group.shipping_override
	Status       string    // e.g., "PENDING", "FULFILLED", "SHIPPED" (from blc_fulfillment_group.status)
	Total        float64   // From blc_fulfillment_group.total
	TotalFeeTax  float64   // From blc_fulfillment_group.total_fee_tax
	TotalFgTax   float64   // From blc_fulfillment_group.total_fg_tax
	TotalItemTax float64   // From blc_fulfillment_group.total_item_tax
	TotalTax     float64   // From blc_fulfillment_group.total_tax

	AddressID          *int64 // Reference to shipping address (from blc_fulfillment_group.address_id)
	FulfillmentOptionID *int64 // From blc_fulfillment_group.fulfillment_option_id
	PersonalMessageID  *int64 // From blc_fulfillment_group.personal_message_id
	PhoneID            *int64 // From blc_fulfillment_group.phone_id

	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewFulfillmentGroup creates a new FulfillmentGroup
func NewFulfillmentGroup(orderID int64, fgType string) (*FulfillmentGroup, error) {
	if orderID == 0 {
		return nil, NewDomainError("OrderID cannot be zero for FulfillmentGroup")
	}
	if fgType == "" {
		return nil, NewDomainError("FulfillmentGroup Type cannot be empty")
	}

	now := time.Now()
	return &FulfillmentGroup{
		OrderID:             orderID,
		Type:                fgType,
		ShippingPrice:       0.0,
		ShippingPriceTaxable: false,
		MerchandiseTotal:    0.0,
		Method:              "",
		IsPrimary:           false,
		ReferenceNumber:     "",
		RetailPrice:         0.0,
		SalePrice:           0.0,
		Sequence:            0,
		Service:             "",
		ShippingOverride:    false,
		Status:              "PENDING",
		Total:               0.0,
		TotalFeeTax:         0.0,
		TotalFgTax:          0.0,
		TotalItemTax:        0.0,
		TotalTax:            0.0,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

// UpdateShippingDetails updates shipping-related fields
func (fg *FulfillmentGroup) UpdateShippingDetails(shippingPrice, merchandiseTotal float64, method, service string, isTaxable, shippingOverride bool) {
	fg.ShippingPrice = shippingPrice
	fg.MerchandiseTotal = merchandiseTotal
	fg.Method = method
	fg.Service = service
	fg.ShippingPriceTaxable = isTaxable
	fg.ShippingOverride = shippingOverride
	fg.UpdatedAt = time.Now()
}

// SetAddress sets the shipping address ID
func (fg *FulfillmentGroup) SetAddress(addressID int64) {
	fg.AddressID = &addressID
	fg.UpdatedAt = time.Now()
}

// SetFulfillmentOption sets the fulfillment option ID
func (fg *FulfillmentGroup) SetFulfillmentOption(optionID int64) {
	fg.FulfillmentOptionID = &optionID
	fg.UpdatedAt = time.Now()
}

// SetPersonalMessage sets the personal message ID
func (fg *FulfillmentGroup) SetPersonalMessage(messageID int64) {
	fg.PersonalMessageID = &messageID
	fg.UpdatedAt = time.Now()
}

// SetPhone sets the phone ID
func (fg *FulfillmentGroup) SetPhone(phoneID int64) {
	fg.PhoneID = &phoneID
	fg.UpdatedAt = time.Now()
}

// UpdateStatus updates the fulfillment group status
func (fg *FulfillmentGroup) UpdateStatus(status string) {
	fg.Status = status
	fg.UpdatedAt = time.Now()
}

// CalculateTotals recalculates totals for the fulfillment group
func (fg *FulfillmentGroup) CalculateTotals(totalItemTax, totalFeeTax, totalFgTax float64) {
	fg.TotalItemTax = totalItemTax
	fg.TotalFeeTax = totalFeeTax
	fg.TotalFgTax = totalFgTax
	fg.TotalTax = totalItemTax + totalFeeTax + totalFgTax
	fg.Total = fg.MerchandiseTotal + fg.ShippingPrice + fg.TotalTax // Simplified total calculation
	fg.UpdatedAt = time.Now()
}
