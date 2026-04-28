package provider

import (
	"context"

	"github.com/complai/complai/services/go/tds-gateway-service/internal/domain"
)

type SandboxTDSProvider interface {
	VerifyPAN(ctx context.Context, req domain.PANVerifyRequest) (*domain.PANVerifyResponse, error)
	VerifyTAN(ctx context.Context, req domain.TANVerifyRequest) (*domain.TANVerifyResponse, error)
	GenerateChallan(ctx context.Context, req domain.ChallanRequest) (*domain.ChallanResponse, error)
	FileForm26Q(ctx context.Context, req domain.Form26QRequest) (*domain.FormFilingResponse, error)
	FileForm24Q(ctx context.Context, req domain.Form24QRequest) (*domain.FormFilingResponse, error)
}
