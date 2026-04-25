-- +goose Up

CREATE TABLE workflow_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    workflow_type VARCHAR(100) NOT NULL,
    description TEXT,
    version INT NOT NULL DEFAULT 1,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    config JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, workflow_type, version)
);

CREATE TABLE workflow_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    workflow_type VARCHAR(100) NOT NULL,
    temporal_workflow_id VARCHAR(255),
    temporal_run_id VARCHAR(255),
    state VARCHAR(20) NOT NULL DEFAULT 'running',
    input JSONB DEFAULT '{}',
    output JSONB,
    error_message TEXT,
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    trace_id VARCHAR(64)
);
CREATE INDEX idx_workflow_instances_tenant ON workflow_instances(tenant_id, state);
CREATE INDEX idx_workflow_instances_temporal ON workflow_instances(temporal_workflow_id);

CREATE TABLE human_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    workflow_instance_id UUID NOT NULL REFERENCES workflow_instances(id),
    task_type VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    assigned_to UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    input JSONB DEFAULT '{}',
    output JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ
);
CREATE INDEX idx_human_tasks_tenant ON human_tasks(tenant_id, status);
CREATE INDEX idx_human_tasks_workflow ON human_tasks(workflow_instance_id);

ALTER TABLE workflow_definitions ENABLE ROW LEVEL SECURITY;
ALTER TABLE workflow_instances ENABLE ROW LEVEL SECURITY;
ALTER TABLE human_tasks ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON workflow_definitions USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON workflow_instances USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON human_tasks USING (tenant_id = current_setting('app.tenant_id')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS human_tasks CASCADE;
DROP TABLE IF EXISTS workflow_instances CASCADE;
DROP TABLE IF EXISTS workflow_definitions CASCADE;
