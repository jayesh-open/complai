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

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": msg})
}

