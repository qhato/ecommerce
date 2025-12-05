-- Create notification table
CREATE TABLE IF NOT EXISTS blc_notification (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    priority VARCHAR(50) NOT NULL DEFAULT 'NORMAL',
    recipient_id VARCHAR(255) NOT NULL,
    recipient_email VARCHAR(255),
    recipient_phone VARCHAR(50),
    subject VARCHAR(500),
    body TEXT NOT NULL,
    template_id BIGINT,
    template_data JSONB,
    scheduled_for TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,
    error_msg TEXT,
    retry_count INT NOT NULL DEFAULT 0,
    max_retries INT NOT NULL DEFAULT 3,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notification_type ON blc_notification(type, status);
CREATE INDEX idx_notification_status ON blc_notification(status, priority);
CREATE INDEX idx_notification_recipient ON blc_notification(recipient_id);
CREATE INDEX idx_notification_scheduled ON blc_notification(scheduled_for) WHERE status = 'PENDING';
CREATE INDEX idx_notification_failed ON blc_notification(status, retry_count) WHERE status = 'FAILED';

-- Create notification template table
CREATE TABLE IF NOT EXISTS blc_notification_template (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,
    subject VARCHAR(500),
    body_template TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_template_name ON blc_notification_template(name);
CREATE INDEX idx_template_type ON blc_notification_template(type, is_active);

COMMENT ON TABLE blc_notification IS 'Notifications to be sent (email, SMS, push)';
COMMENT ON TABLE blc_notification_template IS 'Notification templates for reusable content';
COMMENT ON COLUMN blc_notification.type IS 'Notification type: EMAIL, SMS, PUSH';
COMMENT ON COLUMN blc_notification.status IS 'Status: PENDING, SENDING, SENT, FAILED, CANCELLED';
COMMENT ON COLUMN blc_notification.priority IS 'Priority: LOW, NORMAL, HIGH, URGENT';
