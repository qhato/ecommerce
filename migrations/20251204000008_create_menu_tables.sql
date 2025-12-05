-- Create menu table
CREATE TABLE IF NOT EXISTS blc_menu (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_menu_type ON blc_menu(type, is_active);

-- Create menu items table
CREATE TABLE IF NOT EXISTS blc_menu_item (
    id BIGSERIAL PRIMARY KEY,
    menu_id BIGINT NOT NULL REFERENCES blc_menu(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES blc_menu_item(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(500),
    target VARCHAR(20) NOT NULL DEFAULT '_self',
    icon VARCHAR(100),
    css_class VARCHAR(100),
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_menu_item_menu ON blc_menu_item(menu_id, is_active, sort_order);
CREATE INDEX idx_menu_item_parent ON blc_menu_item(parent_id);

COMMENT ON TABLE blc_menu IS 'Navigation menus (header, footer, sidebar)';
COMMENT ON TABLE blc_menu_item IS 'Menu items with hierarchical structure';
COMMENT ON COLUMN blc_menu.type IS 'Menu type: HEADER, FOOTER, SIDEBAR, MOBILE';
