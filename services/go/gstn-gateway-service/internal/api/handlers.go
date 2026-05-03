package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/gstn-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/gstn-gateway-service/internal/provider"
)

type Handlers struct {
	provider provider.GSTNProvider
}

func NewHandlers(p provider.GSTNProvider) *Handlers {
	return &Handlers{provider: p}
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "gstn-gateway-service"})
}

func (h *Handlers) Authenticate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	resp, err := h.provider.Authenticate(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("authenticate failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      uuid.New().String(),
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: "success",
		},
	})
}

func (h *Handlers) GSTR1Save(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR1SaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GSTR1Save(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr1 save failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR1Get(w http.ResponseWriter, r *http.Request) {
	tenantID, _, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR1GetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = uuid.New().String()
	}

	start := time.Now()
	resp, err := h.provider.GSTR1Get(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr1 get failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR1Reset(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR1ResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GSTR1Reset(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr1 reset failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR1Submit(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR1SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GSTR1Submit(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr1 submit failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR1File(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR1FileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GSTR1File(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr1 file failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR1Status(w http.ResponseWriter, r *http.Request) {
	tenantID, _, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR1StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = uuid.New().String()
	}

	start := time.Now()
	resp, err := h.provider.GSTR1Status(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr1 status failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR2BGet(w http.ResponseWriter, r *http.Request) {
	tenantID, _, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR2BGetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = uuid.New().String()
	}

	start := time.Now()
	resp, err := h.provider.GSTR2BGet(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr2b get failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR2AGet(w http.ResponseWriter, r *http.Request) {
	tenantID, _, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR2AGetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = uuid.New().String()
	}

	start := time.Now()
	resp, err := h.provider.GSTR2AGet(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr2a get failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) IMSGet(w http.ResponseWriter, r *http.Request) {
	tenantID, _, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.IMSGetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = uuid.New().String()
	}

	start := time.Now()
	resp, err := h.provider.IMSGet(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("ims get failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) IMSAction(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.IMSActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.IMSAction(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("ims action failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) IMSBulkAction(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.IMSBulkActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.IMSBulkAction(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("ims bulk action failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR3BSave(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR3BSaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GSTR3BSave(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr3b save failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR3BSubmit(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR3BSubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GSTR3BSubmit(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr3b submit failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR3BFile(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR3BFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GSTR3BFile(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr3b file failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR1SummaryHandler(w http.ResponseWriter, r *http.Request) {
	tenantID, _, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR1SummaryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = uuid.New().String()
	}

	start := time.Now()
	resp, err := h.provider.GSTR1Summary(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr1 summary failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR9Save(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR9SaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GSTR9Save(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr9 save failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR9Submit(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR9SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GSTR9Submit(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr9 submit failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR9File(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR9FileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GSTR9File(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr9 file failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR9Status(w http.ResponseWriter, r *http.Request) {
	tenantID, _, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR9StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = uuid.New().String()
	}

	start := time.Now()
	resp, err := h.provider.GSTR9Status(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr9 status failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR9CSave(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR9CSaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GSTR9CSave(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr9c save failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR9CFile(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR9CFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GSTR9CFile(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr9c file failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func (h *Handlers) GSTR9CStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, _, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTR9CStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = uuid.New().String()
	}

	start := time.Now()
	resp, err := h.provider.GSTR9CStatus(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstr9c status failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: resp.Status,
		},
	})
}

func extractHeaders(r *http.Request) (tenantID, idempotencyKey string, err error) {
	tenantID = r.Header.Get("X-Tenant-Id")
	if tenantID == "" {
		return "", "", fmt.Errorf("missing X-Tenant-Id header")
	}
	if _, err := uuid.Parse(tenantID); err != nil {
		return "", "", fmt.Errorf("invalid X-Tenant-Id: %w", err)
	}
	idempotencyKey = r.Header.Get("X-Idempotency-Key")
	if idempotencyKey == "" {
		idempotencyKey = uuid.New().String()
	}
	return tenantID, idempotencyKey, nil
}
