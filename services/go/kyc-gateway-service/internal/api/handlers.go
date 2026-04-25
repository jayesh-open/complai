package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/kyc-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/kyc-gateway-service/internal/provider"
)

type Handlers struct {
	provider provider.KYCProvider
}

func NewHandlers(p provider.KYCProvider) *Handlers {
	return &Handlers{provider: p}
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "kyc-gateway-service"})
}

func (h *Handlers) VerifyPAN(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.PANVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.VerifyPAN(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("pan verify failed")
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

func (h *Handlers) VerifyGSTIN(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GSTINVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.VerifyGSTIN(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("gstin verify failed")
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

func (h *Handlers) VerifyTAN(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.TANVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.VerifyTAN(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("tan verify failed")
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

func (h *Handlers) VerifyBank(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.BankVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.VerifyBank(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("bank verify failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	status := "success"
	if !resp.Valid {
		status = "invalid"
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      req.RequestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: status,
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
