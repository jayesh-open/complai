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
	"github.com/complai/complai/services/go/kyc-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/kyc-gateway-service/internal/provider"
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
	assert.Equal(t, "kyc-gateway-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: VerifyPAN
// ---------------------------------------------------------------------------

func TestVerifyPAN_Valid(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	rec := httptest.NewRecorder()

	req := domain.PANVerifyRequest{
		PAN:       "AABCA1234A",
		Name:      "Test Company",
		RequestID: uuid.New().String(),
	}
	h.VerifyPAN(rec, postJSON(t, "/v1/gateway/kyc/pan/verify", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.PANVerifyResponse
	meta := parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.True(t, resp.Valid)
	assert.Equal(t, "valid", resp.Status)
	assert.Equal(t, "Company", resp.Category)
	assert.Equal(t, "valid", meta.ProviderStatus)
}

func TestVerifyPAN_Invalid(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	rec := httptest.NewRecorder()

	req := domain.PANVerifyRequest{
		PAN:       "123",
		Name:      "Bad",
		RequestID: uuid.New().String(),
	}
	h.VerifyPAN(rec, postJSON(t, "/v1/gateway/kyc/pan/verify", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.PANVerifyResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.False(t, resp.Valid)
	assert.Equal(t, "invalid", resp.Status)
}

func TestVerifyPAN_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.PANVerifyRequest{PAN: "AABCA1234A"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.VerifyPAN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVerifyPAN_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.VerifyPAN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: VerifyGSTIN
// ---------------------------------------------------------------------------

func TestVerifyGSTIN_Valid(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	rec := httptest.NewRecorder()

	req := domain.GSTINVerifyRequest{
		GSTIN:     "29AABCA1234A1Z5",
		RequestID: uuid.New().String(),
	}
	h.VerifyGSTIN(rec, postJSON(t, "/v1/gateway/kyc/gstin/verify", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.GSTINVerifyResponse
	meta := parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.True(t, resp.Valid)
	assert.Equal(t, "Active", resp.Status)
	assert.Equal(t, "Karnataka", resp.State)
	assert.Equal(t, "AABCA1234A", resp.PAN)
	assert.Equal(t, "Active", meta.ProviderStatus)
}

func TestVerifyGSTIN_Invalid(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	rec := httptest.NewRecorder()

	req := domain.GSTINVerifyRequest{
		GSTIN:     "XX",
		RequestID: uuid.New().String(),
	}
	h.VerifyGSTIN(rec, postJSON(t, "/v1/gateway/kyc/gstin/verify", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.GSTINVerifyResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.False(t, resp.Valid)
	assert.Equal(t, "invalid", resp.Status)
}

func TestVerifyGSTIN_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.GSTINVerifyRequest{GSTIN: "29AABCA1234A1Z5"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.VerifyGSTIN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVerifyGSTIN_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.VerifyGSTIN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: VerifyTAN
// ---------------------------------------------------------------------------

func TestVerifyTAN_Valid(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	rec := httptest.NewRecorder()

	req := domain.TANVerifyRequest{
		TAN:       "BLRA12345B",
		RequestID: uuid.New().String(),
	}
	h.VerifyTAN(rec, postJSON(t, "/v1/gateway/kyc/tan/verify", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.TANVerifyResponse
	meta := parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.True(t, resp.Valid)
	assert.Equal(t, "valid", resp.Status)
	assert.Equal(t, "valid", meta.ProviderStatus)
}

func TestVerifyTAN_Invalid(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	rec := httptest.NewRecorder()

	req := domain.TANVerifyRequest{
		TAN:       "123",
		RequestID: uuid.New().String(),
	}
	h.VerifyTAN(rec, postJSON(t, "/v1/gateway/kyc/tan/verify", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.TANVerifyResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.False(t, resp.Valid)
	assert.Equal(t, "invalid", resp.Status)
}

func TestVerifyTAN_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.TANVerifyRequest{TAN: "BLRA12345B"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.VerifyTAN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVerifyTAN_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.VerifyTAN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: VerifyBank
// ---------------------------------------------------------------------------

func TestVerifyBank_Valid(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	rec := httptest.NewRecorder()

	req := domain.BankVerifyRequest{
		AccountNumber: "1234567890",
		IFSC:          "SBIN0001234",
		RequestID:     uuid.New().String(),
	}
	h.VerifyBank(rec, postJSON(t, "/v1/gateway/kyc/bank/verify", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.BankVerifyResponse
	meta := parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.True(t, resp.Valid)
	assert.Equal(t, "State Bank of India", resp.BankName)
	assert.Equal(t, "success", meta.ProviderStatus)
}

func TestVerifyBank_Invalid(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	rec := httptest.NewRecorder()

	req := domain.BankVerifyRequest{
		AccountNumber: "1234567890",
		IFSC:          "123",
		RequestID:     uuid.New().String(),
	}
	h.VerifyBank(rec, postJSON(t, "/v1/gateway/kyc/bank/verify", req, tenantID))

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.BankVerifyResponse
	meta := parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.False(t, resp.Valid)
	assert.Equal(t, "invalid", meta.ProviderStatus)
}

func TestVerifyBank_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.BankVerifyRequest{AccountNumber: "123", IFSC: "SBIN0001234"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.VerifyBank(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVerifyBank_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.VerifyBank(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Invalid tenant ID format
// ---------------------------------------------------------------------------

func TestVerifyPAN_InvalidTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.PANVerifyRequest{PAN: "AABCA1234A"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", "not-uuid")
	rec := httptest.NewRecorder()
	h.VerifyPAN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Router
// ---------------------------------------------------------------------------

func TestNewRouter(t *testing.T) {
	p := provider.NewMockProvider()
	r := NewRouter(p)
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

func TestRouter_PAN(t *testing.T) {
	p := provider.NewMockProvider()
	r := NewRouter(p)

	body, _ := json.Marshal(domain.PANVerifyRequest{PAN: "AABCA1234A", Name: "Test", RequestID: "r1"})
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/kyc/pan/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
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
