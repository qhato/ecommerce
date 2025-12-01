CREATE TABLE IF NOT EXISTS blc_order_item_adjustment (
    order_item_adjustment_id BIGSERIAL PRIMARY KEY,
    applied_to_sale_price BOOLEAN NULL,
    adjustment_reason VARCHAR(255) NOT NULL,
    adjustment_value NUMERIC(19, 5) NOT NULL,
    offer_id BIGINT NOT NULL,
    order_item_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- No updated_at in Broadleaf schema, but keeping for consistency if needed in domain
    -- updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- CONSTRAINT fk_blc_order_item_adjustment_order_item_id FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id),
    -- CONSTRAINT fk_blc_order_item_adjustment_offer_id FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_order_item_adjustment_order_item_id ON blc_order_item_adjustment (order_item_id);
CREATE INDEX IF NOT EXISTS idx_blc_order_item_adjustment_offer_id ON blc_order_item_adjustment (offer_id);
