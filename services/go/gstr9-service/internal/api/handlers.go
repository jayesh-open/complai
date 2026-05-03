package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/complai/complai/services/go/gstr9-service/internal/domain"
	"github.com/complai/complai/services/go/gstr9-service/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var (
	gstinRe = regexp.MustCompile(`^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z][0-9A-Z][Z][0-9A-Z]$`)
	fyRe    = regexp.MustCompile(`^\d{4}-\d{2}$`)
)

type Handlers struct {
	store     store.Repository
	gstSvcURL string
}

func NewHandlers(s store.Repository, gstSvcURL string) *Handlers {
	return &Handlers{store: s, gstSvcURL: gstSvcURL}
}

func (h *Handlers) CreateAnnualReturn(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}

	var req struct {
		GSTIN         string `json:"gstin"`
		FinancialYear string `json:"financial_year"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !gstinRe.MatchString(req.GSTIN) {
		writeError(w, http.StatusBadRequest, "invalid GSTIN format")
		return
	}
	if !fyRe.MatchString(req.FinancialYear) {
		writeError(w, http.StatusBadRequest, "invalid financial_year format (use YYYY-YY)")
		return
	}

	now := time.Now()
	filing := &domain.GSTR9Filing{
		ID:            uuid.New(),
		TenantID:      tenantID,
		GSTIN:         req.GSTIN,
		FinancialYear: req.FinancialYear,
		Status:        domain.FilingStatusDraft,
		RequestID:     uuid.New(),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := h.store.CreateFiling(r.Context(), tenantID, filing); err != nil {
		if err == domain.ErrDuplicateFiling {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		log.Error().Err(err).Msg("create filing failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	_ = h.store.CreateAuditLog(r.Context(), tenantID, &domain.GSTR9AuditLog{
		ID: uuid.New(), TenantID: tenantID, FilingID: filing.ID,
		Action: "created", Details: "annual return initiated", ActorID: tenantID,
		CreatedAt: now,
	})

	writeJSON(w, http.StatusCreated, filing)
}

func (h *Handlers) GetAnnualReturn(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid filing id")
		return
	}

	filing, err := h.store.GetFiling(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "filing not found")
		return
	}
	writeJSON(w, http.StatusOK, filing)
}

func (h *Handlers) ListAnnualReturns(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	limit, offset := pagination(r)
	gstin := r.URL.Query().Get("gstin")
	fy := r.URL.Query().Get("fy")

	list, total, err := h.store.ListFilings(r.Context(), tenantID, gstin, fy, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("list filings failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"filings": list, "total": total})
}

func (h *Handlers) GetTableData(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid filing id")
		return
	}
	tableNum := chi.URLParam(r, "table")

	all, err := h.store.ListTableData(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	var filtered []domain.GSTR9TableData
	for _, td := range all {
		if td.TableNumber == tableNum {
			filtered = append(filtered, td)
		}
	}
	if len(filtered) == 0 {
		writeError(w, http.StatusNotFound, "table data not found")
		return
	}
	writeJSON(w, http.StatusOK, filtered)
}

func (h *Handlers) SaveAnnualReturn(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid filing id")
		return
	}

	filing, err := h.store.GetFiling(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "filing not found")
		return
	}

	if filing.Status != domain.FilingStatusDraft && filing.Status != domain.FilingStatusAggregated {
		writeError(w, http.StatusBadRequest, "filing cannot be saved in current status")
		return
	}

	if err := h.store.UpdateFilingStatus(r.Context(), tenantID, id, domain.FilingStatusSaved); err != nil {
		log.Error().Err(err).Msg("save filing failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	_ = h.store.CreateAuditLog(r.Context(), tenantID, &domain.GSTR9AuditLog{
		ID: uuid.New(), TenantID: tenantID, FilingID: id,
		Action: "saved", Details: "annual return saved as draft", ActorID: tenantID,
		CreatedAt: time.Now(),
	})

	filing.Status = domain.FilingStatusSaved
	writeJSON(w, http.StatusOK, filing)
}

func (h *Handlers) AggregateAnnualReturn(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid filing id")
		return
	}

	filing, err := h.store.GetFiling(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "filing not found")
		return
	}

	var req struct {
		Months []domain.MonthlyData `json:"months"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.Months) == 0 {
		writeError(w, http.StatusBadRequest, "months data required")
		return
	}

	_ = h.store.DeleteTableData(r.Context(), tenantID, id)

	tables := domain.Aggregate(id, tenantID, req.Months)
	for i := range tables {
		if err := h.store.CreateTableData(r.Context(), tenantID, &tables[i]); err != nil {
			log.Error().Err(err).Msg("store table data failed")
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	turnover := domain.ComputeAggregateTurnover(tables)
	threshold := domain.CheckThreshold(turnover)

	filing.AggregateTurnover = turnover
	filing.IsMandatory = threshold.GSTR9Mandatory
	filing.Status = domain.FilingStatusAggregated
	_ = h.store.UpdateFilingStatus(r.Context(), tenantID, id, domain.FilingStatusAggregated)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"filing":    filing,
		"threshold": threshold,
		"tables":    len(tables),
	})
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

func (h *Handlers) InitiateReconciliation(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	gstr9ID, err := uuid.Parse(chi.URLParam(r, "gstr9Id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid gstr9 filing id")
		return
	}

	gstr9Filing, err := h.store.GetFiling(r.Context(), tenantID, gstr9ID)
	if err != nil {
		writeError(w, http.StatusNotFound, "gstr9 filing not found")
		return
	}

	if _, err := h.store.GetGSTR9CFilingByGSTR9ID(r.Context(), tenantID, gstr9ID); err == nil {
		writeError(w, http.StatusConflict, domain.ErrGSTR9CDuplicate.Error())
		return
	}

	var req domain.AuditedFinancials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tables, err := h.store.ListTableData(r.Context(), tenantID, gstr9ID)
	if err != nil || len(tables) == 0 {
		writeError(w, http.StatusBadRequest, "gstr9 must be aggregated before reconciliation")
		return
	}

	now := time.Now()
	gstr9cFiling := &domain.GSTR9CFiling{
		ID:            uuid.New(),
		TenantID:      tenantID,
		GSTR9FilingID: gstr9ID,
		Status:        domain.GSTR9CStatusDraft,
		AuditedTurnover: req.Turnover,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := h.store.CreateGSTR9CFiling(r.Context(), tenantID, gstr9cFiling); err != nil {
		log.Error().Err(err).Msg("create gstr9c filing failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	mismatches := domain.Reconcile(gstr9Filing, tables, req, tenantID, gstr9cFiling.ID)

	var unreconciledTotal = decimal.Zero
	for i := range mismatches {
		if err := h.store.CreateMismatch(r.Context(), tenantID, &mismatches[i]); err != nil {
			log.Error().Err(err).Msg("create mismatch failed")
		}
		unreconciledTotal = unreconciledTotal.Add(mismatches[i].Difference.Abs())
	}

	gstr9cFiling.UnreconciledAmount = unreconciledTotal
	_ = h.store.UpdateGSTR9CUnreconciled(r.Context(), tenantID, gstr9cFiling.ID, unreconciledTotal)
	_ = h.store.UpdateGSTR9CStatus(r.Context(), tenantID, gstr9cFiling.ID, domain.GSTR9CStatusReconciled)
	gstr9cFiling.Status = domain.GSTR9CStatusReconciled

	result := domain.ReconciliationResult{
		GSTR9CFiling: *gstr9cFiling,
		Mismatches:   mismatches,
		CanSubmit:    domain.CanSubmit(mismatches),
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *Handlers) GetReconciliation(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid reconciliation id")
		return
	}

	filing, err := h.store.GetGSTR9CFiling(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "reconciliation not found")
		return
	}

	mismatches, _ := h.store.ListMismatches(r.Context(), tenantID, id)
	result := domain.ReconciliationResult{
		GSTR9CFiling: *filing,
		Mismatches:   mismatches,
		CanSubmit:    domain.CanSubmit(mismatches),
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handlers) ListReconciliationMismatches(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid reconciliation id")
		return
	}

	if _, err := h.store.GetGSTR9CFiling(r.Context(), tenantID, id); err != nil {
		writeError(w, http.StatusNotFound, "reconciliation not found")
		return
	}

	mismatches, err := h.store.ListMismatches(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, mismatches)
}

func (h *Handlers) ResolveMismatch(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	reconcID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid reconciliation id")
		return
	}
	mismatchID, err := uuid.Parse(chi.URLParam(r, "mismatchId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid mismatch id")
		return
	}

	if _, err := h.store.GetGSTR9CFiling(r.Context(), tenantID, reconcID); err != nil {
		writeError(w, http.StatusNotFound, "reconciliation not found")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Reason == "" {
		writeError(w, http.StatusBadRequest, "reason is required")
		return
	}

	mm, err := h.store.GetMismatch(r.Context(), tenantID, mismatchID)
	if err != nil {
		writeError(w, http.StatusNotFound, "mismatch not found")
		return
	}
	if mm.GSTR9CFilingID != reconcID {
		writeError(w, http.StatusNotFound, "mismatch not found")
		return
	}

	if err := h.store.ResolveMismatch(r.Context(), tenantID, mismatchID, req.Reason, tenantID); err != nil {
		log.Error().Err(err).Msg("resolve mismatch failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	mm.Resolved = true
	mm.ResolvedReason = req.Reason
	writeJSON(w, http.StatusOK, mm)
}

func (h *Handlers) CertifyReconciliation(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid reconciliation id")
		return
	}

	filing, err := h.store.GetGSTR9CFiling(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "reconciliation not found")
		return
	}
	if filing.Status == domain.GSTR9CStatusCertified {
		writeError(w, http.StatusBadRequest, domain.ErrGSTR9CAlreadyCertified.Error())
		return
	}
	if filing.Status != domain.GSTR9CStatusReconciled {
		writeError(w, http.StatusBadRequest, domain.ErrGSTR9CNotReconciled.Error())
		return
	}

	mismatches, _ := h.store.ListMismatches(r.Context(), tenantID, id)
	if !domain.CanSubmit(mismatches) {
		writeError(w, http.StatusBadRequest, domain.ErrUnresolvedMismatches.Error())
		return
	}

	if err := h.store.CertifyGSTR9C(r.Context(), tenantID, id, tenantID); err != nil {
		log.Error().Err(err).Msg("certify gstr9c failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	filing.Status = domain.GSTR9CStatusCertified
	filing.IsSelfCertified = true
	writeJSON(w, http.StatusOK, filing)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": msg})
}

