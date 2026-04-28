-- +goose Up

CREATE TABLE deductees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    vendor_id UUID NOT NULL,
    name TEXT NOT NULL,
    pan TEXT NOT NULL DEFAULT '',
    pan_verified BOOLEAN NOT NULL DEFAULT FALSE,
    pan_status TEXT NOT NULL DEFAULT 'NOT_VERIFIED',
    deductee_type TEXT NOT NULL,
    resident_status TEXT NOT NULL DEFAULT 'RESIDENT',
    section_override TEXT,
    lower_deduction_cert TEXT,
    lower_deduction_rate DECIMAL(5,2),
    lower_deduction_valid_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, vendor_id)
);

CREATE TABLE tds_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    deductee_id UUID NOT NULL REFERENCES deductees(id),
    section TEXT NOT NULL,
    financial_year TEXT NOT NULL,
    quarter TEXT NOT NULL,
    transaction_date DATE NOT NULL,
    payment_date DATE,
    gross_amount DECIMAL(15,2) NOT NULL,
    tds_rate DECIMAL(8,4) NOT NULL,
    tds_amount DECIMAL(15,2) NOT NULL,
    surcharge DECIMAL(15,2) NOT NULL DEFAULT 0,
    cess DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_tax DECIMAL(15,2) NOT NULL,
    invoice_number TEXT,
    invoice_id UUID,
    nature_of_payment TEXT NOT NULL,
    pan_at_deduction TEXT NOT NULL DEFAULT '',
    no_pan_deduction BOOLEAN NOT NULL DEFAULT FALSE,
    lower_cert_applied BOOLEAN NOT NULL DEFAULT FALSE,
    challan_number TEXT,
    challan_date DATE,
    bsr_code TEXT,
    status TEXT NOT NULL DEFAULT 'PENDING',
    remarks TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tds_aggregates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    deductee_id UUID NOT NULL REFERENCES deductees(id),
    section TEXT NOT NULL,
    financial_year TEXT NOT NULL,
    total_paid DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_tds DECIMAL(15,2) NOT NULL DEFAULT 0,
    transaction_count INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, deductee_id, section, financial_year)
);

-- RLS
ALTER TABLE deductees ENABLE ROW LEVEL SECURITY;
ALTER TABLE tds_entries ENABLE ROW LEVEL SECURITY;
ALTER TABLE tds_aggregates ENABLE ROW LEVEL SECURITY;

CREATE POLICY deductees_tenant_isolation ON deductees
    USING (tenant_id = current_setting('app.tenant_id')::UUID);
CREATE POLICY tds_entries_tenant_isolation ON tds_entries
    USING (tenant_id = current_setting('app.tenant_id')::UUID);
CREATE POLICY tds_aggregates_tenant_isolation ON tds_aggregates
    USING (tenant_id = current_setting('app.tenant_id')::UUID);

-- Indexes
CREATE INDEX idx_deductees_tenant ON deductees(tenant_id);
CREATE INDEX idx_deductees_pan ON deductees(tenant_id, pan);
CREATE INDEX idx_tds_entries_tenant ON tds_entries(tenant_id);
CREATE INDEX idx_tds_entries_deductee ON tds_entries(tenant_id, deductee_id);
CREATE INDEX idx_tds_entries_quarter ON tds_entries(tenant_id, financial_year, quarter);
CREATE INDEX idx_tds_entries_section ON tds_entries(tenant_id, section);
CREATE INDEX idx_tds_aggregates_lookup ON tds_aggregates(tenant_id, deductee_id, section, financial_year);

-- +goose Down
DROP TABLE IF EXISTS tds_aggregates;
DROP TABLE IF EXISTS tds_entries;
DROP TABLE IF EXISTS deductees;
