-- Create CMS content table
CREATE TABLE IF NOT EXISTS blc_content (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT',
    body TEXT,
    meta_title VARCHAR(255),
    meta_description TEXT,
    meta_keywords TEXT,
    template VARCHAR(100),
    author_id BIGINT NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE,
    version INT NOT NULL DEFAULT 1,
    parent_id BIGINT REFERENCES blc_content(id) ON DELETE SET NULL,
    sort_order INT NOT NULL DEFAULT 0,
    locale VARCHAR(10) NOT NULL DEFAULT 'en_US',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_content_slug ON blc_content(slug);
CREATE INDEX idx_content_type ON blc_content(type, status);
CREATE INDEX idx_content_status ON blc_content(status);
CREATE INDEX idx_content_author ON blc_content(author_id);
CREATE INDEX idx_content_parent ON blc_content(parent_id);
CREATE INDEX idx_content_locale ON blc_content(locale);

COMMENT ON TABLE blc_content IS 'CMS content items (pages, articles, blocks, widgets)';
COMMENT ON COLUMN blc_content.type IS 'Content type: PAGE, ARTICLE, BANNER, BLOCK, WIDGET';
COMMENT ON COLUMN blc_content.status IS 'Publication status: DRAFT, REVIEW, PUBLISHED, ARCHIVED';
