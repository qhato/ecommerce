CREATE TABLE IF NOT EXISTS blc_offer_price_data (
    offer_price_data_id BIGSERIAL PRIMARY KEY,
    end_date TIMESTAMP NULL,
    start_date TIMESTAMP NULL,
    amount NUMERIC(19, 5) NOT NULL,
    archived BPCHAR(1) NULL,
    discount_type VARCHAR(255) NULL,
    identifier_type VARCHAR(255) NULL,
    identifier_value VARCHAR(255) NULL,
    quantity INT NOT NULL,
    offer_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blc_offer_price_data_offer_id FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_offer_price_data_offer_id ON blc_offer_price_data (offer_id);
