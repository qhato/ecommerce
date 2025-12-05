-- Create checkout_session table
CREATE TABLE IF NOT EXISTS blc_checkout_session (
    id VARCHAR(100) PRIMARY KEY,
    order_id BIGINT NOT NULL,
    customer_id VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    is_guest_checkout BOOLEAN NOT NULL DEFAULT false,
    state VARCHAR(50) NOT NULL,
    current_step INT NOT NULL DEFAULT 1,
    completed_steps TEXT[] NOT NULL DEFAULT '{}',
    shipping_address_id BIGINT,
    billing_address_id BIGINT,
    shipping_method_id VARCHAR(100),
    payment_method_id BIGINT,
    subtotal NUMERIC(19, 5) NOT NULL DEFAULT 0,
    shipping_cost NUMERIC(19, 5) NOT NULL DEFAULT 0,
    tax_amount NUMERIC(19, 5) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(19, 5) NOT NULL DEFAULT 0,
    total_amount NUMERIC(19, 5) NOT NULL DEFAULT 0,
    coupon_codes TEXT[] NOT NULL DEFAULT '{}',
    customer_notes TEXT,
    session_data JSONB,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_activity_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    submitted_at TIMESTAMP WITH TIME ZONE,
    confirmed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_checkout_session_order ON blc_checkout_session(order_id);
CREATE INDEX idx_checkout_session_customer ON blc_checkout_session(customer_id);
CREATE INDEX idx_checkout_session_email ON blc_checkout_session(email);
CREATE INDEX idx_checkout_session_state ON blc_checkout_session(state);
CREATE INDEX idx_checkout_session_expires ON blc_checkout_session(expires_at) WHERE state NOT IN ('CONFIRMED', 'CANCELLED', 'EXPIRED');
CREATE INDEX idx_checkout_session_active ON blc_checkout_session(email, state) WHERE state NOT IN ('CONFIRMED', 'CANCELLED', 'EXPIRED');

-- Create shipping_option table
CREATE TABLE IF NOT EXISTS blc_shipping_option (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    carrier VARCHAR(100) NOT NULL,
    service_code VARCHAR(100),
    speed VARCHAR(50) NOT NULL,
    estimated_days_min INT NOT NULL DEFAULT 1,
    estimated_days_max INT NOT NULL DEFAULT 7,
    base_cost NUMERIC(19, 5) NOT NULL DEFAULT 0,
    cost_per_item NUMERIC(19, 5) NOT NULL DEFAULT 0,
    cost_per_weight NUMERIC(19, 5) NOT NULL DEFAULT 0,
    free_shipping_threshold NUMERIC(19, 5),
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_international BOOLEAN NOT NULL DEFAULT false,
    requires_signature BOOLEAN NOT NULL DEFAULT false,
    allowed_countries TEXT[] NOT NULL DEFAULT '{}',
    excluded_countries TEXT[] NOT NULL DEFAULT '{}',
    allowed_states TEXT[] NOT NULL DEFAULT '{}',
    excluded_states TEXT[] NOT NULL DEFAULT '{}',
    tracking_supported BOOLEAN NOT NULL DEFAULT true,
    insurance_included BOOLEAN NOT NULL DEFAULT false,
    priority INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_shipping_cost_positive CHECK (base_cost >= 0),
    CONSTRAINT chk_estimated_days CHECK (estimated_days_max >= estimated_days_min)
);

CREATE INDEX idx_shipping_option_carrier ON blc_shipping_option(carrier, is_active);
CREATE INDEX idx_shipping_option_active ON blc_shipping_option(is_active, priority);
CREATE INDEX idx_shipping_option_countries ON blc_shipping_option USING GIN(allowed_countries);

-- Add comments
COMMENT ON TABLE blc_checkout_session IS 'Checkout sessions tracking multi-step checkout progress';
COMMENT ON TABLE blc_shipping_option IS 'Available shipping methods with costs and constraints';

COMMENT ON COLUMN blc_checkout_session.state IS 'Current checkout state: INITIATED, CUSTOMER_INFO_ADDED, SHIPPING_INFO_ADDED, etc.';
COMMENT ON COLUMN blc_checkout_session.completed_steps IS 'Array of completed step names';
COMMENT ON COLUMN blc_checkout_session.session_data IS 'Additional flexible session data stored as JSON';
COMMENT ON COLUMN blc_checkout_session.expires_at IS 'Session expiration timestamp (default 24 hours)';

COMMENT ON COLUMN blc_shipping_option.speed IS 'Shipping speed: STANDARD, EXPEDITED, OVERNIGHT, TWO_DAY, SAME_DAY';
COMMENT ON COLUMN blc_shipping_option.cost_per_item IS 'Additional cost per item in cart';
COMMENT ON COLUMN blc_shipping_option.cost_per_weight IS 'Additional cost per weight unit (kg or lb)';
COMMENT ON COLUMN blc_shipping_option.free_shipping_threshold IS 'Minimum order value for free shipping';
COMMENT ON COLUMN blc_shipping_option.allowed_countries IS 'Allowed countries (empty = all countries allowed)';
COMMENT ON COLUMN blc_shipping_option.excluded_countries IS 'Excluded countries';
