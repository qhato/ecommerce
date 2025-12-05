-- Create tax_jurisdiction table
CREATE TABLE IF NOT EXISTS blc_tax_jurisdiction (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    jurisdiction_type VARCHAR(50) NOT NULL,
    parent_id BIGINT REFERENCES blc_tax_jurisdiction(id) ON DELETE SET NULL,
    country VARCHAR(2) NOT NULL,
    state_province VARCHAR(100),
    county VARCHAR(100),
    city VARCHAR(100),
    postal_code VARCHAR(20),
    is_active BOOLEAN NOT NULL DEFAULT true,
    priority INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tax_jurisdiction_code ON blc_tax_jurisdiction(code);
CREATE INDEX idx_tax_jurisdiction_country ON blc_tax_jurisdiction(country, is_active);
CREATE INDEX idx_tax_jurisdiction_location ON blc_tax_jurisdiction(country, state_province, county, city, postal_code) WHERE is_active=true;
CREATE INDEX idx_tax_jurisdiction_parent ON blc_tax_jurisdiction(parent_id);
CREATE INDEX idx_tax_jurisdiction_active ON blc_tax_jurisdiction(is_active, priority);

-- Create tax_rate table
CREATE TABLE IF NOT EXISTS blc_tax_rate (
    id BIGSERIAL PRIMARY KEY,
    jurisdiction_id BIGINT NOT NULL REFERENCES blc_tax_jurisdiction(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    tax_type VARCHAR(50) NOT NULL,
    rate NUMERIC(19, 5) NOT NULL,
    tax_category VARCHAR(50) NOT NULL,
    is_compound BOOLEAN NOT NULL DEFAULT false,
    is_shipping_taxable BOOLEAN NOT NULL DEFAULT true,
    min_threshold NUMERIC(19, 5),
    max_threshold NUMERIC(19, 5),
    priority INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_rate_non_negative CHECK (rate >= 0),
    CONSTRAINT chk_threshold_range CHECK (max_threshold IS NULL OR min_threshold IS NULL OR max_threshold >= min_threshold)
);

CREATE INDEX idx_tax_rate_jurisdiction ON blc_tax_rate(jurisdiction_id);
CREATE INDEX idx_tax_rate_category ON blc_tax_rate(tax_category, is_active);
CREATE INDEX idx_tax_rate_jurisdiction_category ON blc_tax_rate(jurisdiction_id, tax_category, is_active) WHERE is_active=true;
CREATE INDEX idx_tax_rate_active ON blc_tax_rate(is_active, priority);
CREATE INDEX idx_tax_rate_dates ON blc_tax_rate(start_date, end_date) WHERE is_active=true;

-- Create tax_exemption table
CREATE TABLE IF NOT EXISTS blc_tax_exemption (
    id BIGSERIAL PRIMARY KEY,
    customer_id VARCHAR(255) NOT NULL,
    exemption_certificate VARCHAR(255) NOT NULL UNIQUE,
    jurisdiction_id BIGINT REFERENCES blc_tax_jurisdiction(id) ON DELETE SET NULL,
    tax_category VARCHAR(50),
    reason TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_exemption_date_range CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date)
);

CREATE INDEX idx_tax_exemption_customer ON blc_tax_exemption(customer_id, is_active);
CREATE INDEX idx_tax_exemption_certificate ON blc_tax_exemption(exemption_certificate);
CREATE INDEX idx_tax_exemption_jurisdiction ON blc_tax_exemption(jurisdiction_id, is_active);
CREATE INDEX idx_tax_exemption_active ON blc_tax_exemption(customer_id, is_active, start_date, end_date) WHERE is_active=true;

-- Add comments
COMMENT ON TABLE blc_tax_jurisdiction IS 'Tax jurisdictions (federal, state, county, city, district levels)';
COMMENT ON TABLE blc_tax_rate IS 'Tax rates for different jurisdictions and categories';
COMMENT ON TABLE blc_tax_exemption IS 'Tax exemptions for customers with certificates';

COMMENT ON COLUMN blc_tax_jurisdiction.code IS 'Unique jurisdiction code (e.g., US-CA, US-CA-SF)';
COMMENT ON COLUMN blc_tax_jurisdiction.jurisdiction_type IS 'FEDERAL, STATE, COUNTY, CITY, or DISTRICT';
COMMENT ON COLUMN blc_tax_jurisdiction.parent_id IS 'Parent jurisdiction (e.g., state for a city)';
COMMENT ON COLUMN blc_tax_jurisdiction.priority IS 'Order of application (lower = applied first)';

COMMENT ON COLUMN blc_tax_rate.tax_type IS 'PERCENTAGE, FLAT, or COMPOUND';
COMMENT ON COLUMN blc_tax_rate.tax_category IS 'GENERAL, FOOD, CLOTHING, DIGITAL, SHIPPING, SERVICE, EXEMPT';
COMMENT ON COLUMN blc_tax_rate.is_compound IS 'Whether tax is calculated on subtotal + previous taxes';
COMMENT ON COLUMN blc_tax_rate.is_shipping_taxable IS 'Whether this rate applies to shipping charges';
COMMENT ON COLUMN blc_tax_rate.min_threshold IS 'Minimum amount for tax to apply';
COMMENT ON COLUMN blc_tax_rate.max_threshold IS 'Maximum amount for tax to apply';
COMMENT ON COLUMN blc_tax_rate.priority IS 'Order of application (lower = applied first)';

COMMENT ON COLUMN blc_tax_exemption.customer_id IS 'Customer ID who has the exemption';
COMMENT ON COLUMN blc_tax_exemption.exemption_certificate IS 'Unique certificate number';
COMMENT ON COLUMN blc_tax_exemption.jurisdiction_id IS 'Specific jurisdiction (NULL = all jurisdictions)';
COMMENT ON COLUMN blc_tax_exemption.tax_category IS 'Specific category (NULL = all categories)';
COMMENT ON COLUMN blc_tax_exemption.start_date IS 'Date when exemption becomes active';
COMMENT ON COLUMN blc_tax_exemption.end_date IS 'Date when exemption expires';
