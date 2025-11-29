CREATE TABLE IF NOT EXISTS blc_tax_detail (
    tax_detail_id BIGSERIAL PRIMARY KEY,
    amount NUMERIC(19, 5) NULL,
    tax_country VARCHAR(255) NULL,
    jurisdiction_name VARCHAR(255) NULL,
    rate NUMERIC(19, 5) NULL,
    tax_region VARCHAR(255) NULL,
    tax_name VARCHAR(255) NULL,
    type VARCHAR(255) NULL,
    currency_code VARCHAR(255) NULL,
    module_config_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    -- CONSTRAINT fk_blc_tax_detail_module_config_id FOREIGN KEY (module_config_id) REFERENCES blc_module_configuration(module_config_id),
    -- CONSTRAINT fk_blc_tax_detail_currency_code FOREIGN KEY (currency_code) REFERENCES blc_currency(currency_code)
);
