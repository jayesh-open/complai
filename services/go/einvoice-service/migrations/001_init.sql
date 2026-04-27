-- +goose Up

CREATE TABLE einvoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    irn VARCHAR(64) DEFAULT '',
    ack_no VARCHAR(64) DEFAULT '',
    invoice_number VARCHAR(50) NOT NULL,
    invoice_date VARCHAR(10) NOT NULL,
    invoice_type VARCHAR(3) NOT NULL DEFAULT 'INV',
    supplier_gstin VARCHAR(15) NOT NULL,
    supplier_name VARCHAR(255) DEFAULT '',
    buyer_gstin VARCHAR(15) DEFAULT '',
    buyer_name VARCHAR(255) DEFAULT '',
    supply_type VARCHAR(10) NOT NULL DEFAULT 'B2B',
    place_of_supply VARCHAR(2) DEFAULT '',
    reverse_charge BOOLEAN NOT NULL DEFAULT false,
    taxable_value NUMERIC(15,2) NOT NULL DEFAULT 0,
    cgst_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    sgst_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    igst_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    cess_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    irn_generated_at TIMESTAMPTZ,
    irn_cancelled_at TIMESTAMPTZ,
    cancel_reason VARCHAR(255) DEFAULT '',
    signed_invoice TEXT DEFAULT '',
    signed_qr_code TEXT DEFAULT '',
    request_id UUID NOT NULL,
    source_system VARCHAR(30) NOT NULL DEFAULT 'manual',
    source_id VARCHAR(64) DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, supplier_gstin, invoice_number)
);
CREATE INDEX idx_einv_tenant_gstin ON einvoices(tenant_id, supplier_gstin);
CREATE INDEX idx_einv_irn ON einvoices(irn) WHERE irn != '';
CREATE INDEX idx_einv_status ON einvoices(tenant_id, status);

CREATE TABLE einvoice_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES einvoices(id),
    tenant_id UUID NOT NULL,
    line_number INT NOT NULL DEFAULT 1,
    description VARCHAR(500) DEFAULT '',
    hsn_code VARCHAR(8) DEFAULT '',
    quantity NUMERIC(10,3) NOT NULL DEFAULT 0,
    unit VARCHAR(3) DEFAULT 'NOS',
    unit_price NUMERIC(15,4) NOT NULL DEFAULT 0,
    discount NUMERIC(15,2) NOT NULL DEFAULT 0,
    taxable_value NUMERIC(15,2) NOT NULL DEFAULT 0,
    cgst_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    cgst_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    sgst_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    sgst_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    igst_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    igst_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    cess_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    cess_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_einv_li_invoice ON einvoice_line_items(invoice_id, line_number);

CREATE TABLE einvoice_outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    request_id UUID NOT NULL UNIQUE,
    invoice_id UUID NOT NULL REFERENCES einvoices(id),
    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}',
    target_queue VARCHAR(100) NOT NULL DEFAULT 'gov.outbound.irp.queue',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    published_at TIMESTAMPTZ,
    failed_reason VARCHAR(500) DEFAULT '',
    retry_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_einv_outbox_status ON einvoice_outbox(status, created_at);

-- RLS
ALTER TABLE einvoices ENABLE ROW LEVEL SECURITY;
ALTER TABLE einvoice_line_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE einvoice_outbox ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON einvoices USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON einvoice_line_items USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON einvoice_outbox USING (tenant_id = current_setting('app.tenant_id')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS einvoice_outbox CASCADE;
DROP TABLE IF EXISTS einvoice_line_items CASCADE;
DROP TABLE IF EXISTS einvoices CASCADE;
