package application

import (
	"time"

	"github.com/qhato/ecommerce/internal/customer/domain"
)

// CustomerDTO represents a customer data transfer object
type CustomerDTO struct {
	ID                     int64              `json:"id"`
	EmailAddress           string             `json:"email_address"`
	UserName               string             `json:"user_name"`
	FirstName              string             `json:"first_name"`
	LastName               string             `json:"last_name"`
	FullName               string             `json:"full_name"`
	Archived               bool               `json:"archived"`
	Deactivated            bool               `json:"deactivated"`
	IsTaxExempt            bool               `json:"is_tax_exempt"`
	TaxExemptionCode       string             `json:"tax_exemption_code,omitempty"`
	PasswordChangeRequired bool               `json:"password_change_required"`
	ReceiveEmail           bool               `json:"receive_email"`
	IsRegistered           bool               `json:"is_registered"`
	IsActive               bool               `json:"is_active"`
	Attributes             map[string]string  `json:"attributes,omitempty"`
	Roles                  []string           `json:"roles,omitempty"`
	CreatedAt              time.Time          `json:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at"`
}

// AddressDTO represents an address data transfer object
type AddressDTO struct {
	ID                  int64  `json:"id"`
	AddressLine1        string `json:"address_line1"`
	AddressLine2        string `json:"address_line2,omitempty"`
	AddressLine3        string `json:"address_line3,omitempty"`
	City                string `json:"city"`
	StateProvinceRegion string `json:"state_province_region,omitempty"`
	PostalCode          string `json:"postal_code"`
	CountryCode         string `json:"country_code"`
	CompanyName         string `json:"company_name,omitempty"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	PrimaryPhone        string `json:"primary_phone,omitempty"`
}

// CustomerAddressDTO represents a customer address data transfer object
type CustomerAddressDTO struct {
	ID          int64       `json:"id"`
	AddressName string      `json:"address_name"`
	Archived    bool        `json:"archived"`
	Address     *AddressDTO `json:"address"`
}

// ToCustomerDTO converts a domain Customer to CustomerDTO
func ToCustomerDTO(customer *domain.Customer) *CustomerDTO {
	attributes := make(map[string]string)
	for _, attr := range customer.Attributes {
		attributes[attr.Name] = attr.Value
	}

	roles := make([]string, len(customer.Roles))
	for i, role := range customer.Roles {
		roles[i] = role.RoleName
	}

	return &CustomerDTO{
		ID:                     customer.ID,
		EmailAddress:           customer.EmailAddress,
		UserName:               customer.UserName,
		FirstName:              customer.FirstName,
		LastName:               customer.LastName,
		FullName:               customer.GetFullName(),
		Archived:               customer.Archived,
		Deactivated:            customer.Deactivated,
		IsTaxExempt:            customer.IsTaxExempt,
		TaxExemptionCode:       customer.TaxExemptionCode,
		PasswordChangeRequired: customer.PasswordChangeRequired,
		ReceiveEmail:           customer.ReceiveEmail,
		IsRegistered:           customer.IsRegistered,
		IsActive:               customer.IsActive(),
		Attributes:             attributes,
		Roles:                  roles,
		CreatedAt:              customer.CreatedAt,
		UpdatedAt:              customer.UpdatedAt,
	}
}

// ToAddressDTO converts a domain Address to AddressDTO
func ToAddressDTO(address *domain.Address) *AddressDTO {
	return &AddressDTO{
		ID:                  address.ID,
		AddressLine1:        address.AddressLine1,
		AddressLine2:        address.AddressLine2,
		AddressLine3:        address.AddressLine3,
		City:                address.City,
		StateProvinceRegion: address.StateProvinceRegion,
		PostalCode:          address.PostalCode,
		CountryCode:         address.CountryCode,
		CompanyName:         address.CompanyName,
		FirstName:           address.FirstName,
		LastName:            address.LastName,
		PrimaryPhone:        address.PrimaryPhone,
	}
}

// ToCustomerAddressDTO converts a domain CustomerAddress to CustomerAddressDTO
func ToCustomerAddressDTO(customerAddress *domain.CustomerAddress) *CustomerAddressDTO {
	var addressDTO *AddressDTO
	if customerAddress.Address != nil {
		addressDTO = ToAddressDTO(customerAddress.Address)
	}

	return &CustomerAddressDTO{
		ID:          customerAddress.ID,
		AddressName: customerAddress.AddressName,
		Archived:    customerAddress.Archived,
		Address:     addressDTO,
	}
}
