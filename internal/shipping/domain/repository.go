package domain

import "context"

type CarrierConfigRepository interface {
	Create(ctx context.Context, config *CarrierConfig) error
	Update(ctx context.Context, config *CarrierConfig) error
	FindByID(ctx context.Context, id int64) (*CarrierConfig, error)
	FindByCarrier(ctx context.Context, carrier ShippingCarrier) (*CarrierConfig, error)
	FindAll(ctx context.Context, enabledOnly bool) ([]*CarrierConfig, error)
}

type ShippingMethodRepository interface {
	Create(ctx context.Context, method *ShippingMethod) error
	Update(ctx context.Context, method *ShippingMethod) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*ShippingMethod, error)
	FindByCarrier(ctx context.Context, carrier ShippingCarrier) ([]*ShippingMethod, error)
	FindAllEnabled(ctx context.Context) ([]*ShippingMethod, error)
}

type ShippingBandRepository interface {
	Create(ctx context.Context, band *ShippingBand) error
	Delete(ctx context.Context, id int64) error
	FindByMethodID(ctx context.Context, methodID int64) ([]*ShippingBand, error)
	DeleteByMethodID(ctx context.Context, methodID int64) error
}

type ShippingRuleRepository interface {
	Create(ctx context.Context, rule *ShippingRule) error
	Update(ctx context.Context, rule *ShippingRule) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*ShippingRule, error)
	FindAllEnabled(ctx context.Context) ([]*ShippingRule, error)
}
