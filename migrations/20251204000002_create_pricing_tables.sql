-- Create price_list table
CREATE TABLE IF NOT EXISTS blc_price_list (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    price_list_type VARCHAR(50) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    priority INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    description TEXT,
    customer_segments TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_price_list_code ON blc_price_list(code);
CREATE INDEX idx_price_list_active ON blc_price_list(is_active, currency, priority DESC) WHERE is_active=true;
CREATE INDEX idx_price_list_customer_segments ON blc_price_list USING GIN(customer_segments);

-- Create price_list_item table
CREATE TABLE IF NOT EXISTS blc_price_list_item (
    id BIGSERIAL PRIMARY KEY,
    price_list_id BIGINT NOT NULL REFERENCES blc_price_list(id) ON DELETE CASCADE,
    sku_id VARCHAR(255) NOT NULL,
    product_id VARCHAR(255),
    price NUMERIC(19, 5) NOT NULL,
    compare_at_price NUMERIC(19, 5),
    min_quantity INT NOT NULL DEFAULT 1,
    max_quantity INT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_price_positive CHECK (price >= 0),
    CONSTRAINT chk_min_quantity_positive CHECK (min_quantity >= 0),
    CONSTRAINT chk_max_quantity_gte_min CHECK (max_quantity IS NULL OR max_quantity >= min_quantity)
);

CREATE INDEX idx_price_list_item_price_list ON blc_price_list_item(price_list_id);
CREATE INDEX idx_price_list_item_sku ON blc_price_list_item(sku_id);
CREATE INDEX idx_price_list_item_active ON blc_price_list_item(sku_id, is_active) WHERE is_active=true;
CREATE UNIQUE INDEX idx_price_list_item_unique ON blc_price_list_item(price_list_id, sku_id, min_quantity);

-- Create pricing_rule table
CREATE TABLE IF NOT EXISTS blc_pricing_rule (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    rule_type VARCHAR(50) NOT NULL,
    priority INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    condition_expression TEXT,
    action_type VARCHAR(50) NOT NULL,
    action_value NUMERIC(19, 5) NOT NULL,
    applicable_skus TEXT[] NOT NULL DEFAULT '{}',
    applicable_categories TEXT[] NOT NULL DEFAULT '{}',
    customer_segments TEXT[] NOT NULL DEFAULT '{}',
    min_quantity INT NOT NULL DEFAULT 1,
    max_quantity INT,
    min_order_value NUMERIC(19, 5),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_action_value_non_negative CHECK (action_value >= 0)
);

CREATE INDEX idx_pricing_rule_active ON blc_pricing_rule(is_active, priority DESC) WHERE is_active=true;
CREATE INDEX idx_pricing_rule_skus ON blc_pricing_rule USING GIN(applicable_skus);
CREATE INDEX idx_pricing_rule_categories ON blc_pricing_rule USING GIN(applicable_categories);
CREATE INDEX idx_pricing_rule_segments ON blc_pricing_rule USING GIN(customer_segments);

-- Add comments
COMMENT ON TABLE blc_price_list IS 'Price lists for different customer segments and currencies';
COMMENT ON TABLE blc_price_list_item IS 'Individual price entries for SKUs in price lists';
COMMENT ON TABLE blc_pricing_rule IS 'Dynamic pricing rules for automated price adjustments';

COMMENT ON COLUMN blc_price_list.priority IS 'Higher priority lists take precedence (higher number = higher priority)';
COMMENT ON COLUMN blc_price_list.customer_segments IS 'Customer segments this price list applies to (empty = all customers)';
COMMENT ON COLUMN blc_price_list_item.min_quantity IS 'Minimum quantity required for this price (for tiered pricing)';
COMMENT ON COLUMN blc_price_list_item.max_quantity IS 'Maximum quantity for this price (NULL = no maximum)';
COMMENT ON COLUMN blc_pricing_rule.action_type IS 'FIXED_PRICE, PERCENT_DISCOUNT, AMOUNT_DISCOUNT, PERCENT_SURCHARGE, AMOUNT_SURCHARGE';
COMMENT ON COLUMN blc_pricing_rule.rule_type IS 'QUANTITY_TIERED, VOLUME_DISCOUNT, CUSTOMER_SEGMENT, DYNAMIC, TIME_BASED';
