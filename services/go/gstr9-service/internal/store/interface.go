package store

import (
	"context"

	"github.com/complai/complai/services/go/gstr9-service/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

	CreateGSTR9CFiling(ctx context.Context, tenantID uuid.UUID, f *domain.GSTR9CFiling) error
	GetGSTR9CFiling(ctx context.Context, tenantID, id uuid.UUID) (*domain.GSTR9CFiling, error)
	GetGSTR9CFilingByGSTR9ID(ctx context.Context, tenantID, gstr9FilingID uuid.UUID) (*domain.GSTR9CFiling, error)
	UpdateGSTR9CStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.GSTR9CStatus) error
	UpdateGSTR9CUnreconciled(ctx context.Context, tenantID, id uuid.UUID, amount decimal.Decimal) error
	CertifyGSTR9C(ctx context.Context, tenantID, id, certifiedBy uuid.UUID) error

	CreateMismatch(ctx context.Context, tenantID uuid.UUID, m *domain.GSTR9CMismatch) error
	ListMismatches(ctx context.Context, tenantID, gstr9cFilingID uuid.UUID) ([]domain.GSTR9CMismatch, error)
	GetMismatch(ctx context.Context, tenantID, id uuid.UUID) (*domain.GSTR9CMismatch, error)
	ResolveMismatch(ctx context.Context, tenantID, id uuid.UUID, reason string, resolvedBy uuid.UUID) error
	DeleteMismatches(ctx context.Context, tenantID, gstr9cFilingID uuid.UUID) error
}
