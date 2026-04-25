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
