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

// GSTR-2B / IMS / GSTR-3B types

type GSTR2BGetRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
}

type GSTR2BInvoice struct {
	SupplierGSTIN string  `json:"supplier_gstin"`
	InvoiceNumber string  `json:"invoice_number"`
	InvoiceDate   string  `json:"invoice_date"`
	TaxableValue  float64 `json:"taxable_value"`
	CGSTAmount    float64 `json:"cgst_amount"`
	SGSTAmount    float64 `json:"sgst_amount"`
	IGSTAmount    float64 `json:"igst_amount"`
	TotalValue    float64 `json:"total_value"`
	HSN           string  `json:"hsn"`
	ReverseCharge bool    `json:"reverse_charge"`
	IMSAction     string  `json:"ims_action"`
}

type GSTR2BGetResponse struct {
	GSTIN       string          `json:"gstin"`
	RetPeriod   string          `json:"ret_period"`
	Invoices    []GSTR2BInvoice `json:"invoices"`
	TotalCount  int             `json:"total_count"`
	GeneratedOn string          `json:"generated_on"`
	Status      string          `json:"status"`
}

type IMSGetRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
}

type IMSInvoice struct {
	InvoiceID     string  `json:"invoice_id"`
	SupplierGSTIN string  `json:"supplier_gstin"`
	InvoiceNumber string  `json:"invoice_number"`
	InvoiceDate   string  `json:"invoice_date"`
	TaxableValue  float64 `json:"taxable_value"`
	TotalValue    float64 `json:"total_value"`
	CGSTAmount    float64 `json:"cgst_amount"`
	SGSTAmount    float64 `json:"sgst_amount"`
	IGSTAmount    float64 `json:"igst_amount"`
	Action        string  `json:"action"`
}

type IMSGetResponse struct {
	GSTIN     string       `json:"gstin"`
	RetPeriod string       `json:"ret_period"`
	Invoices  []IMSInvoice `json:"invoices"`
	Summary   IMSSummary   `json:"summary"`
	Status    string       `json:"status"`
}

type IMSSummary struct {
	Accepted      int     `json:"accepted"`
	Rejected      int     `json:"rejected"`
	Pending       int     `json:"pending"`
	AcceptedValue float64 `json:"accepted_value"`
	RejectedValue float64 `json:"rejected_value"`
	PendingValue  float64 `json:"pending_value"`
}

type IMSActionRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	InvoiceID string `json:"invoice_id"`
	Action    string `json:"action"`
	Reason    string `json:"reason,omitempty"`
}

type IMSActionResponse struct {
	InvoiceID string `json:"invoice_id"`
	Action    string `json:"action"`
	Status    string `json:"status"`
}

type GSTR1SummaryRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
}

type GSTR1SummaryResponse struct {
	GSTIN     string                 `json:"gstin"`
	RetPeriod string                 `json:"ret_period"`
	Summary   map[string]interface{} `json:"summary"`
	Status    string                 `json:"status"`
}

type GSTR3BSaveRequest struct {
	GSTIN     string      `json:"gstin"`
	RetPeriod string      `json:"ret_period"`
	Data      interface{} `json:"data"`
}

type GSTR3BSubmitRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
}

type GSTR3BFileRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	SignType  string `json:"sign_type"`
	OTP       string `json:"otp,omitempty"`
}

func (c *GSTNClient) GetGSTR2B(ctx context.Context, tenantID uuid.UUID, req GSTR2BGetRequest) (*GSTR2BGetResponse, error) {
	var resp GSTR2BGetResponse
	err := c.postDecode(ctx, tenantID, "/v1/gateway/adaequare/gstr2b/get", req, &resp)
	return &resp, err
}

func (c *GSTNClient) GetIMS(ctx context.Context, tenantID uuid.UUID, req IMSGetRequest) (*IMSGetResponse, error) {
	var resp IMSGetResponse
	err := c.postDecode(ctx, tenantID, "/v1/gateway/adaequare/ims/get", req, &resp)
	return &resp, err
}

func (c *GSTNClient) SendIMSAction(ctx context.Context, tenantID uuid.UUID, req IMSActionRequest) (*IMSActionResponse, error) {
	var resp IMSActionResponse
	err := c.postDecode(ctx, tenantID, "/v1/gateway/adaequare/ims/action", req, &resp)
	return &resp, err
}

func (c *GSTNClient) GetGSTR1Summary(ctx context.Context, tenantID uuid.UUID, req GSTR1SummaryRequest) (*GSTR1SummaryResponse, error) {
	var resp GSTR1SummaryResponse
	err := c.postDecode(ctx, tenantID, "/v1/gateway/adaequare/gstr1/summary", req, &resp)
	return &resp, err
}

func (c *GSTNClient) SaveGSTR3B(ctx context.Context, tenantID uuid.UUID, req GSTR3BSaveRequest) (*GSTNResponse, error) {
	return c.post(ctx, tenantID, "/v1/gateway/adaequare/gstr3b/save", req)
}

func (c *GSTNClient) SubmitGSTR3B(ctx context.Context, tenantID uuid.UUID, req GSTR3BSubmitRequest) (*GSTNResponse, error) {
	return c.post(ctx, tenantID, "/v1/gateway/adaequare/gstr3b/submit", req)
}

func (c *GSTNClient) FileGSTR3B(ctx context.Context, tenantID uuid.UUID, req GSTR3BFileRequest) (*GSTNResponse, error) {
	return c.post(ctx, tenantID, "/v1/gateway/adaequare/gstr3b/file", req)
}

func (c *GSTNClient) postDecode(ctx context.Context, tenantID uuid.UUID, path string, body interface{}, result interface{}) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-Idempotency-Key", uuid.New().String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("gstn gateway call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gstn gateway returned %d", resp.StatusCode)
	}

	var gatewayResp struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&gatewayResp); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return json.Unmarshal(gatewayResp.Data, result)
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
