CREATE TABLE IF NOT EXISTS blc_sku_attribute (
    sku_attr_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    sku_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blc_sku_attribute_sku_id FOREIGN KEY (sku_id) REFERENCES blc_sku(sku_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_sku_attribute_sku_id ON blc_sku_attribute (sku_id);
CREATE INDEX IF NOT EXISTS idx_blc_sku_attribute_name ON blc_sku_attribute (name);
