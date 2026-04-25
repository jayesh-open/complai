-- +goose Up

CREATE TYPE tenancy_tier AS ENUM ('pooled', 'bridge', 'silo', 'on_prem');
CREATE TYPE tenant_status AS ENUM ('active', 'suspended', 'deleted');

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(500) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    tier tenancy_tier NOT NULL DEFAULT 'pooled',
    status tenant_status NOT NULL DEFAULT 'active',
    kms_key_arn VARCHAR(500),
    settings JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE tenant_pans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    pan VARCHAR(10) NOT NULL,
    entity_name VARCHAR(500) NOT NULL,
    pan_type VARCHAR(20) NOT NULL DEFAULT 'company',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, pan)
);

CREATE TABLE tenant_gstins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    pan_id UUID NOT NULL REFERENCES tenant_pans(id) ON DELETE CASCADE,
    gstin VARCHAR(15) NOT NULL,
    trade_name VARCHAR(500),
    state_code VARCHAR(2) NOT NULL,
    registration_type VARCHAR(50) NOT NULL DEFAULT 'regular',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, gstin)
);

CREATE TABLE tenant_tans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    pan_id UUID NOT NULL REFERENCES tenant_pans(id) ON DELETE CASCADE,
    tan VARCHAR(10) NOT NULL,
    deductor_name VARCHAR(500) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, tan)
);

CREATE TABLE tenant_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    setting_key VARCHAR(255) NOT NULL,
    setting_value JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, setting_key)
);

CREATE TABLE tenant_feature_flags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    flag_name VARCHAR(255) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT false,
    config JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, flag_name)
);

-- RLS
ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_pans ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_gstins ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_tans ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_settings ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_feature_flags ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON tenants USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON tenant_pans USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON tenant_gstins USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON tenant_tans USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON tenant_settings USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON tenant_feature_flags USING (tenant_id = current_setting('app.tenant_id')::uuid);

CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_status ON tenants(status) WHERE status != 'deleted';
CREATE INDEX idx_tenant_pans_tenant ON tenant_pans(tenant_id);
CREATE INDEX idx_tenant_gstins_tenant_pan ON tenant_gstins(tenant_id, pan_id);
CREATE INDEX idx_tenant_tans_tenant_pan ON tenant_tans(tenant_id, pan_id);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS tenant_feature_flags CASCADE;
DROP TABLE IF EXISTS tenant_settings CASCADE;
DROP TABLE IF EXISTS tenant_tans CASCADE;
DROP TABLE IF EXISTS tenant_gstins CASCADE;
DROP TABLE IF EXISTS tenant_pans CASCADE;
DROP TABLE IF EXISTS tenants CASCADE;
DROP TYPE IF EXISTS tenant_status;
DROP TYPE IF EXISTS tenancy_tier;
