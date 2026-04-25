-- +goose Up

CREATE TABLE notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    channel VARCHAR(20) NOT NULL DEFAULT 'email',
    subject TEXT,
    body TEXT NOT NULL,
    variables JSONB DEFAULT '[]',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, name)
);

CREATE TABLE notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    email_enabled BOOLEAN NOT NULL DEFAULT true,
    digest_enabled BOOLEAN NOT NULL DEFAULT false,
    quiet_hours_start TIME,
    quiet_hours_end TIME,
    email_address VARCHAR(255),
    email_valid BOOLEAN NOT NULL DEFAULT true,
    bounce_count INT NOT NULL DEFAULT 0,
    unsubscribe_token UUID DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, user_id)
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    template_id UUID REFERENCES notification_templates(id),
    channel VARCHAR(20) NOT NULL DEFAULT 'email',
    subject TEXT,
    body TEXT,
    recipient VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'queued',
    sent_at TIMESTAMPTZ,
    failed_reason TEXT,
    digest_group VARCHAR(100),
    digest_batch_id UUID,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_notifications_tenant ON notifications(tenant_id, created_at);
CREATE INDEX idx_notifications_status ON notifications(tenant_id, status);
CREATE INDEX idx_notifications_digest ON notifications(tenant_id, user_id, digest_group, status);

CREATE TABLE notification_bounces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    notification_id UUID REFERENCES notifications(id),
    bounce_type VARCHAR(50) NOT NULL,
    bounce_subtype VARCHAR(50),
    email_address VARCHAR(255) NOT NULL,
    diagnostic TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE notification_templates ENABLE ROW LEVEL SECURITY;
ALTER TABLE notification_preferences ENABLE ROW LEVEL SECURITY;
ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;
ALTER TABLE notification_bounces ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON notification_templates USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON notification_preferences USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON notifications USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY tenant_isolation ON notification_bounces USING (tenant_id = current_setting('app.tenant_id')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO complai_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO complai_app;

-- +goose Down
DROP TABLE IF EXISTS notification_bounces CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS notification_preferences CASCADE;
DROP TABLE IF EXISTS notification_templates CASCADE;
