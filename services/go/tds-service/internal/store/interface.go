package store

import (
	"context"

	"github.com/complai/complai/services/go/tds-service/internal/domain"
	"github.com/google/uuid"
)

type Repository interface {
	CreateDeductee(ctx context.Context, tenantID uuid.UUID, d *domain.Deductee) error
	GetDeductee(ctx context.Context, tenantID, id uuid.UUID) (*domain.Deductee, error)
	GetDeducteeByVendor(ctx context.Context, tenantID, vendorID uuid.UUID) (*domain.Deductee, error)
	ListDeductees(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.Deductee, int, error)
	UpsertDeductee(ctx context.Context, tenantID uuid.UUID, d *domain.Deductee) error

	CreateEntry(ctx context.Context, tenantID uuid.UUID, e *domain.TDSEntry) error
	GetEntry(ctx context.Context, tenantID, id uuid.UUID) (*domain.TDSEntry, error)
	ListEntries(ctx context.Context, tenantID uuid.UUID, fy, quarter string, limit, offset int) ([]domain.TDSEntry, int, error)

	GetAggregate(ctx context.Context, tenantID, deducteeID uuid.UUID, code domain.PaymentCode, fy string) (*domain.TDSAggregate, error)
	UpsertAggregate(ctx context.Context, tenantID uuid.UUID, agg *domain.TDSAggregate) error

	GetSummary(ctx context.Context, tenantID uuid.UUID, fy string) (*domain.TDSSummary, error)

	CreateFiling(ctx context.Context, tenantID uuid.UUID, f *domain.Filing) error
	GetFiling(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Filing, error)
	GetFilingByIdempotencyKey(ctx context.Context, tenantID uuid.UUID, key string) (*domain.Filing, error)
	UpdateFilingStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status domain.FilingStatus, tokenNumber, ackNumber, errMsg string) error
	ListFilings(ctx context.Context, tenantID uuid.UUID, fy, quarter string, limit, offset int) ([]domain.Filing, int, error)
}
