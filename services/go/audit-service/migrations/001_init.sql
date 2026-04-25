-- +goose Up

CREATE TABLE audit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    action VARCHAR(50) NOT NULL,
    old_value JSONB,
    new_value JSONB,
    status VARCHAR(20) NOT NULL DEFAULT 'success',
    error_message TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    trace_id VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_audit_events_tenant_resource ON audit_events(tenant_id, resource_type, action, created_at);
CREATE INDEX idx_audit_events_tenant_created ON audit_events(tenant_id, created_at);
CREATE INDEX idx_audit_events_resource_id ON audit_events(tenant_id, resource_id);

CREATE TABLE merkle_chains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    hour_bucket TIMESTAMPTZ NOT NULL,
    event_count INT NOT NULL DEFAULT 0,
    hash_payload TEXT NOT NULL,
    previous_hash VARCHAR(64) NOT NULL DEFAULT '',
    computed_hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, hour_bucket)
);
CREATE INDEX idx_merkle_chains_tenant_hour ON merkle_chains(tenant_id, hour_bucket);

ALTER TABLE audit_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE merkle_chains ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON audit_events USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON merkle_chains USING (tenant_id = current_setting('app.tenant_id')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS merkle_chains CASCADE;
DROP TABLE IF EXISTS audit_events CASCADE;
