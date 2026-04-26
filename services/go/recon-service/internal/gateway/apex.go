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

type ApexClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewApexClient(baseURL string) *ApexClient {
	return &ApexClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type apexAPInvoicesRequest struct {
	GSTIN        string `json:"gstin"`
	ReturnPeriod string `json:"return_period"`
	RequestID    string `json:"request_id"`
}

type apexAPInvoice struct {
	ID            string  `json:"id"`
	VendorGSTIN   string  `json:"vendor_gstin"`
	InvoiceNumber string  `json:"invoice_number"`
	InvoiceDate   string  `json:"invoice_date"`
	TaxableValue  float64 `json:"taxable_value"`
	CGSTAmount    float64 `json:"cgst_amount"`
	SGSTAmount    float64 `json:"sgst_amount"`
	IGSTAmount    float64 `json:"igst_amount"`
	TotalAmount   float64 `json:"total_amount"`
	HSN           string  `json:"hsn,omitempty"`
	PlaceOfSupply string  `json:"place_of_supply"`
	ReverseCharge bool    `json:"reverse_charge"`
}

type apexAPInvoicesResponse struct {
	Invoices []apexAPInvoice `json:"invoices"`
	Total    int             `json:"total"`
}

func (c *ApexClient) FetchAPInvoices(ctx context.Context, tenantID uuid.UUID, gstin, returnPeriod string) ([]domain.PurchaseRegisterEntry, error) {
	reqBody := apexAPInvoicesRequest{
		GSTIN:        gstin,
		ReturnPeriod: returnPeriod,
		RequestID:    uuid.New().String(),
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + "/v1/gateway/apex/ap-invoices"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-Idempotency-Key", uuid.New().String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("apex gateway call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apex gateway returned %d", resp.StatusCode)
	}

	var wrapper struct {
		Data struct {
			Data apexAPInvoicesResponse `json:"data"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	entries := make([]domain.PurchaseRegisterEntry, 0, len(wrapper.Data.Data.Invoices))
	for _, inv := range wrapper.Data.Data.Invoices {
		entries = append(entries, domain.PurchaseRegisterEntry{
			InvoiceNumber: inv.InvoiceNumber,
			InvoiceDate:   inv.InvoiceDate,
			VendorGSTIN:   inv.VendorGSTIN,
			TaxableValue:  decimal.NewFromFloat(inv.TaxableValue),
			CGSTAmount:    decimal.NewFromFloat(inv.CGSTAmount),
			SGSTAmount:    decimal.NewFromFloat(inv.SGSTAmount),
			IGSTAmount:    decimal.NewFromFloat(inv.IGSTAmount),
			TotalValue:    decimal.NewFromFloat(inv.TotalAmount),
			HSN:           inv.HSN,
			PlaceOfSupply: inv.PlaceOfSupply,
			ReverseCharge: inv.ReverseCharge,
			SourceID:      inv.ID,
		})
	}

	return entries, nil
}
