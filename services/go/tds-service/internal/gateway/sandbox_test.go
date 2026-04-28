package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSandboxClient_VerifyPAN_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/gateway/tds/pan/verify", r.URL.Path)
		assert.Equal(t, "tenant-1", r.Header.Get("X-Tenant-Id"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]string{
				"pan": "ABCPD1234E", "name": "Test", "status": "VALID", "category": "INDIVIDUAL",
			},
		})
	}))
	defer srv.Close()

	c := NewSandboxClient(srv.URL)
	result, err := c.VerifyPAN(context.Background(), "tenant-1", "ABCPD1234E", "Test")
	require.NoError(t, err)
	assert.Equal(t, "VALID", result.Status)
	assert.Equal(t, "INDIVIDUAL", result.Category)
}

func TestSandboxClient_VerifyPAN_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	c := NewSandboxClient(srv.URL)
	_, err := c.VerifyPAN(context.Background(), "tenant-1", "ABCPD1234E", "Test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "502")
}

func TestSandboxClient_VerifyPAN_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	c := NewSandboxClient(srv.URL)
	_, err := c.VerifyPAN(context.Background(), "tenant-1", "ABCPD1234E", "Test")
	assert.Error(t, err)
}

func TestSandboxClient_VerifyPAN_ConnectionError(t *testing.T) {
	c := NewSandboxClient("http://127.0.0.1:1")
	_, err := c.VerifyPAN(context.Background(), "tenant-1", "ABCPD1234E", "Test")
	assert.Error(t, err)
}

func TestSandboxClient_GenerateChallan_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/gateway/tds/challan/generate", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"challan_number": "CHN-001", "bsr_code": "BSR001",
				"deposit_date": "2025-06-15", "amount": 50000.0, "status": "SUCCESS",
			},
		})
	}))
	defer srv.Close()

	c := NewSandboxClient(srv.URL)
	payload := map[string]interface{}{"tan": "MUMA12345A", "amount": 50000.0}
	result, err := c.GenerateChallan(context.Background(), "tenant-1", payload)
	require.NoError(t, err)
	assert.Equal(t, "SUCCESS", result.Status)
	assert.Equal(t, "CHN-001", result.ChallanNumber)
}

func TestSandboxClient_GenerateChallan_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewSandboxClient(srv.URL)
	_, err := c.GenerateChallan(context.Background(), "tenant-1", map[string]interface{}{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestSandboxClient_GenerateChallan_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{invalid"))
	}))
	defer srv.Close()

	c := NewSandboxClient(srv.URL)
	_, err := c.GenerateChallan(context.Background(), "tenant-1", map[string]interface{}{})
	assert.Error(t, err)
}

func TestNewSandboxClient(t *testing.T) {
	c := NewSandboxClient("http://localhost:8098")
	assert.Equal(t, "http://localhost:8098", c.baseURL)
	assert.NotNil(t, c.httpClient)
}
