package provider

import (
	"context"

	"github.com/complai/complai/services/go/tds-gateway-service/internal/domain"
)

type SandboxTDSProvider interface {
	VerifyPAN(ctx context.Context, req domain.PANVerifyRequest) (*domain.PANVerifyResponse, error)
	VerifyTAN(ctx context.Context, req domain.TANVerifyRequest) (*domain.TANVerifyResponse, error)
	GenerateChallan(ctx context.Context, req domain.ChallanRequest) (*domain.ChallanResponse, error)
	FileForm140(ctx context.Context, req domain.Form140Request) (*domain.FormFilingResponse, error)
	FileForm138(ctx context.Context, req domain.Form138Request) (*domain.FormFilingResponse, error)
	FileForm144(ctx context.Context, req domain.Form144Request) (*domain.FormFilingResponse, error)
}
