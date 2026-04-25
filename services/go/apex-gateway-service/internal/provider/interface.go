package provider

import (
	"context"

	"github.com/complai/complai/services/go/apex-gateway-service/internal/domain"
)

// ApexProvider defines the interface for fetching vendor and AP invoice data from Apex P2P.
type ApexProvider interface {
	FetchVendors(ctx context.Context, req *domain.FetchVendorsRequest) (*domain.FetchVendorsResponse, error)
	FetchAPInvoices(ctx context.Context, req *domain.FetchAPInvoicesRequest) (*domain.FetchAPInvoicesResponse, error)
}
