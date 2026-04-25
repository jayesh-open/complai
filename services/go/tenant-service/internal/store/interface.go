package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/tenant-service/internal/domain"
)

// Repository defines the data-access contract for the tenant service.
type Repository interface {
	CreateTenant(ctx context.Context, t *domain.Tenant) error
	GetTenant(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error)
	ListTenants(ctx context.Context, tenantID uuid.UUID) ([]domain.Tenant, error)
	UpdateTenantKMSKey(ctx context.Context, tenantID uuid.UUID, kmsKeyARN string) error
	UpdateTenantStatus(ctx context.Context, tenantID uuid.UUID, status string) error
	CreatePAN(ctx context.Context, tenantID uuid.UUID, p *domain.TenantPAN) error
	CreateGSTIN(ctx context.Context, tenantID uuid.UUID, g *domain.TenantGSTIN) error
	CreateTAN(ctx context.Context, tenantID uuid.UUID, t *domain.TenantTAN) error
	GetHierarchy(ctx context.Context, tenantID uuid.UUID) (*domain.TenantHierarchy, error)
}
