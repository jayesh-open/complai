-- +goose Up

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    external_id VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, email)
);

CREATE TABLE user_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL DEFAULT 'keycloak',
    provider_user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(provider, provider_user_id)
);

CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_info JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ
);

CREATE TABLE mfa_factors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    factor_type VARCHAR(50) NOT NULL,
    secret_encrypted TEXT,
    phone_number VARCHAR(20),
    verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE step_up_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES user_sessions(id) ON DELETE CASCADE,
    action_class VARCHAR(100) NOT NULL,
    verified_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    mfa_method VARCHAR(50) NOT NULL
);

-- RLS
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_credentials ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE mfa_factors ENABLE ROW LEVEL SECURITY;
ALTER TABLE step_up_events ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON users USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON user_credentials USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON user_sessions USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON mfa_factors USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON step_up_events USING (tenant_id = current_setting('app.tenant_id')::uuid);

CREATE INDEX idx_users_tenant_email ON users(tenant_id, email);
CREATE INDEX idx_users_external_id ON users(external_id);
CREATE INDEX idx_user_sessions_tenant_user ON user_sessions(tenant_id, user_id);
CREATE INDEX idx_step_up_events_tenant_session ON step_up_events(tenant_id, session_id);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS step_up_events CASCADE;
DROP TABLE IF EXISTS mfa_factors CASCADE;
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS user_credentials CASCADE;
DROP TABLE IF EXISTS users CASCADE;
