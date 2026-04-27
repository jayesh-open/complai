package provider

import (
	"context"

	"github.com/complai/complai/services/go/ewb-gateway-service/internal/domain"
)

type EWBProvider interface {
	GenerateEWB(ctx context.Context, req *domain.GenerateEWBRequest) (*domain.GenerateEWBResponse, error)
	CancelEWB(ctx context.Context, req *domain.CancelEWBRequest) (*domain.CancelEWBResponse, error)
	GetEWB(ctx context.Context, ewbNo string) (*domain.GetEWBResponse, error)
	UpdateVehicle(ctx context.Context, req *domain.UpdateVehicleRequest) (*domain.UpdateVehicleResponse, error)
	ExtendValidity(ctx context.Context, req *domain.ExtendValidityRequest) (*domain.ExtendValidityResponse, error)
	ConsolidateEWB(ctx context.Context, req *domain.ConsolidateEWBRequest) (*domain.ConsolidateEWBResponse, error)
}
