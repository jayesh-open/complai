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
	"github.com/complai/complai/services/go/ewb-service/internal/domain"
	"github.com/complai/complai/services/go/ewb-service/internal/gateway"
	"github.com/complai/complai/services/go/ewb-service/internal/store"
)

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

func newMockEWBServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/gateway/ewb/generate", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)
		docNo, _ := req["doc_no"].(string)

		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]interface{}{
					"ewb_no":      fmt.Sprintf("EWB-%s", docNo),
					"ewb_date":    "15/04/2026 10:00:00",
					"valid_until": "18/04/2026 23:59:59",
					"status":      "ACT",
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

	mux.HandleFunc("/v1/gateway/ewb/cancel", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)
		ewbNo, _ := req["ewb_no"].(string)

		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]interface{}{
					"ewb_no":      ewbNo,
					"cancel_date": "15/04/2026 11:00:00",
					"status":      "CNL",
				},
				"meta": map[string]interface{}{
					"request_id": uuid.New().String(),
					"latency_ms": 5, "provider_status": "success",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/v1/gateway/ewb/vehicle", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)
		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]interface{}{
					"ewb_no":      req["ewb_no"],
					"vehicle_no":  req["vehicle_no"],
					"valid_until": "18/04/2026 23:59:59",
					"status":      "ACT",
				},
				"meta": map[string]interface{}{
					"request_id": uuid.New().String(),
					"latency_ms": 5, "provider_status": "success",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/v1/gateway/ewb/extend", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)
		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]interface{}{
					"ewb_no":      req["ewb_no"],
					"valid_until": "22/04/2026 23:59:59",
					"status":      "ACT",
				},
				"meta": map[string]interface{}{
					"request_id": uuid.New().String(),
					"latency_ms": 5, "provider_status": "success",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/v1/gateway/ewb/consolidate", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]interface{}{
					"consolidated_ewb_no": fmt.Sprintf("CEWB-%s", uuid.New().String()[:8]),
					"status":              "ACT",
				},
				"meta": map[string]interface{}{
					"request_id": uuid.New().String(),
					"latency_ms": 5, "provider_status": "success",
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
	ewbServer := newMockEWBServer(t)
	ewbClient := gateway.NewEWBClient(ewbServer.URL)
	h := NewHandlers(mockStore, ewbClient, store.RealClock{})
	return h, mockStore, ewbServer
}

func parseDataResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	var wrapper httputil.SuccessResponse
	wrapper.Data = target
	require.NoError(t, json.Unmarshal(body, &wrapper))
}

func postJSONReq(t *testing.T, path string, body interface{}, tenantID string) *http.Request {
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

func withChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func createActiveEWB(t *testing.T, mockStore *store.MockStore, tenantID uuid.UUID, docNo string) *domain.EWayBill {
	t.Helper()
	ewb := &domain.EWayBill{
		DocType:       "INV",
		DocNumber:     docNo,
		DocDate:       "15/04/2026",
		SupplierGSTIN: "29AABCA1234A1Z5",
		BuyerGSTIN:    "27AABCB5678B1Z3",
		VehicleType:   "R",
		DistanceKM:    400,
		TaxableValue:  decimal.NewFromInt(50000),
		TotalValue:    decimal.NewFromInt(59000),
		SourceSystem:  "test",
	}
	require.NoError(t, mockStore.CreateEWB(nil, tenantID, ewb))
	now := time.Now()
	validUntil := now.Add(48 * time.Hour)
	require.NoError(t, mockStore.UpdateEWBGenerated(nil, tenantID, ewb.ID, "EWB-"+docNo, now, validUntil))
	return ewb
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
	assert.Equal(t, "ewb-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: Generate EWB
// ---------------------------------------------------------------------------

func TestGenerateEWB_Success(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()

	req := domain.GenerateEWBRequest{
		DocType:       "INV",
		DocNumber:     "INV-EWB-001",
		DocDate:       "15/04/2026",
		SupplierGSTIN: "29AABCA1234A1Z5",
		SupplierName:  "Seller Corp",
		BuyerGSTIN:    "27AABCB5678B1Z3",
		BuyerName:     "Buyer Inc",
		VehicleNumber: "KA01AB1234",
		VehicleType:   "R",
		DistanceKM:    450,
		TransportMode: "1",
		TaxableValue:  decimal.NewFromInt(84746),
		IGSTAmount:    decimal.NewFromInt(15254),
		TotalValue:    decimal.NewFromInt(100000),
		SourceSystem:  "aura",
		Items: []domain.ItemRequest{
			{
				ProductName: "Steel Plates", HSNCode: "720241",
				Quantity: decimal.NewFromInt(500), Unit: "KG",
				TaxableValue: decimal.NewFromInt(84746),
				IGSTRate:     decimal.NewFromInt(18),
			},
		},
	}

	rec := httptest.NewRecorder()
	h.GenerateEWB(rec, postJSONReq(t, "/v1/ewb/generate", req, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var resp domain.GenerateEWBResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.NotEqual(t, uuid.Nil, resp.ID)
	assert.Contains(t, resp.EWBNumber, "EWB-")
	assert.Equal(t, domain.EWBStatusActive, resp.Status)
}

func TestGenerateEWB_MissingFields(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()

	req := domain.GenerateEWBRequest{DocNumber: "INV-001"}
	rec := httptest.NewRecorder()
	h.GenerateEWB(rec, postJSONReq(t, "/", req, tenantID))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGenerateEWB_MissingTenantID(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()

	body, _ := json.Marshal(domain.GenerateEWBRequest{})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.GenerateEWB(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Cancel EWB — 24h window + state machine
// ---------------------------------------------------------------------------

func TestCancelEWB_WithinWindow(t *testing.T) {
	mockStore := store.NewMockStore()
	ewbServer := newMockEWBServer(t)
	defer ewbServer.Close()
	ewbClient := gateway.NewEWBClient(ewbServer.URL)

	now := time.Now()
	clock := fixedClock{now: now}
	h := NewHandlers(mockStore, ewbClient, clock)
	tenantID := uuid.New()

	ewb := createActiveEWB(t, mockStore, tenantID, "CANCEL-OK")
	generatedAt := now.Add(-12 * time.Hour)
	mockStore.SetGeneratedAt(ewb.ID, &generatedAt)

	cancelReq := domain.CancelEWBRequest{Reason: "1", Remark: "Duplicate"}
	body, _ := json.Marshal(cancelReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/ewb/"+ewb.ID.String()+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiParam(req, "id", ewb.ID.String())

	rec := httptest.NewRecorder()
	h.CancelEWB(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp domain.CancelEWBResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, domain.EWBStatusCancelled, resp.Status)
}

func TestCancelEWB_ExpiredWindow(t *testing.T) {
	mockStore := store.NewMockStore()
	ewbServer := newMockEWBServer(t)
	defer ewbServer.Close()
	ewbClient := gateway.NewEWBClient(ewbServer.URL)

	now := time.Now()
	clock := fixedClock{now: now}
	h := NewHandlers(mockStore, ewbClient, clock)
	tenantID := uuid.New()

	ewb := createActiveEWB(t, mockStore, tenantID, "CANCEL-EXPIRED")
	expiredTime := now.Add(-25 * time.Hour)
	mockStore.SetGeneratedAt(ewb.ID, &expiredTime)

	cancelReq := domain.CancelEWBRequest{Reason: "1", Remark: "Too late"}
	body, _ := json.Marshal(cancelReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/ewb/"+ewb.ID.String()+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiParam(req, "id", ewb.ID.String())

	rec := httptest.NewRecorder()
	h.CancelEWB(rec, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCancelEWB_PendingStatus(t *testing.T) {
	mockStore := store.NewMockStore()
	ewbServer := newMockEWBServer(t)
	defer ewbServer.Close()
	ewbClient := gateway.NewEWBClient(ewbServer.URL)
	h := NewHandlers(mockStore, ewbClient, store.RealClock{})
	tenantID := uuid.New()

	ewb := &domain.EWayBill{
		DocType: "INV", DocNumber: "PENDING-001",
		SupplierGSTIN: "29AABCA1234A1Z5", TaxableValue: decimal.NewFromInt(10000),
		TotalValue: decimal.NewFromInt(11800), SourceSystem: "test",
	}
	require.NoError(t, mockStore.CreateEWB(nil, tenantID, ewb))

	cancelReq := domain.CancelEWBRequest{Reason: "1", Remark: "test"}
	body, _ := json.Marshal(cancelReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/ewb/"+ewb.ID.String()+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiParam(req, "id", ewb.ID.String())

	rec := httptest.NewRecorder()
	h.CancelEWB(rec, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCancelEWB_MissingReason(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()
	fakeID := uuid.New().String()

	cancelReq := domain.CancelEWBRequest{Remark: "No reason"}
	body, _ := json.Marshal(cancelReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/ewb/"+fakeID+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID)
	req = withChiParam(req, "id", fakeID)

	rec := httptest.NewRecorder()
	h.CancelEWB(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCancelEWB_NotFound(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()
	fakeID := uuid.New().String()

	cancelReq := domain.CancelEWBRequest{Reason: "1", Remark: "test"}
	body, _ := json.Marshal(cancelReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/ewb/"+fakeID+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID)
	req = withChiParam(req, "id", fakeID)

	rec := httptest.NewRecorder()
	h.CancelEWB(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Update Vehicle — state machine + multi-vehicle history
// ---------------------------------------------------------------------------

func TestUpdateVehicle_Success(t *testing.T) {
	h, mockStore, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New()

	ewb := createActiveEWB(t, mockStore, tenantID, "VEH-001")

	vehReq := domain.UpdateVehicleRequest{
		VehicleNumber: "MH02CD5678", FromPlace: "Pune", FromState: "27",
		TransportMode: "1", Reason: "2", Remark: "Transshipment at hub",
	}
	body, _ := json.Marshal(vehReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/ewb/"+ewb.ID.String()+"/vehicle", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiParam(req, "id", ewb.ID.String())

	rec := httptest.NewRecorder()
	h.UpdateVehicle(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp domain.UpdateVehicleResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, domain.EWBStatusVehicleUpdated, resp.Status)
	assert.Equal(t, "MH02CD5678", resp.VehicleNumber)
}

func TestUpdateVehicle_MultipleUpdates(t *testing.T) {
	h, mockStore, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New()

	ewb := createActiveEWB(t, mockStore, tenantID, "VEH-MULTI")

	vehicles := []string{"MH02CD5678", "GJ01EF9012", "RJ14GH3456"}
	for _, veh := range vehicles {
		vehReq := domain.UpdateVehicleRequest{
			VehicleNumber: veh, FromPlace: "Hub", Reason: "2", TransportMode: "1",
		}
		body, _ := json.Marshal(vehReq)
		req := httptest.NewRequest(http.MethodPost, "/v1/ewb/"+ewb.ID.String()+"/vehicle", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Tenant-Id", tenantID.String())
		req = withChiParam(req, "id", ewb.ID.String())

		rec := httptest.NewRecorder()
		h.UpdateVehicle(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	}

	updates, err := mockStore.GetVehicleUpdates(nil, tenantID, ewb.ID)
	require.NoError(t, err)
	assert.Len(t, updates, 3)
	assert.Equal(t, "MH02CD5678", updates[0].VehicleNumber)
	assert.Equal(t, "GJ01EF9012", updates[1].VehicleNumber)
	assert.Equal(t, "RJ14GH3456", updates[2].VehicleNumber)
}

func TestUpdateVehicle_CancelledStatus(t *testing.T) {
	h, mockStore, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New()

	ewb := createActiveEWB(t, mockStore, tenantID, "VEH-CANCELLED")
	require.NoError(t, mockStore.UpdateEWBCancelled(nil, tenantID, ewb.ID, "test"))

	vehReq := domain.UpdateVehicleRequest{VehicleNumber: "XX99YY0000", Reason: "1"}
	body, _ := json.Marshal(vehReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/ewb/"+ewb.ID.String()+"/vehicle", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiParam(req, "id", ewb.ID.String())

	rec := httptest.NewRecorder()
	h.UpdateVehicle(rec, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Extend Validity
// ---------------------------------------------------------------------------

func TestExtendValidity_Success(t *testing.T) {
	h, mockStore, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New()

	ewb := createActiveEWB(t, mockStore, tenantID, "EXT-001")

	extReq := domain.ExtendValidityRequest{
		RemainingDistance: 300, FromPlace: "Pune", FromState: "27",
		ExtendReason: "4", TransitType: "R", ConsignmentStatus: "M",
	}
	body, _ := json.Marshal(extReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/ewb/"+ewb.ID.String()+"/extend", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiParam(req, "id", ewb.ID.String())

	rec := httptest.NewRecorder()
	h.ExtendValidity(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp domain.ExtendValidityResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, domain.EWBStatusExtended, resp.Status)
}

func TestExtendValidity_CancelledStatus(t *testing.T) {
	h, mockStore, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New()

	ewb := createActiveEWB(t, mockStore, tenantID, "EXT-CANCELLED")
	require.NoError(t, mockStore.UpdateEWBCancelled(nil, tenantID, ewb.ID, "test"))

	extReq := domain.ExtendValidityRequest{RemainingDistance: 200}
	body, _ := json.Marshal(extReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/ewb/"+ewb.ID.String()+"/extend", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiParam(req, "id", ewb.ID.String())

	rec := httptest.NewRecorder()
	h.ExtendValidity(rec, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Consolidate
// ---------------------------------------------------------------------------

func TestConsolidateEWB_Success(t *testing.T) {
	h, mockStore, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New()

	ewb1 := createActiveEWB(t, mockStore, tenantID, "CON-001")
	ewb2 := createActiveEWB(t, mockStore, tenantID, "CON-002")
	ewb3 := createActiveEWB(t, mockStore, tenantID, "CON-003")

	conReq := domain.ConsolidateRequest{
		EWBIDS:        []uuid.UUID{ewb1.ID, ewb2.ID, ewb3.ID},
		VehicleNumber: "KA01XX9999",
		FromPlace:     "Bangalore", FromState: "29",
		ToPlace: "Mumbai", ToState: "27",
		TransportMode: "1",
	}
	body, _ := json.Marshal(conReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/ewb/consolidate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())

	rec := httptest.NewRecorder()
	h.ConsolidateEWB(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp domain.ConsolidateResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp.ConsolidatedEWBNumber)
	assert.Equal(t, 3, resp.EWBCount)

	// Verify all 3 EWBs are now CONSOLIDATED
	for _, id := range []uuid.UUID{ewb1.ID, ewb2.ID, ewb3.ID} {
		e, err := mockStore.GetEWB(nil, tenantID, id)
		require.NoError(t, err)
		assert.Equal(t, domain.EWBStatusConsolidated, e.Status)
	}
}

func TestConsolidateEWB_WithCancelledEWB(t *testing.T) {
	h, mockStore, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New()

	ewb1 := createActiveEWB(t, mockStore, tenantID, "CON-FAIL-1")
	ewb2 := createActiveEWB(t, mockStore, tenantID, "CON-FAIL-2")
	require.NoError(t, mockStore.UpdateEWBCancelled(nil, tenantID, ewb2.ID, "test"))

	conReq := domain.ConsolidateRequest{
		EWBIDS:        []uuid.UUID{ewb1.ID, ewb2.ID},
		VehicleNumber: "KA01XX0000",
	}
	body, _ := json.Marshal(conReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/ewb/consolidate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())

	rec := httptest.NewRecorder()
	h.ConsolidateEWB(rec, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestConsolidateEWB_TooFew(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()

	conReq := domain.ConsolidateRequest{
		EWBIDS: []uuid.UUID{uuid.New()}, VehicleNumber: "KA01XX9999",
	}
	body, _ := json.Marshal(conReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/ewb/consolidate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID)

	rec := httptest.NewRecorder()
	h.ConsolidateEWB(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Get, GetByNumber, VehicleHistory, List
// ---------------------------------------------------------------------------

func TestGetEWB_NotFound(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()
	fakeID := uuid.New().String()

	req := getWithTenant(t, "/v1/ewb/"+fakeID, tenantID)
	req = withChiParam(req, "id", fakeID)

	rec := httptest.NewRecorder()
	h.GetEWB(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetEWBByNumber_NotFound(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()

	req := getWithTenant(t, "/v1/ewb/number/NOEXIST", tenantID)
	req = withChiParam(req, "ewbNo", "NOEXIST")

	rec := httptest.NewRecorder()
	h.GetEWBByNumber(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetVehicleHistory(t *testing.T) {
	h, mockStore, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New()

	ewb := createActiveEWB(t, mockStore, tenantID, "VEH-HIST")

	for _, veh := range []string{"KA01AB1234", "MH02CD5678"} {
		vehReq := domain.UpdateVehicleRequest{VehicleNumber: veh, Reason: "2", TransportMode: "1"}
		body, _ := json.Marshal(vehReq)
		req := httptest.NewRequest(http.MethodPost, "/v1/ewb/"+ewb.ID.String()+"/vehicle", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Tenant-Id", tenantID.String())
		req = withChiParam(req, "id", ewb.ID.String())
		rec := httptest.NewRecorder()
		h.UpdateVehicle(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	}

	req := getWithTenant(t, "/v1/ewb/"+ewb.ID.String()+"/vehicles", tenantID.String())
	req = withChiParam(req, "id", ewb.ID.String())
	rec := httptest.NewRecorder()
	h.GetVehicleHistory(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var updates []domain.VehicleUpdate
	parseDataResponse(t, rec.Body.Bytes(), &updates)
	assert.Len(t, updates, 2)
}

func TestListEWBs_MissingGSTIN(t *testing.T) {
	h, _, srv := setupTestHandlers(t)
	defer srv.Close()
	tenantID := uuid.New().String()

	rec := httptest.NewRecorder()
	h.ListEWBs(rec, getWithTenant(t, "/v1/ewb/list", tenantID))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: State Machine (domain.CanTransitionTo)
// ---------------------------------------------------------------------------

func TestCanTransitionTo(t *testing.T) {
	tests := []struct {
		from  domain.EWBStatus
		to    domain.EWBStatus
		valid bool
	}{
		{domain.EWBStatusPending, domain.EWBStatusActive, true},
		{domain.EWBStatusPending, domain.EWBStatusCancelled, false},
		{domain.EWBStatusPending, domain.EWBStatusVehicleUpdated, false},

		{domain.EWBStatusActive, domain.EWBStatusVehicleUpdated, true},
		{domain.EWBStatusActive, domain.EWBStatusExtended, true},
		{domain.EWBStatusActive, domain.EWBStatusCancelled, true},
		{domain.EWBStatusActive, domain.EWBStatusConsolidated, true},
		{domain.EWBStatusActive, domain.EWBStatusPending, false},

		{domain.EWBStatusVehicleUpdated, domain.EWBStatusVehicleUpdated, true},
		{domain.EWBStatusVehicleUpdated, domain.EWBStatusExtended, true},
		{domain.EWBStatusVehicleUpdated, domain.EWBStatusCancelled, true},
		{domain.EWBStatusVehicleUpdated, domain.EWBStatusConsolidated, true},

		{domain.EWBStatusExtended, domain.EWBStatusVehicleUpdated, true},
		{domain.EWBStatusExtended, domain.EWBStatusExtended, true},
		{domain.EWBStatusExtended, domain.EWBStatusCancelled, true},

		{domain.EWBStatusCancelled, domain.EWBStatusActive, false},
		{domain.EWBStatusCancelled, domain.EWBStatusVehicleUpdated, false},
		{domain.EWBStatusCancelled, domain.EWBStatusExtended, false},
		{domain.EWBStatusCancelled, domain.EWBStatusConsolidated, false},

		{domain.EWBStatusConsolidated, domain.EWBStatusActive, false},
		{domain.EWBStatusConsolidated, domain.EWBStatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s->%s", tt.from, tt.to), func(t *testing.T) {
			assert.Equal(t, tt.valid, domain.CanTransitionTo(tt.from, tt.to))
		})
	}
}

// ---------------------------------------------------------------------------
// Tests: ValidityDays (regular + ODC)
// ---------------------------------------------------------------------------

func TestValidityDays_Regular(t *testing.T) {
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
			assert.Equal(t, tt.days, store.ValidityDays(tt.km, false))
		})
	}
}

func TestValidityDays_ODC(t *testing.T) {
	tests := []struct {
		km   int
		days int
	}{
		{0, 1},
		{20, 1},
		{21, 2},
		{100, 5},
		{200, 10},
		{201, 11},
		{400, 20},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("ODC_%dkm", tt.km), func(t *testing.T) {
			assert.Equal(t, tt.days, store.ValidityDays(tt.km, true))
		})
	}
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
// Tests: Router
// ---------------------------------------------------------------------------

func TestNewRouter(t *testing.T) {
	mockStore := store.NewMockStore()
	ewbServer := newMockEWBServer(t)
	defer ewbServer.Close()
	ewbClient := gateway.NewEWBClient(ewbServer.URL)

	r := NewRouter(mockStore, ewbClient, store.RealClock{})
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
