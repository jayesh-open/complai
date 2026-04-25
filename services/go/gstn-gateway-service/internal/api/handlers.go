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
