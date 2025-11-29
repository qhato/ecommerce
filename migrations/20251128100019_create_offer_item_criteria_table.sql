CREATE TABLE IF NOT EXISTS blc_offer_item_criteria (
    offer_item_criteria_id BIGSERIAL PRIMARY KEY,
    order_item_match_rule TEXT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
