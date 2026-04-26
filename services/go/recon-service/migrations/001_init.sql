-- +goose Up

CREATE TABLE recon_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    gstin VARCHAR(15) NOT NULL,
    return_period VARCHAR(6) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'RUNNING',
    pr_count INT NOT NULL DEFAULT 0,
    gstr2b_count INT NOT NULL DEFAULT 0,
    matched INT NOT NULL DEFAULT 0,
    mismatch INT NOT NULL DEFAULT 0,
    partial INT NOT NULL DEFAULT 0,
    missing_2b INT NOT NULL DEFAULT 0,
    missing_pr INT NOT NULL DEFAULT 0,
    duplicate INT NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    request_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_recon_runs_tenant ON recon_runs(tenant_id, gstin, return_period);

CREATE TABLE recon_matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    run_id UUID NOT NULL REFERENCES recon_runs(id),
    gstin VARCHAR(15) NOT NULL,
    return_period VARCHAR(6) NOT NULL,
    pr_invoice_number VARCHAR(64),
    pr_invoice_date VARCHAR(10),
    pr_vendor_gstin VARCHAR(15),
    pr_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    pr_hsn VARCHAR(8),
    pr_source_id VARCHAR(64),
    gstr2b_invoice_number VARCHAR(64),
    gstr2b_invoice_date VARCHAR(10),
    gstr2b_supplier_gstin VARCHAR(15),
    gstr2b_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    gstr2b_hsn VARCHAR(8),
    match_type VARCHAR(20) NOT NULL,
    match_confidence NUMERIC(3,2) NOT NULL DEFAULT 0,
    reason_codes TEXT[] DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'UNREVIEWED',
    accepted_by UUID,
    accepted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_recon_matches_run ON recon_matches(run_id);
CREATE INDEX idx_recon_matches_tenant ON recon_matches(tenant_id, gstin, return_period);
CREATE INDEX idx_recon_matches_type ON recon_matches(match_type);
CREATE INDEX idx_recon_matches_status ON recon_matches(status);

CREATE TABLE ims_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    gstin VARCHAR(15) NOT NULL,
    return_period VARCHAR(6) NOT NULL,
    invoice_id VARCHAR(100) NOT NULL,
    action VARCHAR(20) NOT NULL,
    reason TEXT,
    synced_to_gstn BOOLEAN DEFAULT FALSE,
    synced_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_ims_actions_tenant ON ims_actions(tenant_id, gstin, return_period);

CREATE TABLE outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id VARCHAR(64) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    target_queue VARCHAR(100) NOT NULL,
    request_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ
);
CREATE INDEX idx_outbox_status ON outbox(status, created_at);

-- RLS
ALTER TABLE recon_runs ENABLE ROW LEVEL SECURITY;
ALTER TABLE recon_matches ENABLE ROW LEVEL SECURITY;
ALTER TABLE ims_actions ENABLE ROW LEVEL SECURITY;
ALTER TABLE outbox ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON recon_runs USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON recon_matches USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON ims_actions USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON outbox USING (tenant_id = current_setting('app.tenant_id')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS ims_actions CASCADE;
DROP TABLE IF EXISTS recon_matches CASCADE;
DROP TABLE IF EXISTS recon_runs CASCADE;
DROP TABLE IF EXISTS outbox CASCADE;
