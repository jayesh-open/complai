package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/ewb-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/ewb-gateway-service/internal/provider"
)

type Handlers struct {
	provider provider.EWBProvider
}

func NewHandlers(p provider.EWBProvider) *Handlers {
	return &Handlers{provider: p}
}

type headerInfo struct {
	tenantID       uuid.UUID
	idempotencyKey string
}

func extractHeaders(r *http.Request) (headerInfo, error) {
	tid := r.Header.Get("X-Tenant-Id")
	if tid == "" {
		return headerInfo{}, fmt.Errorf("missing X-Tenant-Id header")
	}
	parsed, err := uuid.Parse(tid)
	if err != nil {
		return headerInfo{}, fmt.Errorf("invalid X-Tenant-Id: %w", err)
	}
	idk := r.Header.Get("X-Idempotency-Key")
	if idk == "" {
		idk = uuid.New().String()
	}
	return headerInfo{tenantID: parsed, idempotencyKey: idk}, nil
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "ewb-gateway-service"})
}

func (h *Handlers) GenerateEWB(w http.ResponseWriter, r *http.Request) {
	hdr, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GenerateEWBRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.GSTIN == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "gstin is required"})
		return
	}
	if req.DocNo == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "doc_no is required"})
		return
	}

	start := time.Now()
	resp, err := h.provider.GenerateEWB(r.Context(), &req)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		log.Error().Err(err).Str("tenant_id", hdr.tenantID.String()).Msg("generate EWB failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      hdr.idempotencyKey,
			LatencyMS:      elapsed,
			ProviderStatus: "success",
		},
	})
}

func (h *Handlers) CancelEWB(w http.ResponseWriter, r *http.Request) {
	hdr, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.CancelEWBRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	start := time.Now()
	resp, err := h.provider.CancelEWB(r.Context(), &req)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		log.Error().Err(err).Str("tenant_id", hdr.tenantID.String()).Msg("cancel EWB failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      hdr.idempotencyKey,
			LatencyMS:      elapsed,
			ProviderStatus: "success",
		},
	})
}

func (h *Handlers) GetEWB(w http.ResponseWriter, r *http.Request) {
	hdr, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	ewbNo := r.URL.Query().Get("ewb_no")
	if ewbNo == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "ewb_no query parameter is required"})
		return
	}

	start := time.Now()
	resp, err := h.provider.GetEWB(r.Context(), ewbNo)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		log.Error().Err(err).Str("tenant_id", hdr.tenantID.String()).Str("ewb_no", ewbNo).Msg("get EWB failed")
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      hdr.idempotencyKey,
			LatencyMS:      elapsed,
			ProviderStatus: "success",
		},
	})
}

func (h *Handlers) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	hdr, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.UpdateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	start := time.Now()
	resp, err := h.provider.UpdateVehicle(r.Context(), &req)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		log.Error().Err(err).Str("tenant_id", hdr.tenantID.String()).Msg("update vehicle failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      hdr.idempotencyKey,
			LatencyMS:      elapsed,
			ProviderStatus: "success",
		},
	})
}

func (h *Handlers) ExtendValidity(w http.ResponseWriter, r *http.Request) {
	hdr, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.ExtendValidityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	start := time.Now()
	resp, err := h.provider.ExtendValidity(r.Context(), &req)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		log.Error().Err(err).Str("tenant_id", hdr.tenantID.String()).Msg("extend validity failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      hdr.idempotencyKey,
			LatencyMS:      elapsed,
			ProviderStatus: "success",
		},
	})
}

func (h *Handlers) ConsolidateEWB(w http.ResponseWriter, r *http.Request) {
	hdr, err := extractHeaders(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.ConsolidateEWBRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	start := time.Now()
	resp, err := h.provider.ConsolidateEWB(r.Context(), &req)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		log.Error().Err(err).Str("tenant_id", hdr.tenantID.String()).Msg("consolidate EWB failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GatewayResponse{
		Data: resp,
		Meta: domain.ResponseMeta{
			RequestID:      hdr.idempotencyKey,
			LatencyMS:      elapsed,
			ProviderStatus: "success",
		},
	})
}
