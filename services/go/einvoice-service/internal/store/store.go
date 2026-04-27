package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/complai/complai/services/go/einvoice-service/internal/domain"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func setTenantID(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID) error {
	_, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID.String()))
	return err
}

func (s *Store) CreateEInvoice(ctx context.Context, tenantID uuid.UUID, inv *domain.EInvoice) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	inv.ID = uuid.New()
	inv.TenantID = tenantID
	inv.RequestID = uuid.New()
	inv.Status = domain.IRNStatusPending

	err = tx.QueryRow(ctx, `
		INSERT INTO einvoices (
			id, tenant_id, irn, ack_no, invoice_number, invoice_date, invoice_type,
			supplier_gstin, supplier_name, buyer_gstin, buyer_name,
			supply_type, place_of_supply, reverse_charge,
			taxable_value, cgst_amount, sgst_amount, igst_amount, cess_amount, total_amount,
			status, request_id, source_system, source_id,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11,
			$12, $13, $14,
			$15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24,
			now(), now()
		) RETURNING created_at, updated_at`,
		inv.ID, inv.TenantID, inv.IRN, inv.AckNo, inv.InvoiceNumber, inv.InvoiceDate, inv.InvoiceType,
		inv.SupplierGSTIN, inv.SupplierName, inv.BuyerGSTIN, inv.BuyerName,
		inv.SupplyType, inv.PlaceOfSupply, inv.ReverseCharge,
		inv.TaxableValue, inv.CGSTAmount, inv.SGSTAmount, inv.IGSTAmount, inv.CessAmount, inv.TotalAmount,
		inv.Status, inv.RequestID, inv.SourceSystem, inv.SourceID,
	).Scan(&inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert einvoice: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) GetEInvoice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EInvoice, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	inv := &domain.EInvoice{}
	err = tx.QueryRow(ctx, `
		SELECT id, tenant_id, irn, ack_no, invoice_number, invoice_date, invoice_type,
			supplier_gstin, supplier_name, buyer_gstin, buyer_name,
			supply_type, place_of_supply, reverse_charge,
			taxable_value, cgst_amount, sgst_amount, igst_amount, cess_amount, total_amount,
			status, irn_generated_at, irn_cancelled_at, cancel_reason,
			signed_invoice, signed_qr_code,
			request_id, source_system, source_id, created_at, updated_at
		FROM einvoices WHERE id = $1`, id).Scan(
		&inv.ID, &inv.TenantID, &inv.IRN, &inv.AckNo, &inv.InvoiceNumber, &inv.InvoiceDate, &inv.InvoiceType,
		&inv.SupplierGSTIN, &inv.SupplierName, &inv.BuyerGSTIN, &inv.BuyerName,
		&inv.SupplyType, &inv.PlaceOfSupply, &inv.ReverseCharge,
		&inv.TaxableValue, &inv.CGSTAmount, &inv.SGSTAmount, &inv.IGSTAmount, &inv.CessAmount, &inv.TotalAmount,
		&inv.Status, &inv.IRNGeneratedAt, &inv.IRNCancelledAt, &inv.CancelReason,
		&inv.SignedInvoice, &inv.SignedQRCode,
		&inv.RequestID, &inv.SourceSystem, &inv.SourceID, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get einvoice: %w", err)
	}

	return inv, tx.Commit(ctx)
}

func (s *Store) GetEInvoiceByIRN(ctx context.Context, tenantID uuid.UUID, irn string) (*domain.EInvoice, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	inv := &domain.EInvoice{}
	err = tx.QueryRow(ctx, `
		SELECT id, tenant_id, irn, ack_no, invoice_number, invoice_date, invoice_type,
			supplier_gstin, supplier_name, buyer_gstin, buyer_name,
			supply_type, place_of_supply, reverse_charge,
			taxable_value, cgst_amount, sgst_amount, igst_amount, cess_amount, total_amount,
			status, irn_generated_at, irn_cancelled_at, cancel_reason,
			signed_invoice, signed_qr_code,
			request_id, source_system, source_id, created_at, updated_at
		FROM einvoices WHERE irn = $1`, irn).Scan(
		&inv.ID, &inv.TenantID, &inv.IRN, &inv.AckNo, &inv.InvoiceNumber, &inv.InvoiceDate, &inv.InvoiceType,
		&inv.SupplierGSTIN, &inv.SupplierName, &inv.BuyerGSTIN, &inv.BuyerName,
		&inv.SupplyType, &inv.PlaceOfSupply, &inv.ReverseCharge,
		&inv.TaxableValue, &inv.CGSTAmount, &inv.SGSTAmount, &inv.IGSTAmount, &inv.CessAmount, &inv.TotalAmount,
		&inv.Status, &inv.IRNGeneratedAt, &inv.IRNCancelledAt, &inv.CancelReason,
		&inv.SignedInvoice, &inv.SignedQRCode,
		&inv.RequestID, &inv.SourceSystem, &inv.SourceID, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get einvoice by irn: %w", err)
	}

	return inv, tx.Commit(ctx)
}

func (s *Store) GetEInvoiceByInvoiceNumber(ctx context.Context, tenantID uuid.UUID, gstin, invoiceNumber string) (*domain.EInvoice, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	inv := &domain.EInvoice{}
	err = tx.QueryRow(ctx, `
		SELECT id, tenant_id, irn, ack_no, invoice_number, invoice_date, invoice_type,
			supplier_gstin, supplier_name, buyer_gstin, buyer_name,
			supply_type, place_of_supply, reverse_charge,
			taxable_value, cgst_amount, sgst_amount, igst_amount, cess_amount, total_amount,
			status, irn_generated_at, irn_cancelled_at, cancel_reason,
			signed_invoice, signed_qr_code,
			request_id, source_system, source_id, created_at, updated_at
		FROM einvoices WHERE supplier_gstin = $1 AND invoice_number = $2`, gstin, invoiceNumber).Scan(
		&inv.ID, &inv.TenantID, &inv.IRN, &inv.AckNo, &inv.InvoiceNumber, &inv.InvoiceDate, &inv.InvoiceType,
		&inv.SupplierGSTIN, &inv.SupplierName, &inv.BuyerGSTIN, &inv.BuyerName,
		&inv.SupplyType, &inv.PlaceOfSupply, &inv.ReverseCharge,
		&inv.TaxableValue, &inv.CGSTAmount, &inv.SGSTAmount, &inv.IGSTAmount, &inv.CessAmount, &inv.TotalAmount,
		&inv.Status, &inv.IRNGeneratedAt, &inv.IRNCancelledAt, &inv.CancelReason,
		&inv.SignedInvoice, &inv.SignedQRCode,
		&inv.RequestID, &inv.SourceSystem, &inv.SourceID, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get einvoice by invoice number: %w", err)
	}

	return inv, tx.Commit(ctx)
}

func (s *Store) ListEInvoices(ctx context.Context, tenantID uuid.UUID, req *domain.ListEInvoicesRequest) ([]domain.EInvoice, int, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, 0, fmt.Errorf("set tenant: %w", err)
	}

	query := `SELECT id, tenant_id, irn, ack_no, invoice_number, invoice_date, invoice_type,
		supplier_gstin, supplier_name, buyer_gstin, buyer_name,
		supply_type, place_of_supply, reverse_charge,
		taxable_value, cgst_amount, sgst_amount, igst_amount, cess_amount, total_amount,
		status, irn_generated_at, irn_cancelled_at, cancel_reason,
		signed_invoice, signed_qr_code,
		request_id, source_system, source_id, created_at, updated_at
		FROM einvoices WHERE supplier_gstin = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := tx.Query(ctx, query, req.GSTIN, req.PageSize, req.PageOffset)
	if err != nil {
		return nil, 0, fmt.Errorf("list einvoices: %w", err)
	}
	defer rows.Close()

	var invoices []domain.EInvoice
	for rows.Next() {
		var inv domain.EInvoice
		err := rows.Scan(
			&inv.ID, &inv.TenantID, &inv.IRN, &inv.AckNo, &inv.InvoiceNumber, &inv.InvoiceDate, &inv.InvoiceType,
			&inv.SupplierGSTIN, &inv.SupplierName, &inv.BuyerGSTIN, &inv.BuyerName,
			&inv.SupplyType, &inv.PlaceOfSupply, &inv.ReverseCharge,
			&inv.TaxableValue, &inv.CGSTAmount, &inv.SGSTAmount, &inv.IGSTAmount, &inv.CessAmount, &inv.TotalAmount,
			&inv.Status, &inv.IRNGeneratedAt, &inv.IRNCancelledAt, &inv.CancelReason,
			&inv.SignedInvoice, &inv.SignedQRCode,
			&inv.RequestID, &inv.SourceSystem, &inv.SourceID, &inv.CreatedAt, &inv.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan einvoice: %w", err)
		}
		invoices = append(invoices, inv)
	}

	var totalCount int
	err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM einvoices WHERE supplier_gstin = $1`, req.GSTIN).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("count einvoices: %w", err)
	}

	return invoices, totalCount, tx.Commit(ctx)
}

func (s *Store) UpdateIRNGenerated(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, irn, ackNo, signedInvoice, signedQR string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE einvoices SET
			irn = $1, ack_no = $2, signed_invoice = $3, signed_qr_code = $4,
			status = $5, irn_generated_at = now(), updated_at = now()
		WHERE id = $6`,
		irn, ackNo, signedInvoice, signedQR, domain.IRNStatusGenerated, id)
	if err != nil {
		return fmt.Errorf("update irn generated: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) UpdateIRNCancelled(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, reason string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE einvoices SET
			status = $1, irn_cancelled_at = now(), cancel_reason = $2, updated_at = now()
		WHERE id = $3`,
		domain.IRNStatusCancelled, reason, id)
	if err != nil {
		return fmt.Errorf("update irn cancelled: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) UpdateIRNFailed(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, reason string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE einvoices SET status = $1, cancel_reason = $2, updated_at = now() WHERE id = $3`,
		domain.IRNStatusFailed, reason, id)
	if err != nil {
		return fmt.Errorf("update irn failed: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) CreateLineItems(ctx context.Context, tenantID uuid.UUID, items []domain.EInvoiceLineItem) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	for i := range items {
		items[i].ID = uuid.New()
		items[i].TenantID = tenantID
		items[i].LineNumber = i + 1
		_, err := tx.Exec(ctx, `
			INSERT INTO einvoice_line_items (
				id, invoice_id, tenant_id, line_number, description, hsn_code,
				quantity, unit, unit_price, discount, taxable_value,
				cgst_rate, cgst_amount, sgst_rate, sgst_amount,
				igst_rate, igst_amount, cess_rate, cess_amount,
				created_at
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,now())`,
			items[i].ID, items[i].InvoiceID, items[i].TenantID, items[i].LineNumber,
			items[i].Description, items[i].HSNCode,
			items[i].Quantity, items[i].Unit, items[i].UnitPrice, items[i].Discount, items[i].TaxableValue,
			items[i].CGSTRate, items[i].CGSTAmount, items[i].SGSTRate, items[i].SGSTAmount,
			items[i].IGSTRate, items[i].IGSTAmount, items[i].CessRate, items[i].CessAmount,
		)
		if err != nil {
			return fmt.Errorf("insert line item %d: %w", i, err)
		}
	}

	return tx.Commit(ctx)
}

func (s *Store) GetLineItems(ctx context.Context, tenantID uuid.UUID, invoiceID uuid.UUID) ([]domain.EInvoiceLineItem, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx, `
		SELECT id, invoice_id, tenant_id, line_number, description, hsn_code,
			quantity, unit, unit_price, discount, taxable_value,
			cgst_rate, cgst_amount, sgst_rate, sgst_amount,
			igst_rate, igst_amount, cess_rate, cess_amount, created_at
		FROM einvoice_line_items WHERE invoice_id = $1 ORDER BY line_number`, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("get line items: %w", err)
	}
	defer rows.Close()

	var items []domain.EInvoiceLineItem
	for rows.Next() {
		var item domain.EInvoiceLineItem
		err := rows.Scan(
			&item.ID, &item.InvoiceID, &item.TenantID, &item.LineNumber, &item.Description, &item.HSNCode,
			&item.Quantity, &item.Unit, &item.UnitPrice, &item.Discount, &item.TaxableValue,
			&item.CGSTRate, &item.CGSTAmount, &item.SGSTRate, &item.SGSTAmount,
			&item.IGSTRate, &item.IGSTAmount, &item.CessRate, &item.CessAmount, &item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan line item: %w", err)
		}
		items = append(items, item)
	}

	return items, tx.Commit(ctx)
}

func (s *Store) GetSummary(ctx context.Context, tenantID uuid.UUID, gstin string) (*domain.EInvoiceSummary, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	summary := &domain.EInvoiceSummary{}
	err = tx.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'GENERATED'),
			COUNT(*) FILTER (WHERE status = 'PENDING'),
			COUNT(*) FILTER (WHERE status = 'CANCELLED'),
			COUNT(*) FILTER (WHERE status = 'FAILED'),
			COALESCE(SUM(total_amount), 0)
		FROM einvoices WHERE supplier_gstin = $1`, gstin).Scan(
		&summary.TotalCount, &summary.GeneratedCount, &summary.PendingCount,
		&summary.CancelledCount, &summary.FailedCount, &summary.TotalValue,
	)
	if err != nil {
		return nil, fmt.Errorf("get summary: %w", err)
	}

	_ = tx.Commit(ctx)
	return summary, nil
}

// Ensure compile-time interface compliance.
var _ Repository = (*Store)(nil)

// Clock abstraction for testable time-based logic.
type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

// CancellationWindowOpen checks if an e-invoice is still within the 24h cancellation window.
func CancellationWindowOpen(generatedAt *time.Time, clock Clock) bool {
	if generatedAt == nil {
		return false
	}
	return clock.Now().Sub(*generatedAt) < 24*time.Hour
}

// ValidityDaysForDistance calculates EWB validity: 1 day per 200 km, min 1.
func ValidityDaysForDistance(distanceKm int) int {
	if distanceKm <= 0 {
		return 1
	}
	days := (distanceKm + 199) / 200
	if days < 1 {
		return 1
	}
	return days
}

// Ensure decimal.Decimal satisfies the scanner — this is a build-time type assertion only.
var _ decimal.Decimal
