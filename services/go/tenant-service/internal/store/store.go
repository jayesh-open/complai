package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/tenant-service/internal/domain"
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

func (s *Store) CreateTenant(ctx context.Context, t *domain.Tenant) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	t.ID = uuid.New()
	t.TenantID = t.ID

	if err := setTenantID(ctx, tx, t.TenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO tenants (id, tenant_id, name, slug, tier, status, kms_key_arn, settings)
		 VALUES ($1, $2, $3, $4, $5::tenancy_tier, $6::tenant_status, $7, $8)
		 RETURNING created_at, updated_at`,
		t.ID, t.TenantID, t.Name, t.Slug, t.Tier, "active", t.KMSKeyARN, t.Settings,
	).Scan(&t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert tenant: %w", err)
	}
	t.Status = "active"
	return tx.Commit(ctx)
}

func (s *Store) GetTenant(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var t domain.Tenant
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, name, slug, tier, status, kms_key_arn, settings, created_at, updated_at, deleted_at
		 FROM tenants WHERE id = $1`, tenantID,
	).Scan(&t.ID, &t.TenantID, &t.Name, &t.Slug, &t.Tier, &t.Status, &t.KMSKeyARN, &t.Settings, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("get tenant: %w", err)
	}
	return &t, tx.Commit(ctx)
}

func (s *Store) ListTenants(ctx context.Context, tenantID uuid.UUID) ([]domain.Tenant, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, name, slug, tier, status, kms_key_arn, settings, created_at, updated_at, deleted_at
		 FROM tenants WHERE deleted_at IS NULL ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list tenants: %w", err)
	}
	defer rows.Close()

	var tenants []domain.Tenant
	for rows.Next() {
		var t domain.Tenant
		if err := rows.Scan(&t.ID, &t.TenantID, &t.Name, &t.Slug, &t.Tier, &t.Status, &t.KMSKeyARN, &t.Settings, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt); err != nil {
			return nil, fmt.Errorf("scan tenant: %w", err)
		}
		tenants = append(tenants, t)
	}
	return tenants, tx.Commit(ctx)
}

func (s *Store) UpdateTenantKMSKey(ctx context.Context, tenantID uuid.UUID, kmsKeyARN string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE tenants SET kms_key_arn = $1, updated_at = now() WHERE id = $2`,
		kmsKeyARN, tenantID)
	if err != nil {
		return fmt.Errorf("update kms key: %w", err)
	}
	return tx.Commit(ctx)
}

func (s *Store) UpdateTenantStatus(ctx context.Context, tenantID uuid.UUID, status string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE tenants SET status = $1::tenant_status, updated_at = now() WHERE id = $2`,
		status, tenantID)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	return tx.Commit(ctx)
}

func (s *Store) CreatePAN(ctx context.Context, tenantID uuid.UUID, p *domain.TenantPAN) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO tenant_pans (tenant_id, pan, entity_name, pan_type)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`,
		tenantID, p.PAN, p.EntityName, p.PANType,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert PAN: %w", err)
	}
	p.TenantID = tenantID
	p.Status = "active"
	return tx.Commit(ctx)
}

func (s *Store) CreateGSTIN(ctx context.Context, tenantID uuid.UUID, g *domain.TenantGSTIN) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO tenant_gstins (tenant_id, pan_id, gstin, trade_name, state_code, registration_type)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`,
		tenantID, g.PANID, g.GSTIN, g.TradeName, g.StateCode, g.RegistrationType,
	).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert GSTIN: %w", err)
	}
	g.TenantID = tenantID
	g.Status = "active"
	return tx.Commit(ctx)
}

func (s *Store) CreateTAN(ctx context.Context, tenantID uuid.UUID, t *domain.TenantTAN) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO tenant_tans (tenant_id, pan_id, tan, deductor_name)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`,
		tenantID, t.PANID, t.TAN, t.DeductorName,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert TAN: %w", err)
	}
	t.TenantID = tenantID
	t.Status = "active"
	return tx.Commit(ctx)
}

func (s *Store) GetHierarchy(ctx context.Context, tenantID uuid.UUID) (*domain.TenantHierarchy, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var t domain.Tenant
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, name, slug, tier, status, kms_key_arn, settings, created_at, updated_at, deleted_at
		 FROM tenants WHERE id = $1`, tenantID,
	).Scan(&t.ID, &t.TenantID, &t.Name, &t.Slug, &t.Tier, &t.Status, &t.KMSKeyARN, &t.Settings, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("get tenant: %w", err)
	}

	panRows, err := tx.Query(ctx,
		`SELECT id, tenant_id, pan, entity_name, pan_type, status, created_at, updated_at
		 FROM tenant_pans ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list PANs: %w", err)
	}
	defer panRows.Close()

	var pans []domain.PANWithSub
	for panRows.Next() {
		var p domain.TenantPAN
		if err := panRows.Scan(&p.ID, &p.TenantID, &p.PAN, &p.EntityName, &p.PANType, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan PAN: %w", err)
		}
		pans = append(pans, domain.PANWithSub{TenantPAN: p})
	}

	for i := range pans {
		gstinRows, err := tx.Query(ctx,
			`SELECT id, tenant_id, pan_id, gstin, trade_name, state_code, registration_type, status, created_at, updated_at
			 FROM tenant_gstins WHERE pan_id = $1 ORDER BY created_at`, pans[i].ID)
		if err != nil {
			return nil, fmt.Errorf("list GSTINs: %w", err)
		}
		for gstinRows.Next() {
			var g domain.TenantGSTIN
			if err := gstinRows.Scan(&g.ID, &g.TenantID, &g.PANID, &g.GSTIN, &g.TradeName, &g.StateCode, &g.RegistrationType, &g.Status, &g.CreatedAt, &g.UpdatedAt); err != nil {
				gstinRows.Close()
				return nil, fmt.Errorf("scan GSTIN: %w", err)
			}
			pans[i].GSTINs = append(pans[i].GSTINs, g)
		}
		gstinRows.Close()

		tanRows, err := tx.Query(ctx,
			`SELECT id, tenant_id, pan_id, tan, deductor_name, status, created_at, updated_at
			 FROM tenant_tans WHERE pan_id = $1 ORDER BY created_at`, pans[i].ID)
		if err != nil {
			return nil, fmt.Errorf("list TANs: %w", err)
		}
		for tanRows.Next() {
			var tn domain.TenantTAN
			if err := tanRows.Scan(&tn.ID, &tn.TenantID, &tn.PANID, &tn.TAN, &tn.DeductorName, &tn.Status, &tn.CreatedAt, &tn.UpdatedAt); err != nil {
				tanRows.Close()
				return nil, fmt.Errorf("scan TAN: %w", err)
			}
			pans[i].TANs = append(pans[i].TANs, tn)
		}
		tanRows.Close()
	}

	return &domain.TenantHierarchy{Tenant: t, PANs: pans}, tx.Commit(ctx)
}
