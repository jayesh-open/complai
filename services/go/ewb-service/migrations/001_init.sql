-- +goose Up

CREATE TABLE ewb (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    ewb_number VARCHAR(20) DEFAULT '',
    ewb_date VARCHAR(10) DEFAULT '',
    doc_type VARCHAR(10) NOT NULL DEFAULT 'INV',
    doc_number VARCHAR(50) NOT NULL,
    doc_date VARCHAR(10) NOT NULL,
    supplier_gstin VARCHAR(15) NOT NULL,
    supplier_name VARCHAR(255) DEFAULT '',
    buyer_gstin VARCHAR(15) DEFAULT '',
    buyer_name VARCHAR(255) DEFAULT '',
    supply_type VARCHAR(10) NOT NULL DEFAULT 'O',
    sub_supply_type VARCHAR(10) DEFAULT '',
    transport_mode VARCHAR(5) NOT NULL DEFAULT '1',
    vehicle_number VARCHAR(20) DEFAULT '',
    vehicle_type VARCHAR(5) DEFAULT 'R',
    transporter_id VARCHAR(15) DEFAULT '',
    transporter_name VARCHAR(255) DEFAULT '',
    from_place VARCHAR(255) DEFAULT '',
    from_state VARCHAR(2) DEFAULT '',
    from_pincode VARCHAR(6) DEFAULT '',
    to_place VARCHAR(255) DEFAULT '',
    to_state VARCHAR(2) DEFAULT '',
    to_pincode VARCHAR(6) DEFAULT '',
    distance_km INT NOT NULL DEFAULT 0,
    taxable_value NUMERIC(15,2) NOT NULL DEFAULT 0,
    cgst_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    sgst_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    igst_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    cess_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_value NUMERIC(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    valid_from TIMESTAMPTZ,
    valid_until TIMESTAMPTZ,
    generated_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    cancel_reason VARCHAR(255) DEFAULT '',
    consolidated_ewb_id UUID,
    request_id UUID NOT NULL,
    source_system VARCHAR(30) NOT NULL DEFAULT 'manual',
    source_id VARCHAR(64) DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, supplier_gstin, doc_number)
);
CREATE INDEX idx_ewb_tenant_gstin ON ewb(tenant_id, supplier_gstin);
CREATE INDEX idx_ewb_number ON ewb(ewb_number) WHERE ewb_number != '';
CREATE INDEX idx_ewb_status ON ewb(tenant_id, status);

CREATE TABLE ewb_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ewb_id UUID NOT NULL REFERENCES ewb(id),
    tenant_id UUID NOT NULL,
    item_number INT NOT NULL DEFAULT 1,
    product_name VARCHAR(500) DEFAULT '',
    product_desc VARCHAR(500) DEFAULT '',
    hsn_code VARCHAR(8) DEFAULT '',
    quantity NUMERIC(10,3) NOT NULL DEFAULT 0,
    unit VARCHAR(3) DEFAULT 'NOS',
    taxable_value NUMERIC(15,2) NOT NULL DEFAULT 0,
    cgst_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    sgst_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    igst_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    cess_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_ewb_items_ewb ON ewb_items(ewb_id, item_number);

CREATE TABLE ewb_vehicle_updates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ewb_id UUID NOT NULL REFERENCES ewb(id),
    tenant_id UUID NOT NULL,
    vehicle_number VARCHAR(20) NOT NULL,
    from_place VARCHAR(255) DEFAULT '',
    from_state VARCHAR(2) DEFAULT '',
    transport_mode VARCHAR(5) NOT NULL DEFAULT '1',
    reason VARCHAR(10) NOT NULL DEFAULT '1',
    remark VARCHAR(255) DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_ewb_vehicle_ewb ON ewb_vehicle_updates(ewb_id, updated_at);

CREATE TABLE ewb_consolidations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    consolidated_ewb_number VARCHAR(20) NOT NULL,
    trip_sheet_number VARCHAR(50) DEFAULT '',
    vehicle_number VARCHAR(20) NOT NULL,
    from_place VARCHAR(255) DEFAULT '',
    from_state VARCHAR(2) DEFAULT '',
    to_place VARCHAR(255) DEFAULT '',
    to_state VARCHAR(2) DEFAULT '',
    transport_mode VARCHAR(5) NOT NULL DEFAULT '1',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    generated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ewb_consolidation_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    consolidation_id UUID NOT NULL REFERENCES ewb_consolidations(id),
    ewb_id UUID NOT NULL REFERENCES ewb(id),
    tenant_id UUID NOT NULL,
    ewb_number VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_ewb_consol_items ON ewb_consolidation_items(consolidation_id);

CREATE TABLE ewb_outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    request_id UUID NOT NULL UNIQUE,
    ewb_id UUID REFERENCES ewb(id),
    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}',
    target_queue VARCHAR(100) NOT NULL DEFAULT 'gov.outbound.ewb.queue',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    published_at TIMESTAMPTZ,
    failed_reason VARCHAR(500) DEFAULT '',
    retry_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_ewb_outbox_status ON ewb_outbox(status, created_at);

-- RLS
ALTER TABLE ewb ENABLE ROW LEVEL SECURITY;
ALTER TABLE ewb_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE ewb_vehicle_updates ENABLE ROW LEVEL SECURITY;
ALTER TABLE ewb_consolidations ENABLE ROW LEVEL SECURITY;
ALTER TABLE ewb_consolidation_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE ewb_outbox ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON ewb USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON ewb_items USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON ewb_vehicle_updates USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON ewb_consolidations USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON ewb_consolidation_items USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON ewb_outbox USING (tenant_id = current_setting('app.tenant_id')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS ewb_consolidation_items CASCADE;
DROP TABLE IF EXISTS ewb_consolidations CASCADE;
DROP TABLE IF EXISTS ewb_vehicle_updates CASCADE;
DROP TABLE IF EXISTS ewb_items CASCADE;
DROP TABLE IF EXISTS ewb_outbox CASCADE;
DROP TABLE IF EXISTS ewb CASCADE;
