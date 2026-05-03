package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/complai/complai/services/go/itr-service/internal/domain"
	"github.com/complai/complai/services/go/itr-service/internal/store"
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

func (h *Handlers) CreateTaxpayer(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	var req struct {
		PAN             string `json:"pan"`
		Name            string `json:"name"`
		DateOfBirth     string `json:"date_of_birth"`
		AssesseeType    string `json:"assessee_type"`
		ResidencyStatus string `json:"residency_status"`
		AadhaarLinked   bool   `json:"aadhaar_linked"`
		Email           string `json:"email"`
		Mobile          string `json:"mobile"`
		Address         string `json:"address"`
		EmployerTAN     string `json:"employer_tan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.PAN) != 10 {
		writeError(w, http.StatusBadRequest, "PAN must be 10 characters")
		return
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date_of_birth (use YYYY-MM-DD)")
		return
	}

	t := &domain.Taxpayer{
		ID:              uuid.New(),
		TenantID:        tenantID,
		PAN:             req.PAN,
		Name:            req.Name,
		DateOfBirth:     dob,
		AssesseeType:    domain.AssesseeType(req.AssesseeType),
		ResidencyStatus: domain.ResidencyStatus(req.ResidencyStatus),
		AadhaarLinked:   req.AadhaarLinked,
		Email:           req.Email,
		Mobile:          req.Mobile,
		Address:         req.Address,
		EmployerTAN:     req.EmployerTAN,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.store.CreateTaxpayer(r.Context(), tenantID, t); err != nil {
		log.Error().Err(err).Msg("create taxpayer failed")
		writeError(w, http.StatusInternalServerError, "failed to create taxpayer")
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func (h *Handlers) GetTaxpayer(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid taxpayer id")
		return
	}
	t, err := h.store.GetTaxpayer(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "taxpayer not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handlers) ListTaxpayers(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	limit, offset := pagination(r)
	list, total, err := h.store.ListTaxpayers(r.Context(), tenantID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("list taxpayers failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"taxpayers": list, "total": total})
}

func (h *Handlers) CreateFiling(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	var req struct {
		TaxpayerID  string `json:"taxpayer_id"`
		PAN         string `json:"pan"`
		TaxYear     string `json:"tax_year"`
		FormType    string `json:"form_type"`
		Regime      string `json:"regime"`
		Form10IEARef string `json:"form_10iea_ref"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	taxpayerID, err := uuid.Parse(req.TaxpayerID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid taxpayer_id")
		return
	}

	formType := domain.ITRFormType(req.FormType)
	if formType != domain.FormITR1 && formType != domain.FormITR2 && formType != domain.FormITR3 {
		writeError(w, http.StatusBadRequest, "form_type must be ITR-1, ITR-2, or ITR-3")
		return
	}

	regime := domain.RegimeType(req.Regime)
	if regime == "" {
		regime = domain.NewRegime
	}
	if regime != domain.NewRegime && regime != domain.OldRegime {
		writeError(w, http.StatusBadRequest, "regime must be NEW_REGIME or OLD_REGIME")
		return
	}
	if regime == domain.OldRegime && req.Form10IEARef == "" {
		writeError(w, http.StatusBadRequest, "Form 10-IEA reference required for old regime opt-out")
		return
	}

	idempKey := tenantID.String() + ":" + req.PAN + ":" + req.TaxYear + ":" + string(formType)

	f := &domain.ITRFiling{
		ID:             uuid.New(),
		TenantID:       tenantID,
		TaxpayerID:     taxpayerID,
		PAN:            req.PAN,
		TaxYear:        req.TaxYear,
		FormType:       formType,
		RegimeSelected: regime,
		Form10IEARef:   req.Form10IEARef,
		Status:         domain.StatusDraft,
		GrossIncome:    decimal.Zero,
		TotalDeductions: decimal.Zero,
		TaxableIncome:  decimal.Zero,
		TaxPayable:     decimal.Zero,
		TDSCredited:    decimal.Zero,
		AdvanceTaxPaid: decimal.Zero,
		SelfAssessmentTax: decimal.Zero,
		RefundDue:      decimal.Zero,
		BalancePayable: decimal.Zero,
		IdempotencyKey: idempKey,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.store.CreateFiling(r.Context(), tenantID, f); err != nil {
		log.Error().Err(err).Msg("create filing failed")
		writeError(w, http.StatusConflict, "filing already exists for this PAN + tax year + form type")
		return
	}
	writeJSON(w, http.StatusCreated, f)
}

func (h *Handlers) GetFiling(w http.ResponseWriter, r *http.Request) {
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
	f, err := h.store.GetFiling(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "filing not found")
		return
	}
	writeJSON(w, http.StatusOK, f)
}

func (h *Handlers) ListFilings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	limit, offset := pagination(r)
	taxYear := r.URL.Query().Get("tax_year")
	list, total, err := h.store.ListFilings(r.Context(), tenantID, taxYear, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("list filings failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"filings": list, "total": total})
}

func (h *Handlers) ComputeTax(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Salary        float64 `json:"salary"`
		HouseProperty float64 `json:"house_property"`
		CapitalGains  float64 `json:"capital_gains"`
		Business      float64 `json:"business"`
		OtherSources  float64 `json:"other_sources"`
		Regime        string  `json:"regime"`
		IsResident    bool    `json:"is_resident"`
		TDSCredits    float64 `json:"tds_credits"`
		AdvanceTax    float64 `json:"advance_tax"`
		Section80C    float64 `json:"section_80c"`
		Section80D    float64 `json:"section_80d"`
		Section24b    float64 `json:"section_24b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	regime := domain.RegimeType(req.Regime)
	if regime == "" {
		regime = domain.NewRegime
	}

	input := domain.TaxComputeInput{
		Income: domain.IncomeBreakdown{
			Salary:        decimal.NewFromFloat(req.Salary),
			HouseProperty: decimal.NewFromFloat(req.HouseProperty),
			CapitalGains:  decimal.NewFromFloat(req.CapitalGains),
			Business:      decimal.NewFromFloat(req.Business),
			OtherSources:  decimal.NewFromFloat(req.OtherSources),
		},
		Deductions: domain.DeductionBreakdown{
			Section80C: decimal.NewFromFloat(req.Section80C),
			Section80D: decimal.NewFromFloat(req.Section80D),
			Section24b: decimal.NewFromFloat(req.Section24b),
		},
		Regime:     regime,
		IsResident: req.IsResident,
		TDSCredits: decimal.NewFromFloat(req.TDSCredits),
		AdvanceTax: decimal.NewFromFloat(req.AdvanceTax),
	}

	result := domain.ComputeTax(input)
	writeJSON(w, http.StatusOK, result)
}

func (h *Handlers) AddIncomeEntry(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	filingID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid filing id")
		return
	}
	var req struct {
		Head        string  `json:"head"`
		SubHead     string  `json:"sub_head"`
		Section     string  `json:"section"`
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Exempt      bool    `json:"exempt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if domain.IsOldSectionRef(req.Section) {
		writeError(w, http.StatusBadRequest, "ITA 1961 section reference not accepted — use ITA 2025 equivalents (e.g., Section 202 instead of 115BAC)")
		return
	}

	e := &domain.IncomeEntry{
		ID:          uuid.New(),
		TenantID:    tenantID,
		FilingID:    filingID,
		Head:        domain.IncomeHead(req.Head),
		SubHead:     req.SubHead,
		Section:     req.Section,
		Description: req.Description,
		Amount:      decimal.NewFromFloat(req.Amount),
		Exempt:      req.Exempt,
		CreatedAt:   time.Now(),
	}
	if err := h.store.CreateIncomeEntry(r.Context(), tenantID, e); err != nil {
		log.Error().Err(err).Msg("create income entry failed")
		writeError(w, http.StatusInternalServerError, "failed to add income entry")
		return
	}
	writeJSON(w, http.StatusCreated, e)
}

func (h *Handlers) ListIncomeEntries(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	filingID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid filing id")
		return
	}
	entries, err := h.store.ListIncomeEntries(r.Context(), tenantID, filingID)
	if err != nil {
		log.Error().Err(err).Msg("list income entries failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func (h *Handlers) AddDeduction(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	filingID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid filing id")
		return
	}
	var req struct {
		Section  string  `json:"section"`
		Label    string  `json:"label"`
		Claimed  float64 `json:"claimed"`
		MaxLimit float64 `json:"max_limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claimed := decimal.NewFromFloat(req.Claimed)
	maxLim := decimal.NewFromFloat(req.MaxLimit)
	allowed := claimed
	if maxLim.IsPositive() && claimed.GreaterThan(maxLim) {
		allowed = maxLim
	}

	d := &domain.Deduction{
		ID:        uuid.New(),
		TenantID:  tenantID,
		FilingID:  filingID,
		Section:   req.Section,
		Label:     req.Label,
		Claimed:   claimed,
		Allowed:   allowed,
		MaxLimit:  maxLim,
		CreatedAt: time.Now(),
	}
	if err := h.store.CreateDeduction(r.Context(), tenantID, d); err != nil {
		log.Error().Err(err).Msg("create deduction failed")
		writeError(w, http.StatusInternalServerError, "failed to add deduction")
		return
	}
	writeJSON(w, http.StatusCreated, d)
}

func (h *Handlers) ListDeductions(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	filingID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid filing id")
		return
	}
	deds, err := h.store.ListDeductions(r.Context(), tenantID, filingID)
	if err != nil {
		log.Error().Err(err).Msg("list deductions failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, deds)
}

func (h *Handlers) GetTaxComputation(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	filingID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid filing id")
		return
	}
	tc, err := h.store.GetTaxComputation(r.Context(), tenantID, filingID)
	if err != nil {
		writeError(w, http.StatusNotFound, "tax computation not found")
		return
	}
	writeJSON(w, http.StatusOK, tc)
}

func (h *Handlers) AddTDSCredit(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	filingID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid filing id")
		return
	}
	var req struct {
		DeductorTAN  string  `json:"deductor_tan"`
		DeductorName string  `json:"deductor_name"`
		Section      string  `json:"section"`
		TDSAmount    float64 `json:"tds_amount"`
		GrossPayment float64 `json:"gross_payment"`
		TaxYear      string  `json:"tax_year"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c := &domain.TDSCredit{
		ID:           uuid.New(),
		TenantID:     tenantID,
		FilingID:     filingID,
		DeductorTAN:  req.DeductorTAN,
		DeductorName: req.DeductorName,
		Section:      req.Section,
		TDSAmount:    decimal.NewFromFloat(req.TDSAmount),
		GrossPayment: decimal.NewFromFloat(req.GrossPayment),
		TaxYear:      req.TaxYear,
		CreatedAt:    time.Now(),
	}
	if err := h.store.CreateTDSCredit(r.Context(), tenantID, c); err != nil {
		log.Error().Err(err).Msg("create tds credit failed")
		writeError(w, http.StatusInternalServerError, "failed to add TDS credit")
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *Handlers) ListTDSCredits(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantFrom(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing X-Tenant-Id")
		return
	}
	filingID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid filing id")
		return
	}
	credits, err := h.store.ListTDSCredits(r.Context(), tenantID, filingID)
	if err != nil {
		log.Error().Err(err).Msg("list tds credits failed")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, credits)
}

func (h *Handlers) ReconcileTDS(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AISEntries []domain.AISEntry       `json:"ais_entries"`
		TDSClaims  []domain.TDSCreditEntry `json:"tds_claims"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result := domain.ReconcileTDS(req.AISEntries, req.TDSClaims)
	writeJSON(w, http.StatusOK, result)
}

func (h *Handlers) CheckITR1Eligibility(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	assesseeType := domain.AssesseeType(q.Get("assessee_type"))
	residency := domain.ResidencyStatus(q.Get("residency"))
	totalIncome := decFromQuery(q.Get("total_income"))
	hpCount, _ := strconv.Atoi(q.Get("hp_count"))
	hasCapGains := q.Get("has_capital_gains_over_112a") == "true"
	ltcg112A := decFromQuery(q.Get("ltcg_112a"))
	hasBusiness := q.Get("has_business") == "true"
	hasForeignAssets := q.Get("has_foreign_assets") == "true"
	hasUnlistedEquity := q.Get("has_unlisted_equity") == "true"
	isDirector := q.Get("is_director") == "true"
	agriIncome := decFromQuery(q.Get("agricultural_income"))

	result := domain.CheckITR1Eligibility(assesseeType, residency, totalIncome, hpCount, hasCapGains, ltcg112A, hasBusiness, hasForeignAssets, hasUnlistedEquity, isDirector, agriIncome)
	writeJSON(w, http.StatusOK, result)
}

func (h *Handlers) CheckITR2Eligibility(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	assesseeType := domain.AssesseeType(q.Get("assessee_type"))
	hasBusiness := q.Get("has_business") == "true"
	result := domain.CheckITR2Eligibility(assesseeType, hasBusiness)
	writeJSON(w, http.StatusOK, result)
}

func (h *Handlers) CheckITR3Eligibility(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	assesseeType := domain.AssesseeType(q.Get("assessee_type"))
	result := domain.CheckITR3Eligibility(assesseeType)
	writeJSON(w, http.StatusOK, result)
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

func decFromQuery(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
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
