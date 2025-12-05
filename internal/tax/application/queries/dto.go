package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/tax/domain"
	"github.com/shopspring/decimal"
)

// AddressDTO represents an address for tax calculation
type AddressDTO struct {
	Country       string `json:"country"`
	StateProvince string `json:"stateProvince"`
	County        string `json:"county"`
	City          string `json:"city"`
	PostalCode    string `json:"postalCode"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
}

// TaxableItemDTO represents an item to calculate taxes for
type TaxableItemDTO struct {
	ItemID      string          `json:"itemId"`
	SKU         string          `json:"sku"`
	Description string          `json:"description"`
	Quantity    int             `json:"quantity"`
	UnitPrice   decimal.Decimal `json:"unitPrice"`
	Subtotal    decimal.Decimal `json:"subtotal"`
	TaxCategory string          `json:"taxCategory"`
	IsExempt    bool            `json:"isExempt"`
}

// CalculateTaxRequest represents a request to calculate taxes
type CalculateTaxRequest struct {
	OrderID         *int64           `json:"orderId,omitempty"`
	CustomerID      *string          `json:"customerId,omitempty"`
	ShippingAddress AddressDTO       `json:"shippingAddress"`
	BillingAddress  *AddressDTO      `json:"billingAddress,omitempty"`
	Items           []TaxableItemDTO `json:"items"`
	ShippingAmount  decimal.Decimal  `json:"shippingAmount"`
}

// TaxCalculationResponse represents the response from a tax calculation
type TaxCalculationResponse struct {
	OrderID           *int64                  `json:"orderId,omitempty"`
	Items             []TaxedItemDTO          `json:"items"`
	ShippingTax       decimal.Decimal         `json:"shippingTax"`
	TotalTax          decimal.Decimal         `json:"totalTax"`
	Subtotal          decimal.Decimal         `json:"subtotal"`
	TotalAmount       decimal.Decimal         `json:"totalAmount"`
	EffectiveTaxRate  decimal.Decimal         `json:"effectiveTaxRate"`
	Breakdowns        []TaxBreakdownDTO       `json:"breakdowns"`
	CalculatedAt      time.Time               `json:"calculatedAt"`
	JurisdictionsUsed []string                `json:"jurisdictionsUsed"`
}

// TaxedItemDTO represents an item with calculated taxes
type TaxedItemDTO struct {
	ItemID      string             `json:"itemId"`
	SKU         string             `json:"sku"`
	Quantity    int                `json:"quantity"`
	UnitPrice   decimal.Decimal    `json:"unitPrice"`
	Subtotal    decimal.Decimal    `json:"subtotal"`
	TaxAmount   decimal.Decimal    `json:"taxAmount"`
	TaxCategory string             `json:"taxCategory"`
	Taxes       []AppliedTaxDTO    `json:"taxes"`
}

// AppliedTaxDTO represents a tax that was applied
type AppliedTaxDTO struct {
	JurisdictionCode string          `json:"jurisdictionCode"`
	JurisdictionName string          `json:"jurisdictionName"`
	TaxRateName      string          `json:"taxRateName"`
	TaxType          string          `json:"taxType"`
	Rate             decimal.Decimal `json:"rate"`
	TaxableAmount    decimal.Decimal `json:"taxableAmount"`
	TaxAmount        decimal.Decimal `json:"taxAmount"`
	IsCompound       bool            `json:"isCompound"`
}

// TaxBreakdownDTO represents taxes grouped by jurisdiction
type TaxBreakdownDTO struct {
	JurisdictionCode string          `json:"jurisdictionCode"`
	JurisdictionName string          `json:"jurisdictionName"`
	JurisdictionType string          `json:"jurisdictionType"`
	TotalTaxAmount   decimal.Decimal `json:"totalTaxAmount"`
	Rates            []AppliedTaxDTO `json:"rates"`
}

// EstimateTaxRequest represents a request to estimate taxes
type EstimateTaxRequest struct {
	Address  AddressDTO      `json:"address"`
	Subtotal decimal.Decimal `json:"subtotal"`
}

// EstimateTaxResponse represents the response from a tax estimate
type EstimateTaxResponse struct {
	EstimatedTax     decimal.Decimal `json:"estimatedTax"`
	EffectiveTaxRate decimal.Decimal `json:"effectiveTaxRate"`
}

// TaxJurisdictionDTO represents a tax jurisdiction
type TaxJurisdictionDTO struct {
	ID               int64     `json:"id"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	JurisdictionType string    `json:"jurisdictionType"`
	ParentID         *int64    `json:"parentId,omitempty"`
	Country          string    `json:"country"`
	StateProvince    *string   `json:"stateProvince,omitempty"`
	County           *string   `json:"county,omitempty"`
	City             *string   `json:"city,omitempty"`
	PostalCode       *string   `json:"postalCode,omitempty"`
	IsActive         bool      `json:"isActive"`
	Priority         int       `json:"priority"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// TaxRateDTO represents a tax rate
type TaxRateDTO struct {
	ID                int64            `json:"id"`
	JurisdictionID    int64            `json:"jurisdictionId"`
	Name              string           `json:"name"`
	TaxType           string           `json:"taxType"`
	Rate              decimal.Decimal  `json:"rate"`
	TaxCategory       string           `json:"taxCategory"`
	IsCompound        bool             `json:"isCompound"`
	IsShippingTaxable bool             `json:"isShippingTaxable"`
	MinThreshold      *decimal.Decimal `json:"minThreshold,omitempty"`
	MaxThreshold      *decimal.Decimal `json:"maxThreshold,omitempty"`
	Priority          int              `json:"priority"`
	IsActive          bool             `json:"isActive"`
	StartDate         *time.Time       `json:"startDate,omitempty"`
	EndDate           *time.Time       `json:"endDate,omitempty"`
	CreatedAt         time.Time        `json:"createdAt"`
	UpdatedAt         time.Time        `json:"updatedAt"`
}

// TaxExemptionDTO represents a tax exemption
type TaxExemptionDTO struct {
	ID                   int64      `json:"id"`
	CustomerID           string     `json:"customerId"`
	ExemptionCertificate string     `json:"exemptionCertificate"`
	JurisdictionID       *int64     `json:"jurisdictionId,omitempty"`
	TaxCategory          *string    `json:"taxCategory,omitempty"`
	Reason               string     `json:"reason"`
	IsActive             bool       `json:"isActive"`
	StartDate            *time.Time `json:"startDate,omitempty"`
	EndDate              *time.Time `json:"endDate,omitempty"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
}

// Mapper functions

// ToAddressDomain converts AddressDTO to domain.Address
func ToAddressDomain(dto AddressDTO) domain.Address {
	return domain.Address{
		Country:       dto.Country,
		StateProvince: dto.StateProvince,
		County:        dto.County,
		City:          dto.City,
		PostalCode:    dto.PostalCode,
		AddressLine1:  dto.AddressLine1,
		AddressLine2:  dto.AddressLine2,
	}
}

// ToTaxableItemDomain converts TaxableItemDTO to domain.TaxableItem
func ToTaxableItemDomain(dto TaxableItemDTO) domain.TaxableItem {
	return domain.TaxableItem{
		ItemID:      dto.ItemID,
		SKU:         dto.SKU,
		Description: dto.Description,
		Quantity:    dto.Quantity,
		UnitPrice:   dto.UnitPrice,
		Subtotal:    dto.Subtotal,
		TaxCategory: domain.TaxCategory(dto.TaxCategory),
		IsExempt:    dto.IsExempt,
	}
}

// ToTaxCalculationRequest converts CalculateTaxRequest to domain.TaxCalculationRequest
func ToTaxCalculationRequest(dto CalculateTaxRequest) *domain.TaxCalculationRequest {
	req := domain.NewTaxCalculationRequest(ToAddressDomain(dto.ShippingAddress))
	req.OrderID = dto.OrderID
	req.CustomerID = dto.CustomerID
	req.ShippingAmount = dto.ShippingAmount

	if dto.BillingAddress != nil {
		req.BillingAddress = ToAddressDomain(*dto.BillingAddress)
	}

	for _, item := range dto.Items {
		req.AddItem(ToTaxableItemDomain(item))
	}

	return req
}

// ToTaxCalculationResponse converts domain.TaxCalculationResult to TaxCalculationResponse
func ToTaxCalculationResponse(result *domain.TaxCalculationResult) TaxCalculationResponse {
	items := make([]TaxedItemDTO, len(result.Items))
	for i, item := range result.Items {
		items[i] = ToTaxedItemDTO(item)
	}

	breakdowns := make([]TaxBreakdownDTO, len(result.Breakdowns))
	for i, breakdown := range result.Breakdowns {
		breakdowns[i] = ToTaxBreakdownDTO(breakdown)
	}

	return TaxCalculationResponse{
		OrderID:           result.OrderID,
		Items:             items,
		ShippingTax:       result.ShippingTax,
		TotalTax:          result.TotalTax,
		Subtotal:          result.Subtotal,
		TotalAmount:       result.TotalAmount,
		EffectiveTaxRate:  result.GetEffectiveTaxRate(),
		Breakdowns:        breakdowns,
		CalculatedAt:      result.CalculatedAt,
		JurisdictionsUsed: result.JurisdictionsUsed,
	}
}

// ToTaxedItemDTO converts domain.TaxedItem to TaxedItemDTO
func ToTaxedItemDTO(item domain.TaxedItem) TaxedItemDTO {
	taxes := make([]AppliedTaxDTO, len(item.Taxes))
	for i, tax := range item.Taxes {
		taxes[i] = ToAppliedTaxDTO(tax)
	}

	return TaxedItemDTO{
		ItemID:      item.ItemID,
		SKU:         item.SKU,
		Quantity:    item.Quantity,
		UnitPrice:   item.UnitPrice,
		Subtotal:    item.Subtotal,
		TaxAmount:   item.TaxAmount,
		TaxCategory: string(item.TaxCategory),
		Taxes:       taxes,
	}
}

// ToAppliedTaxDTO converts domain.AppliedTax to AppliedTaxDTO
func ToAppliedTaxDTO(tax domain.AppliedTax) AppliedTaxDTO {
	return AppliedTaxDTO{
		JurisdictionCode: tax.JurisdictionCode,
		JurisdictionName: tax.JurisdictionName,
		TaxRateName:      tax.TaxRateName,
		TaxType:          string(tax.TaxType),
		Rate:             tax.Rate,
		TaxableAmount:    tax.TaxableAmount,
		TaxAmount:        tax.TaxAmount,
		IsCompound:       tax.IsCompound,
	}
}

// ToTaxBreakdownDTO converts domain.TaxBreakdown to TaxBreakdownDTO
func ToTaxBreakdownDTO(breakdown domain.TaxBreakdown) TaxBreakdownDTO {
	rates := make([]AppliedTaxDTO, len(breakdown.Rates))
	for i, rate := range breakdown.Rates {
		rates[i] = ToAppliedTaxDTO(rate)
	}

	return TaxBreakdownDTO{
		JurisdictionCode: breakdown.JurisdictionCode,
		JurisdictionName: breakdown.JurisdictionName,
		JurisdictionType: string(breakdown.JurisdictionType),
		TotalTaxAmount:   breakdown.TotalTaxAmount,
		Rates:            rates,
	}
}

// ToTaxJurisdictionDTO converts domain.TaxJurisdiction to TaxJurisdictionDTO
func ToTaxJurisdictionDTO(jurisdiction *domain.TaxJurisdiction) TaxJurisdictionDTO {
	return TaxJurisdictionDTO{
		ID:               jurisdiction.ID,
		Code:             jurisdiction.Code,
		Name:             jurisdiction.Name,
		JurisdictionType: string(jurisdiction.JurisdictionType),
		ParentID:         jurisdiction.ParentID,
		Country:          jurisdiction.Country,
		StateProvince:    jurisdiction.StateProvince,
		County:           jurisdiction.County,
		City:             jurisdiction.City,
		PostalCode:       jurisdiction.PostalCode,
		IsActive:         jurisdiction.IsActive,
		Priority:         jurisdiction.Priority,
		CreatedAt:        jurisdiction.CreatedAt,
		UpdatedAt:        jurisdiction.UpdatedAt,
	}
}

// ToTaxRateDTO converts domain.TaxRate to TaxRateDTO
func ToTaxRateDTO(rate *domain.TaxRate) TaxRateDTO {
	return TaxRateDTO{
		ID:                rate.ID,
		JurisdictionID:    rate.JurisdictionID,
		Name:              rate.Name,
		TaxType:           string(rate.TaxType),
		Rate:              rate.Rate,
		TaxCategory:       string(rate.TaxCategory),
		IsCompound:        rate.IsCompound,
		IsShippingTaxable: rate.IsShippingTaxable,
		MinThreshold:      rate.MinThreshold,
		MaxThreshold:      rate.MaxThreshold,
		Priority:          rate.Priority,
		IsActive:          rate.IsActive,
		StartDate:         rate.StartDate,
		EndDate:           rate.EndDate,
		CreatedAt:         rate.CreatedAt,
		UpdatedAt:         rate.UpdatedAt,
	}
}

// ToTaxExemptionDTO converts domain.TaxExemption to TaxExemptionDTO
func ToTaxExemptionDTO(exemption *domain.TaxExemption) TaxExemptionDTO {
	var taxCategory *string
	if exemption.TaxCategory != nil {
		cat := string(*exemption.TaxCategory)
		taxCategory = &cat
	}

	return TaxExemptionDTO{
		ID:                   exemption.ID,
		CustomerID:           *exemption.CustomerID,
		ExemptionCertificate: exemption.ExemptionCertificate,
		JurisdictionID:       exemption.JurisdictionID,
		TaxCategory:          taxCategory,
		Reason:               exemption.Reason,
		IsActive:             exemption.IsActive,
		StartDate:            exemption.StartDate,
		EndDate:              exemption.EndDate,
		CreatedAt:            exemption.CreatedAt,
		UpdatedAt:            exemption.UpdatedAt,
	}
}
