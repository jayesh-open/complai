package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/gst-service/internal/domain"
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

func (s *Store) CreateFiling(ctx context.Context, tenantID uuid.UUID, f *domain.GSTR1Filing) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	f.ID = uuid.New()
	f.TenantID = tenantID
	f.RequestID = uuid.New()

	err = tx.QueryRow(ctx,
		`INSERT INTO gstr1_filings (id, tenant_id, gstin, return_period, status, total_count, error_count, request_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING created_at, updated_at`,
		f.ID, f.TenantID, f.GSTIN, f.ReturnPeriod, f.Status, f.TotalCount, f.ErrorCount, f.RequestID,
	).Scan(&f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert filing: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) GetFiling(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) (*domain.GSTR1Filing, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var f domain.GSTR1Filing
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, gstin, return_period, status, total_count, error_count,
		        arn, filed_at, filed_by, approved_by, approved_at, request_id, created_at, updated_at
		 FROM gstr1_filings WHERE id = $1`, filingID,
	).Scan(&f.ID, &f.TenantID, &f.GSTIN, &f.ReturnPeriod, &f.Status, &f.TotalCount, &f.ErrorCount,
		&f.ARN, &f.FiledAt, &f.FiledBy, &f.ApprovedBy, &f.ApprovedAt, &f.RequestID, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get filing: %w", err)
	}

	return &f, tx.Commit(ctx)
}

func (s *Store) GetFilingByPeriod(ctx context.Context, tenantID uuid.UUID, gstin, period string) (*domain.GSTR1Filing, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var f domain.GSTR1Filing
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, gstin, return_period, status, total_count, error_count,
		        arn, filed_at, filed_by, approved_by, approved_at, request_id, created_at, updated_at
		 FROM gstr1_filings WHERE gstin = $1 AND return_period = $2
		 ORDER BY created_at DESC LIMIT 1`, gstin, period,
	).Scan(&f.ID, &f.TenantID, &f.GSTIN, &f.ReturnPeriod, &f.Status, &f.TotalCount, &f.ErrorCount,
		&f.ARN, &f.FiledAt, &f.FiledBy, &f.ApprovedBy, &f.ApprovedAt, &f.RequestID, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get filing by period: %w", err)
	}

	return &f, tx.Commit(ctx)
}

func (s *Store) UpdateFilingStatus(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, status domain.FilingStatus) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE gstr1_filings SET status = $1, updated_at = now() WHERE id = $2`,
		status, filingID)
	if err != nil {
		return fmt.Errorf("update filing status: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) UpdateFilingARN(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, arn string, filedBy uuid.UUID) error {
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
		`UPDATE gstr1_filings SET status = $1, arn = $2, filed_at = $3, filed_by = $4, updated_at = now()
		 WHERE id = $5`,
		domain.FilingStatusFiled, arn, now, filedBy, filingID)
	if err != nil {
		return fmt.Errorf("update filing ARN: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) ApproveFiling(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, approvedBy uuid.UUID) error {
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
		`UPDATE gstr1_filings SET status = $1, approved_by = $2, approved_at = $3, updated_at = now()
		 WHERE id = $4`,
		domain.FilingStatusApproved, approvedBy, now, filingID)
	if err != nil {
		return fmt.Errorf("approve filing: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) BulkInsertEntries(ctx context.Context, tenantID uuid.UUID, entries []domain.SalesRegisterEntry) (int, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return 0, fmt.Errorf("set tenant: %w", err)
	}

	inserted := 0
	for _, e := range entries {
		_, err := tx.Exec(ctx,
			`INSERT INTO sales_register (id, tenant_id, gstin, return_period, document_number, document_date,
			 document_type, supply_type, reverse_charge, supplier_gstin, buyer_gstin, buyer_name, buyer_state,
			 place_of_supply, hsn, taxable_value, cgst_rate, cgst_amount, sgst_rate, sgst_amount,
			 igst_rate, igst_amount, grand_total, source_system, source_id, section)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26)
			 ON CONFLICT (tenant_id, gstin, document_number) DO NOTHING`,
			e.ID, tenantID, e.GSTIN, e.ReturnPeriod, e.DocumentNumber, e.DocumentDate,
			e.DocumentType, e.SupplyType, e.ReverseCharge, e.SupplierGSTIN, e.BuyerGSTIN, e.BuyerName, e.BuyerState,
			e.PlaceOfSupply, e.HSN, e.TaxableValue, e.CGSTRate, e.CGSTAmount, e.SGSTRate, e.SGSTAmount,
			e.IGSTRate, e.IGSTAmount, e.GrandTotal, e.SourceSystem, e.SourceID, e.Section,
		)
		if err != nil {
			return 0, fmt.Errorf("insert entry: %w", err)
		}
		inserted++
	}

	return inserted, tx.Commit(ctx)
}

func (s *Store) ListEntries(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, section string) ([]domain.SalesRegisterEntry, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var f domain.GSTR1Filing
	err = tx.QueryRow(ctx, `SELECT gstin, return_period FROM gstr1_filings WHERE id = $1`, filingID).
		Scan(&f.GSTIN, &f.ReturnPeriod)
	if err != nil {
		return nil, fmt.Errorf("get filing for entries: %w", err)
	}

	query := `SELECT id, tenant_id, gstin, return_period, document_number, document_date,
	          document_type, supply_type, reverse_charge, supplier_gstin, buyer_gstin, buyer_name, buyer_state,
	          place_of_supply, hsn, taxable_value, cgst_rate, cgst_amount, sgst_rate, sgst_amount,
	          igst_rate, igst_amount, grand_total, source_system, source_id, section, created_at, updated_at
	          FROM sales_register WHERE gstin = $1 AND return_period = $2`
	args := []interface{}{f.GSTIN, f.ReturnPeriod}

	if section != "" {
		query += " AND section = $3"
		args = append(args, section)
	}
	query += " ORDER BY document_date, document_number"

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list entries: %w", err)
	}
	defer rows.Close()

	var entries []domain.SalesRegisterEntry
	for rows.Next() {
		var e domain.SalesRegisterEntry
		if err := rows.Scan(&e.ID, &e.TenantID, &e.GSTIN, &e.ReturnPeriod, &e.DocumentNumber, &e.DocumentDate,
			&e.DocumentType, &e.SupplyType, &e.ReverseCharge, &e.SupplierGSTIN, &e.BuyerGSTIN, &e.BuyerName, &e.BuyerState,
			&e.PlaceOfSupply, &e.HSN, &e.TaxableValue, &e.CGSTRate, &e.CGSTAmount, &e.SGSTRate, &e.SGSTAmount,
			&e.IGSTRate, &e.IGSTAmount, &e.GrandTotal, &e.SourceSystem, &e.SourceID, &e.Section,
			&e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan entry: %w", err)
		}
		entries = append(entries, e)
	}

	return entries, tx.Commit(ctx)
}

func (s *Store) CountEntries(ctx context.Context, tenantID uuid.UUID, gstin, period string) (int, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return 0, fmt.Errorf("set tenant: %w", err)
	}

	var count int
	err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM sales_register WHERE gstin = $1 AND return_period = $2`,
		gstin, period,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count entries: %w", err)
	}

	return count, tx.Commit(ctx)
}

func (s *Store) CreateSections(ctx context.Context, tenantID uuid.UUID, sections []domain.GSTR1Section) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	for _, sec := range sections {
		_, err := tx.Exec(ctx,
			`INSERT INTO gstr1_sections (id, tenant_id, filing_id, section, invoice_count, taxable_value, cgst, sgst, igst, total_tax, status)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			 ON CONFLICT (tenant_id, filing_id, section) DO UPDATE SET
			 invoice_count = EXCLUDED.invoice_count, taxable_value = EXCLUDED.taxable_value,
			 cgst = EXCLUDED.cgst, sgst = EXCLUDED.sgst, igst = EXCLUDED.igst, total_tax = EXCLUDED.total_tax`,
			sec.ID, tenantID, sec.FilingID, sec.Section, sec.InvoiceCount, sec.TaxableValue,
			sec.CGST, sec.SGST, sec.IGST, sec.TotalTax, sec.Status,
		)
		if err != nil {
			return fmt.Errorf("insert section: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (s *Store) ListSections(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.GSTR1Section, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, filing_id, section, invoice_count, taxable_value, cgst, sgst, igst, total_tax, status, created_at
		 FROM gstr1_sections WHERE filing_id = $1 ORDER BY section`, filingID)
	if err != nil {
		return nil, fmt.Errorf("list sections: %w", err)
	}
	defer rows.Close()

	var sections []domain.GSTR1Section
	for rows.Next() {
		var s domain.GSTR1Section
		if err := rows.Scan(&s.ID, &s.TenantID, &s.FilingID, &s.Section, &s.InvoiceCount, &s.TaxableValue,
			&s.CGST, &s.SGST, &s.IGST, &s.TotalTax, &s.Status, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan section: %w", err)
		}
		sections = append(sections, s)
	}

	return sections, tx.Commit(ctx)
}

func (s *Store) CreateValidationErrors(ctx context.Context, tenantID uuid.UUID, errs []domain.ValidationError) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	for _, ve := range errs {
		_, err := tx.Exec(ctx,
			`INSERT INTO validation_errors (id, tenant_id, filing_id, entry_id, field, code, message, severity)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			ve.ID, tenantID, ve.FilingID, ve.EntryID, ve.Field, ve.Code, ve.Message, ve.Severity,
		)
		if err != nil {
			return fmt.Errorf("insert validation error: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (s *Store) ListValidationErrors(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.ValidationError, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, filing_id, entry_id, field, code, message, severity, created_at
		 FROM validation_errors WHERE filing_id = $1 ORDER BY created_at`, filingID)
	if err != nil {
		return nil, fmt.Errorf("list validation errors: %w", err)
	}
	defer rows.Close()

	var errs []domain.ValidationError
	for rows.Next() {
		var ve domain.ValidationError
		if err := rows.Scan(&ve.ID, &ve.TenantID, &ve.FilingID, &ve.EntryID, &ve.Field, &ve.Code, &ve.Message,
			&ve.Severity, &ve.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan validation error: %w", err)
		}
		errs = append(errs, ve)
	}

	return errs, tx.Commit(ctx)
}

func (s *Store) CountValidationErrors(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) (int, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return 0, fmt.Errorf("set tenant: %w", err)
	}

	var count int
	err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM validation_errors WHERE filing_id = $1`, filingID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count validation errors: %w", err)
	}

	return count, tx.Commit(ctx)
}

// GSTR-3B store methods

func (s *Store) CreateGSTR3BFiling(ctx context.Context, tenantID uuid.UUID, f *domain.GSTR3BFiling) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	f.ID = uuid.New()
	f.TenantID = tenantID
	f.RequestID = uuid.New()

	err = tx.QueryRow(ctx,
		`INSERT INTO gstr3b_filings (id, tenant_id, gstin, return_period, status, data_json, request_id, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (tenant_id, gstin, return_period) DO UPDATE SET
		 status = EXCLUDED.status, data_json = EXCLUDED.data_json, updated_at = now()
		 RETURNING id, created_at, updated_at`,
		f.ID, f.TenantID, f.GSTIN, f.ReturnPeriod, f.Status, f.DataJSON, f.RequestID, f.CreatedBy,
	).Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert gstr3b filing: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) GetGSTR3BFiling(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) (*domain.GSTR3BFiling, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var f domain.GSTR3BFiling
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, gstin, return_period, status, data_json,
		        arn, filed_at, filed_by, approved_by, approved_at, request_id, created_at, updated_at
		 FROM gstr3b_filings WHERE id = $1`, filingID,
	).Scan(&f.ID, &f.TenantID, &f.GSTIN, &f.ReturnPeriod, &f.Status, &f.DataJSON,
		&f.ARN, &f.FiledAt, &f.FiledBy, &f.ApprovedBy, &f.ApprovedAt, &f.RequestID, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get gstr3b filing: %w", err)
	}

	return &f, tx.Commit(ctx)
}

func (s *Store) GetGSTR3BFilingByPeriod(ctx context.Context, tenantID uuid.UUID, gstin, period string) (*domain.GSTR3BFiling, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var f domain.GSTR3BFiling
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, gstin, return_period, status, data_json,
		        arn, filed_at, filed_by, approved_by, approved_at, request_id, created_at, updated_at
		 FROM gstr3b_filings WHERE gstin = $1 AND return_period = $2
		 ORDER BY created_at DESC LIMIT 1`, gstin, period,
	).Scan(&f.ID, &f.TenantID, &f.GSTIN, &f.ReturnPeriod, &f.Status, &f.DataJSON,
		&f.ARN, &f.FiledAt, &f.FiledBy, &f.ApprovedBy, &f.ApprovedAt, &f.RequestID, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get gstr3b filing by period: %w", err)
	}

	return &f, tx.Commit(ctx)
}

func (s *Store) UpdateGSTR3BStatus(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, status domain.GSTR3BStatus) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE gstr3b_filings SET status = $1, updated_at = now() WHERE id = $2`,
		status, filingID)
	if err != nil {
		return fmt.Errorf("update gstr3b status: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) UpdateGSTR3BData(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, dataJSON string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE gstr3b_filings SET data_json = $1, updated_at = now() WHERE id = $2`,
		dataJSON, filingID)
	if err != nil {
		return fmt.Errorf("update gstr3b data: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) ApproveGSTR3BFiling(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, approvedBy uuid.UUID) error {
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
		`UPDATE gstr3b_filings SET status = $1, approved_by = $2, approved_at = $3, updated_at = now()
		 WHERE id = $4`,
		domain.GSTR3BStatusApproved, approvedBy, now, filingID)
	if err != nil {
		return fmt.Errorf("approve gstr3b filing: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) UpdateGSTR3BARN(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, arn string, filedBy uuid.UUID) error {
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
		`UPDATE gstr3b_filings SET status = $1, arn = $2, filed_at = $3, filed_by = $4, updated_at = now()
		 WHERE id = $5`,
		domain.GSTR3BStatusFiled, arn, now, filedBy, filingID)
	if err != nil {
		return fmt.Errorf("update gstr3b ARN: %w", err)
	}

	return tx.Commit(ctx)
}
