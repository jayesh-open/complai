package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/gst-service/internal/domain"
	"github.com/complai/complai/services/go/gst-service/internal/gateway"
	"github.com/complai/complai/services/go/gst-service/internal/store"
)

func parseDataResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	var wrapper httputil.SuccessResponse
	wrapper.Data = target
	require.NoError(t, json.Unmarshal(body, &wrapper))
}

func fakeAuraServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		invoices := make([]map[string]interface{}, 5)
		for i := 0; i < 5; i++ {
			invoices[i] = map[string]interface{}{
				"id":              uuid.New().String(),
				"document_number": "INV/2026/000" + string(rune('1'+i)),
				"document_date":   "01/04/2026",
				"document_type":   "INV",
				"supply_type":     "B2B",
				"reverse_charge":  false,
				"supplier":        map[string]string{"gstin": "29AABCA1234A1Z5", "name": "Test", "state_code": "29"},
				"buyer":           map[string]string{"gstin": "27AABCB0001B1Z5", "name": "Buyer", "state_code": "27"},
				"line_items": []map[string]interface{}{{
					"hsn": "9988", "taxable_value": "50000", "cgst_rate": "0", "cgst_amount": "0",
					"sgst_rate": "0", "sgst_amount": "0", "igst_rate": "18", "igst_amount": "9000",
				}},
				"totals": map[string]string{
					"taxable_value": "50000", "cgst": "0", "sgst": "0", "igst": "9000", "grand_total": "59000",
				},
				"place_of_supply": "27",
				"source_system":   "aura",
			}
		}
		resp := httputil.SuccessResponse{Data: map[string]interface{}{
			"invoices":    invoices,
			"total_count": 5,
		}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func fakeGSTNServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"status":  "success",
				"message": "ok",
				"arn":     "AA2904202600001234",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func setupHandlers(t *testing.T) (*Handlers, *store.MockStore, func()) {
	t.Helper()
	auraServer := fakeAuraServer()
	gstnServer := fakeGSTNServer()

	mockStore := store.NewMockStore()
	aura := gateway.NewAuraClient(auraServer.URL)
	gstn := gateway.NewGSTNClient(gstnServer.URL)
	h := NewHandlers(mockStore, aura, gstn, nil)

	return h, mockStore, func() {
		auraServer.Close()
		gstnServer.Close()
	}
}

func postJSON(t *testing.T, handler http.HandlerFunc, url string, body interface{}, tenantID uuid.UUID) *httptest.ResponseRecorder {
	t.Helper()
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()
	handler(rec, req)
	return rec
}

func TestHealth(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	h.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "gst-service", data["service"])
}

func TestIngest_Success(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.NotEqual(t, uuid.Nil, resp.FilingID)
	assert.Equal(t, 5, resp.Ingested)
}

func TestIngest_MissingTenantID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	payload, _ := json.Marshal(domain.IngestRequest{GSTIN: "X", ReturnPeriod: "042026"})
	req := httptest.NewRequest(http.MethodPost, "/v1/gst/gstr1/ingest", bytes.NewReader(payload))
	rec := httptest.NewRecorder()
	h.Ingest(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestIngest_MissingGSTIN(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	body := domain.IngestRequest{ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestIngest_InvalidBody(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/v1/gst/gstr1/ingest", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.Ingest(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestValidate_Success(t *testing.T) {
	h, mockStore, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	// First ingest
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	require.Equal(t, http.StatusOK, rec.Code)

	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	// Validate
	valBody := domain.ValidateRequest{FilingID: ingestResp.FilingID}
	rec = postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)
	assert.Equal(t, http.StatusOK, rec.Code)

	var valResp domain.ValidateResponse
	parseDataResponse(t, rec.Body.Bytes(), &valResp)
	assert.Equal(t, 5, valResp.TotalCount)
	assert.True(t, len(valResp.Sections) > 0)

	// Filing should be validated (no errors expected since mock invoices are valid)
	f, _ := mockStore.GetFiling(nil, tenantID, ingestResp.FilingID)
	assert.Equal(t, domain.FilingStatusValidated, f.Status)
}

func TestValidate_NotFound(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	body := domain.ValidateRequest{FilingID: uuid.New()}
	rec := postJSON(t, h.Validate, "/v1/gst/gstr1/validate", body, tenantID)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestApprove_Success(t *testing.T) {
	h, mockStore, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	// Ingest + Validate
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	valBody := domain.ValidateRequest{FilingID: ingestResp.FilingID}
	postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)

	// Approve (maker-checker)
	approver := uuid.New()
	approveBody := domain.ApproveRequest{FilingID: ingestResp.FilingID, ApprovedBy: approver}
	rec = postJSON(t, h.Approve, "/v1/gst/gstr1/approve", approveBody, tenantID)
	assert.Equal(t, http.StatusOK, rec.Code)

	f, _ := mockStore.GetFiling(nil, tenantID, ingestResp.FilingID)
	assert.Equal(t, domain.FilingStatusApproved, f.Status)
	assert.Equal(t, approver, *f.ApprovedBy)
}

func TestApprove_NotValidated(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	approveBody := domain.ApproveRequest{FilingID: ingestResp.FilingID, ApprovedBy: uuid.New()}
	rec = postJSON(t, h.Approve, "/v1/gst/gstr1/approve", approveBody, tenantID)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestFile_FullLifecycle(t *testing.T) {
	h, mockStore, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	filer := uuid.New()

	// Ingest
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	// Validate
	valBody := domain.ValidateRequest{FilingID: ingestResp.FilingID}
	postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)

	// Approve
	approveBody := domain.ApproveRequest{FilingID: ingestResp.FilingID, ApprovedBy: uuid.New()}
	postJSON(t, h.Approve, "/v1/gst/gstr1/approve", approveBody, tenantID)

	// File
	fileBody := domain.FileRequest{FilingID: ingestResp.FilingID, SignType: "EVC", OTP: "123456", FiledBy: filer}
	rec = postJSON(t, h.File, "/v1/gst/gstr1/file", fileBody, tenantID)
	assert.Equal(t, http.StatusOK, rec.Code)

	var fileResp domain.FileResponse
	parseDataResponse(t, rec.Body.Bytes(), &fileResp)
	assert.Equal(t, domain.FilingStatusFiled, fileResp.Status)
	assert.NotEmpty(t, fileResp.ARN)

	f, _ := mockStore.GetFiling(nil, tenantID, ingestResp.FilingID)
	assert.Equal(t, domain.FilingStatusFiled, f.Status)
	assert.NotEmpty(t, f.ARN)
	assert.Equal(t, filer, *f.FiledBy)
}

func TestFile_NotApproved(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	// Ingest + Validate (but no approve)
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	valBody := domain.ValidateRequest{FilingID: ingestResp.FilingID}
	postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)

	fileBody := domain.FileRequest{FilingID: ingestResp.FilingID, SignType: "EVC", FiledBy: uuid.New()}
	rec = postJSON(t, h.File, "/v1/gst/gstr1/file", fileBody, tenantID)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestFile_InvalidSignType(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	fileBody := domain.FileRequest{FilingID: uuid.New(), SignType: "INVALID", FiledBy: uuid.New()}
	rec := postJSON(t, h.File, "/v1/gst/gstr1/file", fileBody, tenantID)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSummary_Success(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	valBody := domain.ValidateRequest{FilingID: ingestResp.FilingID}
	postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr1/summary?filing_id="+ingestResp.FilingID.String(), nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec = httptest.NewRecorder()
	h.Summary(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var summary domain.GSTR1Summary
	parseDataResponse(t, rec.Body.Bytes(), &summary)
	assert.Equal(t, ingestResp.FilingID, summary.Filing.ID)
	assert.True(t, len(summary.Sections) > 0)
}

func TestSummary_InvalidFilingID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr1/summary?filing_id=bad", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.Summary(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListEntries_Success(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr1/entries?filing_id="+ingestResp.FilingID.String(), nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec = httptest.NewRecorder()
	h.ListEntries(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestListErrors_Empty(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	// Validate to generate any errors
	valBody := domain.ValidateRequest{FilingID: ingestResp.FilingID}
	postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr1/errors?filing_id="+ingestResp.FilingID.String(), nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec = httptest.NewRecorder()
	h.ListErrors(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestNewRouter(t *testing.T) {
	mockStore := store.NewMockStore()
	aura := gateway.NewAuraClient("http://localhost:9999")
	gstn := gateway.NewGSTNClient("http://localhost:9999")
	h := NewHandlers(mockStore, aura, gstn, nil)
	r := NewRouter(h)
	require.NotNil(t, r)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategorization_InTest(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	valBody := domain.ValidateRequest{FilingID: ingestResp.FilingID}
	rec = postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)
	var valResp domain.ValidateResponse
	parseDataResponse(t, rec.Body.Bytes(), &valResp)

	sectionNames := make(map[string]bool)
	for _, s := range valResp.Sections {
		sectionNames[s.Section] = true
		assert.True(t, s.InvoiceCount > 0)
		assert.True(t, s.TaxableValue.IsPositive())
	}
	assert.True(t, sectionNames["b2b"], "should have B2B section")
}

func TestValidate_InvalidBody(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/v1/gst/gstr1/validate", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.Validate(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestValidate_MissingTenantID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	payload, _ := json.Marshal(domain.ValidateRequest{FilingID: uuid.New()})
	req := httptest.NewRequest(http.MethodPost, "/v1/gst/gstr1/validate", bytes.NewReader(payload))
	rec := httptest.NewRecorder()
	h.Validate(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestValidate_WrongState(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	// Ingest + Validate + Approve → try to validate again
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	valBody := domain.ValidateRequest{FilingID: ingestResp.FilingID}
	postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)

	approveBody := domain.ApproveRequest{FilingID: ingestResp.FilingID, ApprovedBy: uuid.New()}
	postJSON(t, h.Approve, "/v1/gst/gstr1/approve", approveBody, tenantID)

	// Now validated→approved, validating again should fail
	rec = postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestApprove_InvalidBody(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/v1/gst/gstr1/approve", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.Approve(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestApprove_MissingTenantID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	payload, _ := json.Marshal(domain.ApproveRequest{FilingID: uuid.New(), ApprovedBy: uuid.New()})
	req := httptest.NewRequest(http.MethodPost, "/v1/gst/gstr1/approve", bytes.NewReader(payload))
	rec := httptest.NewRecorder()
	h.Approve(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestApprove_NotFound(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	approveBody := domain.ApproveRequest{FilingID: uuid.New(), ApprovedBy: uuid.New()}
	rec := postJSON(t, h.Approve, "/v1/gst/gstr1/approve", approveBody, tenantID)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestFile_InvalidBody(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/v1/gst/gstr1/file", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.File(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFile_MissingTenantID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	payload, _ := json.Marshal(domain.FileRequest{FilingID: uuid.New(), SignType: "EVC", FiledBy: uuid.New()})
	req := httptest.NewRequest(http.MethodPost, "/v1/gst/gstr1/file", bytes.NewReader(payload))
	rec := httptest.NewRecorder()
	h.File(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFile_NotFound(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	fileBody := domain.FileRequest{FilingID: uuid.New(), SignType: "DSC", FiledBy: uuid.New()}
	rec := postJSON(t, h.File, "/v1/gst/gstr1/file", fileBody, tenantID)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestSummary_MissingTenantID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr1/summary?filing_id="+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	h.Summary(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSummary_NotFound(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr1/summary?filing_id="+uuid.New().String(), nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.Summary(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListEntries_InvalidFilingID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr1/entries?filing_id=bad", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.ListEntries(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListEntries_MissingTenantID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr1/entries?filing_id="+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	h.ListEntries(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListErrors_InvalidFilingID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr1/errors?filing_id=bad", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.ListErrors(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListErrors_MissingTenantID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr1/errors?filing_id="+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	h.ListErrors(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestValidateEntries_MissingDocNumber(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	entries := []domain.SalesRegisterEntry{{
		ID: uuid.New(), DocumentNumber: "", SupplyType: "B2B", BuyerGSTIN: "29X",
		TaxableValue: decimal.NewFromInt(1000), HSN: "9988", PlaceOfSupply: "29",
	}}
	errs := validateEntries(filingID, tenantID, entries)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "REQUIRED", errs[0].Code)
}

func TestValidateEntries_B2BMissingGSTIN(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	entries := []domain.SalesRegisterEntry{{
		ID: uuid.New(), DocumentNumber: "INV/001", SupplyType: "B2B", BuyerGSTIN: "",
		TaxableValue: decimal.NewFromInt(1000), HSN: "9988", PlaceOfSupply: "29",
	}}
	errs := validateEntries(filingID, tenantID, entries)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "B2B_GSTIN_REQUIRED", errs[0].Code)
}

func TestValidateEntries_NegativeTaxable(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	entries := []domain.SalesRegisterEntry{{
		ID: uuid.New(), DocumentNumber: "INV/001", SupplyType: "B2CS",
		TaxableValue: decimal.NewFromInt(-500), HSN: "9988", PlaceOfSupply: "29",
	}}
	errs := validateEntries(filingID, tenantID, entries)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "NEGATIVE_VALUE", errs[0].Code)
}

func TestValidateEntries_MissingHSN(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	entries := []domain.SalesRegisterEntry{{
		ID: uuid.New(), DocumentNumber: "INV/001", SupplyType: "B2CS",
		TaxableValue: decimal.NewFromInt(1000), HSN: "", PlaceOfSupply: "29",
	}}
	errs := validateEntries(filingID, tenantID, entries)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "HSN_REQUIRED", errs[0].Code)
}

func TestValidateEntries_MissingPOS(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	entries := []domain.SalesRegisterEntry{{
		ID: uuid.New(), DocumentNumber: "INV/001", SupplyType: "B2CS",
		TaxableValue: decimal.NewFromInt(1000), HSN: "9988", PlaceOfSupply: "",
	}}
	errs := validateEntries(filingID, tenantID, entries)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "POS_REQUIRED", errs[0].Code)
}

func TestValidateEntries_Clean(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	entries := []domain.SalesRegisterEntry{{
		ID: uuid.New(), DocumentNumber: "INV/001", SupplyType: "B2CS",
		TaxableValue: decimal.NewFromInt(1000), HSN: "9988", PlaceOfSupply: "29",
	}}
	errs := validateEntries(filingID, tenantID, entries)
	assert.Empty(t, errs)
}

type mockStepUpVerifier struct {
	valid bool
}

func (m *mockStepUpVerifier) HasValidStepUp(_ context.Context, _, _ uuid.UUID, _ string) bool {
	return m.valid
}

func setupHandlersWithStepUp(t *testing.T, stepUp StepUpVerifier) (*Handlers, *store.MockStore, func()) {
	t.Helper()
	auraServer := fakeAuraServer()
	gstnServer := fakeGSTNServer()

	mockStore := store.NewMockStore()
	aura := gateway.NewAuraClient(auraServer.URL)
	gstn := gateway.NewGSTNClient(gstnServer.URL)
	h := NewHandlers(mockStore, aura, gstn, stepUp)

	return h, mockStore, func() {
		auraServer.Close()
		gstnServer.Close()
	}
}

func postJSONWithUser(t *testing.T, handler http.HandlerFunc, url string, body interface{}, tenantID, userID uuid.UUID) *httptest.ResponseRecorder {
	t.Helper()
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", userID.String())
	rec := httptest.NewRecorder()
	handler(rec, req)
	return rec
}

func TestFile_StepUpRequired(t *testing.T) {
	verifier := &mockStepUpVerifier{valid: false}
	h, mockStore, cleanup := setupHandlersWithStepUp(t, verifier)
	defer cleanup()

	tenantID := uuid.New()
	filer := uuid.New()

	// Ingest → Validate → Approve
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	require.Equal(t, http.StatusOK, rec.Code)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	valBody := domain.ValidateRequest{FilingID: ingestResp.FilingID}
	postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)

	approver := uuid.New()
	approveBody := domain.ApproveRequest{FilingID: ingestResp.FilingID, ApprovedBy: approver}
	postJSON(t, h.Approve, "/v1/gst/gstr1/approve", approveBody, tenantID)

	// File without valid step-up → 403
	fileBody := domain.FileRequest{FilingID: ingestResp.FilingID, SignType: "EVC", OTP: "123456", FiledBy: filer}
	rec = postJSONWithUser(t, h.File, "/v1/gst/gstr1/file", fileBody, tenantID, filer)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	var errResp map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &errResp)
	assert.Equal(t, "step_up_required", errResp["error"])

	// Filing should still be approved (not transitioned)
	f, _ := mockStore.GetFiling(nil, tenantID, ingestResp.FilingID)
	assert.Equal(t, domain.FilingStatusApproved, f.Status)
}

func TestFile_StepUpVerified(t *testing.T) {
	verifier := &mockStepUpVerifier{valid: true}
	h, mockStore, cleanup := setupHandlersWithStepUp(t, verifier)
	defer cleanup()

	tenantID := uuid.New()
	filer := uuid.New()

	// Ingest → Validate → Approve
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID)
	require.Equal(t, http.StatusOK, rec.Code)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	valBody := domain.ValidateRequest{FilingID: ingestResp.FilingID}
	postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)

	approver := uuid.New()
	approveBody := domain.ApproveRequest{FilingID: ingestResp.FilingID, ApprovedBy: approver}
	postJSON(t, h.Approve, "/v1/gst/gstr1/approve", approveBody, tenantID)

	// File with valid step-up → 200
	fileBody := domain.FileRequest{FilingID: ingestResp.FilingID, SignType: "EVC", OTP: "123456", FiledBy: filer}
	rec = postJSONWithUser(t, h.File, "/v1/gst/gstr1/file", fileBody, tenantID, filer)
	assert.Equal(t, http.StatusOK, rec.Code)

	var fileResp domain.FileResponse
	parseDataResponse(t, rec.Body.Bytes(), &fileResp)
	assert.Equal(t, domain.FilingStatusFiled, fileResp.Status)
	assert.NotEmpty(t, fileResp.ARN)

	f, _ := mockStore.GetFiling(nil, tenantID, ingestResp.FilingID)
	assert.Equal(t, domain.FilingStatusFiled, f.Status)
	assert.Equal(t, filer, *f.FiledBy)
}

func TestApprove_SelfApprovalDenied(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	creator := uuid.New()

	// Ingest with X-User-Id so CreatedBy is set
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSONWithUser(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID, creator)
	require.Equal(t, http.StatusOK, rec.Code)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	// Validate
	valBody := domain.ValidateRequest{FilingID: ingestResp.FilingID}
	postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)

	// Same user tries to approve → 403 self_approval_denied
	approveBody := domain.ApproveRequest{FilingID: ingestResp.FilingID, ApprovedBy: creator}
	rec = postJSONWithUser(t, h.Approve, "/v1/gst/gstr1/approve", approveBody, tenantID, creator)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	var errResp map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &errResp)
	assert.Equal(t, "self_approval_denied", errResp["error"])
}

func fakeGSTNServerWith3B() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/v1/gateway/adaequare/gstr1/summary":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": gateway.GSTR1SummaryResponse{
					GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", Status: "success",
					Summary: map[string]interface{}{
						"b2b":  map[string]interface{}{"taxable_value": 500000.0, "cgst": 0.0, "sgst": 0.0, "igst": 90000.0},
						"b2cs": map[string]interface{}{"taxable_value": 200000.0, "cgst": 18000.0, "sgst": 18000.0, "igst": 0.0},
						"cdnr": map[string]interface{}{"taxable_value": 50000.0, "cgst": 0.0, "sgst": 0.0, "igst": 9000.0},
					},
				},
			})
		case "/v1/gateway/adaequare/gstr2b/get":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": gateway.GSTR2BGetResponse{
					GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", Status: "success", TotalCount: 4,
					Invoices: []gateway.GSTR2BInvoice{
						{SupplierGSTIN: "27AABCB0001B1Z5", InvoiceNumber: "INV-P001", InvoiceDate: "01/04/2026", TaxableValue: 100000, IGSTAmount: 18000, TotalValue: 118000, HSN: "9988"},
						{SupplierGSTIN: "27AABCB0002B1Z5", InvoiceNumber: "INV-P002", InvoiceDate: "05/04/2026", TaxableValue: 100000, IGSTAmount: 18000, TotalValue: 118000, HSN: "9988"},
						{SupplierGSTIN: "27AABCB0003B1Z5", InvoiceNumber: "INV-P003", InvoiceDate: "10/04/2026", TaxableValue: 100000, IGSTAmount: 18000, TotalValue: 118000, HSN: "9988"},
						{SupplierGSTIN: "29AABCB0004B1Z5", InvoiceNumber: "INV-P004", InvoiceDate: "15/04/2026", TaxableValue: 50000, CGSTAmount: 4500, SGSTAmount: 4500, TotalValue: 59000, HSN: "9988"},
					},
				},
			})
		case "/v1/gateway/adaequare/ims/get":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": gateway.IMSGetResponse{
					GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", Status: "success",
					Invoices: []gateway.IMSInvoice{
						{InvoiceID: "ims-001", SupplierGSTIN: "27AABCB0001B1Z5", InvoiceNumber: "INV-P001", TaxableValue: 100000, TotalValue: 118000, IGSTAmount: 18000, Action: "ACCEPT"},
						{InvoiceID: "ims-002", SupplierGSTIN: "27AABCB0002B1Z5", InvoiceNumber: "INV-P002", TaxableValue: 100000, TotalValue: 118000, IGSTAmount: 18000, Action: "ACCEPT"},
						{InvoiceID: "ims-003", SupplierGSTIN: "27AABCB0003B1Z5", InvoiceNumber: "INV-P003", TaxableValue: 100000, TotalValue: 118000, IGSTAmount: 18000, Action: "REJECT"},
					},
					Summary: gateway.IMSSummary{Accepted: 2, Rejected: 1, Pending: 1, AcceptedValue: 236000, RejectedValue: 118000, PendingValue: 59000},
				},
			})
		default:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"status": "success", "message": "ok", "arn": "AA2904202600001234"},
			})
		}
	}))
}

func setupHandlersWith3B(t *testing.T) (*Handlers, *store.MockStore, func()) {
	t.Helper()
	auraServer := fakeAuraServer()
	gstnServer := fakeGSTNServerWith3B()

	mockStore := store.NewMockStore()
	aura := gateway.NewAuraClient(auraServer.URL)
	gstn := gateway.NewGSTNClient(gstnServer.URL)
	h := NewHandlers(mockStore, aura, gstn, nil)

	return h, mockStore, func() {
		auraServer.Close()
		gstnServer.Close()
	}
}

func TestGSTR3BAutoFill(t *testing.T) {
	h, mockStore, cleanup := setupHandlersWith3B(t)
	defer cleanup()

	tenantID := uuid.New()
	body := domain.GSTR3BAutoFillRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSON(t, h.GSTR3BAutoFill, "/v1/gst/gstr3b/auto-fill", body, tenantID)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp domain.GSTR3BAutoFillResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.NotEqual(t, uuid.Nil, resp.FilingID)

	d := resp.Data

	// Table 1: Outward supply from GSTR-1
	assert.True(t, d.GSTR1Summary.B2B.IGST.Equal(decimal.NewFromInt(90000)), "B2B IGST should be 90000")
	assert.True(t, d.GSTR1Summary.B2B.TaxableValue.Equal(decimal.NewFromInt(500000)), "B2B taxable should be 500000")
	assert.True(t, d.GSTR1Summary.B2CS.CGST.Equal(decimal.NewFromInt(18000)), "B2CS CGST should be 18000")
	assert.True(t, d.GSTR1Summary.B2CS.SGST.Equal(decimal.NewFromInt(18000)), "B2CS SGST should be 18000")
	assert.True(t, d.GSTR1Summary.CreditNote.IGST.Equal(decimal.NewFromInt(9000)), "Credit note IGST should be 9000")

	// Table 2: Inward supply from GSTR-2B (4 invoices, all 18% rate)
	assert.True(t, d.InwardSupply.TotalValue.Equal(decimal.NewFromFloat(413000)), "Inward total value")
	assert.True(t, d.InwardSupply.TaxableAt18.Equal(decimal.NewFromFloat(350000)), "All taxable at 18%")
	assert.True(t, d.InwardSupply.ITCAvailable.IGST.Equal(decimal.NewFromFloat(54000)), "Inward ITC IGST")
	assert.True(t, d.InwardSupply.ITCAvailable.CGST.Equal(decimal.NewFromFloat(4500)), "Inward ITC CGST")

	// Table 3: IMS summary
	assert.Equal(t, 2, d.IMSActions.Accepted)
	assert.Equal(t, 1, d.IMSActions.Rejected)
	assert.Equal(t, 1, d.IMSActions.Pending)

	// Table 4A-D: Eligible ITC (GSTR-2B minus rejected by IMS, minus RCM)
	// INV-P001 (IGST 18000) + INV-P002 (IGST 18000) + INV-P004 (CGST 4500, SGST 4500) = IGST 36000, CGST 4500, SGST 4500
	assert.True(t, d.EligibleITC.Total.IGST.Equal(decimal.NewFromFloat(36000)), "Eligible ITC IGST = 36000")
	assert.True(t, d.EligibleITC.Total.CGST.Equal(decimal.NewFromFloat(4500)), "Eligible ITC CGST = 4500")
	assert.True(t, d.EligibleITC.Total.SGST.Equal(decimal.NewFromFloat(4500)), "Eligible ITC SGST = 4500")

	// Table 5: Gross liability = B2B + B2CS - CreditNote (B2CL, Exports, Advances are zero)
	// CGST: 0 + 18000 - 0 = 18000, SGST: 0 + 18000 - 0 = 18000, IGST: 90000 + 0 - 9000 = 81000
	assert.True(t, d.GrossLiability.CGST.Equal(decimal.NewFromInt(18000)), "Gross CGST = 18000")
	assert.True(t, d.GrossLiability.SGST.Equal(decimal.NewFromInt(18000)), "Gross SGST = 18000")
	assert.True(t, d.GrossLiability.IGST.Equal(decimal.NewFromInt(81000)), "Gross IGST = 81000")

	// Table 6: Net liability = Gross - Eligible ITC
	// CGST: 18000 - 4500 = 13500, SGST: 18000 - 4500 = 13500, IGST: 81000 - 36000 = 45000
	assert.True(t, d.NetLiability.CGST.Equal(decimal.NewFromInt(13500)), "Net CGST = 13500")
	assert.True(t, d.NetLiability.SGST.Equal(decimal.NewFromInt(13500)), "Net SGST = 13500")
	assert.True(t, d.NetLiability.IGST.Equal(decimal.NewFromInt(45000)), "Net IGST = 45000")

	// Flags: 1 pending IMS action, no RCM
	require.Len(t, d.Flags, 1)
	assert.Contains(t, d.Flags[0], "pending IMS action")

	// Filing persisted in store with correct status
	filing, err := mockStore.GetGSTR3BFiling(context.Background(), tenantID, resp.FilingID)
	require.NoError(t, err)
	assert.Equal(t, domain.GSTR3BStatusPopulated, filing.Status)
	assert.Equal(t, "29AABCA1234A1Z5", filing.GSTIN)
	assert.NotEmpty(t, filing.DataJSON)
}

func TestGSTR3BAutoFill_MissingFields(t *testing.T) {
	h, _, cleanup := setupHandlersWith3B(t)
	defer cleanup()

	tenantID := uuid.New()
	body := domain.GSTR3BAutoFillRequest{GSTIN: "", ReturnPeriod: "042026"}
	rec := postJSON(t, h.GSTR3BAutoFill, "/v1/gst/gstr3b/auto-fill", body, tenantID)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestApprove_MakerCheckerFlow(t *testing.T) {
	h, mockStore, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	analyst := uuid.New()
	taxManager := uuid.New()

	// Analyst ingests filing (sets CreatedBy = analyst)
	body := domain.IngestRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"}
	rec := postJSONWithUser(t, h.Ingest, "/v1/gst/gstr1/ingest", body, tenantID, analyst)
	require.Equal(t, http.StatusOK, rec.Code)
	var ingestResp domain.IngestResponse
	parseDataResponse(t, rec.Body.Bytes(), &ingestResp)

	// Validate
	valBody := domain.ValidateRequest{FilingID: ingestResp.FilingID}
	rec = postJSON(t, h.Validate, "/v1/gst/gstr1/validate", valBody, tenantID)
	require.Equal(t, http.StatusOK, rec.Code)

	// Verify filing is validated
	f, _ := mockStore.GetFiling(nil, tenantID, ingestResp.FilingID)
	require.Equal(t, domain.FilingStatusValidated, f.Status)

	// Same analyst tries to approve → 403 self_approval_denied
	approveBody := domain.ApproveRequest{FilingID: ingestResp.FilingID, ApprovedBy: analyst}
	rec = postJSONWithUser(t, h.Approve, "/v1/gst/gstr1/approve", approveBody, tenantID, analyst)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	var errResp map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &errResp)
	assert.Equal(t, "self_approval_denied", errResp["error"])

	// Filing should still be validated (not approved)
	f, _ = mockStore.GetFiling(nil, tenantID, ingestResp.FilingID)
	assert.Equal(t, domain.FilingStatusValidated, f.Status)

	// Different user (tax-manager) approves → 200
	approveBody = domain.ApproveRequest{FilingID: ingestResp.FilingID, ApprovedBy: taxManager}
	rec = postJSONWithUser(t, h.Approve, "/v1/gst/gstr1/approve", approveBody, tenantID, taxManager)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Filing should be approved with tax-manager as approver
	f, _ = mockStore.GetFiling(nil, tenantID, ingestResp.FilingID)
	assert.Equal(t, domain.FilingStatusApproved, f.Status)
	assert.Equal(t, taxManager, *f.ApprovedBy)
}

// ---------------------------------------------------------------------------
// Tests: GSTR3BSummary
// ---------------------------------------------------------------------------

func TestGSTR3BSummary_Success(t *testing.T) {
	h, mockStore, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	filing := &domain.GSTR3BFiling{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026", Status: domain.GSTR3BStatusPopulated}
	_ = mockStore.CreateGSTR3BFiling(nil, tenantID, filing)

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr3b/summary?filing_id="+filing.ID.String(), nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()
	h.GSTR3BSummary(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGSTR3BSummary_MissingTenant(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr3b/summary?filing_id="+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	h.GSTR3BSummary(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR3BSummary_InvalidFilingID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr3b/summary?filing_id=bad", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GSTR3BSummary(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR3BSummary_NotFound(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/gst/gstr3b/summary?filing_id="+uuid.New().String(), nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GSTR3BSummary(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: GSTR3BApprove
// ---------------------------------------------------------------------------

func TestGSTR3BApprove_Success(t *testing.T) {
	h, mockStore, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	maker := uuid.New()
	checker := uuid.New()
	filing := &domain.GSTR3BFiling{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026", Status: domain.GSTR3BStatusPopulated, CreatedBy: &maker}
	_ = mockStore.CreateGSTR3BFiling(nil, tenantID, filing)

	body := domain.GSTR3BApproveRequest{FilingID: filing.ID, ApprovedBy: checker}
	rec := postJSONWithUser(t, h.GSTR3BApprove, "/v1/gst/gstr3b/approve", body, tenantID, checker)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGSTR3BApprove_SelfApprovalDenied(t *testing.T) {
	h, mockStore, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	maker := uuid.New()
	filing := &domain.GSTR3BFiling{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026", Status: domain.GSTR3BStatusPopulated, CreatedBy: &maker}
	_ = mockStore.CreateGSTR3BFiling(nil, tenantID, filing)

	body := domain.GSTR3BApproveRequest{FilingID: filing.ID, ApprovedBy: maker}
	rec := postJSONWithUser(t, h.GSTR3BApprove, "/v1/gst/gstr3b/approve", body, tenantID, maker)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestGSTR3BApprove_WrongStatus(t *testing.T) {
	h, mockStore, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	filing := &domain.GSTR3BFiling{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026", Status: domain.GSTR3BStatusDraft}
	_ = mockStore.CreateGSTR3BFiling(nil, tenantID, filing)

	body := domain.GSTR3BApproveRequest{FilingID: filing.ID, ApprovedBy: uuid.New()}
	rec := postJSONWithUser(t, h.GSTR3BApprove, "/v1/gst/gstr3b/approve", body, tenantID, uuid.New())
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestGSTR3BApprove_NotFound(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	body := domain.GSTR3BApproveRequest{FilingID: uuid.New(), ApprovedBy: uuid.New()}
	rec := postJSONWithUser(t, h.GSTR3BApprove, "/v1/gst/gstr3b/approve", body, uuid.New(), uuid.New())
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGSTR3BApprove_InvalidBody(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/v1/gst/gstr3b/approve", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GSTR3BApprove(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: GSTR3BFile
// ---------------------------------------------------------------------------

func TestGSTR3BFile_Success(t *testing.T) {
	h, mockStore, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	filer := uuid.New()
	filing := &domain.GSTR3BFiling{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026", Status: domain.GSTR3BStatusApproved}
	_ = mockStore.CreateGSTR3BFiling(nil, tenantID, filing)

	body := domain.GSTR3BFileRequest{FilingID: filing.ID, SignType: "EVC", FiledBy: filer}
	rec := postJSONWithUser(t, h.GSTR3BFile, "/v1/gst/gstr3b/file", body, tenantID, filer)
	assert.Equal(t, http.StatusOK, rec.Code)

	var wrapper struct {
		Data domain.GSTR3BFileResponse `json:"data"`
	}
	json.NewDecoder(rec.Body).Decode(&wrapper)
	assert.Equal(t, domain.GSTR3BStatusFiled, wrapper.Data.Status)
	assert.NotEmpty(t, wrapper.Data.ARN)
}

func TestGSTR3BFile_InvalidSignType(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	body := domain.GSTR3BFileRequest{FilingID: uuid.New(), SignType: "INVALID", FiledBy: uuid.New()}
	rec := postJSONWithUser(t, h.GSTR3BFile, "/v1/gst/gstr3b/file", body, uuid.New(), uuid.New())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR3BFile_NotApproved(t *testing.T) {
	h, mockStore, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	filing := &domain.GSTR3BFiling{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026", Status: domain.GSTR3BStatusPopulated}
	_ = mockStore.CreateGSTR3BFiling(nil, tenantID, filing)

	body := domain.GSTR3BFileRequest{FilingID: filing.ID, SignType: "EVC", FiledBy: uuid.New()}
	rec := postJSONWithUser(t, h.GSTR3BFile, "/v1/gst/gstr3b/file", body, tenantID, uuid.New())
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestGSTR3BFile_NotFound(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	body := domain.GSTR3BFileRequest{FilingID: uuid.New(), SignType: "EVC", FiledBy: uuid.New()}
	rec := postJSONWithUser(t, h.GSTR3BFile, "/v1/gst/gstr3b/file", body, uuid.New(), uuid.New())
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGSTR3BFile_InvalidBody(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/v1/gst/gstr3b/file", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GSTR3BFile(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
