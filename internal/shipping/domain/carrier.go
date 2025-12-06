package domain

import "time"

// ShippingCarrier represents a shipping carrier
type ShippingCarrier string

const (
	CarrierUSPS   ShippingCarrier = "USPS"
	CarrierUPS    ShippingCarrier = "UPS"
	CarrierFedEx  ShippingCarrier = "FedEx"
	CarrierDHL    ShippingCarrier = "DHL"
	CarrierCustom ShippingCarrier = "CUSTOM"
)

// CarrierConfig represents carrier configuration
type CarrierConfig struct {
	ID          int64
	Carrier     ShippingCarrier
	Name        string
	IsEnabled   bool
	Priority    int
	APIKey      string
	APISecret   string
	AccountID   string
	Config      map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewCarrierConfig creates a new carrier configuration
func NewCarrierConfig(carrier ShippingCarrier, name string) *CarrierConfig {
	now := time.Now()
	return &CarrierConfig{
		Carrier:   carrier,
		Name:      name,
		IsEnabled: false,
		Priority:  0,
		Config:    make(map[string]string),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Enable enables the carrier
func (c *CarrierConfig) Enable() {
	c.IsEnabled = true
	c.UpdatedAt = time.Now()
}

// Disable disables the carrier
func (c *CarrierConfig) Disable() {
	c.IsEnabled = false
	c.UpdatedAt = time.Now()
}
