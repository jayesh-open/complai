package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/complai/complai/services/go/tds-service/internal/domain"
	"github.com/complai/complai/services/go/tds-service/internal/store"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testTenant = "11111111-1111-1111-1111-111111111111"

func setupTest() (*Handlers, *store.MockStore) {
	ms := store.NewMockStore()
	return NewHandlers(ms), ms
}

func seedDeductee(ms *store.MockStore) *domain.Deductee {
	tenantID := uuid.MustParse(testTenant)
	d := &domain.Deductee{
		ID:             uuid.New(),
		TenantID:       tenantID,
		VendorID:       uuid.New(),
		Name:           "Test Vendor Pvt Ltd",
		PAN:            "ABCCD1234E",
		PANVerified:    true,
		PANStatus:      "VALID",
		DeducteeType:   domain.DeducteeCompany,
		ResidentStatus: domain.Resident,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	ms.CreateDeductee(nil, tenantID, d)
	return d
}

func TestListDeductees_Empty(t *testing.T) {
	h, _ := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/deductees", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.ListDeductees(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListDeductees_WithData(t *testing.T) {
	h, ms := setupTest()
	seedDeductee(ms)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/deductees", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.ListDeductees(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
}

func TestListDeductees_MissingTenant(t *testing.T) {
	h, _ := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/deductees", nil)
	w := httptest.NewRecorder()

	h.ListDeductees(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDeductee_Found(t *testing.T) {
	h, ms := setupTest()
	d := seedDeductee(ms)
	router := NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/deductees/"+d.ID.String(), nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetDeductee_NotFound(t *testing.T) {
	h, _ := setupTest()
	router := NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/deductees/"+uuid.New().String(), nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetDeductee_InvalidID(t *testing.T) {
	h, _ := setupTest()
	router := NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/deductees/not-a-uuid", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCalculateTDS_ContractorOther(t *testing.T) {
	h, _ := setupTest()
	body := `{"payment_code":"1024","gross_amount":100000,"has_valid_pan":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/calculate", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.CalculateTDS(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "0.02", data["rate"])
	assert.Equal(t, "2000", data["tds_amount"])
}

func TestCalculateTDS_Salary(t *testing.T) {
	h, _ := setupTest()
	body := `{"payment_code":"1002","annual_salary":2000000,"has_valid_pan":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/calculate", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.CalculateTDS(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "16033", data["tds_amount"])
}

func TestCalculateTDS_InvalidPaymentCode(t *testing.T) {
	h, _ := setupTest()
	body := `{"payment_code":"9999","gross_amount":100000}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/calculate", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CalculateTDS(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCalculateTDS_InvalidBody(t *testing.T) {
	h, _ := setupTest()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/calculate", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()

	h.CalculateTDS(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCalculateTDS_WithDTAA(t *testing.T) {
	h, _ := setupTest()
	body := `{"payment_code":"1057","gross_amount":1000000,"has_valid_pan":true,"dtaa_rate":0.10}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/calculate", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.CalculateTDS(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateEntry_Success(t *testing.T) {
	h, ms := setupTest()
	d := seedDeductee(ms)

	body, _ := json.Marshal(map[string]interface{}{
		"deductee_id":       d.ID.String(),
		"payment_code":      "1024",
		"financial_year":    "2026-27",
		"quarter":           "Q1",
		"transaction_date":  "2026-06-15",
		"gross_amount":      100000,
		"nature_of_payment": "Contractor payment",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/entries", bytes.NewReader(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.CreateEntry(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "PENDING", data["status"])
	assert.Equal(t, "0.02", data["tds_rate"])
}

func TestCreateEntry_InvalidDeductee(t *testing.T) {
	h, _ := setupTest()
	body, _ := json.Marshal(map[string]interface{}{
		"deductee_id":       uuid.New().String(),
		"payment_code":      "1024",
		"financial_year":    "2026-27",
		"quarter":           "Q1",
		"transaction_date":  "2026-06-15",
		"gross_amount":      100000,
		"nature_of_payment": "test",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/entries", bytes.NewReader(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.CreateEntry(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateEntry_InvalidPaymentCode(t *testing.T) {
	h, ms := setupTest()
	d := seedDeductee(ms)
	body, _ := json.Marshal(map[string]interface{}{
		"deductee_id":       d.ID.String(),
		"payment_code":      "9999",
		"financial_year":    "2026-27",
		"quarter":           "Q1",
		"transaction_date":  "2026-06-15",
		"gross_amount":      100000,
		"nature_of_payment": "test",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/entries", bytes.NewReader(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.CreateEntry(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateEntry_InvalidDate(t *testing.T) {
	h, ms := setupTest()
	d := seedDeductee(ms)
	body, _ := json.Marshal(map[string]interface{}{
		"deductee_id":       d.ID.String(),
		"payment_code":      "1024",
		"financial_year":    "2026-27",
		"quarter":           "Q1",
		"transaction_date":  "bad-date",
		"gross_amount":      100000,
		"nature_of_payment": "test",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/entries", bytes.NewReader(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.CreateEntry(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListEntries_Empty(t *testing.T) {
	h, _ := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/entries?fy=2026-27", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.ListEntries(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetEntry_NotFound(t *testing.T) {
	h, _ := setupTest()
	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/entries/"+uuid.New().String(), nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetSummary(t *testing.T) {
	h, ms := setupTest()
	tenantID := uuid.MustParse(testTenant)
	d := seedDeductee(ms)

	entry := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenantID, DeducteeID: d.ID,
		Section: domain.Section393_1, PaymentCode: domain.CodeContractorOther,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(100000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(2000),
		TotalTax: decimal.NewFromInt(2000), Status: domain.StatusPending,
		NatureOfPayment: "contractor",
	}
	ms.CreateEntry(nil, tenantID, entry)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/summary?fy=2026-27", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.GetSummary(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHealth_Endpoint(t *testing.T) {
	h, _ := setupTest()
	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPagination_Defaults(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	limit, offset := pagination(req)
	assert.Equal(t, 50, limit)
	assert.Equal(t, 0, offset)
}

func TestPagination_Custom(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test?limit=10&offset=20", nil)
	limit, offset := pagination(req)
	assert.Equal(t, 10, limit)
	assert.Equal(t, 20, offset)
}

func TestPagination_MaxLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test?limit=500", nil)
	limit, _ := pagination(req)
	assert.Equal(t, 50, limit, "limit > 200 should fall back to default")
}

func TestCreateEntry_InvalidBody(t *testing.T) {
	h, _ := setupTest()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/entries", bytes.NewBufferString("not json"))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.CreateEntry(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateEntry_MissingTenant(t *testing.T) {
	h, _ := setupTest()
	body, _ := json.Marshal(map[string]interface{}{
		"deductee_id":       uuid.New().String(),
		"payment_code":      "1024",
		"financial_year":    "2026-27",
		"quarter":           "Q1",
		"transaction_date":  "2026-06-15",
		"gross_amount":      100000,
		"nature_of_payment": "test",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/entries", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateEntry(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateEntry_InvalidDeducteeID(t *testing.T) {
	h, _ := setupTest()
	body, _ := json.Marshal(map[string]interface{}{
		"deductee_id":       "not-a-uuid",
		"payment_code":      "1024",
		"financial_year":    "2026-27",
		"quarter":           "Q1",
		"transaction_date":  "2026-06-15",
		"gross_amount":      100000,
		"nature_of_payment": "test",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/entries", bytes.NewReader(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.CreateEntry(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListEntries_MissingTenant(t *testing.T) {
	h, _ := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/entries?fy=2026-27", nil)
	w := httptest.NewRecorder()

	h.ListEntries(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetEntry_MissingTenant(t *testing.T) {
	h, _ := setupTest()
	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/entries/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetEntry_InvalidID(t *testing.T) {
	h, _ := setupTest()
	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/entries/not-a-uuid", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetSummary_MissingTenant(t *testing.T) {
	h, _ := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/summary?fy=2026-27", nil)
	w := httptest.NewRecorder()

	h.GetSummary(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetSummary_DefaultFY(t *testing.T) {
	h, _ := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/summary", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.GetSummary(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetDeductee_MissingTenant(t *testing.T) {
	h, _ := setupTest()
	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/deductees/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCalculateTDS_WithLowerCert(t *testing.T) {
	h, _ := setupTest()
	body := `{"payment_code":"1024","gross_amount":100000,"has_valid_pan":true,"lower_cert_rate":0.005}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tds/calculate", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.CalculateTDS(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "0.005", data["rate"])
}

func TestListEntries_WithQuarterFilter(t *testing.T) {
	h, ms := setupTest()
	tenantID := uuid.MustParse(testTenant)
	d := seedDeductee(ms)
	entry := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenantID, DeducteeID: d.ID,
		Section: domain.Section393_1, PaymentCode: domain.CodeContractorOther,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q2",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(100000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(2000),
		TotalTax: decimal.NewFromInt(2000), Status: domain.StatusPending,
		NatureOfPayment: "test",
	}
	ms.CreateEntry(nil, tenantID, entry)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tds/entries?fy=2026-27&quarter=Q2", nil)
	req.Header.Set("X-Tenant-Id", testTenant)
	w := httptest.NewRecorder()

	h.ListEntries(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
