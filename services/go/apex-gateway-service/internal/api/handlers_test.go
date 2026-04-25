package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/apex-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/apex-gateway-service/internal/provider"
)

func newTestHandlers() *Handlers {
	return NewHandlers(provider.NewMockProvider())
}

func parseGatewayResponse(t *testing.T, body []byte, target interface{}) domain.ResponseMeta {
	t.Helper()
	var wrapper httputil.SuccessResponse
	var gw domain.GatewayResponse
	gw.Data = target
	wrapper.Data = &gw
	require.NoError(t, json.Unmarshal(body, &wrapper))
	return gw.Meta
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
	req.Header.Set("X-Idempotency-Key", uuid.New().String())
	return req
}

// ---------------------------------------------------------------------------
// Tests: Health
// ---------------------------------------------------------------------------

func TestHealth(t *testing.T) {
	h := newTestHandlers()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	h.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "apex-gateway-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: FetchVendors
// ---------------------------------------------------------------------------

func TestFetchVendors_ReturnsAll(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	req := domain.FetchVendorsRequest{RequestID: uuid.New().String()}
	rec := httptest.NewRecorder()
	h.FetchVendors(rec, postJSON(t, "/v1/gateway/apex/vendors", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.FetchVendorsResponse
	meta := parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, 50, resp.Total)
	assert.Len(t, resp.Vendors, 50)
	assert.Equal(t, "success", meta.ProviderStatus)
}

func TestFetchVendors_WithLimit(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	req := domain.FetchVendorsRequest{Limit: 10, RequestID: uuid.New().String()}
	rec := httptest.NewRecorder()
	h.FetchVendors(rec, postJSON(t, "/v1/gateway/apex/vendors", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.FetchVendorsResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.Len(t, resp.Vendors, 10)
	assert.Equal(t, 50, resp.Total)
}

func TestFetchVendors_WithOffset(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	req := domain.FetchVendorsRequest{Offset: 48, RequestID: uuid.New().String()}
	rec := httptest.NewRecorder()
	h.FetchVendors(rec, postJSON(t, "/v1/gateway/apex/vendors", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.FetchVendorsResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.Len(t, resp.Vendors, 2)
	assert.Equal(t, 50, resp.Total)
}

func TestFetchVendors_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.FetchVendorsRequest{RequestID: uuid.New().String()})
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/apex/vendors", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.FetchVendors(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFetchVendors_InvalidTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.FetchVendorsRequest{RequestID: uuid.New().String()})
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/apex/vendors", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", "not-a-uuid")
	rec := httptest.NewRecorder()
	h.FetchVendors(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFetchVendors_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/apex/vendors", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.FetchVendors(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: FetchAPInvoices
// ---------------------------------------------------------------------------

func TestFetchAPInvoices_ReturnsAll(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	req := domain.FetchAPInvoicesRequest{RequestID: uuid.New().String()}
	rec := httptest.NewRecorder()
	h.FetchAPInvoices(rec, postJSON(t, "/v1/gateway/apex/ap-invoices", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.FetchAPInvoicesResponse
	meta := parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.Greater(t, resp.Total, 150)
	assert.Equal(t, "success", meta.ProviderStatus)
}

func TestFetchAPInvoices_FilterByVendor(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	req := domain.FetchAPInvoicesRequest{VendorID: "VND-001", RequestID: uuid.New().String()}
	rec := httptest.NewRecorder()
	h.FetchAPInvoices(rec, postJSON(t, "/v1/gateway/apex/ap-invoices", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.FetchAPInvoicesResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.Greater(t, resp.Total, 0)
	for _, inv := range resp.Invoices {
		assert.Equal(t, "VND-001", inv.VendorID)
	}
}

func TestFetchAPInvoices_FilterByVendor_NoMatch(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	req := domain.FetchAPInvoicesRequest{VendorID: "VND-999", RequestID: uuid.New().String()}
	rec := httptest.NewRecorder()
	h.FetchAPInvoices(rec, postJSON(t, "/v1/gateway/apex/ap-invoices", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.FetchAPInvoicesResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Total)
}

func TestFetchAPInvoices_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.FetchAPInvoicesRequest{RequestID: uuid.New().String()})
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/apex/ap-invoices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.FetchAPInvoices(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFetchAPInvoices_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/apex/ap-invoices", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.FetchAPInvoices(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Router
// ---------------------------------------------------------------------------

func TestNewRouter(t *testing.T) {
	p := provider.NewMockProvider()
	r := NewRouter(p)
	require.NotNil(t, r)

	// Health endpoint
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Ping (heartbeat) endpoint
	req = httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRouter_VendorsEndpoint(t *testing.T) {
	p := provider.NewMockProvider()
	r := NewRouter(p)

	body, _ := json.Marshal(domain.FetchVendorsRequest{Limit: 5, RequestID: uuid.New().String()})
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/apex/vendors", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	req.Header.Set("X-Idempotency-Key", uuid.New().String())
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRouter_APInvoicesEndpoint(t *testing.T) {
	p := provider.NewMockProvider()
	r := NewRouter(p)

	body, _ := json.Marshal(domain.FetchAPInvoicesRequest{VendorID: "VND-001", RequestID: uuid.New().String()})
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/apex/ap-invoices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	req.Header.Set("X-Idempotency-Key", uuid.New().String())
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: extractHeaders
// ---------------------------------------------------------------------------

func TestExtractHeaders_Valid(t *testing.T) {
	tid := uuid.New().String()
	idem := uuid.New().String()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X-Tenant-Id", tid)
	req.Header.Set("X-Idempotency-Key", idem)

	gotTenant, gotIdem, err := extractHeaders(req)
	require.NoError(t, err)
	assert.Equal(t, tid, gotTenant)
	assert.Equal(t, idem, gotIdem)
}

func TestExtractHeaders_MissingTenant(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	_, _, err := extractHeaders(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing X-Tenant-Id")
}

func TestExtractHeaders_InvalidTenant(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X-Tenant-Id", "bad")
	_, _, err := extractHeaders(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid X-Tenant-Id")
}

func TestExtractHeaders_GeneratesIdempotencyKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	_, idem, err := extractHeaders(req)
	require.NoError(t, err)
	assert.NotEmpty(t, idem)
	_, err = uuid.Parse(idem)
	assert.NoError(t, err)
}
