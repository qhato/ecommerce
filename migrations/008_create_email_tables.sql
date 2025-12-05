-- Create emails table
CREATE TABLE IF NOT EXISTS emails (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    priority INTEGER NOT NULL DEFAULT 5,
    from_address VARCHAR(255) NOT NULL,
    to_addresses TEXT[] NOT NULL,
    cc_addresses TEXT[],
    bcc_addresses TEXT[],
    reply_to VARCHAR(255),
    subject VARCHAR(500) NOT NULL,
    body TEXT,
    html_body TEXT,
    template_name VARCHAR(100),
    template_data JSONB,
    headers JSONB,
    max_retries INTEGER NOT NULL DEFAULT 3,
    retry_count INTEGER NOT NULL DEFAULT 0,
    scheduled_at TIMESTAMP,
    sent_at TIMESTAMP,
    failed_at TIMESTAMP,
    error_message TEXT,
    order_id BIGINT,
    customer_id BIGINT,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create email_attachments table
CREATE TABLE IF NOT EXISTS email_attachments (
    id BIGSERIAL PRIMARY KEY,
    email_id BIGINT NOT NULL REFERENCES emails(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    content BYTEA NOT NULL,
    size BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for emails
CREATE INDEX IF NOT EXISTS idx_emails_status ON emails(status);
CREATE INDEX IF NOT EXISTS idx_emails_type ON emails(type);
CREATE INDEX IF NOT EXISTS idx_emails_order_id ON emails(order_id) WHERE order_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_emails_customer_id ON emails(customer_id) WHERE customer_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_emails_scheduled_at ON emails(scheduled_at) WHERE scheduled_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_emails_created_at ON emails(created_at);
CREATE INDEX IF NOT EXISTS idx_emails_priority_status ON emails(priority DESC, status) WHERE status IN ('PENDING', 'QUEUED');

-- Create indexes for email_attachments
CREATE INDEX IF NOT EXISTS idx_email_attachments_email_id ON email_attachments(email_id);

-- Add comments
COMMENT ON TABLE emails IS 'Stores email messages for transactional email service';
COMMENT ON TABLE email_attachments IS 'Stores email attachments';
COMMENT ON COLUMN emails.type IS 'Type of email (ORDER_CONFIRMATION, PASSWORD_RESET, etc.)';
COMMENT ON COLUMN emails.status IS 'Status of email (PENDING, QUEUED, SENDING, SENT, FAILED, RETRYING, CANCELLED)';
COMMENT ON COLUMN emails.priority IS 'Priority of email (1=LOW, 5=NORMAL, 10=HIGH, 20=URGENT)';
COMMENT ON COLUMN emails.template_data IS 'JSON data for template rendering';
COMMENT ON COLUMN emails.headers IS 'Custom email headers as JSON';
