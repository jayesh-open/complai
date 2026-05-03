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
	GSTR1Summary(ctx context.Context, req *domain.GSTR1SummaryRequest) (*domain.GSTR1SummaryResponse, error)
	GSTR2BGet(ctx context.Context, req *domain.GSTR2BGetRequest) (*domain.GSTR2BGetResponse, error)
	GSTR2AGet(ctx context.Context, req *domain.GSTR2AGetRequest) (*domain.GSTR2AGetResponse, error)
	IMSGet(ctx context.Context, req *domain.IMSGetRequest) (*domain.IMSGetResponse, error)
	IMSAction(ctx context.Context, req *domain.IMSActionRequest) (*domain.IMSActionResponse, error)
	IMSBulkAction(ctx context.Context, req *domain.IMSBulkActionRequest) (*domain.IMSBulkActionResponse, error)
	GSTR3BSave(ctx context.Context, req *domain.GSTR3BSaveRequest) (*domain.GSTR3BSaveResponse, error)
	GSTR3BSubmit(ctx context.Context, req *domain.GSTR3BSubmitRequest) (*domain.GSTR3BSubmitResponse, error)
	GSTR3BFile(ctx context.Context, req *domain.GSTR3BFileRequest) (*domain.GSTR3BFileResponse, error)

	GSTR9Save(ctx context.Context, req *domain.GSTR9SaveRequest) (*domain.GSTR9SaveResponse, error)
	GSTR9Submit(ctx context.Context, req *domain.GSTR9SubmitRequest) (*domain.GSTR9SubmitResponse, error)
	GSTR9File(ctx context.Context, req *domain.GSTR9FileRequest) (*domain.GSTR9FileResponse, error)
	GSTR9Status(ctx context.Context, req *domain.GSTR9StatusRequest) (*domain.GSTR9StatusResponse, error)
	GSTR9CSave(ctx context.Context, req *domain.GSTR9CSaveRequest) (*domain.GSTR9CSaveResponse, error)
	GSTR9CFile(ctx context.Context, req *domain.GSTR9CFileRequest) (*domain.GSTR9CFileResponse, error)
	GSTR9CStatus(ctx context.Context, req *domain.GSTR9CStatusRequest) (*domain.GSTR9CStatusResponse, error)
}
