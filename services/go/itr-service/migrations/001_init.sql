-- +goose Up

-- Taxpayer profiles
CREATE TABLE taxpayers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    pan VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    assessee_type VARCHAR(20) NOT NULL DEFAULT 'INDIVIDUAL',
    residency_status VARCHAR(20) NOT NULL DEFAULT 'RESIDENT',
    aadhaar_linked BOOLEAN NOT NULL DEFAULT false,
    email VARCHAR(255),
    mobile VARCHAR(15),
    address TEXT,
    employer_tan VARCHAR(10),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_taxpayer_tenant_pan UNIQUE (tenant_id, pan)
);

CREATE INDEX idx_taxpayers_tenant ON taxpayers(tenant_id);
CREATE INDEX idx_taxpayers_pan ON taxpayers(tenant_id, pan);

ALTER TABLE taxpayers ENABLE ROW LEVEL SECURITY;
CREATE POLICY taxpayers_tenant_isolation ON taxpayers
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- ITR filings
CREATE TABLE itr_filings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    taxpayer_id UUID NOT NULL REFERENCES taxpayers(id),
    pan VARCHAR(10) NOT NULL,
    tax_year VARCHAR(10) NOT NULL,
    form_type VARCHAR(10) NOT NULL,
    regime_selected VARCHAR(20) NOT NULL DEFAULT 'NEW_REGIME',
    form_10iea_ref VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    gross_income DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_deductions DECIMAL(15,2) NOT NULL DEFAULT 0,
    taxable_income DECIMAL(15,2) NOT NULL DEFAULT 0,
    tax_payable DECIMAL(15,2) NOT NULL DEFAULT 0,
    tds_credited DECIMAL(15,2) NOT NULL DEFAULT 0,
    advance_tax_paid DECIMAL(15,2) NOT NULL DEFAULT 0,
    self_assessment_tax DECIMAL(15,2) NOT NULL DEFAULT 0,
    refund_due DECIMAL(15,2) NOT NULL DEFAULT 0,
    balance_payable DECIMAL(15,2) NOT NULL DEFAULT 0,
    verification_method VARCHAR(20),
    arn VARCHAR(50),
    acknowledgement_number VARCHAR(50),
    filed_at TIMESTAMPTZ,
    idempotency_key VARCHAR(255) NOT NULL,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_filing_idempotency UNIQUE (idempotency_key)
);

CREATE INDEX idx_filings_tenant ON itr_filings(tenant_id);
CREATE INDEX idx_filings_tenant_year ON itr_filings(tenant_id, tax_year);
CREATE INDEX idx_filings_pan_year ON itr_filings(tenant_id, pan, tax_year);

ALTER TABLE itr_filings ENABLE ROW LEVEL SECURITY;
CREATE POLICY filings_tenant_isolation ON itr_filings
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- Income heads
CREATE TABLE income_heads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    filing_id UUID NOT NULL REFERENCES itr_filings(id),
    head VARCHAR(30) NOT NULL,
    sub_head VARCHAR(50),
    section VARCHAR(30),
    description VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    exempt BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_income_heads_filing ON income_heads(filing_id);
CREATE INDEX idx_income_heads_tenant ON income_heads(tenant_id);

ALTER TABLE income_heads ENABLE ROW LEVEL SECURITY;
CREATE POLICY income_heads_tenant_isolation ON income_heads
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- Deductions
CREATE TABLE deductions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    filing_id UUID NOT NULL REFERENCES itr_filings(id),
    section VARCHAR(30) NOT NULL,
    label VARCHAR(255) NOT NULL,
    claimed DECIMAL(15,2) NOT NULL DEFAULT 0,
    allowed DECIMAL(15,2) NOT NULL DEFAULT 0,
    max_limit DECIMAL(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_deductions_filing ON deductions(filing_id);
CREATE INDEX idx_deductions_tenant ON deductions(tenant_id);

ALTER TABLE deductions ENABLE ROW LEVEL SECURITY;
CREATE POLICY deductions_tenant_isolation ON deductions
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- Tax computations
CREATE TABLE tax_computations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    filing_id UUID NOT NULL REFERENCES itr_filings(id),
    regime_type VARCHAR(20) NOT NULL,
    gross_income DECIMAL(15,2) NOT NULL DEFAULT 0,
    standard_deduction DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_deductions DECIMAL(15,2) NOT NULL DEFAULT 0,
    taxable_income DECIMAL(15,2) NOT NULL DEFAULT 0,
    base_tax DECIMAL(15,2) NOT NULL DEFAULT 0,
    surcharge DECIMAL(15,2) NOT NULL DEFAULT 0,
    surcharge_rate DECIMAL(5,4) NOT NULL DEFAULT 0,
    health_ed_cess DECIMAL(15,2) NOT NULL DEFAULT 0,
    rebate_87a DECIMAL(15,2) NOT NULL DEFAULT 0,
    gross_tax_payable DECIMAL(15,2) NOT NULL DEFAULT 0,
    tds_credit DECIMAL(15,2) NOT NULL DEFAULT 0,
    advance_tax DECIMAL(15,2) NOT NULL DEFAULT 0,
    self_assessment_tax DECIMAL(15,2) NOT NULL DEFAULT 0,
    net_tax_payable DECIMAL(15,2) NOT NULL DEFAULT 0,
    refund_due DECIMAL(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_computation_filing UNIQUE (filing_id)
);

CREATE INDEX idx_computations_tenant ON tax_computations(tenant_id);

ALTER TABLE tax_computations ENABLE ROW LEVEL SECURITY;
CREATE POLICY computations_tenant_isolation ON tax_computations
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- TDS credits
CREATE TABLE tds_credits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    filing_id UUID NOT NULL REFERENCES itr_filings(id),
    deductor_tan VARCHAR(10) NOT NULL,
    deductor_name VARCHAR(255) NOT NULL,
    section VARCHAR(30) NOT NULL,
    tds_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    gross_payment DECIMAL(15,2) NOT NULL DEFAULT 0,
    tax_year VARCHAR(10) NOT NULL,
    matched_with_ais BOOLEAN NOT NULL DEFAULT false,
    ais_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    discrepancy DECIMAL(15,2) NOT NULL DEFAULT 0,
    discrepancy_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tds_credits_filing ON tds_credits(filing_id);
CREATE INDEX idx_tds_credits_tenant ON tds_credits(tenant_id);
CREATE INDEX idx_tds_credits_tan ON tds_credits(tenant_id, deductor_tan, tax_year);

ALTER TABLE tds_credits ENABLE ROW LEVEL SECURITY;
CREATE POLICY tds_credits_tenant_isolation ON tds_credits
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- AIS reconciliations
CREATE TABLE ais_reconciliations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    filing_id UUID NOT NULL REFERENCES itr_filings(id),
    pan VARCHAR(10) NOT NULL,
    tax_year VARCHAR(10) NOT NULL,
    source_type VARCHAR(30) NOT NULL,
    reported_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    ais_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    discrepancy DECIMAL(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ais_recon_filing ON ais_reconciliations(filing_id);
CREATE INDEX idx_ais_recon_tenant ON ais_reconciliations(tenant_id);
CREATE INDEX idx_ais_recon_pan_year ON ais_reconciliations(tenant_id, pan, tax_year);

ALTER TABLE ais_reconciliations ENABLE ROW LEVEL SECURITY;
CREATE POLICY ais_recon_tenant_isolation ON ais_reconciliations
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- Grant privileges to application role
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS ais_reconciliations;
DROP TABLE IF EXISTS tds_credits;
DROP TABLE IF EXISTS tax_computations;
DROP TABLE IF EXISTS deductions;
DROP TABLE IF EXISTS income_heads;
DROP TABLE IF EXISTS itr_filings;
DROP TABLE IF EXISTS taxpayers;
