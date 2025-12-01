CREATE TABLE IF NOT EXISTS blc_order_item_add_attr (
    order_item_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT blc_order_item_add_attr_pkey PRIMARY KEY (order_item_id, name),
    -- CONSTRAINT fk_blc_order_item_add_attr_order_item_id FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_order_item_add_attr_order_item_id ON blc_order_item_add_attr (order_item_id);
