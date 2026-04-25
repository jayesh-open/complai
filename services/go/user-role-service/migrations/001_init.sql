-- +goose Up

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(500) NOT NULL,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, name)
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    resource VARCHAR(255) NOT NULL,
    action VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, resource, action)
);

CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, role_id, permission_id)
);

CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, user_id, role_id)
);

CREATE TABLE role_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(500) NOT NULL,
    description TEXT,
    permissions JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE approval_workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    resource_type VARCHAR(255) NOT NULL,
    action_type VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_approval',
    requested_by UUID NOT NULL,
    decided_by UUID,
    payload JSONB NOT NULL DEFAULT '{}',
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    decided_at TIMESTAMPTZ
);

-- RLS on tenant-scoped tables
ALTER TABLE roles ENABLE ROW LEVEL SECURITY;
ALTER TABLE permissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE role_permissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;
ALTER TABLE approval_workflows ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON roles USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON permissions USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON role_permissions USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON user_roles USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON approval_workflows USING (tenant_id = current_setting('app.tenant_id')::uuid);

CREATE INDEX idx_roles_tenant ON roles(tenant_id);
CREATE INDEX idx_permissions_tenant_resource ON permissions(tenant_id, resource);
CREATE INDEX idx_user_roles_tenant_user ON user_roles(tenant_id, user_id);
CREATE INDEX idx_approval_workflows_tenant_status ON approval_workflows(tenant_id, status);
CREATE INDEX idx_approval_workflows_requested_by ON approval_workflows(tenant_id, requested_by);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS approval_workflows CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS role_templates CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
