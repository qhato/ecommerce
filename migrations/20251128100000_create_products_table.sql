CREATE TABLE IF NOT EXISTS blc_product (
    product_id BIGSERIAL PRIMARY KEY,
    archived BPCHAR(1) NULL,
    can_sell_without_options BOOLEAN NULL,
    canonical_url VARCHAR(255) NULL,
    display_template VARCHAR(255) NULL,
    enable_default_sku_in_inventory BOOLEAN NULL,
    manufacture VARCHAR(255) NULL,
    meta_desc VARCHAR(255) NULL,
    meta_title VARCHAR(255) NULL,
    model VARCHAR(255) NULL,
    override_generated_url BOOLEAN NULL,
    url VARCHAR(255) NULL,
    url_key VARCHAR(255) NULL,
    default_category_id BIGINT NULL,
    default_sku_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk57aoxhpvwg389v7sx4m153mde FOREIGN KEY (default_category_id) REFERENCES blc_category(category_id) -- Will be added later when blc_category is created
);

-- Add indexes for frequently queried columns
CREATE INDEX IF NOT EXISTS idx_blc_product_url_key ON blc_product (url_key);
CREATE INDEX IF NOT EXISTS idx_blc_product_default_category_id ON blc_product (default_category_id);
CREATE INDEX IF NOT EXISTS idx_blc_product_default_sku_id ON blc_product (default_sku_id);
