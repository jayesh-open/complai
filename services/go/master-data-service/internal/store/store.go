package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/master-data-service/internal/domain"
)

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

// ---------------------------------------------------------------------------
// Vendors
// ---------------------------------------------------------------------------

func (s *Store) CreateVendor(ctx context.Context, tenantID uuid.UUID, v *domain.Vendor) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	v.ID = uuid.New()
	v.TenantID = tenantID

	err = tx.QueryRow(ctx,
		`INSERT INTO vendors (id, tenant_id, name, pan, gstin, email, phone,
		 address_line1, address_line2, city, state_code, pincode,
		 bank_name, bank_account, bank_ifsc, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		 RETURNING kyc_status, compliance_score, status, created_at, updated_at`,
		v.ID, v.TenantID, v.Name, v.PAN, v.GSTIN, v.Email, v.Phone,
		v.AddressLine1, v.AddressLine2, v.City, v.StateCode, v.Pincode,
		v.BankName, v.BankAccount, v.BankIFSC, v.Metadata,
	).Scan(&v.KYCStatus, &v.ComplianceScore, &v.Status, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert vendor: %w", err)
	}
	return tx.Commit(ctx)
}

func (s *Store) GetVendor(ctx context.Context, tenantID uuid.UUID, vendorID uuid.UUID) (*domain.Vendor, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var v domain.Vendor
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, name, pan, gstin, email, phone,
		 address_line1, address_line2, city, state_code, pincode,
		 bank_name, bank_account, bank_ifsc,
		 kyc_status, compliance_score, status, metadata, created_by, created_at, updated_at
		 FROM vendors WHERE id = $1`, vendorID,
	).Scan(&v.ID, &v.TenantID, &v.Name, &v.PAN, &v.GSTIN, &v.Email, &v.Phone,
		&v.AddressLine1, &v.AddressLine2, &v.City, &v.StateCode, &v.Pincode,
		&v.BankName, &v.BankAccount, &v.BankIFSC,
		&v.KYCStatus, &v.ComplianceScore, &v.Status, &v.Metadata, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get vendor: %w", err)
	}
	return &v, tx.Commit(ctx)
}

func (s *Store) ListVendors(ctx context.Context, tenantID uuid.UUID) ([]domain.Vendor, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, name, pan, gstin, email, phone,
		 address_line1, address_line2, city, state_code, pincode,
		 bank_name, bank_account, bank_ifsc,
		 kyc_status, compliance_score, status, metadata, created_by, created_at, updated_at
		 FROM vendors ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list vendors: %w", err)
	}
	defer rows.Close()

	var vendors []domain.Vendor
	for rows.Next() {
		var v domain.Vendor
		if err := rows.Scan(&v.ID, &v.TenantID, &v.Name, &v.PAN, &v.GSTIN, &v.Email, &v.Phone,
			&v.AddressLine1, &v.AddressLine2, &v.City, &v.StateCode, &v.Pincode,
			&v.BankName, &v.BankAccount, &v.BankIFSC,
			&v.KYCStatus, &v.ComplianceScore, &v.Status, &v.Metadata, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan vendor: %w", err)
		}
		vendors = append(vendors, v)
	}
	return vendors, tx.Commit(ctx)
}

func (s *Store) UpdateVendor(ctx context.Context, tenantID uuid.UUID, vendorID uuid.UUID, req *domain.UpdateVendorRequest) (*domain.Vendor, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var v domain.Vendor
	err = tx.QueryRow(ctx,
		`UPDATE vendors SET
		 name = COALESCE($1, name),
		 pan = COALESCE($2, pan),
		 gstin = COALESCE($3, gstin),
		 email = COALESCE($4, email),
		 phone = COALESCE($5, phone),
		 address_line1 = COALESCE($6, address_line1),
		 address_line2 = COALESCE($7, address_line2),
		 city = COALESCE($8, city),
		 state_code = COALESCE($9, state_code),
		 pincode = COALESCE($10, pincode),
		 bank_name = COALESCE($11, bank_name),
		 bank_account = COALESCE($12, bank_account),
		 bank_ifsc = COALESCE($13, bank_ifsc),
		 status = COALESCE($14, status),
		 updated_at = now()
		 WHERE id = $15
		 RETURNING id, tenant_id, name, pan, gstin, email, phone,
		 address_line1, address_line2, city, state_code, pincode,
		 bank_name, bank_account, bank_ifsc,
		 kyc_status, compliance_score, status, metadata, created_by, created_at, updated_at`,
		req.Name, req.PAN, req.GSTIN, req.Email, req.Phone,
		req.AddressLine1, req.AddressLine2, req.City, req.StateCode, req.Pincode,
		req.BankName, req.BankAccount, req.BankIFSC, req.Status, vendorID,
	).Scan(&v.ID, &v.TenantID, &v.Name, &v.PAN, &v.GSTIN, &v.Email, &v.Phone,
		&v.AddressLine1, &v.AddressLine2, &v.City, &v.StateCode, &v.Pincode,
		&v.BankName, &v.BankAccount, &v.BankIFSC,
		&v.KYCStatus, &v.ComplianceScore, &v.Status, &v.Metadata, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("update vendor: %w", err)
	}
	return &v, tx.Commit(ctx)
}

// ---------------------------------------------------------------------------
// Customers
// ---------------------------------------------------------------------------

func (s *Store) CreateCustomer(ctx context.Context, tenantID uuid.UUID, c *domain.Customer) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	c.ID = uuid.New()
	c.TenantID = tenantID

	err = tx.QueryRow(ctx,
		`INSERT INTO customers (id, tenant_id, name, pan, gstin, email, phone,
		 address_line1, city, state_code, pincode, payment_terms_days, credit_limit, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		 RETURNING status, created_at, updated_at`,
		c.ID, c.TenantID, c.Name, c.PAN, c.GSTIN, c.Email, c.Phone,
		c.AddressLine1, c.City, c.StateCode, c.Pincode, c.PaymentTermsDays, c.CreditLimit, c.Metadata,
	).Scan(&c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert customer: %w", err)
	}
	return tx.Commit(ctx)
}

func (s *Store) GetCustomer(ctx context.Context, tenantID uuid.UUID, customerID uuid.UUID) (*domain.Customer, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var c domain.Customer
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, name, pan, gstin, email, phone,
		 address_line1, city, state_code, pincode, payment_terms_days, credit_limit,
		 status, metadata, created_at, updated_at
		 FROM customers WHERE id = $1`, customerID,
	).Scan(&c.ID, &c.TenantID, &c.Name, &c.PAN, &c.GSTIN, &c.Email, &c.Phone,
		&c.AddressLine1, &c.City, &c.StateCode, &c.Pincode, &c.PaymentTermsDays, &c.CreditLimit,
		&c.Status, &c.Metadata, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}
	return &c, tx.Commit(ctx)
}

func (s *Store) ListCustomers(ctx context.Context, tenantID uuid.UUID) ([]domain.Customer, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, name, pan, gstin, email, phone,
		 address_line1, city, state_code, pincode, payment_terms_days, credit_limit,
		 status, metadata, created_at, updated_at
		 FROM customers ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	defer rows.Close()

	var customers []domain.Customer
	for rows.Next() {
		var c domain.Customer
		if err := rows.Scan(&c.ID, &c.TenantID, &c.Name, &c.PAN, &c.GSTIN, &c.Email, &c.Phone,
			&c.AddressLine1, &c.City, &c.StateCode, &c.Pincode, &c.PaymentTermsDays, &c.CreditLimit,
			&c.Status, &c.Metadata, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan customer: %w", err)
		}
		customers = append(customers, c)
	}
	return customers, tx.Commit(ctx)
}

// ---------------------------------------------------------------------------
// Items
// ---------------------------------------------------------------------------

func (s *Store) CreateItem(ctx context.Context, tenantID uuid.UUID, i *domain.Item) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	i.ID = uuid.New()
	i.TenantID = tenantID

	err = tx.QueryRow(ctx,
		`INSERT INTO items (id, tenant_id, name, description, hsn_code, unit_of_measure,
		 unit_price, gst_rate, is_service, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING status, created_at, updated_at`,
		i.ID, i.TenantID, i.Name, i.Description, i.HSNCode, i.UnitOfMeasure,
		i.UnitPrice, i.GSTRate, i.IsService, i.Metadata,
	).Scan(&i.Status, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert item: %w", err)
	}
	return tx.Commit(ctx)
}

func (s *Store) GetItem(ctx context.Context, tenantID uuid.UUID, itemID uuid.UUID) (*domain.Item, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var i domain.Item
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, name, description, hsn_code, unit_of_measure,
		 unit_price, gst_rate, is_service, status, metadata, created_at, updated_at
		 FROM items WHERE id = $1`, itemID,
	).Scan(&i.ID, &i.TenantID, &i.Name, &i.Description, &i.HSNCode, &i.UnitOfMeasure,
		&i.UnitPrice, &i.GSTRate, &i.IsService, &i.Status, &i.Metadata, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}
	return &i, tx.Commit(ctx)
}

func (s *Store) ListItems(ctx context.Context, tenantID uuid.UUID) ([]domain.Item, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, name, description, hsn_code, unit_of_measure,
		 unit_price, gst_rate, is_service, status, metadata, created_at, updated_at
		 FROM items ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var i domain.Item
		if err := rows.Scan(&i.ID, &i.TenantID, &i.Name, &i.Description, &i.HSNCode, &i.UnitOfMeasure,
			&i.UnitPrice, &i.GSTRate, &i.IsService, &i.Status, &i.Metadata, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, i)
	}
	return items, tx.Commit(ctx)
}

// ---------------------------------------------------------------------------
// HSN Codes
// ---------------------------------------------------------------------------

func (s *Store) ListHSNCodes(ctx context.Context, tenantID uuid.UUID) ([]domain.HSNCode, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, code, description, gst_rate, effective_from, effective_to, created_at, updated_at
		 FROM hsn_codes ORDER BY code`)
	if err != nil {
		return nil, fmt.Errorf("list hsn codes: %w", err)
	}
	defer rows.Close()

	var codes []domain.HSNCode
	for rows.Next() {
		var h domain.HSNCode
		if err := rows.Scan(&h.ID, &h.TenantID, &h.Code, &h.Description, &h.GSTRate,
			&h.EffectiveFrom, &h.EffectiveTo, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan hsn code: %w", err)
		}
		codes = append(codes, h)
	}
	return codes, tx.Commit(ctx)
}

func (s *Store) GetHSNCode(ctx context.Context, tenantID uuid.UUID, hsnID uuid.UUID) (*domain.HSNCode, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var h domain.HSNCode
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, code, description, gst_rate, effective_from, effective_to, created_at, updated_at
		 FROM hsn_codes WHERE id = $1`, hsnID,
	).Scan(&h.ID, &h.TenantID, &h.Code, &h.Description, &h.GSTRate,
		&h.EffectiveFrom, &h.EffectiveTo, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get hsn code: %w", err)
	}
	return &h, tx.Commit(ctx)
}

func (s *Store) CreateHSNCode(ctx context.Context, tenantID uuid.UUID, h *domain.HSNCode) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	h.ID = uuid.New()
	h.TenantID = tenantID

	effectiveFrom := h.EffectiveFrom
	if effectiveFrom == "" {
		effectiveFrom = "now()"
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO hsn_codes (id, tenant_id, code, description, gst_rate, effective_from)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING effective_from, created_at, updated_at`,
		h.ID, h.TenantID, h.Code, h.Description, h.GSTRate, effectiveFrom,
	).Scan(&h.EffectiveFrom, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert hsn code: %w", err)
	}
	return tx.Commit(ctx)
}

// ---------------------------------------------------------------------------
// State Codes
// ---------------------------------------------------------------------------

func (s *Store) ListStateCodes(ctx context.Context, tenantID uuid.UUID) ([]domain.StateCode, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, code, name, tin_code FROM state_codes ORDER BY code`)
	if err != nil {
		return nil, fmt.Errorf("list state codes: %w", err)
	}
	defer rows.Close()

	var codes []domain.StateCode
	for rows.Next() {
		var sc domain.StateCode
		if err := rows.Scan(&sc.ID, &sc.TenantID, &sc.Code, &sc.Name, &sc.TINCode); err != nil {
			return nil, fmt.Errorf("scan state code: %w", err)
		}
		codes = append(codes, sc)
	}
	return codes, tx.Commit(ctx)
}
