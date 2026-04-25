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
	"github.com/complai/complai/services/go/gstn-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/gstn-gateway-service/internal/provider"
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

func parseDataResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	var wrapper httputil.SuccessResponse
	wrapper.Data = target
	require.NoError(t, json.Unmarshal(body, &wrapper))
}

func TestHealth(t *testing.T) {
	h := newTestHandlers()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	h.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "gstn-gateway-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: Authenticate
// ---------------------------------------------------------------------------

func TestAuthenticate(t *testing.T) {
	h := newTestHandlers()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/adaequare/authenticate", nil)
	h.Authenticate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var data domain.AuthResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &data)
	assert.Contains(t, data.AccessToken, "mock-gsp-token-")
	assert.Equal(t, "bearer", data.TokenType)
	assert.Equal(t, 86399, data.ExpiresIn)
	assert.Equal(t, "gsp", data.Scope)
}

// ---------------------------------------------------------------------------
// Tests: Full GSTR-1 Lifecycle (save → get → reset → save → submit → file)
// ---------------------------------------------------------------------------

func TestGSTR1_FullLifecycle(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	gstin := "29AABCA1234A1Z5"
	period := "042026"

	// Step 1: Save B2B section
	saveReq := domain.GSTR1SaveRequest{
		GSTIN:     gstin,
		RetPeriod: period,
		Section:   "b2b",
		Data:      map[string]interface{}{"invoices": []interface{}{map[string]interface{}{"inum": "INV001"}}},
		RequestID: uuid.New().String(),
	}
	rec := httptest.NewRecorder()
	h.GSTR1Save(rec, postJSON(t, "/v1/gateway/adaequare/gstr1/save", saveReq, tenantID))
	assert.Equal(t, http.StatusOK, rec.Code)
	var saveResp domain.GSTR1SaveResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &saveResp)
	assert.Equal(t, "success", saveResp.Status)
	assert.NotEmpty(t, saveResp.Token)

	// Step 2: Get saved data
	getReq := domain.GSTR1GetRequest{GSTIN: gstin, RetPeriod: period, RequestID: uuid.New().String()}
	rec = httptest.NewRecorder()
	h.GSTR1Get(rec, postJSON(t, "/v1/gateway/adaequare/gstr1/get", getReq, tenantID))
	assert.Equal(t, http.StatusOK, rec.Code)
	var getResp domain.GSTR1GetResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &getResp)
	assert.Equal(t, "saved", getResp.Status)
	assert.Contains(t, getResp.Data, "b2b")

	// Step 3: Reset
	resetReq := domain.GSTR1ResetRequest{GSTIN: gstin, RetPeriod: period, RequestID: uuid.New().String()}
	rec = httptest.NewRecorder()
	h.GSTR1Reset(rec, postJSON(t, "/v1/gateway/adaequare/gstr1/reset", resetReq, tenantID))
	assert.Equal(t, http.StatusOK, rec.Code)
	var resetResp domain.GSTR1ResetResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resetResp)
	assert.Equal(t, "success", resetResp.Status)

	// Step 4: Save again after reset
	saveReq2 := domain.GSTR1SaveRequest{
		GSTIN:     gstin,
		RetPeriod: period,
		Section:   "b2b",
		Data:      map[string]interface{}{"invoices": []interface{}{map[string]interface{}{"inum": "INV002"}}},
		RequestID: uuid.New().String(),
	}
	rec = httptest.NewRecorder()
	h.GSTR1Save(rec, postJSON(t, "/v1/gateway/adaequare/gstr1/save", saveReq2, tenantID))
	assert.Equal(t, http.StatusOK, rec.Code)

	// Step 5: Submit
	submitReq := domain.GSTR1SubmitRequest{GSTIN: gstin, RetPeriod: period, RequestID: uuid.New().String()}
	rec = httptest.NewRecorder()
	h.GSTR1Submit(rec, postJSON(t, "/v1/gateway/adaequare/gstr1/submit", submitReq, tenantID))
	assert.Equal(t, http.StatusOK, rec.Code)
	var submitResp domain.GSTR1SubmitResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &submitResp)
	assert.Equal(t, "success", submitResp.Status)

	// Step 6: File with EVC
	fileReq := domain.GSTR1FileRequest{
		GSTIN:     gstin,
		RetPeriod: period,
		SignType:  "EVC",
		EVOTP:     "123456",
		PAN:       "AABCA1234A",
		RequestID: uuid.New().String(),
	}
	rec = httptest.NewRecorder()
	h.GSTR1File(rec, postJSON(t, "/v1/gateway/adaequare/gstr1/file", fileReq, tenantID))
	assert.Equal(t, http.StatusOK, rec.Code)
	var fileResp domain.GSTR1FileResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &fileResp)
	assert.Equal(t, "success", fileResp.Status)
	assert.NotEmpty(t, fileResp.ARN)
	assert.Contains(t, fileResp.ARN, "AA29")

	// Step 7: Status — should show filed with ARN
	statusReq := domain.GSTR1StatusRequest{GSTIN: gstin, RetPeriod: period, RequestID: uuid.New().String()}
	rec = httptest.NewRecorder()
	h.GSTR1Status(rec, postJSON(t, "/v1/gateway/adaequare/gstr1/status", statusReq, tenantID))
	assert.Equal(t, http.StatusOK, rec.Code)
	var statusResp domain.GSTR1StatusResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &statusResp)
	assert.Equal(t, "filed", statusResp.Status)
	assert.Equal(t, fileResp.ARN, statusResp.ARN)
	assert.NotNil(t, statusResp.FiledAt)
}

// ---------------------------------------------------------------------------
// Tests: Idempotency
// ---------------------------------------------------------------------------

func TestGSTR1Save_Idempotency(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	requestID := uuid.New().String()

	req := domain.GSTR1SaveRequest{
		GSTIN:     "29AABCA1234A1Z5",
		RetPeriod: "042026",
		Section:   "b2b",
		Data:      map[string]interface{}{"invoices": []interface{}{}},
		RequestID: requestID,
	}

	rec1 := httptest.NewRecorder()
	h.GSTR1Save(rec1, postJSON(t, "/v1/gateway/adaequare/gstr1/save", req, tenantID))
	assert.Equal(t, http.StatusOK, rec1.Code)

	rec2 := httptest.NewRecorder()
	h.GSTR1Save(rec2, postJSON(t, "/v1/gateway/adaequare/gstr1/save", req, tenantID))
	assert.Equal(t, http.StatusOK, rec2.Code)

	assert.Equal(t, rec1.Body.String(), rec2.Body.String())
}

func TestGSTR1Submit_Idempotency(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	gstin := "29AABCA1234A1Z5"
	period := "042026"

	saveReq := domain.GSTR1SaveRequest{
		GSTIN: gstin, RetPeriod: period, Section: "b2b",
		Data: map[string]interface{}{"x": 1}, RequestID: uuid.New().String(),
	}
	rec := httptest.NewRecorder()
	h.GSTR1Save(rec, postJSON(t, "/v1/gateway/adaequare/gstr1/save", saveReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	submitRequestID := uuid.New().String()
	submitReq := domain.GSTR1SubmitRequest{GSTIN: gstin, RetPeriod: period, RequestID: submitRequestID}

	rec1 := httptest.NewRecorder()
	h.GSTR1Submit(rec1, postJSON(t, "/v1/gateway/adaequare/gstr1/submit", submitReq, tenantID))
	assert.Equal(t, http.StatusOK, rec1.Code)

	rec2 := httptest.NewRecorder()
	h.GSTR1Submit(rec2, postJSON(t, "/v1/gateway/adaequare/gstr1/submit", submitReq, tenantID))
	assert.Equal(t, http.StatusOK, rec2.Code)

	assert.Equal(t, rec1.Body.String(), rec2.Body.String())
}

func TestGSTR1File_Idempotency(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	gstin := "29AABCA1234A1Z5"
	period := "042026"

	saveReq := domain.GSTR1SaveRequest{
		GSTIN: gstin, RetPeriod: period, Section: "b2b",
		Data: map[string]interface{}{"x": 1}, RequestID: uuid.New().String(),
	}
	rec := httptest.NewRecorder()
	h.GSTR1Save(rec, postJSON(t, "/v1/gateway/adaequare/gstr1/save", saveReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	submitReq := domain.GSTR1SubmitRequest{GSTIN: gstin, RetPeriod: period, RequestID: uuid.New().String()}
	rec = httptest.NewRecorder()
	h.GSTR1Submit(rec, postJSON(t, "/v1/gateway/adaequare/gstr1/submit", submitReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	fileRequestID := uuid.New().String()
	fileReq := domain.GSTR1FileRequest{
		GSTIN: gstin, RetPeriod: period, SignType: "EVC", EVOTP: "123456",
		PAN: "AABCA1234A", RequestID: fileRequestID,
	}

	rec1 := httptest.NewRecorder()
	h.GSTR1File(rec1, postJSON(t, "/v1/gateway/adaequare/gstr1/file", fileReq, tenantID))
	assert.Equal(t, http.StatusOK, rec1.Code)

	rec2 := httptest.NewRecorder()
	h.GSTR1File(rec2, postJSON(t, "/v1/gateway/adaequare/gstr1/file", fileReq, tenantID))
	assert.Equal(t, http.StatusOK, rec2.Code)

	assert.Equal(t, rec1.Body.String(), rec2.Body.String())
}

// ---------------------------------------------------------------------------
// Tests: Error cases
// ---------------------------------------------------------------------------

func TestGSTR1Save_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.GSTR1Save(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1Save_InvalidTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", "not-uuid")
	rec := httptest.NewRecorder()
	h.GSTR1Save(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1Save_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GSTR1Save(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1File_BeforeSubmit(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	saveReq := domain.GSTR1SaveRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", Section: "b2b",
		Data: map[string]interface{}{"x": 1}, RequestID: uuid.New().String(),
	}
	rec := httptest.NewRecorder()
	h.GSTR1Save(rec, postJSON(t, "/", saveReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	fileReq := domain.GSTR1FileRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", SignType: "EVC", EVOTP: "123456",
		PAN: "AABCA1234A", RequestID: uuid.New().String(),
	}
	rec = httptest.NewRecorder()
	h.GSTR1File(rec, postJSON(t, "/", fileReq, tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

func TestGSTR1File_InvalidSignType(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	gstin := "29AABCA1234A1Z5"
	period := "042026"

	saveReq := domain.GSTR1SaveRequest{
		GSTIN: gstin, RetPeriod: period, Section: "b2b",
		Data: map[string]interface{}{"x": 1}, RequestID: uuid.New().String(),
	}
	rec := httptest.NewRecorder()
	h.GSTR1Save(rec, postJSON(t, "/", saveReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	submitReq := domain.GSTR1SubmitRequest{GSTIN: gstin, RetPeriod: period, RequestID: uuid.New().String()}
	rec = httptest.NewRecorder()
	h.GSTR1Submit(rec, postJSON(t, "/", submitReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	fileReq := domain.GSTR1FileRequest{
		GSTIN: gstin, RetPeriod: period, SignType: "INVALID",
		PAN: "AABCA1234A", RequestID: uuid.New().String(),
	}
	rec = httptest.NewRecorder()
	h.GSTR1File(rec, postJSON(t, "/", fileReq, tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

func TestGSTR1Submit_NoSections(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	submitReq := domain.GSTR1SubmitRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", RequestID: uuid.New().String(),
	}
	rec := httptest.NewRecorder()
	h.GSTR1Submit(rec, postJSON(t, "/", submitReq, tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

func TestGSTR1Reset_AfterFiled(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	gstin := "29AABCA1234A1Z5"
	period := "042026"

	saveReq := domain.GSTR1SaveRequest{
		GSTIN: gstin, RetPeriod: period, Section: "b2b",
		Data: map[string]interface{}{"x": 1}, RequestID: uuid.New().String(),
	}
	rec := httptest.NewRecorder()
	h.GSTR1Save(rec, postJSON(t, "/", saveReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	submitReq := domain.GSTR1SubmitRequest{GSTIN: gstin, RetPeriod: period, RequestID: uuid.New().String()}
	rec = httptest.NewRecorder()
	h.GSTR1Submit(rec, postJSON(t, "/", submitReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	fileReq := domain.GSTR1FileRequest{
		GSTIN: gstin, RetPeriod: period, SignType: "EVC", EVOTP: "123456",
		PAN: "AABCA1234A", RequestID: uuid.New().String(),
	}
	rec = httptest.NewRecorder()
	h.GSTR1File(rec, postJSON(t, "/", fileReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	resetReq := domain.GSTR1ResetRequest{GSTIN: gstin, RetPeriod: period, RequestID: uuid.New().String()}
	rec = httptest.NewRecorder()
	h.GSTR1Reset(rec, postJSON(t, "/", resetReq, tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

func TestGSTR1Get_EmptyFiling(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	getReq := domain.GSTR1GetRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", RequestID: uuid.New().String(),
	}
	rec := httptest.NewRecorder()
	h.GSTR1Get(rec, postJSON(t, "/", getReq, tenantID))
	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.GSTR1GetResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, "draft", resp.Status)
	assert.Empty(t, resp.Data)
}

func TestGSTR1Get_SpecificSection(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	gstin := "29AABCA1234A1Z5"
	period := "042026"

	saveB2B := domain.GSTR1SaveRequest{
		GSTIN: gstin, RetPeriod: period, Section: "b2b",
		Data: map[string]interface{}{"b2b_data": true}, RequestID: uuid.New().String(),
	}
	rec := httptest.NewRecorder()
	h.GSTR1Save(rec, postJSON(t, "/", saveB2B, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	saveHSN := domain.GSTR1SaveRequest{
		GSTIN: gstin, RetPeriod: period, Section: "hsn",
		Data: map[string]interface{}{"hsn_data": true}, RequestID: uuid.New().String(),
	}
	rec = httptest.NewRecorder()
	h.GSTR1Save(rec, postJSON(t, "/", saveHSN, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	getReq := domain.GSTR1GetRequest{GSTIN: gstin, RetPeriod: period, Section: "b2b", RequestID: uuid.New().String()}
	rec = httptest.NewRecorder()
	h.GSTR1Get(rec, postJSON(t, "/", getReq, tenantID))
	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.GSTR1GetResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.Contains(t, resp.Data, "b2b")
	assert.NotContains(t, resp.Data, "hsn")
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

func TestGSTR1File_EVCMissingOTP(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	gstin := "29AABCA1234A1Z5"
	period := "042026"

	saveReq := domain.GSTR1SaveRequest{
		GSTIN: gstin, RetPeriod: period, Section: "b2b",
		Data: map[string]interface{}{"x": 1}, RequestID: uuid.New().String(),
	}
	rec := httptest.NewRecorder()
	h.GSTR1Save(rec, postJSON(t, "/", saveReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	submitReq := domain.GSTR1SubmitRequest{GSTIN: gstin, RetPeriod: period, RequestID: uuid.New().String()}
	rec = httptest.NewRecorder()
	h.GSTR1Submit(rec, postJSON(t, "/", submitReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	fileReq := domain.GSTR1FileRequest{
		GSTIN: gstin, RetPeriod: period, SignType: "EVC", EVOTP: "",
		PAN: "AABCA1234A", RequestID: uuid.New().String(),
	}
	rec = httptest.NewRecorder()
	h.GSTR1File(rec, postJSON(t, "/", fileReq, tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

func TestGSTR1Get_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GSTR1Get(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1Get_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.GSTR1GetRequest{GSTIN: "X", RetPeriod: "042026"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.GSTR1Get(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1Reset_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GSTR1Reset(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1Reset_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.GSTR1ResetRequest{GSTIN: "X", RetPeriod: "042026"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.GSTR1Reset(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1Submit_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GSTR1Submit(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1Submit_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.GSTR1Submit(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1File_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GSTR1File(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1File_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.GSTR1FileRequest{GSTIN: "X", RetPeriod: "042026"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.GSTR1File(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1Status_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GSTR1Status(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1Status_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.GSTR1StatusRequest{GSTIN: "X", RetPeriod: "042026"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.GSTR1Status(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGSTR1Save_AfterSubmitted(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	gstin := "29AABCA1234A1Z5"
	period := "042026"

	saveReq := domain.GSTR1SaveRequest{
		GSTIN: gstin, RetPeriod: period, Section: "b2b",
		Data: map[string]interface{}{"x": 1}, RequestID: uuid.New().String(),
	}
	rec := httptest.NewRecorder()
	h.GSTR1Save(rec, postJSON(t, "/", saveReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	submitReq := domain.GSTR1SubmitRequest{GSTIN: gstin, RetPeriod: period, RequestID: uuid.New().String()}
	rec = httptest.NewRecorder()
	h.GSTR1Submit(rec, postJSON(t, "/", submitReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	saveReq2 := domain.GSTR1SaveRequest{
		GSTIN: gstin, RetPeriod: period, Section: "hsn",
		Data: map[string]interface{}{"y": 2}, RequestID: uuid.New().String(),
	}
	rec = httptest.NewRecorder()
	h.GSTR1Save(rec, postJSON(t, "/", saveReq2, tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}
