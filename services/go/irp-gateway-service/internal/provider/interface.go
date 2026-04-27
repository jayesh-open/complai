package provider

import (
	"context"

	"github.com/complai/complai/services/go/irp-gateway-service/internal/domain"
)

type IRPProvider interface {
	Authenticate(ctx context.Context) (*domain.AuthResponse, error)
	GenerateIRN(ctx context.Context, req *domain.GenerateIRNRequest) (*domain.GenerateIRNResponse, error)
	CancelIRN(ctx context.Context, req *domain.CancelIRNRequest) (*domain.CancelIRNResponse, error)
	GetIRNByIRN(ctx context.Context, req *domain.GetIRNByIRNRequest) (*domain.GetIRNResponse, error)
	GetIRNByDoc(ctx context.Context, req *domain.GetIRNByDocRequest) (*domain.GetIRNResponse, error)
	ValidateGSTIN(ctx context.Context, req *domain.GSTINValidateRequest) (*domain.GSTINValidateResponse, error)
}
