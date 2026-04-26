package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/complai/complai/services/go/recon-service/internal/domain"
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

// GSTR-2B types

type GSTR2BGetRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	RequestID string `json:"request_id"`
}

type gstr2bInvoice struct {
	SupplierGSTIN string  `json:"supplier_gstin"`
	InvoiceNumber string  `json:"invoice_number"`
	InvoiceDate   string  `json:"invoice_date"`
	TaxableValue  float64 `json:"taxable_value"`
	CGSTAmount    float64 `json:"cgst_amount"`
	SGSTAmount    float64 `json:"sgst_amount"`
	IGSTAmount    float64 `json:"igst_amount"`
	TotalValue    float64 `json:"total_value"`
	PlaceOfSupply string  `json:"place_of_supply"`
	ReverseCharge bool    `json:"reverse_charge"`
	HSN           string  `json:"hsn"`
	IMSAction     string  `json:"ims_action"`
}

type gstr2bGetResponse struct {
	GSTIN      string          `json:"gstin"`
	RetPeriod  string          `json:"ret_period"`
	Invoices   []gstr2bInvoice `json:"invoices"`
	TotalCount int             `json:"total_count"`
	Status     string          `json:"status"`
	RequestID  string          `json:"request_id"`
}

func (c *GSTNClient) FetchGSTR2B(ctx context.Context, tenantID uuid.UUID, gstin, returnPeriod string) ([]domain.GSTR2BEntry, error) {
	reqBody := GSTR2BGetRequest{
		GSTIN:     gstin,
		RetPeriod: returnPeriod,
		RequestID: uuid.New().String(),
	}

	resp, err := c.post(ctx, tenantID, "/v1/gateway/adaequare/gstr2b/get", reqBody)
	if err != nil {
		return nil, err
	}

	var getResp gstr2bGetResponse
	if err := json.Unmarshal(resp, &getResp); err != nil {
		return nil, fmt.Errorf("decode gstr2b response: %w", err)
	}

	entries := make([]domain.GSTR2BEntry, 0, len(getResp.Invoices))
	for _, inv := range getResp.Invoices {
		entries = append(entries, domain.GSTR2BEntry{
			SupplierGSTIN: inv.SupplierGSTIN,
			InvoiceNumber: inv.InvoiceNumber,
			InvoiceDate:   inv.InvoiceDate,
			TaxableValue:  decimal.NewFromFloat(inv.TaxableValue),
			CGSTAmount:    decimal.NewFromFloat(inv.CGSTAmount),
			SGSTAmount:    decimal.NewFromFloat(inv.SGSTAmount),
			IGSTAmount:    decimal.NewFromFloat(inv.IGSTAmount),
			TotalValue:    decimal.NewFromFloat(inv.TotalValue),
			HSN:           inv.HSN,
			PlaceOfSupply: inv.PlaceOfSupply,
			ReverseCharge: inv.ReverseCharge,
			IMSAction:     inv.IMSAction,
		})
	}

	return entries, nil
}

// IMS types

type IMSGetRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	RequestID string `json:"request_id"`
}

type IMSGetResponse struct {
	GSTIN      string       `json:"gstin"`
	RetPeriod  string       `json:"ret_period"`
	Invoices   []IMSInvoice `json:"invoices"`
	TotalCount int          `json:"total_count"`
	Status     string       `json:"status"`
	RequestID  string       `json:"request_id"`
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
	ActionAt      string  `json:"action_at,omitempty"`
	ActionBy      string  `json:"action_by,omitempty"`
}

func (c *GSTNClient) FetchIMSState(ctx context.Context, tenantID uuid.UUID, gstin, returnPeriod string) (*IMSGetResponse, error) {
	reqBody := IMSGetRequest{
		GSTIN:     gstin,
		RetPeriod: returnPeriod,
		RequestID: uuid.New().String(),
	}

	resp, err := c.post(ctx, tenantID, "/v1/gateway/adaequare/ims/get", reqBody)
	if err != nil {
		return nil, err
	}

	var imsResp IMSGetResponse
	if err := json.Unmarshal(resp, &imsResp); err != nil {
		return nil, fmt.Errorf("decode ims response: %w", err)
	}

	return &imsResp, nil
}

// IMS Action types

type IMSActionGatewayRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	InvoiceID string `json:"invoice_id"`
	Action    string `json:"action"`
	Reason    string `json:"reason,omitempty"`
	RequestID string `json:"request_id"`
}

type IMSActionGatewayResponse struct {
	InvoiceID string `json:"invoice_id"`
	Action    string `json:"action"`
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
}

func (c *GSTNClient) SendIMSAction(ctx context.Context, tenantID uuid.UUID, gstin, returnPeriod, invoiceID, action, reason string) (*IMSActionGatewayResponse, error) {
	reqBody := IMSActionGatewayRequest{
		GSTIN:     gstin,
		RetPeriod: returnPeriod,
		InvoiceID: invoiceID,
		Action:    action,
		Reason:    reason,
		RequestID: uuid.New().String(),
	}

	resp, err := c.post(ctx, tenantID, "/v1/gateway/adaequare/ims/action", reqBody)
	if err != nil {
		return nil, err
	}

	var actionResp IMSActionGatewayResponse
	if err := json.Unmarshal(resp, &actionResp); err != nil {
		return nil, fmt.Errorf("decode ims action response: %w", err)
	}

	return &actionResp, nil
}

// IMS Bulk Action types

type IMSBulkActionGatewayRequest struct {
	GSTIN      string   `json:"gstin"`
	RetPeriod  string   `json:"ret_period"`
	InvoiceIDs []string `json:"invoice_ids"`
	Action     string   `json:"action"`
	RequestID  string   `json:"request_id"`
}

type IMSBulkActionGatewayResponse struct {
	OperationID   string `json:"operation_id"`
	TotalInvoices int    `json:"total_invoices"`
	Status        string `json:"status"`
	RequestID     string `json:"request_id"`
}

func (c *GSTNClient) SendIMSBulkAction(ctx context.Context, tenantID uuid.UUID, gstin, returnPeriod string, invoiceIDs []string, action string) (*IMSBulkActionGatewayResponse, error) {
	reqBody := IMSBulkActionGatewayRequest{
		GSTIN:      gstin,
		RetPeriod:  returnPeriod,
		InvoiceIDs: invoiceIDs,
		Action:     action,
		RequestID:  uuid.New().String(),
	}

	resp, err := c.post(ctx, tenantID, "/v1/gateway/adaequare/ims/bulk-action", reqBody)
	if err != nil {
		return nil, err
	}

	var bulkResp IMSBulkActionGatewayResponse
	if err := json.Unmarshal(resp, &bulkResp); err != nil {
		return nil, fmt.Errorf("decode ims bulk action response: %w", err)
	}

	return &bulkResp, nil
}

func (c *GSTNClient) post(ctx context.Context, tenantID uuid.UUID, path string, body interface{}) (json.RawMessage, error) {
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
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&gatewayResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return gatewayResp.Data, nil
}
