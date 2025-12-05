-- Create return request table
CREATE TABLE IF NOT EXISTS blc_return_request (
    id BIGSERIAL PRIMARY KEY,
    rma VARCHAR(100) UNIQUE NOT NULL,
    order_id BIGINT NOT NULL,
    customer_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'REQUESTED',
    reason VARCHAR(50) NOT NULL,
    reason_details TEXT,
    refund_amount NUMERIC(19, 5) NOT NULL DEFAULT 0,
    refund_method VARCHAR(50),
    approved_by BIGINT,
    approved_at TIMESTAMP WITH TIME ZONE,
    received_at TIMESTAMP WITH TIME ZONE,
    inspected_at TIMESTAMP WITH TIME ZONE,
    refunded_at TIMESTAMP WITH TIME ZONE,
    tracking_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_return_rma ON blc_return_request(rma);
CREATE INDEX idx_return_order ON blc_return_request(order_id);
CREATE INDEX idx_return_customer ON blc_return_request(customer_id);
CREATE INDEX idx_return_status ON blc_return_request(status);

-- Create return items table
CREATE TABLE IF NOT EXISTS blc_return_item (
    id BIGSERIAL PRIMARY KEY,
    return_id BIGINT NOT NULL REFERENCES blc_return_request(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL,
    sku VARCHAR(100) NOT NULL,
    quantity INT NOT NULL,
    unit_price NUMERIC(19, 5) NOT NULL,
    total_price NUMERIC(19, 5) NOT NULL,
    condition VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_return_item_return ON blc_return_item(return_id);
CREATE INDEX idx_return_item_product ON blc_return_item(product_id);

COMMENT ON TABLE blc_return_request IS 'Product return requests (RMA)';
COMMENT ON TABLE blc_return_item IS 'Items included in return requests';
COMMENT ON COLUMN blc_return_request.status IS 'Status: REQUESTED, APPROVED, REJECTED, RECEIVED, INSPECTED, REFUNDED, CANCELLED';
COMMENT ON COLUMN blc_return_request.reason IS 'Reason: DEFECTIVE, WRONG_ITEM, NOT_AS_DESCRIBED, CHANGED_MIND, OTHER';
