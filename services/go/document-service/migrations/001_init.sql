-- +goose Up

CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    document_number VARCHAR(100),
    file_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL DEFAULT 'application/octet-stream',
    file_size BIGINT NOT NULL DEFAULT 0,
    s3_bucket VARCHAR(255) NOT NULL,
    s3_key VARCHAR(512) NOT NULL,
    encrypted_dek BYTEA,
    kms_key_arn VARCHAR(512),
    encryption_algo VARCHAR(20) DEFAULT 'AES-256-GCM',
    virus_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    ocr_status VARCHAR(20) NOT NULL DEFAULT 'none',
    ocr_result JSONB,
    tags JSONB DEFAULT '[]',
    metadata JSONB DEFAULT '{}',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_documents_tenant ON documents(tenant_id, created_at);
CREATE INDEX idx_documents_tenant_type ON documents(tenant_id, document_type);
CREATE INDEX idx_documents_s3 ON documents(s3_bucket, s3_key);

CREATE TABLE document_lineage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    parent_document_id UUID NOT NULL REFERENCES documents(id),
    child_document_id UUID NOT NULL REFERENCES documents(id),
    relationship VARCHAR(50) NOT NULL DEFAULT 'derived',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE documents ENABLE ROW LEVEL SECURITY;
ALTER TABLE document_lineage ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON documents USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON document_lineage USING (tenant_id = current_setting('app.tenant_id')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS document_lineage CASCADE;
DROP TABLE IF EXISTS documents CASCADE;
