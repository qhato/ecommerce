CREATE TABLE IF NOT EXISTS blc_order_item (
    order_item_id BIGSERIAL PRIMARY KEY,
    created_by BIGINT NULL,
    date_created TIMESTAMP NULL,
    date_updated TIMESTAMP NULL,
    updated_by BIGINT NULL,
    discounts_allowed BOOLEAN NULL,
    has_validation_errors BOOLEAN NULL,
    item_taxable_flag BOOLEAN NULL,
    name VARCHAR(255) NULL,
    order_item_type VARCHAR(255) NULL,
    price NUMERIC(19, 5) NULL,
    quantity INT NOT NULL,
    retail_price NUMERIC(19, 5) NULL,
    retail_price_override BOOLEAN NULL,
    sale_price NUMERIC(19, 5) NULL,
    sale_price_override BOOLEAN NULL,
    total_tax NUMERIC(19, 2) NULL,
    category_id BIGINT NULL,
    gift_wrap_item_id BIGINT NULL,
    order_id BIGINT NULL,
    parent_order_item_id BIGINT NULL,
    personal_message_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- CONSTRAINT fk_blc_order_item_category_id FOREIGN KEY (category_id) REFERENCES blc_category(category_id),
    -- CONSTRAINT fk_blc_order_item_order_id FOREIGN KEY (order_id) REFERENCES blc_order(order_id),
    -- CONSTRAINT fk_blc_order_item_personal_message_id FOREIGN KEY (personal_message_id) REFERENCES blc_personal_message(personal_message_id)
    -- Note: gift_wrap_item_id and parent_order_item_id might reference blc_order_item itself, need to define properly once blc_order_item exists
);

CREATE INDEX IF NOT EXISTS idx_blc_order_item_order_id ON blc_order_item (order_id);
CREATE INDEX IF NOT EXISTS idx_blc_order_item_category_id ON blc_order_item (category_id);
CREATE INDEX IF NOT EXISTS idx_blc_order_item_parent_order_item_id ON blc_order_item (parent_order_item_id);
