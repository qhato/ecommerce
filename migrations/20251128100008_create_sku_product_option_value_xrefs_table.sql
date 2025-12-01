CREATE TABLE IF NOT EXISTS blc_sku_option_value_xref (
    sku_option_value_xref_id BIGSERIAL PRIMARY KEY,
    sku_id BIGINT NOT NULL,
    product_option_value_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blc_sku_option_value_xref_sku_id FOREIGN KEY (sku_id) REFERENCES blc_sku(sku_id),
    -- CONSTRAINT fk_blc_sku_option_value_xref_product_option_value_id FOREIGN KEY (product_option_value_id) REFERENCES blc_product_option_value(product_option_value_id) -- Will be added later when blc_product_option_value is created
    CONSTRAINT uk_blc_sku_option_value_xref UNIQUE (sku_id, product_option_value_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_sku_option_value_xref_sku_id ON blc_sku_option_value_xref (sku_id);
CREATE INDEX IF NOT EXISTS idx_blc_sku_option_value_xref_product_option_value_id ON blc_sku_option_value_xref (product_option_value_id);
