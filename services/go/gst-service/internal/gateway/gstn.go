package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type GSTNClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewGSTNClient(baseURL string) *GSTNClient {
	return &GSTNClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type GSTR1SaveRequest struct {
	GSTIN      string      `json:"gstin"`
	RetPeriod  string      `json:"ret_period"`
	B2B        interface{} `json:"b2b,omitempty"`
	B2CL       interface{} `json:"b2cl,omitempty"`
	B2CS       interface{} `json:"b2cs,omitempty"`
	CDNR       interface{} `json:"cdnr,omitempty"`
	CDNUR      interface{} `json:"cdnur,omitempty"`
	EXP        interface{} `json:"exp,omitempty"`
	NIL        interface{} `json:"nil,omitempty"`
	HSN        interface{} `json:"hsn,omitempty"`
	DOC        interface{} `json:"doc_issue,omitempty"`
}

type GSTR1SubmitRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
}

type GSTR1FileRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	SignType  string `json:"sign_type"`
	OTP       string `json:"otp,omitempty"`
}

type GSTNResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ARN     string `json:"arn,omitempty"`
}

func (c *GSTNClient) SaveGSTR1(ctx context.Context, tenantID uuid.UUID, req GSTR1SaveRequest) (*GSTNResponse, error) {
	return c.post(ctx, tenantID, "/v1/gateway/adaequare/gstr1/save", req)
}

func (c *GSTNClient) SubmitGSTR1(ctx context.Context, tenantID uuid.UUID, req GSTR1SubmitRequest) (*GSTNResponse, error) {
	return c.post(ctx, tenantID, "/v1/gateway/adaequare/gstr1/submit", req)
}

func (c *GSTNClient) FileGSTR1(ctx context.Context, tenantID uuid.UUID, req GSTR1FileRequest) (*GSTNResponse, error) {
	return c.post(ctx, tenantID, "/v1/gateway/adaequare/gstr1/file", req)
}

func (c *GSTNClient) post(ctx context.Context, tenantID uuid.UUID, path string, body interface{}) (*GSTNResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-Idempotency-Key", uuid.New().String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gstn gateway call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gstn gateway returned %d", resp.StatusCode)
	}

	var gatewayResp struct {
		Data GSTNResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&gatewayResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &gatewayResp.Data, nil
}
