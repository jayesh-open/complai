-- +goose Up

-- HSN codes (seeded, shared reference data but tenant-scoped for overrides)
CREATE TABLE hsn_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(8) NOT NULL,
    description TEXT NOT NULL,
    gst_rate NUMERIC(5,2) NOT NULL,
    effective_from DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_to DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, code, effective_from)
);

-- State codes
CREATE TABLE state_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(2) NOT NULL,
    name VARCHAR(100) NOT NULL,
    tin_code VARCHAR(2),
    UNIQUE(tenant_id, code)
);

-- Vendors
CREATE TABLE vendors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    pan VARCHAR(10),
    gstin VARCHAR(15),
    email VARCHAR(255),
    phone VARCHAR(20),
    address_line1 TEXT,
    address_line2 TEXT,
    city VARCHAR(100),
    state_code VARCHAR(2),
    pincode VARCHAR(6),
    bank_name VARCHAR(255),
    bank_account VARCHAR(20),
    bank_ifsc VARCHAR(11),
    kyc_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    compliance_score NUMERIC(5,2) DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    metadata JSONB DEFAULT '{}',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_vendors_tenant ON vendors(tenant_id);
CREATE INDEX idx_vendors_tenant_gstin ON vendors(tenant_id, gstin);
CREATE INDEX idx_vendors_tenant_pan ON vendors(tenant_id, pan);

-- Customers
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    pan VARCHAR(10),
    gstin VARCHAR(15),
    email VARCHAR(255),
    phone VARCHAR(20),
    address_line1 TEXT,
    city VARCHAR(100),
    state_code VARCHAR(2),
    pincode VARCHAR(6),
    payment_terms_days INT DEFAULT 30,
    credit_limit NUMERIC(15,2),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_customers_tenant ON customers(tenant_id);

-- Items (products/services catalog)
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    hsn_code VARCHAR(8) NOT NULL,
    unit_of_measure VARCHAR(10) NOT NULL DEFAULT 'NOS',
    unit_price NUMERIC(15,2),
    gst_rate NUMERIC(5,2),
    is_service BOOLEAN NOT NULL DEFAULT false,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_items_tenant ON items(tenant_id);
CREATE INDEX idx_items_tenant_hsn ON items(tenant_id, hsn_code);

-- Enable RLS on all tables
ALTER TABLE hsn_codes ENABLE ROW LEVEL SECURITY;
ALTER TABLE state_codes ENABLE ROW LEVEL SECURITY;
ALTER TABLE vendors ENABLE ROW LEVEL SECURITY;
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;
ALTER TABLE items ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON hsn_codes USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON state_codes USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON vendors USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON customers USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON items USING (tenant_id = current_setting('app.tenant_id')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS items CASCADE;
DROP TABLE IF EXISTS customers CASCADE;
DROP TABLE IF EXISTS vendors CASCADE;
DROP TABLE IF EXISTS state_codes CASCADE;
DROP TABLE IF EXISTS hsn_codes CASCADE;
