-- Create order adjustment table (order-level discounts)
CREATE TABLE IF NOT EXISTS blc_order_adjustment (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES blc_order(id) ON DELETE CASCADE,
    offer_id BIGINT REFERENCES blc_offer(offer_id) ON DELETE SET NULL,
    offer_name VARCHAR(255) NOT NULL,
    adjustment_value NUMERIC(19, 5) NOT NULL,
    adjustment_reason VARCHAR(255) NOT NULL,
    applied_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_order_adjustment_order_id ON blc_order_adjustment(order_id);
CREATE INDEX IF NOT EXISTS idx_order_adjustment_offer_id ON blc_order_adjustment(offer_id);

-- Create order item adjustment table (item-level discounts)
CREATE TABLE IF NOT EXISTS blc_order_item_adjustment (
    id BIGSERIAL PRIMARY KEY,
    order_item_id BIGINT NOT NULL REFERENCES blc_order_item(id) ON DELETE CASCADE,
    offer_id BIGINT REFERENCES blc_offer(offer_id) ON DELETE SET NULL,
    offer_name VARCHAR(255) NOT NULL,
    adjustment_value NUMERIC(19, 5) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    applied_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_order_item_adjustment_order_item_id ON blc_order_item_adjustment(order_item_id);
CREATE INDEX IF NOT EXISTS idx_order_item_adjustment_offer_id ON blc_order_item_adjustment(offer_id);

-- Create fulfillment group adjustment table (shipping discounts)
CREATE TABLE IF NOT EXISTS blc_fulfillment_group_adjustment (
    id BIGSERIAL PRIMARY KEY,
    fulfillment_group_id BIGINT NOT NULL REFERENCES blc_fulfillment_group(id) ON DELETE CASCADE,
    offer_id BIGINT REFERENCES blc_offer(offer_id) ON DELETE SET NULL,
    offer_name VARCHAR(255) NOT NULL,
    adjustment_value NUMERIC(19, 5) NOT NULL,
    adjustment_reason VARCHAR(255) NOT NULL,
    applied_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fulfillment_adjustment_fulfillment_group_id ON blc_fulfillment_group_adjustment(fulfillment_group_id);
CREATE INDEX IF NOT EXISTS idx_fulfillment_adjustment_offer_id ON blc_fulfillment_group_adjustment(offer_id);

-- Add comments
COMMENT ON TABLE blc_order_adjustment IS 'Order-level adjustments/discounts from offers';
COMMENT ON TABLE blc_order_item_adjustment IS 'Order item-level adjustments/discounts from offers';
COMMENT ON TABLE blc_fulfillment_group_adjustment IS 'Fulfillment group (shipping) adjustments/discounts from offers';
