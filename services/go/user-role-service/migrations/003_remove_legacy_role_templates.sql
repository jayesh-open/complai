-- +goose Up
DELETE FROM role_templates WHERE name IN ('tenant-admin', 'tax-manager', 'ap-clerk', 'tax-analyst', 'viewer');

-- +goose Down
-- Cannot recover legacy templates (data was unstructured, no clean reproduction). Down is no-op.
