CREATE TABLE IF NOT EXISTS blc_product_attribute (
    product_attribute_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NULL,
    product_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blc_product_attribute_product_id FOREIGN KEY (product_id) REFERENCES blc_product(product_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_product_attribute_product_id ON blc_product_attribute (product_id);
CREATE INDEX IF NOT EXISTS idx_blc_product_attribute_name ON blc_product_attribute (name);
