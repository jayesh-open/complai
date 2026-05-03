package provider

import (
	"context"
	"errors"

	"github.com/complai/complai/services/go/gstn-gateway-service/internal/domain"
)

var ErrNotImplemented = errors.New("adaequare provider not implemented; awaiting GSP credentials")

var _ GSTNProvider = (*AdaequareProvider)(nil)

// AdaequareProvider will connect to Adaequare GSP endpoints for real GSTN operations.
// Wire-up deferred to Part 13 when Adaequare auth credentials are resolved.
//
// Adaequare GSTR-9 endpoint reference:
//   POST /enriched/returns/gstr9 (save)
//   POST /enriched/returns/gstr9/submit
//   POST /enriched/returns/gstr9/file
//   GET  /enriched/returns/gstr9/status/{reference}
//
// Adaequare GSTR-9C endpoint reference:
//   POST /enriched/returns/gstr9c (save)
//   POST /enriched/returns/gstr9c/file (DSC mandatory)
//   GET  /enriched/returns/gstr9c/status/{reference}
type AdaequareProvider struct {
	baseURL string
}

func NewAdaequareProvider(baseURL string) *AdaequareProvider {
	return &AdaequareProvider{baseURL: baseURL}
}

func (a *AdaequareProvider) Authenticate(_ context.Context) (*domain.AuthResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR1Save(_ context.Context, _ *domain.GSTR1SaveRequest) (*domain.GSTR1SaveResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR1Get(_ context.Context, _ *domain.GSTR1GetRequest) (*domain.GSTR1GetResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR1Reset(_ context.Context, _ *domain.GSTR1ResetRequest) (*domain.GSTR1ResetResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR1Submit(_ context.Context, _ *domain.GSTR1SubmitRequest) (*domain.GSTR1SubmitResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR1File(_ context.Context, _ *domain.GSTR1FileRequest) (*domain.GSTR1FileResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR1Status(_ context.Context, _ *domain.GSTR1StatusRequest) (*domain.GSTR1StatusResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR1Summary(_ context.Context, _ *domain.GSTR1SummaryRequest) (*domain.GSTR1SummaryResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR2BGet(_ context.Context, _ *domain.GSTR2BGetRequest) (*domain.GSTR2BGetResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR2AGet(_ context.Context, _ *domain.GSTR2AGetRequest) (*domain.GSTR2AGetResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) IMSGet(_ context.Context, _ *domain.IMSGetRequest) (*domain.IMSGetResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) IMSAction(_ context.Context, _ *domain.IMSActionRequest) (*domain.IMSActionResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) IMSBulkAction(_ context.Context, _ *domain.IMSBulkActionRequest) (*domain.IMSBulkActionResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR3BSave(_ context.Context, _ *domain.GSTR3BSaveRequest) (*domain.GSTR3BSaveResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR3BSubmit(_ context.Context, _ *domain.GSTR3BSubmitRequest) (*domain.GSTR3BSubmitResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR3BFile(_ context.Context, _ *domain.GSTR3BFileRequest) (*domain.GSTR3BFileResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR9Save(_ context.Context, _ *domain.GSTR9SaveRequest) (*domain.GSTR9SaveResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR9Submit(_ context.Context, _ *domain.GSTR9SubmitRequest) (*domain.GSTR9SubmitResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR9File(_ context.Context, _ *domain.GSTR9FileRequest) (*domain.GSTR9FileResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR9Status(_ context.Context, _ *domain.GSTR9StatusRequest) (*domain.GSTR9StatusResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR9CSave(_ context.Context, _ *domain.GSTR9CSaveRequest) (*domain.GSTR9CSaveResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR9CFile(_ context.Context, _ *domain.GSTR9CFileRequest) (*domain.GSTR9CFileResponse, error) {
	return nil, ErrNotImplemented
}

func (a *AdaequareProvider) GSTR9CStatus(_ context.Context, _ *domain.GSTR9CStatusRequest) (*domain.GSTR9CStatusResponse, error) {
	return nil, ErrNotImplemented
}
