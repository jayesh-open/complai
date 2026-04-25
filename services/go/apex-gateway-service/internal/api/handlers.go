package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/apex-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/apex-gateway-service/internal/provider"
)

type Handlers struct {
	provider provider.ApexProvider
}

func NewHandlers(p provider.ApexProvider) *Handlers {
	return &Handlers{provider: p}
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "apex-gateway-service"})
}

func (h *Handlers) FetchVendors(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.FetchVendorsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.TenantID == "" {
		req.TenantID = tenantID
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.FetchVendors(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("fetch vendors failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: "success",
		},
	})
}

func (h *Handlers) FetchAPInvoices(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.FetchAPInvoicesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.TenantID == "" {
		req.TenantID = tenantID
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.FetchAPInvoices(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("fetch ap invoices failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: "success",
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
