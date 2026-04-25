-- +goose Up

CREATE TABLE gstr1_filings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    gstin VARCHAR(15) NOT NULL,
    return_period VARCHAR(6) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    total_count INT NOT NULL DEFAULT 0,
    error_count INT NOT NULL DEFAULT 0,
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
CREATE INDEX idx_filings_tenant_gstin ON gstr1_filings(tenant_id, gstin, return_period);

CREATE TABLE sales_register (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    gstin VARCHAR(15) NOT NULL,
    return_period VARCHAR(6) NOT NULL,
    document_number VARCHAR(64) NOT NULL,
    document_date VARCHAR(10) NOT NULL,
    document_type VARCHAR(10) NOT NULL,
    supply_type VARCHAR(10) NOT NULL,
    reverse_charge BOOLEAN NOT NULL DEFAULT false,
    supplier_gstin VARCHAR(15) NOT NULL,
    buyer_gstin VARCHAR(15),
    buyer_name VARCHAR(255),
    buyer_state VARCHAR(2),
    place_of_supply VARCHAR(2) NOT NULL,
    hsn VARCHAR(8),
    taxable_value NUMERIC(15,2) NOT NULL DEFAULT 0,
    cgst_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    cgst_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    sgst_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    sgst_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    igst_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    igst_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    grand_total NUMERIC(15,2) NOT NULL DEFAULT 0,
    source_system VARCHAR(30) NOT NULL DEFAULT 'manual',
    source_id VARCHAR(64),
    section VARCHAR(10),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, gstin, document_number)
);
CREATE INDEX idx_sr_tenant_period ON sales_register(tenant_id, gstin, return_period);
CREATE INDEX idx_sr_section ON sales_register(tenant_id, gstin, return_period, section);

CREATE TABLE gstr1_sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    filing_id UUID NOT NULL REFERENCES gstr1_filings(id),
    section VARCHAR(10) NOT NULL,
    invoice_count INT NOT NULL DEFAULT 0,
    taxable_value NUMERIC(15,2) NOT NULL DEFAULT 0,
    cgst NUMERIC(15,2) NOT NULL DEFAULT 0,
    sgst NUMERIC(15,2) NOT NULL DEFAULT 0,
    igst NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_tax NUMERIC(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, filing_id, section)
);
CREATE INDEX idx_sections_filing ON gstr1_sections(filing_id);

CREATE TABLE validation_errors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    filing_id UUID NOT NULL REFERENCES gstr1_filings(id),
    entry_id UUID NOT NULL,
    field VARCHAR(64) NOT NULL,
    code VARCHAR(32) NOT NULL,
    message TEXT NOT NULL,
    severity VARCHAR(10) NOT NULL DEFAULT 'error',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_valerr_filing ON validation_errors(filing_id);
CREATE INDEX idx_valerr_entry ON validation_errors(entry_id);

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
ALTER TABLE gstr1_filings ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales_register ENABLE ROW LEVEL SECURITY;
ALTER TABLE gstr1_sections ENABLE ROW LEVEL SECURITY;
ALTER TABLE validation_errors ENABLE ROW LEVEL SECURITY;
ALTER TABLE outbox ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON gstr1_filings USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON sales_register USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON gstr1_sections USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON validation_errors USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON outbox USING (tenant_id = current_setting('app.tenant_id')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS validation_errors CASCADE;
DROP TABLE IF EXISTS gstr1_sections CASCADE;
DROP TABLE IF EXISTS sales_register CASCADE;
DROP TABLE IF EXISTS outbox CASCADE;
DROP TABLE IF EXISTS gstr1_filings CASCADE;
