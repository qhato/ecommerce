CREATE TABLE IF NOT EXISTS blc_sku (
    sku_id BIGSERIAL PRIMARY KEY,
    active_end_date TIMESTAMP NULL,
    active_start_date TIMESTAMP NULL,
    available_flag BPCHAR(1) NULL,
    cost NUMERIC(19, 5) NULL,
    description VARCHAR(255) NULL,
    container_shape VARCHAR(255) NULL,
    depth NUMERIC(19, 2) NULL,
    dimension_unit_of_measure VARCHAR(255) NULL,
    girth NUMERIC(19, 2) NULL,
    height NUMERIC(19, 2) NULL,
    container_size VARCHAR(255) NULL,
    width NUMERIC(19, 2) NULL,
    discountable_flag BPCHAR(1) NULL,
    display_template VARCHAR(255) NULL,
    external_id VARCHAR(255) NULL,
    fulfillment_type VARCHAR(255) NULL,
    inventory_type VARCHAR(255) NULL,
    is_machine_sortable BOOLEAN NULL,
    long_description TEXT NULL,
    name VARCHAR(255) NULL,
    quantity_available INT NULL, -- This needs careful handling with blc_sku_availability
    retail_price NUMERIC(19, 5) NULL,
    sale_price NUMERIC(19, 5) NULL,
    tax_code VARCHAR(255) NULL,
    taxable_flag BPCHAR(1) NULL,
    upc VARCHAR(255) NULL,
    url_key VARCHAR(255) NULL,
    weight NUMERIC(19, 2) NULL,
    weight_unit_of_measure VARCHAR(255) NULL,
    currency_code VARCHAR(255) NULL,
    default_product_id BIGINT NULL,
    addl_product_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blc_sku_currency_code FOREIGN KEY (currency_code) REFERENCES blc_currency(currency_code),
    CONSTRAINT fk_blc_sku_addl_product_id FOREIGN KEY (addl_product_id) REFERENCES blc_product(product_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_sku_active_end_date ON blc_sku (active_end_date);
CREATE INDEX IF NOT EXISTS idx_blc_sku_active_start_date ON blc_sku (active_start_date);
CREATE INDEX IF NOT EXISTS idx_blc_sku_available_flag ON blc_sku (available_flag);
CREATE INDEX IF NOT EXISTS idx_blc_sku_discountable_flag ON blc_sku (discountable_flag);
CREATE INDEX IF NOT EXISTS idx_blc_sku_external_id ON blc_sku (external_id);
CREATE INDEX IF NOT EXISTS idx_blc_sku_name ON blc_sku (name);
CREATE INDEX IF NOT EXISTS idx_blc_sku_taxable_flag ON blc_sku (taxable_flag);
CREATE INDEX IF NOT EXISTS idx_blc_sku_upc ON blc_sku (upc);
CREATE INDEX IF NOT EXISTS idx_blc_sku_url_key ON blc_sku (url_key);
