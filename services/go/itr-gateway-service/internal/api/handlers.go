package api

import (
	"encoding/json"
	"net/http"

	"github.com/complai/complai/services/go/itr-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/itr-gateway-service/internal/provider"
	"github.com/rs/zerolog/log"
)

type Handlers struct {
	provider provider.SandboxITRProvider
}

func NewHandlers(p provider.SandboxITRProvider) *Handlers {
	return &Handlers{provider: p}
}

func (h *Handlers) CheckPANAadhaarLink(w http.ResponseWriter, r *http.Request) {
	var req domain.PANAadhaarLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.PAN == "" {
		writeError(w, http.StatusBadRequest, "pan is required")
		return
	}
	req.TenantID = r.Header.Get("X-Tenant-Id")

	resp, err := h.provider.CheckPANAadhaarLink(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("PAN-Aadhaar link check failed")
		writeError(w, http.StatusInternalServerError, "link check failed")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handlers) FetchAIS(w http.ResponseWriter, r *http.Request) {
	var req domain.AISRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.PAN == "" || req.TaxYear == "" {
		writeError(w, http.StatusBadRequest, "pan and tax_year are required")
		return
	}
	req.TenantID = r.Header.Get("X-Tenant-Id")

	resp, err := h.provider.FetchAIS(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("AIS fetch failed")
		writeError(w, http.StatusInternalServerError, "ais fetch failed")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handlers) SubmitITR(w http.ResponseWriter, r *http.Request) {
	var req domain.ITRSubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.PAN == "" || req.TaxYear == "" || req.FormType == "" {
		writeError(w, http.StatusBadRequest, "pan, tax_year, and form_type are required")
		return
	}
	req.TenantID = r.Header.Get("X-Tenant-Id")

	resp, err := h.provider.SubmitITR(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("ITR submission failed")
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handlers) GenerateITRV(w http.ResponseWriter, r *http.Request) {
	var req domain.ITRVRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ARN == "" {
		writeError(w, http.StatusBadRequest, "arn is required")
		return
	}
	req.TenantID = r.Header.Get("X-Tenant-Id")

	resp, err := h.provider.GenerateITRV(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("ITR-V generation failed")
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handlers) CheckEVerification(w http.ResponseWriter, r *http.Request) {
	var req domain.EVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ARN == "" {
		writeError(w, http.StatusBadRequest, "arn is required")
		return
	}
	req.TenantID = r.Header.Get("X-Tenant-Id")

	resp, err := h.provider.CheckEVerification(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("e-verification check failed")
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handlers) CheckRefundStatus(w http.ResponseWriter, r *http.Request) {
	var req domain.RefundStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.PAN == "" || req.TaxYear == "" {
		writeError(w, http.StatusBadRequest, "pan and tax_year are required")
		return
	}
	req.TenantID = r.Header.Get("X-Tenant-Id")

	resp, err := h.provider.CheckRefundStatus(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("refund status check failed")
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": v})
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": msg})
}
