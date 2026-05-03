package store

import (
	"context"

	"github.com/complai/complai/services/go/gstr9-service/internal/domain"
	"github.com/google/uuid"
)

type Repository interface {
	CreateFiling(ctx context.Context, tenantID uuid.UUID, f *domain.GSTR9Filing) error
	GetFiling(ctx context.Context, tenantID, id uuid.UUID) (*domain.GSTR9Filing, error)
	UpdateFilingStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.FilingStatus) error
	ListFilings(ctx context.Context, tenantID uuid.UUID, gstin, fy string, limit, offset int) ([]domain.GSTR9Filing, int, error)

	CreateTableData(ctx context.Context, tenantID uuid.UUID, td *domain.GSTR9TableData) error
	ListTableData(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.GSTR9TableData, error)
	DeleteTableData(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) error

	CreateAuditLog(ctx context.Context, tenantID uuid.UUID, log *domain.GSTR9AuditLog) error
	ListAuditLogs(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.GSTR9AuditLog, error)
}
