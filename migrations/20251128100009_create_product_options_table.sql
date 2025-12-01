CREATE TABLE IF NOT EXISTS blc_product_option (
    product_option_id BIGSERIAL PRIMARY KEY,
    attribute_name VARCHAR(255) NULL,
    display_order INT NULL,
    error_code VARCHAR(255) NULL,
    error_message VARCHAR(255) NULL,
    label VARCHAR(255) NULL,
    long_description TEXT NULL,
    name VARCHAR(255) NULL,
    validation_strategy_type VARCHAR(255) NULL,
    validation_type VARCHAR(255) NULL,
    required BOOLEAN NULL,
    option_type VARCHAR(255) NULL,
    use_in_sku_generation BOOLEAN NULL,
    validation_string VARCHAR(255) NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_blc_product_option_name ON blc_product_option (name);
CREATE INDEX IF NOT EXISTS idx_blc_product_option_option_type ON blc_product_option (option_type);
