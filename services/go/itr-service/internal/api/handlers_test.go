package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/complai/complai/services/go/itr-service/internal/store"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter() (*Handlers, http.Handler) {
	s := store.NewMockStore()
	h := NewHandlers(s)
	return h, NewRouter(h)
}

func TestHealth(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestCreateTaxpayer(t *testing.T) {
	_, router := setupRouter()
	body := `{"pan":"ABCDE1234F","name":"Rahul Sharma","date_of_birth":"1990-01-15","assessee_type":"INDIVIDUAL","residency_status":"RESIDENT","aadhaar_linked":true,"email":"rahul@example.com","mobile":"9876543210"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/taxpayers", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "ABCDE1234F", data["pan"])
	assert.Equal(t, "INDIVIDUAL", data["assessee_type"])
}

func TestCreateTaxpayer_InvalidPAN(t *testing.T) {
	_, router := setupRouter()
	body := `{"pan":"SHORT","name":"Bad PAN","date_of_birth":"1990-01-15","assessee_type":"INDIVIDUAL","residency_status":"RESIDENT"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/taxpayers", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestCreateTaxpayer_MissingTenantID(t *testing.T) {
	_, router := setupRouter()
	body := `{"pan":"ABCDE1234F","name":"No Tenant","date_of_birth":"1990-01-15"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/taxpayers", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestCreateFiling_NewRegime(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	// Create taxpayer first
	tpBody := `{"pan":"ABCDE1234F","name":"Rahul","date_of_birth":"1990-01-15","assessee_type":"INDIVIDUAL","residency_status":"RESIDENT"}`
	tpReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/taxpayers", bytes.NewBufferString(tpBody))
	tpReq.Header.Set("X-Tenant-Id", tenantID)
	tpW := httptest.NewRecorder()
	router.ServeHTTP(tpW, tpReq)
	require.Equal(t, 201, tpW.Code)

	var tpResp map[string]interface{}
	json.Unmarshal(tpW.Body.Bytes(), &tpResp)
	taxpayerID := tpResp["data"].(map[string]interface{})["id"].(string)

	// Create filing
	filingBody, _ := json.Marshal(map[string]string{
		"taxpayer_id": taxpayerID,
		"pan":         "ABCDE1234F",
		"tax_year":    "2026-27",
		"form_type":   "ITR-1",
		"regime":      "NEW_REGIME",
	})
	fReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings", bytes.NewBuffer(filingBody))
	fReq.Header.Set("X-Tenant-Id", tenantID)
	fW := httptest.NewRecorder()
	router.ServeHTTP(fW, fReq)
	assert.Equal(t, 201, fW.Code)

	var fResp map[string]interface{}
	json.Unmarshal(fW.Body.Bytes(), &fResp)
	data := fResp["data"].(map[string]interface{})
	assert.Equal(t, "NEW_REGIME", data["regime_selected"])
	assert.Equal(t, "DRAFT", data["status"])
	assert.Equal(t, "2026-27", data["tax_year"])
}

func TestCreateFiling_OldRegime_RequiresForm10IEA(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	tpBody := `{"pan":"ABCDE1234F","name":"Rahul","date_of_birth":"1990-01-15","assessee_type":"INDIVIDUAL","residency_status":"RESIDENT"}`
	tpReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/taxpayers", bytes.NewBufferString(tpBody))
	tpReq.Header.Set("X-Tenant-Id", tenantID)
	tpW := httptest.NewRecorder()
	router.ServeHTTP(tpW, tpReq)
	require.Equal(t, 201, tpW.Code)

	var tpResp map[string]interface{}
	json.Unmarshal(tpW.Body.Bytes(), &tpResp)
	taxpayerID := tpResp["data"].(map[string]interface{})["id"].(string)

	// Old regime without Form 10-IEA → should fail
	filingBody, _ := json.Marshal(map[string]string{
		"taxpayer_id": taxpayerID,
		"pan":         "ABCDE1234F",
		"tax_year":    "2026-27",
		"form_type":   "ITR-1",
		"regime":      "OLD_REGIME",
	})
	fReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings", bytes.NewBuffer(filingBody))
	fReq.Header.Set("X-Tenant-Id", tenantID)
	fW := httptest.NewRecorder()
	router.ServeHTTP(fW, fReq)
	assert.Equal(t, 400, fW.Code)

	var errResp map[string]interface{}
	json.Unmarshal(fW.Body.Bytes(), &errResp)
	assert.Contains(t, errResp["error"].(string), "Form 10-IEA")
}

func TestCreateFiling_OldRegime_WithForm10IEA(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	tpBody := `{"pan":"ABCDE1234F","name":"Rahul","date_of_birth":"1990-01-15","assessee_type":"INDIVIDUAL","residency_status":"RESIDENT"}`
	tpReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/taxpayers", bytes.NewBufferString(tpBody))
	tpReq.Header.Set("X-Tenant-Id", tenantID)
	tpW := httptest.NewRecorder()
	router.ServeHTTP(tpW, tpReq)
	require.Equal(t, 201, tpW.Code)

	var tpResp map[string]interface{}
	json.Unmarshal(tpW.Body.Bytes(), &tpResp)
	taxpayerID := tpResp["data"].(map[string]interface{})["id"].(string)

	filingBody, _ := json.Marshal(map[string]string{
		"taxpayer_id":    taxpayerID,
		"pan":            "ABCDE1234F",
		"tax_year":       "2026-27",
		"form_type":      "ITR-1",
		"regime":         "OLD_REGIME",
		"form_10iea_ref": "10IEA-2026-ABC123",
	})
	fReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings", bytes.NewBuffer(filingBody))
	fReq.Header.Set("X-Tenant-Id", tenantID)
	fW := httptest.NewRecorder()
	router.ServeHTTP(fW, fReq)
	assert.Equal(t, 201, fW.Code)

	var fResp map[string]interface{}
	json.Unmarshal(fW.Body.Bytes(), &fResp)
	data := fResp["data"].(map[string]interface{})
	assert.Equal(t, "OLD_REGIME", data["regime_selected"])
	assert.Equal(t, "10IEA-2026-ABC123", data["form_10iea_ref"])
}

func TestComputeTax_Endpoint(t *testing.T) {
	_, router := setupRouter()
	body := `{"salary":1275000,"regime":"NEW_REGIME","is_resident":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/compute-tax", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "NEW_REGIME", data["regime"])
	assert.Equal(t, "60000", data["base_tax"])
	assert.Equal(t, "62400", data["gross_tax_payable"])
}

func TestComputeTax_DefaultsToNewRegime(t *testing.T) {
	_, router := setupRouter()
	body := `{"salary":1275000,"is_resident":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/compute-tax", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "NEW_REGIME", data["regime"])
}

func TestAddIncomeEntry_ITA2025Enforcement(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	// Create taxpayer and filing
	tpBody := `{"pan":"ABCDE1234F","name":"Rahul","date_of_birth":"1990-01-15","assessee_type":"INDIVIDUAL","residency_status":"RESIDENT"}`
	tpReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/taxpayers", bytes.NewBufferString(tpBody))
	tpReq.Header.Set("X-Tenant-Id", tenantID)
	tpW := httptest.NewRecorder()
	router.ServeHTTP(tpW, tpReq)
	require.Equal(t, 201, tpW.Code)

	var tpResp map[string]interface{}
	json.Unmarshal(tpW.Body.Bytes(), &tpResp)
	taxpayerID := tpResp["data"].(map[string]interface{})["id"].(string)

	filingBody, _ := json.Marshal(map[string]string{
		"taxpayer_id": taxpayerID, "pan": "ABCDE1234F", "tax_year": "2026-27", "form_type": "ITR-1",
	})
	fReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings", bytes.NewBuffer(filingBody))
	fReq.Header.Set("X-Tenant-Id", tenantID)
	fW := httptest.NewRecorder()
	router.ServeHTTP(fW, fReq)
	require.Equal(t, 201, fW.Code)

	var fResp map[string]interface{}
	json.Unmarshal(fW.Body.Bytes(), &fResp)
	filingID := fResp["data"].(map[string]interface{})["id"].(string)

	// Attempt to add income entry with old section ref "115BAC" → should be rejected
	incBody := `{"head":"SALARY","section":"115BAC","description":"old section ref","amount":500000}`
	incReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings/"+filingID+"/income", bytes.NewBufferString(incBody))
	incReq.Header.Set("X-Tenant-Id", tenantID)
	incW := httptest.NewRecorder()
	router.ServeHTTP(incW, incReq)
	assert.Equal(t, 400, incW.Code)

	var errResp map[string]interface{}
	json.Unmarshal(incW.Body.Bytes(), &errResp)
	assert.Contains(t, errResp["error"].(string), "ITA 1961")
	assert.Contains(t, errResp["error"].(string), "Section 202")
}

func TestAddIncomeEntry_ITA2025_Valid(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	tpBody := `{"pan":"ABCDE1234F","name":"Rahul","date_of_birth":"1990-01-15","assessee_type":"INDIVIDUAL","residency_status":"RESIDENT"}`
	tpReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/taxpayers", bytes.NewBufferString(tpBody))
	tpReq.Header.Set("X-Tenant-Id", tenantID)
	tpW := httptest.NewRecorder()
	router.ServeHTTP(tpW, tpReq)
	require.Equal(t, 201, tpW.Code)

	var tpResp map[string]interface{}
	json.Unmarshal(tpW.Body.Bytes(), &tpResp)
	taxpayerID := tpResp["data"].(map[string]interface{})["id"].(string)

	filingBody, _ := json.Marshal(map[string]string{
		"taxpayer_id": taxpayerID, "pan": "ABCDE1234F", "tax_year": "2026-27", "form_type": "ITR-1",
	})
	fReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings", bytes.NewBuffer(filingBody))
	fReq.Header.Set("X-Tenant-Id", tenantID)
	fW := httptest.NewRecorder()
	router.ServeHTTP(fW, fReq)
	require.Equal(t, 201, fW.Code)

	var fResp map[string]interface{}
	json.Unmarshal(fW.Body.Bytes(), &fResp)
	filingID := fResp["data"].(map[string]interface{})["id"].(string)

	// Valid ITA 2025 section ref
	incBody := `{"head":"SALARY","section":"392","description":"salary from employer","amount":1200000}`
	incReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings/"+filingID+"/income", bytes.NewBufferString(incBody))
	incReq.Header.Set("X-Tenant-Id", tenantID)
	incW := httptest.NewRecorder()
	router.ServeHTTP(incW, incReq)
	assert.Equal(t, 201, incW.Code)
}

func TestReconcileTDS_Endpoint(t *testing.T) {
	_, router := setupRouter()
	body := `{
		"ais_entries": [{"deductor_tan":"MUMB12345A","section":"392","amount":500000,"tds_amount":50000}],
		"tds_claims": [{"deductor_tan":"MUMB12345A","section":"392","gross_payment":500000,"tds_amount":50000}]
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/reconcile-tds", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "0", data["difference"])
}

func TestCheckITR1Eligibility_Endpoint(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/eligibility/itr1?assessee_type=INDIVIDUAL&residency=RESIDENT&total_income=3000000&hp_count=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, true, data["eligible"])
}

func TestCheckITR2Eligibility_Endpoint(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/eligibility/itr2?assessee_type=INDIVIDUAL&has_business=false", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestCheckITR3Eligibility_Endpoint(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/eligibility/itr3?assessee_type=HUF", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestListTaxpayers_Empty(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/taxpayers", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetTaxpayer_NotFound(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/taxpayers/"+uuid.New().String(), nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestGetTaxpayer_InvalidID(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/taxpayers/not-a-uuid", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestCreateFiling_InvalidFormType(t *testing.T) {
	_, router := setupRouter()
	body := `{"taxpayer_id":"` + uuid.New().String() + `","pan":"ABCDE1234F","tax_year":"2026-27","form_type":"ITR-99"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestCreateFiling_InvalidRegime(t *testing.T) {
	_, router := setupRouter()
	body := `{"taxpayer_id":"` + uuid.New().String() + `","pan":"ABCDE1234F","tax_year":"2026-27","form_type":"ITR-1","regime":"INVALID"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestListFilings_Empty(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/filings?tax_year=2026-27", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetFiling_NotFound(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/filings/"+uuid.New().String(), nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestGetTaxComputation_NotFound(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/filings/"+uuid.New().String()+"/computation", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestAddDeduction_Success(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	tpBody := `{"pan":"ABCDE1234F","name":"Rahul","date_of_birth":"1990-01-15","assessee_type":"INDIVIDUAL","residency_status":"RESIDENT"}`
	tpReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/taxpayers", bytes.NewBufferString(tpBody))
	tpReq.Header.Set("X-Tenant-Id", tenantID)
	tpW := httptest.NewRecorder()
	router.ServeHTTP(tpW, tpReq)
	require.Equal(t, 201, tpW.Code)
	var tpResp map[string]interface{}
	json.Unmarshal(tpW.Body.Bytes(), &tpResp)
	taxpayerID := tpResp["data"].(map[string]interface{})["id"].(string)

	filingBody, _ := json.Marshal(map[string]string{
		"taxpayer_id": taxpayerID, "pan": "ABCDE1234F", "tax_year": "2026-27", "form_type": "ITR-1",
	})
	fReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings", bytes.NewBuffer(filingBody))
	fReq.Header.Set("X-Tenant-Id", tenantID)
	fW := httptest.NewRecorder()
	router.ServeHTTP(fW, fReq)
	require.Equal(t, 201, fW.Code)
	var fResp map[string]interface{}
	json.Unmarshal(fW.Body.Bytes(), &fResp)
	filingID := fResp["data"].(map[string]interface{})["id"].(string)

	dedBody := `{"section":"VI-A","label":"Life insurance","claimed":200000,"max_limit":150000}`
	dedReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings/"+filingID+"/deductions", bytes.NewBufferString(dedBody))
	dedReq.Header.Set("X-Tenant-Id", tenantID)
	dedW := httptest.NewRecorder()
	router.ServeHTTP(dedW, dedReq)
	assert.Equal(t, 201, dedW.Code)

	var dedResp map[string]interface{}
	json.Unmarshal(dedW.Body.Bytes(), &dedResp)
	data := dedResp["data"].(map[string]interface{})
	assert.Equal(t, "150000", data["allowed"])

	// List deductions
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/itr/filings/"+filingID+"/deductions", nil)
	listReq.Header.Set("X-Tenant-Id", tenantID)
	listW := httptest.NewRecorder()
	router.ServeHTTP(listW, listReq)
	assert.Equal(t, 200, listW.Code)
}

func TestAddTDSCredit_Success(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	tpBody := `{"pan":"ABCDE1234F","name":"Rahul","date_of_birth":"1990-01-15","assessee_type":"INDIVIDUAL","residency_status":"RESIDENT"}`
	tpReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/taxpayers", bytes.NewBufferString(tpBody))
	tpReq.Header.Set("X-Tenant-Id", tenantID)
	tpW := httptest.NewRecorder()
	router.ServeHTTP(tpW, tpReq)
	require.Equal(t, 201, tpW.Code)
	var tpResp map[string]interface{}
	json.Unmarshal(tpW.Body.Bytes(), &tpResp)
	taxpayerID := tpResp["data"].(map[string]interface{})["id"].(string)

	filingBody, _ := json.Marshal(map[string]string{
		"taxpayer_id": taxpayerID, "pan": "ABCDE1234F", "tax_year": "2026-27", "form_type": "ITR-1",
	})
	fReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings", bytes.NewBuffer(filingBody))
	fReq.Header.Set("X-Tenant-Id", tenantID)
	fW := httptest.NewRecorder()
	router.ServeHTTP(fW, fReq)
	require.Equal(t, 201, fW.Code)
	var fResp map[string]interface{}
	json.Unmarshal(fW.Body.Bytes(), &fResp)
	filingID := fResp["data"].(map[string]interface{})["id"].(string)

	tdsBody := `{"deductor_tan":"MUMB12345A","deductor_name":"Infosys Ltd","section":"392","tds_amount":50000,"gross_payment":500000,"tax_year":"2026-27"}`
	tdsReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings/"+filingID+"/tds-credits", bytes.NewBufferString(tdsBody))
	tdsReq.Header.Set("X-Tenant-Id", tenantID)
	tdsW := httptest.NewRecorder()
	router.ServeHTTP(tdsW, tdsReq)
	assert.Equal(t, 201, tdsW.Code)

	// List TDS credits
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/itr/filings/"+filingID+"/tds-credits", nil)
	listReq.Header.Set("X-Tenant-Id", tenantID)
	listW := httptest.NewRecorder()
	router.ServeHTTP(listW, listReq)
	assert.Equal(t, 200, listW.Code)
}

func TestListIncomeEntries_InvalidFilingID(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/filings/not-uuid/income", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestAddDeduction_InvalidFilingID(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings/not-uuid/deductions", bytes.NewBufferString(`{}`))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestComputeTax_InvalidBody(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/compute-tax", bytes.NewBufferString(`not json`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestCreateTaxpayer_InvalidBody(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/taxpayers", bytes.NewBufferString(`not json`))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestCreateFiling_InvalidBody(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings", bytes.NewBufferString(`not json`))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestPagination_CustomLimitOffset(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/taxpayers?limit=10&offset=5", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestAddIncomeEntry_InvalidBody(t *testing.T) {
	_, router := setupRouter()
	filingID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings/"+filingID+"/income", bytes.NewBufferString(`not json`))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestAddDeduction_InvalidBody(t *testing.T) {
	_, router := setupRouter()
	filingID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings/"+filingID+"/deductions", bytes.NewBufferString(`not json`))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestAddTDSCredit_InvalidBody(t *testing.T) {
	_, router := setupRouter()
	filingID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings/"+filingID+"/tds-credits", bytes.NewBufferString(`not json`))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestReconcileTDS_InvalidBody(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/reconcile-tds", bytes.NewBufferString(`not json`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestCreateTaxpayer_InvalidDOB(t *testing.T) {
	_, router := setupRouter()
	body := `{"pan":"ABCDE1234F","name":"Test","date_of_birth":"not-a-date","assessee_type":"INDIVIDUAL"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/taxpayers", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestCreateFiling_InvalidTaxpayerID(t *testing.T) {
	_, router := setupRouter()
	body := `{"taxpayer_id":"not-uuid","pan":"ABCDE1234F","tax_year":"2026-27","form_type":"ITR-1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestGetFiling_InvalidID(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/filings/not-uuid", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestGetTaxComputation_InvalidFilingID(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/filings/not-uuid/computation", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestAddTDSCredit_InvalidFilingID(t *testing.T) {
	_, router := setupRouter()
	body := `{"deductor_tan":"TEST","section":"392","tds_amount":1000}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings/not-uuid/tds-credits", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestListTDSCredits_InvalidFilingID(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/filings/not-uuid/tds-credits", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestListDeductions_InvalidFilingID(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/filings/not-uuid/deductions", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestAddIncomeEntry_InvalidFilingID(t *testing.T) {
	_, router := setupRouter()
	body := `{"head":"SALARY","description":"test","amount":100}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/filings/not-uuid/income", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestListFilings_MissingTenant(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/filings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestGetFiling_MissingTenant(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/filings/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestReconcileAIS_Success(t *testing.T) {
	_, router := setupRouter()
	body := `{"ais":{"pan":"ABCDE1234F","tax_year":"2026-27","salary_income":1200000,"interest_income":50000},"books":{"salary_income":1200000,"interest_income":50000},"block_on_errors":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/reconcile-ais", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, false, data["has_errors"])
	assert.Equal(t, false, data["submission_blocked"])
}

func TestReconcileAIS_WithMismatches(t *testing.T) {
	_, router := setupRouter()
	body := `{"ais":{"pan":"ABCDE1234F","tax_year":"2026-27","salary_income":1200000,"interest_income":50000},"books":{"salary_income":1100000},"block_on_errors":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/reconcile-ais", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, true, data["has_errors"])
	assert.Equal(t, true, data["submission_blocked"])
}

func TestReconcileAIS_MissingPAN(t *testing.T) {
	_, router := setupRouter()
	body := `{"ais":{},"books":{}}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/reconcile-ais", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestReconcileAIS_InvalidBody(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/reconcile-ais", bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}
