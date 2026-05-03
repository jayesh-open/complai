-- +goose Up

-- Bulk filing batches (employer-initiated batch ITR filing)
CREATE TABLE bulk_filing_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    tax_year VARCHAR(10) NOT NULL,
    employer_tan VARCHAR(10) NOT NULL,
    employer_name VARCHAR(255) NOT NULL,
    total_employees INT NOT NULL DEFAULT 0,
    processed INT NOT NULL DEFAULT 0,
    ready INT NOT NULL DEFAULT 0,
    with_mismatches INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bulk_batches_tenant ON bulk_filing_batches(tenant_id);
CREATE INDEX idx_bulk_batches_tenant_year ON bulk_filing_batches(tenant_id, tax_year);

ALTER TABLE bulk_filing_batches ENABLE ROW LEVEL SECURITY;
CREATE POLICY bulk_batches_tenant_isolation ON bulk_filing_batches
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- Bulk filing employees (per-employee within a batch)
CREATE TABLE bulk_filing_employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    batch_id UUID NOT NULL REFERENCES bulk_filing_batches(id),
    pan VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    gross_salary DECIMAL(15,2) NOT NULL DEFAULT 0,
    tds_deducted DECIMAL(15,2) NOT NULL DEFAULT 0,
    form_type VARCHAR(10) NOT NULL DEFAULT 'ITR-1',
    filing_id UUID REFERENCES itr_filings(id),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING_REVIEW',
    mismatch_count INT NOT NULL DEFAULT 0,
    magic_link_token VARCHAR(64),
    token_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_bulk_emp_batch_pan UNIQUE (batch_id, pan)
);

CREATE INDEX idx_bulk_employees_tenant ON bulk_filing_employees(tenant_id);
CREATE INDEX idx_bulk_employees_batch ON bulk_filing_employees(batch_id);
CREATE INDEX idx_bulk_employees_pan ON bulk_filing_employees(tenant_id, pan);
CREATE INDEX idx_bulk_employees_token ON bulk_filing_employees(magic_link_token) WHERE magic_link_token IS NOT NULL;

ALTER TABLE bulk_filing_employees ENABLE ROW LEVEL SECURITY;
CREATE POLICY bulk_employees_tenant_isolation ON bulk_filing_employees
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- Magic link tokens (standalone table for token lookup without tenant context)
CREATE TABLE magic_link_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    token VARCHAR(64) NOT NULL UNIQUE,
    pan VARCHAR(10) NOT NULL,
    batch_id UUID NOT NULL REFERENCES bulk_filing_batches(id),
    filing_id UUID NOT NULL REFERENCES itr_filings(id),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_magic_tokens_token ON magic_link_tokens(token);
CREATE INDEX idx_magic_tokens_tenant ON magic_link_tokens(tenant_id);

ALTER TABLE magic_link_tokens ENABLE ROW LEVEL SECURITY;
CREATE POLICY magic_tokens_tenant_isolation ON magic_link_tokens
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- Add severity column to ais_reconciliations
ALTER TABLE ais_reconciliations ADD COLUMN severity VARCHAR(10) NOT NULL DEFAULT 'INFO';

-- Add bulk_batch_id FK to itr_filings for batch-originated filings
ALTER TABLE itr_filings ADD COLUMN bulk_batch_id UUID REFERENCES bulk_filing_batches(id);

-- Grant privileges
GRANT SELECT, INSERT, UPDATE, DELETE ON bulk_filing_batches TO complai_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON bulk_filing_employees TO complai_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON magic_link_tokens TO complai_app;

-- +goose Down
ALTER TABLE itr_filings DROP COLUMN IF EXISTS bulk_batch_id;
ALTER TABLE ais_reconciliations DROP COLUMN IF EXISTS severity;
DROP TABLE IF EXISTS magic_link_tokens;
DROP TABLE IF EXISTS bulk_filing_employees;
DROP TABLE IF EXISTS bulk_filing_batches;
