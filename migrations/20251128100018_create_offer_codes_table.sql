CREATE TABLE IF NOT EXISTS blc_offer_code (
    offer_code_id BIGSERIAL PRIMARY KEY,
    archived BPCHAR(1) NULL,
    email_address VARCHAR(255) NULL,
    max_uses INT NULL,
    offer_code VARCHAR(255) NOT NULL,
    end_date TIMESTAMP NULL,
    start_date TIMESTAMP NULL,
    uses INT NULL,
    offer_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blc_offer_code_offer_id FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_offer_code_email_address ON blc_offer_code (email_address);
CREATE INDEX IF NOT EXISTS idx_blc_offer_code_code ON blc_offer_code (offer_code);
CREATE INDEX IF NOT EXISTS idx_blc_offer_code_offer_id ON blc_offer_code (offer_id);
