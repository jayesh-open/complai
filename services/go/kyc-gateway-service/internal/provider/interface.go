package provider

import (
	"context"

	"github.com/complai/complai/services/go/kyc-gateway-service/internal/domain"
)

type KYCProvider interface {
	VerifyPAN(ctx context.Context, req *domain.PANVerifyRequest) (*domain.PANVerifyResponse, error)
	VerifyGSTIN(ctx context.Context, req *domain.GSTINVerifyRequest) (*domain.GSTINVerifyResponse, error)
	VerifyTAN(ctx context.Context, req *domain.TANVerifyRequest) (*domain.TANVerifyResponse, error)
	VerifyBank(ctx context.Context, req *domain.BankVerifyRequest) (*domain.BankVerifyResponse, error)
}
