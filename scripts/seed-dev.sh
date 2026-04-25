#!/bin/bash
# Complai — Development seed script
# Creates 1 platform tenant + 2 customer tenants + 5 users per tenant
# Plus roles, permissions, role templates, and Keycloak users.
# Runs psql via docker exec (no local psql required).

set -euo pipefail

PG_CONTAINER="${PG_CONTAINER:-complai-postgres}"
PG_USER="${PG_USER:-complai}"
KEYCLOAK_URL="${KEYCLOAK_URL:-http://localhost:8080}"

dpsql() {
  local db=$1
  shift
  docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$db" -q "$@"
}

# Fixed UUIDs
PLATFORM_TENANT="00000000-0000-0000-0000-000000000001"
TENANT_A="00000000-0000-0000-0000-000000000002"
TENANT_B="00000000-0000-0000-0000-000000000003"

echo "==> Seeding tenants..."

dpsql tenant_db <<SQL
INSERT INTO tenants (id, tenant_id, name, slug, tier, status, settings)
VALUES
  ('${PLATFORM_TENANT}', '${PLATFORM_TENANT}', 'Complai Platform', 'platform', 'pooled', 'active', '{}'),
  ('${TENANT_A}', '${TENANT_A}', 'Acme Industries', 'acme-industries', 'pooled', 'active', '{}'),
  ('${TENANT_B}', '${TENANT_B}', 'Beta Corp', 'beta-corp', 'pooled', 'active', '{}')
ON CONFLICT (slug) DO NOTHING;
SQL

echo "  -> 3 tenants seeded (platform + 2 customers)"

echo ""
echo "==> Seeding users into identity_db..."

seed_user() {
  local uid=$1 tid=$2 email=$3 first=$4 last=$5
  dpsql identity_db <<SQL
    INSERT INTO users (id, tenant_id, email, email_verified, first_name, last_name, status)
    VALUES ('${uid}', '${tid}', '${email}', true, '${first}', '${last}', 'active')
    ON CONFLICT (tenant_id, email) DO NOTHING;
SQL
}

# Platform tenant
seed_user "10000000-0000-0000-0000-000000000001" "$PLATFORM_TENANT" "admin@platform.complai.dev"   "Platform" "Admin"
seed_user "10000000-0000-0000-0000-000000000002" "$PLATFORM_TENANT" "manager@platform.complai.dev" "Platform" "Manager"
seed_user "10000000-0000-0000-0000-000000000003" "$PLATFORM_TENANT" "analyst@platform.complai.dev" "Platform" "Analyst"
seed_user "10000000-0000-0000-0000-000000000004" "$PLATFORM_TENANT" "clerk@platform.complai.dev"   "Platform" "Clerk"
seed_user "10000000-0000-0000-0000-000000000005" "$PLATFORM_TENANT" "viewer@platform.complai.dev"  "Platform" "Viewer"

# Tenant A
seed_user "20000000-0000-0000-0000-000000000001" "$TENANT_A" "admin@tenanta.complai.dev"   "Acme" "Admin"
seed_user "20000000-0000-0000-0000-000000000002" "$TENANT_A" "manager@tenanta.complai.dev" "Acme" "Manager"
seed_user "20000000-0000-0000-0000-000000000003" "$TENANT_A" "analyst@tenanta.complai.dev" "Acme" "Analyst"
seed_user "20000000-0000-0000-0000-000000000004" "$TENANT_A" "clerk@tenanta.complai.dev"   "Acme" "Clerk"
seed_user "20000000-0000-0000-0000-000000000005" "$TENANT_A" "viewer@tenanta.complai.dev"  "Acme" "Viewer"

# Tenant B
seed_user "30000000-0000-0000-0000-000000000001" "$TENANT_B" "admin@tenantb.complai.dev"   "Beta" "Admin"
seed_user "30000000-0000-0000-0000-000000000002" "$TENANT_B" "manager@tenantb.complai.dev" "Beta" "Manager"
seed_user "30000000-0000-0000-0000-000000000003" "$TENANT_B" "analyst@tenantb.complai.dev" "Beta" "Analyst"
seed_user "30000000-0000-0000-0000-000000000004" "$TENANT_B" "clerk@tenantb.complai.dev"   "Beta" "Clerk"
seed_user "30000000-0000-0000-0000-000000000005" "$TENANT_B" "viewer@tenantb.complai.dev"  "Beta" "Viewer"

echo "  -> 15 users seeded (5 per tenant)"

echo ""
echo "==> Seeding roles and permissions into user_role_db..."

for tid in "$PLATFORM_TENANT" "$TENANT_A" "$TENANT_B"; do
  dpsql user_role_db <<SQL
    INSERT INTO roles (tenant_id, name, display_name, description, is_system) VALUES
      ('${tid}', 'tenant-admin', 'Tenant Administrator', 'Full tenant access', true),
      ('${tid}', 'tax-manager', 'Tax Manager', 'Manage tax filings and returns', false),
      ('${tid}', 'tax-analyst', 'Tax Analyst', 'Prepare tax filings', false),
      ('${tid}', 'ap-clerk', 'AP Clerk', 'Accounts payable data entry', false),
      ('${tid}', 'viewer', 'Viewer', 'Read-only access', false)
    ON CONFLICT (tenant_id, name) DO NOTHING;

    INSERT INTO permissions (tenant_id, resource, action, description) VALUES
      ('${tid}', '*', '*', 'Superadmin wildcard'),
      ('${tid}', 'invoices', 'read', 'Read invoices'),
      ('${tid}', 'invoices', 'write', 'Create/edit invoices'),
      ('${tid}', 'invoices', 'delete', 'Delete invoices'),
      ('${tid}', 'filings', 'read', 'Read filings'),
      ('${tid}', 'filings', 'write', 'Create/submit filings'),
      ('${tid}', 'filings', 'approve', 'Approve filings'),
      ('${tid}', 'users', 'read', 'Read users'),
      ('${tid}', 'users', 'write', 'Manage users'),
      ('${tid}', 'settings', 'read', 'Read settings'),
      ('${tid}', 'settings', 'write', 'Manage settings')
    ON CONFLICT (tenant_id, resource, action) DO NOTHING;

    INSERT INTO role_permissions (tenant_id, role_id, permission_id)
    SELECT '${tid}', r.id, p.id
    FROM roles r, permissions p
    WHERE r.tenant_id = '${tid}' AND p.tenant_id = '${tid}'
      AND r.name = 'tenant-admin' AND p.resource = '*' AND p.action = '*'
    ON CONFLICT (tenant_id, role_id, permission_id) DO NOTHING;

    INSERT INTO role_permissions (tenant_id, role_id, permission_id)
    SELECT '${tid}', r.id, p.id
    FROM roles r, permissions p
    WHERE r.tenant_id = '${tid}' AND p.tenant_id = '${tid}'
      AND r.name = 'tax-manager'
      AND ((p.resource = 'invoices') OR (p.resource = 'filings'))
    ON CONFLICT (tenant_id, role_id, permission_id) DO NOTHING;

    INSERT INTO role_permissions (tenant_id, role_id, permission_id)
    SELECT '${tid}', r.id, p.id
    FROM roles r, permissions p
    WHERE r.tenant_id = '${tid}' AND p.tenant_id = '${tid}'
      AND r.name = 'tax-analyst'
      AND ((p.resource = 'invoices' AND p.action IN ('read','write'))
        OR (p.resource = 'filings' AND p.action IN ('read','write')))
    ON CONFLICT (tenant_id, role_id, permission_id) DO NOTHING;

    INSERT INTO role_permissions (tenant_id, role_id, permission_id)
    SELECT '${tid}', r.id, p.id
    FROM roles r, permissions p
    WHERE r.tenant_id = '${tid}' AND p.tenant_id = '${tid}'
      AND r.name = 'viewer'
      AND p.action = 'read'
    ON CONFLICT (tenant_id, role_id, permission_id) DO NOTHING;
SQL
done

echo "  -> Roles and permissions seeded for all 3 tenants"

echo ""
echo "==> Assigning roles to users in user_role_db..."

assign_role() {
  local tid=$1 uid=$2 role_name=$3
  dpsql user_role_db <<SQL
    INSERT INTO user_roles (tenant_id, user_id, role_id)
    SELECT '${tid}', '${uid}', r.id
    FROM roles r WHERE r.tenant_id = '${tid}' AND r.name = '${role_name}'
    ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING;
SQL
}

# Platform tenant
assign_role "$PLATFORM_TENANT" "10000000-0000-0000-0000-000000000001" "tenant-admin"
assign_role "$PLATFORM_TENANT" "10000000-0000-0000-0000-000000000002" "tax-manager"
assign_role "$PLATFORM_TENANT" "10000000-0000-0000-0000-000000000003" "tax-analyst"
assign_role "$PLATFORM_TENANT" "10000000-0000-0000-0000-000000000004" "ap-clerk"
assign_role "$PLATFORM_TENANT" "10000000-0000-0000-0000-000000000005" "viewer"

# Tenant A
assign_role "$TENANT_A" "20000000-0000-0000-0000-000000000001" "tenant-admin"
assign_role "$TENANT_A" "20000000-0000-0000-0000-000000000002" "tax-manager"
assign_role "$TENANT_A" "20000000-0000-0000-0000-000000000003" "tax-analyst"
assign_role "$TENANT_A" "20000000-0000-0000-0000-000000000004" "ap-clerk"
assign_role "$TENANT_A" "20000000-0000-0000-0000-000000000005" "viewer"

# Tenant B
assign_role "$TENANT_B" "30000000-0000-0000-0000-000000000001" "tenant-admin"
assign_role "$TENANT_B" "30000000-0000-0000-0000-000000000002" "tax-manager"
assign_role "$TENANT_B" "30000000-0000-0000-0000-000000000003" "tax-analyst"
assign_role "$TENANT_B" "30000000-0000-0000-0000-000000000004" "ap-clerk"
assign_role "$TENANT_B" "30000000-0000-0000-0000-000000000005" "viewer"

echo "  -> User-role assignments complete"

echo ""
echo "==> Seeding role templates..."

dpsql user_role_db <<SQL
INSERT INTO role_templates (name, display_name, description, permissions) VALUES
  ('tenant-admin', 'Tenant Administrator', 'Full access to all tenant resources', '[{"resource":"*","action":"*"}]'),
  ('tax-manager', 'Tax Manager', 'Manage all tax filings and invoices', '[{"resource":"invoices","action":"*"},{"resource":"filings","action":"*"}]'),
  ('tax-analyst', 'Tax Analyst', 'Prepare invoices and filings', '[{"resource":"invoices","action":"read"},{"resource":"invoices","action":"write"},{"resource":"filings","action":"read"},{"resource":"filings","action":"write"}]'),
  ('ap-clerk', 'AP Clerk', 'Accounts payable data entry', '[{"resource":"invoices","action":"read"},{"resource":"invoices","action":"write"}]'),
  ('viewer', 'Viewer', 'Read-only access across modules', '[{"resource":"invoices","action":"read"},{"resource":"filings","action":"read"},{"resource":"users","action":"read"},{"resource":"settings","action":"read"}]')
ON CONFLICT (name) DO NOTHING;
SQL

echo "  -> 5 role templates seeded"

echo ""
echo "==> Creating Keycloak users..."

get_admin_token() {
  curl -sf "${KEYCLOAK_URL}/realms/master/protocol/openid-connect/token" \
    -d "grant_type=password&username=admin&password=admin&client_id=admin-cli" \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['access_token'])" 2>/dev/null
}

ADMIN_TOKEN=$(get_admin_token)
if [ -z "$ADMIN_TOKEN" ]; then
  echo "  !! Could not get Keycloak admin token. Is Keycloak running?"
  echo "  !! Skipping Keycloak user creation. Existing realm users still work."
else
  # Keycloak 24 requires custom attributes to be registered in user profile
  PROFILE=$(curl -sf "${KEYCLOAK_URL}/admin/realms/complai/users/profile" \
    -H "Authorization: Bearer ${ADMIN_TOKEN}")
  HAS_TENANT_ID=$(echo "$PROFILE" | python3 -c "
import sys,json
p=json.load(sys.stdin)
print('yes' if any(a['name']=='tenant_id' for a in p.get('attributes',[])) else 'no')
" 2>/dev/null)
  if [ "$HAS_TENANT_ID" = "no" ]; then
    UPDATED=$(echo "$PROFILE" | python3 -c "
import sys,json
p=json.load(sys.stdin)
p['attributes'].append({'name':'tenant_id','displayName':'Tenant ID','permissions':{'view':['admin','user'],'edit':['admin']},'validations':{}})
print(json.dumps(p))
")
    curl -sf -o /dev/null -X PUT "${KEYCLOAK_URL}/admin/realms/complai/users/profile" \
      -H "Authorization: Bearer ${ADMIN_TOKEN}" \
      -H "Content-Type: application/json" -d "$UPDATED"
    echo "  -> Registered tenant_id in Keycloak user profile"
  fi
  create_kc_user() {
    local username=$1 email=$2 first=$3 last=$4 password=$5 tenant_id=$6 role=$7

    local resp
    resp=$(curl -sf -o /dev/null -w "%{http_code}" \
      "${KEYCLOAK_URL}/admin/realms/complai/users" \
      -H "Authorization: Bearer ${ADMIN_TOKEN}" \
      -H "Content-Type: application/json" \
      -d "{
        \"username\": \"${username}\",
        \"email\": \"${email}\",
        \"emailVerified\": true,
        \"enabled\": true,
        \"firstName\": \"${first}\",
        \"lastName\": \"${last}\",
        \"attributes\": {\"tenant_id\": [\"${tenant_id}\"]},
        \"credentials\": [{\"type\": \"password\", \"value\": \"${password}\", \"temporary\": false}]
      }" 2>/dev/null)

    if [ "$resp" = "201" ] || [ "$resp" = "409" ]; then
      local user_id
      user_id=$(curl -sf "${KEYCLOAK_URL}/admin/realms/complai/users?username=${username}&exact=true" \
        -H "Authorization: Bearer ${ADMIN_TOKEN}" \
        | python3 -c "import sys,json; print(json.load(sys.stdin)[0]['id'])" 2>/dev/null)

      if [ -n "$user_id" ]; then
        # Ensure tenant_id attribute is set (Keycloak 24 user profile may strip on create)
        local full_user
        full_user=$(curl -sf "${KEYCLOAK_URL}/admin/realms/complai/users/${user_id}" \
          -H "Authorization: Bearer ${ADMIN_TOKEN}")
        local patched
        patched=$(echo "$full_user" | python3 -c "
import sys,json
u=json.load(sys.stdin)
attrs=u.get('attributes') or {}
attrs['tenant_id']=['${tenant_id}']
u['attributes']=attrs
print(json.dumps(u))
")
        curl -sf -o /dev/null -X PUT \
          "${KEYCLOAK_URL}/admin/realms/complai/users/${user_id}" \
          -H "Authorization: Bearer ${ADMIN_TOKEN}" \
          -H "Content-Type: application/json" -d "$patched"

        if [ -n "$role" ]; then
          local role_json
          role_json=$(curl -sf "${KEYCLOAK_URL}/admin/realms/complai/roles/${role}" \
            -H "Authorization: Bearer ${ADMIN_TOKEN}" 2>/dev/null)

          if [ -n "$role_json" ]; then
            curl -sf -o /dev/null \
              "${KEYCLOAK_URL}/admin/realms/complai/users/${user_id}/role-mappings/realm" \
              -H "Authorization: Bearer ${ADMIN_TOKEN}" \
              -H "Content-Type: application/json" \
              -d "[${role_json}]" 2>/dev/null || true
          fi
        fi
      fi
      echo "    -> ${username} (${role:-no-role})"
    else
      echo "    !! Failed to create ${username} (HTTP ${resp})"
    fi
  }

  # Platform tenant users
  create_kc_user "platform-admin"   "admin@platform.complai.dev"   "Platform" "Admin"   "password" "$PLATFORM_TENANT" "complai-admin"
  create_kc_user "platform-manager" "manager@platform.complai.dev" "Platform" "Manager" "password" "$PLATFORM_TENANT" "reviewer"
  create_kc_user "platform-analyst" "analyst@platform.complai.dev" "Platform" "Analyst" "password" "$PLATFORM_TENANT" "preparer"
  create_kc_user "platform-clerk"   "clerk@platform.complai.dev"   "Platform" "Clerk"   "password" "$PLATFORM_TENANT" "preparer"
  create_kc_user "platform-viewer"  "viewer@platform.complai.dev"  "Platform" "Viewer"  "password" "$PLATFORM_TENANT" "viewer"

  # Tenant A users
  create_kc_user "acme-admin"   "admin@tenanta.complai.dev"   "Acme" "Admin"   "password" "$TENANT_A" "tenant-admin"
  create_kc_user "acme-manager" "manager@tenanta.complai.dev" "Acme" "Manager" "password" "$TENANT_A" "reviewer"
  create_kc_user "acme-analyst" "analyst@tenanta.complai.dev" "Acme" "Analyst" "password" "$TENANT_A" "preparer"
  create_kc_user "acme-clerk"   "clerk@tenanta.complai.dev"   "Acme" "Clerk"   "password" "$TENANT_A" "preparer"
  create_kc_user "acme-viewer"  "viewer@tenanta.complai.dev"  "Acme" "Viewer"  "password" "$TENANT_A" "viewer"

  # Tenant B users
  create_kc_user "beta-admin"   "admin@tenantb.complai.dev"   "Beta" "Admin"   "password" "$TENANT_B" "tenant-admin"
  create_kc_user "beta-manager" "manager@tenantb.complai.dev" "Beta" "Manager" "password" "$TENANT_B" "reviewer"
  create_kc_user "beta-analyst" "analyst@tenantb.complai.dev" "Beta" "Analyst" "password" "$TENANT_B" "preparer"
  create_kc_user "beta-clerk"   "clerk@tenantb.complai.dev"   "Beta" "Clerk"   "password" "$TENANT_B" "preparer"
  create_kc_user "beta-viewer"  "viewer@tenantb.complai.dev"  "Beta" "Viewer"  "password" "$TENANT_B" "viewer"
fi

echo ""
echo "==> Seed complete!"
echo "  Platform tenant: ${PLATFORM_TENANT}"
echo "  Tenant A (Acme): ${TENANT_A}"
echo "  Tenant B (Beta): ${TENANT_B}"
echo "  Users per tenant: 5 (admin, manager, analyst, clerk, viewer)"
echo "  Keycloak password for all seeded users: password"
echo "  Existing Keycloak users: dev-admin/admin, dev-preparer/preparer"
