-- +goose Up
-- Remove legacy kebab-case roles from the 3 dev tenants (UUIDs 0001, 0002, 0003).
-- They will be re-seeded via the new template-based seed flow.
DELETE FROM role_permissions
  WHERE role_id IN (SELECT id FROM roles WHERE name IN ('tenant-admin', 'tax-manager', 'ap-clerk', 'tax-analyst', 'viewer'));
DELETE FROM user_roles
  WHERE role_id IN (SELECT id FROM roles WHERE name IN ('tenant-admin', 'tax-manager', 'ap-clerk', 'tax-analyst', 'viewer'));
DELETE FROM roles WHERE name IN ('tenant-admin', 'tax-manager', 'ap-clerk', 'tax-analyst', 'viewer');

-- +goose Down
-- Cannot recover; legacy data was placeholder.
