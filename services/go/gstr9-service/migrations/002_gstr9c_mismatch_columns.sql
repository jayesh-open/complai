-- +goose Up

ALTER TABLE gstr9c_mismatches ADD COLUMN section VARCHAR(10) NOT NULL DEFAULT 'II';
ALTER TABLE gstr9c_mismatches ADD COLUMN reason TEXT NOT NULL DEFAULT '';
ALTER TABLE gstr9c_mismatches ADD COLUMN suggested_action TEXT NOT NULL DEFAULT '';
ALTER TABLE gstr9c_mismatches ADD COLUMN resolved BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE gstr9c_mismatches ADD COLUMN resolved_reason TEXT NOT NULL DEFAULT '';
ALTER TABLE gstr9c_mismatches ADD COLUMN resolved_at TIMESTAMPTZ;
ALTER TABLE gstr9c_mismatches ADD COLUMN resolved_by UUID;

-- +goose Down

ALTER TABLE gstr9c_mismatches DROP COLUMN IF EXISTS resolved_by;
ALTER TABLE gstr9c_mismatches DROP COLUMN IF EXISTS resolved_at;
ALTER TABLE gstr9c_mismatches DROP COLUMN IF EXISTS resolved_reason;
ALTER TABLE gstr9c_mismatches DROP COLUMN IF EXISTS resolved;
ALTER TABLE gstr9c_mismatches DROP COLUMN IF EXISTS suggested_action;
ALTER TABLE gstr9c_mismatches DROP COLUMN IF EXISTS reason;
ALTER TABLE gstr9c_mismatches DROP COLUMN IF EXISTS section;
