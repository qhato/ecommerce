CREATE TABLE IF NOT EXISTS blc_sku_availability (
    sku_availability_id BIGSERIAL PRIMARY KEY,
    availability_date TIMESTAMP NULL,
    availability_status VARCHAR(255) NULL,
    location_id BIGINT NULL,
    qty_on_hand INT NULL,
    reserve_qty INT NULL,
    sku_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- CONSTRAINT fk_blc_sku_availability_sku_id FOREIGN KEY (sku_id) REFERENCES blc_sku(sku_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_sku_availability_location_id ON blc_sku_availability (location_id);
CREATE INDEX IF NOT EXISTS idx_blc_sku_availability_sku_id ON blc_sku_availability (sku_id);
CREATE INDEX IF NOT EXISTS idx_blc_sku_availability_status ON blc_sku_availability (availability_status);
