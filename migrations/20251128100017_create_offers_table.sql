CREATE TABLE IF NOT EXISTS blc_offer (
    offer_id BIGSERIAL PRIMARY KEY,
    offer_adjustment_type VARCHAR(255) NULL,
    apply_to_child_items BOOLEAN NULL,
    apply_to_sale_price BOOLEAN NULL,
    archived BPCHAR(1) NULL,
    automatically_added BOOLEAN NULL,
    combinable_with_other_offers BOOLEAN NULL,
    offer_description VARCHAR(255) NULL,
    offer_discount_type VARCHAR(255) NULL,
    end_date TIMESTAMP NULL,
    marketing_message VARCHAR(255) NULL,
    max_uses_per_customer BIGINT NULL,
    max_uses INT NULL,
    max_uses_strategy VARCHAR(255) NULL,
    minimum_days_per_usage BIGINT NULL,
    offer_name VARCHAR(255) NOT NULL,
    offer_item_qualifier_rule TEXT NULL,
    offer_item_target_rule TEXT NULL,
    order_min_total NUMERIC(19, 5) NULL,
    offer_priority INT NULL,
    qualifying_item_min_total NUMERIC(19, 5) NULL,
    requires_related_tar_qual BOOLEAN NULL,
    start_date TIMESTAMP NULL,
    target_min_total NUMERIC(19, 5) NULL,
    target_system VARCHAR(255) NULL,
    totalitarian_offer BOOLEAN NULL,
    offer_type VARCHAR(255) NOT NULL,
    use_list_for_discounts BOOLEAN NULL,
    offer_value NUMERIC(19, 5) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_blc_offer_start_date ON blc_offer (start_date);
CREATE INDEX IF NOT EXISTS idx_blc_offer_automatically_added ON blc_offer (automatically_added);
CREATE INDEX IF NOT EXISTS idx_blc_offer_discount_type ON blc_offer (offer_discount_type);
CREATE INDEX IF NOT EXISTS idx_blc_offer_marketing_message ON blc_offer (marketing_message);
CREATE INDEX IF NOT EXISTS idx_blc_offer_name ON blc_offer (offer_name);
CREATE INDEX IF NOT EXISTS idx_blc_offer_type ON blc_offer (offer_type);
