package store

import (
	"context"

	"github.com/complai/complai/services/go/itr-service/internal/domain"
	"github.com/google/uuid"
)

type Repository interface {
	CreateTaxpayer(ctx context.Context, tenantID uuid.UUID, t *domain.Taxpayer) error
	GetTaxpayer(ctx context.Context, tenantID, id uuid.UUID) (*domain.Taxpayer, error)
	GetTaxpayerByPAN(ctx context.Context, tenantID uuid.UUID, pan string) (*domain.Taxpayer, error)
	ListTaxpayers(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.Taxpayer, int, error)

	CreateFiling(ctx context.Context, tenantID uuid.UUID, f *domain.ITRFiling) error
	GetFiling(ctx context.Context, tenantID, id uuid.UUID) (*domain.ITRFiling, error)
	GetFilingByIdempotencyKey(ctx context.Context, tenantID uuid.UUID, key string) (*domain.ITRFiling, error)
	UpdateFilingStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.FilingStatus, arn, ackNumber, errMsg string) error
	ListFilings(ctx context.Context, tenantID uuid.UUID, taxYear string, limit, offset int) ([]domain.ITRFiling, int, error)

	CreateIncomeEntry(ctx context.Context, tenantID uuid.UUID, e *domain.IncomeEntry) error
	ListIncomeEntries(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.IncomeEntry, error)

	CreateDeduction(ctx context.Context, tenantID uuid.UUID, d *domain.Deduction) error
	ListDeductions(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.Deduction, error)

	SaveTaxComputation(ctx context.Context, tenantID uuid.UUID, tc *domain.TaxComputation) error
	GetTaxComputation(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) (*domain.TaxComputation, error)

	CreateTDSCredit(ctx context.Context, tenantID uuid.UUID, c *domain.TDSCredit) error
	ListTDSCredits(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.TDSCredit, error)

	CreateAISReconciliation(ctx context.Context, tenantID uuid.UUID, r *domain.AISReconciliation) error
	ListAISReconciliations(ctx context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.AISReconciliation, error)
}
