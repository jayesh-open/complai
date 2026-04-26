package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/gst-service/internal/domain"
)

type Repository interface {
	CreateFiling(ctx context.Context, tenantID uuid.UUID, f *domain.GSTR1Filing) error
	GetFiling(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) (*domain.GSTR1Filing, error)
	GetFilingByPeriod(ctx context.Context, tenantID uuid.UUID, gstin, period string) (*domain.GSTR1Filing, error)
	UpdateFilingStatus(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, status domain.FilingStatus) error
	UpdateFilingARN(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, arn string, filedBy uuid.UUID) error
	ApproveFiling(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, approvedBy uuid.UUID) error

	BulkInsertEntries(ctx context.Context, tenantID uuid.UUID, entries []domain.SalesRegisterEntry) (int, error)
	ListEntries(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, section string) ([]domain.SalesRegisterEntry, error)
	CountEntries(ctx context.Context, tenantID uuid.UUID, gstin, period string) (int, error)

	CreateSections(ctx context.Context, tenantID uuid.UUID, sections []domain.GSTR1Section) error
	ListSections(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.GSTR1Section, error)

	CreateValidationErrors(ctx context.Context, tenantID uuid.UUID, errs []domain.ValidationError) error
	ListValidationErrors(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.ValidationError, error)
	CountValidationErrors(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) (int, error)

	// GSTR-3B
	CreateGSTR3BFiling(ctx context.Context, tenantID uuid.UUID, f *domain.GSTR3BFiling) error
	GetGSTR3BFiling(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) (*domain.GSTR3BFiling, error)
	GetGSTR3BFilingByPeriod(ctx context.Context, tenantID uuid.UUID, gstin, period string) (*domain.GSTR3BFiling, error)
	UpdateGSTR3BStatus(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, status domain.GSTR3BStatus) error
	UpdateGSTR3BData(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, dataJSON string) error
	ApproveGSTR3BFiling(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, approvedBy uuid.UUID) error
	UpdateGSTR3BARN(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID, arn string, filedBy uuid.UUID) error
}
