-- +goose Up

CREATE TABLE tds_filings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    form_type TEXT NOT NULL,
    financial_year TEXT NOT NULL,
    quarter TEXT NOT NULL,
    tan TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'DRAFT',
    deductee_count INTEGER NOT NULL DEFAULT 0,
    total_tds_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    fvu_content TEXT,
    token_number TEXT,
    acknowledgement_number TEXT,
    filing_date TIMESTAMPTZ,
    error_message TEXT,
    idempotency_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(idempotency_key)
);

ALTER TABLE tds_filings ENABLE ROW LEVEL SECURITY;

CREATE POLICY tds_filings_tenant_isolation ON tds_filings
    USING (tenant_id = current_setting('app.tenant_id')::UUID);

CREATE INDEX idx_tds_filings_tenant ON tds_filings(tenant_id);
CREATE INDEX idx_tds_filings_lookup ON tds_filings(tenant_id, form_type, financial_year, quarter);
CREATE INDEX idx_tds_filings_idempotency ON tds_filings(idempotency_key);

-- +goose Down
DROP TABLE IF EXISTS tds_filings;
