-- Create admin users table
CREATE TABLE IF NOT EXISTS blc_admin_user (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_super BOOLEAN NOT NULL DEFAULT false,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_admin_user_username ON blc_admin_user(username);
CREATE INDEX idx_admin_user_email ON blc_admin_user(email);
CREATE INDEX idx_admin_user_active ON blc_admin_user(is_active);

-- Create roles table
CREATE TABLE IF NOT EXISTS blc_admin_role (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_admin_role_name ON blc_admin_role(name);
CREATE INDEX idx_admin_role_active ON blc_admin_role(is_active);

-- Create permissions table
CREATE TABLE IF NOT EXISTS blc_admin_permission (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_admin_permission_name ON blc_admin_permission(name);
CREATE INDEX idx_admin_permission_resource ON blc_admin_permission(resource);
CREATE INDEX idx_admin_permission_active ON blc_admin_permission(is_active);

-- Create user-role junction table
CREATE TABLE IF NOT EXISTS blc_admin_user_role (
    user_id BIGINT NOT NULL REFERENCES blc_admin_user(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES blc_admin_role(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    assigned_by BIGINT REFERENCES blc_admin_user(id),
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_admin_user_role_user ON blc_admin_user_role(user_id);
CREATE INDEX idx_admin_user_role_role ON blc_admin_user_role(role_id);

-- Create role-permission junction table
CREATE TABLE IF NOT EXISTS blc_admin_role_permission (
    role_id BIGINT NOT NULL REFERENCES blc_admin_role(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES blc_admin_permission(id) ON DELETE CASCADE,
    granted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    granted_by BIGINT REFERENCES blc_admin_user(id),
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX idx_admin_role_permission_role ON blc_admin_role_permission(role_id);
CREATE INDEX idx_admin_role_permission_permission ON blc_admin_role_permission(permission_id);

-- Create audit log table
CREATE TABLE IF NOT EXISTS blc_admin_audit_log (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    username VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100),
    resource_id VARCHAR(255),
    description TEXT,
    severity VARCHAR(50) NOT NULL DEFAULT 'INFO',
    ip_address VARCHAR(45),
    user_agent TEXT,
    details JSONB,
    success BOOLEAN NOT NULL DEFAULT true,
    error_msg TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_log_user ON blc_admin_audit_log(user_id);
CREATE INDEX idx_audit_log_action ON blc_admin_audit_log(action);
CREATE INDEX idx_audit_log_resource ON blc_admin_audit_log(resource, resource_id);
CREATE INDEX idx_audit_log_created ON blc_admin_audit_log(created_at DESC);
CREATE INDEX idx_audit_log_severity ON blc_admin_audit_log(severity) WHERE severity IN ('WARNING', 'ERROR', 'CRITICAL');

-- Add comments
COMMENT ON TABLE blc_admin_user IS 'Administrative users with access to the admin panel';
COMMENT ON TABLE blc_admin_role IS 'Roles for grouping permissions';
COMMENT ON TABLE blc_admin_permission IS 'Fine-grained permissions for resources and actions';
COMMENT ON TABLE blc_admin_user_role IS 'User-Role assignments (many-to-many)';
COMMENT ON TABLE blc_admin_role_permission IS 'Role-Permission assignments (many-to-many)';
COMMENT ON TABLE blc_admin_audit_log IS 'Audit trail of all administrative actions';

-- Insert default super admin user (password: admin123 - should be changed immediately)
-- Password hash generated with bcrypt cost 10
INSERT INTO blc_admin_user (username, email, password_hash, first_name, last_name, is_super)
VALUES ('admin', 'admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Super', 'Admin', true)
ON CONFLICT (username) DO NOTHING;

-- Insert default roles
INSERT INTO blc_admin_role (name, description) VALUES
('Administrator', 'Full system administrator with all permissions'),
('Editor', 'Can create and edit content'),
('Viewer', 'Read-only access to resources')
ON CONFLICT (name) DO NOTHING;

-- Insert default permissions
INSERT INTO blc_admin_permission (name, description, resource, action) VALUES
('PRODUCT_CREATE', 'Create products', 'PRODUCT', 'CREATE'),
('PRODUCT_READ', 'View products', 'PRODUCT', 'READ'),
('PRODUCT_UPDATE', 'Update products', 'PRODUCT', 'UPDATE'),
('PRODUCT_DELETE', 'Delete products', 'PRODUCT', 'DELETE'),
('ORDER_READ', 'View orders', 'ORDER', 'READ'),
('ORDER_UPDATE', 'Update orders', 'ORDER', 'UPDATE'),
('CUSTOMER_READ', 'View customers', 'CUSTOMER', 'READ'),
('CUSTOMER_UPDATE', 'Update customers', 'CUSTOMER', 'UPDATE'),
('ADMIN_ALL', 'All administrative permissions', 'ADMIN', 'ALL')
ON CONFLICT (name) DO NOTHING;
