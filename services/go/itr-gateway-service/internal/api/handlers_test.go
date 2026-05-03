package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/complai/complai/services/go/itr-gateway-service/internal/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup() *http.ServeMux {
	p := provider.NewMockProvider()
	h := NewHandlers(p)
	r := NewRouter(h)
	mux := http.NewServeMux()
	mux.Handle("/", r)
	return mux
}

func post(handler http.Handler, path string, body interface{}) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", "test-tenant")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func TestHealth(t *testing.T) {
	mux := setup()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCheckPANAadhaarLink_Success(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/pan-aadhaar/check", map[string]string{"pan": "ABCDE1234F"})
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, true, data["linked"])
}

func TestCheckPANAadhaarLink_MissingPAN(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/pan-aadhaar/check", map[string]string{})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCheckPANAadhaarLink_InvalidBody(t *testing.T) {
	mux := setup()
	req := httptest.NewRequest(http.MethodPost, "/v1/gateway/itr/pan-aadhaar/check", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFetchAIS_Success(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/ais/fetch", map[string]string{"pan": "ABCDE1234F", "tax_year": "2026-27"})
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "ABCDE1234F", data["pan"])
	entries := data["tds_entries"].([]interface{})
	assert.Len(t, entries, 2)
}

func TestFetchAIS_MissingFields(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/ais/fetch", map[string]string{"pan": "ABCDE1234F"})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFetchAIS_InvalidPAN(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/ais/fetch", map[string]string{"pan": "BAD", "tax_year": "2026-27"})
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSubmitITR_Success(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/submit", map[string]string{
		"pan": "ABCDE1234F", "tax_year": "2026-27", "form_type": "ITR-1", "payload": "{}",
	})
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "SUBMITTED", data["status"])
}

func TestSubmitITR_MissingFields(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/submit", map[string]string{"pan": "ABCDE1234F"})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSubmitITR_UnsupportedForm(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/submit", map[string]string{
		"pan": "ABCDE1234F", "tax_year": "2026-27", "form_type": "ITR-99",
	})
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGenerateITRV_MissingARN(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/itrv/generate", map[string]string{})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGenerateITRV_NotFound(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/itrv/generate", map[string]string{"arn": "BOGUS"})
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCheckEVerification_MissingARN(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/everify/check", map[string]string{})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCheckEVerification_NotFound(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/everify/check", map[string]string{"arn": "BOGUS"})
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCheckRefundStatus_Success(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/refund/status", map[string]string{
		"pan": "ABCDE1234F", "tax_year": "2026-27",
	})
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "PROCESSED", data["status"])
}

func TestCheckRefundStatus_MissingFields(t *testing.T) {
	mux := setup()
	w := post(mux, "/v1/gateway/itr/refund/status", map[string]string{"pan": "ABCDE1234F"})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSubmitAndGenerateITRV_E2E(t *testing.T) {
	mux := setup()

	w := post(mux, "/v1/gateway/itr/submit", map[string]string{
		"pan": "ABCDE1234F", "tax_year": "2026-27", "form_type": "ITR-1", "payload": "{}",
	})
	require.Equal(t, http.StatusOK, w.Code)
	var submitResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &submitResp))
	arn := submitResp["data"].(map[string]interface{})["arn"].(string)

	w2 := post(mux, "/v1/gateway/itr/itrv/generate", map[string]string{"arn": arn})
	assert.Equal(t, http.StatusOK, w2.Code)
	var itrvResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &itrvResp))
	assert.Contains(t, itrvResp["data"].(map[string]interface{})["itrv_url"], arn)

	w3 := post(mux, "/v1/gateway/itr/everify/check", map[string]string{"arn": arn})
	assert.Equal(t, http.StatusOK, w3.Code)
	var evResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w3.Body.Bytes(), &evResp))
	assert.Equal(t, true, evResp["data"].(map[string]interface{})["verified"])
}
