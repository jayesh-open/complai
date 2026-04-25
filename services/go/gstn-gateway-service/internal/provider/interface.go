package provider

import (
	"context"

	"github.com/complai/complai/services/go/gstn-gateway-service/internal/domain"
)

type GSTNProvider interface {
	Authenticate(ctx context.Context) (*domain.AuthResponse, error)
	GSTR1Save(ctx context.Context, req *domain.GSTR1SaveRequest) (*domain.GSTR1SaveResponse, error)
	GSTR1Get(ctx context.Context, req *domain.GSTR1GetRequest) (*domain.GSTR1GetResponse, error)
	GSTR1Reset(ctx context.Context, req *domain.GSTR1ResetRequest) (*domain.GSTR1ResetResponse, error)
	GSTR1Submit(ctx context.Context, req *domain.GSTR1SubmitRequest) (*domain.GSTR1SubmitResponse, error)
	GSTR1File(ctx context.Context, req *domain.GSTR1FileRequest) (*domain.GSTR1FileResponse, error)
	GSTR1Status(ctx context.Context, req *domain.GSTR1StatusRequest) (*domain.GSTR1StatusResponse, error)
}
