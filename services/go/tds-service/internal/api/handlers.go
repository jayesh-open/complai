package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/complai/complai/services/go/tds-service/internal/domain"
	"github.com/complai/complai/services/go/tds-service/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

type Handlers struct {
	store store.Repository
}

func NewHandlers(s store.Repository) *Handlers {
	return &Handlers{store: s}
}

func (h *Handlers) ListDeductees(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	limit, offset := pagination(r)

	list, total, err := h.store.ListDeductees(r.Context(), tenantID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("list deductees failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"deductees": list, "total": total})
}

func (h *Handlers) GetDeductee(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid deductee id")
		return
	}

	d, err := h.store.GetDeductee(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "deductee not found")
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (h *Handlers) CalculateTDS(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Section        domain.Section      `json:"section"`
		GrossAmount    float64             `json:"gross_amount"`
		DeducteeType   domain.DeducteeType `json:"deductee_type"`
		HasValidPAN    bool                `json:"has_valid_pan"`
		ResidentStatus domain.ResidentStatus `json:"resident_status"`
		RentType       domain.RentType     `json:"rent_type"`
		AggregateForFY float64             `json:"aggregate_for_fy"`
		AnnualSalary   float64             `json:"annual_salary"`
		DTAARate       *float64            `json:"dtaa_rate"`
		LowerCertRate  *float64            `json:"lower_cert_rate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !domain.ValidSection(req.Section) {
		writeError(w, http.StatusBadRequest, "invalid section")
		return
	}

	in := domain.CalcInput{
		Section:        req.Section,
		GrossAmount:    decimal.NewFromFloat(req.GrossAmount),
		DeducteeType:   req.DeducteeType,
		HasValidPAN:    req.HasValidPAN,
		ResidentStatus: req.ResidentStatus,
		RentType:       req.RentType,
		AggregateForFY: decimal.NewFromFloat(req.AggregateForFY),
		AnnualSalary:   decimal.NewFromFloat(req.AnnualSalary),
	}
	if req.DTAARate != nil {
		d := decimal.NewFromFloat(*req.DTAARate)
		in.DTAARate = &d
	}
	if req.LowerCertRate != nil {
		d := decimal.NewFromFloat(*req.LowerCertRate)
		in.LowerCertRate = &d
	}

	result := domain.Calculate(in)
	writeJSON(w, http.StatusOK, result)
}

func (h *Handlers) CreateEntry(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	var req struct {
		DeducteeID      string  `json:"deductee_id"`
		Section         string  `json:"section"`
		FinancialYear   string  `json:"financial_year"`
		Quarter         string  `json:"quarter"`
		TransactionDate string  `json:"transaction_date"`
		GrossAmount     float64 `json:"gross_amount"`
		NatureOfPayment string  `json:"nature_of_payment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	deducteeID, err := uuid.Parse(req.DeducteeID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid deductee_id")
		return
	}
	section := domain.Section(req.Section)
	if !domain.ValidSection(section) {
		writeError(w, http.StatusBadRequest, "invalid section")
		return
	}
	txnDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid transaction_date (use YYYY-MM-DD)")
		return
	}

	deductee, err := h.store.GetDeductee(r.Context(), tenantID, deducteeID)
	if err != nil {
		writeError(w, http.StatusNotFound, "deductee not found")
		return
	}

	agg, _ := h.store.GetAggregate(r.Context(), tenantID, deducteeID, section, req.FinancialYear)
	aggTotal := decimal.Zero
	if agg != nil {
		aggTotal = agg.TotalPaid
	}

	calcResult := domain.Calculate(domain.CalcInput{
		Section:        section,
		GrossAmount:    decimal.NewFromFloat(req.GrossAmount),
		DeducteeType:   deductee.DeducteeType,
		HasValidPAN:    deductee.PANVerified,
		ResidentStatus: deductee.ResidentStatus,
		AggregateForFY: aggTotal,
		AnnualSalary:   decimal.NewFromFloat(req.GrossAmount),
	})

	entry := &domain.TDSEntry{
		ID:              uuid.New(),
		TenantID:        tenantID,
		DeducteeID:      deducteeID,
		Section:         section,
		FinancialYear:   req.FinancialYear,
		Quarter:         req.Quarter,
		TransactionDate: txnDate,
		GrossAmount:     decimal.NewFromFloat(req.GrossAmount),
		TDSRate:         calcResult.Rate,
		TDSAmount:       calcResult.TDSAmount,
		Surcharge:       calcResult.Surcharge,
		Cess:            calcResult.Cess,
		TotalTax:        calcResult.TotalTax,
		NatureOfPayment: req.NatureOfPayment,
		PANAtDeduction:  deductee.PAN,
		NoPANDeduction:  calcResult.NoPAN,
		LowerCertApplied: calcResult.LowerCert,
		Status:          domain.StatusPending,
	}

	if err := h.store.CreateEntry(r.Context(), tenantID, entry); err != nil {
		log.Error().Err(err).Msg("create entry failed")
		writeError(w, http.StatusInternalServerError, "failed to create entry")
		return
	}

	newAgg := &domain.TDSAggregate{
		ID:               uuid.New(),
		TenantID:         tenantID,
		DeducteeID:       deducteeID,
		Section:          section,
		FinancialYear:    req.FinancialYear,
		TotalPaid:        aggTotal.Add(entry.GrossAmount),
		TotalTDS:         agg.TotalTDS.Add(entry.TotalTax),
		TransactionCount: agg.TransactionCount + 1,
	}
	h.store.UpsertAggregate(r.Context(), tenantID, newAgg)

	writeJSON(w, http.StatusCreated, entry)
}

func (h *Handlers) ListEntries(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	limit, offset := pagination(r)
	fy := r.URL.Query().Get("fy")
	quarter := r.URL.Query().Get("quarter")

	list, total, err := h.store.ListEntries(r.Context(), tenantID, fy, quarter, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("list entries failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": list, "total": total})
}

func (h *Handlers) GetEntry(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid entry id")
		return
	}

	e, err := h.store.GetEntry(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "entry not found")
		return
	}
	writeJSON(w, http.StatusOK, e)
}

func (h *Handlers) GetSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	fy := r.URL.Query().Get("fy")
	if fy == "" {
		fy = "2025-26"
	}

	sum, err := h.store.GetSummary(r.Context(), tenantID, fy)
	if err != nil {
		log.Error().Err(err).Msg("get summary failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, sum)
}

func tenantFrom(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(r.Header.Get("X-Tenant-Id"))
}

func pagination(r *http.Request) (int, int) {
	limit := 50
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return limit, offset
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
