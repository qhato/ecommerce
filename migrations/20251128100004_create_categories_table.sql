CREATE TABLE IF NOT EXISTS blc_category (
    category_id BIGSERIAL PRIMARY KEY,
    active_end_date TIMESTAMP NULL,
    active_start_date TIMESTAMP NULL,
    archived BPCHAR(1) NULL,
    description VARCHAR(255) NULL,
    display_template VARCHAR(255) NULL,
    external_id VARCHAR(255) NULL,
    fulfillment_type VARCHAR(255) NULL,
    inventory_type VARCHAR(255) NULL,
    long_description TEXT NULL,
    meta_desc VARCHAR(255) NULL,
    meta_title VARCHAR(255) NULL,
    name VARCHAR(255) NOT NULL,
    override_generated_url BOOLEAN NULL,
    product_desc_pattern_override VARCHAR(255) NULL,
    product_title_pattern_override VARCHAR(255) NULL,
    root_display_order NUMERIC(10, 6) NULL,
    tax_code VARCHAR(255) NULL,
    url VARCHAR(255) NULL,
    url_key VARCHAR(255) NULL,
    default_parent_category_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blc_category_default_parent_category_id FOREIGN KEY (default_parent_category_id) REFERENCES blc_category(category_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_category_external_id ON blc_category (external_id);
CREATE INDEX IF NOT EXISTS idx_blc_category_name ON blc_category (name);
CREATE INDEX IF NOT EXISTS idx_blc_category_default_parent_category_id ON blc_category (default_parent_category_id);
CREATE INDEX IF NOT EXISTS idx_blc_category_url ON blc_category (url);
CREATE INDEX IF NOT EXISTS idx_blc_category_url_key ON blc_category (url_key);
