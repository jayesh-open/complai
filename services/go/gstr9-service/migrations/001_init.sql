-- +goose Up

-- GSTR-9 annual return filings
CREATE TABLE IF NOT EXISTS gstr9_filings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    gstin VARCHAR(15) NOT NULL,
    financial_year VARCHAR(9) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    aggregate_turnover NUMERIC(18,2) NOT NULL DEFAULT 0,
    is_mandatory BOOLEAN NOT NULL DEFAULT FALSE,
    arn VARCHAR(50),
    filed_at TIMESTAMPTZ,
    filed_by UUID,
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    request_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, gstin, financial_year)
);

ALTER TABLE gstr9_filings ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_gstr9_filings ON gstr9_filings
    USING (tenant_id::text = current_setting('app.tenant_id', true));

CREATE INDEX idx_gstr9_filings_tenant ON gstr9_filings(tenant_id);
CREATE INDEX idx_gstr9_filings_gstin_fy ON gstr9_filings(gstin, financial_year);

-- GSTR-9 table data (19 tables across 6 parts)
CREATE TABLE IF NOT EXISTS gstr9_table_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    filing_id UUID NOT NULL REFERENCES gstr9_filings(id),
    part_number INT NOT NULL,
    table_number VARCHAR(10) NOT NULL,
    description VARCHAR(200) NOT NULL,
    taxable_value NUMERIC(18,2) NOT NULL DEFAULT 0,
    cgst NUMERIC(18,2) NOT NULL DEFAULT 0,
    sgst NUMERIC(18,2) NOT NULL DEFAULT 0,
    igst NUMERIC(18,2) NOT NULL DEFAULT 0,
    cess NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE gstr9_table_data ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_gstr9_table_data ON gstr9_table_data
    USING (tenant_id::text = current_setting('app.tenant_id', true));

CREATE INDEX idx_gstr9_table_data_filing ON gstr9_table_data(filing_id);

-- GSTR-9 audit logs
CREATE TABLE IF NOT EXISTS gstr9_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    filing_id UUID NOT NULL REFERENCES gstr9_filings(id),
    action VARCHAR(50) NOT NULL,
    details TEXT NOT NULL DEFAULT '',
    actor_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE gstr9_audit_logs ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_gstr9_audit_logs ON gstr9_audit_logs
    USING (tenant_id::text = current_setting('app.tenant_id', true));

CREATE INDEX idx_gstr9_audit_logs_filing ON gstr9_audit_logs(filing_id);

-- GSTR-9C reconciliation filings (for 10c-2)
CREATE TABLE IF NOT EXISTS gstr9c_filings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    gstr9_filing_id UUID NOT NULL REFERENCES gstr9_filings(id),
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    audited_turnover NUMERIC(18,2) NOT NULL DEFAULT 0,
    unreconciled_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    is_self_certified BOOLEAN NOT NULL DEFAULT TRUE,
    certified_at TIMESTAMPTZ,
    certified_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, gstr9_filing_id)
);

ALTER TABLE gstr9c_filings ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_gstr9c_filings ON gstr9c_filings
    USING (tenant_id::text = current_setting('app.tenant_id', true));

CREATE INDEX idx_gstr9c_filings_gstr9 ON gstr9c_filings(gstr9_filing_id);

-- GSTR-9C mismatches (for 10c-2)
CREATE TABLE IF NOT EXISTS gstr9c_mismatches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    gstr9c_filing_id UUID NOT NULL REFERENCES gstr9c_filings(id),
    category VARCHAR(30) NOT NULL,
    description VARCHAR(200) NOT NULL,
    books_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    gstr9_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    difference NUMERIC(18,2) NOT NULL DEFAULT 0,
    severity VARCHAR(10) NOT NULL DEFAULT 'INFO',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE gstr9c_mismatches ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_gstr9c_mismatches ON gstr9c_mismatches
    USING (tenant_id::text = current_setting('app.tenant_id', true));

CREATE INDEX idx_gstr9c_mismatches_filing ON gstr9c_mismatches(gstr9c_filing_id);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS gstr9c_mismatches;
DROP TABLE IF EXISTS gstr9c_filings;
DROP TABLE IF EXISTS gstr9_audit_logs;
DROP TABLE IF EXISTS gstr9_table_data;
DROP TABLE IF EXISTS gstr9_filings;
