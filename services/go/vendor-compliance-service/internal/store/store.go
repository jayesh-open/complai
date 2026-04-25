package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/vendor-compliance-service/internal/domain"
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

func (s *Store) UpsertVendorSnapshot(ctx context.Context, tenantID uuid.UUID, v *domain.VendorSnapshot) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	v.TenantID = tenantID
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO vendor_snapshots (id, tenant_id, vendor_id, name, legal_name, trade_name,
		 pan, gstin, tan, state, state_code, category, registration_status, msme_registered,
		 email, phone, address, synced_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		 ON CONFLICT (tenant_id, vendor_id) DO UPDATE SET
		 name = EXCLUDED.name, legal_name = EXCLUDED.legal_name, trade_name = EXCLUDED.trade_name,
		 pan = EXCLUDED.pan, gstin = EXCLUDED.gstin, tan = EXCLUDED.tan,
		 state = EXCLUDED.state, state_code = EXCLUDED.state_code, category = EXCLUDED.category,
		 registration_status = EXCLUDED.registration_status, msme_registered = EXCLUDED.msme_registered,
		 email = EXCLUDED.email, phone = EXCLUDED.phone, address = EXCLUDED.address,
		 synced_at = EXCLUDED.synced_at, updated_at = now()
		 RETURNING id, created_at, updated_at`,
		v.ID, tenantID, v.VendorID, v.Name, v.LegalName, v.TradeName,
		v.PAN, v.GSTIN, v.TAN, v.State, v.StateCode, v.Category, v.RegistrationStatus, v.MSMERegistered,
		v.Email, v.Phone, v.Address, v.SyncedAt,
	).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert vendor snapshot: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) ListVendorSnapshots(ctx context.Context, tenantID uuid.UUID) ([]domain.VendorSnapshot, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, vendor_id, name, legal_name, trade_name,
		 pan, gstin, tan, state, state_code, category, registration_status, msme_registered,
		 email, phone, address, synced_at, created_at, updated_at
		 FROM vendor_snapshots ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list vendor snapshots: %w", err)
	}
	defer rows.Close()

	var vendors []domain.VendorSnapshot
	for rows.Next() {
		var v domain.VendorSnapshot
		if err := rows.Scan(&v.ID, &v.TenantID, &v.VendorID, &v.Name, &v.LegalName, &v.TradeName,
			&v.PAN, &v.GSTIN, &v.TAN, &v.State, &v.StateCode, &v.Category, &v.RegistrationStatus, &v.MSMERegistered,
			&v.Email, &v.Phone, &v.Address, &v.SyncedAt, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan vendor snapshot: %w", err)
		}
		vendors = append(vendors, v)
	}

	return vendors, tx.Commit(ctx)
}

func (s *Store) GetVendorSnapshot(ctx context.Context, tenantID uuid.UUID, vendorID string) (*domain.VendorSnapshot, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var v domain.VendorSnapshot
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, vendor_id, name, legal_name, trade_name,
		 pan, gstin, tan, state, state_code, category, registration_status, msme_registered,
		 email, phone, address, synced_at, created_at, updated_at
		 FROM vendor_snapshots WHERE vendor_id = $1`, vendorID,
	).Scan(&v.ID, &v.TenantID, &v.VendorID, &v.Name, &v.LegalName, &v.TradeName,
		&v.PAN, &v.GSTIN, &v.TAN, &v.State, &v.StateCode, &v.Category, &v.RegistrationStatus, &v.MSMERegistered,
		&v.Email, &v.Phone, &v.Address, &v.SyncedAt, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get vendor snapshot: %w", err)
	}

	return &v, tx.Commit(ctx)
}

func (s *Store) CreateComplianceScore(ctx context.Context, tenantID uuid.UUID, cs *domain.ComplianceScore) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	cs.TenantID = tenantID
	if cs.ID == uuid.Nil {
		cs.ID = uuid.New()
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO compliance_scores (id, tenant_id, vendor_id, vendor_snapshot_id,
		 total_score, category, risk_level,
		 filing_regularity_score, irn_compliance_score, mismatch_rate_score,
		 payment_behavior_score, document_hygiene_score,
		 filing_regularity_note, irn_compliance_note, mismatch_rate_note,
		 payment_behavior_note, document_hygiene_note, scored_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
		 RETURNING created_at`,
		cs.ID, tenantID, cs.VendorID, cs.VendorSnapshotID,
		cs.TotalScore, cs.Category, cs.RiskLevel,
		cs.FilingRegularityScore, cs.IRNComplianceScore, cs.MismatchRateScore,
		cs.PaymentBehaviorScore, cs.DocumentHygieneScore,
		cs.FilingRegularityNote, cs.IRNComplianceNote, cs.MismatchRateNote,
		cs.PaymentBehaviorNote, cs.DocumentHygieneNote, cs.ScoredAt,
	).Scan(&cs.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert compliance score: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) GetLatestScore(ctx context.Context, tenantID uuid.UUID, vendorID string) (*domain.ComplianceScore, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var cs domain.ComplianceScore
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, vendor_id, vendor_snapshot_id,
		 total_score, category, risk_level,
		 filing_regularity_score, irn_compliance_score, mismatch_rate_score,
		 payment_behavior_score, document_hygiene_score,
		 filing_regularity_note, irn_compliance_note, mismatch_rate_note,
		 payment_behavior_note, document_hygiene_note, scored_at, created_at
		 FROM compliance_scores WHERE vendor_id = $1
		 ORDER BY scored_at DESC LIMIT 1`, vendorID,
	).Scan(&cs.ID, &cs.TenantID, &cs.VendorID, &cs.VendorSnapshotID,
		&cs.TotalScore, &cs.Category, &cs.RiskLevel,
		&cs.FilingRegularityScore, &cs.IRNComplianceScore, &cs.MismatchRateScore,
		&cs.PaymentBehaviorScore, &cs.DocumentHygieneScore,
		&cs.FilingRegularityNote, &cs.IRNComplianceNote, &cs.MismatchRateNote,
		&cs.PaymentBehaviorNote, &cs.DocumentHygieneNote, &cs.ScoredAt, &cs.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get latest score: %w", err)
	}

	return &cs, tx.Commit(ctx)
}

func (s *Store) ListLatestScores(ctx context.Context, tenantID uuid.UUID) ([]domain.ComplianceScore, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT DISTINCT ON (vendor_id) id, tenant_id, vendor_id, vendor_snapshot_id,
		 total_score, category, risk_level,
		 filing_regularity_score, irn_compliance_score, mismatch_rate_score,
		 payment_behavior_score, document_hygiene_score,
		 filing_regularity_note, irn_compliance_note, mismatch_rate_note,
		 payment_behavior_note, document_hygiene_note, scored_at, created_at
		 FROM compliance_scores
		 ORDER BY vendor_id, scored_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list latest scores: %w", err)
	}
	defer rows.Close()

	var scores []domain.ComplianceScore
	for rows.Next() {
		var cs domain.ComplianceScore
		if err := rows.Scan(&cs.ID, &cs.TenantID, &cs.VendorID, &cs.VendorSnapshotID,
			&cs.TotalScore, &cs.Category, &cs.RiskLevel,
			&cs.FilingRegularityScore, &cs.IRNComplianceScore, &cs.MismatchRateScore,
			&cs.PaymentBehaviorScore, &cs.DocumentHygieneScore,
			&cs.FilingRegularityNote, &cs.IRNComplianceNote, &cs.MismatchRateNote,
			&cs.PaymentBehaviorNote, &cs.DocumentHygieneNote, &cs.ScoredAt, &cs.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan compliance score: %w", err)
		}
		scores = append(scores, cs)
	}

	return scores, tx.Commit(ctx)
}

func (s *Store) GetScoreSummary(ctx context.Context, tenantID uuid.UUID) (*domain.ScoreSummary, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var summary domain.ScoreSummary
	err = tx.QueryRow(ctx,
		`SELECT
		 COUNT(*),
		 COALESCE(SUM(CASE WHEN category = 'A' THEN 1 ELSE 0 END), 0),
		 COALESCE(SUM(CASE WHEN category = 'B' THEN 1 ELSE 0 END), 0),
		 COALESCE(SUM(CASE WHEN category = 'C' THEN 1 ELSE 0 END), 0),
		 COALESCE(SUM(CASE WHEN category = 'D' THEN 1 ELSE 0 END), 0),
		 COALESCE(AVG(total_score)::integer, 0)
		 FROM (
		   SELECT DISTINCT ON (vendor_id) category, total_score
		   FROM compliance_scores
		   ORDER BY vendor_id, scored_at DESC
		 ) latest`).Scan(&summary.Total, &summary.CatA, &summary.CatB, &summary.CatC, &summary.CatD, &summary.AvgScore)
	if err != nil {
		return nil, fmt.Errorf("get score summary: %w", err)
	}

	return &summary, tx.Commit(ctx)
}

func (s *Store) CreateSyncStatus(ctx context.Context, tenantID uuid.UUID, ss *domain.SyncStatus) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	ss.TenantID = tenantID
	if ss.ID == uuid.Nil {
		ss.ID = uuid.New()
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO sync_status (id, tenant_id, sync_type, status, vendor_count, scored_count, started_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING created_at`,
		ss.ID, tenantID, ss.SyncType, ss.Status, ss.VendorCount, ss.ScoredCount, ss.StartedAt,
	).Scan(&ss.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert sync status: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) UpdateSyncStatus(ctx context.Context, tenantID uuid.UUID, syncID uuid.UUID, status string, vendorCount, scoredCount int, errMsg string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	now := time.Now().UTC()
	_, err = tx.Exec(ctx,
		`UPDATE sync_status SET status = $1, vendor_count = $2, scored_count = $3,
		 error_message = $4, completed_at = $5 WHERE id = $6`,
		status, vendorCount, scoredCount, errMsg, now, syncID)
	if err != nil {
		return fmt.Errorf("update sync status: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) GetLatestSyncStatus(ctx context.Context, tenantID uuid.UUID) (*domain.SyncStatus, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var ss domain.SyncStatus
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, sync_type, status, vendor_count, scored_count,
		 started_at, completed_at, error_message, created_at
		 FROM sync_status ORDER BY started_at DESC LIMIT 1`,
	).Scan(&ss.ID, &ss.TenantID, &ss.SyncType, &ss.Status, &ss.VendorCount, &ss.ScoredCount,
		&ss.StartedAt, &ss.CompletedAt, &ss.ErrorMessage, &ss.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get latest sync status: %w", err)
	}

	return &ss, tx.Commit(ctx)
}
