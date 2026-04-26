-- +goose Up

CREATE TABLE gstr3b_filings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    gstin VARCHAR(15) NOT NULL,
    return_period VARCHAR(6) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    data_json JSONB,
    arn VARCHAR(64),
    filed_at TIMESTAMPTZ,
    filed_by UUID,
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    created_by UUID,
    request_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, gstin, return_period)
);
CREATE INDEX idx_gstr3b_filings_tenant ON gstr3b_filings(tenant_id, gstin, return_period);

ALTER TABLE gstr3b_filings ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON gstr3b_filings USING (tenant_id = current_setting('app.tenant_id')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON gstr3b_filings TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS gstr3b_filings CASCADE;
