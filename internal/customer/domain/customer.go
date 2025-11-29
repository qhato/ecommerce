package domain

import "time"

// Customer represents a customer entity
type Customer struct {
	ID                     int64
	Archived               bool
	ChallengeAnswer        string
	Deactivated            bool
	EmailAddress           string
	ExternalID             string
	FirstName              string
	IsTaxExempt            bool
	LastName               string
	Password               string
	PasswordChangeRequired bool
	IsPreview              bool
	ReceiveEmail           bool
	IsRegistered           bool
	TaxExemptionCode       string
	UserName               string
	ChallengeQuestionID    *int64
	LocaleCode             string
	Addresses              []CustomerAddress
	Phones                 []CustomerPhone
	Attributes             []CustomerAttribute
	Roles                  []CustomerRole
	CreatedBy              int64
	UpdatedBy              int64
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// CustomerAddress represents a customer address
type CustomerAddress struct {
	ID          int64
	AddressName string
	Archived    bool
	AddressID   int64
	CustomerID  int64
	Address     *Address
}

// Address represents a physical address
type Address struct {
	ID                  int64
	AddressLine1        string
	AddressLine2        string
	AddressLine3        string
	City                string
	CompanyName         string
	County              string
	FirstName           string
	LastName            string
	PrimaryPhone        string
	PostalCode          string
	StateProvinceRegion string
	CountryCode         string
	IsoCountryAlpha2    string
}

// CustomerPhone represents a customer phone number
type CustomerPhone struct {
	ID         int64
	PhoneName  string
	CustomerID int64
	PhoneID    int64
	Phone      *Phone
}

// Phone represents a phone number
type Phone struct {
	ID          int64
	PhoneNumber string
	IsActive    bool
	IsDefault   bool
}

// CustomerAttribute represents a custom attribute
type CustomerAttribute struct {
	ID         int64
	Name       string
	Value      string
	CustomerID int64
}

// CustomerRole represents a customer role
type CustomerRole struct {
	ID         int64
	CustomerID int64
	RoleID     int64
	RoleName   string
}

// NewCustomer creates a new customer
func NewCustomer(emailAddress, userName, password, firstName, lastName string) *Customer {
	now := time.Now()
	return &Customer{
		EmailAddress:           emailAddress,
		UserName:               userName,
		Password:               password,
		FirstName:              firstName,
		LastName:               lastName,
		IsRegistered:           true,
		ReceiveEmail:           true,
		Archived:               false,
		Deactivated:            false,
		PasswordChangeRequired: false,
		CreatedAt:              now,
		UpdatedAt:              now,
		Addresses:              make([]CustomerAddress, 0),
		Phones:                 make([]CustomerPhone, 0),
		Attributes:             make([]CustomerAttribute, 0),
		Roles:                  make([]CustomerRole, 0),
	}
}

// UpdateProfile updates customer profile information
func (c *Customer) UpdateProfile(firstName, lastName, emailAddress string) {
	c.FirstName = firstName
	c.LastName = lastName
	c.EmailAddress = emailAddress
	c.UpdatedAt = time.Now()
}

// ChangePassword updates the customer's password
func (c *Customer) ChangePassword(newPassword string) {
	c.Password = newPassword
	c.PasswordChangeRequired = false
	c.UpdatedAt = time.Now()
}

// Deactivate deactivates the customer account
func (c *Customer) Deactivate() {
	c.Deactivated = true
	c.UpdatedAt = time.Now()
}

// Activate activates the customer account
func (c *Customer) Activate() {
	c.Deactivated = false
	c.UpdatedAt = time.Now()
}

// Archive archives the customer
func (c *Customer) Archive() {
	c.Archived = true
	c.UpdatedAt = time.Now()
}

// AddAttribute adds a custom attribute
func (c *Customer) AddAttribute(name, value string) {
	c.Attributes = append(c.Attributes, CustomerAttribute{
		Name:       name,
		Value:      value,
		CustomerID: c.ID,
	})
	c.UpdatedAt = time.Now()
}

// UpdateAttribute updates an existing attribute
func (c *Customer) UpdateAttribute(name, value string) {
	for i, attr := range c.Attributes {
		if attr.Name == name {
			c.Attributes[i].Value = value
			c.UpdatedAt = time.Now()
			return
		}
	}
	c.AddAttribute(name, value)
}

// GetAttribute retrieves an attribute value by name
func (c *Customer) GetAttribute(name string) (string, bool) {
	for _, attr := range c.Attributes {
		if attr.Name == name {
			return attr.Value, true
		}
	}
	return "", false
}

// AddRole adds a role to the customer
func (c *Customer) AddRole(roleID int64, roleName string) {
	c.Roles = append(c.Roles, CustomerRole{
		CustomerID: c.ID,
		RoleID:     roleID,
		RoleName:   roleName,
	})
	c.UpdatedAt = time.Now()
}

// HasRole checks if customer has a specific role
func (c *Customer) HasRole(roleName string) bool {
	for _, role := range c.Roles {
		if role.RoleName == roleName {
			return true
		}
	}
	return false
}

// IsActive checks if customer is active
func (c *Customer) IsActive() bool {
	return !c.Deactivated && !c.Archived
}

// GetFullName returns the customer's full name
func (c *Customer) GetFullName() string {
	return c.FirstName + " " + c.LastName
}

// SetTaxExempt sets tax exemption status
func (c *Customer) SetTaxExempt(exempt bool, code string) {
	c.IsTaxExempt = exempt
	c.TaxExemptionCode = code
	c.UpdatedAt = time.Now()
}
