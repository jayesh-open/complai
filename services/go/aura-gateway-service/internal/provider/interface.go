package provider

import (
	"context"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/aura-gateway-service/internal/domain"
)

type AuraProvider interface {
	ListARInvoices(ctx context.Context, tenantID uuid.UUID, gstin, period string) (*domain.InvoiceListResponse, error)
}
