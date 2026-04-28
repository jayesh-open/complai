package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ApexClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewApexClient(baseURL string) *ApexClient {
	return &ApexClient{baseURL: baseURL, httpClient: &http.Client{}}
}

type ApexVendor struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	PAN      string `json:"pan"`
	Category string `json:"category"`
}

type FetchVendorsResponse struct {
	Data struct {
		Vendors []ApexVendor `json:"vendors"`
		Total   int          `json:"total"`
	} `json:"data"`
}

func (c *ApexClient) FetchVendors(ctx context.Context, tenantID string) ([]ApexVendor, error) {
	payload := map[string]interface{}{
		"tenant_id": tenantID,
		"limit":     200,
		"offset":    0,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/gateway/apex/vendors", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apex vendors returned %d", resp.StatusCode)
	}

	var result FetchVendorsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Data.Vendors, nil
}
