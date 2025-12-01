CREATE TABLE IF NOT EXISTS blc_order (
    order_id BIGSERIAL PRIMARY KEY,
    created_by BIGINT NULL,
    date_created TIMESTAMP NULL,
    date_updated TIMESTAMP NULL,
    updated_by BIGINT NULL,
    email_address VARCHAR(255) NULL,
    name VARCHAR(255) NULL,
    order_number VARCHAR(255) NULL,
    is_preview BOOLEAN NULL,
    order_status VARCHAR(255) NULL,
    order_subtotal NUMERIC(19, 5) NULL,
    submit_date TIMESTAMP NULL,
    tax_override BOOLEAN NULL,
    order_total NUMERIC(19, 5) NULL,
    total_shipping NUMERIC(19, 5) NULL,
    total_tax NUMERIC(19, 5) NULL,
    currency_code VARCHAR(255) NULL,
    customer_id BIGINT NOT NULL,
    locale_code VARCHAR(255) NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- CONSTRAINT fk_blc_order_locale_code FOREIGN KEY (locale_code) REFERENCES blc_locale(locale_code),
    -- CONSTRAINT fk_blc_order_customer_id FOREIGN KEY (customer_id) REFERENCES blc_customer(customer_id),
    -- CONSTRAINT fk_blc_order_currency_code FOREIGN KEY (currency_code) REFERENCES blc_currency(currency_code)
);

CREATE INDEX IF NOT EXISTS idx_blc_order_customer_id ON blc_order (customer_id);
CREATE INDEX IF NOT EXISTS idx_blc_order_order_number ON blc_order (order_number);
CREATE INDEX IF NOT EXISTS idx_blc_order_order_status ON blc_order (order_status);
