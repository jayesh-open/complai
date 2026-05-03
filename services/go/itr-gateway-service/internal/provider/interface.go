package provider

import (
	"context"

	"github.com/complai/complai/services/go/itr-gateway-service/internal/domain"
)

type SandboxITRProvider interface {
	CheckPANAadhaarLink(ctx context.Context, req domain.PANAadhaarLinkRequest) (*domain.PANAadhaarLinkResponse, error)
	FetchAIS(ctx context.Context, req domain.AISRequest) (*domain.AISResponse, error)
	SubmitITR(ctx context.Context, req domain.ITRSubmitRequest) (*domain.ITRSubmitResponse, error)
	GenerateITRV(ctx context.Context, req domain.ITRVRequest) (*domain.ITRVResponse, error)
	CheckEVerification(ctx context.Context, req domain.EVerifyRequest) (*domain.EVerifyResponse, error)
	CheckRefundStatus(ctx context.Context, req domain.RefundStatusRequest) (*domain.RefundStatusResponse, error)
}
