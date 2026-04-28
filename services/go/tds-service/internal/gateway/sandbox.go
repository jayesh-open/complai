package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SandboxClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewSandboxClient(baseURL string) *SandboxClient {
	return &SandboxClient{baseURL: baseURL, httpClient: &http.Client{}}
}

type PANVerifyResult struct {
	PAN      string `json:"pan"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Category string `json:"category"`
}

func (c *SandboxClient) VerifyPAN(ctx context.Context, tenantID, pan, name string) (*PANVerifyResult, error) {
	body, _ := json.Marshal(map[string]string{"pan": pan, "name": name})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/gateway/tds/pan/verify", bytes.NewReader(body))
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
		return nil, fmt.Errorf("PAN verify returned %d", resp.StatusCode)
	}
	var envelope struct {
		Data PANVerifyResult `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

type ChallanResult struct {
	ChallanNumber string  `json:"challan_number"`
	BSRCode       string  `json:"bsr_code"`
	DepositDate   string  `json:"deposit_date"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
}

func (c *SandboxClient) GenerateChallan(ctx context.Context, tenantID string, payload map[string]interface{}) (*ChallanResult, error) {
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/gateway/tds/challan/generate", bytes.NewReader(body))
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
		return nil, fmt.Errorf("challan generate returned %d", resp.StatusCode)
	}
	var envelope struct {
		Data ChallanResult `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}
