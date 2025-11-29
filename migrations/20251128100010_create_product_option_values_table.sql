CREATE TABLE IF NOT EXISTS blc_product_option_value (
    product_option_value_id BIGSERIAL PRIMARY KEY,
    attribute_value VARCHAR(255) NULL,
    display_order INT NULL,
    price_adjustment NUMERIC(19, 5) NULL,
    product_option_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blc_product_option_value_product_option_id FOREIGN KEY (product_option_id) REFERENCES blc_product_option(product_option_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_product_option_value_product_option_id ON blc_product_option_value (product_option_id);
CREATE INDEX IF NOT EXISTS idx_blc_product_option_value_attribute_value ON blc_product_option_value (attribute_value);
