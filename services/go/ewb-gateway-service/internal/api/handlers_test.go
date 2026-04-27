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
	"github.com/complai/complai/services/go/ewb-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/ewb-gateway-service/internal/provider"
)

func newTestRouter() *testEnv {
	p := provider.NewMockProvider()
	router := NewRouter(p)
	return &testEnv{router: router, provider: p}
}

type testEnv struct {
	router   http.Handler
	provider *provider.MockProvider
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

type gatewayData struct {
	Data json.RawMessage `json:"data"`
	Meta json.RawMessage `json:"meta"`
}

func parseGatewayResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	var outer httputil.SuccessResponse
	var gw gatewayData
	outer.Data = &gw
	require.NoError(t, json.Unmarshal(body, &outer))
	require.NoError(t, json.Unmarshal(gw.Data, target))
}

// ---------------------------------------------------------------------------
// Tests: Health
// ---------------------------------------------------------------------------

func TestHealth(t *testing.T) {
	env := newTestRouter()
	rec := httptest.NewRecorder()
	env.router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Full EWB Lifecycle
// ---------------------------------------------------------------------------

func TestEWB_FullLifecycle(t *testing.T) {
	env := newTestRouter()
	tenantID := uuid.New().String()

	// Generate
	genReq := domain.GenerateEWBRequest{
		GSTIN:         "29AABCA1234A1Z5",
		SupplyType:    "O",
		DocType:       "INV",
		DocNo:         "INV-001",
		DocDate:       "15/04/2026",
		FromGSTIN:     "29AABCA1234A1Z5",
		FromName:      "Seller Corp",
		ToGSTIN:       "27AABCB5678B1Z3",
		ToName:        "Buyer Inc",
		VehicleNo:     "KA01AB1234",
		VehicleType:   "R",
		DistanceKM:    450,
		TransportMode: "1",
		TotalValue:    100000,
		TaxableValue:  84745.76,
		IGSTAmount:    15254.24,
	}
	rec := httptest.NewRecorder()
	env.router.ServeHTTP(rec, postJSON(t, "/v1/gateway/ewb/generate", genReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var genResp domain.GenerateEWBResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &genResp)
	assert.NotEmpty(t, genResp.EWBNumber)
	assert.Equal(t, "ACT", genResp.Status)
	ewbNo := genResp.EWBNumber

	// Get
	rec = httptest.NewRecorder()
	env.router.ServeHTTP(rec, getWithTenant(t, "/v1/gateway/ewb/?ewb_no="+ewbNo, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var getResp domain.GetEWBResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &getResp)
	assert.Equal(t, ewbNo, getResp.EWBNumber)
	assert.Equal(t, "ACT", getResp.Status)
	assert.Equal(t, 450, getResp.DistanceKM)

	// Update Vehicle
	vehReq := domain.UpdateVehicleRequest{
		EWBNo:         ewbNo,
		VehicleNo:     "MH02CD5678",
		FromPlace:     "Pune",
		FromState:     "27",
		Reason:        "2",
		TransportMode: "1",
	}
	rec = httptest.NewRecorder()
	env.router.ServeHTTP(rec, postJSON(t, "/v1/gateway/ewb/vehicle", vehReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var vehResp domain.UpdateVehicleResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &vehResp)
	assert.Equal(t, "MH02CD5678", vehResp.VehicleNo)

	// Extend Validity
	extReq := domain.ExtendValidityRequest{
		EWBNo:             ewbNo,
		FromPlace:         "Pune",
		FromState:         "27",
		RemainingDistance: 300,
		ExtendReason:      "4",
		TransitType:       "R",
		ConsignmentStatus: "M",
	}
	rec = httptest.NewRecorder()
	env.router.ServeHTTP(rec, postJSON(t, "/v1/gateway/ewb/extend", extReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var extResp domain.ExtendValidityResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &extResp)
	assert.Equal(t, "ACT", extResp.Status)

	// Cancel
	cancelReq := domain.CancelEWBRequest{
		EWBNo:  ewbNo,
		Reason: "2",
		Remark: "Duplicate entry",
	}
	rec = httptest.NewRecorder()
	env.router.ServeHTTP(rec, postJSON(t, "/v1/gateway/ewb/cancel", cancelReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var cancelResp domain.CancelEWBResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &cancelResp)
	assert.Equal(t, "CNL", cancelResp.Status)

	// Verify cancelled
	rec = httptest.NewRecorder()
	env.router.ServeHTTP(rec, getWithTenant(t, "/v1/gateway/ewb/?ewb_no="+ewbNo, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)
	parseGatewayResponse(t, rec.Body.Bytes(), &getResp)
	assert.Equal(t, "CNL", getResp.Status)
}

// ---------------------------------------------------------------------------
// Tests: Idempotent Generation
// ---------------------------------------------------------------------------

func TestGenerateEWB_Idempotent(t *testing.T) {
	env := newTestRouter()
	tenantID := uuid.New().String()

	req := domain.GenerateEWBRequest{
		GSTIN: "29AABCA1234A1Z5", DocType: "INV", DocNo: "INV-IDEM",
		DocDate: "15/04/2026", VehicleType: "R", DistanceKM: 100, TotalValue: 50000,
	}

	rec1 := httptest.NewRecorder()
	env.router.ServeHTTP(rec1, postJSON(t, "/v1/gateway/ewb/generate", req, tenantID))
	require.Equal(t, http.StatusOK, rec1.Code)

	rec2 := httptest.NewRecorder()
	env.router.ServeHTTP(rec2, postJSON(t, "/v1/gateway/ewb/generate", req, tenantID))
	require.Equal(t, http.StatusOK, rec2.Code)

	var r1, r2 domain.GenerateEWBResponse
	parseGatewayResponse(t, rec1.Body.Bytes(), &r1)
	parseGatewayResponse(t, rec2.Body.Bytes(), &r2)
	assert.Equal(t, r1.EWBNumber, r2.EWBNumber)
}

// ---------------------------------------------------------------------------
// Tests: Consolidation
// ---------------------------------------------------------------------------

func TestConsolidateEWB(t *testing.T) {
	env := newTestRouter()
	tenantID := uuid.New().String()

	var ewbNos []string
	for i := 0; i < 3; i++ {
		req := domain.GenerateEWBRequest{
			GSTIN: "29AABCA1234A1Z5", DocType: "INV",
			DocNo: "INV-C-" + uuid.New().String()[:6], DocDate: "15/04/2026",
			VehicleType: "R", DistanceKM: 200, TotalValue: 10000,
		}
		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, postJSON(t, "/v1/gateway/ewb/generate", req, tenantID))
		require.Equal(t, http.StatusOK, rec.Code)
		var gr domain.GenerateEWBResponse
		parseGatewayResponse(t, rec.Body.Bytes(), &gr)
		ewbNos = append(ewbNos, gr.EWBNumber)
	}

	conReq := domain.ConsolidateEWBRequest{
		FromGSTIN: "29AABCA1234A1Z5", VehicleNo: "KA01XX9999",
		TransportMode: "1", EWBNumbers: ewbNos,
	}
	rec := httptest.NewRecorder()
	env.router.ServeHTTP(rec, postJSON(t, "/v1/gateway/ewb/consolidate", conReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var conResp domain.ConsolidateEWBResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &conResp)
	assert.NotEmpty(t, conResp.ConsolidatedEWBNo)
	assert.Equal(t, "ACT", conResp.Status)
}

func TestConsolidateEWB_TooFew(t *testing.T) {
	env := newTestRouter()
	tenantID := uuid.New().String()

	conReq := domain.ConsolidateEWBRequest{
		VehicleNo: "KA01XX9999", EWBNumbers: []string{"one"},
	}
	rec := httptest.NewRecorder()
	env.router.ServeHTTP(rec, postJSON(t, "/v1/gateway/ewb/consolidate", conReq, tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Cancel already cancelled
// ---------------------------------------------------------------------------

func TestCancelEWB_AlreadyCancelled(t *testing.T) {
	env := newTestRouter()
	tenantID := uuid.New().String()

	req := domain.GenerateEWBRequest{
		GSTIN: "29AABCA1234A1Z5", DocType: "INV", DocNo: "INV-DUP-CNL",
		DocDate: "15/04/2026", VehicleType: "R", DistanceKM: 100, TotalValue: 50000,
	}
	rec := httptest.NewRecorder()
	env.router.ServeHTTP(rec, postJSON(t, "/v1/gateway/ewb/generate", req, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)
	var gr domain.GenerateEWBResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &gr)

	cancelReq := domain.CancelEWBRequest{EWBNo: gr.EWBNumber, Reason: "1"}
	rec = httptest.NewRecorder()
	env.router.ServeHTTP(rec, postJSON(t, "/v1/gateway/ewb/cancel", cancelReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	rec = httptest.NewRecorder()
	env.router.ServeHTTP(rec, postJSON(t, "/v1/gateway/ewb/cancel", cancelReq, tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Error cases
// ---------------------------------------------------------------------------

func TestGenerateEWB_MissingTenantID(t *testing.T) {
	env := newTestRouter()
	body, _ := json.Marshal(domain.GenerateEWBRequest{})
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/ewb/generate", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	env.router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGenerateEWB_MissingGSTIN(t *testing.T) {
	env := newTestRouter()
	tenantID := uuid.New().String()
	req := domain.GenerateEWBRequest{DocNo: "INV-001"}
	rec := httptest.NewRecorder()
	env.router.ServeHTTP(rec, postJSON(t, "/v1/gateway/ewb/generate", req, tenantID))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGenerateEWB_MissingDocNo(t *testing.T) {
	env := newTestRouter()
	tenantID := uuid.New().String()
	req := domain.GenerateEWBRequest{GSTIN: "29AABCA1234A1Z5"}
	rec := httptest.NewRecorder()
	env.router.ServeHTTP(rec, postJSON(t, "/v1/gateway/ewb/generate", req, tenantID))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetEWB_MissingParam(t *testing.T) {
	env := newTestRouter()
	tenantID := uuid.New().String()
	rec := httptest.NewRecorder()
	env.router.ServeHTTP(rec, getWithTenant(t, "/v1/gateway/ewb/", tenantID))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetEWB_NotFound(t *testing.T) {
	env := newTestRouter()
	tenantID := uuid.New().String()
	rec := httptest.NewRecorder()
	env.router.ServeHTTP(rec, getWithTenant(t, "/v1/gateway/ewb/?ewb_no=999999999999", tenantID))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Router + ExtractHeaders
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

func TestExtractHeaders_Valid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	tid := uuid.New()
	req.Header.Set("X-Tenant-Id", tid.String())
	req.Header.Set("X-Idempotency-Key", "my-key")
	h, err := extractHeaders(req)
	require.NoError(t, err)
	assert.Equal(t, tid, h.tenantID)
	assert.Equal(t, "my-key", h.idempotencyKey)
}

func TestExtractHeaders_MissingTenant(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := extractHeaders(req)
	assert.Error(t, err)
}

func TestExtractHeaders_GeneratesIdempotencyKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	h, err := extractHeaders(req)
	require.NoError(t, err)
	assert.NotEmpty(t, h.idempotencyKey)
}
