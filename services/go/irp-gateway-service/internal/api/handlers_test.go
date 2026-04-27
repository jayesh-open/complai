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
	"github.com/complai/complai/services/go/irp-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/irp-gateway-service/internal/provider"
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

func getWithTenant(t *testing.T, path, tenantID string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("X-Tenant-Id", tenantID)
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
	assert.Equal(t, "irp-gateway-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: Authenticate
// ---------------------------------------------------------------------------

func TestAuthenticate(t *testing.T) {
	h := newTestHandlers()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/irp/authenticate", nil)
	h.Authenticate(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var data domain.AuthResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &data)
	assert.Contains(t, data.AccessToken, "mock-irp-token-")
	assert.Equal(t, "bearer", data.TokenType)
	assert.Equal(t, 86399, data.ExpiresIn)
}

// ---------------------------------------------------------------------------
// Tests: Full IRN Lifecycle (generate → get by IRN → get by doc → cancel)
// ---------------------------------------------------------------------------

func TestIRN_FullLifecycle(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	genReq := domain.GenerateIRNRequest{
		GSTIN: "29AABCA1234A1Z5",
		DocDtls: domain.DocDetails{
			Typ: "INV",
			No:  "INV-001",
			Dt:  "15/04/2026",
		},
		SupDtls: domain.PartyDetail{
			GSTIN: "29AABCA1234A1Z5",
			LglNm: "Supplier Co",
			Stcd:  "29",
		},
		BuyDtls: domain.PartyDetail{
			GSTIN: "27AABCB5678B1Z3",
			LglNm: "Buyer Co",
			Pos:   "27",
		},
		ItemList: []domain.LineItem{
			{
				SlNo: "1", PrdDesc: "Steel Plates", HsnCd: "720241",
				Qty: 100, Unit: "KG", UnitPrice: 250, TaxableAmt: 25000,
				IgstRt: 18, IgstAmt: 4500,
			},
		},
		ValDtls: domain.ValDetails{
			TaxableVal: 25000, IGST: 4500, TotInvVal: 29500,
		},
		RequestID: uuid.New().String(),
	}

	// Step 1: Generate IRN
	rec := httptest.NewRecorder()
	h.GenerateIRN(rec, postJSON(t, "/v1/gateway/irp/invoice", genReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var genResp domain.GenerateIRNResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &genResp)
	assert.NotEmpty(t, genResp.IRN)
	assert.Equal(t, 64, len(genResp.IRN))
	assert.Equal(t, "ACT", genResp.Status)
	assert.NotEmpty(t, genResp.AckNo)
	assert.NotEmpty(t, genResp.SignedInvoice)
	assert.NotEmpty(t, genResp.SignedQRCode)

	// Step 2: Get by IRN
	rec = httptest.NewRecorder()
	h.GetIRNByIRN(rec, getWithTenant(t, "/v1/gateway/irp/invoice/irn?irn="+genResp.IRN, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var getResp domain.GetIRNResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &getResp)
	assert.Equal(t, genResp.IRN, getResp.IRN)
	assert.Equal(t, "ACT", getResp.Status)
	assert.Equal(t, "INV", getResp.DocType)
	assert.Equal(t, "INV-001", getResp.DocNo)

	// Step 3: Get by doc details
	rec = httptest.NewRecorder()
	h.GetIRNByDoc(rec, getWithTenant(t, "/v1/gateway/irp/invoice/irnbydocdetails?doctype=INV&docnum=INV-001&docdate=15/04/2026", tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var docResp domain.GetIRNResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &docResp)
	assert.Equal(t, genResp.IRN, docResp.IRN)

	// Step 4: Cancel IRN
	cancelReq := domain.CancelIRNRequest{
		IRN:       genResp.IRN,
		CnlRsn:   "1",
		CnlRem:   "Duplicate invoice",
		RequestID: uuid.New().String(),
	}
	rec = httptest.NewRecorder()
	h.CancelIRN(rec, postJSON(t, "/v1/gateway/irp/invoice/cancel", cancelReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var cancelResp domain.CancelIRNResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &cancelResp)
	assert.Equal(t, genResp.IRN, cancelResp.IRN)
	assert.Equal(t, "CANC", cancelResp.Status)

	// Step 5: Verify status is CANC after get
	rec = httptest.NewRecorder()
	h.GetIRNByIRN(rec, getWithTenant(t, "/v1/gateway/irp/invoice/irn?irn="+genResp.IRN, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var getAfterCancel domain.GetIRNResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &getAfterCancel)
	assert.Equal(t, "CANC", getAfterCancel.Status)
	assert.NotNil(t, getAfterCancel.CancelledAt)
}

// ---------------------------------------------------------------------------
// Tests: Idempotent IRN generation (same doc → same IRN)
// ---------------------------------------------------------------------------

func TestGenerateIRN_Idempotent(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	req := domain.GenerateIRNRequest{
		GSTIN:   "29AABCA1234A1Z5",
		DocDtls: domain.DocDetails{Typ: "INV", No: "INV-IDEM", Dt: "15/04/2026"},
		SupDtls: domain.PartyDetail{GSTIN: "29AABCA1234A1Z5", LglNm: "X"},
		BuyDtls: domain.PartyDetail{GSTIN: "27AABCB5678B1Z3", LglNm: "Y"},
		ValDtls: domain.ValDetails{TotInvVal: 10000},
	}

	rec1 := httptest.NewRecorder()
	h.GenerateIRN(rec1, postJSON(t, "/", req, tenantID))
	require.Equal(t, http.StatusOK, rec1.Code)

	var resp1 domain.GenerateIRNResponse
	parseGatewayResponse(t, rec1.Body.Bytes(), &resp1)

	rec2 := httptest.NewRecorder()
	h.GenerateIRN(rec2, postJSON(t, "/", req, tenantID))
	require.Equal(t, http.StatusOK, rec2.Code)

	var resp2 domain.GenerateIRNResponse
	parseGatewayResponse(t, rec2.Body.Bytes(), &resp2)

	assert.Equal(t, resp1.IRN, resp2.IRN)
}

// ---------------------------------------------------------------------------
// Tests: Cancel already cancelled IRN
// ---------------------------------------------------------------------------

func TestCancelIRN_AlreadyCancelled(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	genReq := domain.GenerateIRNRequest{
		GSTIN:   "29AABCA1234A1Z5",
		DocDtls: domain.DocDetails{Typ: "INV", No: "INV-CANC-TWICE", Dt: "15/04/2026"},
		SupDtls: domain.PartyDetail{GSTIN: "29AABCA1234A1Z5", LglNm: "X"},
		BuyDtls: domain.PartyDetail{GSTIN: "27AABCB5678B1Z3", LglNm: "Y"},
		ValDtls: domain.ValDetails{TotInvVal: 5000},
	}

	rec := httptest.NewRecorder()
	h.GenerateIRN(rec, postJSON(t, "/", genReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var genResp domain.GenerateIRNResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &genResp)

	cancelReq := domain.CancelIRNRequest{IRN: genResp.IRN, CnlRsn: "1", CnlRem: "Test"}

	rec = httptest.NewRecorder()
	h.CancelIRN(rec, postJSON(t, "/", cancelReq, tenantID))
	assert.Equal(t, http.StatusOK, rec.Code)

	rec = httptest.NewRecorder()
	h.CancelIRN(rec, postJSON(t, "/", cancelReq, tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Cancel non-existent IRN
// ---------------------------------------------------------------------------

func TestCancelIRN_NotFound(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	cancelReq := domain.CancelIRNRequest{IRN: "nonexistent-irn", CnlRsn: "1", CnlRem: "Test"}
	rec := httptest.NewRecorder()
	h.CancelIRN(rec, postJSON(t, "/", cancelReq, tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: GSTIN Validation
// ---------------------------------------------------------------------------

func TestValidateGSTIN_Success(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	rec := httptest.NewRecorder()
	h.ValidateGSTIN(rec, getWithTenant(t, "/v1/gateway/irp/master/gstin?gstin=29AABCA1234A1Z5", tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	var resp domain.GSTINValidateResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, "29AABCA1234A1Z5", resp.GSTIN)
	assert.Equal(t, "Active", resp.Status)
	assert.Equal(t, "29", resp.StateCode)
}

func TestValidateGSTIN_InvalidLength(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	rec := httptest.NewRecorder()
	h.ValidateGSTIN(rec, getWithTenant(t, "/v1/gateway/irp/master/gstin?gstin=INVALID", tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

func TestValidateGSTIN_MissingParam(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	rec := httptest.NewRecorder()
	h.ValidateGSTIN(rec, getWithTenant(t, "/v1/gateway/irp/master/gstin", tenantID))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Error cases — missing headers / bad body
// ---------------------------------------------------------------------------

func TestGenerateIRN_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.GenerateIRNRequest{GSTIN: "X"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.GenerateIRN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGenerateIRN_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GenerateIRN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGenerateIRN_MissingGSTIN(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()
	req := domain.GenerateIRNRequest{
		DocDtls: domain.DocDetails{Typ: "INV", No: "INV-001", Dt: "15/04/2026"},
	}
	rec := httptest.NewRecorder()
	h.GenerateIRN(rec, postJSON(t, "/", req, tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

func TestCancelIRN_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	body, _ := json.Marshal(domain.CancelIRNRequest{IRN: "x"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.CancelIRN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCancelIRN_InvalidBody(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.CancelIRN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetIRNByIRN_MissingTenantID(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodGet, "/v1/gateway/irp/invoice/irn?irn=x", nil)
	rec := httptest.NewRecorder()
	h.GetIRNByIRN(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetIRNByIRN_MissingIRN(t *testing.T) {
	h := newTestHandlers()
	rec := httptest.NewRecorder()
	h.GetIRNByIRN(rec, getWithTenant(t, "/v1/gateway/irp/invoice/irn", uuid.New().String()))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetIRNByIRN_NotFound(t *testing.T) {
	h := newTestHandlers()
	rec := httptest.NewRecorder()
	h.GetIRNByIRN(rec, getWithTenant(t, "/v1/gateway/irp/invoice/irn?irn=nonexistent", uuid.New().String()))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

func TestGetIRNByDoc_MissingParams(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	rec := httptest.NewRecorder()
	h.GetIRNByDoc(rec, getWithTenant(t, "/v1/gateway/irp/invoice/irnbydocdetails?doctype=INV", tenantID))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetIRNByDoc_NotFound(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	rec := httptest.NewRecorder()
	h.GetIRNByDoc(rec, getWithTenant(t, "/v1/gateway/irp/invoice/irnbydocdetails?doctype=INV&docnum=NONE&docdate=01/01/2026", tenantID))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
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

// ---------------------------------------------------------------------------
// Tests: Re-generate IRN after cancellation
// ---------------------------------------------------------------------------

func TestGenerateIRN_AfterCancellation(t *testing.T) {
	h := newTestHandlers()
	tenantID := uuid.New().String()

	genReq := domain.GenerateIRNRequest{
		GSTIN:   "29AABCA1234A1Z5",
		DocDtls: domain.DocDetails{Typ: "INV", No: "INV-REGEN", Dt: "15/04/2026"},
		SupDtls: domain.PartyDetail{GSTIN: "29AABCA1234A1Z5", LglNm: "X"},
		BuyDtls: domain.PartyDetail{GSTIN: "27AABCB5678B1Z3", LglNm: "Y"},
		ValDtls: domain.ValDetails{TotInvVal: 10000},
	}

	// Generate
	rec := httptest.NewRecorder()
	h.GenerateIRN(rec, postJSON(t, "/", genReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)
	var resp1 domain.GenerateIRNResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resp1)

	// Cancel
	cancelReq := domain.CancelIRNRequest{IRN: resp1.IRN, CnlRsn: "2", CnlRem: "Mistake"}
	rec = httptest.NewRecorder()
	h.CancelIRN(rec, postJSON(t, "/", cancelReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)

	// Re-generate (same doc details → new IRN since old is cancelled)
	rec = httptest.NewRecorder()
	h.GenerateIRN(rec, postJSON(t, "/", genReq, tenantID))
	require.Equal(t, http.StatusOK, rec.Code)
	var resp2 domain.GenerateIRNResponse
	parseGatewayResponse(t, rec.Body.Bytes(), &resp2)
	assert.NotEqual(t, resp1.IRN, resp2.IRN)
	assert.Equal(t, "ACT", resp2.Status)
}
