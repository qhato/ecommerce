-- Create rule table
CREATE TABLE IF NOT EXISTS blc_rule (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    priority INT NOT NULL DEFAULT 0,
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rule_type ON blc_rule(type, status);
CREATE INDEX idx_rule_status ON blc_rule(status, priority);
CREATE INDEX idx_rule_dates ON blc_rule(start_date, end_date);

-- Create rule conditions table
CREATE TABLE IF NOT EXISTS blc_rule_condition (
    id BIGSERIAL PRIMARY KEY,
    rule_id BIGINT NOT NULL REFERENCES blc_rule(id) ON DELETE CASCADE,
    field VARCHAR(255) NOT NULL,
    operator VARCHAR(50) NOT NULL,
    value TEXT NOT NULL,
    logic_operator VARCHAR(10) NOT NULL DEFAULT 'AND',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rule_condition_rule ON blc_rule_condition(rule_id, sort_order);

-- Create rule actions table
CREATE TABLE IF NOT EXISTS blc_rule_action (
    id BIGSERIAL PRIMARY KEY,
    rule_id BIGINT NOT NULL REFERENCES blc_rule(id) ON DELETE CASCADE,
    action_type VARCHAR(100) NOT NULL,
    parameters JSONB,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rule_action_rule ON blc_rule_action(rule_id, sort_order);

COMMENT ON TABLE blc_rule IS 'Business rules for pricing, promotions, inventory, etc.';
COMMENT ON TABLE blc_rule_condition IS 'Conditions that must be met for rule to apply';
COMMENT ON TABLE blc_rule_action IS 'Actions to execute when rule conditions are met';
COMMENT ON COLUMN blc_rule.type IS 'Rule type: PRICE, PROMOTION, INVENTORY, TAX, SHIPPING, CUSTOM';
COMMENT ON COLUMN blc_rule_condition.operator IS 'Comparison operator: EQUALS, GREATER_THAN, LESS_THAN, CONTAINS, etc.';
