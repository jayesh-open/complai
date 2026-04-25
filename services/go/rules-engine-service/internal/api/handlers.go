package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/rules-engine-service/internal/domain"
	"github.com/complai/complai/services/go/rules-engine-service/internal/store"
)

type Handlers struct {
	store store.Repository
}

func NewHandlers(s store.Repository) *Handlers {
	return &Handlers{store: s}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "rules-engine-service"})
}

func (h *Handlers) Evaluate(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.EvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	start := time.Now()

	var output domain.TaxOutput
	matched := true

	switch req.Category {
	case "tax_determination":
		output = evaluateTaxDetermination(req.Input)
	case "hsn_validation":
		output = evaluateHSNValidation(req.Input)
	case "tds_applicability":
		output = evaluateTDSApplicability(req.Input)
	default:
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("unknown category: %s", req.Category)})
		return
	}

	executionMs := int(time.Since(start).Milliseconds())

	result := domain.EvaluateResult{
		Category:     req.Category,
		Matched:      matched,
		MatchedRules: []string{req.Category},
		Output:       output,
		ExecutionMs:  executionMs,
	}

	// Log execution asynchronously (best-effort)
	go func() {
		inputJSON, _ := json.Marshal(req.Input)
		outputJSON, _ := json.Marshal(output)
		matchedJSON, _ := json.Marshal(result.MatchedRules)
		matchedStr := string(matchedJSON)
		inputStr := string(inputJSON)
		outputStr := string(outputJSON)

		execLog := &domain.RuleExecutionLog{
			InputData:       inputStr,
			MatchedRules:    &matchedStr,
			Output:          outputStr,
			ExecutionTimeMs: executionMs,
		}
		if err := h.store.CreateExecutionLog(r.Context(), tenantID, execLog); err != nil {
			log.Warn().Err(err).Msg("failed to log rule execution")
		}
	}()

	httputil.JSON(w, http.StatusOK, result)
}

func (h *Handlers) CreateRule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.CreateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	rule := &domain.Rule{
		Category:      req.Category,
		Name:          req.Name,
		Description:   req.Description,
		Version:       1,
		Priority:      req.Priority,
		Conditions:    req.Conditions,
		Actions:       req.Actions,
		EffectiveFrom: time.Now(),
	}

	if err := h.store.CreateRule(r.Context(), tenantID, rule); err != nil {
		log.Error().Err(err).Msg("create rule failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, rule)
}

func (h *Handlers) ListRules(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	category := r.URL.Query().Get("category")

	rules, err := h.store.ListRules(r.Context(), tenantID, category)
	if err != nil {
		log.Error().Err(err).Msg("list rules failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if rules == nil {
		rules = []domain.Rule{}
	}

	httputil.JSON(w, http.StatusOK, rules)
}

func (h *Handlers) GetRule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	ruleID, err := uuid.Parse(r.PathValue("ruleID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid rule_id"})
		return
	}

	rule, err := h.store.GetRule(r.Context(), tenantID, ruleID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "rule not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, rule)
}

// ---------------------------------------------------------------------------
// Built-in evaluation functions
// ---------------------------------------------------------------------------

// hsnRates is a built-in HSN code rate table (common codes).
var hsnRates = map[string]float64{
	"0101": 0,  // Live horses
	"1001": 5,  // Wheat
	"1006": 5,  // Rice
	"2201": 18, // Waters
	"3004": 12, // Medicaments
	"6101": 12, // Knitted overcoats
	"7308": 18, // Structures of iron/steel
	"8471": 18, // Computers
	"8517": 18, // Telephones
	"9401": 18, // Seats
	"9403": 18, // Furniture
	"9503": 12, // Toys
	"9954": 18, // Construction services
	"9971": 18, // Financial services
	"9983": 18, // IT services
	"9988": 18, // Manufacturing services
}

func lookupGSTRate(hsnCode string) float64 {
	// Try exact match first
	if rate, ok := hsnRates[hsnCode]; ok {
		return rate
	}
	// Try first 4 digits
	if len(hsnCode) >= 4 {
		if rate, ok := hsnRates[hsnCode[:4]]; ok {
			return rate
		}
	}
	// Try first 2 digits
	if len(hsnCode) >= 2 {
		if rate, ok := hsnRates[hsnCode[:2]]; ok {
			return rate
		}
	}
	return 18 // Default GST rate
}

func evaluateTaxDetermination(input domain.EvaluateInput) domain.TaxOutput {
	gstRate := lookupGSTRate(input.HSNCode)

	if input.SupplierState == input.PlaceOfSupply {
		// INTRA-STATE: split equally between CGST and SGST
		half := gstRate / 2
		return domain.TaxOutput{
			TaxType: "INTRA_STATE",
			GSTRate: &gstRate,
			CGST:    &half,
			SGST:    &half,
		}
	}
	// INTER-STATE: full IGST
	return domain.TaxOutput{
		TaxType: "INTER_STATE",
		GSTRate: &gstRate,
		IGST:    &gstRate,
	}
}

func evaluateHSNValidation(input domain.EvaluateInput) domain.TaxOutput {
	rate := lookupGSTRate(input.HSNCode)
	_, found := hsnRates[input.HSNCode]
	if !found && len(input.HSNCode) >= 4 {
		_, found = hsnRates[input.HSNCode[:4]]
	}
	valid := found || len(input.HSNCode) >= 4 // Valid if found or at least 4 digits
	desc := "Unknown HSN code"
	if found {
		desc = "Valid HSN code"
	}
	return domain.TaxOutput{
		HSNValid:       &valid,
		HSNDescription: &desc,
		GSTRate:        &rate,
	}
}

// tdsRule holds a TDS section's rate and threshold.
type tdsRule struct {
	rate      float64
	threshold float64
}

// tdsRules maps TDS sections to their rates and thresholds.
var tdsRules = map[string]tdsRule{
	"194C": {1, 30000},     // Contractor (individual/HUF 1%, others 2%)
	"194J": {10, 30000},    // Professional/Technical fees
	"194H": {5, 15000},     // Commission/Brokerage
	"194I": {10, 240000},   // Rent
	"194A": {10, 40000},    // Interest other than securities
	"194Q": {0.1, 5000000}, // Purchase of goods
}

func evaluateTDSApplicability(input domain.EvaluateInput) domain.TaxOutput {
	if rule, ok := tdsRules[input.Section]; ok {
		applicable := input.PaymentAmount > rule.threshold
		amount := 0.0
		if applicable {
			amount = input.PaymentAmount * rule.rate / 100
		}
		section := input.Section
		return domain.TaxOutput{
			TDSApplicable:  &applicable,
			TDSRate:        &rule.rate,
			TDSSection:     &section,
			TDSAmount:      &amount,
			ThresholdLimit:  &rule.threshold,
		}
	}

	notApplicable := false
	return domain.TaxOutput{
		TDSApplicable: &notApplicable,
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}
