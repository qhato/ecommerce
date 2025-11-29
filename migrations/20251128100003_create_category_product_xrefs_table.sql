CREATE TABLE IF NOT EXISTS blc_category_product_xref (
    category_product_id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    default_reference BOOLEAN NULL,
    display_order NUMERIC(10, 6) NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blc_category_product_xref_product_id FOREIGN KEY (product_id) REFERENCES blc_product(product_id),
    -- CONSTRAINT fk_blc_category_product_xref_category_id FOREIGN KEY (category_id) REFERENCES blc_category(category_id) -- Will be added later when blc_category is created
    CONSTRAINT uk_blc_category_product_xref UNIQUE (category_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_category_product_xref_category_id ON blc_category_product_xref (category_id);
CREATE INDEX IF NOT EXISTS idx_blc_category_product_xref_product_id ON blc_category_product_xref (product_id);
