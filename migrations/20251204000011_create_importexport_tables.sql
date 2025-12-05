-- Create import/export job table
CREATE TABLE IF NOT EXISTS blc_import_export_job (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    format VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    file_path VARCHAR(500),
    file_name VARCHAR(255),
    total_records INT NOT NULL DEFAULT 0,
    processed_records INT NOT NULL DEFAULT 0,
    success_records INT NOT NULL DEFAULT 0,
    failed_records INT NOT NULL DEFAULT 0,
    error_log TEXT,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_job_type ON blc_import_export_job(type, status);
CREATE INDEX idx_job_status ON blc_import_export_job(status);
CREATE INDEX idx_job_entity ON blc_import_export_job(entity_type);
CREATE INDEX idx_job_created_by ON blc_import_export_job(created_by);
CREATE INDEX idx_job_created ON blc_import_export_job(created_at DESC);

COMMENT ON TABLE blc_import_export_job IS 'Import/Export jobs for bulk data operations';
COMMENT ON COLUMN blc_import_export_job.type IS 'Job type: IMPORT, EXPORT';
COMMENT ON COLUMN blc_import_export_job.entity_type IS 'Entity type: PRODUCT, CATEGORY, CUSTOMER, ORDER, CONTENT';
COMMENT ON COLUMN blc_import_export_job.format IS 'File format: CSV, JSON, XML';
COMMENT ON COLUMN blc_import_export_job.status IS 'Status: PENDING, PROCESSING, COMPLETED, FAILED, CANCELLED';
