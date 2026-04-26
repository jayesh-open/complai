package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/recon-service/internal/domain"
)

var _ Repository = (*Store)(nil)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Pool() *pgxpool.Pool { return s.pool }

func setTenantID(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID) error {
	_, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID.String()))
	return err
}

func (s *Store) CreateRun(ctx context.Context, tenantID uuid.UUID, run *domain.ReconRun) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	run.ID = uuid.New()
	run.TenantID = tenantID
	run.RequestID = uuid.New()

	err = tx.QueryRow(ctx,
		`INSERT INTO recon_runs (id, tenant_id, gstin, return_period, status, pr_count, gstr2b_count,
		 matched, mismatch, partial, missing_2b, missing_pr, duplicate, started_at, request_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		 RETURNING created_at`,
		run.ID, run.TenantID, run.GSTIN, run.ReturnPeriod, run.Status, run.PRCount, run.GSTR2BCount,
		run.Matched, run.Mismatch, run.Partial, run.Missing2B, run.MissingPR, run.Duplicate,
		run.StartedAt, run.RequestID,
	).Scan(&run.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert run: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) GetRun(ctx context.Context, tenantID uuid.UUID, runID uuid.UUID) (*domain.ReconRun, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var r domain.ReconRun
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, gstin, return_period, status, pr_count, gstr2b_count,
		        matched, mismatch, partial, missing_2b, missing_pr, duplicate,
		        started_at, completed_at, request_id, created_at
		 FROM recon_runs WHERE id = $1`, runID,
	).Scan(&r.ID, &r.TenantID, &r.GSTIN, &r.ReturnPeriod, &r.Status, &r.PRCount, &r.GSTR2BCount,
		&r.Matched, &r.Mismatch, &r.Partial, &r.Missing2B, &r.MissingPR, &r.Duplicate,
		&r.StartedAt, &r.CompletedAt, &r.RequestID, &r.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get run: %w", err)
	}

	return &r, tx.Commit(ctx)
}

func (s *Store) UpdateRun(ctx context.Context, tenantID uuid.UUID, run *domain.ReconRun) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE recon_runs SET status = $1, pr_count = $2, gstr2b_count = $3,
		 matched = $4, mismatch = $5, partial = $6, missing_2b = $7, missing_pr = $8,
		 duplicate = $9, completed_at = $10
		 WHERE id = $11`,
		run.Status, run.PRCount, run.GSTR2BCount,
		run.Matched, run.Mismatch, run.Partial, run.Missing2B, run.MissingPR,
		run.Duplicate, run.CompletedAt, run.ID)
	if err != nil {
		return fmt.Errorf("update run: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) BulkInsertMatches(ctx context.Context, tenantID uuid.UUID, matches []domain.ReconMatch) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	for _, m := range matches {
		_, err := tx.Exec(ctx,
			`INSERT INTO recon_matches (id, tenant_id, run_id, gstin, return_period,
			 pr_invoice_number, pr_invoice_date, pr_vendor_gstin, pr_amount, pr_hsn, pr_source_id,
			 gstr2b_invoice_number, gstr2b_invoice_date, gstr2b_supplier_gstin, gstr2b_amount, gstr2b_hsn,
			 match_type, match_confidence, reason_codes, status)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`,
			m.ID, tenantID, m.RunID, m.GSTIN, m.ReturnPeriod,
			m.PRInvoiceNumber, m.PRInvoiceDate, m.PRVendorGSTIN, m.PRAmount, m.PRHSN, m.PRSourceID,
			m.GSTR2BInvoiceNumber, m.GSTR2BInvoiceDate, m.GSTR2BSupplierGSTIN, m.GSTR2BAmount, m.GSTR2BHSN,
			m.MatchType, m.MatchConfidence, m.ReasonCodes, m.Status,
		)
		if err != nil {
			return fmt.Errorf("insert match: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (s *Store) ListMatches(ctx context.Context, tenantID uuid.UUID, runID uuid.UUID, matchType string, status string, limit, offset int) ([]domain.ReconMatch, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	query := `SELECT id, tenant_id, run_id, gstin, return_period,
	          pr_invoice_number, pr_invoice_date, pr_vendor_gstin, pr_amount, pr_hsn, pr_source_id,
	          gstr2b_invoice_number, gstr2b_invoice_date, gstr2b_supplier_gstin, gstr2b_amount, gstr2b_hsn,
	          match_type, match_confidence, reason_codes, status, accepted_by, accepted_at, created_at, updated_at
	          FROM recon_matches WHERE run_id = $1`
	args := []interface{}{runID}
	argIdx := 2

	if matchType != "" {
		query += fmt.Sprintf(" AND match_type = $%d", argIdx)
		args = append(args, matchType)
		argIdx++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	query += " ORDER BY match_type, created_at"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, limit)
		argIdx++
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, offset)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list matches: %w", err)
	}
	defer rows.Close()

	var matches []domain.ReconMatch
	for rows.Next() {
		var m domain.ReconMatch
		if err := rows.Scan(&m.ID, &m.TenantID, &m.RunID, &m.GSTIN, &m.ReturnPeriod,
			&m.PRInvoiceNumber, &m.PRInvoiceDate, &m.PRVendorGSTIN, &m.PRAmount, &m.PRHSN, &m.PRSourceID,
			&m.GSTR2BInvoiceNumber, &m.GSTR2BInvoiceDate, &m.GSTR2BSupplierGSTIN, &m.GSTR2BAmount, &m.GSTR2BHSN,
			&m.MatchType, &m.MatchConfidence, &m.ReasonCodes, &m.Status,
			&m.AcceptedBy, &m.AcceptedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan match: %w", err)
		}
		matches = append(matches, m)
	}

	return matches, tx.Commit(ctx)
}

func (s *Store) GetMatch(ctx context.Context, tenantID uuid.UUID, matchID uuid.UUID) (*domain.ReconMatch, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var m domain.ReconMatch
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, run_id, gstin, return_period,
		        pr_invoice_number, pr_invoice_date, pr_vendor_gstin, pr_amount, pr_hsn, pr_source_id,
		        gstr2b_invoice_number, gstr2b_invoice_date, gstr2b_supplier_gstin, gstr2b_amount, gstr2b_hsn,
		        match_type, match_confidence, reason_codes, status, accepted_by, accepted_at, created_at, updated_at
		 FROM recon_matches WHERE id = $1`, matchID,
	).Scan(&m.ID, &m.TenantID, &m.RunID, &m.GSTIN, &m.ReturnPeriod,
		&m.PRInvoiceNumber, &m.PRInvoiceDate, &m.PRVendorGSTIN, &m.PRAmount, &m.PRHSN, &m.PRSourceID,
		&m.GSTR2BInvoiceNumber, &m.GSTR2BInvoiceDate, &m.GSTR2BSupplierGSTIN, &m.GSTR2BAmount, &m.GSTR2BHSN,
		&m.MatchType, &m.MatchConfidence, &m.ReasonCodes, &m.Status,
		&m.AcceptedBy, &m.AcceptedAt, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get match: %w", err)
	}

	return &m, tx.Commit(ctx)
}

func (s *Store) UpdateMatchStatus(ctx context.Context, tenantID uuid.UUID, matchID uuid.UUID, status domain.MatchStatus, acceptedBy *uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	var acceptedAt *time.Time
	if status == domain.MatchStatusAccepted {
		now := time.Now().UTC()
		acceptedAt = &now
	}

	_, err = tx.Exec(ctx,
		`UPDATE recon_matches SET status = $1, accepted_by = $2, accepted_at = $3, updated_at = now()
		 WHERE id = $4`,
		status, acceptedBy, acceptedAt, matchID)
	if err != nil {
		return fmt.Errorf("update match status: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) GetBucketSummary(ctx context.Context, tenantID uuid.UUID, runID uuid.UUID) (*domain.BucketSummary, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT match_type, COUNT(*) FROM recon_matches WHERE run_id = $1 GROUP BY match_type`, runID)
	if err != nil {
		return nil, fmt.Errorf("bucket summary: %w", err)
	}
	defer rows.Close()

	summary := &domain.BucketSummary{}
	for rows.Next() {
		var matchType string
		var count int
		if err := rows.Scan(&matchType, &count); err != nil {
			return nil, fmt.Errorf("scan bucket: %w", err)
		}
		switch domain.MatchType(matchType) {
		case domain.MatchTypeDirect:
			summary.Matched = count
		case domain.MatchTypeProbable:
			summary.Mismatch = count
		case domain.MatchTypePartial:
			summary.Partial = count
		case domain.MatchTypeMissing2B:
			summary.Missing2B = count
		case domain.MatchTypeMissingPR:
			summary.MissingPR = count
		case domain.MatchTypeDuplicate:
			summary.Duplicate = count
		}
	}

	return summary, tx.Commit(ctx)
}

func (s *Store) CreateIMSAction(ctx context.Context, tenantID uuid.UUID, action *domain.IMSAction) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	action.ID = uuid.New()
	action.TenantID = tenantID

	err = tx.QueryRow(ctx,
		`INSERT INTO ims_actions (id, tenant_id, gstin, return_period, invoice_id, action, reason, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING created_at`,
		action.ID, action.TenantID, action.GSTIN, action.ReturnPeriod,
		action.InvoiceID, action.Action, action.Reason, action.CreatedBy,
	).Scan(&action.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert ims action: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) ListIMSActions(ctx context.Context, tenantID uuid.UUID, gstin, returnPeriod string) ([]domain.IMSAction, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, gstin, return_period, invoice_id, action, reason,
		        synced_to_gstn, synced_at, created_by, created_at
		 FROM ims_actions WHERE gstin = $1 AND return_period = $2 ORDER BY created_at`, gstin, returnPeriod)
	if err != nil {
		return nil, fmt.Errorf("list ims actions: %w", err)
	}
	defer rows.Close()

	var actions []domain.IMSAction
	for rows.Next() {
		var a domain.IMSAction
		if err := rows.Scan(&a.ID, &a.TenantID, &a.GSTIN, &a.ReturnPeriod, &a.InvoiceID, &a.Action, &a.Reason,
			&a.SyncedToGSTN, &a.SyncedAt, &a.CreatedBy, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan ims action: %w", err)
		}
		actions = append(actions, a)
	}

	return actions, tx.Commit(ctx)
}
