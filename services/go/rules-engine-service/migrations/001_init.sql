-- +goose Up

CREATE TABLE rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    category VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    version INT NOT NULL DEFAULT 1,
    priority INT NOT NULL DEFAULT 100,
    conditions JSONB NOT NULL DEFAULT '{}',
    actions JSONB NOT NULL DEFAULT '{}',
    effective_from DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_to DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, category, name, version)
);
CREATE INDEX idx_rules_tenant_category ON rules(tenant_id, category, status);

CREATE TABLE rule_execution_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    rule_id UUID REFERENCES rules(id),
    input_data JSONB NOT NULL,
    matched_rules JSONB,
    output JSONB NOT NULL,
    execution_time_ms INT NOT NULL DEFAULT 0,
    trace_id VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_rule_exec_tenant ON rule_execution_logs(tenant_id, created_at);

ALTER TABLE rules ENABLE ROW LEVEL SECURITY;
ALTER TABLE rule_execution_logs ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON rules USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON rule_execution_logs USING (tenant_id = current_setting('app.tenant_id')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS rule_execution_logs CASCADE;
DROP TABLE IF EXISTS rules CASCADE;
