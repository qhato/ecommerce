CREATE TABLE IF NOT EXISTS blc_product_option_xref (
    product_option_xref_id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    product_option_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blc_product_option_xref_product_id FOREIGN KEY (product_id) REFERENCES blc_product(product_id),
    -- CONSTRAINT fk_blc_product_option_xref_product_option_id FOREIGN KEY (product_option_id) REFERENCES blc_product_option(product_option_id) -- Will be added later when blc_product_option is created
    CONSTRAINT uk_blc_product_option_xref UNIQUE (product_id, product_option_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_product_option_xref_product_id ON blc_product_option_xref (product_id);
CREATE INDEX IF NOT EXISTS idx_blc_product_option_xref_product_option_id ON blc_product_option_xref (product_option_id);
