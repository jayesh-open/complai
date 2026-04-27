package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/irp-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/irp-gateway-service/internal/provider"
)

type Handlers struct {
	provider provider.IRPProvider
}

func NewHandlers(p provider.IRPProvider) *Handlers {
	return &Handlers{provider: p}
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "irp-gateway-service"})
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

func (h *Handlers) GenerateIRN(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GenerateIRNRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.GenerateIRN(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("generate IRN failed")
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

func (h *Handlers) CancelIRN(w http.ResponseWriter, r *http.Request) {
	tenantID, idempotencyKey, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.CancelIRNRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RequestID == "" {
		req.RequestID = idempotencyKey
	}

	start := time.Now()
	resp, err := h.provider.CancelIRN(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("cancel IRN failed")
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

func (h *Handlers) GetIRNByIRN(w http.ResponseWriter, r *http.Request) {
	tenantID, _, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	irn := r.URL.Query().Get("irn")
	if irn == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "irn query parameter is required"})
		return
	}

	requestID := r.URL.Query().Get("request_id")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	start := time.Now()
	resp, err := h.provider.GetIRNByIRN(r.Context(), &domain.GetIRNByIRNRequest{
		IRN:       irn,
		RequestID: requestID,
	})
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Str("irn", irn).Msg("get IRN by IRN failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      requestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: "success",
		},
	})
}

func (h *Handlers) GetIRNByDoc(w http.ResponseWriter, r *http.Request) {
	tenantID, _, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	docType := r.URL.Query().Get("doctype")
	docNum := r.URL.Query().Get("docnum")
	docDate := r.URL.Query().Get("docdate")
	if docType == "" || docNum == "" || docDate == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "doctype, docnum, and docdate query parameters are required"})
		return
	}

	requestID := r.URL.Query().Get("request_id")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	start := time.Now()
	resp, err := h.provider.GetIRNByDoc(r.Context(), &domain.GetIRNByDocRequest{
		DocType:   docType,
		DocNum:    docNum,
		DocDate:   docDate,
		RequestID: requestID,
	})
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Msg("get IRN by doc failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      requestID,
			LatencyMs:      int(time.Since(start).Milliseconds()),
			ProviderStatus: "success",
		},
	})
}

func (h *Handlers) ValidateGSTIN(w http.ResponseWriter, r *http.Request) {
	tenantID, _, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	gstin := r.URL.Query().Get("gstin")
	if gstin == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "gstin query parameter is required"})
		return
	}

	requestID := r.URL.Query().Get("request_id")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	start := time.Now()
	resp, err := h.provider.ValidateGSTIN(r.Context(), &domain.GSTINValidateRequest{
		GSTIN:     gstin,
		RequestID: requestID,
	})
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID).Str("gstin", gstin).Msg("validate GSTIN failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      requestID,
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
