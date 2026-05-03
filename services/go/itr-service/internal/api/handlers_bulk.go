package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/complai/complai/services/go/itr-service/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

func (h *Handlers) CreateBulkBatch(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	var req struct {
		TaxYear      string `json:"tax_year"`
		EmployerTAN  string `json:"employer_tan"`
		EmployerName string `json:"employer_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.TaxYear == "" || req.EmployerTAN == "" || req.EmployerName == "" {
		writeError(w, http.StatusBadRequest, "tax_year, employer_tan, and employer_name are required")
		return
	}

	b := &domain.BulkFilingBatch{
		ID:           uuid.New(),
		TenantID:     tenantID,
		TaxYear:      req.TaxYear,
		EmployerTAN:  req.EmployerTAN,
		EmployerName: req.EmployerName,
		Status:       domain.BatchPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := h.store.CreateBulkBatch(r.Context(), tenantID, b); err != nil {
		log.Error().Err(err).Msg("create bulk batch failed")
		writeError(w, http.StatusInternalServerError, "failed to create batch")
		return
	}
	writeJSON(w, http.StatusCreated, b)
}

func (h *Handlers) GetBulkBatch(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "batchId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid batch id")
		return
	}
	b, err := h.store.GetBulkBatch(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "batch not found")
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (h *Handlers) ListBulkBatches(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	limit, offset := pagination(r)
	list, total, err := h.store.ListBulkBatches(r.Context(), tenantID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("list bulk batches failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"batches": list, "total": total})
}

func (h *Handlers) AddBulkEmployee(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	batchID, err := uuid.Parse(chi.URLParam(r, "batchId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid batch id")
		return
	}

	batch, err := h.store.GetBulkBatch(r.Context(), tenantID, batchID)
	if err != nil {
		writeError(w, http.StatusNotFound, "batch not found")
		return
	}
	if batch.TotalEmployees >= domain.MaxBulkBatchSize() {
		writeError(w, http.StatusBadRequest, "batch has reached the 1000-employee limit")
		return
	}

	var req struct {
		PAN         string  `json:"pan"`
		Name        string  `json:"name"`
		Email       string  `json:"email"`
		GrossSalary float64 `json:"gross_salary"`
		TDSDeducted float64 `json:"tds_deducted"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.PAN) != 10 {
		writeError(w, http.StatusBadRequest, "PAN must be 10 characters")
		return
	}

	e := &domain.BulkFilingEmployee{
		ID:          uuid.New(),
		TenantID:    tenantID,
		BatchID:     batchID,
		PAN:         req.PAN,
		Name:        req.Name,
		Email:       req.Email,
		GrossSalary: decimal.NewFromFloat(req.GrossSalary),
		TDSDeducted: decimal.NewFromFloat(req.TDSDeducted),
		FormType:    domain.FormITR1,
		Status:      domain.EmpPendingReview,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := h.store.CreateBulkEmployee(r.Context(), tenantID, e); err != nil {
		log.Error().Err(err).Msg("create bulk employee failed")
		writeError(w, http.StatusInternalServerError, "failed to add employee")
		return
	}

	batch.TotalEmployees++
	_ = h.store.UpdateBulkBatchStatus(r.Context(), tenantID, batchID, batch.Status, batch.Processed, batch.Ready, batch.WithMismatches)

	writeJSON(w, http.StatusCreated, e)
}

func (h *Handlers) ListBulkEmployees(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	batchID, err := uuid.Parse(chi.URLParam(r, "batchId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid batch id")
		return
	}
	limit, offset := pagination(r)
	list, total, err := h.store.ListBulkEmployees(r.Context(), tenantID, batchID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("list bulk employees failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"employees": list, "total": total})
}

func (h *Handlers) ProcessBulkBatch(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	batchID, err := uuid.Parse(chi.URLParam(r, "batchId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid batch id")
		return
	}

	batch, err := h.store.GetBulkBatch(r.Context(), tenantID, batchID)
	if err != nil {
		writeError(w, http.StatusNotFound, "batch not found")
		return
	}
	if batch.Status != domain.BatchPending {
		writeError(w, http.StatusConflict, "batch already processed or processing")
		return
	}

	var req struct {
		AIS domain.AISSourceData `json:"ais"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	employees, _, err := h.store.ListBulkEmployees(r.Context(), tenantID, batchID, 1000, 0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list employees")
		return
	}
	if len(employees) == 0 {
		writeError(w, http.StatusBadRequest, "no employees in batch")
		return
	}

	_ = h.store.UpdateBulkBatchStatus(r.Context(), tenantID, batchID, domain.BatchProcessing, 0, 0, 0)

	var processed, ready, mismatches int
	results := make([]domain.BulkProcessResult, 0, len(employees))

	for _, emp := range employees {
		input := domain.BulkProcessInput{
			PAN: emp.PAN, Name: emp.Name, Email: emp.Email,
			GrossSalary: emp.GrossSalary, TDSDeducted: emp.TDSDeducted,
		}
		result := domain.ProcessEmployeeForBulkFiling(input, req.AIS)
		results = append(results, result)

		_ = h.store.UpdateBulkEmployeeStatus(r.Context(), tenantID, emp.ID, result.Status)
		processed++
		if result.Status == domain.EmpPendingReview {
			ready++
		}
		if result.Status == domain.EmpMismatch {
			mismatches++
		}
	}

	_ = h.store.UpdateBulkBatchStatus(r.Context(), tenantID, batchID, domain.BatchCompleted, processed, ready, mismatches)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"batch_id":       batchID,
		"processed":      processed,
		"ready":          ready,
		"with_mismatches": mismatches,
		"results":        results,
	})
}
