-- Create search_synonyms table
CREATE TABLE IF NOT EXISTS search_synonyms (
    id BIGSERIAL PRIMARY KEY,
    term VARCHAR(255) NOT NULL,
    synonyms TEXT[] NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_search_synonyms_term ON search_synonyms(LOWER(term));
CREATE INDEX idx_search_synonyms_active ON search_synonyms(is_active) WHERE is_active=true;

-- Create search_redirects table
CREATE TABLE IF NOT EXISTS search_redirects (
    id BIGSERIAL PRIMARY KEY,
    search_term VARCHAR(255) NOT NULL,
    target_url VARCHAR(1000) NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    activation_date TIMESTAMP,
    expiration_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_search_redirects_term ON search_redirects(LOWER(search_term));
CREATE INDEX idx_search_redirects_active ON search_redirects(is_active, priority) WHERE is_active=true;

-- Create search_facet_configs table
CREATE TABLE IF NOT EXISTS search_facet_configs (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    label VARCHAR(200) NOT NULL,
    field_name VARCHAR(100) NOT NULL,
    facet_type VARCHAR(20) NOT NULL CHECK (facet_type IN ('FIELD', 'RANGE')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    show_in_results BOOLEAN NOT NULL DEFAULT true,
    show_in_navigation BOOLEAN NOT NULL DEFAULT false,
    priority INTEGER NOT NULL DEFAULT 0,
    min_doc_count INTEGER NOT NULL DEFAULT 1,
    max_values INTEGER NOT NULL DEFAULT 10,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_search_facet_configs_active ON search_facet_configs(is_active, priority) WHERE is_active=true;

-- Create indexing_jobs table
CREATE TABLE IF NOT EXISTS indexing_jobs (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL CHECK (type IN ('FULL', 'INCREMENTAL', 'SINGLE')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('PENDING', 'RUNNING', 'COMPLETED', 'FAILED', 'CANCELLED')),
    entity_type VARCHAR(50) NOT NULL,
    total_items INTEGER NOT NULL DEFAULT 0,
    processed_items INTEGER NOT NULL DEFAULT 0,
    failed_items INTEGER NOT NULL DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_by BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_indexing_jobs_status ON indexing_jobs(status, started_at DESC);
CREATE INDEX idx_indexing_jobs_created ON indexing_jobs(created_at DESC);

-- Add comments
COMMENT ON TABLE search_synonyms IS 'Search synonyms for query expansion';
COMMENT ON TABLE search_redirects IS 'Search redirects for specific search terms';
COMMENT ON TABLE search_facet_configs IS 'Configuration for search facets/filters';
COMMENT ON TABLE indexing_jobs IS 'Indexing job tracking for Elasticsearch';
