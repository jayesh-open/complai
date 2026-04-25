package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/aura-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/aura-gateway-service/internal/provider"
)

func parseDataResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	var wrapper httputil.SuccessResponse
	wrapper.Data = target
	require.NoError(t, json.Unmarshal(body, &wrapper))
}

func TestHealth(t *testing.T) {
	h := NewHandlers(provider.NewMockProvider())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	h.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "aura-gateway-service", data["service"])
}

func TestListARInvoices_Success(t *testing.T) {
	h := NewHandlers(provider.NewMockProvider())
	tenantID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/v1/gateway/aura/invoices?gstin=29AABCA1234A1Z5&period=042026", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()
	h.ListARInvoices(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.InvoiceListResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, 100, resp.TotalCount)
	assert.Len(t, resp.Invoices, 100)

	// Verify invoice type mix (NIL-rated are B2B intra with HSN 0101, so included in B2BIntraCount)
	s := resp.Summary
	assert.Equal(t, 35, s.B2BIntraCount, "B2B intra-state (30 regular + 5 NIL-rated)")
	assert.Equal(t, 20, s.B2BInterCount, "B2B inter-state")
	assert.Equal(t, 15, s.B2CSCount, "B2CS")
	assert.Equal(t, 5, s.B2CLCount, "B2CL")
	assert.Equal(t, 5, s.ExportCount, "Export")
	assert.Equal(t, 5, s.RCMCount, "RCM")
	assert.Equal(t, 10, s.CreditNote, "Credit notes")
	assert.Equal(t, 5, s.DebitNote, "Debit notes")

	// Verify NIL-rated invoices exist (HSN 0101, zero tax)
	var nilCount int
	for _, inv := range resp.Invoices {
		if len(inv.LineItems) > 0 && inv.LineItems[0].HSN == "0101" {
			nilCount++
		}
	}
	assert.Equal(t, 5, nilCount, "NIL-rated invoices with HSN 0101")
}

func TestListARInvoices_DefaultPeriod(t *testing.T) {
	h := NewHandlers(provider.NewMockProvider())
	req := httptest.NewRequest(http.MethodGet, "/v1/gateway/aura/invoices?gstin=29AABCA1234A1Z5", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.ListARInvoices(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.InvoiceListResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, 100, resp.TotalCount)
}

func TestListARInvoices_MissingTenantID(t *testing.T) {
	h := NewHandlers(provider.NewMockProvider())
	req := httptest.NewRequest(http.MethodGet, "/v1/gateway/aura/invoices?gstin=X", nil)
	rec := httptest.NewRecorder()
	h.ListARInvoices(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListARInvoices_InvalidTenantID(t *testing.T) {
	h := NewHandlers(provider.NewMockProvider())
	req := httptest.NewRequest(http.MethodGet, "/v1/gateway/aura/invoices?gstin=X", nil)
	req.Header.Set("X-Tenant-Id", "bad")
	rec := httptest.NewRecorder()
	h.ListARInvoices(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListARInvoices_MissingGSTIN(t *testing.T) {
	h := NewHandlers(provider.NewMockProvider())
	req := httptest.NewRequest(http.MethodGet, "/v1/gateway/aura/invoices", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.ListARInvoices(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListARInvoices_InvoiceFields(t *testing.T) {
	h := NewHandlers(provider.NewMockProvider())
	req := httptest.NewRequest(http.MethodGet, "/v1/gateway/aura/invoices?gstin=29AABCA1234A1Z5", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.ListARInvoices(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp domain.InvoiceListResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)

	inv := resp.Invoices[0]
	assert.NotEqual(t, uuid.Nil, inv.ID)
	assert.Equal(t, "INV", inv.DocumentType)
	assert.Equal(t, "B2B", inv.SupplyType)
	assert.Equal(t, "29AABCA1234A1Z5", inv.Supplier.GSTIN)
	assert.Equal(t, "29", inv.Supplier.StateCode)
	assert.Equal(t, "aura", inv.SourceSystem)
	assert.NotEmpty(t, inv.LineItems)
	assert.True(t, inv.Totals.GrandTotal.IsPositive())
}

func TestListARInvoices_CreditNoteFields(t *testing.T) {
	h := NewHandlers(provider.NewMockProvider())
	req := httptest.NewRequest(http.MethodGet, "/v1/gateway/aura/invoices?gstin=29AABCA1234A1Z5", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.ListARInvoices(rec, req)

	var resp domain.InvoiceListResponse
	parseDataResponse(t, rec.Body.Bytes(), &resp)

	var crnCount int
	for _, inv := range resp.Invoices {
		if inv.DocumentType == "CRN" {
			crnCount++
			assert.Contains(t, inv.DocumentNumber, "CRN/")
		}
	}
	assert.Equal(t, 10, crnCount)
}

func TestNewRouter(t *testing.T) {
	r := NewRouter(provider.NewMockProvider())
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

func TestTenantIDFromRequest_Valid(t *testing.T) {
	expected := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-Id", expected.String())
	got, err := tenantIDFromRequest(req)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestTenantIDFromRequest_Missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := tenantIDFromRequest(req)
	require.Error(t, err)
}

func TestTenantIDFromRequest_Invalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-Id", "not-uuid")
	_, err := tenantIDFromRequest(req)
	require.Error(t, err)
}
