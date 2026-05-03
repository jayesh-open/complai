package store

import (
	"context"
	"fmt"

	"github.com/complai/complai/services/go/itr-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgStore struct {
	pool *pgxpool.Pool
}

func NewPgStore(pool *pgxpool.Pool) *PgStore {
	return &PgStore{pool: pool}
}

func (s *PgStore) setTenant(ctx context.Context, tenantID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", tenantID.String())
	return err
}

func (s *PgStore) CreateTaxpayer(ctx context.Context, tenantID uuid.UUID, t *domain.Taxpayer) error {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO taxpayers (id, tenant_id, pan, name, date_of_birth, assessee_type, residency_status, aadhaar_linked, email, mobile, address, employer_tan)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		t.ID, t.TenantID, t.PAN, t.Name, t.DateOfBirth, t.AssesseeType, t.ResidencyStatus, t.AadhaarLinked, t.Email, t.Mobile, t.Address, t.EmployerTAN)
	return err
}

func (s *PgStore) GetTaxpayer(ctx context.Context, tenantID, id uuid.UUID) (*domain.Taxpayer, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, err
	}
	t := &domain.Taxpayer{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, tenant_id, pan, name, date_of_birth, assessee_type, residency_status, aadhaar_linked, email, mobile, address, employer_tan, created_at, updated_at
		 FROM taxpayers WHERE id = $1`, id).Scan(
		&t.ID, &t.TenantID, &t.PAN, &t.Name, &t.DateOfBirth, &t.AssesseeType, &t.ResidencyStatus, &t.AadhaarLinked, &t.Email, &t.Mobile, &t.Address, &t.EmployerTAN, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("taxpayer not found: %w", err)
	}
	return t, nil
}

func (s *PgStore) GetTaxpayerByPAN(ctx context.Context, tenantID uuid.UUID, pan string) (*domain.Taxpayer, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, err
	}
	t := &domain.Taxpayer{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, tenant_id, pan, name, date_of_birth, assessee_type, residency_status, aadhaar_linked, email, mobile, address, employer_tan, created_at, updated_at
		 FROM taxpayers WHERE pan = $1`, pan).Scan(
		&t.ID, &t.TenantID, &t.PAN, &t.Name, &t.DateOfBirth, &t.AssesseeType, &t.ResidencyStatus, &t.AadhaarLinked, &t.Email, &t.Mobile, &t.Address, &t.EmployerTAN, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("taxpayer not found: %w", err)
	}
	return t, nil
}

func (s *PgStore) ListTaxpayers(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.Taxpayer, int, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, 0, err
	}
	var total int
	_ = s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM taxpayers").Scan(&total)

	rows, err := s.pool.Query(ctx,
		`SELECT id, tenant_id, pan, name, date_of_birth, assessee_type, residency_status, aadhaar_linked, email, mobile, address, employer_tan, created_at, updated_at
		 FROM taxpayers ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []domain.Taxpayer
	for rows.Next() {
		var t domain.Taxpayer
		if err := rows.Scan(&t.ID, &t.TenantID, &t.PAN, &t.Name, &t.DateOfBirth, &t.AssesseeType, &t.ResidencyStatus, &t.AadhaarLinked, &t.Email, &t.Mobile, &t.Address, &t.EmployerTAN, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, t)
	}
	return result, total, nil
}

func (s *PgStore) CreateFiling(ctx context.Context, tenantID uuid.UUID, f *domain.ITRFiling) error {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO itr_filings (id, tenant_id, taxpayer_id, pan, tax_year, form_type, regime_selected, form_10iea_ref, status,
		 gross_income, total_deductions, taxable_income, tax_payable, tds_credited, advance_tax_paid, self_assessment_tax, refund_due, balance_payable,
		 idempotency_key)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`,
		f.ID, f.TenantID, f.TaxpayerID, f.PAN, f.TaxYear, f.FormType, f.RegimeSelected, f.Form10IEARef, f.Status,
		f.GrossIncome, f.TotalDeductions, f.TaxableIncome, f.TaxPayable, f.TDSCredited, f.AdvanceTaxPaid, f.SelfAssessmentTax, f.RefundDue, f.BalancePayable,
		f.IdempotencyKey)
	return err
}

func (s *PgStore) GetFiling(ctx context.Context, tenantID, id uuid.UUID) (*domain.ITRFiling, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, err
	}
	f := &domain.ITRFiling{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, tenant_id, taxpayer_id, pan, tax_year, form_type, regime_selected, form_10iea_ref, status,
		 gross_income, total_deductions, taxable_income, tax_payable, tds_credited, advance_tax_paid, self_assessment_tax, refund_due, balance_payable,
		 verification_method, arn, acknowledgement_number, filed_at, idempotency_key, error_message, created_at, updated_at
		 FROM itr_filings WHERE id = $1`, id).Scan(
		&f.ID, &f.TenantID, &f.TaxpayerID, &f.PAN, &f.TaxYear, &f.FormType, &f.RegimeSelected, &f.Form10IEARef, &f.Status,
		&f.GrossIncome, &f.TotalDeductions, &f.TaxableIncome, &f.TaxPayable, &f.TDSCredited, &f.AdvanceTaxPaid, &f.SelfAssessmentTax, &f.RefundDue, &f.BalancePayable,
		&f.VerificationMethod, &f.ARN, &f.AcknowledgementNumber, &f.FiledAt, &f.IdempotencyKey, &f.ErrorMessage, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("filing not found: %w", err)
	}
	return f, nil
}

func (s *PgStore) GetFilingByIdempotencyKey(ctx context.Context, tenantID uuid.UUID, key string) (*domain.ITRFiling, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, err
	}
	f := &domain.ITRFiling{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, tenant_id, taxpayer_id, pan, tax_year, form_type, regime_selected, status, idempotency_key, created_at, updated_at
		 FROM itr_filings WHERE idempotency_key = $1`, key).Scan(
		&f.ID, &f.TenantID, &f.TaxpayerID, &f.PAN, &f.TaxYear, &f.FormType, &f.RegimeSelected, &f.Status, &f.IdempotencyKey, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("filing not found: %w", err)
	}
	return f, nil
}

func (s *PgStore) UpdateFilingStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.FilingStatus, arn, ackNumber, errMsg string) error {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx,
		`UPDATE itr_filings SET status=$1, arn=$2, acknowledgement_number=$3, error_message=$4, updated_at=NOW() WHERE id=$5`,
		status, arn, ackNumber, errMsg, id)
	return err
}

func (s *PgStore) ListFilings(ctx context.Context, tenantID uuid.UUID, taxYear string, limit, offset int) ([]domain.ITRFiling, int, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, 0, err
	}
	var total int
	if taxYear != "" {
		_ = s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM itr_filings WHERE tax_year=$1", taxYear).Scan(&total)
	} else {
		_ = s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM itr_filings").Scan(&total)
	}

	var pgRows pgx.Rows
	var err error
	if taxYear != "" {
		pgRows, err = s.pool.Query(ctx,
			`SELECT id, tenant_id, taxpayer_id, pan, tax_year, form_type, regime_selected, status, idempotency_key, created_at, updated_at
			 FROM itr_filings WHERE tax_year=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, taxYear, limit, offset)
	} else {
		pgRows, err = s.pool.Query(ctx,
			`SELECT id, tenant_id, taxpayer_id, pan, tax_year, form_type, regime_selected, status, idempotency_key, created_at, updated_at
			 FROM itr_filings ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer pgRows.Close()

	var result []domain.ITRFiling
	for pgRows.Next() {
		var f domain.ITRFiling
		if err := pgRows.Scan(&f.ID, &f.TenantID, &f.TaxpayerID, &f.PAN, &f.TaxYear, &f.FormType, &f.RegimeSelected, &f.Status, &f.IdempotencyKey, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, f)
	}
	return result, total, nil
}

func (s *PgStore) CreateIncomeEntry(ctx context.Context, tenantID uuid.UUID, e *domain.IncomeEntry) error {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO income_heads (id, tenant_id, filing_id, head, sub_head, section, description, amount, exempt)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		e.ID, e.TenantID, e.FilingID, e.Head, e.SubHead, e.Section, e.Description, e.Amount, e.Exempt)
	return err
}

func (s *PgStore) ListIncomeEntries(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.IncomeEntry, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, err
	}
	pgRows, err := s.pool.Query(ctx,
		`SELECT id, tenant_id, filing_id, head, sub_head, section, description, amount, exempt, created_at FROM income_heads WHERE filing_id=$1`, filingID)
	if err != nil {
		return nil, err
	}
	defer pgRows.Close()
	var result []domain.IncomeEntry
	for pgRows.Next() {
		var e domain.IncomeEntry
		if err := pgRows.Scan(&e.ID, &e.TenantID, &e.FilingID, &e.Head, &e.SubHead, &e.Section, &e.Description, &e.Amount, &e.Exempt, &e.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (s *PgStore) CreateDeduction(ctx context.Context, tenantID uuid.UUID, d *domain.Deduction) error {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO deductions (id, tenant_id, filing_id, section, label, claimed, allowed, max_limit)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		d.ID, d.TenantID, d.FilingID, d.Section, d.Label, d.Claimed, d.Allowed, d.MaxLimit)
	return err
}

func (s *PgStore) ListDeductions(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.Deduction, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, err
	}
	pgRows, err := s.pool.Query(ctx,
		`SELECT id, tenant_id, filing_id, section, label, claimed, allowed, max_limit, created_at FROM deductions WHERE filing_id=$1`, filingID)
	if err != nil {
		return nil, err
	}
	defer pgRows.Close()
	var result []domain.Deduction
	for pgRows.Next() {
		var d domain.Deduction
		if err := pgRows.Scan(&d.ID, &d.TenantID, &d.FilingID, &d.Section, &d.Label, &d.Claimed, &d.Allowed, &d.MaxLimit, &d.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

func (s *PgStore) SaveTaxComputation(ctx context.Context, tenantID uuid.UUID, tc *domain.TaxComputation) error {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO tax_computations (id, tenant_id, filing_id, regime_type, gross_income, standard_deduction, total_deductions, taxable_income,
		 base_tax, surcharge, surcharge_rate, health_ed_cess, rebate_87a, gross_tax_payable, tds_credit, advance_tax, self_assessment_tax, net_tax_payable, refund_due)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
		 ON CONFLICT (filing_id) DO UPDATE SET
		 regime_type=$4, gross_income=$5, standard_deduction=$6, total_deductions=$7, taxable_income=$8,
		 base_tax=$9, surcharge=$10, surcharge_rate=$11, health_ed_cess=$12, rebate_87a=$13, gross_tax_payable=$14,
		 tds_credit=$15, advance_tax=$16, self_assessment_tax=$17, net_tax_payable=$18, refund_due=$19`,
		tc.ID, tc.TenantID, tc.FilingID, tc.RegimeType, tc.GrossIncome, tc.StandardDeduction, tc.TotalDeductions, tc.TaxableIncome,
		tc.BaseTax, tc.Surcharge, tc.SurchargeRate, tc.HealthEdCess, tc.Rebate87A, tc.GrossTaxPayable,
		tc.TDSCredit, tc.AdvanceTax, tc.SelfAssessmentTax, tc.NetTaxPayable, tc.RefundDue)
	return err
}

func (s *PgStore) GetTaxComputation(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) (*domain.TaxComputation, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, err
	}
	tc := &domain.TaxComputation{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, tenant_id, filing_id, regime_type, gross_income, standard_deduction, total_deductions, taxable_income,
		 base_tax, surcharge, surcharge_rate, health_ed_cess, rebate_87a, gross_tax_payable, tds_credit, advance_tax, self_assessment_tax, net_tax_payable, refund_due, created_at
		 FROM tax_computations WHERE filing_id=$1`, filingID).Scan(
		&tc.ID, &tc.TenantID, &tc.FilingID, &tc.RegimeType, &tc.GrossIncome, &tc.StandardDeduction, &tc.TotalDeductions, &tc.TaxableIncome,
		&tc.BaseTax, &tc.Surcharge, &tc.SurchargeRate, &tc.HealthEdCess, &tc.Rebate87A, &tc.GrossTaxPayable,
		&tc.TDSCredit, &tc.AdvanceTax, &tc.SelfAssessmentTax, &tc.NetTaxPayable, &tc.RefundDue, &tc.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("tax computation not found: %w", err)
	}
	return tc, nil
}

func (s *PgStore) CreateTDSCredit(ctx context.Context, tenantID uuid.UUID, c *domain.TDSCredit) error {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO tds_credits (id, tenant_id, filing_id, deductor_tan, deductor_name, section, tds_amount, gross_payment, tax_year, matched_with_ais, ais_amount, discrepancy, discrepancy_note)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		c.ID, c.TenantID, c.FilingID, c.DeductorTAN, c.DeductorName, c.Section, c.TDSAmount, c.GrossPayment, c.TaxYear, c.MatchedWithAIS, c.AISAmount, c.Discrepancy, c.DiscrepancyNote)
	return err
}

func (s *PgStore) ListTDSCredits(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.TDSCredit, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, err
	}
	pgRows, err := s.pool.Query(ctx,
		`SELECT id, tenant_id, filing_id, deductor_tan, deductor_name, section, tds_amount, gross_payment, tax_year, matched_with_ais, ais_amount, discrepancy, discrepancy_note, created_at
		 FROM tds_credits WHERE filing_id=$1`, filingID)
	if err != nil {
		return nil, err
	}
	defer pgRows.Close()
	var result []domain.TDSCredit
	for pgRows.Next() {
		var c domain.TDSCredit
		if err := pgRows.Scan(&c.ID, &c.TenantID, &c.FilingID, &c.DeductorTAN, &c.DeductorName, &c.Section, &c.TDSAmount, &c.GrossPayment, &c.TaxYear, &c.MatchedWithAIS, &c.AISAmount, &c.Discrepancy, &c.DiscrepancyNote, &c.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (s *PgStore) CreateAISReconciliation(ctx context.Context, tenantID uuid.UUID, r *domain.AISReconciliation) error {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO ais_reconciliations (id, tenant_id, filing_id, pan, tax_year, source_type, reported_amount, ais_amount, discrepancy, status, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		r.ID, r.TenantID, r.FilingID, r.PAN, r.TaxYear, r.SourceType, r.ReportedAmount, r.AISAmount, r.Discrepancy, r.Status, r.Notes)
	return err
}

func (s *PgStore) ListAISReconciliations(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.AISReconciliation, error) {
	if err := s.setTenant(ctx, tenantID); err != nil {
		return nil, err
	}
	pgRows, err := s.pool.Query(ctx,
		`SELECT id, tenant_id, filing_id, pan, tax_year, source_type, reported_amount, ais_amount, discrepancy, status, notes, created_at
		 FROM ais_reconciliations WHERE filing_id=$1`, filingID)
	if err != nil {
		return nil, err
	}
	defer pgRows.Close()
	var result []domain.AISReconciliation
	for pgRows.Next() {
		var r domain.AISReconciliation
		if err := pgRows.Scan(&r.ID, &r.TenantID, &r.FilingID, &r.PAN, &r.TaxYear, &r.SourceType, &r.ReportedAmount, &r.AISAmount, &r.Discrepancy, &r.Status, &r.Notes, &r.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}
