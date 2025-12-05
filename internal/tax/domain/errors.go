package domain

import "errors"

// Tax Jurisdiction Errors
var (
	ErrJurisdictionCodeRequired     = errors.New("jurisdiction code is required")
	ErrJurisdictionNameRequired     = errors.New("jurisdiction name is required")
	ErrCountryRequired              = errors.New("country is required")
	ErrJurisdictionNotFound         = errors.New("jurisdiction not found")
	ErrJurisdictionAlreadyExists    = errors.New("jurisdiction with this code already exists")
	ErrInvalidJurisdictionType      = errors.New("invalid jurisdiction type")
	ErrJurisdictionIDRequired       = errors.New("jurisdiction ID is required")
)

// Tax Rate Errors
var (
	ErrTaxRateNameRequired          = errors.New("tax rate name is required")
	ErrTaxRateCannotBeNegative      = errors.New("tax rate cannot be negative")
	ErrInvalidTaxRateType           = errors.New("invalid tax rate type")
	ErrInvalidTaxCategory           = errors.New("invalid tax category")
	ErrTaxRateNotFound              = errors.New("tax rate not found")
	ErrTaxRateAlreadyExists         = errors.New("tax rate already exists")
	ErrInvalidThresholdRange        = errors.New("max threshold must be greater than min threshold")
	ErrFlatRateRequiresAmount       = errors.New("flat rate type requires an amount value")
	ErrPercentageRateRequiresRate   = errors.New("percentage rate type requires a rate value")
)

// Tax Calculation Errors
var (
	ErrShippingAddressRequired      = errors.New("shipping address is required")
	ErrInvalidAddress               = errors.New("invalid address provided")
	ErrNoItemsToCalculate           = errors.New("no items provided for tax calculation")
	ErrInvalidTaxableItem           = errors.New("invalid taxable item")
	ErrNegativePrice                = errors.New("price cannot be negative")
	ErrNegativeQuantity             = errors.New("quantity cannot be negative")
	ErrCalculationFailed            = errors.New("tax calculation failed")
	ErrNoApplicableJurisdictions    = errors.New("no applicable jurisdictions found for address")
	ErrNoApplicableTaxRates         = errors.New("no applicable tax rates found")
)

// Tax Exemption Errors
var (
	ErrCustomerIDRequired           = errors.New("customer ID is required")
	ErrExemptionCertificateRequired = errors.New("exemption certificate is required")
	ErrExemptionNotFound            = errors.New("tax exemption not found")
	ErrExemptionExpired             = errors.New("tax exemption has expired")
	ErrExemptionNotActive           = errors.New("tax exemption is not active")
	ErrExemptionAlreadyExists       = errors.New("tax exemption already exists for this customer")
	ErrInvalidExemptionDateRange    = errors.New("end date must be after start date")
)

// Repository Errors
var (
	ErrRepositoryOperation          = errors.New("repository operation failed")
	ErrTransactionFailed            = errors.New("database transaction failed")
	ErrDuplicateEntry               = errors.New("duplicate entry")
	ErrConcurrentUpdate             = errors.New("concurrent update detected")
)
