#!/bin/bash
# Complai — PostgreSQL initialisation script
# Creates all logical databases and enables required extensions.
# Mounted at /docker-entrypoint-initdb.d/01-init.sh — runs once on first start.

set -euo pipefail

DATABASES=(
  identity_db
  tenant_db
  user_role_db
  master_data_db
  document_db
  audit_db
  gst_db
  gstr9_db
  einvoice_db
  ewb_db
  tds_db
  itr_db
  vendor_db
  recon_db
  ap_db
  billing_db
  secretarial_db
  workflow_db
  reporting_db
  rules_engine_db
  keycloak_db
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

echo "==> PostgreSQL initialisation complete."
