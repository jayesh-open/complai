package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/complai/complai/services/go/tds-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/tds-gateway-service/internal/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup() (*Handlers, *provider.MockProvider) {
	p := provider.NewMockProvider()
	return NewHandlers(p), p
}

func TestVerifyPAN_Handler_Success(t *testing.T) {
	h, _ := setup()
	body := `{"pan":"ABCPD1234E","name":"Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/pan/verify", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()

	h.VerifyPAN(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "VALID", data["status"])
	assert.Equal(t, "INDIVIDUAL", data["category"])
}

func TestVerifyPAN_Handler_MissingPAN(t *testing.T) {
	h, _ := setup()
	body := `{"name":"Test"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/pan/verify", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.VerifyPAN(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerifyPAN_Handler_InvalidBody(t *testing.T) {
	h, _ := setup()
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/pan/verify", bytes.NewBufferString("not json"))
	w := httptest.NewRecorder()

	h.VerifyPAN(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerifyTAN_Handler_Success(t *testing.T) {
	h, _ := setup()
	body := `{"tan":"MUMA12345A"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/tan/verify", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()

	h.VerifyTAN(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVerifyTAN_Handler_MissingTAN(t *testing.T) {
	h, _ := setup()
	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/tan/verify", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.VerifyTAN(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGenerateChallan_Handler_Success(t *testing.T) {
	h, _ := setup()
	body := `{"tan":"MUMA12345A","section":"194C","amount":50000,"assessment_year":"2026-27"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/challan/generate", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()

	h.GenerateChallan(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "SUCCESS", data["status"])
}

func TestGenerateChallan_Handler_InvalidAmount(t *testing.T) {
	h, _ := setup()
	body := `{"tan":"MUMA12345A","section":"194C","amount":0}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/challan/generate", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.GenerateChallan(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileForm26Q_Handler_Success(t *testing.T) {
	h, _ := setup()
	body := `{"tan":"MUMA12345A","financial_year":"2025-26","quarter":"Q1","deductions":[{"deductee_pan":"ABCPD1234E","section":"194C","amount":50000,"tds_amount":1000}]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/form26q/file", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()

	h.FileForm26Q(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFileForm26Q_Handler_MissingDeductions(t *testing.T) {
	h, _ := setup()
	body := `{"tan":"MUMA12345A","financial_year":"2025-26","quarter":"Q1","deductions":[]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/form26q/file", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.FileForm26Q(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileForm24Q_Handler_Success(t *testing.T) {
	h, _ := setup()
	body := `{"tan":"MUMA12345A","financial_year":"2025-26","quarter":"Q1","employees":[{"pan":"ABCPD1234E","name":"Emp","gross_salary":1200000,"tds_deducted":50000}]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/form24q/file", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()

	h.FileForm24Q(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFileForm24Q_Handler_MissingEmployees(t *testing.T) {
	h, _ := setup()
	body := `{"tan":"MUMA12345A","financial_year":"2025-26","quarter":"Q1","employees":[]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/form24q/file", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.FileForm24Q(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRouter_HealthEndpoint(t *testing.T) {
	h, _ := setup()
	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

type failingProvider struct{}

func (f *failingProvider) VerifyPAN(_ context.Context, _ domain.PANVerifyRequest) (*domain.PANVerifyResponse, error) {
	return nil, fmt.Errorf("provider down")
}
func (f *failingProvider) VerifyTAN(_ context.Context, _ domain.TANVerifyRequest) (*domain.TANVerifyResponse, error) {
	return nil, fmt.Errorf("provider down")
}
func (f *failingProvider) GenerateChallan(_ context.Context, _ domain.ChallanRequest) (*domain.ChallanResponse, error) {
	return nil, fmt.Errorf("provider down")
}
func (f *failingProvider) FileForm26Q(_ context.Context, _ domain.Form26QRequest) (*domain.FormFilingResponse, error) {
	return nil, fmt.Errorf("provider down")
}
func (f *failingProvider) FileForm24Q(_ context.Context, _ domain.Form24QRequest) (*domain.FormFilingResponse, error) {
	return nil, fmt.Errorf("provider down")
}

func setupFailing() *Handlers {
	return NewHandlers(&failingProvider{})
}

func TestVerifyPAN_Handler_ProviderError(t *testing.T) {
	h := setupFailing()
	body := `{"pan":"ABCPD1234E","name":"Test"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/pan/verify", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()
	h.VerifyPAN(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVerifyTAN_Handler_ProviderError(t *testing.T) {
	h := setupFailing()
	body := `{"tan":"MUMA12345A"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/tan/verify", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()
	h.VerifyTAN(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGenerateChallan_Handler_ProviderError(t *testing.T) {
	h := setupFailing()
	body := `{"tan":"MUMA12345A","section":"194C","amount":50000,"assessment_year":"2026-27"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/challan/generate", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()
	h.GenerateChallan(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestFileForm26Q_Handler_ProviderError(t *testing.T) {
	h := setupFailing()
	body := `{"tan":"MUMA12345A","financial_year":"2025-26","quarter":"Q1","deductions":[{"deductee_pan":"ABCPD1234E","section":"194C","amount":50000,"tds_amount":1000}]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/form26q/file", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()
	h.FileForm26Q(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestFileForm24Q_Handler_ProviderError(t *testing.T) {
	h := setupFailing()
	body := `{"tan":"MUMA12345A","financial_year":"2025-26","quarter":"Q1","employees":[{"pan":"ABCPD1234E","name":"Emp","gross_salary":1200000,"tds_deducted":50000}]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/form24q/file", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()
	h.FileForm24Q(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVerifyTAN_Handler_InvalidBody(t *testing.T) {
	h, _ := setup()
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/tan/verify", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()
	h.VerifyTAN(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGenerateChallan_Handler_InvalidBody(t *testing.T) {
	h, _ := setup()
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/challan/generate", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()
	h.GenerateChallan(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileForm26Q_Handler_InvalidBody(t *testing.T) {
	h, _ := setup()
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/form26q/file", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()
	h.FileForm26Q(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileForm24Q_Handler_InvalidBody(t *testing.T) {
	h, _ := setup()
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/form24q/file", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()
	h.FileForm24Q(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRouter_FullRoundtrip(t *testing.T) {
	h, _ := setup()
	router := NewRouter(h)

	body := `{"pan":"ABCPD1234E","name":"Full Test"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/tds/pan/verify", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "VALID", data["status"])
}
