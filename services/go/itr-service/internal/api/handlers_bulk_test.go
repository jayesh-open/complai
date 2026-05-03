package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBulkBatch(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	body := `{"tax_year":"2026-27","employer_tan":"DELC12345A","employer_name":"Acme Corp"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "PENDING", data["status"])
	assert.Equal(t, "DELC12345A", data["employer_tan"])
}

func TestCreateBulkBatch_MissingFields(t *testing.T) {
	_, router := setupRouter()
	body := `{"tax_year":"2026-27"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestCreateBulkBatch_MissingTenant(t *testing.T) {
	_, router := setupRouter()
	body := `{"tax_year":"2026-27","employer_tan":"DELC12345A","employer_name":"Acme"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestGetBulkBatch(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	body := `{"tax_year":"2026-27","employer_tan":"MUMX99999B","employer_name":"Test Ltd"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString(body))
	createReq.Header.Set("X-Tenant-Id", tenantID)
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, 201, createW.Code)

	var createResp map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createResp)
	batchID := createResp["data"].(map[string]interface{})["id"].(string)

	getReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/itr/bulk/batches/%s", batchID), nil)
	getReq.Header.Set("X-Tenant-Id", tenantID)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	assert.Equal(t, 200, getW.Code)
}

func TestGetBulkBatch_NotFound(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/itr/bulk/batches/%s", uuid.New()), nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestListBulkBatches(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	body := `{"tax_year":"2026-27","employer_tan":"AAAA11111A","employer_name":"List Test"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString(body))
	createReq.Header.Set("X-Tenant-Id", tenantID)
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, 201, createW.Code)

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/itr/bulk/batches", nil)
	listReq.Header.Set("X-Tenant-Id", tenantID)
	listW := httptest.NewRecorder()
	router.ServeHTTP(listW, listReq)

	assert.Equal(t, 200, listW.Code)
	var listResp map[string]interface{}
	json.Unmarshal(listW.Body.Bytes(), &listResp)
	data := listResp["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
}

func TestAddBulkEmployee(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	batchBody := `{"tax_year":"2026-27","employer_tan":"BBBB22222B","employer_name":"Emp Test"}`
	bReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString(batchBody))
	bReq.Header.Set("X-Tenant-Id", tenantID)
	bW := httptest.NewRecorder()
	router.ServeHTTP(bW, bReq)
	require.Equal(t, 201, bW.Code)

	var bResp map[string]interface{}
	json.Unmarshal(bW.Body.Bytes(), &bResp)
	batchID := bResp["data"].(map[string]interface{})["id"].(string)

	empBody := `{"pan":"EMPTE1234F","name":"John Doe","email":"john@test.com","gross_salary":800000,"tds_deducted":50000}`
	eReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/itr/bulk/batches/%s/employees", batchID), bytes.NewBufferString(empBody))
	eReq.Header.Set("X-Tenant-Id", tenantID)
	eW := httptest.NewRecorder()
	router.ServeHTTP(eW, eReq)

	assert.Equal(t, 201, eW.Code)
	var eResp map[string]interface{}
	json.Unmarshal(eW.Body.Bytes(), &eResp)
	data := eResp["data"].(map[string]interface{})
	assert.Equal(t, "EMPTE1234F", data["pan"])
	assert.Equal(t, "PENDING_REVIEW", data["status"])
}

func TestAddBulkEmployee_InvalidPAN(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	batchBody := `{"tax_year":"2026-27","employer_tan":"CCCC33333C","employer_name":"PAN Test"}`
	bReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString(batchBody))
	bReq.Header.Set("X-Tenant-Id", tenantID)
	bW := httptest.NewRecorder()
	router.ServeHTTP(bW, bReq)
	require.Equal(t, 201, bW.Code)

	var bResp map[string]interface{}
	json.Unmarshal(bW.Body.Bytes(), &bResp)
	batchID := bResp["data"].(map[string]interface{})["id"].(string)

	empBody := `{"pan":"SHORT","name":"Bad","email":"bad@test.com","gross_salary":100000,"tds_deducted":5000}`
	eReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/itr/bulk/batches/%s/employees", batchID), bytes.NewBufferString(empBody))
	eReq.Header.Set("X-Tenant-Id", tenantID)
	eW := httptest.NewRecorder()
	router.ServeHTTP(eW, eReq)

	assert.Equal(t, 400, eW.Code)
}

func TestProcessBulkBatch(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	batchBody := `{"tax_year":"2026-27","employer_tan":"DDDD44444D","employer_name":"Process Test"}`
	bReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString(batchBody))
	bReq.Header.Set("X-Tenant-Id", tenantID)
	bW := httptest.NewRecorder()
	router.ServeHTTP(bW, bReq)
	require.Equal(t, 201, bW.Code)

	var bResp map[string]interface{}
	json.Unmarshal(bW.Body.Bytes(), &bResp)
	batchID := bResp["data"].(map[string]interface{})["id"].(string)

	empBody := `{"pan":"PROC01234F","name":"Process Emp","email":"proc@test.com","gross_salary":1000000,"tds_deducted":80000}`
	eReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/itr/bulk/batches/%s/employees", batchID), bytes.NewBufferString(empBody))
	eReq.Header.Set("X-Tenant-Id", tenantID)
	eW := httptest.NewRecorder()
	router.ServeHTTP(eW, eReq)
	require.Equal(t, 201, eW.Code)

	processBody := `{"ais":{"pan":"PROC01234F","tax_year":"2026-27","salary_income":"1000000","tds_entries":[{"deductor_tan":"EMPLOYER","section":"392","tds_amount":"80000"}]}}`
	pReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/itr/bulk/batches/%s/process", batchID), bytes.NewBufferString(processBody))
	pReq.Header.Set("X-Tenant-Id", tenantID)
	pW := httptest.NewRecorder()
	router.ServeHTTP(pW, pReq)

	assert.Equal(t, 200, pW.Code)
	var pResp map[string]interface{}
	json.Unmarshal(pW.Body.Bytes(), &pResp)
	data := pResp["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["processed"])
}

func TestProcessBulkBatch_NoEmployees(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	batchBody := `{"tax_year":"2026-27","employer_tan":"EEEE55555E","employer_name":"Empty Batch"}`
	bReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString(batchBody))
	bReq.Header.Set("X-Tenant-Id", tenantID)
	bW := httptest.NewRecorder()
	router.ServeHTTP(bW, bReq)
	require.Equal(t, 201, bW.Code)

	var bResp map[string]interface{}
	json.Unmarshal(bW.Body.Bytes(), &bResp)
	batchID := bResp["data"].(map[string]interface{})["id"].(string)

	processBody := `{"ais":{"pan":"NOBODY0000X","tax_year":"2026-27"}}`
	pReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/itr/bulk/batches/%s/process", batchID), bytes.NewBufferString(processBody))
	pReq.Header.Set("X-Tenant-Id", tenantID)
	pW := httptest.NewRecorder()
	router.ServeHTTP(pW, pReq)

	assert.Equal(t, 400, pW.Code)
}

func TestProcessBulkBatch_NotFound(t *testing.T) {
	_, router := setupRouter()
	processBody := `{"ais":{"pan":"X","tax_year":"2026-27"}}`
	pReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/itr/bulk/batches/%s/process", uuid.New()), bytes.NewBufferString(processBody))
	pReq.Header.Set("X-Tenant-Id", uuid.New().String())
	pW := httptest.NewRecorder()
	router.ServeHTTP(pW, pReq)
	assert.Equal(t, 404, pW.Code)
}

func TestGetBulkBatch_InvalidID(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/bulk/batches/not-a-uuid", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestCreateBulkBatch_InvalidBody(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString("{bad"))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestAddBulkEmployee_BatchNotFound(t *testing.T) {
	_, router := setupRouter()
	empBody := `{"pan":"ABCDE1234F","name":"Nobody","email":"x@x.com","gross_salary":100000,"tds_deducted":5000}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/itr/bulk/batches/%s/employees", uuid.New()), bytes.NewBufferString(empBody))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestAddBulkEmployee_InvalidBody(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	batchBody := `{"tax_year":"2026-27","employer_tan":"ZZZZ99999Z","employer_name":"Body Test"}`
	bReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString(batchBody))
	bReq.Header.Set("X-Tenant-Id", tenantID)
	bW := httptest.NewRecorder()
	router.ServeHTTP(bW, bReq)
	require.Equal(t, 201, bW.Code)

	var bResp map[string]interface{}
	json.Unmarshal(bW.Body.Bytes(), &bResp)
	batchID := bResp["data"].(map[string]interface{})["id"].(string)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/itr/bulk/batches/%s/employees", batchID), bytes.NewBufferString("{bad"))
	req.Header.Set("X-Tenant-Id", tenantID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestAddBulkEmployee_InvalidBatchID(t *testing.T) {
	_, router := setupRouter()
	empBody := `{"pan":"ABCDE1234F","name":"X","email":"x@x.com","gross_salary":100000,"tds_deducted":5000}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches/not-uuid/employees", bytes.NewBufferString(empBody))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestListBulkBatches_MissingTenant(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/bulk/batches", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestListBulkEmployees_InvalidBatchID(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/itr/bulk/batches/not-uuid/employees", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestProcessBulkBatch_InvalidBody(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	batchBody := `{"tax_year":"2026-27","employer_tan":"INVA11111A","employer_name":"Inv Body"}`
	bReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString(batchBody))
	bReq.Header.Set("X-Tenant-Id", tenantID)
	bW := httptest.NewRecorder()
	router.ServeHTTP(bW, bReq)
	require.Equal(t, 201, bW.Code)

	var bResp map[string]interface{}
	json.Unmarshal(bW.Body.Bytes(), &bResp)
	batchID := bResp["data"].(map[string]interface{})["id"].(string)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/itr/bulk/batches/%s/process", batchID), bytes.NewBufferString("{bad"))
	req.Header.Set("X-Tenant-Id", tenantID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestProcessBulkBatch_InvalidBatchID(t *testing.T) {
	_, router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches/not-uuid/process", bytes.NewBufferString(`{"ais":{}}`))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestListBulkEmployees(t *testing.T) {
	_, router := setupRouter()
	tenantID := uuid.New().String()

	batchBody := `{"tax_year":"2026-27","employer_tan":"FFFF66666F","employer_name":"List Emp"}`
	bReq := httptest.NewRequest(http.MethodPost, "/api/v1/itr/bulk/batches", bytes.NewBufferString(batchBody))
	bReq.Header.Set("X-Tenant-Id", tenantID)
	bW := httptest.NewRecorder()
	router.ServeHTTP(bW, bReq)
	require.Equal(t, 201, bW.Code)

	var bResp map[string]interface{}
	json.Unmarshal(bW.Body.Bytes(), &bResp)
	batchID := bResp["data"].(map[string]interface{})["id"].(string)

	empBody := `{"pan":"LISTP1234F","name":"List Emp","email":"list@test.com","gross_salary":500000,"tds_deducted":25000}`
	eReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/itr/bulk/batches/%s/employees", batchID), bytes.NewBufferString(empBody))
	eReq.Header.Set("X-Tenant-Id", tenantID)
	eW := httptest.NewRecorder()
	router.ServeHTTP(eW, eReq)
	require.Equal(t, 201, eW.Code)

	listReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/itr/bulk/batches/%s/employees", batchID), nil)
	listReq.Header.Set("X-Tenant-Id", tenantID)
	listW := httptest.NewRecorder()
	router.ServeHTTP(listW, listReq)

	assert.Equal(t, 200, listW.Code)
	var listResp map[string]interface{}
	json.Unmarshal(listW.Body.Bytes(), &listResp)
	data := listResp["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
}
