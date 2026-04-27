package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/einvoice-service/internal/domain"
	"github.com/complai/complai/services/go/einvoice-service/internal/gateway"
	"github.com/complai/complai/services/go/einvoice-service/internal/store"
)

// fixedClock returns a fixed time for deterministic tests.
type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

func newMockIRPServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/gateway/irp/invoice", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)

		gstin, _ := req["gstin"].(string)
		docDtls, _ := req["doc_dtls"].(map[string]interface{})
		docNo, _ := docDtls["no"].(string)

		irn := fmt.Sprintf("mock-irn-%s-%s", gstin[:6], docNo)

		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]interface{}{
					"irn":            irn,
					"ack_no":         "1234567890",
					"ack_dt":         "15/04/2026 10:30:00",
					"signed_invoice": "c2lnbmVkLWludm9pY2U=",
					"signed_qr_code": "c2lnbmVkLXFy",
					"status":         "ACT",
				},
				"meta": map[string]interface{}{
					"request_id":      uuid.New().String(),
					"latency_ms":      10,
					"provider_status": "success",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/v1/gateway/irp/invoice/cancel", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)
		irn, _ := req["irn"].(string)

		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]interface{}{
					"irn":         irn,
					"cancel_date": "15/04/2026 11:30:00",
					"status":      "CANC",
				},
				"meta": map[string]interface{}{
					"request_id":      uuid.New().String(),
					"latency_ms":      5,
					"provider_status": "success",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	return httptest.NewServer(mux)
}

func setupTestHandlers(t *testing.T) (*Handlers, *store.MockStore, *httptest.Server) {
	t.Helper()
	mockStore := store.NewMockStore()
	irpServer := newMockIRPServer(t)
	irpClient := newTestIRPClient(irpServer.URL)
	h := NewHandlers(mockStore, irpClient, store.RealClock{})
	return h, mockStore, irpServer
}

func newTestIRPClient(baseURL string) *gateway.IRPClient {
	return gateway.NewIRPClient(baseURL)
}

func parseDataResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	var wrapper httputil.SuccessResponse
	wrapper.Data = target
	require.NoError(t, json.Unmarshal(body, &wrapper))
}

func postJSON(t *testing.T, path string, body interface{}, tenantID string) *http.Request {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID)
	return req
}

func getWithTenant(t *testing.T, path, tenantID string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("X-Tenant-Id", tenantID)
	return req
}

func withChiURLParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// ---------------------------------------------------------------------------
// Tests: Health
// ---------------------------------------------------------------------------

func TestHealth(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()

	rec := httptest.NewRecorder()
	h.Health(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	assert.Equal(t, http.StatusOK, rec.Code)

	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "einvoice-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: Generate IRN
// ---------------------------------------------------------------------------

func TestGenerateIRN_Success(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()

	req := domain.GenerateIRNRequest{
		InvoiceNumber: "INV-001",
		InvoiceDate:   "15/04/2026",
		InvoiceType:   domain.InvoiceTypeINV,
		SupplierGSTIN: "29AABCA1234A1Z5",
		SupplierName:  "Test Supplier",
		BuyerGSTIN:    "27AABCB5678B1Z3",
		BuyerName:     "Test Buyer",
		SupplyType:    domain.SupplyTypeB2B,
		PlaceOfSupply: "27",
		TaxableValue:  decimal.NewFromInt(25000),
		IGSTAmount:    decimal.NewFromInt(4500),
		TotalAmount:   decimal.NewFromInt(29500),
		SourceSystem:  "aura",
		LineItems: []domain.LineItemRequest{
			{
				Description:  "Steel Plates",
				HSNCode:      "720241",
				Quantity:     decimal.NewFromInt(100),
				Unit:         "KG",
				UnitPrice:    decimal.NewFromInt(250),
				TaxableValue: decimal.NewFromInt(25000),
				IGSTRate:     decimal.NewFromInt(18),
				IGSTAmount:   decimal.NewFromInt(4500),
			},
		},
	}

	rec := httptest.NewRecorder()
	h.GenerateIRN(rec, postJSON(t, "/v1/einvoice/generate", req, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var resp domain.GenerateIRNResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.NotEqual(t, uuid.Nil, resp.ID)
	assert.Contains(t, resp.IRN, "mock-irn-")
	assert.Equal(t, domain.IRNStatusGenerated, resp.Status)
	assert.NotEmpty(t, resp.SignedInvoice)
	assert.NotEmpty(t, resp.SignedQRCode)
}

func TestGenerateIRN_MissingFields(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()

	req := domain.GenerateIRNRequest{InvoiceNumber: "INV-001"}
	rec := httptest.NewRecorder()
	h.GenerateIRN(rec, postJSON(t, "/", req, tenantID))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGenerateIRN_MissingTenantID(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()

	body, _ := json.Marshal(domain.GenerateIRNRequest{})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.GenerateIRN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGenerateIRN_InvalidBody(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GenerateIRN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Cancel IRN — 24h window enforcement
// ---------------------------------------------------------------------------

func TestCancelIRN_WithinWindow(t *testing.T) {
	mockStore := store.NewMockStore()
	irpServer := newMockIRPServer(t)
	defer irpServer.Close()
	irpClient := newTestIRPClient(irpServer.URL)

	now := time.Now()
	clock := fixedClock{now: now}
	h := NewHandlers(mockStore, irpClient, clock)
	tenantID := uuid.New()

	inv := &domain.EInvoice{
		InvoiceNumber: "INV-CANCEL-001",
		InvoiceDate:   "15/04/2026",
		InvoiceType:   domain.InvoiceTypeINV,
		SupplierGSTIN: "29AABCA1234A1Z5",
		BuyerGSTIN:    "27AABCB5678B1Z3",
		TaxableValue:  decimal.NewFromInt(10000),
		TotalAmount:   decimal.NewFromInt(11800),
		SourceSystem:  "test",
	}
	require.NoError(t, mockStore.CreateEInvoice(nil, tenantID, inv))

	require.NoError(t, mockStore.UpdateIRNGenerated(nil, tenantID, inv.ID, "test-irn-123", "ack123", "signed", "qr"))
	generatedAt := now.Add(-12 * time.Hour)
	mockStore.SetIRNGeneratedAt(inv.ID, &generatedAt)

	cancelReq := domain.CancelIRNRequest{Reason: "1", Remark: "Duplicate"}
	body, _ := json.Marshal(cancelReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/einvoice/"+inv.ID.String()+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", inv.ID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.CancelIRN(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp domain.CancelIRNResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, domain.IRNStatusCancelled, resp.Status)
}

func TestCancelIRN_ExpiredWindow(t *testing.T) {
	mockStore := store.NewMockStore()
	irpServer := newMockIRPServer(t)
	defer irpServer.Close()
	irpClient := newTestIRPClient(irpServer.URL)

	now := time.Now()
	clock := fixedClock{now: now}
	h := NewHandlers(mockStore, irpClient, clock)
	tenantID := uuid.New()

	inv := &domain.EInvoice{
		InvoiceNumber: "INV-EXPIRED-001",
		InvoiceDate:   "14/04/2026",
		InvoiceType:   domain.InvoiceTypeINV,
		SupplierGSTIN: "29AABCA1234A1Z5",
		BuyerGSTIN:    "27AABCB5678B1Z3",
		TaxableValue:  decimal.NewFromInt(10000),
		TotalAmount:   decimal.NewFromInt(11800),
		SourceSystem:  "test",
	}
	require.NoError(t, mockStore.CreateEInvoice(nil, tenantID, inv))
	require.NoError(t, mockStore.UpdateIRNGenerated(nil, tenantID, inv.ID, "test-irn-expired", "ack", "sig", "qr"))
	expiredTime := now.Add(-25 * time.Hour)
	mockStore.SetIRNGeneratedAt(inv.ID, &expiredTime)

	cancelReq := domain.CancelIRNRequest{Reason: "1", Remark: "Too late"}
	body, _ := json.Marshal(cancelReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/einvoice/"+inv.ID.String()+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", inv.ID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.CancelIRN(rec, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCancelIRN_NotGenerated(t *testing.T) {
	mockStore := store.NewMockStore()
	irpServer := newMockIRPServer(t)
	defer irpServer.Close()
	irpClient := newTestIRPClient(irpServer.URL)
	h := NewHandlers(mockStore, irpClient, store.RealClock{})
	tenantID := uuid.New()

	inv := &domain.EInvoice{
		InvoiceNumber: "INV-PENDING-001",
		InvoiceType:   domain.InvoiceTypeINV,
		SupplierGSTIN: "29AABCA1234A1Z5",
		TaxableValue:  decimal.NewFromInt(10000),
		TotalAmount:   decimal.NewFromInt(11800),
		SourceSystem:  "test",
	}
	require.NoError(t, mockStore.CreateEInvoice(nil, tenantID, inv))

	cancelReq := domain.CancelIRNRequest{Reason: "1", Remark: "Not generated"}
	body, _ := json.Marshal(cancelReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/einvoice/"+inv.ID.String()+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", inv.ID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.CancelIRN(rec, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCancelIRN_MissingReason(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()

	cancelReq := domain.CancelIRNRequest{Remark: "No reason given"}
	body, _ := json.Marshal(cancelReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/einvoice/"+uuid.New().String()+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", uuid.New().String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.CancelIRN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCancelIRN_NotFound(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()
	fakeID := uuid.New().String()

	cancelReq := domain.CancelIRNRequest{Reason: "1", Remark: "test"}
	body, _ := json.Marshal(cancelReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/einvoice/"+fakeID+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", fakeID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.CancelIRN(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Get, List, Summary
// ---------------------------------------------------------------------------

func TestGetEInvoice_NotFound(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()
	fakeID := uuid.New().String()

	req := getWithTenant(t, "/v1/einvoice/"+fakeID, tenantID)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", fakeID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.GetEInvoice(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetEInvoiceByIRN_NotFound(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()

	req := getWithTenant(t, "/v1/einvoice/irn/nonexistent", tenantID)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("irn", "nonexistent")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.GetEInvoiceByIRN(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListEInvoices_MissingGSTIN(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()

	rec := httptest.NewRecorder()
	h.ListEInvoices(rec, getWithTenant(t, "/v1/einvoice/list", tenantID))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetSummary_MissingGSTIN(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()

	rec := httptest.NewRecorder()
	h.GetSummary(rec, getWithTenant(t, "/v1/einvoice/summary", tenantID))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: CancellationWindowOpen
// ---------------------------------------------------------------------------

func TestCancellationWindowOpen(t *testing.T) {
	now := time.Now()
	clock := fixedClock{now: now}

	withinWindow := now.Add(-12 * time.Hour)
	assert.True(t, store.CancellationWindowOpen(&withinWindow, clock))

	expired := now.Add(-25 * time.Hour)
	assert.False(t, store.CancellationWindowOpen(&expired, clock))

	assert.False(t, store.CancellationWindowOpen(nil, clock))

	exactBoundary := now.Add(-24 * time.Hour)
	assert.False(t, store.CancellationWindowOpen(&exactBoundary, clock))

	justInside := now.Add(-23*time.Hour - 59*time.Minute)
	assert.True(t, store.CancellationWindowOpen(&justInside, clock))
}

// ---------------------------------------------------------------------------
// Tests: ValidityDaysForDistance (used by EWB, tested here for shared fn)
// ---------------------------------------------------------------------------

func TestValidityDaysForDistance(t *testing.T) {
	tests := []struct {
		km   int
		days int
	}{
		{0, 1},
		{100, 1},
		{200, 1},
		{201, 2},
		{400, 2},
		{401, 3},
		{600, 3},
		{1500, 8},
		{1800, 9},
		{2000, 10},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%dkm", tt.km), func(t *testing.T) {
			assert.Equal(t, tt.days, store.ValidityDaysForDistance(tt.km))
		})
	}
}

// ---------------------------------------------------------------------------
// Tests: Router
// ---------------------------------------------------------------------------

func TestNewRouter(t *testing.T) {
	mockStore := store.NewMockStore()
	irpServer := newMockIRPServer(t)
	defer irpServer.Close()
	irpClient := newTestIRPClient(irpServer.URL)

	r := NewRouter(mockStore, irpClient, store.RealClock{})
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
