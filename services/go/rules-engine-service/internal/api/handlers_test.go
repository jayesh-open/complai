package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/rules-engine-service/internal/domain"
)

// ---------------------------------------------------------------------------
// Mock store
// ---------------------------------------------------------------------------

type mockStore struct {
	createRuleFn       func(ctx context.Context, tenantID uuid.UUID, r *domain.Rule) error
	getRuleFn          func(ctx context.Context, tenantID uuid.UUID, ruleID uuid.UUID) (*domain.Rule, error)
	listRulesFn        func(ctx context.Context, tenantID uuid.UUID, category string) ([]domain.Rule, error)
	updateRuleFn       func(ctx context.Context, tenantID uuid.UUID, r *domain.Rule) error
	deleteRuleFn       func(ctx context.Context, tenantID uuid.UUID, ruleID uuid.UUID) error
	createExecLogFn    func(ctx context.Context, tenantID uuid.UUID, l *domain.RuleExecutionLog) error
}

func (m *mockStore) CreateRule(ctx context.Context, tenantID uuid.UUID, r *domain.Rule) error {
	if m.createRuleFn != nil {
		return m.createRuleFn(ctx, tenantID, r)
	}
	r.ID = uuid.New()
	r.TenantID = tenantID
	r.Status = "active"
	r.Version = 1
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	r.EffectiveFrom = time.Now()
	return nil
}

func (m *mockStore) GetRule(ctx context.Context, tenantID uuid.UUID, ruleID uuid.UUID) (*domain.Rule, error) {
	if m.getRuleFn != nil {
		return m.getRuleFn(ctx, tenantID, ruleID)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) ListRules(ctx context.Context, tenantID uuid.UUID, category string) ([]domain.Rule, error) {
	if m.listRulesFn != nil {
		return m.listRulesFn(ctx, tenantID, category)
	}
	return nil, nil
}

func (m *mockStore) UpdateRule(ctx context.Context, tenantID uuid.UUID, r *domain.Rule) error {
	if m.updateRuleFn != nil {
		return m.updateRuleFn(ctx, tenantID, r)
	}
	return nil
}

func (m *mockStore) DeleteRule(ctx context.Context, tenantID uuid.UUID, ruleID uuid.UUID) error {
	if m.deleteRuleFn != nil {
		return m.deleteRuleFn(ctx, tenantID, ruleID)
	}
	return nil
}

func (m *mockStore) CreateExecutionLog(ctx context.Context, tenantID uuid.UUID, l *domain.RuleExecutionLog) error {
	if m.createExecLogFn != nil {
		return m.createExecLogFn(ctx, tenantID, l)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func parseDataResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	var wrapper httputil.SuccessResponse
	wrapper.Data = target
	require.NoError(t, json.Unmarshal(body, &wrapper))
}

func newRequest(t *testing.T, method, path string, body interface{}, pathValues map[string]string) *http.Request {
	t.Helper()
	var req *http.Request
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	for k, v := range pathValues {
		req.SetPathValue(k, v)
	}
	return req
}

func requestWithTenant(req *http.Request, tenantID uuid.UUID) *http.Request {
	req.Header.Set("X-Tenant-Id", tenantID.String())
	return req
}

// ---------------------------------------------------------------------------
// Tests: Health
// ---------------------------------------------------------------------------

func TestHealth(t *testing.T) {
	h := NewHandlers(&mockStore{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	h.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "rules-engine-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: Evaluate — Tax Determination
// ---------------------------------------------------------------------------

func TestEvaluate_TaxDetermination_IntraState(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{
		Category: "tax_determination",
		Input: domain.EvaluateInput{
			SupplyType:    "B2B",
			PlaceOfSupply: "27",
			SupplierState: "27",
			HSNCode:       "8471",
			TaxableValue:  100000,
		},
	}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.EvaluateResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	assert.Equal(t, "tax_determination", result.Category)
	assert.True(t, result.Matched)
	assert.Equal(t, "INTRA_STATE", result.Output.TaxType)
	require.NotNil(t, result.Output.GSTRate)
	assert.Equal(t, 18.0, *result.Output.GSTRate)
	require.NotNil(t, result.Output.CGST)
	assert.Equal(t, 9.0, *result.Output.CGST)
	require.NotNil(t, result.Output.SGST)
	assert.Equal(t, 9.0, *result.Output.SGST)
	assert.Nil(t, result.Output.IGST)
}

func TestEvaluate_TaxDetermination_InterState(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{
		Category: "tax_determination",
		Input: domain.EvaluateInput{
			SupplyType:    "B2B",
			PlaceOfSupply: "29",
			SupplierState: "27",
			HSNCode:       "8471",
			TaxableValue:  100000,
		},
	}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.EvaluateResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	assert.Equal(t, "tax_determination", result.Category)
	assert.True(t, result.Matched)
	assert.Equal(t, "INTER_STATE", result.Output.TaxType)
	require.NotNil(t, result.Output.GSTRate)
	assert.Equal(t, 18.0, *result.Output.GSTRate)
	require.NotNil(t, result.Output.IGST)
	assert.Equal(t, 18.0, *result.Output.IGST)
	assert.Nil(t, result.Output.CGST)
	assert.Nil(t, result.Output.SGST)
}

func TestEvaluate_TaxDetermination_ZeroRateHSN(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{
		Category: "tax_determination",
		Input: domain.EvaluateInput{
			SupplyType:    "B2B",
			PlaceOfSupply: "27",
			SupplierState: "27",
			HSNCode:       "0101",
			TaxableValue:  50000,
		},
	}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.EvaluateResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	require.NotNil(t, result.Output.GSTRate)
	assert.Equal(t, 0.0, *result.Output.GSTRate)
	require.NotNil(t, result.Output.CGST)
	assert.Equal(t, 0.0, *result.Output.CGST)
	require.NotNil(t, result.Output.SGST)
	assert.Equal(t, 0.0, *result.Output.SGST)
}

func TestEvaluate_TaxDetermination_DefaultRate(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{
		Category: "tax_determination",
		Input: domain.EvaluateInput{
			SupplyType:    "B2B",
			PlaceOfSupply: "29",
			SupplierState: "27",
			HSNCode:       "9999",
			TaxableValue:  100000,
		},
	}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.EvaluateResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	require.NotNil(t, result.Output.GSTRate)
	assert.Equal(t, 18.0, *result.Output.GSTRate) // Default rate
}

// ---------------------------------------------------------------------------
// Tests: Evaluate — HSN Validation
// ---------------------------------------------------------------------------

func TestEvaluate_HSNValidation_ValidCode(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{
		Category: "hsn_validation",
		Input: domain.EvaluateInput{
			HSNCode: "8471",
		},
	}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.EvaluateResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	assert.Equal(t, "hsn_validation", result.Category)
	require.NotNil(t, result.Output.HSNValid)
	assert.True(t, *result.Output.HSNValid)
	require.NotNil(t, result.Output.HSNDescription)
	assert.Equal(t, "Valid HSN code", *result.Output.HSNDescription)
	require.NotNil(t, result.Output.GSTRate)
	assert.Equal(t, 18.0, *result.Output.GSTRate)
}

func TestEvaluate_HSNValidation_UnknownCode(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{
		Category: "hsn_validation",
		Input: domain.EvaluateInput{
			HSNCode: "5555",
		},
	}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.EvaluateResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	require.NotNil(t, result.Output.HSNValid)
	assert.True(t, *result.Output.HSNValid) // valid because len >= 4
	require.NotNil(t, result.Output.HSNDescription)
	assert.Equal(t, "Unknown HSN code", *result.Output.HSNDescription)
	require.NotNil(t, result.Output.GSTRate)
	assert.Equal(t, 18.0, *result.Output.GSTRate) // Default rate
}

func TestEvaluate_HSNValidation_ShortInvalidCode(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{
		Category: "hsn_validation",
		Input: domain.EvaluateInput{
			HSNCode: "99",
		},
	}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.EvaluateResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	require.NotNil(t, result.Output.HSNValid)
	assert.False(t, *result.Output.HSNValid) // not found, len < 4
}

// ---------------------------------------------------------------------------
// Tests: Evaluate — TDS Applicability
// ---------------------------------------------------------------------------

func TestEvaluate_TDS_194C_AboveThreshold(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{
		Category: "tds_applicability",
		Input: domain.EvaluateInput{
			Section:       "194C",
			PaymentAmount: 50000,
			DeducteeType:  "individual",
		},
	}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.EvaluateResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	assert.Equal(t, "tds_applicability", result.Category)
	require.NotNil(t, result.Output.TDSApplicable)
	assert.True(t, *result.Output.TDSApplicable)
	require.NotNil(t, result.Output.TDSRate)
	assert.Equal(t, 1.0, *result.Output.TDSRate)
	require.NotNil(t, result.Output.TDSSection)
	assert.Equal(t, "194C", *result.Output.TDSSection)
	require.NotNil(t, result.Output.TDSAmount)
	assert.Equal(t, 500.0, *result.Output.TDSAmount) // 50000 * 1 / 100
	require.NotNil(t, result.Output.ThresholdLimit)
	assert.Equal(t, 30000.0, *result.Output.ThresholdLimit)
}

func TestEvaluate_TDS_194C_BelowThreshold(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{
		Category: "tds_applicability",
		Input: domain.EvaluateInput{
			Section:       "194C",
			PaymentAmount: 20000,
			DeducteeType:  "individual",
		},
	}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.EvaluateResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	require.NotNil(t, result.Output.TDSApplicable)
	assert.False(t, *result.Output.TDSApplicable)
	require.NotNil(t, result.Output.TDSAmount)
	assert.Equal(t, 0.0, *result.Output.TDSAmount)
	require.NotNil(t, result.Output.ThresholdLimit)
	assert.Equal(t, 30000.0, *result.Output.ThresholdLimit)
}

func TestEvaluate_TDS_194J_Professional(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{
		Category: "tds_applicability",
		Input: domain.EvaluateInput{
			Section:       "194J",
			PaymentAmount: 100000,
			DeducteeType:  "company",
		},
	}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.EvaluateResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	require.NotNil(t, result.Output.TDSApplicable)
	assert.True(t, *result.Output.TDSApplicable)
	require.NotNil(t, result.Output.TDSRate)
	assert.Equal(t, 10.0, *result.Output.TDSRate)
	require.NotNil(t, result.Output.TDSAmount)
	assert.Equal(t, 10000.0, *result.Output.TDSAmount) // 100000 * 10 / 100
}

func TestEvaluate_TDS_UnknownSection(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{
		Category: "tds_applicability",
		Input: domain.EvaluateInput{
			Section:       "194Z",
			PaymentAmount: 100000,
		},
	}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.EvaluateResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	require.NotNil(t, result.Output.TDSApplicable)
	assert.False(t, *result.Output.TDSApplicable)
}

// ---------------------------------------------------------------------------
// Tests: Evaluate — Error cases
// ---------------------------------------------------------------------------

func TestEvaluate_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{Category: "tax_determination"}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	// No X-Tenant-Id header
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEvaluate_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodPost, "/v1/rules/evaluate", bytes.NewReader([]byte("not json")))
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid request", data["error"])
}

func TestEvaluate_UnknownCategory(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	body := domain.EvaluateRequest{
		Category: "unknown_category",
		Input:    domain.EvaluateInput{},
	}
	req := newRequest(t, http.MethodPost, "/v1/rules/evaluate", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.Evaluate(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Contains(t, data["error"], "unknown category")
}

// ---------------------------------------------------------------------------
// Tests: CreateRule
// ---------------------------------------------------------------------------

func TestCreateRule_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms)

	desc := "Test rule"
	body := domain.CreateRuleRequest{
		Category:    "tax_determination",
		Name:        "test-rule",
		Description: &desc,
		Priority:    100,
		Conditions:  `{"type":"B2B"}`,
		Actions:     `{"apply":"gst"}`,
	}
	req := newRequest(t, http.MethodPost, "/v1/rules", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.CreateRule(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var rule domain.Rule
	parseDataResponse(t, rec.Body.Bytes(), &rule)
	assert.Equal(t, "tax_determination", rule.Category)
	assert.Equal(t, "test-rule", rule.Name)
	assert.Equal(t, "active", rule.Status)
	assert.Equal(t, tenantID, rule.TenantID)
	assert.NotEqual(t, uuid.Nil, rule.ID)
}

func TestCreateRule_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodPost, "/v1/rules", bytes.NewReader([]byte("bad json")))
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.CreateRule(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid request", data["error"])
}

func TestCreateRule_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{})

	body := domain.CreateRuleRequest{Category: "tax_determination", Name: "test"}
	req := newRequest(t, http.MethodPost, "/v1/rules", body, nil)
	rec := httptest.NewRecorder()

	h.CreateRule(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateRule_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createRuleFn: func(_ context.Context, _ uuid.UUID, _ *domain.Rule) error {
			return errors.New("db down")
		},
	}
	h := NewHandlers(ms)

	body := domain.CreateRuleRequest{Category: "tax", Name: "test", Priority: 1, Conditions: "{}", Actions: "{}"}
	req := newRequest(t, http.MethodPost, "/v1/rules", body, nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.CreateRule(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "create failed", data["error"])
}

// ---------------------------------------------------------------------------
// Tests: ListRules
// ---------------------------------------------------------------------------

func TestListRules_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listRulesFn: func(_ context.Context, tid uuid.UUID, category string) ([]domain.Rule, error) {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, "tax_determination", category)
			return []domain.Rule{
				{ID: uuid.New(), TenantID: tenantID, Category: "tax_determination", Name: "rule-1", Status: "active"},
			}, nil
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/rules?category=tax_determination", nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.ListRules(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var rules []domain.Rule
	parseDataResponse(t, rec.Body.Bytes(), &rules)
	require.Len(t, rules, 1)
	assert.Equal(t, "rule-1", rules[0].Name)
}

func TestListRules_NilReturnsEmptyArray(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listRulesFn: func(_ context.Context, _ uuid.UUID, _ string) ([]domain.Rule, error) {
			return nil, nil
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/rules", nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.ListRules(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var rules []domain.Rule
	parseDataResponse(t, rec.Body.Bytes(), &rules)
	assert.NotNil(t, rules)
	assert.Len(t, rules, 0)
}

func TestListRules_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/rules", nil)
	rec := httptest.NewRecorder()

	h.ListRules(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListRules_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listRulesFn: func(_ context.Context, _ uuid.UUID, _ string) ([]domain.Rule, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/rules", nil)
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.ListRules(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "internal error", data["error"])
}

// ---------------------------------------------------------------------------
// Tests: GetRule
// ---------------------------------------------------------------------------

func TestGetRule_Success(t *testing.T) {
	tenantID := uuid.New()
	ruleID := uuid.New()
	ms := &mockStore{
		getRuleFn: func(_ context.Context, tid uuid.UUID, rid uuid.UUID) (*domain.Rule, error) {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, ruleID, rid)
			return &domain.Rule{
				ID:       ruleID,
				TenantID: tenantID,
				Category: "tax_determination",
				Name:     "rule-1",
				Status:   "active",
			}, nil
		},
	}
	h := NewHandlers(ms)

	req := newRequest(t, http.MethodGet, "/v1/rules/"+ruleID.String(), nil, map[string]string{
		"ruleID": ruleID.String(),
	})
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.GetRule(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var rule domain.Rule
	parseDataResponse(t, rec.Body.Bytes(), &rule)
	assert.Equal(t, ruleID, rule.ID)
	assert.Equal(t, "rule-1", rule.Name)
}

func TestGetRule_InvalidRuleID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := newRequest(t, http.MethodGet, "/v1/rules/bad", nil, map[string]string{
		"ruleID": "bad",
	})
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.GetRule(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid rule_id", data["error"])
}

func TestGetRule_NotFound(t *testing.T) {
	tenantID := uuid.New()
	ruleID := uuid.New()
	ms := &mockStore{
		getRuleFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.Rule, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewHandlers(ms)

	req := newRequest(t, http.MethodGet, "/v1/rules/"+ruleID.String(), nil, map[string]string{
		"ruleID": ruleID.String(),
	})
	requestWithTenant(req, tenantID)
	rec := httptest.NewRecorder()

	h.GetRule(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetRule_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{})

	ruleID := uuid.New()
	req := newRequest(t, http.MethodGet, "/v1/rules/"+ruleID.String(), nil, map[string]string{
		"ruleID": ruleID.String(),
	})
	rec := httptest.NewRecorder()

	h.GetRule(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: NewRouter
// ---------------------------------------------------------------------------

func TestNewRouter(t *testing.T) {
	ms := &mockStore{}
	r := NewRouter(ms)
	require.NotNil(t, r)

	// Verify health endpoint is reachable through the router
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify ping heartbeat endpoint
	req = httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: tenantIDFromRequest
// ---------------------------------------------------------------------------

func TestTenantIDFromRequest_Valid(t *testing.T) {
	expected := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-Id", expected.String())

	got, err := tenantIDFromRequest(req)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestTenantIDFromRequest_Missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := tenantIDFromRequest(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing X-Tenant-Id")
}

func TestTenantIDFromRequest_Invalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-Id", "not-uuid")

	_, err := tenantIDFromRequest(req)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Tests: lookupGSTRate
// ---------------------------------------------------------------------------

func TestLookupGSTRate_ExactMatch(t *testing.T) {
	assert.Equal(t, 18.0, lookupGSTRate("8471"))
	assert.Equal(t, 12.0, lookupGSTRate("3004"))
	assert.Equal(t, 5.0, lookupGSTRate("1001"))
	assert.Equal(t, 0.0, lookupGSTRate("0101"))
}

func TestLookupGSTRate_PrefixMatch(t *testing.T) {
	// 8-digit code with first 4 matching
	assert.Equal(t, 18.0, lookupGSTRate("84710000"))
}

func TestLookupGSTRate_Default(t *testing.T) {
	assert.Equal(t, 18.0, lookupGSTRate("9999"))
}

// ---------------------------------------------------------------------------
// Tests: evaluateTaxDetermination (unit)
// ---------------------------------------------------------------------------

func TestEvaluateTaxDetermination_IntraState(t *testing.T) {
	result := evaluateTaxDetermination(domain.EvaluateInput{
		SupplierState: "27",
		PlaceOfSupply: "27",
		HSNCode:       "9983",
	})
	assert.Equal(t, "INTRA_STATE", result.TaxType)
	require.NotNil(t, result.CGST)
	assert.Equal(t, 9.0, *result.CGST)
	require.NotNil(t, result.SGST)
	assert.Equal(t, 9.0, *result.SGST)
	assert.Nil(t, result.IGST)
}

func TestEvaluateTaxDetermination_InterState(t *testing.T) {
	result := evaluateTaxDetermination(domain.EvaluateInput{
		SupplierState: "27",
		PlaceOfSupply: "29",
		HSNCode:       "9983",
	})
	assert.Equal(t, "INTER_STATE", result.TaxType)
	require.NotNil(t, result.IGST)
	assert.Equal(t, 18.0, *result.IGST)
	assert.Nil(t, result.CGST)
	assert.Nil(t, result.SGST)
}

// ---------------------------------------------------------------------------
// Tests: evaluateHSNValidation (unit)
// ---------------------------------------------------------------------------

func TestEvaluateHSNValidation_Known(t *testing.T) {
	result := evaluateHSNValidation(domain.EvaluateInput{HSNCode: "8471"})
	require.NotNil(t, result.HSNValid)
	assert.True(t, *result.HSNValid)
	require.NotNil(t, result.HSNDescription)
	assert.Equal(t, "Valid HSN code", *result.HSNDescription)
}

func TestEvaluateHSNValidation_UnknownButLong(t *testing.T) {
	result := evaluateHSNValidation(domain.EvaluateInput{HSNCode: "5555"})
	require.NotNil(t, result.HSNValid)
	assert.True(t, *result.HSNValid) // >= 4 digits
	require.NotNil(t, result.HSNDescription)
	assert.Equal(t, "Unknown HSN code", *result.HSNDescription)
}

func TestEvaluateHSNValidation_Short(t *testing.T) {
	result := evaluateHSNValidation(domain.EvaluateInput{HSNCode: "55"})
	require.NotNil(t, result.HSNValid)
	assert.False(t, *result.HSNValid) // not found and < 4 digits
}

// ---------------------------------------------------------------------------
// Tests: evaluateTDSApplicability (unit)
// ---------------------------------------------------------------------------

func TestEvaluateTDSApplicability_Above(t *testing.T) {
	result := evaluateTDSApplicability(domain.EvaluateInput{
		Section:       "194C",
		PaymentAmount: 50000,
	})
	require.NotNil(t, result.TDSApplicable)
	assert.True(t, *result.TDSApplicable)
	require.NotNil(t, result.TDSAmount)
	assert.Equal(t, 500.0, *result.TDSAmount)
}

func TestEvaluateTDSApplicability_Below(t *testing.T) {
	result := evaluateTDSApplicability(domain.EvaluateInput{
		Section:       "194C",
		PaymentAmount: 20000,
	})
	require.NotNil(t, result.TDSApplicable)
	assert.False(t, *result.TDSApplicable)
	require.NotNil(t, result.TDSAmount)
	assert.Equal(t, 0.0, *result.TDSAmount)
}

func TestEvaluateTDSApplicability_UnknownSection(t *testing.T) {
	result := evaluateTDSApplicability(domain.EvaluateInput{
		Section:       "999X",
		PaymentAmount: 100000,
	})
	require.NotNil(t, result.TDSApplicable)
	assert.False(t, *result.TDSApplicable)
}

func TestEvaluateTDSApplicability_194Q(t *testing.T) {
	result := evaluateTDSApplicability(domain.EvaluateInput{
		Section:       "194Q",
		PaymentAmount: 10000000,
	})
	require.NotNil(t, result.TDSApplicable)
	assert.True(t, *result.TDSApplicable)
	require.NotNil(t, result.TDSRate)
	assert.Equal(t, 0.1, *result.TDSRate)
	require.NotNil(t, result.TDSAmount)
	assert.Equal(t, 10000.0, *result.TDSAmount) // 10000000 * 0.1 / 100
	require.NotNil(t, result.ThresholdLimit)
	assert.Equal(t, 5000000.0, *result.ThresholdLimit)
}
