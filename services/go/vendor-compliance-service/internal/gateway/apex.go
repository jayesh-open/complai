package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/vendor-compliance-service/internal/domain"
	"github.com/complai/complai/services/go/vendor-compliance-service/internal/scorer"
)

type ApexClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewApexClient(baseURL string) *ApexClient {
	return &ApexClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type successWrapper struct {
	Data json.RawMessage `json:"data"`
}

type gatewayWrapper struct {
	Data json.RawMessage `json:"data"`
}

type vendorsPayload struct {
	Vendors []apexVendor `json:"vendors"`
}

type invoicesPayload struct {
	Invoices []apexInvoice `json:"invoices"`
}

type apexVendor struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	LegalName          string `json:"legal_name"`
	TradeName          string `json:"trade_name"`
	PAN                string `json:"pan"`
	GSTIN              string `json:"gstin"`
	TAN                string `json:"tan"`
	State              string `json:"state"`
	StateCode          string `json:"state_code"`
	Category           string `json:"category"`
	RegistrationStatus string `json:"registration_status"`
	MSMERegistered     bool   `json:"msme_registered"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	Address            string `json:"address"`
}

type apexInvoice struct {
	VendorID        string `json:"vendor_id"`
	IRNGenerated    bool   `json:"irn_generated"`
	GSTFilingStatus string `json:"gst_filing_status"`
	MismatchStatus  string `json:"mismatch_status"`
	PaymentStatus   string `json:"payment_status"`
	PaymentDate     string `json:"payment_date"`
	DueDate         string `json:"due_date"`
}

type fetchVendorsReq struct {
	TenantID  string `json:"tenant_id"`
	RequestID string `json:"request_id"`
}

type fetchInvoicesReq struct {
	TenantID  string `json:"tenant_id"`
	VendorID  string `json:"vendor_id"`
	RequestID string `json:"request_id"`
}

func (c *ApexClient) FetchVendors(ctx context.Context, tenantID uuid.UUID) ([]domain.VendorSnapshot, error) {
	reqBody := fetchVendorsReq{
		TenantID:  tenantID.String(),
		RequestID: uuid.New().String(),
	}
	body, _ := json.Marshal(reqBody)

	url := fmt.Sprintf("%s/v1/gateway/apex/vendors", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch vendors: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch vendors: unexpected status %d", resp.StatusCode)
	}

	// Wire format: {"data": {"data": {"vendors": [...]}, "meta": {...}}}
	var outer successWrapper
	if err := json.NewDecoder(resp.Body).Decode(&outer); err != nil {
		return nil, fmt.Errorf("decode outer wrapper: %w", err)
	}
	var gw gatewayWrapper
	if err := json.Unmarshal(outer.Data, &gw); err != nil {
		return nil, fmt.Errorf("decode gateway wrapper: %w", err)
	}
	var payload vendorsPayload
	if err := json.Unmarshal(gw.Data, &payload); err != nil {
		return nil, fmt.Errorf("decode vendors payload: %w", err)
	}

	now := time.Now().UTC()
	snapshots := make([]domain.VendorSnapshot, len(payload.Vendors))
	for i, v := range payload.Vendors {
		snapshots[i] = domain.VendorSnapshot{
			ID:                 uuid.New(),
			TenantID:           tenantID,
			VendorID:           v.ID,
			Name:               v.Name,
			LegalName:          v.LegalName,
			TradeName:          v.TradeName,
			PAN:                v.PAN,
			GSTIN:              v.GSTIN,
			TAN:                v.TAN,
			State:              v.State,
			StateCode:          v.StateCode,
			Category:           v.Category,
			RegistrationStatus: v.RegistrationStatus,
			MSMERegistered:     v.MSMERegistered,
			Email:              v.Email,
			Phone:              v.Phone,
			Address:            v.Address,
			SyncedAt:           now,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
	}

	return snapshots, nil
}

func (c *ApexClient) FetchAPInvoices(ctx context.Context, tenantID uuid.UUID, vendorID string) ([]scorer.APInvoice, error) {
	reqBody := fetchInvoicesReq{
		TenantID:  tenantID.String(),
		VendorID:  vendorID,
		RequestID: uuid.New().String(),
	}
	body, _ := json.Marshal(reqBody)

	url := fmt.Sprintf("%s/v1/gateway/apex/ap-invoices", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch invoices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch invoices: unexpected status %d", resp.StatusCode)
	}

	// Wire format: {"data": {"data": {"invoices": [...]}, "meta": {...}}}
	var outer successWrapper
	if err := json.NewDecoder(resp.Body).Decode(&outer); err != nil {
		return nil, fmt.Errorf("decode outer wrapper: %w", err)
	}
	var gw gatewayWrapper
	if err := json.Unmarshal(outer.Data, &gw); err != nil {
		return nil, fmt.Errorf("decode gateway wrapper: %w", err)
	}
	var payload invoicesPayload
	if err := json.Unmarshal(gw.Data, &payload); err != nil {
		return nil, fmt.Errorf("decode invoices payload: %w", err)
	}

	invoices := make([]scorer.APInvoice, len(payload.Invoices))
	for i, inv := range payload.Invoices {
		invoices[i] = scorer.APInvoice{
			IRNGenerated:    inv.IRNGenerated,
			GSTFilingStatus: inv.GSTFilingStatus,
			MismatchStatus:  inv.MismatchStatus,
			PaymentStatus:   inv.PaymentStatus,
			PaymentDate:     inv.PaymentDate,
			DueDate:         inv.DueDate,
		}
	}

	return invoices, nil
}
