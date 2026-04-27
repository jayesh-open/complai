package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/einvoice-service/internal/domain"
)

type Repository interface {
	CreateEInvoice(ctx context.Context, tenantID uuid.UUID, inv *domain.EInvoice) error
	GetEInvoice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EInvoice, error)
	GetEInvoiceByIRN(ctx context.Context, tenantID uuid.UUID, irn string) (*domain.EInvoice, error)
	GetEInvoiceByInvoiceNumber(ctx context.Context, tenantID uuid.UUID, gstin, invoiceNumber string) (*domain.EInvoice, error)
	ListEInvoices(ctx context.Context, tenantID uuid.UUID, req *domain.ListEInvoicesRequest) ([]domain.EInvoice, int, error)

	UpdateIRNGenerated(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, irn, ackNo, signedInvoice, signedQR string) error
	UpdateIRNCancelled(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, reason string) error
	UpdateIRNFailed(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, reason string) error

	CreateLineItems(ctx context.Context, tenantID uuid.UUID, items []domain.EInvoiceLineItem) error
	GetLineItems(ctx context.Context, tenantID uuid.UUID, invoiceID uuid.UUID) ([]domain.EInvoiceLineItem, error)

	GetSummary(ctx context.Context, tenantID uuid.UUID, gstin string) (*domain.EInvoiceSummary, error)
}
