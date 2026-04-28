package api

import (
	"encoding/json"
	"net/http"

	"github.com/complai/complai/services/go/tds-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/tds-gateway-service/internal/provider"
	"github.com/rs/zerolog/log"
)

type Handlers struct {
	provider provider.SandboxTDSProvider
}

func NewHandlers(p provider.SandboxTDSProvider) *Handlers {
	return &Handlers{provider: p}
}

func (h *Handlers) VerifyPAN(w http.ResponseWriter, r *http.Request) {
	var req domain.PANVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.PAN == "" {
		writeError(w, http.StatusBadRequest, "pan is required")
		return
	}
	req.TenantID = r.Header.Get("X-Tenant-Id")

	resp, err := h.provider.VerifyPAN(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("PAN verification failed")
		writeError(w, http.StatusInternalServerError, "verification failed")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handlers) VerifyTAN(w http.ResponseWriter, r *http.Request) {
	var req domain.TANVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.TAN == "" {
		writeError(w, http.StatusBadRequest, "tan is required")
		return
	}
	req.TenantID = r.Header.Get("X-Tenant-Id")

	resp, err := h.provider.VerifyTAN(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("TAN verification failed")
		writeError(w, http.StatusInternalServerError, "verification failed")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handlers) GenerateChallan(w http.ResponseWriter, r *http.Request) {
	var req domain.ChallanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.TAN == "" || req.Section == "" || req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "tan, section, and positive amount are required")
		return
	}
	req.TenantID = r.Header.Get("X-Tenant-Id")

	resp, err := h.provider.GenerateChallan(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("challan generation failed")
		writeError(w, http.StatusInternalServerError, "challan generation failed")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handlers) FileForm26Q(w http.ResponseWriter, r *http.Request) {
	var req domain.Form26QRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.TAN == "" || req.FinancialYear == "" || req.Quarter == "" || len(req.Deductions) == 0 {
		writeError(w, http.StatusBadRequest, "tan, financial_year, quarter, and deductions are required")
		return
	}
	req.TenantID = r.Header.Get("X-Tenant-Id")

	resp, err := h.provider.FileForm26Q(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("Form 26Q filing failed")
		writeError(w, http.StatusInternalServerError, "filing failed")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handlers) FileForm24Q(w http.ResponseWriter, r *http.Request) {
	var req domain.Form24QRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.TAN == "" || req.FinancialYear == "" || req.Quarter == "" || len(req.Employees) == 0 {
		writeError(w, http.StatusBadRequest, "tan, financial_year, quarter, and employees are required")
		return
	}
	req.TenantID = r.Header.Get("X-Tenant-Id")

	resp, err := h.provider.FileForm24Q(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("Form 24Q filing failed")
		writeError(w, http.StatusInternalServerError, "filing failed")
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
