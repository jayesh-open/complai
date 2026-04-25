package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/vendor-compliance-service/internal/domain"
	"github.com/complai/complai/services/go/vendor-compliance-service/internal/gateway"
	"github.com/complai/complai/services/go/vendor-compliance-service/internal/store"
)

func parseDataResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	var wrapper httputil.SuccessResponse
	wrapper.Data = target
	require.NoError(t, json.Unmarshal(body, &wrapper))
}

func fakeApexVendorServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/gateway/apex/vendors", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		vendorsData := map[string]interface{}{
			"vendors": []map[string]interface{}{
				{
					"id":                  "V001",
					"name":                "Acme Corp",
					"legal_name":          "Acme Corporation Pvt Ltd",
					"trade_name":          "Acme",
					"pan":                 "AABCA1234A",
					"gstin":               "29AABCA1234A1Z5",
					"tan":                 "BLRA12345B",
					"state":               "Karnataka",
					"state_code":          "29",
					"category":            "Regular",
					"registration_status": "Active",
					"msme_registered":     true,
					"email":               "acme@example.com",
					"phone":               "9876543210",
					"address":             "Bangalore, Karnataka",
				},
				{
					"id":                  "V002",
					"name":                "Beta Inc",
					"legal_name":          "Beta Incorporated Pvt Ltd",
					"trade_name":          "Beta",
					"pan":                 "AABCB5678B",
					"gstin":               "27AABCB5678B1Z3",
					"tan":                 "",
					"state":               "Maharashtra",
					"state_code":          "27",
					"category":            "Regular",
					"registration_status": "Active",
					"msme_registered":     false,
					"email":               "",
					"phone":               "9876543211",
					"address":             "Mumbai, Maharashtra",
				},
			},
			"total":      2,
			"request_id": "test-req",
		}

		gwResp := map[string]interface{}{
			"data": vendorsData,
			"meta": map[string]interface{}{
				"request_id":      "test-req",
				"latency_ms":      1,
				"provider_status": "success",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": gwResp})
	})

	mux.HandleFunc("/v1/gateway/apex/ap-invoices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			VendorID string `json:"vendor_id"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		var invoices []map[string]interface{}
		if req.VendorID == "" || req.VendorID == "V001" {
			invoices = append(invoices, []map[string]interface{}{
				{"vendor_id": "V001", "irn_generated": true, "gst_filing_status": "filed", "mismatch_status": "matched", "payment_status": "paid", "payment_date": "2026-03-15", "due_date": "2026-03-20"},
				{"vendor_id": "V001", "irn_generated": true, "gst_filing_status": "filed", "mismatch_status": "matched", "payment_status": "paid", "payment_date": "2026-03-10", "due_date": "2026-03-15"},
				{"vendor_id": "V001", "irn_generated": true, "gst_filing_status": "filed", "mismatch_status": "matched", "payment_status": "paid", "payment_date": "2026-02-28", "due_date": "2026-03-05"},
			}...)
		}
		if req.VendorID == "" || req.VendorID == "V002" {
			invoices = append(invoices, []map[string]interface{}{
				{"vendor_id": "V002", "irn_generated": false, "gst_filing_status": "late", "mismatch_status": "mismatched", "payment_status": "overdue", "payment_date": "", "due_date": "2026-03-01"},
				{"vendor_id": "V002", "irn_generated": false, "gst_filing_status": "not_filed", "mismatch_status": "pending", "payment_status": "unpaid", "payment_date": "", "due_date": "2026-02-15"},
			}...)
		}

		invoicesData := map[string]interface{}{
			"invoices":   invoices,
			"total":      len(invoices),
			"request_id": "test-req",
		}

		gwResp := map[string]interface{}{
			"data": invoicesData,
			"meta": map[string]interface{}{
				"request_id":      "test-req",
				"latency_ms":      1,
				"provider_status": "success",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": gwResp})
	})

	return httptest.NewServer(mux)
}

func setupHandlers(t *testing.T) (*Handlers, *store.MockStore, func()) {
	t.Helper()
	apexServer := fakeApexVendorServer()

	mockStore := store.NewMockStore()
	apex := gateway.NewApexClient(apexServer.URL)
	h := NewHandlers(mockStore, apex)

	return h, mockStore, func() {
		apexServer.Close()
	}
}

func postJSON(t *testing.T, handler http.HandlerFunc, url string, body string, tenantID uuid.UUID) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()
	handler(rec, req)
	return rec
}

func getJSON(t *testing.T, handler http.HandlerFunc, url string, tenantID uuid.UUID) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, url, nil)
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
	assert.Equal(t, "vendor-compliance-service", data["service"])
}

func TestTriggerSync_Success(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	rec := postJSON(t, h.TriggerSync, "/v1/vendor-compliance/sync", `{"request_id":"test"}`, tenantID)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.SyncResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, 2, resp.VendorCount)
	assert.Equal(t, 2, resp.ScoredCount)
	assert.Equal(t, "completed", resp.Status)
	assert.NotEqual(t, uuid.Nil, resp.SyncID)
}

func TestTriggerSync_MissingTenantID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/v1/vendor-compliance/sync", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()
	h.TriggerSync(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListVendors_WithScores(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	postJSON(t, h.TriggerSync, "/v1/vendor-compliance/sync", `{}`, tenantID)

	rec := getJSON(t, h.ListVendors, "/v1/vendor-compliance/vendors", tenantID)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.VendorListResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, 2, resp.Total)
	assert.Len(t, resp.Vendors, 2)
}

func TestListVendors_Empty(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	rec := getJSON(t, h.ListVendors, "/v1/vendor-compliance/vendors", tenantID)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.VendorListResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Total)
}

func TestListVendors_MissingTenantID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/vendor-compliance/vendors", nil)
	rec := httptest.NewRecorder()
	h.ListVendors(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetVendorScore_Success(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	postJSON(t, h.TriggerSync, "/v1/vendor-compliance/sync", `{}`, tenantID)

	router := chi.NewRouter()
	router.Get("/v1/vendor-compliance/vendors/{vendorId}/score", h.GetVendorScore)

	req := httptest.NewRequest(http.MethodGet, "/v1/vendor-compliance/vendors/V001/score", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.VendorScoreResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, "V001", resp.Vendor.VendorID)
	assert.Equal(t, "Acme Corp", resp.Vendor.Name)
	assert.Greater(t, resp.Score.TotalScore, 0)
}

func TestGetVendorScore_NotFound(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()

	router := chi.NewRouter()
	router.Get("/v1/vendor-compliance/vendors/{vendorId}/score", h.GetVendorScore)

	req := httptest.NewRequest(http.MethodGet, "/v1/vendor-compliance/vendors/NONEXISTENT/score", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetVendorScore_MissingTenantID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	router := chi.NewRouter()
	router.Get("/v1/vendor-compliance/vendors/{vendorId}/score", h.GetVendorScore)

	req := httptest.NewRequest(http.MethodGet, "/v1/vendor-compliance/vendors/V001/score", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetSyncStatus_AfterSync(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	postJSON(t, h.TriggerSync, "/v1/vendor-compliance/sync", `{}`, tenantID)

	rec := getJSON(t, h.GetSyncStatus, "/v1/vendor-compliance/sync/status", tenantID)

	assert.Equal(t, http.StatusOK, rec.Code)
	var status domain.SyncStatus
	parseDataResponse(t, rec.Body.Bytes(), &status)
	assert.Equal(t, "completed", status.Status)
	assert.Equal(t, 2, status.VendorCount)
}

func TestGetSyncStatus_NoSync(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	rec := getJSON(t, h.GetSyncStatus, "/v1/vendor-compliance/sync/status", tenantID)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetSyncStatus_MissingTenantID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/vendor-compliance/sync/status", nil)
	rec := httptest.NewRecorder()
	h.GetSyncStatus(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetScoreSummary_AfterSync(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	postJSON(t, h.TriggerSync, "/v1/vendor-compliance/sync", `{}`, tenantID)

	rec := getJSON(t, h.GetScoreSummary, "/v1/vendor-compliance/summary", tenantID)

	assert.Equal(t, http.StatusOK, rec.Code)
	var summary domain.ScoreSummary
	parseDataResponse(t, rec.Body.Bytes(), &summary)
	assert.Equal(t, 2, summary.Total)
	assert.Greater(t, summary.AvgScore, 0)
}

func TestGetScoreSummary_MissingTenantID(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/vendor-compliance/summary", nil)
	rec := httptest.NewRecorder()
	h.GetScoreSummary(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestNewRouter(t *testing.T) {
	mockStore := store.NewMockStore()
	apex := gateway.NewApexClient("http://localhost:9999")
	h := NewHandlers(mockStore, apex)
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

func TestRouter_VendorScoreEndpoint(t *testing.T) {
	h, _, cleanup := setupHandlers(t)
	defer cleanup()
	r := NewRouter(h)

	tenantID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/v1/vendor-compliance/sync", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/v1/vendor-compliance/vendors/V001/score", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestTriggerSync_CreatesVendorsInStore(t *testing.T) {
	h, mockStore, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	postJSON(t, h.TriggerSync, "/v1/vendor-compliance/sync", `{}`, tenantID)

	v1, err := mockStore.GetVendorSnapshot(nil, tenantID, "V001")
	require.NoError(t, err)
	assert.Equal(t, "Acme Corp", v1.Name)
	assert.Equal(t, "29AABCA1234A1Z5", v1.GSTIN)

	v2, err := mockStore.GetVendorSnapshot(nil, tenantID, "V002")
	require.NoError(t, err)
	assert.Equal(t, "Beta Inc", v2.Name)
}

func TestTriggerSync_CreatesScores(t *testing.T) {
	h, mockStore, cleanup := setupHandlers(t)
	defer cleanup()

	tenantID := uuid.New()
	postJSON(t, h.TriggerSync, "/v1/vendor-compliance/sync", `{}`, tenantID)

	s1, err := mockStore.GetLatestScore(nil, tenantID, "V001")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, s1.TotalScore, 80)

	s2, err := mockStore.GetLatestScore(nil, tenantID, "V002")
	require.NoError(t, err)
	assert.Less(t, s2.TotalScore, s1.TotalScore)
}

func TestTriggerSync_ApexServerFails(t *testing.T) {
	failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer failServer.Close()

	mockStore := store.NewMockStore()
	apex := gateway.NewApexClient(failServer.URL)
	h := NewHandlers(mockStore, apex)

	tenantID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/vendor-compliance/sync", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()
	h.TriggerSync(rec, req)

	assert.Equal(t, http.StatusBadGateway, rec.Code)
}
