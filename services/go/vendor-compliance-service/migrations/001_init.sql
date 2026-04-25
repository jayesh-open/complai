-- +goose Up

CREATE TABLE vendor_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    vendor_id VARCHAR(64) NOT NULL,
    name VARCHAR(256) NOT NULL,
    legal_name VARCHAR(256) NOT NULL,
    trade_name VARCHAR(256),
    pan VARCHAR(10) NOT NULL,
    gstin VARCHAR(15) NOT NULL,
    tan VARCHAR(10),
    state VARCHAR(64) NOT NULL,
    state_code VARCHAR(2) NOT NULL,
    category VARCHAR(64) NOT NULL,
    registration_status VARCHAR(32) NOT NULL,
    msme_registered BOOLEAN NOT NULL DEFAULT FALSE,
    email VARCHAR(256),
    phone VARCHAR(20),
    address TEXT,
    synced_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, vendor_id)
);

CREATE TABLE compliance_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    vendor_id VARCHAR(64) NOT NULL,
    vendor_snapshot_id UUID NOT NULL REFERENCES vendor_snapshots(id),
    total_score INTEGER NOT NULL CHECK (total_score >= 0 AND total_score <= 100),
    category VARCHAR(1) NOT NULL CHECK (category IN ('A', 'B', 'C', 'D')),
    risk_level VARCHAR(16) NOT NULL CHECK (risk_level IN ('Low', 'Medium', 'High', 'Critical')),
    filing_regularity_score INTEGER NOT NULL CHECK (filing_regularity_score >= 0 AND filing_regularity_score <= 30),
    irn_compliance_score INTEGER NOT NULL CHECK (irn_compliance_score >= 0 AND irn_compliance_score <= 20),
    mismatch_rate_score INTEGER NOT NULL CHECK (mismatch_rate_score >= 0 AND mismatch_rate_score <= 20),
    payment_behavior_score INTEGER NOT NULL CHECK (payment_behavior_score >= 0 AND payment_behavior_score <= 15),
    document_hygiene_score INTEGER NOT NULL CHECK (document_hygiene_score >= 0 AND document_hygiene_score <= 15),
    filing_regularity_note TEXT,
    irn_compliance_note TEXT,
    mismatch_rate_note TEXT,
    payment_behavior_note TEXT,
    document_hygiene_note TEXT,
    scored_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, vendor_id, scored_at)
);

CREATE TABLE sync_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    sync_type VARCHAR(32) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    vendor_count INTEGER NOT NULL DEFAULT 0,
    scored_count INTEGER NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, sync_type, started_at)
);

-- Indexes
CREATE INDEX idx_vendor_snapshots_tenant ON vendor_snapshots(tenant_id);
CREATE INDEX idx_vendor_snapshots_gstin ON vendor_snapshots(tenant_id, gstin);
CREATE INDEX idx_compliance_scores_tenant ON compliance_scores(tenant_id);
CREATE INDEX idx_compliance_scores_vendor ON compliance_scores(tenant_id, vendor_id);
CREATE INDEX idx_compliance_scores_category ON compliance_scores(tenant_id, category);
CREATE INDEX idx_sync_status_tenant ON sync_status(tenant_id);

-- RLS
ALTER TABLE vendor_snapshots ENABLE ROW LEVEL SECURITY;
ALTER TABLE compliance_scores ENABLE ROW LEVEL SECURITY;
ALTER TABLE sync_status ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON vendor_snapshots USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON compliance_scores USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON sync_status USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- Grants for complai_app role
GRANT SELECT, INSERT, UPDATE, DELETE ON vendor_snapshots TO complai_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON compliance_scores TO complai_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON sync_status TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS sync_status;
DROP TABLE IF EXISTS compliance_scores;
DROP TABLE IF EXISTS vendor_snapshots;
