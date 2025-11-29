CREATE TABLE IF NOT EXISTS blc_category_attribute (
    category_attribute_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NULL,
    category_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blc_category_attribute_category_id FOREIGN KEY (category_id) REFERENCES blc_category(category_id)
);

CREATE INDEX IF NOT EXISTS idx_blc_category_attribute_category_id ON blc_category_attribute (category_id);
CREATE INDEX IF NOT EXISTS idx_blc_category_attribute_name ON blc_category_attribute (name);
