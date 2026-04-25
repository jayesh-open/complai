package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/gst-service/internal/categorizer"
	"github.com/complai/complai/services/go/gst-service/internal/domain"
	"github.com/complai/complai/services/go/gst-service/internal/gateway"
	"github.com/complai/complai/services/go/gst-service/internal/store"
)

type StepUpVerifier interface {
	HasValidStepUp(ctx context.Context, tenantID, userID uuid.UUID, action string) bool
}

type Handlers struct {
	store      store.Repository
	auraClient *gateway.AuraClient
	gstnClient *gateway.GSTNClient
	stepUp     StepUpVerifier
}

func NewHandlers(s store.Repository, aura *gateway.AuraClient, gstn *gateway.GSTNClient, stepUp StepUpVerifier) *Handlers {
	return &Handlers{store: s, auraClient: aura, gstnClient: gstn, stepUp: stepUp}
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "gst-service"})
}

func (h *Handlers) Ingest(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.IngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.GSTIN == "" || req.ReturnPeriod == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "gstin and return_period are required"})
		return
	}

	entries, err := h.auraClient.FetchARInvoices(r.Context(), tenantID, req.GSTIN, req.ReturnPeriod)
	if err != nil {
		log.Error().Err(err).Msg("fetch AR invoices failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "failed to fetch invoices from Aura"})
		return
	}

	for i := range entries {
		entries[i].Section = categorizer.Categorize(&entries[i])
	}

	existing, _ := h.store.CountEntries(r.Context(), tenantID, req.GSTIN, req.ReturnPeriod)

	inserted, err := h.store.BulkInsertEntries(r.Context(), tenantID, entries)
	if err != nil {
		log.Error().Err(err).Msg("bulk insert entries failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to store entries"})
		return
	}

	filing := &domain.GSTR1Filing{
		GSTIN:        req.GSTIN,
		ReturnPeriod: req.ReturnPeriod,
		Status:       domain.FilingStatusIngested,
		TotalCount:   inserted,
	}
	if uid := userIDFromRequest(r); uid != uuid.Nil {
		filing.CreatedBy = &uid
	}

	if err := h.store.CreateFiling(r.Context(), tenantID, filing); err != nil {
		log.Error().Err(err).Msg("create filing failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create filing"})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.IngestResponse{
		FilingID:   filing.ID,
		Ingested:   inserted,
		Duplicates: existing,
	})
}

func (h *Handlers) Validate(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	filing, err := h.store.GetFiling(r.Context(), tenantID, req.FilingID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "filing not found"})
		return
	}

	if filing.Status != domain.FilingStatusIngested && filing.Status != domain.FilingStatusDraft {
		httputil.JSON(w, http.StatusConflict, map[string]string{"error": "filing must be in ingested or draft state to validate"})
		return
	}

	entries, err := h.store.ListEntries(r.Context(), tenantID, filing.ID, "")
	if err != nil {
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list entries"})
		return
	}

	validationErrors := validateEntries(filing.ID, tenantID, entries)

	if len(validationErrors) > 0 {
		if err := h.store.CreateValidationErrors(r.Context(), tenantID, validationErrors); err != nil {
			log.Error().Err(err).Msg("store validation errors failed")
		}
	}

	sections := buildSections(filing.ID, tenantID, entries)
	if err := h.store.CreateSections(r.Context(), tenantID, sections); err != nil {
		log.Error().Err(err).Msg("store sections failed")
	}

	status := domain.FilingStatusValidated
	if len(validationErrors) > 0 {
		status = domain.FilingStatusDraft
	}
	_ = h.store.UpdateFilingStatus(r.Context(), tenantID, filing.ID, status)

	httputil.JSON(w, http.StatusOK, domain.ValidateResponse{
		FilingID:   filing.ID,
		TotalCount: len(entries),
		ErrorCount: len(validationErrors),
		Sections:   sections,
	})
}

func (h *Handlers) Approve(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.ApproveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	filing, err := h.store.GetFiling(r.Context(), tenantID, req.FilingID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "filing not found"})
		return
	}

	if filing.Status != domain.FilingStatusValidated {
		httputil.JSON(w, http.StatusConflict, map[string]string{"error": "filing must be validated before approval"})
		return
	}

	approverID := userIDFromRequest(r)
	if approverID == uuid.Nil {
		approverID = req.ApprovedBy
	}

	if filing.CreatedBy != nil && *filing.CreatedBy == approverID {
		httputil.JSON(w, http.StatusForbidden, map[string]string{
			"error":   "self_approval_denied",
			"message": "Cannot approve your own filing (maker-checker)",
		})
		return
	}

	if err := h.store.ApproveFiling(r.Context(), tenantID, req.FilingID, approverID); err != nil {
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "approval failed"})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "approved", "filing_id": req.FilingID.String()})
}

func (h *Handlers) File(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.FileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.SignType != "DSC" && req.SignType != "EVC" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "sign_type must be DSC or EVC"})
		return
	}

	filerID := userIDFromRequest(r)
	if filerID == uuid.Nil {
		filerID = req.FiledBy
	}

	if h.stepUp != nil && !h.stepUp.HasValidStepUp(r.Context(), tenantID, filerID, "gstr1_file") {
		httputil.JSON(w, http.StatusForbidden, map[string]string{
			"error":       "step_up_required",
			"message":     "Filing requires step-up authentication",
			"step_up_url": "/v1/auth/step-up",
		})
		return
	}

	filing, err := h.store.GetFiling(r.Context(), tenantID, req.FilingID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "filing not found"})
		return
	}

	if filing.Status != domain.FilingStatusApproved {
		httputil.JSON(w, http.StatusConflict, map[string]string{"error": "filing must be approved before filing"})
		return
	}

	_ = h.store.UpdateFilingStatus(r.Context(), tenantID, filing.ID, domain.FilingStatusSaved)

	gstnSaveReq := gateway.GSTR1SaveRequest{
		GSTIN:     filing.GSTIN,
		RetPeriod: filing.ReturnPeriod,
		B2B:       []string{},
	}
	_, err = h.gstnClient.SaveGSTR1(r.Context(), tenantID, gstnSaveReq)
	if err != nil {
		log.Error().Err(err).Msg("GSTN save failed")
		_ = h.store.UpdateFilingStatus(r.Context(), tenantID, filing.ID, domain.FilingStatusFailed)
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "GSTN save failed"})
		return
	}

	gstnSubmitReq := gateway.GSTR1SubmitRequest{
		GSTIN:     filing.GSTIN,
		RetPeriod: filing.ReturnPeriod,
	}
	_, err = h.gstnClient.SubmitGSTR1(r.Context(), tenantID, gstnSubmitReq)
	if err != nil {
		log.Error().Err(err).Msg("GSTN submit failed")
		_ = h.store.UpdateFilingStatus(r.Context(), tenantID, filing.ID, domain.FilingStatusFailed)
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "GSTN submit failed"})
		return
	}

	_ = h.store.UpdateFilingStatus(r.Context(), tenantID, filing.ID, domain.FilingStatusSubmitted)

	gstnFileReq := gateway.GSTR1FileRequest{
		GSTIN:     filing.GSTIN,
		RetPeriod: filing.ReturnPeriod,
		SignType:  req.SignType,
		OTP:       req.OTP,
	}
	fileResp, err := h.gstnClient.FileGSTR1(r.Context(), tenantID, gstnFileReq)
	if err != nil {
		log.Error().Err(err).Msg("GSTN file failed")
		_ = h.store.UpdateFilingStatus(r.Context(), tenantID, filing.ID, domain.FilingStatusFailed)
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "GSTN file failed"})
		return
	}

	if err := h.store.UpdateFilingARN(r.Context(), tenantID, filing.ID, fileResp.ARN, filerID); err != nil {
		log.Error().Err(err).Msg("update ARN failed")
	}

	httputil.JSON(w, http.StatusOK, domain.FileResponse{
		FilingID: filing.ID,
		Status:   domain.FilingStatusFiled,
		ARN:      fileResp.ARN,
	})
}

func (h *Handlers) Summary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	filingID, err := uuid.Parse(r.URL.Query().Get("filing_id"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid filing_id"})
		return
	}

	filing, err := h.store.GetFiling(r.Context(), tenantID, filingID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "filing not found"})
		return
	}

	sections, _ := h.store.ListSections(r.Context(), tenantID, filingID)
	errCount, _ := h.store.CountValidationErrors(r.Context(), tenantID, filingID)

	httputil.JSON(w, http.StatusOK, domain.GSTR1Summary{
		Filing:   *filing,
		Sections: sections,
		Errors:   errCount,
	})
}

func (h *Handlers) ListEntries(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	filingID, err := uuid.Parse(r.URL.Query().Get("filing_id"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid filing_id"})
		return
	}

	section := r.URL.Query().Get("section")

	entries, err := h.store.ListEntries(r.Context(), tenantID, filingID, section)
	if err != nil {
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list entries"})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"entries":     entries,
		"total_count": len(entries),
	})
}

func (h *Handlers) ListErrors(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	filingID, err := uuid.Parse(r.URL.Query().Get("filing_id"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid filing_id"})
		return
	}

	errs, err := h.store.ListValidationErrors(r.Context(), tenantID, filingID)
	if err != nil {
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list errors"})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"errors":      errs,
		"total_count": len(errs),
	})
}

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}

func userIDFromRequest(r *http.Request) uuid.UUID {
	id, err := uuid.Parse(r.Header.Get("X-User-Id"))
	if err != nil {
		return uuid.Nil
	}
	return id
}

func validateEntries(filingID, tenantID uuid.UUID, entries []domain.SalesRegisterEntry) []domain.ValidationError {
	var errs []domain.ValidationError

	for _, e := range entries {
		if e.DocumentNumber == "" {
			errs = append(errs, domain.ValidationError{
				ID: uuid.New(), TenantID: tenantID, FilingID: filingID, EntryID: e.ID,
				Field: "document_number", Code: "REQUIRED", Message: "Document number is required", Severity: "error",
			})
		}

		if e.SupplyType == "B2B" && (e.BuyerGSTIN == "" || e.BuyerGSTIN == "URP") {
			errs = append(errs, domain.ValidationError{
				ID: uuid.New(), TenantID: tenantID, FilingID: filingID, EntryID: e.ID,
				Field: "buyer_gstin", Code: "B2B_GSTIN_REQUIRED", Message: "B2B supply requires buyer GSTIN", Severity: "error",
			})
		}

		if e.TaxableValue.IsNegative() {
			errs = append(errs, domain.ValidationError{
				ID: uuid.New(), TenantID: tenantID, FilingID: filingID, EntryID: e.ID,
				Field: "taxable_value", Code: "NEGATIVE_VALUE", Message: "Taxable value cannot be negative", Severity: "error",
			})
		}

		if e.HSN == "" {
			errs = append(errs, domain.ValidationError{
				ID: uuid.New(), TenantID: tenantID, FilingID: filingID, EntryID: e.ID,
				Field: "hsn", Code: "HSN_REQUIRED", Message: "HSN code is required", Severity: "warning",
			})
		}

		if e.PlaceOfSupply == "" {
			errs = append(errs, domain.ValidationError{
				ID: uuid.New(), TenantID: tenantID, FilingID: filingID, EntryID: e.ID,
				Field: "place_of_supply", Code: "POS_REQUIRED", Message: "Place of supply is required", Severity: "error",
			})
		}
	}

	return errs
}

func buildSections(filingID, tenantID uuid.UUID, entries []domain.SalesRegisterEntry) []domain.GSTR1Section {
	sectionMap := make(map[string]*domain.GSTR1Section)

	for _, e := range entries {
		sec := categorizer.Categorize(&e)
		s, ok := sectionMap[sec]
		if !ok {
			s = &domain.GSTR1Section{
				ID:       uuid.New(),
				TenantID: tenantID,
				FilingID: filingID,
				Section:  sec,
				Status:   "computed",
			}
			sectionMap[sec] = s
		}
		s.InvoiceCount++
		s.TaxableValue = s.TaxableValue.Add(e.TaxableValue)
		s.CGST = s.CGST.Add(e.CGSTAmount)
		s.SGST = s.SGST.Add(e.SGSTAmount)
		s.IGST = s.IGST.Add(e.IGSTAmount)
		s.TotalTax = s.TotalTax.Add(e.CGSTAmount).Add(e.SGSTAmount).Add(e.IGSTAmount)
	}

	sections := make([]domain.GSTR1Section, 0, len(sectionMap))
	for _, s := range sectionMap {
		sections = append(sections, *s)
	}
	return sections
}

