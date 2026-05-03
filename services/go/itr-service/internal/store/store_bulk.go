package store

import (
	"context"
	"fmt"

	"github.com/complai/complai/services/go/itr-service/internal/domain"
	"github.com/google/uuid"
)

func (s *PgStore) CreateBulkBatch(ctx context.Context, tenantID uuid.UUID, b *domain.BulkFilingBatch) error {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO bulk_filing_batches (id, tenant_id, tax_year, employer_tan, employer_name, total_employees, processed, ready, with_mismatches, status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		b.ID, b.TenantID, b.TaxYear, b.EmployerTAN, b.EmployerName, b.TotalEmployees, b.Processed, b.Ready, b.WithMismatches, b.Status)
	return err
}

func (s *PgStore) GetBulkBatch(ctx context.Context, tenantID, id uuid.UUID) (*domain.BulkFilingBatch, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, err
	}
	b := &domain.BulkFilingBatch{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, tenant_id, tax_year, employer_tan, employer_name, total_employees, processed, ready, with_mismatches, status, created_at, updated_at
		 FROM bulk_filing_batches WHERE id = $1`, id).Scan(
		&b.ID, &b.TenantID, &b.TaxYear, &b.EmployerTAN, &b.EmployerName, &b.TotalEmployees, &b.Processed, &b.Ready, &b.WithMismatches, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("batch not found: %w", err)
	}
	return b, nil
}

func (s *PgStore) UpdateBulkBatchStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.BulkBatchStatus, processed, ready, mismatches int) error {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx,
		`UPDATE bulk_filing_batches SET status=$1, processed=$2, ready=$3, with_mismatches=$4, updated_at=NOW() WHERE id=$5`,
		status, processed, ready, mismatches, id)
	return err
}

func (s *PgStore) ListBulkBatches(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.BulkFilingBatch, int, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, 0, err
	}
	var total int
	_ = s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM bulk_filing_batches").Scan(&total)

	pgRows, err := s.pool.Query(ctx,
		`SELECT id, tenant_id, tax_year, employer_tan, employer_name, total_employees, processed, ready, with_mismatches, status, created_at, updated_at
		 FROM bulk_filing_batches ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer pgRows.Close()

	var result []domain.BulkFilingBatch
	for pgRows.Next() {
		var b domain.BulkFilingBatch
		if err := pgRows.Scan(&b.ID, &b.TenantID, &b.TaxYear, &b.EmployerTAN, &b.EmployerName, &b.TotalEmployees, &b.Processed, &b.Ready, &b.WithMismatches, &b.Status, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, b)
	}
	return result, total, nil
}

func (s *PgStore) CreateBulkEmployee(ctx context.Context, tenantID uuid.UUID, e *domain.BulkFilingEmployee) error {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO bulk_filing_employees (id, tenant_id, batch_id, pan, name, email, gross_salary, tds_deducted, form_type, filing_id, status, mismatch_count, magic_link_token, token_expires_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		e.ID, e.TenantID, e.BatchID, e.PAN, e.Name, e.Email, e.GrossSalary, e.TDSDeducted, e.FormType, e.FilingID, e.Status, e.MismatchCount, e.MagicLinkToken, e.TokenExpiresAt)
	return err
}

func (s *PgStore) GetBulkEmployee(ctx context.Context, tenantID, id uuid.UUID) (*domain.BulkFilingEmployee, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, err
	}
	e := &domain.BulkFilingEmployee{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, tenant_id, batch_id, pan, name, email, gross_salary, tds_deducted, form_type, filing_id, status, mismatch_count, magic_link_token, token_expires_at, created_at, updated_at
		 FROM bulk_filing_employees WHERE id = $1`, id).Scan(
		&e.ID, &e.TenantID, &e.BatchID, &e.PAN, &e.Name, &e.Email, &e.GrossSalary, &e.TDSDeducted, &e.FormType, &e.FilingID, &e.Status, &e.MismatchCount, &e.MagicLinkToken, &e.TokenExpiresAt, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}
	return e, nil
}

func (s *PgStore) ListBulkEmployees(ctx context.Context, tenantID uuid.UUID, batchID uuid.UUID, limit, offset int) ([]domain.BulkFilingEmployee, int, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, 0, err
	}
	var total int
	_ = s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM bulk_filing_employees WHERE batch_id=$1", batchID).Scan(&total)

	pgRows, err := s.pool.Query(ctx,
		`SELECT id, tenant_id, batch_id, pan, name, email, gross_salary, tds_deducted, form_type, filing_id, status, mismatch_count, magic_link_token, token_expires_at, created_at, updated_at
		 FROM bulk_filing_employees WHERE batch_id=$1 ORDER BY created_at LIMIT $2 OFFSET $3`, batchID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer pgRows.Close()

	var result []domain.BulkFilingEmployee
	for pgRows.Next() {
		var e domain.BulkFilingEmployee
		if err := pgRows.Scan(&e.ID, &e.TenantID, &e.BatchID, &e.PAN, &e.Name, &e.Email, &e.GrossSalary, &e.TDSDeducted, &e.FormType, &e.FilingID, &e.Status, &e.MismatchCount, &e.MagicLinkToken, &e.TokenExpiresAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, e)
	}
	return result, total, nil
}

func (s *PgStore) UpdateBulkEmployeeStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.EmployeeFilingStatus) error {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx,
		`UPDATE bulk_filing_employees SET status=$1, updated_at=NOW() WHERE id=$2`,
		status, id)
	return err
}
