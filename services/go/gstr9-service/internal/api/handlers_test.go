package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/complai/complai/services/go/gstr9-service/internal/domain"
	"github.com/complai/complai/services/go/gstr9-service/internal/store"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testTenant = "11111111-1111-1111-1111-111111111111"
const otherTenant = "22222222-2222-2222-2222-222222222222"

func setupTest() (*Handlers, *store.MockStore, *http.ServeMux) {
	ms := store.NewMockStore()
	h := NewHandlers(ms, "http://localhost:8093")
	router := NewRouter(h)
	mux := http.NewServeMux()
	mux.Handle("/", router)
	return h, ms, mux
}

func TestCreateAnnualReturn_Success(t *testing.T) {
	_, _, mux := setupTest()
	body := `{"gstin":"27AABCU9603R1ZM","financial_year":"2025-26"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gstr9/annual-return", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "27AABCU9603R1ZM", data["gstin"])
	assert.Equal(t, "2025-26", data["financial_year"])
	assert.Equal(t, "draft", data["status"])
}

func TestCreateAnnualReturn_DuplicateReturnsConflict(t *testing.T) {
	_, _, mux := setupTest()
	body := `{"gstin":"27AABCU9603R1ZM","financial_year":"2025-26"}`

	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/gstr9/annual-return", bytes.NewBufferString(body))
	req1.Header.Set("X-Tenant-Id", testTenant)
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/gstr9/annual-return", bytes.NewBufferString(body))
	req2.Header.Set("X-Tenant-Id", testTenant)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestCreateAnnualReturn_InvalidGSTIN(t *testing.T) {
	_, _, mux := setupTest()
	body := `{"gstin":"BADGSTIN","financial_year":"2025-26"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gstr9/annual-return", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "GSTIN")
}

func TestCreateAnnualReturn_InvalidFY(t *testing.T) {
	_, _, mux := setupTest()
	body := `{"gstin":"27AABCU9603R1ZM","financial_year":"2025"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gstr9/annual-return", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateAnnualReturn_MissingTenant(t *testing.T) {
	_, _, mux := setupTest()
	body := `{"gstin":"27AABCU9603R1ZM","financial_year":"2025-26"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gstr9/annual-return", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateAnnualReturn_InvalidBody(t *testing.T) {
	_, _, mux := setupTest()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gstr9/annual-return", bytes.NewBufferString("not json"))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func createFiling(t *testing.T, mux *http.ServeMux) string {
	t.Helper()
	body := `{"gstin":"27AABCU9603R1ZM","financial_year":"2025-26"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gstr9/annual-return", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	return data["id"].(string)
}

func TestGetAnnualReturn_Success(t *testing.T) {
	_, _, mux := setupTest()
	id := createFiling(t, mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/gstr9/annual-return/"+id, nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, id, data["id"])
}

func TestGetAnnualReturn_NotFound(t *testing.T) {
	_, _, mux := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/gstr9/annual-return/"+uuid.New().String(), nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAnnualReturn_InvalidID(t *testing.T) {
	_, _, mux := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/gstr9/annual-return/not-a-uuid", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAnnualReturn_CrossTenantReturns404(t *testing.T) {
	_, _, mux := setupTest()
	id := createFiling(t, mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/gstr9/annual-return/"+id, nil)
	req.Header.Set("X-Tenant-Id", otherTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "cross-tenant access should return 404")
}

func TestListAnnualReturns_Empty(t *testing.T) {
	_, _, mux := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/gstr9/annual-return", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListAnnualReturns_WithData(t *testing.T) {
	_, _, mux := setupTest()
	createFiling(t, mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/gstr9/annual-return?fy=2025-26", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
}

func TestListAnnualReturns_CrossTenantReturnsZero(t *testing.T) {
	_, _, mux := setupTest()
	createFiling(t, mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/gstr9/annual-return", nil)
	req.Header.Set("X-Tenant-Id", otherTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(0), data["total"])
}

func TestSaveAnnualReturn_Success(t *testing.T) {
	_, _, mux := setupTest()
	id := createFiling(t, mux)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/gstr9/annual-return/"+id+"/save", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "saved", data["status"])
}

func TestSaveAnnualReturn_AlreadySavedRejectsSave(t *testing.T) {
	_, _, mux := setupTest()
	id := createFiling(t, mux)

	req1 := httptest.NewRequest(http.MethodPut, "/api/v1/gstr9/annual-return/"+id+"/save", nil)
	req1.Header.Set("X-Tenant-Id", testTenant)
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	req2 := httptest.NewRequest(http.MethodPut, "/api/v1/gstr9/annual-return/"+id+"/save", nil)
	req2.Header.Set("X-Tenant-Id", testTenant)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code, "already-saved filing should reject re-save")
}

func TestSaveAnnualReturn_NotFound(t *testing.T) {
	_, _, mux := setupTest()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/gstr9/annual-return/"+uuid.New().String()+"/save", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSaveAnnualReturn_CrossTenantReturns404(t *testing.T) {
	_, _, mux := setupTest()
	id := createFiling(t, mux)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/gstr9/annual-return/"+id+"/save", nil)
	req.Header.Set("X-Tenant-Id", otherTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTableData_NotFoundWithoutAggregation(t *testing.T) {
	_, _, mux := setupTest()
	id := createFiling(t, mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/gstr9/annual-return/"+id+"/table/4A", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "no table data before aggregation")
}

func TestAggregateAndGetTable(t *testing.T) {
	_, _, mux := setupTest()
	id := createFiling(t, mux)

	monthsData := buildTestMonths()
	aggBody, _ := json.Marshal(map[string]interface{}{"months": monthsData})
	aggReq := httptest.NewRequest(http.MethodPost, "/api/v1/gstr9/annual-return/"+id+"/aggregate", bytes.NewReader(aggBody))
	aggReq.Header.Set("X-Tenant-Id", testTenant)
	aggW := httptest.NewRecorder()
	mux.ServeHTTP(aggW, aggReq)
	assert.Equal(t, http.StatusOK, aggW.Code)

	var aggResp map[string]interface{}
	json.Unmarshal(aggW.Body.Bytes(), &aggResp)
	data := aggResp["data"].(map[string]interface{})
	assert.Equal(t, float64(27), data["tables"])

	req := httptest.NewRequest(http.MethodGet, "/api/v1/gstr9/annual-return/"+id+"/table/4A", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAggregateAnnualReturn_EmptyMonths(t *testing.T) {
	_, _, mux := setupTest()
	id := createFiling(t, mux)

	body := `{"months":[]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gstr9/annual-return/"+id+"/aggregate", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAggregateAnnualReturn_InvalidBody(t *testing.T) {
	_, _, mux := setupTest()
	id := createFiling(t, mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/gstr9/annual-return/"+id+"/aggregate", bytes.NewBufferString("bad"))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAggregateAnnualReturn_NotFound(t *testing.T) {
	_, _, mux := setupTest()
	body := `{"months":[{"return_period":"202504"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gstr9/annual-return/"+uuid.New().String()+"/aggregate", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTableData_CrossTenantReturns404(t *testing.T) {
	_, ms, mux := setupTest()
	id := createFiling(t, mux)

	filingID := uuid.MustParse(id)
	tenantID := uuid.MustParse(testTenant)
	td := &domain.GSTR9TableData{
		ID: uuid.New(), TenantID: tenantID, FilingID: filingID,
		PartNumber: 1, TableNumber: "4A", Description: "test",
		TaxableValue: decimal.NewFromInt(100), CGST: decimal.NewFromInt(9),
		SGST: decimal.NewFromInt(9), IGST: decimal.NewFromInt(18), Cess: decimal.Zero,
	}
	ms.CreateTableData(nil, tenantID, td)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/gstr9/annual-return/"+id+"/table/4A", nil)
	req.Header.Set("X-Tenant-Id", otherTenant)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHealthEndpoint(t *testing.T) {
	_, _, mux := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPagination_Defaults(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	limit, offset := pagination(req)
	assert.Equal(t, 50, limit)
	assert.Equal(t, 0, offset)
}

func TestPagination_Custom(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test?limit=10&offset=20", nil)
	limit, offset := pagination(req)
	assert.Equal(t, 10, limit)
	assert.Equal(t, 20, offset)
}

func TestPagination_MaxLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test?limit=500", nil)
	limit, _ := pagination(req)
	assert.Equal(t, 50, limit, "limit > 200 should fall back to default")
}

func buildTestMonths() []map[string]interface{} {
	months := make([]map[string]interface{}, 0, 12)
	periods := domain.ReturnPeriodsForFY("2025-26")
	for _, p := range periods {
		months = append(months, map[string]interface{}{
			"ReturnPeriod": p,
			"Outward": map[string]interface{}{
				"taxable_value": 500000, "cgst": 45000, "sgst": 45000, "igst": 25000, "cess": 5000,
			},
			"Inward": map[string]interface{}{
				"taxable_value": 300000, "cgst": 27000, "sgst": 27000, "igst": 12000, "cess": 1500,
			},
			"ITC": map[string]interface{}{
				"cgst": 17500, "sgst": 17500, "igst": 10000, "cess": 5000,
			},
			"TaxPaid": map[string]interface{}{
				"taxable_value": 0, "cgst": 12000, "sgst": 12000, "igst": 4500, "cess": 1500,
			},
		})
	}
	return months
}
