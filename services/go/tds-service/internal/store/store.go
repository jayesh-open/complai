package store

import (
	"context"
	"fmt"

	"github.com/complai/complai/services/go/tds-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type PgStore struct {
	pool *pgxpool.Pool
}

func NewPgStore(pool *pgxpool.Pool) *PgStore {
	return &PgStore{pool: pool}
}

func (s *PgStore) setTenant(ctx context.Context, tenantID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID.String()))
	return err
}

func (s *PgStore) CreateDeductee(ctx context.Context, tenantID uuid.UUID, d *domain.Deductee) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO deductees
		(id, tenant_id, vendor_id, name, pan, pan_verified, pan_status, deductee_type, resident_status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		d.ID, tenantID, d.VendorID, d.Name, d.PAN, d.PANVerified, d.PANStatus,
		d.DeducteeType, d.ResidentStatus)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) GetDeductee(ctx context.Context, tenantID, id uuid.UUID) (*domain.Deductee, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, err
	}

	var d domain.Deductee
	err = tx.QueryRow(ctx, `SELECT id, tenant_id, vendor_id, name, pan, pan_verified, pan_status,
		deductee_type, resident_status, created_at, updated_at
		FROM deductees WHERE id = $1`, id).Scan(
		&d.ID, &d.TenantID, &d.VendorID, &d.Name, &d.PAN, &d.PANVerified, &d.PANStatus,
		&d.DeducteeType, &d.ResidentStatus, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)
	return &d, nil
}

func (s *PgStore) GetDeducteeByVendor(ctx context.Context, tenantID, vendorID uuid.UUID) (*domain.Deductee, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, err
	}

	var d domain.Deductee
	err = tx.QueryRow(ctx, `SELECT id, tenant_id, vendor_id, name, pan, pan_verified, pan_status,
		deductee_type, resident_status, created_at, updated_at
		FROM deductees WHERE vendor_id = $1`, vendorID).Scan(
		&d.ID, &d.TenantID, &d.VendorID, &d.Name, &d.PAN, &d.PANVerified, &d.PANStatus,
		&d.DeducteeType, &d.ResidentStatus, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)
	return &d, nil
}

func (s *PgStore) ListDeductees(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.Deductee, int, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, 0, err
	}

	var total int
	if err := tx.QueryRow(ctx, "SELECT COUNT(*) FROM deductees").Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := tx.Query(ctx, `SELECT id, tenant_id, vendor_id, name, pan, pan_verified, pan_status,
		deductee_type, resident_status, created_at, updated_at
		FROM deductees ORDER BY name LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []domain.Deductee
	for rows.Next() {
		var d domain.Deductee
		if err := rows.Scan(&d.ID, &d.TenantID, &d.VendorID, &d.Name, &d.PAN, &d.PANVerified,
			&d.PANStatus, &d.DeducteeType, &d.ResidentStatus, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, d)
	}
	tx.Commit(ctx)
	return out, total, nil
}

func (s *PgStore) UpsertDeductee(ctx context.Context, tenantID uuid.UUID, d *domain.Deductee) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO deductees
		(id, tenant_id, vendor_id, name, pan, pan_verified, pan_status, deductee_type, resident_status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (tenant_id, vendor_id) DO UPDATE SET
			name = EXCLUDED.name, pan = EXCLUDED.pan, pan_status = EXCLUDED.pan_status,
			deductee_type = EXCLUDED.deductee_type, updated_at = NOW()`,
		d.ID, tenantID, d.VendorID, d.Name, d.PAN, d.PANVerified, d.PANStatus,
		d.DeducteeType, d.ResidentStatus)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) CreateEntry(ctx context.Context, tenantID uuid.UUID, e *domain.TDSEntry) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO tds_entries
		(id, tenant_id, deductee_id, section, financial_year, quarter, transaction_date,
		 gross_amount, tds_rate, tds_amount, surcharge, cess, total_tax,
		 nature_of_payment, pan_at_deduction, no_pan_deduction, lower_cert_applied, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`,
		e.ID, tenantID, e.DeducteeID, e.Section, e.FinancialYear, e.Quarter, e.TransactionDate,
		e.GrossAmount, e.TDSRate, e.TDSAmount, e.Surcharge, e.Cess, e.TotalTax,
		e.NatureOfPayment, e.PANAtDeduction, e.NoPANDeduction, e.LowerCertApplied, e.Status)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) GetEntry(ctx context.Context, tenantID, id uuid.UUID) (*domain.TDSEntry, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, err
	}

	var e domain.TDSEntry
	err = tx.QueryRow(ctx, `SELECT id, tenant_id, deductee_id, section, financial_year, quarter,
		transaction_date, gross_amount, tds_rate, tds_amount, surcharge, cess, total_tax,
		nature_of_payment, pan_at_deduction, no_pan_deduction, lower_cert_applied, status,
		created_at, updated_at
		FROM tds_entries WHERE id = $1`, id).Scan(
		&e.ID, &e.TenantID, &e.DeducteeID, &e.Section, &e.FinancialYear, &e.Quarter,
		&e.TransactionDate, &e.GrossAmount, &e.TDSRate, &e.TDSAmount, &e.Surcharge, &e.Cess, &e.TotalTax,
		&e.NatureOfPayment, &e.PANAtDeduction, &e.NoPANDeduction, &e.LowerCertApplied, &e.Status,
		&e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)
	return &e, nil
}

func (s *PgStore) ListEntries(ctx context.Context, tenantID uuid.UUID, fy, quarter string, limit, offset int) ([]domain.TDSEntry, int, error) {
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
	if fy != "" {
		where += fmt.Sprintf(" AND financial_year = $%d", argN)
		args = append(args, fy)
		argN++
	}
	if quarter != "" {
		where += fmt.Sprintf(" AND quarter = $%d", argN)
		args = append(args, quarter)
		argN++
	}

	var total int
	if err := tx.QueryRow(ctx, "SELECT COUNT(*) FROM tds_entries "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	q := fmt.Sprintf(`SELECT id, tenant_id, deductee_id, section, financial_year, quarter,
		transaction_date, gross_amount, tds_rate, tds_amount, surcharge, cess, total_tax,
		nature_of_payment, pan_at_deduction, no_pan_deduction, lower_cert_applied, status,
		created_at, updated_at
		FROM tds_entries %s ORDER BY transaction_date DESC LIMIT $%d OFFSET $%d`, where, argN, argN+1)

	rows, err := tx.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []domain.TDSEntry
	for rows.Next() {
		var e domain.TDSEntry
		if err := rows.Scan(&e.ID, &e.TenantID, &e.DeducteeID, &e.Section, &e.FinancialYear, &e.Quarter,
			&e.TransactionDate, &e.GrossAmount, &e.TDSRate, &e.TDSAmount, &e.Surcharge, &e.Cess, &e.TotalTax,
			&e.NatureOfPayment, &e.PANAtDeduction, &e.NoPANDeduction, &e.LowerCertApplied, &e.Status,
			&e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, e)
	}
	tx.Commit(ctx)
	return out, total, nil
}

func (s *PgStore) GetAggregate(ctx context.Context, tenantID, deducteeID uuid.UUID, section domain.Section, fy string) (*domain.TDSAggregate, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, err
	}

	var a domain.TDSAggregate
	err = tx.QueryRow(ctx, `SELECT id, tenant_id, deductee_id, section, financial_year,
		total_paid, total_tds, transaction_count, updated_at
		FROM tds_aggregates WHERE deductee_id = $1 AND section = $2 AND financial_year = $3`,
		deducteeID, section, fy).Scan(
		&a.ID, &a.TenantID, &a.DeducteeID, &a.Section, &a.FinancialYear,
		&a.TotalPaid, &a.TotalTDS, &a.TransactionCount, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)
	return &a, nil
}

func (s *PgStore) UpsertAggregate(ctx context.Context, tenantID uuid.UUID, agg *domain.TDSAggregate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO tds_aggregates
		(id, tenant_id, deductee_id, section, financial_year, total_paid, total_tds, transaction_count)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (tenant_id, deductee_id, section, financial_year) DO UPDATE SET
			total_paid = EXCLUDED.total_paid, total_tds = EXCLUDED.total_tds,
			transaction_count = EXCLUDED.transaction_count, updated_at = NOW()`,
		agg.ID, tenantID, agg.DeducteeID, agg.Section, agg.FinancialYear,
		agg.TotalPaid, agg.TotalTDS, agg.TransactionCount)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PgStore) GetSummary(ctx context.Context, tenantID uuid.UUID, fy string) (*domain.TDSSummary, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); err != nil {
		return nil, err
	}

	sum := &domain.TDSSummary{
		EntriesBySection: make(map[domain.Section]int),
		EntriesByStatus:  make(map[domain.EntryStatus]int),
	}
	tx.QueryRow(ctx, "SELECT COUNT(*) FROM deductees").Scan(&sum.TotalDeductees)
	tx.QueryRow(ctx, "SELECT COUNT(*), COALESCE(SUM(total_tax),0) FROM tds_entries WHERE financial_year = $1", fy).Scan(&sum.TotalEntries, &sum.TotalTDSDeducted)

	deposited := decimal.Zero
	tx.QueryRow(ctx, "SELECT COALESCE(SUM(total_tax),0) FROM tds_entries WHERE financial_year = $1 AND status = 'DEPOSITED'", fy).Scan(&deposited)
	sum.TotalTDSDeposited = deposited
	sum.PendingDeposit = sum.TotalTDSDeducted.Sub(deposited)

	rows, _ := tx.Query(ctx, "SELECT section, COUNT(*) FROM tds_entries WHERE financial_year = $1 GROUP BY section", fy)
	if rows != nil {
		for rows.Next() {
			var sec domain.Section
			var cnt int
			rows.Scan(&sec, &cnt)
			sum.EntriesBySection[sec] = cnt
		}
		rows.Close()
	}

	rows2, _ := tx.Query(ctx, "SELECT status, COUNT(*) FROM tds_entries WHERE financial_year = $1 GROUP BY status", fy)
	if rows2 != nil {
		for rows2.Next() {
			var st domain.EntryStatus
			var cnt int
			rows2.Scan(&st, &cnt)
			sum.EntriesByStatus[st] = cnt
		}
		rows2.Close()
	}

	tx.Commit(ctx)
	return sum, nil
}
