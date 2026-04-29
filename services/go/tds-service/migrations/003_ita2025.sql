-- +goose Up
-- ITA 2025 alignment: new columns for payment codes, sub-clauses, tax year.
-- Replaces ITA 1961 section numbering (192/194C/194I/194J/194Q/195)
-- with ITA 2025 sections (392, 393(1), 393(2), 393(3)) and 4-digit payment codes.

-- Rename legacy section column for historical records
ALTER TABLE tds_entries RENAME COLUMN section TO legacy_section_1961;
COMMENT ON COLUMN tds_entries.legacy_section_1961 IS 'DEPRECATED: ITA 1961 section. Retained for historical records only.';

-- Add ITA 2025 columns to tds_entries
ALTER TABLE tds_entries
    ADD COLUMN section_2025 TEXT NOT NULL DEFAULT '393(1)',
    ADD COLUMN payment_code TEXT NOT NULL DEFAULT '1024',
    ADD COLUMN sub_clause TEXT NOT NULL DEFAULT '',
    ADD COLUMN tax_year TEXT NOT NULL DEFAULT '2026-27',
    ADD COLUMN deductee_residency TEXT NOT NULL DEFAULT 'RESIDENT',
    ADD COLUMN deductee_type TEXT NOT NULL DEFAULT 'COMPANY',
    ADD COLUMN cess_rate DECIMAL(8,4) NOT NULL DEFAULT 0,
    ADD COLUMN surcharge_rate DECIMAL(8,4) NOT NULL DEFAULT 0,
    ADD COLUMN base_rate DECIMAL(8,4) NOT NULL DEFAULT 0,
    ADD COLUMN effective_rate DECIMAL(8,4) NOT NULL DEFAULT 0,
    ADD COLUMN dtaa_country_code TEXT,
    ADD COLUMN form_41_filed BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN trc_attached BOOLEAN NOT NULL DEFAULT FALSE;

-- Backfill section_2025 from legacy section
UPDATE tds_entries SET section_2025 = '392' WHERE legacy_section_1961 = '192';
UPDATE tds_entries SET section_2025 = '393(1)' WHERE legacy_section_1961 IN ('194C','194I','194J','194Q');
UPDATE tds_entries SET section_2025 = '393(2)' WHERE legacy_section_1961 = '195';

-- Backfill payment_code from legacy section (best-effort mapping)
UPDATE tds_entries SET payment_code = '1002' WHERE legacy_section_1961 = '192';
UPDATE tds_entries SET payment_code = '1024' WHERE legacy_section_1961 = '194C';
UPDATE tds_entries SET payment_code = '1009' WHERE legacy_section_1961 = '194I';
UPDATE tds_entries SET payment_code = '1027' WHERE legacy_section_1961 = '194J';
UPDATE tds_entries SET payment_code = '1031' WHERE legacy_section_1961 = '194Q';
UPDATE tds_entries SET payment_code = '1057' WHERE legacy_section_1961 = '195';

-- Drop assessment_year / previous_year if they exist (no-op if absent)
-- These columns were never added in prior migrations, so this is a safety net.

-- Rename legacy section column in tds_aggregates
ALTER TABLE tds_aggregates RENAME COLUMN section TO legacy_section_1961;
COMMENT ON COLUMN tds_aggregates.legacy_section_1961 IS 'DEPRECATED: ITA 1961 section. Retained for historical records only.';

ALTER TABLE tds_aggregates
    ADD COLUMN payment_code TEXT NOT NULL DEFAULT '1024';

-- Backfill aggregates payment_code
UPDATE tds_aggregates SET payment_code = '1002' WHERE legacy_section_1961 = '192';
UPDATE tds_aggregates SET payment_code = '1024' WHERE legacy_section_1961 = '194C';
UPDATE tds_aggregates SET payment_code = '1009' WHERE legacy_section_1961 = '194I';
UPDATE tds_aggregates SET payment_code = '1027' WHERE legacy_section_1961 = '194J';
UPDATE tds_aggregates SET payment_code = '1031' WHERE legacy_section_1961 = '194Q';
UPDATE tds_aggregates SET payment_code = '1057' WHERE legacy_section_1961 = '195';

-- Drop old unique constraint and add new one keyed on payment_code
ALTER TABLE tds_aggregates DROP CONSTRAINT IF EXISTS tds_aggregates_tenant_id_deductee_id_section_financial_year_key;
ALTER TABLE tds_aggregates ADD CONSTRAINT tds_aggregates_tenant_deductee_code_fy_key
    UNIQUE(tenant_id, deductee_id, payment_code, financial_year);

-- New indexes for ITA 2025 columns
DROP INDEX IF EXISTS idx_tds_entries_section;
CREATE INDEX idx_tds_entries_section_2025 ON tds_entries(tenant_id, section_2025);
CREATE INDEX idx_tds_entries_payment_code ON tds_entries(tenant_id, payment_code);
CREATE INDEX idx_tds_entries_tax_year ON tds_entries(tenant_id, tax_year, quarter);

DROP INDEX IF EXISTS idx_tds_aggregates_lookup;
CREATE INDEX idx_tds_aggregates_lookup ON tds_aggregates(tenant_id, deductee_id, payment_code, financial_year);

-- +goose Down
DROP INDEX IF EXISTS idx_tds_entries_tax_year;
DROP INDEX IF EXISTS idx_tds_entries_payment_code;
DROP INDEX IF EXISTS idx_tds_entries_section_2025;

ALTER TABLE tds_entries
    DROP COLUMN IF EXISTS section_2025,
    DROP COLUMN IF EXISTS payment_code,
    DROP COLUMN IF EXISTS sub_clause,
    DROP COLUMN IF EXISTS tax_year,
    DROP COLUMN IF EXISTS deductee_residency,
    DROP COLUMN IF EXISTS deductee_type,
    DROP COLUMN IF EXISTS cess_rate,
    DROP COLUMN IF EXISTS surcharge_rate,
    DROP COLUMN IF EXISTS base_rate,
    DROP COLUMN IF EXISTS effective_rate,
    DROP COLUMN IF EXISTS dtaa_country_code,
    DROP COLUMN IF EXISTS form_41_filed,
    DROP COLUMN IF EXISTS trc_attached;

ALTER TABLE tds_entries RENAME COLUMN legacy_section_1961 TO section;

ALTER TABLE tds_aggregates DROP CONSTRAINT IF EXISTS tds_aggregates_tenant_deductee_code_fy_key;
ALTER TABLE tds_aggregates DROP COLUMN IF EXISTS payment_code;
ALTER TABLE tds_aggregates RENAME COLUMN legacy_section_1961 TO section;
ALTER TABLE tds_aggregates ADD CONSTRAINT tds_aggregates_tenant_id_deductee_id_section_financial_year_key
    UNIQUE(tenant_id, deductee_id, section, financial_year);

CREATE INDEX idx_tds_entries_section ON tds_entries(tenant_id, section);
CREATE INDEX idx_tds_aggregates_lookup ON tds_aggregates(tenant_id, deductee_id, section, financial_year);
