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

func TestApexClient_FetchVendors_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/gateway/apex/vendors", r.URL.Path)
		assert.Equal(t, "tenant-1", r.Header.Get("X-Tenant-Id"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FetchVendorsResponse{
			Data: struct {
				Vendors []ApexVendor `json:"vendors"`
				Total   int          `json:"total"`
			}{
				Vendors: []ApexVendor{
					{ID: "v1", TenantID: "tenant-1", Name: "Vendor A", PAN: "ABCCD1234E", Category: "COMPANY"},
					{ID: "v2", TenantID: "tenant-1", Name: "Vendor B", PAN: "XYZPD5678F", Category: "INDIVIDUAL"},
				},
				Total: 2,
			},
		})
	}))
	defer srv.Close()

	c := NewApexClient(srv.URL)
	vendors, err := c.FetchVendors(context.Background(), "tenant-1")
	require.NoError(t, err)
	assert.Len(t, vendors, 2)
	assert.Equal(t, "Vendor A", vendors[0].Name)
	assert.Equal(t, "COMPANY", vendors[0].Category)
}

func TestApexClient_FetchVendors_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	c := NewApexClient(srv.URL)
	_, err := c.FetchVendors(context.Background(), "tenant-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "503")
}

func TestApexClient_FetchVendors_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	c := NewApexClient(srv.URL)
	_, err := c.FetchVendors(context.Background(), "tenant-1")
	assert.Error(t, err)
}

func TestApexClient_FetchVendors_ConnectionError(t *testing.T) {
	c := NewApexClient("http://127.0.0.1:1")
	_, err := c.FetchVendors(context.Background(), "tenant-1")
	assert.Error(t, err)
}

func TestApexClient_FetchVendors_EmptyResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FetchVendorsResponse{})
	}))
	defer srv.Close()

	c := NewApexClient(srv.URL)
	vendors, err := c.FetchVendors(context.Background(), "tenant-1")
	require.NoError(t, err)
	assert.Empty(t, vendors)
}

func TestNewApexClient(t *testing.T) {
	c := NewApexClient("http://localhost:9090")
	assert.Equal(t, "http://localhost:9090", c.baseURL)
	assert.NotNil(t, c.httpClient)
}
