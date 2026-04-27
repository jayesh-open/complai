#!/bin/bash
# Complai — PostgreSQL initialisation script
# Creates all logical databases and enables required extensions.
# Mounted at /docker-entrypoint-initdb.d/01-init.sh — runs once on first start.

set -euo pipefail

DATABASES=(
  # Parts 1–3: Identity, Tenant, Platform services
  identity_db
  tenant_db
  user_role_db
  master_data_db
  document_db
  audit_db
  workflow_db
  notification_db
  rules_engine_db
  keycloak_db

  # Part 5: GST Returns
  gst_db

  # Part 6: Vendor Compliance
  vendor_compliance_db

  # Part 7: Reconciliation
  recon_db

  # Part 8: e-Invoicing + E-Way Bill
  einvoice_db
  ewb_db

  # Part 9: TDS (forward-provisioned)
  tds_db

  # Part 10: ITR + GSTR-9/9C (forward-provisioned)
  gstr9_db
  itr_db

  # Part 14: Reporting (forward-provisioned)
  reporting_db

  # Future: Secretarial (forward-provisioned)
  secretarial_db
)

EXTENSIONS=(
  "uuid-ossp"
  pgcrypto
  pg_trgm
)

echo "==> Creating logical databases..."

for db in "${DATABASES[@]}"; do
  echo "  -> Creating database: ${db}"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-SQL
    SELECT 'CREATE DATABASE ${db}'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '${db}')\gexec
SQL

  for ext in "${EXTENSIONS[@]}"; do
    echo "     -> Enabling extension: ${ext} in ${db}"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "${db}" <<-SQL
      CREATE EXTENSION IF NOT EXISTS "${ext}";
SQL
  done
done

# ---------------------------------------------------------------------------
# Application role (non-superuser) — RLS policies apply to this role
# ---------------------------------------------------------------------------
echo ""
echo "==> Creating application role..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-SQL
  DO \$\$
  BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'complai_app') THEN
      CREATE ROLE complai_app WITH LOGIN PASSWORD 'complai_app_dev';
    END IF;
  END
  \$\$;
SQL

for db in "${DATABASES[@]}"; do
  echo "  -> Granting privileges on ${db} to complai_app"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "${db}" <<-SQL
    GRANT CONNECT ON DATABASE ${db} TO complai_app;
    GRANT USAGE ON SCHEMA public TO complai_app;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO complai_app;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO complai_app;
SQL
done

echo "==> PostgreSQL initialisation complete."
