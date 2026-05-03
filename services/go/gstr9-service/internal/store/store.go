package store

import (
	"context"
	"fmt"

	"github.com/complai/complai/services/go/gstr9-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type PgStore struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PgStore {
	return &PgStore{pool: pool}
}

func (s *PgStore) CreateFiling(ctx context.Context, tenantID uuid.UUID, f *domain.GSTR9Filing) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO gstr9_filings
		(id, tenant_id, gstin, financial_year, status, aggregate_turnover, is_mandatory,
		 arn, filed_at, filed_by, approved_by, approved_at, request_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		f.ID, tenantID, f.GSTIN, f.FinancialYear, f.Status, f.AggregateTurnover, f.IsMandatory,
		f.ARN, f.FiledAt, f.FiledBy, f.ApprovedBy, f.ApprovedAt, f.RequestID, f.CreatedAt, f.UpdatedAt)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) GetFiling(ctx context.Context, tenantID, id uuid.UUID) (*domain.GSTR9Filing, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, err
	}

	var f domain.GSTR9Filing
	err = tx.QueryRow(ctx, `SELECT id, tenant_id, gstin, financial_year, status, aggregate_turnover,
		is_mandatory, arn, filed_at, filed_by, approved_by, approved_at, request_id, created_at, updated_at
		FROM gstr9_filings WHERE id = $1`, id).Scan(
		&f.ID, &f.TenantID, &f.GSTIN, &f.FinancialYear, &f.Status, &f.AggregateTurnover,
		&f.IsMandatory, &f.ARN, &f.FiledAt, &f.FiledBy, &f.ApprovedBy, &f.ApprovedAt,
		&f.RequestID, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)
	return &f, nil
}

func (s *PgStore) UpdateFilingStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.FilingStatus) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE gstr9_filings SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) ListFilings(ctx context.Context, tenantID uuid.UUID, gstin, fy string, limit, offset int) ([]domain.GSTR9Filing, int, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, 0, err
	}

	where := "WHERE 1=1"
	args := []interface{}{}
	argN := 1
	if gstin != "" {
		where += fmt.Sprintf(" AND gstin = $%d", argN)
		args = append(args, gstin)
		argN++
	}
	if fy != "" {
		where += fmt.Sprintf(" AND financial_year = $%d", argN)
		args = append(args, fy)
		argN++
	}

	var total int
	if err := tx.QueryRow(ctx, "SELECT COUNT(*) FROM gstr9_filings "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	q := fmt.Sprintf(`SELECT id, tenant_id, gstin, financial_year, status, aggregate_turnover,
		is_mandatory, arn, filed_at, filed_by, approved_by, approved_at, request_id, created_at, updated_at
		FROM gstr9_filings %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, where, argN, argN+1)

	rows, err := tx.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []domain.GSTR9Filing
	for rows.Next() {
		var f domain.GSTR9Filing
		if err := rows.Scan(&f.ID, &f.TenantID, &f.GSTIN, &f.FinancialYear, &f.Status, &f.AggregateTurnover,
			&f.IsMandatory, &f.ARN, &f.FiledAt, &f.FiledBy, &f.ApprovedBy, &f.ApprovedAt,
			&f.RequestID, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, f)
	}
	tx.Commit(ctx)
	return out, total, nil
}

func (s *PgStore) CreateTableData(ctx context.Context, tenantID uuid.UUID, td *domain.GSTR9TableData) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO gstr9_table_data
		(id, tenant_id, filing_id, part_number, table_number, description,
		 taxable_value, cgst, sgst, igst, cess, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		td.ID, tenantID, td.FilingID, td.PartNumber, td.TableNumber, td.Description,
		td.TaxableValue, td.CGST, td.SGST, td.IGST, td.Cess, td.CreatedAt)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) ListTableData(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.GSTR9TableData, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `SELECT id, tenant_id, filing_id, part_number, table_number, description,
		taxable_value, cgst, sgst, igst, cess, created_at
		FROM gstr9_table_data WHERE filing_id = $1 ORDER BY part_number, table_number`, filingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.GSTR9TableData
	for rows.Next() {
		var td domain.GSTR9TableData
		if err := rows.Scan(&td.ID, &td.TenantID, &td.FilingID, &td.PartNumber, &td.TableNumber, &td.Description,
			&td.TaxableValue, &td.CGST, &td.SGST, &td.IGST, &td.Cess, &td.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, td)
	}
	tx.Commit(ctx)
	return out, nil
}

func (s *PgStore) DeleteTableData(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM gstr9_table_data WHERE filing_id = $1`, filingID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) CreateAuditLog(ctx context.Context, tenantID uuid.UUID, log *domain.GSTR9AuditLog) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO gstr9_audit_logs
		(id, tenant_id, filing_id, action, details, actor_id, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		log.ID, tenantID, log.FilingID, log.Action, log.Details, log.ActorID, log.CreatedAt)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) ListAuditLogs(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.GSTR9AuditLog, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `SELECT id, tenant_id, filing_id, action, details, actor_id, created_at
		FROM gstr9_audit_logs WHERE filing_id = $1 ORDER BY created_at`, filingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.GSTR9AuditLog
	for rows.Next() {
		var l domain.GSTR9AuditLog
		if err := rows.Scan(&l.ID, &l.TenantID, &l.FilingID, &l.Action, &l.Details, &l.ActorID, &l.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	tx.Commit(ctx)
	return out, nil
}

func (s *PgStore) CreateGSTR9CFiling(ctx context.Context, tenantID uuid.UUID, f *domain.GSTR9CFiling) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO gstr9c_filings
		(id, tenant_id, gstr9_filing_id, status, audited_turnover, unreconciled_amount,
		 is_self_certified, certified_at, certified_by, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		f.ID, tenantID, f.GSTR9FilingID, f.Status, f.AuditedTurnover, f.UnreconciledAmount,
		f.IsSelfCertified, f.CertifiedAt, f.CertifiedBy, f.CreatedAt, f.UpdatedAt)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) GetGSTR9CFiling(ctx context.Context, tenantID, id uuid.UUID) (*domain.GSTR9CFiling, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, err
	}

	var f domain.GSTR9CFiling
	err = tx.QueryRow(ctx, `SELECT id, tenant_id, gstr9_filing_id, status, audited_turnover,
		unreconciled_amount, is_self_certified, certified_at, certified_by, created_at, updated_at
		FROM gstr9c_filings WHERE id = $1`, id).Scan(
		&f.ID, &f.TenantID, &f.GSTR9FilingID, &f.Status, &f.AuditedTurnover,
		&f.UnreconciledAmount, &f.IsSelfCertified, &f.CertifiedAt, &f.CertifiedBy,
		&f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)
	return &f, nil
}

func (s *PgStore) GetGSTR9CFilingByGSTR9ID(ctx context.Context, tenantID, gstr9FilingID uuid.UUID) (*domain.GSTR9CFiling, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, err
	}

	var f domain.GSTR9CFiling
	err = tx.QueryRow(ctx, `SELECT id, tenant_id, gstr9_filing_id, status, audited_turnover,
		unreconciled_amount, is_self_certified, certified_at, certified_by, created_at, updated_at
		FROM gstr9c_filings WHERE gstr9_filing_id = $1`, gstr9FilingID).Scan(
		&f.ID, &f.TenantID, &f.GSTR9FilingID, &f.Status, &f.AuditedTurnover,
		&f.UnreconciledAmount, &f.IsSelfCertified, &f.CertifiedAt, &f.CertifiedBy,
		&f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)
	return &f, nil
}

func (s *PgStore) UpdateGSTR9CStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.GSTR9CStatus) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE gstr9c_filings SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) UpdateGSTR9CUnreconciled(ctx context.Context, tenantID, id uuid.UUID, amount decimal.Decimal) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE gstr9c_filings SET unreconciled_amount = $1, updated_at = NOW() WHERE id = $2`,
		amount, id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) CertifyGSTR9C(ctx context.Context, tenantID, id, certifiedBy uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE gstr9c_filings
		SET status = $1, is_self_certified = TRUE, certified_at = NOW(), certified_by = $2, updated_at = NOW()
		WHERE id = $3`, domain.GSTR9CStatusCertified, certifiedBy, id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) CreateMismatch(ctx context.Context, tenantID uuid.UUID, m *domain.GSTR9CMismatch) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO gstr9c_mismatches
		(id, tenant_id, gstr9c_filing_id, section, category, description,
		 books_amount, gstr9_amount, difference, severity, reason, suggested_action,
		 resolved, resolved_reason, resolved_at, resolved_by, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
		m.ID, tenantID, m.GSTR9CFilingID, m.Section, m.Category, m.Description,
		m.BooksAmount, m.GSTR9Amount, m.Difference, m.Severity, m.Reason, m.SuggestedAction,
		m.Resolved, m.ResolvedReason, m.ResolvedAt, m.ResolvedBy, m.CreatedAt)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) ListMismatches(ctx context.Context, tenantID, gstr9cFilingID uuid.UUID) ([]domain.GSTR9CMismatch, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `SELECT id, tenant_id, gstr9c_filing_id, section, category, description,
		books_amount, gstr9_amount, difference, severity, reason, suggested_action,
		resolved, resolved_reason, resolved_at, resolved_by, created_at
		FROM gstr9c_mismatches WHERE gstr9c_filing_id = $1 ORDER BY section, category`, gstr9cFilingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.GSTR9CMismatch
	for rows.Next() {
		var m domain.GSTR9CMismatch
		if err := rows.Scan(&m.ID, &m.TenantID, &m.GSTR9CFilingID, &m.Section, &m.Category, &m.Description,
			&m.BooksAmount, &m.GSTR9Amount, &m.Difference, &m.Severity, &m.Reason, &m.SuggestedAction,
			&m.Resolved, &m.ResolvedReason, &m.ResolvedAt, &m.ResolvedBy, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	tx.Commit(ctx)
	return out, nil
}

func (s *PgStore) GetMismatch(ctx context.Context, tenantID, id uuid.UUID) (*domain.GSTR9CMismatch, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, err
	}

	var m domain.GSTR9CMismatch
	err = tx.QueryRow(ctx, `SELECT id, tenant_id, gstr9c_filing_id, section, category, description,
		books_amount, gstr9_amount, difference, severity, reason, suggested_action,
		resolved, resolved_reason, resolved_at, resolved_by, created_at
		FROM gstr9c_mismatches WHERE id = $1`, id).Scan(
		&m.ID, &m.TenantID, &m.GSTR9CFilingID, &m.Section, &m.Category, &m.Description,
		&m.BooksAmount, &m.GSTR9Amount, &m.Difference, &m.Severity, &m.Reason, &m.SuggestedAction,
		&m.Resolved, &m.ResolvedReason, &m.ResolvedAt, &m.ResolvedBy, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)
	return &m, nil
}

func (s *PgStore) ResolveMismatch(ctx context.Context, tenantID, id uuid.UUID, reason string, resolvedBy uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE gstr9c_mismatches
		SET resolved = TRUE, resolved_reason = $1, resolved_at = NOW(), resolved_by = $2
		WHERE id = $3`, reason, resolvedBy, id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) DeleteMismatches(ctx context.Context, tenantID uuid.UUID, gstr9cFilingID uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM gstr9c_mismatches WHERE gstr9c_filing_id = $1`, gstr9cFilingID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

var _ Repository = (*PgStore)(nil)
